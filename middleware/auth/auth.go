package auth

import (
	"fmt"
	"net/http"
	"strings"

	log "github.com/Sirupsen/logrus"
	"github.com/factionlabs/quorra/accounts"
	"github.com/gorilla/sessions"
)

func defaultDeniedHostHandler(w http.ResponseWriter, r *http.Request) {
	http.Error(w, "unauthorized", http.StatusUnauthorized)
}

type AuthRequired struct {
	deniedHostHandler http.Handler
	manager           *accounts.Manager
	store             *sessions.CookieStore
	storeKey          string
}

func NewAuthRequired(m *accounts.Manager, store *sessions.CookieStore, storeKey string) *AuthRequired {
	return &AuthRequired{
		deniedHostHandler: http.HandlerFunc(defaultDeniedHostHandler),
		store:             store,
		storeKey:          storeKey,
		manager:           m,
	}
}

func (a *AuthRequired) Handler(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		err := a.handleRequest(w, r)
		if err != nil {
			log.Warnf("unauthorized request for %s from %s", r.URL.Path, r.RemoteAddr)
			return
		}
		h.ServeHTTP(w, r)
	})
}

func (a *AuthRequired) handleRequest(w http.ResponseWriter, r *http.Request) error {
	valid := false
	authHeader := r.Header.Get("X-Auth-Token")
	parts := strings.Split(authHeader, ":")
	if len(parts) == 2 {
		// validate
		user := parts[0]
		token := parts[1]
		if err := a.manager.VerifyAuthToken(user, token); err == nil {
			valid = true
			// set current user
			session, _ := a.store.Get(r, a.storeKey)
			session.Values["username"] = user
			session.Save(r, w)
		}
	}

	if !valid {
		a.deniedHostHandler.ServeHTTP(w, r)
		return fmt.Errorf("unauthorized %s", r.RemoteAddr)
	}

	return nil
}

func (a *AuthRequired) HandlerFuncWithNext(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	err := a.handleRequest(w, r)

	if err != nil {
		log.Warnf("unauthorized request for %s from %s", r.URL.Path, r.RemoteAddr)
		return
	}

	if next != nil {
		next(w, r)
	}
}
