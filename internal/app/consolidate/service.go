package consolidate

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"math"
	"regexp"
	"strings"
	"sync"
	"time"
	"worker-consol/internal/abstraction"
	"worker-consol/internal/factory"
	kafkaproducer "worker-consol/internal/kafka/producer"
	"worker-consol/internal/model"
	"worker-consol/internal/repository"
	"worker-consol/pkg/constant"
	utilDate "worker-consol/pkg/util/date"
	"worker-consol/pkg/util/trxmanager"

	"github.com/Knetic/govaluate"
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
	}
}

var errs []error
var errMsg string

type DataToConsolidate struct {
	ConsolidatedID     int
	MasterID           int
	ListDataID         []int
	ListConsolidatedID []int
}

type JsonData struct {
	CompanyName string
	Type        string
	Period      string
	Versions    int
	DataID      int
	Errors      string
}

var payloadConsol *DataToConsolidate

func (s *service) Process(ctx *abstraction.Context, payload *abstraction.JsonData) {
	var tmpData DataToConsolidate
	if err := json.Unmarshal([]byte(payload.Data), &tmpData); err != nil {
		fmt.Printf("Error unmarshalling. Error: %s", err.Error())
		errs = append(errs, errors.New("error unmashalling"))
		return
	}
	start := time.Now()
	payloadConsol = &tmpData
	errMsg = ""
	switch payload.Name {
	case "CONSOLIDATION":
		s.Consolidate(ctx, payload)

	case "DUPLICATE":
		s.DuplicateConsolidate(ctx, payload)

	case "COMBINE", "EDIT_COMBINE":
		s.CombineConsolidate(ctx, payload)
	}
	duration := time.Since(start)
	fmt.Println("done in", int(math.Ceil(duration.Seconds())), "seconds")
}

func (s *service) SendNotification(ctx *abstraction.Context, payload *abstraction.JsonData, consolidationData *model.ConsolidationEntityModel) {
	var jsonData JsonData

	notifData := model.NotificationEntityModel{}
	notifData.Context = ctx
	switch payload.Name {
	case "DUPLICATE":
		notifData.Description = "Proses Duplikat Konsolidasi telah selesai."
	case "COMBINE", "EDIT_COMBINE":
		notifData.Description = "Proses Kombinasi Konsolidasi telah selesai."
	default:
		notifData.Description = "Proses Konsolidasi telah selesai."
	}

	if errMsg != "" {
		notifData.Description = errMsg
		jsonErr, err := json.Marshal(errs)
		if err != nil {
			log.Println(err)
			return
		}
		notifData.Data = string(jsonErr)
	}
	tmpfalse := false
	notifData.IsOpened = &tmpfalse
	notifData.CreatedBy = ctx.Auth.ID
	notifData.CreatedAt = *utilDate.DateTodayLocal()
	notifData.Data = "{}"

	if err := trxmanager.New(s.Db).WithTrx(ctx, func(ctx *abstraction.Context) error {
		trialBalance, err := s.TrialBalanceRepository.FindByID(ctx, &payloadConsol.MasterID)
		if err != nil {
			log.Println(err)
			return err
		}

		periodTB, err := time.Parse(time.RFC3339, trialBalance.Period)
		if err != nil {
			return err
		}
		tmpPeriod := periodTB.Format("2006-01-02")

		if payload.Name == "CONSOLIDATION" && errMsg != "" {
			consol, err := s.ConsolidationRepository.FindByID(ctx, &payloadConsol.ConsolidatedID)
			if err != nil {
				log.Println(err)
				return err
			}

			updateStatusConsol := model.ConsolidationEntityModel{}
			updateStatusConsol.Status = constant.CONSOLIDATION_STATUS_DRAFT
			updateStatusConsol.Context = ctx
			err = s.ConsolidationRepository.Update(ctx, &consol.ID, &updateStatusConsol)
			if err != nil {
				errMsg = "Gagal Memperbarui Data Konsolidasi. Error: " + err.Error()
				return err
			}
		}

		jsonData.CompanyName = trialBalance.Company.Name
		jsonData.Type = "consolidation"
		jsonData.Versions = trialBalance.Versions
		jsonData.Period = tmpPeriod
		if consolidationData != nil {
			jsonData.Versions = consolidationData.ConsolidationVersions
			jsonData.Period = consolidationData.Period
			jsonData.DataID = consolidationData.ID
		}

		jsonStr, err := json.Marshal(jsonData)
		if err != nil {
			log.Println(err)
			return err
		}

		notifData.Data = string(jsonStr)
		_, err = s.NotificationRepository.Create(ctx, &notifData)
		if err != nil {
			log.Println(err)
			return err
		}

		return nil
	}); err != nil {
		fmt.Printf("Error validationing data. Detail: %s", err.Error())
		return
	}

	waktu := time.Now()
	map1 := kafkaproducer.JsonData{
		UserID:    ctx.Auth.ID,
		CompanyID: ctx.Auth.CompanyID,
		Name:      "consolidation",
		Timestamp: &waktu,
		Data:      notifData.Description,
		Filter: struct {
			Period   string
			Versions int
		}{payload.Filter.Period, payload.Filter.Versions},
	}

	jsonStr, err := json.Marshal(map1)
	if err != nil {
		log.Println(err)
		return
	}

	go kafkaproducer.NewProducer("NOTIFICATION").SendMessage("NOTIFICATION", string(jsonStr))
}

