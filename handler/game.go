package handler

import (
	"errors"
	"fmt"
	"html/template"
	"net/http"
	"strconv"

	"github.com/google/uuid"
	"github.com/julienschmidt/httprouter"
	"github.com/sirupsen/logrus"

	"github.com/gerbenjacobs/millwheat/game"
	gamedata "github.com/gerbenjacobs/millwheat/game/data"
)

type GameData struct {
	PageUser
	Town         *game.Town
	Buildings    game.Buildings
	Items        game.Items
	WarriorTypes []game.WarriorType
	WarriorCosts map[game.WarriorType]game.ItemSetSlice

	Warehouse            map[game.ItemID]game.WarehouseItem
	WarehouseList        []game.ItemID
	WarehouseBreakpoints []game.ItemID

	QueuedJobs      map[uuid.UUID][]*game.Job
	QueuedBuildings []*game.Job

	Season         *game.Season
	LastBattle     *game.Battle
	UpcomingBattle *game.Battle
	MyWarriors     []game.Warrior

	// /game/building/:buildingID
	CurrentBuilding     game.Building
	CurrentTownBuilding game.TownBuilding
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

	// Battle data
	season, err := h.BattleSvc.Season(r.Context())
	if err != nil {
		logrus.Errorf("failed to get season: %v", err)
		error500(w, errors.New("failed to load season"))
		return
	}
	lastBattle, err := h.BattleSvc.LastBattle(r.Context())
	if err != nil {
		logrus.Errorf("failed to get lastBattle: %v", err)
		error500(w, errors.New("failed to load lastBattle"))
		return
	}
	upcomingBattle, err := h.BattleSvc.UpcomingBattle(r.Context())
	if err != nil {
		logrus.Errorf("failed to get upcoming battle: %v", err)
		error500(w, errors.New("failed to load upcoming battle"))
		return
	}
	warriors, err := h.BattleSvc.MyWarriors(r.Context())
	if err != nil {
		logrus.Errorf("failed to get my warriors: %v", err)
		error500(w, errors.New("failed to load warriors"))
		return
	}

	tmpl, _ := template.New("layout.html").Funcs(funcs).ParseFiles(
		"handler/templates/layout.html",
		"handler/templates/game.html",
		"handler/templates/partials/town.html",
		"handler/templates/partials/warehouse.html",
		"handler/templates/partials/buildqueue.html",
		"handler/templates/partials/barracks.html",
	)

	if err := tmpl.Execute(w, GameData{
		PageUser:     data,
		Town:         currentTown,
		Buildings:    h.Buildings,
		Items:        h.Items,
		WarriorTypes: game.WarriorTypes,
		WarriorCosts: game.WarriorCosts,

		Warehouse:            warehouse,
		WarehouseList:        gamedata.WarehouseOrder,
		WarehouseBreakpoints: gamedata.WarehouseOrderBreakpoints,

		QueuedJobs:      h.ProductionSvc.QueuedJobs(r.Context()),
		QueuedBuildings: h.ProductionSvc.QueuedBuildings(r.Context()),

		Season:         season,
		LastBattle:     lastBattle,
		UpcomingBattle: upcomingBattle,
		MyWarriors:     warriors,
	}); err != nil {
		logrus.Errorf("failed to execute layout: %v", err)
		error500(w, errors.New("failed to create layout"))
		return
	}
}

func (h *Handler) produce(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	// handle form data
	err := r.ParseForm()
	if err != nil {
		_ = storeAndSaveFlash(r, w, "error|Failed to submit your data")
		http.Redirect(w, r, redirectPage(r), http.StatusFound)
		return
	}

	buildingID, err := uuid.Parse(r.Form.Get("building"))
	if err != nil {
		_ = storeAndSaveFlash(r, w, "error|Invalid building provided")
		http.Redirect(w, r, redirectPage(r), http.StatusFound)
		return
	}
	qty, err := strconv.Atoi(r.Form.Get("quantity"))
	if err != nil || qty <= 0 {
		_ = storeAndSaveFlash(r, w, "info|You have supplied an invalid number")
		http.Redirect(w, r, redirectPage(r), http.StatusFound)
		return
	}
	itemSet := game.ItemSet{
		ItemID:   game.ItemID(r.Form.Get("product")),
		Quantity: qty,
	}

	// actually produce the items
	if err := h.GameSvc.Produce(r.Context(), buildingID, itemSet); err != nil {
		logrus.Errorf("failed to produce: %s", err)
		_ = storeAndSaveFlash(r, w, "error|Failed to produce your item: "+err.Error())
		http.Redirect(w, r, redirectPage(r), http.StatusFound)
		return
	}

	_ = storeAndSaveFlash(r, w, "success|Item has been queued")
	http.Redirect(w, r, redirectPage(r), http.StatusFound)
}

