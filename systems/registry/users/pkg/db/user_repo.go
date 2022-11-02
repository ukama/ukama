package db

import (
	"github.com/google/uuid"
	"github.com/ukama/ukama/systems/common/sql"
	"gorm.io/gorm"
)

// UserRepo declares an interface to manage users.
// It is also useful for mocking purposes.
type UserRepo interface {
	Add(user *User, nestedFunc func(*User, *gorm.DB) error) (*User, error)
	Get(uuid uuid.UUID) (*User, error)
	Update(user *User) (*User, error)
	Delete(uuid uuid.UUID, nestedFunc func(uuid.UUID, *gorm.DB) error) error
}

type userRepo struct {
	Db sql.Db
}

func NewUserRepo(db sql.Db) *userRepo {
	return &userRepo{
		Db: db,
	}
}

// Add adds a user. Add nestedFunc to execute an action inside a transation
// if one of the nestedFunc returns error then Add action is rolled back
func (u *userRepo) Add(user *User, nestedFunc func(*User, *gorm.DB) error) (*User, error) {
	err := u.Db.GetGormDb().Transaction(func(tx *gorm.DB) error {
		d := tx.Create(user)

		if d.Error != nil {
			return d.Error
		}

		if nestedFunc != nil {
			nestErr := nestedFunc(user, tx)
			if nestErr != nil {
				return nestErr
			}
		}

		return nil
	})

	return user, err
}

func (u *userRepo) Get(uuid uuid.UUID) (*User, error) {
	var user User

	result := u.Db.GetGormDb().Where("uuid = ?", user.Uuid).First(&user)
	if result.Error != nil {
		return nil, result.Error
	}

	return &user, nil
}

// Update user modified non-empty fields provided by user struct
// Returned fields are those that were updated
func (u *userRepo) Update(user *User) (*User, error) {
	d := u.Db.GetGormDb().Where("uuid = ?", user.Uuid).UpdateColumns(user)
	if d.RowsAffected == 0 {
		return nil, gorm.ErrRecordNotFound
	}

	if d.Error != nil {
		return nil, d.Error
	}

	return user, nil
}

func (u *userRepo) Delete(userID uuid.UUID, nestedFunc func(uuid.UUID, *gorm.DB) error) error {
	err := u.Db.GetGormDb().Transaction(func(tx *gorm.DB) error {
		result := tx.Where(&User{Uuid: userID}).Delete(&User{})
		if result.Error != nil {
			return result.Error
		}

		if nestedFunc != nil {
			nestErr := nestedFunc(userID, tx)
			if nestErr != nil {
				return nestErr
			}
		}

		return nil
	})

	return err
}
