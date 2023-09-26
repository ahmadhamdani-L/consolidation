package mutasipersediaan

import (
	"errors"
	"fmt"
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
	"strings"
	"time"

	"github.com/xuri/excelize/v2"
	"gorm.io/gorm"
)

type service struct {
	Repository                       repository.MutasiPersediaan
	MutasiPersediaanDetailRepository repository.MutasiPersediaanDetail
	ParameterRepository              repository.Parameter
	FormatterRepository              repository.Formatter
	FormatterBridgesRepository       repository.FormatterBridges
	Db                               *gorm.DB
}

type Service interface {
	Find(ctx *abstraction.Context, payload *dto.MutasiPersediaanGetRequest) (*dto.MutasiPersediaanGetResponse, error)
	FindByID(ctx *abstraction.Context, payload *dto.MutasiPersediaanGetByIDRequest) (*dto.MutasiPersediaanGetByIDResponse, error)
	Create(ctx *abstraction.Context, payload *dto.MutasiPersediaanCreateRequest) (*dto.MutasiPersediaanCreateResponse, error)
	Update(ctx *abstraction.Context, payload *dto.MutasiPersediaanUpdateRequest) (*dto.MutasiPersediaanUpdateResponse, error)
	Delete(ctx *abstraction.Context, payload *dto.MutasiPersediaanDeleteRequest) (*dto.MutasiPersediaanDeleteResponse, error)
	Import(ctx *abstraction.Context, payload *dto.MutasiPersediaanImportRequest) (*dto.MutasiPersediaanImportResponse, error)
	Export(ctx *abstraction.Context, payload *dto.MutasiPersediaanExportRequest) (*dto.MutasiPersediaanExportResponse, error)
	GetVersion(ctx *abstraction.Context, payload *dto.GetVersionRequest) (*dto.GetVersionResponse, error)
}

func NewService(f *factory.Factory) *service {
	repository := f.MutasiPersediaanRepository
	mutasiPersediaanDetail := f.MutasiPersediaanDetailRepository
	parameterRRepository := f.ParameterRepository
	formatterRepository := f.FormatterRepository
	formatterBridges := f.FormatterBridgesRepository
	db := f.Db
	return &service{
		Repository:                       repository,
		MutasiPersediaanDetailRepository: mutasiPersediaanDetail,
		ParameterRepository:              parameterRRepository,
		FormatterRepository:              formatterRepository,
		FormatterBridgesRepository:       formatterBridges,
		Db:                               db,
	}
}

func (s *service) Find(ctx *abstraction.Context, payload *dto.MutasiPersediaanGetRequest) (*dto.MutasiPersediaanGetResponse, error) {
	data, info, err := s.Repository.Find(ctx, &payload.MutasiPersediaanFilterModel, &payload.Pagination)
	if err != nil {
		return &dto.MutasiPersediaanGetResponse{}, helper.ErrorHandler(err)
	}

	result := &dto.MutasiPersediaanGetResponse{
		Datas:          *data,
		PaginationInfo: *info,
	}
	return result, nil
}

func (s *service) FindByID(ctx *abstraction.Context, payload *dto.MutasiPersediaanGetByIDRequest) (*dto.MutasiPersediaanGetByIDResponse, error) {
	data, err := s.Repository.FindByID(ctx, &payload.ID)
	if err != nil {
		return &dto.MutasiPersediaanGetByIDResponse{}, helper.ErrorHandler(err)
	}

	allowed := helper.CompanyValidation(ctx.Auth.ID, data.CompanyID)
	if !allowed {
		return nil, response.ErrorBuilder(&response.ErrorConstant.UnauthorizedAccess, errors.New("Not Allowed"))
	}

	result := &dto.MutasiPersediaanGetByIDResponse{
		MutasiPersediaanEntityModel: *data,
	}
	return result, nil
}

func (s *service) Create(ctx *abstraction.Context, payload *dto.MutasiPersediaanCreateRequest) (*dto.MutasiPersediaanCreateResponse, error) {
	var data model.MutasiPersediaanEntityModel

	allowed := helper.CompanyValidation(ctx.Auth.ID, payload.CompanyID)
	if !allowed {
		return nil, response.ErrorBuilder(&response.ErrorConstant.UnauthorizedAccess, errors.New("Not Allowed"))
	}

	if err := trxmanager.New(s.Db).WithTrx(ctx, func(ctx *abstraction.Context) error {
		data.Context = ctx
		data.MutasiPersediaanEntity = payload.MutasiPersediaanEntity
		data.MutasiPersediaanEntity.Status = 1

		result, err := s.Repository.Create(ctx, &data)
		if err != nil {
			return helper.ErrorHandler(err)
		}
		data = *result
		return nil
	}); err != nil {
		return &dto.MutasiPersediaanCreateResponse{}, err
	}

	result := &dto.MutasiPersediaanCreateResponse{
		MutasiPersediaanEntityModel: data,
	}
	return result, nil
}

