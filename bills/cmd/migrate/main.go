package main

import (
	"bills/db"
	"database/sql"
	"flag"
	"fmt"
	"log"
	"os"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

func main() {
	// Parse command line flags
	down := flag.Bool("down", false, "Drop all tables")
	dbPath := flag.String("db", "bills.db", "Database path")
	check := flag.Bool("check", false, "Check database content")
	flag.Parse()

	// Open database
	sqlDB, err := sql.Open("sqlite3", *dbPath)
	if err != nil {
		log.Fatal(err)
	}
	defer sqlDB.Close()

	// Enable foreign keys
	_, err = sqlDB.Exec("PRAGMA foreign_keys = ON")
	if err != nil {
		log.Fatal(err)
	}

	if *check {
		// Query bills
		rows, err := sqlDB.Query(`
			SELECT b.id, b.amount, b.due_date, b.paid, b.issuer_id, b.receiver_id
			FROM bills b
			ORDER BY b.id DESC
		`)
		if err != nil {
			log.Fatal(err)
		}
		defer rows.Close()

		fmt.Println("\nBills:")
		fmt.Println("----------------------------------------")
		for rows.Next() {
			var (
				id, issuerID, receiverID int64
				amount                   float64
				dueDate                  time.Time
				paid                     bool
			)
			if err := rows.Scan(&id, &amount, &dueDate, &paid, &issuerID, &receiverID); err != nil {
				log.Fatal(err)
			}
			fmt.Printf("Bill ID: %d\n", id)
			fmt.Printf("Amount: %.2f\n", amount)
			fmt.Printf("Due Date: %s\n", dueDate.Format("2006-01-02"))
			fmt.Printf("Paid: %v\n", paid)
			fmt.Printf("Issuer ID: %d\n", issuerID)
			fmt.Printf("Receiver ID: %d\n", receiverID)

			// Query bill items for this bill
			itemRows, err := sqlDB.Query(`
				SELECT a.id, a.item_id, a.quantity, a.unit_price, a.subtotal,
					   i.description
				FROM bill_item_assignments a
				LEFT JOIN bill_items i ON a.item_id = i.id
				WHERE a.bill_id = ?
			`, id)
			if err != nil {
				log.Fatal(err)
			}
			defer itemRows.Close()

			fmt.Println("\nBill Items:")
			fmt.Println("----------------------------------------")
			for itemRows.Next() {
				var (
					assignID, itemID    int64
					quantity            int
					unitPrice, subtotal float64
					description         string
				)
				if err := itemRows.Scan(&assignID, &itemID, &quantity, &unitPrice, &subtotal, &description); err != nil {
					log.Fatal(err)
				}
				fmt.Printf("Assignment ID: %d\n", assignID)
				fmt.Printf("Item ID: %d\n", itemID)
				fmt.Printf("Description: %s\n", description)
				fmt.Printf("Quantity: %d\n", quantity)
				fmt.Printf("Unit Price: %.2f\n", unitPrice)
				fmt.Printf("Subtotal: %.2f\n", subtotal)
			}
			fmt.Println("----------------------------------------\n")
		}
		os.Exit(0)
	}

	if *down {
		// Drop all tables
		if err := db.DropDB(sqlDB, *dbPath); err != nil {
			log.Fatal(err)
		}
		log.Println("Successfully dropped all tables")
	} else {
		// Apply migrations
		if err := db.MigrateDB(sqlDB, *dbPath); err != nil {
			log.Fatal(err)
		}
		log.Println("Successfully applied all migrations")
	}

	os.Exit(0)
}
