package mutasifa

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
	Repository                 repository.MutasiFa
	ParameterRepository        repository.Parameter
	FormatterRepository        repository.Formatter
	MutasiFaDetailRepository   repository.MutasiFaDetail
	FormatterBridgesRepository repository.FormatterBridges
	Db                         *gorm.DB
}

type Service interface {
	Find(ctx *abstraction.Context, payload *dto.MutasiFaGetRequest) (*dto.MutasiFaGetResponse, error)
	FindByID(ctx *abstraction.Context, payload *dto.MutasiFaGetByIDRequest) (*dto.MutasiFaGetByIDResponse, error)
	Create(ctx *abstraction.Context, payload *dto.MutasiFaCreateRequest) (*dto.MutasiFaCreateResponse, error)
	Update(ctx *abstraction.Context, payload *dto.MutasiFaUpdateRequest) (*dto.MutasiFaUpdateResponse, error)
	Delete(ctx *abstraction.Context, payload *dto.MutasiFaDeleteRequest) (*dto.MutasiFaDeleteResponse, error)
	Export(ctx *abstraction.Context) (*dto.MutasiFaExportResponse, error)
	GetVersion(ctx *abstraction.Context, payload *dto.GetVersionRequest) (*dto.GetVersionResponse, error)
}

func NewService(f *factory.Factory) *service {
	repository := f.MutasiFaRepository
	parameter := f.ParameterRepository
	formatter := f.FormatterRepository
	repositoryDetail := f.MutasiFaDetailRepository
	formatterBridges := f.FormatterBridgesRepository
	db := f.Db
	return &service{
		Repository:                 repository,
		ParameterRepository:        parameter,
		FormatterRepository:        formatter,
		MutasiFaDetailRepository:   repositoryDetail,
		FormatterBridgesRepository: formatterBridges,
		Db:                         db,
	}
}

func (s *service) Find(ctx *abstraction.Context, payload *dto.MutasiFaGetRequest) (*dto.MutasiFaGetResponse, error) {
	data, info, err := s.Repository.Find(ctx, &payload.MutasiFaFilterModel, &payload.Pagination)
	if err != nil {
		return &dto.MutasiFaGetResponse{}, response.ErrorBuilder(&response.ErrorConstant.InternalServerError, err)
	}

	result := &dto.MutasiFaGetResponse{
		Datas:          *data,
		PaginationInfo: *info,
	}
	return result, nil
}

func (s *service) FindByID(ctx *abstraction.Context, payload *dto.MutasiFaGetByIDRequest) (*dto.MutasiFaGetByIDResponse, error) {
	data, err := s.Repository.FindByID(ctx, &payload.ID)
	if err != nil {
		return &dto.MutasiFaGetByIDResponse{}, helper.ErrorHandler(err)
	}
	allowed := helper.CompanyValidation(ctx.Auth.ID, data.CompanyID)
	if !allowed {
		return nil, response.ErrorBuilder(&response.ErrorConstant.UnauthorizedAccess, errors.New("not allowed"))
	}
	result := &dto.MutasiFaGetByIDResponse{
		MutasiFaEntityModel: *data,
	}
	return result, nil
}

func (s *service) Create(ctx *abstraction.Context, payload *dto.MutasiFaCreateRequest) (*dto.MutasiFaCreateResponse, error) {
	var data model.MutasiFaEntityModel

	allowed := helper.CompanyValidation(ctx.Auth.ID, payload.CompanyID)
	if !allowed {
		return nil, response.ErrorBuilder(&response.ErrorConstant.UnauthorizedAccess, errors.New("not allowed"))
	}

	if err := trxmanager.New(s.Db).WithTrx(ctx, func(ctx *abstraction.Context) error {
		data.Context = ctx
		data.MutasiFaEntity = payload.MutasiFaEntity
		data.MutasiFaEntity.Status = 1

		result, err := s.Repository.Create(ctx, &data)
		if err != nil {
			return helper.ErrorHandler(err)
		}
		data = *result
		return nil
	}); err != nil {
		return &dto.MutasiFaCreateResponse{}, err
	}

	result := &dto.MutasiFaCreateResponse{
		MutasiFaEntityModel: data,
	}
	return result, nil
}

