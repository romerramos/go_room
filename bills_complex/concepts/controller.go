package concepts

import "github.com/gin-gonic/gin"

// Controller defines the interface that all feature controllers must implement
type Controller interface {
	// Templates returns a map of template names to their file paths
	// The first file should be the layout, followed by the content template
	Templates() map[string][]string

	// RegisterRoutes sets up all routes for this controller
	RegisterRoutes(engine *gin.Engine)
}
