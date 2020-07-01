package handler

import (
	"errors"
	"fmt"
	"html/template"
	"net/http"
	"strconv"
	"time"

	"github.com/google/uuid"
	"github.com/julienschmidt/httprouter"
	"github.com/sirupsen/logrus"

	"github.com/gerbenjacobs/millwheat/game"
	gamedata "github.com/gerbenjacobs/millwheat/game/data"
	"github.com/gerbenjacobs/millwheat/services"
)

type GameData struct {
	PageUser
	Town                 *game.Town
	Buildings            game.Buildings
	Items                game.Items
	Warehouse            map[game.ItemID]game.WarehouseItem
	WarehouseList        []game.ItemID
	WarehouseBreakpoints []game.ItemID
	QueuedJobs           map[uuid.UUID][]*game.Job
	QueuedBuildings      []*game.Job
}

var funcs = template.FuncMap{
	"hasItemID": func(haystack []game.ItemID, needle game.ItemID) bool {
		for _, v := range haystack {
			if v == needle {
				return true
			}
		}
		return false
	},
}

func (h *Handler) game(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	// create PageUser
	data, err := h.getUserAndState(r, w, "Game &#x2694;&#xfe0f; Millwheat")
	if err != nil {
		_ = storeAndSaveFlash(r, w, "error|Failed to load your information")
		http.Redirect(w, r, "/", http.StatusFound)
		return
	}

	// create game state
	currentTown, err := h.TownSvc.Town(r.Context(), data.CurrentTown)
	if err != nil {
		logrus.Errorf("failed to get current town: %v", err)
		error500(w, errors.New("failed to load town"))
		return
	}
	warehouse, err := h.TownSvc.Warehouse(r.Context(), currentTown.ID)
	if err != nil {
		logrus.Errorf("failed to get warehouse: %v", err)
		error500(w, errors.New("failed to load warehouse"))
		return
	}

	tmpl, _ := template.New("layout.html").Funcs(funcs).ParseFiles(
		"handler/templates/layout.html",
		"handler/templates/game.html",
		"handler/templates/partials/town.html",
		"handler/templates/partials/warehouse.html",
		"handler/templates/partials/buildqueue.html",
	)

	if err := tmpl.Execute(w, GameData{
		PageUser:             data,
		Town:                 currentTown,
		Buildings:            h.Buildings,
		Items:                h.Items,
		Warehouse:            warehouse,
		WarehouseList:        gamedata.WarehouseOrder,
		WarehouseBreakpoints: gamedata.WarehouseOrderBreakpoints,
		QueuedJobs:           h.ProductionSvc.QueuedJobs(r.Context()),
		QueuedBuildings:      h.ProductionSvc.QueuedBuildings(r.Context()),
	}); err != nil {
		logrus.Errorf("failed to execute layout: %v", err)
		error500(w, errors.New("failed to create layout"))
		return
	}
}

