package loans

import "github.com/labstack/echo/v4"

func Routes(e *echo.Echo, handler Handler) {
	e.POST("/loans", handler.Create)
	e.GET("/loans/:id", handler.Read)
	e.PUT("/loans/:id", handler.Update)
	e.DELETE("/loans/:id", handler.Delete)
	e.GET("/customers/:customerId/loans", handler.GetByCustomerId)
	e.GET("/mortgages/:mortgageId/loan", handler.GetByMortgageId)
}
