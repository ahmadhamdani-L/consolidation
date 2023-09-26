package validation

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"math"
	"strconv"
	"strings"
	"sync"
	"time"
	"worker-validation/internal/abstraction"
	"worker-validation/internal/factory"
	kafkaproducer "worker-validation/internal/kafka/producer"
	"worker-validation/internal/model"
	"worker-validation/internal/repository"
	"worker-validation/pkg/constant"
	utilDate "worker-validation/pkg/util/date"
	"worker-validation/pkg/util/trxmanager"

	"golang.org/x/text/language"
	"golang.org/x/text/message"
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
	ValidationRepository                 repository.Validation
	ControllerRepository                 repository.Controller
	EmployeeBenefitRepository            repository.EmployeeBenefit
	EmployeeBenefitDetailRepository      repository.EmployeeBenefitDetail
}

type Service interface {
	Validation(ctx *abstraction.Context, payload *abstraction.JsonData)
	InitiateData(ctx *abstraction.Context, payload *abstraction.JsonData) error
	Validate(ctx *abstraction.Context, jsonData *abstraction.JsonData, validationData *model.ValidationDetailEntityModel, wg *sync.WaitGroup, validasiModul *string)
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
	validationRepo := f.ValidationRepository
	controllerRepo := f.ControllerRepository
	employeeBenefitRepo := f.EmployeeBenefitRepository
	employeeBenefitDetailRepo := f.EmployeeBenefitDetailRepository
	db := f.Db

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
		ValidationRepository:                 validationRepo,
		ControllerRepository:                 controllerRepo,
		EmployeeBenefitRepository:            employeeBenefitRepo,
		EmployeeBenefitDetailRepository:      employeeBenefitDetailRepo,
	}
}

type DataToValidate struct {
	MasterID   int
	ListDataID []int
	ModulID    map[int]string
}

type JsonData struct {
	Company  string `json:"company"`
	Name     string `json:"name"`
	Period   string `json:"period"`
	Versions int    `json:"versions"`
	DataID   int    `json:"data_id"`
	Errors   string `json:"errors"`
}

var hasErr bool
var msg string

func (s *service) SendNotification(ctx *abstraction.Context, trialBalanceID *int) error {
	var jsonData JsonData
	trialBalance, err := s.TrialBalanceRepository.FindByID(ctx, trialBalanceID)
	if err != nil {
		log.Println(err)
		return err
	}

	if trialBalance != nil {
		jsonData.Versions = trialBalance.Versions
		jsonData.Period = trialBalance.Period
		jsonData.DataID = trialBalance.ID
	}

	jsonData.Name = "validation"
	jsonData.Company = trialBalance.Company.Name
	jsonStr, err := json.Marshal(jsonData)
	if err != nil {
		log.Println(err)
		return err
	}

	notifData := model.NotificationEntityModel{}
	notifData.Context = ctx
	notifData.Description = "Proses Validasi telah selesai."
	if hasErr {
		notifData.Description += " Silakan cek hasil validasi karena terdapat beberapa modul yang tidak valid."
	}
	tmpfalse := false
	notifData.IsOpened = &tmpfalse
	notifData.CreatedBy = ctx.Auth.ID
	notifData.CreatedAt = *utilDate.DateTodayLocal()
	notifData.Data = string(jsonStr)

	_, err = s.NotificationRepository.Create(ctx, &notifData)
	if err != nil {
		log.Println(err)
		return err
	}

	waktu := time.Now()
	map1 := kafkaproducer.JsonData{
		UserID:    ctx.Auth.ID,
		CompanyID: ctx.Auth.CompanyID,
		Name:      "validation",
		Timestamp: &waktu,
		Data:      notifData.Description,
		Filter: struct {
			Period   string
			Versions int
		}{trialBalance.Period, trialBalance.Versions},
	}

	jsonStr, err = json.Marshal(map1)
	if err != nil {
		log.Println(err)
		return err
	}

	go kafkaproducer.NewProducer("NOTIFICATION").SendMessage("NOTIFICATION", string(jsonStr))
	return nil
}

func (s *service) Validation(ctx *abstraction.Context, payload *abstraction.JsonData) {
	// runtime.GOMAXPROCS(runtime.NumCPU())
	wg := new(sync.WaitGroup)
	hasErr = false
	var tmpData DataToValidate
	if err := json.Unmarshal([]byte(payload.Data), &tmpData); err != nil {
		fmt.Printf("Error unmarshalling. Error: %s", err.Error())
		return
	}

	modulToValidate := make(map[int]string)
	for k, v := range tmpData.ModulID {
		if _, ok := modulToValidate[k]; !ok {
			modulToValidate[k] = v
		}
	}

	{
		for _, vID := range tmpData.ListDataID {
			if vID == 0 {
				continue
			}
			trialBalance, err := s.TrialBalanceRepository.FindByID(ctx, &vID)
			if err != nil {
				log.Println(err)
				return
			}
			if trialBalance.Status == constant.MODUL_STATUS_VALIDATED {
				continue
			}
			if len(tmpData.ModulID) == 0 {
				if err := trxmanager.New(s.Db).WithTrx(ctx, func(ctx *abstraction.Context) error {
					tmpDetailData, err := s.InitiateData(ctx, payload, &trialBalance.ID)
					if err != nil {
						log.Println(err)
						return err
					}
					for k, v := range tmpDetailData {
						if _, ok := modulToValidate[k]; !ok {
							modulToValidate[k] = v
						}
					}

					return nil
				}); err != nil {
					fmt.Printf("Error validationing data. Detail: %s", err.Error())
					return
				}
			}

			// wg.Add(len(modulToValidate))
			var tmpErr []error
			for id, modul := range modulToValidate {
				// go func(ctx *abstraction.Context, id int, modul string, tmpErr []error) {
				// defer wg.Done()
				newContext := new(abstraction.Context)
				newContext = ctx
				if err := trxmanager.New(s.Db).WithTrx(newContext, func(newContext *abstraction.Context) error {
					tmpdetailValidationData, err := s.ValidationRepository.FindByID(newContext, &id)
					if err != nil {
						log.Println(err)
						tmpErr = append(tmpErr, err)
						return err
					}

					payload.CompanyID = tmpdetailValidationData.CompanyID
					payload.Filter.Period = tmpdetailValidationData.Period
					payload.Filter.Versions = tmpdetailValidationData.Versions
					err = s.Validate(newContext, payload, tmpdetailValidationData, wg, &modul)
					if err != nil {
						log.Println(err)
						tmpErr = append(tmpErr, err)
						if err == gorm.ErrRecordNotFound || err.Error() == constant.VALIDATION_NOTE_NOT_BALANCE {
							return nil
						}
						return err
					}
					return nil
				}); err != nil {
					log.Println(err)
					log.Println("313Error validationing data.")
					return
				}
				// }(ctx, id, modul, tmpErr)
			}

			// wg.Wait()
			newContext := new(abstraction.Context)
			newContext = ctx
			if err := trxmanager.New(s.Db).WithTrx(newContext, func(ctx *abstraction.Context) error {

				dataTB := model.TrialBalanceEntityModel{}
				dataTB.Context = ctx
				dataTB.ValidationNote = constant.VALIDATION_NOTE_BALANCE
				// dataTB.Status= constant.MODUL_STATUS_VALIDATED
				if len(tmpErr) > 0 {
					hasErr = true
					dataTB.ValidationNote = constant.VALIDATION_NOTE_NOT_BALANCE
					dataTB.Status= constant.MODUL_STATUS_DRAFT
				}

				if _, err := s.TrialBalanceRepository.Update(ctx, &trialBalance.ID, &dataTB); err != nil {
					log.Println(err)
				}

				if err := s.SendNotification(ctx, &trialBalance.ID); err != nil {
					log.Println(err)
					return err
				}
				return nil
			}); err != nil {
				log.Println("Error validationing data.")
				return
			}
		}
	}
}

func (s *service) InitiateData(ctx *abstraction.Context, payload *abstraction.JsonData, trialBalanceID *int) (map[int]string, error) {
	var tmpData DataToValidate
	initData := make(map[string]int)
	result := make(map[int]string)
	if err := json.Unmarshal([]byte(payload.Data), &tmpData); err != nil {
		return nil, err
	}
	// ==== TrialBalance =====
	trialBalance, err := s.TrialBalanceRepository.FindByID(ctx, trialBalanceID)
	if err != nil {
		return nil, err
	}

	if trialBalance.Status != constant.MODUL_STATUS_DRAFT {
		return nil, errors.New("data tidak ditemukan")
	}
	initData["TRIAL-BALANCE"] = trialBalance.ID

	// Criteria Data
	criteria := model.FilterData{}
	criteria.Period = trialBalance.Period
	criteria.CompanyID = trialBalance.CompanyID
	criteria.Status = trialBalance.Status
	criteria.Versions = trialBalance.Versions

	// ==== Aging Utang Piutang =====
	agingUtangPiutang, err := s.AgingUtangPiutangRepository.FindByCriteria(ctx, &criteria)
	if err != nil {
		return nil, err
	}
	fmt.Println(agingUtangPiutang.ID)
	initData["AGING-UTANG-PIUTANG"] = agingUtangPiutang.ID

	// ==== Mutasi Persediaan =====
	mutasiPersediaan, err := s.MutasiPersediaanRepository.FindByCriteria(ctx, &criteria)
	if err != nil {
		return nil, err
	}
	fmt.Println(mutasiPersediaan.ID)
	initData["MUTASI-PERSEDIAAN"] = mutasiPersediaan.ID

	// ==== Mutasi Fa =====
	mutasiFa, err := s.MutasiFaRepository.FindByCriteria(ctx, &criteria)
	if err != nil {
		return nil, err
	}
	fmt.Println(mutasiFa.ID)
	initData["MUTASI-FA"] = mutasiFa.ID

	// ==== Mutasi Dta =====
	mutasiDta, err := s.MutasiFaRepository.FindByCriteria(ctx, &criteria)
	if err != nil {
		return nil, err
	}
	fmt.Println(mutasiDta.ID)
	initData["MUTASI-DTA"] = mutasiDta.ID

	// ==== Mutasi Ia =====
	mutasiIa, err := s.MutasiIaRepository.FindByCriteria(ctx, &criteria)
	if err != nil {
		return nil, err
	}
	fmt.Println(mutasiIa.ID)
	initData["MUTASI-IA"] = mutasiIa.ID

	// ==== Mutasi Rua =====
	mutasiRua, err := s.MutasiFaRepository.FindByCriteria(ctx, &criteria)
	if err != nil {
		return nil, err
	}
	fmt.Println(mutasiRua.ID)
	initData["MUTASI-RUA"] = mutasiRua.ID

	// ==== Employee Benefit =====
	employeeBenefit, err := s.EmployeeBenefitRepository.FindByCriteria(ctx, &criteria)
	if err != nil {
		return nil, err
	}
	fmt.Println(employeeBenefit.ID)
	// initData["EMPLOYEE-BENEFIT"] = employeeBenefit.ID
	for name := range initData {
		validationDetail := model.ValidationDetailEntityModel{}
		validationDetail.Context = ctx
		validationDetail.Name = name
		validationDetail.Status = constant.VALIDATION_STATUS_ON_PROCESS
		validationDetail.Note = constant.VALIDATION_NOTE_ON_PROCESS
		validationDetail.Period = trialBalance.Period
		validationDetail.CompanyID = trialBalance.CompanyID
		validationDetail.Versions = trialBalance.Versions
		validationDetail.ValidateBy = ctx.Auth.ID
		// validationDetail.SourceID = refID
		validation, err := s.ValidationRepository.FirstOrCreate(ctx, &validationDetail)
		if err != nil {
			return nil, err
		}
		if validation.Status != constant.VALIDATION_STATUS_BALANCE {
			result[validation.ID] = name
		}
	}
	return result, nil
}

