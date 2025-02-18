package main

import (
	"bills/internal/models"
	"bills/internal/repository"
	"database/sql"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

func main() {
	// Open database
	db, err := sql.Open("sqlite3", "bills.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// Enable foreign keys
	_, err = db.Exec("PRAGMA foreign_keys = ON")
	if err != nil {
		log.Fatal(err)
	}

	// Initialize repositories
	issuerRepo := repository.NewSQLiteIssuerRepository(db)
	receiverRepo := repository.NewSQLiteReceiverRepository(db)
	billItemRepo := repository.NewSQLiteBillItemRepository(db)

	// Seed issuers
	issuers := []*models.Issuer{
		models.NewIssuer(
			"Tech Solutions Inc.",
			"US123456789",
			"123 Silicon Valley",
			"San Francisco",
			"CA",
			"94105",
			"USA",
		),
		models.NewIssuer(
			"European Digital GmbH",
			"DE987654321",
			"456 Tech Strasse",
			"Berlin",
			"Berlin",
			"10115",
			"Germany",
		),
		models.NewIssuer(
			"Asian Tech Co., Ltd.",
			"JP456789123",
			"789 Tech Street",
			"Tokyo",
			"Tokyo",
			"100-0001",
			"Japan",
		),
	}

	log.Println("Seeding issuers...")
	for _, issuer := range issuers {
		if err := issuerRepo.Create(issuer); err != nil {
			log.Printf("Error creating issuer %s: %v\n", issuer.Name, err)
		} else {
			log.Printf("Created issuer: %s\n", issuer.Name)
		}
	}

	// Seed receivers
	receivers := []*models.Receiver{
		models.NewReceiver(
			"Global Corp",
			"UK789123456",
			"321 Business Ave",
			"London",
			"England",
			"EC1A 1BB",
			"UK",
		),
		models.NewReceiver(
			"Nordic Innovations AS",
			"NO123789456",
			"654 Innovation Road",
			"Oslo",
			"Oslo",
			"0150",
			"Norway",
		),
		models.NewReceiver(
			"Mediterranean Trade Ltd",
			"ES456123789",
			"987 Costa Street",
			"Barcelona",
			"Catalonia",
			"08001",
			"Spain",
		),
	}

	log.Println("Seeding receivers...")
	for _, receiver := range receivers {
		if err := receiverRepo.Create(receiver); err != nil {
			log.Printf("Error creating receiver %s: %v\n", receiver.Name, err)
		} else {
			log.Printf("Created receiver: %s\n", receiver.Name)
		}
	}

	// Seed bill items
	billItems := []*models.BillItem{
		models.NewBillItem(
			"Software Development Services",
			150.00,
			"EUR",
		),
		models.NewBillItem(
			"Cloud Infrastructure Setup",
			500.00,
			"USD",
		),
		models.NewBillItem(
			"Technical Consultation",
			200.00,
			"EUR",
		),
		models.NewBillItem(
			"System Maintenance",
			75.00,
			"EUR",
		),
		models.NewBillItem(
			"Data Migration Service",
			300.00,
			"USD",
		),
		models.NewBillItem(
			"Security Audit",
			450.00,
			"EUR",
		),
	}

	log.Println("Seeding bill items...")
	for _, item := range billItems {
		if err := billItemRepo.Create(item); err != nil {
			log.Printf("Error creating bill item %s: %v\n", item.Name, err)
		} else {
			log.Printf("Created bill item: %s\n", item.Name)
		}
	}

	log.Println("Seeding completed successfully!")
}
