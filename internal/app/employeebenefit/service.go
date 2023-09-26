package employeebenefit

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
	Repository                      repository.EmployeeBenefit
	Db                              *gorm.DB
	FormatterRepository             repository.Formatter
	EmployeeBenefitDetailRepository repository.EmployeeBenefitDetail
	ParameterRepository             repository.Parameter
	FormatterBridgesRepository      repository.FormatterBridges
}

type Service interface {
	Find(ctx *abstraction.Context, payload *dto.EmployeeBenefitGetRequest) (*dto.EmployeeBenefitGetResponse, error)
	FindByID(ctx *abstraction.Context, payload *dto.EmployeeBenefitGetByIDRequest) (*dto.EmployeeBenefitGetByIDResponse, error)
	Create(ctx *abstraction.Context, payload *dto.EmployeeBenefitCreateRequest) (*dto.EmployeeBenefitCreateResponse, error)
	Update(ctx *abstraction.Context, payload *dto.EmployeeBenefitUpdateRequest) (*dto.EmployeeBenefitUpdateResponse, error)
	Delete(ctx *abstraction.Context, payload *dto.EmployeeBenefitDeleteRequest) (*dto.EmployeeBenefitDeleteResponse, error)
	Import(ctx *abstraction.Context, payload *dto.EmployeeBenefitImportRequest, datas *[]model.EmployeeBenefitDetailEntity) (*dto.EmployeeBenefitImportResponse, error)
	Export(ctx *abstraction.Context, payload *dto.EmployeeBenefitExportRequest) (*dto.EmployeeBenefitExportResponse, error)
	GetVersion(ctx *abstraction.Context, payload *dto.GetVersionRequest) (*dto.GetVersionResponse, error)
}

func NewService(f *factory.Factory) *service {
	repository := f.EmployeeBenefitRepository
	formatterRepository := f.FormatterRepository
	employeeBenefitDetailRepository := f.EmployeeBenefitDetailRepository
	parameterRepository := f.ParameterRepository
	formatterBridgesRepo := f.FormatterBridgesRepository
	db := f.Db
	return &service{
		Repository:                      repository,
		Db:                              db,
		FormatterRepository:             formatterRepository,
		EmployeeBenefitDetailRepository: employeeBenefitDetailRepository,
		ParameterRepository:             parameterRepository,
		FormatterBridgesRepository:      formatterBridgesRepo,
	}
}

func (s *service) Find(ctx *abstraction.Context, payload *dto.EmployeeBenefitGetRequest) (*dto.EmployeeBenefitGetResponse, error) {
	data, info, err := s.Repository.Find(ctx, &payload.EmployeeBenefitFilterModel, &payload.Pagination)
	if err != nil {
		return &dto.EmployeeBenefitGetResponse{}, helper.ErrorHandler(err)
	}

	result := &dto.EmployeeBenefitGetResponse{
		Datas:          *data,
		PaginationInfo: *info,
	}
	return result, nil
}

func (s *service) FindByID(ctx *abstraction.Context, payload *dto.EmployeeBenefitGetByIDRequest) (*dto.EmployeeBenefitGetByIDResponse, error) {
	data, err := s.Repository.FindByID(ctx, &payload.ID)
	if err != nil {
		return &dto.EmployeeBenefitGetByIDResponse{}, helper.ErrorHandler(err)
	}

	allowed := helper.CompanyValidation(ctx.Auth.ID, data.CompanyID)
	if !allowed {
		return nil, response.ErrorBuilder(&response.ErrorConstant.UnauthorizedAccess, errors.New("Not Allowed"))
	}

	result := &dto.EmployeeBenefitGetByIDResponse{
		EmployeeBenefitEntityModel: *data,
	}
	return result, nil
}

