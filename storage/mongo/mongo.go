package mongo

import (
	"strings"

	mgo "gopkg.in/mgo.v2"
)

//NewDB creates new database connection
func NewDB(conn string, db string) (*DB, error) {
	session, err := mgo.Dial(conn)
	if err != nil {
		return nil, err
	}
	m := &DB{S: session, DB: db}
	//TODO: ensure all indexes
	return m, nil
}

//DB is database connection structure
type DB struct {
	S  *mgo.Session
	DB string
}

//Session creates new database connection session
//don't dorget to close it, defer is commonly used way to do it
func (db *DB) Session(col string) *Session {
	s := Session{}
	s.S = db.S.Clone()
	s.C = s.S.DB(db.DB).C(col)
	return &s
}

//Close closes database connection
func (db *DB) Close() {
	db.S.Close()
}

//Session implements one single session connection to database
type Session struct {
	C *mgo.Collection
	S *mgo.Session
}

//Close closes connection to current session
func (s *Session) Close() {
	s.S.Close()
}

// EnsureIndex wraps mongo EnsureIndex function. Purpose is to try recreate index if something went wrong.
func (s *Session) EnsureIndex(index mgo.Index) error {
	err := s.C.EnsureIndex(index)

	if err != nil && strings.Contains(err.Error(), "already exists with different options") {
		if err := s.C.DropIndex(index.Key...); err != nil {
			return err
		}
		return s.C.EnsureIndex(index)
	}
	return err
}
