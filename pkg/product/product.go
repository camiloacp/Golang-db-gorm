package product

import (
	"errors"
	"fmt"
	"strings"

	"gorm.io/gorm"
)

var ErrIDRequired = errors.New("id is required")

// Model of product
type Model struct {
	gorm.Model
	Name         string  `gorm:"type:varchar(100);not null"`
	Observations *string `gorm:"type:varchar(100)"`
	Price        int     `gorm:"not null"`
}

// TableName overrides the default GORM table name
func (Model) TableName() string { return "products" }

func (m *Model) String() string {
	obs := ""
	if m.Observations != nil {
		obs = *m.Observations
	}
	return fmt.Sprintf("%02d | %-20s | %-20s | %5d | %10s | %10s",
		m.ID, m.Name, obs, m.Price,
		m.CreatedAt.Format("2006-01-02 15:04:05"), m.UpdatedAt.Format("2006-01-02 15:04:05"))
}

// Models slice of Model
type Models []*Model

func (m Models) String() string {
	builder := strings.Builder{}
	builder.WriteString(fmt.Sprintf("%02s | %-20s | %-20s | %5s | %10s | %10s\n",
		"ID", "Name", "Observations", "Price", "CreatedAt", "UpdatedAt"))
	for _, model := range m {
		builder.WriteString(model.String() + "\n")
	}
	return builder.String()
}

// Storage interface that must be implemented by the storage layer
type Storage interface {
	Create(*Model) error
	Update(*Model) error
	GetAll() (Models, error)
	GetByID(uint) (*Model, error)
	Delete(uint) error
}

// Service of product
type Service struct {
	storage Storage
}

// NewService return a new pointer of Service
func NewService(s Storage) *Service {
	return &Service{storage: s}
}

// Create a new product
func (s *Service) Create(m *Model) error {
	return s.storage.Create(m)
}

// GetAll returns all products
func (s *Service) GetAll() (Models, error) {
	return s.storage.GetAll()
}

// GetByID returns a product by ID
func (s *Service) GetByID(id uint) (*Model, error) {
	return s.storage.GetByID(id)
}

// Update a product
func (s *Service) Update(m *Model) error {
	if m.ID == 0 {
		return ErrIDRequired
	}
	return s.storage.Update(m)
}

// Delete a product by ID
func (s *Service) Delete(id uint) error {
	return s.storage.Delete(id)
}
