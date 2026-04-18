package model

import (
	"fmt"
	"strings"

	"gorm.io/gorm"
)

type Product struct {
	gorm.Model
	Name         string  `gorm:"type:varchar(100);not null"`
	Observations *string `gorm:"type:varchar(100)"`
	Price        int     `gorm:"not null"`
	InvoiceItems []*InvoiceItem
}

func (Product) TableName() string { return "products" }

func (m *Product) String() string {
	obs := ""
	if m.Observations != nil {
		obs = *m.Observations
	}
	return fmt.Sprintf("%02d | %-20s | %-20s | %5d | %10s | %10s",
		m.ID, m.Name, obs, m.Price,
		m.CreatedAt.Format("2006-01-02 15:04:05"), m.UpdatedAt.Format("2006-01-02 15:04:05"))
}

type Products []*Product

func (m Products) String() string {
	builder := strings.Builder{}
	builder.WriteString(fmt.Sprintf("%02s | %-20s | %-20s | %5s | %10s | %10s\n",
		"ID", "Name", "Observations", "Price", "CreatedAt", "UpdatedAt"))
	for _, model := range m {
		builder.WriteString(model.String() + "\n")
	}
	return builder.String()
}

type InvoiceHeader struct {
	gorm.Model
	Client       string        `gorm:"type:varchar(100);not null"`
	InvoiceItems []*InvoiceItem
}

func (InvoiceHeader) TableName() string { return "invoice_headers" }

type InvoiceItem struct {
	gorm.Model
	InvoiceHeaderID uint `gorm:"not null;constraint:OnDelete:RESTRICT,OnUpdate:RESTRICT"`
	ProductID       uint `gorm:"not null;constraint:OnDelete:RESTRICT,OnUpdate:RESTRICT"`
}

func (InvoiceItem) TableName() string { return "invoice_items" }

type InvoiceItems []*InvoiceItem
