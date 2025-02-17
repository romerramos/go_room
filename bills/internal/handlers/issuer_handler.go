package handlers

import (
	"bills/internal/models"
	"bills/internal/repository"
	"html/template"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
)

// IssuerHandler handles HTTP requests for issuers
type IssuerHandler struct {
	repo repository.IssuerRepository
	tmpl *template.Template
}

// NewIssuerHandler creates a new IssuerHandler instance
func NewIssuerHandler(repo repository.IssuerRepository, tmpl *template.Template) *IssuerHandler {
	return &IssuerHandler{
		repo: repo,
		tmpl: tmpl,
	}
}

// RenderIssuers renders the issuers list template
func (h *IssuerHandler) RenderIssuers(c echo.Context) error {
	issuers, err := h.repo.GetAll()
	if err != nil {
		return err
	}

	return c.Render(http.StatusOK, "issuers.html", map[string]interface{}{
		"Issuers": issuers,
	})
}

// GetIssuersList returns the issuers list partial for HTMX updates
func (h *IssuerHandler) GetIssuersList(c echo.Context) error {
	issuers, err := h.repo.GetAll()
	if err != nil {
		return err
	}

	return c.Render(http.StatusOK, "issuers-list.html", map[string]interface{}{
		"Issuers": issuers,
	})
}

// GetIssuersSelect returns a select dropdown with issuers for HTMX updates
func (h *IssuerHandler) GetIssuersSelect(c echo.Context) error {
	issuers, err := h.repo.GetAll()
	if err != nil {
		return err
	}

	return c.Render(http.StatusOK, "issuers-select.html", map[string]interface{}{
		"Issuers": issuers,
	})
}

// CreateIssuer handles the creation of a new issuer
func (h *IssuerHandler) CreateIssuer(c echo.Context) error {
	issuer := models.NewIssuer(
		c.FormValue("name"),
		c.FormValue("vat_number"),
		c.FormValue("street"),
		c.FormValue("city"),
		c.FormValue("state"),
		c.FormValue("zip_code"),
		c.FormValue("country"),
	)

	if err := h.repo.Create(issuer); err != nil {
		return err
	}

	// If it's an HTMX request, return the updated list
	if c.Request().Header.Get("HX-Request") == "true" {
		return h.GetIssuersList(c)
	}

	return c.Redirect(http.StatusSeeOther, "/issuers")
}

// UpdateIssuer handles updating an issuer
func (h *IssuerHandler) UpdateIssuer(c echo.Context) error {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		return err
	}

	issuer, err := h.repo.GetByID(id)
	if err != nil {
		return err
	}

	issuer.Name = c.FormValue("name")
	issuer.VATNumber = c.FormValue("vat_number")
	issuer.Street = c.FormValue("street")
	issuer.City = c.FormValue("city")
	issuer.State = c.FormValue("state")
	issuer.ZipCode = c.FormValue("zip_code")
	issuer.Country = c.FormValue("country")

	if err := h.repo.Update(issuer); err != nil {
		return err
	}

	return h.GetIssuersList(c)
}

// DeleteIssuer handles the deletion of an issuer
func (h *IssuerHandler) DeleteIssuer(c echo.Context) error {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		return err
	}

	if err := h.repo.Delete(id); err != nil {
		return err
	}

	return h.GetIssuersList(c)
}