func (s *service) Consol(ctx *abstraction.Context, payload *abstraction.JsonData) (*model.ConsolidationEntityModel, error) {
	var consolidationData *model.ConsolidationEntityModel
	var wg sync.WaitGroup
	trialBalance, err := s.TrialBalanceRepository.FindByID(ctx, &payloadConsol.MasterID)
	if err != nil {
		errMsg = "Data Trial Balance yang dijadikan referensi induk tidak ditemukan. Error: " + err.Error()
		return nil, err
	}

	status := constant.MODUL_STATUS_VALIDATED
	if trialBalance.Status < status {
		errMsg = "Data Trial Balance yang dijadikan referensi induk belum divalidasi."
		return nil, errors.New(errMsg)
	}

	if payloadConsol.ConsolidatedID != 0 {
		consolidationData, err = s.ConsolidationRepository.FindByID(ctx, &payloadConsol.ConsolidatedID)
		if err != nil {
			errMsg = "Data Induk Konsolidasi tidak ditemukan"
			return nil, err
		}

		if payload.Name == "DUPLICATE" {
			//create new consolidation
			dataConsolidation := model.ConsolidationEntityModel{}
			dataConsolidation.Context = ctx
			dataConsolidation.CompanyID = consolidationData.CompanyID
			dataConsolidation.Period = consolidationData.Period
			dataConsolidation.Versions = trialBalance.Versions
			dataConsolidation.CreatedBy = ctx.Auth.ID
			dataConsolidation.CreatedAt = *utilDate.DateTodayLocal()
			newConsolidationData, err := s.ConsolidationRepository.Create(ctx, &dataConsolidation)
			if err != nil {
				errMsg = "Terdapat error pada saat membuat data baru."
				return nil, err
			}
			consolidationData = newConsolidationData
		} else {
			if consolidationData.Status == constant.CONSOLIDATION_STATUS_CONSOLIDATE {
				errMsg = "Data sudah terkonsolidasi"
				return nil, err
			}
			//update ketika versi tidak sama
			if consolidationData.Versions != trialBalance.Versions {
				if trialBalance.ConsolidationID != 0 {
					errMsg = "Data Trial Balance tidak bisa dijadikan referensi induk karena sudah dikombinasi."
					return nil, errors.New(errMsg)
				}

				criteriaTB := model.TrialBalanceFilterModel{}
				criteriaTB.CompanyID = &consolidationData.CompanyID
				criteriaTB.Period = &consolidationData.Period
				criteriaTB.Versions = &consolidationData.Versions
				oldTB, err := s.TrialBalanceRepository.FindByCriteria(ctx, &criteriaTB)
				if err != nil {
					errMsg = "Gagal memperbarui data konsolidasi induk perusahaan. Data trial balance tidak ditemukan"
					return nil, err
				}
				_, err = s.TrialBalanceRepository.SetConsolID(ctx, &oldTB.ID, nil)
				if err != nil {
					errMsg = "Gagal memperbarui data konsolidasi induk perusahaan"
					return nil, err
				}

				consolidationData.Versions = trialBalance.Versions
				err = s.ConsolidationRepository.Update(ctx, &consolidationData.ID, consolidationData)
				if err != nil {
					errMsg = "Gagal mengupdate versi yang digunakan untuk konsolidasi"
					return nil, err
				}
				_, err = s.TrialBalanceRepository.SetConsolID(ctx, &trialBalance.ID, &consolidationData.ID)
				if err != nil {
					errMsg = "Gagal memperbarui data konsolidasi induk perusahaan"
					return nil, err
				}
			}

			criteriaConsolBridge := model.ConsolidationBridgeFilterModel{}
			criteriaConsolBridge.ConsolidationID = &consolidationData.ID
			listID, err := s.ConsolidationBridgeRepository.FindListConsolBridge(ctx, &criteriaConsolBridge)
			if err != nil {
				errMsg = "Data Konsolidasi Anak Perusahaan tidak ditemukan"
				return nil, err
			}

			listTB, err := s.ConsolidationBridgeRepository.FindTBByConsolBridgeID(ctx, &listID)
			if err != nil {
				errMsg = "Data Konsolidasi Anak Perusahaan tidak ditemukan"
				return nil, err
			}

			for _, vTB := range *listTB {
				if vTB.Status == constant.MODUL_STATUS_CONSOLIDATE {
					continue
				}
				// if vTB.ConsolidationID != 0 && vTB.ConsolidationID != consolidationData.ID {
				// 	errMsg = "Gagal Memperbarui Data Konsolidasi Anak Perusahaan karena data telah digunakan untuk konsolidasi"
				// 	return nil, err
				// }
				_, err := s.TrialBalanceRepository.SetConsolID(ctx, &vTB.ID, nil)
				if err != nil {
					errMsg = "Gagal Memperbarui Data Konsolidasi Anak Perusahaan"
					return nil, err
				}
			}

			// delete consolidation bridge detail
			err = s.ConsolidationBridgeDetailRepository.DeleteByListBridgeID(ctx, &listID)
			if err != nil && err != gorm.ErrRecordNotFound {
				errMsg = "Gagal memperbarui data anak perusahaan yang dijadikan referensi untuk konsolidasi"
				return nil, err
			}

			//delete consolidation bridge
			err = s.ConsolidationBridgeRepository.DeleteByConsolID(ctx, &consolidationData.ID)
			if err != nil && err != gorm.ErrRecordNotFound {
				errMsg = "Gagal Memperbarui Data Konsolidasi Anak Perusahaan"
				return nil, err
			}

			//delete consolidation detail
			err = s.ConsolidationDetailRepository.DeleteByConsolID(ctx, &consolidationData.ID)
			if err != nil && err != gorm.ErrRecordNotFound {
				errMsg = "Gagal Memperbarui Data Konsolidasi Detil Anak Perusahaan"
				return nil, err
			}
		}
	} else {
		if trialBalance.ConsolidationID != 0 {
			errMsg = "Data Trial Balance tidak bisa dijadikan referensi induk karena sudah dikombinasi."
			return nil, errors.New(errMsg)
		}
		consolModel := model.ConsolidationEntityModel{}
		consolModel.Context = ctx
		consolModel.Period = trialBalance.Period
		consolModel.Versions = trialBalance.Versions
		consolModel.CompanyID = trialBalance.CompanyID
		consolModel.Status = constant.CONSOLIDATION_STATUS_PROCESS
		consolModel.CreatedBy = ctx.Auth.ID
		tmpFalse := false
		consolModel.IsDuplicated = &tmpFalse
		consolidationData, err = s.ConsolidationRepository.Create(ctx, &consolModel)
		if err != nil {
			errMsg = "Gagal Menambah Data Konsolidasi. Error: " + err.Error()
			return nil, err
		}
		_, err = s.TrialBalanceRepository.SetConsolID(ctx, &trialBalance.ID, &consolidationData.ID)
		if err != nil {
			errMsg = "Gagal Memperbarui Data Konsolidasi Anak Perusahaan"
			return nil, err
		}
	}
	var insertedConsolBridge []*model.ConsolidationBridgeEntityModel
	for _, vID := range payloadConsol.ListDataID {
		trialBalanceCb, err := s.TrialBalanceRepository.FindByID(ctx, &vID)
		if err != nil {
			errMsg = "Data Trial Balance anak perusahaan yang dijadikan referensi tidak ditemukan. Error: " + err.Error()
			return nil, err
		}
		if trialBalance.Period != trialBalanceCb.Period {
			errMsg = "Data Trial Balance anak perusahaan yang dijadikan referensi tidak sama dengan periode data konsolidasi. Error: " + err.Error()
			return nil, errors.New("not match trial balance data period")
		}
		if (trialBalanceCb.ConsolidationID != 0) && trialBalanceCb.ConsolidationID != consolidationData.ID {
			errMsg = "Data Trial Balance anak perusahaan sudah tidak bisa dijadikan referensi karena sudah terkait dengan data consolidation yang lain."
			return nil, errors.New("cannot use trial balance data")
		}

		_, err = s.TrialBalanceRepository.SetConsolID(ctx, &trialBalanceCb.ID, &consolidationData.ID)
		if err != nil {
			errMsg = "Gagal Memperbarui Data Konsolidasi Anak Perusahaan"
			return nil, err
		}

		modelConsolidationBridge := model.ConsolidationBridgeEntityModel{}
		modelConsolidationBridge.Versions = trialBalanceCb.Versions
		modelConsolidationBridge.CompanyID = trialBalanceCb.CompanyID
		modelConsolidationBridge.ConsolidationID = consolidationData.ID
		modelConsolidationBridge.Period = trialBalanceCb.Period
		modelConsolidationBridge.Context = ctx
		insertedConsolidationBridge, err := s.ConsolidationBridgeRepository.Create(ctx, &modelConsolidationBridge)
		if err != nil {
			errMsg = "Gagal Memperbarui Data Konsolidasi Anak Perusahaan. Error: " + err.Error()
			return nil, err
		}
		insertedConsolBridge = append(insertedConsolBridge, insertedConsolidationBridge)
	}

	for _, vID := range payloadConsol.ListConsolidatedID {
		consolData, err := s.ConsolidationRepository.FindByID(ctx, &vID)
		if err != nil {
			errMsg = "Data Konsolidasi anak perusahaan yang dijadikan referensi tidak ditemukan. Error: " + err.Error()
			return nil, err
		}
		if consolData.Period != consolidationData.Period {
			errMsg = "Data Konsolidasi anak perusahaan yang dijadikan referensi tidak sama dengan periode data konsolidasi. Error: " + err.Error()
			return nil, errors.New("not match consolidation data period")
		}
		modelConsolidationBridge := model.ConsolidationBridgeEntityModel{}
		modelConsolidationBridge.CompanyID = consolData.CompanyID
		modelConsolidationBridge.Versions = consolData.Versions
		modelConsolidationBridge.ConsolidationVersions = consolData.ConsolidationVersions
		modelConsolidationBridge.ConsolidationID = consolidationData.ID
		modelConsolidationBridge.Period = consolData.Period
		modelConsolidationBridge.Context = ctx
		insertedConsolidationBridge, err := s.ConsolidationBridgeRepository.Create(ctx, &modelConsolidationBridge)
		if err != nil {
			errMsg = "Gagal Memperbarui Data Konsolidasi Anak Perusahaan. Error: " + err.Error()
			return nil, err
		}

		insertedConsolBridge = append(insertedConsolBridge, insertedConsolidationBridge)
		// wg.Add(len(insertedConsolBridge))
	}

	for _, vConsolBridge := range insertedConsolBridge {
		s.ConsolidateBridge(ctx, consolidationData, vConsolBridge, &wg)
	}

	if payload.Name == "COMBINE" {
		s.ConsolidateTrialBalanceWithFmt(ctx, consolidationData, trialBalance, &wg)
		// s.ConsolidateTrialBalance(ctx, consolidationData, trialBalance, &wg)
	} else {
		s.CombineConsolidateTrialBalance(ctx, consolidationData, trialBalance, &wg)
	}

	return consolidationData, nil
}

func (s *service) CombineConsolidate(ctx *abstraction.Context, payload *abstraction.JsonData) {
	if err := trxmanager.New(s.Db).WithTrx(ctx, func(ctx *abstraction.Context) error {
		consolidationData, err := s.Consol(ctx, payload)
		if err != nil {
			return err
		}
		if errMsg != "" {
			return errors.New("Terjadi Kesalahan! Error: " + errMsg)
		}
		updateStatusConsol := model.ConsolidationEntityModel{}
		updateStatusConsol.Status = constant.CONSOLIDATION_STATUS_DRAFT
		updateStatusConsol.Context = ctx
		err = s.ConsolidationRepository.Update(ctx, &consolidationData.ID, &updateStatusConsol)
		if err != nil {
			errMsg = "Gagal Memperbarui Data Konsolidasi. Error: " + err.Error()
			return err
		}

		// trialBalance, err := s.TrialBalanceRepository.FindByID(ctx, &payloadConsol.MasterID)
		// if err != nil {
		// 	return err
		// }
		// updateTB := model.TrialBalanceEntityModel{}
		// updateTB.ConsolidationID = consolidationData.ID
		// updateTB.Context = ctx
		// _, err = s.TrialBalanceRepository.Update(ctx, &trialBalance.ID, &updateTB)
		// if err != nil {
		// 	return err
		// }

		// criteriaConsolBridge := model.ConsolidationBridgeFilterModel{}
		// criteriaConsolBridge.ConsolidationID = &consolidationData.ID
		// listConsolBridge, err := s.ConsolidationBridgeRepository.Find(ctx, &criteriaConsolBridge)
		// for _, vConsolBridge := range *listConsolBridge {
		// 	err := s.ConsolidationRepository.UpdateStatusModul(ctx, &vConsolBridge.CompanyID, &vConsolBridge.Period, &vConsolBridge.Versions)
		// 	if err != nil {
		// 		errMsg = "Gagal Memperbarui Status Modul Konsolidasi"
		// 		errs = append(errs, err)
		// 		return err
		// 	}
		// }

		s.SendNotification(ctx, payload, consolidationData)
		return nil
	}); err != nil {
		log.Println("combine consolidate error: ", err)
		if errMsg == "" {
			errMsg = "Terdapat error pada saat proses kombinasi konsolidasi."
		}
		s.SendNotification(ctx, payload, nil)
		return
	}
}

