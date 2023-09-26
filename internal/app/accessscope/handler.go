package accessscope

import (
	"mcash-finance-console-core/internal/abstraction"
	"mcash-finance-console-core/internal/dto"
	"mcash-finance-console-core/internal/factory"
	modelhelper "mcash-finance-console-core/internal/model/helper"
	"mcash-finance-console-core/pkg/util/response"
	"net/http"

	"github.com/labstack/echo/v4"
)

type handler struct {
	service *service
}

func NewHandler(f *factory.Factory) *handler {
	return &handler{
		service: NewService(f),
	}
}

// Get
// @Summary Get Access Scope
// @Description Get Access Scope
// @Tags Access Scope
// @Accept json
// @Produce json
// @Security BearerAuth
// @param request query dto.AccessScopeGetRequest true "request query"
// @Success 200 {object} dto.AccessScopeGetResponseDoc
// @Failure 400 {object} response.errorResponse
// @Failure 404 {object} response.errorResponse
// @Failure 500 {object} response.errorResponse
// @Router /access-scope [get]
func (h *handler) Get(c echo.Context) error {
	cc := c.(*abstraction.Context)
	payload := new(dto.AccessScopeGetRequest)
	if err := c.Bind(payload); err != nil {
		return response.ErrorBuilder(&response.ErrorConstant.BadRequest, err).Send(c)
	}
	if err := c.Validate(payload); err != nil {
		return response.ErrorBuilder(&response.ErrorConstant.Validation, err).Send(c)
	}

	companyPayload, err := modelhelper.MultiCompanyFilter(c.Request().URL.Query())
	if err != nil {
		return response.ErrorBuilder(&response.ErrorConstant.BadRequest, err).Send(c)
	}
	if companyPayload.ArrCompanyID != nil && len(*companyPayload.ArrCompanyID) > 0 {
		payload.CompanyCustomFilter = companyPayload
	}
	result, err := h.service.Find(cc, payload)
	if err != nil {
		return response.ErrorResponse(err).Send(c)
	}
	return response.CustomSuccessBuilder(http.StatusOK, result.Datas, "Get Data Success", &result.PaginationInfo).Send(c)
}

// Get By ID
// @Summary Get Access Scope by id
// @Description Get Access Scope by id
// @Tags Access Scope
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "id path"
// @Success 200 {object} dto.AccessScopeGetByIDResponseDoc
// @Failure 400 {object} response.errorResponse
// @Failure 404 {object} response.errorResponse
// @Failure 500 {object} response.errorResponse
// @Router /access-scope/{id} [get]
func (h *handler) GetByID(c echo.Context) error {
	cc := c.(*abstraction.Context)
	payload := new(dto.AccessScopeGetByIDRequest)

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
// @Summary Create Access Scope
// @Description Create Access Scope
// @Tags Access Scope
// @Accept  json
// @Produce  json
// @Security BearerAuth
// @Param request body dto.AccessScopeCreateRequest true "request body"
// @Success 200 {object} dto.AccessScopeCreateResponseDoc
// @Failure 400 {object} response.errorResponse
// @Failure 404 {object} response.errorResponse
// @Failure 500 {object} response.errorResponse
// @Router /access-scope [post]
func (h *handler) Create(c echo.Context) error {
	cc := c.(*abstraction.Context)
	payload := new(dto.AccessScopeCreateRequest)

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
// @Summary Update Access Scope
// @Description Update Access Scope
// @Tags Access Scope
// @Accept  json
// @Produce  json
// @Security BearerAuth
// @Param id path int true "id path"
// @Param request body dto.AccessScopeUpdateRequest true "request body"
// @Success 200 {object} dto.AccessScopeUpdateResponseDoc
// @Failure 400 {object} response.errorResponse
// @Failure 404 {object} response.errorResponse
// @Failure 500 {object} response.errorResponse
// @Router /access-scope/{id} [patch]
func (h *handler) Update(c echo.Context) error {
	cc := c.(*abstraction.Context)
	payload := new(dto.AccessScopeUpdateRequest)

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

func (h *handler) Delete(c echo.Context) error {
	return response.ErrorBuilder(&response.ErrorConstant.RouteNotFound, nil).Send(c)
	// cc := c.(*abstraction.Context)
	// payload := new(dto.AccessScopeDeleteRequest)

	// if err := c.Bind(payload); err != nil {
	// 	return response.ErrorBuilder(&response.ErrorConstant.BadRequest, err).Send(c)
	// }
	// if err := c.Validate(payload); err != nil {
	// 	return response.ErrorBuilder(&response.ErrorConstant.Validation, err).Send(c)
	// }

	// result, err := h.service.Delete(cc, payload)
	// if err != nil {
	// 	return response.ErrorResponse(err).Send(c)
	// }

	// return response.SuccessResponse(result).Send(c)
}

// Reset godoc
// @Summary Reset Access Scope
// @Description Reset Access Scope
// @Tags Access Scope
// @Accept  json
// @Produce  json
// @Security BearerAuth
// @Param id path int true "id path"
// @Success 200 {object} dto.AccessScopeDeleteResponseDoc
// @Failure 400 {object} response.errorResponse
// @Failure 404 {object} response.errorResponse
// @Failure 500 {object} response.errorResponse
// @Router /access-scope/{id} [delete]
func (h *handler) Reset(c echo.Context) error {
	cc := c.(*abstraction.Context)
	payload := new(dto.AccessScopeDeleteRequest)

	if err := c.Bind(payload); err != nil {
		return response.ErrorBuilder(&response.ErrorConstant.BadRequest, err).Send(c)
	}
	if err := c.Validate(payload); err != nil {
		return response.ErrorBuilder(&response.ErrorConstant.Validation, err).Send(c)
	}

	result, err := h.service.Reset(cc, payload)
	if err != nil {
		return response.ErrorResponse(err).Send(c)
	}

	return response.SuccessResponse(result).Send(c)
}
