package invoiceitem

import (
	"go-db-gorm/model"

	"gorm.io/gorm"
)

type Storage interface {
	CreateTx(*gorm.DB, uint, model.InvoiceItems) error
}

type Service struct {
	storage Storage
}

func NewService(s Storage) *Service {
	return &Service{storage: s}
}
