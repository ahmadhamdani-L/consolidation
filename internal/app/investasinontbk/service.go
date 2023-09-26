package investasinontbk

import (
	"errors"
	"fmt"
	"log"
	"mcash-finance-console-core/internal/abstraction"
	"mcash-finance-console-core/internal/dto"
	"mcash-finance-console-core/internal/factory"
	"mcash-finance-console-core/internal/model"
	"mcash-finance-console-core/internal/repository"
	"mcash-finance-console-core/pkg/util/helper"
	"mcash-finance-console-core/pkg/util/response"
	"mcash-finance-console-core/pkg/util/trxmanager"
	"os"
	"strconv"
	"time"

	"github.com/xuri/excelize/v2"
	"gorm.io/gorm"
)

type service struct {
	Repository                      repository.InvestasiNonTbk
	InvestasiNonTbkDetailRepository repository.InvestasiNonTbkDetail
	FormatterRepository             repository.Formatter
	ParameterRepository             repository.Parameter
	Db                              *gorm.DB
}

type Service interface {
	Find(ctx *abstraction.Context, payload *dto.InvestasiNonTbkGetRequest) (*dto.InvestasiNonTbkGetResponse, error)
	FindByID(ctx *abstraction.Context, payload *dto.InvestasiNonTbkGetByIDRequest) (*dto.InvestasiNonTbkGetByIDResponse, error)
	Create(ctx *abstraction.Context, payload *dto.InvestasiNonTbkCreateRequest) (*dto.InvestasiNonTbkCreateResponse, error)
	Update(ctx *abstraction.Context, payload *dto.InvestasiNonTbkUpdateRequest) (*dto.InvestasiNonTbkUpdateResponse, error)
	Delete(ctx *abstraction.Context, payload *dto.InvestasiNonTbkDeleteRequest) (*dto.InvestasiNonTbkDeleteResponse, error)
	Export(ctx *abstraction.Context, payload *dto.InvestasiNonTbkExportRequest) (*dto.InvestasiNonTbkExportResponse, error)
	GetVersion(ctx *abstraction.Context, payload *dto.GetVersionRequest) (*dto.GetVersionResponse, error)
}

func NewService(f *factory.Factory) *service {
	repository := f.InvestasiNonTbkRepository
	investasiNonTbkDetailRepository := f.InvestasiNonTbkDetailRepository
	parameter := f.ParameterRepository
	formatter := f.FormatterRepository
	db := f.Db
	return &service{
		Repository:                      repository,
		InvestasiNonTbkDetailRepository: investasiNonTbkDetailRepository,
		ParameterRepository:             parameter,
		FormatterRepository:             formatter,
		Db:                              db,
	}
}

func (s *service) Find(ctx *abstraction.Context, payload *dto.InvestasiNonTbkGetRequest) (*dto.InvestasiNonTbkGetResponse, error) {
	data, info, err := s.Repository.Find(ctx, &payload.InvestasiNonTbkFilterModel, &payload.Pagination)
	if err != nil {
		return &dto.InvestasiNonTbkGetResponse{}, helper.ErrorHandler(err)
	}

	result := &dto.InvestasiNonTbkGetResponse{
		Datas:          *data,
		PaginationInfo: *info,
	}
	return result, nil
}

func (s *service) FindByID(ctx *abstraction.Context, payload *dto.InvestasiNonTbkGetByIDRequest) (*dto.InvestasiNonTbkGetByIDResponse, error) {
	data, err := s.Repository.FindByID(ctx, &payload.ID)
	if err != nil {
		return &dto.InvestasiNonTbkGetByIDResponse{}, helper.ErrorHandler(err)
	}
	allowed := helper.CompanyValidation(ctx.Auth.ID, data.CompanyID)
	if !allowed {
		return nil, response.ErrorBuilder(&response.ErrorConstant.UnauthorizedAccess, errors.New("Not Allowed"))
	}
	result := &dto.InvestasiNonTbkGetByIDResponse{
		InvestasiNonTbkEntityModel: *data,
	}
	return result, nil
}

