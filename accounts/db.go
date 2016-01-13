package accounts

import (
	"crypto/rand"
	"encoding/base64"
	"errors"
	"fmt"
	"strings"
	"time"

	log "github.com/Sirupsen/logrus"
	r "github.com/dancannon/gorethink"
	"github.com/factionlabs/quorra/auth"
	"github.com/factionlabs/quorra/auth/builtin"
)

const (
	TableName = "accounts"
)

var (
	ErrAccountExists    = errors.New("account already exists")
	ErrInvalidAuthToken = errors.New("invalid auth token")
)

type Manager struct {
	dbName        string
	session       *r.Session
	authenticator auth.Authenticator
}

func New(dbName string, session *r.Session) (*Manager, error) {
	// TODO: configurable salt
	authenticator := builtin.NewAuthenticator("meld")

	m := &Manager{
		dbName:        dbName,
		session:       session,
		authenticator: authenticator,
	}

	if err := m.init(); err != nil {
		return nil, err
	}

	return m, nil
}

// init initializes the datastore for the specified table
func (m *Manager) init() error {
	_, err := m.t().Run(m.session)

	// create table
	if err != nil {
		if strings.Contains(err.Error(), "does not exist") {
			_, err := r.DB(m.dbName).TableCreate(TableName).RunWrite(m.session)
			if err != nil {
				return err
			}

			log.Debugf("created table: name=%s", TableName)

			// create initial user
			rb := make([]byte, 12)
			if _, err := rand.Read(rb); err != nil {
				return err
			}

			password := base64.StdEncoding.EncodeToString(rb)

			account := &Account{
				Username:  "admin",
				FirstName: "Lighthouse",
				LastName:  "Admin",
				Password:  password,
				IsAdmin:   true,
			}

			if err := m.Save(account); err != nil {
				return err
			}

			log.Debugf("created admin user: username=admin password=%s", password)
		} else {
			return err
		}
	}

	return nil
}

// t is a func that returns an r.Term for the table.
func (m *Manager) t() r.Term {
	return r.DB(m.dbName).Table(TableName)
}

// All returns all Accounts.
func (m *Manager) All() ([]*Account, error) {
	log.Debug("query: accounts.All")

	res, err := m.t().Run(m.session)
	if err != nil {
		return nil, err
	}
	defer res.Close()

	data := []*Account{}
	if err := res.All(&data); err != nil {
		return nil, err
	}

	return data, nil
}

// Get returns the Account by ID
func (m *Manager) Get(id string) (*Account, error) {
	log.Debug("query: accounts.Get")

	res, err := m.t().Get(id).Run(m.session)
	if err != nil {
		return nil, err
	}
	defer res.Close()

	data := &Account{}
	if err := res.One(&data); err != nil {
		return nil, err
	}

	return data, nil
}

// GetByUsername returns a Account by username.
func (m *Manager) GetByUsername(username string) (*Account, error) {
	log.Debug("query: accounts.GetByUsername")
	res, err := m.t().Filter(map[string]string{"username": username}).Run(m.session)
	if err != nil {
		return nil, err
	}
	defer res.Close()

	data := &Account{}
	if err := res.One(&data); err != nil {
		return nil, err
	}

	return data, nil
}

// Save inserts or updates a Account.
func (m *Manager) Save(account *Account) error {
	log.Debug("query: accounts.Save")
	var (
		hash string
		err  error
	)

	// hash password
	if account.Password != "" {
		h, err := auth.Hash(account.Password)
		if err != nil {
			return err
		}

		hash = h
	}

	if account.Id != "" {
		updates := map[string]interface{}{
			"firstName": account.FirstName,
			"lastName":  account.LastName,
			"email":     account.Email,
			"isAdmin":   account.IsAdmin,
		}

		if account.Password != "" {
			updates["password"] = hash
		}

		_, err = m.t().Get(account.Id).Update(updates).RunWrite(m.session)
		log.Debugf("updated account: id=%s username=%s", account.Id, account.Username)
	} else {
		// set creation date
		account.Created = time.Now().Unix()

		// check for existing account when creating a new account
		existing, _ := m.GetByUsername(account.Username)
		if existing != nil {
			return ErrAccountExists
		}

		account.Password = hash
		_, err = m.t().Insert(account).RunWrite(m.session)
		log.Debugf("saved account: username=%s", account.Username)
	}

	if err != nil {
		return err
	}

	return nil
}

