package bills

import (
	"net/http"

	"bills/concepts"

	"github.com/gin-gonic/gin"
)

// Controller manages bill-related HTTP requests and responses
type Controller struct {
	// We'll add service dependencies here later
}

// Ensure Controller implements the concepts.Controller interface
var _ concepts.Controller = (*Controller)(nil)

// NewController creates a new instance of Controller
func NewController() concepts.Controller {
	return &Controller{}
}

// Templates returns the template mappings for this controller
func (c *Controller) Templates() map[string][]string {
	return map[string][]string{
		"bills_index": {"layouts/application.html", "concepts/bills/templates/index.html"},
		"bills_new":   {"layouts/application.html", "concepts/bills/templates/new.html"},
	}
}

// RegisterRoutes registers all bill-related routes
func (c *Controller) RegisterRoutes(engine *gin.Engine) {
	billsGroup := engine.Group("/bills")
	{
		billsGroup.GET("/", c.Index)
		billsGroup.GET("/new", c.New)
	}
}

// Index renders the main bills page
func (c *Controller) Index(ctx *gin.Context) {
	ctx.HTML(http.StatusOK, "bills_index", gin.H{
		"title": "Bills Manager",
	})
}

// New renders the new bill form
func (c *Controller) New(ctx *gin.Context) {
	// Check if this is an HTMX request
	if ctx.GetHeader("HX-Request") == "true" {
		ctx.HTML(http.StatusOK, "bills_new", gin.H{
			"title": "New Bill",
		})
	} else {
		ctx.HTML(http.StatusOK, "bills_new", gin.H{
			"title": "New Bill",
		})
	}
}
