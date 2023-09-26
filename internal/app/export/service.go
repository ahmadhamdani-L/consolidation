package export

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"time"

	"mcash-finance-console-core/internal/abstraction"
	"mcash-finance-console-core/internal/dto"
	"mcash-finance-console-core/internal/factory"
	"mcash-finance-console-core/internal/model"
	"mcash-finance-console-core/internal/repository"
	"mcash-finance-console-core/pkg/kafka"
	"mcash-finance-console-core/pkg/util/helper"
	"mcash-finance-console-core/pkg/util/response"

	"github.com/xuri/excelize/v2"
	"gorm.io/gorm"
)

type service struct {
	NotificationRepository                     repository.Notification
	ImportedWorksheetRepository                repository.ImportedWorksheet
	ConsolidationRepository                    repository.Consolidation
	AgingUtangPiutangRepository                repository.AgingUtangPiutang
	AgingUtangPiutangDetailRepository          repository.AgingUtangPiutangDetail
	FormatterRepository                        repository.Formatter
	FormatterBridgesRepository                 repository.FormatterBridges
	EmployeeBenefitRepository                  repository.EmployeeBenefit
	EmployeeBenefitDetailRepository            repository.EmployeeBenefitDetail
	AdjustmentRepository                       repository.Adjustment
	AdjustmentDetailRepository                 repository.AdjustmentDetail
	InvestasiNonTbkRepository                  repository.InvestasiNonTbk
	InvestasiNonTbkDetailRepository            repository.InvestasiNonTbkDetail
	InvestasiTbkRepository                     repository.InvestasiTbk
	InvestasiTbkDetailRepository               repository.InvestasiTbkDetail
	MutasiDtaRepository                        repository.MutasiDta
	MutasiDtaDetailRepository                  repository.MutasiDtaDetail
	MutasiFaRepository                         repository.MutasiFa
	MutasiFaDetailRepository                   repository.MutasiFaDetail
	MutasiIaRepository                         repository.MutasiIa
	MutasiIaDetailRepository                   repository.MutasiIaDetail
	MutasiPersediaanRepository                 repository.MutasiPersediaan
	MutasiPersediaanDetailRepository           repository.MutasiPersediaanDetail
	MutasiRuaRepository                        repository.MutasiRua
	MutasiRuaDetailRepository                  repository.MutasiRuaDetail
	PembelianPenjualanBerelasiRepository       repository.PembelianPenjualanBerelasi
	PembelianPenjualanBerelasiDetailRepository repository.PembelianPenjualanBerelasiDetail
	TrialBalanceRepository                     repository.TrialBalance
	TrialBalanceDetailRepository               repository.TrialBalanceDetail
	FormatterDetailRepository                  repository.FormatterDetail

	Db *gorm.DB
}

type Service interface {
	RequestExport(ctx *abstraction.Context, payload dto.ExportRequest) (dto.ExportResponse, error)
	GetExport(ctx *abstraction.Context, payload *dto.GetExportRequest) (*string, error)
	ExportConsol(ctx *abstraction.Context, payload *dto.ExportConsolRequest) (*string, error)
}

func NewService(f *factory.Factory) *service {
	return &service{
		NotificationRepository:                     f.NotificationRepository,
		ImportedWorksheetRepository:                f.ImportedWorksheetRepository,
		ConsolidationRepository:                    f.ConsolidationRepository,
		AgingUtangPiutangRepository:                f.AgingUtangPiutangRepository,
		AgingUtangPiutangDetailRepository:          f.AgingUtangPiutangDetailRepository,
		FormatterRepository:                        f.FormatterRepository,
		FormatterBridgesRepository:                 f.FormatterBridgesRepository,
		EmployeeBenefitRepository:                  f.EmployeeBenefitRepository,
		EmployeeBenefitDetailRepository:            f.EmployeeBenefitDetailRepository,
		AdjustmentRepository:                       f.AdjustmentRepository,
		AdjustmentDetailRepository:                 f.AdjustmentDetailRepository,
		InvestasiNonTbkRepository:                  f.InvestasiNonTbkRepository,
		InvestasiNonTbkDetailRepository:            f.InvestasiNonTbkDetailRepository,
		InvestasiTbkRepository:                     f.InvestasiTbkRepository,
		InvestasiTbkDetailRepository:               f.InvestasiTbkDetailRepository,
		MutasiDtaRepository:                        f.MutasiDtaRepository,
		MutasiDtaDetailRepository:                  f.MutasiDtaDetailRepository,
		MutasiFaRepository:                         f.MutasiFaRepository,
		MutasiFaDetailRepository:                   f.MutasiFaDetailRepository,
		MutasiIaRepository:                         f.MutasiIaRepository,
		MutasiIaDetailRepository:                   f.MutasiIaDetailRepository,
		MutasiPersediaanRepository:                 f.MutasiPersediaanRepository,
		MutasiPersediaanDetailRepository:           f.MutasiPersediaanDetailRepository,
		MutasiRuaRepository:                        f.MutasiRuaRepository,
		MutasiRuaDetailRepository:                  f.MutasiRuaDetailRepository,
		PembelianPenjualanBerelasiRepository:       f.PembelianPenjualanBerelasiRepository,
		PembelianPenjualanBerelasiDetailRepository: f.PembelianPenjualanBerelasiDetailRepository,
		TrialBalanceRepository:                     f.TrialBalanceRepository,
		TrialBalanceDetailRepository:               f.TrialBalanceDetailRepository,
		FormatterDetailRepository:                  f.FormatterDetailRepository,

		Db: f.Db,
	}
}

func (s *service) RequestExport(ctx *abstraction.Context, payload *dto.ExportRequest) (dto.ExportResponse, error) {
	importedWorksheet, err := s.ImportedWorksheetRepository.FindByID(ctx, &payload.ImportID)
	if err != nil {
		return dto.ExportResponse{}, helper.ErrorHandler(err)
	}
	allowed := helper.CompanyValidation(ctx.Auth.ID, importedWorksheet.CompanyID)
	if !allowed {
		return dto.ExportResponse{}, response.ErrorBuilder(&response.ErrorConstant.UnauthorizedAccess, errors.New("not allowed"))
	}

	datePeriod, err := time.Parse(time.RFC3339, importedWorksheet.Period)
	if err != nil {
		return dto.ExportResponse{}, err
	}
	period := datePeriod.Format("2006-01-02")

	waktu := time.Now()
	msg := kafka.JsonData{
		UserID:    ctx.Auth.ID,
		CompanyID: importedWorksheet.CompanyID,
		Timestamp: &waktu,
		// Name:      ctx.Auth.Name,
		Filter: struct {
			Period   string
			Versions int
			Request  string
		}{period, importedWorksheet.Versions, payload.Request},
	}

	jsonStr, err := json.Marshal(msg)
	if err != nil {
		fmt.Printf("Error: %s", err.Error())
		return dto.ExportResponse{}, err
	}

	go kafka.NewService("EXPORT").SendMessage("EXPORT", string(jsonStr))
	return dto.ExportResponse{
		Message: "Sukses Request Data",
	}, nil
}

func (s *service) GetExport(ctx *abstraction.Context, payload *dto.GetExportRequest) (*string, error) {
	type NotifData struct {
		File string
	}
	var jsonData NotifData
	notifData, err := s.NotificationRepository.FindByID(ctx, &payload.NotificationID)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal([]byte(notifData.Data), &jsonData)
	if err != nil {
		return nil, err
	}
	return &jsonData.File, nil
}

func (s *service) ExportConsol(ctx *abstraction.Context, payload *dto.ExportConsolRequest) (*dto.ExportResponse, error) {
	consolidationData, err := s.ConsolidationRepository.FindByID(ctx, &payload.ConsolidationID)
	if err != nil {
		return nil, helper.ErrorHandler(err)
	}
	allowed := helper.CompanyValidation(ctx.Auth.ID, consolidationData.CompanyID)
	if !allowed {
		return nil, response.ErrorBuilder(&response.ErrorConstant.UnauthorizedAccess, errors.New("not allowed"))
	}

	datePeriod, err := time.Parse(time.RFC3339, consolidationData.Period)
	if err != nil {
		return nil, err
	}
	period := datePeriod.Format("2006-01-02")

	waktu := time.Now()
	msg := kafka.JsonData{
		UserID:    ctx.Auth.ID,
		CompanyID: consolidationData.CompanyID,
		Timestamp: &waktu,
		Data:      fmt.Sprintf("%d", consolidationData.ID),
		Filter: struct {
			Period   string
			Versions int
			Request  string
		}{
			Period:   period,
			Versions: consolidationData.ConsolidationVersions,
			Request:  "",
		},
	}
	jsonStr, err := json.Marshal(&msg)
	if err != nil {
		return nil, err
	}
	go kafka.NewService("EXPORT").SendMessage("EXPORT_CONSOLIDATION", string(jsonStr))
	return &dto.ExportResponse{
		Message: "Sukses Request export consolidation",
	}, nil
}

