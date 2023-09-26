package auth

import (
	"mcash-finance-console-core/internal/middleware"

	"github.com/labstack/echo/v4"
)

func (h *handler) Route(g *echo.Group) {
	g.POST("/login", h.Login)
	// g.POST("/register", h.Register)
	g.GET("/checkauth", h.CheckAuth)
	g.PATCH("/change-password", h.ChangePassword, middleware.Authentication)
	g.GET("/notification-token", h.GetNotificationToken, middleware.Authentication)
}
