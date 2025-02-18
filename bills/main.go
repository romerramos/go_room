package main

import (
	"bills/db"
	"bills/internal/handlers"
	"bills/internal/repository"
	"database/sql"
	"html/template"
	"io"
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	_ "github.com/mattn/go-sqlite3"
)

// Template renderer for Echo
type Template struct {
	templates *template.Template
}

func (t *Template) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	// List of partial templates that should be rendered directly
	partials := map[string]bool{
		"bills-list":        true,
		"bill-items-list":   true,
		"issuers-list":      true,
		"receivers-list":    true,
		"bill-items-select": true,
		"issuers-select":    true,
		"receivers-select":  true,
	}

	// If it's a partial template, render it directly
	if partials[name] {
		return t.templates.ExecuteTemplate(w, name, data)
	}

	// For full pages, render the layout with the data
	return t.templates.ExecuteTemplate(w, name, data)
}

func main() {
	// Load environment variables
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found")
	}

	// Initialize database
	dbPath := "bills.db"
	sqlDB, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		log.Fatal(err)
	}
	defer sqlDB.Close()

	// Run migrations
	if err := db.MigrateDB(sqlDB, dbPath); err != nil {
		log.Fatal(err)
	}

	// Initialize repositories
	billRepo := repository.NewSQLiteBillRepository(sqlDB)
	receiverRepo := repository.NewSQLiteReceiverRepository(sqlDB)
	issuerRepo := repository.NewSQLiteIssuerRepository(sqlDB)
	billItemRepo := repository.NewSQLiteBillItemRepository(sqlDB)
	billItemAssignmentRepo := repository.NewSQLiteBillItemAssignmentRepository(sqlDB)

	// Initialize Echo
	e := echo.New()

	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.CORS())

	// Initialize templates - each template is now a complete HTML file
	t := &Template{
		templates: template.Must(template.New("").Funcs(template.FuncMap{
			"now": time.Now,
			"date": func(format string, t time.Time) string {
				return t.Format(format)
			},
		}).ParseFiles(
			"templates/bills.html",
			"templates/bills-list.html",
			"templates/bill-form.html",
			"templates/bill-items-select.html",
			"templates/bill-items.html",
			"templates/bill-items-list.html",
			"templates/issuers.html",
			"templates/issuers-list.html",
			"templates/issuers-select.html",
			"templates/receivers.html",
			"templates/receivers-list.html",
			"templates/receivers-select.html",
		)),
	}
	e.Renderer = t

	// Initialize handlers
	billHandler := handlers.NewBillHandler(billRepo, receiverRepo, issuerRepo, billItemRepo, billItemAssignmentRepo, t.templates)
	receiverHandler := handlers.NewReceiverHandler(receiverRepo, t.templates)
	issuerHandler := handlers.NewIssuerHandler(issuerRepo, t.templates)
	billItemHandler := handlers.NewBillItemHandler(billItemRepo, t.templates)

	// Bill routes
	e.GET("/", billHandler.RenderBills)
	e.POST("/bills", billHandler.CreateBill)
	e.GET("/bills", billHandler.RenderBills)
	e.POST("/bills/:id/toggle", billHandler.TogglePaid)
	e.DELETE("/bills/:id", billHandler.DeleteBill)

	// Receiver routes
	e.GET("/receivers", receiverHandler.RenderReceivers)
	e.POST("/receivers", receiverHandler.CreateReceiver)
	e.GET("/receivers/list", receiverHandler.GetReceiversList)
	e.GET("/receivers/select", receiverHandler.GetReceiversSelect)
	e.DELETE("/receivers/:id", receiverHandler.DeleteReceiver)

	// Issuer routes
	e.GET("/issuers", issuerHandler.RenderIssuers)
	e.POST("/issuers", issuerHandler.CreateIssuer)
	e.GET("/issuers/list", issuerHandler.GetIssuersList)
	e.GET("/issuers/select", issuerHandler.GetIssuersSelect)
	e.DELETE("/issuers/:id", issuerHandler.DeleteIssuer)

	// Bill Item routes
	e.GET("/bill-items", billItemHandler.RenderBillItems)
	e.POST("/bill-items", billItemHandler.CreateBillItem)
	e.GET("/bill-items/list", billItemHandler.GetBillItemsList)
	e.GET("/bill-items/select", billItemHandler.GetBillItemsSelect)
	e.DELETE("/bill-items/:id", billItemHandler.DeleteBillItem)

	// Start server
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	e.Logger.Fatal(e.Start(":" + port))
}