func (s *service) ExportModul(ctx *abstraction.Context, payload *dto.ExportRequest) (*string, error) {
	importData, err := s.ImportedWorksheetRepository.FindByID(ctx, &payload.ImportID)
	if err != nil {
		return nil, helper.ErrorHandler(err)
	}
	allowed := helper.CompanyValidation(ctx.Auth.ID, importData.CompanyID)
	if !allowed {
		return nil, response.ErrorBuilder(&response.ErrorConstant.UnauthorizedAccess, errors.New("not allowed"))
	}

	if importData.Status == 0 {
		return nil, response.ErrorBuilder(&response.ErrorConstant.NotFound, errors.New("data not found"))
	}

	criteriaTB := model.TrialBalanceFilterModel{}
	criteriaTB.CompanyID = &importData.CompanyID
	criteriaTB.Period = &importData.Period
	criteriaTB.Versions = &importData.Versions
	trialBalance, err := s.TrialBalanceRepository.FindByCriteria(ctx, &criteriaTB)
	if err != nil {
		return nil, helper.ErrorHandler(err)
	}

	f := excelize.NewFile()
	currentSheet := f.GetSheetName(f.GetActiveSheetIndex())
	f.DeleteSheet(currentSheet)
	wg := sync.WaitGroup{}
	errorMessages := []string{}
	wg.Add(10)

	resultF, err := s.ExportAdjustments(ctx, f, trialBalance)
	if err != nil {
		errorMessages = append(errorMessages, fmt.Sprintf("Error Export Adjustment. Error: %s", err.Error()))
		return nil, err
	}
	f = resultF

	resultF, err = s.ExportTrialBalances(ctx, f, trialBalance)
	if err != nil {
		errorMessages = append(errorMessages, fmt.Sprintf("Error Export Trial Balance. Error: %s", err.Error()))
		return nil, err
	}
	f = resultF

	go func() {
		defer wg.Done()
		filterAgingUtangPiutang := model.AgingUtangPiutangFilterModel{}
		filterAgingUtangPiutang.CompanyID = &importData.CompanyID
		filterAgingUtangPiutang.Period = &importData.Period
		filterAgingUtangPiutang.Versions = &importData.Versions
		agingUtangPiutang, err := s.AgingUtangPiutangRepository.FindByCriteria(ctx, &filterAgingUtangPiutang)
		if err != nil {
			errorMessages = append(errorMessages, fmt.Sprintf("Error Export Aging Utang Piutang. Error: %s", err.Error()))
			return
		}

		resultF, err := s.ExportAgingUtangPiutangs(ctx, f, agingUtangPiutang)
		if err != nil {
			errorMessages = append(errorMessages, fmt.Sprintf("Error Export Aging Utang Piutang. Error: %s", err.Error()))
			return
		}
		f = resultF
	}()

	go func() {
		defer wg.Done()
		filterEmployeeBenefit := model.EmployeeBenefitFilterModel{}
		filterEmployeeBenefit.CompanyID = &importData.CompanyID
		filterEmployeeBenefit.Period = &importData.Period
		filterEmployeeBenefit.Versions = &importData.Versions
		employeeBenefit, err := s.EmployeeBenefitRepository.FindByCriteria(ctx, &filterEmployeeBenefit)
		if err != nil {
			errorMessages = append(errorMessages, fmt.Sprintf("Error Export Employee Benefit. Error: %s", err.Error()))
			return
		}

		resultF, err := s.ExportEmployeeBenefits(ctx, f, employeeBenefit)
		if err != nil {
			errorMessages = append(errorMessages, fmt.Sprintf("Error Export Employee Benefit. Error: %s", err.Error()))
			return
		}
		f = resultF
	}()

	go func() {
		defer wg.Done()
		filterInvestasiNonTbk := model.InvestasiNonTbkFilterModel{}
		filterInvestasiNonTbk.CompanyID = &importData.CompanyID
		filterInvestasiNonTbk.Period = &importData.Period
		filterInvestasiNonTbk.Versions = &importData.Versions
		investasiNonTbk, err := s.InvestasiNonTbkRepository.FindByCriteria(ctx, &filterInvestasiNonTbk)
		if err != nil {
			errorMessages = append(errorMessages, fmt.Sprintf("Error Export Investasi Non Tbk. Error: %s", err.Error()))
			return
		}

		resultF, err := s.ExportInvestasiNonTbks(ctx, f, investasiNonTbk)
		if err != nil {
			errorMessages = append(errorMessages, fmt.Sprintf("Error Export Investasi Non Tbk. Error: %s", err.Error()))
			return
		}
		f = resultF
	}()

	go func() {
		defer wg.Done()
		filterInvestasiTbk := model.InvestasiTbkFilterModel{}
		filterInvestasiTbk.CompanyID = &importData.CompanyID
		filterInvestasiTbk.Period = &importData.Period
		filterInvestasiTbk.Versions = &importData.Versions
		investasiTbk, err := s.InvestasiTbkRepository.FindByCriteria(ctx, &filterInvestasiTbk)
		if err != nil {
			errorMessages = append(errorMessages, fmt.Sprintf("Error Export Investasi Tbk. Error: %s", err.Error()))
			return
		}

		resultF, err := s.ExportInvestasiTbks(ctx, f, investasiTbk)
		if err != nil {
			errorMessages = append(errorMessages, fmt.Sprintf("Error Export Investasi Tbk. Error: %s", err.Error()))
			return
		}
		f = resultF
	}()

	go func() {
		defer wg.Done()
		filterMutasiDta := model.MutasiDtaFilterModel{}
		filterMutasiDta.CompanyID = &importData.CompanyID
		filterMutasiDta.Period = &importData.Period
		filterMutasiDta.Versions = &importData.Versions
		mutasiDta, err := s.MutasiDtaRepository.FindByCriteria(ctx, &filterMutasiDta)
		if err != nil {
			errorMessages = append(errorMessages, fmt.Sprintf("Error Export Mutasi Dta. Error: %s", err.Error()))
			return
		}

		resultF, err := s.ExportMutasiDtas(ctx, f, mutasiDta)
		if err != nil {
			errorMessages = append(errorMessages, fmt.Sprintf("Error Export Mutasi Dta. Error: %s", err.Error()))
			return
		}
		f = resultF
	}()

	go func() {
		defer wg.Done()
		filterMutasiFa := model.MutasiFaFilterModel{}
		filterMutasiFa.CompanyID = &importData.CompanyID
		filterMutasiFa.Period = &importData.Period
		filterMutasiFa.Versions = &importData.Versions
		mutasiFa, err := s.MutasiFaRepository.FindByCriteria(ctx, &filterMutasiFa)
		if err != nil {
			errorMessages = append(errorMessages, fmt.Sprintf("Error Export Mutasi Fa. Error: %s", err.Error()))
			return
		}

		resultF, err := s.ExportMutasiFas(ctx, f, mutasiFa)
		if err != nil {
			errorMessages = append(errorMessages, fmt.Sprintf("Error Export Mutasi Fa. Error: %s", err.Error()))
			return
		}
		f = resultF
	}()

	go func() {
		defer wg.Done()
		filterMutasiIa := model.MutasiIaFilterModel{}
		filterMutasiIa.CompanyID = &importData.CompanyID
		filterMutasiIa.Period = &importData.Period
		filterMutasiIa.Versions = &importData.Versions
		mutasiIa, err := s.MutasiIaRepository.FindByCriteria(ctx, &filterMutasiIa)
		if err != nil {
			errorMessages = append(errorMessages, fmt.Sprintf("Error Export Mutasi Ia. Error: %s", err.Error()))
			return
		}

		resultF, err := s.ExportMutasiIas(ctx, f, mutasiIa)
		if err != nil {
			errorMessages = append(errorMessages, fmt.Sprintf("Error Export Mutasi Ia. Error: %s", err.Error()))
			return
		}
		f = resultF
	}()

	go func() {
		defer wg.Done()
		filterMutasiPersediaan := model.MutasiPersediaanFilterModel{}
		filterMutasiPersediaan.CompanyID = &importData.CompanyID
		filterMutasiPersediaan.Period = &importData.Period
		filterMutasiPersediaan.Versions = &importData.Versions
		mutasiIa, err := s.MutasiPersediaanRepository.FindByCriteria(ctx, &filterMutasiPersediaan)
		if err != nil {
			errorMessages = append(errorMessages, fmt.Sprintf("Error Export Mutasi Persediaan. Error: %s", err.Error()))
			return
		}

		resultF, err := s.ExportMutasiPersediaans(ctx, f, mutasiIa)
		if err != nil {
			errorMessages = append(errorMessages, fmt.Sprintf("Error Export Mutasi Persediaan. Error: %s", err.Error()))
			return
		}
		f = resultF
	}()

	go func() {
		defer wg.Done()
		filterMutasiRua := model.MutasiRuaFilterModel{}
		filterMutasiRua.CompanyID = &importData.CompanyID
		filterMutasiRua.Period = &importData.Period
		filterMutasiRua.Versions = &importData.Versions
		mutasiRua, err := s.MutasiRuaRepository.FindByCriteria(ctx, &filterMutasiRua)
		if err != nil {
			errorMessages = append(errorMessages, fmt.Sprintf("Error Export Mutasi Rua. Error: %s", err.Error()))
			return
		}

		resultF, err := s.ExportMutasiRuas(ctx, f, mutasiRua)
		if err != nil {
			errorMessages = append(errorMessages, fmt.Sprintf("Error Export Mutasi Rua. Error: %s", err.Error()))
			return
		}
		f = resultF
	}()

	go func() {
		defer wg.Done()
		filterPembelianPenjualanBerelasi := model.PembelianPenjualanBerelasiFilterModel{}
		filterPembelianPenjualanBerelasi.CompanyID = &importData.CompanyID
		filterPembelianPenjualanBerelasi.Period = &importData.Period
		filterPembelianPenjualanBerelasi.Versions = &importData.Versions
		pembelianPenjualanBerelasi, err := s.PembelianPenjualanBerelasiRepository.FindByCriteria(ctx, &filterPembelianPenjualanBerelasi)
		if err != nil {
			errorMessages = append(errorMessages, fmt.Sprintf("Error Export Pembelian Penjualan Berelasi. Error: %s", err.Error()))
			return
		}

		resultF, err := s.ExportPembelianPenjualanBerelasis(ctx, f, pembelianPenjualanBerelasi)
		if err != nil {
			errorMessages = append(errorMessages, fmt.Sprintf("Error Export Pembelian Penjualan Berelasi. Error: %s", err.Error()))
			return
		}
		f = resultF
	}()

	wg.Wait()

	if len(errorMessages) > 0 {
		return nil, response.ErrorBuilder(&response.ErrorConstant.InternalServerError, errors.New(strings.Join(errorMessages, ", ")))
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
	datePeriod, err := time.Parse(time.RFC3339, importData.Period)
	if err != nil {
		return nil, err
	}
	period := datePeriod.Format("2006-01-02")
	fileName := fmt.Sprintf("%s_%s__Ver-%d.xlsx", trialBalance.Company.Name, period, importData.Versions)
	fileLoc := fmt.Sprintf("assets/%s", fileName)
	err = f.SaveAs(fileLoc)
	if err != nil {
		return nil, err
	}

	return &fileLoc, nil
}

func (s *service) ExportAdjustments(ctx *abstraction.Context, f *excelize.File, trialBalance *model.TrialBalanceEntityModel) (*excelize.File, error) {
	datas, err := s.AdjustmentRepository.ExportAll(ctx, &trialBalance.ID)
	if err != nil {
		return nil, err
	}

	datePeriod, err := time.Parse(time.RFC3339, trialBalance.Period)
	if err != nil {
		return nil, err
	}

	tb, err := s.TrialBalanceRepository.Get(ctx, trialBalance.ID)
	if err != nil {
		return nil, err
	}

	sheet := "ADJUSTMENT"
	f.NewSheet(sheet)
	arrStyleColWidth := []map[string]interface{}{
		{"COL": "A", "WIDTH": 6.50},
		{"COL": "B", "WIDTH": 15.38},
		{"COL": "C", "WIDTH": 7.25},
		{"COL": "D", "WIDTH": 2.14},
		{"COL": "E", "WIDTH": 57.45},
		{"COL": "F", "WIDTH": 7.25},
		{"COL": "G", "WIDTH": 15.38},
		{"COL": "H", "WIDTH": 15.38},
		{"COL": "I", "WIDTH": 15.38},
		{"COL": "J", "WIDTH": 15.38},
	}

	for _, v := range arrStyleColWidth {
		tmpColWidth := fmt.Sprintf("%f", v["WIDTH"])
		colWidth, _ := strconv.ParseFloat(tmpColWidth, 64)
		f.SetColWidth(sheet, fmt.Sprintf("%s", v["COL"]), fmt.Sprintf("%s", v["COL"]), colWidth)
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

	f.SetColStyle(sheet, "A:Z", styleDefault)
	stylingBorderLROnly, err := f.NewStyle(&excelize.Style{
		Border: []excelize.Border{
			{Type: "right", Color: "000000", Style: 2},
			{Type: "left", Color: "000000", Style: 2},
		},
	})
	if err != nil {
		return nil, err
	}

	stylingBorderAll, err := f.NewStyle(&excelize.Style{
		Border: []excelize.Border{
			{Type: "top", Color: "000000", Style: 2},
			{Type: "right", Color: "000000", Style: 2},
			{Type: "left", Color: "000000", Style: 2},
			{Type: "bottom", Color: "000000", Style: 2},
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
			{Type: "top", Color: "000000", Style: 2},
			{Type: "bottom", Color: "000000", Style: 2},
			{Type: "left", Color: "000000", Style: 2},
			{Type: "right", Color: "000000", Style: 2},
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

	f.MergeCell(sheet, "A6", "A7")
	f.MergeCell(sheet, "B6", "B7")
	f.MergeCell(sheet, "C6", "C7")
	f.MergeCell(sheet, "D6", "E7")
	f.MergeCell(sheet, "F6", "F7")
	f.MergeCell(sheet, "G6", "H6")
	f.MergeCell(sheet, "I6", "J6")

	stylingCurrency, err := f.NewStyle(&excelize.Style{
		NumFmt: 7,
		Border: []excelize.Border{
			{Type: "left", Color: "000000", Style: 1},
			{Type: "right", Color: "000000", Style: 1},
		},
	})
	if err != nil {
		return nil, err
	}

	stylingSubTotal, err := f.NewStyle(&excelize.Style{
		Border: []excelize.Border{
			{Type: "top", Color: "000000", Style: 2},
			{Type: "bottom", Color: "000000", Style: 2},
			{Type: "left", Color: "000000", Style: 2},
			{Type: "right", Color: "000000", Style: 2},
		},
		Fill: excelize.Fill{
			Type:    "pattern",
			Pattern: 1,
			Color:   []string{"#00ff00"},
		},
		Font: &excelize.Font{
			Bold: true,
		},
	})
	if err != nil {
		return nil, err
	}

	stylingSubTotal2, err := f.NewStyle(&excelize.Style{
		Border: []excelize.Border{
			{Type: "top", Color: "000000", Style: 2},
			{Type: "bottom", Color: "000000", Style: 2},
			{Type: "left", Color: "000000", Style: 2},
			{Type: "right", Color: "000000", Style: 2},
		},
		Fill: excelize.Fill{
			Type:    "pattern",
			Pattern: 1,
			Color:   []string{"#99ff33"},
		},
		Font: &excelize.Font{
			Bold: true,
		},
	})
	if err != nil {
		return nil, err
	}

	f.SetCellValue(sheet, "A2", "Company")
	f.SetCellValue(sheet, "B2", ": "+tb.Company.CompanyEntity.Name)
	f.SetCellValue(sheet, "A3", "Date")
	f.SetCellValue(sheet, "B3", ": "+datePeriod.Format("02-Jan-06"))
	f.SetCellValue(sheet, "A4", "Subject")
	f.SetCellValue(sheet, "B4", ": AJE")

	f.SetCellStyle(sheet, "A6", "J7", stylingHeader)
	f.SetCellValue(sheet, "A6", "NO")
	f.SetCellValue(sheet, "B6", "COA")
	f.SetCellValue(sheet, "C6", "AJE")
	f.SetCellValue(sheet, "D6", "DESCRIPTION")
	f.SetCellValue(sheet, "F6", "WP")
	f.SetCellValue(sheet, "F7", "REFF")
	f.SetCellValue(sheet, "G6", "Balance Sheet")
	f.SetCellValue(sheet, "G7", "DR")
	f.SetCellValue(sheet, "H7", "CR")
	f.SetCellValue(sheet, "I6", "Income Stat")
	f.SetCellValue(sheet, "I7", "DR")
	f.SetCellValue(sheet, "J7", "CR")

	row := 8
	counterRef := 1
	rowBefore := row

	reffN := make(map[string]int)
	for i, v := range *datas {
		f.SetCellValue(sheet, fmt.Sprintf("A%d", row), i+1)
		f.SetCellValue(sheet, fmt.Sprintf("B%d", row), v.CoaCode)
		f.SetCellValue(sheet, fmt.Sprintf("C%d", row), *v.ReffNumber)
		f.SetCellValue(sheet, fmt.Sprintf("E%d", row), *v.Description)
		f.SetCellValue(sheet, fmt.Sprintf("E%d", row+1), *v.Note)
		if v.BalanceSheetDr != nil && *v.BalanceSheetDr != 0 {
			f.SetCellValue(sheet, fmt.Sprintf("G%d", row), *v.BalanceSheetDr)
		}
		if v.BalanceSheetCr != nil && *v.BalanceSheetCr != 0 {
			f.SetCellValue(sheet, fmt.Sprintf("H%d", row), *v.BalanceSheetCr)
		}
		if v.IncomeStatementDr != nil && *v.IncomeStatementDr != 0 {
			f.SetCellValue(sheet, fmt.Sprintf("I%d", row), *v.IncomeStatementDr)
		}
		if v.IncomeStatementCr != nil && *v.IncomeStatementCr != 0 {
			f.SetCellValue(sheet, fmt.Sprintf("J%d", row), *v.IncomeStatementCr)
		}
		if reffN[*v.ReffNumber] == 0 {
			row += 1
			if *v.Note != "" {
				row += 1
			}
		}
		if reffN[*v.ReffNumber] != 0 {
			row += 2
			if *v.Note != "" {
				row += 1
			}
		}
		reffN[*v.ReffNumber] = counterRef
		counterRef++
	}

	if len(*datas) == 0 {
		row += 1
	}

	f.SetCellStyle(sheet, "A8", fmt.Sprintf("B%d", row), stylingBorderLROnly)
	f.SetCellStyle(sheet, "C8", fmt.Sprintf("C%d", row), stylingBorderLROnly)
	f.SetCellStyle(sheet, "F8", fmt.Sprintf("F%d", row), stylingBorderLROnly)

	f.SetCellStyle(sheet, "G8", fmt.Sprintf("J%d", row), stylingCurrency)

	f.SetCellStyle(sheet, fmt.Sprintf("A%d", row), fmt.Sprintf("F%d", row+1), stylingBorderAll)

	f.MergeCell(sheet, fmt.Sprintf("D%d", row), fmt.Sprintf("E%d", row))
	f.MergeCell(sheet, fmt.Sprintf("D%d", row+1), fmt.Sprintf("E%d", row+1))

	f.SetSheetFormatPr(sheet, excelize.DefaultRowHeight(12.85))

	f.SetCellValue(sheet, fmt.Sprintf("D%d", row), "Total PM")

	f.SetCellFormula(sheet, fmt.Sprintf("G%d", row), fmt.Sprintf("=SUM(G%d:G%d)", rowBefore, row-1))
	f.SetCellFormula(sheet, fmt.Sprintf("H%d", row), fmt.Sprintf("=SUM(H%d:H%d)", rowBefore, row-1))
	f.SetCellFormula(sheet, fmt.Sprintf("I%d", row), fmt.Sprintf("=SUM(I%d:I%d)", rowBefore, row-1))
	f.SetCellFormula(sheet, fmt.Sprintf("J%d", row), fmt.Sprintf("=SUM(J%d:J%d)", rowBefore, row-1))

	f.SetCellFormula(sheet, fmt.Sprintf("I%d", row+1), fmt.Sprintf("=G%d+I%d", row, row))
	f.SetCellFormula(sheet, fmt.Sprintf("J%d", row+1), fmt.Sprintf("=H%d+J%d", row, row))
	f.SetCellFormula(sheet, fmt.Sprintf("J%d", row+2), fmt.Sprintf("=I%d-J%d", row+1, row+1))

	f.SetCellStyle(sheet, fmt.Sprintf("D%d", row), fmt.Sprintf("J%d", row), stylingSubTotal)
	f.SetCellStyle(sheet, fmt.Sprintf("G%d", row+1), fmt.Sprintf("J%d", row+1), stylingSubTotal2)

	f.SetDefaultFont("Arial")

	return f, nil
}

func (s *service) ExportMutasiFas(ctx *abstraction.Context, f *excelize.File, mutasiFa *model.MutasiFaEntityModel) (*excelize.File, error) {
	sheet := "MUTASI_FA"
	f.NewSheet(sheet)

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
	stylingControl, err := f.NewStyle(&excelize.Style{
		NumFmt: 41,
		Border: []excelize.Border{
			{Type: "top", Color: "000000", Style: 1},
			{Type: "bottom", Color: "000000", Style: 1},
			{Type: "left", Color: "000000", Style: 1},
			{Type: "right", Color: "000000", Style: 1},
		},
		Fill: excelize.Fill{
			Type:    "pattern",
			Pattern: 1,
			Color:   []string{"#FFFF00"},
		},
	})
	if err != nil {
		return nil, err
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
		NumFmt: 41,
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
		NumFmt: 41,
	})
	if err != nil {
		return nil, err
	}

	datePeriod, err := time.Parse(time.RFC3339, mutasiFa.Period)
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
	f.SetCellValue(sheet, "B4", mutasiFa.Company.Name)
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
			return nil, response.ErrorBuilder(&response.ErrorConstant.InternalServerError, err)
		}

		tmpStr := "MUTASI-FA"
		criteriaBridge := model.FormatterBridgesFilterModel{}
		criteriaBridge.FormatterBridgesFilter.Source = &tmpStr
		criteriaBridge.FormatterBridgesFilter.FormatterID = &data.ID
		criteriaBridge.FormatterBridgesFilter.TrxRefID = &mutasiFa.ID

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
			if v.ControlFormula != "" {

				if v.AutoSummary != nil && *v.AutoSummary {
					f.SetCellStyle(sheet, fmt.Sprintf("B%d", row), fmt.Sprintf("C%d", row), styleLabelTotal)
					f.SetCellStyle(sheet, fmt.Sprintf("D%d", row), fmt.Sprintf("K%d", row), styleCurrencyTotal)
					f.MergeCell(sheet, fmt.Sprintf("B%d", row), fmt.Sprintf("C%d", row))
					if strings.ToUpper(v.Code) == "NET_BOOK_VALUE" {
						f.SetCellFormula(sheet, fmt.Sprintf("D%d", row), fmt.Sprintf("=SUM(D%d:D%d)", partRowStart, row-1))
						f.SetCellFormula(sheet, fmt.Sprintf("J%d", row), fmt.Sprintf("=SUM(J%d:J%d)", partRowStart, row-1))
					} else {
						for chr := 'D'; chr <= 'J'; chr++ {
							f.SetCellFormula(sheet, fmt.Sprintf("%c%d", chr, row), fmt.Sprintf("=SUM(%c%d:%c%d)", chr, partRowStart, chr, row-1))
						}
					}

				}
				f.SetCellStyle(sheet, fmt.Sprintf("K%d", row), fmt.Sprintf("K%d", row), stylingControl)

				formula := v.FxSummary
				reg := regexp.MustCompile(`[A-Za-z0-9_~#:()]+|[0-9]+\d{3,}`)

				match := reg.FindAllString(formula, -1)
				for _, vMatch := range match {
					// cari jml berdasarkan code
					if _, ok := rowCode[vMatch]; ok {
						formula = strings.ReplaceAll(formula, vMatch, fmt.Sprintf("J%d", rowCode[vMatch]))
					}
					if _, ok := tbRowCodes[vMatch]; ok {
						formula = strings.ReplaceAll(formula, vMatch, fmt.Sprintf("L%d", tbRowCodes[vMatch]))
						f.SetCellValue("TRIAL_BALANCE", fmt.Sprintf("M%d", tbRowCodes[vMatch]), "control")
					}

				}
				f.SetCellFormula(sheet, fmt.Sprintf("K%d", row), fmt.Sprintf("=%s", formula))

			}
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
				if strings.ToUpper(v.Code) == "CONTROL_1" {
					arrChr = []string{"F"}
					f.SetCellStyle(sheet, fmt.Sprintf("B%d", row), fmt.Sprintf("K%d", row), styleCurrency)
					f.SetCellStyle(sheet, fmt.Sprintf("F%d", row), fmt.Sprintf("F%d", row), stylingControl)
				}
				if strings.ToUpper(v.Code) == "CONTROL_2" {
					arrChr = []string{"J"}
					f.SetCellStyle(sheet, fmt.Sprintf("B%d", row), fmt.Sprintf("K%d", row), styleCurrencyWoBorder)

				}

				for _, chr := range arrChr {
					formula := v.FxSummary
					reg := regexp.MustCompile(`[A-Za-z0-9_~#:()]+|[0-9]+\d{3,}`)
					match := reg.FindAllString(formula, -1)
					for _, vMatch := range match {
						// cari jml berdasarkan code
						if rowCode[vMatch] != 0 {
							formula = strings.ReplaceAll(formula, vMatch, fmt.Sprintf("%s%d", chr, rowCode[vMatch]))
						}
						if _, ok := tbRowCodes[vMatch]; ok {
							formula = strings.ReplaceAll(formula, vMatch, fmt.Sprintf("L%d", tbRowCodes[vMatch]))
							f.SetCellValue("TRIAL_BALANCE", fmt.Sprintf("M%d", tbRowCodes[vMatch]), "control")
						}

					}
					if strings.ToUpper(v.Code) == "CONTROL_2" {
						row = row - 1
						f.SetCellStyle(sheet, fmt.Sprintf("K%d", row), fmt.Sprintf("K%d", row), stylingControl)
						f.SetCellFormula(sheet, fmt.Sprintf("%s%d", "K", row), fmt.Sprintf("=%s", formula))
						continue

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
					for chr := 'D'; chr <= 'J'; chr++ {
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
			criteriaMF.MutasiFaID = &mutasiFa.ID

			paginationMF := abstraction.Pagination{}
			pagesize := 10000
			paginationMF.PageSize = &pagesize

			MutasiFaDetail, _, err := s.MutasiFaDetailRepository.Find(ctx, &criteriaMF, &paginationMF)
			if err != nil && v.Code != "" {
				continue
			}
			if len(*MutasiFaDetail) == 0 {
				continue
			}
			for _, vv := range *MutasiFaDetail {
				if vv.Code == "COST:" || vv.Code == "ACCUMULATED_DEPRECIATION" {
					f.SetCellValue(sheet, fmt.Sprintf("B%d", row), vv.Description)
					continue
				}
				valueKsong := 0.0
				f.SetCellValue(sheet, fmt.Sprintf("B%d", row), vv.Description)
				f.SetCellValue(sheet, fmt.Sprintf("D%d", row), valueKsong)
				f.SetCellValue(sheet, fmt.Sprintf("E%d", row), valueKsong)
				f.SetCellValue(sheet, fmt.Sprintf("F%d", row), valueKsong)
				f.SetCellValue(sheet, fmt.Sprintf("G%d", row), valueKsong)
				f.SetCellValue(sheet, fmt.Sprintf("H%d", row), valueKsong)
				f.SetCellValue(sheet, fmt.Sprintf("I%d", row), valueKsong)
				f.SetCellStyle(sheet, fmt.Sprintf("D%d", row), fmt.Sprintf("I%d", row), styleCurrency)

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

				f.SetCellFormula(sheet, fmt.Sprintf("J%d", row), fmt.Sprintf("=D%d+E%d+F%d-G%d+H%d+I%d", row, row, row, row, row, row))
				f.SetCellStyle(sheet, fmt.Sprintf("J%d", row), fmt.Sprintf("J%d", row), styleCurrencyWoBorder)

			}

			row++
		}

		rowStart = row
		row = rowStart
	}

	// Penambahan detail pengurangan
	row += 2
	var criteria dto.FormatterGetRequest
	tmpStr := "MUTASI-DETAIL-PENGURANGAN"
	criteria.FormatterFilterModel.FormatterFor = &tmpStr

	data, err := s.FormatterRepository.FindWithDetail(ctx, &criteria.FormatterFilterModel)
	if err != nil {
		return nil, response.ErrorBuilder(&response.ErrorConstant.InternalServerError, err)
	}

	criteriaBridge := model.FormatterBridgesFilterModel{}
	criteriaBridge.FormatterBridgesFilter.Source = &tmpStr
	criteriaBridge.FormatterBridgesFilter.FormatterID = &data.ID
	criteriaBridge.FormatterBridgesFilter.TrxRefID = &mutasiFa.ID

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
			if strings.ToUpper(v.Code) == "CONTROL_1" {
				arrChr = []string{"E"}
				f.SetCellStyle(sheet, fmt.Sprintf("E%d", row), fmt.Sprintf("E%d", row), styleCurrency)
				f.SetCellStyle(sheet, fmt.Sprintf("E%d", row), fmt.Sprintf("E%d", row), stylingControl)
			}

			if strings.ToUpper(v.Code) == "CONTROL_2" {
				arrChr = []string{"F"}
			}
			for _, chr := range arrChr {
				formula := v.FxSummary
				reg := regexp.MustCompile(`[A-Za-z0-9_~#:()]+|[0-9]+\d{3,}`)
				match := reg.FindAllString(formula, -1)
				for _, vMatch := range match {
					// cari jml berdasarkan code
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
					if _, ok := tbRowCodes[vMatch]; ok {
						formula = strings.ReplaceAll(formula, vMatch, fmt.Sprintf("L%d", tbRowCodes[vMatch]))
						f.SetCellValue("TRIAL_BALANCE", fmt.Sprintf("M%d", tbRowCodes[vMatch]), "control")
					}
				}
				if strings.ToUpper(v.Code) == "CONTROL_2" {
					row = row - 1
					f.SetCellStyle(sheet, fmt.Sprintf("F%d", row), fmt.Sprintf("F%d", row), styleCurrency)
					f.SetCellStyle(sheet, fmt.Sprintf("F%d", row), fmt.Sprintf("F%d", row), stylingControl)
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
		criteriaMF.MutasiFaID = &mutasiFa.ID
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
			valueKosong := 0.0
			f.SetCellValue(sheet, fmt.Sprintf("B%d", row), vv.Description)
			f.SetCellValue(sheet, fmt.Sprintf("E%d", row), valueKosong)
			f.SetCellValue(sheet, fmt.Sprintf("F%d", row), valueKosong)
			if vv.Deductions != nil && *vv.Deductions != 0 {
				f.SetCellValue(sheet, fmt.Sprintf("G%d", row), *vv.Deductions)
			}
			if vv.Revaluation != nil && *vv.Revaluation != 0 {
				f.SetCellValue(sheet, fmt.Sprintf("I%d", row), *vv.Revaluation)
			}
		}
		row++
	}

	return f, nil
}

func (s *service) ExportMutasiIas(ctx *abstraction.Context, f *excelize.File, mutasiIa *model.MutasiIaEntityModel) (*excelize.File, error) {
	sheet := "MUTASI_IA"
	f.NewSheet(sheet)

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
	stylingControl, err := f.NewStyle(&excelize.Style{
		NumFmt: 41,
		Border: []excelize.Border{
			{Type: "top", Color: "000000", Style: 1},
			{Type: "bottom", Color: "000000", Style: 1},
			{Type: "left", Color: "000000", Style: 1},
			{Type: "right", Color: "000000", Style: 1},
		},
		Fill: excelize.Fill{
			Type:    "pattern",
			Pattern: 1,
			Color:   []string{"#FFFF00"},
		},
	})
	if err != nil {
		return nil, err
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
		NumFmt: 41,
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
		NumFmt: 41,
	})
	if err != nil {
		return nil, err
	}

	allowed := helper.CompanyValidation(ctx.Auth.ID, mutasiIa.CompanyID)
	if !allowed {
		return nil, response.ErrorBuilder(&response.ErrorConstant.UnauthorizedAccess, errors.New("not allowed"))
	}

	datePeriod, err := time.Parse(time.RFC3339, mutasiIa.Period)
	if err != nil {
		return nil, err
	}

	formatterCode := []string{"MUTASI-IA-COST", "MUTASI-IA-ACCUMULATED-DEPRECATION"}
	formatterTitle := []string{"Mutasi Intangible Assets (IA)", ""}
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
	f.SetCellValue(sheet, "B4", mutasiIa.Company.Name)
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
			return nil, helper.ErrorHandler(err)
		}

		tmpStr := "MUTASI-IA"
		criteriaBridge := model.FormatterBridgesFilterModel{}
		criteriaBridge.FormatterBridgesFilter.Source = &tmpStr
		criteriaBridge.FormatterBridgesFilter.FormatterID = &data.ID
		criteriaBridge.FormatterBridgesFilter.TrxRefID = &mutasiIa.ID

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
			if v.ControlFormula != "" {

				if v.AutoSummary != nil && *v.AutoSummary {
					f.SetCellStyle(sheet, fmt.Sprintf("B%d", row), fmt.Sprintf("C%d", row), styleLabelTotal)
					f.SetCellStyle(sheet, fmt.Sprintf("D%d", row), fmt.Sprintf("K%d", row), styleCurrencyTotal)
					f.MergeCell(sheet, fmt.Sprintf("B%d", row), fmt.Sprintf("C%d", row))
					if strings.ToUpper(v.Code) == "NET_BOOK_VALUE" {
						f.SetCellFormula(sheet, fmt.Sprintf("D%d", row), fmt.Sprintf("=SUM(D%d:D%d)", partRowStart, row-1))
						f.SetCellFormula(sheet, fmt.Sprintf("J%d", row), fmt.Sprintf("=SUM(J%d:J%d)", partRowStart, row-1))
					} else {
						for chr := 'D'; chr <= 'J'; chr++ {
							f.SetCellFormula(sheet, fmt.Sprintf("%c%d", chr, row), fmt.Sprintf("=SUM(%c%d:%c%d)", chr, partRowStart, chr, row-1))
						}
					}

				}
				f.SetCellStyle(sheet, fmt.Sprintf("K%d", row), fmt.Sprintf("K%d", row), stylingControl)

				formula := v.FxSummary
				reg := regexp.MustCompile(`[A-Za-z0-9_~#:()]+|[0-9]+\d{3,}`)

				match := reg.FindAllString(formula, -1)
				for _, vMatch := range match {
					// cari jml berdasarkan code
					if _, ok := rowCode[vMatch]; ok {
						formula = strings.ReplaceAll(formula, vMatch, fmt.Sprintf("J%d", rowCode[vMatch]))
					}
					if _, ok := tbRowCodes[vMatch]; ok {
						formula = strings.ReplaceAll(formula, vMatch, fmt.Sprintf("L%d", tbRowCodes[vMatch]))
						f.SetCellValue("TRIAL_BALANCE", fmt.Sprintf("M%d", tbRowCodes[vMatch]), "control")
					}

				}
				f.SetCellFormula(sheet, fmt.Sprintf("K%d", row), fmt.Sprintf("=%s", formula))

			}
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
				if strings.ToUpper(v.Code) == "CONTROL_1" {
					arrChr = []string{"F"}
					f.SetCellStyle(sheet, fmt.Sprintf("B%d", row), fmt.Sprintf("K%d", row), styleCurrency)
					f.SetCellStyle(sheet, fmt.Sprintf("F%d", row), fmt.Sprintf("F%d", row), stylingControl)
				}
				if strings.ToUpper(v.Code) == "CONTROL_2" {
					arrChr = []string{"J"}
					f.SetCellStyle(sheet, fmt.Sprintf("B%d", row), fmt.Sprintf("K%d", row), styleCurrencyWoBorder)

				}

				for _, chr := range arrChr {
					formula := v.FxSummary
					reg := regexp.MustCompile(`[A-Za-z0-9_~#:()]+|[0-9]+\d{3,}`)
					match := reg.FindAllString(formula, -1)
					for _, vMatch := range match {
						// cari jml berdasarkan code
						if rowCode[vMatch] != 0 {
							formula = strings.ReplaceAll(formula, vMatch, fmt.Sprintf("%s%d", chr, rowCode[vMatch]))
						}
						if _, ok := tbRowCodes[vMatch]; ok {
							formula = strings.ReplaceAll(formula, vMatch, fmt.Sprintf("L%d", tbRowCodes[vMatch]))
							f.SetCellValue("TRIAL_BALANCE", fmt.Sprintf("M%d", tbRowCodes[vMatch]), "control")
						}

					}
					if strings.ToUpper(v.Code) == "CONTROL_2" {
						row = row - 1
						f.SetCellStyle(sheet, fmt.Sprintf("K%d", row), fmt.Sprintf("K%d", row), stylingControl)
						f.SetCellFormula(sheet, fmt.Sprintf("%s%d", "K", row), fmt.Sprintf("=%s", formula))
						continue

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
					for chr := 'D'; chr <= 'J'; chr++ {
						f.SetCellFormula(sheet, fmt.Sprintf("%c%d", chr, row), fmt.Sprintf("=SUM(%c%d:%c%d)", chr, partRowStart, chr, row-1))
					}
				}
				row++
				partRowStart = row
				continue
			}

			criteriaMI := model.MutasiIaDetailFilterModel{}
			criteriaMI.Code = &v.Code
			criteriaMI.FormatterBridgesID = &bridges.ID
			criteriaMI.MutasiIaID = &mutasiIa.ID

			paginationMI := abstraction.Pagination{}
			pagesize := 10000
			paginationMI.PageSize = &pagesize

			MutasiIaDetail, _, err := s.MutasiIaDetailRepository.Find(ctx, &criteriaMI, &paginationMI)
			if err != nil {
				return nil, helper.ErrorHandler(err)
			}
			if len(*MutasiIaDetail) == 0 {
				continue
			}
			for _, vv := range *MutasiIaDetail {
				if vv.Code == "COST:" || vv.Code == "ACCUMULATED_DEPRECIATION" {
					f.SetCellValue(sheet, fmt.Sprintf("B%d", row), vv.Description)
					continue
				}
				valueKsong := 0.0
				f.SetCellValue(sheet, fmt.Sprintf("B%d", row), vv.Description)
				f.SetCellValue(sheet, fmt.Sprintf("D%d", row), valueKsong)
				f.SetCellValue(sheet, fmt.Sprintf("E%d", row), valueKsong)
				f.SetCellValue(sheet, fmt.Sprintf("F%d", row), valueKsong)
				f.SetCellValue(sheet, fmt.Sprintf("G%d", row), valueKsong)
				f.SetCellValue(sheet, fmt.Sprintf("H%d", row), valueKsong)
				f.SetCellValue(sheet, fmt.Sprintf("I%d", row), valueKsong)
				f.SetCellStyle(sheet, fmt.Sprintf("D%d", row), fmt.Sprintf("I%d", row), styleCurrency)

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
				f.SetCellFormula(sheet, fmt.Sprintf("J%d", row), fmt.Sprintf("=D%d+E%d+F%d-G%d+H%d+I%d", row, row, row, row, row, row))
				if vv.Control != nil && *vv.Control != 0 {
					f.SetCellValue(sheet, fmt.Sprintf("J%d", row), *vv.Control)
				}
			}

			row++
		}
		rowStart = row
		row = rowStart
	}

	// Penambahan detail pengurangan
	row += 2
	var criteria dto.FormatterGetRequest
	tmpStr := "MUTASI-DETAIL-PENGURANGAN"
	criteria.FormatterFilterModel.FormatterFor = &tmpStr

	data, err := s.FormatterRepository.FindWithDetail(ctx, &criteria.FormatterFilterModel)
	if err != nil {
		return nil, err
	}

	criteriaBridge := model.FormatterBridgesFilterModel{}
	criteriaBridge.FormatterBridgesFilter.Source = &tmpStr
	criteriaBridge.FormatterBridgesFilter.FormatterID = &data.ID
	criteriaBridge.FormatterBridgesFilter.TrxRefID = &mutasiIa.ID

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
					// cari jml berdasarkan code
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

		criteriaMI := model.MutasiIaDetailFilterModel{}
		criteriaMI.Code = &v.Code
		criteriaMI.FormatterBridgesID = &bridges.ID
		criteriaMI.MutasiIaID = &mutasiIa.ID
		paginationMI := abstraction.Pagination{}
		pagesize := 10000
		paginationMI.PageSize = &pagesize
		MutasiIaDetail, _, err := s.MutasiIaDetailRepository.Find(ctx, &criteriaMI, &paginationMI)
		if err != nil {
			row++
			continue
		}
		if len(*MutasiIaDetail) == 0 {
			row++
			continue
		}
		for _, vv := range *MutasiIaDetail {
			valueKosong := 0.0
			f.SetCellValue(sheet, fmt.Sprintf("B%d", row), vv.Description)
			f.SetCellValue(sheet, fmt.Sprintf("E%d", row), valueKosong)
			f.SetCellValue(sheet, fmt.Sprintf("F%d", row), valueKosong)
			if vv.Deductions != nil && *vv.Deductions != 0 {
				f.SetCellValue(sheet, fmt.Sprintf("G%d", row), *vv.Deductions)
			}
			if vv.Revaluation != nil && *vv.Revaluation != 0 {
				f.SetCellValue(sheet, fmt.Sprintf("I%d", row), *vv.Revaluation)
			}
		}
		row++
	}

	return f, nil
}

