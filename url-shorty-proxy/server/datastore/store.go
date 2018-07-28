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

type PSQLStore struct {
	db *gorm.DB
}

func NewPSQLStore() (*Store, error) {
	conn, err := DBConnection()
	if err != nil {
		return nil, err
	}

	store := &Store{Conn: PSQLStore{db: conn}}
	return store, nil
}

func DBConnection() (*gorm.DB, error) {

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

	// Logs, among other things raw SQL statements
	db.LogMode(true)

	// Creates a table based on the Shorty struct. Also adds created, updated, deleted
	// timestamps
	db.AutoMigrate(&Shorty{})

	// Ensures uniqueness of the shorty column
	db.Model(&Shorty{}).AddUniqueIndex("idx_shorty_shorty", "shorty")

	return db, err
}

func (st *PSQLStore) Create(a Shorty) error {
	err := st.db.Create(a).Error
	return err
}

func (st *PSQLStore) Delete(s string) error {
	var abbr Shorty
	err := st.db.Where("shorty=?", s).Delete(&abbr).Error
	return err
}

func (st *PSQLStore) Get(id uint) (*Shorty, error) {
	var abbr Shorty
	abbr.Model.ID = id

	err := st.db.First(&abbr).Error
	return &abbr, err
}

func (st *PSQLStore) GetByAbbr(shorty string) (*Shorty, error) {
	var abbr Shorty

	err := st.db.First(&abbr, "shorty = ?", shorty).Error

	return &abbr, err
}
