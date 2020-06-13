package handler

import (
	"errors"
	"html/template"
	"math/rand"
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/sirupsen/logrus"

	"github.com/gerbenjacobs/millwheat/game"
	gamedata "github.com/gerbenjacobs/millwheat/game/data"
)

type GameData struct {
	PageUser
	Town          *game.Town
	Buildings     game.Buildings
	Items         game.Items
	Warehouse     map[game.ItemID]game.WarehouseItem
	WarehouseList []game.ItemID
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
		PageUser:      data,
		Town:          currentTown,
		Buildings:     h.Buildings,
		Items:         h.Items,
		Warehouse:     warehouse,
		WarehouseList: gamedata.WarehouseOrder,
	}); err != nil {
		logrus.Errorf("failed to execute layout: %v", err)
		error500(w, errors.New("failed to create layout"))
		return
	}
}
