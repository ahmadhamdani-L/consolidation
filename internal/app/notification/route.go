package notification

import (
	"mcash-finance-console-core/internal/middleware"

	"github.com/labstack/echo/v4"
)

func (h *handler) Route(g *echo.Group) {
	g.GET("", h.Get, middleware.Authentication)
	g.GET("/:id", h.GetByID, middleware.Authentication)
	g.PATCH("/read", h.MarkAsRead, middleware.Authentication)
	g.GET("/tes", h.Test, middleware.Authentication)
}