func (s *service) Validate(ctx *abstraction.Context, jsonData *abstraction.JsonData, validationData *model.ValidationDetailEntityModel, wg *sync.WaitGroup, validasiModul *string) error {
	// defer wg.Done()
	criteriaData := model.FilterData{}
	criteriaData.Period = jsonData.Filter.Period
	criteriaData.CompanyID = jsonData.CompanyID
	criteriaData.Versions = jsonData.Filter.Versions
	criteriaData.Status = constant.MODUL_STATUS_DRAFT

	dataValidate := model.ValidationDetailEntityModel{}
	dataValidate.Context = ctx
	dataValidate.ID = validationData.ID
	dataValidate.ValidateBy = ctx.Auth.ID
	dataValidate.Status = constant.MODUL_STATUS_DRAFT
	dataValidate.Note = "Data sudah tervalidasi!"
	dataValidate.ModifiedAt = utilDate.DateTodayLocal()
	refID := 0

	switch *validasiModul {
	case "AGING-UTANG-PIUTANG":
		tmpData, err := s.AgingUtangPiutangRepository.FindByCriteria(ctx, &criteriaData)
		if err != nil || tmpData.ID == 0 {
			log.Println(err)
			msg := "Data tidak ditemukan"
			s.ErrorCause(ctx, &validationData.ID, &msg)
			return err
		}
		refID = tmpData.ID
	case "MUTASI-FA":
		tmpData, err := s.MutasiFaRepository.FindByCriteria(ctx, &criteriaData)
		if err != nil || tmpData.ID == 0 {
			log.Println(err)
			msg := "Data tidak ditemukan"
			s.ErrorCause(ctx, &validationData.ID, &msg)
			return err
		}
		refID = tmpData.ID
	case "MUTASI-IA":
		tmpData, err := s.MutasiIaRepository.FindByCriteria(ctx, &criteriaData)
		if err != nil || tmpData.ID == 0 {
			log.Println(err)
			msg := "Data tidak ditemukan"
			s.ErrorCause(ctx, &validationData.ID, &msg)
			return err
		}
		refID = tmpData.ID
	case "MUTASI-DTA":
		tmpData, err := s.MutasiDtaRepository.FindByCriteria(ctx, &criteriaData)
		if err != nil || tmpData.ID == 0 {
			log.Println(err)
			msg := "Data tidak ditemukan"
			s.ErrorCause(ctx, &validationData.ID, &msg)
			return err
		}
		refID = tmpData.ID
	case "MUTASI-PERSEDIAAN":
		tmpData, err := s.MutasiPersediaanRepository.FindByCriteria(ctx, &criteriaData)
		if err != nil || tmpData.ID == 0 {
			log.Println(err)
			msg := "Data tidak ditemukan"
			s.ErrorCause(ctx, &validationData.ID, &msg)
			return err
		}
		refID = tmpData.ID
	case "MUTASI-RUA":
		tmpData, err := s.MutasiRuaRepository.FindByCriteria(ctx, &criteriaData)
		if err != nil || tmpData.ID == 0 {
			log.Println(err)
			msg := "Data tidak ditemukan"
			s.ErrorCause(ctx, &validationData.ID, &msg)
			return err
		}
		refID = tmpData.ID
	case "EMPLOYEE-BENEFIT":
		tmpData, err := s.EmployeeBenefitRepository.FindByCriteria(ctx, &criteriaData)
		if err != nil || tmpData.ID == 0 {
			log.Println(err)
			msg := "Data tidak ditemukan"
			s.ErrorCause(ctx, &validationData.ID, &msg)
			return err
		}
		refID = tmpData.ID
	case "TRIAL-BALANCE":
		tmpData, err := s.TrialBalanceRepository.FindByCriteria(ctx, &criteriaData)
		if err != nil || tmpData.ID == 0 {
			log.Println(err)
			msg := "Data tidak ditemukan"
			s.ErrorCause(ctx, &validationData.ID, &msg)
			return err
		}
		refID = tmpData.ID
	}
	tmpData, err := s.TrialBalanceRepository.FindByCriteria(ctx, &criteriaData)
	if err != nil || tmpData.ID == 0 {
		log.Println(err)
		msg := "Data tidak ditemukan"
		s.ErrorCause(ctx, &validationData.ID, &msg)
		return err
	}
	criteriaFmtB := model.FormatterBridgesFilterModel{}
	criteriaFmtB.Source = validasiModul
	criteriaFmtB.TrxRefID = &refID
	fmtBridges, err := s.FormatterBridgesRepository.FindWithCriteria(ctx, &criteriaFmtB)

	if err != nil {
		msg := "Data tidak ditemukan"
		s.ErrorCause(ctx, &validationData.ID, &msg)
		return err
	}
	sourceTB := "TRIAL-BALANCE"
	criteriaFmtB1 := model.FormatterBridgesFilterModel{}
	criteriaFmtB1.Source = &sourceTB
	criteriaFmtB1.TrxRefID = &tmpData.ID
	fmtBridgesTB, err := s.FormatterBridgesRepository.FindWithCriterias(ctx, &criteriaFmtB1)
	if err != nil || fmtBridgesTB.ID == 0 {
		msg := "Data tidak ditemukan"
		s.ErrorCause(ctx, &validationData.ID, &msg)
		return err
	}
	balance := true
	for _, fB := range *fmtBridges {
		fmtData, err := s.FormatterRepository.FindByID(ctx, &fB.FormatterID)
		if err != nil || fmtData.ID == 0 {
			msg := "Data Formatter tidak ditemukan"
			s.ErrorCause(ctx, &validationData.ID, &msg)
			return err
		}
		// if fB.FormatterID == 21 {
		// 	continue
		// }
		tmpT := true
		criteriaFmtDetail := model.FormatterDetailFilterModel{}
		criteriaFmtDetail.IsControl = &tmpT
		criteriaFmtDetail.FormatterID = &fmtData.ID
		fmtDetail, err := s.FormatterDetailRepository.Find(ctx, &criteriaFmtDetail)
		if err != nil  {
			continue
		}
		
		for _, vFmt := range *fmtDetail {
			criteriaController := model.ControllerFilterModel{}
			criteriaController.FormatterID = &vFmt.FormatterID
			criteriaController.Code = &vFmt.Code

			controller, err := s.ControllerRepository.FindByCriteria(ctx, &criteriaController)
			if err != nil || len(*controller) == 0 {
				msg := fmt.Sprintf("Controller dengan kode %s tidak ditemukan", vFmt.Code)
				s.ErrorCause(ctx, &validationData.ID, &msg)
				return err
			}
			isValid := true
			for _, vControl := range *controller {
				tmpCodeCoa1 := vControl.Coa1
				criteriaSummaryCoa1 := model.TrialBalanceDetailFilterModel{}
				criteriaSummaryCoa1.Code = &tmpCodeCoa1
				criteriaSummaryCoa1.FormatterBridgesID = &fmtBridgesTB.ID

				summaryCoa1, err := s.TrialBalanceDetailRepository.FindSummary(ctx, &criteriaSummaryCoa1)
				if err != nil {
					msg := "Controller error on coa1. Silakan cek kembali"
					s.ErrorCause(ctx, &validationData.ID, &msg)
					return err
				}
				dataControls := 0.0
				
				if vControl.Coa1 == "" {
						// hasilCoa2 := 0.0
						splitCoa2 := strings.Split(vControl.Coa2, ".")
						for _, coa2 := range splitCoa2 {
							tmpCodeCoa1 := coa2
							criteriaSummaryCoa1 := model.TrialBalanceDetailFilterModel{}
							criteriaSummaryCoa1.Code = &tmpCodeCoa1
							criteriaSummaryCoa1.FormatterBridgesID = &fmtBridgesTB.ID
	
							summaryCoa1, err := s.TrialBalanceDetailRepository.FindSummary(ctx, &criteriaSummaryCoa1)
							if err != nil {
								msg := "Controller error on coa1. Silakan cek kembali"
								s.ErrorCause(ctx, &validationData.ID, &msg)
								return err
							}
							dataControls -= *summaryCoa1.AmountAfterAje
							
	
						}
						// dataControls = *summaryCoa1.AmountAfterAje - hasilCoa2
						// tmp = dataControls
					
				}
				splitControllerCommand := strings.Split(vControl.ControllerCommand, ".")
				if len(splitControllerCommand) < 2 || len(splitControllerCommand) > 4 {
					msg := "Terdapat format Controller Command yang salah"
					s.ErrorCause(ctx, &validationData.ID, &msg)
					return err
				}
				hasil1 := 0.0
				hasil2 := 0.0
				hasilTotal := 0.0
				if vControl.ControllerType == 1 { // coa vs coa					splitControllerCommand := strings.Split(vControl.ControllerCommand, ".")
					if len(splitControllerCommand) < 2 || len(splitControllerCommand) > 4 {
						log.Println(err)
						msg := "Controller Command Tidak Ditemukan"
						s.ErrorCause(ctx, &validationData.ID, &msg)
						return err
					}
					//find table

					// isValidatedProcess := 0
					switch strings.ToLower(splitControllerCommand[0]) {
					case "aging_utang_piutang_detail":
						dataAUP, err := s.AgingUtangPiutangRepository.FindByCriteria(ctx, &criteriaData) // data Aging Utang Piutang
						if err != nil || dataAUP.ID == 0 {
							log.Println(err)
							msg := "Data Aging Utang Piutang Tidak Ditemukan"
							s.ErrorCause(ctx, &validationData.ID, &msg)
							return err
						}
						tmpStrAUP := "AGING-UTANG-PIUTANG"
						criteriaFmtBAUP := model.FormatterBridgesFilterModel{}
						criteriaFmtBAUP.Source = &tmpStrAUP
						criteriaFmtBAUP.TrxRefID = &dataAUP.ID

						// dataFmtBAUP, err := s.FormatterBridgesRepository.FindWithCriteria(ctx, &criteriaFmtBAUP) // data Formatter Bridges Untuk Aging Utang Piutang
						// if err != nil || dataFmtBAUP.ID == 0 {
						// 	log.Println(err)
						// 	msg := "Data Aging Utang Piutang Tidak Ditemukan"
						// 	s.ErrorCause(ctx, &validationData.ID, &msg)
						// 	return err
						// }

						criteriaAUPD := model.AgingUtangPiutangDetailFilterModel{}
						criteriaAUPD.Code = &vFmt.Code
						criteriaAUPD.FormatterBridgesID = &fB.ID
						dataAUPD, err := s.AgingUPDetailRepository.FindByCriteria(ctx, &criteriaAUPD) // data Aging Utang Piutang Detail
						if err != nil || dataAUPD.ID == 0 {
							log.Println(err)
							msg := "Data Aging Utang Piutang Tidak Ditemukan"
							s.ErrorCause(ctx, &validationData.ID, &msg)
							return err
						}

						switch strings.ToLower(splitControllerCommand[1]) {
						case "piutangusaha_3rdparty":
							hasil1 = *dataAUPD.Piutangusaha3rdparty

						case "piutangusaha_berelasi":
							hasil1 = *dataAUPD.PiutangusahaBerelasi

						case "piutanglainshortterm_3rdparty":
							hasil1 = *dataAUPD.Piutanglainshortterm3rdparty

						case "piutanglainshortterm_berelasi":
							hasil1 = *dataAUPD.PiutanglainshorttermBerelasi
							
						case "piutangberelasishortterm":
							hasil1 = *dataAUPD.Piutangberelasishortterm

						case "piutanglainlongterm_3rdparty":
							hasil1 = *dataAUPD.Piutanglainshortterm3rdparty

						case "piutanglainlongterm_berelasi":
							hasil1 = *dataAUPD.PiutanglainlongtermBerelasi

						case "piutangberelasilongterm":
							hasil1 = *dataAUPD.Piutangberelasilongterm

						case "utangusaha_3rdparty":
							hasil1 = *dataAUPD.Utangusaha3rdparty

						case "utangusaha_berelasi":
							hasil1 = *dataAUPD.UtangusahaBerelasi

						case "utanglainshortterm_3rdparty":
							hasil1 = *dataAUPD.Utanglainshortterm3rdparty

						case "utanglainshortterm_berelasi":
							hasil1 = *dataAUPD.UtanglainshorttermBerelasi

						case "utangberelasishortterm":
							hasil1 = *dataAUPD.Utangberelasishortterm

						case "utanglainlongterm_3rdparty":
							hasil1 = *dataAUPD.Utanglainlongterm3rdparty

						case "utanglainlongterm_berelasi":
							hasil1 = *dataAUPD.UtanglainlongtermBerelasi

						case "utangberelasilongterm":
							hasil1 = *dataAUPD.Utangberelasilongterm
						}
					case "mutasi_fa_detail":
						dataMFA, err := s.MutasiFaRepository.FindByCriteria(ctx, &criteriaData) // data MutasiFa
						if err != nil || dataMFA.ID == 0 {
							msg := "Data Mutasi Fa Tidak Ditemukan"
							s.ErrorCause(ctx, &validationData.ID, &msg)
							return err
						}
						tmpStrMFA := "MUTASI-FA"
						criteriaFmtBMFA := model.FormatterBridgesFilterModel{}
						criteriaFmtBMFA.Source = &tmpStrMFA
						criteriaFmtBMFA.TrxRefID = &dataMFA.ID

						// dataFmtBMFA, err := s.FormatterBridgesRepository.FindWithCriteria(ctx, &criteriaFmtBMFA) // data Formatter Bridges Untuk Aging Utang Piutang
						// if err != nil || dataFmtBMFA.ID == 0 {
						// 	log.Println(err)
						// 	msg := "Data Mutasi Fa Tidak Ditemukan"
						// 	s.ErrorCause(ctx, &validationData.ID, &msg)
						// 	return err
						// }

						criteriaMFAD := model.MutasiFaDetailFilterModel{}
						criteriaMFAD.Code = &vFmt.Code
						criteriaMFAD.FormatterBridgesID = &fB.ID
						dataMFAD, err := s.MutasiFaDetailRepository.FindByCriteria(ctx, &criteriaMFAD) // data Aging Utang Piutang Detail
						if err != nil || dataMFAD.ID == 0 {
							log.Println(err)
							msg := "Data Mutasi Fa Tidak Ditemukan"
							s.ErrorCause(ctx, &validationData.ID, &msg)
							return err
						}

						switch strings.ToLower(splitControllerCommand[1]) {
						case "beginning_balance":
							hasil1 = *dataMFAD.BeginningBalance
						case "acquisition_of_subsidiary":
							hasil1 = *dataMFAD.AcquisitionOfSubsidiary
						case "additions":
							hasil1 = *dataMFAD.Additions
						case "deductions":
							hasil1 = *dataMFAD.Deductions
						case "reclassification":
							hasil1 = *dataMFAD.Reclassification
						case "revaluation":
							hasil1 = *dataMFAD.Revaluation
						case "ending_balance":
							hasil1 = *dataMFAD.EndingBalance
						}
					case "mutasi_ia_detail":
						dataMIA, err := s.MutasiIaRepository.FindByCriteria(ctx, &criteriaData) // data Mutasi Ia
						if err != nil || dataMIA.ID == 0 {
							log.Println(err)
							msg := "Data Mutasi Ia Tidak Ditemukan"
							s.ErrorCause(ctx, &validationData.ID, &msg)
							return err
						}
						tmpStrMIA := "MUTASI-IA"
						criteriaFmtBMIA := model.FormatterBridgesFilterModel{}
						criteriaFmtBMIA.Source = &tmpStrMIA
						criteriaFmtBMIA.TrxRefID = &dataMIA.ID

						// dataFmtBMIA, err := s.FormatterBridgesRepository.FindWithCriteria(ctx, &criteriaFmtBMIA) // data Formatter Bridges Untuk MutasiIa
						// if err != nil || dataFmtBMIA.ID == 0 {
						// 	log.Println(err)
						// 	msg := "Data Mutasi Ia Tidak Ditemukan"
						// 	s.ErrorCause(ctx, &validationData.ID, &msg)
						// 	return err
						// }

						criteriaMIAD := model.MutasiIaDetailFilterModel{}
						criteriaMIAD.Code = &vFmt.Code
						criteriaMIAD.FormatterBridgesID = &fB.ID
						dataMIAD, err := s.MutasiIaDetailRepository.FindByCriteria(ctx, &criteriaMIAD) // data Mutasi Ia Detail
						if err != nil || dataMIAD.ID == 0 {
							log.Println(err)
							msg := "Data Mutasi Ia Tidak Ditemukan"
							s.ErrorCause(ctx, &validationData.ID, &msg)
							return err
						}

						switch strings.ToLower(splitControllerCommand[1]) {
						case "beginning_balance":
							hasil1 = *dataMIAD.BeginningBalance
						case "acquisition_of_subsidiary":
							hasil1 = *dataMIAD.AcquisitionOfSubsidiary
						case "additions":
							hasil1 = *dataMIAD.Additions
						case "deductions":
							hasil1 = *dataMIAD.Deductions
						case "reclassification":
							hasil1 = *dataMIAD.Reclassification
						case "revaluation":
							hasil1 = *dataMIAD.Revaluation
						case "ending_balance":
							hasil1 = *dataMIAD.EndingBalance
						}
					case "mutasi_dta_detail":
						dataDTA, err := s.MutasiDtaRepository.FindByCriteria(ctx, &criteriaData) // data Mutasi Dta
						if err != nil || dataDTA.ID == 0 {
							log.Println(err)
							msg := "Data Mutasi Dta Tidak Ditemukan"
							s.ErrorCause(ctx, &validationData.ID, &msg)
							return err
						}
						tmpStrDTA := "MUTASI-DTA"
						criteriaFmtBDTA := model.FormatterBridgesFilterModel{}
						criteriaFmtBDTA.Source = &tmpStrDTA
						criteriaFmtBDTA.TrxRefID = &dataDTA.ID

						// dataFmtBDTA, err := s.FormatterBridgesRepository.FindWithCriteria(ctx, &criteriaFmtBDTA) // data Formatter Bridges Untuk Mutasi Dta
						// if err != nil || dataFmtBDTA.ID == 0 {
						// 	log.Println(err)
						// 	msg := "Data Mutasi Dta Tidak Ditemukan"
						// 	s.ErrorCause(ctx, &validationData.ID, &msg)
						// 	return err
						// }

						criteriaMDTAD := model.MutasiDtaDetailFilterModel{}
						criteriaMDTAD.Code = &vFmt.Code
						criteriaMDTAD.FormatterBridgesID = &fB.ID
						dataDTAD, err := s.MutasiDtaDetailRepository.FindByCriteria(ctx, &criteriaMDTAD) // data Mutasi Dta Detail
						if err != nil || dataDTAD.ID == 0 {
							log.Println(err)
							msg := "Data Mutasi Dta Tidak Ditemukan"
							s.ErrorCause(ctx, &validationData.ID, &msg)
							return err
						}

						switch strings.ToLower(splitControllerCommand[1]) {
						case "saldo_awal":
							hasil1 = *dataDTAD.SaldoAwal
						case "manfaat_beban_pajak":
							hasil1 = *dataDTAD.ManfaatBebanPajak
						case "oci":
							hasil1 = *dataDTAD.Oci
						case "akuisisi_entitas_anak":
							hasil1 = *dataDTAD.AkuisisiEntitasAnak
						case "dibebankan_ke_lr":
							hasil1 = *dataDTAD.DibebankanKeLr
						case "dibebankan_ke_oci":
							hasil1 = *dataDTAD.DibebankanKeOci
						case "saldo_akhir":
							hasil1 = *dataDTAD.SaldoAkhir
						}
					case "mutasi_persediaan_detail":
						dataPersediaan, err := s.MutasiPersediaanRepository.FindByCriteria(ctx, &criteriaData) // data Mutasi Persediaan
						if err != nil || dataPersediaan.ID == 0 {
							log.Println(err)
							msg := "Data Mutasi Persediaan Tidak Ditemukan"
							s.ErrorCause(ctx, &validationData.ID, &msg)
							return err
						}

						criteriaMPD := model.MutasiPersediaanDetailFilterModel{}
						criteriaMPD.Code = &vFmt.Code
						criteriaMPD.FormatterBridgesID = &fB.ID
						dataPersediaanD, err := s.MutasiPersediaanDetailRepository.FindByCriteria(ctx, &criteriaMPD) // data Mutasi Persediaan Detail
						if err != nil || dataPersediaanD.ID == 0 {
							log.Println(err)
							msg := "Data Mutasi Persediaan Tidak Ditemukan"
							s.ErrorCause(ctx, &validationData.ID, &msg)
							return err
						}

						switch strings.ToLower(splitControllerCommand[1]) {
						case "amount":
							hasil1 = *dataPersediaanD.Amount
						}
						switch strings.ToLower(splitControllerCommand[2]) {
						case "4":

							index, err := strconv.Atoi(splitControllerCommand[2])
							if err != nil {
								fmt.Println("Terjadi kesalahan dalam mengonversi string menjadi int:", err)
								return err
							}
							tmpStrPersediaan := "MUTASI-PERSEDIAAN"
							criteriaFmtBPersediaan := model.FormatterBridgesFilterModel{}
							criteriaFmtBPersediaan.Source = &tmpStrPersediaan
							criteriaFmtBPersediaan.TrxRefID = &dataPersediaan.ID
							criteriaFmtBPersediaan.FormatterID = &index

							dataFmtBPersediaan, err := s.FormatterBridgesRepository.FindWithCriterias(ctx, &criteriaFmtBPersediaan) // data Formatter Bridges Untuk Mutasi Persediaan
							if err != nil {
								log.Println(err)
								msg := "Data Mutasi Persediaan Tidak Ditemukan"
								s.ErrorCause(ctx, &validationData.ID, &msg)
								return err
							}
							criteriaMPD := model.MutasiPersediaanDetailFilterModel{}
							criteriaMPD.Code = &vControl.Coa2
							criteriaMPD.FormatterBridgesID = &dataFmtBPersediaan.ID
							dataPersediaanD, err := s.MutasiPersediaanDetailRepository.FindByCriteria(ctx, &criteriaMPD) // data Mutasi Persediaan Detail
							if err != nil || dataPersediaanD.ID == 0 {
								log.Println(err)
								msg := "Data Mutasi Persediaan Tidak Ditemukan"
								s.ErrorCause(ctx, &validationData.ID, &msg)
								return err
							}
							hasil2 = *dataPersediaanD.Amount
						}
						switch strings.ToLower(splitControllerCommand[2]) {
						case "5":
							indexStr := splitControllerCommand[2][19 : len(splitControllerCommand[2])-1] // Mengambil indeks dari string

							index, err := strconv.Atoi(indexStr)
							if err != nil {
								fmt.Println("Terjadi kesalahan dalam mengonversi string menjadi int:", err)
								return err
							}
							tmpStrPersediaan := "MUTASI-PERSEDIAAN"
							criteriaFmtBPersediaan := model.FormatterBridgesFilterModel{}
							criteriaFmtBPersediaan.Source = &tmpStrPersediaan
							criteriaFmtBPersediaan.TrxRefID = &dataPersediaan.ID
							criteriaFmtBPersediaan.FormatterID = &index

							dataFmtBPersediaan, err := s.FormatterBridgesRepository.FindWithCriterias(ctx, &criteriaFmtBPersediaan) // data Formatter Bridges Untuk Mutasi Persediaan
							if err != nil {
								log.Println(err)
								msg := "Data Mutasi Persediaan Tidak Ditemukan"
								s.ErrorCause(ctx, &validationData.ID, &msg)
								return err
							}
							criteriaMPD := model.MutasiPersediaanDetailFilterModel{}
							criteriaMPD.Code = &vControl.Coa2
							criteriaMPD.FormatterBridgesID = &dataFmtBPersediaan.ID
							dataPersediaanD, err := s.MutasiPersediaanDetailRepository.FindByCriteria(ctx, &criteriaMPD) // data Mutasi Persediaan Detail
							if err != nil || dataPersediaanD.ID == 0 {
								log.Println(err)
								msg := "Data Mutasi Persediaan Tidak Ditemukan"
								s.ErrorCause(ctx, &validationData.ID, &msg)
								return err
							}
							hasil2 = *dataPersediaanD.Amount
						}
						switch strings.ToLower(splitControllerCommand[3]) {
						case "-":
							hasilTotal = hasil1 - hasil2
						}
					case "mutasi_rua_detail":
						dataMRUA, err := s.MutasiRuaRepository.FindByCriteria(ctx, &criteriaData) // data Mutasi Rua
						if err != nil || dataMRUA.ID == 0 {
							log.Println(err)
							msg := "Data Mutasi Rua Tidak Ditemukan"
							s.ErrorCause(ctx, &validationData.ID, &msg)
							return err
						}
						tmpStrMRUA := "MUTASI-RUA"
						criteriaFmtBMRUA := model.FormatterBridgesFilterModel{}
						criteriaFmtBMRUA.Source = &tmpStrMRUA
						criteriaFmtBMRUA.TrxRefID = &dataMRUA.ID

						// dataFmtBMRUA, err := s.FormatterBridgesRepository.FindWithCriteria(ctx, &criteriaFmtBMRUA) // data Formatter Bridges Untuk Mutasi Rua
						// if err != nil || dataFmtBMRUA.ID == 0 {
						// 	log.Println(err)
						// 	msg := "Data Mutasi Rua Tidak Ditemukan"
						// 	s.ErrorCause(ctx, &validationData.ID, &msg)
						// 	return err
						// }

						criteriaMRUAD := model.MutasiRuaDetailFilterModel{}
						criteriaMRUAD.Code = &vFmt.Code
						criteriaMRUAD.FormatterBridgesID = &fB.ID
						dataMRUAD, err := s.MutasiRuaDetailRepository.FindByCriteria(ctx, &criteriaMRUAD) // data Mutasi Rua Detail
						if err != nil || dataMRUAD.ID == 0 {
							log.Println(err)
							msg := "Data Mutasi Rua Tidak Ditemukan"
							s.ErrorCause(ctx, &validationData.ID, &msg)
							return err
						}

						switch strings.ToLower(splitControllerCommand[1]) {
						case "beginning_balance":
							hasil1 = *dataMRUAD.BeginningBalance
						case "acquisition_of_subsidiary":
							hasil1 = *dataMRUAD.AcquisitionOfSubsidiary
						case "additions":
							hasil1 = *dataMRUAD.Additions
						case "deductions":
							hasil1 = *dataMRUAD.Deductions
						case "reclassification":
							hasil1 = *dataMRUAD.Reclassification
						case "remeasurement":
							hasil1 = *dataMRUAD.Remeasurement
						case "ending_balance":
							hasil1 = *dataMRUAD.EndingBalance
						}
					case "employee_benefit_detail":
						criteriaEBDCOA1 := model.EmployeeBenefitDetailFilterModel{}
						criteriaEBDCOA1.Code = &vControl.Coa1
						criteriaEBDCOA1.FormatterBridgesID = &fB.ID
						dataEBD1, err := s.EmployeeBenefitDetailRepository.FindByCriteria(ctx, &criteriaEBDCOA1)
						if err != nil || dataEBD1.ID == 0 {
							log.Println(err)
							msg := "Data Employee Benefit Tidak Ditemukan"
							s.ErrorCause(ctx, &validationData.ID, &msg)
							return err
						}

						criteriaEBD := model.EmployeeBenefitDetailFilterModel{}
						criteriaEBD.Code = &vControl.Coa2
						criteriaEBD.FormatterBridgesID = &fB.ID
						dataEBD, err := s.EmployeeBenefitDetailRepository.FindByCriteria(ctx, &criteriaEBD) // data Employee Benefit Detail
						if err != nil || dataEBD.ID == 0 {
							log.Println(err)
							msg := "Data Employee Benefit Tidak Ditemukan"
							s.ErrorCause(ctx, &validationData.ID, &msg)
							return err
						}
						hasilTotal = *dataEBD1.Amount - *dataEBD.Amount
					case "trial_balance_detail":
						tmpStrTB := "TRIAL-BALANCE"
						criteriaFmtBTB := model.FormatterBridgesFilterModel{}
						criteriaFmtBTB.Source = &tmpStrTB
						criteriaFmtBTB.TrxRefID = &refID

						// dataFmtBTB, err := s.FormatterBridgesRepository.FindWithCriteria(ctx, &criteriaFmtBTB) // data Formatter Bridges Untuk Trial Balance
						// if err != nil || dataFmtBTB.ID == 0 {
						// 	log.Println(err)
						// 	msg := "Data Trial Balance Tidak Ditemukan"
						// 	s.ErrorCause(ctx, &validationData.ID, &msg)
						// 	return err
						// }

						criteriaTBD := model.TrialBalanceDetailFilterModel{}
						criteriaTBD.Code = &vFmt.Code
						criteriaTBD.FormatterBridgesID = &fB.ID
						dataTBD, err := s.TrialBalanceDetailRepository.FindByCriteria(ctx, &criteriaTBD) // data Trial Balance Detail
						if err != nil || dataTBD.ID == 0 {
							log.Println(err)
							msg := "Data Trial Balance Tidak Ditemukan"
							s.ErrorCause(ctx, &validationData.ID, &msg)
							return err
						}

						switch strings.ToLower(splitControllerCommand[1]) {
						case "amount_before_aje":
							hasil1 = *dataTBD.AmountBeforeAje
						case "amount_aje_cr":
							hasil1 = *dataTBD.AmountAjeCr
						case "amount_aje_dr":
							hasil1 = *dataTBD.AmountAjeDr
						case "amount_after_aje":
							hasil1 = *dataTBD.AmountAfterAje
						}
					}

					// criteriaSummaryCoa2 := model.TrialBalanceDetailFilterModel{}
					// criteriaSummaryCoa2.Code = &vControl.Coa2
					// criteriaSummaryCoa2.FormatterBridgesID = &fmtBridgesTB.ID

					// summaryCoa2, err := s.TrialBalanceDetailRepository.FindSummary(ctx, &criteriaSummaryCoa2)
					// if err != nil {
					// 	log.Println(err)
					// 	msg := "Controller error di coa2. Silakan cek kembali"
					// 	s.ErrorCause(ctx, &validationData.ID, &msg)
					// 	return err
					// }
					tmp := 0.0

					if tmp != hasilTotal {
						// error insert to db and dont continue process
						p := message.NewPrinter(language.English)
						hasil1wSeparator := p.Sprintf("%2.f", hasil1)
						tmpwSeparator := p.Sprintf("%2.f", hasil2)
						dataValidate.Status = constant.VALIDATION_STATUS_NOT_BALANCE
						dataValidate.Note = fmt.Sprintf("Gagal memvalidasi. Terdapat ketidaksamaan nominal untuk %s (%s vs %s)", vControl.Name, hasil1wSeparator, tmpwSeparator)
						break
					}
				} else if vControl.ControllerType == 2 { // coa vs table
				if fB.Source == "MUTASI-IA" && vFmt.Code == "GAIN_(LOSS)" {
					continue
				}
				if fB.Source == "MUTASI-RUA" && vFmt.Code == "GAIN_(LOSS)" {
					continue
				}
					splitControllerCommand := strings.Split(vControl.ControllerCommand, ".")
					if len(splitControllerCommand) < 2 || len(splitControllerCommand) > 4 {
						log.Println(err)
						msg := ""
						s.ErrorCause(ctx, &validationData.ID, &msg)
						return err
					}
					// find table

					switch strings.ToLower(splitControllerCommand[0]) {
					case "aging_utang_piutang_detail":
						dataAUP, err := s.AgingUtangPiutangRepository.FindByCriteria(ctx, &criteriaData) // data Aging Utang Piutang
						if err != nil || dataAUP.ID == 0 {
							log.Println(err)
							msg := "Data Aging Utang Piutang Tidak Ditemukan"
							s.ErrorCause(ctx, &validationData.ID, &msg)
							return err
						}
						tmpStrAUP := "AGING-UTANG-PIUTANG"
						criteriaFmtBAUP := model.FormatterBridgesFilterModel{}
						criteriaFmtBAUP.Source = &tmpStrAUP
						criteriaFmtBAUP.TrxRefID = &dataAUP.ID

						// dataFmtBAUP, err := s.FormatterBridgesRepository.FindWithCriteria(ctx, &criteriaFmtBAUP) // data Formatter Bridges Untuk Aging Utang Piutang
						// if err != nil || dataFmtBAUP.ID == 0 {
						// 	log.Println(err)
						// 	msg := "Data Aging Utang Piutang Tidak Ditemukan"
						// 	s.ErrorCause(ctx, &validationData.ID, &msg)
						// 	return err
						// }

						criteriaAUPD := model.AgingUtangPiutangDetailFilterModel{}
						criteriaAUPD.Code = &vFmt.Code
						criteriaAUPD.FormatterBridgesID = &fB.ID
						dataAUPD, err := s.AgingUPDetailRepository.FindByCriteria(ctx, &criteriaAUPD) // data Aging Utang Piutang Detail
						if err != nil || dataAUPD.ID == 0 {
							log.Println(err)
							msg := "Data Aging Utang Piutang Tidak Ditemukan"
							s.ErrorCause(ctx, &validationData.ID, &msg)
							return err
						}

						switch strings.ToLower(splitControllerCommand[1]) {
						case "piutangusaha_3rdparty":
							hasil1 = *dataAUPD.Piutangusaha3rdparty
						case "piutangusaha_berelasi":
							hasil1 = *dataAUPD.PiutangusahaBerelasi
						case "piutanglainshortterm_3rdparty":
							hasil1 = *dataAUPD.Piutanglainshortterm3rdparty
						case "piutanglainshortterm_berelasi":
							hasil1 = *dataAUPD.PiutanglainshorttermBerelasi
						case "piutangberelasishortterm":
							hasil1 = *dataAUPD.Piutangberelasishortterm
						case "piutanglainlongterm_3rdparty":
							hasil1 = *dataAUPD.Piutanglainlongterm3rdparty
						case "piutanglainlongterm_berelasi":
							hasil1 = *dataAUPD.PiutanglainlongtermBerelasi
						case "piutangberelasilongterm":
							hasil1 = *dataAUPD.Piutangberelasilongterm
						case "utangusaha_3rdparty":
							hasil1 = *dataAUPD.Utangusaha3rdparty
						case "utangusaha_berelasi":
							hasil1 = *dataAUPD.UtangusahaBerelasi
						case "utanglainshortterm_3rdparty":
							hasil1 = *dataAUPD.Utanglainshortterm3rdparty
						case "utanglainshortterm_berelasi":
							hasil1 = *dataAUPD.UtanglainshorttermBerelasi
						case "utangberelasishortterm":
							hasil1 = *dataAUPD.Utangberelasishortterm
						case "utanglainlongterm_3rdparty":
							hasil1 = *dataAUPD.Utanglainlongterm3rdparty
						case "utanglainlongterm_berelasi":
							hasil1 = *dataAUPD.UtanglainlongtermBerelasi
						case "utangberelasilongterm":
							hasil1 = *dataAUPD.Utangberelasilongterm
						}
					case "mutasi_fa_detail":
						dataMFA, err := s.MutasiFaRepository.FindByCriteria(ctx, &criteriaData) // data MutasiFa
						if err != nil || dataMFA.ID == 0 {
							log.Println(err)
							msg := "Data Mutasi Fa Tidak Ditemukan"
							s.ErrorCause(ctx, &validationData.ID, &msg)
							return err
						}
						tmpStrMFA := "MUTASI-FA"
						criteriaFmtBMFA := model.FormatterBridgesFilterModel{}
						criteriaFmtBMFA.Source = &tmpStrMFA
						criteriaFmtBMFA.TrxRefID = &dataMFA.ID

						// dataFmtBMFA, err := s.FormatterBridgesRepository.FindWithCriteria(ctx, &criteriaFmtBMFA) // data Formatter Bridges Untuk Aging Utang Piutang
						// if err != nil || dataFmtBMFA.ID == 0 {
						// 	log.Println(err)
						// 	msg := "Data Mutasi Fa Tidak Ditemukan"
						// 	s.ErrorCause(ctx, &validationData.ID, &msg)
						// 	return err
						// }

						criteriaMFAD := model.MutasiFaDetailFilterModel{}
						criteriaMFAD.Code = &vFmt.Code
						criteriaMFAD.FormatterBridgesID = &fB.ID
						dataMFAD, err := s.MutasiFaDetailRepository.FindByCriteria(ctx, &criteriaMFAD) // data Aging Utang Piutang Detail
						if err != nil || dataMFAD.ID == 0 {
							log.Println(err)
							msg := "Data Mutasi Fa Tidak Ditemukan"
							s.ErrorCause(ctx, &validationData.ID, &msg)
							return err
						}

						switch strings.ToLower(splitControllerCommand[1]) {
						case "beginning_balance":
							hasil1 = *dataMFAD.BeginningBalance
						case "acquisition_of_subsidiary":
							hasil1 = *dataMFAD.AcquisitionOfSubsidiary
						case "additions":
							hasil1 = *dataMFAD.Additions
						case "deductions":
							hasil1 = *dataMFAD.Deductions
						case "reclassification":
							hasil1 = *dataMFAD.Reclassification
						case "revaluation":
							hasil1 = *dataMFAD.Revaluation
						case "ending_balance":
							hasil1 = *dataMFAD.EndingBalance
						}
					case "mutasi_ia_detail":
						if vFmt.Code == "GAIN_(LOSS)" {
							continue
						}
						dataMIA, err := s.MutasiIaRepository.FindByCriteria(ctx, &criteriaData) // data Mutasi Ia
						if err != nil || dataMIA.ID == 0 {
							log.Println(err)
							msg := "Data Mutasi Ia Tidak Ditemukan"
							s.ErrorCause(ctx, &validationData.ID, &msg)
							return err
						}
						tmpStrMIA := "MUTASI-IA"
						criteriaFmtBMIA := model.FormatterBridgesFilterModel{}
						criteriaFmtBMIA.Source = &tmpStrMIA
						criteriaFmtBMIA.TrxRefID = &dataMIA.ID

						// dataFmtBMIA, err := s.FormatterBridgesRepository.FindWithCriteria(ctx, &criteriaFmtBMIA) // data Formatter Bridges Untuk MutasiIa
						// if err != nil || dataFmtBMIA.ID == 0 {
						// 	log.Println(err)
						// 	msg := "Data Mutasi Ia Tidak Ditemukan"
						// 	s.ErrorCause(ctx, &validationData.ID, &msg)
						// 	return err
						// }

						criteriaMIAD := model.MutasiIaDetailFilterModel{}
						criteriaMIAD.Code = &vFmt.Code
						criteriaMIAD.FormatterBridgesID = &fB.ID
						dataMIAD, err := s.MutasiIaDetailRepository.FindByCriteria(ctx, &criteriaMIAD) // data Mutasi Ia Detail
						if err != nil || dataMIAD.ID == 0 {
							log.Println(err)
							msg := "Data Mutasi Ia Tidak Ditemukan"
							s.ErrorCause(ctx, &validationData.ID, &msg)
							return err
						}

						switch strings.ToLower(splitControllerCommand[1]) {
						case "beginning_balance":
							hasil1 = *dataMIAD.BeginningBalance
						case "acquisition_of_subsidiary":
							hasil1 = *dataMIAD.AcquisitionOfSubsidiary
						case "additions":
							hasil1 = *dataMIAD.Additions
						case "deductions":
							hasil1 = *dataMIAD.Deductions
						case "reclassification":
							hasil1 = *dataMIAD.Reclassification
						case "revaluation":
							hasil1 = *dataMIAD.Revaluation
						case "ending_balance":
							hasil1 = *dataMIAD.EndingBalance
						}
					case "mutasi_dta_detail":
						dataDTA, err := s.MutasiDtaRepository.FindByCriteria(ctx, &criteriaData) // data Mutasi Dta
						if err != nil || dataDTA.ID == 0 {
							log.Println(err)
							msg := "Data Mutasi Dta Tidak Ditemukan"
							s.ErrorCause(ctx, &validationData.ID, &msg)
							return err
						}
						tmpStrDTA := "MUTASI-DTA"
						criteriaFmtBDTA := model.FormatterBridgesFilterModel{}
						criteriaFmtBDTA.Source = &tmpStrDTA
						criteriaFmtBDTA.TrxRefID = &dataDTA.ID

						// dataFmtBDTA, err := s.FormatterBridgesRepository.FindWithCriteria(ctx, &criteriaFmtBDTA) // data Formatter Bridges Untuk Mutasi Dta
						// if err != nil || dataFmtBDTA.ID == 0 {
						// 	log.Println(err)
						// 	msg := "Data Mutasi Dta Tidak Ditemukan"
						// 	s.ErrorCause(ctx, &validationData.ID, &msg)
						// 	return err
						// }

						criteriaMDTAD := model.MutasiDtaDetailFilterModel{}
						criteriaMDTAD.Code = &vFmt.Code
						criteriaMDTAD.FormatterBridgesID = &fB.ID
						dataDTAD, err := s.MutasiDtaDetailRepository.FindByCriteria(ctx, &criteriaMDTAD) // data Mutasi Dta Detail
						if err != nil || dataDTAD.ID == 0 {
							log.Println(err)
							msg := "Data Mutasi Dta Tidak Ditemukan"
							s.ErrorCause(ctx, &validationData.ID, &msg)
							return err
						}

						switch strings.ToLower(splitControllerCommand[1]) {
						case "saldo_awal":
							hasil1 = *dataDTAD.SaldoAwal
						case "manfaat_beban_pajak":
							hasil1 = *dataDTAD.ManfaatBebanPajak
						case "oci":
							hasil1 = *dataDTAD.Oci
						case "akuisisi_entitas_anak":
							hasil1 = *dataDTAD.AkuisisiEntitasAnak
						case "dibebankan_ke_lr":
							hasil1 = *dataDTAD.DibebankanKeLr
						case "dibebankan_ke_oci":
							hasil1 = *dataDTAD.DibebankanKeOci
						case "saldo_akhir":
							hasil1 = *dataDTAD.SaldoAkhir
						}
					case "mutasi_persediaan_detail":
						dataPersediaan, err := s.MutasiPersediaanRepository.FindByCriteria(ctx, &criteriaData) // data Mutasi Persediaan
						if err != nil || dataPersediaan.ID == 0 {
							log.Println(err)
							msg := "Data Mutasi Persediaan Tidak Ditemukan"
							s.ErrorCause(ctx, &validationData.ID, &msg)
							return err
						}
						tmpStrPersediaan := "MUTASI-PERSEDIAAN"
						criteriaFmtBPersediaan := model.FormatterBridgesFilterModel{}
						criteriaFmtBPersediaan.Source = &tmpStrPersediaan
						criteriaFmtBPersediaan.TrxRefID = &dataPersediaan.ID

						// dataFmtBPersediaan, err := s.FormatterBridgesRepository.FindWithCriteria(ctx, &criteriaFmtBPersediaan) // data Formatter Bridges Untuk Mutasi Persediaan
						// if err != nil || dataFmtBPersediaan.ID == 0 {
						// 	log.Println(err)
						// 	msg := "Data Mutasi Persediaan Tidak Ditemukan"
						// 	s.ErrorCause(ctx, &validationData.ID, &msg)
						// 	return err
						// }

						criteriaMPD := model.MutasiPersediaanDetailFilterModel{}
						criteriaMPD.Code = &vFmt.Code
						criteriaMPD.FormatterBridgesID = &fB.ID
						dataPersediaanD, err := s.MutasiPersediaanDetailRepository.FindByCriteria(ctx, &criteriaMPD) // data Mutasi Persediaan Detail
						if err != nil || dataPersediaanD.ID == 0 {
							log.Println(err)
							msg := "Data Mutasi Persediaan Tidak Ditemukan"
							s.ErrorCause(ctx, &validationData.ID, &msg)
							return err
						}

						switch strings.ToLower(splitControllerCommand[1]) {
						case "amount":
							hasil1 = *dataPersediaanD.Amount
						}
					case "mutasi_rua_detail":
						if vFmt.Code == "GAIN_(LOSS)" {
							continue
						}
						dataMRUA, err := s.MutasiRuaRepository.FindByCriteria(ctx, &criteriaData) // data Mutasi Rua
						if err != nil || dataMRUA.ID == 0 {
							log.Println(err)
							msg := "Data Mutasi Rua Tidak Ditemukan"
							s.ErrorCause(ctx, &validationData.ID, &msg)
							return err
						}
						// tmpStrMRUA := "MUTASI-RUA"
						// criteriaFmtBMRUA := model.FormatterBridgesFilterModel{}
						// criteriaFmtBMRUA.Source = &tmpStrMRUA
						// criteriaFmtBMRUA.TrxRefID = &dataMRUA.ID

						// dataFmtBMRUA, err := s.FormatterBridgesRepository.FindWithCriteria(ctx, &criteriaFmtBMRUA) // data Formatter Bridges Untuk Mutasi Rua
						// if err != nil || dataFmtBMRUA.ID == 0 {
						// 	log.Println(err)
						// 	msg := "Data Mutasi Rua Tidak Ditemukan"
						// 	s.ErrorCause(ctx, &validationData.ID, &msg)
						// 	return err
						// }

						criteriaMRUAD := model.MutasiRuaDetailFilterModel{}
						criteriaMRUAD.Code = &vFmt.Code
						criteriaMRUAD.FormatterBridgesID = &fB.ID
						dataMRUAD, err := s.MutasiRuaDetailRepository.FindByCriteria(ctx, &criteriaMRUAD) // data Mutasi Rua Detail
						if err != nil || dataMRUAD.ID == 0 {
							log.Println(err)
							msg := "Data Mutasi Rua Tidak Ditemukan"
							s.ErrorCause(ctx, &validationData.ID, &msg)
							return err
						}

						switch strings.ToLower(splitControllerCommand[1]) {
						case "beginning_balance":
							hasil1 = *dataMRUAD.BeginningBalance
						case "acquisition_of_subsidiary":
							hasil1 = *dataMRUAD.AcquisitionOfSubsidiary
						case "additions":
							hasil1 = *dataMRUAD.Additions
						case "deductions":
							hasil1 = *dataMRUAD.Deductions
						case "reclassification":
							hasil1 = *dataMRUAD.Reclassification
						case "remeasurement":
							hasil1 = *dataMRUAD.Remeasurement
						case "ending_balance":
							hasil1 = *dataMRUAD.EndingBalance
						}
					case "employee_benefit_detail":
						dataEB, err := s.EmployeeBenefitRepository.FindByCriteria(ctx, &criteriaData) // data Employee Benefit
						if err != nil || dataEB.ID == 0 {
							log.Println(err)
							msg := "Data Employee Benefit Tidak Ditemukan"
							s.ErrorCause(ctx, &validationData.ID, &msg)
							return err
						}
						tmpStrEB := "EMPLOYEE-BENEFIT"
						criteriaFmtBEB := model.FormatterBridgesFilterModel{}
						criteriaFmtBEB.Source = &tmpStrEB
						criteriaFmtBEB.TrxRefID = &dataEB.ID

						dataFmtBEB, err := s.FormatterBridgesRepository.FindWithCriterias(ctx, &criteriaFmtBEB) // data Formatter Bridges Untuk Employee Benefit
						if err != nil || dataFmtBEB.ID == 0 {
							log.Println(err)
							msg := "Data Employee Benefit Tidak Ditemukan"
							s.ErrorCause(ctx, &validationData.ID, &msg)
							return err
						}

						criteriaEBD := model.EmployeeBenefitDetailFilterModel{}
						criteriaEBD.FormatterBridgesID = &dataFmtBEB.ID
						criteriaEBD.Code = &vFmt.Code
						dataEBD, err := s.EmployeeBenefitDetailRepository.FindWithCodes(ctx, &fB.ID ,&vFmt.Code) // data Employee Benefit Detail
						if err != nil {
							log.Println(err)
							msg := "Data Employee Benefit Tidak Ditemukan"
							s.ErrorCause(ctx, &validationData.ID, &msg)
							return err
						}

						switch strings.ToLower(splitControllerCommand[1]) {
						case "amount":
							hasil1 = *dataEBD.Amount
						}
					case "trial_balance_detail":
						tmpStrTB := "TRIAL-BALANCE"
						criteriaFmtBTB := model.FormatterBridgesFilterModel{}
						criteriaFmtBTB.Source = &tmpStrTB
						criteriaFmtBTB.TrxRefID = &refID

						// dataFmtBTB, err := s.FormatterBridgesRepository.FindWithCriteria(ctx, &criteriaFmtBTB) // data Formatter Bridges Untuk Trial Balance
						// if err != nil || dataFmtBTB.ID == 0 {
						// 	log.Println(err)
						// 	msg := "Data Trial Balance Tidak Ditemukan"
						// 	s.ErrorCause(ctx, &validationData.ID, &msg)
						// 	return err
						// }

						criteriaTBD := model.TrialBalanceDetailFilterModel{}
						criteriaTBD.Code = &vFmt.Code
						criteriaTBD.FormatterBridgesID = &fB.ID
						dataTBD, err := s.TrialBalanceDetailRepository.FindByCriteria(ctx, &criteriaTBD) // data Trial Balance Detail
						if err != nil || dataTBD.ID == 0 {
							log.Println(err)
							msg := "Data Trial Balance Tidak Ditemukan"
							s.ErrorCause(ctx, &validationData.ID, &msg)
							return err
						}

						switch strings.ToLower(splitControllerCommand[1]) {
						case "amount_before_aje":
							hasil1 = *dataTBD.AmountBeforeAje
						case "amount_aje_cr":
							hasil1 = *dataTBD.AmountAjeCr
						case "amount_aje_dr":
							hasil1 = *dataTBD.AmountAjeDr
						case "amount_after_aje":
							hasil1 = *dataTBD.AmountAfterAje
						}
					}
				} else { // coa vs table
					splitControllerCommand := strings.Split(vControl.ControllerCommand, ".")
					if len(splitControllerCommand) < 2 || len(splitControllerCommand) > 4 {
						log.Println(err)
						msg := ""
						s.ErrorCause(ctx, &validationData.ID, &msg)
						return err
					}
					// find table

					switch strings.ToLower(splitControllerCommand[0]) {
					case "aging_utang_piutang_detail":
						dataAUP, err := s.AgingUtangPiutangRepository.FindByCriteria(ctx, &criteriaData) // data Aging Utang Piutang
						if err != nil || dataAUP.ID == 0 {
							log.Println(err)
							msg := "Data Aging Utang Piutang Tidak Ditemukan"
							s.ErrorCause(ctx, &validationData.ID, &msg)
							return err
						}
						tmpStrAUP := "AGING-UTANG-PIUTANG"
						criteriaFmtBAUP := model.FormatterBridgesFilterModel{}
						criteriaFmtBAUP.Source = &tmpStrAUP
						criteriaFmtBAUP.TrxRefID = &dataAUP.ID

						// dataFmtBAUP, err := s.FormatterBridgesRepository.FindWithCriteria(ctx, &criteriaFmtBAUP) // data Formatter Bridges Untuk Aging Utang Piutang
						// if err != nil || dataFmtBAUP.ID == 0 {
						// 	log.Println(err)
						// 	msg := "Data Aging Utang Piutang Tidak Ditemukan"
						// 	s.ErrorCause(ctx, &validationData.ID, &msg)
						// 	return err
						// }

						criteriaAUPD := model.AgingUtangPiutangDetailFilterModel{}
						criteriaAUPD.Code = &vFmt.Code
						criteriaAUPD.FormatterBridgesID = &fB.ID
						dataAUPD, err := s.AgingUPDetailRepository.FindByCriteria(ctx, &criteriaAUPD) // data Aging Utang Piutang Detail
						if err != nil || dataAUPD.ID == 0 {
							log.Println(err)
							msg := "Data Aging Utang Piutang Tidak Ditemukan"
							s.ErrorCause(ctx, &validationData.ID, &msg)
							return err
						}

						switch strings.ToLower(splitControllerCommand[1]) {
						case "piutangusaha_3rdparty":
							hasil1 = *dataAUPD.Piutangusaha3rdparty
						case "piutangusaha_berelasi":
							hasil1 = *dataAUPD.PiutangusahaBerelasi
						case "piutanglainshortterm_3rdparty":
							hasil1 = *dataAUPD.Piutanglainshortterm3rdparty
						case "piutanglainshortterm_berelasi":
							hasil1 = *dataAUPD.PiutanglainshorttermBerelasi
						case "piutangberelasishortterm":
							hasil1 = *dataAUPD.Piutangberelasishortterm
						case "piutanglainlongterm_3rdparty":
							hasil1 = *dataAUPD.Piutanglainshortterm3rdparty
						case "piutanglainlongterm_berelasi":
							hasil1 = *dataAUPD.PiutanglainlongtermBerelasi
						case "piutangberelasilongterm":
							hasil1 = *dataAUPD.Piutangberelasilongterm
						case "utangusaha_3rdparty":
							hasil1 = *dataAUPD.UtangusahaBerelasi
						case "utangusaha_berelasi":
							hasil1 = *dataAUPD.UtangusahaBerelasi
						case "utanglainshortterm_3rdparty":
							hasil1 = *dataAUPD.Utanglainshortterm3rdparty
						case "utanglainshortterm_berelasi":
							hasil1 = *dataAUPD.UtanglainshorttermBerelasi
						case "utangberelasishortterm":
							hasil1 = *dataAUPD.Utangberelasishortterm
						case "utanglainlongterm_3rdparty":
							hasil1 = *dataAUPD.Utanglainlongterm3rdparty
						case "utanglainlongterm_berelasi":
							hasil1 = *dataAUPD.UtanglainlongtermBerelasi
						case "utangberelasilongterm":
							hasil1 = *dataAUPD.Utangberelasilongterm
						}
					case "mutasi_fa_detail":
						dataMFA, err := s.MutasiFaRepository.FindByCriteria(ctx, &criteriaData) // data MutasiFa
						if err != nil || dataMFA.ID == 0 {
							log.Println(err)
							msg := "Data Mutasi Fa Tidak Ditemukan"
							s.ErrorCause(ctx, &validationData.ID, &msg)
							return err
						}
						tmpStrMFA := "MUTASI-FA"
						criteriaFmtBMFA := model.FormatterBridgesFilterModel{}
						criteriaFmtBMFA.Source = &tmpStrMFA
						criteriaFmtBMFA.TrxRefID = &dataMFA.ID


						splitControllerCommandCoa := strings.Split(vControl.Coa2, ".")
						for _, scoa := range splitControllerCommandCoa {
							if scoa == "+" || scoa =="-" {
								continue
							}
							tmpCodeCoa1 := scoa
							criteriaSummaryCoa1 := model.TrialBalanceDetailFilterModel{}
							criteriaSummaryCoa1.Code = &tmpCodeCoa1
							criteriaSummaryCoa1.FormatterBridgesID = &fmtBridgesTB.ID

							sC, err := s.TrialBalanceDetailRepository.FindSummary(ctx, &criteriaSummaryCoa1)
							if err != nil {
								msg := "Controller error on coa1. Silakan cek kembali"
								s.ErrorCause(ctx, &validationData.ID, &msg)
								return err
							}
							switch strings.ToLower(splitControllerCommandCoa[0]) {
							case "+":
								hasil2 = +*sC.AmountAfterAje
							case "-":
								hasil2 = -*sC.AmountAfterAje
							}
							
						}
						criteriaMFAD := model.MutasiFaDetailFilterModel{}
						criteriaMFAD.Code = &vFmt.Code
						criteriaMFAD.FormatterBridgesID = &fB.ID
						dataMFAD, err := s.MutasiFaDetailRepository.FindByCriteria(ctx, &criteriaMFAD) // data Aging Utang Piutang Detail
						if err != nil || dataMFAD.ID == 0 {
							log.Println(err)
							msg := "Data Mutasi Fa Tidak Ditemukan"
							s.ErrorCause(ctx, &validationData.ID, &msg)
							return err
						}

						switch strings.ToLower(splitControllerCommand[1]) {
						case "beginning_balance":
							hasil1 = *dataMFAD.BeginningBalance
						case "acquisition_of_subsidiary":
							hasil1 = *dataMFAD.AcquisitionOfSubsidiary
						case "additions":
							hasil1 = *dataMFAD.Additions
						case "deductions":
							hasil1 = *dataMFAD.Deductions
						case "reclassification":
							hasil1 = *dataMFAD.Reclassification
						case "revaluation":
							hasil1 = *dataMFAD.Revaluation
						case "ending_balance":
							hasil1 = *dataMFAD.EndingBalance
						}
					case "mutasi_ia_detail":
						dataMIA, err := s.MutasiIaRepository.FindByCriteria(ctx, &criteriaData) // data Mutasi Ia
						if err != nil || dataMIA.ID == 0 {
							log.Println(err)
							msg := "Data Mutasi Ia Tidak Ditemukan"
							s.ErrorCause(ctx, &validationData.ID, &msg)
							return err
						}
						tmpStrMIA := "MUTASI-IA"
						criteriaFmtBMIA := model.FormatterBridgesFilterModel{}
						criteriaFmtBMIA.Source = &tmpStrMIA
						criteriaFmtBMIA.TrxRefID = &dataMIA.ID

						// dataFmtBMIA, err := s.FormatterBridgesRepository.FindWithCriteria(ctx, &criteriaFmtBMIA) // data Formatter Bridges Untuk MutasiIa
						// if err != nil || dataFmtBMIA.ID == 0 {
						// 	log.Println(err)
						// 	msg := "Data Mutasi Ia Tidak Ditemukan"
						// 	s.ErrorCause(ctx, &validationData.ID, &msg)
						// 	return err
						// }

						criteriaMIAD := model.MutasiIaDetailFilterModel{}
						criteriaMIAD.Code = &vFmt.Code
						criteriaMIAD.FormatterBridgesID = &fB.ID
						dataMIAD, err := s.MutasiIaDetailRepository.FindByCriteria(ctx, &criteriaMIAD) // data Mutasi Ia Detail
						if err != nil || dataMIAD.ID == 0 {
							log.Println(err)
							msg := "Data Mutasi Ia Tidak Ditemukan"
							s.ErrorCause(ctx, &validationData.ID, &msg)
							return err
						}

						switch strings.ToLower(splitControllerCommand[1]) {
						case "beginning_balance":
							hasil1 = *dataMIAD.BeginningBalance
						case "acquisition_of_subsidiary":
							hasil1 = *dataMIAD.AcquisitionOfSubsidiary
						case "additions":
							hasil1 = *dataMIAD.Additions
						case "deductions":
							hasil1 = *dataMIAD.Deductions
						case "reclassification":
							hasil1 = *dataMIAD.Reclassification
						case "revaluation":
							hasil1 = *dataMIAD.Revaluation
						case "ending_balance":
							hasil1 = *dataMIAD.EndingBalance
						}
					case "mutasi_dta_detail":
						dataDTA, err := s.MutasiDtaRepository.FindByCriteria(ctx, &criteriaData) // data Mutasi Dta
						if err != nil || dataDTA.ID == 0 {
							log.Println(err)
							msg := "Data Mutasi Dta Tidak Ditemukan"
							s.ErrorCause(ctx, &validationData.ID, &msg)
							return err
						}
						tmpStrDTA := "MUTASI-DTA"
						criteriaFmtBDTA := model.FormatterBridgesFilterModel{}
						criteriaFmtBDTA.Source = &tmpStrDTA
						criteriaFmtBDTA.TrxRefID = &dataDTA.ID

						splitControllerCommandCoa := strings.Split(vControl.Coa2, ".")
						for _, scoa := range splitControllerCommandCoa {
							if scoa == "+" || scoa =="-" {
								continue
							}
							tmpCodeCoa1 := scoa
							criteriaSummaryCoa1 := model.TrialBalanceDetailFilterModel{}
							criteriaSummaryCoa1.Code = &tmpCodeCoa1
							criteriaSummaryCoa1.FormatterBridgesID = &fmtBridgesTB.ID

							sC, err := s.TrialBalanceDetailRepository.FindWithCodes(ctx, &fmtBridgesTB.ID, &tmpCodeCoa1)
							if err != nil {
								msg := "Controller error on coa1. Silakan cek kembali"
								s.ErrorCause(ctx, &validationData.ID, &msg)
								return err
							}
							
							switch strings.ToLower(splitControllerCommandCoa[0]) {
							case "+":
								hasil2 = hasil2+*sC.AmountAfterAje
							case "-":
								nilaiString := strconv.FormatFloat(*sC.AmountAfterAje, 'f', 2, 64)
								split := strings.Split(nilaiString, "-")
								if len(split) > 0 {
								hasil2 = hasil2-*sC.AmountAfterAje*-1
								break
								}else{
									hasil2 = hasil2-*sC.AmountAfterAje
								}
							}
							
						}

						criteriaMDTAD := model.MutasiDtaDetailFilterModel{}
						criteriaMDTAD.Code = &vFmt.Code
						criteriaMDTAD.FormatterBridgesID = &fB.ID
						dataDTAD, err := s.MutasiDtaDetailRepository.FindByCriteria(ctx, &criteriaMDTAD) // data Mutasi Dta Detail
						if err != nil || dataDTAD.ID == 0 {
							log.Println(err)
							msg := "Data Mutasi Dta Tidak Ditemukan"
							s.ErrorCause(ctx, &validationData.ID, &msg)
							return err
						}

						switch strings.ToLower(splitControllerCommand[1]) {
						case "saldo_awal":
							hasil1 = *dataDTAD.SaldoAwal
						case "manfaat_beban_pajak":
							hasil1 = *dataDTAD.ManfaatBebanPajak
						case "oci":
							hasil1 = *dataDTAD.Oci
						case "akuisisi_entitas_anak":
							hasil1 = *dataDTAD.AkuisisiEntitasAnak
						case "dibebankan_ke_lr":
							hasil1 = *dataDTAD.DibebankanKeLr
						case "dibebankan_ke_oci":
							hasil1 = *dataDTAD.DibebankanKeOci
						case "saldo_akhir":
							hasil1 = *dataDTAD.SaldoAkhir
						}	
					case "mutasi_persediaan_detail":
						dataPersediaan, err := s.MutasiPersediaanRepository.FindByCriteria(ctx, &criteriaData) // data Mutasi Persediaan
						if err != nil || dataPersediaan.ID == 0 {
							log.Println(err)
							msg := "Data Mutasi Persediaan Tidak Ditemukan"
							s.ErrorCause(ctx, &validationData.ID, &msg)
							return err
						}
						tmpStrPersediaan := "MUTASI-PERSEDIAAN"
						criteriaFmtBPersediaan := model.FormatterBridgesFilterModel{}
						criteriaFmtBPersediaan.Source = &tmpStrPersediaan
						criteriaFmtBPersediaan.TrxRefID = &dataPersediaan.ID

						// dataFmtBPersediaan, err := s.FormatterBridgesRepository.FindWithCriteria(ctx, &criteriaFmtBPersediaan) // data Formatter Bridges Untuk Mutasi Persediaan
						// if err != nil || dataFmtBPersediaan.ID == 0 {
						// 	log.Println(err)
						// 	msg := "Data Mutasi Persediaan Tidak Ditemukan"
						// 	s.ErrorCause(ctx, &validationData.ID, &msg)
						// 	return err
						// }

						criteriaMPD := model.MutasiPersediaanDetailFilterModel{}
						criteriaMPD.Code = &vFmt.Code
						criteriaMPD.FormatterBridgesID = &fB.ID
						dataPersediaanD, err := s.MutasiPersediaanDetailRepository.FindByCriteria(ctx, &criteriaMPD) // data Mutasi Persediaan Detail
						if err != nil || dataPersediaanD.ID == 0 {
							log.Println(err)
							msg := "Data Mutasi Persediaan Tidak Ditemukan"
							s.ErrorCause(ctx, &validationData.ID, &msg)
							return err
						}

						switch strings.ToLower(splitControllerCommand[1]) {
						case "amount":
							hasil1 = *dataPersediaanD.Amount
						}
					case "mutasi_rua_detail":
						dataMRUA, err := s.MutasiRuaRepository.FindByCriteria(ctx, &criteriaData) // data Mutasi Rua
						if err != nil || dataMRUA.ID == 0 {
							log.Println(err)
							msg := "Data Mutasi Rua Tidak Ditemukan"
							s.ErrorCause(ctx, &validationData.ID, &msg)
							return err
						}
						// tmpStrMRUA := "MUTASI-RUA"
						// criteriaFmtBMRUA := model.FormatterBridgesFilterModel{}
						// criteriaFmtBMRUA.Source = &tmpStrMRUA
						// criteriaFmtBMRUA.TrxRefID = &dataMRUA.ID

						// dataFmtBMRUA, err := s.FormatterBridgesRepository.FindWithCriteria(ctx, &criteriaFmtBMRUA) // data Formatter Bridges Untuk Mutasi Rua
						// if err != nil || dataFmtBMRUA.ID == 0 {
						// 	log.Println(err)
						// 	msg := "Data Mutasi Rua Tidak Ditemukan"
						// 	s.ErrorCause(ctx, &validationData.ID, &msg)
						// 	return err
						// }

						criteriaMRUAD := model.MutasiRuaDetailFilterModel{}
						criteriaMRUAD.Code = &vFmt.Code
						criteriaMRUAD.FormatterBridgesID = &fB.ID
						dataMRUAD, err := s.MutasiRuaDetailRepository.FindByCriteria(ctx, &criteriaMRUAD) // data Mutasi Rua Detail
						if err != nil || dataMRUAD.ID == 0 {
							log.Println(err)
							msg := "Data Mutasi Rua Tidak Ditemukan"
							s.ErrorCause(ctx, &validationData.ID, &msg)
							return err
						}

						switch strings.ToLower(splitControllerCommand[1]) {
						case "beginning_balance":
							hasil1 = *dataMRUAD.BeginningBalance
						case "acquisition_of_subsidiary":
							hasil1 = *dataMRUAD.AcquisitionOfSubsidiary
						case "additions":
							hasil1 = *dataMRUAD.Additions
						case "deductions":
							hasil1 = *dataMRUAD.Deductions
						case "reclassification":
							hasil1 = *dataMRUAD.Reclassification
						case "remeasurement":
							hasil1 = *dataMRUAD.Remeasurement
						case "ending_balance":
							hasil1 = *dataMRUAD.EndingBalance
						}
					case "employee_benefit_detail":
						dataEB, err := s.EmployeeBenefitRepository.FindByCriteria(ctx, &criteriaData) // data Employee Benefit
						if err != nil || dataEB.ID == 0 {
							log.Println(err)
							msg := "Data Employee Benefit Tidak Ditemukan"
							s.ErrorCause(ctx, &validationData.ID, &msg)
							return err
						}
						tmpStrEB := "EMPLOYEE-BENEFIT"
						criteriaFmtBEB := model.FormatterBridgesFilterModel{}
						criteriaFmtBEB.Source = &tmpStrEB
						criteriaFmtBEB.TrxRefID = &dataEB.ID

						// dataFmtBEB, err := s.FormatterBridgesRepository.FindWithCriteria(ctx, &criteriaFmtBEB) // data Formatter Bridges Untuk Employee Benefit
						// if err != nil || dataFmtBEB.ID == 0 {
						// 	log.Println(err)
						// 	msg := "Data Employee Benefit Tidak Ditemukan"
						// 	s.ErrorCause(ctx, &validationData.ID, &msg)
						// 	return err
						// }

						criteriaEBD := model.EmployeeBenefitDetailFilterModel{}
						criteriaEBD.Code = &vFmt.Code
						criteriaEBD.FormatterBridgesID = &fB.ID
						dataEBD, err := s.EmployeeBenefitDetailRepository.FindByCriteria(ctx, &criteriaEBD) // data Employee Benefit Detail
						if err != nil || dataEBD.ID == 0 {
							log.Println(err)
							msg := "Data Employee Benefit Tidak Ditemukan"
							s.ErrorCause(ctx, &validationData.ID, &msg)
							return err
						}

						switch strings.ToLower(splitControllerCommand[1]) {
						case "amount":
							hasil1 = *dataEBD.Amount
						}
					case "trial_balance_detail":
						tmpStrTB := "TRIAL-BALANCE"
						criteriaFmtBTB := model.FormatterBridgesFilterModel{}
						criteriaFmtBTB.Source = &tmpStrTB
						criteriaFmtBTB.TrxRefID = &refID

						// dataFmtBTB, err := s.FormatterBridgesRepository.FindWithCriteria(ctx, &criteriaFmtBTB) // data Formatter Bridges Untuk Trial Balance
						// if err != nil || dataFmtBTB.ID == 0 {
						// 	log.Println(err)
						// 	msg := "Data Trial Balance Tidak Ditemukan"
						// 	s.ErrorCause(ctx, &validationData.ID, &msg)
						// 	return err
						// }

						criteriaTBD := model.TrialBalanceDetailFilterModel{}
						criteriaTBD.Code = &vFmt.Code
						criteriaTBD.FormatterBridgesID = &fB.ID
						dataTBD, err := s.TrialBalanceDetailRepository.FindByCriteria(ctx, &criteriaTBD) // data Trial Balance Detail
						if err != nil || dataTBD.ID == 0 {
							log.Println(err)
							msg := "Data Trial Balance Tidak Ditemukan"
							s.ErrorCause(ctx, &validationData.ID, &msg)
							return err
						}

						switch strings.ToLower(splitControllerCommand[1]) {
						case "amount_before_aje":
							hasil1 = *dataTBD.AmountBeforeAje
						case "amount_aje_cr":
							hasil1 = *dataTBD.AmountAjeCr
						case "amount_aje_dr":
							hasil1 = *dataTBD.AmountAjeDr
						case "amount_after_aje":
							hasil1 = *dataTBD.AmountAfterAje
						}
					}
					if hasil1 != hasil2 {
						// error insert to db and dont continue process
						p := message.NewPrinter(language.English)
						hasil1wSeparator := p.Sprintf("%2.f", hasil1)
						tmpwSeparator := p.Sprintf("%2.f", hasil2)
						dataValidate.Status = constant.VALIDATION_STATUS_NOT_BALANCE
						dataValidate.Note = fmt.Sprintf("Gagal memvalidasi. Terdapat ketidaksamaan nominal untuk %s (%s vs %s)", vControl.Name, hasil1wSeparator, tmpwSeparator)
						break
					}
				}
				tmp := 0.0
			if summaryCoa1.AmountAfterAje != nil && *summaryCoa1.AmountAfterAje != 0 || dataControls != 0{
				if vControl.Coa2 != "" {
					math.Abs(dataControls)
					tmp = dataControls
				}else {
					tmp = *summaryCoa1.AmountAfterAje
				}
				
			}
			dataValidate.Status = constant.VALIDATION_STATUS_BALANCE
			dataValidate.Note = "Data telah tervalidasi"

			hasil1 = math.Abs(hasil1)
			tmp = math.Abs(tmp)
				if hasil1 != tmp {
					// error insert to db and dont continue process
					p := message.NewPrinter(language.English)
					hasil1wSeparator := p.Sprintf("%2.f", hasil1)
					tmpwSeparator := p.Sprintf("%2.f", tmp)
					dataValidate.Status = constant.VALIDATION_STATUS_NOT_BALANCE
					dataValidate.Note = fmt.Sprintf("Gagal memvalidasi. Terdapat ketidaksamaan nominal untuk %s (%s vs %s)", vControl.Name, hasil1wSeparator, tmpwSeparator)
					isValid = false
					break
				}	
			}
			if !isValid {
				balance = false
				break
			}
		}
		
		_, err = s.ValidationRepository.Update(ctx, &validationData.ID, &dataValidate)
	if err != nil {
		log.Println(err)
		msg := ""
		s.ErrorCause(ctx, &validationData.ID, &msg)
	}
	if !balance {
		return errors.New(constant.VALIDATION_NOTE_NOT_BALANCE)
	}
	
	}
	return nil
	
}