func (s *service) Create(ctx *abstraction.Context, payload *dto.EmployeeBenefitCreateRequest) (*dto.EmployeeBenefitCreateResponse, error) {
	var data model.EmployeeBenefitEntityModel

	allowed := helper.CompanyValidation(ctx.Auth.ID, payload.CompanyID)
	if !allowed {
		return nil, response.ErrorBuilder(&response.ErrorConstant.UnauthorizedAccess, errors.New("Not Allowed"))
	}

	if err := trxmanager.New(s.Db).WithTrx(ctx, func(ctx *abstraction.Context) error {
		data.Context = ctx
		data.EmployeeBenefitEntity = payload.EmployeeBenefitEntity
		data.EmployeeBenefitEntity.Status = 1

		result, err := s.Repository.Create(ctx, &data)
		if err != nil {
			return helper.ErrorHandler(err)
		}
		data = *result
		return nil
	}); err != nil {
		return &dto.EmployeeBenefitCreateResponse{}, err
	}

	result := &dto.EmployeeBenefitCreateResponse{
		EmployeeBenefitEntityModel: data,
	}
	return result, nil
}

func (s *service) Update(ctx *abstraction.Context, payload *dto.EmployeeBenefitUpdateRequest) (*dto.EmployeeBenefitUpdateResponse, error) {
	var data model.EmployeeBenefitEntityModel
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
		data.EmployeeBenefitEntity = payload.EmployeeBenefitEntity

		result, err := s.Repository.Update(ctx, &payload.ID, &data)
		if err != nil {
			return helper.ErrorHandler(err)
		}
		data = *result
		return nil
	}); err != nil {
		return &dto.EmployeeBenefitUpdateResponse{}, err
	}
	result := &dto.EmployeeBenefitUpdateResponse{
		EmployeeBenefitEntityModel: data,
	}
	return result, nil
}

func (s *service) Delete(ctx *abstraction.Context, payload *dto.EmployeeBenefitDeleteRequest) (*dto.EmployeeBenefitDeleteResponse, error) {
	var data model.EmployeeBenefitEntityModel
	if err := trxmanager.New(s.Db).WithTrx(ctx, func(ctx *abstraction.Context) error {
		existing, err := s.Repository.FindByID(ctx, &payload.ID)
		if err != nil {
			return helper.ErrorHandler(err)
		}

		allowed := helper.CompanyValidation(ctx.Auth.ID, existing.CompanyID)
		if !allowed {
			return helper.ErrorHandler(err)
		}

		data.Context = ctx
		result, err := s.Repository.Delete(ctx, &payload.ID, &data)
		if err != nil {
			return helper.ErrorHandler(err)
		}
		data = *result
		return nil
	}); err != nil {
		return &dto.EmployeeBenefitDeleteResponse{}, helper.ErrorHandler(err)
	}
	result := &dto.EmployeeBenefitDeleteResponse{
		// EmployeeBenefitEntityModel: data,
	}
	return result, nil
}

