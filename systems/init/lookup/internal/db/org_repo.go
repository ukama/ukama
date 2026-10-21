/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */

package db

import (
	"context"

	"github.com/ukama/ukama/systems/common/sql"
	"github.com/ukama/ukama/systems/common/uuid"
)

type OrgRepo interface {
	Add(ctx context.Context, org *Org) error
	Update(ctx context.Context, org *Org) error
	GetAll(ctx context.Context) ([]Org, error)
	GetByName(ctx context.Context, name string) (*Org, error)
	GetById(ctx context.Context, id uuid.UUID) (*Org, error)
}

type orgRepo struct {
	Db sql.Db
}

func NewOrgRepo(db sql.Db) *orgRepo {
	return &orgRepo{
		Db: db,
	}
}

func (r *orgRepo) Add(ctx context.Context, org *Org) error {

	d := r.Db.GetGormDb().WithContext(ctx).Create(org)
	return d.Error
}

func (r *orgRepo) Update(ctx context.Context, org *Org) error {
	d := r.Db.GetGormDb().WithContext(ctx).Where(&Org{Name: org.Name}).Updates(org)
	return d.Error
}

func (r *orgRepo) GetByName(ctx context.Context, name string) (*Org, error) {
	org := &Org{}
	d := r.Db.GetGormDb().WithContext(ctx).Where(&Org{Name: name}).First(&org)
	return org, d.Error
}

func (r *orgRepo) GetAll(ctx context.Context) ([]Org, error) {
	var orgs []Org
	d := r.Db.GetGormDb().WithContext(ctx).Find(&orgs)
	return orgs, d.Error
}

func (r *orgRepo) GetById(ctx context.Context, id uuid.UUID) (*Org, error) {
	org := &Org{}
	d := r.Db.GetGormDb().WithContext(ctx).Where(&Org{OrgId: id}).First(&org)
	return org, d.Error
}
