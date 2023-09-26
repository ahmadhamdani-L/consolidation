package export

import (
	"mcash-finance-console-core/internal/middleware"

	"github.com/labstack/echo/v4"
)

func (h *handler) Route(g *echo.Group) {
	g.GET("/request", h.Request, middleware.Authentication)
	g.GET("/export_modul/:notification_id", h.GetRequest, middleware.Authentication)
	g.GET("/request_consol", h.RequestExportConsol, middleware.Authentication)
	g.GET("/export_consol/:notification_id", h.GetRequestConsol, middleware.Authentication)
	g.GET("", h.ExportAll, middleware.Authentication)
}
