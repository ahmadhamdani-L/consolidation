package jctedetail

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
// @Summary Get Jcte Detail
// @Description Get Jcte Detail
// @Tags Jcte Detail
// @Accept json
// @Produce json
// @Security BearerAuth
// @param request query dto.JcteDetailGetRequest true "request query"
// @Success 200 {object} dto.JcteDetailGetResponseDoc
// @Failure 400 {object} response.errorResponse
// @Failure 404 {object} response.errorResponse
// @Failure 500 {object} response.errorResponse
// @Router /jcte-detail [get]
func (h *handler) Get(c echo.Context) error {
	cc := c.(*abstraction.Context)

	payload := new(dto.JcteDetailGetRequest)
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
// @Summary Get Jcte Detail by id
// @Description Get Jcte Detail by id
// @Tags Jcte Detail
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "id path"
// @Success 200 {object} dto.JcteDetailGetByIDResponseDoc
// @Failure 400 {object} res.errorResponse
// @Failure 404 {object} res.errorResponse
// @Failure 500 {object} res.errorResponse
// @Router /jcte-detail/{id} [get]
func (h *handler) GetByID(c echo.Context) error {
	cc := c.(*abstraction.Context)

	payload := new(dto.JcteDetailGetByIDRequest)
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
// @Summary Update Jcte Detail
// @Description Update Jcte Detail
// @Tags Jcte Detail
// @Accept json
// @Produce json
// @Security BearerAuth
// @param id path int true "id path"
// @param request body dto.JcteDetailUpdateRequest true "request body"
// @Success 200 {object} dto.JcteDetailUpdateResponseDoc
// @Failure 400 {object} response.errorResponse
// @Failure 404 {object} response.errorResponse
// @Failure 500 {object} response.errorResponse
// @Router /jcte-detail/{id} [patch]
func (h *handler) Update(c echo.Context) error {
	cc := c.(*abstraction.Context)
	payload := new(dto.JcteDetailUpdateRequest)

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
// @Summary Get Jcte Detail
// @Description Create Jcte Detail
// @Tags Jcte Detail
// @Accept json
// @Produce json
// @Security BearerAuth
// @param request body dto.JcteDetailCreateRequest true "request body"
// @Success 200 {object} dto.JcteDetailCreateResponseDoc
// @Failure 400 {object} response.errorResponse
// @Failure 404 {object} response.errorResponse
// @Failure 500 {object} response.errorResponse
// @Router /jcte-detail [post]
func (h *handler) Create(c echo.Context) error {
	cc := c.(*abstraction.Context)

	payload := new(dto.JcteDetailCreateRequest)

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
