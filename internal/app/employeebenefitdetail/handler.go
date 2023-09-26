package employeebenefitdetail

import (
	"mcash-finance-console-core/internal/abstraction"
	"mcash-finance-console-core/internal/dto"
	"mcash-finance-console-core/internal/factory"
	"mcash-finance-console-core/pkg/util/response"

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
// @Summary Get Employee Benefit Detail
// @Description Get Employee Benefit Detail
// @Tags Employee Benefit Detail
// @Accept json
// @Produce json
// @Security BearerAuth
// @param request query dto.EmployeeBenefitDetailGetRequest true "request query"
// @Success 200 {object} dto.EmployeeBenefitDetailGetResponseDoc
// @Failure 400 {object} response.errorResponse
// @Failure 404 {object} response.errorResponse
// @Failure 500 {object} response.errorResponse
// @Router /employee-benefit-detail [get]
func (h *handler) Get(c echo.Context) error {
	cc := c.(*abstraction.Context)
	payload := new(dto.EmployeeBenefitDetailGetRequest)
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
	return response.SuccessResponse(result.Datas).Send(c)
}

// Get By ID
// @Summary Get Employee Benefit Detail By ID
// @Description Get Employee Benefit Detail By ID
// @Tags Employee Benefit Detail
// @Accept json
// @Produce json
// @Security BearerAuth
// @param id path int true "id path"
// @Success 200 {object} dto.EmployeeBenefitDetailGetByIDResponseDoc
// @Failure 400 {object} response.errorResponse
// @Failure 404 {object} response.errorResponse
// @Failure 500 {object} response.errorResponse
// @Router /employee-benefit-detail/{id} [get]
func (h *handler) GetByID(c echo.Context) error {
	cc := c.(*abstraction.Context)
	payload := new(dto.EmployeeBenefitDetailGetByIDRequest)

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

func (h *handler) Create(c echo.Context) error {
	cc := c.(*abstraction.Context)
	payload := new(dto.EmployeeBenefitDetailCreateRequest)

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
// @Summary Update Employee Benefit Detail
// @Description Update Employee Benefit Detail
// @Tags Employee Benefit Detail
// @Accept json
// @Produce json
// @Security BearerAuth
// @param id path int true "id path"
// @param request query dto.EmployeeBenefitDetailUpdateRequest true "request query"
// @Success 200 {object} dto.EmployeeBenefitDetailUpdateResponseDoc
// @Failure 400 {object} response.errorResponse
// @Failure 404 {object} response.errorResponse
// @Failure 500 {object} response.errorResponse
// @Router /employee-benefit-detail/{id} [patch]
func (h *handler) Update(c echo.Context) error {
	cc := c.(*abstraction.Context)
	payload := new(dto.EmployeeBenefitDetailUpdateRequest)

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
// @Summary Delete Employee Benefit Detail
// @Description Delete Employee Benefit Detail
// @Tags Employee Benefit Detail
// @Accept json
// @Produce json
// @Security BearerAuth
// @param id path int true "id path"
// @Success 200 {object} dto.EmployeeBenefitDetailDeleteResponseDoc
// @Failure 400 {object} response.errorResponse
// @Failure 404 {object} response.errorResponse
// @Failure 500 {object} response.errorResponse
// @Router /employee-benefit-detail/{id} [delete]
func (h *handler) Delete(c echo.Context) error {
	cc := c.(*abstraction.Context)
	payload := new(dto.EmployeeBenefitDetailDeleteRequest)

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