// Delete removes a Account.
func (m *Manager) Delete(id string) error {
	log.Debug("query: accounts.Delete")
	if _, err := m.t().Get(id).Delete().RunWrite(m.session); err != nil {
		return err
	}

	log.Debugf("deleted account: id=%s", id)

	return nil
}

// VerifyAuthToken verifies a user auth token.  This returns an error if invalid.
func (m *Manager) VerifyAuthToken(username, token string) error {
	account, err := m.GetByUsername(username)
	if err != nil {
		return err
	}

	for _, t := range account.Tokens {
		if token == t.Token {
			return nil
		}
	}

	return ErrInvalidAuthToken
}

// Authenticate validates the user and password.
func (m *Manager) Authenticate(username, password string) (bool, error) {
	// only get the account to get the hashed password if using the builtin auth
	passwordHash := ""
	if m.authenticator.Name() == "builtin" {
		account, err := m.GetByUsername(username)
		if err != nil {
			log.Error(err)
			return false, err
		}

		passwordHash = account.Password
	}

	return m.authenticator.Authenticate(username, password, passwordHash)
}

// NewAuthToken generates a new auth token for the specified user.  The auth
// token contains a user agent to allow for multiple logins.
func (m *Manager) NewAuthToken(username string, userAgent string) (*auth.AuthToken, error) {
	tk, err := m.authenticator.GenerateToken()
	if err != nil {
		return nil, err
	}

	account, err := m.GetByUsername(username)
	if err != nil {
		return nil, err
	}

	token := &auth.AuthToken{}
	tokens := account.Tokens
	found := false

	for _, t := range tokens {
		if t.UserAgent == userAgent {
			found = true
			t.Token = tk
			token = t
			break
		}
	}

	if !found {
		token = &auth.AuthToken{
			UserAgent: userAgent,
			Token:     tk,
		}

		tokens = append(tokens, token)
	}

	// delete token
	res, err := m.t().Filter(map[string]string{"username": username}).Filter(r.Row.Field("userAgent").Eq(userAgent)).Delete().Run(m.session)
	if err != nil {
		return nil, err
	}
	defer res.Close()

	// add new token
	if _, err := m.t().Filter(map[string]string{"username": username}).Update(map[string]interface{}{"tokens": tokens}).RunWrite(m.session); err != nil {
		return nil, err
	}

	return token, nil
}

// ClearAuthTokens removes the auth tokens for the user and agent.
// If user agent is an empty string, all tokens are removed.
func (m *Manager) ClearAuthTokens(username string, userAgent string) error {
	log.Debugf("clearing tokens: username=%s agent=%q", username, userAgent)

	account, err := m.GetByUsername(username)
	if err != nil {
		return err
	}

	tokens := []*auth.AuthToken{}

	// if userAgent is specified, filter tokens to keep
	if userAgent != "" {
		for _, t := range account.Tokens {
			if t.UserAgent != userAgent {
				tokens = append(tokens, t)
			}
		}
	}

	// update tokens
	if _, err := m.t().Filter(map[string]string{"username": username}).Update(map[string]interface{}{"tokens": tokens}).RunWrite(m.session); err != nil {
		return err
	}

	return nil
}

// ChangePassword changes the password for the specified user.  Password
// must not be hashed as it will be hashed before save.
func (m *Manager) ChangePassword(username, password string) error {
	if !m.authenticator.IsUpdateSupported() {
		return fmt.Errorf("not supported for authenticator: %s", m.authenticator.Name())
	}

	hash, err := auth.Hash(password)
	if err != nil {
		return err
	}

	res, err := m.t().Filter(map[string]string{"username": username}).Update(map[string]string{"password": hash}).Run(m.session)
	if err != nil {
		return err
	}
	defer res.Close()

	log.Debugf("password updated: username=%s", username)

	return nil
}