func (s *service) Create(ctx *abstraction.Context, payload *dto.InvestasiNonTbkCreateRequest) (*dto.InvestasiNonTbkCreateResponse, error) {
	var data model.InvestasiNonTbkEntityModel

	allowed := helper.CompanyValidation(ctx.Auth.ID, payload.CompanyID)
	if !allowed {
		return nil, response.ErrorBuilder(&response.ErrorConstant.UnauthorizedAccess, errors.New("Not Allowed"))
	}

	if err := trxmanager.New(s.Db).WithTrx(ctx, func(ctx *abstraction.Context) error {
		data.Context = ctx
		data.InvestasiNonTbkEntity = payload.InvestasiNonTbkEntity
		data.InvestasiNonTbkEntity.Status = 1

		result, err := s.Repository.Create(ctx, &data)
		if err != nil {
			return helper.ErrorHandler(err)
		}
		data = *result
		return nil
	}); err != nil {
		return &dto.InvestasiNonTbkCreateResponse{}, err
	}

	result := &dto.InvestasiNonTbkCreateResponse{
		InvestasiNonTbkEntityModel: data,
	}
	return result, nil
}

func (s *service) Update(ctx *abstraction.Context, payload *dto.InvestasiNonTbkUpdateRequest) (*dto.InvestasiNonTbkUpdateResponse, error) {
	var data model.InvestasiNonTbkEntityModel
	if err := trxmanager.New(s.Db).WithTrx(ctx, func(ctx *abstraction.Context) error {
		existing, err := s.Repository.FindByID(ctx, &payload.ID)
		if err != nil {
			return helper.ErrorHandler(err)
		}

		if existing.Status != 1 {
			return response.ErrorBuilder(&response.ErrorConstant.DataValidated, errors.New("Cannot update data"))
		}

		allowed := helper.CompanyValidation(ctx.Auth.ID, existing.CompanyID)
		allowed2 := helper.CompanyValidation(ctx.Auth.ID, payload.CompanyID)
		if !allowed || !allowed2 {
			return response.ErrorBuilder(&response.ErrorConstant.UnauthorizedAccess, errors.New("Not Allowed"))
		}

		data.Context = ctx
		data.InvestasiNonTbkEntity = payload.InvestasiNonTbkEntity

		result, err := s.Repository.Update(ctx, &payload.ID, &data)
		if err != nil {
			return helper.ErrorHandler(err)
		}
		data = *result
		return nil
	}); err != nil {
		return &dto.InvestasiNonTbkUpdateResponse{}, err
	}
	result := &dto.InvestasiNonTbkUpdateResponse{
		InvestasiNonTbkEntityModel: data,
	}
	return result, nil
}

func (s *service) Delete(ctx *abstraction.Context, payload *dto.InvestasiNonTbkDeleteRequest) (*dto.InvestasiNonTbkDeleteResponse, error) {
	var data model.InvestasiNonTbkEntityModel
	if err := trxmanager.New(s.Db).WithTrx(ctx, func(ctx *abstraction.Context) error {
		existing, err := s.Repository.FindByID(ctx, &payload.ID)
		if err != nil {
			return helper.ErrorHandler(err)
		}

		allowed := helper.CompanyValidation(ctx.Auth.ID, existing.CompanyID)
		if !allowed {
			return response.ErrorBuilder(&response.ErrorConstant.UnauthorizedAccess, errors.New("Not Allowed"))
		}

		data.Context = ctx
		result, err := s.Repository.Delete(ctx, &payload.ID, &data)
		if err != nil {
			return helper.ErrorHandler(err)
		}
		data = *result
		return nil
	}); err != nil {
		return &dto.InvestasiNonTbkDeleteResponse{}, helper.ErrorHandler(err)
	}
	result := &dto.InvestasiNonTbkDeleteResponse{
		// InvestasiNonTbkEntityModel: data,
	}
	return result, nil
}

