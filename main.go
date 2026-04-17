package main

import (
	"log"

	"go-db-gorm/pkg/invoice"
	"go-db-gorm/pkg/invoiceheader"
	"go-db-gorm/pkg/invoiceitem"
	"go-db-gorm/storage"
)

func main() {
	storage.New(storage.MySQL)

	if err := storage.Migrate(); err != nil {
		log.Fatalf("migration failed: %v", err)
	}

	storageHeader := storage.NewGormInvoiceHeader(storage.DB())
	storageItem := storage.NewGormInvoiceItem(storage.DB())
	storageInvoice := storage.NewGormInvoice(storage.DB(), storageHeader, storageItem)

	svcInvoice := invoice.NewService(storageInvoice)

	inv := &invoice.Model{
		Header: &invoiceheader.Model{Client: "Katherine Sanchez"},
		Items: invoiceitem.Models{
			{ProductID: 1},
			{ProductID: 2},
			{ProductID: 3},
		},
	}

	if err := svcInvoice.Create(inv); err != nil {
		log.Fatalf("create invoice failed: %v", err)
	}
}
