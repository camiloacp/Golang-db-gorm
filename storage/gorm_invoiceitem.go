package storage

import (
	"go-db-gorm/model"

	"gorm.io/gorm"
)

type GormInvoiceItem struct {
	db *gorm.DB
}

func NewGormInvoiceItem(db *gorm.DB) *GormInvoiceItem {
	return &GormInvoiceItem{db: db}
}

func (g *GormInvoiceItem) CreateTx(tx *gorm.DB, headerID uint, ms model.InvoiceItems) error {
	for _, m := range ms {
		m.InvoiceHeaderID = headerID
		if err := tx.Create(m).Error; err != nil {
			return err
		}
	}
	return nil
}
