package invoiceitem

import (
	"gorm.io/gorm"
)

// Model of invoiceitem
type Model struct {
	gorm.Model
	InvoiceHeaderID uint `gorm:"not null;constraint:OnDelete:RESTRICT,OnUpdate:RESTRICT"`
	ProductID       uint `gorm:"not null;constraint:OnDelete:RESTRICT,OnUpdate:RESTRICT"`
}

// TableName overrides the default GORM table name
func (Model) TableName() string { return "invoice_items" }

// Models slice of Model
type Models []*Model

// Storage interface that must be implemented by the storage layer
type Storage interface {
	CreateTx(*gorm.DB, uint, Models) error
}

// Service of invoiceitem
type Service struct {
	storage Storage
}

// NewService return a new pointer of Service
func NewService(s Storage) *Service {
	return &Service{storage: s}
}