func (s *service) ExportMutasiRuas(ctx *abstraction.Context, f *excelize.File, mutasiRua *model.MutasiRuaEntityModel) (*excelize.File, error) {
	sheet := "MUTASI_RUA"
	f.NewSheet(sheet)

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
	stylingControl, err := f.NewStyle(&excelize.Style{
		NumFmt: 41,
		Border: []excelize.Border{
			{Type: "top", Color: "000000", Style: 1},
			{Type: "bottom", Color: "000000", Style: 1},
			{Type: "left", Color: "000000", Style: 1},
			{Type: "right", Color: "000000", Style: 1},
		},
		Fill: excelize.Fill{
			Type:    "pattern",
			Pattern: 1,
			Color:   []string{"#FFFF00"},
		},
	})
	if err != nil {
		return nil, err
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
		NumFmt: 41,
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
		NumFmt: 41,
	})
	if err != nil {
		return nil, err
	}
	datePeriod, err := time.Parse(time.RFC3339, mutasiRua.Period)
	if err != nil {
		return nil, err
	}

	formatterCode := []string{"MUTASI-RUA-COST", "MUTASI-RUA-ACCUMULATED-DEPRECATION"}
	formatterTitle := []string{"Mutasi Right of Used Assets (RUA)"}
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
	f.SetCellValue(sheet, "B4", mutasiRua.Company.Name)
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
			return nil, helper.ErrorHandler(err)
		}

		tmpStr := "MUTASI-RUA"
		criteriaBridge := model.FormatterBridgesFilterModel{}
		criteriaBridge.FormatterBridgesFilter.Source = &tmpStr
		criteriaBridge.FormatterBridgesFilter.FormatterID = &data.ID
		criteriaBridge.FormatterBridgesFilter.TrxRefID = &mutasiRua.ID

		bridges, err := s.FormatterBridgesRepository.FindWithCriteria(ctx, &criteriaBridge)
		if err != nil {
			return nil, helper.ErrorHandler(err)
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
			if v.ControlFormula != "" {

				if v.AutoSummary != nil && *v.AutoSummary {
					f.SetCellStyle(sheet, fmt.Sprintf("B%d", row), fmt.Sprintf("C%d", row), styleLabelTotal)
					f.SetCellStyle(sheet, fmt.Sprintf("D%d", row), fmt.Sprintf("K%d", row), styleCurrencyTotal)
					f.MergeCell(sheet, fmt.Sprintf("B%d", row), fmt.Sprintf("C%d", row))
					if strings.ToUpper(v.Code) == "NET_BOOK_VALUE" {
						f.SetCellFormula(sheet, fmt.Sprintf("D%d", row), fmt.Sprintf("=SUM(D%d:D%d)", partRowStart, row-1))
						f.SetCellFormula(sheet, fmt.Sprintf("J%d", row), fmt.Sprintf("=SUM(J%d:J%d)", partRowStart, row-1))
					} else {
						for chr := 'D'; chr <= 'J'; chr++ {
							f.SetCellFormula(sheet, fmt.Sprintf("%c%d", chr, row), fmt.Sprintf("=SUM(%c%d:%c%d)", chr, partRowStart, chr, row-1))
						}
					}

				}
				f.SetCellStyle(sheet, fmt.Sprintf("K%d", row), fmt.Sprintf("K%d", row), stylingControl)
				// f.SetCellStyle(sheet, fmt.Sprintf("K%d", row), fmt.Sprintf("K%d", row), styleCurrency)
				formula := v.FxSummary
				reg := regexp.MustCompile(`[A-Za-z0-9_~#:()]+|[0-9]+\d{3,}`)

				match := reg.FindAllString(formula, -1)
				for _, vMatch := range match {
					// cari jml berdasarkan code
					if _, ok := rowCode[vMatch]; ok {
						formula = strings.ReplaceAll(formula, vMatch, fmt.Sprintf("J%d", rowCode[vMatch]))
					}
					if _, ok := tbRowCodes[vMatch]; ok {
						formula = strings.ReplaceAll(formula, vMatch, fmt.Sprintf("L%d", tbRowCodes[vMatch]))
						f.SetCellValue("TRIAL_BALANCE", fmt.Sprintf("M%d", tbRowCodes[vMatch]), "control")
					}

				}
				f.SetCellFormula(sheet, fmt.Sprintf("K%d", row), fmt.Sprintf("=%s", formula))

			}
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
				if strings.ToUpper(v.Code) == "CONTROL_1" {
					arrChr = []string{"F"}
					f.SetCellStyle(sheet, fmt.Sprintf("B%d", row), fmt.Sprintf("K%d", row), styleCurrency)
					f.SetCellStyle(sheet, fmt.Sprintf("F%d", row), fmt.Sprintf("F%d", row), stylingControl)
				}
				if strings.ToUpper(v.Code) == "CONTROL_2" {
					arrChr = []string{"J"}
					f.SetCellStyle(sheet, fmt.Sprintf("B%d", row), fmt.Sprintf("K%d", row), styleCurrencyWoBorder)

				}

				for _, chr := range arrChr {
					formula := v.FxSummary
					reg := regexp.MustCompile(`[A-Za-z0-9_~#:()]+|[0-9]+\d{3,}`)
					match := reg.FindAllString(formula, -1)
					for _, vMatch := range match {
						// cari jml berdasarkan code
						if rowCode[vMatch] != 0 {
							formula = strings.ReplaceAll(formula, vMatch, fmt.Sprintf("%s%d", chr, rowCode[vMatch]))
						}
						if _, ok := tbRowCodes[vMatch]; ok {
							formula = strings.ReplaceAll(formula, vMatch, fmt.Sprintf("L%d", tbRowCodes[vMatch]))
							f.SetCellValue("TRIAL_BALANCE", fmt.Sprintf("M%d", tbRowCodes[vMatch]), "control")
						}

					}
					if strings.ToUpper(v.Code) == "CONTROL_2" {
						row = row - 1
						f.SetCellStyle(sheet, fmt.Sprintf("K%d", row), fmt.Sprintf("K%d", row), stylingControl)
						f.SetCellFormula(sheet, fmt.Sprintf("%s%d", "K", row), fmt.Sprintf("=%s", formula))
						continue

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
					for chr := 'D'; chr <= 'J'; chr++ {
						f.SetCellFormula(sheet, fmt.Sprintf("%c%d", chr, row), fmt.Sprintf("=SUM(%c%d:%c%d)", chr, partRowStart, chr, row-1))
					}
				}
				row++
				partRowStart = row
				continue
			}

			criteriaMR := model.MutasiRuaDetailFilterModel{}
			criteriaMR.Code = &v.Code
			criteriaMR.FormatterBridgesID = &bridges.ID
			criteriaMR.MutasiRuaID = &mutasiRua.ID
			paginationMR := abstraction.Pagination{}
			pagesize := 1
			paginationMR.PageSize = &pagesize

			MRuaDetail, _, err := s.MutasiRuaDetailRepository.Find(ctx, &criteriaMR, &paginationMR)
			if err != nil && v.Code != "" {
				continue
			}
			if len(*MRuaDetail) == 0 {
				continue
			}
			for _, vv := range *MRuaDetail {
				if vv.Code == "COST:" || vv.Code == "ACCUMULATED_DEPRECIATION" {
					f.SetCellValue(sheet, fmt.Sprintf("B%d", row), vv.Description)
					continue
				}
				valueKsong := 0.0
				f.SetCellValue(sheet, fmt.Sprintf("B%d", row), vv.Description)
				f.SetCellValue(sheet, fmt.Sprintf("D%d", row), valueKsong)
				f.SetCellValue(sheet, fmt.Sprintf("E%d", row), valueKsong)
				f.SetCellValue(sheet, fmt.Sprintf("F%d", row), valueKsong)
				f.SetCellValue(sheet, fmt.Sprintf("G%d", row), valueKsong)
				f.SetCellValue(sheet, fmt.Sprintf("H%d", row), valueKsong)
				f.SetCellValue(sheet, fmt.Sprintf("I%d", row), valueKsong)
				f.SetCellStyle(sheet, fmt.Sprintf("D%d", row), fmt.Sprintf("I%d", row), styleCurrency)

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
				if vv.Remeasurement != nil && *vv.Remeasurement != 0 {
					f.SetCellValue(sheet, fmt.Sprintf("I%d", row), *vv.Remeasurement)
				}
				f.SetCellFormula(sheet, fmt.Sprintf("J%d", row), fmt.Sprintf("=D%d+E%d+F%d-G%d+H%d+I%d", row, row, row, row, row, row))
				if vv.Control != nil && *vv.Control != 0 {
					f.SetCellValue(sheet, fmt.Sprintf("J%d", row), *vv.Control)
				}
			}
			row++
		}
		rowStart = row
		row = rowStart
	}

	// Penambahan detail pengurangan
	row += 2
	var criteria dto.FormatterGetRequest
	tmpStr := "MUTASI-DETAIL-PENGURANGAN"
	criteria.FormatterFilterModel.FormatterFor = &tmpStr

	data, err := s.FormatterRepository.FindWithDetail(ctx, &criteria.FormatterFilterModel)
	if err != nil {
		return nil, err
	}

	criteriaBridge := model.FormatterBridgesFilterModel{}
	criteriaBridge.FormatterBridgesFilter.Source = &tmpStr
	criteriaBridge.FormatterBridgesFilter.FormatterID = &data.ID
	criteriaBridge.FormatterBridgesFilter.TrxRefID = &mutasiRua.ID

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
					// cari jml berdasarkan code
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

		criteriaMR := model.MutasiRuaDetailFilterModel{}
		criteriaMR.Code = &v.Code
		criteriaMR.FormatterBridgesID = &bridges.ID
		criteriaMR.MutasiRuaID = &mutasiRua.ID
		paginationMR := abstraction.Pagination{}
		pagesize := 10000
		paginationMR.PageSize = &pagesize
		MutasiRuaDetail, _, err := s.MutasiRuaDetailRepository.Find(ctx, &criteriaMR, &paginationMR)
		if err != nil {
			row++
			continue
		}
		if len(*MutasiRuaDetail) == 0 {
			row++
			continue
		}
		for _, vv := range *MutasiRuaDetail {
			valueKosong := 0.0
			f.SetCellValue(sheet, fmt.Sprintf("B%d", row), vv.Description)
			f.SetCellValue(sheet, fmt.Sprintf("E%d", row), valueKosong)
			f.SetCellValue(sheet, fmt.Sprintf("F%d", row), valueKosong)
			if vv.Deductions != nil && *vv.Deductions != 0 {
				f.SetCellValue(sheet, fmt.Sprintf("G%d", row), *vv.Deductions)
			}
			if vv.Remeasurement != nil && *vv.Remeasurement != 0 {
				f.SetCellValue(sheet, fmt.Sprintf("I%d", row), *vv.Remeasurement)
			}
		}
		row++
	}

	return f, nil
}

