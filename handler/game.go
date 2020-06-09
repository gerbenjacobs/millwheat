package handler

import (
	"errors"
	"html/template"
	"math/rand"
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/sirupsen/logrus"

	"github.com/gerbenjacobs/millwheat/game"
)

type GameData struct {
	PageUser
	*game.Town
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

	tmpl, _ := template.New("layout.html").Funcs(template.FuncMap{"rand": rand.Float64}).ParseFiles(
		"handler/templates/layout.html",
		"handler/templates/game.html",
	)

	if err := tmpl.Execute(w, GameData{
		PageUser: data,
		Town:     currentTown,
	}); err != nil {
		logrus.Errorf("failed to execute layout: %v", err)
		error500(w, errors.New("failed to create layout"))
		return
	}
}
