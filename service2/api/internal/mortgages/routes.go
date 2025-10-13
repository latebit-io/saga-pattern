package mortgages

import "github.com/labstack/echo/v4"

func Routes(e *echo.Echo, handler Handler) {
	e.POST("/applications", handler.Create)
	e.GET("/applications/:id", handler.Read)
	e.PUT("/applications/:id", handler.Update)
	e.DELETE("/applications/:id", handler.Delete)
	e.GET("/customers/:customerId/applications", handler.GetByCustomerId)
}
