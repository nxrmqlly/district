package app

import (
	"net/http"
	"html/template"

	"github.com/nxrmqlly/district/store"
)

// Router is a http.Handler like object that handles routes and templates
type Router struct {
	mux       *http.ServeMux
	queries   *store.Queries
	templates *template.Template
}

func NewRouter(queries *store.Queries) (*Router, error) {
	tmpl, err := parseTemplates()
	if err != nil {
		return nil, err
	}

	ro := Router{
		mux:       http.NewServeMux(),
		queries:   queries,
		templates: tmpl,
	}
	ro.routes()

	return &ro, nil
}

func (ro *Router) routes() {
	ro.mux.HandleFunc("GET /", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-Type", "text/plain")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("hello, world"))
	})
}

func (ro *Router) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ro.mux.ServeHTTP(w, r)
}

func (ro *Router) Handle(pattern string, f http.HandlerFunc, mws ...Middleware) {
	ro.mux.Handle(pattern, Chain(f, mws...))
}
