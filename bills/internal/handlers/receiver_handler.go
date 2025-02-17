package handlers

import (
	"bills/internal/models"
	"bills/internal/repository"
	"html/template"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
)

// ReceiverHandler handles HTTP requests for receivers
type ReceiverHandler struct {
	repo repository.ReceiverRepository
	tmpl *template.Template
}

// NewReceiverHandler creates a new ReceiverHandler instance
func NewReceiverHandler(repo repository.ReceiverRepository, tmpl *template.Template) *ReceiverHandler {
	return &ReceiverHandler{
		repo: repo,
		tmpl: tmpl,
	}
}

// RenderReceivers renders the receivers list template
func (h *ReceiverHandler) RenderReceivers(c echo.Context) error {
	receivers, err := h.repo.GetAll()
	if err != nil {
		return err
	}

	return c.Render(http.StatusOK, "receivers.html", map[string]interface{}{
		"Receivers": receivers,
	})
}

// GetReceiversList returns the receivers list partial for HTMX updates
func (h *ReceiverHandler) GetReceiversList(c echo.Context) error {
	receivers, err := h.repo.GetAll()
	if err != nil {
		return err
	}

	return c.Render(http.StatusOK, "receivers-list.html", map[string]interface{}{
		"Receivers": receivers,
	})
}

// GetReceiversSelect returns a select dropdown with receivers for HTMX updates
func (h *ReceiverHandler) GetReceiversSelect(c echo.Context) error {
	receivers, err := h.repo.GetAll()
	if err != nil {
		return err
	}

	return c.Render(http.StatusOK, "receivers-select.html", map[string]interface{}{
		"Receivers": receivers,
	})
}

// CreateReceiver handles the creation of a new receiver
func (h *ReceiverHandler) CreateReceiver(c echo.Context) error {
	receiver := models.NewReceiver(
		c.FormValue("name"),
		c.FormValue("vat_number"),
		c.FormValue("street"),
		c.FormValue("city"),
		c.FormValue("state"),
		c.FormValue("zip_code"),
		c.FormValue("country"),
	)

	if err := h.repo.Create(receiver); err != nil {
		return err
	}

	// If it's an HTMX request, return the updated list
	if c.Request().Header.Get("HX-Request") == "true" {
		return h.GetReceiversList(c)
	}

	return c.Redirect(http.StatusSeeOther, "/receivers")
}

// UpdateReceiver handles updating a receiver
func (h *ReceiverHandler) UpdateReceiver(c echo.Context) error {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		return err
	}

	receiver, err := h.repo.GetByID(id)
	if err != nil {
		return err
	}

	receiver.Name = c.FormValue("name")
	receiver.VATNumber = c.FormValue("vat_number")
	receiver.Street = c.FormValue("street")
	receiver.City = c.FormValue("city")
	receiver.State = c.FormValue("state")
	receiver.ZipCode = c.FormValue("zip_code")
	receiver.Country = c.FormValue("country")

	if err := h.repo.Update(receiver); err != nil {
		return err
	}

	return h.GetReceiversList(c)
}

// DeleteReceiver handles the deletion of a receiver
func (h *ReceiverHandler) DeleteReceiver(c echo.Context) error {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		return err
	}

	if err := h.repo.Delete(id); err != nil {
		return err
	}

	return h.GetReceiversList(c)
}
