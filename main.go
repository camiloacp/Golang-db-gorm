package main

import (
	"log"

	"go-db-gorm/model"
	"go-db-gorm/pkg/invoice"
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
		Header: &model.InvoiceHeader{Client: "Katherine Sanchez"},
		Items: model.InvoiceItems{
			{ProductID: 1},
			{ProductID: 2},
			{ProductID: 3},
		},
	}

	if err := svcInvoice.Create(inv); err != nil {
		log.Fatalf("create invoice failed: %v", err)
	}
}
