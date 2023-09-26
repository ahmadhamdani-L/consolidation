package trialbalance

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
	Repository                 repository.TrialBalance
	TBDetailRepository         repository.TrialBalanceDetail
	FormatterRepository        repository.Formatter
	CoaRepository              repository.Coa
	FormatterBridgesRepository repository.FormatterBridges
	FormatterDetailRepository  repository.FormatterDetail
	AjeRepository              repository.Adjustment
	Db                         *gorm.DB
}

type Service interface {
	Find(ctx *abstraction.Context, payload *dto.TrialBalanceGetRequest) (*dto.TrialBalanceGetResponse, error)
	FindByID(ctx *abstraction.Context, payload *dto.TrialBalanceGetByIDRequest) (*dto.TrialBalanceGetByIDResponse, error)
	Create(ctx *abstraction.Context, payload *dto.TrialBalanceCreateRequest) (*dto.TrialBalanceCreateResponse, error)
	Update(ctx *abstraction.Context, payload *dto.TrialBalanceUpdateRequest) (*dto.TrialBalanceUpdateResponse, error)
	Delete(ctx *abstraction.Context, payload *dto.TrialBalanceDeleteRequest) (*dto.TrialBalanceDeleteResponse, error)
	Export(ctx *abstraction.Context, payload *dto.TrialBalanceExportRequest) (*dto.TrialBalanceExportResponse, error)
	GetVersion(ctx *abstraction.Context, payload *dto.GetVersionRequest) (*dto.GetVersionResponse, error)
}

func NewService(f *factory.Factory) *service {
	repository := f.TrialBalanceRepository
	TBDetailRepository := f.TrialBalanceDetailRepository
	FormatterRepository := f.FormatterRepository
	FormatterDetailRepository := f.FormatterDetailRepository
	CoaRepository := f.CoaRepository
	formatterBridgesRepo := f.FormatterBridgesRepository
	ajeRepo := f.AdjustmentRepository
	db := f.Db
	return &service{
		Repository:                 repository,
		TBDetailRepository:         TBDetailRepository,
		FormatterRepository:        FormatterRepository,
		FormatterDetailRepository:  FormatterDetailRepository,
		CoaRepository:              CoaRepository,
		FormatterBridgesRepository: formatterBridgesRepo,
		AjeRepository:              ajeRepo,
		Db:                         db,
	}
}

func (s *service) Find(ctx *abstraction.Context, payload *dto.TrialBalanceGetRequest) (*dto.TrialBalanceGetResponse, error) {
	data, info, err := s.Repository.Find(ctx, &payload.TrialBalanceFilterModel, &payload.Pagination)
	if err != nil {
		return &dto.TrialBalanceGetResponse{}, helper.ErrorHandler(err)
	}
	result := &dto.TrialBalanceGetResponse{
		Datas:          *data,
		PaginationInfo: *info,
	}
	return result, nil
}

func (s *service) FindByID(ctx *abstraction.Context, payload *dto.TrialBalanceGetByIDRequest) (*dto.TrialBalanceGetByIDResponse, error) {
	data, err := s.Repository.FindByID(ctx, &payload.ID)
	if err != nil {
		return &dto.TrialBalanceGetByIDResponse{}, helper.ErrorHandler(err)
	}
	allowed := helper.CompanyValidation(ctx.Auth.ID, data.CompanyID)
	if !allowed {
		return nil, response.ErrorBuilder(&response.ErrorConstant.UnauthorizedAccess, errors.New("Not Allowed"))
	}
	result := &dto.TrialBalanceGetByIDResponse{
		TrialBalanceEntityModel: *data,
	}
	return result, nil
}

func (s *service) Create(ctx *abstraction.Context, payload *dto.TrialBalanceCreateRequest) (*dto.TrialBalanceCreateResponse, error) {
	var data model.TrialBalanceEntityModel

	allowed := helper.CompanyValidation(ctx.Auth.ID, payload.CompanyID)
	if !allowed {
		return nil, response.ErrorBuilder(&response.ErrorConstant.UnauthorizedAccess, errors.New("Not Allowed"))
	}

	if err := trxmanager.New(s.Db).WithTrx(ctx, func(ctx *abstraction.Context) error {
		data.Context = ctx
		data.TrialBalanceEntity = payload.TrialBalanceEntity
		data.TrialBalanceEntity.Status = 1
		result, err := s.Repository.Create(ctx, &data)
		if err != nil {
			return helper.ErrorHandler(err)
		}
		data = *result
		return nil
	}); err != nil {
		return &dto.TrialBalanceCreateResponse{}, err
	}
	result := &dto.TrialBalanceCreateResponse{
		TrialBalanceEntityModel: data,
	}
	return result, nil
}

