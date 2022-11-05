package db

import (
	"github.com/google/uuid"
	"github.com/ukama/ukama/systems/common/sql"
	"gorm.io/gorm"
)

type UserRepo interface {
	Add(user *User) error
	Get(uuid uuid.UUID) (*User, error)
	// Deactivate(uuid uuid.UUID) error
	Delete(uuid uuid.UUID) error
}

type userRepo struct {
	Db sql.Db
}

func NewUserRepo(db sql.Db) *userRepo {
	return &userRepo{
		Db: db,
	}
}

func (u *userRepo) Add(user *User) error {
	d := u.Db.GetGormDb().Create(user)

	return d.Error
}

func (u *userRepo) Get(uuid uuid.UUID) (*User, error) {
	var user User

	result := u.Db.GetGormDb().First(&user, uuid)
	if result.Error != nil {
		return nil, result.Error
	}

	return &user, nil
}

func (u *userRepo) Delete(userUUID uuid.UUID, nestedFunc func(uuid.UUID, *gorm.DB) error) error {
	err := u.Db.GetGormDb().Transaction(func(tx *gorm.DB) error {
		result := tx.Where(&User{Uuid: userUUID}).Delete(&User{})
		if result.Error != nil {
			return result.Error
		}

		if nestedFunc != nil {
			nestErr := nestedFunc(userUUID, tx)
			if nestErr != nil {
				return nestErr
			}
		}

		return nil
	})

	return err
}