func (s *service) DuplicateConsolidate(ctx *abstraction.Context, payload *abstraction.JsonData) {
	if err := trxmanager.New(s.Db).WithTrx(ctx, func(ctx *abstraction.Context) error {
		consolidationData, err := s.Consol(ctx, payload)
		if err != nil {
			return err
		}
		if errMsg != "" {
			return errors.New("Terjadi Kesalahan! Error: " + errMsg)
		}
		updateStatusConsol := model.ConsolidationEntityModel{}
		updateStatusConsol.Status = constant.CONSOLIDATION_STATUS_DRAFT
		updateStatusConsol.Context = ctx
		err = s.ConsolidationRepository.Update(ctx, &consolidationData.ID, &updateStatusConsol)
		if err != nil {
			errMsg = "Gagal Memperbarui Data Konsolidasi. Error: " + err.Error()
			return err
		}
		// err = s.ConsolidationRepository.UpdateStatusModul(ctx, &consolidationData.CompanyID, &consolidationData.Period, &consolidationData.Versions)
		// if err != nil {
		// 	errMsg = "Gagal Memperbarui Modul Data Konsolidasi. Error: " + err.Error()
		// 	return err
		// }

		// criteriaConsolBridge := model.ConsolidationBridgeFilterModel{}
		// criteriaConsolBridge.ConsolidationID = &consolidationData.ID
		// listConsolBridge, err := s.ConsolidationBridgeRepository.Find(ctx, &criteriaConsolBridge)
		// for _, vConsolBridge := range *listConsolBridge {
		// 	err := s.ConsolidationRepository.UpdateStatusModul(ctx, &vConsolBridge.CompanyID, &vConsolBridge.Period, &vConsolBridge.Versions)
		// 	if err != nil {
		// 		errMsg = "Gagal Memperbarui Status Modul Konsolidasi"
		// 		errs = append(errs, err)
		// 		return err
		// 	}
		// }

		s.SendNotification(ctx, payload, consolidationData)
		return nil
	}); err != nil {
		log.Println("duplicate consolidate error: ", err)
		if errMsg == "" {
			errMsg = "Terdapat error pada saat proses duplikat konsolidasi."
		}
		s.SendNotification(ctx, payload, nil)
		return
	}
}

func (s *service) Consolidate(ctx *abstraction.Context, payload *abstraction.JsonData) {
	if err := trxmanager.New(s.Db).WithTrx(ctx, func(ctx *abstraction.Context) error {
		consolidationData, err := s.Consol(ctx, payload)
		if err != nil {
			return err
		}
		if errMsg != "" {
			return errors.New("Terjadi Kesalahan! Error: " + errMsg)
		}
		updateConsol := model.ConsolidationEntityModel{}
		updateConsol.Status = constant.CONSOLIDATION_STATUS_CONSOLIDATE
		updateConsol.Context = ctx
		err = s.ConsolidationRepository.Update(ctx, &consolidationData.ID, &updateConsol)
		if err != nil {
			errMsg = "Gagal Memperbarui Status Konsolidasi"
			return err
		}
		err = s.ConsolidationRepository.UpdateStatusModul(ctx, &consolidationData.CompanyID, &consolidationData.Period, &consolidationData.Versions)
		if err != nil {
			errMsg = "Gagal Memperbarui Modul Data Konsolidasi. Error: " + err.Error()
			return err
		}

		for _, tbID := range payloadConsol.ListDataID {
			trialBalance, err := s.TrialBalanceRepository.FindByID(ctx, &tbID)
			if err != nil {
				errMsg = "Gagal memperbarui data tb parent only karena data tidak ditemukan"
				return err
			}
			if trialBalance.Status == constant.MODUL_STATUS_VALIDATED {
				updateTBParentOnly := model.TrialBalanceEntityModel{}
				updateTBParentOnly.Status = constant.MODUL_STATUS_CONSOLIDATE
				updateTBParentOnly.Context = ctx
				_, err = s.TrialBalanceRepository.Update(ctx, &trialBalance.ID, &updateTBParentOnly)
				if err != nil {
					errMsg = "Gagal memperbarui data tb parent only"
					return err
				}
			}
		}

		// criteriaConsolBridge := model.ConsolidationBridgeFilterModel{}
		// criteriaConsolBridge.ConsolidationID = &consolidationData.ID
		// listConsolBridge, err := s.ConsolidationBridgeRepository.Find(ctx, &criteriaConsolBridge)
		// for _, vConsolBridge := range *listConsolBridge {
		// 	err := s.ConsolidationRepository.UpdateStatusModul(ctx, &vConsolBridge.CompanyID, &vConsolBridge.Period, &vConsolBridge.Versions)
		// 	if err != nil {
		// 		errMsg = "Gagal Memperbarui Status Modul Konsolidasi"
		// 		errs = append(errs, err)
		// 		return err
		// 	}
		// }

		err = s.ConsolidationRepository.UpdateStatusJurnal(ctx, &consolidationData.ID)
		if err != nil {
			errMsg = "Gagal Memperbarui Modul Data Konsolidasi. Error: " + err.Error()
			return err
		}
		s.SendNotification(ctx, payload, consolidationData)
		return nil
	}); err != nil {
		log.Println("error: ", err)
		if errMsg == "" {
			errMsg = "Terdapat error pada saat proses kombinasi konsolidasi."
		}
		s.SendNotification(ctx, payload, nil)
		return
	}
}

