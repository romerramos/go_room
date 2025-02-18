package handlers

import (
	"bills/internal/models"
	"bills/internal/repository"
	"html/template"
	"net/http"
	"strconv"
	"time"

	"github.com/labstack/echo/v4"
)

// BillHandler handles HTTP requests for bills
type BillHandler struct {
	repo               repository.BillRepository
	receiverRepo       repository.ReceiverRepository
	issuerRepo         repository.IssuerRepository
	billItemRepo       repository.BillItemRepository
	billItemAssignRepo repository.BillItemAssignmentRepository
	tmpl               *template.Template
}

// NewBillHandler creates a new BillHandler instance
func NewBillHandler(
	repo repository.BillRepository,
	receiverRepo repository.ReceiverRepository,
	issuerRepo repository.IssuerRepository,
	billItemRepo repository.BillItemRepository,
	billItemAssignRepo repository.BillItemAssignmentRepository,
	tmpl *template.Template,
) *BillHandler {
	return &BillHandler{
		repo:               repo,
		receiverRepo:       receiverRepo,
		issuerRepo:         issuerRepo,
		billItemRepo:       billItemRepo,
		billItemAssignRepo: billItemAssignRepo,
		tmpl:               tmpl,
	}
}

// RenderBills renders the bills list template
func (h *BillHandler) RenderBills(c echo.Context) error {
	bills, err := h.repo.GetAll()
	if err != nil {
		return err
	}

	receivers, err := h.receiverRepo.GetAll()
	if err != nil {
		return err
	}

	issuers, err := h.issuerRepo.GetAll()
	if err != nil {
		return err
	}

	billItems, err := h.billItemRepo.GetAll()
	if err != nil {
		return err
	}

	return c.Render(http.StatusOK, "bills.html", map[string]interface{}{
		"Bills":               bills,
		"Receivers":           receivers,
		"Issuers":             issuers,
		"Items":               billItems,
		"Today":               time.Now().Format("2006-01-02"),
		"SupportedCurrencies": models.SupportedCurrencies(),
		"DefaultCurrency":     models.DefaultCurrency(),
	})
}

// GetBillsList returns the bills list partial for HTMX updates
func (h *BillHandler) GetBillsList(c echo.Context) error {
	bills, err := h.repo.GetAll()
	if err != nil {
		return err
	}

	return c.Render(http.StatusOK, "bills-list.html", map[string]interface{}{
		"Bills": bills,
	})
}

// CreateBill handles the creation of a new bill
func (h *BillHandler) CreateBill(c echo.Context) error {
	// Parse form data
	dueDate, err := time.Parse("2006-01-02", c.FormValue("due_date"))
	if err != nil {
		return err
	}

	issuerID, err := strconv.ParseInt(c.FormValue("issuer_id"), 10, 64)
	if err != nil {
		return err
	}

	receiverID, err := strconv.ParseInt(c.FormValue("receiver_id"), 10, 64)
	if err != nil {
		return err
	}

	// Create the bill
	bill := models.NewBill(dueDate, issuerID, receiverID)

	// Parse bill item assignments
	itemIDs := c.Request().Form["item_ids[]"]
	quantities := c.Request().Form["quantities[]"]
	prices := c.Request().Form["prices[]"]
	currencies := c.Request().Form["currencies[]"]
	exchangeRates := c.Request().Form["exchange_rates[]"]

	// Track unique currencies
	uniqueCurrencies := make(map[string]bool)

	// Create bill item assignments
	for i := range itemIDs {
		itemID, err := strconv.ParseInt(itemIDs[i], 10, 64)
		if err != nil {
			continue
		}

		quantity, err := strconv.Atoi(quantities[i])
		if err != nil {
			continue
		}

		price, err := strconv.ParseFloat(prices[i], 64)
		if err != nil {
			continue
		}

		currency := currencies[i]
		if currency == "" || !models.IsSupportedCurrency(currency) {
			currency = models.DefaultCurrency()
		}
		uniqueCurrencies[currency] = true

		exchangeRate, err := strconv.ParseFloat(exchangeRates[i], 64)
		if err != nil || exchangeRate <= 0 {
			if models.IsDefaultCurrency(currency) {
				exchangeRate = 1.0
			} else {
				// TODO: Get exchange rate from service
				exchangeRate = 1.0
			}
		}

		assignment := models.NewBillItemAssignment(0, itemID, quantity, price, currency, exchangeRate)
		bill.Items = append(bill.Items, assignment)
	}

	// Set bill currency based on items
	if len(uniqueCurrencies) == 1 {
		// If all items have the same currency, use that
		for currency := range uniqueCurrencies {
			bill.Currency = currency
		}
	} else {
		// If items have different currencies, use EUR
		bill.Currency = models.DefaultCurrency()
	}

	// Calculate totals
	bill.CalculateTotals()

	// Save the bill
	if err := h.repo.Create(bill); err != nil {
		return err
	}

	// Get issuer and receiver names
	issuer, err := h.issuerRepo.GetByID(issuerID)
	if err != nil {
		return err
	}
	receiver, err := h.receiverRepo.GetByID(receiverID)
	if err != nil {
		return err
	}
	bill.IssuerName = issuer.Name
	bill.ReceiverName = receiver.Name

	// If it's an HTMX request, return the updated list
	if c.Request().Header.Get("HX-Request") == "true" {
		return h.GetBillsList(c)
	}

	return c.Redirect(http.StatusSeeOther, "/")
}

// TogglePaid handles toggling the paid status of a bill
func (h *BillHandler) TogglePaid(c echo.Context) error {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		return err
	}

	bill, err := h.repo.GetByID(id)
	if err != nil {
		return err
	}

	bill.Paid = !bill.Paid
	if err := h.repo.Update(bill); err != nil {
		return err
	}

	return h.GetBillsList(c)
}

// DeleteBill handles the deletion of a bill
func (h *BillHandler) DeleteBill(c echo.Context) error {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		return err
	}

	// Delete bill item assignments first
	if err := h.billItemAssignRepo.DeleteByBillID(id); err != nil {
		return err
	}

	// Then delete the bill
	if err := h.repo.Delete(id); err != nil {
		return err
	}

	// If it's an HTMX request, return the updated list
	if c.Request().Header.Get("HX-Request") == "true" {
		return h.GetBillsList(c)
	}

	// For regular requests, redirect to the root URL
	return c.Redirect(http.StatusSeeOther, "/")
}
