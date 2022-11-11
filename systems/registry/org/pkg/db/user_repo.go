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

func NewUserRepo(db sql.Db) UserRepo {
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

	result := u.Db.GetGormDb().Where("uuid = ?", uuid).First(&user)
	if result.Error != nil {
		return nil, result.Error
	}

	return &user, nil
}

func (u *userRepo) Delete(userUUID uuid.UUID) error {
	err := u.Db.GetGormDb().Transaction(func(tx *gorm.DB) error {
		result := tx.Where(&User{Uuid: userUUID}).Delete(&User{})

		if result.Error != nil {
			return result.Error
		}

		return nil
	})

	return err
}
