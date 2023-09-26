package imports

import (
	"mcash-finance-console-core/internal/middleware"

	"github.com/labstack/echo/v4"
)

func (h *handler) Route(g *echo.Group) {
	g.POST("/importasync", h.ImportAsync, middleware.Authentication)
	g.POST("/import/:id", h.Import, middleware.Authentication)
	g.GET("/:id", h.Download)
	g.GET("/download-template", h.DownloadTemplate, middleware.Authentication)
	g.GET("/download-all-worksheet/:id", h.DownloadAll)
	g.GET("/template", h.UploadTemplate)
	g.POST("/jurnal/:jurnal", h.UploadJurnal)
	g.POST("/request-jurnal", h.ImportJurnalAsync)
	
}
