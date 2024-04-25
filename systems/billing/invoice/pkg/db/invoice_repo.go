/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */

package db

import (
	"gorm.io/gorm"
	"gorm.io/gorm/clause"

	"github.com/ukama/ukama/systems/common/sql"
	"github.com/ukama/ukama/systems/common/uuid"
)

type InvoiceRepo interface {
	Add(invoice *Invoice, nestedFunc func(*Invoice, *gorm.DB) error) error
	Get(id uuid.UUID) (*Invoice, error)
	List(invoiceeId string, invoiceeType InvoiceeType, networkId string,
		isPaid bool, count uint32, sort bool) ([]Invoice, error)

	// Update(orgId uint, network *Network) error
	Delete(invoiceId uuid.UUID, nestedFunc func(uuid.UUID, *gorm.DB) error) error
}

type invoiceRepo struct {
	Db sql.Db
}

func NewInvoiceRepo(db sql.Db) InvoiceRepo {
	return &invoiceRepo{
		Db: db,
	}
}

func (i *invoiceRepo) Add(invoice *Invoice, nestedFunc func(invoice *Invoice, tx *gorm.DB) error) error {
	err := i.Db.GetGormDb().Transaction(func(tx *gorm.DB) error {
		if nestedFunc != nil {
			nestErr := nestedFunc(invoice, tx)
			if nestErr != nil {
				return nestErr
			}
		}

		result := tx.Create(invoice)
		if result.Error != nil {
			return result.Error
		}

		return nil
	})

	return err
}

func (i *invoiceRepo) Get(id uuid.UUID) (*Invoice, error) {
	var inv Invoice

	result := i.Db.GetGormDb().First(&inv, id)
	if result.Error != nil {
		return nil, result.Error
	}

	return &inv, nil
}

func (i *invoiceRepo) List(invoiceeId string, invoiceeType InvoiceeType, networkId string,
	isPaid bool, count uint32, sort bool) ([]Invoice, error) {
	invoices := []Invoice{}

	tx := i.Db.GetGormDb().Preload(clause.Associations)

	if invoiceeId != "" {
		tx = tx.Where("invoicee_id = ?", invoiceeId)
	}

	if invoiceeType != InvoiceeTypeUnknown {
		tx = tx.Where("invoicee_type = ?", invoiceeType)
	}

	if networkId != "" {
		tx = tx.Where("network_id = ?", networkId)
	}

	if isPaid {
		tx = tx.Where("is_paid = ?", isPaid)
	}

	if sort {
		tx = tx.Order("time DESC")
	}

	if count > 0 {
		tx = tx.Limit(int(count))
	}

	result := tx.Find(&invoices)
	if result.Error != nil {
		return nil, result.Error
	}

	if result.RowsAffected == 0 {
		return nil, gorm.ErrRecordNotFound
	}

	return invoices, nil
}

func (i *invoiceRepo) Delete(invoiceId uuid.UUID, nestedFunc func(uuid.UUID, *gorm.DB) error) error {
	err := i.Db.GetGormDb().Transaction(func(tx *gorm.DB) error {
		result := tx.Delete(&Invoice{}, invoiceId)
		if result.Error != nil {
			return result.Error
		}

		if nestedFunc != nil {
			nestErr := nestedFunc(invoiceId, tx)
			if nestErr != nil {
				return nestErr
			}
		}

		return nil
	})

	return err
}
