package handlers

import (
	"bills/internal/models"
	"bills/internal/repository"
	"html/template"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
)

// BillItemHandler handles HTTP requests for bill items
type BillItemHandler struct {
	repo repository.BillItemRepository
	tmpl *template.Template
}

// NewBillItemHandler creates a new BillItemHandler instance
func NewBillItemHandler(repo repository.BillItemRepository, tmpl *template.Template) *BillItemHandler {
	return &BillItemHandler{
		repo: repo,
		tmpl: tmpl,
	}
}

// RenderBillItems renders the bill items list template
func (h *BillItemHandler) RenderBillItems(c echo.Context) error {
	items, err := h.repo.GetAll()
	if err != nil {
		return err
	}

	return c.Render(http.StatusOK, "bill-items.html", map[string]interface{}{
		"Items": items,
	})
}

// GetBillItemsList returns the bill items list partial for HTMX updates
func (h *BillItemHandler) GetBillItemsList(c echo.Context) error {
	items, err := h.repo.GetAll()
	if err != nil {
		return err
	}

	return c.Render(http.StatusOK, "bill-items-list.html", map[string]interface{}{
		"Items": items,
	})
}

// GetBillItemsSelect returns a select dropdown with bill items for HTMX updates
func (h *BillItemHandler) GetBillItemsSelect(c echo.Context) error {
	items, err := h.repo.GetAll()
	if err != nil {
		return err
	}

	return c.Render(http.StatusOK, "bill-items-select.html", map[string]interface{}{
		"Items": items,
	})
}

// CreateBillItem handles the creation of a new bill item
func (h *BillItemHandler) CreateBillItem(c echo.Context) error {
	defaultPrice, err := strconv.ParseFloat(c.FormValue("default_price"), 64)
	if err != nil {
		return err
	}

	item := models.NewBillItem(
		c.FormValue("description"),
		defaultPrice,
	)

	if err := h.repo.Create(item); err != nil {
		return err
	}

	// If it's an HTMX request, return the updated list
	if c.Request().Header.Get("HX-Request") == "true" {
		return h.GetBillItemsList(c)
	}

	return c.Redirect(http.StatusSeeOther, "/bill-items")
}

// UpdateBillItem handles updating a bill item
func (h *BillItemHandler) UpdateBillItem(c echo.Context) error {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		return err
	}

	item, err := h.repo.GetByID(id)
	if err != nil {
		return err
	}

	defaultPrice, err := strconv.ParseFloat(c.FormValue("default_price"), 64)
	if err != nil {
		return err
	}

	item.Description = c.FormValue("description")
	item.DefaultPrice = defaultPrice

	if err := h.repo.Update(item); err != nil {
		return err
	}

	return h.GetBillItemsList(c)
}

// DeleteBillItem handles the deletion of a bill item
func (h *BillItemHandler) DeleteBillItem(c echo.Context) error {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		return err
	}

	if err := h.repo.Delete(id); err != nil {
		return err
	}

	return h.GetBillItemsList(c)
}
