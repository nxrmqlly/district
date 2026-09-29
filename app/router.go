package app

import (
	"embed"
	"html/template"
	"io/fs"
	"net/http"

	"github.com/nxrmqlly/district/auth"
	"github.com/nxrmqlly/district/store"
)

// Router handles routes and templates. It implements http.Handler
type Router struct {
	mux       *http.ServeMux
	queries   *store.Queries
	templates *template.Template
	auth      *auth.Service

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
		auth:      auth.New(queries),
	}

	err = ro.routes()
	if err != nil {
		return nil, err
	}

	return &ro, nil
}

//go:embed static
var staticFS embed.FS

func (ro *Router) routes() error {
	csrf := http.NewCrossOriginProtection()
	ro.GlobalMws = append(ro.GlobalMws,
		ro.Logging,        // the csrf handler should ideally go after logging, so caught 403s
		csrf.Handler,      // are logged.
		ro.Authentication, //
	)

	staticFS, err := fs.Sub(staticFS, "static")
	if err != nil {
		return err
	}

	ro.mux.Handle("GET /static/", http.StripPrefix("/static/", http.FileServerFS(staticFS)))
	ro.Handle("GET  /{$}", ro.handleHome)
	ro.Handle("GET  /submit", ro.handleSubmitView)
	ro.Handle("POST /submit", ro.handleSubmitCreate)
	ro.Handle("GET  /p/{id}", ro.handlePostGet)
	ro.Handle("GET  /login", ro.handleLoginView)
	ro.Handle("POST /login", ro.handleLogin)
	ro.Handle("GET  /register", ro.handleRegisterView)
	ro.Handle("POST /register", ro.handleRegister)

	return nil
}

func (ro *Router) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ro.mux.ServeHTTP(w, r)
}

func (ro *Router) Handle(pattern string, f http.HandlerFunc, mws ...Middleware) {
	ro.mux.Handle(pattern, Chain(f, append(ro.GlobalMws, mws...)...))
}
