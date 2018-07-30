package datastore

import (
	"fmt"
	"log"
	"os"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

var (
	host     = "localhost:5432"
	user     = "postgres"
	password = ""
	database = "url_shortener"
)

type PSQLShorty struct {
	gorm.Model
	Shorty, Url string
}

type PSQLStore struct {
	db *gorm.DB
}

func NewPSQLStore() (*Store, error) {
	conn, err := psqlConnection()
	if err != nil {
		return nil, err
	}

	store := &Store{Conn: PSQLStore{db: conn}}
	return store, nil
}

func psqlConnection() (*gorm.DB, error) {

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
	db.LogMode(true)

	// Creates a table based on the Shorty struct. Also adds created, updated, deleted
	// timestamps
	db.AutoMigrate(&PSQLShorty{})

	// Ensures uniqueness of the shorty column
	db.Model(&PSQLShorty{}).AddUniqueIndex("idx_shorty_shorty", "shorty")

	return db, err
}

func (st *PSQLStore) Create(a PSQLShorty) error {
	err := st.db.Create(&a).Error
	return err
}

func (st *PSQLStore) Delete(s string) error {
	var abbr *PSQLShorty

	abbr, err := st.GetByAbbr(s)
	if err != nil {
		return err
	}
	err = st.db.Unscoped().Delete(&abbr).Error
	return err
}

func (st *PSQLStore) Get(id uint) (*PSQLShorty, error) {
	var abbr PSQLShorty
	abbr.Model.ID = id

	err := st.db.First(&abbr).Error
	return &abbr, err
}

func (st *PSQLStore) GetByAbbr(shorty string) (*PSQLShorty, error) {
	var abbr PSQLShorty

	err := st.db.First(&abbr, "shorty = ?", shorty).Error

	return &abbr, err
}