func (h *Handler) produce(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	err := r.ParseForm()
	if err != nil {
		_ = storeAndSaveFlash(r, w, "error|Failed to submit your data")
		http.Redirect(w, r, "/game", http.StatusFound)
		return
	}

	// get town
	currentTown, err := h.TownSvc.Town(r.Context(), services.TownFromContext(r.Context()))
	if err != nil {
		logrus.Errorf("failed to get current town: %v", err)
		error500(w, errors.New("failed to load town"))
		return
	}

	// get building
	var townBuilding *game.TownBuilding
	var building *game.Building
	for _, tb := range currentTown.Buildings {
		if r.Form.Get("building") == tb.ID.String() {
			townB := tb
			b := h.Buildings[tb.Type]

			townBuilding = &townB
			building = &b
		}
	}
	if townBuilding == nil || building == nil {
		_ = storeAndSaveFlash(r, w, "info|This building is not found")
		http.Redirect(w, r, "/game", http.StatusFound)
		return
	}

	// check if product has any effect in this building
	qty, err := strconv.Atoi(r.Form.Get("quantity"))
	if err != nil || qty <= 0 {
		_ = storeAndSaveFlash(r, w, "info|You have supplied an invalid number")
		http.Redirect(w, r, "/game", http.StatusFound)
		return
	}

	product := game.ItemID(r.Form.Get("product"))
	if !building.CanDealWith(product) {
		_ = storeAndSaveFlash(r, w, "info|This product can not be made here")
		http.Redirect(w, r, "/game", http.StatusFound)
		return
	}

	// produce item
	productionResult, err := building.CreateProduct(product, qty, townBuilding.CurrentLevel)
	if err != nil {
		_ = storeAndSaveFlash(r, w, "error|Failed to create your product")
		http.Redirect(w, r, "/game", http.StatusFound)
		return
	}

	// check if items are in warehouse
	if !h.TownSvc.ItemsInWarehouse(r.Context(), productionResult.Consumption) {
		_ = storeAndSaveFlash(r, w, "error|You don't have the required products; "+gamedata.ItemSetSlice(productionResult.Consumption).String())
		http.Redirect(w, r, "/game", http.StatusFound)
		return
	}

	// queue job
	job := &game.InputJob{
		Type: game.JobTypeProduct,
		ProductJob: &game.ProductJob{
			BuildingID:  townBuilding.ID,
			Production:  productionResult.Production,
			Consumption: productionResult.Consumption,
		},
		Duration: time.Duration(productionResult.Hours) * time.Minute,
	}
	if err := h.ProductionSvc.CreateJob(r.Context(), job); err != nil {
		_ = storeAndSaveFlash(r, w, "error|Failed to queue your job")
		http.Redirect(w, r, "/game", http.StatusFound)
		return
	}

	// extract consumption items from warehouse
	if err := h.TownSvc.TakeFromWarehouse(r.Context(), productionResult.Consumption); err != nil {
		// TODO undo job without returning items
		_ = storeAndSaveFlash(r, w, "error|Failed to gather items needed")
		http.Redirect(w, r, "/game", http.StatusFound)
		return
	}

	_ = storeAndSaveFlash(r, w, "success|Item has been queued")
	http.Redirect(w, r, "/game", http.StatusFound)
}

func (h *Handler) upgradeBuilding(w http.ResponseWriter, r *http.Request, buildingID *uuid.UUID, buildingType game.BuildingType, level int) {
	building, ok := gamedata.Buildings[buildingType]
	if !ok {
		_ = storeAndSaveFlash(r, w, "error|Failed to submit your data")
		http.Redirect(w, r, "/game", http.StatusFound)
		return
	}
	// get production requirements for building
	productionResult, err := game.CreateBuilding(building, level)
	if err != nil {
		_ = storeAndSaveFlash(r, w, "error|Failed to create your building")
		http.Redirect(w, r, "/game", http.StatusFound)
		return
	}
	// check if items are in warehouse
	if !h.TownSvc.ItemsInWarehouse(r.Context(), productionResult.Consumption) {
		_ = storeAndSaveFlash(r, w, "error|You don't have the required products; "+gamedata.ItemSetSlice(productionResult.Consumption).String())
		http.Redirect(w, r, "/game", http.StatusFound)
		return
	}

	// extract consumption items from warehouse
	if err := h.TownSvc.TakeFromWarehouse(r.Context(), productionResult.Consumption); err != nil {
		_ = storeAndSaveFlash(r, w, "error|Failed to gather items needed")
		http.Redirect(w, r, "/game", http.StatusFound)
		return
	}

	// set or create building id
	bID := uuid.New()
	if buildingID != nil {
		bID = *buildingID
	}

	// queue building job
	if err := h.ProductionSvc.CreateJob(r.Context(), &game.InputJob{
		Type: game.JobTypeBuilding,
		BuildingJob: &game.BuildingJob{
			ID:    bID,
			Type:  buildingType,
			Level: level,
		},
		Duration: 20 * time.Second, // TODO: fix time
	}); err != nil {
		// TODO rollback warehouse items
		_ = storeAndSaveFlash(r, w, "error|Failed to queue your building")
		http.Redirect(w, r, "/game", http.StatusFound)
		return
	}

	_ = storeAndSaveFlash(r, w, "success|Building has been queued")
	http.Redirect(w, r, "/game", http.StatusFound)
}

