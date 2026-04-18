package storage

import (
	"go-db-gorm/model"

	"gorm.io/gorm"
)

type GormInvoiceHeader struct {
	db *gorm.DB
}

func NewGormInvoiceHeader(db *gorm.DB) *GormInvoiceHeader {
	return &GormInvoiceHeader{db: db}
}

func (g *GormInvoiceHeader) CreateTx(tx *gorm.DB, m *model.InvoiceHeader) error {
	return tx.Create(m).Error
}
