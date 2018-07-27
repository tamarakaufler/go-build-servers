package datastore

import "github.com/jinzhu/gorm"

type Datastore interface {
	Create(Shorty) error
	Delete(shortURL string) error
	Get(id int) (*Shorty, error)
	GetByAbbr(shortURL string) (*Shorty, error)
}

type Shorty struct {
	gorm.Model
	Id          uint
	Shorty, Url string
}

type Store struct {
	Conn PSQLStore
}
