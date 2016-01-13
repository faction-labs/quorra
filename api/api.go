package api

import (
	"net/http"
	"strings"

	log "github.com/Sirupsen/logrus"
	r "github.com/dancannon/gorethink"
	"github.com/factionlabs/quorra"
	"github.com/factionlabs/quorra/accounts"
	mAuth "github.com/factionlabs/quorra/middleware/auth"
	"github.com/gorilla/context"
	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
)

type APIConfig struct {
	ListenAddr string
	PublicDir  string
	DBAddr     string
	DBName     string
	StoreKey   string
}

type API struct {
	listenAddr string
	publicDir  string
	dbAddr     string
	dbName     string
	session    *r.Session
	store      *sessions.CookieStore
	storeKey   string
}

func NewAPI(c *APIConfig) (*API, error) {
	session, err := meld.DB(c.DBAddr, c.DBName)
	if err != nil {
		return nil, err
	}

	a := &API{
		listenAddr: c.ListenAddr,
		publicDir:  c.PublicDir,
		dbAddr:     c.DBAddr,
		dbName:     c.DBName,
		session:    session,
		store:      sessions.NewCookieStore([]byte(c.StoreKey)),
		storeKey:   c.StoreKey,
	}

	// initdb
	if err := a.initDB(); err != nil {
		return nil, err
	}

	return a, nil
}

func (a *API) initDB() error {
	_, err := r.DB(a.dbName).Run(a.session)

	// create db
	if strings.Contains(err.Error(), "does not exist") {
		if _, err := r.DBCreate(a.dbName).RunWrite(a.session); err != nil {
			return err
		}

		log.Debugf("database created: name=%s", a.dbName)
	}

	return nil
}

func (a *API) db() r.Term {
	return r.DB(a.dbName)
}

func (a *API) Run() error {
	globalMux := http.NewServeMux()

	memberManager, err := accounts.New(a.dbName, a.session)
	if err != nil {
		return err
	}

	// api router; public read; auth-required write
	authRequired := mAuth.NewAuthRequired(memberManager, a.store, a.storeKey)

	apiRouter := mux.NewRouter()
	apiRouter.HandleFunc("/api/accounts", a.getAccounts).Methods("GET")
	globalMux.Handle("/api/", apiRouter)

	// login router; public
	loginRouter := mux.NewRouter()
	loginRouter.HandleFunc("/auth/login", a.login).Methods("POST")
	globalMux.Handle("/auth/", loginRouter)

	// member router; auth required
	accountRouter := mux.NewRouter()
	accountRouter.Handle("/accounts/changepassword", authRequired.Handler(http.HandlerFunc(a.changePassword))).Methods("POST")
	accountRouter.Handle("/accounts/logout", authRequired.Handler(http.HandlerFunc(a.logout))).Methods("GET")
	globalMux.Handle("/accounts/", accountRouter)

	// static handler
	globalMux.Handle("/", http.FileServer(http.Dir(a.publicDir)))

	s := &http.Server{
		Addr:    a.listenAddr,
		Handler: context.ClearHandler(globalMux),
	}

	log.Infof("api serving: addr=%s", a.listenAddr)
	return s.ListenAndServe()
}
