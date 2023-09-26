package export

import (
	"archive/zip"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"path"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"time"
	"worker/configs"
	"worker/internal/abstraction"
	"worker/internal/factory"
	kafkaproducer "worker/internal/kafka/producer"
	"worker/internal/model"
	"worker/internal/repository"
	utilDate "worker/pkg/util/date"

	"github.com/xuri/excelize/v2"
	"gorm.io/gorm"
)

type service struct {
	AgingUtangPiutangRepository          repository.AgingUtangPiutang
	Db                                   *gorm.DB
	FormatterRepository                  repository.Formatter
	CompanyRepository                    repository.Company
	AgingUPDetailRepository              repository.AgingUtangPiutangDetail
	ParameterRepository                  repository.Parameter
	InvestasiNonTbkRepository            repository.InvestasiNonTbk
	InvestasiNonTbkDetailRepository      repository.InvestasiNonTbkDetail
	InvestasiTbkRepository               repository.InvestasiTbk
	InvestasiTbkDetailRepository         repository.InvestasiTbkDetail
	MutasiDtaRepository                  repository.MutasiDta
	MutasiDtaDetailRepository            repository.MutasiDtaDetail
	MutasiFaRepository                   repository.MutasiFa
	MutasiFaDetailRepository             repository.MutasiFaDetail
	MutasiIaRepository                   repository.MutasiIa
	MutasiIaDetailRepository             repository.MutasiIaDetail
	MutasiPersediaanRepository           repository.MutasiPersediaan
	MutasiPersediaanDetailRepository     repository.MutasiPersediaanDetail
	MutasiRuaRepository                  repository.MutasiRua
	MutasiRuaDetailRepository            repository.MutasiRuaDetail
	PembelianPenjualanBerelasiRepository repository.PembelianPenjualanBerelasi
	PPBerelasiDetailRepository           repository.PembelianPenjualanBerelasiDetail
	TrialBalanceRepository               repository.TrialBalance
	TrialBalanceDetailRepository         repository.TrialBalanceDetail
	CoaRepository                        repository.Coa
	FormatterBridgesRepository           repository.FormatterBridges
	FormatterDetailRepository            repository.FormatterDetail
	JpmRepository                        repository.Jpm
	JpmDetailRepository                  repository.JpmDetail
	AjeRepository                        repository.Adjustment
	AjeDetailRepository                  repository.AdjustmentDetail
	JelimRepository                      repository.Jelim
	JelimDetailRepository                repository.JelimDetail
	JcteRepository                       repository.Jcte
	JcteDetailRepository                 repository.JcteDetail
	NotificationRepository               repository.Notification
	EmployeeBenefitRepository            repository.EmployeeBenefit
	EmployeeBenefitDetailRepository      repository.EmployeeBenefitDetail
	ConsolidationRepository              repository.Consolidation
	ConsolidationDetailRepository        repository.ConsolidationDetail
	ConsolidationBridgeRepository        repository.ConsolidationBridge
	ConsolidationBridgeDetailRepository  repository.ConsolidationBridgeDetail
}

type Service interface {
	ExportAll(ctx *abstraction.Context, payload *abstraction.JsonData)
	ExportAgingUtangPiutang(ctx *abstraction.Context, payload *abstraction.JsonData, wg *sync.WaitGroup)
	ExportInvestasiNonTbk(ctx *abstraction.Context, payload *abstraction.JsonData, wg *sync.WaitGroup)
	ExportInvestasiTbk(ctx *abstraction.Context, payload *abstraction.JsonData, wg *sync.WaitGroup)
	ExportMutasiDta(ctx *abstraction.Context, payload *abstraction.JsonData, wg *sync.WaitGroup)
	ExportMutasiFa(ctx *abstraction.Context, payload *abstraction.JsonData, wg *sync.WaitGroup)
	ExportMutasiIa(ctx *abstraction.Context, payload *abstraction.JsonData, wg *sync.WaitGroup)
	ExportMutasiPersediaan(ctx *abstraction.Context, payload *abstraction.JsonData, wg *sync.WaitGroup)
	ExportMutasiRua(ctx *abstraction.Context, payload *abstraction.JsonData, wg *sync.WaitGroup)
	ExportPembelianPenjualanBerelasi(ctx *abstraction.Context, payload *abstraction.JsonData, wg *sync.WaitGroup)
	ExportTrialBalance(ctx *abstraction.Context, payload *abstraction.JsonData, wg *sync.WaitGroup)
	ExportEmployeeBenefit(ctx *abstraction.Context, payload *abstraction.JsonData, wg *sync.WaitGroup)
	ExportConsolidation(ctx *abstraction.Context, payload *abstraction.JsonData)
}

func NewService(f *factory.Factory) *service {
	agingUtangPiutangRepository := f.AgingUtangPiutangRepository
	formatterRepository := f.FormatterRepository
	companyRepository := f.CompanyRepository
	agingUPDetailRepository := f.AgingUtangPiutangDetailRepository
	parameterRepository := f.ParameterRepository
	investasiNonTbkRepository := f.InvestasiNonTbkRepository
	investasiNonTbkDetailRepository := f.InvestasiNonTbkDetailRepository
	investasiTbkRepository := f.InvestasiTbkRepository
	investasiTbkDetailRepository := f.InvestasiTbkDetailRepository
	mutasiDtaRepository := f.MutasiDtaRepository
	mutasiDtaDetailRepository := f.MutasiDtaDetailRepository
	mutasiFaRepository := f.MutasiFaRepository
	mutasiFaDetailRepository := f.MutasiFaDetailRepository
	mutasiIaRepository := f.MutasiIaRepository
	mutasiIaDetailRepository := f.MutasiIaDetailRepository
	mutasiPersediaanRepository := f.MutasiPersediaanRepository
	mutasiPersediaanDetailRepository := f.MutasiPersediaanDetailRepository
	mutasiRuaRepository := f.MutasiRuaRepository
	mutasiRuaDetailRepository := f.MutasiRuaDetailRepository
	pembelianPenjualanBerelasiRepository := f.PembelianPenjualanBerelasiRepository
	pembelianPenjualanBerelasiDetailRepository := f.PembelianPenjualanBerelasiDetailRepository
	trialBalanceRepository := f.TrialBalanceRepository
	trialBalanceDetailRepository := f.TrialBalanceDetailRepository
	coaRepository := f.CoaRepository
	formatterBridgesRepository := f.FormatterBridgesRepository
	formatterDetailRepository := f.FormatterDetailRepository
	jpmRepo := f.JpmRepository
	jpmDetailRepo := f.JpmDetailRepository
	ajeRepo := f.AjeRepository
	ajeDetailRepo := f.AjeDetailRepository
	jcteRepo := f.JcteRepository
	jcteDetailRepo := f.JcteDetailRepository
	jelimRepo := f.JelimRepository
	jelimDetailRepo := f.JelimDetailRepository
	notificationRepo := f.NotificationRepository
	consolRepo := f.ConsolidationRepository
	consolDetailRepo := f.ConsolidationDetailRepository
	consolBridgeRepo := f.ConsolidationBridgeRepository
	consolBridgeDetailRepo := f.ConsolidationBridgeDetailRepository
	db := f.Db
	employeeBenefitRepo := f.EmployeeBenefitRepository
	employeeBenefitDetailRepo := f.EmployeeBenefitDetailRepository

	return &service{
		AgingUtangPiutangRepository:          agingUtangPiutangRepository,
		ParameterRepository:                  parameterRepository,
		Db:                                   db,
		FormatterRepository:                  formatterRepository,
		CompanyRepository:                    companyRepository,
		AgingUPDetailRepository:              agingUPDetailRepository,
		InvestasiNonTbkRepository:            investasiNonTbkRepository,
		InvestasiNonTbkDetailRepository:      investasiNonTbkDetailRepository,
		InvestasiTbkRepository:               investasiTbkRepository,
		InvestasiTbkDetailRepository:         investasiTbkDetailRepository,
		MutasiDtaRepository:                  mutasiDtaRepository,
		MutasiDtaDetailRepository:            mutasiDtaDetailRepository,
		MutasiFaRepository:                   mutasiFaRepository,
		MutasiFaDetailRepository:             mutasiFaDetailRepository,
		MutasiIaRepository:                   mutasiIaRepository,
		MutasiIaDetailRepository:             mutasiIaDetailRepository,
		MutasiPersediaanRepository:           mutasiPersediaanRepository,
		MutasiPersediaanDetailRepository:     mutasiPersediaanDetailRepository,
		MutasiRuaRepository:                  mutasiRuaRepository,
		MutasiRuaDetailRepository:            mutasiRuaDetailRepository,
		PembelianPenjualanBerelasiRepository: pembelianPenjualanBerelasiRepository,
		PPBerelasiDetailRepository:           pembelianPenjualanBerelasiDetailRepository,
		TrialBalanceRepository:               trialBalanceRepository,
		TrialBalanceDetailRepository:         trialBalanceDetailRepository,
		CoaRepository:                        coaRepository,
		FormatterBridgesRepository:           formatterBridgesRepository,
		FormatterDetailRepository:            formatterDetailRepository,
		JpmRepository:                        jpmRepo,
		JpmDetailRepository:                  jpmDetailRepo,
		JelimRepository:                      jelimRepo,
		JelimDetailRepository:                jelimDetailRepo,
		JcteRepository:                       jcteRepo,
		JcteDetailRepository:                 jcteDetailRepo,
		AjeRepository:                        ajeRepo,
		AjeDetailRepository:                  ajeDetailRepo,
		NotificationRepository:               notificationRepo,
		ConsolidationRepository:              consolRepo,
		ConsolidationDetailRepository:        consolDetailRepo,
		ConsolidationBridgeRepository:        consolBridgeRepo,
		ConsolidationBridgeDetailRepository:  consolBridgeDetailRepo,
		EmployeeBenefitRepository:            employeeBenefitRepo,
		EmployeeBenefitDetailRepository:      employeeBenefitDetailRepo,
	}
}

var errs = make(map[string]string)

type JsonData struct {
	CompanyName string
	Type        string
	Period      string
	Versions    int
	DataID      int
	Errors      string
	File        string
}

func (s *service) ExportAll(ctx *abstraction.Context, payload *abstraction.JsonData) {
	date := time.Now().Format("20060102_150405")
	saveAsset := fmt.Sprintf("assets/%d/REQUEST_DATA_%s", payload.UserID, date)
	tmpFolder := path.Join(configs.App().StoragePath(), saveAsset)
	_, err := os.Stat(tmpFolder)
	if err != nil {
		if os.IsNotExist(err) {
			if err := os.MkdirAll(tmpFolder, os.ModePerm); err != nil {
				errs["General"] = "Failed to create Zip"
			}
		} else {
			errs["General"] = "Failed to create Zip error on server"
		}
	}
	var request []string
	if payload.Filter.Request != "" {
		request = strings.Split(strings.ToUpper(payload.Filter.Request), ",")
	}

	datePeriod, err := time.Parse("2006-01-02", payload.Filter.Period)
	if err != nil {
		errs["General"] = "Invalid Period Format"
		log.Println(err)
		return
	}
	period := datePeriod.Format("2006-01-02")

	var arrFileName []string
	wg := new(sync.WaitGroup)
	if len(request) > 0 {
		wg.Add(len(request))
		for _, v := range request {
			switch v {
			case "TB":
				arrFileName = append(arrFileName, fmt.Sprintf("TrialBalance_%s.xlsx", period))
				go s.ExportTrialBalance(ctx, payload, wg, tmpFolder)
			case "AUP":
				arrFileName = append(arrFileName, fmt.Sprintf("AgingUtangPiutang_%s.xlsx", period))
				go s.ExportAgingUtangPiutang(ctx, payload, wg, tmpFolder)
			case "MP":
				arrFileName = append(arrFileName, fmt.Sprintf("MutasiPersediaan_%s.xlsx", period))
				go s.ExportMutasiPersediaan(ctx, payload, wg, tmpFolder)
			case "MRUA":
				arrFileName = append(arrFileName, fmt.Sprintf("MutasiRua_%s.xlsx", period))
				go s.ExportMutasiRua(ctx, payload, wg, tmpFolder)
			case "EB":
				arrFileName = append(arrFileName, fmt.Sprintf("EmployeeBenefit_%s.xlsx", period))
				go s.ExportEmployeeBenefit(ctx, payload, wg, tmpFolder)
			// case "AJE":
			// 	arrFileName = append(arrFileName, fmt.Sprintf("Adjustment_%s.xlsx", period))
			// 	go s.ExportAje(ctx, payload, wg, tmpFolder)
			// case "JCTE":
			// 	arrFileName = append(arrFileName, fmt.Sprintf("JurnalCTE_%s.xlsx", period))
			// 	go s.ExportJcte(ctx, payload, wg, tmpFolder)
			// case "JELIM":
			// 	arrFileName = append(arrFileName, fmt.Sprintf("JurnalEliminasi_%s.xlsx", period))
			// 	go s.ExportJelim(ctx, payload, wg, tmpFolder)
			// case "JPM":
			// 	arrFileName = append(arrFileName, fmt.Sprintf("JurnalProformaModal_%s.xlsx", period))
			// 	go s.ExportJpm(ctx, payload, wg, tmpFolder)
			case "PPB":
				arrFileName = append(arrFileName, fmt.Sprintf("PembelianPenjualanBerelasi_%s.xlsx", period))
				go s.ExportPembelianPenjualanBerelasi(ctx, payload, wg, tmpFolder)
			case "MFA":
				arrFileName = append(arrFileName, fmt.Sprintf("MutasiFa_%s.xlsx", period))
				go s.ExportMutasiFa(ctx, payload, wg, tmpFolder)
			case "MDTA":
				arrFileName = append(arrFileName, fmt.Sprintf("MutasiDta_%s.xlsx", period))
				go s.ExportMutasiDta(ctx, payload, wg, tmpFolder)
			case "MIA":
				arrFileName = append(arrFileName, fmt.Sprintf("MutasiIa_%s.xlsx", period))
				go s.ExportMutasiIa(ctx, payload, wg, tmpFolder)
			case "IT":
				arrFileName = append(arrFileName, fmt.Sprintf("InvestasiTbk_%s.xlsx", period))
				go s.ExportInvestasiTbk(ctx, payload, wg, tmpFolder)
			case "INT":
				arrFileName = append(arrFileName, fmt.Sprintf("InvestasiNonTbk_%s.xlsx", period))
				go s.ExportInvestasiNonTbk(ctx, payload, wg, tmpFolder)
			}
		}
	} else {
		// arrFileName = []string{"TrialBalance.xlsx", "AgingUtangPiutang.xlsx", "MutasiPersediaan.xlsx", "MutasiRua.xlsx", "PembelianPenjualanBerelasi.xlsx", "MutasiFa.xlsx", "MutasiDta.xlsx", "MutasiIa.xlsx", "InvestasiTbk.xlsx", "InvestasiNonTbk.xlsx", "Adjustment.xlsx", "JurnalCTE.xlsx", "JurnalEliminasi.xlsx", "JurnalProformaModal.xlsx"}
		arrFileName = []string{fmt.Sprintf("TrialBalance_%s.xlsx", period), fmt.Sprintf("AgingUtangPiutang_%s.xlsx", period), fmt.Sprintf("MutasiPersediaan_%s.xlsx", period), fmt.Sprintf("MutasiRua_%s.xlsx", period), fmt.Sprintf("PembelianPenjualanBerelasi_%s.xlsx", period), fmt.Sprintf("MutasiFa_%s.xlsx", period), fmt.Sprintf("MutasiDta_%s.xlsx", period), fmt.Sprintf("MutasiIa_%s.xlsx", period), fmt.Sprintf("InvestasiTbk_%s.xlsx", period), fmt.Sprintf("InvestasiNonTbk_%s.xlsx", period), fmt.Sprintf("EmployeeBenefit_%s.xlsx", period)}
		wg.Add(len(arrFileName))
		go s.ExportAgingUtangPiutang(ctx, payload, wg, tmpFolder)
		go s.ExportInvestasiNonTbk(ctx, payload, wg, tmpFolder)
		go s.ExportInvestasiTbk(ctx, payload, wg, tmpFolder)
		go s.ExportMutasiDta(ctx, payload, wg, tmpFolder)
		go s.ExportMutasiFa(ctx, payload, wg, tmpFolder)
		go s.ExportMutasiIa(ctx, payload, wg, tmpFolder)
		go s.ExportMutasiPersediaan(ctx, payload, wg, tmpFolder)
		go s.ExportMutasiRua(ctx, payload, wg, tmpFolder)
		go s.ExportPembelianPenjualanBerelasi(ctx, payload, wg, tmpFolder)
		go s.ExportTrialBalance(ctx, payload, wg, tmpFolder)
		go s.ExportEmployeeBenefit(ctx, payload, wg, tmpFolder)
		// go s.ExportAje(ctx, payload, wg, tmpFolder)
		// go s.ExportJcte(ctx, payload, wg, tmpFolder)
		// go s.ExportJelim(ctx, payload, wg, tmpFolder)
		// go s.ExportJelim(ctx, payload, wg, tmpFolder)
	}

	tmpLoc := fmt.Sprintf("assets/%d/RequestExports%s.zip", payload.UserID, date)
	loc := path.Join(configs.App().StoragePath(), tmpLoc)
	archive, err := os.Create(loc)
	if err != nil {
		errs["General"] = "Failed to create ZIP"
		log.Println(err)
	}
	wg.Wait()

	zipWriter := zip.NewWriter(archive)
	for _, file := range arrFileName {
		// source := fmt.Sprintf("%s/%s", tmpFolder, file)
		f1, err := os.Open(fmt.Sprintf("%s/%s", tmpFolder, file))
		if err != nil {
			errs["General"] = "Failed to create ZIP"
			log.Println(err)
		}

		w1, err := zipWriter.Create(file)
		if err != nil {
			errs["General"] = "Failed to create ZIP"
			log.Println(err)
		}
		if _, err := io.Copy(w1, f1); err != nil {
			errs["General"] = "Failed to create ZIP"
			log.Println(err)
		}
		f1.Close()
	}
	defer func() {
		if err := zipWriter.Close(); err != nil {
			errs["General"] = "Failed to create ZIP"
			fmt.Println(err)
		}
		if err := archive.Close(); err != nil {
			errs["General"] = "Failed to create ZIP"
			fmt.Println(err)
		}
		if err := os.RemoveAll(tmpFolder); err != nil {
			errs["General"] = "Failed to remove temporary data"
			log.Println(err)
		}
	}()

	notifData := model.NotificationEntityModel{}
	notifData.Context = ctx
	waktu := time.Now()
	map1 := kafkaproducer.JsonData{
		FileLoc:   tmpLoc,
		UserID:    ctx.Auth.ID,
		CompanyID: ctx.Auth.CompanyID,
		Name:      "export",
		Timestamp: &waktu,
	}

	company, err := s.CompanyRepository.FindByID(ctx, &ctx.Auth.CompanyID)
	if err != nil {
		log.Println(err)
		return
	}

	notifData.Description = "Proses Export Berhasil"
	notifData.Data = "{}"
	datas := JsonData{
		CompanyName: company.Name,
		Period:      payload.Filter.Period,
		Versions:    payload.Filter.Versions,
		File:        tmpLoc,
		Type:        "export",
	}
	checkErr := 0
	for _, err := range errs {
		if err != "" {
			checkErr++
		}
	}
	if checkErr > 0 {
		notifData.Description = "Proses Export Gagal!"
		tmpMap := []string{}
		for _, err := range errs {
			if err != "" {
				tmpMap = append(tmpMap, err)
			}
		}
		jsonErr, err := json.Marshal(tmpMap)
		if err != nil {
			log.Println(err)
			return
		}
		datas.Errors = string(jsonErr)
		// go kafkaproducer.NewProducer("NOTIFICATION").SendMessage("NOTIFICATION", string(jsonStr))
		// return
	}
	jsonData, err := json.Marshal(datas)
	if err != nil {
		log.Println(err)
		return
	}
	notifData.Data = string(jsonData)
	map1.Data = notifData.Description
	tmpfalse := false
	notifData.IsOpened = &tmpfalse
	notifData.CreatedBy = ctx.Auth.ID
	notifData.CreatedAt = *utilDate.DateTodayLocal()
	notificationData, err := s.NotificationRepository.Create(ctx, &notifData)
	if err != nil {
		log.Println(err)
		return
	}
	map1.ID = notificationData.ID
	jsonStr, err := json.Marshal(map1)
	if err != nil {
		errs["General"] = "Failed to create Notification"
		log.Println(err)
		return
	}

	go kafkaproducer.NewProducer("NOTIFICATION").SendMessage("NOTIFICATION", string(jsonStr))
}

