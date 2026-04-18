package product

import (
	"errors"

	"go-db-gorm/model"
)

var ErrIDRequired = errors.New("id is required")

type Storage interface {
	Create(*model.Product) error
	Update(*model.Product) error
	GetAll() (model.Products, error)
	GetByID(uint) (*model.Product, error)
	Delete(uint) error
}

type Service struct {
	storage Storage
}

func NewService(s Storage) *Service {
	return &Service{storage: s}
}

func (s *Service) Create(m *model.Product) error {
	return s.storage.Create(m)
}

func (s *Service) GetAll() (model.Products, error) {
	return s.storage.GetAll()
}

func (s *Service) GetByID(id uint) (*model.Product, error) {
	return s.storage.GetByID(id)
}

func (s *Service) Update(m *model.Product) error {
	if m.ID == 0 {
		return ErrIDRequired
	}
	return s.storage.Update(m)
}

func (s *Service) Delete(id uint) error {
	return s.storage.Delete(id)
}