func (s *service) Update(ctx *abstraction.Context, payload *dto.MutasiPersediaanUpdateRequest) (*dto.MutasiPersediaanUpdateResponse, error) {
	var data model.MutasiPersediaanEntityModel
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
		data.MutasiPersediaanEntity = payload.MutasiPersediaanEntity

		result, err := s.Repository.Update(ctx, &payload.ID, &data)
		if err != nil {
			return helper.ErrorHandler(err)
		}
		data = *result
		return nil
	}); err != nil {
		return &dto.MutasiPersediaanUpdateResponse{}, err
	}
	result := &dto.MutasiPersediaanUpdateResponse{
		MutasiPersediaanEntityModel: data,
	}
	return result, nil
}

func (s *service) Delete(ctx *abstraction.Context, payload *dto.MutasiPersediaanDeleteRequest) (*dto.MutasiPersediaanDeleteResponse, error) {
	var data model.MutasiPersediaanEntityModel
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
		return &dto.MutasiPersediaanDeleteResponse{}, err
	}
	result := &dto.MutasiPersediaanDeleteResponse{
		// MutasiPersediaanEntityModel: data,
	}
	return result, nil
}

func (s *service) Export(ctx *abstraction.Context, payload *dto.MutasiPersediaanExportRequest) (*dto.MutasiPersediaanExportResponse, error) {
	f := excelize.NewFile()
	sheet := f.GetSheetName(f.GetActiveSheetIndex())
	f.SetColWidth(sheet, "B", "B", 31.01)
	f.SetColWidth(sheet, "C", "C", 9.15)

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
			Color:   []string{"#f8cbad"},
		},
		Alignment: &excelize.Alignment{
			WrapText: true,
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
			Color:   []string{"#99ff66"},
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
			Color:   []string{"#99ff66"},
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

	formatterCode := []string{"MUTASI-PERSEDIAAN", "MUTASI-CADANGAN-PENGHAPUSAN-PERSEDIAAN"}
	formatterTitle := []string{"Mutasi Persediaan", "Mutasi Cadangan penghapusan persediaan"}
	row, rowStart := 4, 4

	// criteria := model.MutasiPersediaanFilterModel{}
	// criteria.CompanyID = &payload.CompanyID
	// version := payload.Versions
	// criteria.Versions = &version
	// criteria.Period = &payload.Period
	// criteria.FormatterID = &data.ID

	mutasipersediaan, err := s.Repository.FindByID(ctx, &payload.MutasiPersediaanID)
	if err != nil {
		return nil, helper.ErrorHandler(err)
	}

	allowed := helper.CompanyValidation(ctx.Auth.ID, mutasipersediaan.CompanyID)
	if !allowed {
		return nil, response.ErrorBuilder(&response.ErrorConstant.UnauthorizedAccess, errors.New("Not Allowed"))
	}

	for i, formatter := range formatterCode {
		f.SetCellStyle(sheet, fmt.Sprintf("B%d", (rowStart-1)), fmt.Sprintf("C%d", (rowStart-1)), styleHeader)
		f.SetCellValue(sheet, fmt.Sprintf("B%d", (rowStart-2)), formatterTitle[i])
		f.SetCellValue(sheet, fmt.Sprintf("B%d", (rowStart-1)), "Description")
		f.SetCellValue(sheet, fmt.Sprintf("C%d", (rowStart-1)), "Amount")

		var criteria model.FormatterFilterModel
		criteria.FormatterFor = &formatter

		data, err := s.FormatterRepository.FindWithDetail(ctx, &criteria)
		if err != nil {
			return nil, helper.ErrorHandler(err)
		}

		tmpStr := "MUTASI-PERSEDIAAN"
		criteriaBridge := model.FormatterBridgesFilterModel{}
		criteriaBridge.FormatterBridgesFilter.Source = &tmpStr
		criteriaBridge.FormatterBridgesFilter.FormatterID = &data.ID
		criteriaBridge.FormatterBridgesFilter.TrxRefID = &mutasipersediaan.ID

		bridges, err := s.FormatterBridgesRepository.FindWithCriteria(ctx, &criteriaBridge)
		if err != nil {
			return nil, helper.ErrorHandler(err)
		}
		rowCode := make(map[string]int)
		partRowStart := row
		for _, v := range data.FormatterDetail {
			rowCode[v.Code] = row
			f.SetCellStyle(sheet, fmt.Sprintf("B%d", row), fmt.Sprintf("B%d", row), styleLabel)
			f.SetCellStyle(sheet, fmt.Sprintf("C%d", row), fmt.Sprintf("C%d", row), styleCurrency)

			if strings.Contains(strings.ToLower(v.Code), "blank") {
				row++
				continue
			}

			f.SetCellValue(sheet, fmt.Sprintf("B%d", row), v.Description)

			if v.IsTotal != nil && *v.IsTotal == true {
				f.SetCellStyle(sheet, fmt.Sprintf("B%d", row), fmt.Sprintf("B%d", row), styleLabelTotal)
				f.SetCellStyle(sheet, fmt.Sprintf("C%d", row), fmt.Sprintf("C%d", row), styleCurrencyTotal)
				if v.FxSummary == "" {
					row++
					continue
				}
				formula := v.FxSummary
				reg := regexp.MustCompile(`[A-Za-z_~]+`)
				match := reg.FindAllString(formula, -1)
				for _, vMatch := range match {
					//cari jml berdasarkan code
					if rowCode[vMatch] != 0 {
						formula = strings.ReplaceAll(formula, vMatch, fmt.Sprintf("C%d", rowCode[vMatch]))
					}

				}
				f.SetCellFormula(sheet, fmt.Sprintf("C%d", row), fmt.Sprintf("=%s", formula))
				row++
				continue
			}

			if v.AutoSummary != nil && *v.AutoSummary == true {
				f.SetCellStyle(sheet, fmt.Sprintf("B%d", row), fmt.Sprintf("B%d", row), styleLabelTotal)
				f.SetCellStyle(sheet, fmt.Sprintf("C%d", row), fmt.Sprintf("C%d", row), styleCurrencyTotal)
				// f.MergeCell(sheet, fmt.Sprintf("B%d", row), fmt.Sprintf("C%d", row))

				f.SetCellFormula(sheet, fmt.Sprintf("C%d", row), fmt.Sprintf("=SUM(C%d:C%d)", partRowStart, row-1))
				row++
				partRowStart = row
				continue
			}

			criteriaMP := model.MutasiPersediaanDetailFilterModel{}
			criteriaMP.Code = &v.Code
			criteriaMP.FormatterBridgesID = &bridges.ID
			criteriaMP.MutasiPersediaanID = &mutasipersediaan.ID

			paginationMP := abstraction.Pagination{}
			pagesize := 1
			paginationMP.PageSize = &pagesize

			mutasiPersediaanDetail, _, err := s.MutasiPersediaanDetailRepository.Find(ctx, &criteriaMP, &paginationMP)
			if err != nil {
				return nil, helper.ErrorHandler(err)
			}
			for _, vv := range *mutasiPersediaanDetail {
				if vv.Amount != nil && *vv.Amount != 0 {
					f.SetCellValue(sheet, fmt.Sprintf("C%d", row), *vv.Amount)
				}
			}
			row++
		}
		rowStart = row + 4
		row = rowStart
	}
	tmpFolder := fmt.Sprintf("assets/%d", ctx.Auth.ID)
	_, err = os.Stat(tmpFolder)
	if err != nil {
		if os.IsNotExist(err) {
			if err := os.MkdirAll(tmpFolder, os.ModePerm); err != nil {
				return nil, err
			}
		} else {
			return nil, err
		}
	}
	datePeriod, err := time.Parse(time.RFC3339, mutasipersediaan.Period)
	if err != nil {
		return nil, err
	}
	period := datePeriod.Format("2006-01-02")
	fileName := fmt.Sprintf("MutasiPersediaan_%s.xlsx", period)
	fileLoc := fmt.Sprintf("assets/%d/%s", ctx.Auth.ID, fileName)
	err = f.SaveAs(fileLoc)
	if err != nil {
		return nil, err
	}

	result := &dto.MutasiPersediaanExportResponse{
		FileName: fileName,
		Path:     fileLoc,
	}
	return result, nil
}

func (s *service) GetVersion(ctx *abstraction.Context, payload *dto.GetVersionRequest) (*dto.GetVersionResponse, error) {
	filter := model.MutasiPersediaanFilterModel{
		CompanyCustomFilter: model.CompanyCustomFilter{
			CompanyID:          payload.CompanyID,
			ArrCompanyID:       payload.ArrCompanyID,
			ArrCompanyString:   payload.ArrCompanyString,
			ArrCompanyOperator: payload.ArrCompanyOperator,
		},
	}
	filter.Period = payload.Period
	filter.Status = payload.Status
	data, err := s.Repository.GetVersion(ctx, &filter)
	if err != nil {
		return &dto.GetVersionResponse{}, helper.ErrorHandler(err)
	}
	result := &dto.GetVersionResponse{
		Data: *data,
	}
	return result, nil
}
