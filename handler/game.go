package handler

import (
	"errors"
	"html/template"
	"math/rand"
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/sirupsen/logrus"
)

func (h *Handler) game(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	data, err := h.getUserAndState(r, w, "Game &#x2694;&#xfe0f; Millwheat")
	if err != nil {
		_ = storeAndSaveFlash(r, w, "error|Failed to load your information")
		http.Redirect(w, r, "/", http.StatusFound)
		return
	}

	tmpl, _ := template.New("layout.html").Funcs(template.FuncMap{"rand": rand.Float64}).ParseFiles(
		"handler/templates/layout.html",
		"handler/templates/game.html",
	)

	if err := tmpl.Execute(w, data); err != nil {
		logrus.Errorf("failed to execute layout: %v", err)
		error500(w, errors.New("failed to create layout"))
		return
	}
}
