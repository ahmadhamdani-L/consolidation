package consolidation

import (
	"mcash-finance-console-core/internal/middleware"

	"github.com/labstack/echo/v4"
)

func (h *handler) Route(g *echo.Group) {
	g.GET("/list_company/:company_id", h.GetListAvailable, middleware.Authentication)
	g.GET("", h.Get, middleware.Authentication)
	g.POST("/combaine", h.Combaine, middleware.Authentication)
	g.POST("/edit-combaine", h.EditCombain, middleware.Authentication)
	g.POST("/duplicate", h.Duplicate, middleware.Authentication)
	g.POST("/consolidation", h.Consolidation, middleware.Authentication)
	g.GET("/list_company_duplicate/:consolidation_id", h.GetListDuplicateAvailable, middleware.Authentication)
	g.GET("/get-version", h.GetVersion, middleware.Authentication)
	g.GET("/get-company-list-combaine", h.FindListCompanyCreateNewCombine, middleware.Authentication)
	g.DELETE("/:id", h.Delete, middleware.Authentication)
	g.GET("/:id", h.GetByID, middleware.Authentication)
	g.GET("/control/:consolidation_id", h.GetControl, middleware.Authentication)
}
