package company

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
// @Summary Get Company
// @Description Get Company
// @Tags Company
// @Accept json
// @Produce json
// @Security BearerAuth
// @param request query dto.CompanyGetRequest true "request query"
// @Success 200 {object} dto.CompanyGetResponseDoc
// @Failure 400 {object} response.errorResponse
// @Failure 404 {object} response.errorResponse
// @Failure 500 {object} response.errorResponse
// @Router /company [get]
func (h *handler) Get(c echo.Context) error {
	cc := c.(*abstraction.Context)
	payload := new(dto.CompanyGetRequest)

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
// @Summary Get Company By ID
// @Description Get Company By ID
// @Tags Company
// @Accept json
// @Produce json
// @Security BearerAuth
// @param id path int true "id path"
// @Success 200 {object} dto.CompanyGetByIDResponseDoc
// @Failure 400 {object} response.errorResponse
// @Failure 404 {object} response.errorResponse
// @Failure 500 {object} response.errorResponse
// @Router /company/{id} [get]
func (h *handler) GetByID(c echo.Context) error {
	cc := c.(*abstraction.Context)
	payload := new(dto.CompanyGetByIDRequest)

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
// @Summary Create Company
// @Description Create Company
// @Tags Company
// @Accept json
// @Produce json
// @Security BearerAuth
// @param request body dto.CompanyCreateRequest true "request body"
// @Success 200 {object} dto.CompanyCreateResponseDoc
// @Failure 400 {object} response.errorResponse
// @Failure 404 {object} response.errorResponse
// @Failure 500 {object} response.errorResponse
// @Router /company [post]
func (h *handler) Create(c echo.Context) error {
	cc := c.(*abstraction.Context)
	payload := new(dto.CompanyCreateRequest)

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
// @Summary Update Company
// @Description Update Company
// @Tags Company
// @Accept json
// @Produce json
// @Security BearerAuth
// @param id path int true "id path"
// @param request body dto.CompanyUpdateRequest true "request body"
// @Success 200 {object} dto.CompanyUpdateResponseDoc
// @Failure 400 {object} response.errorResponse
// @Failure 404 {object} response.errorResponse
// @Failure 500 {object} response.errorResponse
// @Router /company{id} [patch]
func (h *handler) Update(c echo.Context) error {
	cc := c.(*abstraction.Context)
	payload := new(dto.CompanyUpdateRequest)

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
// @Summary Delete Company
// @Description Delete Company
// @Tags Company
// @Accept json
// @Produce json
// @Security BearerAuth
// @param id path int true "id path"
// @Success 200 {object} dto.CompanyDeleteResponseDoc
// @Failure 400 {object} response.errorResponse
// @Failure 404 {object} response.errorResponse
// @Failure 500 {object} response.errorResponse
// @Router /company [delete]
func (h *handler) Delete(c echo.Context) error {
	cc := c.(*abstraction.Context)
	payload := new(dto.CompanyDeleteRequest)

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

// Get Treeview
// @Summary Get Treeview Company
// @Description Get Treeview Company
// @Tags Company
// @Accept json
// @Produce json
// @Security BearerAuth
// @param request query dto.CompanyGetRequest true "request query"
// @Success 200 {object} dto.CompanyGetResponseDoc
// @Failure 400 {object} response.errorResponse
// @Failure 404 {object} response.errorResponse
// @Failure 500 {object} response.errorResponse
// @Router /company/treeview [get]
func (h *handler) GetTreeview(c echo.Context) error {
	cc := c.(*abstraction.Context)

	result, err := h.service.FindTree(cc)
	if err != nil {
		return response.ErrorResponse(err).Send(c)
	}
	return response.SuccessResponse(result.Datas).Send(c)
}

// Get List Filter Modul
// @Summary Get List Filter Modul
// @Description Get List for filtering in modul
// @Tags Company
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} dto.CompanyGetResponseDoc
// @Failure 400 {object} response.errorResponse
// @Failure 404 {object} response.errorResponse
// @Failure 500 {object} response.errorResponse
// @Router /company/get-filter [get]
func (h *handler) GetFilterList(c echo.Context) error {
	cc := c.(*abstraction.Context)

	result, err := h.service.FindListFilter(cc)
	if err != nil {
		return response.ErrorResponse(err).Send(c)
	}
	return response.SuccessResponse(result.Datas).Send(c)
}
