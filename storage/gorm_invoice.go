package storage

import (
	"fmt"

	"go-db-gorm/pkg/invoice"
	"go-db-gorm/pkg/invoiceheader"
	"go-db-gorm/pkg/invoiceitem"

	"gorm.io/gorm"
)

// GormInvoice implements invoice.Storage using GORM
type GormInvoice struct {
	db            *gorm.DB
	storageHeader invoiceheader.Storage
	storageItem   invoiceitem.Storage
}

// NewGormInvoice returns a new pointer of GormInvoice
func NewGormInvoice(db *gorm.DB, h invoiceheader.Storage, i invoiceitem.Storage) *GormInvoice {
	return &GormInvoice{
		db:            db,
		storageHeader: h,
		storageItem:   i,
	}
}

// Create inserts a full invoice (header + items) inside a single transaction.
// On any error GORM automatically rolls back the transaction.
func (g *GormInvoice) Create(m *invoice.Model) error {
	return g.db.Transaction(func(tx *gorm.DB) error {
		if err := g.storageHeader.CreateTx(tx, m.Header); err != nil {
			return fmt.Errorf("header: %w", err)
		}

		if err := g.storageItem.CreateTx(tx, m.Header.ID, m.Items); err != nil {
			return fmt.Errorf("items: %w", err)
		}
		return nil
	})
}
