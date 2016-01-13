package api

import (
	"encoding/json"
	"net/http"

	"github.com/factionlabs/quorra/accounts"
)

func (a *API) getAccounts(w http.ResponseWriter, r *http.Request) {
	m, err := accounts.New(a.dbName, a.session)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	allAccounts, err := m.All()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err := json.NewEncoder(w).Encode(allAccounts); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
