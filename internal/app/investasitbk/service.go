package investasitbk

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
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/xuri/excelize/v2"
	"gorm.io/gorm"
)

type service struct {
	Repository                   repository.InvestasiTbk
	InvestasiTbkDetailRepository repository.InvestasiTbkDetail
	FormatterRepository          repository.Formatter
	ParameterRepository          repository.Parameter
	FormatterBridgesRepository   repository.FormatterBridges
	Db                           *gorm.DB
}

type Service interface {
	Find(ctx *abstraction.Context, payload *dto.InvestasiTbkGetRequest) (*dto.InvestasiTbkGetResponse, error)
	FindByID(ctx *abstraction.Context, payload *dto.InvestasiTbkGetByIDRequest) (*dto.InvestasiTbkGetByIDResponse, error)
	Create(ctx *abstraction.Context, payload *dto.InvestasiTbkCreateRequest) (*dto.InvestasiTbkCreateResponse, error)
	Update(ctx *abstraction.Context, payload *dto.InvestasiTbkUpdateRequest) (*dto.InvestasiTbkUpdateResponse, error)
	Delete(ctx *abstraction.Context, payload *dto.InvestasiTbkDeleteRequest) (*dto.InvestasiTbkDeleteResponse, error)
	Import(ctx *abstraction.Context, payload *dto.InvestasiTbkImportRequest) (*dto.InvestasiTbkImportResponse, error)
	Export(ctx *abstraction.Context, payload *dto.InvestasiTbkExportRequest) (*dto.InvestasiTbkExportResponse, error)
	GetVersion(ctx *abstraction.Context, payload *dto.GetVersionRequest) (*dto.GetVersionResponse, error)
}

func NewService(f *factory.Factory) *service {
	repository := f.InvestasiTbkRepository
	investasiTbkDetailRepository := f.InvestasiTbkDetailRepository
	formattter := f.FormatterRepository
	parameter := f.ParameterRepository
	formatterBridgesRepo := f.FormatterBridgesRepository
	db := f.Db
	return &service{
		Repository:                   repository,
		InvestasiTbkDetailRepository: investasiTbkDetailRepository,
		FormatterRepository:          formattter,
		ParameterRepository:          parameter,
		FormatterBridgesRepository:   formatterBridgesRepo,
		Db:                           db,
	}
}

func (s *service) Find(ctx *abstraction.Context, payload *dto.InvestasiTbkGetRequest) (*dto.InvestasiTbkGetResponse, error) {
	data, info, err := s.Repository.Find(ctx, &payload.InvestasiTbkFilterModel, &payload.Pagination)
	if err != nil {
		return &dto.InvestasiTbkGetResponse{}, response.ErrorBuilder(&response.ErrorConstant.InternalServerError, err)
	}

	result := &dto.InvestasiTbkGetResponse{
		Datas:          *data,
		PaginationInfo: *info,
	}
	return result, nil
}

func (s *service) FindByID(ctx *abstraction.Context, payload *dto.InvestasiTbkGetByIDRequest) (*dto.InvestasiTbkGetByIDResponse, error) {
	data, err := s.Repository.FindByID(ctx, &payload.ID)
	if err != nil {
		return &dto.InvestasiTbkGetByIDResponse{}, response.ErrorBuilder(&response.ErrorConstant.InternalServerError, err)
	}
	allowed := helper.CompanyValidation(ctx.Auth.ID, data.CompanyID)
	if !allowed {
		return nil, response.ErrorBuilder(&response.ErrorConstant.UnauthorizedAccess, errors.New("Not Allowed"))
	}
	result := &dto.InvestasiTbkGetByIDResponse{
		InvestasiTbkEntityModel: *data,
	}
	return result, nil
}

func (s *service) Create(ctx *abstraction.Context, payload *dto.InvestasiTbkCreateRequest) (*dto.InvestasiTbkCreateResponse, error) {
	var data model.InvestasiTbkEntityModel

	if err := trxmanager.New(s.Db).WithTrx(ctx, func(ctx *abstraction.Context) error {
		data.Context = ctx
		data.InvestasiTbkEntity = payload.InvestasiTbkEntity
		data.InvestasiTbkEntity.Status = 1

		result, err := s.Repository.Create(ctx, &data)
		if err != nil {
			return response.ErrorBuilder(&response.ErrorConstant.UnprocessableEntity, err)
		}
		data = *result
		return nil
	}); err != nil {
		return &dto.InvestasiTbkCreateResponse{}, err
	}

	result := &dto.InvestasiTbkCreateResponse{
		InvestasiTbkEntityModel: data,
	}
	return result, nil
}

