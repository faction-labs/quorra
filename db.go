package meld

import (
	r "github.com/dancannon/gorethink"
)

func DB(dbAddr, dbName string) (*r.Session, error) {
	session, err := r.Connect(r.ConnectOpts{
		Address:  dbAddr,
		Database: dbName,
	})
	if err != nil {
		return nil, err
	}

	return session, err
}
