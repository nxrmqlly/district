package app

import (
	"database/sql"
	"errors"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
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

	if err != nil {
		log.Println(err)

		if errors.Is(err, pgx.ErrNoRows) || errors.Is(err, sql.ErrNoRows) {
			httpx.ErrorJSON(w, http.StatusBadRequest, "invalid parent comment")
			return
		}

		httpx.ErrorJSON(w, http.StatusInternalServerError, "could not create comment")
		return
	}

	httpx.WriteJSON(w, http.StatusCreated, comment)
}

type CommentTree struct {
	ID             int64
	PostID         int64
	ParentID       *int64 // nil for root
	Body           string
	AuthorUsername string
	AuthorID       uuid.UUID
	CreatedAt      time.Time
	DeletedAt      *time.Time

	Children []*CommentTree
}

func translateCT(c store.GetCommentsByPostRow) *CommentTree {
	return &CommentTree{
		ID:             c.ID,
		PostID:         c.PostID,
		ParentID:       c.ParentID,
		Body:           c.Body,
		AuthorUsername: c.AuthorUsername,
		AuthorID:       c.AuthorID,
		CreatedAt:      c.CreatedAt,
		DeletedAt:      c.DeletedAt,
	}
}

func buildCommentTree(comments []store.GetCommentsByPostRow) []*CommentTree {
	byID := make(map[int64]*CommentTree, len(comments))
	for _, c := range comments {
		byID[c.ID] = translateCT(c)
	}

	var roots []*CommentTree

	for _, c := range comments {
		tree := byID[c.ID]

		// root comment
		if c.ParentID == nil {
			roots = append(roots, tree)
			continue
		}

		parent, ok := byID[*c.ParentID]
		if !ok {
			// parent hard deleted or missing = put into root
			roots = append(roots, tree)
			continue
		}

		parent.Children = append(parent.Children, tree)
	}

	return roots
}
