package route

import (
	controller "erp/internal/api/controllers"
	"erp/internal/lib"
)

type UserRoutes struct {
	handler *lib.Handler
}

func NewUserRoutes(handler *lib.Handler, controller *controller.UserController) *UserRoutes {
	r := handler.Group("/users")

	r.GET("/", controller.ListUsers)

	return &UserRoutes{
		handler: handler,
	}
}
