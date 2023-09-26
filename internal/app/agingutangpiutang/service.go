package agingutangpiutang

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
	Repository                 repository.AgingUtangPiutang
	Db                         *gorm.DB
	FormatterRepository        repository.Formatter
	AgingUPDetailRepository    repository.AgingUtangPiutangDetail
	ParameterRepository        repository.Parameter
	FormatterBridgesRepository repository.FormatterBridges
}

type Service interface {
	Find(ctx *abstraction.Context, payload *dto.AgingUtangPiutangGetRequest) (*dto.AgingUtangPiutangGetResponse, error)
	FindByID(ctx *abstraction.Context, payload *dto.AgingUtangPiutangGetByIDRequest) (*dto.AgingUtangPiutangGetByIDResponse, error)
	Create(ctx *abstraction.Context, payload *dto.AgingUtangPiutangCreateRequest) (*dto.AgingUtangPiutangCreateResponse, error)
	Update(ctx *abstraction.Context, payload *dto.AgingUtangPiutangUpdateRequest) (*dto.AgingUtangPiutangUpdateResponse, error)
	Delete(ctx *abstraction.Context, payload *dto.AgingUtangPiutangDeleteRequest) (*dto.AgingUtangPiutangDeleteResponse, error)
	Import(ctx *abstraction.Context, payload *dto.AgingUtangPiutangImportRequest, datas *[]model.AgingUtangPiutangDetailEntity) (*dto.AgingUtangPiutangImportResponse, error)
	Export(ctx *abstraction.Context, payload *dto.AgingUtangPiutangExportRequest) (*dto.AgingUtangPiutangExportResponse, error)
	GetVersion(ctx *abstraction.Context, payload *dto.GetVersionRequest) (*dto.GetVersionResponse, error)
}

func NewService(f *factory.Factory) *service {
	repository := f.AgingUtangPiutangRepository
	formatterRepository := f.FormatterRepository
	agingUPDetailRepository := f.AgingUtangPiutangDetailRepository
	parameterRepository := f.ParameterRepository
	formatterBridgesRepo := f.FormatterBridgesRepository
	db := f.Db
	return &service{
		Repository:                 repository,
		Db:                         db,
		FormatterRepository:        formatterRepository,
		AgingUPDetailRepository:    agingUPDetailRepository,
		ParameterRepository:        parameterRepository,
		FormatterBridgesRepository: formatterBridgesRepo,
	}
}

func (s *service) Find(ctx *abstraction.Context, payload *dto.AgingUtangPiutangGetRequest) (*dto.AgingUtangPiutangGetResponse, error) {
	data, info, err := s.Repository.Find(ctx, &payload.AgingUtangPiutangFilterModel, &payload.Pagination)
	if err != nil {
		return &dto.AgingUtangPiutangGetResponse{}, response.ErrorBuilder(&response.ErrorConstant.InternalServerError, err)
	}

	result := &dto.AgingUtangPiutangGetResponse{
		Datas:          *data,
		PaginationInfo: *info,
	}
	return result, nil
}

func (s *service) FindByID(ctx *abstraction.Context, payload *dto.AgingUtangPiutangGetByIDRequest) (*dto.AgingUtangPiutangGetByIDResponse, error) {
	data, err := s.Repository.FindByID(ctx, &payload.ID)
	if err != nil {
		return &dto.AgingUtangPiutangGetByIDResponse{}, response.ErrorBuilder(&response.ErrorConstant.InternalServerError, err)
	}
	allowed := helper.CompanyValidation(ctx.Auth.ID, data.CompanyID)
	if !allowed {
		return &dto.AgingUtangPiutangGetByIDResponse{}, response.ErrorBuilder(&response.ErrorConstant.UnauthorizedAccess, errors.New("Not Allowed"))
	}
	result := &dto.AgingUtangPiutangGetByIDResponse{
		AgingUtangPiutangEntityModel: *data,
	}
	return result, nil
}