func (s *service) Update(ctx *abstraction.Context, payload *dto.MutasiFaUpdateRequest) (*dto.MutasiFaUpdateResponse, error) {
	var data model.MutasiFaEntityModel
	if err := trxmanager.New(s.Db).WithTrx(ctx, func(ctx *abstraction.Context) error {
		existing, err := s.Repository.FindByID(ctx, &payload.ID)
		if err != nil {
			return helper.ErrorHandler(err)
		}

		if existing.Status != 1 {
			return response.ErrorBuilder(&response.ErrorConstant.DataValidated, errors.New("cannot update data"))
		}

		allowed := helper.CompanyValidation(ctx.Auth.ID, existing.CompanyID)
		allowed2 := helper.CompanyValidation(ctx.Auth.ID, payload.CompanyID)
		if !allowed || !allowed2 {
			return response.ErrorBuilder(&response.ErrorConstant.UnauthorizedAccess, errors.New("not allowed"))
		}

		data.Context = ctx
		data.MutasiFaEntity = payload.MutasiFaEntity

		result, err := s.Repository.Update(ctx, &payload.ID, &data)
		if err != nil {
			return helper.ErrorHandler(err)
		}
		data = *result
		return nil
	}); err != nil {
		return &dto.MutasiFaUpdateResponse{}, err
	}
	result := &dto.MutasiFaUpdateResponse{
		MutasiFaEntityModel: data,
	}
	return result, nil
}

func (s *service) Delete(ctx *abstraction.Context, payload *dto.MutasiFaDeleteRequest) (*dto.MutasiFaDeleteResponse, error) {
	var data model.MutasiFaEntityModel
	if err := trxmanager.New(s.Db).WithTrx(ctx, func(ctx *abstraction.Context) error {
		existing, err := s.Repository.FindByID(ctx, &payload.ID)
		if err != nil {
			return helper.ErrorHandler(err)
		}

		allowed := helper.CompanyValidation(ctx.Auth.ID, existing.CompanyID)
		if !allowed {
			return response.ErrorBuilder(&response.ErrorConstant.UnauthorizedAccess, errors.New("not allowed"))
		}

		data.Context = ctx
		result, err := s.Repository.Delete(ctx, &payload.ID, &data)
		if err != nil {
			return helper.ErrorHandler(err)
		}
		data = *result
		return nil
	}); err != nil {
		return &dto.MutasiFaDeleteResponse{}, err
	}
	result := &dto.MutasiFaDeleteResponse{
		// MutasiFaEntityModel: data,
	}
	return result, nil
}

