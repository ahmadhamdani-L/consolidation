package user

import (
	"mcash-finance-console-core/internal/middleware"

	"github.com/labstack/echo/v4"
)

func (h *handler) Route(g *echo.Group) {
	g.GET("", h.Get, middleware.Authentication)
	g.GET("/:id", h.GetByID, middleware.Authentication)
	g.POST("", h.Create, middleware.Authentication)
	g.PATCH("/:id", h.Update, middleware.Authentication)
	g.DELETE("/:id", h.Delete, middleware.Authentication)
	g.POST("/reset-password", h.ForgotPassword)
	g.PATCH("/reset-password/:resetToken", h.ResetPassword, middleware.AuthenticationResetPassword)
	g.PATCH("/change-status/:user_id", h.UserActive, middleware.Authentication)
}
