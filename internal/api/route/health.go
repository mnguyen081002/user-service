package route

import (
	"erp/internal/api/controllers"
	"erp/internal/lib"
)

type HealthRoutes struct {
	handler *lib.Handler
}

func NewHealthRoutes(handler *lib.Handler, controller *controller.HealthController) *HealthRoutes {
	g := handler.Group("/health")

	g.GET("/", controller.Health)

	return &HealthRoutes{
		handler: handler,
	}
}