func (s *service) ErrorCause(ctx *abstraction.Context, validationID *int, errorMsg *string) {
	updateValidationData := model.ValidationDetailEntityModel{}
	updateValidationData.Context = ctx
	updateValidationData.Note = *errorMsg
	updateValidationData.ValidateBy = ctx.Auth.ID
	updateValidationData.Status = constant.VALIDATION_STATUS_NOT_BALANCE
	hasErr = true
	_, err := s.ValidationRepository.Update(ctx, validationID, &updateValidationData)
	if err != nil {
		log.Printf("Error Update Message for validation. Detail: %s", err.Error())
		return
	}

	// criteriaData := model.FilterData{}
	// criteriaData.Period = validationData.Period
	// criteriaData.Versions = validationData.Versions
	// criteriaData.CompanyID = validationData.CompanyID
	// trialBalance, err := s.TrialBalanceRepository.FindByCriteria(ctx, &criteriaData)
	// if err != nil {
	// 	log.Printf("Error Update TB for validation. Detail: %s", err.Error())
	// 	return
	// }

	// tableName := "trial_balance"
	// tmpStatus := constant.MODUL_STATUS_DRAFT
	// err = s.ValidationRepository.UpdateStatus(ctx, &tableName, &trialBalance.ID, &tmpStatus)
	// if err != nil {
	// 	log.Printf("Error Update TB for validation. Detail: %s", err.Error())
	// 	return
	// }
}

