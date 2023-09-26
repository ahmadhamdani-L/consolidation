package validation

import (
	"mcash-finance-console-core/internal/middleware"

	"github.com/labstack/echo/v4"
)

func (h *handler) Route(g *echo.Group) {
	g.GET("", h.Get, middleware.Authentication)
	g.GET("/list_company", h.GetListAvailable, middleware.Authentication)
	g.GET("/:id", h.GetByID, middleware.Authentication)
	g.POST("", h.Validate, middleware.Authentication)
	g.POST("/:validation_id", h.ValidateModul, middleware.Authentication)
	// g.PATCH("/:id", h.Re, middleware.Authentication)
}