func (s *service) Update(ctx *abstraction.Context, payload *dto.TrialBalanceUpdateRequest) (*dto.TrialBalanceUpdateResponse, error) {
	var data model.TrialBalanceEntityModel

	if err := trxmanager.New(s.Db).WithTrx(ctx, func(ctx *abstraction.Context) error {
		existing, err := s.Repository.FindByID(ctx, &payload.ID)
		if err != nil {
			return helper.ErrorHandler(err)
		}

		if existing.Status != 1 {
			return response.ErrorBuilder(&response.ErrorConstant.DataValidated, errors.New("Cannot update data that has been validated"))
		}

		if existing.ValidationNote != "" {
			return response.ErrorBuilder(&response.ErrorConstant.DataValidated, errors.New("Cannot update data because the data is on progress validation"))
		}

		allowed := helper.CompanyValidation(ctx.Auth.ID, existing.CompanyID)
		allowed2 := helper.CompanyValidation(ctx.Auth.ID, payload.CompanyID)
		if !allowed || !allowed2 {
			return response.ErrorBuilder(&response.ErrorConstant.UnauthorizedAccess, errors.New("Not Allowed"))
		}

		data.Context = ctx
		data.TrialBalanceEntity = payload.TrialBalanceEntity
		result, err := s.Repository.Update(ctx, &payload.ID, &data)
		if err != nil {
			return helper.ErrorHandler(err)
		}
		data = *result
		return nil
	}); err != nil {
		return &dto.TrialBalanceUpdateResponse{}, err
	}
	result := &dto.TrialBalanceUpdateResponse{
		TrialBalanceEntityModel: data,
	}
	return result, nil
}

func (s *service) Delete(ctx *abstraction.Context, payload *dto.TrialBalanceDeleteRequest) (*dto.TrialBalanceDeleteResponse, error) {
	var data model.TrialBalanceEntityModel
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
		return &dto.TrialBalanceDeleteResponse{}, err
	}
	result := &dto.TrialBalanceDeleteResponse{
		TrialBalanceEntityModel: data,
	}
	return result, nil
}

// func (s *service) Import(ctx *abstraction.Context, payload *dto.TrialBalanceImportRequest, datas *[]model.TrialBalanceDetailEntity) (*dto.TrialBalanceImportResponse, error) {
// 	var dataTB model.TrialBalanceEntityModel
// 	currentYear, currentMonth, _ := time.Now().Date()
// 	currentLocation := time.Now().Location()
// 	firstOfMonth := time.Date(currentYear, currentMonth, 1, 0, 0, 0, 0, currentLocation)
// 	lastOfMonth := firstOfMonth.AddDate(0, 1, -1)
// 	period := lastOfMonth.Format("2006-01-02")

// 	if err := trxmanager.New(s.Db).WithTrx(ctx, func(ctx *abstraction.Context) error {
// 		//cek company berdasarkan user
// 		//belum ada
// 		//skip

// 		criteriaFormatter := model.FormatterFilterModel{}
// 		tmpStr := "TRIAL-BALANCE"
// 		criteriaFormatter.FormatterFor = &tmpStr

// 		getFormatter, _, err := s.FormatterRepository.Find(ctx, &criteriaFormatter, &abstraction.Pagination{})
// 		if err != nil {
// 			return helper.ErrorHandler(err)
// 		}

// 		formatterID := 0
// 		for _, tmpFormatter := range *getFormatter {
// 			formatterID = tmpFormatter.ID
// 		}
// 		if formatterID == 0 {
// 			fmt.Println("No Formatter found")
// 			return helper.ErrorHandler(err)
// 		}

// 		criteriaTB := model.TrialBalanceFilterModel{}
// 		criteriaTB.Period = &period
// 		criteriaTB.FormatterID = &formatterID

// 		getTrialBalance, err := s.Repository.GetCount(ctx, &criteriaTB)
// 		if err != nil {
// 			return helper.ErrorHandler(err)
// 		}
// 		version := *getTrialBalance + 1

