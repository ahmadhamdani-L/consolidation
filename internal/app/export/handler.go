package export

import (
	"errors"
	"fmt"
	"mcash-finance-console-core/internal/abstraction"
	"mcash-finance-console-core/internal/dto"
	"mcash-finance-console-core/internal/factory"
	"mcash-finance-console-core/pkg/util/response"
	"net/http"
	"os"
	"path"
	"strings"

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

// Export Bulk
// @Summary Request Export
// @Description Export Bulk
// @Tags Export Bulk
// @Accept json
// @Produce json
// @Security BearerAuth
// @param request query dto.ExportRequest true "request query"
// @Success 200 {object} dto.ExportResponseDoc
// @Failure 400 {object} response.errorResponse
// @Failure 404 {object} response.errorResponse
// @Failure 500 {object} response.errorResponse
// @Router /export/request [get]
func (h *handler) Request(c echo.Context) error {
	cc := c.(*abstraction.Context)
	payload := new(dto.ExportRequest)
	if err := c.Bind(payload); err != nil {
		return response.ErrorBuilder(&response.ErrorConstant.BadRequest, err).Send(c)
	}

	if err := c.Validate(payload); err != nil {
		return response.ErrorBuilder(&response.ErrorConstant.Validation, err).Send(c)
	}

	if payload.Request != "" {
		request := strings.Split(strings.ToUpper(payload.Request), ",")
		for _, v := range request {
			switch v {
			case "TB", "AUP", "MP", "MRUA", "PPB", "MFA", "MIA", "MDTA", "INT", "IT", "EB":
				// case "AUP":
				// case "MP":
				// case "MRUA":
				// case "AJE":
				// case "JCTE":
				// case "JELIM":
				// case "PPB":
				// case "MFA":
				// case "MIA":
				// case "MDTA":
				// case "INT":
				// case "IT":
				// case "JPM":
				continue
			default:
				return response.ErrorBuilder(&response.ErrorConstant.BadRequest, errors.New(fmt.Sprintf("%s not a valid file request", v))).Send(c)
			}
		}
	}

	result, err := h.service.RequestExport(cc, payload)
	if err != nil {
		return response.ErrorResponse(err).Send(c)
	}
	return response.SuccessResponse(result).Send(c)
}

// Get Export Request
// @Summary Export Bulk
// @Description Export Bulk
// @Tags Export Bulk
// @Accept json
// @Produce application/zip
// @Security BearerAuth
// @param request query dto.GetExportRequest true "request query"
// @Success 200 {object} dto.ExportResponseDoc
// @Failure 400 {object} response.errorResponse
// @Failure 404 {object} response.errorResponse
// @Failure 500 {object} response.errorResponse
// @Router /export/export_modul/{notification_id} [get]
func (h *handler) GetRequest(c echo.Context) error {
	cc := c.(*abstraction.Context)
	payload := new(dto.GetExportRequest)
	if err := c.Bind(payload); err != nil {
		return response.ErrorBuilder(&response.ErrorConstant.BadRequest, err).Send(c)
	}

	if err := c.Validate(payload); err != nil {
		return response.ErrorBuilder(&response.ErrorConstant.Validation, err).Send(c)
	}

	result, err := h.service.GetExport(cc, payload)
	if err != nil {
		return response.ErrorResponse(err).Send(c)
	}
	splitPath := strings.Split(*result, "/")
	f, err := os.Open(*result)
	if err != nil {
		return err
	}
	c.Response().Header().Set(echo.HeaderContentDisposition, fmt.Sprintf("attachment; filename=%s", splitPath[len(splitPath)-1]))
	return c.Stream(http.StatusOK, "application/zip", f)
}

// Request Export Consol
// @Summary Request Export Consol
// @Description Request Export Consol
// @Tags Export Bulk
// @Accept json
// @Produce application/json
// @Security BearerAuth
// @param request query dto.ExportConsolRequest true "request query"
// @Success 200 {object} dto.ExportResponseDoc
// @Failure 400 {object} response.errorResponse
// @Failure 404 {object} response.errorResponse
// @Failure 500 {object} response.errorResponse
// @Router /export/request_consol [get]
func (h *handler) RequestExportConsol(c echo.Context) error {
	cc := c.(*abstraction.Context)
	payload := new(dto.ExportConsolRequest)
	if err := c.Bind(payload); err != nil {
		return response.ErrorBuilder(&response.ErrorConstant.BadRequest, err).Send(c)
	}

	if err := c.Validate(payload); err != nil {
		return response.ErrorBuilder(&response.ErrorConstant.Validation, err).Send(c)
	}

	result, err := h.service.ExportConsol(cc, payload)
	if err != nil {
		return response.ErrorResponse(err).Send(c)
	}

	return response.SuccessResponse(result).Send(c)
}

// Get Export Consol
// @Summary Export Bulk
// @Description Export Bulk
// @Tags Export Bulk
// @Accept json
// @Produce application/vnd.openxmlformats-officedocument.spreadsheetml.sheet
// @Security BearerAuth
// @param request query dto.GetExportRequest true "request query"
// @Success 200 {object} dto.ExportResponseDoc
// @Failure 400 {object} response.errorResponse
// @Failure 404 {object} response.errorResponse
// @Failure 500 {object} response.errorResponse
// @Router /export/export_consol/{notification_id} [get]
func (h *handler) GetRequestConsol(c echo.Context) error {
	cc := c.(*abstraction.Context)
	payload := new(dto.GetExportRequest)
	if err := c.Bind(payload); err != nil {
		return response.ErrorBuilder(&response.ErrorConstant.BadRequest, err).Send(c)
	}

	if err := c.Validate(payload); err != nil {
		return response.ErrorBuilder(&response.ErrorConstant.Validation, err).Send(c)
	}

	result, err := h.service.GetExport(cc, payload)
	if err != nil {
		return response.ErrorResponse(err).Send(c)
	}
	splitPath := strings.Split(*result, "/")
	f, err := os.Open(*result)
	if err != nil {
		return err
	}

	c.Response().Header().Set(echo.HeaderContentDisposition, fmt.Sprintf("attachment; filename=%s", splitPath[len(splitPath)-1]))
	return c.Stream(http.StatusOK, "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet", f)
}

// Get Export All
// @Summary Export All
// @Description Export All
// @Tags Export Bulk
// @Accept json
// @Produce application/vnd.openxmlformats-officedocument.spreadsheetml.sheet
// @Security BearerAuth
// @param request query dto.GetExportRequest true "request query"
// @Success 200 {object} dto.ExportResponseDoc
// @Failure 400 {object} response.errorResponse
// @Failure 404 {object} response.errorResponse
// @Failure 500 {object} response.errorResponse
// @Router /export [get]
func (h *handler) ExportAll(c echo.Context) error {
	cc := c.(*abstraction.Context)
	payload := new(dto.ExportRequest)
	if err := c.Bind(payload); err != nil {
		return response.ErrorBuilder(&response.ErrorConstant.BadRequest, err).Send(c)
	}

	if err := c.Validate(payload); err != nil {
		return response.ErrorBuilder(&response.ErrorConstant.Validation, err).Send(c)
	}

	result, err := h.service.ExportModul(cc, payload)
	if err != nil {
		return response.ErrorResponse(err).Send(c)
	}
	f, err := os.Open(*result)
	if err != nil {
		return err
	}

	c.Response().Header().Set(echo.HeaderContentDisposition, fmt.Sprintf("attachment; filename=%s", path.Base(*result)))
	return c.Stream(http.StatusOK, "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet", f)
}
