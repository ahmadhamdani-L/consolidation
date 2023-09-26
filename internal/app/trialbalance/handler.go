package trialbalance

import (
	"encoding/json"
	"errors"
	"fmt"
	"mcash-finance-console-core/internal/abstraction"
	"mcash-finance-console-core/internal/dto"
	"mcash-finance-console-core/internal/factory"
	"mcash-finance-console-core/internal/model"
	modelhelper "mcash-finance-console-core/internal/model/helper"
	"mcash-finance-console-core/pkg/kafka"
	"mcash-finance-console-core/pkg/util/helper"
	"mcash-finance-console-core/pkg/util/response"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/xuri/excelize/v2"
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
// @Summary Get Trial Balance
// @Description Get Trial Balance
// @Tags Trial Balance
// @Accept json
// @Produce json
// @Security BearerAuth
// @param request query dto.TrialBalanceGetRequest true "request query"
// @Success 200 {object} dto.TrialBalanceGetResponseDoc
// @Failure 400 {object} response.errorResponse
// @Failure 404 {object} response.errorResponse
// @Failure 500 {object} response.errorResponse
// @Router /trial-balance [get]
func (h *handler) Get(c echo.Context) error {
	cc := c.(*abstraction.Context)
	payload := new(dto.TrialBalanceGetRequest)

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
// @Summary Get Trial Balance By ID
// @Description Get Trial Balance By ID
// @Tags Trial Balance
// @Accept json
// @Produce json
// @Security BearerAuth
// @param id path int true "id path"
// @Success 200 {object} dto.TrialBalanceGetByIDResponseDoc
// @Failure 400 {object} response.errorResponse
// @Failure 404 {object} response.errorResponse
// @Failure 500 {object} response.errorResponse
// @Router /trial-balance/{id} [get]
func (h *handler) GetByID(c echo.Context) error {
	cc := c.(*abstraction.Context)
	payload := new(dto.TrialBalanceGetByIDRequest)

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
	payload := new(dto.TrialBalanceCreateRequest)

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
// @Summary Update Trial Balance
// @Description Update Trial Balance
// @Tags Trial Balance
// @Accept json
// @Produce json
// @Security BearerAuth
// @param id path int true "id path"
// @param request body dto.TrialBalanceUpdateRequest true "request body"
// @Success 200 {object} dto.TrialBalanceUpdateResponseDoc
// @Failure 400 {object} response.errorResponse
// @Failure 404 {object} response.errorResponse
// @Failure 500 {object} response.errorResponse
// @Router /trial-balance/{id} [patch]
func (h *handler) Update(c echo.Context) error {
	cc := c.(*abstraction.Context)
	payload := new(dto.TrialBalanceUpdateRequest)

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
// @Summary Delete Trial Balance
// @Description Delete Trial Balance
// @Tags Trial Balance
// @Accept json
// @Produce json
// @Security BearerAuth
// @param id path int true "id path"
// @Success 200 {object} dto.TrialBalanceDeleteResponseDoc
// @Failure 400 {object} response.errorResponse
// @Failure 404 {object} response.errorResponse
// @Failure 500 {object} response.errorResponse
// @Router /trial-balance/{id} [delete]
func (h *handler) Delete(c echo.Context) error {
	cc := c.(*abstraction.Context)
	payload := new(dto.TrialBalanceDeleteRequest)

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
// @Summary Export Trial Balance
// @Description Export Trial Balance
// @Tags Trial Balance
// @Accept json
// @Produce json
// @Security BearerAuth
// @param trial_balance_id query int true "trial_balance_id query"
// @Success 200 {object} dto.TrialBalanceExportResponseDoc
// @Failure 400 {object} response.errorResponse
// @Failure 404 {object} response.errorResponse
// @Failure 500 {object} response.errorResponse
// @Router /trial-balance/export [get]
func (h *handler) Export(c echo.Context) error {
	cc := c.(*abstraction.Context)
	payload := new(dto.TrialBalanceExportRequest)

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

func (h *handler) Import(c echo.Context) error {
	cc := c.(*abstraction.Context)
	payload := new(dto.TrialBalanceImportRequest)

	file, err := c.FormFile("trial")
	if err != nil {
		return response.ErrorBuilder(&response.ErrorConstant.UnprocessableEntity, err).Send(c)
	}

	src, err := file.Open()
	if err != nil {
		return nil
	}
	defer src.Close()

	f, err := excelize.OpenReader(src)
	if err != nil {
		return nil
	}

	defer func() {
		if err := f.Close(); err != nil {
			fmt.Println(err)
		}
	}()

	sheet := f.GetSheetName(f.GetActiveSheetIndex())
	rows, err := f.GetRows(sheet)
	if err != nil {
		return response.ErrorBuilder(&response.ErrorConstant.BadRequest, err).Send(c)
	}

	head1, err := f.GetCellValue(sheet, "B6")
	if err != nil {
		return response.ErrorBuilder(&response.ErrorConstant.BadRequest, err).Send(c)
	}

	head2, err := f.GetCellValue(sheet, "C6")
	if err != nil {
		return response.ErrorBuilder(&response.ErrorConstant.BadRequest, err).Send(c)
	}

	head3, err := f.GetCellValue(sheet, "F6")
	if err != nil {
		return response.ErrorBuilder(&response.ErrorConstant.BadRequest, err).Send(c)
	}

	head4, err := f.GetCellValue(sheet, "G7")
	if err != nil {
		return response.ErrorBuilder(&response.ErrorConstant.BadRequest, err).Send(c)
	}

	head5, err := f.GetCellValue(sheet, "H6")
	if err != nil {
		return response.ErrorBuilder(&response.ErrorConstant.BadRequest, err).Send(c)
	}

	if strings.ToLower(head1) != "no akun" || strings.ToLower(head2) != "keterangan" || strings.ToLower(head3) != "wp reff" || strings.ToLower(head4) != "unaudited" || strings.ToLower(head5) != "adjustment journal entry" {
		return response.ErrorBuilder(&response.ErrorConstant.BadRequest, err).Send(c)
	}
	rows = rows[8:][:]

	datas := []model.TrialBalanceDetailEntity{}
	for _, row := range rows {
		if len(row) > 1 && len(row[1]) == 9 {
			if len(row) < 6 {
				continue
			}
			nominal, err := strconv.ParseFloat(row[6], 64)
			if err != nil {
				continue
			}
			coa := row[1]
			nominalBeforeAje := nominal
			data := model.TrialBalanceDetailEntity{
				Code:            coa,
				AmountBeforeAje: &nominalBeforeAje,
				Description:     &row[4],
			}
			datas = append(datas, data)
		}
	}

	payload = &dto.TrialBalanceImportRequest{
		UserID: cc.Auth.ID,
	}
	fmt.Println(payload)
	/* result, err := h.service.Import(cc, payload, &datas)
	if err != nil {
		return response.ErrorResponse(err).Send(c)
	} */

	return response.SuccessResponse("result").Send(c)
}

func (h *handler) ExportAsync(c echo.Context) error {
	cc := c.(*abstraction.Context)
	payload := new(dto.TrialBalanceExportAsyncRequest)
	fmt.Println(cc)
	if err := c.Bind(payload); err != nil {
		return response.ErrorBuilder(&response.ErrorConstant.BadRequest, err).Send(c)
	}
	if payload.Period != "" {
		_, err := time.Parse("2006-01-02", payload.Period)
		if err != nil {
			return response.ErrorBuilder(&response.ErrorConstant.Validation, err).Send(c)
		}
	}
	if err := c.Validate(payload); err != nil {
		return response.ErrorBuilder(&response.ErrorConstant.Validation, err).Send(c)
	}

	waktu := time.Now()
	msg := kafka.JsonData{
		UserID:    1,
		CompanyID: 1,
		Timestamp: &waktu,
		// Name:      cc.Auth.Name,
		Filter: struct {
			Period   string
			Versions int
			Request  string
		}{payload.Period, payload.Versions, ""},
	}
	jsonStr, err := json.Marshal(msg)
	if err != nil {
		fmt.Printf("Error: %s", err.Error())
	}

	go kafka.NewService("EXPORT").SendMessage("EXPORT", string(jsonStr))

	return response.SuccessResponse("SUKSES").Send(c)
}

// Get Version
// @Summary Get Trial Balance Version
// @Description Get Trial Balance Version
// @Tags Trial Balance
// @Accept json
// @Produce json
// @Security BearerAuth
// @param request query dto.GetVersionRequest true "request query"
// @Success 200 {object} dto.GetVersionResponseDoc
// @Failure 400 {object} response.errorResponse
// @Failure 404 {object} response.errorResponse
// @Failure 500 {object} response.errorResponse
// @Router /trial-balance/get-version [get]
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
