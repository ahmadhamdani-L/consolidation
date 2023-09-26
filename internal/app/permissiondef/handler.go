package permissiondef

import (
	"mcash-finance-console-core/internal/abstraction"
	"mcash-finance-console-core/internal/dto"
	"mcash-finance-console-core/internal/factory"
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
// @Summary Get Permission Def
// @Description Get Permission Def
// @Tags Permission Def
// @Accept json
// @Produce json
// @Security BearerAuth
// @param request query dto.PermissionDefGetRequest true "request query"
// @Success 200 {object} dto.PermissionDefGetResponseDoc
// @Failure 400 {object} response.errorResponse
// @Failure 404 {object} response.errorResponse
// @Failure 500 {object} response.errorResponse
// @Router /permission-def [get]
func (h *handler) Get(c echo.Context) error {
	cc := c.(*abstraction.Context)
	payload := new(dto.PermissionDefGetRequest)
	if err := c.Bind(payload); err != nil {
		return response.ErrorBuilder(&response.ErrorConstant.BadRequest, err).Send(c)
	}
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
// @Summary Get Permission Def by id
// @Description Get Permission Def by id
// @Tags Permission Def
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "id path"
// @Success 200 {object} dto.PermissionDefGetByIDResponseDoc
// @Failure 400 {object} response.errorResponse
// @Failure 404 {object} response.errorResponse
// @Failure 500 {object} response.errorResponse
// @Router /permission-def/{id} [get]
func (h *handler) GetByID(c echo.Context) error {
	cc := c.(*abstraction.Context)
	payload := new(dto.PermissionDefGetByIDRequest)

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
// @Summary Create Permission Def
// @Description Create Permission Def
// @Tags Permission Def
// @Accept  json
// @Produce  json
// @Security BearerAuth
// @Param request body dto.PermissionDefCreateRequest true "request body"
// @Success 200 {object} dto.PermissionDefCreateResponseDoc
// @Failure 400 {object} response.errorResponse
// @Failure 404 {object} response.errorResponse
// @Failure 500 {object} response.errorResponse
// @Router /permission-def [post]
func (h *handler) Create(c echo.Context) error {
	cc := c.(*abstraction.Context)
	payload := new(dto.PermissionDefCreateRequest)

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
// @Summary Update Permission Def
// @Description Update Permission Def
// @Tags Permission Def
// @Accept  json
// @Produce  json
// @Security BearerAuth
// @Param id path int true "id path"
// @Param request body dto.PermissionDefUpdateRequest true "request body"
// @Success 200 {object} dto.PermissionDefUpdateResponseDoc
// @Failure 400 {object} response.errorResponse
// @Failure 404 {object} response.errorResponse
// @Failure 500 {object} response.errorResponse
// @Router /permission-def/{id} [patch]
func (h *handler) Update(c echo.Context) error {
	cc := c.(*abstraction.Context)
	payload := new(dto.PermissionDefUpdateRequest)

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
// @Summary Delete Permission Def
// @Description Delete Permission Def
// @Tags Permission Def
// @Accept  json
// @Produce  json
// @Security BearerAuth
// @Param id path int true "id path"
// @Success 200 {object} dto.PermissionDefDeleteResponseDoc
// @Failure 400 {object} response.errorResponse
// @Failure 404 {object} response.errorResponse
// @Failure 500 {object} response.errorResponse
// @Router /permission-def/{id} [delete]
func (h *handler) Delete(c echo.Context) error {
	cc := c.(*abstraction.Context)
	payload := new(dto.PermissionDefDeleteRequest)

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
