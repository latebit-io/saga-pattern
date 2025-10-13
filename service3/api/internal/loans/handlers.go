package loans

import (
	"net/http"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

type Handler struct {
	service Service
}

func NewLoanHandler(service Service) Handler {
	return Handler{service}
}

func (h *Handler) Create(c echo.Context) error {
	loan := new(Loan)
	if err := c.Bind(loan); err != nil {
		return err
	}

	loan.Id = uuid.New()
	if loan.Status == "" {
		loan.Status = "active"
	}
	if err := h.service.Create(c.Request().Context(), *loan); err != nil {
		return err
	}

	return c.JSON(http.StatusCreated, loan)
}

func (h *Handler) Read(c echo.Context) error {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return err
	}

	loan, err := h.service.Read(c.Request().Context(), id)
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, loan)
}

func (h *Handler) Update(c echo.Context) error {
	id := c.Param("id")
	loan := new(Loan)
	if err := c.Bind(loan); err != nil {
		return err
	}
	var err error
	loan.Id, err = uuid.Parse(id)
	if err != nil {
		return err
	}
	if err := h.service.Update(c.Request().Context(), *loan); err != nil {
		return err
	}
	return c.JSON(http.StatusOK, loan)
}

func (h *Handler) Delete(c echo.Context) error {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return err
	}
	if err := h.service.Delete(c.Request().Context(), id); err != nil {
		return err
	}
	return c.NoContent(http.StatusNoContent)
}

func (h *Handler) GetByCustomerId(c echo.Context) error {
	customerId, err := uuid.Parse(c.Param("customerId"))
	if err != nil {
		return err
	}

	loans, err := h.service.GetByCustomerId(c.Request().Context(), customerId)
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, loans)
}

func (h *Handler) GetByMortgageId(c echo.Context) error {
	mortgageId, err := uuid.Parse(c.Param("mortgageId"))
	if err != nil {
		return err
	}

	loan, err := h.service.GetByMortgageId(c.Request().Context(), mortgageId)
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, loan)
}
