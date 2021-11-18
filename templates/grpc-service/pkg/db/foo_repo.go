// This is an example of a repository
//
package db

import (
	"log"

	"github.com/ukama/ukamaX/common/sql"
	"gorm.io/gorm/clause"
)

// declare interface so that we can mock it
type FooRepo interface {
	AddOrUpdate(foo *Foo) error
	Get(id int) (*Foo, error)
}

type fooRepo struct {
	Db sql.Db
}

func NewFooRepo(db sql.Db) *fooRepo {
	return &fooRepo{
		Db: db,
	}
}

func (r *fooRepo) Add(foo *Foo) error {
	d := r.Db.GetGormDb().Create(foo)
	log.Fatal("Not implemented")
	return d.Error
}

func (r *fooRepo) Get(id int) (*Foo, error) {
	var foo Foo
	result := r.Db.GetGormDb().Preload(clause.Associations).First(&foo, id)
	if result.Error != nil {
		return nil, result.Error
	}

	log.Fatal("Not implemented")
	return &foo, nil
}