func (s *service) CombineConsolidateTrialBalance(ctx *abstraction.Context, consolData *model.ConsolidationEntityModel, trialBalance *model.TrialBalanceEntityModel, wg *sync.WaitGroup) {
	// defer wg.Done()
	strTB := "TRIAL-BALANCE"
	criteriaFormatterBridge := model.FormatterBridgesFilterModel{}
	criteriaFormatterBridge.TrxRefID = &trialBalance.ID
	criteriaFormatterBridge.Source = &strTB
	formatterBridge, err := s.FormatterBridgesRepository.FindWithCriteria(ctx, &criteriaFormatterBridge)
	if err != nil {
		errMsg = "Data Penghubung Trial Balance tidak ditemukan"
		return
	}

	criteriaFormatter := model.FormatterDetailFilterModel{}
	criteriaFormatter.FormatterID = &formatterBridge.FormatterID
	tmpT := true
	criteriaFormatter.IsShowView = &tmpT
	data, err := s.FormatterDetailRepository.Find(ctx, &criteriaFormatter)
	if err != nil {
		errMsg = "Data Formatter Detail tidak ditemukan"
		return
	}

	rowCode := make(map[string]map[string]float64)
	isAutoSum := make(map[string]bool)
	for _, v := range *data {
		if v.AutoSummary != nil && *v.AutoSummary {
			isAutoSum[v.Code] = true
		}
		if v.IsCoa != nil && *v.IsCoa {
			tbdetails, err := s.TrialBalanceDetailRepository.FindWithCode(ctx, &formatterBridge.ID, &v.Code)
			if err != nil {
				errMsg = "Data Detail Trial Balance tidak ditemukan. Kriteria data: " + v.Code
				return
			}
			rowCode[v.Code] = make(map[string]float64)
			rowCode[v.Code]["AmountBeforeJPM"] = 0.0
			rowCode[v.Code]["AmountJPMDr"] = 0.0
			rowCode[v.Code]["AmountJPMCr"] = 0.0
			rowCode[v.Code]["AmountAfterJPM"] = 0.0
			rowCode[v.Code]["AmountJCTEDr"] = 0.0
			rowCode[v.Code]["AmountJCTECr"] = 0.0
			rowCode[v.Code]["AmountAfterJCTE"] = 0.0
			rowCode[v.Code]["AmountCombine"] = 0.0
			rowCode[v.Code]["AmountJELIMDr"] = 0.0
			rowCode[v.Code]["AmountJELIMCr"] = 0.0
			rowCode[v.Code]["AmountConsole"] = 0.0
			rowCode[v.Code]["AmountAfterJELIM"] = 0.0
			for _, tbd := range *tbdetails {
				if strings.Contains(strings.ToUpper(tbd.Code), "_SUBTOTAL") {
					continue
				}

				tmpHeadCoa := fmt.Sprintf("%c", tbd.Code[0])
				if tmpHeadCoa == "9" {
					tmpHeadCoa = tbd.Code[:1]
				}

				modelConsolDetail := model.ConsolidationDetailEntityModel{}
				modelConsolDetail.Context = ctx
				modelConsolDetail.ConsolidationID = consolData.ID
				modelConsolDetail.Code = tbd.Code
				modelConsolDetail.Description = *tbd.Description
				modelConsolDetail.SortID = v.SortID
				modelConsolDetail.AmountBeforeJpm = tbd.AmountAfterAje
				// amountZero := 0.0
				amountAfterAje := 0.0
				if tbd.AmountAfterAje != nil {
					amountAfterAje = *tbd.AmountAfterAje
				}
				jpm, err := s.JpmDetailRepository.FindSummary(ctx, &consolData.ID, &tbd.Code)
				if err != nil && err != gorm.ErrRecordNotFound {
					errMsg = "Terjadi kesalahan pada saat mencari data JPM Detail. Error: " + err.Error()
					return
				}
				amountJpmCr := 0.0
				if jpm.BalanceSheetCr != nil && *jpm.BalanceSheetCr != 0 {
					amountJpmCr = *jpm.BalanceSheetCr
				} else if jpm.IncomeStatementCr != nil && *jpm.IncomeStatementCr != 0 {
					amountJpmCr = *jpm.IncomeStatementCr
				}

				amountJpmDr := 0.0
				if jpm.BalanceSheetDr != nil && *jpm.BalanceSheetDr != 0 {
					amountJpmDr = *jpm.BalanceSheetDr
				} else if jpm.IncomeStatementDr != nil && *jpm.IncomeStatementDr != 0 {
					amountJpmDr = *jpm.IncomeStatementDr
				}

				tmpJpm := sumByHeadCoa(tmpHeadCoa, amountAfterAje, amountJpmDr, amountJpmCr)

				modelConsolDetail.AmountJpmCr = &amountJpmCr
				modelConsolDetail.AmountJpmDr = &amountJpmDr
				// tmpJpm := amountAfterAje + amountJpmDr - amountJpmCr
				modelConsolDetail.AmountAfterJpm = &tmpJpm

				jcte, err := s.JcteDetailRepository.FindSummary(ctx, &consolData.ID, &tbd.Code)
				if err != nil && err != gorm.ErrRecordNotFound {
					errMsg = "Terjadi kesalahan pada saat mencari data JCTE Detail. Error: " + err.Error()
					return
				}

				amountJcteCr := 0.0
				if jcte.BalanceSheetCr != nil && *jcte.BalanceSheetCr != 0 {
					amountJcteCr = *jcte.BalanceSheetCr
				} else if jcte.IncomeStatementCr != nil && *jcte.IncomeStatementCr != 0 {
					amountJcteCr = *jcte.IncomeStatementCr
				}
				modelConsolDetail.AmountJcteCr = &amountJcteCr

				amountJcteDr := 0.0
				if jcte.BalanceSheetDr != nil && *jcte.BalanceSheetDr != 0 {
					amountJcteDr = *jcte.BalanceSheetDr
				} else if jcte.IncomeStatementDr != nil && *jcte.IncomeStatementDr != 0 {
					amountJcteDr = *jcte.IncomeStatementDr
				}
				modelConsolDetail.AmountJcteDr = &amountJcteDr

				// tmpJcte := tmpJpm + amountJcteDr - amountJcteCr
				tmpJcte := sumByHeadCoa(tmpHeadCoa, tmpJpm, amountJcteDr, amountJcteCr)
				modelConsolDetail.AmountAfterJcte = &tmpJcte

				criteriaConsolBridge := model.ConsolidationBridgeFilterModel{}
				criteriaConsolBridge.ConsolidationID = &consolData.ID

				listConsolidationBridge, err := s.ConsolidationBridgeRepository.FindListConsolBridge(ctx, &criteriaConsolBridge)
				if err != nil && err != gorm.ErrRecordNotFound {
					errMsg = "Terjadi kesalahan pada saat mencari data Anak Usaha. Error: " + err.Error()
					return
				}
				amountCombine := 0.0
				if listConsolidationBridge != "" {
					criteriaTBDCompany := model.TrialBalanceDetailFilterModel{}
					criteriaTBDCompany.FormatterBridgesID = &formatterBridge.ID
					sumDataCompany, err := s.ConsolidationBridgeDetailRepository.FindSummary(ctx, &listConsolidationBridge, &tbd.Code)
					if err != nil && err != gorm.ErrRecordNotFound {
						errMsg = "Terjadi kesalahan pada saat mencari jumlah anak usaha. Error: " + err.Error()
						return
					}
					if sumDataCompany.Amount != nil && *sumDataCompany.Amount != 0 {
						amountCombine = *sumDataCompany.Amount
					}
				}
				tmpCombine := tmpJcte + amountCombine
				// tmpCombine := amountCombine
				modelConsolDetail.AmountCombineSubsidiary = &tmpCombine

				jelim, err := s.JelimDetailRepository.FindSummary(ctx, &consolData.ID, &tbd.Code)
				if err != nil && err != gorm.ErrRecordNotFound {
					errMsg = "Terjadi kesalahan pada saat mencari data JELIM Detail. Error: " + err.Error()
					return
				}

				amountJelimCr := 0.0
				if jelim.BalanceSheetCr != nil && *jelim.BalanceSheetCr != 0 {
					amountJelimCr = *jelim.BalanceSheetCr
				} else if jelim.IncomeStatementCr != nil && *jelim.IncomeStatementCr != 0 {
					amountJelimCr = *jelim.IncomeStatementCr
				}

				amountJelimDr := 0.0
				if jelim.BalanceSheetDr != nil && *jelim.BalanceSheetDr != 0 {
					amountJelimDr = *jelim.BalanceSheetDr
				} else if jelim.IncomeStatementDr != nil && *jelim.IncomeStatementDr != 0 {
					amountJelimDr = *jelim.IncomeStatementDr
				}

				modelConsolDetail.AmountJelimCr = &amountJelimCr
				modelConsolDetail.AmountJelimDr = &amountJelimDr
				// tmpConsol := amountCombine + *modelConsolDetail.AmountJelimDr - *modelConsolDetail.AmountJelimCr
				tmpConsol := sumByHeadCoa(tmpHeadCoa, tmpCombine, amountJelimDr, amountJelimCr)
				modelConsolDetail.AmountConsole = &tmpConsol

				insertedConsolidationDetail, err := s.ConsolidationDetailRepository.Create(ctx, &modelConsolDetail)
				if err != nil {
					errMsg = "Terjadi kesalahan pada saat memperbarui data consol detail. Error: " + err.Error()
					return
				}
				rowCode[v.Code]["AmountBeforeJPM"] += *insertedConsolidationDetail.AmountBeforeJpm
				rowCode[v.Code]["AmountJPMDr"] += *insertedConsolidationDetail.AmountJpmDr
				rowCode[v.Code]["AmountJPMCr"] += *insertedConsolidationDetail.AmountJpmCr
				rowCode[v.Code]["AmountAfterJPM"] += *insertedConsolidationDetail.AmountAfterJpm
				rowCode[v.Code]["AmountJCTEDr"] += *insertedConsolidationDetail.AmountJcteDr
				rowCode[v.Code]["AmountJCTECr"] += *insertedConsolidationDetail.AmountJcteCr
				rowCode[v.Code]["AmountAfterJCTE"] += *insertedConsolidationDetail.AmountAfterJcte
				rowCode[v.Code]["AmountCombine"] += *insertedConsolidationDetail.AmountCombineSubsidiary
				rowCode[v.Code]["AmountJELIMDr"] += *insertedConsolidationDetail.AmountJelimDr
				rowCode[v.Code]["AmountJELIMCr"] += *insertedConsolidationDetail.AmountJelimCr
				rowCode[v.Code]["AmountConsole"] += *insertedConsolidationDetail.AmountConsole
			}

			if v.AutoSummary != nil && *v.AutoSummary {
				modelConsolDetail := model.ConsolidationDetailEntityModel{}
				modelConsolDetail.Context = ctx
				// rowCode[v.Code] = make(map[string]float64)
				tmp := rowCode[v.Code]["AmountBeforeJPM"]
				modelConsolDetail.AmountBeforeJpm = &tmp
				tmp1 := rowCode[v.Code]["AmountJPMCr"]
				modelConsolDetail.AmountJpmCr = &tmp1
				tmp2 := rowCode[v.Code]["AmountJPMDr"]
				modelConsolDetail.AmountJpmDr = &tmp2
				tmp3 := rowCode[v.Code]["AmountAfterJPM"]
				modelConsolDetail.AmountAfterJpm = &tmp3
				tmp4 := rowCode[v.Code]["AmountJCTECr"]
				modelConsolDetail.AmountJcteCr = &tmp4
				tmp5 := rowCode[v.Code]["AmountJCTEDr"]
				modelConsolDetail.AmountJcteDr = &tmp5
				tmp6 := rowCode[v.Code]["AmountAfterJCTE"]
				modelConsolDetail.AmountAfterJcte = &tmp6
				tmp7 := rowCode[v.Code]["AmountCombine"]
				modelConsolDetail.AmountCombineSubsidiary = &tmp7
				tmp8 := rowCode[v.Code]["AmountJELIMDr"]
				modelConsolDetail.AmountJelimDr = &tmp8
				tmp9 := rowCode[v.Code]["AmountJELIMCr"]
				modelConsolDetail.AmountJelimCr = &tmp9
				tmp10 := rowCode[v.Code]["AmountConsole"]
				modelConsolDetail.AmountConsole = &tmp10
				modelConsolDetail.Code = fmt.Sprintf("%s_SUBTOTAL", v.Code)
				modelConsolDetail.Description = "SUB TOTAL"
				modelConsolDetail.SortID = v.SortID
				modelConsolDetail.ConsolidationID = consolData.ID
				insertedConsolidationDetail, err := s.ConsolidationDetailRepository.Create(ctx, &modelConsolDetail)
				if err != nil {
					errMsg = "Terjadi kesalahan pada saat memperbarui perhitungan summary coa. Kriteria Data: " + v.Code + ". Error: " + err.Error()
					return
				}

				rowCode[fmt.Sprintf("%s_SUBTOTAL", v.Code)] = make(map[string]float64)
				rowCode[fmt.Sprintf("%s_SUBTOTAL", v.Code)]["AmountBeforeJPM"] += *insertedConsolidationDetail.AmountBeforeJpm
				rowCode[fmt.Sprintf("%s_SUBTOTAL", v.Code)]["AmountJPMDr"] += *insertedConsolidationDetail.AmountJpmDr
				rowCode[fmt.Sprintf("%s_SUBTOTAL", v.Code)]["AmountJPMCr"] += *insertedConsolidationDetail.AmountJpmCr
				rowCode[fmt.Sprintf("%s_SUBTOTAL", v.Code)]["AmountAfterJPM"] += *insertedConsolidationDetail.AmountAfterJpm
				rowCode[fmt.Sprintf("%s_SUBTOTAL", v.Code)]["AmountJCTEDr"] += *insertedConsolidationDetail.AmountJcteDr
				rowCode[fmt.Sprintf("%s_SUBTOTAL", v.Code)]["AmountJCTECr"] += *insertedConsolidationDetail.AmountJcteCr
				rowCode[fmt.Sprintf("%s_SUBTOTAL", v.Code)]["AmountAfterJCTE"] += *insertedConsolidationDetail.AmountAfterJcte
				rowCode[fmt.Sprintf("%s_SUBTOTAL", v.Code)]["AmountCombine"] += *insertedConsolidationDetail.AmountCombineSubsidiary
				rowCode[fmt.Sprintf("%s_SUBTOTAL", v.Code)]["AmountJELIMDr"] += *insertedConsolidationDetail.AmountJelimDr
				rowCode[fmt.Sprintf("%s_SUBTOTAL", v.Code)]["AmountJELIMCr"] += *insertedConsolidationDetail.AmountJelimCr
				rowCode[fmt.Sprintf("%s_SUBTOTAL", v.Code)]["AmountConsole"] += *insertedConsolidationDetail.AmountConsole
			}
		}

		if v.IsTotal != nil && *v.IsTotal {
			//masukan dulu ke db dengan nilai amountnya 0
			total := make(map[string]float64)
			for _, vSubtotal := range []string{"AmountBeforeJPM", "AmountJPMDr", "AmountJPMCr", "AmountAfterJPM", "AmountJCTEDr", "AmountJCTECr", "AmountAfterJCTE", "AmountCombine", "AmountJELIMDr", "AmountJELIMCr", "AmountConsole"} {
				total[vSubtotal] = 0.0
				if v.FxSummary != "" {
					formula := v.FxSummary
					reg := regexp.MustCompile(`[A-Za-z0-9_~:()]+|[0-9]+\d{3,}`)
					match := reg.FindAllString(formula, -1)
					for _, vMatch := range match {
						if _, ok := rowCode[fmt.Sprintf("%s_SUBTOTAL", vMatch)][vSubtotal]; ok {
							formula = strings.ReplaceAll(formula, vMatch, fmt.Sprintf("(%.2f)", rowCode[fmt.Sprintf("%s_SUBTOTAL", vMatch)][vSubtotal]))
						} else {
							if _, ok := rowCode[vMatch][vSubtotal]; ok {
								formula = strings.ReplaceAll(formula, vMatch, fmt.Sprintf("(%.2f)", rowCode[vMatch][vSubtotal]))
							}
						}
					}

					expression, err := govaluate.NewEvaluableExpression(formula)
					if err != nil {
						errMsg = "Terjadi kesalahan pada saat memperbarui perhitungan total. Kriteria Data: " + v.Code + ". Error: " + err.Error()
						return
					}
					result, err := expression.Evaluate(nil)
					if err != nil {
						errMsg = "Terjadi kesalahan pada saat memperbarui perhitungan total. Kriteria Data: " + v.Code + ". Error: " + err.Error()
						return
					}
					f, ok := result.(float64) // type assertion
					if !ok {
						fmt.Println("Interface value is not a float64")
						return
					}
					total[vSubtotal] = float64(int(f*100)) / 100
				}

			}

			modelConsolDetail := model.ConsolidationDetailEntityModel{}
			modelConsolDetail.Context = ctx
			modelConsolDetail.Code = v.Code
			modelConsolDetail.Description = v.Description
			modelConsolDetail.SortID = v.SortID
			modelConsolDetail.ConsolidationID = consolData.ID
			tmpamount1 := total["AmountBeforeJPM"]
			modelConsolDetail.AmountBeforeJpm = &tmpamount1
			tmpamount2 := total["AmountJPMDr"]
			modelConsolDetail.AmountJpmDr = &tmpamount2
			tmpamount3 := total["AmountJPMCr"]
			modelConsolDetail.AmountJpmCr = &tmpamount3
			tmpamount4 := total["AmountAfterJPM"]
			modelConsolDetail.AmountAfterJpm = &tmpamount4
			tmpamount5 := total["AmountJCTEDr"]
			modelConsolDetail.AmountJcteDr = &tmpamount5
			tmpamount6 := total["AmountJCTECr"]
			modelConsolDetail.AmountJcteCr = &tmpamount6
			tmpamount7 := total["AmountAfterJCTE"]
			modelConsolDetail.AmountAfterJcte = &tmpamount7
			tmpamount8 := total["AmountCombine"]
			modelConsolDetail.AmountCombineSubsidiary = &tmpamount8
			tmpamount9 := total["AmountJELIMDr"]
			modelConsolDetail.AmountJelimDr = &tmpamount9
			tmpamount10 := total["AmountJELIMCr"]
			modelConsolDetail.AmountJelimCr = &tmpamount10
			tmpamount11 := total["AmountConsole"]
			modelConsolDetail.AmountConsole = &tmpamount11
			_, err = s.ConsolidationDetailRepository.Create(ctx, &modelConsolDetail)
			if err != nil {
				errMsg = "Terjadi kesalahan pada saat menginput perhitungan total. Kriteria Data: " + v.Code + ". Error: " + err.Error()
				return
			}
			rowCode[v.Code] = make(map[string]float64)
			rowCode[v.Code]["AmountBeforeJPM"] = total["AmountBeforeJPM"]
			rowCode[v.Code]["AmountJPMCr"] = total["AmountJPMCr"]
			rowCode[v.Code]["AmountJPMDr"] = total["AmountJPMDr"]
			rowCode[v.Code]["AmountAfterJPM"] = total["AmountAfterJPM"]
			rowCode[v.Code]["AmountJCTEDr"] = total["AmountJCTEDr"]
			rowCode[v.Code]["AmountJCTECr"] = total["AmountJCTECr"]
			rowCode[v.Code]["AmountAfterJCTE"] = total["AmountAfterJCTE"]
			rowCode[v.Code]["AmountCombine"] = total["AmountCombine"]
			rowCode[v.Code]["AmountJELIMDr"] = total["AmountJELIMDr"]
			rowCode[v.Code]["AmountJELIMCr"] = total["AmountJELIMCr"]
			rowCode[v.Code]["AmountConsole"] = total["AmountConsole"]
		}
		if v.IsLabel != nil && *v.IsLabel && v.IsTotal != nil && !*v.IsTotal && v.IsCoa != nil && !*v.IsCoa {
			modelConsolDetail := model.ConsolidationDetailEntityModel{}
			modelConsolDetail.Context = ctx
			modelConsolDetail.Code = v.Code
			modelConsolDetail.Description = v.Description
			modelConsolDetail.SortID = v.SortID
			modelConsolDetail.ConsolidationID = consolData.ID
			_, err := s.ConsolidationDetailRepository.Create(ctx, &modelConsolDetail)
			if err != nil {
				errMsg = "Terjadi kesalahan pada saat memperbarui perhitungan summary coa. Kriteria Data: " + v.Code + ". Error: " + err.Error()
				return
			}
		}
	}

}

