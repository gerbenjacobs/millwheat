package handler

import (
	"context"
	"io"
	"net/http"
	"os"

	"github.com/google/uuid"
	"github.com/gorilla/handlers"
	"github.com/julienschmidt/httprouter"
	"github.com/sirupsen/logrus"

	"github.com/gerbenjacobs/millwheat/services"
)

func (h *Handler) AuthMiddleware(f httprouter.Handle) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		u, err := h.Auth.ReadFromRequest(r)
		if err != nil || u == nil || u.Valid() != nil {
			// Not logged in
			_ = storeAndSaveFlash(r, w, "error|Please log in")
			http.Redirect(w, r, "/login", http.StatusFound)
			return
		}

		// check user
		data, err := h.UserSvc.User(r.Context(), uuid.MustParse(u.UserID))
		if err != nil {
			// most likely old cookie
			http.Redirect(w, r, "/logout", http.StatusFound)
			return
		}

		// Add identifiers to context
		r = r.WithContext(context.WithValue(r.Context(), services.CtxKeyUserID, u.UserID))
		r = r.WithContext(context.WithValue(r.Context(), services.CtxKeyTownID, data.CurrentTown))
		f(w, r, p)
	}
}

func customLoggingMiddleware(handler http.Handler) http.Handler {
	return handlers.CustomLoggingHandler(os.Stdout, handler, func(_ io.Writer, p handlers.LogFormatterParams) {
		if p.StatusCode < 200 || p.StatusCode > 299 && p.StatusCode != 304 {
			logrus.Debugf("%d %s \"%s %s\" %d \"%s\"", p.StatusCode, p.Request.Proto, p.Request.Method, p.URL.String(), p.Size, p.Request.Header.Get("User-Agent"))
		}
	})
}