func (s *service) Export(ctx *abstraction.Context, payload *dto.MutasiFaExportRequest) (*dto.MutasiFaExportResponse, error) {
	f := excelize.NewFile()
	sheet := f.GetSheetName(f.GetActiveSheetIndex())

	arrStyleColWidth := []map[string]interface{}{
		{"COL": "A", "WIDTH": 8.21},
		{"COL": "B", "WIDTH": 3.74},
		{"COL": "C", "WIDTH": 33.57},
		{"COL": "D", "WIDTH": 18.21},
		{"COL": "E", "WIDTH": 16.60},
		{"COL": "F", "WIDTH": 16.60},
		{"COL": "G", "WIDTH": 19.64},
		{"COL": "H", "WIDTH": 19.64},
		{"COL": "I", "WIDTH": 18.21},
		{"COL": "J", "WIDTH": 17.68},
		{"COL": "K", "WIDTH": 16.78},
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
			Bold: true,
		},
		Fill: excelize.Fill{
			Type:    "pattern",
			Pattern: 1,
			Color:   []string{"#66ff33"},
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
			// {Type: "top", Color: "000000", Style: 1},
			// {Type: "bottom", Color: "000000", Style: 1},
			{Type: "left", Color: "000000", Style: 1},
			{Type: "right", Color: "000000", Style: 1},
		},
		CustomNumFmt: &numberFormat,
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
			Bold: true,
		},
		Fill: excelize.Fill{
			Type:    "pattern",
			Pattern: 1,
			Color:   []string{"#66ff33"},
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
			Bold: true,
		},
		Fill: excelize.Fill{
			Type:    "pattern",
			Pattern: 1,
			Color:   []string{"#66ff33"},
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
	})
	if err != nil {
		return nil, err
	}

	styleCurrencyWoBorder, err := f.NewStyle(&excelize.Style{
		CustomNumFmt: &numberFormat,
	})
	if err != nil {
		return nil, err
	}

	// criteria := model.MutasiFaFilterModel{}
	// criteria.CompanyID = &payload.CompanyID
	// version := payload.Versions
	// criteria.Versions = &version
	// criteria.Period = &payload.Period
	// criteria.FormatterID = &data.ID

	mutasifa, err := s.Repository.FindByID(ctx, &payload.MutasiFaID)
	if err != nil {
		return nil, err
	}

	if mutasifa.ID == 0 {
		return nil, errors.New("data not found")
	}

	allowed := helper.CompanyValidation(ctx.Auth.ID, mutasifa.CompanyID)
	if !allowed {
		return nil, response.ErrorBuilder(&response.ErrorConstant.UnauthorizedAccess, errors.New("not allowed"))
	}

	datePeriod, err := time.Parse(time.RFC3339, mutasifa.Period)
	if err != nil {
		return nil, err
	}

	formatterCode := []string{"MUTASI-FA-COST", "MUTASI-FA-ACCUMULATED-DEPRECATION"}
	formatterTitle := []string{"Mutasi Fixed Assets (FA)", ""}
	row, rowStart := 7, 7

	f.SetCellValue(sheet, "B2", formatterTitle[0])

	f.MergeCell(sheet, "B4", "C6")
	f.MergeCell(sheet, "D4", "K4")
	f.MergeCell(sheet, "D5", "D6")
	f.MergeCell(sheet, "E5", "E6")
	f.MergeCell(sheet, "F5", "F6")
	f.MergeCell(sheet, "G5", "G6")
	f.MergeCell(sheet, "H5", "H6")
	f.MergeCell(sheet, "I5", "I6")
	f.MergeCell(sheet, "J5", "J6")
	f.MergeCell(sheet, "K5", "K6")

	f.SetCellStyle(sheet, "B4", "K6", styleHeader)
	f.SetCellValue(sheet, "B4", mutasifa.Company.Name)
	f.SetCellValue(sheet, "D4", datePeriod.Format("02 January 2006"))
	f.SetCellValue(sheet, "D5", "Beginning Balance")
	f.SetCellValue(sheet, "E5", "Acquisition of Subsidiary")
	f.SetCellValue(sheet, "F5", "Additions (+)")
	f.SetCellValue(sheet, "G5", "Deductions (-)")
	f.SetCellValue(sheet, "H5", "Reclassification")
	f.SetCellValue(sheet, "I5", "Revaluation")
	f.SetCellValue(sheet, "J5", "Ending balance")
	f.SetCellValue(sheet, "K5", "Control")
	rowCode := make(map[string]int)
	for _, formatter := range formatterCode {

		var criteria dto.FormatterGetRequest
		criteria.FormatterFilterModel.FormatterFor = &formatter

		data, err := s.FormatterRepository.FindWithDetail(ctx, &criteria.FormatterFilterModel)
		if err != nil {
			return &dto.MutasiFaExportResponse{}, response.ErrorBuilder(&response.ErrorConstant.InternalServerError, err)
		}

		tmpStr := "MUTASI-FA"
		criteriaBridge := model.FormatterBridgesFilterModel{}
		criteriaBridge.FormatterBridgesFilter.Source = &tmpStr
		criteriaBridge.FormatterBridgesFilter.FormatterID = &data.ID
		criteriaBridge.FormatterBridgesFilter.TrxRefID = &mutasifa.ID

		bridges, err := s.FormatterBridgesRepository.FindWithCriteria(ctx, &criteriaBridge)
		if err != nil {
			return nil, err
		}
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

			if v.IsTotal != nil && *v.IsTotal {
				f.SetCellStyle(sheet, fmt.Sprintf("B%d", row), fmt.Sprintf("C%d", row), styleLabelTotal)
				f.SetCellStyle(sheet, fmt.Sprintf("D%d", row), fmt.Sprintf("K%d", row), styleCurrencyTotal)
				if v.FxSummary == "" {
					row++
					continue
				}
				arrChr := []string{"D", "E", "F", "G", "H", "I", "J", "K"}
				if strings.ToUpper(v.Code) == "NET_BOOK_VALUE" {
					arrChr = []string{"D", "J"}
				}
				for _, chr := range arrChr {
					formula := v.FxSummary
					reg := regexp.MustCompile(`[A-Za-z_~]+`)
					match := reg.FindAllString(formula, -1)
					for _, vMatch := range match {
						//cari jml berdasarkan code
						if rowCode[vMatch] != 0 {
							formula = strings.ReplaceAll(formula, vMatch, fmt.Sprintf("%s%d", chr, rowCode[vMatch]))
						}

					}
					f.SetCellFormula(sheet, fmt.Sprintf("%s%d", chr, row), fmt.Sprintf("=%s", formula))
				}
				row++
				continue
			}

			if v.AutoSummary != nil && *v.AutoSummary {
				f.SetCellStyle(sheet, fmt.Sprintf("B%d", row), fmt.Sprintf("C%d", row), styleLabelTotal)
				f.SetCellStyle(sheet, fmt.Sprintf("D%d", row), fmt.Sprintf("K%d", row), styleCurrencyTotal)
				f.MergeCell(sheet, fmt.Sprintf("B%d", row), fmt.Sprintf("C%d", row))
				if strings.ToUpper(v.Code) == "NET_BOOK_VALUE" {
					f.SetCellFormula(sheet, fmt.Sprintf("D%d", row), fmt.Sprintf("=SUM(D%d:D%d)", partRowStart, row-1))
					f.SetCellFormula(sheet, fmt.Sprintf("J%d", row), fmt.Sprintf("=SUM(J%d:J%d)", partRowStart, row-1))
				} else {
					for chr := 'D'; chr <= 'K'; chr++ {
						f.SetCellFormula(sheet, fmt.Sprintf("%c%d", chr, row), fmt.Sprintf("=SUM(%c%d:%c%d)", chr, partRowStart, chr, row-1))
					}
				}
				row++
				partRowStart = row
				continue
			}

			criteriaMF := model.MutasiFaDetailFilterModel{}
			criteriaMF.Code = &v.Code
			criteriaMF.FormatterBridgesID = &bridges.ID
			criteriaMF.MutasiFaID = &mutasifa.ID

			paginationMF := abstraction.Pagination{}
			pagesize := 10000
			// sortBy := "id"
			// sort := "asc"
			paginationMF.PageSize = &pagesize
			// paginationMF.SortBy = &sortBy
			// paginationMF.Sort = &sort

			MutasiFaDetail, _, err := s.MutasiFaDetailRepository.Find(ctx, &criteriaMF, &paginationMF)
			if err != nil && v.Code != "" {
				continue
			}
			if len(*MutasiFaDetail) == 0 {
				continue
			}
			for _, vv := range *MutasiFaDetail {
				f.SetCellValue(sheet, fmt.Sprintf("B%d", row), vv.Description)
				if vv.BeginningBalance != nil && *vv.BeginningBalance != 0 {
					f.SetCellValue(sheet, fmt.Sprintf("D%d", row), *vv.BeginningBalance)
				}
				if vv.AcquisitionOfSubsidiary != nil && *vv.AcquisitionOfSubsidiary != 0 {
					f.SetCellValue(sheet, fmt.Sprintf("E%d", row), *vv.AcquisitionOfSubsidiary)
				}
				if vv.Additions != nil && *vv.Additions != 0 {
					f.SetCellValue(sheet, fmt.Sprintf("F%d", row), *vv.Additions)
				}
				if vv.Deductions != nil && *vv.Deductions != 0 {
					f.SetCellValue(sheet, fmt.Sprintf("G%d", row), *vv.Deductions)
				}
				if vv.Reclassification != nil && *vv.Reclassification != 0 {
					f.SetCellValue(sheet, fmt.Sprintf("H%d", row), *vv.Reclassification)
				}
				if vv.Revaluation != nil && *vv.Revaluation != 0 {
					f.SetCellValue(sheet, fmt.Sprintf("I%d", row), *vv.Revaluation)
				}
				// f.SetCellValue(sheet, fmt.Sprintf("J%d", row), *vv.EndingBalance)
				f.SetCellFormula(sheet, fmt.Sprintf("J%d", row), fmt.Sprintf("=D%d+E%d+F%d-G%d+H%d+I%d", row, row, row, row, row, row))
				if vv.Control != nil && *vv.Control != 0 {
					f.SetCellValue(sheet, fmt.Sprintf("K%d", row), *vv.Control)
				}
			}
			row++
		}
		rowStart = row
		row = rowStart
	}
	//Penambahan detail pengurangan
	row += 2
	var criteria dto.FormatterGetRequest
	tmpStr := "MUTASI-DETAIL-PENGURANGAN"
	criteria.FormatterFilterModel.FormatterFor = &tmpStr

	data, err := s.FormatterRepository.FindWithDetail(ctx, &criteria.FormatterFilterModel)
	if err != nil {
		return &dto.MutasiFaExportResponse{}, response.ErrorBuilder(&response.ErrorConstant.InternalServerError, err)
	}

	criteriaBridge := model.FormatterBridgesFilterModel{}
	criteriaBridge.FormatterBridgesFilter.Source = &tmpStr
	criteriaBridge.FormatterBridgesFilter.FormatterID = &data.ID
	criteriaBridge.FormatterBridgesFilter.TrxRefID = &mutasifa.ID

	bridges, err := s.FormatterBridgesRepository.FindWithCriteria(ctx, &criteriaBridge)
	if err != nil {
		return nil, err
	}
	partRowStart := row
	f.SetCellValue(sheet, fmt.Sprintf("E%d", row), "Penjualan")
	f.SetCellValue(sheet, fmt.Sprintf("F%d", row), "Penghapusan")
	for _, v := range data.FormatterDetail {
		rowCode[v.Code] = row
		f.SetCellStyle(sheet, fmt.Sprintf("E%d", row), fmt.Sprintf("F%d", row), styleCurrencyWoBorder)

		if strings.Contains(strings.ToLower(v.Code), "blank") {
			row++
			continue
		}

		f.SetCellValue(sheet, fmt.Sprintf("B%d", row), v.Description)

		if v.IsTotal != nil && *v.IsTotal {
			if v.FxSummary == "" {
				row++
				continue
			}
			arrChr := []string{"E", "F"}
			for _, chr := range arrChr {
				formula := v.FxSummary
				reg := regexp.MustCompile(`[A-Za-z_~]+`)
				match := reg.FindAllString(formula, -1)
				for _, vMatch := range match {
					//cari jml berdasarkan code
					if rowCode[vMatch] != 0 {
						if v.IsCoa != nil && *v.IsCoa {
							if chr == "E" {
								formula = strings.ReplaceAll(formula, vMatch, fmt.Sprintf("G%d", rowCode[vMatch]))
							} else {
								formula = strings.ReplaceAll(formula, vMatch, fmt.Sprintf("I%d", rowCode[vMatch]))
							}
						} else {
							formula = strings.ReplaceAll(formula, vMatch, fmt.Sprintf("%s%d", chr, rowCode[vMatch]))
						}
					}

				}
				err = f.SetCellFormula(sheet, fmt.Sprintf("%s%d", chr, row), fmt.Sprintf("=%s", formula))
				if err != nil {
					fmt.Println(err)
				}
			}
			row++
			continue
		}

		if v.AutoSummary != nil && *v.AutoSummary {
			for chr := 'E'; chr <= 'F'; chr++ {
				f.SetCellFormula(sheet, fmt.Sprintf("%c%d", chr, row), fmt.Sprintf("=SUM(%c%d:%c%d)", chr, partRowStart, chr, row-1))
			}
			row++
			partRowStart = row
			continue
		}

		criteriaMF := model.MutasiFaDetailFilterModel{}
		criteriaMF.Code = &v.Code
		criteriaMF.FormatterBridgesID = &bridges.ID
		criteriaMF.MutasiFaID = &mutasifa.ID
		paginationMF := abstraction.Pagination{}
		pagesize := 10000
		paginationMF.PageSize = &pagesize
		MutasiFaDetail, _, err := s.MutasiFaDetailRepository.Find(ctx, &criteriaMF, &paginationMF)
		if err != nil {
			row++
			continue
		}
		if len(*MutasiFaDetail) == 0 {
			row++
			continue
		}
		for _, vv := range *MutasiFaDetail {
			f.SetCellValue(sheet, fmt.Sprintf("B%d", row), vv.Description)
			if vv.Deductions != nil && *vv.Deductions != 0 {
				f.SetCellValue(sheet, fmt.Sprintf("G%d", row), *vv.Deductions)
			}
			if vv.Revaluation != nil && *vv.Revaluation != 0 {
				f.SetCellValue(sheet, fmt.Sprintf("I%d", row), *vv.Revaluation)
			}
		}
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
			return nil, err
		}
	}

	period := datePeriod.Format("2006-01-02")
	fileName := fmt.Sprintf("MutasiFa_%s.xlsx", period)
	fileLoc := fmt.Sprintf("assets/%d/%s", ctx.Auth.ID, fileName)
	err = f.SaveAs(fileLoc)
	if err != nil {
		return nil, err
	}

	result := &dto.MutasiFaExportResponse{
		FileName: fileName,
		Path:     fileLoc,
	}
	return result, nil
}

func (s *service) GetVersion(ctx *abstraction.Context, payload *dto.GetVersionRequest) (*dto.GetVersionResponse, error) {
	filter := model.MutasiFaFilterModel{
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