func (s *service) ConsolidateTrialBalance(ctx *abstraction.Context, consolData *model.ConsolidationEntityModel, trialBalance *model.TrialBalanceEntityModel, wg *sync.WaitGroup) {
	// defer wg.Done()
	strTB := "TRIAL-BALANCE"
	criteriaFormatterBridge := model.FormatterBridgesFilterModel{}
	criteriaFormatterBridge.TrxRefID = &trialBalance.ID
	criteriaFormatterBridge.Source = &strTB
	formatterBridge, err := s.FormatterBridgesRepository.FindWithCriteria(ctx, &criteriaFormatterBridge)
	if err != nil {
		errMsg = "Data penghubung trial balance tidak ditemukan"
		return
	}

	criteriaTBD := model.TrialBalanceDetailFilterModel{}
	criteriaTBD.FormatterBridgesID = &formatterBridge.ID

	tbdData, err := s.TrialBalanceDetailRepository.Find(ctx, &criteriaTBD)
	if err != nil {
		errMsg = "Data trial balance tidak ditemukan"
		return
	}
	for _, v := range *tbdData {
		modelConsolDetail := model.ConsolidationDetailEntityModel{}
		modelConsolDetail.Context = ctx
		modelConsolDetail.ConsolidationID = consolData.ID
		modelConsolDetail.Code = v.Code
		modelConsolDetail.Description = *v.Description
		modelConsolDetail.SortID = v.SortID
		modelConsolDetail.AmountBeforeJpm = v.AmountAfterAje
		amountZero := 0.0

		modelConsolDetail.AmountJpmCr = &amountZero
		modelConsolDetail.AmountJpmDr = &amountZero
		modelConsolDetail.AmountAfterJpm = v.AmountAfterAje
		modelConsolDetail.AmountJcteCr = &amountZero
		modelConsolDetail.AmountJcteDr = &amountZero
		modelConsolDetail.AmountAfterJcte = modelConsolDetail.AmountAfterJpm

		criteriaConsolBridge := model.ConsolidationBridgeFilterModel{}
		criteriaConsolBridge.ConsolidationID = &consolData.ID
		listConsolidationBridge, err := s.ConsolidationBridgeRepository.FindListConsolBridge(ctx, &criteriaConsolBridge)
		if err != nil && err != gorm.ErrRecordNotFound {
			errMsg = "Terjadi kesalahan pada saat mencari data Anak Usaha. Error: " + err.Error()
			return
		}
		amountCombine := 0.0
		if listConsolidationBridge != "" {
			criteriaTBDCompany := model.TrialBalanceDetailFilterModel{}
			criteriaTBDCompany.FormatterBridgesID = &formatterBridge.ID
			sumDataCompany, err := s.ConsolidationBridgeDetailRepository.FindSummary(ctx, &listConsolidationBridge, &v.Code)
			if err != nil && err != gorm.ErrRecordNotFound {
				errMsg = "Terjadi kesalahan pada saat mencari jumlah anak usaha. Error: " + err.Error()
				return
			}
			if sumDataCompany.Amount != nil && *sumDataCompany.Amount != 0 {
				amountCombine = *sumDataCompany.Amount
			}
		}

		tmpCombine := amountCombine
		if modelConsolDetail.AmountAfterJcte != nil && *modelConsolDetail.AmountAfterJcte != 0 {
			tmpCombine += *modelConsolDetail.AmountAfterJcte
		}
		modelConsolDetail.AmountCombineSubsidiary = &tmpCombine

		modelConsolDetail.AmountJelimCr = &amountZero
		modelConsolDetail.AmountJelimDr = &amountZero
		// tmpJelim := *modelConsolDetail.AmountAfterJpm + *modelConsolDetail.AmountJelimDr - *modelConsolDetail.AmountJelimCr
		modelConsolDetail.AmountConsole = &tmpCombine

		_, err = s.ConsolidationDetailRepository.Create(ctx, &modelConsolDetail)
		if err != nil {
			errMsg = "Terjadi kesalahan pada saat memperbarui data detil konsolidasi. Kriteria Data: " + v.Code + ". Error: " + err.Error()
			return
		}

	}
}

