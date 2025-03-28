/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */

package db

import (
	"github.com/ukama/ukama/systems/common/uuid"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"

	"github.com/ukama/ukama/systems/common/sql"
)

type MemberRepo interface {

	/* Members */
	AddMember(member *Member, orgId string, nestedFunc func(string, string) error) error
	GetMember(memberId uuid.UUID) (*Member, error)
	GetMemberByUserId(userId uuid.UUID) (*Member, error)
	GetMembers() ([]Member, error)
	UpdateMember(member *Member) error
	RemoveMember(memberId uuid.UUID, orgId string, nestedFunc func(string, string) error) error
	GetMemberCount() (int64, int64, error)
}

type memberRepo struct {
	Db sql.Db
}

func NewMemberRepo(db sql.Db) MemberRepo {
	return &memberRepo{
		Db: db,
	}
}

func (r *memberRepo) AddMember(member *Member, orgId string, nestedFunc func(string, string) error) error {
	err := r.Db.GetGormDb().Transaction(func(tx *gorm.DB) error {
		if nestedFunc != nil {
			nestErr := nestedFunc(orgId, member.UserId.String())
			if nestErr != nil {
				return nestErr
			}
		}

		d := tx.Create(member)
		if d.Error != nil {
			return d.Error
		}

		return nil
	})

	return err
}

func (r *memberRepo) GetMember(memberId uuid.UUID) (*Member, error) {
	var member Member

	result := r.Db.GetGormDb().
		Where("member_id = ?", memberId).First(&member)

	if result.Error != nil {
		return nil, result.Error
	}

	return &member, nil
}

func (r *memberRepo) GetMemberByUserId(userId uuid.UUID) (*Member, error) {
	var member Member

	result := r.Db.GetGormDb().
		Where("user_id = ?", userId).First(&member)

	if result.Error != nil {
		return nil, result.Error
	}

	return &member, nil
}

func (r *memberRepo) GetMembers() ([]Member, error) {
	var members []Member

	result := r.Db.GetGormDb().Find(&members)
	if result.Error != nil {
		return nil, result.Error
	}

	return members, nil
}

func (r *memberRepo) UpdateMember(member *Member) error {
	d := r.Db.GetGormDb().Clauses(clause.Returning{}).
		Where("member_id = ?", member.MemberId).Updates(member)

	if d.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}

	return d.Error
}

func (r *memberRepo) RemoveMember(memberId uuid.UUID, orgId string, nestedFunc func(string, string) error) error {
	err := r.Db.GetGormDb().Transaction(func(tx *gorm.DB) error {
		//TOTO: Were causing error Orgs relation not found
		// if nestedFunc != nil {
		// 	nestErr := nestedFunc(orgId, memberId.String())
		// 	if nestErr != nil {
		// 		return nestErr
		// 	}
		// }

		d := tx.Where("member_id = ?", memberId).Delete(&Member{})
		if d.Error != nil {
			return d.Error
		}

		if d.RowsAffected == 0 {
			return gorm.ErrRecordNotFound
		}

		return nil
	})

	return err
}

func (r *memberRepo) GetMemberCount() (int64, int64, error) {
	var activeMemberCount int64
	var deactiveMemberCount int64

	result := r.Db.GetGormDb().Model(&Member{}).
		Where("deactivated = ?", false).
		Count(&activeMemberCount)

	if result.Error != nil {
		return 0, 0, result.Error
	}

	result = r.Db.GetGormDb().Model(&Member{}).
		Where("deactivated = ?", true).
		Count(&deactiveMemberCount)

	if result.Error != nil {
		return 0, 0, result.Error
	}

	return activeMemberCount, deactiveMemberCount, nil
}
