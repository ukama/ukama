package db

import (
	uuid "github.com/satori/go.uuid"
	"github.com/ukama/ukamaX/common/sql"
	"gorm.io/gorm/clause"
)

// declare interface so that we can mock it
type UserRepo interface {
	Add(user *User) (*User, error)
	Get(uuid uuid.UUID) (*User, error)
	Delete(uuid uuid.UUID) error
	GetByOrg(orgName string) ([]User, error)
}

type userRepo struct {
	Db sql.Db
}

func NewUserRepo(db sql.Db) *userRepo {
	return &userRepo{
		Db: db,
	}
}

func (u *userRepo) Add(user *User) (*User, error) {
	d := u.Db.GetGormDb().Create(user)
	return user, d.Error
}

func (u *userRepo) Get(uuid uuid.UUID) (*User, error) {
	user := User{}
	result := u.Db.GetGormDb().Preload(clause.Associations).Where(&User{UUID: uuid}).First(&user)
	if result.Error != nil {
		return nil, result.Error
	}

	return &user, nil
}

func (u *userRepo) Delete(uuid uuid.UUID) error {
	result := u.Db.GetGormDb().Where(&User{UUID: uuid}).Delete(&User{})
	if result.Error != nil {
		return result.Error
	}

	return nil
}

func (u *userRepo) GetByOrg(orgName string) ([]User, error) {
	var users []User
	result := u.Db.GetGormDb().Raw("select u.* from users u "+
		"inner join imsis i on i.user_uuid = u.uuid "+
		"inner join orgs o on o.id = i.org_id "+
		"where o.name=? and o.deleted_at is null", orgName).Scan(&users)

	if result.Error != nil {
		return nil, result.Error
	}
	return users, nil
}