func (s *service) ExportConsolidation(ctx *abstraction.Context, payload *abstraction.JsonData) {
	defer fmt.Println("Done Consolidation")
	criterieConsol := model.ConsolidationFilterModel{}
	criterieConsol.CompanyID = &payload.CompanyID
	criterieConsol.ConsolidationVersions = &payload.Filter.Versions
	criterieConsol.Period = &payload.Filter.Period
	consolidationData, err := s.ConsolidationRepository.FindByCriteria(ctx, &criterieConsol)
	if err != nil {
		errs["CONSOL"] = "Error: Record not found"
		fmt.Println(err)
		return
	}

	f := excelize.NewFile()
	tmp := f.GetSheetName(f.GetActiveSheetIndex())
	f.SetSheetName(tmp, "Consolidation")
	sheet := f.GetSheetName(f.GetActiveSheetIndex())

	errJournal := []error{}
	wg := sync.WaitGroup{}
	wg.Add(3)
	// Export Jurnal
	go func() {
		defer wg.Done()
		journalJpmFile, err := s.ExportJpmWorksheet(ctx, consolidationData, f)
		if err != nil {
			errJournal = append(errJournal, err)
			return
		}
		f = journalJpmFile
	}()

	go func() {
		defer wg.Done()
		journalJcteFile, err := s.ExportJcteWorksheet(ctx, consolidationData, f)
		if err != nil {
			errJournal = append(errJournal, err)
			return
		}
		f = journalJcteFile
	}()

	go func() {
		defer wg.Done()
		journalJelimFile, err := s.ExportJelimWorksheet(ctx, consolidationData, f)
		if err != nil {
			errJournal = append(errJournal, err)
			return
		}
		f = journalJelimFile
	}()

	// go func() {
	// 	defer wg.Done()
	// 	mutasiFaNew, err := s.ExportMutasiFaNew(ctx, consolidationData, f)
	// 	if err != nil {
	// 		errJournal = append(errJournal, err)
	// 		return
	// 	}
	// 	f = mutasiFaNew
	// }()

	formatterID := 3
	var criteriaFormatter model.FormatterDetailFilterModel
	criteriaFormatter.FormatterID = &formatterID
	t := true
	criteriaFormatter.IsShowExport = &t
	data, err := s.FormatterDetailRepository.Find(ctx, &criteriaFormatter)
	if err != nil {
		errs["CONSOL"] = "Error: Get Formatter Detail"
		return
	}

	datePeriod, err := time.Parse(time.RFC3339, consolidationData.Period)
	if err != nil {
		errs["CONSOL"] = "Error: Invalid Period Format"
		return
	}

	lastCol := 18
	companyCol := lastCol + len(consolidationData.ConsolidationBridge)

	lastColCoord, err := excelize.CoordinatesToCellName(lastCol, 1)
	if err != nil {
		errs["CONSOL"] = "Error: Generate Configuration File Excel"
		fmt.Println(err)
		return
	}
	lastColCoord = strings.ReplaceAll(lastColCoord, "1", "")

	lastCompanyColCoord, err := excelize.CoordinatesToCellName(companyCol-1, 1)
	if err != nil {
		errs["CONSOL"] = "Error: Generate Configuration File Excel"
		fmt.Println(err)
		return
	}
	lastCompanyColCoord = strings.ReplaceAll(lastCompanyColCoord, "1", "")

	arrStyleColWidth := []map[string]interface{}{
		{"COL": "A", "WIDTH": 0.83},
		{"COL": "B", "WIDTH": 14.38},
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
		{"COL": "M", "WIDTH": 9.90},
		{"COL": "N", "WIDTH": 16.83},
		{"COL": "O", "WIDTH": 10.10},
		{"COL": "P", "WIDTH": 17.65},
		{"COL": "Q", "WIDTH": 22.14},
	}
	for _, v := range arrStyleColWidth {
		tmpColWidth := fmt.Sprintf("%f", v["WIDTH"])
		colWidth, _ := strconv.ParseFloat(tmpColWidth, 64)
		err = f.SetColWidth(sheet, fmt.Sprintf("%s", v["COL"]), fmt.Sprintf("%s", v["COL"]), colWidth)
		if err != nil {
			errs["CONSOL"] = "Error: Generate Configuration File Excel"
			log.Println(err)
			return
		}
	}

	lastColWidth := []float64{18.26, 9.90, 16.83, 10.10, 17.65, 22.14}

	for c := lastCol; c <= companyCol; c++ {
		cellCol, err := excelize.CoordinatesToCellName(c, 1)
		if err != nil {
			errs["CONSOL"] = "Error: Generate Configuration File Excel"
			fmt.Println(err)
			return
		}
		cellCol = strings.ReplaceAll(cellCol, "1", "")

		err = f.SetColWidth(sheet, cellCol, cellCol, 22.14)
		if err != nil {
			fmt.Println(err)
			return
		}
	}

	for i, cw := range lastColWidth {
		cellCol, err := excelize.CoordinatesToCellName(companyCol+(i+1), 1)
		if err != nil {
			errs["CONSOL"] = "Error: Generate Configuration File Excel"
			fmt.Println(err)
			return
		}
		cellCol = strings.ReplaceAll(cellCol, "1", "")

		err = f.SetColWidth(sheet, cellCol, cellCol, cw)
		if err != nil {
			errs["CONSOL"] = "Error: Generate Configuration File Excel"
			fmt.Println(err)
			return
		}
	}

	styleDefault, err := f.NewStyle(&excelize.Style{
		Font: &excelize.Font{
			Family: "Arial",
			Size:   10,
		},
	})
	if err != nil {
		errs["CONSOL"] = "Error: Generate Configuration File Excel"
		fmt.Println(err)
		return
	}
	err = f.SetColStyle(sheet, "A:Z", styleDefault)
	if err != nil {
		errs["CONSOL"] = "Error: Generate Configuration File Excel"
		fmt.Println(err)
		return
	}

	stylingBorderRightOnly, err := f.NewStyle(&excelize.Style{
		Border: []excelize.Border{
			{Type: "right", Color: "000000", Style: 1},
		},
	})
	if err != nil {
		errs["CONSOL"] = "Error: Generate Configuration File Excel"
		fmt.Println(err)
		return
	}
	stylingBorderTopOnly, err := f.NewStyle(&excelize.Style{
		Border: []excelize.Border{
			{Type: "top", Color: "000000", Style: 1},
		},
	})
	if err != nil {
		errs["CONSOL"] = "Error: Generate Configuration File Excel"
		fmt.Println(err)
		return
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
		errs["CONSOL"] = "Error: Generate Configuration File Excel"
		fmt.Println(err)
		return
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
		errs["CONSOL"] = "Error: Generate Configuration File Excel"
		fmt.Println(err)
		return
	}
	// numberFormat := "#,##"
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
		errs["CONSOL"] = "Error: Generate Configuration File Excel"
		fmt.Println(err)
		return
	}
	stylingSubTotalCurrency, err := f.NewStyle(&excelize.Style{
		Border: []excelize.Border{
			{Type: "left", Color: "000000", Style: 1},
			{Type: "right", Color: "000000", Style: 1},
			{Type: "top", Color: "000000", Style: 1},
			{Type: "bottom", Color: "000000", Style: 1},
		},
		Font: &excelize.Font{
			Bold: true,
		},
		NumFmt: 41,
	})
	if err != nil {
		errs["CONSOL"] = "Error: Generate Configuration File Excel"
		fmt.Println(err)
		return
	}

	stylingTotal, err := f.NewStyle(&excelize.Style{
		Border: []excelize.Border{
			{Type: "top", Color: "000000", Style: 1},
			{Type: "bottom", Color: "000000", Style: 1},
		},
		Font: &excelize.Font{
			Bold: true,
		},
		Fill: excelize.Fill{
			Type:    "pattern",
			Pattern: 1,
			Color:   []string{"#ccff33"},
		},
	})
	if err != nil {
		errs["CONSOL"] = "Error: Generate Configuration File Excel"
		fmt.Println(err)
		return
	}

	stylingTotalCurrency, err := f.NewStyle(&excelize.Style{
		
		Border: []excelize.Border{
			{Type: "left", Color: "000000", Style: 1},
			{Type: "right", Color: "000000", Style: 1},
			{Type: "top", Color: "000000", Style: 1},
			{Type: "bottom", Color: "000000", Style: 1},
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
		errs["CONSOL"] = "Error: Generate Configuration File Excel"
		fmt.Println(err)
		return
	}

	stylingTotalControl, err := f.NewStyle(&excelize.Style{
		// NumFmt: 7,
		Border: []excelize.Border{
			{Type: "top", Color: "000000", Style: 1},
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
		errs["CONSOL"] = "Error: Generate Configuration File Excel"
		fmt.Println(err)
		return
	}

	stylingTotalLvl1, err := f.NewStyle(&excelize.Style{
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
			Color:   []string{"#FF00FF"},
		},
		NumFmt: 41,
	})
	if err != nil {
		errs["CONSOL"] = "Error: Generate Configuration File Excel"
		fmt.Println(err)
		return
	}
	stylingTotalLvl2, err := f.NewStyle(&excelize.Style{
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
			Color:   []string{"#18C5E8"},
		},
		NumFmt: 41,
	})
	if err != nil {
		errs["CONSOL"] = "Error: Generate Configuration File Excel"
		fmt.Println(err)
		return
	}
	stylingTotalLvl3, err := f.NewStyle(&excelize.Style{
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
			Color:   []string{"#3ADA24"},
		},
		NumFmt: 41,
	})
	if err != nil {
		errs["CONSOL"] = "Error: Generate Configuration File Excel"
		fmt.Println(err)
		return
	}
	
	stylingTotalLvl4, err := f.NewStyle(&excelize.Style{
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
			Color:   []string{"#CCFF33"},
		},
		NumFmt: 41,
	})
	if err != nil {
		errs["CONSOL"] = "Error: Generate Configuration File Excel"
		fmt.Println(err)
		return
	}
	stylingTotalLvl5, err := f.NewStyle(&excelize.Style{
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
			Color:   []string{"#FF0000"},
		},
		NumFmt: 41,
	})
	if err != nil {
		errs["CONSOL"] = "Error: Generate Configuration File Excel"
		fmt.Println(err)
		return
	}
	styleColorTextLvl1, err := f.NewStyle(&excelize.Style{
		Font: &excelize.Font{
			Bold: true,
			// Color:   "#FF00FF",
		},
	})
	if err != nil {
		errs["CONSOL"] = "Error: Generate Configuration File Excel"
		fmt.Println(err)
		return
	}
	styleColorTextLvl2, err := f.NewStyle(&excelize.Style{
		Font: &excelize.Font{
			Bold: true,
			// Color:   "#00CCFF",
		},
	})
	if err != nil {
		errs["CONSOL"] = "Error: Generate Configuration File Excel"
		fmt.Println(err)
		return
	}
	styleColorTextLvl3, err := f.NewStyle(&excelize.Style{
		Font: &excelize.Font{
			Bold: true,
			// Color:   "#00FF00",
		},
	})
	if err != nil {
		errs["CONSOL"] = "Error: Generate Configuration File Excel"
		fmt.Println(err)
		return
	}
	styleColorTextLvl4, err := f.NewStyle(&excelize.Style{
		Font: &excelize.Font{
			Bold: true,
			// Color:   "#FFFF00",
		},
	})
	if err != nil {
		errs["CONSOL"] = "Error: Generate Configuration File Excel"
		fmt.Println(err)
		return
	}
	styleColorTextLvl5, err := f.NewStyle(&excelize.Style{
		Font: &excelize.Font{
			Bold: true,
			// Color:   "#FF0000",
		},
	})
	if err != nil {
		errs["CONSOL"] = "Error: Generate Configuration File Excel"
		fmt.Println(err)
		return
	}
	err = f.MergeCell(sheet, "C6", "E8")
	if err != nil {
		errs["CONSOL"] = "Error: Generate Configuration File Excel"
		fmt.Println(err)
		return
	}
	err = f.MergeCell(sheet, "B6", "B8")
	if err != nil {
		errs["CONSOL"] = "Error: Generate Configuration File Excel"
		fmt.Println(err)
		return
	}
	err = f.MergeCell(sheet, "F6", "F8")
	if err != nil {
		errs["CONSOL"] = "Error: Generate Configuration File Excel"
		fmt.Println(err)
		return
	}
	err = f.MergeCell(sheet, "H6", "K7")
	if err != nil {
		errs["CONSOL"] = "Error: Generate Configuration File Excel"
		fmt.Println(err)
		return
	}
	err = f.MergeCell(sheet, "M6", "P7")
	if err != nil {
		errs["CONSOL"] = "Error: Generate Configuration File Excel"
		fmt.Println(err)
		return
	}
	err = f.MergeCell(sheet, "H6", "K7")
	if err != nil {
		errs["CONSOL"] = "Error: Generate Configuration File Excel"
		fmt.Println(err)
		return
	}

	stylingCurrency, err := f.NewStyle(&excelize.Style{
		// NumFmt: 7,
		Border: []excelize.Border{
			{Type: "left", Color: "000000", Style: 1},
			{Type: "right", Color: "000000", Style: 1},
		},
		NumFmt: 41,
	})
	if err != nil {
		errs["CONSOL"] = "Error: Generate Configuration File Excel"
		fmt.Println(err)
		return
	}

	err = f.SetCellStyle(sheet, "F6", "F8", stylingHeader2)
	if err != nil {
		errs["CONSOL"] = "Error: Generate Configuration File Excel"
		fmt.Println(err)
		return
	}

	err = f.SetCellValue(sheet, "B2", "Company")
	if err != nil {
		errs["CONSOL"] = "Error: Generate Configuration File Excel"
		fmt.Println(err)
		return
	}
	err = f.SetCellValue(sheet, "B3", "Date")
	if err != nil {
		errs["CONSOL"] = "Error: Generate Configuration File Excel"
		fmt.Println(err)
		return
	}
	err = f.SetCellValue(sheet, "B4", "Subject")
	if err != nil {
		errs["CONSOL"] = "Error: Generate Configuration File Excel"
		fmt.Println(err)
		return
	}
	err = f.SetCellValue(sheet, "B6", "No Akun")
	if err != nil {
		errs["CONSOL"] = "Error: Generate Configuration File Excel"
		fmt.Println(err)
		return
	}
	err = f.SetCellValue(sheet, "C2", ":")
	if err != nil {
		errs["CONSOL"] = "Error: Generate Configuration File Excel"
		fmt.Println(err)
		return
	}
	err = f.SetCellValue(sheet, "C3", ":")
	if err != nil {
		errs["CONSOL"] = "Error: Generate Configuration File Excel"
		fmt.Println(err)
		return
	}
	err = f.SetCellValue(sheet, "C4", ":")
	if err != nil {
		errs["CONSOL"] = "Error: Generate Configuration File Excel"
		fmt.Println(err)
		return
	}
	err = f.SetCellValue(sheet, "C6", "Keterangan")
	if err != nil {
		errs["CONSOL"] = "Error: Generate Configuration File Excel"
		fmt.Println(err)
		return
	}
	err = f.SetCellValue(sheet, "D2", consolidationData.Company.Name)
	if err != nil {
		errs["CONSOL"] = "Error: Generate Configuration File Excel"
		fmt.Println(err)
		return
	}
	err = f.SetCellValue(sheet, "D3", datePeriod.Format("02-Jan-06"))
	if err != nil {
		errs["CONSOL"] = "Error: Generate Configuration File Excel"
		fmt.Println(err)
		return
	}
	err = f.SetCellValue(sheet, "D4", "Detail Aset, Liabilitas, Ekuitas")
	if err != nil {
		errs["CONSOL"] = "Error: Generate Configuration File Excel"
		fmt.Println(err)
		return
	}
	err = f.SetCellValue(sheet, "F6", "WP Reff")
	if err != nil {
		errs["CONSOL"] = "Error: Generate Configuration File Excel"
		fmt.Println(err)
		return
	}
	err = f.SetCellValue(sheet, "G6", consolidationData.Company.Code)
	if err != nil {
		errs["CONSOL"] = "Error: Generate Configuration File Excel"
		fmt.Println(err)
		return
	}
	err = f.SetCellValue(sheet, "G7", "Unaudited")
	if err != nil {
		errs["CONSOL"] = "Error: Generate Configuration File Excel"
		fmt.Println(err)
		return
	}
	err = f.SetCellValue(sheet, "G8", datePeriod.Format("02-Jan-06"))
	if err != nil {
		errs["CONSOL"] = "Error: Generate Configuration File Excel"
		fmt.Println(err)
		return
	}
	err = f.SetCellValue(sheet, "H6", "Proforma Modal")
	if err != nil {
		errs["CONSOL"] = "Error: Generate Configuration File Excel"
		fmt.Println(err)
		return
	}
	err = f.SetCellValue(sheet, "I8", "Debet")
	if err != nil {
		errs["CONSOL"] = "Error: Generate Configuration File Excel"
		fmt.Println(err)
		return
	}
	err = f.SetCellValue(sheet, "K8", "Kredit")
	if err != nil {
		errs["CONSOL"] = "Error: Generate Configuration File Excel"
		fmt.Println(err)
		return
	}
	err = f.SetCellFormula(sheet, "L6", "=G6")
	if err != nil {
		errs["CONSOL"] = "Error: Generate Configuration File Excel"
		fmt.Println(err)
		return
	}
	err = f.SetCellFormula(sheet, "L7", "=G7")
	if err != nil {
		errs["CONSOL"] = "Error: Generate Configuration File Excel"
		fmt.Println(err)
		return
	}
	err = f.SetCellFormula(sheet, "L8", "=G8")
	if err != nil {
		errs["CONSOL"] = "Error: Generate Configuration File Excel"
		fmt.Println(err)
		return
	}
	err = f.SetCellValue(sheet, "M6", "Cost to Equity")
	if err != nil {
		errs["CONSOL"] = "Error: Generate Configuration File Excel"
		fmt.Println(err)
		return
	}
	err = f.SetCellValue(sheet, "N8", "Debet")
	if err != nil {
		errs["CONSOL"] = "Error: Generate Configuration File Excel"
		fmt.Println(err)
		return
	}
	err = f.SetCellValue(sheet, "P8", "Kredit")
	if err != nil {
		errs["CONSOL"] = "Error: Generate Configuration File Excel"
		fmt.Println(err)
		return
	}
	err = f.SetCellFormula(sheet, "Q6", "=G6")
	if err != nil {
		errs["CONSOL"] = "Error: Generate Configuration File Excel"
		fmt.Println(err)
		return
	}
	err = f.SetCellFormula(sheet, "Q7", "=G7")
	if err != nil {
		errs["CONSOL"] = "Error: Generate Configuration File Excel"
		fmt.Println(err)
		return
	}
	err = f.SetCellFormula(sheet, "Q8", "=G8")
	if err != nil {
		errs["CONSOL"] = "Error: Generate Configuration File Excel"
		fmt.Println(err)
		return
	}

	colCombine, err := excelize.CoordinatesToCellName(companyCol, 1)
	if err != nil {
		errs["CONSOL"] = "Error: Generate Configuration File Excel"
		fmt.Println(err)
		return
	}
	colCombine = strings.ReplaceAll(colCombine, "1", "")
	colElimination, err := excelize.CoordinatesToCellName(companyCol+1, 1)
	if err != nil {
		errs["CONSOL"] = "Error: Generate Configuration File Excel"
		fmt.Println(err)
		return
	}
	colElimination = strings.ReplaceAll(colElimination, "1", "")
	colEliminationReffKredit, err := excelize.CoordinatesToCellName(companyCol+3, 1)
	if err != nil {
		errs["CONSOL"] = "Error: Generate Configuration File Excel"
		fmt.Println(err)
		return
	}
	colEliminationKredit, err := excelize.CoordinatesToCellName(companyCol+4, 1)
	if err != nil {
		errs["CONSOL"] = "Error: Generate Configuration File Excel"
		fmt.Println(err)
		return
	}
	colEliminationReffKredit = strings.ReplaceAll(colEliminationReffKredit, "1", "")
	colEliminationKredit = strings.ReplaceAll(colEliminationKredit, "1", "")

	colEliminationReffDebit, err := excelize.CoordinatesToCellName(companyCol+1, 1)
	if err != nil {
		errs["CONSOL"] = "Error: Generate Configuration File Excel"
		fmt.Println(err)
		return
	}
	colEliminationDebit, err := excelize.CoordinatesToCellName(companyCol+2, 1)
	if err != nil {
		errs["CONSOL"] = "Error: Generate Configuration File Excel"
		fmt.Println(err)
		return
	}
	colEliminationReffDebit = strings.ReplaceAll(colEliminationReffDebit, "1", "")
	colEliminationDebit = strings.ReplaceAll(colEliminationDebit, "1", "")
	colEliminationLast, err := excelize.CoordinatesToCellName(companyCol+3, 1)
	if err != nil {
		errs["CONSOL"] = "Error: Generate Configuration File Excel"
		fmt.Println(err)
		return
	}
	colEliminationLast = strings.ReplaceAll(colEliminationLast, "1", "")
	colKonsolidasi, err := excelize.CoordinatesToCellName(companyCol+5, 1)
	if err != nil {
		errs["CONSOL"] = "Error: Generate Configuration File Excel"
		fmt.Println(err)
		return
	}
	colKonsolidasi = strings.ReplaceAll(colKonsolidasi, "1", "")
	err = f.SetCellValue(sheet, fmt.Sprintf("%s6", colCombine), "Combine")
	if err != nil {
		errs["CONSOL"] = "Error: Generate Configuration File Excel"
		fmt.Println(err)
		return
	}
	err = f.SetCellFormula(sheet, fmt.Sprintf("%s7", colCombine), "=Q7")
	if err != nil {
		errs["CONSOL"] = "Error: Generate Configuration File Excel"
		fmt.Println(err)
		return
	}
	err = f.SetCellFormula(sheet, fmt.Sprintf("%s8", colCombine), "=Q8")
	if err != nil {
		errs["CONSOL"] = "Error: Generate Configuration File Excel"
		fmt.Println(err)
		return
	}

	err = f.MergeCell(sheet, fmt.Sprintf("%s6", colElimination), fmt.Sprintf("%s7", colEliminationKredit))
	if err != nil {
		errs["CONSOL"] = "Error: Generate Configuration File Excel"
		fmt.Println(err)
		return
	}

	err = f.SetCellValue(sheet, fmt.Sprintf("%s6", colElimination), "Elimination")
	if err != nil {
		errs["CONSOL"] = "Error: Generate Configuration File Excel"
		fmt.Println(err)
		return
	}
	err = f.SetCellValue(sheet, fmt.Sprintf("%s8", colEliminationDebit), "Debit")
	if err != nil {
		errs["CONSOL"] = "Error: Generate Configuration File Excel"
		fmt.Println(err)
		return
	}
	err = f.SetCellValue(sheet, fmt.Sprintf("%s8", colEliminationKredit), "Kredit")
	if err != nil {
		errs["CONSOL"] = "Error: Generate Configuration File Excel"
		fmt.Println(err)
		return
	}

	err = f.SetCellValue(sheet, fmt.Sprintf("%s6", colKonsolidasi), "Konsolidasi")
	if err != nil {
		errs["CONSOL"] = "Error: Generate Configuration File Excel"
		fmt.Println(err)
		return
	}
	err = f.SetCellFormula(sheet, fmt.Sprintf("%s7", colKonsolidasi), "=G7")
	if err != nil {
		errs["CONSOL"] = "Error: Generate Configuration File Excel"
		fmt.Println(err)
		return
	}
	err = f.SetCellFormula(sheet, fmt.Sprintf("%s8", colKonsolidasi), "=G8")
	if err != nil {
		errs["CONSOL"] = "Error: Generate Configuration File Excel"
		fmt.Println(err)
		return
	}

	err = f.SetCellStyle(sheet, "B6", fmt.Sprintf("%s8", colKonsolidasi), stylingHeader)
	if err != nil {
		errs["CONSOL"] = "Error: Generate Configuration File Excel"
		fmt.Println(err)
		return
	}

	for i, brdge := range consolidationData.ConsolidationBridge {
		vCompany := brdge.Company
		colCompany, err := excelize.CoordinatesToCellName(lastCol+i, 1)
		if err != nil {
			fmt.Println(err)
			return
		}
		colCompany = strings.ReplaceAll(colCompany, "1", "")
		err = f.SetCellValue(sheet, fmt.Sprintf("%s6", colCompany), vCompany.Code)
		if err != nil {
			errs["CONSOL"] = "Error: Generate Configuration File Excel"
			fmt.Println(err)
			return
		}
		err = f.SetCellFormula(sheet, fmt.Sprintf("%s7", colCompany), "=G7")
		if err != nil {
			errs["CONSOL"] = "Error: Generate Configuration File Excel"
			fmt.Println(err)
			return
		}
		err = f.SetCellFormula(sheet, fmt.Sprintf("%s8", colCompany), "=G8")
		if err != nil {
			errs["CONSOL"] = "Error: Generate Configuration File Excel"
			fmt.Println(err)
			return
		}
	}

	// find summary journal
	summaryJCTE, err := s.JcteRepository.FindSummary(ctx, &consolidationData.ID)
	if err != nil {
		return
	}

	summaryJPM, err := s.JpmRepository.FindSummary(ctx, &consolidationData.ID)
	if err != nil {
		return
	}

	summaryJELIM, err := s.JelimRepository.FindSummary(ctx, &consolidationData.ID)
	if err != nil {
		return
	}

	row := 9
	var summary []map[string]interface{}
	// var total []map[string]interface{}
	rowCode := make(map[string]int)
	isAutoSum := make(map[string]bool)
	tbRowCode := make(map[string]int)
	customRow := make(map[string]string)

	//jpm
	fmlJpmDr := make(map[string]string)
	fmlJpmCr := make(map[string]string)
	reffJpmDr := make(map[string]string)
	reffJpmCr := make(map[string]string)
	sheetJpm := "JPM"
	rowsJpm, err := f.GetRows(sheetJpm)
	if err != nil {
		return
	}
	var lineJpm []string

	for _, jpm := range rowsJpm {
		lineJpm = append(lineJpm, "line")

		if len(jpm) == 0 {
			continue
		}
		if len(jpm) == 2 {
			continue
		}
		if len(jpm) == 5 {
			continue
		}
		// if jpm[1] != vTbDetail.Code {
		// 	continue
		// }
		if jpm[6] != "" {
			if fmlJpmDr[jpm[1]] != "" {
				rumusAllJpmDr := strings.Replace(fmlJpmDr[jpm[1]], "=", "+", 1)
				rumusJpmDr := "=" + sheetJpm + "!G" + strconv.Itoa(len(lineJpm)) + rumusAllJpmDr
				fmt.Sprintln(rumusJpmDr)
				rumusReffAllJpmDr := reffJpmDr[jpm[1]]
				rumusReffJpmDr := rumusReffAllJpmDr + "&" + `","` + "&" + sheetJpm + "!C" + strconv.Itoa(len(lineJpm))
				// if jpm[1] == jpm[2] {
				// 	f.SetCellFormula(sheet, fmt.Sprintf("I%d", jpmw), rumus)
				// 	f.SetCellFormula(sheet, fmt.Sprintf("H%d", jpmw), rumusReff)
				// }
				fmlJpmDr[jpm[1]] = rumusJpmDr
				reffJpmDr[jpm[1]] = rumusReffJpmDr
			}
			if fmlJpmDr[jpm[1]] == "" {
				rumusJpmDr := "=" + sheetJpm + "!G" + strconv.Itoa(len(lineJpm))
				rumusReffJpmDr := "=" + sheetJpm + "!C" + strconv.Itoa(len(lineJpm))
				// if jpm[1] == jpm[2] {
				// 	f.SetCellFormula(sheet, fmt.Sprintf("I%d", jpmw), rumus)
				// 	f.SetCellFormula(sheet, fmt.Sprintf("H%d", jpmw), rumusReffJpm)
				// }
				fmlJpmDr[jpm[1]] = rumusJpmDr
				reffJpmDr[jpm[1]] = rumusReffJpmDr
			}
		}
		if len(jpm) == 8 && jpm[7] != "" {
			if fmlJpmCr[jpm[1]] != "" {
				rumusAllJpmCr := strings.Replace(fmlJpmCr[jpm[1]], "=", "+", 1)
				fmt.Sprintln(rumusAllJpmCr)
				rumusJpmCr := "=" + sheetJpm + "!H" + strconv.Itoa(len(lineJpm)) + rumusAllJpmCr
				rumusReffAllJpmCr := reffJpmCr[jpm[1]]
				rumusReffJpmCr := rumusReffAllJpmCr + "&" + `","` + "&" + sheetJpm + "!C" + strconv.Itoa(len(lineJpm))
				fmt.Sprintln(rumusJpmCr)
				// if jpm[1] == jpm[2] {
				// 	f.SetCellFormula(sheet, fmt.Sprintf("K%d", jpmw), rumus)
				// 	f.SetCellFormula(sheet, fmt.Sprintf("J%d", jpmw), rumusReff)
				// }
				fmlJpmCr[jpm[1]] = rumusJpmCr
				reffJpmCr[jpm[1]] = rumusReffJpmCr
			}
			if fmlJpmCr[jpm[1]] == "" {
				rumusJpmCr := "=" + sheetJpm + "!H" + strconv.Itoa(len(lineJpm))
				rumusReffJpmCr := "=" + sheetJpm + "!C" + strconv.Itoa(len(lineJpm))
				// if jpm[1] == jpm[2] {
				// 	f.SetCellFormula(sheet, fmt.Sprintf("K%d", jpmw), rumus)
				// 	f.SetCellFormula(sheet, fmt.Sprintf("J%d", jpmw), rumusReffJpm)
				// }
				fmlJpmCr[jpm[1]] = rumusJpmCr
				reffJpmCr[jpm[1]] = rumusReffJpmCr
			}

		}

		if len(jpm) == 9 && jpm[8] != "" {
			if fmlJpmDr[jpm[1]] != "" {
				rumusAllJpmDr := strings.Replace(fmlJpmDr[jpm[1]], "=", "+", 1)
				fmt.Sprintln(rumusAllJpmDr)
				rumusJpmDr := "=" + sheetJpm + "!I" + strconv.Itoa(len(lineJpm)) + rumusAllJpmDr
				rumusReffAllJpmDr := reffJpmDr[jpm[1]]
				rumusReffJpmDr := rumusReffAllJpmDr + "&" + `","` + "&" + sheetJpm + "!C" + strconv.Itoa(len(lineJpm))
				fmt.Sprintln(rumusJpmDr)
				// if jpm[1] == jpm[2] {
				// 	f.SetCellFormula(sheet, fmt.Sprintf("I%d", jpmw), rumus)
				// 	f.SetCellFormula(sheet, fmt.Sprintf("H%d", jpmw), rumusReff)
				// }
				fmlJpmDr[jpm[1]] = rumusJpmDr
				reffJpmDr[jpm[1]] = rumusReffJpmDr
			}
			if fmlJpmDr[jpm[1]] == "" {
				rumusJpmDr := "=" + sheetJpm + "!I" + strconv.Itoa(len(lineJpm))
				rumusReffJpmDr := "=" + sheetJpm + "!C" + strconv.Itoa(len(lineJpm))
				// if jpm[1] == jpm[2] {
				// 	f.SetCellFormula(sheet, fmt.Sprintf("I%d", jpmw), rumus)
				// 	f.SetCellFormula(sheet, fmt.Sprintf("H%d", jpmw), rumusReff)
				// }
				fmlJpmDr[jpm[1]] = rumusJpmDr
				reffJpmDr[jpm[1]] = rumusReffJpmDr
			}

		}
		if len(jpm) == 10 && jpm[9] != "" {
			if fmlJpmCr[jpm[1]] != "" {
				rumusAllJpmCr := strings.Replace(fmlJpmCr[jpm[1]], "=", "+", 1)
				fmt.Sprintln(rumusAllJpmCr)
				rumusJpmCr := "=" + sheetJpm + "!J" + strconv.Itoa(len(lineJpm)) + rumusAllJpmCr
				rumusReffAllJpmCr := reffJpmCr[jpm[1]]
				rumusReffJpmCr := rumusReffAllJpmCr + "&" + `","` + "&" + sheetJpm + "!C" + strconv.Itoa(len(lineJpm))
				fmt.Sprintln(rumusJpmCr)
				// if jpm[1] == jpm[2] {
				// 	f.SetCellFormula(sheet, fmt.Sprintf("K%d", jpmw), rumus)
				// 	f.SetCellFormula(sheet, fmt.Sprintf("J%d", jpmw), rumusReff)
				// }
				fmlJpmCr[jpm[1]] = rumusJpmCr
				reffJpmCr[jpm[1]] = rumusReffJpmCr
			}
			if fmlJpmCr[jpm[1]] == "" {
				rumusReffCr := "=" + sheetJpm + "!C" + strconv.Itoa(len(lineJpm))
				rumusJpmCr := "=" + sheetJpm + "!J" + strconv.Itoa(len(lineJpm))
				// if jpm[1] == jpm[2] {
				// 	f.SetCellFormula(sheet, fmt.Sprintf("K%d", jpmw), rumus)
				// 	f.SetCellFormula(sheet, fmt.Sprintf("J%d", jpmw), rumusReff)
				// }
				fmlJpmCr[jpm[1]] = rumusJpmCr
				reffJpmCr[jpm[1]] = rumusReffCr
			}
		}
		if len(jpm) == 10 && jpm[8] == "" && jpm[9] == "" {
			// rumusReffjpmCr := "=" + sheetjpm + "!C" + strconv.Itoa(len(linejpm))
			rumusJpmDr := "=" + sheetJpm + "!I" + strconv.Itoa(len(lineJpm))
			rumusJpmCr := "=" + sheetJpm + "!J" + strconv.Itoa(len(lineJpm))
			// if jpm[1] == jpm[2] {
			// 	f.SetCellFormula(sheet, fmt.Sprintf("K%d", jpmw), rumus)
			// 	f.SetCellFormula(sheet, fmt.Sprintf("J%d", jpmw), rumusReff)
			// }
			fmlJpmDr["CONTROL_TO_ADJUSTMENT_SHEET"] = rumusJpmDr
			fmlJpmCr["CONTROL_TO_ADJUSTMENT_SHEET"] = rumusJpmCr
			// reffjpmCr[jpm[1]] = rumusReffjpmCr
		
		}
	}
	//jcte
	fmlJcteDr := make(map[string]string)
	reffJcteDr := make(map[string]string)
	fmlJcteCr := make(map[string]string)
	reffJcteCr := make(map[string]string)
	sheetJcte := "JCTE"
	rowsJcte, err := f.GetRows(sheetJcte)
	if err != nil {
		return
	}
	var lineJcte []string

	for _, Jcte := range rowsJcte {
		lineJcte = append(lineJcte, "lineJcte")

		if len(Jcte) == 0 {
			continue
		}
		if len(Jcte) == 2 {
			continue
		}
		if len(Jcte) == 5 {
			continue
		}
		// if Jcte[1] != vTbDetail.Code {
		// 	continue
		// }
		if Jcte[6] != "" {
			if fmlJcteDr[Jcte[1]] != "" {
				rumusAllJcteDr := strings.Replace(fmlJcteDr[Jcte[1]], "=", "+", 1)
				rumusJcteDr := "=" + sheetJcte + "!G" + strconv.Itoa(len(lineJcte)) + rumusAllJcteDr
				fmt.Sprintln(rumusJcteDr)
				rumusReffAllJcteDr := reffJcteDr[Jcte[1]]
				rumusReffJcteDr := rumusReffAllJcteDr + "&" + `","` + "&" + sheetJcte + "!C" + strconv.Itoa(len(lineJcte))
				// if Jcte[1] == Jcte[2] {
				// 	f.SetCellFormula(sheet, fmt.Sprintf("I%d", Jctew), rumus)
				// 	f.SetCellFormula(sheet, fmt.Sprintf("H%d", Jctew), rumusReff)
				// }
				fmlJcteDr[Jcte[1]] = rumusJcteDr
				reffJcteDr[Jcte[1]] = rumusReffJcteDr
			}
			if fmlJcteDr[Jcte[1]] == "" {
				rumusJcteDr := "=" + sheetJcte + "!G" + strconv.Itoa(len(lineJcte))
				rumusReffJcteDr := "=" + sheetJcte + "!C" + strconv.Itoa(len(lineJcte))
				// if Jcte[1] == Jcte[2] {
				// 	f.SetCellFormula(sheet, fmt.Sprintf("I%d", Jctew), rumus)
				// 	f.SetCellFormula(sheet, fmt.Sprintf("H%d", Jctew), rumusReff)
				// }
				fmlJcteDr[Jcte[1]] = rumusJcteDr
				reffJcteDr[Jcte[1]] = rumusReffJcteDr
			}
		}
		if len(Jcte) == 8 && Jcte[7] != "" {
			if fmlJcteCr[Jcte[1]] != "" {
				rumusAllJcteCr := strings.Replace(fmlJcteCr[Jcte[1]], "=", "+", 1)
				fmt.Sprintln(rumusAllJcteCr)
				rumusJcteCr := "=" + sheetJcte + "!H" + strconv.Itoa(len(lineJcte)) + rumusAllJcteCr
				rumusReffAllJcteCr := reffJcteCr[Jcte[1]]
				rumusReffJcteCr := rumusReffAllJcteCr + "&" + `","` + "&" + sheetJcte + "!C" + strconv.Itoa(len(lineJcte))
				fmt.Sprintln(rumusJcteCr)
				// if Jcte[1] == Jcte[2] {
				// 	f.SetCellFormula(sheet, fmt.Sprintf("K%d", Jctew), rumus)
				// 	f.SetCellFormula(sheet, fmt.Sprintf("J%d", Jctew), rumusReff)
				// }
				fmlJcteCr[Jcte[1]] = rumusJcteCr
				reffJcteCr[Jcte[1]] = rumusReffJcteCr
			}
			if fmlJcteCr[Jcte[1]] == "" {
				rumusJcteCr := "=" + sheetJcte + "!H" + strconv.Itoa(len(lineJcte))
				rumusReffJcteCr := "=" + sheetJcte + "!C" + strconv.Itoa(len(lineJcte))
				// if Jcte[1] == Jcte[2] {
				// 	f.SetCellFormula(sheet, fmt.Sprintf("K%d", Jctew), rumus)
				// 	f.SetCellFormula(sheet, fmt.Sprintf("J%d", Jctew), rumusReff)
				// }
				fmlJcteCr[Jcte[1]] = rumusJcteCr
				reffJcteCr[Jcte[1]] = rumusReffJcteCr
			}

		}

		if len(Jcte) == 9 && Jcte[8] != "" {
			if fmlJcteDr[Jcte[1]] != "" {
				rumusAllJcteDr := strings.Replace(fmlJcteDr[Jcte[1]], "=", "+", 1)
				fmt.Sprintln(rumusAllJcteDr)
				rumusJcteDr := "=" + sheetJcte + "!I" + strconv.Itoa(len(lineJcte)) + rumusAllJcteDr
				rumusReffAllJcteDr := reffJcteDr[Jcte[1]]
				rumusReffJcteDr := rumusReffAllJcteDr + "&" + `","` + "&" + sheetJcte + "!C" + strconv.Itoa(len(lineJcte))
				fmt.Sprintln(rumusJcteDr)
				// if jpm[1] == jpm[2] {
				// 	f.SetCellFormula(sheet, fmt.Sprintf("I%d", jpmw), rumus)
				// 	f.SetCellFormula(sheet, fmt.Sprintf("H%d", jpmw), rumusReff)
				// }
				fmlJcteDr[Jcte[1]] = rumusJcteDr
				reffJcteDr[Jcte[1]] = rumusReffJcteDr
			}
			if fmlJcteDr[Jcte[1]] == "" {
				rumusJcteDr := "=" + sheetJcte + "!I" + strconv.Itoa(len(lineJcte))
				rumusReffJcteDr := "=" + sheetJcte + "!C" + strconv.Itoa(len(lineJcte))
				// if jpm[1] == jpm[2] {
				// 	f.SetCellFormula(sheet, fmt.Sprintf("I%d", jpmw), rumus)
				// 	f.SetCellFormula(sheet, fmt.Sprintf("H%d", jpmw), rumusReff)
				// }
				fmlJcteDr[Jcte[1]] = rumusJcteDr
				reffJcteDr[Jcte[1]] = rumusReffJcteDr
			}

		}
		if len(Jcte) == 10 && Jcte[9] != "" {
			if fmlJcteCr[Jcte[1]] != "" {
				rumusAllJcteCr := strings.Replace(fmlJcteCr[Jcte[1]], "=", "+", 1)
				fmt.Sprintln(rumusAllJcteCr)
				rumusJcteCr := "=" + sheetJcte + "!J" + strconv.Itoa(len(lineJcte)) + rumusAllJcteCr
				rumusReffAllJcteCr := reffJcteCr[Jcte[1]]
				rumusReffJcteCr := rumusReffAllJcteCr + "&" + `","` + "&" + sheetJcte + "!C" + strconv.Itoa(len(lineJcte))
				fmt.Sprintln(rumusJcteCr)
				// if Jcte[1] == Jcte[2] {
				// 	f.SetCellFormula(sheet, fmt.Sprintf("K%d", Jctew), rumusJcte)
				// 	f.SetCellFormula(sheet, fmt.Sprintf("J%d", Jctew), rumusJcteReff)
				// }
				fmlJcteCr[Jcte[1]] = rumusJcteCr
				reffJcteCr[Jcte[1]] = rumusReffJcteCr
			}
			if fmlJcteCr[Jcte[1]] == "" {
				rumusReffJcteCr := "=" + sheetJcte + "!C" + strconv.Itoa(len(lineJcte))
				rumusJcteCr := "=" + sheetJcte + "!J" + strconv.Itoa(len(lineJcte))
				// if Jcte[1] == Jcte[2] {
				// 	f.SetCellFormula(sheet, fmt.Sprintf("K%d", Jctew), rumus)
				// 	f.SetCellFormula(sheet, fmt.Sprintf("J%d", Jctew), rumusReff)
				// }
				fmlJcteCr[Jcte[1]] = rumusJcteCr
				reffJcteCr[Jcte[1]] = rumusReffJcteCr
			}
		}
		if len(Jcte) == 10 && Jcte[8] == "" && Jcte[9] == "" {
			// rumusReffJcteCr := "=" + sheetJcte + "!C" + strconv.Itoa(len(lineJcte))
			rumusJcteDr := "=" + sheetJcte + "!I" + strconv.Itoa(len(lineJcte))
			rumusJcteCr := "=" + sheetJcte + "!J" + strconv.Itoa(len(lineJcte))
			// if jpm[1] == jpm[2] {
			// 	f.SetCellFormula(sheet, fmt.Sprintf("K%d", jpmw), rumus)
			// 	f.SetCellFormula(sheet, fmt.Sprintf("J%d", jpmw), rumusReff)
			// }
			fmlJcteDr["CONTROL_TO_ADJUSTMENT_SHEET"] = rumusJcteDr
			fmlJcteCr["CONTROL_TO_ADJUSTMENT_SHEET"] = rumusJcteCr
			// reffJcteCr[Jcte[1]] = rumusReffJcteCr
		
		}
	}

	//jelim
	fmlJelimCr := make(map[string]string)
	reffJelimCr := make(map[string]string)
	fmlJelimDr := make(map[string]string)
	reffJelimDr := make(map[string]string)
	sheetJelim := "JELIM"
	rowsJelim, err := f.GetRows(sheetJelim)
	if err != nil {
		return
	}
	var lineJelim []string

	for _, Jelim := range rowsJelim {
		lineJelim = append(lineJelim, "lineJelim")

		
		if len(Jelim) == 0 {
			continue
		}
		if len(Jelim) == 2 {
			continue
		}
		if len(Jelim) == 5 {
			continue
		}
	
		if Jelim[6] != "" {
			if fmlJelimDr[Jelim[1]] != "" {
				rumusAllJelimDr := strings.Replace(fmlJelimDr[Jelim[1]], "=", "+", 1)
				rumusJelimDr := "=" + sheetJelim + "!G" + strconv.Itoa(len(lineJelim)) + rumusAllJelimDr
				fmt.Sprintln(rumusJelimDr)
				rumusReffAllJelimDr := reffJelimDr[Jelim[1]]
				rumusReffJelimDr := rumusReffAllJelimDr + "&" + `","` + "&" + sheetJelim + "!C" + strconv.Itoa(len(lineJelim))
				// if Jelim[1] == Jelim[2] {
				// 	f.SetCellFormula(sheet, fmt.Sprintf("I%d", Jelimw), rumus)
				// 	f.SetCellFormula(sheet, fmt.Sprintf("H%d", Jelimw), rumusReff)
				// }
				fmlJelimDr[Jelim[1]] = rumusJelimDr
				reffJelimDr[Jelim[1]] = rumusReffJelimDr
			}
			if fmlJelimDr[Jelim[1]] == "" {
				rumusJelimDr := "=" + sheetJelim + "!G" + strconv.Itoa(len(lineJelim))
				rumusReffJelimDr := "=" + sheetJelim + "!C" + strconv.Itoa(len(lineJelim))
				// if Jelim[1] == Jelim[2] {
				// 	f.SetCellFormula(sheet, fmt.Sprintf("I%d", Jelimw), rumus)
				// 	f.SetCellFormula(sheet, fmt.Sprintf("H%d", Jelimw), rumusReff)
				// }
				fmlJelimDr[Jelim[1]] = rumusJelimDr
				reffJelimDr[Jelim[1]] = rumusReffJelimDr
			}
		}
		if len(Jelim) == 8 && Jelim[7] != "" {
			if fmlJelimCr[Jelim[1]] != "" {
				rumusAllJelimCr := strings.Replace(fmlJelimCr[Jelim[1]], "=", "+", 1)
				fmt.Sprintln(rumusAllJelimCr)
				rumusJelimCr := "=" + sheetJelim + "!H" + strconv.Itoa(len(lineJelim)) + rumusAllJelimCr
				rumusReffAllJelimCr := reffJelimCr[Jelim[1]]
				rumusReffJelimCr := rumusReffAllJelimCr + "&" + `","` + "&" + sheetJelim + "!C" + strconv.Itoa(len(lineJelim))
				fmt.Sprintln(rumusJelimCr)
				// if Jelim[1] == Jelim[2] {
				// 	f.SetCellFormula(sheet, fmt.Sprintf("K%d", Jelimw), rumus)
				// 	f.SetCellFormula(sheet, fmt.Sprintf("J%d", Jelimw), rumusReff)
				// }
				fmlJelimCr[Jelim[1]] = rumusJelimCr
				reffJelimCr[Jelim[1]] = rumusReffJelimCr
			}
			if fmlJelimCr[Jelim[1]] == "" {
				rumusJelimCr := "=" + sheetJelim + "!H" + strconv.Itoa(len(lineJelim))
				rumusReffJelimCr := "=" + sheetJelim + "!C" + strconv.Itoa(len(lineJelim))
				// if Jelim[1] == Jelim[2] {
				// 	f.SetCellFormula(sheet, fmt.Sprintf("K%d", Jelimw), rumus)
				// 	f.SetCellFormula(sheet, fmt.Sprintf("J%d", Jelimw), rumusReff)
				// }
				fmlJelimCr[Jelim[1]] = rumusJelimCr
				reffJelimCr[Jelim[1]] = rumusReffJelimCr
			}

		}
		if len(Jelim) == 9 && Jelim[8] != "" {
			if fmlJelimDr[Jelim[1]] != "" {
				rumusAllJelimDr := strings.Replace(fmlJelimDr[Jelim[1]], "=", "+", 1)
				fmt.Sprintln(rumusAllJelimDr)
				rumusJelimDr := "=" + sheetJelim + "!I" + strconv.Itoa(len(lineJelim)) + rumusAllJelimDr
				rumusReffAllJelimDr := reffJelimDr[Jelim[1]]
				rumusReffJelimDr := rumusReffAllJelimDr + "&" + `","` + "&" + sheetJelim + "!C" + strconv.Itoa(len(lineJelim))
				fmt.Sprintln(rumusJelimDr)
				// if jpm[1] == jpm[2] {
				// 	f.SetCellFormula(sheet, fmt.Sprintf("I%d", jpmw), rumus)
				// 	f.SetCellFormula(sheet, fmt.Sprintf("H%d", jpmw), rumusReff)
				// }
				fmlJelimDr[Jelim[1]] = rumusJelimDr
				reffJelimDr[Jelim[1]] = rumusReffJelimDr
			}
			if fmlJelimDr[Jelim[1]] == "" {
				rumusJelimDr := "=" + sheetJelim + "!I" + strconv.Itoa(len(lineJelim))
				rumusReffJelimDr := "=" + sheetJelim + "!C" + strconv.Itoa(len(lineJelim))
				// if jpm[1] == jpm[2] {
				// 	f.SetCellFormula(sheet, fmt.Sprintf("I%d", jpmw), rumus)
				// 	f.SetCellFormula(sheet, fmt.Sprintf("H%d", jpmw), rumusReff)
				// }
				fmlJelimDr[Jelim[1]] = rumusJelimDr
				reffJelimDr[Jelim[1]] = rumusReffJelimDr
			}
		}
		if len(Jelim) == 10 && Jelim[9] != "" {
			if fmlJelimCr[Jelim[1]] != "" {
				rumusAllJelimCr := strings.Replace(fmlJelimCr[Jelim[1]], "=", "+", 1)
				fmt.Sprintln(rumusAllJelimCr)
				rumusJelimCr := "=" + sheetJelim + "!J" + strconv.Itoa(len(lineJelim)) + rumusAllJelimCr
				rumusReffAllJelimCr := reffJelimCr[Jelim[1]]
				rumusReffJelimCr := rumusReffAllJelimCr + "&" + `","` + "&" + sheetJelim + "!C" + strconv.Itoa(len(lineJelim))
				fmt.Sprintln(rumusJelimCr)
				// if jpm[1] == jpm[2] {
				// 	f.SetCellFormula(sheet, fmt.Sprintf("K%d", jpmw), rumus)
				// 	f.SetCellFormula(sheet, fmt.Sprintf("J%d", jpmw), rumusReff)
				// }
				fmlJelimCr[Jelim[1]] = rumusJelimCr
				reffJelimCr[Jelim[1]] = rumusReffJelimCr
			}
			if fmlJelimCr[Jelim[1]] == "" {
				rumusReffJelimCr := "=" + sheetJelim + "!C" + strconv.Itoa(len(lineJelim))
				rumusJelimCr := "=" + sheetJelim + "!J" + strconv.Itoa(len(lineJelim))
				// if jpm[1] == jpm[2] {
				// 	f.SetCellFormula(sheet, fmt.Sprintf("K%d", jpmw), rumus)
				// 	f.SetCellFormula(sheet, fmt.Sprintf("J%d", jpmw), rumusReff)
				// }
				fmlJelimCr[Jelim[1]] = rumusJelimCr
				reffJelimCr[Jelim[1]] = rumusReffJelimCr
			}
		}
		if len(Jelim) == 10 && Jelim[8] == "" && Jelim[9] == "" {
				// rumusReffJelimCr := "=" + sheetJelim + "!C" + strconv.Itoa(len(lineJelim))
				rumusJelimDr := "=" + sheetJelim + "!I" + strconv.Itoa(len(lineJelim))
				rumusJelimCr := "=" + sheetJelim + "!J" + strconv.Itoa(len(lineJelim))
				// if jpm[1] == jpm[2] {
				// 	f.SetCellFormula(sheet, fmt.Sprintf("K%d", jpmw), rumus)
				// 	f.SetCellFormula(sheet, fmt.Sprintf("J%d", jpmw), rumusReff)
				// }
				fmlJelimDr["CONTROL_TO_ADJUSTMENT_SHEET"] = rumusJelimDr
				fmlJelimCr["CONTROL_TO_ADJUSTMENT_SHEET"] = rumusJelimCr
				// reffJelimCr[Jelim[1]] = rumusReffJelimCr
			
		}
	}
	satu := 1
	dua := 2
	tiga := 3
	empat := 4
	lima := 5
	for _, v := range *data {
		rowCode[v.Code] = row
		if v.AutoSummary != nil && *v.AutoSummary {
			isAutoSum[v.Code] = true
		}
		if err = f.SetCellStyle(sheet, fmt.Sprintf("G%d", row), fmt.Sprintf("%s%d", colKonsolidasi, row), stylingCurrency); err != nil {
			errs["CONSOL"] = "Error: Generate Configuration File Excel"
			fmt.Println(err)
			return
		}
		// var codeCoa string

		if !(v.IsTotal != nil && *v.IsTotal) && v.IsLabel != nil && *v.IsLabel {
			if err = f.SetCellStyle(sheet, fmt.Sprintf("G%d", row), fmt.Sprintf("%s%d", colKonsolidasi, row), stylingCurrency); err != nil {
				errs["TB"] = "Error: Generating Configuration File Excel"
				log.Println(err)
				return
			}
			f.SetCellValue(sheet, fmt.Sprintf("C%d", row), v.Description)
			if *v.Level == satu {
				f.SetCellStyle(sheet, fmt.Sprintf("C%d", row), fmt.Sprintf("C%d", row), styleColorTextLvl1)
			}
			if *v.Level == dua {
				f.SetCellStyle(sheet, fmt.Sprintf("C%d", row), fmt.Sprintf("C%d", row), styleColorTextLvl2)
			}
			if *v.Level == tiga {
				f.SetCellStyle(sheet, fmt.Sprintf("C%d", row), fmt.Sprintf("C%d", row), styleColorTextLvl3)
			}
			if *v.Level == empat {
				f.SetCellStyle(sheet, fmt.Sprintf("C%d", row), fmt.Sprintf("C%d", row), styleColorTextLvl4)
			}
			if *v.Level == lima {
				f.SetCellStyle(sheet, fmt.Sprintf("C%d", row), fmt.Sprintf("C%d", row), styleColorTextLvl5)
			}
		}

		if v.IsCoa != nil && *v.IsCoa {
			// kasih note error pada saat import excel, seperti error pada saat import baris x, misal coa tidak terdaftar
			rowBefore := row
			consolidationDetailData, err := s.ConsolidationRepository.FindDetailByCode(ctx, &consolidationData.ID, &v.Code)
			if err != nil {
				errs["CONSOL"] = "Error: Consolidation Detail Data Not Found"
				fmt.Println(err)
				return
			}
			for _, vDetailData := range *consolidationDetailData {
				if err = f.SetCellStyle(sheet, fmt.Sprintf("G%d", row), fmt.Sprintf("%s%d", colKonsolidasi, row), stylingCurrency); err != nil {
					errs["CONSOL"] = "Error: Generate Configuration File Excel"
					fmt.Println(err)
					return
				}
				if strings.Contains(strings.ToLower(vDetailData.Code), "subtotal") {
					continue
				}
				tbRowCode[vDetailData.Code] = row
				f.SetCellValue(sheet, fmt.Sprintf("B%d", row), vDetailData.Code)
				f.SetCellValue(sheet, fmt.Sprintf("E%d", row), vDetailData.Description)
				f.SetCellValue(sheet, fmt.Sprintf("G%d", row), *vDetailData.AmountBeforeJpm)
				f.SetCellValue(sheet, fmt.Sprintf("I%d", row), *vDetailData.AmountJpmDr)
				f.SetCellValue(sheet, fmt.Sprintf("K%d", row), *vDetailData.AmountJpmCr)
				// f.SetCellFormula(sheet, fmt.Sprintf("L%d", row), fmt.Sprintf("=G%d+I%d-K%d", row, row, row))

				flt := float64(0)
				if _, ok := fmlJpmDr[vDetailData.Code]; ok {
					if *vDetailData.AmountJpmDr > flt || *vDetailData.AmountJpmDr < flt{
						f.SetCellFormula(sheet, fmt.Sprintf("H%d", row), reffJpmDr[vDetailData.Code])
						f.SetCellFormula(sheet, fmt.Sprintf("I%d", row), fmlJpmDr[vDetailData.Code])
					}
				}
				if _, ok := fmlJpmCr[vDetailData.Code]; ok {
					if *vDetailData.AmountJpmCr > flt || *vDetailData.AmountJpmCr < flt{
						f.SetCellFormula(sheet, fmt.Sprintf("J%d", row), reffJpmCr[vDetailData.Code])
						f.SetCellFormula(sheet, fmt.Sprintf("K%d", row), fmlJpmCr[vDetailData.Code])
					}
				}

				f.SetCellValue(sheet, fmt.Sprintf("N%d", row), *vDetailData.AmountJcteDr)
				f.SetCellValue(sheet, fmt.Sprintf("P%d", row), *vDetailData.AmountJcteCr)
				// f.SetCellFormula(sheet, fmt.Sprintf("Q%d", row), fmt.Sprintf("=L%d+N%d-P%d", row, row, row))

				if _, ok := fmlJcteDr[vDetailData.Code]; ok {
					if *vDetailData.AmountJcteDr > flt  || *vDetailData.AmountJcteDr < flt{
						f.SetCellFormula(sheet, fmt.Sprintf("M%d", row), reffJcteDr[vDetailData.Code])
						f.SetCellFormula(sheet, fmt.Sprintf("N%d", row), fmlJcteDr[vDetailData.Code])
					}
				}
				if _, ok := fmlJcteCr[vDetailData.Code]; ok {
					if *vDetailData.AmountJcteCr > flt || *vDetailData.AmountJcteCr < flt {
						f.SetCellFormula(sheet, fmt.Sprintf("O%d", row), reffJcteCr[vDetailData.Code])
						f.SetCellFormula(sheet, fmt.Sprintf("P%d", row), fmlJcteCr[vDetailData.Code])
					}
				}

				for i := 0; i < len(consolidationData.ConsolidationBridge); i++ {
					colCompany, err := excelize.CoordinatesToCellName(lastCol+i, 1)
					if err != nil {
						errs["CONSOL"] = "Error: Generate Configuration File Excel"
						return
					}
					colCompany = strings.ReplaceAll(colCompany, "1", "")
					tmpAmount := 0.0
					amountData, err := s.ConsolidationBridgeDetailRepository.GetWithCode(ctx, &consolidationData.ConsolidationBridge[i].ID, &vDetailData.Code)
					if err != nil && err.Error() != "record not found" {
						errs["CONSOL"] = "Error: Consolidation Bridge Detail Data Not Found"
						return
					}
					if amountData != nil && amountData.Amount != 0 {
						tmpAmount = amountData.Amount
					}
					err = f.SetCellValue(sheet, fmt.Sprintf("%s%d", colCompany, row), tmpAmount)
					if err != nil {
						errs["CONSOL"] = "Error: Generate Configuration File Excel Input Amount"
						fmt.Println(err)
						return
					}
				}

				f.SetCellFormula(sheet, fmt.Sprintf("%s%d", colCombine, row), fmt.Sprintf("=SUM(Q%d:%s%d)", row, lastCompanyColCoord, row))

				f.SetCellValue(sheet, fmt.Sprintf("%s%d", colEliminationDebit, row), *vDetailData.AmountJelimDr)
				f.SetCellValue(sheet, fmt.Sprintf("%s%d", colEliminationKredit, row), *vDetailData.AmountJelimCr)
				// f.SetCellFormula(sheet, fmt.Sprintf("%s%d", colKonsolidasi, row), fmt.Sprintf("=%s%d+%s%d-%s%d", colCombine, row, colEliminationDebit, row, colEliminationKredit, row))
				if _, ok := fmlJelimDr[vDetailData.Code]; ok {
					if *vDetailData.AmountJelimDr > flt || *vDetailData.AmountJelimDr < flt {
						f.SetCellFormula(sheet, fmt.Sprintf("%s%d", colEliminationReffDebit, row), reffJelimDr[vDetailData.Code])
						f.SetCellFormula(sheet, fmt.Sprintf("%s%d", colEliminationDebit, row), fmlJelimDr[vDetailData.Code])
					}
				}

				if _, ok := fmlJelimCr[vDetailData.Code]; ok {
					if *vDetailData.AmountJelimCr > flt || *vDetailData.AmountJelimCr < flt {
						f.SetCellFormula(sheet, fmt.Sprintf("%s%d", colEliminationReffKredit, row), reffJelimCr[vDetailData.Code])
						f.SetCellFormula(sheet, fmt.Sprintf("%s%d", colEliminationKredit, row), fmlJelimCr[vDetailData.Code])
					}
				}
				tmpHeadCoa := fmt.Sprintf("%c", vDetailData.Code[0])
				if tmpHeadCoa == "9" {
					tmpHeadCoa = vDetailData.Code[:1]
				}
				switch tmpHeadCoa {
				case "1", "5", "6", "7", "91", "92":
					f.SetCellFormula(sheet, fmt.Sprintf("L%d", row), fmt.Sprintf("=G%d+I%d-K%d", row, row, row))
					f.SetCellFormula(sheet, fmt.Sprintf("Q%d", row), fmt.Sprintf("=L%d+N%d-P%d", row, row, row))
					f.SetCellFormula(sheet, fmt.Sprintf("%s%d", colKonsolidasi, row), fmt.Sprintf("=%s%d+%s%d-%s%d", colCombine, row, colEliminationDebit, row, colEliminationKredit, row))
				default:
					f.SetCellFormula(sheet, fmt.Sprintf("L%d", row), fmt.Sprintf("=G%d-I%d+K%d", row, row, row))
					f.SetCellFormula(sheet, fmt.Sprintf("Q%d", row), fmt.Sprintf("=L%d-N%d+P%d", row, row, row))
					f.SetCellFormula(sheet, fmt.Sprintf("%s%d", colKonsolidasi, row), fmt.Sprintf("=%s%d-%s%d+%s%d", colCombine, row, colEliminationDebit, row, colEliminationKredit, row))
				}

				row++
			}

			rowAfter := row - 1
			rowConsol := len(*consolidationDetailData)
			if v.AutoSummary != nil && *v.AutoSummary && rowConsol > 1 {
				var tmp = map[string]interface{}{"code": v.Code, "row": row}
				summary = append(summary, tmp)
				f.SetCellValue(sheet, fmt.Sprintf("D%d", row), "Subtotal")
				f.SetCellFormula(sheet, fmt.Sprintf("G%d", row), fmt.Sprintf("=SUM(G%d:G%d)", rowBefore, rowAfter))
				f.SetCellFormula(sheet, fmt.Sprintf("I%d", row), fmt.Sprintf("=SUM(I%d:I%d)", rowBefore, rowAfter))
				f.SetCellFormula(sheet, fmt.Sprintf("K%d", row), fmt.Sprintf("=SUM(K%d:K%d)", rowBefore, rowAfter))
				f.SetCellFormula(sheet, fmt.Sprintf("L%d", row), fmt.Sprintf("=SUM(L%d:L%d)", rowBefore, rowAfter))
				f.SetCellFormula(sheet, fmt.Sprintf("N%d", row), fmt.Sprintf("=SUM(N%d:N%d)", rowBefore, rowAfter))
				f.SetCellFormula(sheet, fmt.Sprintf("P%d", row), fmt.Sprintf("=SUM(P%d:P%d)", rowBefore, rowAfter))
				f.SetCellFormula(sheet, fmt.Sprintf("Q%d", row), fmt.Sprintf("=SUM(Q%d:Q%d)", rowBefore, rowAfter))
				// f.SetCellFormula(sheet, fmt.Sprintf("S%d", row), fmt.Sprintf("=SUM(S%d:S%d)", rowBefore, rowAfter))
				for i := 0; i < len(consolidationData.ConsolidationBridge); i++ {
					colCompany, err := excelize.CoordinatesToCellName(lastCol+i, 1)
					if err != nil {
						errs["CONSOL"] = "Error: Generate Configuration File Excel"
						return
					}
					colCompany = strings.ReplaceAll(colCompany, "1", "")
					err = f.SetCellFormula(sheet, fmt.Sprintf("%s%d", colCompany, row), fmt.Sprintf("=SUM(%s%d:%s%d)", colCompany, rowBefore, colCompany, rowAfter))
					if err != nil {
						errs["CONSOL"] = "Error: Generate Configuration File Excel on Set Formula"
						fmt.Println(err)
						return
					}
					f.SetCellStyle(sheet, fmt.Sprintf("%s%d", colCompany, row), fmt.Sprintf("%s%d", colCompany, row), stylingSubTotalCurrency)
				}

				if err = f.SetCellStyle(sheet, fmt.Sprintf("G%d", row), fmt.Sprintf("%s%d", colKonsolidasi, row), stylingCurrency); err != nil {
					errs["CONSOL"] = "Error: Generate Configuration File Excel"
					fmt.Println(err)
					return
				}

				err = f.SetCellFormula(sheet, fmt.Sprintf("%s%d", colCombine, row), fmt.Sprintf("=SUM(%s%d:%s%d)", colCombine, rowBefore, colCombine, rowAfter))
				if err != nil {
					errs["CONSOL"] = "Error: Generate Configuration File Excel on Set Formula"
					fmt.Println(err)
					return
				}

				err = f.SetCellFormula(sheet, fmt.Sprintf("%s%d", colEliminationDebit, row), fmt.Sprintf("=SUM(%s%d:%s%d)", colEliminationDebit, rowBefore, colEliminationDebit, rowAfter))
				if err != nil {
					errs["CONSOL"] = "Error: Generate Configuration File Excel on Set Formula"
					fmt.Println(err)
					return
				}

				err = f.SetCellFormula(sheet, fmt.Sprintf("%s%d", colEliminationKredit, row), fmt.Sprintf("=SUM(%s%d:%s%d)", colEliminationKredit, rowBefore, colEliminationKredit, rowAfter))
				if err != nil {
					errs["CONSOL"] = "Error: Generate Configuration File Excel on Set Formula"
					fmt.Println(err)
					return
				}

				err = f.SetCellFormula(sheet, fmt.Sprintf("%s%d", colKonsolidasi, row), fmt.Sprintf("=SUM(%s%d:%s%d)", colKonsolidasi, rowBefore, colKonsolidasi, rowAfter))
				if err != nil {
					errs["CONSOL"] = "Error: Generate Configuration File Excel on Set Formula"
					fmt.Println(err)
					return
				}
				// f.SetCellStyle(sheet, fmt.Sprintf("%s%d", colCombine, row), fmt.Sprintf("%s%d", colCombine, row), stylingSubTotalCurrency)
				// f.SetCellStyle(sheet, fmt.Sprintf("%s%d", colEliminationDebit, row), fmt.Sprintf("%s%d", colEliminationDebit, row), stylingSubTotalCurrency)
				// f.SetCellStyle(sheet, fmt.Sprintf("%s%d", colEliminationKredit, row), fmt.Sprintf("%s%d", colEliminationKredit, row), stylingSubTotalCurrency)
				// f.SetCellStyle(sheet, fmt.Sprintf("%s%d", colKonsolidasi, row), fmt.Sprintf("%s%d", colKonsolidasi, row), stylingSubTotalCurrency)
				if err = f.SetCellStyle(sheet, fmt.Sprintf("G%d", row), fmt.Sprintf("%s%d", colKonsolidasi, row), stylingSubTotalCurrency); err != nil {
					errs["CONSOL"] = "Error: Generate Configuration File Excel"
					fmt.Println(err)
					return
				}

				f.SetCellStyle(sheet, fmt.Sprintf("G%d", row), fmt.Sprintf("Q%d", row), stylingSubTotalCurrency)
				f.SetCellStyle(sheet, fmt.Sprintf("D%d", row), fmt.Sprintf("F%d", row), stylingSubTotal)
				rowCode[fmt.Sprintf("%s_SUBTOTAL", v.Code)] = row
				row++
				if err = f.SetCellStyle(sheet, fmt.Sprintf("G%d", row), fmt.Sprintf("%s%d", colKonsolidasi, row), stylingCurrency); err != nil {
					errs["CONSOL"] = "Error: Generate Configuration File Excel"
					fmt.Println(err)
					return
				}
			}
		}

		if v.IsTotal != nil && *v.IsTotal {
			tbRowCode[v.Code] = row
			f.SetCellValue(sheet, fmt.Sprintf("C%d", row), v.Description)
			if v.Code != "TOTAL_JOURNAL_IN_WP" && v.Code != "CONTROL_TO_ADJUSTMENT_SHEET" && v.Code != "CONTROL" && v.Code != "CONTROL_TO_WBS_1" {
				if v.IsControl != nil && *v.IsControl {
					f.SetCellStyle(sheet, fmt.Sprintf("B%d", row), fmt.Sprintf("%s%d", colKonsolidasi, row), stylingTotalControl)
					if err = f.SetCellStyle(sheet, fmt.Sprintf("G%d", row), fmt.Sprintf("%s%d", colKonsolidasi, row), stylingTotalControl); err != nil {
						errs["CONSOL"] = "Error: Generate Configuration File Excel"
						fmt.Println(err)
						return
					}
					if *v.Level == satu {
						f.SetCellStyle(sheet, fmt.Sprintf("B%d", row), fmt.Sprintf("%s%d", colKonsolidasi, row), stylingTotalLvl1)
					}
					if *v.Level == dua {
						f.SetCellStyle(sheet, fmt.Sprintf("B%d", row), fmt.Sprintf("%s%d", colKonsolidasi, row), stylingTotalLvl2)
					}
					if *v.Level == tiga {
						f.SetCellStyle(sheet, fmt.Sprintf("B%d", row), fmt.Sprintf("%s%d", colKonsolidasi, row), stylingTotalLvl3)
					}
					if *v.Level == empat {
						f.SetCellStyle(sheet, fmt.Sprintf("B%d", row), fmt.Sprintf("%s%d", colKonsolidasi, row), stylingTotalLvl4)
					}
					if *v.Level == lima {
						f.SetCellStyle(sheet, fmt.Sprintf("B%d", row), fmt.Sprintf("%s%d", colKonsolidasi, row), stylingTotalLvl5)
					}
				} else {
					f.SetCellStyle(sheet, fmt.Sprintf("B%d", row), fmt.Sprintf("%s%d", colKonsolidasi, row), stylingTotal)
					if err = f.SetCellStyle(sheet, fmt.Sprintf("G%d", row), fmt.Sprintf("%s%d", colKonsolidasi, row), stylingTotalCurrency); err != nil {
						errs["CONSOL"] = "Error: Generate Configuration File Excel"
						fmt.Println(err)
						return
					}
					if *v.Level == satu {
						f.SetCellStyle(sheet, fmt.Sprintf("B%d", row), fmt.Sprintf("%s%d", colKonsolidasi, row), stylingTotalLvl1)
					}
					if *v.Level == dua {
						f.SetCellStyle(sheet, fmt.Sprintf("B%d", row), fmt.Sprintf("%s%d", colKonsolidasi, row), stylingTotalLvl2)
					}
					if *v.Level == tiga {
						f.SetCellStyle(sheet, fmt.Sprintf("B%d", row), fmt.Sprintf("%s%d", colKonsolidasi, row), stylingTotalLvl3)
					}
					if *v.Level == empat {
						f.SetCellStyle(sheet, fmt.Sprintf("B%d", row), fmt.Sprintf("%s%d", colKonsolidasi, row), stylingTotalLvl4)
					}
					if *v.Level == lima {
						f.SetCellStyle(sheet, fmt.Sprintf("B%d", row), fmt.Sprintf("%s%d", colKonsolidasi, row), stylingTotalLvl5)
					}
				}
			}
			if v.Code == "TOTAL_LIABILITAS_DAN_EKUITAS" {
				rowAset := row
				if _, ok := rowCode["TOTAL_ASET"]; ok {
					rowAset = rowCode["TOTAL_ASET"]
				}
				f.SetCellFormula(sheet, "G5", fmt.Sprintf("=G%d-G%d", rowAset, row))
				// f.SetCellFormula(sheet, "I5", fmt.Sprintf("=I%d-I%d", rowAset, row))
				// f.SetCellFormula(sheet, "K5", fmt.Sprintf("=K%d-K%d", rowAset, row))
				f.SetCellFormula(sheet, "L5", fmt.Sprintf("=L%d-L%d", rowAset, row))
				// f.SetCellFormula(sheet, "N5", fmt.Sprintf("=N%d-N%d", rowAset, row))
				// f.SetCellFormula(sheet, "P5", fmt.Sprintf("=P%d-P%d", rowAset, row))
				f.SetCellFormula(sheet, "Q5", fmt.Sprintf("=Q%d-Q%d", rowAset, row))
				// tmpRow := ""
				for i := lastCol; i < lastCol+len(consolidationData.ConsolidationBridge)+6; i++ {
					if i == companyCol+1 || i == companyCol+3 {
						continue
					}
					colCompany, err := excelize.CoordinatesToCellName(i, 1)
					if err != nil {
						errs["CONSOL"] = "Error: Generate Configuration File Excel"
						return
					}
					colCompany = strings.ReplaceAll(colCompany, "1", "")
					f.SetCellFormula(sheet, fmt.Sprintf("%s5", colCompany), fmt.Sprintf("=%s%d-%s%d", colCompany, rowAset, colCompany, row))
					// tmpRow = colCompany
				}
				// f.SetCellFormula(sheet, fmt.Sprintf("%s%d", colCombine, row), fmt.Sprintf("=SUM(%s%d:%s%d)", colCombine, rowBefore, colCombine, rowAfter))
				// f.SetCellFormula(sheet, fmt.Sprintf("%s%d", colEliminationDebit, row), fmt.Sprintf("=SUM(%s%d:%s%d)", colEliminationDebit, rowBefore, colEliminationDebit, rowAfter))
				// f.SetCellFormula(sheet, fmt.Sprintf("%s%d", colEliminationKredit, row), fmt.Sprintf("=SUM(%s%d:%s%d)", colEliminationKredit, rowBefore, colEliminationKredit, rowAfter))
				// f.SetCellFormula(sheet, fmt.Sprintf("%s%d", colKonsolidasi, row), fmt.Sprintf("=SUM(%s%d:%s%d)", colKonsolidasi, rowBefore, colKonsolidasi, rowAfter))
			}

			//show control aje
			if v.Code == "CONTROL_TO_ADJUSTMENT_SHEET" {
				f.SetCellValue(sheet, fmt.Sprintf("C%d", row), v.Description)
				dbt := 0.0
				if summaryJPM.IncomeStatementDr != nil && *summaryJPM.IncomeStatementDr != 0 {
					dbt = *summaryJPM.IncomeStatementDr
				}
				if summaryJPM.BalanceSheetDr != nil && *summaryJPM.BalanceSheetDr != 0 {
					dbt += *summaryJPM.BalanceSheetDr
				}
				cdt := 0.0
				if summaryJPM.IncomeStatementCr != nil && *summaryJPM.IncomeStatementCr != 0 {
					cdt = *summaryJPM.IncomeStatementCr
				}
				if summaryJPM.BalanceSheetCr != nil && *summaryJPM.BalanceSheetCr != 0 {
					cdt += *summaryJPM.BalanceSheetCr
				}
				f.SetCellValue(sheet, fmt.Sprintf("I%d", row), dbt)
				f.SetCellValue(sheet, fmt.Sprintf("K%d", row), cdt)
				f.SetCellFormula(sheet, fmt.Sprintf("I%d", row), fmlJpmDr["CONTROL_TO_ADJUSTMENT_SHEET"])
				f.SetCellFormula(sheet, fmt.Sprintf("K%d", row), fmlJpmCr["CONTROL_TO_ADJUSTMENT_SHEET"])
				f.SetCellFormula(sheet, fmt.Sprintf("J%d", row), fmt.Sprintf("=I%d-K%d", row, row))

				dbt2 := 0.0
				if summaryJCTE.IncomeStatementDr != nil && *summaryJCTE.IncomeStatementDr != 0 {
					dbt2 = *summaryJCTE.IncomeStatementDr
				}
				if summaryJCTE.BalanceSheetDr != nil && *summaryJCTE.BalanceSheetDr != 0 {
					dbt2 += *summaryJCTE.BalanceSheetDr
				}
				cdt2 := 0.0
				if summaryJCTE.IncomeStatementCr != nil && *summaryJCTE.IncomeStatementCr != 0 {
					cdt2 = *summaryJCTE.IncomeStatementCr
				}
				if summaryJCTE.BalanceSheetCr != nil && *summaryJCTE.BalanceSheetCr != 0 {
					cdt2 += *summaryJCTE.BalanceSheetCr
				}
				f.SetCellValue(sheet, fmt.Sprintf("N%d", row), dbt2)
				f.SetCellValue(sheet, fmt.Sprintf("P%d", row), cdt2)
				f.SetCellFormula(sheet, fmt.Sprintf("N%d", row), fmlJcteDr["CONTROL_TO_ADJUSTMENT_SHEET"])
				f.SetCellFormula(sheet, fmt.Sprintf("P%d", row), fmlJcteCr["CONTROL_TO_ADJUSTMENT_SHEET"])
				f.SetCellFormula(sheet, fmt.Sprintf("O%d", row), fmt.Sprintf("=N%d-P%d", row, row))

				dbt3 := 0.0
				if summaryJELIM.IncomeStatementDr != nil && *summaryJELIM.IncomeStatementDr != 0 {
					dbt3 = *summaryJELIM.IncomeStatementDr
				}
				if summaryJELIM.BalanceSheetDr != nil && *summaryJELIM.BalanceSheetDr != 0 {
					dbt3 += *summaryJELIM.BalanceSheetDr
				}
				cdt3 := 0.0
				if summaryJELIM.IncomeStatementCr != nil && *summaryJELIM.IncomeStatementCr != 0 {
					cdt3 = *summaryJELIM.IncomeStatementCr
				}
				if summaryJELIM.BalanceSheetCr != nil && *summaryJELIM.BalanceSheetCr != 0 {
					cdt3 += *summaryJELIM.BalanceSheetCr
				}
				colJelimDb, err := excelize.CoordinatesToCellName(companyCol+2, row)
				if err != nil {
					errs["CONSOL"] = "Error: Generate Configuration File Excel Setting Formula"
					return
				}
				colJelimCr, err := excelize.CoordinatesToCellName(companyCol+4, row)
				if err != nil {
					errs["CONSOL"] = "Error: Generate Configuration File Excel Setting Formula"
					return
				}
				colMidJelim, err := excelize.CoordinatesToCellName(companyCol+3, row)
				if err != nil {
					errs["CONSOL"] = "Error: Generate Configuration File Excel Setting Formula"
					return
				}
				// if _, ok := fmlJelimDr[v.Code]; ok {
				// 		f.SetCellFormula(sheet, fmt.Sprintf("%s%d", colJelimDb, row), fmlJelimDr[v.Code])
				// 		f.SetCellFormula(sheet, fmt.Sprintf("%s%d", colJelimCr, row), fmlJelimDr[v.Code])
				// }
				// f.SetCellValue(sheet, colJelimDb, dbt3)
				// f.SetCellValue(sheet, colJelimCr, cdt3)
				f.SetCellFormula(sheet, colJelimDb, fmlJelimDr["CONTROL_TO_ADJUSTMENT_SHEET"])
				f.SetCellFormula(sheet, colJelimCr, fmlJelimCr["CONTROL_TO_ADJUSTMENT_SHEET"])
				f.SetCellFormula(sheet, colMidJelim, fmt.Sprintf("=%s-%s", colJelimDb, colJelimCr))
			}
			if v.Code == "CONTROL" {
				f.SetCellFormula(sheet, "I5", fmt.Sprintf("=I%d", row))
				f.SetCellFormula(sheet, "K5", fmt.Sprintf("=K%d", row))
				f.SetCellFormula(sheet, "N5", fmt.Sprintf("=N%d", row))
				f.SetCellFormula(sheet, "P5", fmt.Sprintf("=P%d", row))
				// tmpRow := ""
				tmpl := len(consolidationData.ConsolidationBridge)
				for i := tmpl + 2; i < tmpl+5; i++ {
					if i == tmpl+3 {
						continue
					}
					colCompany, err := excelize.CoordinatesToCellName(lastCol+i, 1)
					if err != nil {
						errs["CONSOL"] = "Error: Generate Configuration File Excel"
						return
					}
					colCompany = strings.ReplaceAll(colCompany, "1", "")
					f.SetCellFormula(sheet, fmt.Sprintf("%s5", colCompany), fmt.Sprintf("=%s%d", colCompany, row))
					// tmpRow = colCompany
				}
			}
			if v.Code == "TOTAL_JOURNAL_IN_WP" {
				f.SetCellFormula(sheet, fmt.Sprintf("J%d", row), fmt.Sprintf("=I%d-K%d", row, row))
				f.SetCellFormula(sheet, fmt.Sprintf("O%d", row), fmt.Sprintf("=N%d-P%d", row, row))
				colMidJelim, err := excelize.CoordinatesToCellName(companyCol+3, row)
				if err != nil {
					errs["CONSOL"] = "Error: Generate Configuration File Excel Setting Formula"
					return
				}
				tmpMidJelim := strings.ReplaceAll(colMidJelim, fmt.Sprintf("%d", row), "")
				f.SetCellFormula(sheet, colMidJelim, fmt.Sprintf("=%s%d-%s%d", tmpMidJelim, row, tmpMidJelim, row))
			}

			if v.FxSummary == "" {
				row++
				continue
			}
			for tmpChr := 7; tmpChr <= (companyCol + 5); tmpChr++ {
				colKonsolidasi, err := excelize.CoordinatesToCellName(tmpChr, 1)
				if err != nil {
					errs["CONSOL"] = "Error: Generate Configuration File Excel Setting Formula"
					return
				}
				chr := strings.ReplaceAll(colKonsolidasi, "1", "")
				if chr == "H" || chr == "J" || chr == "M" || chr == "O" || tmpChr == companyCol+1 || tmpChr == companyCol+3 || ((chr == "G" || chr == "L" || chr == "Q" || (tmpChr >= lastCol && tmpChr <= companyCol) || tmpChr == companyCol+5) && (v.Code == "TOTAL_JOURNAL_IN_WP" || v.Code == "CONTROL_TO_ADJUSTMENT_SHEET")) || (v.Code == "CONTROL_TO_WBS_1" && (chr == "I" || chr == "K" || chr == "N" || chr == "P" || tmpChr == companyCol+2 || tmpChr == companyCol+4)) {
					continue
				}
				formula := v.FxSummary
				reg := regexp.MustCompile(`[A-Za-z0-9_~:()]+|[0-9]+\d{2,}`)
				match := reg.FindAllString(formula, -1)
				for _, vMatch := range match {
					if len(vMatch) < 3 {
						continue
					}
					//cari jml berdasarkan code
					if isAutoSum[vMatch] {
						if rowCode[fmt.Sprintf("%s_SUBTOTAL", vMatch)] != 0 {
							formula = strings.ReplaceAll(formula, vMatch, fmt.Sprintf("%s%d", chr, rowCode[fmt.Sprintf("%s_SUBTOTAL", vMatch)]))
						} else {
							formula = strings.ReplaceAll(formula, vMatch, fmt.Sprintf("%s%d", chr, rowCode[vMatch]))
						}
					} else {
						if rowCode[vMatch] != 0 {
							formula = strings.ReplaceAll(formula, vMatch, fmt.Sprintf("%s%d", chr, rowCode[vMatch]))
						}
					}
				}
				f.SetCellFormula(sheet, fmt.Sprintf("%s%d", chr, row), fmt.Sprintf("=%s", formula))
				// if fmt.Sprintf("%c", chr) == "I" || fmt.Sprintf("%c", chr) == "K" || fmt.Sprintf("%c", chr) == "N" || fmt.Sprintf("%c", chr) == "P" || fmt.Sprintf("%c", chr) == "U" || fmt.Sprintf("%c", chr) == "W" {
				// 	newStr := strings.ReplaceAll(formula, "-", "+")
				// 	f.SetCellFormula(sheet, fmt.Sprintf("%c%d", chr, row), fmt.Sprintf("=%s", newStr))
				// }
				if chr == "I" || chr == "K" || chr == "N" || chr == "P"  || chr == colEliminationDebit || chr == colEliminationKredit {
					if v.Code == "CONTROL" {
						continue
					}
					newStr := strings.ReplaceAll(formula, "-", "+")
					f.SetCellFormula(sheet, fmt.Sprintf("%s%d", chr, row), fmt.Sprintf("=%s", newStr))
				}
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
			f.SetCellFormula(sheet, fmt.Sprintf("N%d", val), strings.ReplaceAll(vCustomRow, "@", "N"))
			f.SetCellFormula(sheet, fmt.Sprintf("P%d", val), strings.ReplaceAll(vCustomRow, "@", "P"))
			// f.SetCellFormula(sheet, fmt.Sprintf("R%d", val), strings.ReplaceAll(vCustomRow, "@", "R"))
			// f.SetCellFormula(sheet, fmt.Sprintf("U%d", val), strings.ReplaceAll(vCustomRow, "@", "U"))
			// f.SetCellFormula(sheet, fmt.Sprintf("W%d", val), strings.ReplaceAll(vCustomRow, "@", "W"))
			// f.SetCellFormula(sheet, fmt.Sprintf("V%d", val), strings.ReplaceAll(vCustomRow, "@", "V"))
			for i := 0; i < len(consolidationData.ConsolidationBridge); i++ {
				colCompany, err := excelize.CoordinatesToCellName(lastCol+i, 1)
				if err != nil {
					errs["CONSOL"] = "Error: Generate Configuration File Excel"
					return
				}
				colCompany = strings.ReplaceAll(colCompany, "1", "")
				// tmpAmount := 0.0
				// amountData, err := s.ConsolidationBridgeDetailRepository.GetWithCode(ctx, &consolidationData.ConsolidationBridge[i].ID, &vDetailData.Code)
				// if err != nil && err.Error() != "record not found" {
				// 	errs["CONSOL"] = "Error: Consolidation Bridge Detail Data Not Found"
				// 	return
				// }
				// if amountData != nil && amountData.Amount != 0 {
				// 	tmpAmount = amountData.Amount
				// }
				// err = f.SetCellValue(sheet, fmt.Sprintf("%s%d", colCompany, row), tmpAmount)
				f.SetCellFormula(sheet, fmt.Sprintf("%s%d",colCompany, val), strings.ReplaceAll(vCustomRow, "@", colCompany))
				if err != nil {
					errs["CONSOL"] = "Error: Generate Configuration File Excel Input Amount"
					fmt.Println(err)
					return
				}
			}
			f.SetCellFormula(sheet, fmt.Sprintf("%s%d",colEliminationDebit, val), strings.ReplaceAll(vCustomRow, "@", colEliminationDebit))
			f.SetCellFormula(sheet, fmt.Sprintf("%s%d",colEliminationKredit, val), strings.ReplaceAll(vCustomRow, "@", colEliminationKredit))
			
			
		}
	}

	if err = f.SetCellStyle(sheet, fmt.Sprintf("B%d", row), fmt.Sprintf("%s%d", colKonsolidasi, row), stylingBorderTopOnly); err != nil {
		errs["CONSOL"] = "Error: Generate Configuration File Excel"
		fmt.Println(err)
		return
	}

	if err = f.SetCellStyle(sheet, "A9", fmt.Sprintf("A%d", row-1), stylingBorderRightOnly); err != nil {
		errs["CONSOL"] = "Error: Generate Configuration File Excel"
		fmt.Println(err)
		return
	}

	if err = f.SetSheetFormatPr(sheet, excelize.DefaultRowHeight(12.85)); err != nil {
		errs["CONSOL"] = "Error: Generate Configuration File Excel"
		fmt.Println(err)
		return
	}
	f.SetDefaultFont("Arial")
	wg.Wait()
	date := time.Now().Format("20060102_150405")
	period := datePeriod.Format("2006-01-02")
	saveAsset := fmt.Sprintf("assets/%d/%s", payload.UserID, date)
	tmpFolder := path.Join(configs.App().StoragePath(), saveAsset)
	// tmpFolder := fmt.Sprintf("/%d/%s", payload.UserID, date)
	_, err = os.Stat(tmpFolder)
	if err != nil {
		if os.IsNotExist(err) {
			if err := os.MkdirAll(tmpFolder, os.ModePerm); err != nil {
				log.Fatal(err)
			}
		} else {
			fmt.Println(err)
		}
	}
	loc := fmt.Sprintf("%s/Konsolidasi_%s.xlsx", tmpFolder, period)
	err = f.SaveAs(loc)
	if err != nil {
		errs["CONSOL"] = "Error: Saving File Excel"
		fmt.Println(err)
		return
	}
	f, err = excelize.OpenFile(loc)
	if err != nil {
		return
	}

	rows, err := f.GetRows("Consolidation")
	if err != nil {
		fmt.Println("Error reading rows:", err)
		return
	}
	// startRow := 9
	// endRow := 15
	
	// fa := false
	
	var criteriaFormatterGrouping model.FormatterDetailFilterModel
	criteriaFormatterGrouping.FormatterID = &formatterID
	criteriaFormatterGrouping.IsLabel = &t
	// criteriaFormatterGrouping.IsTotal = &fa
	pagesize := 100000
	tmpStr := "sort_id"
	tmpStr1 := "ASC"
	paginationTB := abstraction.Pagination{
		PageSize: &pagesize,
		SortBy:   &tmpStr,
		Sort:     &tmpStr1,
	}
	dataGrouping, _, err := s.FormatterDetailRepository.FindGroup(ctx, &criteriaFormatterGrouping, &paginationTB)
	if err != nil {
		return
	}
	

	lineCode := make(map[string]int)
	for i, row := range rows {
		
		if len(row) < 2 {
			cellValue := "test"
			lineCode[cellValue] = i
		}else {
			cellValue := row[2]
			lineCode[cellValue] = i
		}
		i++
		// Simpan nomor baris dalam map dengan kunci (key) nilai sel
	}
	for i, d := range *dataGrouping {

		if *d.IsTotal == t {
			continue
		}
		if _,ok := lineCode[d.Description]; ok{
		
			startRow := lineCode[d.Description]
			startRow = startRow + 2

			if i+1 < len(*dataGrouping) {
				secondElement := (*dataGrouping)[i+1]
				endRow := lineCode[secondElement.Description]
				endRow = endRow - 1

				for row := startRow; row <= endRow; row++ {
					if err := f.SetRowOutlineLevel("Consolidation", row, 1); err != nil {
						return
					}
				}
			}
		}
	}
	// Simpan perubahan ke dalam file Excel
	err = f.SaveAs(loc)
	if err != nil {
		errs["CONSOL"] = "Error: Saving File Excel"
		fmt.Println(err)
		return
	}
	f.Close()
	notifData := model.NotificationEntityModel{}
	notifData.Context = ctx
	waktu := time.Now()
	map1 := kafkaproducer.JsonData{
		FileLoc:   loc,
		UserID:    ctx.Auth.ID,
		CompanyID: ctx.Auth.CompanyID,
		Name:      "export",
		Timestamp: &waktu,
	}

	map1.Data = "Proses Export Konsolidasi Berhasil!"
	if len(errs) != 0 {
		log.Print(errs)
		map1.Data = "Proses Export Konsolidasi Gagal!"
	}

	datas := JsonData{
		CompanyName: consolidationData.Company.Name,
		Period:      payload.Filter.Period,
		Versions:    payload.Filter.Versions,
		File:        loc,
		Type:        "export",
	}

	checkErr := 0
	for _, err := range errs {
		if err != "" {
			checkErr++
		}
	}
	if checkErr > 0 {
		notifData.Description = "Proses Export Gagal!"
		tmpMap := []string{}
		for _, err := range errs {
			if err != "" {
				tmpMap = append(tmpMap, err)
			}
		}
		jsonErr, err := json.Marshal(tmpMap)
		if err != nil {
			log.Println(err)
			return
		}
		datas.Errors = string(jsonErr)
		// go kafkaproducer.NewProducer("NOTIFICATION").SendMessage("NOTIFICATION", string(jsonStr))
		// return
	}
	jsonData, err := json.Marshal(datas)
	if err != nil {
		log.Println(err)
		return
	}
	notifData.Description = map1.Data
	tmpfalse := false
	notifData.IsOpened = &tmpfalse
	notifData.CreatedBy = ctx.Auth.ID
	notifData.CreatedAt = *utilDate.DateTodayLocal()
	notifData.Data = string(jsonData)

	notificationData, err := s.NotificationRepository.Create(ctx, &notifData)
	if err != nil {
		fmt.Println(err)
	}
	map1.ID = notificationData.ID
	jsonStr, err := json.Marshal(map1)
	if err != nil {
		fmt.Println(err)
	}

	go kafkaproducer.NewProducer("NOTIFICATION").SendMessage("NOTIFICATION", string(jsonStr))
}

