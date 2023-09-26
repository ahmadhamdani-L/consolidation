package jpmdetail

import (
	"fmt"
	"mcash-finance-console-core/internal/abstraction"
	"mcash-finance-console-core/internal/dto"
	"mcash-finance-console-core/internal/factory"
	"mcash-finance-console-core/pkg/util/response"
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

// Get
// @Summary Get Jpm Detail
// @Description Get Jpm Detail
// @Tags Jpm Detail
// @Accept json
// @Produce json
// @Security BearerAuth
// @param request query dto.JpmDetailGetRequest true "request query"
// @Success 200 {object} dto.JpmDetailGetResponseDoc
// @Failure 400 {object} response.errorResponse
// @Failure 404 {object} response.errorResponse
// @Failure 500 {object} response.errorResponse
// @Router /jpm-detail [get]
func (h *handler) Get(c echo.Context) error {
	cc := c.(*abstraction.Context)

	payload := new(dto.JpmDetailGetRequest)
	if err := c.Bind(payload); err != nil {
		return res.ErrorBuilder(&res.ErrorConstant.BadRequest, err).Send(c)
	}
	if err = c.Validate(payload); err != nil {
		return res.ErrorBuilder(&res.ErrorConstant.Validation, err).Send(c)
	}

	result, err := h.service.Find(cc, payload)
	if err != nil {
		return res.ErrorResponse(err).Send(c)
	}

	return res.CustomSuccessBuilder(200, result.Datas, "Get datas success", &result.PaginationInfo).Send(c)
}

// Get By ID
// @Summary Get Jpm Detail by id
// @Description Get Jpm Detail by id
// @Tags Jpm Detail
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "id path"
// @Success 200 {object} dto.JpmDetailGetByIDResponseDoc
// @Failure 400 {object} res.errorResponse
// @Failure 404 {object} res.errorResponse
// @Failure 500 {object} res.errorResponse
// @Router /jpm-detail/{id} [get]
func (h *handler) GetByID(c echo.Context) error {
	cc := c.(*abstraction.Context)

	payload := new(dto.JpmDetailGetByIDRequest)
	if err = c.Bind(payload); err != nil {
		return res.ErrorBuilder(&res.ErrorConstant.BadRequest, err).Send(c)
	}
	if err = c.Validate(payload); err != nil {
		response := res.ErrorBuilder(&res.ErrorConstant.Validation, err)
		return response.Send(c)
	}

	fmt.Printf("%+v", payload)

	result, err := h.service.FindByID(cc, payload)
	if err != nil {
		return res.ErrorResponse(err).Send(c)
	}

	return res.SuccessResponse(result).Send(c)
}

// Update godoc
// @Summary Update Jpm Detail
// @Description Update Jpm Detail
// @Tags Jpm Detail
// @Accept json
// @Produce json
// @Security BearerAuth
// @param id path int true "id path"
// @param request body dto.JpmDetailUpdateRequest true "request body"
// @Success 200 {object} dto.JpmDetailUpdateResponseDoc
// @Failure 400 {object} response.errorResponse
// @Failure 404 {object} response.errorResponse
// @Failure 500 {object} response.errorResponse
// @Router /jpm-detail/{id} [patch]
func (h *handler) Update(c echo.Context) error {
	cc := c.(*abstraction.Context)
	payload := new(dto.JpmDetailUpdateRequest)

	if err := c.Bind(payload); err != nil {
		return response.ErrorBuilder(&response.ErrorConstant.BadRequest, err).Send(c)
	}
	if err := c.Validate(payload); err != nil {
		return response.ErrorBuilder(&response.ErrorConstant.Validation, err).Send(c)
	}

	result, err := h.service.Update(cc, payload)
	if err != nil {
		return response.ErrorResponse(err).Send(c)
	}

	return response.SuccessResponse(result).Send(c)
}

// Create godoc
// @Summary Get jpm Detail
// @Description Create jpm Detail
// @Tags jpm Detail
// @Accept json
// @Produce json
// @Security BearerAuth
// @param request body dto.JpmDetailCreateRequest true "request body"
// @Success 200 {object} dto.JpmDetailCreateResponseDoc
// @Failure 400 {object} response.errorResponse
// @Failure 404 {object} response.errorResponse
// @Failure 500 {object} response.errorResponse
// @Router /jpm-detail [post]
func (h *handler) Create(c echo.Context) error {
	cc := c.(*abstraction.Context)

	payload := new(dto.JpmDetailCreateRequest)

	if err := c.Bind(payload); err != nil {
		return res.ErrorBuilder(&res.ErrorConstant.BadRequest, err).Send(c)
	}
	if err := c.Validate(payload); err != nil {
		return res.ErrorBuilder(&res.ErrorConstant.Validation, err).Send(c)
	}

	result, err := h.service.Create(cc, payload)
	if err != nil {
		return res.ErrorResponse(err).Send(c)
	}

	return res.SuccessResponse(result).Send(c)
}
