package app

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"net/mail"
	"regexp"
	"strings"

	"github.com/google/uuid"
	"github.com/nxrmqlly/district/app/httpx"
	"github.com/nxrmqlly/district/auth"
)

var usernameRE = regexp.MustCompile(`^[A-Za-z0-9_]{3,24}$`)

// SessionFromContext derives the auth session from the request context.
// ok is true for authenticated sessions and false for anonymous sessions
func SessionFromContext(ctx context.Context) (*auth.AuthSession, bool) {
	sess, ok := ctx.Value(sessionKey).(*auth.AuthSession)
	return sess, ok
}

// create a new session and set the session cookie to the writer
func (ro *Router) newSessionCookie(w http.ResponseWriter, r *http.Request, userId uuid.UUID) {
	sessToken, err := ro.auth.CreateSession(r.Context(), userId)
	if err != nil {
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}
	http.SetCookie(w, &http.Cookie{
		Name:     "district_session",
		Value:    sessToken,
		Path:     "/",
		HttpOnly: true,
		Secure:   false, // FIXME: secure cookie should be config driven
		SameSite: http.SameSiteLaxMode,
	})
}

type LoginPageData struct {
	Login string
	Error string
}

func (ro *Router) handleLoginView(w http.ResponseWriter, r *http.Request) {
	// already logged in
	if _, ok := SessionFromContext(r.Context()); ok {
		fmt.Println("/login: user already logged in, redirecting...")
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	ro.RenderPage(w, r, "login", "login", LoginPageData{})
}

func (ro *Router) handleLogin(w http.ResponseWriter, r *http.Request) {
	if _, ok := SessionFromContext(r.Context()); ok {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	if err := r.ParseForm(); err != nil {
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}

	login := strings.TrimSpace(r.FormValue("login"))
	password := r.FormValue("password")

	user, err := ro.auth.Authenticate(r.Context(), login, password)
	if err != nil {
		if errors.Is(err, auth.ErrInvalidCredentials) {
			ro.RenderPage(w, r, "login", "login", LoginPageData{
				Error: "invalid login or password",
				Login: login,
			})
		} else {
			log.Printf("err: /login: %s", err.Error())
			ro.RenderPage(w, r, "login", "login", LoginPageData{
				Error: "internal server error",
				Login: login,
			})
		}
		return
	}

	ro.newSessionCookie(w, r, user.ID)
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

type RegisterPageData struct {
	Email    string
	Username string
	Error    string
}

func (ro *Router) handleRegisterView(w http.ResponseWriter, r *http.Request) {
	// already logged in
	if _, ok := SessionFromContext(r.Context()); ok {
		fmt.Println("/register: user already logged in, redirecting...")
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	ro.RenderPage(w, r, "register", "register", RegisterPageData{})
}

func (ro *Router) handleRegister(w http.ResponseWriter, r *http.Request) {
	if _, ok := SessionFromContext(r.Context()); ok {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	if err := r.ParseForm(); err != nil {
		httpx.ErrorJSON(w, http.StatusBadRequest, "invalid form")
		return
	}

	email := strings.TrimSpace(r.FormValue("email"))
	username := strings.TrimSpace(r.FormValue("username"))
	password := r.FormValue("password")

	verrs := ""

	if !usernameRE.MatchString(username) {
		verrs += "username must contain only characters, numbers or underscore; min 3, max 24 characters\n"
	}

	if len(email) > 254 {
		verrs += "invalid email\n"
	}

	if addr, err := mail.ParseAddress(email); err != nil || addr.Address != email {
		verrs += "invalid email\n"
	}

	// TODO: better passwd security
	if len(password) < 8 || len(password) > 128 {
		verrs += "password must be 8 to 128 characters long\n"
	}

	if verrs != "" {
		ro.RenderPage(w, r, "register", "register", RegisterPageData{
			Username: username,
			Email:    email,
			Error:    verrs,
		})
		return
	}

	// * all validations passed

	user, err := ro.auth.RegisterUser(r.Context(), email, username, password)
	if err != nil {
		switch {
		case errors.Is(err, auth.ErrUsernameTaken):
			ro.RenderPage(w, r, "register", "register", RegisterPageData{
				Username: username,
				Email:    email,
				Error:    "username is already taken",
			})
		case errors.Is(err, auth.ErrEmailInUse):
			ro.RenderPage(w, r, "register", "register", RegisterPageData{
				Username: username,
				Email:    email,
				Error:    "email already in use",
			})
		default:
			http.Error(w, "internal server error", http.StatusInternalServerError)
		}
		return
	}

	ro.newSessionCookie(w, r, user.ID)
	http.Redirect(w, r, "/", http.StatusSeeOther)
}
