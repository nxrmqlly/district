package app

import (
	"embed"
	"html/template"
	"io/fs"
	"net/http"
	"strings"

	"github.com/nxrmqlly/district/app/httpx"
)

//go:embed templates
var templatesFS embed.FS

func parseTemplates() (*template.Template, error) {
	templates := template.New("")

	if err := fs.WalkDir(templatesFS, ".", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if d.IsDir() || !strings.HasSuffix(path, ".html") {
			return nil
		}

		// skip pages so we can parse them at request time
		if strings.HasPrefix(path, "templates/pages/") {
			return nil
		}

		_, err = templates.ParseFS(templatesFS, path)
		return err

	}); err != nil {
		return nil, err
	}

	return templates, nil
}

func (ro *Router) RenderPage(w http.ResponseWriter, r *http.Request, page string, data any) {
	tmpl, err := ro.templates.Clone()
	if err != nil {
		httpx.ErrorJSON(w, http.StatusInternalServerError, "internal server error")
		return
	}

	_, err = tmpl.ParseFS(templatesFS, "templates/pages/"+page+".html")
	if err != nil {
		httpx.ErrorJSON(w, http.StatusInternalServerError, "internal server error")
		return
	}

	if err := tmpl.ExecuteTemplate(w, "layout", data); err != nil {
		httpx.ErrorJSON(w, http.StatusInternalServerError, "internal server error")
		return
	}
}
