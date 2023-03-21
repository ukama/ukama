package db

import (
	"github.com/ukama/ukama/systems/common/sql"
	"github.com/ukama/ukama/systems/common/uuid"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type UserRepo interface {
	Add(user *User, nestedFunc func(*User, *gorm.DB) error) error
	Get(uuid uuid.UUID) (*User, error)
	Update(user *User, nestedFunc func(*User, *gorm.DB) error) error
	Delete(uuid uuid.UUID, nestedFunc func(uuid.UUID, *gorm.DB) error) error
	GetUserCount() (int64, int64, error)
}

type userRepo struct {
	Db sql.Db
}

func NewUserRepo(db sql.Db) *userRepo {
	return &userRepo{
		Db: db,
	}
}

func (u *userRepo) Add(user *User, nestedFunc func(user *User, tx *gorm.DB) error) error {
	err := u.Db.GetGormDb().Transaction(func(tx *gorm.DB) error {
		result := tx.Create(user)
		if result.Error != nil {
			return result.Error
		}

		if nestedFunc != nil {
			nestErr := nestedFunc(user, tx)
			if nestErr != nil {
				return nestErr
			}
		}

		return nil
	})

	return err
}

func (u *userRepo) Get(uuid uuid.UUID) (*User, error) {
	var user User

	result := u.Db.GetGormDb().Where("uuid = ?", uuid).First(&user)
	if result.Error != nil {
		return nil, result.Error
	}

	return &user, nil
}

// Update user modified non-empty fields provided by user struct
func (u *userRepo) Update(user *User, nestedFunc func(*User, *gorm.DB) error) error {
	err := u.Db.GetGormDb().Transaction(func(tx *gorm.DB) error {
		result := tx.Clauses(clause.Returning{}).Where("uuid = ?", user.Uuid).Updates(user)

		if result.RowsAffected == 0 {
			return gorm.ErrRecordNotFound
		}

		if result.Error != nil {
			return result.Error
		}

		if nestedFunc != nil {
			nestErr := nestedFunc(user, tx)
			if nestErr != nil {
				return nestErr
			}
		}

		return nil
	})

	return err
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

func (u *userRepo) GetUserCount() (int64, int64, error) {
	var userCount int64
	var deactiveUserCount int64

	result := u.Db.GetGormDb().Model(&User{}).Count(&userCount)
	if result.Error != nil {
		return 0, 0, result.Error
	}

	result = u.Db.GetGormDb().Model(&User{}).Where("deactivated = ?", true).Count(&deactiveUserCount)
	if result.Error != nil {
		return 0, 0, result.Error
	}

	return userCount, deactiveUserCount, nil
}