func (s *service) ExportAgingUtangPiutang(ctx *abstraction.Context, payload *abstraction.JsonData, wg *sync.WaitGroup, tmpFolder string) {
	defer wg.Done()
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
		errs["AGING"] = "Error: Generate Configuration File Excel"
		log.Println(err)
		return
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
		errs["AGING"] = "Error: Generate Configuration File Excel"
		log.Println(err)
		return
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
		errs["AGING"] = "Error: Generate Configuration File Excel"
		log.Println(err)
		return
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
		errs["AGING"] = "Error: Generate Configuration File Excel"
		log.Println(err)
		return
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
		errs["AGING"] = "Error: Generate Configuration File Excel"
		log.Println(err)
		return
	}

	criteria := model.AgingUtangPiutangFilterModel{}
	criteria.CompanyID = &payload.CompanyID
	criteria.Versions = &payload.Filter.Versions
	criteria.Period = &payload.Filter.Period
	// criteria.FormatterID = &data.ID

	agingutangpiutang, err := s.AgingUtangPiutangRepository.Find(ctx, &criteria)
	if err != nil {
		errs["AGING"] = "Error: Data Aging Utang Piutang Not Found"
		log.Println(err)
		return
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
		f.SetCellValue(sheet, fmt.Sprintf("J%d", (rowStart-2)), "Piutang pihak berelasi jangka panjang (8)")
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

		var criteria model.FormatterFilterModel
		criteria.FormatterFor = &formatter

		data, err := s.FormatterRepository.FindWithDetail(ctx, &criteria)
		if err != nil {
			errs["AGING"] = "Error: Formatter Data for Aging Utang Piutang Not Found"
			log.Println(err)
			return
		}

		tmpStr := "AGING-UTANG-PIUTANG"
		criteriaBridge := model.FormatterBridgesFilterModel{}
		criteriaBridge.FormatterBridgesFilter.Source = &tmpStr
		criteriaBridge.FormatterBridgesFilter.FormatterID = &data.ID
		for _, valaup := range *agingutangpiutang {
			criteriaBridge.FormatterBridgesFilter.TrxRefID = &valaup.ID
		}

		bridges, err := s.FormatterBridgesRepository.FindWithCriteria(ctx, &criteriaBridge)
		if err != nil {
			errs["AGING"] = "Error: Data Aging Utang Piutang Not Found"
			log.Println(err)
			return
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
			for _, valaup := range *agingutangpiutang {
				criteriaAUP.AgingUtangPiutangID = &valaup.ID
			}

			paginationAUP := abstraction.Pagination{}
			pagesize := 10000
			// sortBy := "id"
			// sort := "asc"
			paginationAUP.PageSize = &pagesize
			// paginationAUP.SortBy = &sortBy
			// paginationAUP.Sort = &sort

			AgingUPDetail, err := s.AgingUPDetailRepository.Find(ctx, &criteriaAUP)
			if err != nil {
				errs["AGING"] = "Error: Data Aging Utang Piutang Not Found"
				log.Println(err)
				return
			}
			if len(*AgingUPDetail) == 0 {
				errs["AGING"] = "Error: Data Aging Utang Piutang Not Found"
				log.Println("Data Not Found")
				return
			}
			for _, vv := range *AgingUPDetail {
				f.SetCellValue(sheet, fmt.Sprintf("C%d", row), *vv.Piutangusaha3rdparty)
				f.SetCellValue(sheet, fmt.Sprintf("D%d", row), *vv.PiutangusahaBerelasi)
				f.SetCellValue(sheet, fmt.Sprintf("E%d", row), *vv.Piutanglainshortterm3rdparty)
				f.SetCellValue(sheet, fmt.Sprintf("F%d", row), *vv.PiutanglainshorttermBerelasi)
				f.SetCellValue(sheet, fmt.Sprintf("G%d", row), *vv.Piutangberelasishortterm)
				f.SetCellValue(sheet, fmt.Sprintf("H%d", row), *vv.Piutanglainlongterm3rdparty)
				f.SetCellValue(sheet, fmt.Sprintf("I%d", row), *vv.PiutanglainlongtermBerelasi)
				f.SetCellValue(sheet, fmt.Sprintf("J%d", row), *vv.Piutangberelasilongterm)
				f.SetCellValue(sheet, fmt.Sprintf("L%d", row), *vv.Utangusaha3rdparty)
				f.SetCellValue(sheet, fmt.Sprintf("M%d", row), *vv.UtangusahaBerelasi)
				f.SetCellValue(sheet, fmt.Sprintf("N%d", row), *vv.Utanglainshortterm3rdparty)
				f.SetCellValue(sheet, fmt.Sprintf("O%d", row), *vv.UtanglainshorttermBerelasi)
				f.SetCellValue(sheet, fmt.Sprintf("P%d", row), *vv.Utangberelasishortterm)
				f.SetCellValue(sheet, fmt.Sprintf("Q%d", row), *vv.Utanglainlongterm3rdparty)
				f.SetCellValue(sheet, fmt.Sprintf("R%d", row), *vv.UtanglainlongtermBerelasi)
				f.SetCellValue(sheet, fmt.Sprintf("S%d", row), *vv.Utangberelasilongterm)
			}
			row++
		}
		rowStart = row + 5
		row = rowStart
	}
	datePeriod, err := time.Parse("2006-01-02", payload.Filter.Period)
	if err != nil {
		errs["AGING"] = "Error: Invalid Date Period"
		log.Println(err)
		return
	}
	period := datePeriod.Format("2006-01-02")
	err = f.SaveAs(fmt.Sprintf("%s/AgingUtangPiutang_%s.xlsx", tmpFolder, period))
	if err != nil {
		errs["AGING"] = "Error: Saving File Excel"
		log.Println(err)
		return
	}
	f.Close()

}

