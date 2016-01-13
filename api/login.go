package api

import (
	"encoding/json"
	"net/http"

	log "github.com/Sirupsen/logrus"
	"github.com/factionlabs/quorra/accounts"
)

type Credentials struct {
	Username string `json:"username,omitempty"`
	Password string `json:"password,omitempty"`
}

func (a *API) login(w http.ResponseWriter, r *http.Request) {
	var creds *Credentials
	if err := json.NewDecoder(r.Body).Decode(&creds); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	m, err := accounts.New(a.dbName, a.session)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	loginSuccessful, err := m.Authenticate(creds.Username, creds.Password)
	if err != nil {
		log.Errorf("error during login for %s from %s: %s", creds.Username, r.RemoteAddr, err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if !loginSuccessful {
		log.Warnf("invalid login for %s from %s", creds.Username, r.RemoteAddr)
		http.Error(w, "invalid username/password", http.StatusForbidden)
		return
	}

	// return token
	token, err := m.NewAuthToken(creds.Username, r.UserAgent())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if err := json.NewEncoder(w).Encode(token); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (a *API) logout(w http.ResponseWriter, r *http.Request) {
	session, _ := a.store.Get(r, a.storeKey)
	username := session.Values["username"].(string)
	if username == "" {
		http.Error(w, "unauthorized", http.StatusInternalServerError)
		return
	}

	m, err := accounts.New(a.dbName, a.session)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err := m.ClearAuthTokens(username, r.UserAgent()); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (a *API) changePassword(w http.ResponseWriter, r *http.Request) {
	session, _ := a.store.Get(r, a.storeKey)
	var creds *Credentials
	if err := json.NewDecoder(r.Body).Decode(&creds); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	username := session.Values["username"].(string)
	if username == "" {
		http.Error(w, "unauthorized", http.StatusInternalServerError)
		return
	}

	m, err := accounts.New(a.dbName, a.session)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err := m.ChangePassword(username, creds.Password); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
