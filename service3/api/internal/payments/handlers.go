package payments

import (
	"net/http"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

type Handler struct {
	service Service
}

func NewPaymentHandler(service Service) Handler {
	return Handler{service}
}

func (h *Handler) Create(c echo.Context) error {
	payment := new(Payment)
	if err := c.Bind(payment); err != nil {
		return err
	}

	payment.Id = uuid.New()
	if payment.PaymentType == "" {
		payment.PaymentType = "regular"
	}
	if err := h.service.Create(c.Request().Context(), *payment); err != nil {
		return err
	}

	return c.JSON(http.StatusCreated, payment)
}

func (h *Handler) Read(c echo.Context) error {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return err
	}

	payment, err := h.service.Read(c.Request().Context(), id)
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, payment)
}

func (h *Handler) GetByLoanId(c echo.Context) error {
	loanId, err := uuid.Parse(c.Param("loanId"))
	if err != nil {
		return err
	}

	payments, err := h.service.GetByLoanId(c.Request().Context(), loanId)
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, payments)
}

func (h *Handler) GetByCustomerId(c echo.Context) error {
	customerId, err := uuid.Parse(c.Param("customerId"))
	if err != nil {
		return err
	}

	payments, err := h.service.GetByCustomerId(c.Request().Context(), customerId)
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, payments)
}
