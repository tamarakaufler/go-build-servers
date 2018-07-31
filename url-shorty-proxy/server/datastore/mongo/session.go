package mongo

import (
	"log"
	"os"
	"time"

	"gopkg.in/mgo.v2"
)

var (
	host       = "localhost:27017"
	user       = "shorty_user"
	password   = "shortypass"
	database   = "url_shortener"
	collection = "shorties"
)

type Store struct {
	Conn MGOStore
}

type MGOStore struct {
	session              *mgo.Session
	database, collection string
}

func NewMGOStore(logger *log.Logger) (*Store, error) {

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

	//mgo.SetDebug(true)
	mgo.SetLogger(logger)

	// To establish an authenticated session to the database.
	// Provides cluster info.
	connInfo := &mgo.DialInfo{
		Addrs:    []string{host},
		Timeout:  60 * time.Second,
		Database: database,
		Username: user,
		Password: password,
	}

	// Create a session which maintains a pool of socket connections
	// to MongoDB cluster (Addrs)
	db, err := mgo.DialWithInfo(connInfo)
	if err != nil {
		return nil, err
	}

	// Read from a slave if possible
	db.SetMode(mgo.Monotonic, true)

	// Create an index on the most searched field
	shortiesIndex := mgo.Index{
		Key:        []string{"shorty"},
		Unique:     true,
		DropDups:   true,
		Background: true,
		Sparse:     true,
	}
	err = db.DB(database).C(collection).EnsureIndex(shortiesIndex)
	if err != nil {
		return nil, err
	}

	mgoStore := MGOStore{
		session:    db,
		database:   database,
		collection: collection,
	}
	store := &Store{
		Conn: mgoStore,
	}
	return store, nil
}
