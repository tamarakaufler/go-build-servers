package mongo

import (
	"log"
	"os"

	"gopkg.in/mgo.v2"
)

var (
	host       = "localhost:27017"
	user       = "shorty_user"
	password   = "shortypass"
	database   = "url_shortener"
	collection = "shorties"
)

type MGOStore struct {
	session    *mgo.Session
	collection *mgo.Collection
}

type Store struct {
	Conn MGOStore
}

func NewMGOStore() (*Store, error) {

	// host is in the form of: localhost:21017
	if os.Getenv("DB_HOST") != "" {
		host = os.Getenv("DB_HOST")
	}
	if os.Getenv("DB_NAME") != "" {
		database = os.Getenv("DB_NAME")
	}
	if os.Getenv("COLL_NAME") != "" {
		collection = os.Getenv("COLL_NAME")
	}
	if os.Getenv("DB_USER") != "" {
		user = os.Getenv("DB_USER")
	}
	if os.Getenv("DB_PASS") != "" {
		password = os.Getenv("DB_PASS")
	}

	log.Printf("Connecting to database using: %s - %s - %s - %s\n", host, database, user, password)

	sess, err := mgo.Dial(host)
	if err != nil {
		return nil, err
	}

	coll := sess.DB(database).C(collection)
	coll.EnsureIndex(mgo.Index{
		Key:        []string{"shorty"},
		Unique:     true,
		DropDups:   true,
		Background: true,
		Sparse:     true,
	})

	mgoStore := MGOStore{
		session:    sess,
		collection: coll,
	}
	store := &Store{
		Conn: mgoStore,
	}

	return store, nil
}

func (s *MGOStore) Copy() *mgo.Session {
	return s.session.Copy()
}

func (s *MGOStore) GetCollection(db string, col string) *mgo.Collection {
	return s.session.DB(db).C(col)
}

func (s *MGOStore) Close() {
	if s.session != nil {
		s.session.Close()
	}
}
