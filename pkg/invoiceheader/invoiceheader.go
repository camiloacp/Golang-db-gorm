package invoiceheader

import (
	"gorm.io/gorm"
)

// Model of invoiceheader
type Model struct {
	gorm.Model
	Client string `gorm:"type:varchar(100);not null"`
}

// TableName overrides the default GORM table name
func (Model) TableName() string { return "invoice_headers" }

// Storage interface that must be implemented by the storage layer
type Storage interface {
	CreateTx(*gorm.DB, *Model) error
}

// Service of invoiceheader
type Service struct {
	storage Storage
}

// NewService return a new pointer of Service
func NewService(s Storage) *Service {
	return &Service{storage: s}
}
