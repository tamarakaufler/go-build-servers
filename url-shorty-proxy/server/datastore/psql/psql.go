package psql

import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

type PSQLShorty struct {
	gorm.Model
	Shorty, Url string
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

func (st *PSQLStore) GetAll() ([]PSQLShorty, error) {
	var abbrs []PSQLShorty
	err := st.db.Find(&abbrs).Error

	return abbrs, err
}