func (s *service) ExportInvestasiNonTbk(ctx *abstraction.Context, payload *abstraction.JsonData, wg *sync.WaitGroup, tmpFolder string) {
	defer wg.Done()
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
			errs["NONTBK"] = "Error: Generating Configuration File Excel"
			log.Println(err)
			return
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
		errs["NONTBK"] = "Error: Generating Configuration File Excel"
		log.Println(err)
		return
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
		errs["NONTBK"] = "Error: Generating Configuration File Excel"
		log.Println(err)
		return
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
		errs["NONTBK"] = "Error: Generating Configuration File Excel"
		log.Println(err)
		return
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
		errs["NONTBK"] = "Error: Generating Configuration File Excel"
		log.Println(err)
		return
	}

	row, rowStart := 5, 5

	criteria := model.InvestasiNonTbkFilterModel{}
	criteria.CompanyID = &payload.CompanyID
	criteria.Versions = &payload.Filter.Versions
	criteria.Period = &payload.Filter.Period

	investasinontbkDatas, err := s.InvestasiNonTbkRepository.Find(ctx, &criteria)
	if err != nil {
		errs["NONTBK"] = "Error: Data Investasi Non TBK not found"
		log.Println(err)
		return
	}

	criteriaDetail := model.InvestasiNonTbkDetailFilterModel{}
	investasinontbk := model.InvestasiNonTbkEntityModel{}
	for _, valINT := range *investasinontbkDatas {
		criteriaDetail.InvestasiNonTbkID = &valINT.ID
		investasinontbk = valINT
	}

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

	detail, err := s.InvestasiNonTbkDetailRepository.Find(ctx, &criteriaDetail)
	if err != nil {
		errs["NONTBK"] = "Error: Formatter Data for Investasi Non TBK not found"
		log.Println(err)
		return
	}

	for _, v := range *detail {
		f.SetCellStyle(sheet, fmt.Sprintf("B%d", row), fmt.Sprintf("C%d", row), styleLabel)
		f.SetCellStyle(sheet, fmt.Sprintf("D%d", row), fmt.Sprintf("E%d", row), styleCurrency)
		f.SetCellStyle(sheet, fmt.Sprintf("F%d", row), fmt.Sprintf("F%d", row), styleCurrencyPercentage)
		f.SetCellStyle(sheet, fmt.Sprintf("G%d", row), fmt.Sprintf("J%d", row), styleCurrency)

		f.SetCellValue(sheet, fmt.Sprintf("B%d", row), v.Code)

		f.SetCellValue(sheet, fmt.Sprintf("B%d", row), v.SortId)
		f.SetCellValue(sheet, fmt.Sprintf("C%d", row), v.Code)
		f.SetCellValue(sheet, fmt.Sprintf("D%d", row), *v.LbrSahamOwnership)
		f.SetCellValue(sheet, fmt.Sprintf("E%d", row), *v.TotalLbrSaham)
		f.SetCellFormula(sheet, fmt.Sprintf("F%d", row), fmt.Sprintf("=D%d/E%d", row, row))
		hargaPar := 0.0
		if v.HargaPar != nil {
			hargaPar = *v.HargaPar
		}
		hargaBeli := 0.0
		if v.HargaBeli != nil {
			hargaBeli = *v.HargaPar
		}
		f.SetCellValue(sheet, fmt.Sprintf("G%d", row), hargaPar)
		f.SetCellFormula(sheet, fmt.Sprintf("H%d", row), fmt.Sprintf("=D%d*G%d", row, row))
		f.SetCellValue(sheet, fmt.Sprintf("I%d", row), hargaBeli)
		f.SetCellFormula(sheet, fmt.Sprintf("J%d", row), fmt.Sprintf("=D%d*I%d", row, row))
		row++
	}
	// rowStart = row
	row = rowStart
	datePeriod, err := time.Parse(time.RFC3339, investasinontbk.Period)
	if err != nil {
		errs["NONTBK"] = "Error: Invalid Date Period"
		log.Println(err)
		return
	}
	period := datePeriod.Format("2006-01-02")
	err = f.SaveAs(fmt.Sprintf("%s/InvestasiNonTbk_%s.xlsx", tmpFolder, period))
	if err != nil {
		errs["NONTBK"] = "Error: Saving File Excel"
		log.Println(err)
		return
	}
	f.Close()

}

