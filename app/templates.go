package app

import (
	"bytes"
	"embed"
	"fmt"
	"html/template"
	"io/fs"
	"net/http"
	"strings"

	"github.com/nxrmqlly/district/auth"
)

//go:embed templates
var templatesFS embed.FS

type PageContext struct {
	Page    string
	Title   string
	Session *auth.AuthSession
	Site    any // TODO: Config driven site Data
	Data    any
}

// for the custom {{ render .X .X }} directive
func renderDirective(root *template.Template) func(string, any) (template.HTML, error) {
	return func(name string, data any) (template.HTML, error) {
		var buf bytes.Buffer
		tmpl := root.Lookup(name)

		if tmpl == nil {
			return "", fmt.Errorf("template %q not found", name)
		}
		if err := tmpl.Execute(&buf, data); err != nil {
			return "", err
		}
		return template.HTML(buf.String()), nil
	}
}

func parseTemplates() (*template.Template, error) {
	root := template.New("")

	root.Funcs(template.FuncMap{
		"render": renderDirective(root),
	})

	if err := fs.WalkDir(templatesFS, ".", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if d.IsDir() || !strings.HasSuffix(path, ".html") {
			return nil
		}

		_, err = root.ParseFS(templatesFS, path)
		return err

	}); err != nil {
		return nil, err
	}

	return root, nil
}

func (ro *Router) RenderPage(w http.ResponseWriter, r *http.Request, page, title string, data any) {
	sess, _ := SessionFromContext(r.Context())
	if err := ro.templates.ExecuteTemplate(w, "layout", PageContext{
		Page:    page,
		Title:   title,
		Site:    nil,
		Session: sess,
		Data:    data,
	}); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
