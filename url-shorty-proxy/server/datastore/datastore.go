package datastore

import "github.com/jinzhu/gorm"

type Datastore interface {
	Create(Abbreviation) error
	Delete(shortURL string) error
	Get(id int) (*Abbreviation, error)
	GetByAbbr(shortURL string) (*Abbreviation, error)
}

type Abbreviation struct {
	gorm.Model
	Id          uint
	Shorty, Url string
}

type Store struct {
	Conn PSQLStore
}