func (s *service) ErrorCauseOnSuccessValidation(ctx *abstraction.Context, trialBalance *model.TrialBalanceEntityModel, errorMsg *string) {
	updateValidationData := model.ValidationDetailEntityModel{}
	updateValidationData.Context = ctx
	updateValidationData.Note = *errorMsg
	updateValidationData.ValidateBy = ctx.Auth.ID
	updateValidationData.Status = constant.VALIDATION_STATUS_NOT_BALANCE
	hasErr = true
	criteriaUpdate := model.ValidationDetailFilterModel{}
	criteriaUpdate.ValidationDetailFilter.CompanyID = &trialBalance.CompanyID
	criteriaUpdate.ValidationDetailFilter.Period = &trialBalance.Period
	criteriaUpdate.ValidationDetailFilter.Versions = &trialBalance.Versions
	err := s.ValidationRepository.UpdateByCriteria(ctx, &criteriaUpdate, &updateValidationData)
	if err != nil {
		log.Printf("Error Update Message for validation. Detail: %s", err.Error())
	}

	tableName := "trial_balance"
	tmpStatus := constant.MODUL_STATUS_DRAFT
	err = s.ValidationRepository.UpdateStatus(ctx, &tableName, &trialBalance.ID, &tmpStatus)
	if err != nil {
		log.Printf("Error Update TB for validation. Detail: %s", err.Error())
		return
	}
}

