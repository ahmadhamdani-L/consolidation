package consolidationdetail

import (
	"mcash-finance-console-core/internal/middleware"

	"github.com/labstack/echo/v4"
)

func (h *handler) Route(g *echo.Group) {
	g.GET("", h.View, middleware.Authentication)
	g.GET("/view-all", h.ViewAll, middleware.Authentication)
}
