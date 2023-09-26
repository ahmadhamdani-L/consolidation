package mutasidta

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
	Repository                 repository.MutasiDta
	MutasiDtaDetailRepository  repository.MutasiDtaDetail
	FormatterRepository        repository.Formatter
	ParameterRepository        repository.Parameter
	FormatterBridgesRepository repository.FormatterBridges
	Db                         *gorm.DB
}

type Service interface {
	Find(ctx *abstraction.Context, payload *dto.MutasiDtaGetRequest) (*dto.MutasiDtaGetResponse, error)
	FindByID(ctx *abstraction.Context, payload *dto.MutasiDtaGetByIDRequest) (*dto.MutasiDtaGetByIDResponse, error)
	Create(ctx *abstraction.Context, payload *dto.MutasiDtaCreateRequest) (*dto.MutasiDtaCreateResponse, error)
	Update(ctx *abstraction.Context, payload *dto.MutasiDtaUpdateRequest) (*dto.MutasiDtaUpdateResponse, error)
	Delete(ctx *abstraction.Context, payload *dto.MutasiDtaDeleteRequest) (*dto.MutasiDtaDeleteResponse, error)
	Export(ctx *abstraction.Context, payload *dto.MutasiDtaExportRequest) (*dto.MutasiDtaExportResponse, error)
	GetVersion(ctx *abstraction.Context, payload *dto.GetVersionRequest) (*dto.GetVersionResponse, error)
}

func NewService(f *factory.Factory) *service {
	repository := f.MutasiDtaRepository
	mutasiDtaDetail := f.MutasiDtaDetailRepository
	formatter := f.FormatterRepository
	parameter := f.ParameterRepository
	formaterBridgesRepo := f.FormatterBridgesRepository
	db := f.Db
	return &service{
		Repository:                 repository,
		MutasiDtaDetailRepository:  mutasiDtaDetail,
		FormatterRepository:        formatter,
		ParameterRepository:        parameter,
		FormatterBridgesRepository: formaterBridgesRepo,
		Db:                         db,
	}
}

func (s *service) Find(ctx *abstraction.Context, payload *dto.MutasiDtaGetRequest) (*dto.MutasiDtaGetResponse, error) {
	data, info, err := s.Repository.Find(ctx, &payload.MutasiDtaFilterModel, &payload.Pagination)
	if err != nil {
		return &dto.MutasiDtaGetResponse{}, helper.ErrorHandler(err)
	}

	result := &dto.MutasiDtaGetResponse{
		Datas:          *data,
		PaginationInfo: *info,
	}
	return result, nil
}

func (s *service) FindByID(ctx *abstraction.Context, payload *dto.MutasiDtaGetByIDRequest) (*dto.MutasiDtaGetByIDResponse, error) {
	data, err := s.Repository.FindByID(ctx, &payload.ID)
	if err != nil {
		return &dto.MutasiDtaGetByIDResponse{}, helper.ErrorHandler(err)
	}
	allowed := helper.CompanyValidation(ctx.Auth.ID, data.CompanyID)
	if !allowed {
		return nil, response.ErrorBuilder(&response.ErrorConstant.UnauthorizedAccess, errors.New("Not Allowed"))
	}
	result := &dto.MutasiDtaGetByIDResponse{
		MutasiDtaEntityModel: *data,
	}
	return result, nil
}

func (s *service) Create(ctx *abstraction.Context, payload *dto.MutasiDtaCreateRequest) (*dto.MutasiDtaCreateResponse, error) {
	var data model.MutasiDtaEntityModel

	allowed := helper.CompanyValidation(ctx.Auth.ID, payload.CompanyID)
	if !allowed {
		return nil, response.ErrorBuilder(&response.ErrorConstant.UnauthorizedAccess, errors.New("Not Allowed"))
	}

	if err := trxmanager.New(s.Db).WithTrx(ctx, func(ctx *abstraction.Context) error {
		data.Context = ctx
		data.MutasiDtaEntity = payload.MutasiDtaEntity
		data.MutasiDtaEntity.Status = 1

		result, err := s.Repository.Create(ctx, &data)
		if err != nil {
			return helper.ErrorHandler(err)
		}
		data = *result
		return nil
	}); err != nil {
		return &dto.MutasiDtaCreateResponse{}, err
	}

	result := &dto.MutasiDtaCreateResponse{
		MutasiDtaEntityModel: data,
	}
	return result, nil
}

