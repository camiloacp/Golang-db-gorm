package storage

import (
	"fmt"

	"go-db-gorm/model"

	"gorm.io/gorm"
)

type GormProduct struct {
	db *gorm.DB
}

func NewGormProduct(db *gorm.DB) *GormProduct {
	return &GormProduct{db: db}
}

func (g *GormProduct) Create(m *model.Product) error {
	return g.db.Create(m).Error
}

func (g *GormProduct) GetAll() (model.Products, error) {
	var ms model.Products
	if err := g.db.Find(&ms).Error; err != nil {
		return nil, err
	}
	return ms, nil
}

func (g *GormProduct) GetByID(id uint) (*model.Product, error) {
	m := &model.Product{}
	if err := g.db.First(m, id).Error; err != nil {
		return nil, err
	}
	return m, nil
}

func (g *GormProduct) Update(m *model.Product) error {
	result := g.db.Model(m).Updates(m)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return fmt.Errorf("not exist product with ID: %d", m.ID)
	}
	return nil
}

func (g *GormProduct) Delete(id uint) error {
	result := g.db.Delete(&model.Product{}, id)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return fmt.Errorf("not exist product with ID: %d", id)
	}
	return nil
}