func (s *service) IsSuccessValidate(ctx *abstraction.Context, payload *abstraction.JsonData, trialBalance *model.TrialBalanceEntityModel) error {
	criteriaData := model.FilterData{}
	criteriaData.Period = trialBalance.Period
	criteriaData.CompanyID = trialBalance.CompanyID
	criteriaData.Versions = trialBalance.Versions
	criteriaData.Status = constant.MODUL_STATUS_VALIDATED
	tmpStatus := constant.MODUL_STATUS_VALIDATED

	agingUtangPiutang, err := s.AgingUtangPiutangRepository.FindByCriteria(ctx, &criteriaData)
	if err != nil {
		return err
	}

	tableName := "aging_utang_piutang"
	err = s.ValidationRepository.UpdateStatus(ctx, &tableName, &agingUtangPiutang.ID, &tmpStatus)
	if err != nil {
		return err
	}
	mutasiFa, err := s.MutasiFaRepository.FindByCriteria(ctx, &criteriaData)
	if err != nil {
		return err
	}

	tableName = "mutasi_fa"
	err = s.ValidationRepository.UpdateStatus(ctx, &tableName, &mutasiFa.ID, &tmpStatus)
	if err != nil {
		return err
	}
	mutasiIa, err := s.MutasiIaRepository.FindByCriteria(ctx, &criteriaData)
	if err != nil {
		return err
	}

	tableName = "mutasi_ia"
	err = s.ValidationRepository.UpdateStatus(ctx, &tableName, &mutasiIa.ID, &tmpStatus)
	if err != nil {
		return err
	}
	mutasiDta, err := s.MutasiDtaRepository.FindByCriteria(ctx, &criteriaData)
	if err != nil {
		return err
	}

	tableName = "mutasi_dta"
	err = s.ValidationRepository.UpdateStatus(ctx, &tableName, &mutasiDta.ID, &tmpStatus)
	if err != nil {
		return err
	}
	mutasiPersediaan, err := s.MutasiPersediaanRepository.FindByCriteria(ctx, &criteriaData)
	if err != nil {
		return err
	}

	tableName = "mutasi_Persediaan"
	err = s.ValidationRepository.UpdateStatus(ctx, &tableName, &mutasiPersediaan.ID, &tmpStatus)
	if err != nil {
		return err
	}
	mutasiRua, err := s.MutasiRuaRepository.FindByCriteria(ctx, &criteriaData)
	if err != nil {
		return err
	}

	tableName = "mutasi_rua"
	err = s.ValidationRepository.UpdateStatus(ctx, &tableName, &mutasiRua.ID, &tmpStatus)
	if err != nil {
		return err
	}
	employeeBenefit, err := s.EmployeeBenefitRepository.FindByCriteria(ctx, &criteriaData)
	if err != nil {
		return err
	}

	tableName = "employee_benefit"
	err = s.ValidationRepository.UpdateStatus(ctx, &tableName, &employeeBenefit.ID, &tmpStatus)
	if err != nil {
		return err
	}

	tableName = "trial_balance"
	err = s.ValidationRepository.UpdateStatus(ctx, &tableName, &trialBalance.ID, &tmpStatus)
	if err != nil {
		return err
	}

	tableName = "investasi_tbk"
	it, err := s.InvestasiTbkRepository.FindByCriteria(ctx, &criteriaData)
	if err != nil {
		msg = "Investasi TBK not found"
		s.ErrorCauseOnSuccessValidation(ctx, trialBalance, &msg)
		return err
	}

	err = s.ValidationRepository.UpdateStatus(ctx, &tableName, &it.ID, &tmpStatus)
	if err != nil && err != gorm.ErrRecordNotFound {
		msg = "Failed update Investasi TBK"
		s.ErrorCauseOnSuccessValidation(ctx, trialBalance, &msg)
		return err
	}

	tableName = "investasi_non_tbk"
	investasinontbk, err := s.InvestasiNonTbkRepository.FindByCriteria(ctx, &criteriaData)
	if err != nil {
		msg = "Investasi Non TBK not found"
		s.ErrorCauseOnSuccessValidation(ctx, trialBalance, &msg)
		return err
	}

	err = s.ValidationRepository.UpdateStatus(ctx, &tableName, &investasinontbk.ID, &tmpStatus)
	if err != nil && err != gorm.ErrRecordNotFound {
		msg = "Failed update Investasi TBK"
		s.ErrorCauseOnSuccessValidation(ctx, trialBalance, &msg)
		return err
	}

	tableName = "pembelian_penjualan_berelasi"
	ppb, err := s.PembelianPenjualanBerelasiRepository.FindByCriteria(ctx, &criteriaData)
	if err != nil {
		msg = "Pembelian Penjualan Berelasi not found"
		s.ErrorCauseOnSuccessValidation(ctx, trialBalance, &msg)
		return err
	}

	err = s.ValidationRepository.UpdateStatus(ctx, &tableName, &ppb.ID, &tmpStatus)
	if err != nil && err != gorm.ErrRecordNotFound {
		msg = "Failed update Pembelian Penjualan Berelasi"
		s.ErrorCauseOnSuccessValidation(ctx, trialBalance, &msg)
		return err
	}

	criteriaJurnal := model.AdjustmentFilterModel{}
	criteriaJurnal.TrialBalanceID = &trialBalance.ID

	dataUpdate := model.AdjustmentEntityModel{}
	dataUpdate.Status = tmpStatus
	dataUpdate.Context = ctx
	err = s.AjeRepository.UpdateByCriteria(ctx, &criteriaJurnal, &dataUpdate)
	if err != nil && err != gorm.ErrRecordNotFound {
		msg = "Failed update jurnal adjustment"
		s.ErrorCauseOnSuccessValidation(ctx, trialBalance, &msg)
		return err
	}

	return nil
}

