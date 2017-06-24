package router

import (
	"net/http"

	"github.com/gorilla/context"
	"github.com/julienschmidt/httprouter"
)

type router struct {
	*httprouter.Router
}

// New returns a router object ready to use.
func New() *router {
	return &router{httprouter.New()}
}

func (r *router) POST(path string, h http.Handler) {
	r.Handle("POST", path, wrapHandler(h))
}

func wrapHandler(h http.Handler) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		context.Set(r, "params", ps)
		h.ServeHTTP(w, r)
	}
}
