package app

import (
	"bytes"
	"embed"
	"fmt"
	"html/template"
	"io/fs"
	"net/http"
	"strings"
	"time"

	"github.com/nxrmqlly/district/auth"
	"github.com/nxrmqlly/district/config"
)

//go:embed templates
var templatesFS embed.FS

type PageContext struct {
	Page    string
	Title   string
	Session *auth.AuthSession
	Site    config.Site
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

func TimeAgo(t time.Time) string {
	d := time.Since(t)
	switch {
	case d < time.Minute:
		return "just now"
	case d < time.Hour:
		return fmt.Sprintf("%dm ago", int(d.Minutes()))
	case d < 24*time.Hour:
		return fmt.Sprintf("%dh ago", int(d.Hours()))
	case d < 30*24*time.Hour:
		return fmt.Sprintf("%dd ago", int(d.Hours()/24))
	default:
		return t.Format("2 Jan 2006")
	}
}

func parseTemplates() (*template.Template, error) {
	root := template.New("")

	root.Funcs(template.FuncMap{
		"render":  renderDirective(root),
		"timeago": TimeAgo,
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
		Site:    config.Get().Site,
		Session: sess,
		Data:    data,
	}); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
