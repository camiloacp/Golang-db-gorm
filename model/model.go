package model

import (
	"gorm.io/gorm"
)

// gorm.Model es un struct embebido que agrega automáticamente 4 campos a la tabla:
//
//	ID        uint           → primary key autoincremental (SERIAL en Postgres, AUTO_INCREMENT en MySQL)
//	CreatedAt time.Time      → GORM lo setea solo al hacer Create()
//	UpdatedAt time.Time      → GORM lo actualiza solo al hacer Save() / Updates()
//	DeletedAt gorm.DeletedAt → habilita soft delete: en vez de borrar el registro,
//	                           GORM escribe la fecha en este campo y lo excluye de
//	                           todas las queries futuras (WHERE deleted_at IS NULL).
//	                           Para borrar físicamente se usa db.Unscoped().Delete()
//
// Ventaja principal: no necesitas declarar ni gestionar estos campos en ningún modelo,
// GORM los maneja de forma transparente en cada operación.

// Product representa la tabla "products"
type Product struct {
	gorm.Model
	Name         string  `gorm:"type:varchar(100);not null"`
	Observations *string `gorm:"type:varchar(100)"`
	Price        int     `gorm:"not null"`
	InvoiceItems []InvoiceItem
}

// InvoiceHeader representa la tabla "invoice_headers"
type InvoiceHeader struct {
	gorm.Model
	Client       string `gorm:"type:varchar(100);not null"`
	InvoiceItems []InvoiceItem
}

// InvoiceItem representa la tabla "invoice_items"
type InvoiceItem struct {
	gorm.Model
	InvoiceHeaderID uint `gorm:"constraint:OnDelete:RESTRICT,OnUpdate:RESTRICT"`
	ProductID       uint `gorm:"constraint:OnDelete:RESTRICT,OnUpdate:RESTRICT"`
}
