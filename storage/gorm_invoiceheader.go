package storage

import (
	"go-db-gorm/pkg/invoiceheader"

	"gorm.io/gorm"
)

// GormInvoiceHeader implements invoiceheader.Storage using GORM
type GormInvoiceHeader struct {
	db *gorm.DB
}

// NewGormInvoiceHeader returns a new pointer of GormInvoiceHeader
func NewGormInvoiceHeader(db *gorm.DB) *GormInvoiceHeader {
	return &GormInvoiceHeader{db: db}
}

// CreateTx inserts a new invoice header inside an active transaction.
// GORM sets m.ID and m.CreatedAt automatically after the insert.
func (g *GormInvoiceHeader) CreateTx(tx *gorm.DB, m *invoiceheader.Model) error {
	return tx.Create(m).Error
}
