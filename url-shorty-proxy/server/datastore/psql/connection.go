package psql

import (
	"fmt"
	"log"
	"os"

	"github.com/jinzhu/gorm"
)

var (
	host     = "localhost:5432"
	user     = "postgres"
	password = ""
	database = "url_shortener"
)

type Store struct {
	Logger *log.Logger
	Conn   PSQLStore
}

type PSQLStore struct {
	db *gorm.DB
}

func NewPSQLStore(logger *log.Logger) (*Store, error) {
	conn, err := psqlConnection(logger)
	if err != nil {
		return nil, err
	}

	store := &Store{Conn: PSQLStore{db: conn}}
	return store, nil
}

func psqlConnection(logger *log.Logger) (*gorm.DB, error) {

	// host is in the form of: localhost:5432
	if os.Getenv("DB_HOST") != "" {
		host = os.Getenv("DB_HOST")
	}
	if os.Getenv("DB_NAME") != "" {
		database = os.Getenv("DB_NAME")
	}
	if os.Getenv("DB_USER") != "" {
		user = os.Getenv("DB_USER")
	}
	if os.Getenv("DB_PASS") != "" {
		password = os.Getenv("DB_PASS")
	}

	log.Printf("Connecting to database using: %s - %s - %s - %s\n", host, database, user, password)

	db, err := gorm.Open(
		"postgres",
		fmt.Sprintf(
			"postgres://%s:%s@%s/%s?sslmode=disable",
			user, password, host, database,
		),
	)
	if err != nil {
		return nil, err
	}

	// Logs, among other things, raw SQL statements
	db.SetLogger(logger)
	db.LogMode(true)

	// Creates a table based on the Shorty struct. Also adds created, updated, deleted
	// timestamps
	db.AutoMigrate(&PSQLShorty{})

	// Ensures uniqueness of the shorty column
	db.Model(&PSQLShorty{}).AddUniqueIndex("idx_shorty_shorty", "shorty")

	return db, err
}