func (h *Handler) collect(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	// handle form data
	err := r.ParseForm()
	if err != nil {
		_ = storeAndSaveFlash(r, w, "error|Failed to submit your data")
		http.Redirect(w, r, redirectPage(r), http.StatusFound)
		return
	}

	buildingID, err := uuid.Parse(r.Form.Get("building"))
	if err != nil {
		_ = storeAndSaveFlash(r, w, "error|Invalid building provided")
		http.Redirect(w, r, redirectPage(r), http.StatusFound)
		return
	}

	// actually collect the produce
	if err := h.GameSvc.Collect(r.Context(), buildingID); err != nil {
		_ = storeAndSaveFlash(r, w, "error|Failed to collect: "+err.Error())
		http.Redirect(w, r, redirectPage(r), http.StatusFound)
		return
	}

	_ = storeAndSaveFlash(r, w, "success|Items have been stored in your warehouse")
	http.Redirect(w, r, redirectPage(r), http.StatusFound)
}

func (h *Handler) queue(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	// handle the form data
	err := r.ParseForm()
	if err != nil {
		_ = storeAndSaveFlash(r, w, "error|Failed to submit your data")
		http.Redirect(w, r, redirectPage(r), http.StatusFound)
		return
	}

	formBuilding, err := strconv.Atoi(r.Form.Get("building"))
	if err != nil {
		_ = storeAndSaveFlash(r, w, "error|Failed to submit your data")
		http.Redirect(w, r, redirectPage(r), http.StatusFound)
		return
	}
	buildingType := game.BuildingType(formBuilding)

	// actually queue the building
	if err := h.GameSvc.AddBuilding(r.Context(), buildingType); err != nil {
		logrus.Errorf("failed to queue building: %s", err)
		_ = storeAndSaveFlash(r, w, "error|Failed to queue building: "+err.Error())
		http.Redirect(w, r, redirectPage(r), http.StatusFound)
		return
	}

	_ = storeAndSaveFlash(r, w, "success|Building has been queued")
	http.Redirect(w, r, "/game#buildqueue", http.StatusFound)
}

func (h *Handler) upgrade(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	// handle the form data
	err := r.ParseForm()
	if err != nil {
		_ = storeAndSaveFlash(r, w, "error|Failed to submit your data")
		http.Redirect(w, r, redirectPage(r), http.StatusFound)
		return
	}

	id := r.Form.Get("building")
	if id == "" {
		_ = storeAndSaveFlash(r, w, "error|Failed to submit your data")
		http.Redirect(w, r, redirectPage(r), http.StatusFound)
		return
	}

	buildingID, err := uuid.Parse(id)
	if err != nil {
		_ = storeAndSaveFlash(r, w, "error|Failed to submit your data")
		http.Redirect(w, r, redirectPage(r), http.StatusFound)
		return
	}

	// actually upgrade the building
	if err := h.GameSvc.UpgradeBuilding(r.Context(), buildingID); err != nil {
		logrus.Errorf("failed to queue building: %s", err)
		_ = storeAndSaveFlash(r, w, "error|Failed to queue building: "+err.Error())
		http.Redirect(w, r, redirectPage(r), http.StatusFound)
		return
	}

	_ = storeAndSaveFlash(r, w, "success|Building has been queued for upgrading")
	http.Redirect(w, r, "/game#buildqueue", http.StatusFound)
}

func (h *Handler) demolish(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	// handle form data
	err := r.ParseForm()
	if err != nil {
		_ = storeAndSaveFlash(r, w, "error|Failed to submit your data")
		http.Redirect(w, r, redirectPage(r), http.StatusFound)
		return
	}

	id := r.Form.Get("building")
	if id == "" {
		_ = storeAndSaveFlash(r, w, "error|Failed to submit your data")
		http.Redirect(w, r, redirectPage(r), http.StatusFound)
		return
	}

	buildingID, err := uuid.Parse(id)
	if err != nil {
		_ = storeAndSaveFlash(r, w, "error|Failed to submit your data")
		http.Redirect(w, r, redirectPage(r), http.StatusFound)
		return
	}

	// actually demolish the building
	if err := h.GameSvc.DemolishBuilding(r.Context(), buildingID); err != nil {
		logrus.Errorf("failed to demolish building: %s", err)
		_ = storeAndSaveFlash(r, w, "error|Failed to demolish building: "+err.Error())
		http.Redirect(w, r, redirectPage(r), http.StatusFound)
		return
	}

	_ = storeAndSaveFlash(r, w, "success|Building has been demolished")
	http.Redirect(w, r, "/game", http.StatusFound)
}