func (s *service) ExportInvestasiTbk(ctx *abstraction.Context, payload *abstraction.JsonData, wg *sync.WaitGroup, tmpFolder string) {
	defer wg.Done()
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
			errs["NONTBK"] = "Error: Generating Configuration file Excel"
			log.Println(err)
			return
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
		errs["NONTBK"] = "Error: Generating Configuration file Excel"
		log.Println(err)
		return
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
		errs["NONTBK"] = "Error: Generating Configuration file Excel"
		log.Println(err)
		return
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
		errs["NONTBK"] = "Error: Generating Configuration file Excel"
		log.Println(err)
		return
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
		errs["NONTBK"] = "Error: Generating Configuration file Excel"
		log.Println(err)
		return
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
		errs["NONTBK"] = "Error: Generating Configuration file Excel"
		log.Println(err)
		return
	}

	criteria := model.InvestasiTbkFilterModel{}
	criteria.CompanyID = &payload.CompanyID
	criteria.Versions = &payload.Filter.Versions
	criteria.Period = &payload.Filter.Period

	investasitbk, err := s.InvestasiTbkRepository.Find(ctx, &criteria)
	if err != nil {
		errs["NONTBK"] = "Error: Data Not Found"
		log.Println(err)
		return
	}

	datePeriod, err := time.Parse("2006-01-02", payload.Filter.Period)
	if err != nil {
		errs["NONTBK"] = "Error: Invalid Date Period"
		log.Println(err)
		return
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

		var criteria model.FormatterFilterModel
		criteria.FormatterFor = &formatter

		data, err := s.FormatterRepository.FindWithDetail(ctx, &criteria)
		if err != nil {
			errs["NONTBK"] = "Error: Formatter Data for Investasi Tbk Not Found"
			log.Println(err)
			return
		}

		tmpStr := "INVESTASI-TBK"
		criteriaBridge := model.FormatterBridgesFilterModel{}
		criteriaBridge.FormatterBridgesFilter.Source = &tmpStr
		criteriaBridge.FormatterBridgesFilter.FormatterID = &data.ID
		for _, valmia := range *investasitbk {
			criteriaBridge.FormatterBridgesFilter.TrxRefID = &valmia.ID
		}

		bridges, err := s.FormatterBridgesRepository.FindWithCriteria(ctx, &criteriaBridge)
		if err != nil {
			errs["NONTBK"] = "Error: Data Investasi Tbk Not Found"
			log.Println(err)
			return
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

				for chr := 'D'; chr <= 'K'; chr++ {
					f.SetCellFormula(sheet, fmt.Sprintf("%c%d", chr, row), fmt.Sprintf("=SUM(%c%d:%c%d)", chr, partRowStart, chr, row-1))
				}
				row++
				partRowStart = row
				continue
			}

			if v.IsLabel != nil && *v.IsLabel {
				row++
				continue
			}

			criteriaIT := model.InvestasiTbkDetailFilterModel{}
			criteriaIT.Stock = &v.Code
			criteriaIT.FormatterBridgesID = &bridges.ID
			for _, invest := range *investasitbk {
				criteriaIT.InvestasiTbkID = &invest.ID
			}

			InvestasiTbkDetail, err := s.InvestasiTbkDetailRepository.Find(ctx, &criteriaIT)
			if err != nil {
				errs["NONTBK"] = "Error: Data Investasi Tbk Not Found"
				log.Println(err)
				return
			}

			if len(*InvestasiTbkDetail) == 0 {
				errs["NONTBK"] = "Error: Data Investasi Tbk Not Found"
				log.Println("Data Not Found")
				return
			}

			for _, vv := range *InvestasiTbkDetail {
				f.SetCellValue(sheet, fmt.Sprintf("B%d", row), vv.SortId)
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
	period := datePeriod.Format("2006-01-02")
	savein := fmt.Sprintf("%s/InvestasiTbk_%s.xlsx", tmpFolder, period)
	fmt.Println(savein)
	err = f.SaveAs(fmt.Sprintf("%s/InvestasiTbk_%s.xlsx", tmpFolder, period))
	if err != nil {
		errs["NONTBK"] = "Error: Saving File Excel"
		log.Println(err)
		return
	}
	f.Close()

}

func (s *service) ExportMutasiDta(ctx *abstraction.Context, payload *abstraction.JsonData, wg *sync.WaitGroup, tmpFolder string) {
	defer wg.Done()
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
			errs["DTA"] = "Error: Generating Configuration Excel"
			log.Println(err)
			return
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
		errs["DTA"] = "Error: Generating Configuration Excel"
		log.Println(err)
		return
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
		errs["DTA"] = "Error: Generating Configuration Excel"
		log.Println(err)
		return
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
		errs["DTA"] = "Error: Generating Configuration Excel"
		log.Println(err)
		return
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
		errs["DTA"] = "Error: Generating Configuration Excel"
		log.Println(err)
		return
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
		errs["DTA"] = "Error: Generating Configuration Excel"
		log.Println(err)
		return
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

	criteria := model.MutasiDtaFilterModel{}
	criteria.CompanyID = &payload.CompanyID
	criteria.Versions = &payload.Filter.Versions
	criteria.Period = &payload.Filter.Period

	mutasidtaDatas, err := s.MutasiDtaRepository.Find(ctx, &criteria)
	if err != nil {
		errs["DTA"] = "Error: Data Mutasi Dta Not Found"
		log.Println(err)
		return
	}
	mutasidta := model.MutasiDtaEntityModel{}
	for _, valmutasidta := range *mutasidtaDatas {
		mutasidta = valmutasidta
	}

	datePeriod, err := time.Parse(time.RFC3339, mutasidta.Period)
	if err != nil {
		errs["DTA"] = "Error: Invalid Date Period"
		log.Println(err)
		return
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

		var criteria model.FormatterFilterModel
		criteria.FormatterFor = &formatter

		data, err := s.FormatterRepository.FindWithDetail(ctx, &criteria)
		if err != nil {
			errs["DTA"] = "Error: Formatter Data for Mutasi Dta Not Found"
			log.Println(err)
			return
		}

		tmpStr := "MUTASI-DTA"
		criteriaBridge := model.FormatterBridgesFilterModel{}
		criteriaBridge.FormatterBridgesFilter.Source = &tmpStr
		criteriaBridge.FormatterBridgesFilter.FormatterID = &data.ID
		criteriaBridge.FormatterBridgesFilter.TrxRefID = &mutasidta.ID

		bridges, err := s.FormatterBridgesRepository.FindWithCriteria(ctx, &criteriaBridge)
		if err != nil {
			errs["DTA"] = "Error: Data Mutasi Dta Not Found"
			log.Println(err)
			return
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

			if v.IsTotal != nil && *v.IsTotal {
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

			if v.AutoSummary != nil && *v.AutoSummary {
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

			mDtaDetail, err := s.MutasiDtaDetailRepository.Find(ctx, &criteriaMD)
			if err != nil && v.Code != "" {
				continue
			}
			if len(*mDtaDetail) == 0 {
				continue
			}
			for _, vv := range *mDtaDetail {
				f.SetCellValue(sheet, fmt.Sprintf("B%d", row), vv.SortId)
				f.SetCellValue(sheet, fmt.Sprintf("C%d", row), vv.Description)
				f.SetCellValue(sheet, fmt.Sprintf("D%d", row), *vv.SaldoAwal)
				f.SetCellValue(sheet, fmt.Sprintf("E%d", row), *vv.ManfaatBebanPajak)
				f.SetCellValue(sheet, fmt.Sprintf("F%d", row), *vv.Oci)
				f.SetCellValue(sheet, fmt.Sprintf("G%d", row), *vv.AkuisisiEntitasAnak)
				f.SetCellValue(sheet, fmt.Sprintf("H%d", row), *vv.DibebankanKeLr)
				f.SetCellValue(sheet, fmt.Sprintf("I%d", row), *vv.DibebankanKeOci)
				f.SetCellFormula(sheet, fmt.Sprintf("J%d", row), fmt.Sprintf("=SUM(D%d:I%d)", row, row))
			}
			row++
		}
		rowStart = row
		row = rowStart
	}
	period := datePeriod.Format("2006-01-02")
	err = f.SaveAs(fmt.Sprintf("%s/MutasiDta_%s.xlsx", tmpFolder, period))
	if err != nil {
		errs["DTA"] = "Error: Saving File Excel"
		log.Println(err)
		return
	}
	f.Close()
}

func (s *service) ExportMutasiFa(ctx *abstraction.Context, payload *abstraction.JsonData, wg *sync.WaitGroup, tmpFolder string) {
	defer wg.Done()
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
			errs["FA"] = "Error: Generating Configuration File Excel"
			log.Println(err)
			return
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
		errs["FA"] = "Error: Generating Configuration File Excel"
		log.Println(err)
		return
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
		errs["FA"] = "Error: Generating Configuration File Excel"
		log.Println(err)
		return
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
		errs["FA"] = "Error: Generating Configuration File Excel"
		log.Println(err)
		return
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
		errs["FA"] = "Error: Generating Configuration File Excel"
		log.Println(err)
		return
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
		errs["FA"] = "Error: Generating Configuration File Excel"
		log.Println(err)
		return
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

	criteria := model.MutasiFaFilterModel{}
	criteria.CompanyID = &payload.CompanyID
	criteria.Versions = &payload.Filter.Versions
	criteria.Period = &payload.Filter.Period

	mutasifaDatas, err := s.MutasiFaRepository.Find(ctx, &criteria)
	if err != nil {
		errs["FA"] = "Error: Data Mutasi Fa Not Found"
		log.Println(err)
		return
	}

	mutasifa := model.MutasiFaEntityModel{}
	for _, vMutasiFa := range *mutasifaDatas {
		mutasifa = vMutasiFa
	}

	datePeriod, err := time.Parse(time.RFC3339, mutasifa.Period)
	if err != nil {
		errs["FA"] = "Error: Invalid Date Period"
		log.Println(err)
		return
	}

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

		var criteria model.FormatterFilterModel
		criteria.FormatterFor = &formatter

		data, err := s.FormatterRepository.FindWithDetail(ctx, &criteria)
		if err != nil {
			errs["FA"] = "Error: Formatter Data for Mutasi Fa Not Found"
			log.Println(err)
			return
		}

		tmpStr := "MUTASI-FA"
		criteriaBridge := model.FormatterBridgesFilterModel{}
		criteriaBridge.FormatterBridgesFilter.Source = &tmpStr
		criteriaBridge.FormatterBridgesFilter.FormatterID = &data.ID
		criteriaBridge.FormatterBridgesFilter.TrxRefID = &mutasifa.ID

		bridges, err := s.FormatterBridgesRepository.FindWithCriteria(ctx, &criteriaBridge)
		if err != nil {
			errs["FA"] = "Error: Data Mutasi Fa Not Found"
			log.Println(err)
			return
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

			criteriaMIA := model.MutasiFaDetailFilterModel{}
			criteriaMIA.Code = &v.Code
			criteriaMIA.FormatterBridgesID = &bridges.ID
			criteriaMIA.MutasiFaID = &mutasifa.ID

			MFaDetail, err := s.MutasiFaDetailRepository.Find(ctx, &criteriaMIA)
			if err != nil && v.Code != "" {
				continue
			}
			if len(*MFaDetail) == 0 {
				continue
			}
			for _, vv := range *MFaDetail {
				f.SetCellValue(sheet, fmt.Sprintf("B%d", row), vv.Description)
				f.SetCellValue(sheet, fmt.Sprintf("D%d", row), *vv.BeginningBalance)
				f.SetCellValue(sheet, fmt.Sprintf("E%d", row), *vv.AcquisitionOfSubsidiary)
				f.SetCellValue(sheet, fmt.Sprintf("F%d", row), *vv.Additions)
				f.SetCellValue(sheet, fmt.Sprintf("G%d", row), *vv.Deductions)
				f.SetCellValue(sheet, fmt.Sprintf("H%d", row), *vv.Reclassification)
				f.SetCellValue(sheet, fmt.Sprintf("I%d", row), *vv.Revaluation)
				// f.SetCellValue(sheet, fmt.Sprintf("J%d", row), *vv.EndingBalance)
				f.SetCellFormula(sheet, fmt.Sprintf("J%d", row), fmt.Sprintf("=D%d+E%d+F%d+G%d+H%d+I%d", row, row, row, row, row, row))
				f.SetCellValue(sheet, fmt.Sprintf("K%d", row), *vv.Control)
			}
			row++
		}
		rowStart = row
		row = rowStart
	}
	period := datePeriod.Format("2006-01-02")
	err = f.SaveAs(fmt.Sprintf("%s/MutasiFa_%s.xlsx", tmpFolder, period))
	if err != nil {
		errs["FA"] = "Error: Saving File Excel"
		log.Println(err)
		return
	}
	f.Close()
}

func (s *service) ExportMutasiIa(ctx *abstraction.Context, payload *abstraction.JsonData, wg *sync.WaitGroup, tmpFolder string) {
	defer wg.Done()
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
			errs["IA"] = "Error: Generating Configuration File Excel"
			log.Println(err)
			return
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
		errs["IA"] = "Error: Generating Configuration File Excel"
		log.Println(err)
		return
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
		errs["IA"] = "Error: Generating Configuration File Excel"
		log.Println(err)
		return
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
		errs["IA"] = "Error: Generating Configuration File Excel"
		log.Println(err)
		return
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
		errs["IA"] = "Error: Generating Configuration File Excel"
		log.Println(err)
		return
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
		errs["IA"] = "Error: Generating Configuration File Excel"
		log.Println(err)
		return
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

	criteria := model.MutasiIaFilterModel{}
	criteria.CompanyID = &payload.CompanyID
	criteria.Versions = &payload.Filter.Versions
	criteria.Period = &payload.Filter.Period

	mutasiiaDatas, err := s.MutasiIaRepository.Find(ctx, &criteria)
	if err != nil {
		errs["IA"] = "Error: Data Mutasi IA Not Found"
		log.Println(err)
		return
	}

	mutasiia := model.MutasiIaEntityModel{}
	for _, vMutasiIa := range *mutasiiaDatas {
		mutasiia = vMutasiIa
	}

	datePeriod, err := time.Parse(time.RFC3339, mutasiia.Period)
	if err != nil {
		errs["IA"] = "Error: Invalid Date Period"
		log.Println(err)
		return
	}

	f.SetCellStyle(sheet, "B4", "K6", styleHeader)
	f.SetCellValue(sheet, "B4", mutasiia.Company.Name)
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

		var criteria model.FormatterFilterModel
		criteria.FormatterFor = &formatter

		data, err := s.FormatterRepository.FindWithDetail(ctx, &criteria)
		if err != nil {
			errs["IA"] = "Error: Formatter Data for Mutasi IA Not Found"
			log.Println(err)
			return
		}

		tmpStr := "MUTASI-IA"
		criteriaBridge := model.FormatterBridgesFilterModel{}
		criteriaBridge.FormatterBridgesFilter.Source = &tmpStr
		criteriaBridge.FormatterBridgesFilter.FormatterID = &data.ID
		criteriaBridge.FormatterBridgesFilter.TrxRefID = &mutasiia.ID

		bridges, err := s.FormatterBridgesRepository.FindWithCriteria(ctx, &criteriaBridge)
		if err != nil {
			errs["IA"] = "Error: Data Mutasi IA Not Found"
			log.Println(err)
			return
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

			criteriaMIA := model.MutasiIaDetailFilterModel{}
			criteriaMIA.Code = &v.Code
			criteriaMIA.FormatterBridgesID = &bridges.ID
			criteriaMIA.MutasiIaID = &mutasiia.ID

			MIaDetail, err := s.MutasiIaDetailRepository.Find(ctx, &criteriaMIA)
			if err != nil && v.Code != "" {
				continue
			}
			if len(*MIaDetail) == 0 {
				continue
			}
			for _, vv := range *MIaDetail {
				f.SetCellValue(sheet, fmt.Sprintf("B%d", row), vv.Description)
				f.SetCellValue(sheet, fmt.Sprintf("D%d", row), *vv.BeginningBalance)
				f.SetCellValue(sheet, fmt.Sprintf("E%d", row), *vv.AcquisitionOfSubsidiary)
				f.SetCellValue(sheet, fmt.Sprintf("F%d", row), *vv.Additions)
				f.SetCellValue(sheet, fmt.Sprintf("G%d", row), *vv.Deductions)
				f.SetCellValue(sheet, fmt.Sprintf("H%d", row), *vv.Reclassification)
				f.SetCellValue(sheet, fmt.Sprintf("I%d", row), *vv.Revaluation)
				// f.SetCellValue(sheet, fmt.Sprintf("J%d", row), *vv.EndingBalance)
				f.SetCellFormula(sheet, fmt.Sprintf("J%d", row), fmt.Sprintf("=D%d+E%d+F%d+G%d+H%d+I%d", row, row, row, row, row, row))
				f.SetCellValue(sheet, fmt.Sprintf("K%d", row), *vv.Control)
			}
			row++
		}
		rowStart = row
		row = rowStart
	}
	period := datePeriod.Format("2006-01-02")
	err = f.SaveAs(fmt.Sprintf("%s/MutasiIa_%s.xlsx", tmpFolder, period))
	if err != nil {
		errs["IA"] = "Error: Saving File Excel"
		log.Println(err)
		return
	}
	f.Close()
}

func (s *service) ExportMutasiPersediaan(ctx *abstraction.Context, payload *abstraction.JsonData, wg *sync.WaitGroup, tmpFolder string) {
	defer wg.Done()
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
		errs["PERSEDIAAN"] = "Error: Generating Configuration File Excel"
		log.Println(err)
		return
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
		errs["PERSEDIAAN"] = "Error: Generating Configuration File Excel"
		log.Println(err)
		return
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
		errs["PERSEDIAAN"] = "Error: Generating Configuration File Excel"
		log.Println(err)
		return
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
		errs["PERSEDIAAN"] = "Error: Generating Configuration File Excel"
		log.Println(err)
		return
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
		errs["PERSEDIAAN"] = "Error: Generating Configuration File Excel"
		log.Println(err)
		return
	}

	criteria := model.MutasiPersediaanFilterModel{}
	criteria.CompanyID = &payload.CompanyID
	criteria.Versions = &payload.Filter.Versions
	criteria.Period = &payload.Filter.Period

	mutasipersediaan, err := s.MutasiPersediaanRepository.Find(ctx, &criteria)
	if err != nil {
		errs["PERSEDIAAN"] = "Error: Data Mutasi Persediaan Not Found"
		log.Println(err)
		return
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
			errs["PERSEDIAAN"] = "Error: Formatter Data for Mutasi Persediaan Not Found"
			log.Println(err)
			return
		}

		tmpStr := "MUTASI-PERSEDIAAN"
		criteriaBridge := model.FormatterBridgesFilterModel{}
		criteriaBridge.FormatterBridgesFilter.Source = &tmpStr
		criteriaBridge.FormatterBridgesFilter.FormatterID = &data.ID
		for _, valmp := range *mutasipersediaan {
			criteriaBridge.FormatterBridgesFilter.TrxRefID = &valmp.ID
		}

		bridges, err := s.FormatterBridgesRepository.FindWithCriteria(ctx, &criteriaBridge)
		if err != nil {
			errs["PERSEDIAAN"] = "Error: Data Mutasi Persediaan Not Found"
			log.Println(err)
			return
		}

		rowCode := make(map[string]int)
		partRowStart := row
		for _, v := range data.FormatterDetail {
			f.SetCellStyle(sheet, fmt.Sprintf("B%d", row), fmt.Sprintf("B%d", row), styleLabel)
			f.SetCellStyle(sheet, fmt.Sprintf("C%d", row), fmt.Sprintf("C%d", row), styleCurrency)

			if strings.Contains(strings.ToLower(v.Code), "blank") {
				row++
				continue
			}

			f.SetCellValue(sheet, fmt.Sprintf("B%d", row), v.Description)
			if v.IsTotal != nil && *v.IsTotal {
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

			if v.AutoSummary != nil && *v.AutoSummary {
				f.SetCellStyle(sheet, fmt.Sprintf("B%d", row), fmt.Sprintf("B%d", row), styleLabelTotal)
				f.SetCellStyle(sheet, fmt.Sprintf("C%d", row), fmt.Sprintf("C%d", row), styleCurrencyTotal)
				// f.MergeCell(sheet, fmt.Sprintf("B%d", row), fmt.Sprintf("C%d", row))

				f.SetCellFormula(sheet, fmt.Sprintf("C%d", row), fmt.Sprintf("=SUM(C%d:C%d)", partRowStart, row-1))
				row++
				partRowStart = row
				continue
			}

			if v.IsLabel != nil && *v.IsLabel {
				row++
				continue
			}

			criteriaMP := model.MutasiPersediaanDetailFilterModel{}
			criteriaMP.Code = &v.Code
			criteriaMP.FormatterBridgesID = &bridges.ID
			for _, valmutasipersediaan := range *mutasipersediaan {
				criteriaMP.MutasiPersediaanID = &valmutasipersediaan.ID
			}

			mutasiPersediaanDetail, err := s.MutasiPersediaanDetailRepository.Find(ctx, &criteriaMP)
			if err != nil {
				errs["PERSEDIAAN"] = "Error: Data Mutasi Persediaan Not Found"
				log.Println(err)
				return
			}
			if len(*mutasiPersediaanDetail) == 0 {
				errs["PERSEDIAAN"] = "Error: Data Mutasi Persediaan Not Found"
				log.Println("Data Not Found")
				return
			}
			for _, vv := range *mutasiPersediaanDetail {
				f.SetCellValue(sheet, fmt.Sprintf("C%d", row), *vv.Amount)
			}
			row++
		}
		rowStart = row + 4
		row = rowStart
	}
	datePeriod, err := time.Parse("2006-01-02", payload.Filter.Period)
	if err != nil {
		errs["PERSEDIAAN"] = "Error: Invalid Date Period"
		log.Println(err)
		return
	}
	period := datePeriod.Format("2006-01-02")
	err = f.SaveAs(fmt.Sprintf("%s/MutasiPersediaan_%s.xlsx", tmpFolder, period))
	if err != nil {
		errs["PERSEDIAAN"] = "Error: Saving File Excel"
		log.Println(err)
		return
	}
	f.Close()
}

func (s *service) ExportMutasiRua(ctx *abstraction.Context, payload *abstraction.JsonData, wg *sync.WaitGroup, tmpFolder string) {
	defer wg.Done()
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
			errs["RUA"] = "Error: Generating Configuration File Excel"
			log.Println(err)
			return
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
		errs["RUA"] = "Error: Generating Configuration File Excel"
		log.Println(err)
		return
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
		errs["RUA"] = "Error: Generating Configuration File Excel"
		log.Println(err)
		return
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
		errs["RUA"] = "Error: Generating Configuration File Excel"
		log.Println(err)
		return
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
		errs["RUA"] = "Error: Generating Configuration File Excel"
		log.Println(err)
		return
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
		errs["RUA"] = "Error: Generating Configuration File Excel"
		log.Println(err)
		return
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

	criteria := model.MutasiRuaFilterModel{}
	criteria.CompanyID = &payload.CompanyID
	criteria.Versions = &payload.Filter.Versions
	criteria.Period = &payload.Filter.Period

	mutasiruaDatas, err := s.MutasiRuaRepository.Find(ctx, &criteria)
	if err != nil {
		errs["RUA"] = "Error: Data Mutasi Rua Not Found"
		log.Println(err)
		return
	}

	if len(*mutasiruaDatas) == 0 {
		errs["RUA"] = "Error: Data Mutasi Rua Not Found"
		log.Println("Data Not Found!")
		return
	}

	mutasirua := model.MutasiRuaEntityModel{}
	for _, vMutasiRua := range *mutasiruaDatas {
		mutasirua = vMutasiRua
	}

	datePeriod, err := time.Parse(time.RFC3339, mutasirua.Period)
	if err != nil {
		errs["RUA"] = "Error: Invalid Date Period"
		log.Println("Data Not Found!")
		return
	}

	f.SetCellStyle(sheet, "B4", "K6", styleHeader)
	f.SetCellValue(sheet, "B4", mutasirua.Company.Name)
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

		var criteria model.FormatterFilterModel
		criteria.FormatterFor = &formatter

		data, err := s.FormatterRepository.FindWithDetail(ctx, &criteria)
		if err != nil {
			errs["RUA"] = "Error: Formatter Data for Mutasi Rua Not Found"
			log.Println(err)
			return
		}

		tmpStr := "MUTASI-RUA"
		criteriaBridge := model.FormatterBridgesFilterModel{}
		criteriaBridge.FormatterBridgesFilter.Source = &tmpStr
		criteriaBridge.FormatterBridgesFilter.FormatterID = &data.ID
		criteriaBridge.FormatterBridgesFilter.TrxRefID = &mutasirua.ID

		bridges, err := s.FormatterBridgesRepository.FindWithCriteria(ctx, &criteriaBridge)
		if err != nil {
			errs["RUA"] = "Error: Data Mutasi Rua Not Found"
			log.Println(err)
			return
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

			criteriaMR := model.MutasiRuaDetailFilterModel{}
			criteriaMR.Code = &v.Code
			criteriaMR.FormatterBridgesID = &bridges.ID
			criteriaMR.MutasiRuaID = &mutasirua.ID

			MutasiRuaDetail, err := s.MutasiRuaDetailRepository.Find(ctx, &criteriaMR)
			if err != nil {
				errs["RUA"] = "Error: Data Mutasi Rua Not Found"
				log.Println(err)
				return
			}
			if len(*MutasiRuaDetail) == 0 {
				errs["RUA"] = "Error: Data Mutasi Rua Not Found"
				log.Println(err)
				return
			}
			for _, vv := range *MutasiRuaDetail {
				f.SetCellValue(sheet, fmt.Sprintf("B%d", row), vv.Description)
				f.SetCellValue(sheet, fmt.Sprintf("D%d", row), *vv.BeginningBalance)
				f.SetCellValue(sheet, fmt.Sprintf("E%d", row), *vv.AcquisitionOfSubsidiary)
				f.SetCellValue(sheet, fmt.Sprintf("F%d", row), *vv.Additions)
				f.SetCellValue(sheet, fmt.Sprintf("G%d", row), *vv.Deductions)
				f.SetCellValue(sheet, fmt.Sprintf("H%d", row), *vv.Reclassification)
				f.SetCellValue(sheet, fmt.Sprintf("I%d", row), *vv.Remeasurement)
				f.SetCellFormula(sheet, fmt.Sprintf("J%d", row), fmt.Sprintf("=D%d+E%d+F%d+G%d+H%d+I%d", row, row, row, row, row, row))
				f.SetCellValue(sheet, fmt.Sprintf("K%d", row), *vv.Control)
			}
			row++
		}
		rowStart = row
		row = rowStart
	}
	period := datePeriod.Format("2006-01-02")
	err = f.SaveAs(fmt.Sprintf("%s/MutasiRua_%s.xlsx", tmpFolder, period))
	if err != nil {
		errs["RUA"] = "Error: Saving File Excel"
		log.Println(err)
		return
	}
	f.Close()
}

func (s *service) ExportPembelianPenjualanBerelasi(ctx *abstraction.Context, payload *abstraction.JsonData, wg *sync.WaitGroup, tmpFolder string) {
	defer wg.Done()
	f := excelize.NewFile()
	sheet := f.GetSheetName(f.GetActiveSheetIndex())
	f.SetColWidth(sheet, "D", "D", 66)

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
		errs["PPB"] = "Error: Generating Configuration File Excel"
		log.Println(err)
		return
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
		errs["PPB"] = "Error: Generating Configuration File Excel"
		log.Println(err)
		return
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
		errs["PPB"] = "Error: Generating Configuration File Excel"
		log.Println(err)
		return
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
		errs["PPB"] = "Error: Generating Configuration File Excel"
		log.Println(err)
		return
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
		errs["PPB"] = "Error: Generating Configuration File Excel"
		log.Println(err)
		return
	}

	f.SetCellStyle(sheet, "B4", "F4", styleHeader)
	f.SetCellValue(sheet, "C2", "List pembelian dan penjualan berelasi")
	f.SetCellValue(sheet, "B4", "NO")
	f.SetCellValue(sheet, "C4", "CODE")
	f.SetCellValue(sheet, "D4", "COMPANY")
	f.SetCellValue(sheet, "E4", "PEMBELIAN")
	f.SetCellValue(sheet, "F4", "PENJUALAN")

	criteriapbb := model.PembelianPenjualanBerelasiFilterModel{}
	criteriapbb.Period = &payload.Filter.Period
	criteriapbb.CompanyID = &payload.CompanyID
	criteriapbb.Versions = &payload.Filter.Versions

	pbbDatas, err := s.PembelianPenjualanBerelasiRepository.Find(ctx, &criteriapbb)
	if err != nil {
		errs["PPB"] = "Error: Data Pembelian Penjualan Berelasi Not Found"
		log.Println(err)
		return
	}

	tmpID := 0
	for _, vPbb := range *pbbDatas {
		tmpID = vPbb.ID
	}

	pembelianpenjualanberelasi, err := s.PembelianPenjualanBerelasiRepository.FindByID(ctx, &tmpID)
	if err != nil {
		errs["PPB"] = "Error: Data Pembelian Penjualan Berelasi Not Found"
		log.Println(err)
		return
	}

	criteria := model.PembelianPenjualanBerelasiFilterModel{}
	criteria.CompanyID = &pembelianpenjualanberelasi.CompanyID
	criteria.Period = &pembelianpenjualanberelasi.Period
	criteria.Versions = &pembelianpenjualanberelasi.Versions
	data, err := s.PembelianPenjualanBerelasiRepository.Export(ctx, &criteria)
	if err != nil {
		errs["PPB"] = "Error: Data Pembelian Penjualan Berelasi Not Found"
		log.Println(err)
		return
	}

	row, rowStart := 5, 5
	for i, v := range data.PembelianPenjualanBerelasiDetail {
		f.SetCellValue(sheet, fmt.Sprintf("B%d", row), (i + 1))
		f.SetCellValue(sheet, fmt.Sprintf("C%d", row), v.Company.Code)
		f.SetCellValue(sheet, fmt.Sprintf("D%d", row), v.Company.Name)
		f.SetCellValue(sheet, fmt.Sprintf("E%d", row), *v.BoughtAmount)
		f.SetCellValue(sheet, fmt.Sprintf("F%d", row), *v.SalesAmount)
		row++
	}

	f.SetCellStyle(sheet, "B5", fmt.Sprintf("D%d", row), styleLabel)
	f.SetCellStyle(sheet, "E5", fmt.Sprintf("F%d", row), styleCurrency)
	f.SetCellStyle(sheet, fmt.Sprintf("B%d", row+1), fmt.Sprintf("D%d", row+1), styleLabelTotal)
	f.SetCellStyle(sheet, fmt.Sprintf("E%d", row+1), fmt.Sprintf("F%d", row+1), styleCurrencyTotal)

	f.SetCellValue(sheet, fmt.Sprintf("D%d", row+1), "Total")
	f.SetCellFormula(sheet, fmt.Sprintf("E%d", row+1), fmt.Sprintf("=SUM(E%d:E%d)", rowStart, row))
	f.SetCellFormula(sheet, fmt.Sprintf("F%d", row+1), fmt.Sprintf("=SUM(F%d:F%d)", rowStart, row))
	datePeriod, err := time.Parse(time.RFC3339, pembelianpenjualanberelasi.Period)
	if err != nil {
		errs["PPB"] = "Error: Invali Date Period"
		log.Println(err)
		return
	}
	period := datePeriod.Format("2006-01-02")
	err = f.SaveAs(fmt.Sprintf("%s/PembelianPenjualanBerelasi_%s.xlsx", tmpFolder, period))
	if err != nil {
		errs["PPB"] = "Error: Saving File Excel"
		log.Println(err)
		return
	}
	f.Close()
}