func (s *service) ConsolidateTrialBalanceWithFmt(ctx *abstraction.Context, consolData *model.ConsolidationEntityModel, trialBalance *model.TrialBalanceEntityModel, wg *sync.WaitGroup) {
	// defer wg.Done()
	strTB := "TRIAL-BALANCE"
	criteriaFormatterBridge := model.FormatterBridgesFilterModel{}
	criteriaFormatterBridge.TrxRefID = &trialBalance.ID
	criteriaFormatterBridge.Source = &strTB
	formatterBridge, err := s.FormatterBridgesRepository.FindWithCriteria(ctx, &criteriaFormatterBridge)
	if err != nil {
		errMsg = "Data Penghubung Trial Balance tidak ditemukan"
		return
	}

	criteriaFormatter := model.FormatterDetailFilterModel{}
	criteriaFormatter.FormatterID = &formatterBridge.FormatterID
	data, err := s.FormatterDetailRepository.Find(ctx, &criteriaFormatter)
	if err != nil {
		errMsg = "Data Formatter Detail tidak ditemukan"
		return
	}

	rowCode := make(map[string]map[string]float64)
	isAutoSum := make(map[string]bool)
	for _, v := range *data {
		if v.AutoSummary != nil && *v.AutoSummary {
			isAutoSum[v.Code] = true
		}
		if v.IsCoa != nil && *v.IsCoa {
			tbdetails, err := s.TrialBalanceDetailRepository.FindWithCode(ctx, &formatterBridge.ID, &v.Code)
			if err != nil {
				errMsg = "Data Detail Trial Balance tidak ditemukan. Kriteria data: " + v.Code
				return
			}
			rowCode[v.Code] = make(map[string]float64)
			rowCode[v.Code]["AmountBeforeJPM"] = 0.0
			rowCode[v.Code]["AmountJPMDr"] = 0.0
			rowCode[v.Code]["AmountJPMCr"] = 0.0
			rowCode[v.Code]["AmountAfterJPM"] = 0.0
			rowCode[v.Code]["AmountJCTEDr"] = 0.0
			rowCode[v.Code]["AmountJCTECr"] = 0.0
			rowCode[v.Code]["AmountAfterJCTE"] = 0.0
			rowCode[v.Code]["AmountCombine"] = 0.0
			rowCode[v.Code]["AmountJELIMDr"] = 0.0
			rowCode[v.Code]["AmountJELIMCr"] = 0.0
			rowCode[v.Code]["AmountConsole"] = 0.0
			rowCode[v.Code]["AmountAfterJELIM"] = 0.0
			for _, tbd := range *tbdetails {
				if strings.Contains(strings.ToUpper(tbd.Code), "_SUBTOTAL") {
					continue
				}

				tmpHeadCoa := fmt.Sprintf("%c", tbd.Code[0])
				if tmpHeadCoa == "9" {
					tmpHeadCoa = tbd.Code[:1]
				}

				modelConsolDetail := model.ConsolidationDetailEntityModel{}
				modelConsolDetail.Context = ctx
				modelConsolDetail.ConsolidationID = consolData.ID
				modelConsolDetail.Code = tbd.Code
				modelConsolDetail.Description = *tbd.Description
				modelConsolDetail.SortID = v.SortID
				modelConsolDetail.AmountBeforeJpm = tbd.AmountAfterAje
				amount := 0.0
				if tbd.AmountAfterAje != nil && *tbd.AmountAfterAje != 0 {
					amount = *tbd.AmountAfterAje
				}
				amountZero := 0.0

				modelConsolDetail.AmountJpmCr = &amountZero
				modelConsolDetail.AmountJpmDr = &amountZero
				// tmpJpm := amountAfterAje + amountJpmDr - amountJpmCr
				modelConsolDetail.AmountAfterJpm = &amount

				modelConsolDetail.AmountJcteCr = &amountZero
				modelConsolDetail.AmountJcteDr = &amountZero
				modelConsolDetail.AmountAfterJcte = &amount

				criteriaConsolBridge := model.ConsolidationBridgeFilterModel{}
				criteriaConsolBridge.ConsolidationID = &consolData.ID

				listConsolidationBridge, err := s.ConsolidationBridgeRepository.FindListConsolBridge(ctx, &criteriaConsolBridge)
				if err != nil && err != gorm.ErrRecordNotFound {
					errMsg = "Terjadi kesalahan pada saat mencari data Anak Usaha. Error: " + err.Error()
					return
				}
				amountCombine := 0.0
				if listConsolidationBridge != "" {
					criteriaTBDCompany := model.TrialBalanceDetailFilterModel{}
					criteriaTBDCompany.FormatterBridgesID = &formatterBridge.ID
					sumDataCompany, err := s.ConsolidationBridgeDetailRepository.FindSummary(ctx, &listConsolidationBridge, &tbd.Code)
					if err != nil && err != gorm.ErrRecordNotFound {
						errMsg = "Terjadi kesalahan pada saat mencari jumlah anak usaha. Error: " + err.Error()
						return
					}
					if sumDataCompany.Amount != nil && *sumDataCompany.Amount != 0 {
						amountCombine = *sumDataCompany.Amount
					}
				}
				tmpCombine := amount + amountCombine
				modelConsolDetail.AmountCombineSubsidiary = &tmpCombine

				modelConsolDetail.AmountJelimCr = &amountZero
				modelConsolDetail.AmountJelimDr = &amountZero
				modelConsolDetail.AmountConsole = &tmpCombine

				insertedConsolidationDetail, err := s.ConsolidationDetailRepository.Create(ctx, &modelConsolDetail)
				if err != nil {
					errMsg = "Terjadi kesalahan pada saat memperbarui data consol detail. Error: " + err.Error()
					return
				}
				rowCode[v.Code]["AmountBeforeJPM"] += *insertedConsolidationDetail.AmountBeforeJpm
				rowCode[v.Code]["AmountJPMDr"] += *insertedConsolidationDetail.AmountJpmDr
				rowCode[v.Code]["AmountJPMCr"] += *insertedConsolidationDetail.AmountJpmCr
				rowCode[v.Code]["AmountAfterJPM"] += *insertedConsolidationDetail.AmountAfterJpm
				rowCode[v.Code]["AmountJCTEDr"] += *insertedConsolidationDetail.AmountJcteDr
				rowCode[v.Code]["AmountJCTECr"] += *insertedConsolidationDetail.AmountJcteCr
				rowCode[v.Code]["AmountAfterJCTE"] += *insertedConsolidationDetail.AmountAfterJcte
				rowCode[v.Code]["AmountCombine"] += *insertedConsolidationDetail.AmountCombineSubsidiary
				rowCode[v.Code]["AmountJELIMDr"] += *insertedConsolidationDetail.AmountJelimDr
				rowCode[v.Code]["AmountJELIMCr"] += *insertedConsolidationDetail.AmountJelimCr
				rowCode[v.Code]["AmountConsole"] += *insertedConsolidationDetail.AmountConsole
			}

			if v.AutoSummary != nil && *v.AutoSummary {
				modelConsolDetail := model.ConsolidationDetailEntityModel{}
				modelConsolDetail.Context = ctx
				// rowCode[v.Code] = make(map[string]float64)
				tmp := rowCode[v.Code]["AmountBeforeJPM"]
				modelConsolDetail.AmountBeforeJpm = &tmp
				tmp1 := rowCode[v.Code]["AmountJPMCr"]
				modelConsolDetail.AmountJpmCr = &tmp1
				tmp2 := rowCode[v.Code]["AmountJPMDr"]
				modelConsolDetail.AmountJpmDr = &tmp2
				tmp3 := rowCode[v.Code]["AmountAfterJPM"]
				modelConsolDetail.AmountAfterJpm = &tmp3
				tmp4 := rowCode[v.Code]["AmountJCTECr"]
				modelConsolDetail.AmountJcteCr = &tmp4
				tmp5 := rowCode[v.Code]["AmountJCTEDr"]
				modelConsolDetail.AmountJcteDr = &tmp5
				tmp6 := rowCode[v.Code]["AmountAfterJCTE"]
				modelConsolDetail.AmountAfterJcte = &tmp6
				tmp7 := rowCode[v.Code]["AmountCombine"]
				modelConsolDetail.AmountCombineSubsidiary = &tmp7
				tmp8 := rowCode[v.Code]["AmountJELIMDr"]
				modelConsolDetail.AmountJelimDr = &tmp8
				tmp9 := rowCode[v.Code]["AmountJELIMCr"]
				modelConsolDetail.AmountJelimCr = &tmp9
				tmp10 := rowCode[v.Code]["AmountConsole"]
				modelConsolDetail.AmountConsole = &tmp10
				modelConsolDetail.Code = fmt.Sprintf("%s_SUBTOTAL", v.Code)
				modelConsolDetail.Description = "SUB TOTAL"
				modelConsolDetail.SortID = v.SortID
				modelConsolDetail.ConsolidationID = consolData.ID
				insertedConsolidationDetail, err := s.ConsolidationDetailRepository.Create(ctx, &modelConsolDetail)
				if err != nil {
					errMsg = "Terjadi kesalahan pada saat memperbarui perhitungan summary coa. Kriteria Data: " + v.Code + ". Error: " + err.Error()
					return
				}

				rowCode[fmt.Sprintf("%s_SUBTOTAL", v.Code)] = make(map[string]float64)
				rowCode[fmt.Sprintf("%s_SUBTOTAL", v.Code)]["AmountBeforeJPM"] += *insertedConsolidationDetail.AmountBeforeJpm
				rowCode[fmt.Sprintf("%s_SUBTOTAL", v.Code)]["AmountJPMDr"] += *insertedConsolidationDetail.AmountJpmDr
				rowCode[fmt.Sprintf("%s_SUBTOTAL", v.Code)]["AmountJPMCr"] += *insertedConsolidationDetail.AmountJpmCr
				rowCode[fmt.Sprintf("%s_SUBTOTAL", v.Code)]["AmountAfterJPM"] += *insertedConsolidationDetail.AmountAfterJpm
				rowCode[fmt.Sprintf("%s_SUBTOTAL", v.Code)]["AmountJCTEDr"] += *insertedConsolidationDetail.AmountJcteDr
				rowCode[fmt.Sprintf("%s_SUBTOTAL", v.Code)]["AmountJCTECr"] += *insertedConsolidationDetail.AmountJcteCr
				rowCode[fmt.Sprintf("%s_SUBTOTAL", v.Code)]["AmountAfterJCTE"] += *insertedConsolidationDetail.AmountAfterJcte
				rowCode[fmt.Sprintf("%s_SUBTOTAL", v.Code)]["AmountCombine"] += *insertedConsolidationDetail.AmountCombineSubsidiary
				rowCode[fmt.Sprintf("%s_SUBTOTAL", v.Code)]["AmountJELIMDr"] += *insertedConsolidationDetail.AmountJelimDr
				rowCode[fmt.Sprintf("%s_SUBTOTAL", v.Code)]["AmountJELIMCr"] += *insertedConsolidationDetail.AmountJelimCr
				rowCode[fmt.Sprintf("%s_SUBTOTAL", v.Code)]["AmountConsole"] += *insertedConsolidationDetail.AmountConsole
			}
		}

		if v.IsTotal != nil && *v.IsTotal {
			//masukan dulu ke db dengan nilai amountnya 0
			total := make(map[string]float64)
			for _, vSubtotal := range []string{"AmountBeforeJPM", "AmountJPMDr", "AmountJPMCr", "AmountAfterJPM", "AmountJCTEDr", "AmountJCTECr", "AmountAfterJCTE", "AmountCombine", "AmountJELIMDr", "AmountJELIMCr", "AmountConsole"} {
				total[vSubtotal] = 0.0
				if v.FxSummary != "" {
					formula := v.FxSummary
					reg := regexp.MustCompile(`[A-Za-z0-9_~:()]+|[0-9]+\d{3,}`)
					match := reg.FindAllString(formula, -1)
					for _, vMatch := range match {
						if _, ok := rowCode[fmt.Sprintf("%s_SUBTOTAL", vMatch)][vSubtotal]; ok {
							formula = strings.ReplaceAll(formula, vMatch, fmt.Sprintf("(%.2f)", rowCode[fmt.Sprintf("%s_SUBTOTAL", vMatch)][vSubtotal]))
						} else {
							if _, ok := rowCode[vMatch][vSubtotal]; ok {
								formula = strings.ReplaceAll(formula, vMatch, fmt.Sprintf("(%.2f)", rowCode[vMatch][vSubtotal]))
							}
						}
					}

					expression, err := govaluate.NewEvaluableExpression(formula)
					if err != nil {
						errMsg = "Terjadi kesalahan pada saat memperbarui perhitungan total. Kriteria Data: " + v.Code + ". Error: " + err.Error()
						return
					}
					result, err := expression.Evaluate(nil)
					if err != nil {
						errMsg = "Terjadi kesalahan pada saat memperbarui perhitungan total. Kriteria Data: " + v.Code + ". Error: " + err.Error()
						return
					}
					f, ok := result.(float64) // type assertion
					if !ok {
						fmt.Println("Interface value is not a float64")
						return
					}
					total[vSubtotal] = float64(int(f*100)) / 100
				}

			}

			modelConsolDetail := model.ConsolidationDetailEntityModel{}
			modelConsolDetail.Context = ctx
			modelConsolDetail.Code = v.Code
			modelConsolDetail.Description = v.Description
			modelConsolDetail.SortID = v.SortID
			modelConsolDetail.ConsolidationID = consolData.ID
			tmpamount := total["AmountBeforeJPM"]
			modelConsolDetail.AmountBeforeJpm = &tmpamount
			tmpamount1 := total["AmountJPMDr"]
			modelConsolDetail.AmountJpmDr = &tmpamount1
			tmpamount2 := total["AmountJPMCr"]
			modelConsolDetail.AmountJpmCr = &tmpamount2
			tmpamount3 := total["AmountAfterJPM"]
			modelConsolDetail.AmountAfterJpm = &tmpamount3
			tmpamount4 := total["AmountJCTEDr"]
			modelConsolDetail.AmountJcteDr = &tmpamount4
			tmpamount5 := total["AmountJCTECr"]
			modelConsolDetail.AmountJcteCr = &tmpamount5
			tmpamount6 := total["AmountAfterJCTE"]
			modelConsolDetail.AmountAfterJcte = &tmpamount6
			tmpamount7 := total["AmountCombine"]
			modelConsolDetail.AmountCombineSubsidiary = &tmpamount7
			tmpamount8 := total["AmountJELIMDr"]
			modelConsolDetail.AmountJelimDr = &tmpamount8
			tmpamount9 := total["AmountJELIMCr"]
			modelConsolDetail.AmountJelimCr = &tmpamount9
			tmpamount10 := total["AmountConsole"]
			modelConsolDetail.AmountConsole = &tmpamount10
			_, err = s.ConsolidationDetailRepository.Create(ctx, &modelConsolDetail)
			if err != nil {
				errMsg = "Terjadi kesalahan pada saat menginput perhitungan total. Kriteria Data: " + v.Code + ". Error: " + err.Error()
				return
			}
			rowCode[v.Code] = make(map[string]float64)
			rowCode[v.Code]["AmountBeforeJPM"] = total["AmountBeforeJPM"]
			rowCode[v.Code]["AmountJPMCr"] = total["AmountJPMCr"]
			rowCode[v.Code]["AmountJPMDr"] = total["AmountJPMDr"]
			rowCode[v.Code]["AmountAfterJPM"] = total["AmountAfterJPM"]
			rowCode[v.Code]["AmountJCTEDr"] = total["AmountJCTEDr"]
			rowCode[v.Code]["AmountJCTECr"] = total["AmountJCTECr"]
			rowCode[v.Code]["AmountAfterJCTE"] = total["AmountAfterJCTE"]
			rowCode[v.Code]["AmountCombine"] = total["AmountCombine"]
			rowCode[v.Code]["AmountJELIMDr"] = total["AmountJELIMDr"]
			rowCode[v.Code]["AmountJELIMCr"] = total["AmountJELIMCr"]
			rowCode[v.Code]["AmountConsole"] = total["AmountConsole"]
		}
		if v.IsLabel != nil && *v.IsLabel && v.IsTotal != nil && !*v.IsTotal && v.IsCoa != nil && !*v.IsCoa {
			modelConsolDetail := model.ConsolidationDetailEntityModel{}
			modelConsolDetail.Context = ctx
			modelConsolDetail.Code = v.Code
			modelConsolDetail.Description = v.Description
			modelConsolDetail.SortID = v.SortID
			modelConsolDetail.ConsolidationID = consolData.ID
			_, err := s.ConsolidationDetailRepository.Create(ctx, &modelConsolDetail)
			if err != nil {
				errMsg = "Terjadi kesalahan pada saat memperbarui perhitungan summary coa. Kriteria Data: " + v.Code + ". Error: " + err.Error()
				return
			}
		}
	}
}