func (h *Handler) queue(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	err := r.ParseForm()
	if err != nil {
		_ = storeAndSaveFlash(r, w, "error|Failed to submit your data")
		http.Redirect(w, r, "/game", http.StatusFound)
		return
	}

	formBuilding, err := strconv.Atoi(r.Form.Get("building"))
	buildingType := game.BuildingType(formBuilding)

	// do the building upgrading work
	h.upgradeBuilding(w, r, nil, buildingType, 1)
}

func (h *Handler) collect(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	err := r.ParseForm()
	if err != nil {
		_ = storeAndSaveFlash(r, w, "error|Failed to submit your data")
		http.Redirect(w, r, "/game", http.StatusFound)
		return
	}

	// get town
	currentTown, err := h.TownSvc.Town(r.Context(), services.TownFromContext(r.Context()))
	if err != nil {
		logrus.Errorf("failed to get current town: %v", err)
		_ = storeAndSaveFlash(r, w, "error|Failed to load your town")
		http.Redirect(w, r, "/game", http.StatusFound)
		return
	}

	// get building
	var townBuilding *game.TownBuilding
	var building *game.Building
	for _, tb := range currentTown.Buildings {
		if r.Form.Get("building") == tb.ID.String() {
			townB := tb
			b := h.Buildings[tb.Type]

			townBuilding = &townB
			building = &b
		}
	}
	if townBuilding == nil || building == nil {
		_ = storeAndSaveFlash(r, w, "info|This building is not found")
		http.Redirect(w, r, "/game", http.StatusFound)
		return
	}

	cp, err := townBuilding.GetCurrentProduction(*building)
	if err != nil {
		logrus.Warnf("failed to collect resources for %s: %s", townBuilding.ID, err)
		_ = storeAndSaveFlash(r, w, "error|Can't collect items: "+err.Error())
		http.Redirect(w, r, "/game", http.StatusFound)
		return
	}
	if err = h.TownSvc.GiveToWarehouse(r.Context(), []game.ItemSet{*cp}); err != nil {
		_ = storeAndSaveFlash(r, w, "error|Failed to put items in warehouse")
		http.Redirect(w, r, "/game", http.StatusFound)
	}

	townBuilding.CurrentProduction = 0
	townBuilding.LastCollection = time.Now().UTC()
	currentTown.Buildings[townBuilding.ID] = *townBuilding

	_ = storeAndSaveFlash(r, w, fmt.Sprintf("success|%d %s has been stored in your warehouse", cp.Quantity, cp.ItemID))
	http.Redirect(w, r, "/game", http.StatusFound)
}

func (h *Handler) cancel(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	err := r.ParseForm()
	if err != nil {
		_ = storeAndSaveFlash(r, w, "error|Failed to submit your data")
		http.Redirect(w, r, "/game", http.StatusFound)
		return
	}

	id := r.Form.Get("job")
	if id == "" {
		id = r.Form.Get("building")
	}
	if id == "" {
		_ = storeAndSaveFlash(r, w, "error|Failed to submit your data")
		http.Redirect(w, r, "/game", http.StatusFound)
		return
	}

	jobID, err := uuid.Parse(id)
	if err != nil {
		_ = storeAndSaveFlash(r, w, "error|Failed to submit your data")
		http.Redirect(w, r, "/game", http.StatusFound)
		return
	}

	// Collect returnable resources
	resources, err := h.ProductionSvc.RevertJobResources(r.Context(), jobID)
	if err != nil {
		_ = storeAndSaveFlash(r, w, "error|Failed to get your items back: "+err.Error())
		http.Redirect(w, r, "/game", http.StatusFound)
		return
	}

	// Cancel job && reshuffle
	if err = h.ProductionSvc.CancelJob(r.Context(), jobID); err != nil {
		_ = storeAndSaveFlash(r, w, "error|Failed to cancel your job: "+err.Error())
		http.Redirect(w, r, "/game", http.StatusFound)
		return
	}
	h.ProductionSvc.ReshuffleQueue(r.Context())

	// Apply resources to warehouse
	if err = h.TownSvc.GiveToWarehouse(r.Context(), resources); err != nil {
		_ = storeAndSaveFlash(r, w, "error|Failed to put your items in the warehouse: "+err.Error())
		http.Redirect(w, r, "/game", http.StatusFound)
		return
	}

	_ = storeAndSaveFlash(r, w, "success|Job has been cancelled")
	http.Redirect(w, r, "/game", http.StatusFound)
}

