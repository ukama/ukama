package db

import (
	"github.com/google/uuid"
	"github.com/ukama/ukama/systems/common/sql"
	"gorm.io/gorm"
)

type SimRepo interface {
	Add(sim *Sim, nestedFunc func(*Sim, *gorm.DB) error) error
	Get(simID uuid.UUID) (*Sim, error)
	GetByUser(userID uuid.UUID) ([]Sim, error)
	GetByNetwork(NetworkID uuid.UUID) ([]Sim, error)
	Update(sim *Sim, nestedFunc func(*Sim, *gorm.DB) error) error
	Delete(simID uuid.UUID, nestedFunc func(uuid.UUID, *gorm.DB) error) error
}

type simRepo struct {
	Db sql.Db
}

func NewSimRepo(db sql.Db) *simRepo {
	return &simRepo{
		Db: db,
	}
}
