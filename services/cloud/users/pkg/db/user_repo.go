package db

import (
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"github.com/ukama/ukamaX/common/sql"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// declare interface so that we can mock it
type UserRepo interface {
	Add(user *User, orgName string, nestedFunc func(*User, *gorm.DB) error) (*User, error)
	Get(uuid uuid.UUID) (*User, error)
	Delete(uuid uuid.UUID, nestedFunc func(uuid.UUID, *gorm.DB) error) error
	GetByOrg(orgName string) ([]User, error)
	IsOverTheLimit(orgName string) (bool, error)
	// Update user modifiable fields
	// user ID here is used as identifier of a user that should be updated
	Update(user *User) (*User, error)
}

type userRepo struct {
	Db sql.Db
}

func NewUserRepo(db sql.Db) *userRepo {
	return &userRepo{
		Db: db,
	}
}

func (r *userRepo) Add(user *User, orgName string, nestedFunc func(*User, *gorm.DB) error) (*User, error) {
	org, err := makeUserOrgExist(r.Db.GetGormDb(), orgName)
	if err != nil {
		return nil, err
	}
	user.Org = org

	err = r.Db.GetGormDb().Transaction(func(tx *gorm.DB) error {
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
	user := User{}
	result := u.Db.GetGormDb().Preload(clause.Associations).Where(&User{Uuid: uuid}).First(&user)
	if result.Error != nil {
		return nil, result.Error
	}

	return &user, nil
}

func (u *userRepo) Delete(userId uuid.UUID, nestedFunc func(uuid.UUID, *gorm.DB) error) error {

	err := u.Db.GetGormDb().Transaction(func(tx *gorm.DB) error {
		result := tx.Where(&User{Uuid: userId}).Delete(&User{})
		if result.Error != nil {
			return result.Error
		}

		if result.Error != nil {
			return result.Error
		}

		if nestedFunc != nil {
			nestErr := nestedFunc(userId, tx)
			if nestErr != nil {
				return nestErr
			}
		}

		return nil
	})

	return err
}

func (u *userRepo) GetByOrg(orgName string) ([]User, error) {
	var users []User
	result := u.Db.GetGormDb().Raw("select u.* from users u "+
		"inner join orgs o on o.id = u.org_id "+
		"where o.name=? and u.deleted_at is null", orgName).Scan(&users)

	if result.Error != nil {
		return nil, result.Error
	}
	return users, nil
}

// check if number of users limit reached
func (u *userRepo) IsOverTheLimit(org string) (bool, error) {
	logrus.Debugf("Checking if user limit reached for org %s", org)
	limit := []struct {
		Limit *int `gorm:"column:limit"`
	}{}

	d := u.Db.GetGormDb().Raw(`select o.user_limit - count(s.*) as limit from orgs o
										inner join users u on o.id = u.org_id
										inner join simcards s on u.id = s.user_id
									where o.name = ?
									group by o.id`, org).Scan(&limit)

	if d.Error != nil {
		logrus.Debugf("Error checking user limit for org %s: %s", org, d.Error)
		return false, d.Error
	}

	if len(limit) == 0 {
		logrus.Debugf("No rows returned")
		return false, nil
	}

	if limit[0].Limit != nil && *limit[0].Limit <= 0 {
		logrus.Infoln("limit reached")
		return true, nil
	}

	if limit[0].Limit != nil {
		logrus.Infof("User places left : %v", limit[0].Limit)
	}

	return false, nil
}

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
