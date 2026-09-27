package app

import (
	"html/template"
	"io/fs"
	"net/http"

	"github.com/nxrmqlly/district/store"
)

// Router handles routes and templates. It implements http.Handler
type Router struct {
	mux       *http.ServeMux
	queries   *store.Queries
	templates *template.Template

	GlobalMws []Middleware
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

	err = ro.routes()
	if err != nil {
		return nil, err
	}

	return &ro, nil
}

func (ro *Router) routes() error {
	ro.GlobalMws = append(ro.GlobalMws, ro.Logging)

	staticFS, err := fs.Sub(staticFS, "static")
	if err != nil {
		return err
	}
	ro.mux.Handle("GET /static/", http.StripPrefix("/static/", http.FileServerFS(staticFS)))

	ro.mux.HandleFunc("GET /hello", func(w http.ResponseWriter, r *http.Request) {
		// ro.RenderPage(w, r, "home", nil)
		w.Write([]byte("hey"))
	})

	ro.mux.HandleFunc("GET /", func(w http.ResponseWriter, r *http.Request) {
		ro.RenderPage(w, r, "home", nil)
	})

	return nil
}

func (ro *Router) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ro.mux.ServeHTTP(w, r)
}

func (ro *Router) Handle(pattern string, f http.HandlerFunc, mws ...Middleware) {
	ro.mux.Handle(pattern, Chain(f, append(ro.GlobalMws, mws...)...))
}