func (s *service) Export(ctx *abstraction.Context, payload *dto.InvestasiNonTbkExportRequest) (*dto.InvestasiNonTbkExportResponse, error) {
	f := excelize.NewFile()
	sheet := f.GetSheetName(f.GetActiveSheetIndex())

	arrStyleColWidth := []map[string]interface{}{
		{"COL": "A", "WIDTH": 8.21},
		{"COL": "B", "WIDTH": 3.50},
		{"COL": "C", "WIDTH": 13.71},
		{"COL": "D", "WIDTH": 13.71},
		{"COL": "E", "WIDTH": 13.71},
		{"COL": "F", "WIDTH": 13.71},
		{"COL": "G", "WIDTH": 13.71},
		{"COL": "H", "WIDTH": 13.71},
		{"COL": "I", "WIDTH": 13.71},
	}
	for _, v := range arrStyleColWidth {
		tmpColWidth := fmt.Sprintf("%f", v["WIDTH"])
		colWidth, _ := strconv.ParseFloat(tmpColWidth, 64)
		err := f.SetColWidth(sheet, fmt.Sprintf("%s", v["COL"]), fmt.Sprintf("%s", v["COL"]), colWidth)
		if err != nil {
			return nil, response.ErrorBuilder(&response.ErrorConstant.InternalServerError, errors.New("Internal Server Error"))
		}
	}

	styleHeader, err := f.NewStyle(&excelize.Style{
		Border: []excelize.Border{
			{Type: "top", Color: "000000", Style: 1},
			{Type: "bottom", Color: "000000", Style: 1},
			{Type: "left", Color: "000000", Style: 1},
			{Type: "right", Color: "000000", Style: 1},
		},
		Font: &excelize.Font{
			Bold:   true,
			Family: "Arial",
			Size:   10,
		},
		Fill: excelize.Fill{
			Type:    "pattern",
			Pattern: 1,
			Color:   []string{"#fac090"},
		},
		Alignment: &excelize.Alignment{
			WrapText:   true,
			Horizontal: "center",
			Vertical:   "center",
		},
	})
	if err != nil {
		return nil, response.ErrorBuilder(&response.ErrorConstant.InternalServerError, errors.New("Internal Server Error"))
	}
	numberFormat := "#,##0.000"
	styleCurrencyPercentage, err := f.NewStyle(&excelize.Style{
		Border: []excelize.Border{
			{Type: "top", Color: "000000", Style: 1},
			{Type: "bottom", Color: "000000", Style: 1},
			{Type: "left", Color: "000000", Style: 1},
			{Type: "right", Color: "000000", Style: 1},
		},
		CustomNumFmt: &numberFormat,
		Font: &excelize.Font{
			Family: "Arial",
			Size:   10,
		},
	})
	if err != nil {
		return nil, response.ErrorBuilder(&response.ErrorConstant.InternalServerError, errors.New("Internal Server Error"))
	}
	numberFormat = "#,##"
	styleCurrency, err := f.NewStyle(&excelize.Style{
		Border: []excelize.Border{
			{Type: "top", Color: "000000", Style: 1},
			{Type: "bottom", Color: "000000", Style: 1},
			{Type: "left", Color: "000000", Style: 1},
			{Type: "right", Color: "000000", Style: 1},
		},
		CustomNumFmt: &numberFormat,
		Font: &excelize.Font{
			Family: "Arial",
			Size:   10,
		},
	})
	if err != nil {
		return nil, response.ErrorBuilder(&response.ErrorConstant.InternalServerError, errors.New("Internal Server Error"))
	}

	styleLabel, err := f.NewStyle(&excelize.Style{
		Border: []excelize.Border{
			{Type: "top", Color: "000000", Style: 1},
			{Type: "bottom", Color: "000000", Style: 1},
			{Type: "left", Color: "000000", Style: 1},
			{Type: "right", Color: "000000", Style: 1},
		},
		Font: &excelize.Font{
			Family: "Arial",
			Size:   10,
		},
	})
	if err != nil {
		return nil, response.ErrorBuilder(&response.ErrorConstant.InternalServerError, errors.New("Internal Server Error"))
	}

	// criteria := model.InvestasiNonTbkFilterModel{}
	// criteria.CompanyID = &payload.CompanyID
	// version := payload.Versions
	// criteria.Versions = &version
	// criteria.Period = &payload.Period

	investasinontbk, err := s.Repository.FindByID(ctx, &payload.InvestasiNonTbkID)
	if err != nil {
		return nil, helper.ErrorHandler(err)
	}

	allowed := helper.CompanyValidation(ctx.Auth.ID, investasinontbk.CompanyID)
	if !allowed {
		return nil, response.ErrorBuilder(&response.ErrorConstant.UnauthorizedAccess, errors.New("Not Allowed"))
	}

	datePeriod, err := time.Parse(time.RFC3339, investasinontbk.Period)
	if err != nil {
		return nil, err
	}

	row := 5

	f.SetCellValue(sheet, "B2", "Detail investasi anak usaha Non TBK")

	f.SetCellStyle(sheet, "B4", "J4", styleHeader)
	f.SetCellValue(sheet, "B4", "No")
	f.SetCellValue(sheet, "C4", "Code")
	f.SetCellValue(sheet, "D4", "Lembar saham dimiliki")
	f.SetCellValue(sheet, "E4", "Total lembar saham")
	f.SetCellValue(sheet, "F4", "% Ownership")
	f.SetCellValue(sheet, "G4", "Harga Par")
	f.SetCellValue(sheet, "H4", "Total harga Par")
	f.SetCellValue(sheet, "I4", "Harga beli")
	f.SetCellValue(sheet, "J4", "Total Harga beli")

	criteriaDetail := model.InvestasiNonTbkDetailFilterModel{}

	criteriaDetail.InvestasiNonTbkID = &investasinontbk.ID

	paginationDetail := abstraction.Pagination{}
	pagesize := 10000
	// sortBy := "id"
	// sort := "asc"
	paginationDetail.PageSize = &pagesize
	// paginationDetail.SortBy = &sortBy
	// paginationDetail.Sort = &sort

	detail, _, err := s.InvestasiNonTbkDetailRepository.Find(ctx, &criteriaDetail, &paginationDetail)
	if err != nil {
		return nil, helper.ErrorHandler(err)
	}

	for _, v := range *detail {
		f.SetCellStyle(sheet, fmt.Sprintf("B%d", row), fmt.Sprintf("C%d", row), styleLabel)
		f.SetCellStyle(sheet, fmt.Sprintf("D%d", row), fmt.Sprintf("E%d", row), styleCurrency)
		f.SetCellStyle(sheet, fmt.Sprintf("F%d", row), fmt.Sprintf("F%d", row), styleCurrencyPercentage)
		f.SetCellStyle(sheet, fmt.Sprintf("G%d", row), fmt.Sprintf("J%d", row), styleCurrency)

		f.SetCellValue(sheet, fmt.Sprintf("B%d", row), v.Code)

		f.SetCellValue(sheet, fmt.Sprintf("B%d", row), v.SortID)
		f.SetCellValue(sheet, fmt.Sprintf("C%d", row), v.Code)
		f.SetCellValue(sheet, fmt.Sprintf("D%d", row), *v.LbrSahamOwnership)
		f.SetCellValue(sheet, fmt.Sprintf("E%d", row), *v.TotalLbrSaham)
		f.SetCellFormula(sheet, fmt.Sprintf("F%d", row), fmt.Sprintf("=D%d/E%d", row, row))
		f.SetCellValue(sheet, fmt.Sprintf("G%d", row), *v.HargaPar)
		f.SetCellFormula(sheet, fmt.Sprintf("H%d", row), fmt.Sprintf("=D%d*G%d", row, row))
		f.SetCellValue(sheet, fmt.Sprintf("I%d", row), *v.HargaBeli)
		f.SetCellFormula(sheet, fmt.Sprintf("J%d", row), fmt.Sprintf("=D%d*I%d", row, row))
		row++
	}

	tmpFolder := fmt.Sprintf("assets/%d", ctx.Auth.ID)
	_, err = os.Stat(tmpFolder)
	if err != nil {
		if os.IsNotExist(err) {
			if err := os.MkdirAll(tmpFolder, os.ModePerm); err != nil {
				log.Fatal(err)
			}
		} else {
			return nil, response.ErrorBuilder(&response.ErrorConstant.InternalServerError, errors.New("Internal Server Error"))
		}
	}
	period := datePeriod.Format("2006-01-02")
	fileName := fmt.Sprintf("InvestasiNonTbk_%s.xlsx", period)
	fileLoc := fmt.Sprintf("assets/%d/%s", ctx.Auth.ID, fileName)
	err = f.SaveAs(fileLoc)
	if err != nil {
		return nil, err
	}

	result := &dto.InvestasiNonTbkExportResponse{
		FileName: fileName,
		Path:     fileLoc,
	}
	return result, nil
}

func (s *service) GetVersion(ctx *abstraction.Context, payload *dto.GetVersionRequest) (*dto.GetVersionResponse, error) {
	filter := model.InvestasiNonTbkFilterModel{
		CompanyCustomFilter: model.CompanyCustomFilter{
			CompanyID:          payload.CompanyID,
			ArrCompanyID:       payload.ArrCompanyID,
			ArrCompanyString:   payload.ArrCompanyString,
			ArrCompanyOperator: payload.ArrCompanyOperator,
		},
	}
	filter.Period = payload.Period

	data, err := s.Repository.GetVersion(ctx, &filter)
	if err != nil {
		return &dto.GetVersionResponse{}, helper.ErrorHandler(err)
	}
	result := &dto.GetVersionResponse{
		Data: *data,
	}
	return result, nil
}
