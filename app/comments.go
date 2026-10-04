package app

import (
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/nxrmqlly/district/app/httpx"
	"github.com/nxrmqlly/district/store"
)

func (ro *Router) handleCommentCreate(w http.ResponseWriter, r *http.Request) {
	postID, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err != nil {
		log.Println(err)
		http.Error(w, "error parsing post id", http.StatusBadRequest)
		return
	}

	sess, _ := SessionFromContext(r.Context())
	if err := r.ParseForm(); err != nil {
		httpx.ErrorJSON(w, http.StatusBadRequest, "invalid form")
		return
	}

	var parentID *int64

	rawPID := strings.TrimSpace(r.Form.Get("parent_id"))
	
	// parent id is nil = top lvl comment
	if rawPID != "" {
		pid, err := strconv.ParseInt(rawPID, 10, 64)
		if err != nil {
			http.Error(w, "error parsing parent_id", http.StatusBadRequest)
			return
		}
		parentID = &pid
	} else {
		parentID = nil
	}

	body := strings.TrimSpace(r.FormValue("body"))
	if body == "" {
		http.Error(w, "body must not be empty", http.StatusBadRequest)
		return
	}

	comment, err := ro.queries.CreateComment(r.Context(), store.CreateCommentParams{
		PostID:   postID,
		ParentID: parentID,
		AuthorID: sess.UserID,
		Body:     body,
	})

	httpx.WriteJSON(w, 200, comment)
}
