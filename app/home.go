package app

import (
	"net/http"
)

func (ro *Router) handleHome(w http.ResponseWriter, r *http.Request) {
	ro.RenderPage(w, r, "home", "", nil)
}
