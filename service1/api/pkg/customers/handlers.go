package customers

import (
	"net/http"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

type Handler struct {
	service Service
}

func NewCustomersHandler(service Service) Handler {
	return Handler{service}
}

func (h *Handler) Create(c echo.Context) error {
	customer := new(Customer)
	if err := c.Bind(customer); err != nil {
		return err
	}

	customer.Id = uuid.New()
	if err := h.service.Create(c.Request().Context(), *customer); err != nil {
		return err
	}

	return c.JSON(http.StatusCreated, customer)
}

func (h *Handler) Read(c echo.Context) error {
	id := c.Param("id")
	customer, err := h.service.Read(c.Request().Context(), uuid.MustParse(id))
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, customer)
}

func (h *Handler) Update(c echo.Context) error {
	id := c.Param("id")
	customer := new(Customer)
	if err := c.Bind(customer); err != nil {
		return err
	}
	customer.Id = uuid.MustParse(id)
	if err := h.service.Update(c.Request().Context(), *customer); err != nil {
		return err
	}
	return c.JSON(http.StatusOK, customer)
}

func (h *Handler) Delete(c echo.Context) error {
	id := c.Param("id")
	if err := h.service.Delete(c.Request().Context(), uuid.MustParse(id)); err != nil {
		return err
	}
	return c.NoContent(http.StatusNoContent)
}
