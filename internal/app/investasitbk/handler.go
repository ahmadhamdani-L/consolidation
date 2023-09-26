package investasitbk

import (
	"errors"
	"fmt"
	"mcash-finance-console-core/internal/abstraction"
	"mcash-finance-console-core/internal/dto"
	"mcash-finance-console-core/internal/factory"
	modelhelper "mcash-finance-console-core/internal/model/helper"
	"mcash-finance-console-core/pkg/util/helper"
	"mcash-finance-console-core/pkg/util/response"
	"net/http"
	"os"
	"time"

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
// @Summary Get Investasi TBK
// @Description Get Investasi TBK
// @Tags Investasi TBK
// @Accept json
// @Produce json
// @Security BearerAuth
// @param request query dto.InvestasiTbkGetRequest true "request query"
// @Success 200 {object} dto.InvestasiTbkGetResponseDoc
// @Failure 400 {object} response.errorResponse
// @Failure 404 {object} response.errorResponse
// @Failure 500 {object} response.errorResponse
// @Router /investasi-tbk [get]
func (h *handler) Get(c echo.Context) error {
	cc := c.(*abstraction.Context)
	payload := new(dto.InvestasiTbkGetRequest)
	if err := c.Bind(payload); err != nil {
		return response.ErrorBuilder(&response.ErrorConstant.BadRequest, err).Send(c)
	}

	if payload.CompanyID != nil {
		allowed := helper.CompanyValidation(cc.Auth.ID, *payload.CompanyID)
		if !allowed {
			return response.ErrorBuilder(&response.ErrorConstant.BadRequest, errors.New("Not Allowed")).Send(c)
		}
	}

	if payload.CompanyCustomFilter.CompanyID != nil {
		versionPayload, err := helper.MultiVersionFilter(c.Request().URL.Query())
		if err != nil {
			return response.ErrorBuilder(&response.ErrorConstant.BadRequest, err).Send(c)
		}
		payload.ArrVersions = &versionPayload
	} else {
		companyPayload, err := modelhelper.MultiCompanyFilter(c.Request().URL.Query())
		if err != nil {
			return response.ErrorBuilder(&response.ErrorConstant.BadRequest, err).Send(c)
		}
		payload.CompanyCustomFilter = companyPayload
	}
	// currentYear, currentMonth, _ := time.Now().Date()
	// if payload.Period != nil {
	// 	datePeriod, err := time.Parse("2006-01-02", *payload.Period)
	// 	if err != nil {
	// 		return response.ErrorBuilder(&response.ErrorConstant.Validation, err).Send(c)
	// 	}
	// 	currentYear, currentMonth, _ = datePeriod.Date()
	// }
	// currentLocation := time.Now().Location()
	// firstOfMonth := time.Date(currentYear, currentMonth, 1, 0, 0, 0, 0, currentLocation)
	// lastOfMonth := firstOfMonth.AddDate(0, 1, -1)
	// period := lastOfMonth.Format("2006-01-02")
	// payload.Period = &period

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
// @Summary Get Investasi TBK By ID
// @Description Get Investasi TBK By ID
// @Tags Investasi TBK
// @Accept json
// @Produce json
// @Security BearerAuth
// @param id path int true "id path"
// @Success 200 {object} dto.InvestasiTbkGetByIDResponseDoc
// @Failure 400 {object} response.errorResponse
// @Failure 404 {object} response.errorResponse
// @Failure 500 {object} response.errorResponse
// @Router /investasi-tbk/{id} [get]
func (h *handler) GetByID(c echo.Context) error {
	cc := c.(*abstraction.Context)
	payload := new(dto.InvestasiTbkGetByIDRequest)

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
	payload := new(dto.InvestasiTbkCreateRequest)

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
// @Summary Update Investasi TBK
// @Description Update Investasi TBK
// @Tags Investasi TBK
// @Accept json
// @Produce json
// @Security BearerAuth
// @param id path int true "id path"
// @param request body dto.InvestasiTbkUpdateRequest true "request body"
// @Success 200 {object} dto.InvestasiTbkUpdateResponseDoc
// @Failure 400 {object} response.errorResponse
// @Failure 404 {object} response.errorResponse
// @Failure 500 {object} response.errorResponse
// @Router /investasi-tbk/{id} [patch]
func (h *handler) Update(c echo.Context) error {
	cc := c.(*abstraction.Context)
	payload := new(dto.InvestasiTbkUpdateRequest)

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
// @Summary Delete Investasi TBK
// @Description Delete Investasi TBK
// @Tags Investasi TBK
// @Accept json
// @Produce json
// @Security BearerAuth
// @param id path int true "id path"
// @Success 200 {object} dto.InvestasiTbkDeleteResponseDoc
// @Failure 400 {object} response.errorResponse
// @Failure 404 {object} response.errorResponse
// @Failure 500 {object} response.errorResponse
// @Router /investasi-tbk/{id} [delete]
func (h *handler) Delete(c echo.Context) error {
	cc := c.(*abstraction.Context)
	payload := new(dto.InvestasiTbkDeleteRequest)

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

// Export
// @Summary Export Investasi TBK
// @Description Export Investasi TBK
// @Tags Investasi TBK
// @Accept json
// @Produce json
// @Security BearerAuth
// @param investasi_tbk_id query int true "investasi_tbk_id query"
// @Success 200 {object} dto.InvestasiTbkExportResponseDoc
// @Failure 400 {object} response.errorResponse
// @Failure 404 {object} response.errorResponse
// @Failure 500 {object} response.errorResponse
// @Router /investasi-tbk/export [get]
func (h *handler) Export(c echo.Context) error {
	cc := c.(*abstraction.Context)
	payload := new(dto.InvestasiTbkExportRequest)

	if err := c.Bind(payload); err != nil {
		return response.ErrorBuilder(&response.ErrorConstant.BadRequest, err).Send(c)
	}
	if err := c.Validate(payload); err != nil {
		return response.ErrorBuilder(&response.ErrorConstant.Validation, err).Send(c)
	}

	result, err := h.service.Export(cc, payload)
	if err != nil {
		return response.ErrorResponse(err).Send(c)
	}

	f, err := os.Open(result.Path)
	if err != nil {
		return err
	}

	c.Response().Header().Set(echo.HeaderContentDisposition, fmt.Sprintf("attachment; filename=%s", result.FileName))
	return c.Stream(http.StatusOK, "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet", f)
}

// Get Version
// @Summary Get Investasi Tbk Version
// @Description Get Investasi Tbk Version
// @Tags Investasi TBK
// @Accept json
// @Produce json
// @Security BearerAuth
// @param request query dto.GetVersionRequest true "request query"
// @Success 200 {object} dto.GetVersionResponseDoc
// @Failure 400 {object} response.errorResponse
// @Failure 404 {object} response.errorResponse
// @Failure 500 {object} response.errorResponse
// @Router /investasi-tbk/get-version [get]
func (h *handler) GetVersion(c echo.Context) error {
	cc := c.(*abstraction.Context)
	payload := new(dto.GetVersionRequest)

	if err := c.Bind(payload); err != nil {
		return response.ErrorBuilder(&response.ErrorConstant.BadRequest, err).Send(c)
	}

	companyPayload, err := modelhelper.MultiCompanyFilter(c.Request().URL.Query())
	if err != nil {
		return response.ErrorBuilder(&response.ErrorConstant.BadRequest, err).Send(c)
	}
	companyPayload.CompanyID = payload.CompanyID
	payload.CompanyCustomFilter = companyPayload

	currentYear, currentMonth, _ := time.Now().Date()
	if payload.Period != nil {
		datePeriod, err := time.Parse("2006-01-02", *payload.Period)
		if err != nil {
			return response.ErrorBuilder(&response.ErrorConstant.Validation, err).Send(c)
		}
		currentYear, currentMonth, _ = datePeriod.Date()
	}
	currentLocation := time.Now().Location()
	firstOfMonth := time.Date(currentYear, currentMonth, 1, 0, 0, 0, 0, currentLocation)
	lastOfMonth := firstOfMonth.AddDate(0, 1, -1)
	period := lastOfMonth.Format("2006-01-02")
	payload.Period = &period

	if payload.CompanyID != nil {
		allowed := helper.CompanyValidation(cc.Auth.ID, *payload.CompanyID)
		if !allowed {
			return response.ErrorBuilder(&response.ErrorConstant.BadRequest, errors.New("Not Allowed")).Send(c)
		}
	}

	if err := c.Validate(payload); err != nil {
		return response.ErrorBuilder(&response.ErrorConstant.Validation, err).Send(c)
	}

	result, err := h.service.GetVersion(cc, payload)
	if err != nil {
		return response.ErrorResponse(err).Send(c)
	}
	return response.SuccessResponse(result.Data).Send(c)
}