func (s *service) Update(ctx *abstraction.Context, payload *dto.MutasiDtaUpdateRequest) (*dto.MutasiDtaUpdateResponse, error) {
	var data model.MutasiDtaEntityModel
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
		data.MutasiDtaEntity = payload.MutasiDtaEntity

		result, err := s.Repository.Update(ctx, &payload.ID, &data)
		if err != nil {
			return helper.ErrorHandler(err)
		}
		data = *result
		return nil
	}); err != nil {
		return &dto.MutasiDtaUpdateResponse{}, err
	}
	result := &dto.MutasiDtaUpdateResponse{
		MutasiDtaEntityModel: data,
	}
	return result, nil
}

func (s *service) Delete(ctx *abstraction.Context, payload *dto.MutasiDtaDeleteRequest) (*dto.MutasiDtaDeleteResponse, error) {
	var data model.MutasiDtaEntityModel
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
		return &dto.MutasiDtaDeleteResponse{}, helper.ErrorHandler(err)
	}
	result := &dto.MutasiDtaDeleteResponse{
		// MutasiDtaEntityModel: data,
	}
	return result, nil
}

func (s *service) Export(ctx *abstraction.Context, payload *dto.MutasiDtaExportRequest) (*dto.MutasiDtaExportResponse, error) {
	f := excelize.NewFile()
	sheet := f.GetSheetName(f.GetActiveSheetIndex())

	arrStyleColWidth := []map[string]interface{}{
		{"COL": "A", "WIDTH": 8.21},
		{"COL": "B", "WIDTH": 3.50},
		{"COL": "C", "WIDTH": 33.39},
		{"COL": "D", "WIDTH": 15.35},
		{"COL": "E", "WIDTH": 15.35},
		{"COL": "F", "WIDTH": 13.74},
		{"COL": "G", "WIDTH": 13.74},
		{"COL": "H", "WIDTH": 13.74},
		{"COL": "I", "WIDTH": 15.71},
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
			Color:   []string{"#fac090"},
		},
		Alignment: &excelize.Alignment{
			// WrapText:   true,
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
	})
	if err != nil {
		return nil, err
	}

	styleLabel, err := f.NewStyle(&excelize.Style{
		Border: []excelize.Border{
			// {Type: "top", Color: "000000", Style: 1},
			// {Type: "bottom", Color: "000000", Style: 1},
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

	formatterCode := []string{"MUTASI-DTA"}
	formatterTitle := []string{"Mutasi DTA"}
	row, rowStart := 7, 7

	f.SetCellValue(sheet, "B2", formatterTitle[0])

	f.MergeCell(sheet, "B4", "B6")
	f.MergeCell(sheet, "C4", "C6")
	f.MergeCell(sheet, "D4", "D6")
	f.MergeCell(sheet, "E4", "G4")
	f.MergeCell(sheet, "E5", "E6")
	f.MergeCell(sheet, "F5", "F6")
	f.MergeCell(sheet, "G5", "G6")
	f.MergeCell(sheet, "H4", "I4")
	f.MergeCell(sheet, "H5", "H6")
	f.MergeCell(sheet, "I5", "I6")
	f.MergeCell(sheet, "J4", "J6")

	// criteria := model.MutasiDtaFilterModel{}
	// criteria.CompanyID = &payload.CompanyID
	// version := payload.Versions
	// criteria.Versions = &version
	// criteria.Period = &payload.Period
	mutasidta, err := s.Repository.FindByID(ctx, &payload.MutasiDtaID)
	if err != nil {
		return nil, helper.ErrorHandler(err)
	}

	allowed := helper.CompanyValidation(ctx.Auth.ID, mutasidta.CompanyID)
	if !allowed {
		return nil, response.ErrorBuilder(&response.ErrorConstant.UnauthorizedAccess, errors.New("Not Allowed"))
	}

	datePeriod, err := time.Parse(time.RFC3339, mutasidta.Period)
	if err != nil {
		return nil, err
	}

	f.SetCellStyle(sheet, "B4", "J6", styleHeader)
	f.SetCellValue(sheet, "B4", "NO")
	f.SetCellValue(sheet, "C4", "Description")
	// f.SetCellValue(sheet, "D4", "PT xxx Saldo Awal 01.01.21")
	f.SetCellValue(sheet, "D4", fmt.Sprintf("%s Saldo Awal %s", mutasidta.Company.Name, datePeriod.Format("02.01.06")))
	f.SetCellValue(sheet, "E4", "Penambahan (Pengurangan)")
	f.SetCellValue(sheet, "E5", "Manfaat (beban) pajak tangguhan")
	f.SetCellValue(sheet, "F5", "OCI")
	f.SetCellValue(sheet, "G5", "Akuisisi Entitas anak")
	f.SetCellValue(sheet, "H4", "Dampak perubahan tariff pajak")
	f.SetCellValue(sheet, "H5", "Dibebankan ke laba rugi")
	f.SetCellValue(sheet, "I5", "Dibebankan ke OCI")
	// f.SetCellValue(sheet, "J4", "PT xxx\nSaldo Akhir\n31.12.21")
	f.SetCellValue(sheet, "J4", fmt.Sprintf("%s Saldo Akhir %s", mutasidta.Company.Name, datePeriod.Format("02.01.06")))

	for _, formatter := range formatterCode {

		var criteria dto.FormatterGetRequest
		criteria.FormatterFilterModel.FormatterFor = &formatter

		data, err := s.FormatterRepository.FindWithDetail(ctx, &criteria.FormatterFilterModel)
		if err != nil {
			return &dto.MutasiDtaExportResponse{}, helper.ErrorHandler(err)
		}

		tmpStr := "MUTASI-DTA"
		criteriaBridge := model.FormatterBridgesFilterModel{}
		criteriaBridge.FormatterBridgesFilter.Source = &tmpStr
		criteriaBridge.FormatterBridgesFilter.FormatterID = &data.ID
		criteriaBridge.FormatterBridgesFilter.TrxRefID = &mutasidta.ID

		bridges, err := s.FormatterBridgesRepository.FindWithCriteria(ctx, &criteriaBridge)
		if err != nil {
			return nil, helper.ErrorHandler(err)
		}
		rowCode := make(map[string]int)
		partRowStart := row
		for _, v := range data.FormatterDetail {
			rowCode[v.Code] = row
			f.SetCellStyle(sheet, fmt.Sprintf("B%d", row), fmt.Sprintf("C%d", row), styleLabel)
			f.SetCellStyle(sheet, fmt.Sprintf("D%d", row), fmt.Sprintf("J%d", row), styleCurrency)

			if strings.Contains(strings.ToLower(v.Code), "blank") {
				row++
				continue
			}

			f.SetCellValue(sheet, fmt.Sprintf("B%d", row), v.Description)

			if v.IsTotal != nil && *v.IsTotal == true {
				f.SetCellStyle(sheet, fmt.Sprintf("B%d", row), fmt.Sprintf("C%d", row), styleLabelTotal)
				f.SetCellStyle(sheet, fmt.Sprintf("D%d", row), fmt.Sprintf("J%d", row), styleCurrencyTotal)
				f.MergeCell(sheet, fmt.Sprintf("B%d", row), fmt.Sprintf("C%d", row))
				if v.FxSummary == "" {
					row++
					continue
				}
				for chr := 'D'; chr <= 'J'; chr++ {
					formula := v.FxSummary
					reg := regexp.MustCompile(`[A-Za-z_~]+`)
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

			if v.AutoSummary != nil && *v.AutoSummary == true {
				f.SetCellStyle(sheet, fmt.Sprintf("B%d", row), fmt.Sprintf("C%d", row), styleLabelTotal)
				f.SetCellStyle(sheet, fmt.Sprintf("D%d", row), fmt.Sprintf("J%d", row), styleCurrencyTotal)
				f.MergeCell(sheet, fmt.Sprintf("B%d", row), fmt.Sprintf("C%d", row))

				for chr := 'D'; chr <= 'J'; chr++ {
					f.SetCellFormula(sheet, fmt.Sprintf("%c%d", chr, row), fmt.Sprintf("=SUM(%c%d:%c%d)", chr, partRowStart, chr, row-1))
				}
				row++
				partRowStart = row
				continue
			}

			criteriaMD := model.MutasiDtaDetailFilterModel{}
			criteriaMD.Code = &v.Code
			criteriaMD.FormatterBridgesID = &bridges.ID
			criteriaMD.MutasiDtaID = &mutasidta.ID

			paginationMD := abstraction.Pagination{}
			pagesize := 1
			// sortBy := "id"
			// sort := "asc"
			paginationMD.PageSize = &pagesize
			// paginationMD.SortBy = &sortBy
			// paginationMD.Sort = &sort

			mDtaDetail, _, err := s.MutasiDtaDetailRepository.Find(ctx, &criteriaMD, &paginationMD)
			if err != nil && v.Code != "" {
				continue
			}
			if len(*mDtaDetail) == 0 {
				continue
			}
			for _, vv := range *mDtaDetail {
				f.SetCellValue(sheet, fmt.Sprintf("B%d", row), vv.SortID)
				f.SetCellValue(sheet, fmt.Sprintf("C%d", row), vv.Description)
				if vv.SaldoAwal != nil && *vv.SaldoAwal != 0 {
					f.SetCellValue(sheet, fmt.Sprintf("D%d", row), *vv.SaldoAwal)
				}
				if vv.ManfaatBebanPajak != nil && *vv.ManfaatBebanPajak != 0 {
					f.SetCellValue(sheet, fmt.Sprintf("E%d", row), *vv.ManfaatBebanPajak)
				}
				if vv.Oci != nil && *vv.Oci != 0 {
					f.SetCellValue(sheet, fmt.Sprintf("F%d", row), *vv.Oci)
				}
				if vv.AkuisisiEntitasAnak != nil && *vv.AkuisisiEntitasAnak != 0 {
					f.SetCellValue(sheet, fmt.Sprintf("G%d", row), *vv.AkuisisiEntitasAnak)
				}
				if vv.DibebankanKeLr != nil && *vv.DibebankanKeLr != 0 {
					f.SetCellValue(sheet, fmt.Sprintf("H%d", row), *vv.DibebankanKeLr)
				}
				if vv.DibebankanKeOci != nil && *vv.DibebankanKeOci != 0 {
					f.SetCellValue(sheet, fmt.Sprintf("I%d", row), *vv.DibebankanKeOci)
				}
				f.SetCellFormula(sheet, fmt.Sprintf("J%d", row), fmt.Sprintf("=SUM(D%d:I%d)", row, row))
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
				log.Println(err)
				return nil, err
			}
		} else {
			return nil, err
		}
	}

	period := datePeriod.Format("2006-01-02")
	fileName := fmt.Sprintf("MutasiDta_%s.xlsx", period)
	fileLoc := fmt.Sprintf("assets/%d/%s", ctx.Auth.ID, fileName)
	err = f.SaveAs(fileLoc)
	if err != nil {
		return nil, err
	}

	result := &dto.MutasiDtaExportResponse{
		FileName: fileName,
		Path:     fileLoc,
	}
	return result, nil
}

func (s *service) GetVersion(ctx *abstraction.Context, payload *dto.GetVersionRequest) (*dto.GetVersionResponse, error) {
	filter := model.MutasiDtaFilterModel{
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
