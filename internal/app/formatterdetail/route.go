package formatterdetail

import (
	"mcash-finance-console-core/internal/middleware"

	"github.com/labstack/echo/v4"
)

func (h *handler) Route(g *echo.Group) {
	g.GET("/view-all", h.ViewAll, middleware.Authentication)
	g.POST("/drag-and-drop/:parent_id", h.DragAndDrop, middleware.Authentication)
	g.POST("/createsubcoa", h.Create, middleware.Authentication)
	g.GET("/find-coa-formatterxmcoa", h.GetCoa, middleware.Authentication)
	g.POST("/createcoa", h.CreateCoaxFmt, middleware.Authentication)
	g.DELETE("/:id", h.Delete, middleware.Authentication)
	g.GET("/export", h.ExportFmtDev, middleware.Authentication)
	
}
