package auth

import (
	"mcash-finance-console-core/internal/abstraction"
	"mcash-finance-console-core/internal/dto"
	"mcash-finance-console-core/internal/factory"
	res "mcash-finance-console-core/pkg/util/response"

	"github.com/labstack/echo/v4"
)

type handler struct {
	service *service
}

var err error

func NewHandler(f *factory.Factory) *handler {
	service := NewService(f)
	return &handler{service}
}

// Login
// @Summary Login user
// @Description Login user
// @Tags auth
// @Accept json
// @Produce json
// @Param request body dto.AuthLoginRequest true "request body"
// @Success 200 {object} dto.AuthLoginResponseDoc
// @Failure 400 {object} res.errorResponse
// @Failure 404 {object} res.errorResponse
// @Failure 500 {object} res.errorResponse
// @Router /auth/login [post]
func (h *handler) Login(c echo.Context) error {
	cc := c.(*abstraction.Context)

	payload := new(dto.AuthLoginRequest)
	if err = c.Bind(payload); err != nil {
		return res.ErrorBuilder(&res.ErrorConstant.BadRequest, err).Send(c)
	}
	if err = c.Validate(payload); err != nil {
		return res.ErrorBuilder(&res.ErrorConstant.Validation, err).Send(c)
	}

	data, err := h.service.Login(cc, payload)
	if err != nil {
		return res.ErrorResponse(err).Send(c)
	}

	return res.SuccessResponse(data).Send(c)
}

// Register
// @Summary Register user
// @Description Register user
// @Tags auth
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body dto.AuthRegisterRequest true "request body"
// @Success 200 {object} dto.AuthRegisterResponseDoc
// @Failure 400 {object} res.errorResponse
// @Failure 404 {object} res.errorResponse
// @Failure 500 {object} res.errorResponse
// @Router /auth/register [post]
// func (h *handler) Register(c echo.Context) error {
// 	cc := c.(*abstraction.Context)

// 	payload := new(dto.AuthRegisterRequest)
// 	if err = c.Bind(payload); err != nil {
// 		return res.ErrorBuilder(&res.ErrorConstant.BadRequest, err).Send(c)
// 	}
// 	if err = c.Validate(payload); err != nil {
// 		return res.ErrorBuilder(&res.ErrorConstant.Validation, err).Send(c)
// 	}

// 	data, err := h.service.Register(cc, payload)
// 	if err != nil {
// 		return res.ErrorResponse(err).Send(c)
// 	}

// 	return res.SuccessResponse(data).Send(c)
// }

func (h *handler) CheckAuth(c echo.Context) error {
	cc := c.(*abstraction.Context)
	authToken := c.Request().Header.Get("Authorization")
	if authToken == "" {
		return res.ErrorBuilder(&res.ErrorConstant.Unauthorized, nil).Send(c)
	}
	data, err := h.service.CheckAuth(cc, authToken)
	if err != nil {
		return res.ErrorResponse(err).Send(c)
	}

	return res.SuccessResponse(data).Send(c)
}

// Change Password
// @Summary Change Password User
// @Description Change Password User
// @Tags auth
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body dto.AuthRegisterRequest true "request body"
// @Success 200 {string} string "Password has been change!"
// @Failure 400 {object} res.errorResponse
// @Failure 404 {object} res.errorResponse
// @Failure 500 {object} res.errorResponse
// @Router /auth/change-password [patch]
func (h *handler) ChangePassword(c echo.Context) error {
	cc := c.(*abstraction.Context)

	payload := new(dto.ChangePasswordRequest)
	if err = c.Bind(payload); err != nil {
		return res.ErrorBuilder(&res.ErrorConstant.BadRequest, err).Send(c)
	}
	if err = c.Validate(payload); err != nil {
		return res.ErrorBuilder(&res.ErrorConstant.Validation, err).Send(c)
	}

	err := h.service.ChangePassword(cc, payload)
	if err != nil {
		return res.ErrorResponse(err).Send(c)
	}

	return res.SuccessResponse("Password has been change!").Send(c)
}

func (h *handler) GetNotificationToken(c echo.Context) error {
	cc := c.(*abstraction.Context)
	data, err := h.service.GetNotificationToken(cc)
	if err != nil {
		return res.ErrorResponse(err).Send(c)
	}

	return res.SuccessResponse(data).Send(c)
}