func (s *service) ExportTrialBalance(ctx *abstraction.Context, payload *abstraction.JsonData, wg *sync.WaitGroup, tmpFolder string) {
	defer wg.Done()
	var (
		criteria          model.TrialBalanceFilterModel
		criteriaFormatter model.FormatterDetailFilterModel
	)

	criteria.CompanyID = &payload.CompanyID
	criteria.Versions = &payload.Filter.Versions
	criteria.Period = &payload.Filter.Period

	tb, err := s.TrialBalanceRepository.Get(ctx, &criteria)
	if err != nil {
		errs["TB"] = "Error: Data Trial Balance Not Found"
		log.Println(err)
		return
	}
	if len(tb.FormatterBridges) == 0 {
		tmpErr := errors.New("no formatter found")
		errs["TB"] = "Error: Data Pembelian Penjualan Berelasi Not Found"
		log.Println(tmpErr)
		return
	}
	var formatterID int
	for _, fmtbridges := range tb.FormatterBridges {
		formatterID = fmtbridges.FormatterID
	}
	criteriaFormatter.FormatterID = &formatterID

	data, err := s.FormatterDetailRepository.Find(ctx, &criteriaFormatter)
	if err != nil {
		errs["TB"] = "Error: Formatter Data for Trial Balance Not Found"
		log.Println(err)
		return
	}

	datePeriod, err := time.Parse(time.RFC3339, tb.Period)
	if err != nil {
		errs["TB"] = "Error: Invalid Date Period"
		log.Println(err)
		return
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
			errs["TB"] = "Error: Generating Configuration File Excel"
			log.Println(err)
			return
		}
	}

	styleDefault, err := f.NewStyle(&excelize.Style{
		Font: &excelize.Font{
			Family: "Arial",
			Size:   10,
		},
	})
	if err != nil {
		errs["TB"] = "Error: Generating Configuration File Excel"
		log.Println(err)
		return
	}
	err = f.SetColStyle(sheet, "A:Z", styleDefault)
	if err != nil {
		errs["TB"] = "Error: Generating Configuration File Excel"
		log.Println(err)
		return
	}

	stylingBorderLROnly, err := f.NewStyle(&excelize.Style{
		Border: []excelize.Border{
			{Type: "right", Color: "000000", Style: 1},
			{Type: "left", Color: "000000", Style: 1},
		},
	})
	if err != nil {
		errs["TB"] = "Error: Generating Configuration File Excel"
		log.Println(err)
		return
	}

	stylingBorderRightOnly, err := f.NewStyle(&excelize.Style{
		Border: []excelize.Border{
			{Type: "right", Color: "000000", Style: 1},
		},
	})
	if err != nil {
		errs["TB"] = "Error: Generating Configuration File Excel"
		log.Println(err)
		return
	}
	stylingBorderTopOnly, err := f.NewStyle(&excelize.Style{
		Border: []excelize.Border{
			{Type: "top", Color: "000000", Style: 1},
		},
	})
	if err != nil {
		errs["TB"] = "Error: Generating Configuration File Excel"
		log.Println(err)
		return
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
		errs["TB"] = "Error: Generating Configuration File Excel"
		log.Println(err)
		return
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
		errs["TB"] = "Error: Generating Configuration File Excel"
		log.Println(err)
		return
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
	stylingSubTotalCurrency, err := f.NewStyle(&excelize.Style{
		Border: []excelize.Border{
			{Type: "top", Color: "000000", Style: 1},
			{Type: "bottom", Color: "000000", Style: 1},
		},
		Font: &excelize.Font{
			Bold: true,
		},
		NumFmt:       7,
		CustomNumFmt: &numberFormat,
	})
	if err != nil {
		errs["TB"] = "Error: Generating Configuration File Excel"
		log.Println(err)
		return
	}

	stylingTotal, err := f.NewStyle(&excelize.Style{
		Border: []excelize.Border{
			{Type: "top", Color: "000000", Style: 1},
			{Type: "bottom", Color: "000000", Style: 1},
		},
		Font: &excelize.Font{
			Bold: true,
		},
		Fill: excelize.Fill{
			Type:    "pattern",
			Pattern: 1,
			Color:   []string{"#ccff33"},
		},
	})
	// stylingTotalCurrency, err := f.NewStyle(&excelize.Style{
	// 	Border: []excelize.Border{
	// 		{Type: "top", Color: "000000", Style: 1},
	// 		{Type: "bottom", Color: "000000", Style: 1},
	// 	},
	// 	Font: &excelize.Font{
	// 		Bold: true,
	// 	},
	// 	Fill: excelize.Fill{
	// 		Type:    "pattern",
	// 		Pattern: 1,
	// 		Color:   []string{"#ccff33"},
	// 	},
	// 	NumFmt:       7,
	// 	CustomNumFmt: &numberFormat,
	// })
	// if err != nil {
	// 	errs["TB"] = "Error: Generating Configuration File Excel"
	// 	log.Println(err)
	// 	return
	// }

	stylingTotalControl, err := f.NewStyle(&excelize.Style{
		Border: []excelize.Border{
			{Type: "top", Color: "000000", Style: 1},
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
	})
	// stylingTotalControlCurrency, err := f.NewStyle(&excelize.Style{
	// 	Border: []excelize.Border{
	// 		{Type: "top", Color: "000000", Style: 1},
	// 		{Type: "bottom", Color: "000000", Style: 1},
	// 	},
	// 	Font: &excelize.Font{
	// 		Bold: true,
	// 	},
	// 	Fill: excelize.Fill{
	// 		Type:    "pattern",
	// 		Pattern: 1,
	// 		Color:   []string{"#3ada24"},
	// 	},
	// 	NumFmt:       7,
	// 	CustomNumFmt: &numberFormat,
	// })
	// if err != nil {
	// 	errs["TB"] = "Error: Generating Configuration File Excel"
	// 	log.Println(err)
	// 	return
	// }

	err = f.MergeCell(sheet, "B6", "B8")
	if err != nil {
		errs["TB"] = "Error: Generating Configuration File Excel"
		log.Println(err)
		return
	}
	err = f.MergeCell(sheet, "C6", "E8")
	if err != nil {
		errs["TB"] = "Error: Generating Configuration File Excel"
		log.Println(err)
		return
	}
	err = f.MergeCell(sheet, "F6", "F8")
	if err != nil {
		errs["TB"] = "Error: Generating Configuration File Excel"
		log.Println(err)
		return
	}
	err = f.MergeCell(sheet, "H6", "K7")
	if err != nil {
		errs["TB"] = "Error: Generating Configuration File Excel"
		log.Println(err)
		return
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
		errs["TB"] = "Error: Generating Configuration File Excel"
		log.Println(err)
		return
	}

	err = f.SetCellStyle(sheet, "B6", "L8", stylingHeader)
	if err != nil {
		errs["TB"] = "Error: Generating Configuration File Excel"
		log.Println(err)
		return
	}
	err = f.SetCellStyle(sheet, "F6", "F8", stylingHeader2)
	if err != nil {
		errs["TB"] = "Error: Generating Configuration File Excel"
		log.Println(err)
		return
	}

	err = f.SetCellValue(sheet, "B6", "No Akun")
	if err != nil {
		errs["TB"] = "Error: Generating Configuration File Excel"
		log.Println(err)
		return
	}
	err = f.SetCellValue(sheet, "C6", "Keterangan")
	if err != nil {
		errs["TB"] = "Error: Generating Configuration File Excel"
		log.Println(err)
		return
	}
	err = f.SetCellValue(sheet, "F6", "WP Reff")
	if err != nil {
		errs["TB"] = "Error: Generating Configuration File Excel"
		log.Println(err)
		return
	}
	err = f.SetCellValue(sheet, "G6", tb.Company.Name)
	if err != nil {
		errs["TB"] = "Error: Generating Configuration File Excel"
		log.Println(err)
		return
	}
	err = f.SetCellValue(sheet, "G7", "Unaudited")
	if err != nil {
		errs["TB"] = "Error: Generating Configuration File Excel"
		log.Println(err)
		return
	}
	err = f.SetCellValue(sheet, "G8", datePeriod.Format("02-Jan-06"))
	if err != nil {
		errs["TB"] = "Error: Generating Configuration File Excel"
		log.Println(err)
		return
	}
	err = f.SetCellValue(sheet, "H6", "Adjustment Journal Entry")
	if err != nil {
		errs["TB"] = "Error: Generating Configuration File Excel"
		log.Println(err)
		return
	}
	err = f.SetCellValue(sheet, "I8", "Debet")
	if err != nil {
		errs["TB"] = "Error: Generating Configuration File Excel"
		log.Println(err)
		return
	}
	err = f.SetCellValue(sheet, "K8", "Kredit")
	if err != nil {
		errs["TB"] = "Error: Generating Configuration File Excel"
		log.Println(err)
		return
	}
	err = f.SetCellFormula(sheet, "L6", "=G6")
	if err != nil {
		errs["TB"] = "Error: Generating Configuration File Excel"
		log.Println(err)
		return
	}
	err = f.SetCellFormula(sheet, "L7", "=G7")
	if err != nil {
		errs["TB"] = "Error: Generating Configuration File Excel"
		log.Println(err)
		return
	}
	err = f.SetCellFormula(sheet, "L8", "=G8")
	if err != nil {
		errs["TB"] = "Error: Generating Configuration File Excel"
		log.Println(err)
		return
	}

	criteriaBridge := model.FormatterBridgesFilterModel{}
	tmpStr := "TRIAL-BALANCE"
	criteriaBridge.FormatterBridgesFilter.Source = &tmpStr
	criteriaBridge.FormatterBridgesFilter.FormatterID = &formatterID
	criteriaBridge.FormatterBridgesFilter.TrxRefID = &tb.ID
	bridges, err := s.FormatterBridgesRepository.FindWithCriteria(ctx, &criteriaBridge)
	if err != nil {
		errs["TB"] = "Error: Data Trial Balance Not Found"
		log.Println(err)
		return
	}

	//find summary aje
	summaryAJE, err := s.AjeRepository.FindSummary(ctx, &tb.ID)
	if err != nil {
		errs["TB"] = "Error: Error when find summary aje"
		log.Println(err)
		return
	}

	row := 9
	// var total []map[string]interface{}
	rowCode := make(map[string]int)
	isAutoSum := make(map[string]bool)
	tbRowCode := make(map[string]int)
	customRow := make(map[string]string)

	for _, v := range *data {
		rowCode[v.Code] = row
		if v.AutoSummary != nil && *v.AutoSummary {
			isAutoSum[v.Code] = true
		}

		if err = f.SetCellStyle(sheet, fmt.Sprintf("G%d", row), fmt.Sprintf("L%d", row), stylingCurrency); err != nil {
			errs["TB"] = "Error: Generating Configuration File Excel"
			log.Println(err)
			return
		}
		// var codeCoa string
		if !(v.IsTotal != nil && *v.IsTotal) && v.IsLabel != nil && *v.IsLabel {
			if err = f.SetCellStyle(sheet, fmt.Sprintf("G%d", row), fmt.Sprintf("L%d", row), stylingCurrency); err != nil {
				errs["TB"] = "Error: Generating Configuration File Excel"
				log.Println(err)
				return
			}
			if err = f.SetCellStyle(sheet, fmt.Sprintf("F%d", row), fmt.Sprintf("F%d", row), stylingBorderLROnly); err != nil {
				errs["TB"] = "Error: Generating Configuration File Excel"
				log.Println(err)
				return
			}
			f.SetCellValue(sheet, fmt.Sprintf("C%d", row), v.Description)
		}
		if v.IsCoa != nil && *v.IsCoa {
			rowBefore := row
			tbdetails, err := s.TrialBalanceDetailRepository.FindToExport(ctx, &v.Code, &bridges.ID)
			if err != nil {
				errs["TB"] = "Error: Data Trial Balance Detail Not Found for code: " + v.Code
				return
			}
			for _, vTbDetail := range *tbdetails {
				if err = f.SetCellStyle(sheet, fmt.Sprintf("G%d", row), fmt.Sprintf("L%d", row), stylingCurrency); err != nil {
					errs["TB"] = "Error: Generating Configuration File Excel"
					log.Println(err)
					return
				}
				if err = f.SetCellStyle(sheet, fmt.Sprintf("F%d", row), fmt.Sprintf("F%d", row), stylingBorderLROnly); err != nil {
					errs["TB"] = "Error: Generating Configuration File Excel"
					log.Println(err)
					return
				}
				if strings.Contains(strings.ToUpper(vTbDetail.Code), "SUBTOTAL") {
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
					amountBeforeAje = *vTbDetail.AmountAjeDr
				}
				amountAjeCr := 0.0
				if vTbDetail.AmountAjeCr != nil {
					amountBeforeAje = *vTbDetail.AmountAjeCr
				}
				f.SetCellValue(sheet, fmt.Sprintf("G%d", row), amountBeforeAje)
				f.SetCellValue(sheet, fmt.Sprintf("I%d", row), amountAjeDr)
				f.SetCellValue(sheet, fmt.Sprintf("K%d", row), amountAjeCr)
				// f.SetCellFormula(sheet, fmt.Sprintf("L%d", row), fmt.Sprintf("=G%d+I%d-K%d", row, row, row))
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
			rowTB := len(*tbdetails)
			rowAfter := row - 1
			if v.AutoSummary != nil && *v.AutoSummary && rowTB > 0 {
				f.SetCellValue(sheet, fmt.Sprintf("D%d", row), "Subtotal")
				f.SetCellFormula(sheet, fmt.Sprintf("G%d", row), fmt.Sprintf("=SUM(G%d:G%d)", rowBefore, rowAfter))
				f.SetCellFormula(sheet, fmt.Sprintf("I%d", row), fmt.Sprintf("=SUM(I%d:I%d)", rowBefore, rowAfter))
				f.SetCellFormula(sheet, fmt.Sprintf("K%d", row), fmt.Sprintf("=SUM(K%d:K%d)", rowBefore, rowAfter))
				f.SetCellFormula(sheet, fmt.Sprintf("L%d", row), fmt.Sprintf("=SUM(L%d:L%d)", rowBefore, rowAfter))
				f.SetCellStyle(sheet, fmt.Sprintf("D%d", row), fmt.Sprintf("F%d", row), stylingSubTotal)
				f.SetCellStyle(sheet, fmt.Sprintf("G%d", row), fmt.Sprintf("L%d", row), stylingSubTotalCurrency)
				rowCode[fmt.Sprintf("%s_SUBTOTAL", v.Code)] = row
				row++
				f.SetCellStyle(sheet, fmt.Sprintf("G%d", row), fmt.Sprintf("L%d", row), stylingCurrency)
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
				f.SetCellFormula(sheet, "I5", fmt.Sprintf("=I%d-I%d", rowAset, row))
				f.SetCellFormula(sheet, "K5", fmt.Sprintf("=K%d-K%d", rowAset, row))
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
				f.SetCellFormula(sheet, fmt.Sprintf("H%d", row), fmt.Sprintf("=I%d-K%d", row, row))
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
				if fmt.Sprintf("%c", chr) == "H" || fmt.Sprintf("%c", chr) == "J" {
					continue
				}
				formula := v.FxSummary
				reg := regexp.MustCompile(`[A-Za-z0-9_~:()]+|[0-9]+\d{2,}`)
				match := reg.FindAllString(formula, -1)
				for _, vMatch := range match {
					if len(vMatch) < 4 {
						continue
					}
					//cari jml berdasarkan code
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
	customRow["310402002"] = "=TOTAL_PENGHASILAN_KOMPREHENSIF_LAIN_~_BS-SUM(310501002,310502002,310503002)"
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
			f.SetCellFormula(sheet, fmt.Sprintf("L%d", val), strings.ReplaceAll(vCustomRow, "@", "L"))
		}
	}

	if err = f.SetCellStyle(sheet, fmt.Sprintf("B%d", row), fmt.Sprintf("L%d", row), stylingBorderTopOnly); err != nil {
		errs["TB"] = "Error: Generating Configuration File Excel"
		log.Println(err)
		return
	}

	// if err = f.SetCellStyle(sheet, "G9", fmt.Sprintf("L%d", row-1), stylingCurrency); err != nil {
	// 	errs["TB"] = "Error: Generating Configuration File Excel"
	// 	log.Println(err)
	// 	return
	// }

	if err = f.SetCellStyle(sheet, "A9", fmt.Sprintf("A%d", row-1), stylingBorderRightOnly); err != nil {
		errs["TB"] = "Error: Generating Configuration File Excel"
		log.Println(err)
		return
	}

	if err = f.SetSheetFormatPr(sheet, excelize.DefaultRowHeight(12.85)); err != nil {
		errs["TB"] = "Error: Generating Configuration File Excel"
		log.Println(err)
		return
	}
	f.SetDefaultFont("Arial")
	period := datePeriod.Format("2006-01-02")
	err = f.SaveAs(fmt.Sprintf("%s/TrialBalance_%s.xlsx", tmpFolder, period))
	if err != nil {
		errs["TB"] = "Error: Saving File Excel"
		log.Println(err)
		return
	}
	f.Close()
}

func (s *service) ExportJpmWorksheet(ctx *abstraction.Context, consolidationData *model.ConsolidationEntityModel, f *excelize.File) (*excelize.File, error) {
	datas, err := s.JpmRepository.ExportAll(ctx, &consolidationData.ID)
	if err != nil {
		return nil, err
	}

	datePeriod, err := time.Parse(time.RFC3339, consolidationData.Period)
	if err != nil {
		return nil, err
	}

	sheet := "JPM"
	f.NewSheet(sheet)
	arrStyleColWidth := []map[string]interface{}{
		{"COL": "A", "WIDTH": 6.43},
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
		err = f.SetColWidth(sheet, fmt.Sprintf("%s", v["COL"]), fmt.Sprintf("%s", v["COL"]), colWidth)
		if err != nil {
			return nil, err
		}
	}

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

	f.MergeCell(sheet, "A6", "A7")
	f.MergeCell(sheet, "B6", "B7")
	f.MergeCell(sheet, "C6", "C7")
	f.MergeCell(sheet, "D6", "E7")
	f.MergeCell(sheet, "F6", "F7")
	f.MergeCell(sheet, "G6", "H6")
	f.MergeCell(sheet, "I6", "j6")

	f.SetCellValue(sheet, "A2", "Company")
	f.SetCellValue(sheet, "B2", ": "+consolidationData.Company.Name)
	f.SetCellValue(sheet, "A3", "Date")
	f.SetCellValue(sheet, "B3", ": "+datePeriod.Format("02-Jan-06"))
	f.SetCellValue(sheet, "A4", "Subject")
	f.SetCellValue(sheet, "B4", ": Proforma Modal")

	f.SetCellStyle(sheet, "A6", "J7", stylingHeader)
	f.SetCellValue(sheet, "A6", "NO")
	f.SetCellValue(sheet, "B6", "COA")
	f.SetCellValue(sheet, "C6", "JPM")
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
	reffJpm := make(map[string]int)
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
		if reffJpm[*v.ReffNumber] == 0 {
			row += 1
			if *v.Note != "" {
				row += 1
			}
		}
		if reffJpm[*v.ReffNumber] != 0 {
			row += 2
			if *v.Note != "" {
				row += 1
			}
		}
		reffJpm[*v.ReffNumber] = counterRef
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
	// f.SetCellFormula(sheet, fmt.Sprintf("J%d", row+2), fmt.Sprintf("=I%d-J%d", row+1, row+1))

	f.SetCellStyle(sheet, fmt.Sprintf("D%d", row), fmt.Sprintf("J%d", row), stylingSubTotal)
	f.SetCellStyle(sheet, fmt.Sprintf("G%d", row+1), fmt.Sprintf("J%d", row+1), stylingSubTotal2)

	f.SetDefaultFont("Arial")

	return f, nil
}

func (s *service) ExportJcteWorksheet(ctx *abstraction.Context, consolidationData *model.ConsolidationEntityModel, f *excelize.File) (*excelize.File, error) {
	datas, err := s.JcteRepository.ExportAll(ctx, &consolidationData.ID)
	if err != nil {
		return nil, err
	}

	datePeriod, err := time.Parse(time.RFC3339, consolidationData.Period)
	if err != nil {
		return nil, err
	}

	sheet := "JCTE"
	f.NewSheet(sheet)

	arrStyleColWidth := []map[string]interface{}{
		{"COL": "A", "WIDTH": 6.43},
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

	err = f.MergeCell(sheet, "A6", "A7")
	if err != nil {
		return nil, err
	}
	err = f.MergeCell(sheet, "B6", "B7")
	if err != nil {
		return nil, err
	}
	err = f.MergeCell(sheet, "C6", "C7")
	if err != nil {
		return nil, err
	}
	err = f.MergeCell(sheet, "D6", "E7")
	if err != nil {
		return nil, err
	}

	err = f.MergeCell(sheet, "F6", "F7")
	if err != nil {
		return nil, err
	}

	err = f.MergeCell(sheet, "G6", "H6")
	if err != nil {
		return nil, err
	}

	err = f.MergeCell(sheet, "I6", "j6")
	if err != nil {
		return nil, err
	}

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
	f.SetCellValue(sheet, "B2", ": "+consolidationData.Company.Name)
	f.SetCellValue(sheet, "A3", "Date")
	f.SetCellValue(sheet, "B3", ": "+datePeriod.Format("02-Jan-06"))
	f.SetCellValue(sheet, "A4", "Subject")
	f.SetCellValue(sheet, "B4", ": Cost to Equity")

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
	reffJcte := make(map[string]int)
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

		if reffJcte[*v.ReffNumber] == 0 {
			row += 1
			if *v.Note != "" {
				row += 1
			}
		}
		if reffJcte[*v.ReffNumber] != 0 {
			row += 2
			if *v.Note != "" {
				row += 1
			}
		}
		reffJcte[*v.ReffNumber] = counterRef
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
	// f.SetCellFormula(sheet, fmt.Sprintf("J%d", row+2), fmt.Sprintf("=I%d-J%d", row+1, row+1))

	f.SetCellStyle(sheet, fmt.Sprintf("D%d", row), fmt.Sprintf("J%d", row), stylingSubTotal)
	f.SetCellStyle(sheet, fmt.Sprintf("G%d", row+1), fmt.Sprintf("J%d", row+1), stylingSubTotal2)

	f.SetDefaultFont("Arial")

	return f, nil
}

func (s *service) ExportJelimWorksheet(ctx *abstraction.Context, consolidationData *model.ConsolidationEntityModel, f *excelize.File) (*excelize.File, error) {
	datas, err := s.JelimRepository.ExportAll(ctx, &consolidationData.ID)
	if err != nil {
		return nil, err
	}

	datePeriod, err := time.Parse(time.RFC3339, consolidationData.Period)
	if err != nil {
		return nil, err
	}

	// if len(*datas) == 0 {
	// 	errs = append(errs, errors.New("Data Not Found"))
	// 	log.Println(err)
	// 	return
	// }

	sheet := "JELIM"
	f.NewSheet(sheet)

	arrStyleColWidth := []map[string]interface{}{
		{"COL": "A", "WIDTH": 6.43},
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

	err = f.MergeCell(sheet, "A6", "A7")
	if err != nil {
		return nil, err
	}
	err = f.MergeCell(sheet, "B6", "B7")
	if err != nil {
		return nil, err
	}
	err = f.MergeCell(sheet, "C6", "C7")
	if err != nil {
		return nil, err
	}
	err = f.MergeCell(sheet, "D6", "E7")
	if err != nil {
		return nil, err
	}

	err = f.MergeCell(sheet, "F6", "F7")
	if err != nil {
		return nil, err
	}

	err = f.MergeCell(sheet, "G6", "H6")
	if err != nil {
		return nil, err
	}

	err = f.MergeCell(sheet, "I6", "j6")
	if err != nil {
		return nil, err
	}

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
	f.SetCellValue(sheet, "B2", ": "+consolidationData.Company.Name)
	f.SetCellValue(sheet, "A3", "Date")
	f.SetCellValue(sheet, "B3", ": "+datePeriod.Format("02-Jan-06"))
	f.SetCellValue(sheet, "A4", "Subject")
	f.SetCellValue(sheet, "B4", ": Elimination Journal Entries")

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
	reffJelim := make(map[string]int)
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

		if reffJelim[*v.ReffNumber] == 0 {
			row += 1
			if *v.Note != "" {
				row += 1
			}
		}
		if reffJelim[*v.ReffNumber] != 0 {
			row += 2
			if *v.Note != "" {
				row += 1
			}
		}
		reffJelim[*v.ReffNumber] = counterRef
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
	// f.SetCellFormula(sheet, fmt.Sprintf("J%d", row+2), fmt.Sprintf("=I%d-J%d", row+1, row+1))

	f.SetCellStyle(sheet, fmt.Sprintf("D%d", row), fmt.Sprintf("J%d", row), stylingSubTotal)
	f.SetCellStyle(sheet, fmt.Sprintf("G%d", row+1), fmt.Sprintf("J%d", row+1), stylingSubTotal2)

	f.SetDefaultFont("Arial")
	return f, nil
}

func (s *service) ExportJpm(ctx *abstraction.Context, payload *abstraction.JsonData, wg *sync.WaitGroup, tmpFolder string) {
	defer wg.Done()
	var (
		criteria       model.JpmFilterModel
		criteriaConsol model.ConsolidationFilterModel
	)

	criteriaConsol.CompanyID = &payload.CompanyID
	criteriaConsol.ConsolidationVersions = &payload.Filter.Versions
	criteriaConsol.Period = &payload.Filter.Period

	consolidationData, err := s.ConsolidationRepository.FindByCriteria(ctx, &criteriaConsol)
	if err != nil {
		errs["JPM"] = "Error: Data Not Found"
		log.Println(err)
		return
	}

	criteria.CompanyID = &consolidationData.CompanyID
	criteria.Period = &consolidationData.Period
	datas, err := s.JpmRepository.ExportOne(ctx, &criteria)
	if err != nil {
		errs["JPM"] = "Error: Data Not Found"
		log.Println(err)
		return
	}

	datePeriod, err := time.Parse(time.RFC3339, datas.Period)
	if err != nil {
		errs["JPM"] = "Error: Data Not Found"
		log.Println(err)
		return
	}

	f := excelize.NewFile()
	sheet := f.GetSheetName(f.GetActiveSheetIndex())
	arrStyleColWidth := []map[string]interface{}{
		{"COL": "A", "WIDTH": 6.43},
		{"COL": "B", "WIDTH": 7.25},
		{"COL": "C", "WIDTH": 2.14},
		{"COL": "D", "WIDTH": 57.45},
		{"COL": "E", "WIDTH": 7.25},
		{"COL": "F", "WIDTH": 15.38},
		{"COL": "G", "WIDTH": 15.38},
		{"COL": "H", "WIDTH": 15.38},
		{"COL": "I", "WIDTH": 15.38},
	}

	for _, v := range arrStyleColWidth {
		tmpColWidth := fmt.Sprintf("%f", v["WIDTH"])
		colWidth, _ := strconv.ParseFloat(tmpColWidth, 64)
		err = f.SetColWidth(sheet, fmt.Sprintf("%s", v["COL"]), fmt.Sprintf("%s", v["COL"]), colWidth)
		if err != nil {
			errs["JPM"] = "Error: Generating Configuration File Excel"
			log.Println(err)
			return
		}
	}

	stylingBorderLROnly, err := f.NewStyle(&excelize.Style{
		Border: []excelize.Border{
			{Type: "right", Color: "000000", Style: 2},
			{Type: "left", Color: "000000", Style: 2},
		},
	})
	if err != nil {
		errs["JPM"] = "Error: Generating Configuration File Excel"
		log.Println(err)
		return
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
		errs["JPM"] = "Error: Generating Configuration File Excel"
		log.Println(err)
		return
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
		errs["JPM"] = "Error: Generating Configuration File Excel"
		log.Println(err)
		return
	}

	stylingCurrency, err := f.NewStyle(&excelize.Style{
		NumFmt: 7,
		Border: []excelize.Border{
			{Type: "left", Color: "000000", Style: 1},
			{Type: "right", Color: "000000", Style: 1},
		},
	})
	if err != nil {
		errs["JPM"] = "Error: Generating Configuration File Excel"
		log.Println(err)
		return
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
		errs["JPM"] = "Error: Generating Configuration File Excel"
		log.Println(err)
		return
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
		errs["JPM"] = "Error: Generating Configuration File Excel"
		log.Println(err)
		return
	}

	f.MergeCell(sheet, "A6", "A7")
	f.MergeCell(sheet, "B6", "B7")
	f.MergeCell(sheet, "E6", "E7")
	f.MergeCell(sheet, "C6", "D7")
	f.MergeCell(sheet, "F6", "G6")
	f.MergeCell(sheet, "H6", "I6")

	f.SetCellValue(sheet, "A2", "Company")
	f.SetCellValue(sheet, "B2", ": "+datas.Company.Name)
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

	rowBefore := row

	for i, v := range datas.JpmDetail {

		f.SetCellValue(sheet, fmt.Sprintf("A%d", row), i+1)
		f.SetCellValue(sheet, fmt.Sprintf("B%d", row), v.CoaCode)
		f.SetCellValue(sheet, fmt.Sprintf("E%d", row), *v.Description)
		f.SetCellValue(sheet, fmt.Sprintf("E%d", row+1), *v.Note)
		f.SetCellValue(sheet, fmt.Sprintf("G%d", row), *v.BalanceSheetDr)
		f.SetCellValue(sheet, fmt.Sprintf("H%d", row), *v.BalanceSheetCr)
		f.SetCellValue(sheet, fmt.Sprintf("I%d", row), *v.IncomeStatementDr)
		f.SetCellValue(sheet, fmt.Sprintf("J%d", row), *v.IncomeStatementCr)

		row += 2
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
	err = f.SaveAs(fmt.Sprintf("%s/JurnalProformaModal.xlsx", tmpFolder))
	if err != nil {
		errs["JPM"] = "Error: Saving File Excel"
		log.Println(err)
		return
	}
}

func (s *service) ExportAje(ctx *abstraction.Context, payload *abstraction.JsonData, wg *sync.WaitGroup, tmpFolder string) {
	defer wg.Done()
	var (
		criteria   model.AdjustmentFilterModel
		criteriaTB model.TrialBalanceFilterModel
	)

	criteriaTB.CompanyID = &payload.CompanyID
	criteriaTB.Versions = &payload.Filter.Versions
	criteriaTB.Period = &payload.Filter.Period

	dataTB, err := s.TrialBalanceRepository.Find(ctx, &criteriaTB)
	if err != nil {
		errs["JPM"] = "Error: Data Not Found"
		log.Println(err)
		return
	}

	if len(*dataTB) == 0 {
		errs["JPM"] = "Error: Data Not Found"
		log.Println(err)
		return
	}

	for _, v := range *dataTB {
		criteria.TrialBalanceID = &v.ID
	}

	criteria.CompanyID = &payload.CompanyID
	criteria.Period = &payload.Filter.Period
	datas, err := s.AjeRepository.Export(ctx, &criteria)
	if err != nil {
		errs["JPM"] = "Error: Data Not Found"
		log.Println(err)
		return
	}

	datePeriod, err := time.Parse(time.RFC3339, datas.Period)
	if err != nil {
		errs["JPM"] = "Error: Invalid Date Perio"
		log.Println(err)
		return
	}

	// if len(*datas) == 0 {
	// 	errs = append(errs, errors.New("Data Not Found"))
	// 	log.Println(err)
	// 	return
	// }

	f := excelize.NewFile()
	sheet := f.GetSheetName(f.GetActiveSheetIndex())
	arrStyleColWidth := []map[string]interface{}{
		{"COL": "A", "WIDTH": 6.43},
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
		err = f.SetColWidth(sheet, fmt.Sprintf("%s", v["COL"]), fmt.Sprintf("%s", v["COL"]), colWidth)
		if err != nil {
			errs["JPM"] = "Error: Generating Configuration File Excel"
			log.Println(err)
			return
		}
	}

	styleDefault, err := f.NewStyle(&excelize.Style{
		Font: &excelize.Font{
			Family: "Arial",
			Size:   10,
		},
	})
	if err != nil {
		errs["JPM"] = "Error: Generating Configuration File Excel"
		log.Println(err)
		return
	}
	err = f.SetColStyle(sheet, "A:Z", styleDefault)
	if err != nil {
		errs["JPM"] = "Error: Generating Configuration File Excel"
		log.Println(err)
		return
	}

	stylingBorderLROnly, err := f.NewStyle(&excelize.Style{
		Border: []excelize.Border{
			{Type: "right", Color: "000000", Style: 2},
			{Type: "left", Color: "000000", Style: 2},
		},
	})
	if err != nil {
		errs["JPM"] = "Error: Generating Configuration File Excel"
		log.Println(err)
		return
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
		errs["JPM"] = "Error: Generating Configuration File Excel"
		log.Println(err)
		return
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
		errs["JPM"] = "Error: Generating Configuration File Excel"
		log.Println(err)
		return
	}

	err = f.MergeCell(sheet, "A6", "A7")
	if err != nil {
		errs["JPM"] = "Error: Generating Configuration File Excel"
		log.Println(err)
		return
	}
	err = f.MergeCell(sheet, "B6", "B7")
	if err != nil {
		errs["JPM"] = "Error: Generating Configuration File Excel"
		log.Println(err)
		return
	}
	err = f.MergeCell(sheet, "C6", "C7")
	if err != nil {
		errs["JPM"] = "Error: Generating Configuration File Excel"
		log.Println(err)
		return
	}
	err = f.MergeCell(sheet, "D6", "E7")
	if err != nil {
		errs["JPM"] = "Error: Generating Configuration File Excel"
		log.Println(err)
		return
	}

	err = f.MergeCell(sheet, "F6", "F7")
	if err != nil {
		errs["JPM"] = "Error: Generating Configuration File Excel"
		log.Println(err)
		return
	}

	err = f.MergeCell(sheet, "G6", "H6")
	if err != nil {
		errs["JPM"] = "Error: Generating Configuration File Excel"
		log.Println(err)
		return
	}

	err = f.MergeCell(sheet, "I6", "j6")
	if err != nil {
		errs["JPM"] = "Error: Generating Configuration File Excel"
		log.Println(err)
		return
	}

	stylingCurrency, err := f.NewStyle(&excelize.Style{
		NumFmt: 7,
		Border: []excelize.Border{
			{Type: "left", Color: "000000", Style: 1},
			{Type: "right", Color: "000000", Style: 1},
		},
	})
	if err != nil {
		errs["JPM"] = "Error: Generating Configuration File Excel"
		log.Println(err)
		return
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
		errs["JPM"] = "Error: Generating Configuration File Excel"
		log.Println(err)
		return
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
		errs["JPM"] = "Error: Generating Configuration File Excel"
		log.Println(err)
		return
	}

	err = f.SetCellStyle(sheet, "A6", "J7", stylingHeader)
	if err != nil {
		errs["JPM"] = "Error: Generating Configuration File Excel"
		log.Println(err)
		return
	}
	f.SetCellValue(sheet, "A2", "Company")
	f.SetCellValue(sheet, "B2", ": "+datas.Company.Name)
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
	rowBefore := row
	for i, v := range datas.AdjustmentDetail {

		f.SetCellValue(sheet, fmt.Sprintf("A%d", row), i+1)
		f.SetCellValue(sheet, fmt.Sprintf("B%d", row), v.CoaCode)
		f.SetCellValue(sheet, fmt.Sprintf("E%d", row), *v.Description)
		f.SetCellValue(sheet, fmt.Sprintf("E%d", row+1), *v.Note)
		f.SetCellValue(sheet, fmt.Sprintf("G%d", row), *v.BalanceSheetDr)
		f.SetCellValue(sheet, fmt.Sprintf("H%d", row), *v.BalanceSheetCr)
		f.SetCellValue(sheet, fmt.Sprintf("I%d", row), *v.IncomeStatementDr)
		f.SetCellValue(sheet, fmt.Sprintf("J%d", row), *v.IncomeStatementCr)

		row += 2
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

	period := datePeriod.Format("2006-01-02")
	err = f.SaveAs(fmt.Sprintf("assets/%d/Adjustment_%s.xlsx", ctx.Auth.ID, period))
	if err != nil {
		errs["JPM"] = "Error: Saving File Excel"
		log.Println(err)
		return
	}
}

func (s *service) ExportJcte(ctx *abstraction.Context, payload *abstraction.JsonData, wg *sync.WaitGroup, tmpFolder string) {
	defer wg.Done()
	var (
		criteria       model.JcteFilterModel
		criteriaConsol model.ConsolidationFilterModel
	)

	criteriaConsol.CompanyID = &payload.CompanyID
	criteriaConsol.ConsolidationVersions = &payload.Filter.Versions
	criteriaConsol.Period = &payload.Filter.Period
	consolidationData, err := s.ConsolidationRepository.FindByCriteria(ctx, &criteriaConsol)
	if err != nil {
		errs["JCTE"] = "Error: Data Not Found"
		log.Println(err)
		return
	}

	criteria.CompanyID = &consolidationData.CompanyID
	criteria.Period = &consolidationData.Period
	criteria.ConsolidationID = &consolidationData.ID
	datas, err := s.JcteRepository.ExportOne(ctx, &criteria)
	if err != nil {
		errs["JCTE"] = "Error: Data Not Found"
		log.Println(err)
		return
	}

	datePeriod, err := time.Parse(time.RFC3339, datas.Period)
	if err != nil {
		errs["JCTE"] = "Error: Invalid Date Period"
		log.Println(err)
		return
	}

	f := excelize.NewFile()
	sheet := f.GetSheetName(f.GetActiveSheetIndex())
	arrStyleColWidth := []map[string]interface{}{
		{"COL": "A", "WIDTH": 6.43},
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
		err = f.SetColWidth(sheet, fmt.Sprintf("%s", v["COL"]), fmt.Sprintf("%s", v["COL"]), colWidth)
		if err != nil {
			errs["JCTE"] = "Error: Generating Configuration File Excel"
			log.Println(err)
			return
		}
	}

	styleDefault, err := f.NewStyle(&excelize.Style{
		Font: &excelize.Font{
			Family: "Arial",
			Size:   10,
		},
	})
	if err != nil {
		errs["JCTE"] = "Error: Generating Configuration File Excel"
		log.Println(err)
		return
	}
	err = f.SetColStyle(sheet, "A:Z", styleDefault)
	if err != nil {
		errs["JCTE"] = "Error: Generating Configuration File Excel"
		log.Println(err)
		return
	}

	stylingBorderLROnly, err := f.NewStyle(&excelize.Style{
		Border: []excelize.Border{
			{Type: "right", Color: "000000", Style: 2},
			{Type: "left", Color: "000000", Style: 2},
		},
	})
	if err != nil {
		errs["JCTE"] = "Error: Generating Configuration File Excel"
		log.Println(err)
		return
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
		errs["JCTE"] = "Error: Generating Configuration File Excel"
		log.Println(err)
		return
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
		errs["JCTE"] = "Error: Generating Configuration File Excel"
		log.Println(err)
		return
	}

	err = f.MergeCell(sheet, "A6", "A7")
	if err != nil {
		errs["JCTE"] = "Error: Generating Configuration File Excel"
		log.Println(err)
		return
	}
	err = f.MergeCell(sheet, "B6", "B7")
	if err != nil {
		errs["JCTE"] = "Error: Generating Configuration File Excel"
		log.Println(err)
		return
	}
	err = f.MergeCell(sheet, "C6", "C7")
	if err != nil {
		errs["JCTE"] = "Error: Generating Configuration File Excel"
		log.Println(err)
		return
	}
	err = f.MergeCell(sheet, "D6", "E7")
	if err != nil {
		errs["JCTE"] = "Error: Generating Configuration File Excel"
		log.Println(err)
		return
	}

	err = f.MergeCell(sheet, "F6", "F7")
	if err != nil {
		errs["JCTE"] = "Error: Generating Configuration File Excel"
		log.Println(err)
		return
	}

	err = f.MergeCell(sheet, "G6", "H6")
	if err != nil {
		errs["JCTE"] = "Error: Generating Configuration File Excel"
		log.Println(err)
		return
	}

	err = f.MergeCell(sheet, "I6", "j6")
	if err != nil {
		errs["JCTE"] = "Error: Generating Configuration File Excel"
		log.Println(err)
		return
	}

	stylingCurrency, err := f.NewStyle(&excelize.Style{
		NumFmt: 7,
		Border: []excelize.Border{
			{Type: "left", Color: "000000", Style: 1},
			{Type: "right", Color: "000000", Style: 1},
		},
	})
	if err != nil {
		errs["JCTE"] = "Error: Generating Configuration File Excel"
		log.Println(err)
		return
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
		errs["JCTE"] = "Error: Generating Configuration File Excel"
		log.Println(err)
		return
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
		errs["JCTE"] = "Error: Generating Configuration File Excel"
		log.Println(err)
		return
	}

	f.SetCellValue(sheet, "A2", "Company")
	f.SetCellValue(sheet, "B2", ": "+datas.Company.Name)
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
	rowBefore := row
	for i, v := range datas.JcteDetail {

		f.SetCellValue(sheet, fmt.Sprintf("A%d", row), i+1)
		f.SetCellValue(sheet, fmt.Sprintf("B%d", row), v.CoaCode)
		f.SetCellValue(sheet, fmt.Sprintf("E%d", row), *v.Description)
		f.SetCellValue(sheet, fmt.Sprintf("E%d", row+1), *v.Note)
		f.SetCellValue(sheet, fmt.Sprintf("G%d", row), *v.BalanceSheetDr)
		f.SetCellValue(sheet, fmt.Sprintf("H%d", row), *v.BalanceSheetCr)
		f.SetCellValue(sheet, fmt.Sprintf("I%d", row), *v.IncomeStatementDr)
		f.SetCellValue(sheet, fmt.Sprintf("J%d", row), *v.IncomeStatementCr)

		row += 2
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
	period := datePeriod.Format("2006-01-02")
	err = f.SaveAs(fmt.Sprintf("%s/JurnalCTE_%s.xlsx", tmpFolder, period))
	if err != nil {
		errs["JCTE"] = "Error: Saving File Excel"
		log.Println(err)
		return
	}
}

func (s *service) ExportJelim(ctx *abstraction.Context, payload *abstraction.JsonData, wg *sync.WaitGroup, tmpFolder string) {
	defer wg.Done()
	var (
		criteria       model.JelimFilterModel
		criteriaConsol model.ConsolidationFilterModel
	)

	criteriaConsol.CompanyID = &payload.CompanyID
	criteriaConsol.ConsolidationVersions = &payload.Filter.Versions
	criteriaConsol.Period = &payload.Filter.Period
	consolidationData, err := s.ConsolidationRepository.FindByCriteria(ctx, &criteriaConsol)
	if err != nil {
		errs["JELIM"] = "Error: Data Not Found"
		log.Println(err)
		return
	}

	criteria.CompanyID = &consolidationData.CompanyID
	criteria.Period = &consolidationData.Period
	criteria.ConsolidationID = &consolidationData.ID
	datas, err := s.JelimRepository.ExportOne(ctx, &criteria)
	if err != nil {
		errs["JELIM"] = "Error: Data Not Found"
		log.Println(err)
		return
	}

	datePeriod, err := time.Parse(time.RFC3339, datas.Period)
	if err != nil {
		errs["JELIM"] = "Error: Invalid Date Period"
		log.Println(err)
		return
	}

	// if len(*datas) == 0 {
	// 	errs = append(errs, errors.New("Data Not Found"))
	// 	log.Println(err)
	// 	return
	// }

	f := excelize.NewFile()
	sheet := f.GetSheetName(f.GetActiveSheetIndex())
	arrStyleColWidth := []map[string]interface{}{
		{"COL": "A", "WIDTH": 6.43},
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
		err = f.SetColWidth(sheet, fmt.Sprintf("%s", v["COL"]), fmt.Sprintf("%s", v["COL"]), colWidth)
		if err != nil {
			errs["JELIM"] = "Error: Generating Configuration File Excel"
			log.Println(err)
			return
		}
	}

	styleDefault, err := f.NewStyle(&excelize.Style{
		Font: &excelize.Font{
			Family: "Arial",
			Size:   10,
		},
	})
	if err != nil {
		errs["JELIM"] = "Error: Generating Configuration File Excel"
		log.Println(err)
		return
	}
	err = f.SetColStyle(sheet, "A:Z", styleDefault)
	if err != nil {
		errs["JELIM"] = "Error: Generating Configuration File Excel"
		log.Println(err)
		return
	}

	stylingBorderLROnly, err := f.NewStyle(&excelize.Style{
		Border: []excelize.Border{
			{Type: "right", Color: "000000", Style: 2},
			{Type: "left", Color: "000000", Style: 2},
		},
	})
	if err != nil {
		errs["JELIM"] = "Error: Generating Configuration File Excel"
		log.Println(err)
		return
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
		errs["JELIM"] = "Error: Generating Configuration File Excel"
		log.Println(err)
		return
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
		errs["JELIM"] = "Error: Generating Configuration File Excel"
		log.Println(err)
		return
	}

	err = f.MergeCell(sheet, "A6", "A7")
	if err != nil {
		errs["JELIM"] = "Error: Generating Configuration File Excel"
		log.Println(err)
		return
	}
	err = f.MergeCell(sheet, "B6", "B7")
	if err != nil {
		errs["JELIM"] = "Error: Generating Configuration File Excel"
		log.Println(err)
		return
	}
	err = f.MergeCell(sheet, "C6", "C7")
	if err != nil {
		errs["JELIM"] = "Error: Generating Configuration File Excel"
		log.Println(err)
		return
	}
	err = f.MergeCell(sheet, "D6", "E7")
	if err != nil {
		errs["JELIM"] = "Error: Generating Configuration File Excel"
		log.Println(err)
		return
	}

	err = f.MergeCell(sheet, "F6", "F7")
	if err != nil {
		errs["JELIM"] = "Error: Generating Configuration File Excel"
		log.Println(err)
		return
	}

	err = f.MergeCell(sheet, "G6", "H6")
	if err != nil {
		errs["JELIM"] = "Error: Generating Configuration File Excel"
		log.Println(err)
		return
	}

	err = f.MergeCell(sheet, "I6", "j6")
	if err != nil {
		errs["JELIM"] = "Error: Generating Configuration File Excel"
		log.Println(err)
		return
	}

	stylingCurrency, err := f.NewStyle(&excelize.Style{
		NumFmt: 7,
		Border: []excelize.Border{
			{Type: "left", Color: "000000", Style: 1},
			{Type: "right", Color: "000000", Style: 1},
		},
	})
	if err != nil {
		errs["JELIM"] = "Error: Generating Configuration File Excel"
		log.Println(err)
		return
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
		errs["JELIM"] = "Error: Generating Configuration File Excel"
		log.Println(err)
		return
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
		errs["JELIM"] = "Error: Generating Configuration File Excel"
		log.Println(err)
		return
	}

	f.SetCellValue(sheet, "A2", "Company")
	f.SetCellValue(sheet, "B2", ": "+datas.Company.Name)
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
	rowBefore := row
	for i, v := range datas.JelimDetail {

		f.SetCellValue(sheet, fmt.Sprintf("A%d", row), i+1)
		f.SetCellValue(sheet, fmt.Sprintf("B%d", row), v.CoaCode)
		f.SetCellValue(sheet, fmt.Sprintf("E%d", row), *v.Description)
		f.SetCellValue(sheet, fmt.Sprintf("E%d", row+1), *v.Note)
		f.SetCellValue(sheet, fmt.Sprintf("G%d", row), *v.BalanceSheetDr)
		f.SetCellValue(sheet, fmt.Sprintf("H%d", row), *v.BalanceSheetCr)
		f.SetCellValue(sheet, fmt.Sprintf("I%d", row), *v.IncomeStatementDr)
		f.SetCellValue(sheet, fmt.Sprintf("J%d", row), *v.IncomeStatementCr)

		row += 2
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
	period := datePeriod.Format("2006-01-02")
	err = f.SaveAs(fmt.Sprintf("%s/JurnalEliminasi_%s.xlsx", tmpFolder, period))
	if err != nil {
		errs["JELIM"] = "Error: Saving Excel File"
		log.Println(err)
		return
	}
}

func (s *service) ExportEmployeeBenefit(ctx *abstraction.Context, payload *abstraction.JsonData, wg *sync.WaitGroup, tmpFolder string) {
	defer wg.Done()
	f := excelize.NewFile()
	sheet := f.GetSheetName(f.GetActiveSheetIndex())
	f.SetColWidth(sheet, "A", "A", 4.29)
	f.SetColWidth(sheet, "B", "B", 3.57)
	f.SetColWidth(sheet, "C", "E", 8.39)
	f.SetColWidth(sheet, "F", "F", 21.60)
	f.SetColWidth(sheet, "G", "G", 19.10)

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
		errs["EB"] = "Error: Generating Configuration File"
		log.Println(err)
		return
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
		errs["EB"] = "Error: Generating Configuration File Excel"
		log.Println(err)
		return
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
		errs["EB"] = "Error: Generating Configuration File Excel"
		log.Println(err)
		return
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
		errs["EB"] = "Error: Generating Configuration File Excel"
		log.Println(err)
		return
	}

	styleLabel, err := f.NewStyle(&excelize.Style{
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
		errs["EB"] = "Error: Generating Configuration File Excel"
		log.Println(err)
		return
	}

	// styleLabel, err := f.NewStyle(&excelize.Style{
	// 	Border: []excelize.Border{
	// 		{Type: "top", Color: "000000", Style: 1},
	// 		{Type: "bottom", Color: "000000", Style: 1},
	// 		{Type: "left", Color: "000000", Style: 1},
	// 		{Type: "right", Color: "000000", Style: 1},
	// 	},
	// })
	// if err != nil {
	// 	errs = append(errs, err)
	// 	log.Println(err)
	// 	return
	// }

	criteria := model.EmployeeBenefitFilterModel{}
	criteria.CompanyID = &payload.CompanyID
	criteria.Versions = &payload.Filter.Versions
	criteria.Period = &payload.Filter.Period
	// criteria.FormatterID = &data.ID

	employeeBenefit, err := s.EmployeeBenefitRepository.Find(ctx, &criteria)
	if err != nil {
		errs["EB"] = "Error: Data Not Found"
		log.Println(err)
		return
	}

	datePeriod, err := time.Parse("2006-01-02", payload.Filter.Period)
	if err != nil {
		errs["EB"] = "Error: Invalid Date Period"
		log.Println(err)
		return
	}

	formatterCode := []string{"EMPLOYEE-BENEFIT-ASUMSI", "EMPLOYEE-BENEFIT-REKONSILIASI", "EMPLOYEE-BENEFIT-RINCIAN-LAPORAN", "EMPLOYEE-BENEFIT-RINCIAN-EKUITAS", "EMPLOYEE-BENEFIT-MUTASI", "EMPLOYEE-BENEFIT-INFORMASI", "EMPLOYEE-BENEFIT-ANALISIS"}
	formatterTitle := []string{"Asumsi-asumsi yang digunakan:", "Rekonsiliasi jumlah liabilitas imbalan kerja karyawan pada laporan posisi keuangan adalah sebagai berikut:", "Rincian beban imbalan kerja karyawan yang diakui dalam laporan laba rugi dan penghasilan komprehensif lain adalah sebagai berikut:", "Rincian beban imbalan kerja karyawan yang diakui pada ekuitas dalam penghasilan komprehensif lain adalah sebagai berikut:", "Mutasi liabilitas imbalan kerja karyawan adalah sebagai berikut:", "Informasi historis dari nilai kini liabilitas imbalan pasti, nilai wajar aset program dan penyesuaian adalah sebagai berikut:", "Analisis sensitivitas dari perubahan asumsi-asumsi utama terhadap liabilitas imbalan kerja"}
	row, rowStart := 7, 7
	rowCode := make(map[string]int)
	for i, formatter := range formatterCode {
		f.SetCellValue(sheet, fmt.Sprintf("A%d", (rowStart-3)), formatterTitle[i])
		f.MergeCell(sheet, fmt.Sprintf("B%d", (rowStart-2)), fmt.Sprintf("F%d", (rowStart-1)))
		f.SetCellStyle(sheet, fmt.Sprintf("B%d", (rowStart-2)), fmt.Sprintf("G%d", (rowStart-1)), styleHeader)
		f.SetCellStyle(sheet, fmt.Sprintf("A%d", (rowStart-3)), fmt.Sprintf("A%d", (rowStart-3)), styleLabel)

		f.SetCellValue(sheet, fmt.Sprintf("B%d", (rowStart-2)), "Description")
		f.SetCellValue(sheet, fmt.Sprintf("G%d", (rowStart-2)), datePeriod.Format("02-Jan-06"))
		for _, vEB := range *employeeBenefit {
			f.SetCellValue(sheet, fmt.Sprintf("G%d", (rowStart-1)), vEB.Company.Name)
		}

		var criteria model.FormatterFilterModel
		criteria.FormatterFor = &formatter

		data, err := s.FormatterRepository.FindWithDetail(ctx, &criteria)
		if err != nil {
			errs["EB"] = "Error: Formatter Data for Jurnal Elimination Not Found"
			log.Println(err)
			return
		}

		tmpStr := "EMPLOYEE-BENEFIT"
		criteriaBridge := model.FormatterBridgesFilterModel{}
		criteriaBridge.FormatterBridgesFilter.Source = &tmpStr
		criteriaBridge.FormatterBridgesFilter.FormatterID = &data.ID
		for _, valaup := range *employeeBenefit {
			criteriaBridge.FormatterBridgesFilter.TrxRefID = &valaup.ID
		}

		bridges, err := s.FormatterBridgesRepository.FindWithCriteria(ctx, &criteriaBridge)
		if err != nil {
			errs["EB"] = "Error: Data Not Found"
			log.Println(err)
			return
		}

		partRowStart := row
		for _, v := range data.FormatterDetail {
			// f.SetCellStyle(sheet, fmt.Sprintf("B%d", row), fmt.Sprintf("F%d", row), styleLabel)
			if strings.Contains(strings.ToLower(v.Code), "blank") {
				row++
				continue
			}
			if rowCode[v.Code] == 0 {
				rowCode[v.Code] = row
			}
			f.SetCellValue(sheet, fmt.Sprintf("B%d", row), v.Description)

			if v.IsTotal != nil && *v.IsTotal {
				// f.SetCellStyle(sheet, fmt.Sprintf("B%d", row), fmt.Sprintf("B%d", row), styleLabelTotal)
				f.SetCellStyle(sheet, fmt.Sprintf("G%d", row), fmt.Sprintf("G%d", row), styleCurrencyTotal)
				if v.FxSummary == "" {
					row++
					continue
				}
				formula := v.FxSummary
				reg := regexp.MustCompile(`[A-Za-z_~]+`)
				match := reg.FindAllString(formula, -1)
				for _, vMatch := range match {
					//cari jml berdasarkan code
					// if _, ok := rowCode[vMatch]; !ok {
					// 	formula = strings.ReplaceAll(formula, vMatch, fmt.Sprintf("G%d", rowCode[vMatch]))
					// }
					if rowCode[vMatch] != 0 {
						formula = strings.ReplaceAll(formula, vMatch, fmt.Sprintf("G%d", rowCode[vMatch]))
					}
				}
				f.SetCellFormula(sheet, fmt.Sprintf("G%d", row), fmt.Sprintf("=%s", formula))
				row++
				continue
			}

			if v.AutoSummary != nil && *v.AutoSummary {
				f.SetCellStyle(sheet, fmt.Sprintf("B%d", row), fmt.Sprintf("B%d", row), styleLabelTotal)
				f.SetCellStyle(sheet, fmt.Sprintf("G%d", row), fmt.Sprintf("G%d", row), styleCurrencyTotal)
				f.SetCellFormula(sheet, fmt.Sprintf("G%d", row), fmt.Sprintf("=SUM(G%d:G%d)", partRowStart, row-1))
				row++
				partRowStart = row
				continue
			}

			if v.IsLabel != nil && *v.IsLabel {
				row++
				continue
			}

			criteriaAUP := model.EmployeeBenefitDetailFilterModel{}
			criteriaAUP.Code = &v.Code
			criteriaAUP.FormatterBridgesID = &bridges.ID
			for _, valaup := range *employeeBenefit {
				criteriaAUP.EmployeeBenefitID = &valaup.ID
			}

			employeeBenefitDetail, err := s.EmployeeBenefitDetailRepository.Find(ctx, &criteriaAUP)
			if err != nil {
				errs["EB"] = "Error: Data Not Found"
				log.Println(err)
				return
			}
			if len(*employeeBenefitDetail) == 0 {
				errs["EB"] = "Error: Data Not Found"
				log.Println("Data Not Found")
				return
			}
			for _, vv := range *employeeBenefitDetail {
				if v.IsLabel == nil || !*v.IsLabel {
					if vv.IsValue != nil && *vv.IsValue {
						f.SetCellValue(sheet, fmt.Sprintf("G%d", row), vv.Value)
						continue
					}
					f.SetCellStyle(sheet, fmt.Sprintf("G%d", row), fmt.Sprintf("G%d", row), styleCurrency)
					amount := 0.0
					if vv.Amount != nil {
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
	period := datePeriod.Format("2006-01-02")
	err = f.SaveAs(fmt.Sprintf("%s/EmployeeBenefit_%s.xlsx", tmpFolder, period))
	if err != nil {
		errs["EB"] = "Error: Saving File Excel"
		log.Println(err)
		return
	}
	f.Close()
}

func (s *service) ExportMutasiFaNew(ctx *abstraction.Context, consolidationData *model.ConsolidationEntityModel, f *excelize.File) (*excelize.File, error) {
	sheet := "MUTASI FA"
	f.NewSheet(sheet)
	err := f.SetPanes(sheet, `{"freeze":true,"split":false,"x_split":1,"y_split":0,"top_left_cell":"B1","active_pane":"bottomRight","panes":[{"sqref":"B1","active_cell":"B2"}]}`)
	if err != nil {
		log.Fatal(err)
	}
	formatterCode := []string{"MUTASI-FA-COST", "MUTASI-FA-ACCUMULATED-DEPRECATION"}

	
	
	var CharCombineBeginingBalance []string
	for _, formatter := range formatterCode {

		var criteria model.FormatterFilterModel
		criteria.FormatterFor = &formatter

		// rowStart = row
		// row = rowStart
		rowIndex2 := 2
		rowColumn2 := 3
		rowPlus := 0
		CParent, _ := excelize.CoordinatesToCellName(rowIndex2, rowColumn2-1)
		DParent, _ := excelize.CoordinatesToCellName(rowIndex2, rowColumn2)
		rowIndex2 = rowIndex2 + 1
		EParent, _ := excelize.CoordinatesToCellName(rowIndex2, rowColumn2)
		rowIndex2 = rowIndex2 + 1
		FParent, _ := excelize.CoordinatesToCellName(rowIndex2, rowColumn2)
		rowIndex2 = rowIndex2 + 1
		GParent, _ := excelize.CoordinatesToCellName(rowIndex2, rowColumn2)
		rowIndex2 = rowIndex2 + 1
		HParent, _ := excelize.CoordinatesToCellName(rowIndex2, rowColumn2)
		rowIndex2 = rowIndex2 + 1
		IParent, _ := excelize.CoordinatesToCellName(rowIndex2, rowColumn2)
		rowIndex2 = rowIndex2 + 1
		JParent, _ := excelize.CoordinatesToCellName(rowIndex2, rowColumn2)
		rowIndex2 = rowIndex2 + 1
		// K4, _ := excelize.CoordinatesToCellName(rowIndex2, rowColumn2)
		// rowIndex2 = rowIndex2 + 1
		// rowIndex2 = rowIndex2+1

		f.SetCellValue(sheet, CParent, "Parent")
		f.SetCellValue(sheet, DParent, "Beginning Balance")
		f.SetCellValue(sheet, EParent, "Acquisition of Subsidiary")
		f.SetCellValue(sheet, FParent, "Additions (+)")
		f.SetCellValue(sheet, GParent, "Deductions (-)")
		f.SetCellValue(sheet, HParent, "Reclassification")
		f.SetCellValue(sheet, IParent, "Revaluation")
		f.SetCellValue(sheet, JParent, "Ending balance")
		rowColumnParent := 5
		var TotalCostParent string
		var TotalAccumParent string
		var WorkInProsesParent string
		var EndTotalCostParent string
		var EndTotalAccumParent string
		var EndWorkInProsesParent string
		filterMutasiFa := model.MutasiFaFilterModel{}
		filterMutasiFa.CompanyID = &consolidationData.CompanyID
		filterMutasiFa.Period = &consolidationData.Period
		filterMutasiFa.Versions = &consolidationData.Versions
			if consolidationData.Versions == 0 {
				filterMutasiFa.Versions = &consolidationData.ConsolidationVersions
			}
			mutasiFa, err := s.MutasiFaRepository.FindByCriteria(ctx, &filterMutasiFa)
			if err != nil {
				return nil, err
			}
			tmpStr := "MUTASI-FA"
			criteriaBridge := model.FormatterBridgesFilterModel{}
			criteriaBridge.FormatterBridgesFilter.Source = &tmpStr
			// criteriaBridge.FormatterBridgesFilter.FormatterID = &data.ID
			criteriaBridge.FormatterBridgesFilter.TrxRefID = &mutasiFa.ID

			bridges, err := s.FormatterBridgesRepository.FindWithCriteriaNew(ctx, &criteriaBridge)
			if err != nil {
				return nil, err
			}
			for _, brid := range *bridges {

				if brid.FormatterID == 21 {
					continue
				}
				criteriaMF := model.MutasiFaDetailFilterModel{}
				// criteriaMF.Code = &v.Code
				criteriaMF.FormatterBridgesID = &brid.ID
				criteriaMF.MutasiFaID = &mutasiFa.ID

				MutasiFaDetail, err := s.MutasiFaDetailRepository.Find(ctx, &criteriaMF)
				if err != nil {
					continue
				}
				
				for _, vv := range *MutasiFaDetail {
					rowIndex := 2
					
					if vv.Code == "ACCUMULATED_DEPRECIATION"  {
						rowColumnParent = rowColumnParent+1
					}
					
						BeginningBalance, _ := excelize.CoordinatesToCellName(rowIndex, rowColumnParent)
						f.SetCellValue(sheet, BeginningBalance, 0)
						if vv.Code != "" && vv.Code != "TOTAL_COST" && vv.Code != "TOTAL_ACCUMULATED_DEPRECIATION" && vv.Code != "WORK_IN_PROCESS" && vv.Code != "NET_BOOK_VALUE" {
							CharCombineBeginingBalance = append(CharCombineBeginingBalance,BeginningBalance)
						}

						rowIndex = rowIndex + 1
						AcquisitionOfSubsidiary, _ := excelize.CoordinatesToCellName(rowIndex, rowColumnParent)
						f.SetCellValue(sheet, AcquisitionOfSubsidiary, 0)

						rowIndex = rowIndex + 1
						Additions, _ := excelize.CoordinatesToCellName(rowIndex, rowColumnParent)
						f.SetCellValue(sheet, Additions, 0)

						rowIndex = rowIndex + 1
						Deductions, _ := excelize.CoordinatesToCellName(rowIndex, rowColumnParent)
						f.SetCellValue(sheet, Deductions, 0)

						rowIndex = rowIndex + 1
						Reclassification, _ := excelize.CoordinatesToCellName(rowIndex, rowColumnParent)
						f.SetCellValue(sheet, Reclassification, 0)

						rowIndex = rowIndex + 1
						Revaluation, _ := excelize.CoordinatesToCellName(rowIndex, rowColumnParent)
						f.SetCellValue(sheet, Revaluation, 0)
						rowIndex = rowIndex + 1
						EndingBalance, _ := excelize.CoordinatesToCellName(rowIndex, rowColumnParent)
						f.SetCellFormula(sheet, EndingBalance, fmt.Sprintf("=%s+%s+%s+%s+%s+%s", BeginningBalance, AcquisitionOfSubsidiary, Additions, Deductions, Reclassification, Revaluation))

						

						if vv.BeginningBalance != nil && *vv.BeginningBalance != 0 {
							f.SetCellValue(sheet, BeginningBalance, *vv.BeginningBalance)
						}
						if vv.AcquisitionOfSubsidiary != nil && *vv.AcquisitionOfSubsidiary != 0 {
							f.SetCellValue(sheet, AcquisitionOfSubsidiary, *vv.AcquisitionOfSubsidiary)
						}
						if vv.Additions != nil && *vv.Additions != 0 {
							f.SetCellValue(sheet, Additions, *vv.Additions)
						}
						if vv.Deductions != nil && *vv.Deductions != 0 {
							f.SetCellValue(sheet, Deductions, *vv.Deductions)
						}
						if vv.Reclassification != nil && *vv.Reclassification != 0 {
							f.SetCellValue(sheet, Reclassification, *vv.Reclassification)
						}
						if vv.Revaluation != nil && *vv.Revaluation != 0 {
							f.SetCellValue(sheet, Revaluation, *vv.Revaluation)
						}
				
						if vv.Code == ""  {
							if err := f.SetCellValue(sheet, BeginningBalance, ""); err != nil {
								log.Fatal(err)
							}
							if err := f.SetCellValue(sheet, AcquisitionOfSubsidiary, ""); err != nil {
								log.Fatal(err)
							}
							if err := f.SetCellValue(sheet, Additions, ""); err != nil {
								log.Fatal(err)
							}
							if err := f.SetCellValue(sheet, Deductions, ""); err != nil {
								log.Fatal(err)
							}
							if err := f.SetCellValue(sheet, Reclassification, ""); err != nil {
								log.Fatal(err)
							}
							if err := f.SetCellValue(sheet, Revaluation, ""); err != nil {
								log.Fatal(err)
							}
							if err := f.SetCellValue(sheet, EndingBalance, ""); err != nil {
								log.Fatal(err)
							}
						}
						if vv.Code == "TOTAL_COST" {

							if err := f.SetCellValue(sheet, BeginningBalance, ""); err != nil {
								log.Fatal(err)
							}
							var BeginningBalanceStartRow string
							var BeginningBalanceEndRow string
							if len(BeginningBalance) > 3 {
								BeginningBalanceStartRow = BeginningBalance[:2] + "6"
								BeginningBalanceEndRow = BeginningBalance[:2] + "12"
							} else {
								BeginningBalanceStartRow = BeginningBalance[:1] + "6"
								BeginningBalanceEndRow = BeginningBalance[:1] + "12"
							}
							rowColumnParent = rowColumnParent+1
							if len(BeginningBalance) > 3 {
								BeginningBalance = BeginningBalance[:2] + strconv.Itoa(rowColumnParent)
							} else {
								BeginningBalance = BeginningBalance[:1] + strconv.Itoa(rowColumnParent)
							}
						
							TotalCostParent = BeginningBalance
							f.SetCellFormula(sheet, BeginningBalance, fmt.Sprintf("=SUM(%s:%s)", BeginningBalanceStartRow, BeginningBalanceEndRow))
							CharCombineBeginingBalance = append(CharCombineBeginingBalance,BeginningBalance)
							rowIndex = rowIndex + 1
							if err := f.SetCellValue(sheet, AcquisitionOfSubsidiary, ""); err != nil {
								log.Fatal(err)
							}
						
							var AcquisitionOfSubsidiaryStartRow string
							var AcquisitionOfSubsidiaryEndRow string
							if len(AcquisitionOfSubsidiary) > 3 {
								AcquisitionOfSubsidiaryStartRow = AcquisitionOfSubsidiary[:2] + "6"
								AcquisitionOfSubsidiaryEndRow = AcquisitionOfSubsidiary[:2] + "12"
							} else {
								AcquisitionOfSubsidiaryStartRow = AcquisitionOfSubsidiary[:1] + "6"
								AcquisitionOfSubsidiaryEndRow = AcquisitionOfSubsidiary[:1] + "12"
							}
						
							if len(AcquisitionOfSubsidiary) > 3 {
								AcquisitionOfSubsidiary = AcquisitionOfSubsidiary[:2] + strconv.Itoa(rowColumnParent)
							} else {
								AcquisitionOfSubsidiary = AcquisitionOfSubsidiary[:1] + strconv.Itoa(rowColumnParent)
							}
						
							f.SetCellFormula(sheet, AcquisitionOfSubsidiary, fmt.Sprintf("=SUM(%s:%s)", AcquisitionOfSubsidiaryStartRow, AcquisitionOfSubsidiaryEndRow))
						
							rowIndex = rowIndex + 1
							if err := f.SetCellValue(sheet, Additions, ""); err != nil {
								log.Fatal(err)
							}
							var AdditionsStartRow string
							var AdditionsEndRow string
							if len(Additions) > 3 {
								AdditionsStartRow = Additions[:2] + "6"
								AdditionsEndRow = Additions[:2] + "12"
							} else {
								AdditionsStartRow = Additions[:1] + "6"
								AdditionsEndRow = Additions[:1] + "12"
							}
						
							if len(Additions) > 3 {
								Additions = Additions[:2] + strconv.Itoa(rowColumnParent)
							} else {
								Additions = Additions[:1] + strconv.Itoa(rowColumnParent)
							}
							f.SetCellFormula(sheet, Additions, fmt.Sprintf("=SUM(%s:%s)", AdditionsStartRow, AdditionsEndRow))
						
							rowIndex = rowIndex + 1
							if err := f.SetCellValue(sheet, Deductions, ""); err != nil {
								log.Fatal(err)
							}
							var DeductionsStartRow string
							var DeductionsEndRow string
							if len(Deductions) > 3 {
								DeductionsStartRow = Deductions[:2] + "6"
								DeductionsEndRow = Deductions[:2] + "12"
							} else {
								DeductionsStartRow = Deductions[:1] + "6"
								DeductionsEndRow = Deductions[:1] + "12"
							}
						
							if len(Deductions) > 3 {
								Deductions = Deductions[:2] + strconv.Itoa(rowColumnParent)
							} else {
								Deductions = Deductions[:1] + strconv.Itoa(rowColumnParent)
							}
							f.SetCellFormula(sheet, Deductions, fmt.Sprintf("=SUM(%s:%s)", DeductionsStartRow, DeductionsEndRow))
						
							rowIndex = rowIndex + 1
							if err := f.SetCellValue(sheet, Reclassification, ""); err != nil {
								log.Fatal(err)
							}
							var ReclassificationStartRow string
							var ReclassificationEndRow string
							if len(Reclassification) > 3 {
								ReclassificationStartRow = Reclassification[:2] + "6"
								ReclassificationEndRow = Reclassification[:2] + "12"
							} else {
								ReclassificationStartRow = Reclassification[:1] + "6"
								ReclassificationEndRow = Reclassification[:1] + "12"
							}
						
							if len(Reclassification) > 3 {
								Reclassification = Reclassification[:2] + strconv.Itoa(rowColumnParent)
							} else {
								Reclassification = Reclassification[:1] + strconv.Itoa(rowColumnParent)
							}
							f.SetCellFormula(sheet, Reclassification, fmt.Sprintf("=SUM(%s:%s)", ReclassificationStartRow, ReclassificationEndRow))
						
							rowIndex = rowIndex + 1
							if err := f.SetCellValue(sheet, Revaluation, ""); err != nil {
								log.Fatal(err)
							}
							var RevaluationStartRow string
							var RevaluationEndRow string
							if len(Revaluation) > 3 {
								RevaluationStartRow = Revaluation[:2] + "6"
								RevaluationEndRow = Revaluation[:2] + "12"
							} else {
								RevaluationStartRow = Revaluation[:1] + "6"
								RevaluationEndRow = Revaluation[:1] + "12"
							}
						
							if len(Revaluation) > 3 {
								Revaluation = Revaluation[:2] + strconv.Itoa(rowColumnParent)
							} else {
								Revaluation = Revaluation[:1] + strconv.Itoa(rowColumnParent)
							}
							f.SetCellFormula(sheet, Revaluation, fmt.Sprintf("=SUM(%s:%s)", RevaluationStartRow, RevaluationEndRow))
						
							if err := f.SetCellValue(sheet, EndingBalance, ""); err != nil {
								log.Fatal(err)
							}
						
							//
							if len(EndingBalance) > 3 {
								EndingBalance = EndingBalance[:2] + strconv.Itoa(rowColumnParent)
							} else {
								EndingBalance = EndingBalance[:1] + strconv.Itoa(rowColumnParent)
							}
							EndTotalCostParent = EndingBalance
							f.SetCellFormula(sheet, EndingBalance, fmt.Sprintf("=%s+%s+%s+%s+%s+%s", BeginningBalance, AcquisitionOfSubsidiary, Additions, Deductions, Reclassification, Revaluation))
						}
						if vv.Code == "TOTAL_ACCUMULATED_DEPRECIATION" {
						
							if err := f.SetCellValue(sheet, BeginningBalance, ""); err != nil {
								log.Fatal(err)
							}
							var BeginningBalanceStartRow string
							var BeginningBalanceEndRow string
							if len(BeginningBalance) > 3 {
								BeginningBalanceStartRow = BeginningBalance[:2] + "17"
								BeginningBalanceEndRow = BeginningBalance[:2] + "23"
							} else {
								BeginningBalanceStartRow = BeginningBalance[:1] + "17"
								BeginningBalanceEndRow = BeginningBalance[:1] + "23"
							}
							rowColumnParent = rowColumnParent+1
							if len(BeginningBalance) > 3 {
								BeginningBalance = BeginningBalance[:2] + strconv.Itoa(rowColumnParent)
							} else {
								BeginningBalance = BeginningBalance[:1] + strconv.Itoa(rowColumnParent)
							}
						
							TotalAccumParent = BeginningBalance
							f.SetCellFormula(sheet, BeginningBalance, fmt.Sprintf("=SUM(%s:%s)", BeginningBalanceStartRow, BeginningBalanceEndRow))
							CharCombineBeginingBalance = append(CharCombineBeginingBalance,BeginningBalance)
							rowIndex = rowIndex + 1
							if err := f.SetCellValue(sheet, AcquisitionOfSubsidiary, ""); err != nil {
								log.Fatal(err)
							}
						
							var AcquisitionOfSubsidiaryStartRow string
							var AcquisitionOfSubsidiaryEndRow string
							if len(AcquisitionOfSubsidiary) > 3 {
								AcquisitionOfSubsidiaryStartRow = AcquisitionOfSubsidiary[:2] + "17"
								AcquisitionOfSubsidiaryEndRow = AcquisitionOfSubsidiary[:2] + "23"
							} else {
								AcquisitionOfSubsidiaryStartRow = AcquisitionOfSubsidiary[:1] + "17"
								AcquisitionOfSubsidiaryEndRow = AcquisitionOfSubsidiary[:1] + "23"
							}
						
							if len(AcquisitionOfSubsidiary) > 3 {
								AcquisitionOfSubsidiary = AcquisitionOfSubsidiary[:2] + strconv.Itoa(rowColumnParent)
							} else {
								AcquisitionOfSubsidiary = AcquisitionOfSubsidiary[:1] + strconv.Itoa(rowColumnParent)
							}
						
							f.SetCellFormula(sheet, AcquisitionOfSubsidiary, fmt.Sprintf("=SUM(%s:%s)", AcquisitionOfSubsidiaryStartRow, AcquisitionOfSubsidiaryEndRow))
						
							rowIndex = rowIndex + 1
							if err := f.SetCellValue(sheet, Additions, ""); err != nil {
								log.Fatal(err)
							}
							var AdditionsStartRow string
							var AdditionsEndRow string
							if len(Additions) > 3 {
								AdditionsStartRow = Additions[:2] + "17"
								AdditionsEndRow = Additions[:2] + "23"
							} else {
								AdditionsStartRow = Additions[:1] + "17"
								AdditionsEndRow = Additions[:1] + "23"
							}
						
							if len(Additions) > 3 {
								Additions = Additions[:2] + strconv.Itoa(rowColumnParent)
							} else {
								Additions = Additions[:1] + strconv.Itoa(rowColumnParent)
							}
							f.SetCellFormula(sheet, Additions, fmt.Sprintf("=SUM(%s:%s)", AdditionsStartRow, AdditionsEndRow))
						
							rowIndex = rowIndex + 1
							if err := f.SetCellValue(sheet, Deductions, ""); err != nil {
								log.Fatal(err)
							}
							var DeductionsStartRow string
							var DeductionsEndRow string
							if len(Deductions) > 3 {
								DeductionsStartRow = Deductions[:2] + "17"
								DeductionsEndRow = Deductions[:2] + "23"
							} else {
								DeductionsStartRow = Deductions[:1] + "17"
								DeductionsEndRow = Deductions[:1] + "23"
							}
						
							if len(Deductions) > 3 {
								Deductions = Deductions[:2] + strconv.Itoa(rowColumnParent)
							} else {
								Deductions = Deductions[:1] + strconv.Itoa(rowColumnParent)
							}
							f.SetCellFormula(sheet, Deductions, fmt.Sprintf("=SUM(%s:%s)", DeductionsStartRow, DeductionsEndRow))
						
							rowIndex = rowIndex + 1
							if err := f.SetCellValue(sheet, Reclassification, ""); err != nil {
								log.Fatal(err)
							}
							var ReclassificationStartRow string
							var ReclassificationEndRow string
							if len(Reclassification) > 3 {
								ReclassificationStartRow = Reclassification[:2] + "17"
								ReclassificationEndRow = Reclassification[:2] + "23"
							} else {
								ReclassificationStartRow = Reclassification[:1] + "17"
								ReclassificationEndRow = Reclassification[:1] + "23"
							}
						
							if len(Reclassification) > 3 {
								Reclassification = Reclassification[:2] + strconv.Itoa(rowColumnParent)
							} else {
								Reclassification = Reclassification[:1] + strconv.Itoa(rowColumnParent)
							}
							f.SetCellFormula(sheet, Reclassification, fmt.Sprintf("=SUM(%s:%s)", ReclassificationStartRow, ReclassificationEndRow))
						
							rowIndex = rowIndex + 1
							if err := f.SetCellValue(sheet, Revaluation, ""); err != nil {
								log.Fatal(err)
							}
							var RevaluationStartRow string
							var RevaluationEndRow string
							if len(Revaluation) > 3 {
								RevaluationStartRow = Revaluation[:2] + "17"
								RevaluationEndRow = Revaluation[:2] + "23"
							} else {
								RevaluationStartRow = Revaluation[:1] + "17"
								RevaluationEndRow = Revaluation[:1] + "23"
							}
						
							if len(Revaluation) > 3 {
								Revaluation = Revaluation[:2] + strconv.Itoa(rowColumnParent)
							} else {
								Revaluation = Revaluation[:1] + strconv.Itoa(rowColumnParent)
							}
							f.SetCellFormula(sheet, Revaluation, fmt.Sprintf("=SUM(%s:%s)", RevaluationStartRow, RevaluationEndRow))
						
							if err := f.SetCellValue(sheet, EndingBalance, ""); err != nil {
								log.Fatal(err)
							}
						
							//
							if len(EndingBalance) > 3 {
								EndingBalance = EndingBalance[:2] + strconv.Itoa(rowColumnParent)
							} else {
								EndingBalance = EndingBalance[:1] + strconv.Itoa(rowColumnParent)
							}
							EndTotalAccumParent = EndingBalance
							f.SetCellFormula(sheet, EndingBalance, fmt.Sprintf("=%s+%s+%s+%s+%s+%s", BeginningBalance, AcquisitionOfSubsidiary, Additions, Deductions, Reclassification, Revaluation))
						}
						if vv.Code == "WORK_IN_PROCESS"  {
							// rowColumnParent = rowColumnParent+1
							// BeginningBalance = BeginningBalance[:1]+strconv.Itoa(rowColumnParent)
							WorkInProsesParent = BeginningBalance
							CharCombineBeginingBalance = append(CharCombineBeginingBalance,BeginningBalance)
							// EndingBalance = EndingBalance[:1]+strconv.Itoa(rowColumnParent)
							EndWorkInProsesParent = EndingBalance
						}
						if vv.Code == "NET_BOOK_VALUE"  {
							
							if err := f.SetCellValue(sheet, BeginningBalance, ""); err != nil {
								log.Fatal(err)
							}
							if err := f.SetCellValue(sheet, AcquisitionOfSubsidiary, ""); err != nil {
								log.Fatal(err)
							}
							if err := f.SetCellValue(sheet, Additions, ""); err != nil {
								log.Fatal(err)
							}
							if err := f.SetCellValue(sheet, Deductions, ""); err != nil {
								log.Fatal(err)
							}
							if err := f.SetCellValue(sheet, Reclassification, ""); err != nil {
								log.Fatal(err)
							}
							if err := f.SetCellValue(sheet, Revaluation, ""); err != nil {
								log.Fatal(err)
							}
							if err := f.SetCellValue(sheet, EndingBalance, ""); err != nil {
								log.Fatal(err)
							}
							rowColumnParent = rowColumnParent+1
							if len(BeginningBalance) > 3 {
								BeginningBalance = BeginningBalance[:2]+strconv.Itoa(rowColumnParent)
								
							} else {
								BeginningBalance = BeginningBalance[:1]+strconv.Itoa(rowColumnParent)
								
							}
							CharCombineBeginingBalance = append(CharCombineBeginingBalance,BeginningBalance)
							if len(EndingBalance) > 3 {
								EndingBalance = EndingBalance[:2]+strconv.Itoa(rowColumnParent)
								
							} else {
								EndingBalance = EndingBalance[:1]+strconv.Itoa(rowColumnParent)
								
							}
							f.SetCellFormula(sheet, BeginningBalance, fmt.Sprintf("=%s-%s+%s", TotalCostParent, TotalAccumParent, WorkInProsesParent))
							f.SetCellFormula(sheet, EndingBalance, fmt.Sprintf("=%s-%s+%s", EndTotalCostParent, EndTotalAccumParent, EndWorkInProsesParent))	
						}
					rowColumnParent++
				}
			}
			

		rowIndex2 = rowIndex2 + 1
		for _, cd := range consolidationData.ConsolidationBridge {
			
			C4, _ := excelize.CoordinatesToCellName(rowIndex2, 2)
			D4, _ := excelize.CoordinatesToCellName(rowIndex2, rowColumn2)
			rowIndex2 = rowIndex2 + 1
			E4, _ := excelize.CoordinatesToCellName(rowIndex2, rowColumn2)
			rowIndex2 = rowIndex2 + 1
			F4, _ := excelize.CoordinatesToCellName(rowIndex2, rowColumn2)
			rowIndex2 = rowIndex2 + 1
			G4, _ := excelize.CoordinatesToCellName(rowIndex2, rowColumn2)
			rowIndex2 = rowIndex2 + 1
			H4, _ := excelize.CoordinatesToCellName(rowIndex2, rowColumn2)
			rowIndex2 = rowIndex2 + 1
			I4, _ := excelize.CoordinatesToCellName(rowIndex2, rowColumn2)
			rowIndex2 = rowIndex2 + 1
			J4, _ := excelize.CoordinatesToCellName(rowIndex2, rowColumn2)
			rowIndex2 = rowIndex2 + 1
			
			f.SetCellValue(sheet, C4, cd.Company.Name)
			f.SetCellValue(sheet, D4, "Beginning Balance")
			f.SetCellValue(sheet, E4, "Acquisition of Subsidiary")
			f.SetCellValue(sheet, F4, "Additions (+)")
			f.SetCellValue(sheet, G4, "Deductions (-)")
			f.SetCellValue(sheet, H4, "Reclassification")
			f.SetCellValue(sheet, I4, "Revaluation")
			f.SetCellValue(sheet, J4, "Ending balance")
			
			rowIndex2++

			filterMutasiFa := model.MutasiFaFilterModel{}
			filterMutasiFa.CompanyID = &cd.CompanyID
			filterMutasiFa.Period = &cd.Period
			filterMutasiFa.Versions = &cd.Versions
			if cd.Versions == 0 {
				filterMutasiFa.Versions = &cd.ConsolidationVersions
			}
			mutasiFa, err := s.MutasiFaRepository.FindByCriteria(ctx, &filterMutasiFa)
			if err != nil {
				return nil, err
			}
			tmpStr := "MUTASI-FA"
			criteriaBridge := model.FormatterBridgesFilterModel{}
			criteriaBridge.FormatterBridgesFilter.Source = &tmpStr
			// criteriaBridge.FormatterBridgesFilter.FormatterID = &data.ID
			criteriaBridge.FormatterBridgesFilter.TrxRefID = &mutasiFa.ID

			bridges, err := s.FormatterBridgesRepository.FindWithCriteriaNew(ctx, &criteriaBridge)
			if err != nil {
				return nil, err
			}

			// var list []int
			rowColumn := 5
			var TotalCost string
			var TotalAccum string
			var WorkInProses string
			var EndTotalCost string
			var EndTotalAccum string
			var EndWorkInProses string
			for _, brid := range *bridges {

				if brid.FormatterID == 21 {
					continue
				}
				criteriaMF := model.MutasiFaDetailFilterModel{}
				// criteriaMF.Code = &v.Code
				criteriaMF.FormatterBridgesID = &brid.ID
				criteriaMF.MutasiFaID = &mutasiFa.ID

				MutasiFaDetail, err := s.MutasiFaDetailRepository.Find(ctx, &criteriaMF)
				if err != nil {
					continue
				}
				
				for _, vv := range *MutasiFaDetail {
					rowIndex := 10 + rowPlus
					
					if vv.Code == "ACCUMULATED_DEPRECIATION"  {
						rowColumn = rowColumn+1
					}
					
						BeginningBalance, _ := excelize.CoordinatesToCellName(rowIndex, rowColumn)
						f.SetCellValue(sheet, BeginningBalance, 0)

						rowIndex = rowIndex + 1
						AcquisitionOfSubsidiary, _ := excelize.CoordinatesToCellName(rowIndex, rowColumn)
						f.SetCellValue(sheet, AcquisitionOfSubsidiary, 0)

						rowIndex = rowIndex + 1
						Additions, _ := excelize.CoordinatesToCellName(rowIndex, rowColumn)
						f.SetCellValue(sheet, Additions, 0)

						rowIndex = rowIndex + 1
						Deductions, _ := excelize.CoordinatesToCellName(rowIndex, rowColumn)
						f.SetCellValue(sheet, Deductions, 0)

						rowIndex = rowIndex + 1
						Reclassification, _ := excelize.CoordinatesToCellName(rowIndex, rowColumn)
						f.SetCellValue(sheet, Reclassification, 0)

						rowIndex = rowIndex + 1
						Revaluation, _ := excelize.CoordinatesToCellName(rowIndex, rowColumn)
						f.SetCellValue(sheet, Revaluation, 0)
						rowIndex = rowIndex + 1
						EndingBalance, _ := excelize.CoordinatesToCellName(rowIndex, rowColumn)
						f.SetCellFormula(sheet, EndingBalance, fmt.Sprintf("=%s+%s+%s+%s+%s+%s", BeginningBalance, AcquisitionOfSubsidiary, Additions, Deductions, Reclassification, Revaluation))

						

						if vv.BeginningBalance != nil && *vv.BeginningBalance != 0 {
							f.SetCellValue(sheet, BeginningBalance, *vv.BeginningBalance)
						}
						if vv.AcquisitionOfSubsidiary != nil && *vv.AcquisitionOfSubsidiary != 0 {
							f.SetCellValue(sheet, AcquisitionOfSubsidiary, *vv.AcquisitionOfSubsidiary)
						}
						if vv.Additions != nil && *vv.Additions != 0 {
							f.SetCellValue(sheet, Additions, *vv.Additions)
						}
						if vv.Deductions != nil && *vv.Deductions != 0 {
							f.SetCellValue(sheet, Deductions, *vv.Deductions)
						}
						if vv.Reclassification != nil && *vv.Reclassification != 0 {
							f.SetCellValue(sheet, Reclassification, *vv.Reclassification)
						}
						if vv.Revaluation != nil && *vv.Revaluation != 0 {
							f.SetCellValue(sheet, Revaluation, *vv.Revaluation)
						}
				
						if vv.Code == ""  {
							if err := f.SetCellValue(sheet, BeginningBalance, ""); err != nil {
								log.Fatal(err)
							}
							if err := f.SetCellValue(sheet, AcquisitionOfSubsidiary, ""); err != nil {
								log.Fatal(err)
							}
							if err := f.SetCellValue(sheet, Additions, ""); err != nil {
								log.Fatal(err)
							}
							if err := f.SetCellValue(sheet, Deductions, ""); err != nil {
								log.Fatal(err)
							}
							if err := f.SetCellValue(sheet, Reclassification, ""); err != nil {
								log.Fatal(err)
							}
							if err := f.SetCellValue(sheet, Revaluation, ""); err != nil {
								log.Fatal(err)
							}
							if err := f.SetCellValue(sheet, EndingBalance, ""); err != nil {
								log.Fatal(err)
							}
						}
						if vv.Code == "TOTAL_COST" {

							if err := f.SetCellValue(sheet, BeginningBalance, ""); err != nil {
								log.Fatal(err)
							}
							var BeginningBalanceStartRow string
							var BeginningBalanceEndRow string
							if len(BeginningBalance) > 3 {
								BeginningBalanceStartRow = BeginningBalance[:2] + "6"
								BeginningBalanceEndRow = BeginningBalance[:2] + "12"
							} else {
								BeginningBalanceStartRow = BeginningBalance[:1] + "6"
								BeginningBalanceEndRow = BeginningBalance[:1] + "12"
							}
							rowColumn = rowColumn+1
							if len(BeginningBalance) > 3 {
								BeginningBalance = BeginningBalance[:2] + strconv.Itoa(rowColumn)
							} else {
								BeginningBalance = BeginningBalance[:1] + strconv.Itoa(rowColumn)
							}
						
							TotalCost = BeginningBalance
							f.SetCellFormula(sheet, BeginningBalance, fmt.Sprintf("=SUM(%s:%s)", BeginningBalanceStartRow, BeginningBalanceEndRow))
						
							rowIndex = rowIndex + 1
							if err := f.SetCellValue(sheet, AcquisitionOfSubsidiary, ""); err != nil {
								log.Fatal(err)
							}
						
							var AcquisitionOfSubsidiaryStartRow string
							var AcquisitionOfSubsidiaryEndRow string
							if len(AcquisitionOfSubsidiary) > 3 {
								AcquisitionOfSubsidiaryStartRow = AcquisitionOfSubsidiary[:2] + "6"
								AcquisitionOfSubsidiaryEndRow = AcquisitionOfSubsidiary[:2] + "12"
							} else {
								AcquisitionOfSubsidiaryStartRow = AcquisitionOfSubsidiary[:1] + "6"
								AcquisitionOfSubsidiaryEndRow = AcquisitionOfSubsidiary[:1] + "12"
							}
						
							if len(AcquisitionOfSubsidiary) > 3 {
								AcquisitionOfSubsidiary = AcquisitionOfSubsidiary[:2] + strconv.Itoa(rowColumn)
							} else {
								AcquisitionOfSubsidiary = AcquisitionOfSubsidiary[:1] + strconv.Itoa(rowColumn)
							}
						
							f.SetCellFormula(sheet, AcquisitionOfSubsidiary, fmt.Sprintf("=SUM(%s:%s)", AcquisitionOfSubsidiaryStartRow, AcquisitionOfSubsidiaryEndRow))
						
							rowIndex = rowIndex + 1
							if err := f.SetCellValue(sheet, Additions, ""); err != nil {
								log.Fatal(err)
							}
							var AdditionsStartRow string
							var AdditionsEndRow string
							if len(Additions) > 3 {
								AdditionsStartRow = Additions[:2] + "6"
								AdditionsEndRow = Additions[:2] + "12"
							} else {
								AdditionsStartRow = Additions[:1] + "6"
								AdditionsEndRow = Additions[:1] + "12"
							}
						
							if len(Additions) > 3 {
								Additions = Additions[:2] + strconv.Itoa(rowColumn)
							} else {
								Additions = Additions[:1] + strconv.Itoa(rowColumn)
							}
							f.SetCellFormula(sheet, Additions, fmt.Sprintf("=SUM(%s:%s)", AdditionsStartRow, AdditionsEndRow))
						
							rowIndex = rowIndex + 1
							if err := f.SetCellValue(sheet, Deductions, ""); err != nil {
								log.Fatal(err)
							}
							var DeductionsStartRow string
							var DeductionsEndRow string
							if len(Deductions) > 3 {
								DeductionsStartRow = Deductions[:2] + "6"
								DeductionsEndRow = Deductions[:2] + "12"
							} else {
								DeductionsStartRow = Deductions[:1] + "6"
								DeductionsEndRow = Deductions[:1] + "12"
							}
						
							if len(Deductions) > 3 {
								Deductions = Deductions[:2] + strconv.Itoa(rowColumn)
							} else {
								Deductions = Deductions[:1] + strconv.Itoa(rowColumn)
							}
							f.SetCellFormula(sheet, Deductions, fmt.Sprintf("=SUM(%s:%s)", DeductionsStartRow, DeductionsEndRow))
						
							rowIndex = rowIndex + 1
							if err := f.SetCellValue(sheet, Reclassification, ""); err != nil {
								log.Fatal(err)
							}
							var ReclassificationStartRow string
							var ReclassificationEndRow string
							if len(Reclassification) > 3 {
								ReclassificationStartRow = Reclassification[:2] + "6"
								ReclassificationEndRow = Reclassification[:2] + "12"
							} else {
								ReclassificationStartRow = Reclassification[:1] + "6"
								ReclassificationEndRow = Reclassification[:1] + "12"
							}
						
							if len(Reclassification) > 3 {
								Reclassification = Reclassification[:2] + strconv.Itoa(rowColumn)
							} else {
								Reclassification = Reclassification[:1] + strconv.Itoa(rowColumn)
							}
							f.SetCellFormula(sheet, Reclassification, fmt.Sprintf("=SUM(%s:%s)", ReclassificationStartRow, ReclassificationEndRow))
						
							rowIndex = rowIndex + 1
							if err := f.SetCellValue(sheet, Revaluation, ""); err != nil {
								log.Fatal(err)
							}
							var RevaluationStartRow string
							var RevaluationEndRow string
							if len(Revaluation) > 3 {
								RevaluationStartRow = Revaluation[:2] + "6"
								RevaluationEndRow = Revaluation[:2] + "12"
							} else {
								RevaluationStartRow = Revaluation[:1] + "6"
								RevaluationEndRow = Revaluation[:1] + "12"
							}
						
							if len(Revaluation) > 3 {
								Revaluation = Revaluation[:2] + strconv.Itoa(rowColumn)
							} else {
								Revaluation = Revaluation[:1] + strconv.Itoa(rowColumn)
							}
							f.SetCellFormula(sheet, Revaluation, fmt.Sprintf("=SUM(%s:%s)", RevaluationStartRow, RevaluationEndRow))
						
							if err := f.SetCellValue(sheet, EndingBalance, ""); err != nil {
								log.Fatal(err)
							}
						
							//
							if len(EndingBalance) > 3 {
								EndingBalance = EndingBalance[:2] + strconv.Itoa(rowColumn)
							} else {
								EndingBalance = EndingBalance[:1] + strconv.Itoa(rowColumn)
							}
							EndTotalCost = EndingBalance
							f.SetCellFormula(sheet, EndingBalance, fmt.Sprintf("=%s+%s+%s+%s+%s+%s", BeginningBalance, AcquisitionOfSubsidiary, Additions, Deductions, Reclassification, Revaluation))
						}
						if vv.Code == "TOTAL_ACCUMULATED_DEPRECIATION" {
						
							if err := f.SetCellValue(sheet, BeginningBalance, ""); err != nil {
								log.Fatal(err)
							}
							var BeginningBalanceStartRow string
							var BeginningBalanceEndRow string
							if len(BeginningBalance) > 3 {
								BeginningBalanceStartRow = BeginningBalance[:2] + "17"
								BeginningBalanceEndRow = BeginningBalance[:2] + "23"
							} else {
								BeginningBalanceStartRow = BeginningBalance[:1] + "17"
								BeginningBalanceEndRow = BeginningBalance[:1] + "23"
							}
							rowColumn = rowColumn+1
							if len(BeginningBalance) > 3 {
								BeginningBalance = BeginningBalance[:2] + strconv.Itoa(rowColumn)
							} else {
								BeginningBalance = BeginningBalance[:1] + strconv.Itoa(rowColumn)
							}
						
							TotalAccum = BeginningBalance
							f.SetCellFormula(sheet, BeginningBalance, fmt.Sprintf("=SUM(%s:%s)", BeginningBalanceStartRow, BeginningBalanceEndRow))
						
							rowIndex = rowIndex + 1
							if err := f.SetCellValue(sheet, AcquisitionOfSubsidiary, ""); err != nil {
								log.Fatal(err)
							}
						
							var AcquisitionOfSubsidiaryStartRow string
							var AcquisitionOfSubsidiaryEndRow string
							if len(AcquisitionOfSubsidiary) > 3 {
								AcquisitionOfSubsidiaryStartRow = AcquisitionOfSubsidiary[:2] + "17"
								AcquisitionOfSubsidiaryEndRow = AcquisitionOfSubsidiary[:2] + "23"
							} else {
								AcquisitionOfSubsidiaryStartRow = AcquisitionOfSubsidiary[:1] + "17"
								AcquisitionOfSubsidiaryEndRow = AcquisitionOfSubsidiary[:1] + "23"
							}
						
							if len(AcquisitionOfSubsidiary) > 3 {
								AcquisitionOfSubsidiary = AcquisitionOfSubsidiary[:2] + strconv.Itoa(rowColumn)
							} else {
								AcquisitionOfSubsidiary = AcquisitionOfSubsidiary[:1] + strconv.Itoa(rowColumn)
							}
						
							f.SetCellFormula(sheet, AcquisitionOfSubsidiary, fmt.Sprintf("=SUM(%s:%s)", AcquisitionOfSubsidiaryStartRow, AcquisitionOfSubsidiaryEndRow))
						
							rowIndex = rowIndex + 1
							if err := f.SetCellValue(sheet, Additions, ""); err != nil {
								log.Fatal(err)
							}
							var AdditionsStartRow string
							var AdditionsEndRow string
							if len(Additions) > 3 {
								AdditionsStartRow = Additions[:2] + "17"
								AdditionsEndRow = Additions[:2] + "23"
							} else {
								AdditionsStartRow = Additions[:1] + "17"
								AdditionsEndRow = Additions[:1] + "23"
							}
						
							if len(Additions) > 3 {
								Additions = Additions[:2] + strconv.Itoa(rowColumn)
							} else {
								Additions = Additions[:1] + strconv.Itoa(rowColumn)
							}
							f.SetCellFormula(sheet, Additions, fmt.Sprintf("=SUM(%s:%s)", AdditionsStartRow, AdditionsEndRow))
						
							rowIndex = rowIndex + 1
							if err := f.SetCellValue(sheet, Deductions, ""); err != nil {
								log.Fatal(err)
							}
							var DeductionsStartRow string
							var DeductionsEndRow string
							if len(Deductions) > 3 {
								DeductionsStartRow = Deductions[:2] + "17"
								DeductionsEndRow = Deductions[:2] + "23"
							} else {
								DeductionsStartRow = Deductions[:1] + "17"
								DeductionsEndRow = Deductions[:1] + "23"
							}
						
							if len(Deductions) > 3 {
								Deductions = Deductions[:2] + strconv.Itoa(rowColumn)
							} else {
								Deductions = Deductions[:1] + strconv.Itoa(rowColumn)
							}
							f.SetCellFormula(sheet, Deductions, fmt.Sprintf("=SUM(%s:%s)", DeductionsStartRow, DeductionsEndRow))
						
							rowIndex = rowIndex + 1
							if err := f.SetCellValue(sheet, Reclassification, ""); err != nil {
								log.Fatal(err)
							}
							var ReclassificationStartRow string
							var ReclassificationEndRow string
							if len(Reclassification) > 3 {
								ReclassificationStartRow = Reclassification[:2] + "17"
								ReclassificationEndRow = Reclassification[:2] + "23"
							} else {
								ReclassificationStartRow = Reclassification[:1] + "17"
								ReclassificationEndRow = Reclassification[:1] + "23"
							}
						
							if len(Reclassification) > 3 {
								Reclassification = Reclassification[:2] + strconv.Itoa(rowColumn)
							} else {
								Reclassification = Reclassification[:1] + strconv.Itoa(rowColumn)
							}
							f.SetCellFormula(sheet, Reclassification, fmt.Sprintf("=SUM(%s:%s)", ReclassificationStartRow, ReclassificationEndRow))
						
							rowIndex = rowIndex + 1
							if err := f.SetCellValue(sheet, Revaluation, ""); err != nil {
								log.Fatal(err)
							}
							var RevaluationStartRow string
							var RevaluationEndRow string
							if len(Revaluation) > 3 {
								RevaluationStartRow = Revaluation[:2] + "17"
								RevaluationEndRow = Revaluation[:2] + "23"
							} else {
								RevaluationStartRow = Revaluation[:1] + "17"
								RevaluationEndRow = Revaluation[:1] + "23"
							}
						
							if len(Revaluation) > 3 {
								Revaluation = Revaluation[:2] + strconv.Itoa(rowColumn)
							} else {
								Revaluation = Revaluation[:1] + strconv.Itoa(rowColumn)
							}
							f.SetCellFormula(sheet, Revaluation, fmt.Sprintf("=SUM(%s:%s)", RevaluationStartRow, RevaluationEndRow))
						
							if err := f.SetCellValue(sheet, EndingBalance, ""); err != nil {
								log.Fatal(err)
							}
						
							//
							if len(EndingBalance) > 3 {
								EndingBalance = EndingBalance[:2] + strconv.Itoa(rowColumn)
							} else {
								EndingBalance = EndingBalance[:1] + strconv.Itoa(rowColumn)
							}
							EndTotalAccum = EndingBalance
							f.SetCellFormula(sheet, EndingBalance, fmt.Sprintf("=%s+%s+%s+%s+%s+%s", BeginningBalance, AcquisitionOfSubsidiary, Additions, Deductions, Reclassification, Revaluation))
						}
						if vv.Code == "WORK_IN_PROCESS"  {
							// rowColumn = rowColumn+1
							// BeginningBalance = BeginningBalance[:1]+strconv.Itoa(rowColumn)
							WorkInProses = BeginningBalance
							// EndingBalance = EndingBalance[:1]+strconv.Itoa(rowColumn)
							EndWorkInProses = EndingBalance
						}
						if vv.Code == "NET_BOOK_VALUE"  {
							
							if err := f.SetCellValue(sheet, BeginningBalance, ""); err != nil {
								log.Fatal(err)
							}
							if err := f.SetCellValue(sheet, AcquisitionOfSubsidiary, ""); err != nil {
								log.Fatal(err)
							}
							if err := f.SetCellValue(sheet, Additions, ""); err != nil {
								log.Fatal(err)
							}
							if err := f.SetCellValue(sheet, Deductions, ""); err != nil {
								log.Fatal(err)
							}
							if err := f.SetCellValue(sheet, Reclassification, ""); err != nil {
								log.Fatal(err)
							}
							if err := f.SetCellValue(sheet, Revaluation, ""); err != nil {
								log.Fatal(err)
							}
							if err := f.SetCellValue(sheet, EndingBalance, ""); err != nil {
								log.Fatal(err)
							}
							rowColumn = rowColumn+1
							if len(BeginningBalance) > 3 {
								BeginningBalance = BeginningBalance[:2]+strconv.Itoa(rowColumn)
								
							} else {
								BeginningBalance = BeginningBalance[:1]+strconv.Itoa(rowColumn)
								
							}
							if len(EndingBalance) > 3 {
								EndingBalance = EndingBalance[:2]+strconv.Itoa(rowColumn)
								
							} else {
								EndingBalance = EndingBalance[:1]+strconv.Itoa(rowColumn)
								
							}
							f.SetCellFormula(sheet, BeginningBalance, fmt.Sprintf("=%s-%s+%s", TotalCost, TotalAccum, WorkInProses))
							f.SetCellFormula(sheet, EndingBalance, fmt.Sprintf("=%s-%s+%s", EndTotalCost, EndTotalAccum, EndWorkInProses))	
						}
					rowColumn++
				}
			}
			rowPlus = rowPlus + 8
		}
		rowJelim := rowIndex2
		CJelim, _ := excelize.CoordinatesToCellName(rowIndex2, 2)
		DJelim, _ := excelize.CoordinatesToCellName(rowIndex2, rowColumn2)
		rowIndex2 = rowIndex2 + 1
		EJelim, _ := excelize.CoordinatesToCellName(rowIndex2, rowColumn2)
		rowIndex2 = rowIndex2 + 1
		FJelim, _ := excelize.CoordinatesToCellName(rowIndex2, rowColumn2)
		rowIndex2 = rowIndex2 + 1
		GJelim, _ := excelize.CoordinatesToCellName(rowIndex2, rowColumn2)
		rowIndex2 = rowIndex2 + 1
		HJelim, _ := excelize.CoordinatesToCellName(rowIndex2, rowColumn2)
		rowIndex2 = rowIndex2 + 1
		IJelim, _ := excelize.CoordinatesToCellName(rowIndex2, rowColumn2)
		rowIndex2 = rowIndex2 + 1
		JJelim, _ := excelize.CoordinatesToCellName(rowIndex2, rowColumn2)
		
		f.SetCellValue(sheet, CJelim, "JELIM")
		f.SetCellValue(sheet, DJelim, "Beginning Balance")
		f.SetCellValue(sheet, EJelim, "Acquisition of Subsidiary")
		f.SetCellValue(sheet, FJelim, "Additions (+)")
		f.SetCellValue(sheet, GJelim, "Deductions (-)")
		f.SetCellValue(sheet, HJelim, "Reclassification")
		f.SetCellValue(sheet, IJelim, "Revaluation")
		f.SetCellValue(sheet, JJelim, "Ending balance")

		JelimCoba, _ := excelize.CoordinatesToCellName(rowJelim, rowColumn2+1)
		f.SetCellValue(sheet, JelimCoba, 99)

		rowIndex2 = rowIndex2+2
		rowCombine := rowJelim+8
		CCombaine, _ := excelize.CoordinatesToCellName(rowIndex2, 2)
		DCombaine, _ := excelize.CoordinatesToCellName(rowIndex2, rowColumn2)
		rowIndex2 = rowIndex2 + 1
		ECombaine, _ := excelize.CoordinatesToCellName(rowIndex2, rowColumn2)
		rowIndex2 = rowIndex2 + 1
		FCombaine, _ := excelize.CoordinatesToCellName(rowIndex2, rowColumn2)
		rowIndex2 = rowIndex2 + 1
		GCombaine, _ := excelize.CoordinatesToCellName(rowIndex2, rowColumn2)
		rowIndex2 = rowIndex2 + 1
		HCombaine, _ := excelize.CoordinatesToCellName(rowIndex2, rowColumn2)
		rowIndex2 = rowIndex2 + 1
		ICombaine, _ := excelize.CoordinatesToCellName(rowIndex2, rowColumn2)
		rowIndex2 = rowIndex2 + 1
		JCombaine, _ := excelize.CoordinatesToCellName(rowIndex2, rowColumn2)

		f.SetCellValue(sheet, CCombaine, "COMBAINE")
		f.SetCellValue(sheet, DCombaine, "Beginning Balance")
		f.SetCellValue(sheet, ECombaine, "Acquisition of Subsidiary")
		f.SetCellValue(sheet, FCombaine, "Additions (+)")
		f.SetCellValue(sheet, GCombaine, "Deductions (-)")
		f.SetCellValue(sheet, HCombaine, "Reclassification")
		f.SetCellValue(sheet, ICombaine, "Revaluation")
		f.SetCellValue(sheet, JCombaine, "Ending balance")

		rowColumn2 = rowColumn2+1
		var cBgining []string
		cordinatParent := 2
		CmbnBeginingBalance, _ := excelize.CoordinatesToCellName(rowCombine, rowColumn2)
		for i := cordinatParent; i < rowIndex2; i++ {
			if i%8 == 2 && i != rowIndex2-cordinatParent  {
				searchCharBegining, _ := excelize.CoordinatesToCellName(i, rowColumn2)
				cBgining = append(cBgining, searchCharBegining)
			}
		}
		fmlBeginingBalance := strings.Join(cBgining, "+")
		f.SetCellFormula(sheet, CmbnBeginingBalance, fmt.Sprintf("=%s", fmlBeginingBalance))


		rowIndex2 = rowIndex2 + 2
		Control, _ := excelize.CoordinatesToCellName(rowIndex2, rowColumn2)
		f.SetCellValue(sheet, Control, "Control")
	}

	return f, nil
}