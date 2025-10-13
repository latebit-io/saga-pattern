package mortgages

import (
	"net/http"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

type Handler struct {
	service Service
}

func NewMortgageHandler(service Service) Handler {
	return Handler{service}
}

func (h *Handler) Create(c echo.Context) error {
	application := new(MortgageApplication)
	if err := c.Bind(application); err != nil {
		return err
	}

	application.Id = uuid.New()
	if application.Status == "" {
		application.Status = "pending"
	}
	if err := h.service.Create(c.Request().Context(), *application); err != nil {
		return err
	}

	return c.JSON(http.StatusCreated, application)
}

func (h *Handler) Read(c echo.Context) error {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return err
	}

	application, err := h.service.Read(c.Request().Context(), id)
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, application)
}

func (h *Handler) Update(c echo.Context) error {
	id := c.Param("id")
	application := new(MortgageApplication)
	if err := c.Bind(application); err != nil {
		return err
	}
	var err error
	application.Id, err = uuid.Parse(id)
	if err != nil {
		return err
	}
	if err := h.service.Update(c.Request().Context(), *application); err != nil {
		return err
	}
	return c.JSON(http.StatusOK, application)
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

	applications, err := h.service.GetByCustomerId(c.Request().Context(), customerId)
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, applications)
}
