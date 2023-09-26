package company

import (
	"mcash-finance-console-core/internal/middleware"

	"github.com/labstack/echo/v4"
)

func (h *handler) Route(g *echo.Group) {
	g.GET("", h.Get, middleware.Authentication)
	g.GET("/:id", h.GetByID, middleware.Authentication)
	g.GET("/get-filter", h.GetFilterList, middleware.Authentication)
	g.POST("", h.Create, middleware.Authentication)
	g.PATCH("/:id", h.Update, middleware.Authentication)
	g.DELETE("/:id", h.Delete, middleware.Authentication)
	g.GET("/treeview", h.GetTreeview, middleware.Authentication)
}
