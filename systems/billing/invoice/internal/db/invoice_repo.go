package db

import (
	"github.com/ukama/ukama/systems/common/sql"
	uuid "github.com/ukama/ukama/systems/common/uuid"
	"gorm.io/gorm"
)

type InvoiceRepo interface {
	Add(invoice *Invoice, nestedFunc func(*Invoice, *gorm.DB) error) error
	Get(id uuid.UUID) (*Invoice, error)
	GetBySubscriber(subscriberId uuid.UUID) ([]Invoice, error)

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

func (i *invoiceRepo) GetBySubscriber(subscriberId uuid.UUID) ([]Invoice, error) {
	db := i.Db.GetGormDb()
	var invoices []Invoice

	result := db.Where(&Invoice{SubscriberId: subscriberId}).Find(&invoices)
	if result.Error != nil {
		return nil, result.Error
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
