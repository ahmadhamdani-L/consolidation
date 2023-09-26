package consolidation

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
// @Summary Get Consolidation
// @Description Get Consolidation
// @Tags Consolidation
// @Accept json
// @Produce json
// @Security BearerAuth
// @param request query dto.ConsolidationGetRequest true "request query"
// @Success 200 {object} dto.ConsolidationGetResponseDoc
// @Failure 400 {object} response.errorResponse
// @Failure 404 {object} response.errorResponse
// @Failure 500 {object} response.errorResponse
// @Router /consolidation [get]
func (h *handler) Get(c echo.Context) error {
	cc := c.(*abstraction.Context)
	payload := new(dto.ConsolidationGetRequest)

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

// Get List Available
// @Summary Get Consolidation
// @Description Get Consolidation
// @Tags Consolidation
// @Accept json
// @Produce json
// @Security BearerAuth
// @param trial_balance_id path int true "trial_balance_id path"
// @Success 200 {object} dto.ConsolidationGetListAvailableResponseDoc
// @Failure 400 {object} response.errorResponse
// @Failure 404 {object} response.errorResponse
// @Failure 500 {object} response.errorResponse
// @Router /consolidation/list_company/{company_id} [get]
func (h *handler) GetListAvailable(c echo.Context) error {
	cc := c.(*abstraction.Context)
	payload := new(dto.ConsolidationGetListAvailable)

	if err := c.Bind(payload); err != nil {
		return response.ErrorBuilder(&response.ErrorConstant.BadRequest, err).Send(c)
	}

	if err := c.Validate(payload); err != nil {
		return response.ErrorBuilder(&response.ErrorConstant.Validation, err).Send(c)
	}

	result, err := h.service.ListAvailable(cc, payload)
	if err != nil {
		return response.ErrorResponse(err).Send(c)
	}
	return response.CustomSuccessBuilder(http.StatusOK, result, "Get Data Success", &result.PaginationInfo).Send(c)
}

// Get List Dupliacte Available
// @Summary Get Dupliacte Consolidation
// @Description Get Dupliacte Consolidation
// @Tags Consolidation
// @Accept json
// @Produce json
// @Security BearerAuth
// @param conolidation_id path int true "consolidation_id path"
// @param company_id path int true "company_id path"
// @Success 200 {object} dto.ConsolidationGetListDuplicateAvailableResponseDoc
// @Failure 400 {object} response.errorResponse
// @Failure 404 {object} response.errorResponse
// @Failure 500 {object} response.errorResponse
// @Router /consolidation/list_company_duplicate/{consolidation_id} [get]
func (h *handler) GetListDuplicateAvailable(c echo.Context) error {
	cc := c.(*abstraction.Context)
	payload := new(dto.ConsolidationGetListDuplicateAvailable)

	if err := c.Bind(payload); err != nil {
		return response.ErrorBuilder(&response.ErrorConstant.BadRequest, err).Send(c)
	}

	if err := c.Validate(payload); err != nil {
		return response.ErrorBuilder(&response.ErrorConstant.Validation, err).Send(c)
	}

	result, err := h.service.ListDuplicateAvailable(cc, payload)
	if err != nil {
		return response.ErrorResponse(err).Send(c)
	}
	return response.CustomSuccessBuilder(http.StatusOK, result, "Get Data Success", &result.PaginationInfo).Send(c)
}

// Request To Consolidation
// @Summary Request To Consolidation
// @Description Request To Consolidation
// @Tags Consolidation
// @Accept json
// @Produce json
// @Security BearerAuth
// @param request body dto.ConsolidationConsolidateRequest true "request body"
// @Success 200 {object} dto.ConsolidationConsolidateResponseDoc
// @Failure 400 {object} response.errorResponse
// @Failure 404 {object} response.errorResponse
// @Failure 500 {object} response.errorResponse
// @Router /consolidation [post]
func (h *handler) Combaine(c echo.Context) error {
	cc := c.(*abstraction.Context)
	payload := new(dto.ConsolidationCombaineRequest)

	if err := c.Bind(payload); err != nil {
		return response.ErrorBuilder(&response.ErrorConstant.BadRequest, err).Send(c)
	}

	if err := c.Validate(payload); err != nil {
		return response.ErrorBuilder(&response.ErrorConstant.Validation, err).Send(c)
	}

	err := h.service.RequestToCombaine(cc, payload)
	if err != nil {
		return response.ErrorResponse(err).Send(c)
	}

	return response.SuccessResponse("Success to request consolidation").Send(c)
}

// Request To Consolidation
// @Summary Request To Consolidation
// @Description Request To Consolidation
// @Tags Consolidation
// @Accept json
// @Produce json
// @Security BearerAuth
// @param request body dto.ConsolidationConsolidateRequest true "request body"
// @Success 200 {object} dto.ConsolidationConsolidateResponseDoc
// @Failure 400 {object} response.errorResponse
// @Failure 404 {object} response.errorResponse
// @Failure 500 {object} response.errorResponse
// @Router /consolidation [post]
func (h *handler) Duplicate(c echo.Context) error {
	cc := c.(*abstraction.Context)
	payload := new(dto.ConsolidationConsolidateRequest)

	if err := c.Bind(payload); err != nil {
		return response.ErrorBuilder(&response.ErrorConstant.BadRequest, err).Send(c)
	}

	if err := c.Validate(payload); err != nil {
		return response.ErrorBuilder(&response.ErrorConstant.Validation, err).Send(c)
	}

	err := h.service.RequestToDuplicate(cc, payload)
	if err != nil {
		return response.ErrorResponse(err).Send(c)
	}

	return response.SuccessResponse("Success to request consolidation").Send(c)
}

// Request To Consolidation
// @Summary Request To Consolidation
// @Description Request To Consolidation
// @Tags Consolidation
// @Accept json
// @Produce json
// @Security BearerAuth
// @param request body dto.ConsolidationConsolidateRequest true "request body"
// @Success 200 {object} dto.ConsolidationConsolidateResponseDoc
// @Failure 400 {object} response.errorResponse
// @Failure 404 {object} response.errorResponse
// @Failure 500 {object} response.errorResponse
// @Router /consolidation [post]
func (h *handler) EditCombain(c echo.Context) error {
	cc := c.(*abstraction.Context)
	payload := new(dto.ConsolidationConsolidateRequest)

	if err := c.Bind(payload); err != nil {
		return response.ErrorBuilder(&response.ErrorConstant.BadRequest, err).Send(c)
	}

	if err := c.Validate(payload); err != nil {
		return response.ErrorBuilder(&response.ErrorConstant.Validation, err).Send(c)
	}

	err := h.service.RequestToEditCombain(cc, payload)
	if err != nil {
		return response.ErrorResponse(err).Send(c)
	}

	return response.SuccessResponse("Success to request consolidation").Send(c)
}

// Request To Consolidation
// @Summary Request To Consolidation
// @Description Request To Consolidation
// @Tags Consolidation
// @Accept json
// @Produce json
// @Security BearerAuth
// @param request body dto.ConsolidationConsolidateRequest true "request body"
// @Success 200 {object} dto.ConsolidationConsolidateResponseDoc
// @Failure 400 {object} response.errorResponse
// @Failure 404 {object} response.errorResponse
// @Failure 500 {object} response.errorResponse
// @Router /consolidation [post]
func (h *handler) Consolidation(c echo.Context) error {
	cc := c.(*abstraction.Context)
	payload := new(dto.ConsolidationConsolidateRequest)

	if err := c.Bind(payload); err != nil {
		return response.ErrorBuilder(&response.ErrorConstant.BadRequest, err).Send(c)
	}

	if err := c.Validate(payload); err != nil {
		return response.ErrorBuilder(&response.ErrorConstant.Validation, err).Send(c)
	}

	err := h.service.RequestToConsolidation(cc, payload)
	if err != nil {
		return response.ErrorResponse(err).Send(c)
	}

	return response.SuccessResponse("Success to request consolidation").Send(c)
}

// Get Version
// @Summary Get Consolidation Version
// @Description Get Consolidation Version
// @Tags Consolidation
// @Accept json
// @Produce json
// @Security BearerAuth
// @param request query dto.GetVersionRequest true "request query"
// @Success 200 {object} dto.GetVersionResponseDoc
// @Failure 400 {object} response.errorResponse
// @Failure 404 {object} response.errorResponse
// @Failure 500 {object} response.errorResponse
// @Router /consolidation/get-version [get]
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

	statusPayload, err := helper.MultiStatusFilter(c.Request().URL.Query())
	if err != nil {
		return response.ErrorBuilder(&response.ErrorConstant.BadRequest, err).Send(c)
	}
	payload.ArrStatus = &statusPayload

	if err := c.Validate(payload); err != nil {
		return response.ErrorBuilder(&response.ErrorConstant.Validation, err).Send(c)
	}

	result, err := h.service.GetVersion(cc, payload)
	if err != nil {
		return response.ErrorResponse(err).Send(c)
	}
	return response.SuccessResponse(result.Data).Send(c)
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
func (h *handler) FindListCompanyCreateNewCombine(c echo.Context) error {
	cc := c.(*abstraction.Context)
	payload := new(dto.FindListCompanyCreateNewCombineGetRequest)

	if err := c.Bind(payload); err != nil {
		return response.ErrorBuilder(&response.ErrorConstant.BadRequest, err).Send(c)
	}
	if err := c.Validate(payload); err != nil {
		return response.ErrorBuilder(&response.ErrorConstant.Validation, err).Send(c)
	}

	result, err := h.service.FindListCompanyCreateNewCombine(cc, payload)
	if err != nil {
		return response.ErrorResponse(err).Send(c)
	}
	return response.CustomSuccessBuilder(http.StatusOK, result.Datas, "Get Data Success", &result.PaginationInfo).Send(c)
}

// Delete godoc
// @Summary Delete Consolidation
// @Description Delete Consolidation
// @Tags Consolidation
// @Accept json
// @Produce json
// @Security BearerAuth
// @param id path int true "id path"
// @Success 200 {object} dto.ConsolidationDeleteResponseDoc
// @Failure 400 {object} response.errorResponse
// @Failure 404 {object} response.errorResponse
// @Failure 500 {object} response.errorResponse
// @Router /consolidation/{id} [delete]
func (h *handler) Delete(c echo.Context) error {
	cc := c.(*abstraction.Context)
	payload := new(dto.ConsolidationDeleteRequest)

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

func (h *handler) GetByID(c echo.Context) error {
	cc := c.(*abstraction.Context)

	payload := new(dto.ConsolidationGetByIDRequest)
	if err := c.Bind(payload); err != nil {
		return response.ErrorBuilder(&response.ErrorConstant.BadRequest, err).Send(c)
	}
	if err := c.Validate(payload); err != nil {
		response := response.ErrorBuilder(&response.ErrorConstant.Validation, err)
		return response.Send(c)
	}

	fmt.Printf("%+v", payload)

	result, err := h.service.FindByID(cc, payload)
	if err != nil {
		return response.ErrorResponse(err).Send(c)
	}

	return response.SuccessResponse(result).Send(c)
}

// Get Control
// @Summary Get Control Consolidation
// @Description Get Control Consolidation
// @Tags Consolidation
// @Accept json
// @Produce json
// @Security BearerAuth
// @param request query dto.ConsolidationGetControlRequest true "request query"
// @Success 200 {object} dto.ConsolidationGetControlResponseDoc
// @Failure 400 {object} response.errorResponse
// @Failure 404 {object} response.errorResponse
// @Failure 500 {object} response.errorResponse
// @Router /consolidation/control/{consolidation_id} [get]
func (h *handler) GetControl(c echo.Context) error {
	cc := c.(*abstraction.Context)
	payload := new(dto.ConsolidationGetControlRequest)

	if err := c.Bind(payload); err != nil {
		return response.ErrorBuilder(&response.ErrorConstant.BadRequest, err).Send(c)
	}

	if err := c.Validate(payload); err != nil {
		return response.ErrorBuilder(&response.ErrorConstant.Validation, err).Send(c)
	}

	result, err := h.service.GetControl(cc, payload)
	if err != nil {
		return response.ErrorResponse(err).Send(c)
	}
	return response.SuccessResponse(result).Send(c)
}
