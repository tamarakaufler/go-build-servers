package datastore

type Datastore interface {
	Create(Shorty) error
	Delete(shortURL string) error
	Get(id int) (*Shorty, error)
	GetByAbbr(shortURL string) (*Shorty, error)
}

type Shorty struct {
	Shorty, Url string
}
