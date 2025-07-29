package customers

import "github.com/labstack/echo/v4"

func Routes(e *echo.Echo, handler Handler) {
	e.POST("/customers", handler.Create)
	e.GET("/customers/:id", handler.Read)
	e.PUT("/customers/:id", handler.Update)
	e.DELETE("/customers/:id", handler.Delete)
}
