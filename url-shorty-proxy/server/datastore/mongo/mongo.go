package mongo

import (
	"gopkg.in/mgo.v2/bson"
)

type MGOShorty struct {
	Id          bson.ObjectId `bson:"_id,omitempty"`
	Shorty, Url string
}

func (st *MGOStore) Create(a MGOShorty) error {
	sess := st.session.Copy()
	defer sess.Close()

	return sess.DB(st.database).C(st.collection).Insert(&a)
}

func (st *MGOStore) Delete(shorty string) error {
	sess := st.session.Copy()
	defer sess.Close()

	return sess.DB(st.database).C(st.collection).Remove(bson.M{"shorty": shorty})
}

func (st *MGOStore) Get(id string) (*MGOShorty, error) {
	sess := st.session.Copy()
	defer sess.Close()

	var abbr MGOShorty

	err := sess.DB(st.database).C(st.collection).FindId(bson.ObjectIdHex(id)).One(&abbr)
	return &abbr, err
}

func (st *MGOStore) GetByAbbr(shorty string) (*MGOShorty, error) {
	sess := st.session.Copy()
	defer sess.Close()

	var abbr MGOShorty

	err := sess.DB(st.database).C(st.collection).Find(bson.M{"shorty": shorty}).One(&abbr)
	return &abbr, err
}

func (st *MGOStore) GetAll() ([]MGOShorty, error) {
	sess := st.session.Copy()
	defer sess.Close()

	var abbrs []MGOShorty

	err := sess.DB(st.database).C(st.collection).Find(nil).Sort("shorty").All(&abbrs)

	return abbrs, err

}
