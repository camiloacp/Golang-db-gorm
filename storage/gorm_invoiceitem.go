package storage

import (
	"go-db-gorm/pkg/invoiceitem"

	"gorm.io/gorm"
)

// GormInvoiceItem implements invoiceitem.Storage using GORM
type GormInvoiceItem struct {
	db *gorm.DB
}

// NewGormInvoiceItem returns a new pointer of GormInvoiceItem
func NewGormInvoiceItem(db *gorm.DB) *GormInvoiceItem {
	return &GormInvoiceItem{db: db}
}

// CreateTx inserts each item linked to headerID inside an active transaction.
// GORM sets m.ID and m.CreatedAt automatically after each insert.
func (g *GormInvoiceItem) CreateTx(tx *gorm.DB, headerID uint, ms invoiceitem.Models) error {
	for _, m := range ms {
		m.InvoiceHeaderID = headerID
		if err := tx.Create(m).Error; err != nil {
			return err
		}
	}
	return nil
}
