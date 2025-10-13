package payments

import "github.com/labstack/echo/v4"

func Routes(e *echo.Echo, handler Handler) {
	e.POST("/payments", handler.Create)
	e.GET("/payments/:id", handler.Read)
	e.GET("/loans/:loanId/payments", handler.GetByLoanId)
	e.GET("/customers/:customerId/payments", handler.GetByCustomerId)
}
