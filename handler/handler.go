package handler

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/justinas/alice"

	app "github.com/gerbenjacobs/millwheat"
	"github.com/gerbenjacobs/millwheat/services"
)

// Handler is your dependency container
type Handler struct {
	mux http.Handler
	Dependencies
}

// Dependencies contains all the dependencies your application and its services require
type Dependencies struct {
	UserSvc services.UserService
	Auth    *services.Auth
}

// New creates a new handler given a set of dependencies
func New(dependencies Dependencies) *Handler {
	h := &Handler{
		Dependencies: dependencies,
	}

	r := httprouter.New()
	r.ServeFiles("/css/*filepath", http.Dir("resources/css"))
	r.ServeFiles("/images/*filepath", http.Dir("resources/images"))
	r.GET("/", h.index)
	r.GET("/join", h.join)
	r.POST("/join-now", h.joinNow)
	r.GET("/login", h.login)
	r.POST("/login-now", h.loginNow)
	r.GET("/logout", h.logout)
	r.GET("/health", health)

	r.NotFound = http.HandlerFunc(h.errorHandler(app.ErrPageNotFound))
	r.MethodNotAllowed = http.HandlerFunc(h.errorHandler(app.ErrMethodNotAllowed))

	// create chained list of middleware
	// and wrap our router with it
	mw := alice.New(customLoggingMiddleware)
	h.mux = mw.Then(r)
	return h
}

// ServeHTTP makes sure Handler implements the http.Handler interface
// this keeps the underlying mux private
func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h.mux.ServeHTTP(w, r)
}
