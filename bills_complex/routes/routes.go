package routes

import (
	"bills/concepts"
	"bills/concepts/bills"

	"github.com/gin-contrib/multitemplate"
	"github.com/gin-gonic/gin"
)

// SetupRoutes configures all application routes
func SetupRoutes(engine *gin.Engine) {
	// Initialize all controllers
	controllers := []concepts.Controller{
		bills.NewController(),
		// Add new controllers here:
		// users.NewController(),
		// categories.NewController(),
	}

	// Initialize renderer with all templates
	renderer := multitemplate.NewRenderer()

	// Register templates from all controllers
	for _, ctrl := range controllers {
		for name, files := range ctrl.Templates() {
			renderer.AddFromFiles(name, files...)
		}
	}

	// Set up the renderer
	engine.HTMLRender = renderer

	// Register routes from all controllers
	for _, ctrl := range controllers {
		ctrl.RegisterRoutes(engine)
	}

	// Set up static files
	engine.Static("/static", "./static")
}