func (s *service) Create(ctx *abstraction.Context, payload *dto.AgingUtangPiutangCreateRequest) (*dto.AgingUtangPiutangCreateResponse, error) {
	var data model.AgingUtangPiutangEntityModel

	allowed := helper.CompanyValidation(ctx.Auth.ID, payload.CompanyID)
	if !allowed {
		return &dto.AgingUtangPiutangCreateResponse{}, response.ErrorBuilder(&response.ErrorConstant.UnauthorizedAccess, errors.New("Not Allowed"))
	}

	if err := trxmanager.New(s.Db).WithTrx(ctx, func(ctx *abstraction.Context) error {
		data.Context = ctx
		data.AgingUtangPiutangEntity = payload.AgingUtangPiutangEntity
		data.AgingUtangPiutangEntity.Status = 1

		result, err := s.Repository.Create(ctx, &data)
		if err != nil {
			return response.ErrorBuilder(&response.ErrorConstant.UnprocessableEntity, err)
		}
		data = *result
		return nil
	}); err != nil {
		return &dto.AgingUtangPiutangCreateResponse{}, err
	}

	result := &dto.AgingUtangPiutangCreateResponse{
		AgingUtangPiutangEntityModel: data,
	}
	return result, nil
}

func (s *service) Update(ctx *abstraction.Context, payload *dto.AgingUtangPiutangUpdateRequest) (*dto.AgingUtangPiutangUpdateResponse, error) {
	var data model.AgingUtangPiutangEntityModel

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
		data.AgingUtangPiutangEntity = payload.AgingUtangPiutangEntity

		result, err := s.Repository.Update(ctx, &payload.ID, &data)
		if err != nil {
			return helper.ErrorHandler(err)
		}
		data = *result
		return nil
	}); err != nil {
		return &dto.AgingUtangPiutangUpdateResponse{}, err
	}
	result := &dto.AgingUtangPiutangUpdateResponse{
		AgingUtangPiutangEntityModel: data,
	}
	return result, nil
}

func (s *service) Delete(ctx *abstraction.Context, payload *dto.AgingUtangPiutangDeleteRequest) (*dto.AgingUtangPiutangDeleteResponse, error) {
	var data model.AgingUtangPiutangEntityModel
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
		return &dto.AgingUtangPiutangDeleteResponse{}, err
	}
	result := &dto.AgingUtangPiutangDeleteResponse{
		// AgingUtangPiutangEntityModel: data,
	}
	return result, nil
}

