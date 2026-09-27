package app

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/google/uuid"
	"github.com/nxrmqlly/district/app/httpx"
	"github.com/nxrmqlly/district/store"
)

type SubmitPageData struct {
	Title string
	Body  string
	Error string
}

func (ro *Router) handleSubmitView(w http.ResponseWriter, r *http.Request) {
	ro.RenderPage(w, r, "submit", "new post", SubmitPageData{})
}

func (ro *Router) handleSubmitCreate(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		httpx.ErrorJSON(w, http.StatusBadRequest, "invalid form")
		return
	}

	title := strings.TrimSpace(r.FormValue("title"))
	body := strings.TrimSpace(r.FormValue("body"))

	if title == "" || body == "" {
		ro.RenderPage(w, r, "submit", "new post", SubmitPageData{
			Title: title,
			Body:  body,
			Error: "title and body are required",
		})
		return
	}

	post, err := ro.queries.NewPost(r.Context(), store.NewPostParams{
		AuthorID: uuid.MustParse("01a0e4a6-c902-75f5-b55d-2d20c4eb1468"),
		Title:    title,
		Body:     body,
		EmbedUrl: "",
	})

	if err != nil {
		httpx.ErrorJSON(w, http.StatusInternalServerError, "could not create post")
		return
	}

	http.Redirect(w, r, fmt.Sprintf("/p/%d", post.ID), http.StatusSeeOther)
}

func (ro *Router) handlePostGet(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err != nil {
		log.Println(err)
		http.NotFound(w, r)
		return
	}

	post, err := ro.queries.GetPost(r.Context(), id)
	if err != nil {
		httpx.ErrorJSON(w, http.StatusInternalServerError, "uh oh.")
		return
	}

	httpx.WriteJSON(w, http.StatusOK, post)
}