func (s *service) Export(ctx *abstraction.Context, payload *dto.EmployeeBenefitExportRequest) (*dto.EmployeeBenefitExportResponse, error) {
	f := excelize.NewFile()
	sheet := f.GetSheetName(f.GetActiveSheetIndex())

	arrStyleColWidth := []map[string]interface{}{
		{"COL": "A", "WIDTH": 4.30},
		{"COL": "B", "WIDTH": 3.71},
		{"COL": "C", "WIDTH": 8.43},
		{"COL": "D", "WIDTH": 8.43},
		{"COL": "E", "WIDTH": 8.43},
		{"COL": "F", "WIDTH": 21.71},
		{"COL": "G", "WIDTH": 19.14},
	}
	for _, v := range arrStyleColWidth {
		tmpColWidth := fmt.Sprintf("%f", v["WIDTH"])
		colWidth, _ := strconv.ParseFloat(tmpColWidth, 64)
		err := f.SetColWidth(sheet, fmt.Sprintf("%s", v["COL"]), fmt.Sprintf("%s", v["COL"]), colWidth)
		if err != nil {
			panic(err)
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
			Color:   []string{"#f8cbad"},
		},
		Alignment: &excelize.Alignment{
			WrapText:   true,
			Horizontal: "center",
			Vertical:   "center",
		},
	})
	if err != nil {
		log.Println(err)
		return nil, err
	}
	numberFormat := "#,##"
	styleCurrency, err := f.NewStyle(&excelize.Style{
		Border: []excelize.Border{
			// {Type: "top", Color: "000000", Style: 1},
			// {Type: "bottom", Color: "000000", Style: 1},
			// {Type: "left", Color: "000000", Style: 1},
			// {Type: "right", Color: "000000", Style: 1},
		},
		CustomNumFmt: &numberFormat,
	})
	if err != nil {
		log.Println(err)
		return nil, err
	}

	styleCurrencyTotal, err := f.NewStyle(&excelize.Style{
		Border: []excelize.Border{
			// {Type: "top", Color: "000000", Style: 1},
			// {Type: "bottom", Color: "000000", Style: 1},
			// {Type: "left", Color: "000000", Style: 1},
			// {Type: "right", Color: "000000", Style: 1},
		},
		Font: &excelize.Font{
			Bold: true,
		},
		Fill: excelize.Fill{
			// Type:    "pattern",
			// Pattern: 1,
			// Color:   []string{"#99ff66"},
		},
		CustomNumFmt: &numberFormat,
	})
	if err != nil {
		log.Println(err)
		return nil, err
	}

	styleLabelTotal, err := f.NewStyle(&excelize.Style{
		Border: []excelize.Border{
			// {Type: "top", Color: "000000", Style: 1},
			// {Type: "bottom", Color: "000000", Style: 1},
			// {Type: "left", Color: "000000", Style: 1},
			// {Type: "right", Color: "000000", Style: 1},
		},
		Font: &excelize.Font{
			Bold: true,
		},
		Fill: excelize.Fill{
			// Type:    "pattern",
			// Pattern: 1,
			// Color:   []string{"#99ff66"},
		},
	})
	if err != nil {
		log.Println(err)
		return nil, err
	}

	styleLabel, err := f.NewStyle(&excelize.Style{
		Border: []excelize.Border{
			// {Type: "top", Color: "000000", Style: 1},
			// {Type: "bottom", Color: "000000", Style: 1},
			// {Type: "left", Color: "000000", Style: 1},
			// {Type: "right", Color: "000000", Style: 1},
		},
		Font: &excelize.Font{
			// Bold: true,
		},
		Fill: excelize.Fill{
			// Type:    "pattern",
			// Pattern: 1,
			// Color:   []string{"#99ff66"},
		},
	})
	if err != nil {
		log.Println(err)
		return nil, err
	}

	employeebenefit, err := s.Repository.FindByID(ctx, &payload.EmployeeBenefitID)
	if err != nil {
		return nil, err
	}

	if employeebenefit.ID == 0 {
		return nil, errors.New("Data Not Found")
	}

	allowed := helper.CompanyValidation(ctx.Auth.ID, employeebenefit.CompanyID)
	if !allowed {
		return nil, response.ErrorBuilder(&response.ErrorConstant.UnauthorizedAccess, errors.New("Not Allowed"))
	}

	datePeriod, err := time.Parse(time.RFC3339, employeebenefit.Period)
	if err != nil {
		return nil, err
	}
	period := datePeriod.Format("02-Jan-06")

	formatterCode := []string{"EMPLOYEE-BENEFIT-ASUMSI", "EMPLOYEE-BENEFIT-REKONSILIASI", "EMPLOYEE-BENEFIT-RINCIAN-LAPORAN", "EMPLOYEE-BENEFIT-RINCIAN-EKUITAS", "EMPLOYEE-BENEFIT-MUTASI", "EMPLOYEE-BENEFIT-INFORMASI", "EMPLOYEE-BENEFIT-ANALISIS"}
	formatterTitle := []string{"Asumsi-asumsi yang digunakan:", "Rekonsiliasi jumlah liabilitas imbalan kerja karyawan pada laporan posisi keuangan adalah sebagai berikut:", "Rincian beban imbalan kerja karyawan yang diakui dalam laporan laba rugi dan penghasilan komprehensif lain adalah sebagai berikut:", "Rincian beban imbalan kerja karyawan yang diakui pada ekuitas dalam penghasilan komprehensif lain adalah sebagai berikut:", "Mutasi liabilitas imbalan kerja karyawan adalah sebagai berikut:", "Informasi historis dari nilai kini liabilitas imbalan pasti, nilai wajar aset program dan penyesuaian adalah sebagai berikut:", "Analisis sensitivitas dari perubahan asumsi-asumsi utama terhadap liabilitas imbalan kerja", ""}
	row, rowStart := 7, 7
	rowCode := make(map[string]int)
	for i, formatter := range formatterCode {
		f.SetCellValue(sheet, fmt.Sprintf("B%d", (rowStart-3)), formatterTitle[i])
		f.MergeCell(sheet, fmt.Sprintf("B%d", (rowStart-2)), fmt.Sprintf("F%d", (rowStart-1)))
		f.SetCellStyle(sheet, fmt.Sprintf("B%d", (rowStart-2)), fmt.Sprintf("G%d", (rowStart-1)), styleHeader)
		f.SetCellStyle(sheet, fmt.Sprintf("A%d", (rowStart-3)), fmt.Sprintf("A%d", (rowStart-3)), styleLabel)

		f.SetCellValue(sheet, fmt.Sprintf("B%d", (rowStart-2)), "Description")
		f.SetCellValue(sheet, fmt.Sprintf("G%d", (rowStart-2)), period)
		f.SetCellValue(sheet, fmt.Sprintf("G%d", (rowStart-1)), employeebenefit.Company.Name)

		var criteria dto.FormatterGetRequest
		criteria.FormatterFilterModel.FormatterFor = &formatter

		data, err := s.FormatterRepository.FindWithDetail(ctx, &criteria.FormatterFilterModel)
		if err != nil {
			return &dto.EmployeeBenefitExportResponse{}, helper.ErrorHandler(err)
		}

		tmpStr := "EMPLOYEE-BENEFIT"
		criteriaBridge := model.FormatterBridgesFilterModel{}
		criteriaBridge.FormatterBridgesFilter.Source = &tmpStr
		criteriaBridge.FormatterBridgesFilter.FormatterID = &data.ID
		criteriaBridge.FormatterBridgesFilter.TrxRefID = &employeebenefit.ID

		bridges, err := s.FormatterBridgesRepository.FindWithCriteria(ctx, &criteriaBridge)
		if err != nil {
			return nil, err
		}
		partRowStart := row
		for _, v := range data.FormatterDetail {
			if strings.Contains(strings.ToLower(v.Code), "blank") {
				row++
				continue
			}
			if rowCode[v.Code] == 0 {
				rowCode[v.Code] = row
			}
			f.SetCellValue(sheet, fmt.Sprintf("B%d", row), v.Description)
			f.SetCellStyle(sheet, fmt.Sprintf("B%d", row), fmt.Sprintf("B%d", row), styleLabel)
			f.SetCellStyle(sheet, fmt.Sprintf("G%d", row), fmt.Sprintf("G%d", row), styleCurrency)

			if v.IsTotal != nil && *v.IsTotal == true {
				// f.SetCellStyle(sheet, fmt.Sprintf("B%d", row), fmt.Sprintf("B%d", row), styleLabelTotal)
				f.SetCellStyle(sheet, fmt.Sprintf("G%d", row), fmt.Sprintf("G%d", row), styleCurrencyTotal)
				if v.FxSummary == "" {
					row++
					continue
				}
				formula := v.FxSummary
				reg := regexp.MustCompile(`[A-Za-z0-9_~#:()]+|[0-9]+\d{3,}`)
				match := reg.FindAllString(formula, -1)
				for _, vMatch := range match {
					//cari jml berdasarkan code
					if _, ok := rowCode[vMatch]; ok {
						formula = strings.ReplaceAll(formula, vMatch, fmt.Sprintf("G%d", rowCode[vMatch]))
					}

				}
				f.SetCellFormula(sheet, fmt.Sprintf("G%d", row), fmt.Sprintf("=%s", formula))
				row++
				continue
			}

			if v.AutoSummary != nil && *v.AutoSummary == true {
				f.SetCellStyle(sheet, fmt.Sprintf("B%d", row), fmt.Sprintf("B%d", row), styleLabelTotal)
				f.SetCellStyle(sheet, fmt.Sprintf("G%d", row), fmt.Sprintf("G%d", row), styleCurrencyTotal)
				f.SetCellFormula(sheet, fmt.Sprintf("G%d", row), fmt.Sprintf("=SUM($G%d:$G%d)", partRowStart, row-1))
				row++
				partRowStart = row
				continue
			}

			if v.IsLabel != nil && *v.IsLabel == true {
				row++
				continue
			}

			criteriaEB := model.EmployeeBenefitDetailFilterModel{}
			criteriaEB.Code = &v.Code
			criteriaEB.FormatterBridgesID = &bridges.ID
			criteriaEB.EmployeeBenefitID = &employeebenefit.ID

			paginationEB := abstraction.Pagination{}
			pagesize := 10000
			paginationEB.PageSize = &pagesize

			employeeBenefitDetail, _, err := s.EmployeeBenefitDetailRepository.Find(ctx, &criteriaEB, &paginationEB)
			if err != nil {
				return nil, err
			}
			if len(*employeeBenefitDetail) == 0 {
				continue
			}
			for _, vv := range *employeeBenefitDetail {
				if v.IsLabel == nil || (v.IsLabel != nil && *v.IsLabel == false) {
					if vv.IsValue != nil && *vv.IsValue == true {
						f.SetCellValue(sheet, fmt.Sprintf("G%d", row), vv.Value)
						continue
					}
					f.SetCellStyle(sheet, fmt.Sprintf("G%d", row), fmt.Sprintf("G%d", row), styleCurrency)
					amount := 0.0
					if vv.Amount != nil && *vv.Amount != 0 {
						amount = *vv.Amount
					}
					f.SetCellValue(sheet, fmt.Sprintf("G%d", row), amount)
				}
			}
			row++
		}
		rowStart = row + 5
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
	period = datePeriod.Format("2006-01-02")
	fileName := fmt.Sprintf("EmployeeBenefit_%s.xlsx", period)
	fileLoc := fmt.Sprintf("assets/%d/%s", ctx.Auth.ID, fileName)
	err = f.SaveAs(fileLoc)
	if err != nil {
		return nil, err
	}

	result := &dto.EmployeeBenefitExportResponse{
		FileName: fileName,
		Path:     fileLoc,
	}
	return result, nil
}

func (s *service) GetVersion(ctx *abstraction.Context, payload *dto.GetVersionRequest) (*dto.GetVersionResponse, error) {
	filter := model.EmployeeBenefitFilterModel{
		CompanyCustomFilter: model.CompanyCustomFilter{
			CompanyID:          payload.CompanyID,
			ArrCompanyID:       payload.ArrCompanyID,
			ArrCompanyString:   payload.ArrCompanyString,
			ArrCompanyOperator: payload.ArrCompanyOperator,
		},
	}
	filter.Period = payload.Period
	// filter.ArrStatus = payload.ArrStatus
	data, err := s.Repository.GetVersion(ctx, &filter)
	if err != nil {
		return &dto.GetVersionResponse{}, helper.ErrorHandler(err)
	}
	result := &dto.GetVersionResponse{
		Data: *data,
	}
	return result, nil
}