func (s *service) Export(ctx *abstraction.Context, payload *dto.AgingUtangPiutangExportRequest) (*dto.AgingUtangPiutangExportResponse, error) {
	f := excelize.NewFile()
	sheet := f.GetSheetName(f.GetActiveSheetIndex())
	f.SetColWidth(sheet, "A", "A", 8.29)
	f.SetColWidth(sheet, "B", "B", 26.43)
	f.SetColWidth(sheet, "C", "S", 17.86)

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

	f.SetColWidth(sheet, "K", "K", 2)

	agingutangpiutang, err := s.Repository.FindByID(ctx, &payload.AgingUtangPiutangID)
	if err != nil {
		return nil, err
	}

	if agingutangpiutang.ID == 0 {
		return nil, errors.New("Data Not Found")
	}

	allowed := helper.CompanyValidation(ctx.Auth.ID, agingutangpiutang.CompanyID)
	if !allowed {
		return nil, response.ErrorBuilder(&response.ErrorConstant.UnauthorizedAccess, errors.New("Not Allowed"))
	}

	formatterCode := []string{"AGING-UTANG-PIUTANG", "AGING-UTANG-PIUTANG-MUTASI-ECL"}
	formatterTitle := []string{"Detail Aging", "Mutasi ECL's"}
	row, rowStart := 5, 5
	for i, formatter := range formatterCode {
		f.SetCellValue(sheet, fmt.Sprintf("B%d", (rowStart-3)), formatterTitle[i])
		f.SetRowHeight(sheet, (rowStart - 2), 43.50)
		f.MergeCell(sheet, fmt.Sprintf("B%d", (rowStart-2)), fmt.Sprintf("B%d", (rowStart-1)))
		f.MergeCell(sheet, fmt.Sprintf("G%d", (rowStart-2)), fmt.Sprintf("G%d", (rowStart-1)))
		f.MergeCell(sheet, fmt.Sprintf("J%d", (rowStart-2)), fmt.Sprintf("J%d", (rowStart-1)))
		f.MergeCell(sheet, fmt.Sprintf("P%d", (rowStart-2)), fmt.Sprintf("P%d", (rowStart-1)))
		f.MergeCell(sheet, fmt.Sprintf("S%d", (rowStart-2)), fmt.Sprintf("S%d", (rowStart-1)))

		f.MergeCell(sheet, fmt.Sprintf("C%d", (rowStart-2)), fmt.Sprintf("D%d", (rowStart-2)))
		f.MergeCell(sheet, fmt.Sprintf("E%d", (rowStart-2)), fmt.Sprintf("F%d", (rowStart-2)))
		f.MergeCell(sheet, fmt.Sprintf("H%d", (rowStart-2)), fmt.Sprintf("I%d", (rowStart-2)))
		f.MergeCell(sheet, fmt.Sprintf("L%d", (rowStart-2)), fmt.Sprintf("M%d", (rowStart-2)))
		f.MergeCell(sheet, fmt.Sprintf("N%d", (rowStart-2)), fmt.Sprintf("O%d", (rowStart-2)))
		f.MergeCell(sheet, fmt.Sprintf("Q%d", (rowStart-2)), fmt.Sprintf("R%d", (rowStart-2)))
		f.SetCellStyle(sheet, fmt.Sprintf("B%d", (rowStart-2)), fmt.Sprintf("J%d", (rowStart-1)), styleHeader)
		f.SetCellStyle(sheet, fmt.Sprintf("L%d", (rowStart-2)), fmt.Sprintf("S%d", (rowStart-1)), styleHeader)

		f.SetCellValue(sheet, fmt.Sprintf("B%d", (rowStart-2)), "Description")
		f.SetCellValue(sheet, fmt.Sprintf("C%d", (rowStart-2)), "Piutang Usaha")
		f.SetCellValue(sheet, fmt.Sprintf("E%d", (rowStart-2)), "Piutang lain-lain jangka pendek")
		f.SetCellValue(sheet, fmt.Sprintf("G%d", (rowStart-2)), "Piutang pihak berelasi jangka pendek")
		f.SetCellValue(sheet, fmt.Sprintf("H%d", (rowStart-2)), "Piutang lain-lain jangka panjang")
		f.SetCellValue(sheet, fmt.Sprintf("J%d", (rowStart-2)), "Piutang pihak berelasi jangka panjang")
		f.SetCellValue(sheet, fmt.Sprintf("L%d", (rowStart-2)), "Utang usaha")
		f.SetCellValue(sheet, fmt.Sprintf("N%d", (rowStart-2)), "Utang lain-lain jangka pendek")
		f.SetCellValue(sheet, fmt.Sprintf("P%d", (rowStart-2)), "Utang pihak berelasi jangka pendek")
		f.SetCellValue(sheet, fmt.Sprintf("Q%d", (rowStart-2)), "Utang lain-lain jangka panjang")
		f.SetCellValue(sheet, fmt.Sprintf("S%d", (rowStart-2)), "Utang pihak berelasi jangka panjang")

		header1 := []string{"C", "D", "E", "F", "H", "I", "L", "M", "N", "O", "Q", "R"}
		for i, v := range header1 {
			if (i+1)%2 == 0 {
				f.SetCellValue(sheet, fmt.Sprintf("%s%d", v, (rowStart-1)), "Pihak Berelasi")
			} else {
				f.SetCellValue(sheet, fmt.Sprintf("%s%d", v, (rowStart-1)), "Pihak Ketiga")
			}
		}

		var criteria dto.FormatterGetRequest
		criteria.FormatterFilterModel.FormatterFor = &formatter

		data, err := s.FormatterRepository.FindWithDetail(ctx, &criteria.FormatterFilterModel)
		if err != nil {
			return &dto.AgingUtangPiutangExportResponse{}, response.ErrorBuilder(&response.ErrorConstant.InternalServerError, err)
		}

		tmpStr := "AGING-UTANG-PIUTANG"
		criteriaBridge := model.FormatterBridgesFilterModel{}
		criteriaBridge.FormatterBridgesFilter.Source = &tmpStr
		criteriaBridge.FormatterBridgesFilter.FormatterID = &data.ID
		criteriaBridge.FormatterBridgesFilter.TrxRefID = &agingutangpiutang.ID

		bridges, err := s.FormatterBridgesRepository.FindWithCriteria(ctx, &criteriaBridge)
		if err != nil {
			return nil, err
		}
		rowCode := make(map[string]int)
		partRowStart := row
		for _, v := range data.FormatterDetail {
			rowCode[v.Code] = row
			f.SetCellStyle(sheet, fmt.Sprintf("B%d", row), fmt.Sprintf("B%d", row), styleLabel)
			f.SetCellStyle(sheet, fmt.Sprintf("C%d", row), fmt.Sprintf("J%d", row), styleCurrency)
			f.SetCellStyle(sheet, fmt.Sprintf("L%d", row), fmt.Sprintf("S%d", row), styleCurrency)

			if strings.Contains(strings.ToLower(v.Code), "blank") {
				row++
				continue
			}

			f.SetCellValue(sheet, fmt.Sprintf("B%d", row), v.Description)
			if v.IsTotal != nil && *v.IsTotal == true {
				f.SetCellStyle(sheet, fmt.Sprintf("B%d", row), fmt.Sprintf("B%d", row), styleLabelTotal)
				f.SetCellStyle(sheet, fmt.Sprintf("C%d", row), fmt.Sprintf("J%d", row), styleCurrencyTotal)
				f.SetCellStyle(sheet, fmt.Sprintf("L%d", row), fmt.Sprintf("S%d", row), styleCurrencyTotal)
				if v.FxSummary == "" {
					row++
					continue
				}
				for chr := 'C'; chr <= 'S'; chr++ {
					formula := v.FxSummary
					if chr == 'K' {
						continue
					}
					reg := regexp.MustCompile(`[A-Za-z_~]+`)
					match := reg.FindAllString(formula, -1)
					for _, vMatch := range match {
						//cari jml berdasarkan code
						if _, ok := rowCode[vMatch]; ok {
							formula = strings.ReplaceAll(formula, vMatch, fmt.Sprintf("%c%d", chr, rowCode[vMatch]))
						}

					}
					f.SetCellFormula(sheet, fmt.Sprintf("%c%d", chr, row), fmt.Sprintf("=%s", formula))
				}
				row++
				continue
			}

			if v.AutoSummary != nil && *v.AutoSummary == true {
				f.SetCellStyle(sheet, fmt.Sprintf("B%d", row), fmt.Sprintf("B%d", row), styleLabelTotal)
				f.SetCellStyle(sheet, fmt.Sprintf("C%d", row), fmt.Sprintf("J%d", row), styleCurrencyTotal)
				f.SetCellStyle(sheet, fmt.Sprintf("L%d", row), fmt.Sprintf("S%d", row), styleCurrencyTotal)

				for chr := 'C'; chr <= 'S'; chr++ {
					if chr == 'K' {
						continue
					}
					f.SetCellFormula(sheet, fmt.Sprintf("%c%d", chr, row), fmt.Sprintf("=SUM(%c%d:%c%d)", chr, partRowStart, chr, row-1))
				}
				row++
				partRowStart = row
				continue
			}

			criteriaAUP := model.AgingUtangPiutangDetailFilterModel{}
			criteriaAUP.Code = &v.Code
			criteriaAUP.FormatterBridgesID = &bridges.ID
			criteriaAUP.AgingUtangPiutangID = &agingutangpiutang.ID

			paginationAUP := abstraction.Pagination{}
			pagesize := 10000
			paginationAUP.PageSize = &pagesize

			AgingUPDetail, _, err := s.AgingUPDetailRepository.Find(ctx, &criteriaAUP, &paginationAUP)
			if err != nil && v.Code != "" {
				continue
			}
			if len(*AgingUPDetail) == 0 {
				continue
			}
			for _, vv := range *AgingUPDetail {
				if vv.Piutangusaha3rdparty != nil && *vv.Piutangusaha3rdparty != 0 {
					f.SetCellValue(sheet, fmt.Sprintf("C%d", row), *vv.Piutangusaha3rdparty)
				}
				if vv.PiutangusahaBerelasi != nil && *vv.PiutangusahaBerelasi != 0 {
					f.SetCellValue(sheet, fmt.Sprintf("D%d", row), *vv.PiutangusahaBerelasi)
				}
				if vv.Piutanglainshortterm3rdparty != nil && *vv.Piutanglainshortterm3rdparty != 0 {
					f.SetCellValue(sheet, fmt.Sprintf("E%d", row), *vv.Piutanglainshortterm3rdparty)
				}
				if vv.PiutanglainshorttermBerelasi != nil && *vv.PiutanglainshorttermBerelasi != 0 {
					f.SetCellValue(sheet, fmt.Sprintf("F%d", row), *vv.PiutanglainshorttermBerelasi)
				}
				if vv.Piutangberelasishortterm != nil && *vv.Piutangberelasishortterm != 0 {
					f.SetCellValue(sheet, fmt.Sprintf("G%d", row), *vv.Piutangberelasishortterm)
				}
				if vv.Piutanglainlongterm3rdparty != nil && *vv.Piutanglainlongterm3rdparty != 0 {
					f.SetCellValue(sheet, fmt.Sprintf("H%d", row), *vv.Piutanglainlongterm3rdparty)
				}
				if vv.PiutanglainlongtermBerelasi != nil && *vv.PiutanglainlongtermBerelasi != 0 {
					f.SetCellValue(sheet, fmt.Sprintf("I%d", row), *vv.PiutanglainlongtermBerelasi)
				}
				if vv.Piutangberelasilongterm != nil && *vv.Piutangberelasilongterm != 0 {
					f.SetCellValue(sheet, fmt.Sprintf("J%d", row), *vv.Piutangberelasilongterm)
				}
				if vv.Utangusaha3rdparty != nil && *vv.Utangusaha3rdparty != 0 {
					f.SetCellValue(sheet, fmt.Sprintf("L%d", row), *vv.Utangusaha3rdparty)
				}
				if vv.UtangusahaBerelasi != nil && *vv.UtangusahaBerelasi != 0 {
					f.SetCellValue(sheet, fmt.Sprintf("M%d", row), *vv.UtangusahaBerelasi)
				}
				if vv.Utanglainshortterm3rdparty != nil && *vv.Utanglainshortterm3rdparty != 0 {
					f.SetCellValue(sheet, fmt.Sprintf("N%d", row), *vv.Utanglainshortterm3rdparty)
				}
				if vv.UtanglainshorttermBerelasi != nil && *vv.UtanglainshorttermBerelasi != 0 {
					f.SetCellValue(sheet, fmt.Sprintf("O%d", row), *vv.UtanglainshorttermBerelasi)
				}
				if vv.Utangberelasishortterm != nil && *vv.Utangberelasishortterm != 0 {
					f.SetCellValue(sheet, fmt.Sprintf("P%d", row), *vv.Utangberelasishortterm)
				}
				if vv.Utanglainlongterm3rdparty != nil && *vv.Utanglainlongterm3rdparty != 0 {
					f.SetCellValue(sheet, fmt.Sprintf("Q%d", row), *vv.Utanglainlongterm3rdparty)
				}
				if vv.UtanglainlongtermBerelasi != nil && *vv.UtanglainlongtermBerelasi != 0 {
					f.SetCellValue(sheet, fmt.Sprintf("R%d", row), *vv.UtanglainlongtermBerelasi)
				}
				if vv.Utangberelasilongterm != nil && *vv.Utangberelasilongterm != 0 {
					f.SetCellValue(sheet, fmt.Sprintf("S%d", row), *vv.Utangberelasilongterm)
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
	datePeriod, err := time.Parse(time.RFC3339, agingutangpiutang.Period)
	if err != nil {
		return nil, err
	}
	period := datePeriod.Format("2006-01-02")
	fileName := fmt.Sprintf("AgingUtangPiutang_%s.xlsx", period)
	fileLoc := fmt.Sprintf("assets/%d/%s", ctx.Auth.ID, fileName)
	err = f.SaveAs(fileLoc)
	if err != nil {
		return nil, err
	}

	result := &dto.AgingUtangPiutangExportResponse{
		FileName: fileName,
		Path:     fileLoc,
	}
	return result, nil
}

func (s *service) GetVersion(ctx *abstraction.Context, payload *dto.GetVersionRequest) (*dto.GetVersionResponse, error) {
	filter := model.AgingUtangPiutangFilterModel{
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
		return &dto.GetVersionResponse{}, response.ErrorBuilder(&response.ErrorConstant.InternalServerError, err)
	}
	result := &dto.GetVersionResponse{
		Data: *data,
	}
	return result, nil
}