func (s *service) IsFailedValidate(ctx *abstraction.Context, payload *abstraction.JsonData, trialBalance *model.TrialBalanceEntityModel) error {
	criteriaData := model.FilterData{}
	criteriaData.Period = trialBalance.Period
	criteriaData.CompanyID = trialBalance.CompanyID
	criteriaData.Versions = trialBalance.Versions

	trialBalance, err := s.TrialBalanceRepository.FindByCriteria(ctx, &criteriaData)
	if err != nil {
		log.Printf("Error Update TB for validation. Detail: %s", err.Error())
		return err
	}

	tableName := "trial_balance"
	tmpStatus := constant.MODUL_STATUS_DRAFT
	err = s.ValidationRepository.UpdateStatus(ctx, &tableName, &trialBalance.ID, &tmpStatus)
	if err != nil {
		log.Printf("Error Update TB for validation. Detail: %s", err.Error())
		return err
	}

	agingUtangPiutang, err := s.AgingUtangPiutangRepository.FindByCriteria(ctx, &criteriaData)
	if err != nil {
		log.Printf("Error Update Aging Utang Piutang for validation. Detail: %s", err.Error())
		return err
	}

	tableName = "aging_utang_piutang"
	err = s.ValidationRepository.UpdateStatus(ctx, &tableName, &agingUtangPiutang.ID, &tmpStatus)
	if err != nil {
		log.Printf("Error Update Aging Utang Piutang for validation. Detail: %s", err.Error())
		return err
	}
	mutasiFa, err := s.MutasiFaRepository.FindByCriteria(ctx, &criteriaData)
	if err != nil {
		return err
	}

	tableName = "mutasi_fa"
	err = s.ValidationRepository.UpdateStatus(ctx, &tableName, &mutasiFa.ID, &tmpStatus)
	if err != nil {
		return err
	}
	mutasiIa, err := s.MutasiIaRepository.FindByCriteria(ctx, &criteriaData)
	if err != nil {
		return err
	}

	tableName = "mutasi_ia"
	err = s.ValidationRepository.UpdateStatus(ctx, &tableName, &mutasiIa.ID, &tmpStatus)
	if err != nil {
		return err
	}
	mutasiDta, err := s.MutasiDtaRepository.FindByCriteria(ctx, &criteriaData)
	if err != nil {
		return err
	}

	tableName = "mutasi_dta"
	err = s.ValidationRepository.UpdateStatus(ctx, &tableName, &mutasiDta.ID, &tmpStatus)
	if err != nil {
		return err
	}
	mutasiPersediaan, err := s.MutasiPersediaanRepository.FindByCriteria(ctx, &criteriaData)
	if err != nil {
		return err
	}

	tableName = "mutasi_persediaan"
	err = s.ValidationRepository.UpdateStatus(ctx, &tableName, &mutasiPersediaan.ID, &tmpStatus)
	if err != nil {
		return err
	}
	mutasiRua, err := s.MutasiRuaRepository.FindByCriteria(ctx, &criteriaData)
	if err != nil {
		return err
	}

	tableName = "mutasi_rua"
	err = s.ValidationRepository.UpdateStatus(ctx, &tableName, &mutasiRua.ID, &tmpStatus)
	if err != nil {
		return err
	}
	employeeBenefit, err := s.EmployeeBenefitRepository.FindByCriteria(ctx, &criteriaData)
	if err != nil {
		return err
	}

	tableName = "employee_benefit"
	err = s.ValidationRepository.UpdateStatus(ctx, &tableName, &employeeBenefit.ID, &tmpStatus)
	if err != nil {
		return err
	}

	tableName = "trial_balance"
	err = s.ValidationRepository.UpdateStatus(ctx, &tableName, &trialBalance.ID, &tmpStatus)
	if err != nil {
		return err
	}

	tableName = "investasi_tbk"
	it, err := s.InvestasiTbkRepository.FindByCriteria(ctx, &criteriaData)
	if err != nil {
		msg = "Investasi TBK not found"
		s.ErrorCauseOnSuccessValidation(ctx, trialBalance, &msg)
		return err
	}

	err = s.ValidationRepository.UpdateStatus(ctx, &tableName, &it.ID, &tmpStatus)
	if err != nil && err != gorm.ErrRecordNotFound {
		msg = "Failed update Investasi TBK"
		s.ErrorCauseOnSuccessValidation(ctx, trialBalance, &msg)
		return err
	}

	tableName = "investasi_non_tbk"
	investasinontbk, err := s.InvestasiNonTbkRepository.FindByCriteria(ctx, &criteriaData)
	if err != nil {
		msg = "Investasi Non TBK not found"
		s.ErrorCauseOnSuccessValidation(ctx, trialBalance, &msg)
		return err
	}

	err = s.ValidationRepository.UpdateStatus(ctx, &tableName, &investasinontbk.ID, &tmpStatus)
	if err != nil && err != gorm.ErrRecordNotFound {
		msg = "Failed update Investasi TBK"
		s.ErrorCauseOnSuccessValidation(ctx, trialBalance, &msg)
		return err
	}

	tableName = "pembelian_penjualan_berelasi"
	ppb, err := s.PembelianPenjualanBerelasiRepository.FindByCriteria(ctx, &criteriaData)
	if err != nil {
		msg = "Pembelian Penjualan Berelasi not found"
		s.ErrorCauseOnSuccessValidation(ctx, trialBalance, &msg)
		return err
	}

	err = s.ValidationRepository.UpdateStatus(ctx, &tableName, &ppb.ID, &tmpStatus)
	if err != nil && err != gorm.ErrRecordNotFound {
		msg = "Failed update Pembelian Penjualan Berelasi"
		s.ErrorCauseOnSuccessValidation(ctx, trialBalance, &msg)
		return err
	}

	criteriaJurnal := model.AdjustmentFilterModel{}
	criteriaJurnal.TrialBalanceID = &trialBalance.ID

	dataUpdate := model.AdjustmentEntityModel{}
	dataUpdate.Status = tmpStatus
	dataUpdate.Context = ctx
	err = s.AjeRepository.UpdateByCriteria(ctx, &criteriaJurnal, &dataUpdate)
	if err != nil && err != gorm.ErrRecordNotFound {
		msg = "Failed update jurnal adjustment"
		s.ErrorCauseOnSuccessValidation(ctx, trialBalance, &msg)
		return err
	}

	return nil
}
