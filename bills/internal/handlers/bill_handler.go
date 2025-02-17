package handlers

import (
	"bills/internal/models"
	"bills/internal/repository"
	"fmt"
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
		"Bills":     bills,
		"Receivers": receivers,
		"Issuers":   issuers,
		"Items":     billItems,
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
	// Log form data
	fmt.Printf("\nForm Data:\n")
	fmt.Printf("due_date: %s\n", c.FormValue("due_date"))
	fmt.Printf("issuer_id: %s\n", c.FormValue("issuer_id"))
	fmt.Printf("receiver_id: %s\n", c.FormValue("receiver_id"))
	fmt.Printf("item_ids: %v\n", c.Request().Form["item_ids[]"])
	fmt.Printf("quantities: %v\n", c.Request().Form["quantities[]"])
	fmt.Printf("unit_prices: %v\n", c.Request().Form["unit_prices[]"])

	dueDate, err := time.Parse("2006-01-02", c.FormValue("due_date"))
	if err != nil {
		fmt.Printf("Error parsing due date: %v\n", err)
		return err
	}

	issuerID, err := strconv.ParseInt(c.FormValue("issuer_id"), 10, 64)
	if err != nil {
		fmt.Printf("Error parsing issuer_id: %v\n", err)
		return err
	}

	receiverID, err := strconv.ParseInt(c.FormValue("receiver_id"), 10, 64)
	if err != nil {
		fmt.Printf("Error parsing receiver_id: %v\n", err)
		return err
	}

	// Create the bill
	bill := models.NewBill(dueDate, issuerID, receiverID)

	// Parse bill item assignments
	itemIDs := c.Request().Form["item_ids[]"]
	quantities := c.Request().Form["quantities[]"]
	unitPrices := c.Request().Form["unit_prices[]"]

	fmt.Printf("\nProcessing %d items\n", len(itemIDs))

	// Create bill item assignments
	for i := range itemIDs {
		itemID, err := strconv.ParseInt(itemIDs[i], 10, 64)
		if err != nil {
			fmt.Printf("Error parsing item_id[%d]: %v\n", i, err)
			continue
		}

		quantity, err := strconv.Atoi(quantities[i])
		if err != nil {
			fmt.Printf("Error parsing quantity[%d]: %v\n", i, err)
			continue
		}

		unitPrice, err := strconv.ParseFloat(unitPrices[i], 64)
		if err != nil {
			fmt.Printf("Error parsing unit_price[%d]: %v\n", i, err)
			continue
		}

		assignment := models.NewBillItemAssignment(0, itemID, quantity, unitPrice)
		bill.Items = append(bill.Items, assignment)
		fmt.Printf("Added item: ID=%d, Quantity=%d, UnitPrice=%.2f\n", itemID, quantity, unitPrice)
	}

	// Calculate total amount from items
	bill.CalculateAmount()
	fmt.Printf("\nCalculated total amount: %.2f\n", bill.Amount)

	// Save the bill
	if err := h.repo.Create(bill); err != nil {
		fmt.Printf("Error creating bill: %v\n", err)
		return err
	}

	fmt.Printf("Successfully created bill with ID: %d\n", bill.ID)

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

	return h.GetBillsList(c)
}