func (s *service) ConsolidateBridge(ctx *abstraction.Context, consolData *model.ConsolidationEntityModel, consollBridgeData *model.ConsolidationBridgeEntityModel, wg *sync.WaitGroup) {
	// defer wg.Done()
	if consollBridgeData.ConsolidationVersions != 0 {
		criteriaConsol := model.ConsolidationFilterModel{}
		criteriaConsol.CompanyID = &consollBridgeData.CompanyID
		criteriaConsol.Period = &consollBridgeData.Period
		criteriaConsol.Versions = &consollBridgeData.Versions
		criteriaConsol.ConsolidationVersions = &consollBridgeData.ConsolidationVersions

		consolidation, err := s.ConsolidationRepository.FindByCriteria(ctx, &criteriaConsol)
		if err != nil {
			errMsg = "Data penghubung konsolidasi dengan anak usaha tidak ditemukan"
			return
		}
		criteriaConsolDetail := model.ConsolidationDetailFilterModel{}
		criteriaConsolDetail.ConsolidationID = &consolidation.ID
		consolDetail, err := s.ConsolidationDetailRepository.Find(ctx, &criteriaConsolDetail)
		if err != nil {
			errMsg = "Data detil konsolidasi anak usaha tidak ditemukan"
			return
		}
		for _, vConsol := range *consolDetail {
			// if vConsol.AmountConsole == nil || (vConsol.AmountConsole != nil && *vConsol.AmountCombineSubsidiary == 0) {
			// 	continue
			// }
			modelConsolBridgeDetail := model.ConsolidationBridgeDetailEntityModel{}
			modelConsolBridgeDetail.Code = vConsol.Code
			modelConsolBridgeDetail.ConsolidationBridgeID = consollBridgeData.ID
			modelConsolBridgeDetail.Amount = vConsol.AmountConsole
			_, err := s.ConsolidationBridgeDetailRepository.Create(ctx, &modelConsolBridgeDetail)
			if err != nil {
				errMsg = "Terjadi kesalahan pada saat memperbarui data detil konsolidasi anak usaha. Error: " + err.Error()
				return
			}
		}
	} else {
		criteriaTB := model.TrialBalanceFilterModel{}
		criteriaTB.CompanyID = &consollBridgeData.CompanyID
		criteriaTB.Versions = &consollBridgeData.Versions
		criteriaTB.Period = &consolData.Period
		trialBalance, err := s.TrialBalanceRepository.FindByCriteria(ctx, &criteriaTB)
		if err != nil {
			errMsg = "Data trial balance tidak ditemukan"
			return
		}

		criteriaFmtBridge := model.FormatterBridgesFilterModel{}
		strTB := "TRIAL-BALANCE"
		criteriaFmtBridge.Source = &strTB
		criteriaFmtBridge.TrxRefID = &trialBalance.ID
		fmtBridge, err := s.FormatterBridgesRepository.FindWithCriteria(ctx, &criteriaFmtBridge)
		if err != nil {
			errMsg = "Data formatter bridge tidak ditemukan"
			return
		}

		criteriaTBD := model.TrialBalanceDetailFilterModel{}
		criteriaTBD.FormatterBridgesID = &fmtBridge.ID
		tbdetails, err := s.TrialBalanceDetailRepository.Find(ctx, &criteriaTBD)
		if err != nil {
			errMsg = "Data trial balance detail tidak ditemukan"
			return
		}
		for _, tbd := range *tbdetails {
			if tbd.AmountAfterAje == nil || (tbd.AmountAfterAje != nil && *tbd.AmountAfterAje == 0) {
				continue
			}
			modelConsolBridgeDetail := model.ConsolidationBridgeDetailEntityModel{}
			modelConsolBridgeDetail.Code = tbd.Code
			modelConsolBridgeDetail.ConsolidationBridgeID = consollBridgeData.ID
			modelConsolBridgeDetail.Amount = tbd.AmountAfterAje
			_, err := s.ConsolidationBridgeDetailRepository.Create(ctx, &modelConsolBridgeDetail)
			if err != nil {
				errMsg = "Terjadi kesalahan pada saat memperbarui data detil konsolidasi anak usaha. Error: " + err.Error()
				return
			}
		}
	}
}

func sumByHeadCoa(headCoa string, amountBefore, amountDr, amountCr float64) float64 {
	tmp := 0.0
	switch headCoa {
	case "1", "5", "6", "7", "91", "92":
		tmp = amountBefore + amountDr - amountCr
	default:
		tmp = amountBefore - amountDr + amountCr
	}

	return tmp
}
