package handler

import (
	"errors"
	"html/template"
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/sirupsen/logrus"
)

type helpPage struct {
	file  string
	title string
}

func (h *Handler) helpPages(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	page := ps.ByName("page")
	helpPage := helpPageByURL(page)

	tmpl := template.Must(template.ParseFiles(
		"handler/templates/layout.html",
		"handler/templates/help/layout.html",
		"handler/templates/help/menu.html",
		"handler/templates/help/"+helpPage.file+".html",
	))
	data, err := h.getUserAndState(r, w, helpPage.title+" ⚔️ Millwheat")
	if err != nil {
		data = Game{
			Page: Page{
				Title: helpPage.title + " ⚔️ Millwheat",
			},
		}
	}

	if err := tmpl.Execute(w, data); err != nil {
		logrus.Errorf("failed to execute layout: %v", err)
		error500(w, errors.New("failed to create layout"))
		return
	}
}

func helpPageByURL(path string) helpPage {
	pages := map[string]helpPage{
		"/":          {file: "index", title: "Help pages"},
		"/buildings": {file: "buildings", title: "Buildings ⚔️ Help"},
	}

	p, ok := pages[path]
	if !ok {
		return helpPage{
			file:  "notfound",
			title: "Help page not found",
		}
	}

	return p
}
