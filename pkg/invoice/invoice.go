package invoice

import (
	"go-db-gorm/model"
)

type Model struct {
	Header *model.InvoiceHeader
	Items  model.InvoiceItems
}

type Storage interface {
	Create(*Model) error
}

type Service struct {
	storage Storage
}

func NewService(s Storage) *Service {
	return &Service{s}
}

func (s *Service) Create(m *Model) error {
	return s.storage.Create(m)
}
