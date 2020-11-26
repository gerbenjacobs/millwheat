package handler

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/justinas/alice"

	app "github.com/gerbenjacobs/millwheat"
	"github.com/gerbenjacobs/millwheat/game"
	"github.com/gerbenjacobs/millwheat/services"
)

// Handler is your dependency container
type Handler struct {
	mux http.Handler
	Dependencies
}

// Dependencies contains all the dependencies your application and its services require
type Dependencies struct {
	// web services
	Auth    *services.Auth
	UserSvc services.UserService

	// game services
	GameSvc       services.GameService
	TownSvc       services.TownService
	ProductionSvc services.ProductionService
	BattleSvc     services.BattleService

	// game data
	Items     game.Items
	Buildings game.Buildings
}

// New creates a new handler given a set of dependencies
func New(dependencies Dependencies) *Handler {
	h := &Handler{
		Dependencies: dependencies,
	}

	r := httprouter.New()
	r.ServeFiles("/css/*filepath", http.Dir("resources/css"))
	r.ServeFiles("/js/*filepath", http.Dir("resources/js"))
	r.ServeFiles("/images/*filepath", http.Dir("resources/images"))
	r.GET("/", h.index)

	r.GET("/join", h.join)
	r.POST("/join-now", h.joinNow)
	r.GET("/login", h.login)
	r.POST("/login-now", h.loginNow)
	r.GET("/logout", h.logout)
	r.GET("/lore", h.lore)

	r.GET("/game", h.AuthMiddleware(h.game))
	r.POST("/game/produce", h.AuthMiddleware(h.produce))
	r.POST("/game/queue", h.AuthMiddleware(h.queue))
	r.POST("/game/collect", h.AuthMiddleware(h.collect))
	r.POST("/game/cancel", h.AuthMiddleware(h.cancel))
	r.POST("/game/upgrade", h.AuthMiddleware(h.upgrade))
	r.POST("/game/demolish", h.AuthMiddleware(h.demolish))
	r.POST("/game/warriors", h.AuthMiddleware(h.warriors))
	r.GET("/game/building/:buildingID", h.AuthMiddleware(h.building))

	r.GET("/help/*page", h.helpPages)

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