// 		dataTB.Context = ctx
// 		dataTB.TrialBalanceEntity = model.TrialBalanceEntity{
// 			Versions:    int(version),
// 			Period:      period,
// 			FormatterID: formatterID,
// 			CompanyID:   1,
// 			Status:      1,
// 		}
// 		resultTB, err := s.Repository.Create(ctx, &dataTB)
// 		if err != nil {
// 			return helper.ErrorHandler(err)
// 		}

// 		var arrDataTBD []model.TrialBalanceDetailEntityModel
// 		for _, v := range *datas {
// 			dataTBD := model.TrialBalanceDetailEntityModel{
// 				Context:                  ctx,
// 				TrialBalanceDetailEntity: v,
// 			}
// 			dataTBD.TrialBalanceID = resultTB.ID
// 			arrDataTBD = append(arrDataTBD, dataTBD)
// 		}
// 		_, err = s.TBDetailRepository.Import(ctx, &arrDataTBD)
// 		if err != nil {
// 			return helper.ErrorHandler(err)
// 		}

// 		dataTB = *resultTB
// 		return nil
// 	}); err != nil {
// 		return &dto.TrialBalanceImportResponse{}, err
// 	}
// 	result := &dto.TrialBalanceImportResponse{
// 		Data: dataTB,
// 	}

// 	return result, nil
// }