func (s *service) Update(ctx *abstraction.Context, payload *dto.InvestasiTbkUpdateRequest) (*dto.InvestasiTbkUpdateResponse, error) {
	var data model.InvestasiTbkEntityModel
	if err := trxmanager.New(s.Db).WithTrx(ctx, func(ctx *abstraction.Context) error {
		existing, err := s.Repository.FindByID(ctx, &payload.ID)
		if err != nil {
			return helper.ErrorHandler(err)
		}

		if existing.Status != 1 {
			return response.ErrorBuilder(&response.ErrorConstant.DataValidated, errors.New("Cannot Update Data"))
		}

		allowed := helper.CompanyValidation(ctx.Auth.ID, existing.CompanyID)
		allowed2 := helper.CompanyValidation(ctx.Auth.ID, payload.CompanyID)
		if !allowed || !allowed2 {
			return response.ErrorBuilder(&response.ErrorConstant.UnauthorizedAccess, errors.New("Not Allowed"))
		}

		data.Context = ctx
		data.InvestasiTbkEntity = payload.InvestasiTbkEntity

		result, err := s.Repository.Update(ctx, &payload.ID, &data)
		if err != nil {
			return helper.ErrorHandler(err)
		}
		data = *result
		return nil
	}); err != nil {
		return &dto.InvestasiTbkUpdateResponse{}, err
	}
	result := &dto.InvestasiTbkUpdateResponse{
		InvestasiTbkEntityModel: data,
	}
	return result, nil
}

func (s *service) Delete(ctx *abstraction.Context, payload *dto.InvestasiTbkDeleteRequest) (*dto.InvestasiTbkDeleteResponse, error) {
	var data model.InvestasiTbkEntityModel
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
		return &dto.InvestasiTbkDeleteResponse{}, response.ErrorBuilder(&response.ErrorConstant.InternalServerError, err)
	}
	result := &dto.InvestasiTbkDeleteResponse{
		// InvestasiTbkEntityModel: data,
	}
	return result, nil
}