var tbRowCodes = make(map[string]int)

func (s *service) ExportTrialBalances(ctx *abstraction.Context, f *excelize.File, trialBalance *model.TrialBalanceEntityModel) (*excelize.File, error) {
	sheet := "TRIAL_BALANCE"
	indexSheet := f.NewSheet(sheet)
	f.SetActiveSheet(indexSheet)

	tb, err := s.TrialBalanceRepository.Get(ctx, trialBalance.ID)
	if err != nil {
		return nil, err
	}

	if tb.ID == 0 {
		return nil, errors.New("no data found")
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

	t := true
	var criteriaFormatter model.FormatterDetailFilterModel
	criteriaFormatter.FormatterID = &formatterID
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

	f.SetCellValue(sheet, "B2", "Company")
	f.SetCellValue(sheet, "D2", ": "+tb.Company.CompanyEntity.Name)
	f.SetCellValue(sheet, "B3", "Date")
	f.SetCellValue(sheet, "D3", ": "+datePeriod.Format("02-Jan-06"))
	f.SetCellValue(sheet, "B4", "Subject")
	f.SetCellValue(sheet, "D4", ": DETAIL ASET, LIABILITAS & EKUITAS")

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
		if err = f.SetColWidth(sheet, fmt.Sprintf("%s", v["COL"]), fmt.Sprintf("%s", v["COL"]), colWidth); err != nil {
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

	if err = f.SetColStyle(sheet, "A:Z", styleDefault); err != nil {
		return nil, err
	}

	stylingBorderLeftOnly, err := f.NewStyle(&excelize.Style{
		Border: []excelize.Border{
			{Type: "left", Color: "000000", Style: 1},
		},
	})
	if err != nil {
		return nil, err
	}

	stylingBorderLROnly, err := f.NewStyle(&excelize.Style{
		Border: []excelize.Border{
			{Type: "right", Color: "000000", Style: 1},
			{Type: "left", Color: "000000", Style: 1},
		},
	})
	if err != nil {
		return nil, err
	}

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
		NumFmt: 41,
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
		NumFmt: 41,
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
		NumFmt: 41,
	})
	if err != nil {
		return nil, err
	}

	if err = f.MergeCell(sheet, "B6", "B8"); err != nil {
		return nil, err
	}
	if err = f.MergeCell(sheet, "C6", "E8"); err != nil {
		return nil, err
	}
	if err = f.MergeCell(sheet, "F6", "F8"); err != nil {
		return nil, err
	}
	if err = f.MergeCell(sheet, "H6", "K7"); err != nil {
		return nil, err
	}

	panes := `{
				"freeze":true, 
				"top_left_cell":"A1",
				"active_pane":"bottomRight",
				"panes":[{"active_cell":"G9","sqref":"G9","pane":"bottomRight"}]
			}`
	if err = f.SetPanes(sheet, panes); err != nil {
		fmt.Println("Error setting freeze panes:", err)
		return nil, err
	}

	stylingBorderRightOnly, err := f.NewStyle(&excelize.Style{
		Border: []excelize.Border{
			{Type: "right", Color: "000000", Style: 1},
		},
	})
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

	stylingCurrency2, err := f.NewStyle(&excelize.Style{
		NumFmt: 41,
		Border: []excelize.Border{
			{Type: "left", Color: "000000", Style: 1},
			{Type: "right", Color: "000000", Style: 1},
		},
	})
	if err != nil {
		return nil, err
	}

	if err = f.SetCellStyle(sheet, "B6", "L8", stylingHeader); err != nil {
		return nil, err
	}
	if err = f.SetCellValue(sheet, "B6", "No Akun"); err != nil {
		return nil, err
	}
	if err = f.SetCellValue(sheet, "C6", "Keterangan"); err != nil {
		return nil, err
	}
	if err = f.SetCellValue(sheet, "F6", "WP Reff"); err != nil {
		return nil, err
	}
	if err = f.SetCellValue(sheet, "G6", tb.Company.CompanyEntity.Code); err != nil {
		return nil, err
	}
	if err = f.SetCellValue(sheet, "G7", "Unaudited"); err != nil {
		return nil, err
	}
	if err = f.SetCellValue(sheet, "G8", datePeriod.Format("02-Jan-06")); err != nil {
		return nil, err
	}
	if err = f.SetCellValue(sheet, "H6", "Adjustment Journal Entry"); err != nil {
		return nil, err
	}
	if err = f.SetCellValue(sheet, "I8", "Debet"); err != nil {
		return nil, err
	}
	if err = f.SetCellValue(sheet, "K8", "Kredit"); err != nil {
		return nil, err
	}
	if err = f.SetCellFormula(sheet, "L6", "=G6"); err != nil {
		return nil, err
	}
	if err = f.SetCellFormula(sheet, "L7", "=G7"); err != nil {
		return nil, err
	}
	if err = f.SetCellFormula(sheet, "L8", "=G8"); err != nil {
		return nil, err
	}

	tmpStr2 := "TRIAL-BALANCE"
	criteriaBridge := model.FormatterBridgesFilterModel{}
	criteriaBridge.FormatterBridgesFilter.Source = &tmpStr2
	criteriaBridge.FormatterBridgesFilter.FormatterID = &formatterID
	criteriaBridge.FormatterBridgesFilter.TrxRefID = &tb.ID
	bridges, err := s.FormatterBridgesRepository.FindWithCriteria(ctx, &criteriaBridge)
	if err != nil {
		return nil, err
	}

	// find summary aje
	summaryAJE, err := s.AdjustmentRepository.FindSummary(ctx, &tb.ID)
	if err != nil {
		return nil, err
	}

	row := 9

	rowCode := make(map[string]int)
	isAutoSum := make(map[string]bool)
	customRow := make(map[string]string)
	fml := make(map[string]string)
	reff := make(map[string]string)

	sheet1 := "ADJUSTMENT"
	rows, err := f.GetRows(sheet1)
	if err != nil {
		return nil, err
	}

	var line []string
	for _, ro := range rows {
		line = append(line, "line")

		if len(ro) == 0 {
			continue
		}
		if len(ro) == 2 {
			continue
		}
		if len(ro) == 5 {
			continue
		}
		if ro[6] != "" {
			if fml[ro[1]] != "" {
				rumusAll := strings.Replace(fml[ro[1]], "=", "+", 1)
				rumus := "=" + sheet1 + "!G" + strconv.Itoa(len(line)) + rumusAll
				fmt.Sprintln(rumus)
				rumusReffAll := reff[ro[1]]
				rumusReff := rumusReffAll + "&" + `","` + "&" + sheet1 + "!C" + strconv.Itoa(len(line))
				// if ro[1] == ro[2] {
				// 	f.SetCellFormula(sheet, fmt.Sprintf("I%d", row), rumus)
				// 	f.SetCellFormula(sheet, fmt.Sprintf("H%d", row), rumusReff)
				// }
				fml[ro[1]] = rumus
				reff[ro[1]] = rumusReff
			}
			if fml[ro[1]] == "" {
				rumus := "=" + sheet1 + "!G" + strconv.Itoa(len(line))
				rumusReff := "=" + sheet1 + "!C" + strconv.Itoa(len(line))
				// if ro[1] == ro[2] {
				// 	f.SetCellFormula(sheet, fmt.Sprintf("I%d", row), rumus)
				// 	f.SetCellFormula(sheet, fmt.Sprintf("H%d", row), rumusReff)
				// }
				fml[ro[1]] = rumus
				reff[ro[1]] = rumusReff
			}
		}
		if len(ro) == 8 && ro[7] != "" {
			if fml[ro[1]] != "" {
				rumusAll := strings.Replace(fml[ro[1]], "=", "+", 1)
				fmt.Sprintln(rumusAll)
				rumus := "=" + sheet1 + "!H" + strconv.Itoa(len(line)) + rumusAll
				rumusReffAll := reff[ro[1]]
				rumusReff := rumusReffAll + "&" + `","` + "&" + sheet1 + "!C" + strconv.Itoa(len(line))
				fmt.Sprintln(rumus)
				// if ro[1] == ro[2] {
				// 	f.SetCellFormula(sheet, fmt.Sprintf("K%d", row), rumus)
				// 	f.SetCellFormula(sheet, fmt.Sprintf("J%d", row), rumusReff)
				// }
				fml[ro[1]] = rumus
				reff[ro[1]] = rumusReff
			}
			if fml[ro[1]] == "" {
				rumus := "=" + sheet1 + "!H" + strconv.Itoa(len(line))
				rumusReff := "=" + sheet1 + "!C" + strconv.Itoa(len(line))
				// if ro[1] == ro[2] {
				// 	f.SetCellFormula(sheet, fmt.Sprintf("K%d", row), rumus)
				// 	f.SetCellFormula(sheet, fmt.Sprintf("J%d", row), rumusReff)
				// }
				fml[ro[1]] = rumus
				reff[ro[1]] = rumusReff
			}
		}
		if len(ro) == 9 && ro[8] != "" {
			if fml[ro[1]] != "" {
				rumusAll := strings.Replace(fml[ro[1]], "=", "+", 1)
				fmt.Sprintln(rumusAll)
				rumus := "=" + sheet1 + "!I" + strconv.Itoa(len(line)) + rumusAll
				rumusReffAll := reff[ro[1]]
				rumusReff := rumusReffAll + "&" + `","` + "&" + sheet1 + "!C" + strconv.Itoa(len(line))
				fmt.Sprintln(rumus)
				// if ro[1] == ro[2] {
				// 	f.SetCellFormula(sheet, fmt.Sprintf("I%d", row), rumus)
				// 	f.SetCellFormula(sheet, fmt.Sprintf("H%d", row), rumusReff)
				// }
				fml[ro[1]] = rumus
				reff[ro[1]] = rumusReff
			}
			if fml[ro[1]] == "" {
				rumus := "=" + sheet1 + "!1" + strconv.Itoa(len(line))
				rumusReff := "=" + sheet1 + "!C" + strconv.Itoa(len(line))
				// if ro[1] == ro[2] {
				// 	f.SetCellFormula(sheet, fmt.Sprintf("I%d", row), rumus)
				// 	f.SetCellFormula(sheet, fmt.Sprintf("H%d", row), rumusReff)
				// }
				fml[ro[1]] = rumus
				reff[ro[1]] = rumusReff
			}
		}
		if len(ro) == 10 && ro[9] != "" {
			if fml[ro[1]] != "" {
				rumusAll := strings.Replace(fml[ro[1]], "=", "+", 1)
				fmt.Sprintln(rumusAll)
				rumus := "=" + sheet1 + "!J" + strconv.Itoa(len(line)) + rumusAll
				rumusReffAll := reff[ro[1]]
				rumusReff := rumusReffAll + "&" + `","` + "&" + sheet1 + "!C" + strconv.Itoa(len(line))
				fmt.Sprintln(rumus)
				// if ro[1] == ro[2] {
				// 	f.SetCellFormula(sheet, fmt.Sprintf("K%d", row), rumus)
				// 	f.SetCellFormula(sheet, fmt.Sprintf("J%d", row), rumusReff)
				// }
				fml[ro[1]] = rumus
				reff[ro[1]] = rumusReff
			}
			if fml[ro[1]] == "" {
				rumusReff := "=" + sheet1 + "!C" + strconv.Itoa(len(line))
				rumus := "=" + sheet1 + "!J" + strconv.Itoa(len(line))
				// if ro[1] == ro[2] {
				// 	f.SetCellFormula(sheet, fmt.Sprintf("K%d", row), rumus)
				// 	f.SetCellFormula(sheet, fmt.Sprintf("J%d", row), rumusReff)
				// }
				fml[ro[1]] = rumus
				reff[ro[1]] = rumusReff
			}
		}
	}

	for _, v := range *data {
		rowCode[v.Code] = row
		// tbRowCodes[v.Code] = row
		if err = f.SetCellStyle(sheet, fmt.Sprintf("F%d", row), fmt.Sprintf("F%d", row), stylingBorderLROnly); err != nil {
			return nil, err
		}
		if err = f.SetCellStyle(sheet, fmt.Sprintf("G%d", row), fmt.Sprintf("G%d", row), stylingCurrency2); err != nil {
			return nil, err
		}
		if err = f.SetCellStyle(sheet, fmt.Sprintf("I%d", row), fmt.Sprintf("I%d", row), stylingCurrency); err != nil {
			return nil, err
		}
		if err = f.SetCellStyle(sheet, fmt.Sprintf("K%d", row), fmt.Sprintf("K%d", row), stylingCurrency); err != nil {
			return nil, err
		}
		if err = f.SetCellStyle(sheet, fmt.Sprintf("L%d", row), fmt.Sprintf("L%d", row), stylingCurrency2); err != nil {
			return nil, err
		}
		if v.AutoSummary != nil && *v.AutoSummary {
			isAutoSum[v.Code] = true
		}
		if (v.IsTotal == nil || (v.IsTotal != nil && !*v.IsTotal)) && (v.IsLabel != nil && *v.IsLabel) {
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
			tbdetails, err := s.TrialBalanceDetailRepository.FindToExport(ctx, &v.Code, &bridges.ID)
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

				tbRowCodes[vTbDetail.Code] = row
				if vTbDetail.Code == "310401004" || vTbDetail.Code == "310402002" || vTbDetail.Code == "310501002" || vTbDetail.Code == "310502002" || vTbDetail.Code == "310503002" {
					f.SetCellValue(sheet, fmt.Sprintf("M%d", row), "control")
				}

				var (
					description     string
					amountBeforeAJE float64
					amountAfterAJE  float64
					amountDebet     float64
					amountCredit    float64
				)

				if vTbDetail.Description != nil {
					description = *vTbDetail.Description
				}
				if vTbDetail.AmountBeforeAje != nil {
					amountBeforeAJE = *vTbDetail.AmountBeforeAje
				}
				if vTbDetail.AmountAjeDr != nil {
					amountDebet = *vTbDetail.AmountAjeDr
				}
				if vTbDetail.AmountAjeCr != nil {
					amountCredit = *vTbDetail.AmountAjeCr
				}
				if vTbDetail.AmountAfterAje != nil {
					amountAfterAJE = *vTbDetail.AmountAfterAje
				}

				f.SetCellValue(sheet, fmt.Sprintf("B%d", row), vTbDetail.Code)
				f.SetCellValue(sheet, fmt.Sprintf("E%d", row), description)

				f.SetCellValue(sheet, fmt.Sprintf("G%d", row), amountBeforeAJE)
				f.SetCellValue(sheet, fmt.Sprintf("I%d", row), amountDebet)
				f.SetCellValue(sheet, fmt.Sprintf("K%d", row), amountCredit)
				f.SetCellValue(sheet, fmt.Sprintf("L%d", row), amountAfterAJE)

				if err = f.SetCellStyle(sheet, fmt.Sprintf("G%d", row), fmt.Sprintf("G%d", row), stylingCurrency2); err != nil {
					return nil, err
				}
				if err = f.SetCellStyle(sheet, fmt.Sprintf("I%d", row), fmt.Sprintf("I%d", row), stylingCurrency); err != nil {
					return nil, err
				}
				if err = f.SetCellStyle(sheet, fmt.Sprintf("K%d", row), fmt.Sprintf("K%d", row), stylingCurrency); err != nil {
					return nil, err
				}
				if err = f.SetCellStyle(sheet, fmt.Sprintf("L%d", row), fmt.Sprintf("L%d", row), stylingCurrency2); err != nil {
					return nil, err
				}

				if _, ok := fml[vTbDetail.Code]; ok {
					if vTbDetail.AmountAjeDr != nil && *vTbDetail.AmountAjeDr > 0 {
						f.SetCellFormula(sheet, fmt.Sprintf("H%d", row), reff[vTbDetail.Code])
						f.SetCellFormula(sheet, fmt.Sprintf("I%d", row), fml[vTbDetail.Code])
					}
					if vTbDetail.AmountAjeCr != nil && *vTbDetail.AmountAjeCr > 0 {
						f.SetCellFormula(sheet, fmt.Sprintf("J%d", row), reff[vTbDetail.Code])
						f.SetCellFormula(sheet, fmt.Sprintf("K%d", row), fml[vTbDetail.Code])
					}
				}

				tmpHeadCoa := fmt.Sprintf("%c", vTbDetail.Code[0])
				if tmpHeadCoa == "9" {
					tmpHeadCoa = vTbDetail.Code[:2]
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
			if v.AutoSummary != nil && *v.AutoSummary && len(*tbdetails) > 1 {
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
				tbRowCodes[fmt.Sprintf("%s_SUBTOTAL", v.Code)] = row
				row++
				if err = f.SetCellStyle(sheet, fmt.Sprintf("F%d", row), fmt.Sprintf("L%d", row), stylingBorderLROnly); err != nil {
					return nil, err
				}
			}
		}
		if v.IsTotal != nil && *v.IsTotal {
			tbRowCodes[v.Code] = row
			f.SetCellValue(sheet, fmt.Sprintf("C%d", row), v.Description)
			if v.Code == "TOTAL_INVESTASI_JANGKA_PENDEK" || v.Code == "TOTAL_PIUTANG_LAIN~LAIN_~_PIHAK_KETIGA~PIUTANG_USAHA" || v.Code == "TOTAL_CIP~ASET_TETAP" || v.Code == "TOTAL_PIUTANG_LAIN~LAIN" || v.Code == "TOTAL_CIP~ASET_TAK_BERWUJUD" || v.Code == "NCA_~_RUA_~_AKUMULASI_PENYUSUTAN_~_LAIN~LAIN" || v.Code == "TOTAL_INVESTASI_JANGKA_PANJANG" || v.Code == "TOTAL_PIUTANG_LAIN~LAIN_~_JANGKA_PANJANG" || v.Code == "TOTAL_ASET" || v.Code == "TOTAL_UTANG_USAHA" || v.Code == "TOTAL_UTANG_LAIN~LAIN_JANGKA_PENDEK" || v.Code == "TOTAL_LIABILITAS_IMBALAN_KERJA_~_JANGKA_PENDEK" || v.Code == "TOTAL_UTANG_LAIN~LAIN_JANGKA_PANJANG" {
				f.SetCellValue(sheet, fmt.Sprintf("M%d", row), "control")
			}
			if v.Code == "TOTAL_LIABILITAS_DAN_EKUITAS" {
				rowAset := row
				if _, ok := rowCode["TOTAL_ASET"]; ok {
					rowAset = rowCode["TOTAL_ASET"]
				}
				f.SetCellFormula(sheet, "G5", fmt.Sprintf("=G%d-G%d", rowAset, row))
				f.SetCellFormula(sheet, "L5", fmt.Sprintf("=L%d-L%d", rowAset, row))
			}

			// show control aje
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
	for key, nRow := range tbRowCodes {
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
		if val, ok := tbRowCodes[key]; ok {
			f.SetCellFormula(sheet, fmt.Sprintf("G%d", val), strings.ReplaceAll(vCustomRow, "@", "G"))
			f.SetCellFormula(sheet, fmt.Sprintf("I%d", val), strings.ReplaceAll(vCustomRow, "@", "I"))
			f.SetCellFormula(sheet, fmt.Sprintf("K%d", val), strings.ReplaceAll(vCustomRow, "@", "K"))
			// f.SetCellFormula(sheet, fmt.Sprintf("L%d", val), strings.ReplaceAll(vCustomRow, "@", "L"))
		}
	}

	if err = f.SetCellStyle(sheet, fmt.Sprintf("B%d", row), fmt.Sprintf("L%d", row), stylingBorderTopOnly); err != nil {
		return nil, err
	}
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

	return f, nil
}

func (s *service) ExportMutasiDtas(ctx *abstraction.Context, f *excelize.File, mutasiDta *model.MutasiDtaEntityModel) (*excelize.File, error) {
	sheet := "MUTASI_DTA"
	f.NewSheet(sheet)

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
	stylingControl, err := f.NewStyle(&excelize.Style{
		NumFmt: 41,
		Border: []excelize.Border{
			{Type: "top", Color: "000000", Style: 1},
			{Type: "bottom", Color: "000000", Style: 1},
			{Type: "left", Color: "000000", Style: 1},
			{Type: "right", Color: "000000", Style: 1},
		},
		Fill: excelize.Fill{
			Type:    "pattern",
			Pattern: 1,
			Color:   []string{"#FFFF00"},
		},
	})
	if err != nil {
		return nil, err
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
		NumFmt: 41,
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

	stylingDefault, err := f.NewStyle(&excelize.Style{
		NumFmt: 41,
		Fill: excelize.Fill{
			Type:    "pattern",
			Pattern: 1,
		},
		CustomNumFmt: &numberFormat,
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

	allowed := helper.CompanyValidation(ctx.Auth.ID, mutasiDta.CompanyID)
	if !allowed {
		return nil, response.ErrorBuilder(&response.ErrorConstant.UnauthorizedAccess, errors.New("not allowed"))
	}

	datePeriod, err := time.Parse(time.RFC3339, mutasiDta.Period)
	if err != nil {
		return nil, err
	}

	f.SetCellStyle(sheet, "B4", "J6", styleHeader)
	f.SetCellValue(sheet, "B4", "NO")
	f.SetCellValue(sheet, "C4", "Description")
	f.SetCellValue(sheet, "D4", fmt.Sprintf("%s Saldo Awal %s", mutasiDta.Company.Code, datePeriod.Format("02.01.06")))
	f.SetCellValue(sheet, "E4", "Penambahan (Pengurangan)")
	f.SetCellValue(sheet, "E5", "Manfaat (beban) pajak")
	f.SetCellValue(sheet, "F5", "OCI")
	f.SetCellValue(sheet, "G5", "Akuisisi Entitas anak")
	f.SetCellValue(sheet, "H4", "Dampak perubahan tariff pajak")
	f.SetCellValue(sheet, "H5", "Dibebankan ke laba rugi")
	f.SetCellValue(sheet, "I5", "Dibebankan ke OCI")
	f.SetCellValue(sheet, "J4", fmt.Sprintf("%s Saldo Akhir %s", mutasiDta.Company.Code, datePeriod.Format("02.01.06")))

	for _, formatter := range formatterCode {

		var criteria dto.FormatterGetRequest
		criteria.FormatterFilterModel.FormatterFor = &formatter

		data, err := s.FormatterRepository.FindWithDetail(ctx, &criteria.FormatterFilterModel)
		if err != nil {
			return nil, helper.ErrorHandler(err)
		}

		tmpStr := "MUTASI-DTA"
		criteriaBridge := model.FormatterBridgesFilterModel{}
		criteriaBridge.FormatterBridgesFilter.Source = &tmpStr
		criteriaBridge.FormatterBridgesFilter.FormatterID = &data.ID
		criteriaBridge.FormatterBridgesFilter.TrxRefID = &mutasiDta.ID

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

			f.SetCellValue(sheet, fmt.Sprintf("C%d", row), v.Description)

			if v.IsTotal != nil && *v.IsTotal {
				f.SetCellStyle(sheet, fmt.Sprintf("B%d", row), fmt.Sprintf("C%d", row), stylingDefault)
				f.SetCellStyle(sheet, fmt.Sprintf("D%d", row), fmt.Sprintf("J%d", row), stylingDefault)
				if v.FxSummary == "" {
					row++
					continue
				}
				arrChr := []string{"D", "E", "F", "G", "H", "I", "J", "K"}

				if strings.ToUpper(v.Code) == "CONTROL_1" {
					arrChr = []string{"E"}
					f.SetCellStyle(sheet, fmt.Sprintf("B%d", row), fmt.Sprintf("D%d", row), stylingDefault)
					f.SetCellStyle(sheet, fmt.Sprintf("G%d", row), fmt.Sprintf("I%d", row), stylingDefault)
				}
				if strings.ToUpper(v.Code) == "CONTROL_2" {
					arrChr = []string{"F"}
				}
				if strings.ToUpper(v.Code) == "CONTROL_3" {
					arrChr = []string{"J"}
				}
				for _, chr := range arrChr {
					formula := v.FxSummary
					reg := regexp.MustCompile(`[A-Za-z0-9_~#:()]+|[0-9]+\d{3,}`)
					match := reg.FindAllString(formula, -1)
					for _, vMatch := range match {
						// cari jml berdasarkan code
						if rowCode[vMatch] != 0 {
							formula = strings.ReplaceAll(formula, vMatch, fmt.Sprintf("%s%d", chr, rowCode[vMatch]))
						}
						if _, ok := tbRowCodes[vMatch]; ok {
							formula = strings.ReplaceAll(formula, vMatch, fmt.Sprintf("L%d", tbRowCodes[vMatch]))
							f.SetCellValue("TRIAL_BALANCE", fmt.Sprintf("M%d", tbRowCodes[vMatch]), "control")
						}

					}
					if strings.ToUpper(v.Code) == "CONTROL_2" {
						row = row - 1
						f.SetCellStyle(sheet, fmt.Sprintf("F%d", row), fmt.Sprintf("F%d", row), stylingControl)
					}
					if strings.ToUpper(v.Code) == "CONTROL_3" {
						row = row - 1

						f.SetCellStyle(sheet, fmt.Sprintf("J%d", row), fmt.Sprintf("J%d", row), stylingControl)
					}
					f.SetCellFormula(sheet, fmt.Sprintf("%s%d", chr, row), fmt.Sprintf("=%s", formula))
				}
				row++
				continue
			}

			if v.AutoSummary != nil && *v.AutoSummary {
				f.SetCellStyle(sheet, fmt.Sprintf("C%d", row), fmt.Sprintf("B%d", row), styleLabelTotal)
				f.SetCellStyle(sheet, fmt.Sprintf("D%d", row), fmt.Sprintf("J%d", row), styleCurrencyTotal)
				f.SetCellValue(sheet, fmt.Sprintf("C%d", row), v.Description)

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
			criteriaMD.MutasiDtaID = &mutasiDta.ID

			paginationMD := abstraction.Pagination{}
			pagesize := 1
			paginationMD.PageSize = &pagesize

			mDtaDetail, _, err := s.MutasiDtaDetailRepository.Find(ctx, &criteriaMD, &paginationMD)
			if err != nil && v.Code != "" {
				continue
			}
			if len(*mDtaDetail) == 0 {
				continue
			}
			for _, vv := range *mDtaDetail {
				valueKosong := 0.00
				f.SetCellValue(sheet, fmt.Sprintf("B%d", row), vv.SortID)
				f.SetCellValue(sheet, fmt.Sprintf("C%d", row), vv.Description)
				f.SetCellValue(sheet, fmt.Sprintf("D%d", row), valueKosong)
				f.SetCellValue(sheet, fmt.Sprintf("E%d", row), valueKosong)
				f.SetCellValue(sheet, fmt.Sprintf("F%d", row), valueKosong)
				f.SetCellValue(sheet, fmt.Sprintf("G%d", row), valueKosong)
				f.SetCellValue(sheet, fmt.Sprintf("H%d", row), valueKosong)
				f.SetCellValue(sheet, fmt.Sprintf("I%d", row), valueKosong)
				f.SetCellValue(sheet, fmt.Sprintf("J%d", row), valueKosong)
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
				f.SetCellStyle(sheet, fmt.Sprintf("J%d", row), fmt.Sprintf("J%d", row), styleCurrencyTotal)
			}
			row++
		}
		rowStart = row
		row = rowStart
	}

	return f, nil
}

func (s *service) ExportAgingUtangPiutangs(ctx *abstraction.Context, f *excelize.File, agingutangpiutang *model.AgingUtangPiutangEntityModel) (*excelize.File, error) {
	sheet := "AGING_UTANG_PIUTANG"
	f.NewSheet(sheet)
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
		NumFmt: 41,
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
	stylingControl, err := f.NewStyle(&excelize.Style{
		NumFmt: 41,
		Border: []excelize.Border{
			{Type: "top", Color: "000000", Style: 1},
			{Type: "bottom", Color: "000000", Style: 1},
			{Type: "left", Color: "000000", Style: 1},
			{Type: "right", Color: "000000", Style: 1},
		},
		Fill: excelize.Fill{
			Type:    "pattern",
			Pattern: 1,
			Color:   []string{"#FFFF00"},
		},
	})
	if err != nil {
		return nil, err
	}

	stylingDefault, err := f.NewStyle(&excelize.Style{
		NumFmt: 41,
		Fill: excelize.Fill{
			Type:    "pattern",
			Pattern: 1,
		},
		CustomNumFmt: &numberFormat,
	})
	if err != nil {
		return nil, err
	}

	f.SetColWidth(sheet, "K", "K", 2)

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
			return nil, response.ErrorBuilder(&response.ErrorConstant.InternalServerError, err)
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
			if v.IsTotal != nil && *v.IsTotal {
				if v.FxSummary == "" {
					row++
					continue
				}
				f.SetCellStyle(sheet, fmt.Sprintf("B%d", row), fmt.Sprintf("J%d", row), stylingDefault)
				f.SetCellStyle(sheet, fmt.Sprintf("L%d", row), fmt.Sprintf("S%d", row), stylingDefault)

				arrChr := []string{"E", "F"}
				if strings.ToUpper(v.Code) == "CONTROL_1" {
					arrChr = []string{"C"}
					f.SetCellStyle(sheet, fmt.Sprintf("C%d", row), fmt.Sprintf("C%d", row), stylingControl)

				}
				if strings.ToUpper(v.Code) == "CONTROL_17" {
					arrChr = []string{"C"}

					f.SetCellStyle(sheet, fmt.Sprintf("C%d", row), fmt.Sprintf("C%d", row), stylingControl)
				}
				if strings.ToUpper(v.Code) == "CONTROL_18" {
					arrChr = []string{"C"}

					f.SetCellStyle(sheet, fmt.Sprintf("C%d", row), fmt.Sprintf("C%d", row), stylingControl)
				}
				if strings.ToUpper(v.Code) == "CONTROL_2" {
					arrChr = []string{"D"}
				}
				if strings.ToUpper(v.Code) == "CONTROL_3" {
					arrChr = []string{"E"}
				}
				if strings.ToUpper(v.Code) == "CONTROL_4" {
					arrChr = []string{"F"}
				}
				if strings.ToUpper(v.Code) == "CONTROL_5" {
					arrChr = []string{"G"}
				}
				if strings.ToUpper(v.Code) == "CONTROL_6" {
					arrChr = []string{"H"}
				}
				if strings.ToUpper(v.Code) == "CONTROL_7" {
					arrChr = []string{"I"}
				}
				if strings.ToUpper(v.Code) == "CONTROL_8" {
					arrChr = []string{"J"}
				}
				if strings.ToUpper(v.Code) == "CONTROL_9" {
					arrChr = []string{"L"}
				}
				if strings.ToUpper(v.Code) == "CONTROL_10" {
					arrChr = []string{"M"}
				}
				if strings.ToUpper(v.Code) == "CONTROL_11" {
					arrChr = []string{"N"}
				}
				if strings.ToUpper(v.Code) == "CONTROL_12" {
					arrChr = []string{"O"}
				}
				if strings.ToUpper(v.Code) == "CONTROL_13" {
					arrChr = []string{"P"}
				}
				if strings.ToUpper(v.Code) == "CONTROL_14" {
					arrChr = []string{"Q"}
				}
				if strings.ToUpper(v.Code) == "CONTROL_15" {
					arrChr = []string{"R"}
				}
				if strings.ToUpper(v.Code) == "CONTROL_16" {
					arrChr = []string{"S"}
				}
				if strings.ToUpper(v.Code) == "CONTROL_17" {
					arrChr = []string{"C"}
				}
				if strings.ToUpper(v.Code) == "CONTROL_18" {
					arrChr = []string{"C"}
				}
				for _, chr := range arrChr {
					formula := v.FxSummary
					reg := regexp.MustCompile(`[A-Za-z0-9_~#:'()]+|[0-9]+\d{3,}`)
					match := reg.FindAllString(formula, -1)
					for _, vMatch := range match {
						// cari jml berdasarkan code
						if rowCode[vMatch] != 0 {
							if rowCode[vMatch] != 0 {
								formula = strings.ReplaceAll(formula, vMatch, fmt.Sprintf("%s%d", chr, rowCode[vMatch]))
							}
						}
						if _, ok := tbRowCodes[vMatch]; ok {
							formula = strings.ReplaceAll(formula, vMatch, fmt.Sprintf("L%d", tbRowCodes[vMatch]))
							f.SetCellValue("TRIAL_BALANCE", fmt.Sprintf("M%d", tbRowCodes[vMatch]), "control")
						}
					}

					if strings.ToUpper(v.Code) == "CONTROL_2" {
						row = row - 1

						f.SetCellStyle(sheet, fmt.Sprintf("D%d", row), fmt.Sprintf("D%d", row), stylingControl)
					}
					if strings.ToUpper(v.Code) == "CONTROL_3" {
						row = row - 1

						f.SetCellStyle(sheet, fmt.Sprintf("E%d", row), fmt.Sprintf("E%d", row), stylingControl)
					}
					if strings.ToUpper(v.Code) == "CONTROL_4" {
						row = row - 1

						f.SetCellStyle(sheet, fmt.Sprintf("F%d", row), fmt.Sprintf("F%d", row), stylingControl)
					}
					if strings.ToUpper(v.Code) == "CONTROL_5" {
						row = row - 1

						f.SetCellStyle(sheet, fmt.Sprintf("G%d", row), fmt.Sprintf("G%d", row), stylingControl)
					}
					if strings.ToUpper(v.Code) == "CONTROL_6" {
						row = row - 1

						f.SetCellStyle(sheet, fmt.Sprintf("H%d", row), fmt.Sprintf("H%d", row), stylingControl)
					}
					if strings.ToUpper(v.Code) == "CONTROL_7" {
						row = row - 1

						f.SetCellStyle(sheet, fmt.Sprintf("I%d", row), fmt.Sprintf("I%d", row), stylingControl)
					}
					if strings.ToUpper(v.Code) == "CONTROL_8" {
						row = row - 1

						f.SetCellStyle(sheet, fmt.Sprintf("J%d", row), fmt.Sprintf("J%d", row), stylingControl)
					}
					if strings.ToUpper(v.Code) == "CONTROL_9" {
						row = row - 1

						f.SetCellStyle(sheet, fmt.Sprintf("L%d", row), fmt.Sprintf("L%d", row), stylingControl)
					}
					if strings.ToUpper(v.Code) == "CONTROL_10" {
						row = row - 1

						f.SetCellStyle(sheet, fmt.Sprintf("M%d", row), fmt.Sprintf("M%d", row), stylingControl)
					}
					if strings.ToUpper(v.Code) == "CONTROL_11" {
						row = row - 1

						f.SetCellStyle(sheet, fmt.Sprintf("N%d", row), fmt.Sprintf("N%d", row), stylingControl)
					}
					if strings.ToUpper(v.Code) == "CONTROL_12" {
						row = row - 1

						f.SetCellStyle(sheet, fmt.Sprintf("O%d", row), fmt.Sprintf("O%d", row), stylingControl)
					}
					if strings.ToUpper(v.Code) == "CONTROL_13" {
						row = row - 1

						f.SetCellStyle(sheet, fmt.Sprintf("P%d", row), fmt.Sprintf("P%d", row), stylingControl)
					}
					if strings.ToUpper(v.Code) == "CONTROL_14" {
						row = row - 1

						f.SetCellStyle(sheet, fmt.Sprintf("Q%d", row), fmt.Sprintf("Q%d", row), stylingControl)
					}
					if strings.ToUpper(v.Code) == "CONTROL_15" {
						row = row - 1

						f.SetCellStyle(sheet, fmt.Sprintf("R%d", row), fmt.Sprintf("R%d", row), stylingControl)
					}
					if strings.ToUpper(v.Code) == "CONTROL_16" {
						row = row - 1

						f.SetCellStyle(sheet, fmt.Sprintf("S%d", row), fmt.Sprintf("S%d", row), stylingControl)
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

			AgingUPDetail, _, err := s.AgingUtangPiutangDetailRepository.Find(ctx, &criteriaAUP, &paginationAUP)
			if err != nil && v.Code != "" {
				continue
			}
			if len(*AgingUPDetail) == 0 {
				continue
			}
			valueKosong := 0.00
			f.SetCellValue(sheet, fmt.Sprintf("C%d", row), valueKosong)
			f.SetCellValue(sheet, fmt.Sprintf("D%d", row), valueKosong)
			f.SetCellValue(sheet, fmt.Sprintf("E%d", row), valueKosong)
			f.SetCellValue(sheet, fmt.Sprintf("F%d", row), valueKosong)
			f.SetCellValue(sheet, fmt.Sprintf("G%d", row), valueKosong)
			f.SetCellValue(sheet, fmt.Sprintf("H%d", row), valueKosong)
			f.SetCellValue(sheet, fmt.Sprintf("I%d", row), valueKosong)
			f.SetCellValue(sheet, fmt.Sprintf("J%d", row), valueKosong)
			f.SetCellValue(sheet, fmt.Sprintf("L%d", row), valueKosong)
			f.SetCellValue(sheet, fmt.Sprintf("M%d", row), valueKosong)
			f.SetCellValue(sheet, fmt.Sprintf("N%d", row), valueKosong)
			f.SetCellValue(sheet, fmt.Sprintf("O%d", row), valueKosong)
			f.SetCellValue(sheet, fmt.Sprintf("P%d", row), valueKosong)
			f.SetCellValue(sheet, fmt.Sprintf("Q%d", row), valueKosong)
			f.SetCellValue(sheet, fmt.Sprintf("R%d", row), valueKosong)
			f.SetCellValue(sheet, fmt.Sprintf("S%d", row), valueKosong)
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

	return f, nil
}

func (s *service) ExportPembelianPenjualanBerelasis(ctx *abstraction.Context, f *excelize.File, pembelianPenjualanBerelasi *model.PembelianPenjualanBerelasiEntityModel) (*excelize.File, error) {
	sheet := "PEMBELIAN_PENJUALAN_BERELASI"
	f.NewSheet(sheet)
	err := f.SetColWidth(sheet, "D", "D", 66)
	if err != nil {
		return nil, err
	}
	err = f.SetColWidth(sheet, "E", "E", 13)
	if err != nil {
		return nil, err
	}
	err = f.SetColWidth(sheet, "F", "F", 13)
	if err != nil {
		return nil, err
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
		NumFmt: 41,
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

	f.SetCellStyle(sheet, "B4", "F4", styleHeader)
	f.SetCellValue(sheet, "C2", "List pembelian dan penjualan berelasi")
	f.SetCellValue(sheet, "B4", "NO")
	f.SetCellValue(sheet, "C4", "CODE")
	f.SetCellValue(sheet, "D4", "COMPANY")
	f.SetCellValue(sheet, "E4", "PEMBELIAN")
	f.SetCellValue(sheet, "F4", "PENJUALAN")

	criteria := model.PembelianPenjualanBerelasiFilterModel{}
	criteria.CompanyID = &pembelianPenjualanBerelasi.CompanyID
	criteria.Period = &pembelianPenjualanBerelasi.Period
	criteria.Versions = &pembelianPenjualanBerelasi.Versions
	data, err := s.PembelianPenjualanBerelasiRepository.Export(ctx, &criteria)
	if err != nil {
		return nil, err
	}

	row, rowStart := 5, 5
	for i, v := range data.PembelianPenjualanBerelasiDetail {
		valueKosong := 0.0
		f.SetCellValue(sheet, fmt.Sprintf("B%d", row), (i + 1))
		f.SetCellValue(sheet, fmt.Sprintf("C%d", row), v.Code)
		f.SetCellValue(sheet, fmt.Sprintf("D%d", row), v.Name)
		f.SetCellValue(sheet, fmt.Sprintf("E%d", row), valueKosong)
		f.SetCellValue(sheet, fmt.Sprintf("F%d", row), valueKosong)
		if v.BoughtAmount != nil && *v.BoughtAmount != 0 {
			f.SetCellValue(sheet, fmt.Sprintf("E%d", row), *v.BoughtAmount)
		}
		if v.SalesAmount != nil && *v.SalesAmount != 0 {
			f.SetCellValue(sheet, fmt.Sprintf("F%d", row), *v.SalesAmount)
		}
		row++
	}

	f.SetCellStyle(sheet, "B5", fmt.Sprintf("D%d", row), styleLabel)
	f.SetCellStyle(sheet, "E5", fmt.Sprintf("F%d", row), styleCurrency)
	f.SetCellStyle(sheet, fmt.Sprintf("B%d", row+1), fmt.Sprintf("D%d", row+1), styleLabelTotal)
	f.SetCellStyle(sheet, fmt.Sprintf("E%d", row+1), fmt.Sprintf("F%d", row+1), styleCurrencyTotal)

	f.SetCellValue(sheet, fmt.Sprintf("D%d", row+1), "Total")
	f.SetCellFormula(sheet, fmt.Sprintf("E%d", row+1), fmt.Sprintf("=SUM(E%d:E%d)", rowStart, row))
	f.SetCellFormula(sheet, fmt.Sprintf("F%d", row+1), fmt.Sprintf("=SUM(F%d:F%d)", rowStart, row))

	return f, nil
}

var rowCodePersediaans = make(map[string]int)

func (s *service) ExportMutasiPersediaans(ctx *abstraction.Context, f *excelize.File, mutasiPersediaan *model.MutasiPersediaanEntityModel) (*excelize.File, error) {
	sheet := "MUTASI_PERSEDIAAN"
	f.NewSheet(sheet)
	f.SetColWidth(sheet, "B", "B", 31.01)
	f.SetColWidth(sheet, "C", "C", 13)
	f.SetColWidth(sheet, "D", "D", 13)

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
	stylingControl, err := f.NewStyle(&excelize.Style{
		NumFmt: 41,
		Border: []excelize.Border{
			{Type: "top", Color: "000000", Style: 1},
			{Type: "bottom", Color: "000000", Style: 1},
			{Type: "left", Color: "000000", Style: 1},
			{Type: "right", Color: "000000", Style: 1},
		},
		Fill: excelize.Fill{
			Type:    "pattern",
			Pattern: 1,
			Color:   []string{"#FFFF00"},
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
			Bold: true,
		},
		Fill: excelize.Fill{
			Type:    "pattern",
			Pattern: 1,
			Color:   []string{"#99ff66"},
		},
		NumFmt: 41,
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
		criteriaBridge.FormatterBridgesFilter.TrxRefID = &mutasiPersediaan.ID

		bridges, err := s.FormatterBridgesRepository.FindWithCriteria(ctx, &criteriaBridge)
		if err != nil {
			return nil, helper.ErrorHandler(err)
		}

		partRowStart := row
		for _, v := range data.FormatterDetail {
			rowCodePersediaans[v.Code] = row
			f.SetCellStyle(sheet, fmt.Sprintf("B%d", row), fmt.Sprintf("B%d", row), styleLabel)
			f.SetCellStyle(sheet, fmt.Sprintf("C%d", row), fmt.Sprintf("C%d", row), styleCurrency)

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
				arrChr := []string{"D", "E", "F", "G", "H", "I", "J", "K"}

				if strings.ToUpper(v.Code) == "CONTROL_1" {
					arrChr = []string{"C"}
					f.SetCellStyle(sheet, fmt.Sprintf("C%d", row), fmt.Sprintf("C%d", row), stylingControl)
				}
				for _, chr := range arrChr {
					formula := v.FxSummary
					reg := regexp.MustCompile(`[A-Za-z0-9_~#:()]+|[0-9]+\d{3,}`)
					match := reg.FindAllString(formula, -1)
					for _, vMatch := range match {
						// cari jml berdasarkan code
						if rowCodePersediaans[vMatch] != 0 {
							formula = strings.ReplaceAll(formula, vMatch, fmt.Sprintf("%s%d", chr, rowCodePersediaans[vMatch]))
						}
						if _, ok := tbRowCodes[vMatch]; ok {
							formula = strings.ReplaceAll(formula, vMatch, fmt.Sprintf("L%d", tbRowCodes[vMatch]))
							f.SetCellValue("TRIAL_BALANCE", fmt.Sprintf("M%d", tbRowCodes[vMatch]), "control")
						}

					}
					f.SetCellFormula(sheet, fmt.Sprintf("%s%d", chr, row), fmt.Sprintf("=%s", formula))
				}
				row++
				continue
			}

			if v.AutoSummary != nil && *v.AutoSummary {
				f.SetCellStyle(sheet, fmt.Sprintf("B%d", row), fmt.Sprintf("B%d", row), styleLabelTotal)
				f.SetCellStyle(sheet, fmt.Sprintf("C%d", row), fmt.Sprintf("C%d", row), styleCurrencyTotal)

				f.SetCellFormula(sheet, fmt.Sprintf("C%d", row), fmt.Sprintf("=SUM(C%d:C%d)", partRowStart, row-1))
				row++
				partRowStart = row
				continue
			}
			if v.ControlFormula != "" {

				// if v.FxSummary == "" {
				// 	row++
				// 	continue
				// }
				f.SetCellStyle(sheet, fmt.Sprintf("D%d", row), fmt.Sprintf("D%d", row), stylingControl)
				formula := v.FxSummary
				reg := regexp.MustCompile(`[A-Za-z0-9_~#:()]+|[0-9]+\d{3,}`)
				match := reg.FindAllString(formula, -1)
				for _, vMatch := range match {
					// cari jml berdasarkan code
					if _, ok := rowCodePersediaans[vMatch]; ok {
						formula = strings.ReplaceAll(formula, vMatch, fmt.Sprintf("C%d", rowCodePersediaans[vMatch]))
					}
					if _, ok := tbRowCodes[vMatch]; ok {
						formula = strings.ReplaceAll(formula, vMatch, fmt.Sprintf("L%d", tbRowCodes[vMatch]))
						f.SetCellValue("TRIAL_BALANCE", fmt.Sprintf("M%d", tbRowCodes[vMatch]), "control")
					}

				}
				f.SetCellFormula(sheet, fmt.Sprintf("D%d", row), fmt.Sprintf("=%s", formula))
			}

			criteriaMP := model.MutasiPersediaanDetailFilterModel{}
			criteriaMP.Code = &v.Code
			criteriaMP.FormatterBridgesID = &bridges.ID
			criteriaMP.MutasiPersediaanID = &mutasiPersediaan.ID

			paginationMP := abstraction.Pagination{}
			pagesize := 1
			paginationMP.PageSize = &pagesize

			mutasiPersediaanDetail, _, err := s.MutasiPersediaanDetailRepository.Find(ctx, &criteriaMP, &paginationMP)
			if err != nil {
				return nil, helper.ErrorHandler(err)
			}
			for _, vv := range *mutasiPersediaanDetail {
				valueKosong := 0.0
				f.SetCellValue(sheet, fmt.Sprintf("C%d", row), valueKosong)
				f.SetCellStyle(sheet, fmt.Sprintf("C%d", row), fmt.Sprintf("C%d", row), styleCurrency)
				if vv.Amount != nil && *vv.Amount != 0 {
					f.SetCellValue(sheet, fmt.Sprintf("C%d", row), *vv.Amount)
				}
			}
			row++
		}
		rowStart = row + 4
		row = rowStart
	}

	return f, nil
}

var rowCodeEms = make(map[string]int)

func (s *service) ExportEmployeeBenefits(ctx *abstraction.Context, f *excelize.File, employeeBenefit *model.EmployeeBenefitEntityModel) (*excelize.File, error) {
	sheet := "EMPLOYEE_BENEFIT"
	f.NewSheet(sheet)

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
		CustomNumFmt: &numberFormat,
	})
	if err != nil {
		return nil, err
	}
	stylingControl, err := f.NewStyle(&excelize.Style{
		NumFmt: 41,
		Border: []excelize.Border{
			{Type: "top", Color: "000000", Style: 1},
			{Type: "bottom", Color: "000000", Style: 1},
			{Type: "left", Color: "000000", Style: 1},
			{Type: "right", Color: "000000", Style: 1},
		},
		Fill: excelize.Fill{
			Type:    "pattern",
			Pattern: 1,
			Color:   []string{"#FFFF00"},
		},
	})
	if err != nil {
		return nil, err
	}
	styleCurrencyTotal, err := f.NewStyle(&excelize.Style{
		Font: &excelize.Font{
			Bold: true,
		},
		CustomNumFmt: &numberFormat,
	})
	if err != nil {
		return nil, err
	}

	styleLabelTotal, err := f.NewStyle(&excelize.Style{
		Font: &excelize.Font{
			Bold: true,
		},
	})
	if err != nil {
		return nil, err
	}

	datePeriod, err := time.Parse(time.RFC3339, employeeBenefit.Period)
	if err != nil {
		return nil, err
	}
	period := datePeriod.Format("02-Jan-06")

	formatterCode := []string{"EMPLOYEE-BENEFIT-ASUMSI", "EMPLOYEE-BENEFIT-REKONSILIASI", "EMPLOYEE-BENEFIT-RINCIAN-LAPORAN", "EMPLOYEE-BENEFIT-RINCIAN-EKUITAS", "EMPLOYEE-BENEFIT-MUTASI", "EMPLOYEE-BENEFIT-INFORMASI", "EMPLOYEE-BENEFIT-ANALISIS"}
	formatterTitle := []string{"Asumsi-asumsi yang digunakan:", "Rekonsiliasi jumlah liabilitas imbalan kerja karyawan pada laporan posisi keuangan adalah sebagai berikut:", "Rincian beban imbalan kerja karyawan yang diakui dalam laporan laba rugi dan penghasilan komprehensif lain adalah sebagai berikut:", "Rincian beban imbalan kerja karyawan yang diakui pada ekuitas dalam penghasilan komprehensif lain adalah sebagai berikut:", "Mutasi liabilitas imbalan kerja karyawan adalah sebagai berikut:", "Informasi historis dari nilai kini liabilitas imbalan pasti, nilai wajar aset program dan penyesuaian adalah sebagai berikut:", "Analisis sensitivitas dari perubahan asumsi-asumsi utama terhadap liabilitas imbalan kerja", ""}
	row, rowStart := 7, 7

	for i, formatter := range formatterCode {
		f.SetCellValue(sheet, fmt.Sprintf("B%d", (rowStart-3)), formatterTitle[i])
		f.MergeCell(sheet, fmt.Sprintf("B%d", (rowStart-2)), fmt.Sprintf("F%d", (rowStart-1)))
		f.SetCellStyle(sheet, fmt.Sprintf("B%d", (rowStart-2)), fmt.Sprintf("G%d", (rowStart-1)), styleHeader)
		// f.SetCellStyle(sheet, fmt.Sprintf("A%d", (rowStart-3)), fmt.Sprintf("A%d", (rowStart-3)), styleLabel)

		f.SetCellValue(sheet, fmt.Sprintf("B%d", (rowStart-2)), "Description")
		f.SetCellValue(sheet, fmt.Sprintf("G%d", (rowStart-2)), period)
		f.SetCellValue(sheet, fmt.Sprintf("G%d", (rowStart-1)), employeeBenefit.Company.Name)

		var criteria dto.FormatterGetRequest
		criteria.FormatterFilterModel.FormatterFor = &formatter

		data, err := s.FormatterRepository.FindWithDetail(ctx, &criteria.FormatterFilterModel)
		if err != nil {
			return nil, helper.ErrorHandler(err)
		}

		tmpStr := "EMPLOYEE-BENEFIT"
		criteriaBridge := model.FormatterBridgesFilterModel{}
		criteriaBridge.FormatterBridgesFilter.Source = &tmpStr
		criteriaBridge.FormatterBridgesFilter.FormatterID = &data.ID
		criteriaBridge.FormatterBridgesFilter.TrxRefID = &employeeBenefit.ID

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
			if rowCodeEms[v.Code] == 0 {
				rowCodeEms[v.Code] = row
			}
			rowKosong := 0.0
			f.SetCellValue(sheet, fmt.Sprintf("B%d", row), v.Description)
			// f.SetCellStyle(sheet, fmt.Sprintf("B%d", row), fmt.Sprintf("B%d", row), styleLabel)
			f.SetCellStyle(sheet, fmt.Sprintf("G%d", row), fmt.Sprintf("G%d", row), styleCurrency)
			f.SetCellValue(sheet, fmt.Sprintf("G%d", row), rowKosong)

			if v.ControlFormula != "" {
				if v.AutoSummary != nil && *v.AutoSummary {

					f.SetCellStyle(sheet, fmt.Sprintf("B%d", row), fmt.Sprintf("B%d", row), styleLabelTotal)
					f.SetCellStyle(sheet, fmt.Sprintf("G%d", row), fmt.Sprintf("G%d", row), styleCurrencyTotal)
					f.SetCellFormula(sheet, fmt.Sprintf("G%d", row), fmt.Sprintf("=SUM($G%d:$G%d)", partRowStart, row-1))
					partRowStart = row
				}
				if v.FxSummary == "" {
					row++
					continue
				}
				f.SetCellStyle(sheet, fmt.Sprintf("H%d", row), fmt.Sprintf("H%d", row), stylingControl)
				formula := v.FxSummary
				reg := regexp.MustCompile(`[A-Za-z0-9_~#:()]+|[0-9]+\d{3,}`)
				match := reg.FindAllString(formula, -1)
				for _, vMatch := range match {
					// cari jml berdasarkan code
					if _, ok := rowCodeEms[vMatch]; ok {
						formula = strings.ReplaceAll(formula, vMatch, fmt.Sprintf("G%d", rowCodeEms[vMatch]))
					}
					if _, ok := tbRowCodes[vMatch]; ok {
						formula = strings.ReplaceAll(formula, vMatch, fmt.Sprintf("L%d", tbRowCodes[vMatch]))
						f.SetCellValue("TRIAL_BALANCE", fmt.Sprintf("M%d", tbRowCodes[vMatch]), "control")
					}

				}
				f.SetCellFormula(sheet, fmt.Sprintf("H%d", row), fmt.Sprintf("=%s", formula))
				row++
				continue
			}
			if v.IsTotal != nil && *v.IsTotal {

				// if v.FxSummary == "" {
				// 	row++
				// 	continue
				// }

				arrChr := []string{"G"}

				if strings.ToUpper(v.Code) == "CONTROL_1" {
					arrChr = []string{"G"}
					f.SetCellStyle(sheet, fmt.Sprintf("G%d", row), fmt.Sprintf("G%d", row), styleCurrency)
					f.SetCellStyle(sheet, fmt.Sprintf("G%d", row), fmt.Sprintf("G%d", row), stylingControl)
				}
				if strings.ToUpper(v.Code) == "CONTROL_2" {
					arrChr = []string{"G"}
					f.SetCellStyle(sheet, fmt.Sprintf("G%d", row), fmt.Sprintf("G%d", row), styleCurrency)
					f.SetCellStyle(sheet, fmt.Sprintf("G%d", row), fmt.Sprintf("G%d", row), stylingControl)
				}
				for _, chr := range arrChr {
					formula := v.FxSummary
					reg := regexp.MustCompile(`[A-Za-z0-9_~#:()]+|[0-9]+\d{3,}`)
					match := reg.FindAllString(formula, -1)
					for _, vMatch := range match {
						// cari jml berdasarkan code
						if rowCodeEms[vMatch] != 0 {
							formula = strings.ReplaceAll(formula, vMatch, fmt.Sprintf("%s%d", chr, rowCodeEms[vMatch]))
						}
						if _, ok := tbRowCodes[vMatch]; ok {
							formula = strings.ReplaceAll(formula, vMatch, fmt.Sprintf("L%d", tbRowCodes[vMatch]))
							f.SetCellValue("TRIAL_BALANCE", fmt.Sprintf("M%d", tbRowCodes[vMatch]), "control")
						}

					}
					f.SetCellFormula(sheet, fmt.Sprintf("%s%d", chr, row), fmt.Sprintf("=%s", formula))
				}
				row++
				continue
			}
			if v.AutoSummary != nil && *v.AutoSummary {

				f.SetCellStyle(sheet, fmt.Sprintf("B%d", row), fmt.Sprintf("B%d", row), styleLabelTotal)
				f.SetCellStyle(sheet, fmt.Sprintf("G%d", row), fmt.Sprintf("G%d", row), styleCurrencyTotal)
				f.SetCellFormula(sheet, fmt.Sprintf("G%d", row), fmt.Sprintf("=SUM($G%d:$G%d)", partRowStart, row-1))
				row++
				partRowStart = row
				continue
			}

			criteriaEB := model.EmployeeBenefitDetailFilterModel{}
			criteriaEB.Code = &v.Code
			criteriaEB.FormatterBridgesID = &bridges.ID
			criteriaEB.EmployeeBenefitID = &employeeBenefit.ID

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
				valueKosong := 0.0
				f.SetCellValue(sheet, fmt.Sprintf("G%d", row), valueKosong)
				f.SetCellStyle(sheet, fmt.Sprintf("G%d", row), fmt.Sprintf("G%d", row), styleCurrency)
				if v.IsLabel == nil || (v.IsLabel != nil && !*v.IsLabel) {
					if vv.IsValue != nil && *vv.IsValue {
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

	return f, nil
}

func (s *service) ExportInvestasiNonTbks(ctx *abstraction.Context, f *excelize.File, investasiNonTbk *model.InvestasiNonTbkEntityModel) (*excelize.File, error) {
	sheet := "INVESTASI_NON_TBK"
	f.NewSheet(sheet)

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
			return nil, response.ErrorBuilder(&response.ErrorConstant.InternalServerError, errors.New("internal server error"))
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
		return nil, response.ErrorBuilder(&response.ErrorConstant.InternalServerError, errors.New("internal server error"))
	}
	numberFormat := "#,##"
	styleCurrencyPercentage, err := f.NewStyle(&excelize.Style{
		Border: []excelize.Border{
			{Type: "top", Color: "000000", Style: 1},
			{Type: "bottom", Color: "000000", Style: 1},
			{Type: "left", Color: "000000", Style: 1},
			{Type: "right", Color: "000000", Style: 1},
		},
		NumFmt: 41,
		Font: &excelize.Font{
			Family: "Arial",
			Size:   10,
		},
	})
	if err != nil {
		return nil, response.ErrorBuilder(&response.ErrorConstant.InternalServerError, errors.New("internal server error"))
	}
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
		return nil, response.ErrorBuilder(&response.ErrorConstant.InternalServerError, errors.New("internal server error"))
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
		return nil, response.ErrorBuilder(&response.ErrorConstant.InternalServerError, errors.New("internal server error"))
	}

	allowed := helper.CompanyValidation(ctx.Auth.ID, investasiNonTbk.CompanyID)
	if !allowed {
		return nil, response.ErrorBuilder(&response.ErrorConstant.UnauthorizedAccess, errors.New("not allowed"))
	}

	// datePeriod, err := time.Parse(time.RFC3339, investasiNonTbk.Period)
	// if err != nil {
	// 	return nil, err
	// }

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

	criteriaDetail.InvestasiNonTbkID = &investasiNonTbk.ID

	paginationDetail := abstraction.Pagination{}
	pagesize := 10000
	paginationDetail.PageSize = &pagesize

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
		tmp1 := 0.0
		if v.LbrSahamOwnership != nil && *v.LbrSahamOwnership != 0 {
			tmp1 = *v.LbrSahamOwnership
		}
		f.SetCellValue(sheet, fmt.Sprintf("D%d", row), tmp1)
		tmp2 := 0.0
		if v.TotalLbrSaham != nil && *v.TotalLbrSaham != 0 {
			tmp2 = *v.TotalLbrSaham
		}
		f.SetCellValue(sheet, fmt.Sprintf("E%d", row), tmp2)
		tmp3 := 0.0
		if v.HargaPar != nil && *v.HargaPar != 0 {
			tmp3 = *v.HargaPar
		}
		f.SetCellValue(sheet, fmt.Sprintf("G%d", row), tmp3)
		tmp4 := 0.0
		if v.HargaBeli != nil && *v.HargaBeli != 0 {
			tmp4 = *v.HargaBeli
		}
		f.SetCellValue(sheet, fmt.Sprintf("I%d", row), tmp4)

		f.SetCellFormula(sheet, fmt.Sprintf("F%d", row), fmt.Sprintf("=IFERROR(D%d/E%d,%d)", row, row, 0))
		f.SetCellFormula(sheet, fmt.Sprintf("H%d", row), fmt.Sprintf("=D%d*G%d", row, row))
		f.SetCellFormula(sheet, fmt.Sprintf("J%d", row), fmt.Sprintf("=D%d*I%d", row, row))
		row++
	}

	return f, nil
}

func (s *service) ExportInvestasiTbks(ctx *abstraction.Context, f *excelize.File, investasiTbk *model.InvestasiTbkEntityModel) (*excelize.File, error) {
	sheet := "INVESTASI_TBK"
	f.NewSheet(sheet)

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
	styleCurrencySum, err := f.NewStyle(&excelize.Style{
		Border: []excelize.Border{
			{Type: "top", Color: "000000", Style: 1},
			{Type: "bottom", Color: "000000", Style: 1},
			{Type: "left", Color: "000000", Style: 1},
			{Type: "right", Color: "000000", Style: 1},
		},
		NumFmt: 41,
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
		NumFmt: 41,
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

	allowed := helper.CompanyValidation(ctx.Auth.ID, investasiTbk.CompanyID)
	if !allowed {
		return nil, response.ErrorBuilder(&response.ErrorConstant.UnauthorizedAccess, errors.New("not allowed"))
	}

	datePeriod, err := time.Parse(time.RFC3339, investasiTbk.Period)
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
			return nil, response.ErrorBuilder(&response.ErrorConstant.InternalServerError, err)
		}

		tmpStr := "INVESTASI-TBK"
		criteriaBridge := model.FormatterBridgesFilterModel{}
		criteriaBridge.FormatterBridgesFilter.Source = &tmpStr
		criteriaBridge.FormatterBridgesFilter.FormatterID = &data.ID
		criteriaBridge.FormatterBridgesFilter.TrxRefID = &investasiTbk.ID

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

			f.SetCellValue(sheet, fmt.Sprintf("C%d", row), v.Description)

			if v.IsTotal != nil && *v.IsTotal {
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
						// cari jml berdasarkan code
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
			criteriaIT.InvestasiTbkID = &investasiTbk.ID

			paginationIT := abstraction.Pagination{}
			pagesize := 10000
			paginationIT.PageSize = &pagesize

			investasiTbkDetail, _, err := s.InvestasiTbkDetailRepository.Find(ctx, &criteriaIT, &paginationIT)
			if err != nil && v.Code != "" {
				continue
			}
			if len(*investasiTbkDetail) == 0 {
				continue
			}
			for _, vv := range *investasiTbkDetail {
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

				f.SetCellFormula(sheet, fmt.Sprintf("F%d", row), fmt.Sprintf("=D%d*E%d", row, row))
				f.SetCellFormula(sheet, fmt.Sprintf("H%d", row), fmt.Sprintf("=G%d*D%d", row, row))
				f.SetCellFormula(sheet, fmt.Sprintf("I%d", row), fmt.Sprintf("=H%d-F%d", row, row))
				f.SetCellStyle(sheet, fmt.Sprintf("F%d", row), fmt.Sprintf("F%d", row), styleCurrencySum)
				f.SetCellStyle(sheet, fmt.Sprintf("H%d", row), fmt.Sprintf("H%d", row), styleCurrencySum)
				f.SetCellStyle(sheet, fmt.Sprintf("I%d", row), fmt.Sprintf("I%d", row), styleCurrencySum)
			}
			row++
		}
		rowStart = row
		row = rowStart
	}

	return f, nil
}
