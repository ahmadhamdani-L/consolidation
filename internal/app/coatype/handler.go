package coatype

import (
	"mcash-finance-console-core/internal/abstraction"
	"mcash-finance-console-core/internal/dto"
	"mcash-finance-console-core/internal/factory"
	"mcash-finance-console-core/pkg/util/helper"
	"mcash-finance-console-core/pkg/util/response"
	"net/http"

	"github.com/labstack/echo/v4"
)

type handler struct {
	service *service
}

var err error

func NewHandler(f *factory.Factory) *handler {
	return &handler{
		service: NewService(f),
	}
}

// Get
// @Summary Get COA Type
// @Description Get COA Type
// @Tags COA Type
// @Accept json
// @Produce json
// @Security BearerAuth
// @param request query dto.CoaTypeGetRequest true "request query"
// @Success 200 {object} dto.CoaTypeGetResponseDoc
// @Failure 400 {object} response.errorResponse
// @Failure 404 {object} response.errorResponse
// @Failure 500 {object} response.errorResponse
// @Router /coa-type [get]
func (h *handler) Get(c echo.Context) error {
	cc := c.(*abstraction.Context)
	payload := new(dto.CoaTypeGetRequest)

	if err := c.Bind(payload); err != nil {
		return response.ErrorBuilder(&response.ErrorConstant.BadRequest, err).Send(c)
	}
	coaGroupPayload, err := helper.MultiCoaGroupIDFilter(c.Request().URL.Query())
	if err != nil {
		return response.ErrorBuilder(&response.ErrorConstant.BadRequest, err).Send(c)
	}
	payload.ArrCoaGroupID = &coaGroupPayload
	if err := c.Validate(payload); err != nil {
		return response.ErrorBuilder(&response.ErrorConstant.Validation, err).Send(c)
	}

	result, err := h.service.Find(cc, payload)
	if err != nil {
		return response.ErrorResponse(err).Send(c)
	}
	return response.CustomSuccessBuilder(http.StatusOK, result.Datas, "Get Data Success", &result.PaginationInfo).Send(c)
}

// Get By ID
// @Summary Get COA Type By ID
// @Description Get COA Type By ID
// @Tags COA Type
// @Accept json
// @Produce json
// @Security BearerAuth
// @param id path int true "id path"
// @Success 200 {object} dto.CoaGetByIDResponseDoc
// @Failure 400 {object} response.errorResponse
// @Failure 404 {object} response.errorResponse
// @Failure 500 {object} response.errorResponse
// @Router /samples [get]
func (h *handler) GetByID(c echo.Context) error {
	cc := c.(*abstraction.Context)
	payload := new(dto.CoaTypeGetByIDRequest)

	if err := c.Bind(payload); err != nil {
		return response.ErrorBuilder(&response.ErrorConstant.BadRequest, err).Send(c)
	}
	if err := c.Validate(payload); err != nil {
		return response.ErrorBuilder(&response.ErrorConstant.Validation, err).Send(c)
	}

	result, err := h.service.FindByID(cc, payload)
	if err != nil {
		return response.ErrorResponse(err).Send(c)
	}

	return response.SuccessResponse(result).Send(c)
}

// Create godoc
// @Summary Create COA Type
// @Description Create COA Type
// @Tags COA Type
// @Accept json
// @Produce json
// @Security BearerAuth
// @param request body dto.CoaTypeCreateRequest true "request body"
// @Success 200 {object} dto.CoaTypeCreateResponseDoc
// @Failure 400 {object} response.errorResponse
// @Failure 404 {object} response.errorResponse
// @Failure 500 {object} response.errorResponse
// @Router /coa-type [post]
func (h *handler) Create(c echo.Context) error {
	cc := c.(*abstraction.Context)
	payload := new(dto.CoaTypeCreateRequest)

	if err := c.Bind(payload); err != nil {
		return response.ErrorBuilder(&response.ErrorConstant.BadRequest, err).Send(c)
	}

	if err := c.Validate(payload); err != nil {
		return response.ErrorBuilder(&response.ErrorConstant.Validation, err).Send(c)
	}

	result, err := h.service.Create(cc, payload)
	if err != nil {
		return response.ErrorResponse(err).Send(c)
	}

	return response.SuccessResponse(result).Send(c)
}

// Update godoc
// @Summary Update COA GROUP
// @Description Update COA GROUP
// @Tags COA GROUP
// @Accept json
// @Produce json
// @Security BearerAuth
// @param id path int true "id path"
// @param request body dto.CoaTypeUpdateRequest true "request query"
// @Success 200 {object} dto.CoaTypeUpdateResponseDoc
// @Failure 400 {object} response.errorResponse
// @Failure 404 {object} response.errorResponse
// @Failure 500 {object} response.errorResponse
// @Router /coa-type/{id} [patch]
func (h *handler) Update(c echo.Context) error {
	cc := c.(*abstraction.Context)
	payload := new(dto.CoaTypeUpdateRequest)

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

// Delete godoc
// @Summary Delete COA Type
// @Description Delete COA Type
// @Tags COA Type
// @Accept json
// @Produce json
// @Security BearerAuth
// @param id path int true "id path"
// @Success 200 {object} dto.CoaDeleteResponseDoc
// @Failure 400 {object} response.errorResponse
// @Failure 404 {object} response.errorResponse
// @Failure 500 {object} response.errorResponse
// @Router /coa-type [delete]
func (h *handler) Delete(c echo.Context) error {
	cc := c.(*abstraction.Context)
	payload := new(dto.CoaTypeDeleteRequest)

	if err := c.Bind(payload); err != nil {
		return response.ErrorBuilder(&response.ErrorConstant.BadRequest, err).Send(c)
	}
	if err := c.Validate(payload); err != nil {
		return response.ErrorBuilder(&response.ErrorConstant.Validation, err).Send(c)
	}

	result, err := h.service.Delete(cc, payload)
	if err != nil {
		return response.ErrorResponse(err).Send(c)
	}

	return response.SuccessResponse(result).Send(c)
}