func (s *service) Export(ctx *abstraction.Context, payload *dto.InvestasiTbkExportRequest) (*dto.InvestasiTbkExportResponse, error) {
	f := excelize.NewFile()
	sheet := f.GetSheetName(f.GetActiveSheetIndex())

	arrStyleColWidth := []map[string]interface{}{
		{"COL": "A", "WIDTH": 8.21},
		{"COL": "B", "WIDTH": 4.10},
		{"COL": "C", "WIDTH": 33.39},
		{"COL": "D", "WIDTH": 11.78},
		{"COL": "E", "WIDTH": 12.31},
		{"COL": "F", "WIDTH": 15.71},
		{"COL": "G", "WIDTH": 12.78},
		{"COL": "H", "WIDTH": 15.35},
		{"COL": "I", "WIDTH": 15.35},
		{"COL": "J", "WIDTH": 15.71},
		{"COL": "K", "WIDTH": 14.10},
	}
	for _, v := range arrStyleColWidth {
		tmpColWidth := fmt.Sprintf("%f", v["WIDTH"])
		colWidth, _ := strconv.ParseFloat(tmpColWidth, 64)
		err := f.SetColWidth(sheet, fmt.Sprintf("%s", v["COL"]), fmt.Sprintf("%s", v["COL"]), colWidth)
		if err != nil {
			return nil, err
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
			Color:   []string{"#fcd5b4"},
		},
		Alignment: &excelize.Alignment{
			WrapText:   true,
			Horizontal: "center",
			Vertical:   "center",
		},
	})
	if err != nil {
		return nil, err
	}
	numberFormat := "#,##"
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
		return nil, err
	}

	styleCurrencyTotal, err := f.NewStyle(&excelize.Style{
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
			Color:   []string{"#99ff33"},
		},
		CustomNumFmt: &numberFormat,
	})
	if err != nil {
		return nil, err
	}

	styleLabelTotal, err := f.NewStyle(&excelize.Style{
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
		Alignment: &excelize.Alignment{
			Horizontal: "center",
			Vertical:   "center",
		},
		Fill: excelize.Fill{
			Type:    "pattern",
			Pattern: 1,
			Color:   []string{"#99ff33"},
		},
	})
	if err != nil {
		return nil, err
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
		return nil, err
	}

	investasitbk, err := s.Repository.FindByID(ctx, &payload.InvestasiTbkID)
	if err != nil {
		return nil, err
	}

	if investasitbk.ID == 0 {
		return nil, errors.New("Data Not Found")
	}

	allowed := helper.CompanyValidation(ctx.Auth.ID, investasitbk.CompanyID)
	if !allowed {
		return nil, response.ErrorBuilder(&response.ErrorConstant.UnauthorizedAccess, errors.New("Not Allowed"))
	}

	datePeriod, err := time.Parse(time.RFC3339, investasitbk.Period)
	if err != nil {
		return nil, err
	}

	f.SetCellStyle(sheet, "B4", "K6", styleHeader)
	f.SetCellValue(sheet, "B4", "No")
	f.SetCellValue(sheet, "C4", "Stock")
	f.SetCellValue(sheet, "D4", "Ending Share")
	f.SetCellValue(sheet, "E4", "AVG Price")
	f.SetCellValue(sheet, "F4", "Amount (Cost)")
	f.SetCellValue(sheet, "G4", fmt.Sprintf("Closing Price (%s)", datePeriod.Format("02.01.06")))
	f.SetCellValue(sheet, "H4", "Amount (FV)")
	f.SetCellValue(sheet, "I4", "Unrealized Gain(oss)")
	f.SetCellValue(sheet, "J4", "Realized Gain(loss)")
	f.SetCellValue(sheet, "K4", "Fee")

	formatterCode := []string{"INVESTASI-TBK"}
	formatterTitle := []string{"Summary Investasi Tbk"}
	row, rowStart := 5, 5

	f.SetCellValue(sheet, "B2", formatterTitle[0])

	for _, formatter := range formatterCode {

		var criteria dto.FormatterGetRequest
		criteria.FormatterFilterModel.FormatterFor = &formatter

		data, err := s.FormatterRepository.FindWithDetail(ctx, &criteria.FormatterFilterModel)
		if err != nil {
			return &dto.InvestasiTbkExportResponse{}, response.ErrorBuilder(&response.ErrorConstant.InternalServerError, err)
		}

		tmpStr := "INVESTASI-TBK"
		criteriaBridge := model.FormatterBridgesFilterModel{}
		criteriaBridge.FormatterBridgesFilter.Source = &tmpStr
		criteriaBridge.FormatterBridgesFilter.FormatterID = &data.ID
		criteriaBridge.FormatterBridgesFilter.TrxRefID = &investasitbk.ID

		bridges, err := s.FormatterBridgesRepository.FindWithCriteria(ctx, &criteriaBridge)
		if err != nil {
			return nil, err
		}
		rowCode := make(map[string]int)
		partRowStart := row
		for _, v := range data.FormatterDetail {
			rowCode[v.Code] = row
			f.SetCellStyle(sheet, fmt.Sprintf("B%d", row), fmt.Sprintf("C%d", row), styleLabel)
			f.SetCellStyle(sheet, fmt.Sprintf("D%d", row), fmt.Sprintf("K%d", row), styleCurrency)

			if strings.Contains(strings.ToLower(v.Code), "blank") {
				row++
				continue
			}

			f.SetCellValue(sheet, fmt.Sprintf("B%d", row), v.Description)

			if v.IsTotal != nil && *v.IsTotal == true {
				f.SetCellStyle(sheet, fmt.Sprintf("B%d", row), fmt.Sprintf("C%d", row), styleLabelTotal)
				f.SetCellStyle(sheet, fmt.Sprintf("D%d", row), fmt.Sprintf("K%d", row), styleCurrencyTotal)
				if v.FxSummary == "" {
					row++
					continue
				}
				for chr := 'D'; chr <= 'K'; chr++ {
					formula := v.FxSummary
					reg := regexp.MustCompile(`[0-9]+\d{3,}`)
					match := reg.FindAllString(formula, -1)
					for _, vMatch := range match {
						//cari jml berdasarkan code
						if rowCode[vMatch] != 0 {
							formula = strings.ReplaceAll(formula, vMatch, fmt.Sprintf("%c%d", chr, rowCode[vMatch]))
						}

					}
					f.SetCellFormula(sheet, fmt.Sprintf("%c%d", chr, row), fmt.Sprintf("=%s", formula))
				}
				row++
				continue
			}

			if v.AutoSummary != nil && *v.AutoSummary {
				f.SetCellStyle(sheet, fmt.Sprintf("B%d", row), fmt.Sprintf("C%d", row), styleLabelTotal)
				f.SetCellStyle(sheet, fmt.Sprintf("D%d", row), fmt.Sprintf("K%d", row), styleCurrencyTotal)
				f.SetCellFormula(sheet, fmt.Sprintf("F%d", row), fmt.Sprintf("=SUM(F%d:F%d)", partRowStart, row-1))
				for chr := 'H'; chr <= 'K'; chr++ {
					f.SetCellFormula(sheet, fmt.Sprintf("%c%d", chr, row), fmt.Sprintf("=SUM(%c%d:%c%d)", chr, partRowStart, chr, row-1))
				}
				row++
				partRowStart = row
				continue
			}

			criteriaIT := model.InvestasiTbkDetailFilterModel{}
			criteriaIT.Stock = &v.Code
			criteriaIT.FormatterBridgesID = &bridges.ID
			criteriaIT.InvestasiTbkID = &investasitbk.ID

			paginationIT := abstraction.Pagination{}
			pagesize := 10000
			// sortBy := "id"
			// sort := "asc"
			paginationIT.PageSize = &pagesize
			// paginationIT.SortBy = &sortBy
			// paginationIT.Sort = &sort

			AgingUPDetail, _, err := s.InvestasiTbkDetailRepository.Find(ctx, &criteriaIT, &paginationIT)
			if err != nil && v.Code != "" {
				continue
			}
			if len(*AgingUPDetail) == 0 {
				continue
			}
			for _, vv := range *AgingUPDetail {
				f.SetCellValue(sheet, fmt.Sprintf("B%d", row), vv.SortID)
				f.SetCellValue(sheet, fmt.Sprintf("C%d", row), vv.Stock)
				f.SetCellValue(sheet, fmt.Sprintf("D%d", row), *vv.EndingShares)
				f.SetCellValue(sheet, fmt.Sprintf("E%d", row), *vv.AvgPrice)
				f.SetCellValue(sheet, fmt.Sprintf("F%d", row), *vv.AmountCost)
				f.SetCellValue(sheet, fmt.Sprintf("G%d", row), *vv.ClosingPrice)
				f.SetCellValue(sheet, fmt.Sprintf("H%d", row), *vv.AmountFv)
				f.SetCellValue(sheet, fmt.Sprintf("I%d", row), *vv.UnrealizedGain)
				f.SetCellValue(sheet, fmt.Sprintf("J%d", row), *vv.RealizedGain)
				f.SetCellValue(sheet, fmt.Sprintf("K%d", row), *vv.Fee)
			}
			row++
		}
		rowStart = row
		row = rowStart
	}
	tmpFolder := fmt.Sprintf("assets/%d", ctx.Auth.ID)
	_, err = os.Stat(tmpFolder)
	if err != nil {
		if os.IsNotExist(err) {
			if err := os.MkdirAll(tmpFolder, os.ModePerm); err != nil {
				log.Fatal(err)
			}
		} else {
			return nil, err
		}
	}

	period := datePeriod.Format("2006-01-02")
	fileName := fmt.Sprintf("InvestasiTbk_%s.xlsx", period)
	fileLoc := fmt.Sprintf("assets/%d/%s", ctx.Auth.ID, fileName)
	err = f.SaveAs(fileLoc)
	if err != nil {
		return nil, err
	}

	result := &dto.InvestasiTbkExportResponse{
		FileName: fileName,
		Path:     fileLoc,
	}
	return result, nil
}

func (s *service) GetVersion(ctx *abstraction.Context, payload *dto.GetVersionRequest) (*dto.GetVersionResponse, error) {
	filter := model.InvestasiTbkFilterModel{
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
		return &dto.GetVersionResponse{}, response.ErrorBuilder(&response.ErrorConstant.InternalServerError, err)
	}
	result := &dto.GetVersionResponse{
		Data: *data,
	}
	return result, nil
}
