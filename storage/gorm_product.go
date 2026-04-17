package storage

import (
	"fmt"

	"go-db-gorm/pkg/product"

	"gorm.io/gorm"
)

// GormProduct implements product.Storage using GORM
type GormProduct struct {
	db *gorm.DB
}

// NewGormProduct returns a new pointer of GormProduct
func NewGormProduct(db *gorm.DB) *GormProduct {
	return &GormProduct{db: db}
}

// Create inserts a new product record
func (g *GormProduct) Create(m *product.Model) error {
	return g.db.Create(m).Error
}

// GetAll returns every product in the table
func (g *GormProduct) GetAll() (product.Models, error) {
	var ms product.Models
	if err := g.db.Find(&ms).Error; err != nil {
		return nil, err
	}
	return ms, nil
}

// GetByID returns a single product by primary key
func (g *GormProduct) GetByID(id uint) (*product.Model, error) {
	m := &product.Model{}
	if err := g.db.First(m, id).Error; err != nil {
		return nil, err
	}
	return m, nil
}

// Update saves only the non-zero fields of the given product
func (g *GormProduct) Update(m *product.Model) error {
	result := g.db.Model(m).Updates(m)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return fmt.Errorf("not exist product with ID: %d", m.ID)
	}
	return nil
}

// Delete soft-deletes a product by primary key
func (g *GormProduct) Delete(id uint) error {
	result := g.db.Delete(&product.Model{}, id)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return fmt.Errorf("not exist product with ID: %d", id)
	}
	return nil
}
