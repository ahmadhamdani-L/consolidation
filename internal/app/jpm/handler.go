package jpm

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

var err error

func NewHandler(f *factory.Factory) *handler {
	service := NewService(f)
	return &handler{service}
}

// Get
// @Summary Get Jpm
// @Description Get Jpm
// @Tags Jpm
// @Accept json
// @Produce json
// @Security BearerAuth
// @param request query dto.JpmGetRequest true "request query"
// @Success 200 {object} dto.JpmGetResponseDoc
// @Failure 400 {object} response.errorResponse
// @Failure 404 {object} response.errorResponse
// @Failure 500 {object} response.errorResponse
// @Router /jpm [get]
func (h *handler) Get(c echo.Context) error {
	cc := c.(*abstraction.Context)

	payload := new(dto.JpmGetRequest)
	if err := c.Bind(payload); err != nil {
		return response.ErrorBuilder(&response.ErrorConstant.BadRequest, err).Send(c)
	}
	if payload.CompanyID != nil {
		allowed := helper.CompanyValidation(cc.Auth.ID, *payload.CompanyID)
		if !allowed {
			return response.ErrorBuilder(&response.ErrorConstant.Unauthorized, errors.New("Not Allowed")).Send(c)
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
	if err = c.Validate(payload); err != nil {
		return response.ErrorBuilder(&response.ErrorConstant.Validation, err).Send(c)
	}

	result, err := h.service.Find(cc, payload)
	if err != nil {
		return response.ErrorResponse(err).Send(c)
	}

	return response.CustomSuccessBuilder(200, result.Datas, "Get datas success", &result.PaginationInfo).Send(c)
}

// Get By ID
// @Summary Get Jpm by id
// @Description Get Jpm by id
// @Tags Jpm
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "id path"
// @Success 200 {object} dto.JpmGetByIDResponseDoc
// @Failure 400 {object} response.errorResponse
// @Failure 404 {object} response.errorResponse
// @Failure 500 {object} response.errorResponse
// @Router /jpm/{id} [get]
func (h *handler) GetByID(c echo.Context) error {
	cc := c.(*abstraction.Context)

	payload := new(dto.JpmGetByIDRequest)
	if err = c.Bind(payload); err != nil {
		return response.ErrorBuilder(&response.ErrorConstant.BadRequest, err).Send(c)
	}
	if err = c.Validate(payload); err != nil {
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

func (h *handler) Update(c echo.Context) error {
	cc := c.(*abstraction.Context)

	payload := new(dto.JpmUpdateRequest)
	if err = c.Bind(payload); err != nil {
		return response.ErrorBuilder(&response.ErrorConstant.BadRequest, err).Send(c)
	}
	if err = c.Validate(payload); err != nil {
		response := response.ErrorBuilder(&response.ErrorConstant.Validation, err)
		return response.Send(c)
	}
	result, err := h.service.Update(cc, payload)
	if err != nil {
		return response.ErrorResponse(err).Send(c)
	}

	return response.SuccessResponse(result).Send(c)
}

// Create godoc
// @Summary Get Jpm
// @Description Create Jpm
// @Tags Jpm
// @Accept json
// @Produce json
// @Security BearerAuth
// @param request body dto.JpmCreateRequest true "request body"
// @Success 200 {object} dto.JpmCreateResponseDoc
// @Failure 400 {object} response.errorResponse
// @Failure 404 {object} response.errorResponse
// @Failure 500 {object} response.errorResponse
// @Router /jpm [post]
func (h *handler) Create(c echo.Context) error {
	cc := c.(*abstraction.Context)

	payload := new(dto.JpmCreateRequest)

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

// Delete godoc
// @Summary Delete Jpm
// @Description Delete Jpm
// @Tags Jpm
// @Accept  json
// @Produce  json
// @Security BearerAuth
// @Param id path int true "id path"
// @Success 200 {object}  dto.JpmDeleteResponseDoc
// @Failure 400 {object} response.errorResponse
// @Failure 404 {object} response.errorResponse
// @Failure 500 {object} response.errorResponse
// @Router /jpm/{id} [delete]
func (h *handler) Delete(c echo.Context) error {
	cc := c.(*abstraction.Context)

	payload := new(dto.JpmDeleteRequest)
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
// @Summary Export Jpm
// @Description Export Jpm
// @Tags Jpm
// @Accept  json
// @Produce  json
// @Security BearerAuth
// @Success 200 {object} dto.CoaExportResponseDoc
// @Failure 400 {object} response.errorResponse
// @Failure 404 {object} response.errorResponse
// @Failure 500 {object} response.errorResponse
// @Router /jpm/export [get]
func (h *handler) Export(c echo.Context) error {
	cc := c.(*abstraction.Context)
	payload := new(dto.JpmExportRequest)

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

// 	return response.SuccessResponse(result).Send(c)
// }

// func (h *handler) ExportAsync(c echo.Context) error {
// 	cc := c.(*abstraction.Context)
// 	payload := new(dto.jpmGetRequest)

// 	nolimit := 100000
// 	page := 1
// 	payload.PageSize = &nolimit
// 	payload.Page = &page

// 	if payload.Sort == nil {
// 		coa := "code"
// 		asc := "asc"
// 		payload.SortBy = &coa
// 		payload.Sort = &asc
// 	}

// 	waktu := time.Now()
// 	testing := fmt.Sprintf("Jurnaljpm.xlsx")
// 	map1 := kafka.JsonData{
// 		FileLoc:   testing,
// 		UserID:    1,
// 		CompanyID: 1,
// 		Timestamp: &waktu,
// 		Name:      cc.Auth.Name,
// 	}
// 	jsonStr, err := json.Marshal(map1)
// 	if err != nil {
// 		fmt.Printf("Error: %s", err.Error())
// 	}

// 	kafka.NewService("AJE").SendMessage("EXPORT", string(jsonStr))

// 	return response.SuccessResponse("SUKSES").Send(c)
// }

// Get Version
// @Summary Get Jpm Version
// @Description Get Jpm Version
// @Tags Jpm
// @Accept json
// @Produce json
// @Security BearerAuth
// @param request query dto.GetVersionRequest true "request query"
// @Success 200 {object} dto.GetVersionResponseDoc
// @Failure 400 {object} response.errorResponse
// @Failure 404 {object} response.errorResponse
// @Failure 500 {object} response.errorResponse
// @Router /jpm/get-version [get]
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

	if err := c.Validate(payload); err != nil {
		return response.ErrorBuilder(&response.ErrorConstant.Validation, err).Send(c)
	}

	result, err := h.service.GetVersion(cc, payload)
	if err != nil {
		return response.ErrorResponse(err).Send(c)
	}
	return response.SuccessResponse(result.Data).Send(c)
}