func (h *Handler) upgrade(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	err := r.ParseForm()
	if err != nil {
		_ = storeAndSaveFlash(r, w, "error|Failed to submit your data")
		http.Redirect(w, r, "/game", http.StatusFound)
		return
	}

	id := r.Form.Get("building")
	if id == "" {
		_ = storeAndSaveFlash(r, w, "error|Failed to submit your data")
		http.Redirect(w, r, "/game", http.StatusFound)
		return
	}

	buildingID, err := uuid.Parse(id)
	if err != nil {
		_ = storeAndSaveFlash(r, w, "error|Failed to submit your data")
		http.Redirect(w, r, "/game", http.StatusFound)
		return
	}

	// get town
	currentTown, err := h.TownSvc.Town(r.Context(), services.TownFromContext(r.Context()))
	if err != nil {
		logrus.Errorf("failed to get current town: %v", err)
		_ = storeAndSaveFlash(r, w, "error|Failed to load your town")
		http.Redirect(w, r, "/game", http.StatusFound)
		return
	}

	building, ok := currentTown.Buildings[buildingID]
	if !ok {
		_ = storeAndSaveFlash(r, w, "error|Building not found")
		http.Redirect(w, r, "/game", http.StatusFound)
		return
	}

	// do the building upgrading work
	h.upgradeBuilding(w, r, &buildingID, building.Type, building.CurrentLevel+1)
}

func (h *Handler) demolish(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	err := r.ParseForm()
	if err != nil {
		_ = storeAndSaveFlash(r, w, "error|Failed to submit your data")
		http.Redirect(w, r, "/game", http.StatusFound)
		return
	}

	id := r.Form.Get("building")
	if id == "" {
		_ = storeAndSaveFlash(r, w, "error|Failed to submit your data")
		http.Redirect(w, r, "/game", http.StatusFound)
		return
	}

	buildingID, err := uuid.Parse(id)
	if err != nil {
		_ = storeAndSaveFlash(r, w, "error|Failed to submit your data")
		http.Redirect(w, r, "/game", http.StatusFound)
		return
	}

	// fetch recovered resources
	currentTown, err := h.TownSvc.Town(r.Context(), services.TownFromContext(r.Context()))
	if err != nil {
		logrus.Errorf("failed to get current town: %v", err)
		error500(w, errors.New("failed to load town"))
		return
	}
	tb, ok := currentTown.Buildings[buildingID]
	if !ok {
		_ = storeAndSaveFlash(r, w, "error|Failed to locate building")
		http.Redirect(w, r, "/game", http.StatusFound)
		return
	}

	// demolish building
	if err = h.TownSvc.RemoveBuilding(r.Context(), buildingID); err != nil {
		_ = storeAndSaveFlash(r, w, "error|Failed to demolish your building")
		http.Redirect(w, r, "/game", http.StatusFound)
		return
	}

	// give recovered resources to warehouse
	b, ok := gamedata.Buildings[tb.Type]
	pr, err := game.RecoverBuilding(b, tb.CurrentLevel)
	if !ok || err != nil {
		_ = storeAndSaveFlash(r, w, "error|Failed to recover some of the buildings resources")
		http.Redirect(w, r, "/game", http.StatusFound)
		return
	}
	if err = h.TownSvc.GiveToWarehouse(r.Context(), pr.Consumption); err != nil {
		_ = storeAndSaveFlash(r, w, "error|Failed to put your items in the warehouse: "+err.Error())
		http.Redirect(w, r, "/game", http.StatusFound)
		return
	}

	_ = storeAndSaveFlash(r, w, "success|Building has been demolished")
	http.Redirect(w, r, "/game", http.StatusFound)
}
