package handler

import (
	"errors"
	"html/template"
	"math/rand"
	"net/http"
	"strconv"

	"github.com/davecgh/go-spew/spew"
	"github.com/julienschmidt/httprouter"
	"github.com/sirupsen/logrus"

	"github.com/gerbenjacobs/millwheat/game"
	gamedata "github.com/gerbenjacobs/millwheat/game/data"
)

type GameData struct {
	PageUser
	Town                 *game.Town
	Buildings            game.Buildings
	Items                game.Items
	Warehouse            map[game.ItemID]game.WarehouseItem
	WarehouseList        []game.ItemID
	WarehouseBreakpoints []game.ItemID
}

var funcs = template.FuncMap{
	"rand": rand.Float64,
	"each": func(interval, n, max int) bool {
		// prevents each to occur when list is empty
		if n == max-1 {
			return false
		}
		return n%interval == 0
	},
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
	)

	if err := tmpl.Execute(w, GameData{
		PageUser:             data,
		Town:                 currentTown,
		Buildings:            h.Buildings,
		Items:                h.Items,
		Warehouse:            warehouse,
		WarehouseList:        gamedata.WarehouseOrder,
		WarehouseBreakpoints: gamedata.WarehouseOrderBreakpoints,
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

	// get user
	data, err := h.getUserAndState(r, w, "Game &#x2694;&#xfe0f; Millwheat")
	if err != nil {
		_ = storeAndSaveFlash(r, w, "error|Failed to load your information")
		http.Redirect(w, r, "/game", http.StatusFound)
		return
	}

	// get town
	currentTown, err := h.TownSvc.Town(r.Context(), data.CurrentTown)
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
	product := game.ItemID(r.Form.Get("product"))
	if !building.CanDealWith(product) {
		_ = storeAndSaveFlash(r, w, "info|This product can not be made here")
		http.Redirect(w, r, "/game", http.StatusFound)
		return
	}

	// produce item
	qty, err := strconv.Atoi(r.Form.Get("quantity"))
	if err != nil {
		_ = storeAndSaveFlash(r, w, "info|You have supplied an invalid number")
		http.Redirect(w, r, "/game", http.StatusFound)
		return
	}
	productionResult, err := building.CreateProduct(product, qty, townBuilding.CurrentLevel)
	if err != nil {
		_ = storeAndSaveFlash(r, w, "error|Failed to create your product")
		http.Redirect(w, r, "/game", http.StatusFound)
		return
	}

	// queue job for warehouse
	// TODO

	spew.Dump(productionResult)
	http.Redirect(w, r, "/game", http.StatusFound)
}
