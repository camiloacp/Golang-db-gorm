package invoiceheader

import (
	"go-db-gorm/model"

	"gorm.io/gorm"
)

type Storage interface {
	CreateTx(*gorm.DB, *model.InvoiceHeader) error
}

type Service struct {
	storage Storage
}

func NewService(s Storage) *Service {
	return &Service{storage: s}
}
