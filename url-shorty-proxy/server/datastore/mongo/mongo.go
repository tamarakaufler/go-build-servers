package mongo

import (
	"gopkg.in/mgo.v2/bson"
)

type MGOShorty struct {
	Id          bson.ObjectId `bson:"_id,omitempty"`
	Shorty, Url string
}

func (st *MGOStore) Create(a MGOShorty) error {
	return st.collection.Insert(&a)
}

func (st *MGOStore) Delete(s string) error {
	var abbr *MGOShorty

	abbr, err := st.GetByAbbr(s)
	if err != nil {
		return err
	}
	return st.collection.Remove(&abbr)
}

func (st *MGOStore) Get(id string) (*MGOShorty, error) {
	var abbr MGOShorty

	err := st.collection.FindId(bson.ObjectIdHex(id)).One(&abbr)
	return &abbr, err
}

func (st *MGOStore) GetByAbbr(shorty string) (*MGOShorty, error) {
	var abbr MGOShorty

	err := st.collection.Find(bson.M{"shorty": shorty}).One(&abbr)
	return &abbr, err
}
