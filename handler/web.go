package handler

import (
	"errors"
	"html/template"
	"net/http"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/sessions"
	"github.com/julienschmidt/httprouter"
	"github.com/sirupsen/logrus"

	app "github.com/gerbenjacobs/millwheat"
	"github.com/gerbenjacobs/millwheat/services"
)

const (
	MillWheatCookie = "millwheat-cookie"
)

var store = sessions.NewCookieStore([]byte("RI9bL47khbsUao&L"))

type Page struct {
	Title      string
	Menu       string
	Flashes    map[string]string
	Attributes map[string]interface{}
}

type Game struct {
	Page
	*app.User
}

func (h *Handler) index(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	data, err := h.getUserAndState(r, w, "Millwheat")
	if err != nil {
		_ = storeAndSaveFlash(r, w, "error|Failed to load your information")
		http.Redirect(w, r, "/", http.StatusFound)
		return
	}

	tmpl := template.Must(template.ParseFiles(
		"handler/templates/layout.html",
		"handler/templates/index.html",
	))
	if err := tmpl.Execute(w, data); err != nil {
		logrus.Errorf("failed to execute layout: %v", err)
		error500(w, errors.New("failed to create layout"))
		return
	}
}

func (h *Handler) joinNow(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	if err := r.ParseForm(); err != nil {
		_ = storeAndSaveFlash(r, w, "error|Invalid form")
		http.Redirect(w, r, "/join", http.StatusFound)
		return
	}
	e := r.Form.Get("email")
	p := r.Form.Get("password")

	if e == "" || p == "" {
		_ = storeAndSaveFlash(r, w, "error|Email and/or password was empty")
		http.Redirect(w, r, "/join", http.StatusFound)
		return
	}

	user := &app.User{
		Email:     e,
		Password:  p,
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
	}
	err := h.UserSvc.Add(r.Context(), user)
	switch {
	case err == app.ErrUserEmailUniqueness:
		logrus.Info("duplicate email during sign up")
		_ = storeAndSaveFlash(r, w, "info|This e-mail address is already in use")
		http.Redirect(w, r, "/join", http.StatusFound)
		return
	case err != nil:
		logrus.Errorf("failed to create user: %v", err)
		_ = storeAndSaveFlash(r, w, "error|Failed to create your account, please try again")
		http.Redirect(w, r, "/join", http.StatusFound)
		return
	}

	// log the user in
	http.SetCookie(w, &http.Cookie{
		Name:     services.CookieName,
		Value:    user.Token,
		HttpOnly: true,
	})

	http.Redirect(w, r, "/", http.StatusFound) // TODO: redirect to game
}

func (h *Handler) join(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	flashes, _ := getFlashes(r, w)
	tmpl := template.Must(template.ParseFiles(
		"handler/templates/layout.html",
		"handler/templates/join.html",
	))
	err := tmpl.Execute(w, Game{
		Page: Page{
			Title:   "Register for free -- Millwheat",
			Flashes: flashes,
		},
	})
	if err != nil {
		logrus.Errorf("failed to execute layout: %v", err)
		error500(w, errors.New("failed to create layout"))
		return
	}
}

func (h *Handler) login(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	flashes, _ := getFlashes(r, w)
	tmpl := template.Must(template.ParseFiles(
		"handler/templates/layout.html",
		"handler/templates/login.html",
	))
	err := tmpl.Execute(w, Game{
		Page: Page{
			Title:   "Log in -- Millwheat",
			Flashes: flashes,
		},
	})
	if err != nil {
		logrus.Errorf("failed to execute layout: %v", err)
		error500(w, errors.New("failed to create layout"))
		return
	}
}

func (h *Handler) loginNow(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	if err := r.ParseForm(); err != nil {
		_ = storeAndSaveFlash(r, w, "error|Invalid form")
		http.Redirect(w, r, "/login", http.StatusFound)
		return
	}
	e := r.Form.Get("email")
	p := r.Form.Get("password")

	if e == "" || p == "" {
		_ = storeAndSaveFlash(r, w, "error|Email and/or password was empty")
		http.Redirect(w, r, "/login", http.StatusFound)
		return
	}

	user, err := h.UserSvc.Login(r.Context(), e, p)
	if err != nil {
		logrus.Debugf("failed to login: %v", err)
		_ = storeAndSaveFlash(r, w, "error|Failed to log in, please try again")
		http.Redirect(w, r, "/login", http.StatusFound)
		return
	}

	// log the user in
	http.SetCookie(w, &http.Cookie{
		Name:     services.CookieName,
		Value:    user.Token,
		HttpOnly: true,
	})

	// update last login time
	user.UpdatedAt = time.Now().UTC()
	_, _ = h.UserSvc.Update(r.Context(), user)

	http.Redirect(w, r, "/", http.StatusFound) // TODO: redirect to game
}

func (h *Handler) logout(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	http.SetCookie(w, &http.Cookie{
		Name:     services.CookieName,
		Value:    "",
		HttpOnly: true,
		MaxAge:   -1,
	})
	_ = storeAndSaveFlash(r, w, "success|Successfully logged out")
	http.Redirect(w, r, "/", http.StatusFound)
	return
}

func (h *Handler) errorHandler(localErr error) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		data, err := h.getUserAndState(r, w, localErr.Error()+"-- Oops!")
		if err != nil {
			_ = storeAndSaveFlash(r, w, "error|Failed to load your information")
			http.Redirect(w, r, "/", http.StatusFound)
			return
		}

		tmpl := template.Must(template.ParseFiles(
			"handler/templates/layout.html",
			"handler/templates/error.html",
		))
		switch {
		case errors.Is(localErr, app.ErrPageNotFound):
			w.WriteHeader(404)
		}

		if err := tmpl.Execute(w, data); err != nil {
			logrus.Errorf("failed to execute layout: %v", err)
			error500(w, errors.New("failed to create layout"))
			return
		}
	}
}

func storeAndSaveFlash(r *http.Request, w http.ResponseWriter, msg string) error {
	session, _ := store.Get(r, MillWheatCookie)
	session.AddFlash(msg)
	return session.Save(r, w)
}

func getFlashes(r *http.Request, w http.ResponseWriter) (map[string]string, error) {
	session, _ := store.Get(r, MillWheatCookie)
	flashes := session.Flashes()

	m := map[string]string{}
	for f := range flashes {
		fs := strings.SplitN(flashes[f].(string), "|", 2)
		if len(fs) == 2 {
			m[fs[0]] = fs[1]
		}
	}
	return m, session.Save(r, w)
}

func (h *Handler) getUserAndState(r *http.Request, w http.ResponseWriter, title string) (Game, error) {
	u, _ := h.Auth.ReadFromRequest(r)
	loggedIn := false
	if u != nil {
		loggedIn = u.Valid() == nil
	}
	flashes, _ := getFlashes(r, w)

	data := Game{
		Page: Page{
			Title:   title,
			Flashes: flashes,
			Attributes: map[string]interface{}{
				"logged_in": loggedIn,
			},
		},
	}

	if loggedIn {
		user, err := h.UserSvc.User(r.Context(), uuid.MustParse(u.UserID))
		if err == nil {
			data.User = user
		}
	}
	return data, nil
}
