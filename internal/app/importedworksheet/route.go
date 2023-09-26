package importedworksheet

import (
	"mcash-finance-console-core/internal/middleware"

	"github.com/labstack/echo/v4"
)

func (h *handler) Route(g *echo.Group) {
	g.GET("", h.Get, middleware.Authentication)
	g.GET("/:id", h.GetByID, middleware.Authentication)
	g.DELETE("/:id", h.Delete, middleware.Authentication)
}
