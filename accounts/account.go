package accounts

import (
	"github.com/factionlabs/quorra/auth"
)

type Account struct {
	Id        string            `json:"id,omitempty" gorethink:"id,omitempty"`
	Username  string            `json:"username" gorethink:"username"`
	FirstName string            `json:"firstName" gorethink:"firstName"`
	LastName  string            `json:"lastName" gorethink:"lastName"`
	Email     string            `json:"email" gorethink:"email"`
	Password  string            `json:"-" gorethink:"password"`
	Tokens    []*auth.AuthToken `json:"-" gorethink:"tokens"`
	Created   int64             `json:"created" gorethink:"created"`
	IsAdmin   bool              `json:"isAdmin" gorethink:"isAdmin"`
	Labels    []string          `json:"labels" gorethink:"labels"`
}
