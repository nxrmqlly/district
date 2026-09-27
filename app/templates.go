package app

import (
	"embed"
	"net/http"
)

//go:embed templates
var templatesFS embed.FS

func RenderPage(w http.ResponseWriter, r *http.Request, template string)