func (s *service) Export(ctx *abstraction.Context, payload *dto.TrialBalanceExportRequest) (*dto.TrialBalanceExportResponse, error) {
	var (
		criteriaFormatter model.FormatterDetailFilterModel
	)

	tb, err := s.Repository.Get(ctx, payload.TrialBalanceID)
	if err != nil {
		return nil, err
	}

	if tb.ID == 0 {
		return nil, errors.New("no data found")
	}

	allowed := helper.CompanyValidation(ctx.Auth.ID, tb.CompanyID)
	if !allowed {
		return nil, response.ErrorBuilder(&response.ErrorConstant.UnauthorizedAccess, errors.New("Not Allowed"))
	}

	datePeriod, err := time.Parse(time.RFC3339, tb.Period)
	if err != nil {
		return nil, err
	}

	if len(tb.FormatterBridges) == 0 {
		return nil, err
	}
	var formatterID int
	for _, fmtbridges := range tb.FormatterBridges {
		formatterID = fmtbridges.FormatterID
	}
	criteriaFormatter.FormatterID = &formatterID
	t := true
	criteriaFormatter.IsShowExport = &t
	pagesize := 100000
	tmpStr := "sort_id"
	tmpStr1 := "ASC"
	paginationTB := abstraction.Pagination{
		PageSize: &pagesize,
		SortBy:   &tmpStr,
		Sort:     &tmpStr1,
	}

	data, _, err := s.FormatterDetailRepository.Find(ctx, &criteriaFormatter, &paginationTB)
	if err != nil {
		return nil, err
	}

	f := excelize.NewFile()
	sheet := f.GetSheetName(f.GetActiveSheetIndex())
	arrStyleColWidth := []map[string]interface{}{
		{"COL": "A", "WIDTH": 0.83},
		{"COL": "B", "WIDTH": 15.38},
		{"COL": "C", "WIDTH": 2.14},
		{"COL": "D", "WIDTH": 2.14},
		{"COL": "E", "WIDTH": 57.45},
		{"COL": "F", "WIDTH": 6.43},
		{"COL": "G", "WIDTH": 17.65},
		{"COL": "H", "WIDTH": 10.71},
		{"COL": "I", "WIDTH": 16.83},
		{"COL": "J", "WIDTH": 10.10},
		{"COL": "K", "WIDTH": 17.65},
		{"COL": "L", "WIDTH": 22.14},
	}
	for _, v := range arrStyleColWidth {
		tmpColWidth := fmt.Sprintf("%f", v["WIDTH"])
		colWidth, _ := strconv.ParseFloat(tmpColWidth, 64)
		err = f.SetColWidth(sheet, fmt.Sprintf("%s", v["COL"]), fmt.Sprintf("%s", v["COL"]), colWidth)
		if err != nil {
			return nil, err
		}
	}

	styleDefault, err := f.NewStyle(&excelize.Style{
		Font: &excelize.Font{
			Family: "Arial",
			Size:   10,
		},
	})
	if err != nil {
		return nil, err
	}
	err = f.SetColStyle(sheet, "A:Z", styleDefault)
	if err != nil {
		return nil, err
	}

	stylingBorderRightOnly, _ := f.NewStyle(&excelize.Style{
		Border: []excelize.Border{
			{Type: "right", Color: "000000", Style: 1},
		},
	})

	stylingBorderLeftOnly, _ := f.NewStyle(&excelize.Style{
		Border: []excelize.Border{
			{Type: "left", Color: "000000", Style: 1},
		},
	})

	stylingBorderLROnly, _ := f.NewStyle(&excelize.Style{
		Border: []excelize.Border{
			{Type: "right", Color: "000000", Style: 1},
			{Type: "left", Color: "000000", Style: 1},
		},
	})

	stylingBorderTopOnly, err := f.NewStyle(&excelize.Style{
		Border: []excelize.Border{
			{Type: "top", Color: "000000", Style: 1},
		},
	})
	if err != nil {
		return nil, err
	}

	stylingHeader, err := f.NewStyle(&excelize.Style{
		Font: &excelize.Font{
			Bold: true,
		},
		Border: []excelize.Border{
			{Type: "top", Color: "000000", Style: 1},
			{Type: "bottom", Color: "000000", Style: 1},
			{Type: "left", Color: "000000", Style: 1},
			{Type: "right", Color: "000000", Style: 1},
		},
		Fill: excelize.Fill{
			Type:    "pattern",
			Pattern: 1,
			Color:   []string{"#fac090"},
		},
		Alignment: &excelize.Alignment{
			Horizontal: "center",
			Vertical:   "center",
		},
	})
	if err != nil {
		return nil, err
	}
	stylingHeader2, err := f.NewStyle(&excelize.Style{
		Font: &excelize.Font{
			Color: "#cc00d1",
			Bold:  true,
		},
		Alignment: &excelize.Alignment{
			WrapText:   true,
			Horizontal: "center",
			Vertical:   "center",
		},
		Fill: excelize.Fill{
			Type:    "pattern",
			Pattern: 1,
			Color:   []string{"#dbdbdb"},
		},
	})
	if err != nil {
		return nil, err
	}
	numberFormat := "#,##"
	stylingSubTotal, err := f.NewStyle(&excelize.Style{
		Border: []excelize.Border{
			{Type: "top", Color: "000000", Style: 1},
			{Type: "bottom", Color: "000000", Style: 1},
		},
		Font: &excelize.Font{
			Bold: true,
		},
	})
	if err != nil {
		return nil, err
	}

	stylingSubTotalCurrency, err := f.NewStyle(&excelize.Style{
		Border: []excelize.Border{
			{Type: "top", Color: "000000", Style: 1},
			{Type: "bottom", Color: "000000", Style: 1},
		},
		Font: &excelize.Font{
			Bold: true,
		},
		CustomNumFmt: &numberFormat,
	})
	if err != nil {
		return nil, err
	}

	stylingTotal, err := f.NewStyle(&excelize.Style{
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
			Color:   []string{"#ccff33"},
		},
		CustomNumFmt: &numberFormat,
	})
	if err != nil {
		return nil, err
	}

	stylingTotalControl, err := f.NewStyle(&excelize.Style{
		Border: []excelize.Border{
			{Type: "top", Color: "000000", Style: 1},
			{Type: "left", Color: "000000", Style: 1},
			{Type: "right", Color: "000000", Style: 1},
			{Type: "bottom", Color: "000000", Style: 1},
		},
		Font: &excelize.Font{
			Bold: true,
		},
		Fill: excelize.Fill{
			Type:    "pattern",
			Pattern: 1,
			Color:   []string{"#3ada24"},
		},
		CustomNumFmt: &numberFormat,
	})
	if err != nil {
		return nil, err
	}

	err = f.MergeCell(sheet, "B6", "B8")
	if err != nil {
		return nil, err
	}
	err = f.MergeCell(sheet, "C6", "E8")
	if err != nil {
		return nil, err
	}
	err = f.MergeCell(sheet, "F6", "F8")
	if err != nil {
		return nil, err
	}
	err = f.MergeCell(sheet, "H6", "K7")
	if err != nil {
		return nil, err
	}

	stylingCurrency, err := f.NewStyle(&excelize.Style{
		NumFmt: 7,
		Border: []excelize.Border{
			{Type: "left", Color: "000000", Style: 1},
			{Type: "right", Color: "000000", Style: 1},
		},
		CustomNumFmt: &numberFormat,
	})
	if err != nil {
		return nil, err
	}

	err = f.SetCellStyle(sheet, "B6", "L8", stylingHeader)
	if err != nil {
		return nil, err
	}
	err = f.SetCellStyle(sheet, "F6", "F8", stylingHeader2)
	if err != nil {
		return nil, err
	}

	err = f.SetCellValue(sheet, "B6", "No Akun")
	if err != nil {
		return nil, err
	}
	err = f.SetCellValue(sheet, "C6", "Keterangan")
	if err != nil {
		return nil, err
	}
	err = f.SetCellValue(sheet, "F6", "WP Reff")
	if err != nil {
		return nil, err
	}
	err = f.SetCellValue(sheet, "G6", "PT xxx")
	if err != nil {
		return nil, err
	}
	err = f.SetCellValue(sheet, "G7", "Unaudited")
	if err != nil {
		return nil, err
	}
	err = f.SetCellValue(sheet, "G8", "31-Dec-21")
	if err != nil {
		return nil, err
	}
	err = f.SetCellValue(sheet, "H6", "Adjustment Journal Entry")
	if err != nil {
		return nil, err
	}
	err = f.SetCellValue(sheet, "I8", "Debet")
	if err != nil {
		return nil, err
	}
	err = f.SetCellValue(sheet, "K8", "Kredit")
	if err != nil {
		return nil, err
	}
	err = f.SetCellFormula(sheet, "L6", "=G6")
	if err != nil {
		return nil, err
	}
	err = f.SetCellFormula(sheet, "L7", "=G7")
	if err != nil {
		return nil, err
	}
	err = f.SetCellFormula(sheet, "L8", "=G8")
	if err != nil {
		return nil, err
	}

	criteriaBridge := model.FormatterBridgesFilterModel{}
	tmpStr2 := "TRIAL-BALANCE"
	criteriaBridge.FormatterBridgesFilter.Source = &tmpStr2
	criteriaBridge.FormatterBridgesFilter.FormatterID = &formatterID
	criteriaBridge.FormatterBridgesFilter.TrxRefID = &tb.ID
	bridges, err := s.FormatterBridgesRepository.FindWithCriteria(ctx, &criteriaBridge)
	if err != nil {
		return nil, err
	}

	//find summary aje
	summaryAJE, err := s.AjeRepository.FindSummary(ctx, &tb.ID)
	if err != nil {
		return nil, err
	}

	row := 9
	var summary []map[string]interface{}

	rowCode := make(map[string]int)
	isAutoSum := make(map[string]bool)
	tbRowCode := make(map[string]int)
	customRow := make(map[string]string)
	for _, v := range *data {
		rowCode[v.Code] = row
		if err = f.SetCellStyle(sheet, fmt.Sprintf("G%d", row), fmt.Sprintf("L%d", row), stylingCurrency); err != nil {
			return nil, err
		}
		if err = f.SetCellStyle(sheet, fmt.Sprintf("F%d", row), fmt.Sprintf("F%d", row), stylingBorderLROnly); err != nil {
			return nil, err
		}
		if v.AutoSummary != nil && *v.AutoSummary {
			isAutoSum[v.Code] = true
		}
		if !(v.IsTotal != nil && *v.IsTotal) && v.IsLabel != nil && *v.IsLabel {
			if err = f.SetCellStyle(sheet, fmt.Sprintf("G%d", row), fmt.Sprintf("L%d", row), stylingCurrency); err != nil {
				return nil, err
			}
			if err = f.SetCellStyle(sheet, fmt.Sprintf("F%d", row), fmt.Sprintf("F%d", row), stylingBorderLROnly); err != nil {
				return nil, err
			}
			f.SetCellValue(sheet, fmt.Sprintf("C%d", row), v.Description)
		}

		if v.IsCoa != nil && *v.IsCoa {
			rowBefore := row
			tbdetails, err := s.TBDetailRepository.FindToExport(ctx, &v.Code, &bridges.ID)
			if err != nil {
				return nil, err
			}
			for _, vTbDetail := range *tbdetails {
				if err = f.SetCellStyle(sheet, fmt.Sprintf("G%d", row), fmt.Sprintf("L%d", row), stylingCurrency); err != nil {
					return nil, err
				}
				if err = f.SetCellStyle(sheet, fmt.Sprintf("F%d", row), fmt.Sprintf("F%d", row), stylingBorderLROnly); err != nil {
					return nil, err
				}
				if strings.Contains(strings.ToLower(vTbDetail.Code), "subtotal") {
					continue
				}
				tbRowCode[vTbDetail.Code] = row
				f.SetCellValue(sheet, fmt.Sprintf("B%d", row), vTbDetail.Code)
				f.SetCellValue(sheet, fmt.Sprintf("E%d", row), *vTbDetail.Description)
				amountBeforeAje := 0.0
				if vTbDetail.AmountBeforeAje != nil {
					amountBeforeAje = *vTbDetail.AmountBeforeAje
				}
				amountAjeDr := 0.0
				if vTbDetail.AmountAjeDr != nil {
					amountAjeDr = *vTbDetail.AmountAjeDr
				}
				amountAjeCr := 0.0
				if vTbDetail.AmountAjeCr != nil {
					amountAjeCr = *vTbDetail.AmountAjeCr
				}
				f.SetCellValue(sheet, fmt.Sprintf("G%d", row), amountBeforeAje)
				f.SetCellValue(sheet, fmt.Sprintf("I%d", row), amountAjeDr)
				f.SetCellValue(sheet, fmt.Sprintf("K%d", row), amountAjeCr)
				tmpHeadCoa := fmt.Sprintf("%c", vTbDetail.Code[0])
				if tmpHeadCoa == "9" {
					tmpHeadCoa = vTbDetail.Code[:1]
				}
				switch tmpHeadCoa {
				case "1", "5", "6", "7", "91", "92":
					f.SetCellFormula(sheet, fmt.Sprintf("L%d", row), fmt.Sprintf("=G%d+I%d-K%d", row, row, row))
				default:
					f.SetCellFormula(sheet, fmt.Sprintf("L%d", row), fmt.Sprintf("=G%d-I%d+K%d", row, row, row))
				}
				row++
			}
			rowAfter := row - 1
			rowTB := len(*tbdetails)
			if v.AutoSummary != nil && *v.AutoSummary && rowTB > 1 {
				var tmp = map[string]interface{}{"code": v.Code, "row": row}
				summary = append(summary, tmp)
				f.SetCellValue(sheet, fmt.Sprintf("D%d", row), "Subtotal")
				f.SetCellFormula(sheet, fmt.Sprintf("G%d", row), fmt.Sprintf("=SUM(G%d:G%d)", rowBefore, rowAfter))
				f.SetCellFormula(sheet, fmt.Sprintf("I%d", row), fmt.Sprintf("=SUM(I%d:I%d)", rowBefore, rowAfter))
				f.SetCellFormula(sheet, fmt.Sprintf("K%d", row), fmt.Sprintf("=SUM(K%d:K%d)", rowBefore, rowAfter))
				f.SetCellFormula(sheet, fmt.Sprintf("L%d", row), fmt.Sprintf("=SUM(L%d:L%d)", rowBefore, rowAfter))
				f.SetCellStyle(sheet, fmt.Sprintf("D%d", row), fmt.Sprintf("F%d", row), stylingSubTotal)
				if err = f.SetCellStyle(sheet, fmt.Sprintf("G%d", row), fmt.Sprintf("L%d", row), stylingSubTotalCurrency); err != nil {
					return nil, err
				}
				rowCode[fmt.Sprintf("%s_SUBTOTAL", v.Code)] = row
				row++
				// if err = f.SetCellStyle(sheet, fmt.Sprintf("G%d", row), fmt.Sprintf("L%d", row), stylingSubTotalCurrency); err != nil {
				// 	return nil, err
				// }
				if err = f.SetCellStyle(sheet, fmt.Sprintf("F%d", row), fmt.Sprintf("L%d", row), stylingBorderLROnly); err != nil {
					return nil, err
				}
			}
		}

		if v.IsTotal != nil && *v.IsTotal {
			tbRowCode[v.Code] = row
			f.SetCellValue(sheet, fmt.Sprintf("C%d", row), v.Description)
			if v.Code == "TOTAL_LIABILITAS_DAN_EKUITAS" {
				rowAset := row
				if _, ok := rowCode["TOTAL_ASET"]; ok {
					rowAset = rowCode["TOTAL_ASET"]
				}
				f.SetCellFormula(sheet, "G5", fmt.Sprintf("=G%d-G%d", rowAset, row))
				f.SetCellFormula(sheet, "L5", fmt.Sprintf("=L%d-L%d", rowAset, row))
			}

			//show control aje
			if v.Code == "CONTROL" {
				f.SetCellFormula(sheet, "I5", fmt.Sprintf("=I%d", row))
				f.SetCellFormula(sheet, "K5", fmt.Sprintf("=K%d", row))
			}
			if v.Code == "CONTROL_TO_ADJUSTMENT_SHEET" {
				f.SetCellValue(sheet, fmt.Sprintf("C%d", row), v.Description)
				dbt := 0.0
				if summaryAJE.IncomeStatementDr != nil && *summaryAJE.IncomeStatementDr != 0 {
					dbt = *summaryAJE.IncomeStatementDr
				}
				if summaryAJE.BalanceSheetDr != nil && *summaryAJE.BalanceSheetDr != 0 {
					dbt += *summaryAJE.BalanceSheetDr
				}
				cdt := 0.0
				if summaryAJE.IncomeStatementCr != nil && *summaryAJE.IncomeStatementCr != 0 {
					cdt = *summaryAJE.IncomeStatementCr
				}
				if summaryAJE.BalanceSheetCr != nil && *summaryAJE.BalanceSheetCr != 0 {
					cdt += *summaryAJE.BalanceSheetCr
				}
				f.SetCellValue(sheet, fmt.Sprintf("I%d", row), dbt)
				f.SetCellValue(sheet, fmt.Sprintf("K%d", row), cdt)
				f.SetCellFormula(sheet, fmt.Sprintf("J%d", row), fmt.Sprintf("=I%d-K%d", row, row))
			}
			if v.Code == "TOTAL_JOURNAL_IN_WP" {
				f.SetCellFormula(sheet, fmt.Sprintf("J%d", row), fmt.Sprintf("=I%d-K%d", row, row))
			}

			if v.Code != "TOTAL_JOURNAL_IN_WP" && v.Code != "CONTROL_TO_ADJUSTMENT_SHEET" && v.Code != "CONTROL" && v.Code != "CONTROL_TO_WBS_1" {
				if v.IsControl != nil && *v.IsControl {
					f.SetCellStyle(sheet, fmt.Sprintf("B%d", row), fmt.Sprintf("L%d", row), stylingTotalControl)
				} else {
					f.SetCellStyle(sheet, fmt.Sprintf("B%d", row), fmt.Sprintf("L%d", row), stylingTotal)
				}
			}
			if v.FxSummary == "" {
				row++
				continue
			}
			for chr := 'G'; chr <= 'L'; chr++ {
				if fmt.Sprintf("%c", chr) == "H" || fmt.Sprintf("%c", chr) == "J" || ((fmt.Sprintf("%c", chr) == "G" || fmt.Sprintf("%c", chr) == "L") && (v.Code == "TOTAL_JOURNAL_IN_WP" || v.Code == "CONTROL_TO_ADJUSTMENT_SHEET" || v.Code == "CONTROL")) {
					continue
				}
				formula := v.FxSummary
				// reg := regexp.MustCompile(`[A-Za-z_]+|[0-9]+\d{3,}`)
				reg := regexp.MustCompile(`[A-Za-z0-9_~:()]+|[0-9]+\d{2,}`)
				match := reg.FindAllString(formula, -1)
				for _, vMatch := range match {
					if len(vMatch) < 3 {
						continue
					}
					if isAutoSum[vMatch] {
						if rowCode[fmt.Sprintf("%s_SUBTOTAL", vMatch)] != 0 {
							formula = strings.ReplaceAll(formula, vMatch, fmt.Sprintf("%c%d", chr, rowCode[fmt.Sprintf("%s_SUBTOTAL", vMatch)]))
						} else {
							formula = strings.ReplaceAll(formula, vMatch, fmt.Sprintf("%c%d", chr, rowCode[vMatch]))
						}
					} else {
						if rowCode[vMatch] != 0 {
							formula = strings.ReplaceAll(formula, vMatch, fmt.Sprintf("%c%d", chr, rowCode[vMatch]))
						}
					}
				}
				f.SetCellFormula(sheet, fmt.Sprintf("%c%d", chr, row), fmt.Sprintf("=%s", formula))
			}
			row++
			continue
		}
		row++
	}

	customRow["310401004"] = "=LABA_BERSIH"
	customRow["310402002"] = "=TOTAL_PENGHASILAN_KOMPREHENSIF_LAIN~BS-SUM(310501002,310502002,310503002)"
	customRow["310501002"] = "=950101001"
	customRow["310502002"] = "=950301001+950301002"
	customRow["310503002"] = "=950401001+950401002"
	for key, nRow := range tbRowCode {
		if strings.Contains(customRow["310401004"], key) {
			customRow["310401004"] = strings.ReplaceAll(customRow["310401004"], key, fmt.Sprintf("@%d", nRow))
		}
		if strings.Contains(customRow["310402002"], key) && key != "RE" {
			customRow["310402002"] = strings.ReplaceAll(customRow["310402002"], key, fmt.Sprintf("@%d", nRow))
		}
		if strings.Contains(customRow["310501002"], key) {
			customRow["310501002"] = strings.ReplaceAll(customRow["310501002"], key, fmt.Sprintf("@%d", nRow))
		}
		if strings.Contains(customRow["310502002"], key) {
			customRow["310502002"] = strings.ReplaceAll(customRow["310502002"], key, fmt.Sprintf("@%d", nRow))
		}
		if strings.Contains(customRow["310503002"], key) {
			customRow["310503002"] = strings.ReplaceAll(customRow["310503002"], key, fmt.Sprintf("@%d", nRow))
		}
	}

	for key, vCustomRow := range customRow {
		if val, ok := tbRowCode[key]; ok {
			f.SetCellFormula(sheet, fmt.Sprintf("G%d", val), strings.ReplaceAll(vCustomRow, "@", "G"))
			f.SetCellFormula(sheet, fmt.Sprintf("I%d", val), strings.ReplaceAll(vCustomRow, "@", "I"))
			f.SetCellFormula(sheet, fmt.Sprintf("K%d", val), strings.ReplaceAll(vCustomRow, "@", "K"))
			// f.SetCellFormula(sheet, fmt.Sprintf("L%d", val), strings.ReplaceAll(vCustomRow, "@", "L"))
		}
	}

	if err = f.SetCellStyle(sheet, fmt.Sprintf("B%d", row), fmt.Sprintf("L%d", row), stylingBorderTopOnly); err != nil {
		return nil, err
	}

	/* if err = f.SetCellStyle(sheet, "G9", fmt.Sprintf("L%d", row-1), stylingCurrency); err != nil {
		return nil, err
	} */

	if err = f.SetCellStyle(sheet, "A9", fmt.Sprintf("A%d", row-1), stylingBorderRightOnly); err != nil {
		return nil, err
	}
	if err = f.SetCellStyle(sheet, "M9", fmt.Sprintf("M%d", row-1), stylingBorderLeftOnly); err != nil {
		return nil, err
	}

	if err = f.SetSheetFormatPr(sheet, excelize.DefaultRowHeight(12.85)); err != nil {
		return nil, err
	}
	f.SetDefaultFont("Arial")
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
	fileName := fmt.Sprintf("TrialBalance_%s.xlsx", period)
	fileLoc := fmt.Sprintf("assets/%d/%s", ctx.Auth.ID, fileName)
	err = f.SaveAs(fileLoc)
	if err != nil {
		return nil, err
	}

	result := &dto.TrialBalanceExportResponse{
		FileName: fileName,
		Path:     fileLoc,
	}
	return result, nil
}

func (s *service) GetVersion(ctx *abstraction.Context, payload *dto.GetVersionRequest) (*dto.GetVersionResponse, error) {
	filter := model.TrialBalanceFilterModel{
		CompanyCustomFilter: model.CompanyCustomFilter{
			CompanyID:          payload.CompanyID,
			ArrCompanyID:       payload.ArrCompanyID,
			ArrCompanyString:   payload.ArrCompanyString,
			ArrCompanyOperator: payload.ArrCompanyOperator,
		},
	}
	filter.ArrStatus = payload.ArrStatus
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
