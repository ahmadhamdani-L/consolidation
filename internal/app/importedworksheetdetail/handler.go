package importedworksheetdetail

import (
	"fmt"
	"mcash-finance-console-core/internal/abstraction"
	"mcash-finance-console-core/internal/dto"
	"mcash-finance-console-core/internal/factory"
	res "mcash-finance-console-core/pkg/util/response"

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
// @Summary Get ImportedWorksheet Detail
// @Description Get ImportedWorksheet Detail
// @Tags ImportedWorksheet Detail
// @Accept json
// @Produce json
// @Security BearerAuth
// @param request query dto.ImportedWorksheetDetailGetRequest true "request query"
// @Success 200 {object} dto.ImportedWorksheetDetailGetResponseDoc
// @Failure 400 {object} response.errorResponse
// @Failure 404 {object} response.errorResponse
// @Failure 500 {object} response.errorResponse
// @Router /imported-worksheet-detail [get]
func (h *handler) Get(c echo.Context) error {
	cc := c.(*abstraction.Context)

	payload := new(dto.ImportedWorksheetDetailGetRequest)
	if err := c.Bind(payload); err != nil {
		return res.ErrorBuilder(&res.ErrorConstant.BadRequest, err).Send(c)
	}
	if err = c.Validate(payload); err != nil {
		return res.ErrorBuilder(&res.ErrorConstant.Validation, err).Send(c)
	}

	result, err := h.service.Find(cc, payload)
	if err != nil {
		return res.ErrorResponse(err).Send(c)
	}

	return res.CustomSuccessBuilder(200, result.Datas, "Get datas success", &result.PaginationInfo).Send(c)
}

// Get By ID
// @Summary Get ImportedWorksheet Detail by id
// @Description Get ImportedWorksheet Detail by id
// @Tags ImportedWorksheet Detail
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "id path"
// @Success 200 {object} dto.ImportedWorksheetDetailGetByIDResponseDoc
// @Failure 400 {object} res.errorResponse
// @Failure 404 {object} res.errorResponse
// @Failure 500 {object} res.errorResponse
// @Router /imported-worksheet-detail/{id} [get]
func (h *handler) GetByID(c echo.Context) error {
	cc := c.(*abstraction.Context)

	payload := new(dto.ImportedWorksheetDetailGetByIDRequest)
	if err = c.Bind(payload); err != nil {
		return res.ErrorBuilder(&res.ErrorConstant.BadRequest, err).Send(c)
	}
	if err = c.Validate(payload); err != nil {
		response := res.ErrorBuilder(&res.ErrorConstant.Validation, err)
		return response.Send(c)
	}

	fmt.Printf("%+v", payload)

	result, err := h.service.FindByID(cc, payload)
	if err != nil {
		return res.ErrorResponse(err).Send(c)
	}

	return res.SuccessResponse(result).Send(c)
}