func (h *Handler) cancel(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	// handle form data
	err := r.ParseForm()
	if err != nil {
		_ = storeAndSaveFlash(r, w, "error|Failed to submit your data")
		http.Redirect(w, r, redirectPage(r), http.StatusFound)
		return
	}

	id := r.Form.Get("job")
	if id == "" {
		id = r.Form.Get("building")
	}
	if id == "" {
		_ = storeAndSaveFlash(r, w, "error|Failed to submit your data")
		http.Redirect(w, r, redirectPage(r), http.StatusFound)
		return
	}

	jobID, err := uuid.Parse(id)
	if err != nil {
		_ = storeAndSaveFlash(r, w, "error|Failed to submit your data")
		http.Redirect(w, r, redirectPage(r), http.StatusFound)
		return
	}

	// actually cancel the job
	if err := h.GameSvc.CancelJob(r.Context(), jobID); err != nil {
		logrus.Errorf("failed to cancel job: %s", err)
		_ = storeAndSaveFlash(r, w, "error|Failed to cancel job: "+err.Error())
		http.Redirect(w, r, redirectPage(r), http.StatusFound)
		return
	}

	_ = storeAndSaveFlash(r, w, "success|Job has been cancelled")
	http.Redirect(w, r, redirectPage(r), http.StatusFound)
}

func (h *Handler) warriors(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	// handle form data
	err := r.ParseForm()
	if err != nil {
		_ = storeAndSaveFlash(r, w, "error|Failed to submit your data")
		http.Redirect(w, r, "/game", http.StatusFound)
		return
	}

	warriorType, err := game.WarriorTypeFromString(r.Form.Get("warriorType"))
	if err != nil {
		_ = storeAndSaveFlash(r, w, "error|Invalid warrior type provided")
		http.Redirect(w, r, "/game", http.StatusFound)
		return
	}
	qty, err := strconv.Atoi(r.Form.Get("quantity"))
	if err != nil || qty <= 0 {
		_ = storeAndSaveFlash(r, w, "info|You have supplied an invalid number")
		http.Redirect(w, r, "/game", http.StatusFound)
		return
	}

	// actually produce the items
	if err := h.GameSvc.CreateWarriors(r.Context(), warriorType, qty); err != nil {
		logrus.Errorf("failed to produce: %s", err)
		_ = storeAndSaveFlash(r, w, "error|Failed to produce your warrior: "+err.Error())
		http.Redirect(w, r, "/game", http.StatusFound)
		return
	}

	_ = storeAndSaveFlash(r, w, "success|Warrior has been trained")
	http.Redirect(w, r, "/game#barracks", http.StatusFound)
}

func (h *Handler) building(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
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

	// create Current(Town)Building
	var currentTownBuilding *game.TownBuilding
	for _, t := range currentTown.Buildings {
		if t.ID.String() == params.ByName("buildingID") {
			currentTownBuilding = &t
			break
		}
	}
	if currentTownBuilding == nil {
		logrus.Errorf("failed to locate building: %s", params.ByName("buildingID"))
		_ = storeAndSaveFlash(r, w, "error|Failed to find building")
		http.Redirect(w, r, "/game", http.StatusFound)
		return
	}
	currentBuilding := h.Buildings[currentTownBuilding.Type]

	tmpl, _ := template.New("layout.html").Funcs(funcs).ParseFiles(
		"handler/templates/layout.html",
		"handler/templates/building.html",
		"handler/templates/partials/warehouse.html",
	)

	// Overwrite title
	data.Title = template.HTML(fmt.Sprintf("%s in %s &#x2694;&#xfe0f; Millwheat",
		currentBuilding.Name, currentTown.Name))
	if err := tmpl.Execute(w, GameData{
		PageUser:     data,
		Town:         currentTown,
		Buildings:    h.Buildings,
		Items:        h.Items,
		WarriorTypes: game.WarriorTypes,
		WarriorCosts: game.WarriorCosts,

		Warehouse:            warehouse,
		WarehouseList:        gamedata.WarehouseOrder,
		WarehouseBreakpoints: gamedata.WarehouseOrderBreakpoints,

		QueuedJobs:      h.ProductionSvc.QueuedJobs(r.Context()),
		QueuedBuildings: h.ProductionSvc.QueuedBuildings(r.Context()),

		CurrentBuilding:     currentBuilding,
		CurrentTownBuilding: *currentTownBuilding,
	}); err != nil {
		logrus.Errorf("failed to execute layout: %v", err)
		error500(w, errors.New("failed to create layout"))
		return
	}
}
