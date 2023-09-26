package reupload

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"math"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"
	"worker/internal/abstraction"
	"worker/internal/dto"
	"worker/internal/factory"
	kafkaproducer "worker/internal/kafka/producer"
	"worker/internal/model"
	"worker/internal/repository"
	"worker/pkg/util/trxmanager"

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
	ImportedWorksheetRepository          repository.ImportedWorksheet
	ImportedWorksheetDetailRepository    repository.ImportedWorksheetDetail
	FormatterBridgesRepository           repository.FormatterBridges
	EmployeeBenefitRepository            repository.EmployeeBenefit
	EmployeeBenefitDetailRepository      repository.EmployeeBenefitDetail
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
	importedWorksheetRepository := f.ImportedWorksheetRepository
	importedWorksheetDetailRepository := f.ImportedWorksheetDetailRepository
	formatterBridgesRepository := f.FormatterBridgesRepository
	employeeBenefitRepository := f.EmployeeBenefitRepository
	employeeBenefitDetailRepository := f.EmployeeBenefitDetailRepository
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
		ImportedWorksheetRepository:          importedWorksheetRepository,
		ImportedWorksheetDetailRepository:    importedWorksheetDetailRepository,
		FormatterBridgesRepository:           formatterBridgesRepository,
		EmployeeBenefitRepository:            employeeBenefitRepository,
		EmployeeBenefitDetailRepository:      employeeBenefitDetailRepository,
	}
}

var errs []error

type MsgTrialBalance struct {
	Errmsg  string `json:"msg"`
	LineMsg string `json:"line_msg"`
}

var lineMsg string
var msgErrTrialBalance string
var errTrialBalance string

func (s *service) ReUpload(ctx *abstraction.Context, payload *abstraction.JsonDataReUpload) string {

	if payload.TrialBalance != "" {
		var dataTb []model.TrialBalanceDetailEntity
		if _, err := s.ImportTrialBalance(ctx, payload, dataTb); err != nil {
			if trxmanager.New(s.Db).WithTrx(ctx, func(ctx *abstraction.Context) error {

				msgTrialBalance := MsgTrialBalance{msgErrTrialBalance, lineMsg}
				jsonErr, err := json.Marshal(msgTrialBalance)
				if err != nil {
					fmt.Println(err)
					return err
				}
				var data model.ImportedWorksheetDetailEntityModel
				_, err = s.ImportedWorksheetDetailRepository.FindByID(ctx, &payload.IDWorksheetDetailTrialBalance)
				if err != nil {
					return err
				}
				data.Context = ctx
				data.ImportedWorksheetDetailEntity = model.ImportedWorksheetDetailEntity{
					Status:      1,
					Code:        "TRIAL-BALANCE",
					Name:        "Trial Balance",
					FileName:    payload.FNTrialBalance,
					Note:        payload.TrialBalance,
					ErrMessages: string(jsonErr),
				}

				_, err = s.ImportedWorksheetDetailRepository.Update(ctx, &payload.IDWorksheetDetailTrialBalance, &data)
				if err != nil {
					return err
				}
				return nil
			}); err != nil {
			}
		}
	}
	if payload.MutasiFA != "" {
		var dataDta []model.MutasiFaDetailEntity
		if _, err := s.ImportMutasiFA(ctx, payload, dataDta); err != nil {
			if trxmanager.New(s.Db).WithTrx(ctx, func(ctx *abstraction.Context) error {
				var data model.ImportedWorksheetDetailEntityModel
				_, err = s.ImportedWorksheetDetailRepository.FindByID(ctx, &payload.IDWorksheetDetailMutasiFA)
				if err != nil {
					return err
				}
				data.Context = ctx
				data.ImportedWorksheetDetailEntity = model.ImportedWorksheetDetailEntity{
					// Code: payload.MutasiDta,
					Status: 1,
					Note:   payload.MutasiFA,
				}

				_, err = s.ImportedWorksheetDetailRepository.Update(ctx, &payload.IDWorksheetDetailMutasiFA, &data)
				if err != nil {
					return err
				}
				return nil
			}); err != nil {
			}
		}
	}
	if payload.MutasiDta != "" {
		var dataDta []model.MutasiDtaDetailEntity
		if _, err := s.ImportMutasiDta(ctx, payload, dataDta); err != nil {
			if trxmanager.New(s.Db).WithTrx(ctx, func(ctx *abstraction.Context) error {
				var data model.ImportedWorksheetDetailEntityModel
				_, err = s.ImportedWorksheetDetailRepository.FindByID(ctx, &payload.IDWorksheetDetailMutasiDta)
				if err != nil {
					return err
				}
				data.Context = ctx
				data.ImportedWorksheetDetailEntity = model.ImportedWorksheetDetailEntity{
					// Code: payload.MutasiDta,
					Status: 1,
					Note:   payload.MutasiDta,
				}

				_, err = s.ImportedWorksheetDetailRepository.Update(ctx, &payload.IDWorksheetDetailMutasiDta, &data)
				if err != nil {
					return err
				}
				return nil
			}); err != nil {
			}
		}
	}
	if payload.InvestasiTbk != "" {
		var dataIt []model.InvestasiTbkDetailEntity
		if _, err := s.ImportInvestasiTbk(ctx, payload, dataIt); err != nil {
			if trxmanager.New(s.Db).WithTrx(ctx, func(ctx *abstraction.Context) error {
				var data model.ImportedWorksheetDetailEntityModel
				_, err = s.ImportedWorksheetDetailRepository.FindByID(ctx, &payload.IDWorksheetDetailInvestasiTbk)
				if err != nil {
					return err
				}
				data.Context = ctx
				data.ImportedWorksheetDetailEntity = model.ImportedWorksheetDetailEntity{
					// Code: payload.InvestasiTbk,
					Status: 1,
					Note:   payload.InvestasiTbk,
				}

				_, err = s.ImportedWorksheetDetailRepository.Update(ctx, &payload.IDWorksheetDetailInvestasiTbk, &data)
				if err != nil {
					return err
				}
				return nil
			}); err != nil {
			}
		}
	}
	if payload.AgingUtangPiutang != "" {
		var dataAup []model.AgingUtangPiutangDetailEntity
		if _, err := s.ImportAgingUtangPiutang(ctx, payload, dataAup); err != nil {
			if trxmanager.New(s.Db).WithTrx(ctx, func(ctx *abstraction.Context) error {
				var data model.ImportedWorksheetDetailEntityModel
				_, err = s.ImportedWorksheetDetailRepository.FindByID(ctx, &payload.IDWorksheetDetailAgingUtangPiutang)
				if err != nil {
					return err
				}
				data.Context = ctx
				data.ImportedWorksheetDetailEntity = model.ImportedWorksheetDetailEntity{
					// Code: payload.InvestasiTbk,
					Status: 1,
					Note:   payload.AgingUtangPiutang,
				}

				_, err = s.ImportedWorksheetDetailRepository.Update(ctx, &payload.IDWorksheetDetailAgingUtangPiutang, &data)
				if err != nil {
					return err
				}
				return nil
			}); err != nil {
			}
		}
	}
	if payload.InvestasiNonTbk != "" {
		var dataInt []model.InvestasiNonTbkDetailEntity
		if _, err := s.ImportInvestasiNonTbk(ctx, payload, dataInt); err != nil {
			if trxmanager.New(s.Db).WithTrx(ctx, func(ctx *abstraction.Context) error {
				var data model.ImportedWorksheetDetailEntityModel
				_, err = s.ImportedWorksheetDetailRepository.FindByID(ctx, &payload.IDWorksheetDetailInvestasiNonTbk)
				if err != nil {
					return err
				}
				data.Context = ctx
				data.ImportedWorksheetDetailEntity = model.ImportedWorksheetDetailEntity{
					// Code: payload.InvestasiNonTbk,
					Status: 1,
					Note:   payload.InvestasiNonTbk,
				}

				_, err = s.ImportedWorksheetDetailRepository.Update(ctx, &payload.IDWorksheetDetailInvestasiNonTbk, &data)
				if err != nil {
					return err
				}
				return nil
			}); err != nil {
			}
		}
	}
	if payload.MutasiRua != "" {
		var dataRua []model.MutasiRuaDetailEntity
		if _, err := s.ImportMutasiRua(ctx, payload, dataRua); err != nil {
			if trxmanager.New(s.Db).WithTrx(ctx, func(ctx *abstraction.Context) error {
				var data model.ImportedWorksheetDetailEntityModel
				_, err = s.ImportedWorksheetDetailRepository.FindByID(ctx, &payload.IDWorksheetDetailMutasiRua)
				if err != nil {
					return err
				}
				data.Context = ctx
				data.ImportedWorksheetDetailEntity = model.ImportedWorksheetDetailEntity{
					// Code: payload.MutasiRua,
					Status: 1,
					Note:   payload.MutasiRua,
				}

				_, err = s.ImportedWorksheetDetailRepository.Update(ctx, &payload.IDWorksheetDetailMutasiRua, &data)
				if err != nil {
					return err
				}
				return nil
			}); err != nil {
			}
		}
	}
	if payload.MutasiIa != "" {
		var dataIa []model.MutasiIaDetailEntity
		if _, err := s.ImportMutasiIa(ctx, payload, dataIa); err != nil {
			if trxmanager.New(s.Db).WithTrx(ctx, func(ctx *abstraction.Context) error {
				var data model.ImportedWorksheetDetailEntityModel
				_, err = s.ImportedWorksheetDetailRepository.FindByID(ctx, &payload.IDWorksheetDetailMutasiIa)
				if err != nil {
					return err
				}
				data.Context = ctx
				data.ImportedWorksheetDetailEntity = model.ImportedWorksheetDetailEntity{
					// Code: payload.MutasiIa,
					Status: 1,
					Note:   payload.MutasiIa,
				}

				_, err = s.ImportedWorksheetDetailRepository.Update(ctx, &payload.IDWorksheetDetailMutasiIa, &data)
				if err != nil {
					return err
				}
				return nil
			}); err != nil {
			}
		}
	}
	if payload.MutasiPersediaan != "" {
		var dataMp []model.MutasiPersediaanDetailEntity
		if _, err := s.ImportMutasiPersediaan(ctx, payload, dataMp); err != nil {
			if trxmanager.New(s.Db).WithTrx(ctx, func(ctx *abstraction.Context) error {
				var data model.ImportedWorksheetDetailEntityModel
				_, err = s.ImportedWorksheetDetailRepository.FindByID(ctx, &payload.IDWorksheetDetailMutasiPersediaan)
				if err != nil {
					return err
				}
				data.Context = ctx
				data.ImportedWorksheetDetailEntity = model.ImportedWorksheetDetailEntity{
					// Code: payload.MutasiPersediaan,
					Status: 1,
					Note:   payload.MutasiPersediaan,
				}

				_, err = s.ImportedWorksheetDetailRepository.Update(ctx, &payload.IDWorksheetDetailMutasiPersediaan, &data)
				if err != nil {
					return err
				}
				return nil
			}); err != nil {
			}
		}
	}
	if payload.PembelianPenjualanBerelasi != "" {
		var dataPPB []model.PembelianPenjualanBerelasiDetailEntity
		if _, err := s.ImportPembelianPenjualanBerelasi(ctx, payload, dataPPB); err != nil {
			if trxmanager.New(s.Db).WithTrx(ctx, func(ctx *abstraction.Context) error {
				var data model.ImportedWorksheetDetailEntityModel
				_, err = s.ImportedWorksheetDetailRepository.FindByID(ctx, &payload.IDWorksheetDetailPembelianPenjualanBerelasi)
				if err != nil {
					return err
				}
				data.Context = ctx
				data.ImportedWorksheetDetailEntity = model.ImportedWorksheetDetailEntity{
					// Code: payload.PembelianPenjualanBerelasi,
					Status: 1,
					Note:   payload.PembelianPenjualanBerelasi,
				}

				_, err = s.ImportedWorksheetDetailRepository.Update(ctx, &payload.IDWorksheetDetailPembelianPenjualanBerelasi, &data)
				if err != nil {
					return err
				}
				return nil
			}); err != nil {
			}
		}
	}
	if payload.EmployeeBenefit != "" {
		var dataEB []model.EmployeeBenefitDetailEntity
		if _, err := s.ImportEmployeeBenefit(ctx, payload, dataEB); err != nil {
			if trxmanager.New(s.Db).WithTrx(ctx, func(ctx *abstraction.Context) error {
				var data model.ImportedWorksheetDetailEntityModel
				_, err = s.ImportedWorksheetDetailRepository.FindByID(ctx, &payload.IDWorksheetDetailEmployeeBenefit)
				if err != nil {
					return err
				}
				data.Context = ctx
				data.ImportedWorksheetDetailEntity = model.ImportedWorksheetDetailEntity{
					// Code: payload.PembelianPenjualanBerelasi,
					Status: 1,
					Note:   payload.EmployeeBenefit,
				}

				_, err = s.ImportedWorksheetDetailRepository.Update(ctx, &payload.IDWorksheetDetailEmployeeBenefit, &data)
				if err != nil {
					return err
				}
				return nil
			}); err != nil {
			}
		}
	}
	waktu := time.Now()
	if err := trxmanager.New(s.Db).WithTrx(ctx, func(ctx *abstraction.Context) error {

		var dataImportedWorksheetDetail model.ImportedWorksheetDetailEntityModel
		dataImportedWorksheetDetail.Context = ctx
		dataImportedWorksheetDetail.ImportedWorksheetDetailEntity = model.ImportedWorksheetDetailEntity{
			ImportedWorksheetID: payload.ImportedWorkSheetID,
			Status:              2,
		}
		getStatusImportedWorkSheetDetailSucces, err := s.ImportedWorksheetDetailRepository.GetCountStatus(ctx, &dataImportedWorksheetDetail)
		if err != nil {
			fmt.Println(err)
			return err
		}
		dataImportedWorksheetDetail.Context = ctx
		dataImportedWorksheetDetail.ImportedWorksheetDetailEntity = model.ImportedWorksheetDetailEntity{
			ImportedWorksheetID: payload.ImportedWorkSheetID,
			Status:              1,
		}
		getStatusImportedWorkSheetDetailFailed, err := s.ImportedWorksheetDetailRepository.GetCountStatus(ctx, &dataImportedWorksheetDetail)
		if err != nil {
			fmt.Println(err)
			return err
		}
		// &Uncompleated := 0
		Compleated := 1
		if len(*getStatusImportedWorkSheetDetailSucces) == 11 {
			var data model.ImportedWorksheetEntityModel
			worksheet, err := s.ImportedWorksheetRepository.FindByID(ctx, &payload.ImportedWorkSheetID)
			if err != nil {
				fmt.Println(err)
				return err
			}
			data.Context = ctx
			data.ImportedWorksheetEntity = model.ImportedWorksheetEntity{
				Status: 2,
			}

			_, err = s.ImportedWorksheetRepository.Update(ctx, &payload.ImportedWorkSheetID, &data)
			if err != nil {
				fmt.Println(err)
				return err
			}

			datePeriod, err := time.Parse(time.RFC3339, worksheet.Period)
			if err != nil {
				return err
			}
			period := datePeriod.Format("2006-01-02")
			//Update Trial Balance
			var dataTrialBalance model.TrialBalanceEntityModel
			dataTb, err := s.TrialBalanceRepository.FindByID(ctx, &payload.Version, &payload.CompanyID, &period)
			if err != nil {
				fmt.Println(err)
				return err
			}
			dataTrialBalance.Context = ctx

			dataTrialBalance.TrialBalanceEntity = model.TrialBalanceEntity{
				Status: &Compleated,
			}

			_, err = s.TrialBalanceRepository.Update(ctx, &dataTb.ID, &dataTrialBalance)
			if err != nil {
				fmt.Println(err)
				return err
			}

			//Update Aging Utang Piutang
			var dataAgingUtangPiutang model.AgingUtangPiutangEntityModel
			dataAup, err := s.AgingUtangPiutangRepository.FindByID(ctx, &payload.Version, &payload.CompanyID, &period)
			if err != nil {
				fmt.Println(err)
				return err
			}
			dataAgingUtangPiutang.Context = ctx
			dataAgingUtangPiutang.AgingUtangPiutangEntity = model.AgingUtangPiutangEntity{
				Status: &Compleated,
			}
			_, err = s.AgingUtangPiutangRepository.Update(ctx, &dataAup.ID, &dataAgingUtangPiutang)
			if err != nil {
				fmt.Println(err)
				return err
			}

			//Update Mutasi DTA
			var dataMutasiDta model.MutasiDtaEntityModel
			dataMDta, err := s.MutasiDtaRepository.FindByID(ctx, &payload.Version, &payload.CompanyID, &period)
			if err != nil {
				fmt.Println(err)
				return err
			}
			dataMutasiDta.Context = ctx
			dataMutasiDta.MutasiDtaEntity = model.MutasiDtaEntity{
				Status: &Compleated,
			}
			_, err = s.MutasiDtaRepository.Update(ctx, &dataMDta.ID, &dataMutasiDta)
			if err != nil {
				fmt.Println(err)
				return err
			}

			//Update Mutasi FA
			var dataMutasiFa model.MutasiFaEntityModel
			dataMFa, err := s.MutasiFaRepository.FindByID(ctx, &payload.Version, &payload.CompanyID, &period)
			if err != nil {
				fmt.Println(err)
				return err
			}
			dataMutasiFa.Context = ctx
			dataMutasiFa.MutasiFaEntity = model.MutasiFaEntity{
				Status: &Compleated,
			}
			_, err = s.MutasiFaRepository.Update(ctx, &dataMFa.ID, &dataMutasiFa)
			if err != nil {
				fmt.Println(err)
				return err
			}

			//Update Mutasi IA
			var dataMutasiIa model.MutasiIaEntityModel
			dataMIa, err := s.MutasiIaRepository.FindByID(ctx, &payload.Version, &payload.CompanyID, &period)
			if err != nil {
				fmt.Println(err)
				return err
			}
			dataMutasiIa.Context = ctx
			dataMutasiIa.MutasiIaEntity = model.MutasiIaEntity{
				Status: &Compleated,
			}
			_, err = s.MutasiIaRepository.Update(ctx, &dataMIa.ID, &dataMutasiIa)
			if err != nil {
				fmt.Println(err)
				return err
			}

			//Update Mutasi Persediaan
			var dataMutasiP model.MutasiPersediaanEntityModel
			dataMP, err := s.MutasiPersediaanRepository.FindByID(ctx, &payload.Version, &payload.CompanyID, &period)
			if err != nil {
				fmt.Println(err)
				return err
			}
			dataMutasiP.Context = ctx
			dataMutasiP.MutasiPersediaanEntity = model.MutasiPersediaanEntity{
				Status: &Compleated,
			}
			_, err = s.MutasiPersediaanRepository.Update(ctx, &dataMP.ID, &dataMutasiP)
			if err != nil {
				fmt.Println(err)
				return err
			}

			//Update Mutasi RUA
			var dataMutasiR model.MutasiRuaEntityModel
			dataMR, err := s.MutasiRuaRepository.FindByID(ctx, &payload.Version, &payload.CompanyID, &period)
			if err != nil {
				fmt.Println(err)
				return err
			}
			dataMutasiR.Context = ctx
			dataMutasiR.MutasiRuaEntity = model.MutasiRuaEntity{
				Status: &Compleated,
			}
			_, err = s.MutasiRuaRepository.Update(ctx, &dataMR.ID, &dataMutasiR)
			if err != nil {
				fmt.Println(err)
				return err
			}

			//Update Pembelian & Penjualan Berelasi
			var dataPberelasi model.PembelianPenjualanBerelasiEntityModel
			dataPb, err := s.PembelianPenjualanBerelasiRepository.FindByID(ctx, &payload.Version, &payload.CompanyID, &period)
			if err != nil {
				fmt.Println(err)
				return err
			}
			dataPberelasi.Context = ctx
			dataPberelasi.PembelianPenjualanBerelasiEntity = model.PembelianPenjualanBerelasiEntity{
				Status: &Compleated,
			}
			_, err = s.PembelianPenjualanBerelasiRepository.Update(ctx, &dataPb.ID, &dataPberelasi)
			if err != nil {
				fmt.Println(err)
				return err
			}

			//Update Investasi Non Tbk
			var dataIntbk model.InvestasiNonTbkEntityModel
			dataInt, err := s.InvestasiNonTbkRepository.FindByID(ctx, &payload.Version, &payload.CompanyID, &period)
			if err != nil {
				fmt.Println(err)
				return err
			}
			dataIntbk.Context = ctx
			dataIntbk.InvestasiNonTbkEntity = model.InvestasiNonTbkEntity{
				Status: &Compleated,
			}
			_, err = s.InvestasiNonTbkRepository.Update(ctx, &dataInt.ID, &dataIntbk)
			if err != nil {
				fmt.Println(err)
				return err
			}

			//Update Investasi Tbk
			var dataItbk model.InvestasiTbkEntityModel
			dataIt, err := s.InvestasiTbkRepository.FindByID(ctx, &payload.Version, &payload.CompanyID, &period)
			if err != nil {
				fmt.Println(err)
				return err
			}
			dataItbk.Context = ctx
			dataItbk.InvestasiTbkEntity = model.InvestasiTbkEntity{
				Status: &Compleated,
			}
			_, err = s.InvestasiTbkRepository.Update(ctx, &dataIt.ID, &dataItbk)
			if err != nil {
				fmt.Println(err)
				return err
			}

			//Update Employee Benefit
			var dataEb model.EmployeeBenefitEntityModel
			dataEbt, err := s.EmployeeBenefitRepository.FindByID(ctx, &payload.Version, &payload.CompanyID, &period)
			if err != nil {
				fmt.Println(err)
				return err
			}
			dataEb.Context = ctx
			dataEb.EmployeeBenefitEntity = model.EmployeeBenefitEntity{
				Status: &Compleated,
			}
			_, err = s.EmployeeBenefitRepository.Update(ctx, &dataEbt.ID, &dataEb)
			if err != nil {
				fmt.Println(err)
				return err
			}
		}
		if len(*getStatusImportedWorkSheetDetailSucces) != 11 {

			var data model.ImportedWorksheetEntityModel
			worksheet, err := s.ImportedWorksheetRepository.FindByID(ctx, &payload.ImportedWorkSheetID)
			if err != nil {
				fmt.Println(err)
				return err
			}
			data.Context = ctx
			data.ImportedWorksheetEntity = model.ImportedWorksheetEntity{
				Status: 1,
			}

			_, err = s.ImportedWorksheetRepository.Update(ctx, &payload.ImportedWorkSheetID, &data)
			if err != nil {
				fmt.Println(err)
				return err
			}
			datePeriod, err := time.Parse(time.RFC3339, worksheet.Period)
			if err != nil {
				return err
			}
			period := datePeriod.Format("2006-01-02")
			Uncompleated := 0
			//Update Trial Balance
			var dataTrialBalance model.TrialBalanceEntityModel
			dataTb, err := s.TrialBalanceRepository.FindByID(ctx, &payload.Version, &payload.CompanyID, &period)
			if err != nil {
				fmt.Println(err)
				return err
			}
			dataTrialBalance.Context = ctx
			dataTrialBalance.TrialBalanceEntity = model.TrialBalanceEntity{
				Status: &Uncompleated,
			}

			_, err = s.TrialBalanceRepository.Update(ctx, &dataTb.ID, &dataTrialBalance)
			if err != nil {
				fmt.Println(err)
				return err
			}

			//Update Aging Utang Piutang
			var dataAgingUtangPiutang model.AgingUtangPiutangEntityModel
			dataAup, err := s.AgingUtangPiutangRepository.FindByID(ctx, &payload.Version, &payload.CompanyID, &period)
			if err != nil {
				fmt.Println(err)
				return err
			}

			dataAgingUtangPiutang.Context = ctx
			dataAgingUtangPiutang.AgingUtangPiutangEntity = model.AgingUtangPiutangEntity{
				Status: &Uncompleated,
			}
			_, err = s.AgingUtangPiutangRepository.Update(ctx, &dataAup.ID, &dataAgingUtangPiutang)
			if err != nil {
				fmt.Println(err)
				return err
			}

			//Update Mutasi DTA
			var dataMutasiDta model.MutasiDtaEntityModel
			dataMDta, err := s.MutasiDtaRepository.FindByID(ctx, &payload.Version, &payload.CompanyID, &period)
			if err != nil {
				fmt.Println(err)
				return err
			}
			dataMutasiDta.Context = ctx
			dataMutasiDta.MutasiDtaEntity = model.MutasiDtaEntity{
				Status: &Uncompleated,
			}
			_, err = s.MutasiDtaRepository.Update(ctx, &dataMDta.ID, &dataMutasiDta)
			if err != nil {
				fmt.Println(err)
				return err
			}

			//Update Mutasi FA
			var dataMutasiFa model.MutasiFaEntityModel
			dataMFa, err := s.MutasiFaRepository.FindByID(ctx, &payload.Version, &payload.CompanyID, &period)
			if err != nil {
				fmt.Println(err)
				return err
			}
			dataMutasiFa.Context = ctx
			dataMutasiFa.MutasiFaEntity = model.MutasiFaEntity{
				Status: &Uncompleated,
			}
			_, err = s.MutasiFaRepository.Update(ctx, &dataMFa.ID, &dataMutasiFa)
			if err != nil {
				fmt.Println(err)
				return err
			}

			//Update Mutasi IA
			var dataMutasiIa model.MutasiIaEntityModel
			dataMIa, err := s.MutasiIaRepository.FindByID(ctx, &payload.Version, &payload.CompanyID, &period)
			if err != nil {
				fmt.Println(err)
				return err
			}
			dataMutasiIa.Context = ctx
			dataMutasiIa.MutasiIaEntity = model.MutasiIaEntity{
				Status: &Uncompleated,
			}
			_, err = s.MutasiIaRepository.Update(ctx, &dataMIa.ID, &dataMutasiIa)
			if err != nil {
				fmt.Println(err)
				return err
			}

			//Update Mutasi Persediaan
			var dataMutasiP model.MutasiPersediaanEntityModel
			dataMP, err := s.MutasiPersediaanRepository.FindByID(ctx, &payload.Version, &payload.CompanyID, &period)
			if err != nil {
				fmt.Println(err)
				return err
			}
			dataMutasiP.Context = ctx
			dataMutasiP.MutasiPersediaanEntity = model.MutasiPersediaanEntity{
				Status: &Uncompleated,
			}
			_, err = s.MutasiPersediaanRepository.Update(ctx, &dataMP.ID, &dataMutasiP)
			if err != nil {
				fmt.Println(err)
				return err
			}

			//Update Mutasi RUA
			var dataMutasiR model.MutasiRuaEntityModel
			dataMR, err := s.MutasiRuaRepository.FindByID(ctx, &payload.Version, &payload.CompanyID, &period)
			if err != nil {
				fmt.Println(err)
				return err
			}
			dataMutasiR.Context = ctx
			dataMutasiR.MutasiRuaEntity = model.MutasiRuaEntity{
				Status: &Uncompleated,
			}
			_, err = s.MutasiRuaRepository.Update(ctx, &dataMR.ID, &dataMutasiR)
			if err != nil {
				fmt.Println(err)
				return err
			}

			//Update Pembelian & Penjualan Berelasi
			var dataPberelasi model.PembelianPenjualanBerelasiEntityModel
			dataPb, err := s.PembelianPenjualanBerelasiRepository.FindByID(ctx, &payload.Version, &payload.CompanyID, &period)
			if err != nil {
				fmt.Println(err)
				return err
			}
			dataPberelasi.Context = ctx

			dataPberelasi.PembelianPenjualanBerelasiEntity = model.PembelianPenjualanBerelasiEntity{
				Status: &Uncompleated,
			}
			_, err = s.PembelianPenjualanBerelasiRepository.Update(ctx, &dataPb.ID, &dataPberelasi)
			if err != nil {
				fmt.Println(err)
				return err
			}

			//Update Investasi Non Tbk
			var dataIntbk model.InvestasiNonTbkEntityModel
			dataInt, err := s.InvestasiNonTbkRepository.FindByID(ctx, &payload.Version, &payload.CompanyID, &period)
			if err != nil {
				fmt.Println(err)
				return err
			}
			dataIntbk.Context = ctx

			dataIntbk.InvestasiNonTbkEntity = model.InvestasiNonTbkEntity{
				Status: &Uncompleated,
			}
			_, err = s.InvestasiNonTbkRepository.Update(ctx, &dataInt.ID, &dataIntbk)
			if err != nil {
				fmt.Println(err)
				return err
			}

			//Update Investasi Tbk
			var dataItbk model.InvestasiTbkEntityModel
			dataIt, err := s.InvestasiTbkRepository.FindByID(ctx, &payload.Version, &payload.CompanyID, &period)
			if err != nil {
				fmt.Println(err)
				return err
			}
			dataItbk.Context = ctx
			dataItbk.InvestasiTbkEntity = model.InvestasiTbkEntity{
				Status: &Uncompleated,
			}
			_, err = s.InvestasiTbkRepository.Update(ctx, &dataIt.ID, &dataItbk)
			if err != nil {
				fmt.Println(err)
				return err
			}

			//Update Employee Benefit
			var dataEb model.EmployeeBenefitEntityModel
			dataEbt, err := s.EmployeeBenefitRepository.FindByID(ctx, &payload.Version, &payload.CompanyID, &period)
			if err != nil {
				fmt.Println(err)
				return err
			}
			dataEb.Context = ctx
			dataEb.EmployeeBenefitEntity = model.EmployeeBenefitEntity{
				Status: &Uncompleated,
			}
			_, err = s.EmployeeBenefitRepository.Update(ctx, &dataEbt.ID, &dataEb)
			if err != nil {
				fmt.Println(err)
				return err
			}
		}
		map1 := kafkaproducer.JsonData{
			UserID:    ctx.Auth.ID,
			CompanyID: ctx.Auth.CompanyID,
			Name:      ctx.Auth.Name,
			Timestamp: &waktu,
			Versions:  payload.Version,
			Berhasil:  len(*getStatusImportedWorkSheetDetailSucces),
			Gagal:     len(*getStatusImportedWorkSheetDetailFailed),
		}
		if len(errs) != 0 {
			map1.Data = "Proses Export Gagal!"
			jsonStr, err := json.Marshal(map1)
			if err != nil {
				log.Fatalln(err)
			}
			go kafkaproducer.NewProducer("NOTIFICATION").SendMessage("NOTIFICATION", string(jsonStr))

		}
		map1.Data = "Proses Import Telah Selesai"
		jsonStr, err := json.Marshal(map1)
		if err != nil {
			log.Fatalln(err)
		}
		go kafkaproducer.NewProducer("NOTIFICATION").SendMessage("NOTIFICATION", string(jsonStr))

		return nil

	}); err != nil {
		return "prosesGagal"
	}
	result := "proses Upload telah selesai"
	return result
	// tmpFolder := fmt.Sprintf("/mnt/d/Core-development/uploaded/%s",payload.TrialBalance)
	// defer func() {
	// 	if err := os.Remove(tmpFolder); err != nil {
	// 		panic(err)
	// 	}
	// }()
}
func (s *service) ImportTrialBalance(ctx *abstraction.Context, payload *abstraction.JsonDataReUpload, datas []model.TrialBalanceDetailEntity) (*dto.TrialBalanceImportResponse, error) {

	tmpFolder := payload.TrialBalance
	_, err := os.Stat(tmpFolder)
	if err != nil {
		if os.IsNotExist(err) {
			if err := os.MkdirAll(tmpFolder, os.ModePerm); err != nil {
				fmt.Println(err)
				return nil, err
			}
		} else {
			fmt.Println(err)
			return nil, err
		}
	}
	f, err := excelize.OpenFile(tmpFolder)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	sheet := f.GetSheetName(f.GetActiveSheetIndex())
	rows, err := f.GetRows(sheet)

	head1, err := f.GetCellValue(sheet, "B6")

	head2, err := f.GetCellValue(sheet, "C6")

	head3, err := f.GetCellValue(sheet, "F6")

	head4, err := f.GetCellValue(sheet, "G7")

	head5, err := f.GetCellValue(sheet, "H6")

	head6, err := f.GetCellValue(sheet, "I7")

	head7, err := f.GetCellValue(sheet, "J7")

	head8, err := f.GetCellValue(sheet, "K6")

	if strings.ToLower(head1) != "no akun" || strings.ToLower(head2) != "keterangan" || strings.ToLower(head3) != "wp reff" || strings.ToLower(head4) != "unaudited" || strings.ToLower(head5) != "adjustment journal entry" || strings.ToLower(head6) != "debet" || strings.ToLower(head7) != "kredit" || strings.ToLower(head8) != "unaudited" {

	}

	if err := trxmanager.New(s.Db).WithTrx(ctx, func(ctx *abstraction.Context) error {

		rows = rows[8:][:]

		datas = []model.TrialBalanceDetailEntity{}
		start := time.Now()

		var line [][]string
		var line1 [][]string
		for _, row := range rows {
			if len(row) == 0 || len(row) == 1 {
				line = append(line, row)
			}
			if len(row) > 1 && row[1] == "" {
				line1 = append(line1, row)
			}
			if len(row) > 1 && len(row) == 12 && row[1] != "Subtotal" && row[1] != "No Akun" && row[6] == "" && len(row[1]) > 1 {
				if (row[11]) == "" || (row[11]) == " " {
					row[11] = strings.Replace(strings.ToUpper(row[11]), "", "0", -1)
					row[11] = strings.Replace(strings.ToUpper(row[11]), " ", "0", -1)
				}
				amountAfterAje, err := strconv.ParseFloat(row[11], 64)
				if err != nil {
					errTrialBalance = err.Error()
					coa := strings.Replace(strings.ToUpper(errTrialBalance), "STRCONV.PARSEFLOAT: PARSING ", " ", -1)
					coas := strings.Replace(strings.ToUpper(coa), ": INVALID SYNTAX", " ", -1)
					msgErrTrialBalance = `Invalid Amount ` + `"` + coas + `"` + ` , Tolong Masukkan Inputan Yang Sesuai`
					jmlline := len(datas) + len(line) + len(line1) + 9
					lineMsg = "Template Excel row : " + strconv.Itoa(jmlline)

					return err
				}
				coass := strings.Replace(strings.ToUpper(row[1]), "+", "#", -1)
				coas := strings.Replace(strings.ToUpper(coass), " ", "_", -1)
				coa := strings.Replace(strings.ToUpper(coas), "-", "~", -1)

				getFormatterDetail, err := s.FormatterRepository.FindWithDetailFormatter(ctx, &coa)
				if err != nil {
					msgErrTrialBalance = "Invalid COA " + coa + ", Tolong Masukkan Inputan Yang Sesuai"
					jmlline := len(datas) + len(line) + len(line1) + 9
					lineMsg = "Template Excel row : " + strconv.Itoa(jmlline)

					return err
				}
				data := model.TrialBalanceDetailEntity{
					Code:           coa,
					Description:    &row[4],
					AmountAfterAje: &amountAfterAje,
					SortID:         getFormatterDetail.SortId,
				}
				datas = append(datas, data)
			}
			if len(row) > 1 && len(row) != 2 && row[1] != "Subtotal" && row[1] != "No Akun" && len(row[1]) > 1 {
				if len(row) < 11 {
					continue
				}
				if (row[6]) == "" {
					continue
				}

				nominal, err := strconv.ParseFloat(row[6], 64)
				if err != nil {
					errTrialBalance = err.Error()
					coa := strings.Replace(strings.ToUpper(errTrialBalance), "STRCONV.PARSEFLOAT: PARSING ", " ", -1)
					coas := strings.Replace(strings.ToUpper(coa), ": INVALID SYNTAX", " ", -1)
					msgErrTrialBalance = `Invalid Amount ` + `"` + coas + `"` + ` , Tolong Masukkan Inputan Yang Sesuai`
					jmlline := len(datas) + len(line) + len(line1) + 9
					lineMsg = "Template Excel row : " + strconv.Itoa(jmlline)

					return err
				}
				if (row[8]) == "" || row[8] == " " {
					row[8] = strings.Replace(strings.ToUpper(row[8]), "", "0", -1)
					row[8] = strings.Replace(strings.ToUpper(row[8]), " ", "0", -1)
				}

				if err != nil {
					fmt.Println("817")
					return err
				}
				if (row[10]) == "" || (row[10]) == " " {
					row[10] = strings.Replace(strings.ToUpper(row[10]), "", "0", -1)
					row[10] = strings.Replace(strings.ToUpper(row[10]), " ", "0", -1)
				}
				if err != nil {
					fmt.Println("826")
					return err
				}
				if (row[11]) == "" || (row[11]) == " " {
					row[11] = strings.Replace(strings.ToUpper(row[11]), "", "0", -1)
					row[11] = strings.Replace(strings.ToUpper(row[11]), " ", "0", -1)
				}

				if err != nil {
					fmt.Println("834")
					return err
				}
				rcr := " "
				rdr := " "

				cr, err := strconv.ParseFloat(row[8], 64)
				if err != nil {
					errTrialBalance = err.Error()
					coa := strings.Replace(strings.ToUpper(errTrialBalance), "STRCONV.PARSEFLOAT: PARSING ", " ", -1)
					coas := strings.Replace(strings.ToUpper(coa), ": INVALID SYNTAX", " ", -1)
					msgErrTrialBalance = `Invalid Amount ` + `"` + coas + `"` + ` , Tolong Masukkan Inputan Yang Sesuai`
					jmlline := len(datas) + len(line) + len(line1) + 9
					lineMsg = "Template Excel row : " + strconv.Itoa(jmlline)

					return err
				}
				dr, err := strconv.ParseFloat(row[10], 64)
				if err != nil {
					errTrialBalance = err.Error()
					coa := strings.Replace(strings.ToUpper(errTrialBalance), "STRCONV.PARSEFLOAT: PARSING ", " ", -1)
					coas := strings.Replace(strings.ToUpper(coa), ": INVALID SYNTAX", " ", -1)
					msgErrTrialBalance = `Invalid Amount ` + `"` + coas + `"` + ` , Tolong Masukkan Inputan Yang Sesuai`
					jmlline := len(datas) + len(line) + len(line1) + 9
					lineMsg = "Template Excel row : " + strconv.Itoa(jmlline)

					return err
				}
				amountAfterAje, err := strconv.ParseFloat(row[11], 64)
				if err != nil {
					errTrialBalance = err.Error()
					coa := strings.Replace(strings.ToUpper(errTrialBalance), "STRCONV.PARSEFLOAT: PARSING ", " ", -1)
					coas := strings.Replace(strings.ToUpper(coa), ": INVALID SYNTAX", " ", -1)
					msgErrTrialBalance = `Invalid Amount ` + `"` + coas + `"` + ` , Tolong Masukkan Inputan Yang Sesuai`
					jmlline := len(datas) + len(line) + len(line1) + 9
					lineMsg = "Template Excel row : " + strconv.Itoa(jmlline)

					return err
				}
				coass := strings.Replace(strings.ToUpper(row[1]), "+", "#", -1)
				coas := strings.Replace(strings.ToUpper(coass), " ", "_", -1)
				coa := strings.Replace(strings.ToUpper(coas), "-", "~", -1)
				b := coa

				if len(coa) == 9 && row[4] != "" {

					getFormatterDetail, err := s.FormatterRepository.FindWithDetailFormatter(ctx, &b)

					if err != nil && err != gorm.ErrRecordNotFound {
						fmt.Println(err)
						return err
					}
					if err == gorm.ErrRecordNotFound {
						b = coa[0:6]
						getFormatterDetail, err = s.FormatterRepository.FindWithDetailFormatter(ctx, &b)
						if err != nil && err != gorm.ErrRecordNotFound {
							msgErrTrialBalance = "Invalid COA " + coa + ", Tolong Masukkan Inputan Yang Sesuai"
							jmlline := len(datas) + len(line) + len(line1) + 9
							lineMsg = "Template Excel row : " + strconv.Itoa(jmlline)

							return err
						}
						if err == gorm.ErrRecordNotFound {
							b = coa[0:4]
							getFormatterDetail, err = s.FormatterRepository.FindWithDetailFormatter(ctx, &b)
							if err != nil && err != gorm.ErrRecordNotFound {
								msgErrTrialBalance = "Invalid COA " + coa + ", Tolong Masukkan Inputan Yang Sesuai"
								jmlline := len(datas) + len(line) + len(line1) + 9
								lineMsg = "Template Excel row : " + strconv.Itoa(jmlline)

								return err
							}
							if err == gorm.ErrRecordNotFound {
								b = coa[0:3]
								getFormatterDetail, err = s.FormatterRepository.FindWithDetailFormatter(ctx, &b)
								if err != nil {
									msgErrTrialBalance = "Invalid COA " + coa + ", Tolong Masukkan Inputan Yang Sesuai"
									jmlline := len(datas) + len(line) + len(line1) + 9
									lineMsg = "Template Excel row : " + strconv.Itoa(jmlline)

									return err
								}

							}

						}

					}
					// if row[4] == "" {
					// 	row[4] == row[1]
					// }
					nominalBeforeAje := nominal
					data := model.TrialBalanceDetailEntity{
						Code:            coa,
						AmountBeforeAje: &nominalBeforeAje,
						AmountAjeCr:     &dr,
						AmountAjeDr:     &cr,
						AmountAfterAje:  &amountAfterAje,
						Description:     &row[4],
						ReffAjeDr:       &rcr,
						ReffAjeCr:       &rdr,
						SortID:          getFormatterDetail.SortId,
					}
					datas = append(datas, data)
				}
				if len(coa) == 9 && row[4] == "" {

					getFormatterDetail, err := s.FormatterRepository.FindWithDetailFormatter(ctx, &b)

					if err != nil && err != gorm.ErrRecordNotFound {
						fmt.Println(err)
						return err
					}
					if err == gorm.ErrRecordNotFound {
						b = coa[0:6]
						getFormatterDetail, err = s.FormatterRepository.FindWithDetailFormatter(ctx, &b)
						if err != nil && err != gorm.ErrRecordNotFound {
							msgErrTrialBalance = "Invalid COA " + coa + ", Tolong Masukkan Inputan Yang Sesuai"
							jmlline := len(datas) + len(line) + len(line1) + 9
							lineMsg = "Template Excel row : " + strconv.Itoa(jmlline)

							return err
						}
						if err == gorm.ErrRecordNotFound {
							b = coa[0:4]
							getFormatterDetail, err = s.FormatterRepository.FindWithDetailFormatter(ctx, &b)
							if err != nil && err != gorm.ErrRecordNotFound {
								msgErrTrialBalance = "Invalid COA " + coa + ", Tolong Masukkan Inputan Yang Sesuai"
								jmlline := len(datas) + len(line) + len(line1) + 9
								lineMsg = "Template Excel row : " + strconv.Itoa(jmlline)

								return err
							}
							if err == gorm.ErrRecordNotFound {
								b = coa[0:3]
								getFormatterDetail, err = s.FormatterRepository.FindWithDetailFormatter(ctx, &b)
								if err != nil {
									msgErrTrialBalance = "Invalid COA " + coa + ", Tolong Masukkan Inputan Yang Sesuai"
									jmlline := len(datas) + len(line) + len(line1) + 9
									lineMsg = "Template Excel row : " + strconv.Itoa(jmlline)

									return err
								}

							}

						}

					}
					// if row[4] == "" {
					// 	row[4] == row[1]
					// }
					nominalBeforeAje := nominal
					data := model.TrialBalanceDetailEntity{
						Code:            coa,
						AmountBeforeAje: &nominalBeforeAje,
						AmountAjeCr:     &dr,
						AmountAjeDr:     &cr,
						AmountAfterAje:  &amountAfterAje,
						Description:     &row[1],
						ReffAjeDr:       &rcr,
						ReffAjeCr:       &rdr,
						SortID:          getFormatterDetail.SortId,
					}
					datas = append(datas, data)
				}
			
			
				if len(coa) > 9 && row[4] != ""{
					coass := strings.Replace(strings.ToUpper(row[1]), "+", "#", -1)
					coas := strings.Replace(strings.ToUpper(coass), " ", "_", -1)
					coa := strings.Replace(strings.ToUpper(coas), "-", "~", -1)
					getFormatterDetail, err := s.FormatterRepository.FindWithDetailFormatter(ctx, &coa)
					if err != nil && err != gorm.ErrRecordNotFound {
						msgErrTrialBalance = "Invalid COA " + coa + ", Tolong Masukkan Inputan Yang Sesuai"
						jmlline := len(datas) + len(line) + len(line1) + 9
						lineMsg = "Template Excel row : " + strconv.Itoa(jmlline)

						return err
					}
					if err == gorm.ErrRecordNotFound {
						j := coa[0:6]
						getFormatterDetail, err = s.FormatterRepository.FindWithDetailFormatter(ctx, &j)
						if err != nil && err != gorm.ErrRecordNotFound {
							fmt.Println(err)
							return err
						}
						if err == gorm.ErrRecordNotFound {
							j = coa[0:4]
							getFormatterDetail, err = s.FormatterRepository.FindWithDetailFormatter(ctx, &j)
							if err != nil && err != gorm.ErrRecordNotFound {
								msgErrTrialBalance = "Invalid COA " + coa + ", Tolong Masukkan Inputan Yang Sesuai"
								jmlline := len(datas) + len(line) + len(line1) + 9
								lineMsg = "Template Excel row : " + strconv.Itoa(jmlline)

								return err
							}
							if err == gorm.ErrRecordNotFound {
								j = coa[0:3]
								getFormatterDetail, err = s.FormatterRepository.FindWithDetailFormatter(ctx, &j)
								if err != nil {
									msgErrTrialBalance = "Invalid COA " + coa + ", Tolong Masukkan Inputan Yang Sesuai"
									jmlline := len(datas) + len(line) + len(line1) + 9
									lineMsg = "Template Excel row : " + strconv.Itoa(jmlline)

									return err
								}

							}
						}

					}
					nominalBeforeAje := nominal
					data := model.TrialBalanceDetailEntity{
						Code:            coa,
						AmountBeforeAje: &nominalBeforeAje,
						AmountAjeCr:     &dr,
						AmountAjeDr:     &cr,
						AmountAfterAje:  &amountAfterAje,
						Description:     &row[4],
						ReffAjeDr:       &rcr,
						ReffAjeCr:       &rdr,
						SortID:          getFormatterDetail.SortId,
					}
					datas = append(datas, data)
				}
				if len(coa) > 9 && row[4] == ""{
					coass := strings.Replace(strings.ToUpper(row[1]), "+", "#", -1)
					coas := strings.Replace(strings.ToUpper(coass), " ", "_", -1)
					coa := strings.Replace(strings.ToUpper(coas), "-", "~", -1)
					getFormatterDetail, err := s.FormatterRepository.FindWithDetailFormatter(ctx, &coa)
					if err != nil && err != gorm.ErrRecordNotFound {
						msgErrTrialBalance = "Invalid COA " + coa + ", Tolong Masukkan Inputan Yang Sesuai"
						jmlline := len(datas) + len(line) + len(line1) + 9
						lineMsg = "Template Excel row : " + strconv.Itoa(jmlline)

						return err
					}
					if err == gorm.ErrRecordNotFound {
						j := coa[0:6]
						getFormatterDetail, err = s.FormatterRepository.FindWithDetailFormatter(ctx, &j)
						if err != nil && err != gorm.ErrRecordNotFound {
							fmt.Println(err)
							return err
						}
						if err == gorm.ErrRecordNotFound {
							j = coa[0:4]
							getFormatterDetail, err = s.FormatterRepository.FindWithDetailFormatter(ctx, &j)
							if err != nil && err != gorm.ErrRecordNotFound {
								msgErrTrialBalance = "Invalid COA " + coa + ", Tolong Masukkan Inputan Yang Sesuai"
								jmlline := len(datas) + len(line) + len(line1) + 9
								lineMsg = "Template Excel row : " + strconv.Itoa(jmlline)

								return err
							}
							if err == gorm.ErrRecordNotFound {
								j = coa[0:3]
								getFormatterDetail, err = s.FormatterRepository.FindWithDetailFormatter(ctx, &j)
								if err != nil {
									msgErrTrialBalance = "Invalid COA " + coa + ", Tolong Masukkan Inputan Yang Sesuai"
									jmlline := len(datas) + len(line) + len(line1) + 9
									lineMsg = "Template Excel row : " + strconv.Itoa(jmlline)

									return err
								}

							}
						}

					}
					nominalBeforeAje := nominal
					data := model.TrialBalanceDetailEntity{
						Code:            coa,
						AmountBeforeAje: &nominalBeforeAje,
						AmountAjeCr:     &dr,
						AmountAjeDr:     &cr,
						AmountAfterAje:  &amountAfterAje,
						Description:     &row[1],
						ReffAjeDr:       &rcr,
						ReffAjeCr:       &rdr,
						SortID:          getFormatterDetail.SortId,
					}
					datas = append(datas, data)
				}
				if len(coa) < 9  && row[4] == "" {
					coass := strings.Replace(strings.ToUpper(row[1]), "+", "#", -1)
					coas := strings.Replace(strings.ToUpper(coass), " ", "_", -1)
					coa := strings.Replace(strings.ToUpper(coas), "-", "~", -1)
					getFormatterDetail, err := s.FormatterRepository.FindWithDetailFormatter(ctx, &coa)
					if err != nil && err != gorm.ErrRecordNotFound {
						msgErrTrialBalance = "Invalid COA " + coa + ", Tolong Masukkan Inputan Yang Sesuai"
						jmlline := len(datas) + len(line) + len(line1) + 9
						lineMsg = "Template Excel row : " + strconv.Itoa(jmlline)

						return err
					}
					if err == gorm.ErrRecordNotFound {
						j := coa[0:6]
						getFormatterDetail, err = s.FormatterRepository.FindWithDetailFormatter(ctx, &j)
						if err != nil && err != gorm.ErrRecordNotFound {
							fmt.Println(err)
							return err
						}
						if err == gorm.ErrRecordNotFound {
							j = coa[0:4]
							getFormatterDetail, err = s.FormatterRepository.FindWithDetailFormatter(ctx, &j)
							if err != nil && err != gorm.ErrRecordNotFound {
								msgErrTrialBalance = "Invalid COA " + coa + ", Tolong Masukkan Inputan Yang Sesuai"
								jmlline := len(datas) + len(line) + len(line1) + 9
								lineMsg = "Template Excel row : " + strconv.Itoa(jmlline)

								return err
							}
							if err == gorm.ErrRecordNotFound {
								j = coa[0:3]
								getFormatterDetail, err = s.FormatterRepository.FindWithDetailFormatter(ctx, &j)
								if err != nil {
									msgErrTrialBalance = "Invalid COA " + coa + ", Tolong Masukkan Inputan Yang Sesuai"
									jmlline := len(datas) + len(line) + len(line1) + 9
									lineMsg = "Template Excel row : " + strconv.Itoa(jmlline)

									return err
								}

							}
						}

					}
					nominalBeforeAje := nominal
					data := model.TrialBalanceDetailEntity{
						Code:            coa,
						AmountBeforeAje: &nominalBeforeAje,
						AmountAjeCr:     &dr,
						AmountAjeDr:     &cr,
						AmountAfterAje:  &amountAfterAje,
						Description:     &row[1],
						ReffAjeDr:       &rcr,
						ReffAjeCr:       &rdr,
						SortID:          getFormatterDetail.SortId,
					}
					datas = append(datas, data)
				}
				if len(coa) < 9  && row[4] != "" {
					coass := strings.Replace(strings.ToUpper(row[1]), "+", "#", -1)
					coas := strings.Replace(strings.ToUpper(coass), " ", "_", -1)
					coa := strings.Replace(strings.ToUpper(coas), "-", "~", -1)
					getFormatterDetail, err := s.FormatterRepository.FindWithDetailFormatter(ctx, &coa)
					if err != nil && err != gorm.ErrRecordNotFound {
						msgErrTrialBalance = "Invalid COA " + coa + ", Tolong Masukkan Inputan Yang Sesuai"
						jmlline := len(datas) + len(line) + len(line1) + 9
						lineMsg = "Template Excel row : " + strconv.Itoa(jmlline)

						return err
					}
					if err == gorm.ErrRecordNotFound {
						j := coa[0:6]
						getFormatterDetail, err = s.FormatterRepository.FindWithDetailFormatter(ctx, &j)
						if err != nil && err != gorm.ErrRecordNotFound {
							fmt.Println(err)
							return err
						}
						if err == gorm.ErrRecordNotFound {
							j = coa[0:4]
							getFormatterDetail, err = s.FormatterRepository.FindWithDetailFormatter(ctx, &j)
							if err != nil && err != gorm.ErrRecordNotFound {
								msgErrTrialBalance = "Invalid COA " + coa + ", Tolong Masukkan Inputan Yang Sesuai"
								jmlline := len(datas) + len(line) + len(line1) + 9
								lineMsg = "Template Excel row : " + strconv.Itoa(jmlline)

								return err
							}
							if err == gorm.ErrRecordNotFound {
								j = coa[0:3]
								getFormatterDetail, err = s.FormatterRepository.FindWithDetailFormatter(ctx, &j)
								if err != nil {
									msgErrTrialBalance = "Invalid COA " + coa + ", Tolong Masukkan Inputan Yang Sesuai"
									jmlline := len(datas) + len(line) + len(line1) + 9
									lineMsg = "Template Excel row : " + strconv.Itoa(jmlline)

									return err
								}

							}
						}

					}
					nominalBeforeAje := nominal
					data := model.TrialBalanceDetailEntity{
						Code:            coa,
						AmountBeforeAje: &nominalBeforeAje,
						AmountAjeCr:     &dr,
						AmountAjeDr:     &cr,
						AmountAfterAje:  &amountAfterAje,
						Description:     &row[4],
						ReffAjeDr:       &rcr,
						ReffAjeCr:       &rdr,
						SortID:          getFormatterDetail.SortId,
					}
					datas = append(datas, data)
				}
			}
			if len(row) > 1 && len(row) != 2 && row[1] != "Subtotal" && row[1] != "No Akun" && row[1] != "" && row[4] == "" && row[6] == "" && row[8] == "" {
				if len(row) < 11 {
					continue
				}
				if row[11] == "" {
					row[11] = "0"
				}
				amountAfterAje, err := strconv.ParseFloat(row[11], 64)
				if err != nil {
					errTrialBalance = err.Error()
					coa := strings.Replace(strings.ToUpper(errTrialBalance), "STRCONV.PARSEFLOAT: PARSING ", " ", -1)
					coas := strings.Replace(strings.ToUpper(coa), ": INVALID SYNTAX", " ", -1)
					msgErrTrialBalance = `Invalid Amount ` + `"` + coas + `"` + ` , Tolong Masukkan Inputan Yang Sesuai`
					jmlline := len(datas) + len(line) + len(line1) + 9
					lineMsg = "Template Excel row : " + strconv.Itoa(jmlline)

					return err
				}

				coass := strings.Replace(strings.ToUpper(row[1]), "+", "#", -1)
				coas := strings.Replace(strings.ToUpper(coass), " ", "_", -1)
				coa := strings.Replace(strings.ToUpper(coas), "-", "~", -1)
				getFormatterDetail, err := s.FormatterRepository.FindWithDetailFormatter(ctx, &coa)
				if err != nil && err != gorm.ErrRecordNotFound {
					fmt.Println(err)
					return err
				}
				if getFormatterDetail == nil {
					coa = coa[0:6]
					getFormatterDetail, err = s.FormatterRepository.FindWithDetailFormatter(ctx, &coa)
					if err != nil {
						msgErrTrialBalance = "Invalid COA " + coa + ", Tolong Masukkan Inputan Yang Sesuai"
						jmlline := len(datas) + len(line) + len(line1) + 9
						lineMsg = "Template Excel row : " + strconv.Itoa(jmlline)

						return err
					}
				}

					data := model.TrialBalanceDetailEntity{
						Code:           coa,
						AmountAfterAje: &amountAfterAje,
						Description:    &row[1],
						SortID:         getFormatterDetail.SortId,
					}
					datas = append(datas, data)
			}
			if len(row) > 1 && len(row) == 2 && row[1] != "Subtotal" && row[1] != "No Akun" {
				coass := strings.Replace(strings.ToUpper(row[1]), "+", "#", -1)
				coas := strings.Replace(strings.ToUpper(coass), " ", "_", -1)
				coa := strings.Replace(strings.ToUpper(coas), "-", "~", -1)
				// if len(coa) == 9 {
				// 	coa = coa[0:6]
				// }
				getFormatterDetail, err := s.FormatterRepository.FindWithDetailFormatter(ctx, &coa)
				if err != nil && err != gorm.ErrRecordNotFound {
					fmt.Println(err)
					return err
				}
				if err == gorm.ErrRecordNotFound {
					b := coa[0:6]
					getFormatterDetail, err = s.FormatterRepository.FindWithDetailFormatter(ctx, &b)
					if err != nil && err != gorm.ErrRecordNotFound {
						msgErrTrialBalance = "Invalid COA " + coa + ", Tolong Masukkan Inputan Yang Sesuai"
						jmlline := len(datas) + len(line) + len(line1) + 9
						lineMsg = "Template Excel row : " + strconv.Itoa(jmlline)

						return err
					}
					if err == gorm.ErrRecordNotFound {
						b = coa[0:4]
						getFormatterDetail, err = s.FormatterRepository.FindWithDetailFormatter(ctx, &b)
						if err != nil && err != gorm.ErrRecordNotFound {
							msgErrTrialBalance = "Invalid COA " + coa + ", Tolong Masukkan Inputan Yang Sesuai"
							jmlline := len(datas) + len(line) + len(line1) + 9
							lineMsg = "Template Excel row : " + strconv.Itoa(jmlline)

							return err
						}
						if err == gorm.ErrRecordNotFound {
							b = coa[0:3]
							getFormatterDetail, err = s.FormatterRepository.FindWithDetailFormatter(ctx, &b)
							if err != nil {
								msgErrTrialBalance = "Invalid COA " + coa + ", Tolong Masukkan Inputan Yang Sesuai"
								jmlline := len(datas) + len(line) + len(line1) + 9
								lineMsg = "Template Excel row : " + strconv.Itoa(jmlline)

								return err
							}

						}

					}

				}

				data := model.TrialBalanceDetailEntity{
					Code:        coa,
					Description: &row[1],
					SortID:      getFormatterDetail.SortId,
				}
				datas = append(datas, data)
			}
			if len(row) > 1 && (row[1]) == "Subtotal" {
				nominal, err := strconv.ParseFloat(row[6], 64)
				if err != nil {
					fmt.Println("1000")
					return err
				}

				if (row[6]) == "" {
					continue
				}

				if (row[8]) == "-" || (row[8]) == "" || (row[8]) == " " {
					strings.Replace(strings.ToUpper(row[8]), "-", "0", -1)
					strings.Replace(strings.ToUpper(row[8]), " ", "0", -1)
					strings.Replace(strings.ToUpper(row[8]), "", "0", -1)
				}
				if err != nil {
					return err
				}
				if (row[10]) == "-" || (row[10]) == "" || (row[10]) == " " {
					strings.Replace(strings.ToUpper(row[10]), "-", "0", -1)
					strings.Replace(strings.ToUpper(row[10]), "", "0", -1)
					strings.Replace(strings.ToUpper(row[10]), " ", "0", -1)
				}
				if err != nil {
					return err
				}
				if (row[11]) == "-" || (row[11]) == "" || (row[11]) == " " {
					strings.Replace(strings.ToUpper(row[11]), "-", "0", -1)
					strings.Replace(strings.ToUpper(row[11]), "", "0", -1)
					strings.Replace(strings.ToUpper(row[11]), " ", "0", -1)
				}

				if err != nil {
					return err
				}
				rcr := " "
				rdr := " "

				cr, err := strconv.ParseFloat(row[8], 64)
				if err != nil {
					errTrialBalance = err.Error()
					coa := strings.Replace(strings.ToUpper(errTrialBalance), "STRCONV.PARSEFLOAT: PARSING ", " ", -1)
					coas := strings.Replace(strings.ToUpper(coa), ": INVALID SYNTAX", " ", -1)
					msgErrTrialBalance = `Invalid Amount ` + `"` + coas + `"` + ` , Tolong Masukkan Inputan Yang Sesuai`
					jmlline := len(datas) + len(line) + len(line1) + 9
					lineMsg = "Template Excel row : " + strconv.Itoa(jmlline)

					return err
				}
				dr, err := strconv.ParseFloat(row[10], 64)
				if err != nil {
					errTrialBalance = err.Error()
					coa := strings.Replace(strings.ToUpper(errTrialBalance), "STRCONV.PARSEFLOAT: PARSING ", " ", -1)
					coas := strings.Replace(strings.ToUpper(coa), ": INVALID SYNTAX", " ", -1)
					msgErrTrialBalance = `Invalid Amount ` + `"` + coas + `"` + ` , Tolong Masukkan Inputan Yang Sesuai`
					jmlline := len(datas) + len(line) + len(line1) + 9
					lineMsg = "Template Excel row : " + strconv.Itoa(jmlline)

					return err
				}
				amountAfterAje, err := strconv.ParseFloat(row[11], 64)
				if err != nil {
					errTrialBalance = err.Error()
					coa := strings.Replace(strings.ToUpper(errTrialBalance), "STRCONV.PARSEFLOAT: PARSING ", " ", -1)
					coas := strings.Replace(strings.ToUpper(coa), ": INVALID SYNTAX", " ", -1)
					msgErrTrialBalance = `Invalid Amount ` + `"` + coas + `"` + ` , Tolong Masukkan Inputan Yang Sesuai`
					jmlline := len(datas) + len(line) + len(line1) + 9
					lineMsg = "Template Excel row : " + strconv.Itoa(jmlline)

					return err
				}
				coa := row[1]
				c := len(datas) - 1
				d := datas[c].Code
				j := d[0:6]
				coa = j + "_" + coa
				nominalBeforeAje := nominal
				desc := "Sub Total"

				getFormatterDetail, err := s.FormatterRepository.FindWithDetailFormatter(ctx, &j)

				if err == gorm.ErrRecordNotFound {
					j = coa[0:4]
					getFormatterDetail, err = s.FormatterRepository.FindWithDetailFormatter(ctx, &j)
					if err != nil && err != gorm.ErrRecordNotFound {
						fmt.Println(err)
						return err
					}
					if err == gorm.ErrRecordNotFound {
						j = coa[0:3]
						getFormatterDetail, err = s.FormatterRepository.FindWithDetailFormatter(ctx, &j)
						if err != nil && err != gorm.ErrRecordNotFound {
							// msgErrTrialBalance = "Invalid COA " + coa + ", Tolong Masukkan Inputan Yang Sesuai"
							// jmlline := len(datas) + len(line) + len(line1) + 9
							// lineMsg = "Template Excel row : " + strconv.Itoa(jmlline)

							return err
						}

					}

				}

				data := model.TrialBalanceDetailEntity{
					Code:            coa,
					AmountBeforeAje: &nominalBeforeAje,
					AmountAjeCr:     &dr,
					AmountAjeDr:     &cr,
					AmountAfterAje:  &amountAfterAje,
					ReffAjeDr:       &rcr,
					ReffAjeCr:       &rdr,
					Description:     &desc,
					SortID:          getFormatterDetail.SortId,
				}
				datas = append(datas, data)
			}
		}

		duration := time.Since(start)
		fmt.Println("done in", int(math.Ceil(duration.Seconds())), "seconds")

		criteriaFormatter := model.FormatterFilterModel{}
		tmpStr := "TB-CONSOLIDATION"
		criteriaFormatter.FormatterFor = &tmpStr

		getFormatterID, err := s.FormatterRepository.FindWithDetail(ctx, &criteriaFormatter)
		if err != nil {
			return err
		}

		var dataFormatterBridgeds model.FormatterBridgesEntityModel

		criteriaFB := model.FormatterBridgesFilterModel{}
		criteriaFB.Source = &tmpStr
		criteriaFB.FormatterID = &getFormatterID.ID
		criteriaFB.TrxRefID = &payload.IDTrialBalance

		dataFormatterBridgeds.Context = ctx
		dataFormatterBridgeds.FormatterBridgesEntity = model.FormatterBridgesEntity{
			TrxRefID:    payload.IDTrialBalance,
			FormatterID: getFormatterID.ID,
			Source:      "TRIAL-BALANCE",
		}
		source := "TRIAL-BALANCE"
		trxRefId := &payload.IDTrialBalance
		resultFormatterBridge, err := s.FormatterBridgesRepository.FindByIDTrx(ctx, &source, trxRefId)
		if err != nil {
			return err
		}
		var arrDataTBD []model.TrialBalanceDetailEntityModel

		for _, v := range datas {
			dataTBD := model.TrialBalanceDetailEntityModel{
				Context:                  ctx,
				TrialBalanceDetailEntity: v,
			}
			dataTBD.FormatterBridgesID = resultFormatterBridge.ID
			// dataTBD.SortID = getFormatterID.FormatterDetail[0].
			arrDataTBD = append(arrDataTBD, dataTBD)
		}
		_, err = s.TrialBalanceDetailRepository.Import(ctx, &arrDataTBD)
		if err != nil {
			return err
		}

		// dataTB = *resultTB
		var data model.ImportedWorksheetDetailEntityModel
		_, err = s.ImportedWorksheetDetailRepository.FindByID(ctx, &payload.IDWorksheetDetailTrialBalance)
		if err != nil {
			return err
		}
		jsonEr, err := json.Marshal("")
		if err != nil {
			fmt.Println(err)
			return err
		}
		data.Context = ctx
		data.ImportedWorksheetDetailEntity = model.ImportedWorksheetDetailEntity{
			Status: 2,
			ErrMessages: string(jsonEr),
		}

		_, err = s.ImportedWorksheetDetailRepository.Update(ctx, &payload.IDWorksheetDetailTrialBalance, &data)
		if err != nil {
			return err
		}
		return nil
	}); err != nil {
		return &dto.TrialBalanceImportResponse{}, err
	}
	result := &dto.TrialBalanceImportResponse{}
	return result, nil
}
func (s *service) ImportAgingUtangPiutang(ctx *abstraction.Context, payload *abstraction.JsonDataReUpload, datas []model.AgingUtangPiutangDetailEntity) (*dto.AgingUtangPiutangImportResponse, error) {
	tmpFolder := payload.AgingUtangPiutang
	_, err := os.Stat(tmpFolder)
	if err != nil {
		if os.IsNotExist(err) {
			if err := os.MkdirAll(tmpFolder, os.ModePerm); err != nil {
				log.Fatal(err)
			}
		} else {
			panic(err)
		}
	}

	f, err := excelize.OpenFile(tmpFolder)

	defer func() {
		if err := f.Close(); err != nil {
			fmt.Println(err)
		}
	}()

	sheet := f.GetSheetName(f.GetActiveSheetIndex())
	rows, err := f.GetRows(sheet)
	if err != nil {
		panic(err)
	}

	confRowStart := []string{"AGING-UTANG-PIUTANG-IMPORT-ROW-START", "AGING-UTANG-PIUTANG-MUTASI-ECL-IMPORT-ROW-START"}
	confRowEnd := []string{"AGING-UTANG-PIUTANG-IMPORT-ROW-END", "AGING-UTANG-PIUTANG-MUTASI-ECL-IMPORT-ROW-END"}
	allData := [][]model.AgingUtangPiutangDetailEntity{}
	counter := 1

	if err := trxmanager.New(s.Db).WithTrx(ctx, func(ctx *abstraction.Context) error {
		c2, err := f.GetCellValue(sheet, "B2")
		if c2 != "Detail aging" {
			return errors.New("caannot")
		}
		for i, conf := range confRowStart {
			confCode := conf
			importConfigStartRow, err := s.ParameterRepository.FindWithCode(ctx, &confCode)
			if err != nil || len(*importConfigStartRow) == 0 {
				errs = append(errs, err)
				log.Fatalln(err)
				return err
			}

			startRow := 0

			for _, v := range *importConfigStartRow {
				startRow, err = strconv.Atoi(*v.Value)
				if err != nil {
					startRow = 0
				}
			}

			confCode = confRowEnd[i]
			importConfigEndRow, err := s.ParameterRepository.FindWithCode(ctx, &confCode)
			if err != nil || len(*importConfigEndRow) == 0 {
				errs = append(errs, err)
				log.Fatalln(err)
				return err
			}

			endRow := 0

			for _, v := range *importConfigEndRow {
				endRow, err = strconv.Atoi(*v.Value)
				if err != nil {
					endRow = 0
				}
			}

			rowsData := rows[startRow-1 : endRow]
			datas := []model.AgingUtangPiutangDetailEntity{}
			for _, row := range rowsData {
				if len(row) == 0 {
					continue
				}

				tmp := make([]float64, len(row))
				for i := 2; i < len(row); i++ {
					if i == 10 {
						continue
					}
					tmpi, err := strconv.ParseFloat(row[i], 64)
					if err != nil {
						tmpi = 0
					}
					tmp[i-2] = tmpi
				}

				data := model.AgingUtangPiutangDetailEntity{
					Description:                  row[1],
					Code:                         strings.Replace(strings.ToUpper(row[1]), " ", "_", -1),
					Piutangusaha3rdparty:         &tmp[0],
					PiutangusahaBerelasi:         &tmp[1],
					Piutanglainshortterm3rdparty: &tmp[2],
					PiutanglainshorttermBerelasi: &tmp[3],
					Piutangberelasishortterm:     &tmp[4],
					Piutanglainlongterm3rdparty:  &tmp[5],
					PiutanglainlongtermBerelasi:  &tmp[6],
					Piutangberelasilongterm:      &tmp[7],
					Utangusaha3rdparty:           &tmp[9],
					UtangusahaBerelasi:           &tmp[10],
					Utanglainshortterm3rdparty:   &tmp[11],
					UtanglainshorttermBerelasi:   &tmp[12],
					Utangberelasishortterm:       &tmp[13],
					Utanglainlongterm3rdparty:    &tmp[14],
					UtanglainlongtermBerelasi:    &tmp[15],
					Utangberelasilongterm:        &tmp[16],

					SortID: counter,
				}
				datas = append(datas, data)
				counter++
			}
			allData = append(allData, datas)
		}

		var dataTB []model.AgingUtangPiutangEntityModel
		currentYear, currentMonth, _ := time.Now().Date()
		currentLocation := time.Now().Location()
		firstOfMonth := time.Date(currentYear, currentMonth, 1, 0, 0, 0, 0, currentLocation)
		lastOfMonth := firstOfMonth.AddDate(0, 1, -1)
		period := lastOfMonth.Format("2006-01-02")
		//cek company berdasarkan user
		//belum ada
		//skip

		criteriaAUP := model.AgingUtangPiutangFilterModel{}
		criteriaAUP.Period = &period

		for i, datas := range allData {
			criteriaFormatter := model.FormatterFilterModel{}
			tmpConf := strings.Split(confRowStart[i], "-IMPORT-")
			tmpStr := tmpConf[0]
			criteriaFormatter.FormatterFor = &tmpStr

			getFormatter, err := s.FormatterRepository.Find(ctx, &criteriaFormatter)
			if err != nil {
				errs = append(errs, err)
				log.Fatalln(err)
				return err
			}

			getFormatterID, err := s.FormatterRepository.FindWithDetail(ctx, &criteriaFormatter)
			if err != nil {
				return err
			}
			formatterID := 0
			for _, tmpFormatter := range *getFormatter {
				formatterID = tmpFormatter.ID
			}
			if formatterID == 0 {
				fmt.Println("No Formatter found")
				errs = append(errs, err)
				log.Fatalln(err)
				return err
			}

			temptAging := "AGING-UTANG-PIUTANG"
			var dataFormatterBridgeds model.FormatterBridgesEntityModel

			dataFormatterBridgeds.Context = ctx
			dataFormatterBridgeds.FormatterBridgesEntity = model.FormatterBridgesEntity{
				TrxRefID:    payload.IDAgingUtangPiutang,
				FormatterID: getFormatterID.ID,
				Source:      temptAging,
			}
			resultFormatterBridge, err := s.FormatterBridgesRepository.Create(ctx, &dataFormatterBridgeds)
			if err != nil {
				return err
			}

			for _, v := range datas {
				dataTBD := model.AgingUtangPiutangDetailEntityModel{
					Context:                       ctx,
					AgingUtangPiutangDetailEntity: v,
				}
				dataTBD.FormatterBridgesID = resultFormatterBridge.ID
				_, err := s.AgingUPDetailRepository.Create(ctx, &dataTBD)
				if err != nil {
					return err
				}
			}
			dataTB = append(dataTB)
		}
		var data model.ImportedWorksheetDetailEntityModel
		_, err = s.ImportedWorksheetDetailRepository.FindByID(ctx, &payload.IDWorksheetDetailAgingUtangPiutang)
		if err != nil {
			return err
		}
		data.Context = ctx
		data.ImportedWorksheetDetailEntity = model.ImportedWorksheetDetailEntity{
			Status: 2,
		}

		_, err = s.ImportedWorksheetDetailRepository.Update(ctx, &payload.IDWorksheetDetailAgingUtangPiutang, &data)
		if err != nil {
			return err
		}
		return nil

	}); err != nil {
		return &dto.AgingUtangPiutangImportResponse{}, err
	}
	result := &dto.AgingUtangPiutangImportResponse{}
	return result, nil
}
func (s *service) ImportMutasiFA(ctx *abstraction.Context, payload *abstraction.JsonDataReUpload, datas []model.MutasiFaDetailEntity) (*dto.MutasiFaImportResponse, error) {
	tmpFolder := payload.MutasiFA
	_, err := os.Stat(tmpFolder)
	if err != nil {
		if os.IsNotExist(err) {
			if err := os.MkdirAll(tmpFolder, os.ModePerm); err != nil {
				fmt.Println(err)
				return nil, err
			}
		} else {
			fmt.Println(err)
			return nil, err
		}
	}

	f, err := excelize.OpenFile(tmpFolder)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	sheet := f.GetSheetName(f.GetActiveSheetIndex())
	rows, err := f.GetRows(sheet)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	confRowStart := []string{"MUTASI-FA-COST-IMPORT-ROW-START", "MUTASI-FA-ACCUMULATED-DEPRECATION-IMPORT-ROW-START"}
	confRowEnd := []string{"MUTASI-FA-COST-IMPORT-ROW-END", "MUTASI-FA-ACCUMULATED-DEPRECATION-IMPORT-ROW-END"}
	allData := [][]model.MutasiFaDetailEntity{}
	counter := 1

	if err := trxmanager.New(s.Db).WithTrx(ctx, func(ctx *abstraction.Context) error {
		c2, err := f.GetCellValue(sheet, "B2")
		if c2 != "Mutasi Fixed Assets (FA)" {
			return errors.New("caannot")
		}
		for i, conf := range confRowStart {
			confCode := conf
			importConfigStartRow, err := s.ParameterRepository.FindWithCode(ctx, &confCode)
			if err != nil || len(*importConfigStartRow) == 0 {
				errs = append(errs, err)
				log.Fatalln(err)
				return err
			}

			startRow := 0

			for _, v := range *importConfigStartRow {
				startRow, err = strconv.Atoi(*v.Value)
				if err != nil {
					startRow = 0
				}
			}

			confCode = confRowEnd[i]
			importConfigEndRow, err := s.ParameterRepository.FindWithCode(ctx, &confCode)
			if err != nil || len(*importConfigEndRow) == 0 {
				errs = append(errs, err)
				log.Fatalln(err)
				return err
			}

			endRow := 0

			for _, v := range *importConfigEndRow {
				endRow, err = strconv.Atoi(*v.Value)
				if err != nil {
					endRow = 0
				}
			}

			rowsData := rows[startRow-1 : endRow]
			datas := []model.MutasiFaDetailEntity{}
			for _, row := range rowsData {
				if len(row) == 0 {
					continue
				}
				tmp := make([]float64, len(row))
				for i := 3; i < len(row); i++ {
					if len(row[i]) == 0 {
						tmp[i-3] = 0
						continue
					}
					tmpi, err := strconv.ParseFloat(row[i], 64)
					if err != nil {
						tmpi = 0
					}
					tmp[i-3] = tmpi
				}
		if row[1] == "Cost:" || row[1] == "Accumulated Depreciation" {
			coa := strings.Replace(strings.ToUpper(row[1]), " ", "_", -1)
			codecoa := strings.Replace(strings.ToUpper(coa), "-", "~", -1)
			data := model.MutasiFaDetailEntity{
				Description:             row[1],
				Code:                    codecoa,
				// BeginningBalance:        &tmp[0],
				// AcquisitionOfSubsidiary: &tmp[1],
				// Additions:               &tmp[2],
				// Deductions:              &tmp[3],
				// Reclassification:        &tmp[4],
				// Revaluation:             &tmp[5],
				// EndingBalance:           &tmp[6],
				// Control:                 &tmp[7],
				SortId:                  counter,
			}
			datas = append(datas, data)
			counter++
		}

		if row[1] != "Cost:" && row[1] != "Accumulated Depreciation" {
			coa := strings.Replace(strings.ToUpper(row[1]), " ", "_", -1)
			codecoa := strings.Replace(strings.ToUpper(coa), "-", "~", -1)
			data := model.MutasiFaDetailEntity{
				Description:             row[1],
				Code:                    codecoa,
				BeginningBalance:        &tmp[0],
				AcquisitionOfSubsidiary: &tmp[1],
				Additions:               &tmp[2],
				Deductions:              &tmp[3],
				Reclassification:        &tmp[4],
				Revaluation:             &tmp[5],
				EndingBalance:           &tmp[6],
				Control:                 &tmp[7],
				SortId:                  counter,
			}
			datas = append(datas, data)
			counter++
		}
			
			}
			allData = append(allData, datas)
		}

		var dataTB []model.MutasiFaEntityModel
		currentYear, currentMonth, _ := time.Now().Date()
		currentLocation := time.Now().Location()
		firstOfMonth := time.Date(currentYear, currentMonth, 1, 0, 0, 0, 0, currentLocation)
		lastOfMonth := firstOfMonth.AddDate(0, 1, -1)
		period := lastOfMonth.Format("2006-01-02")

		//cek company berdasarkan user
		//belum ada
		//skip

		companyID := 1
		criteriaTB := model.MutasiFaFilterModel{}
		criteriaTB.Period = &period
		criteriaTB.CompanyID = &companyID

		for i, datas := range allData {
			criteriaFormatter := model.FormatterFilterModel{}
			tmpConf := strings.Split(confRowStart[i], "-IMPORT-")
			tmpStr := tmpConf[0]
			criteriaFormatter.FormatterFor = &tmpStr

			getFormatter, err := s.FormatterRepository.Find(ctx, &criteriaFormatter)
			if err != nil {
				errs = append(errs, err)
				log.Fatalln(err)
				return err
			}

			getFormatterID, err := s.FormatterRepository.FindWithDetail(ctx, &criteriaFormatter)
			if err != nil {
				return err
			}

			formatterID := 0
			for _, tmpFormatter := range *getFormatter {
				formatterID = tmpFormatter.ID
			}
			if formatterID == 0 {
				fmt.Println("No Formatter found")
				errs = append(errs, err)
				log.Fatalln(err)
				return err
			}

			var dataFormatterBridgeds model.FormatterBridgesEntityModel

			criteriaFB := model.FormatterBridgesFilterModel{}
			criteriaFB.Source = &tmpStr
			criteriaFB.FormatterID = &formatterID
			criteriaFB.TrxRefID = &payload.IDMutasiFA

			temptMFA := "MUTASI-FA"
			dataFormatterBridgeds.Context = ctx
			dataFormatterBridgeds.FormatterBridgesEntity = model.FormatterBridgesEntity{
				TrxRefID:    payload.IDMutasiFA,
				FormatterID: getFormatterID.ID,
				Source:      temptMFA,
			}
			resultFormatterBridge, err := s.FormatterBridgesRepository.Create(ctx, &dataFormatterBridgeds)
			if err != nil {
				return err
			}

			for _, v := range datas {
				dataTBD := model.MutasiFaDetailEntityModel{
					Context:              ctx,
					MutasiFaDetailEntity: v,
				}
				dataTBD.FormatterBridgesID = resultFormatterBridge.ID
				_, err = s.MutasiFaDetailRepository.Create(ctx, &dataTBD)
				if err != nil {
					errs = append(errs, err)
					log.Fatalln(err)
					return err
				}
			}
			dataTB = append(dataTB)
		}
		var data model.ImportedWorksheetDetailEntityModel
		_, err = s.ImportedWorksheetDetailRepository.FindByID(ctx, &payload.IDWorksheetDetailMutasiFA)
		if err != nil {
			return err
		}
		data.Context = ctx
		data.ImportedWorksheetDetailEntity = model.ImportedWorksheetDetailEntity{
			Status: 2,
		}

		_, err = s.ImportedWorksheetDetailRepository.Update(ctx, &payload.IDWorksheetDetailMutasiFA, &data)
		if err != nil {
			return err
		}
		return nil
	}); err != nil {
		return &dto.MutasiFaImportResponse{}, err
	}
	result := &dto.MutasiFaImportResponse{}
	return result, nil
}
func (s *service) ImportInvestasiTbk(ctx *abstraction.Context, payload *abstraction.JsonDataReUpload, datas []model.InvestasiTbkDetailEntity) (*dto.InvestasiTbkImportResponse, error) {

	tmpFolder := payload.InvestasiTbk
	_, err := os.Stat(tmpFolder)
	if err != nil {
		if os.IsNotExist(err) {
			if err := os.MkdirAll(tmpFolder, os.ModePerm); err != nil {
				log.Fatal(err)
			}
		} else {
			panic(err)
		}
	}

	f, err := excelize.OpenFile(tmpFolder)

	defer func() {
		if err := f.Close(); err != nil {
			fmt.Println(err)
		}
	}()

	sheet := f.GetSheetName(f.GetActiveSheetIndex())
	rows, err := f.GetRows(sheet)
	if err != nil {
		panic(err)
	}

	confRowStart := []string{"INVESTASI-TBK-IMPORT-ROW-START"}
	confRowEnd := []string{"INVESTASI-TBK-IMPORT-ROW-END"}
	allData := [][]model.InvestasiTbkDetailEntity{}
	counter := 1

	if err := trxmanager.New(s.Db).WithTrx(ctx, func(ctx *abstraction.Context) error {
		c2, err := f.GetCellValue(sheet, "B2")
		if c2 != "Summary Investasi TBK" {
			return errors.New("caannot")
		}
		for i, conf := range confRowStart {
			confCode := conf
			importConfigStartRow, err := s.ParameterRepository.FindWithCode(ctx, &confCode)
			if err != nil || len(*importConfigStartRow) == 0 {
				panic(err)
			}

			startRow := 0

			for _, v := range *importConfigStartRow {
				startRow, err = strconv.Atoi(*v.Value)
				if err != nil {
					startRow = 0
				}
			}

			confCode = confRowEnd[i]
			importConfigEndRow, err := s.ParameterRepository.FindWithCode(ctx, &confCode)
			if err != nil || len(*importConfigEndRow) == 0 {
				panic(err)
			}

			endRow := 0

			for _, v := range *importConfigEndRow {
				endRow, err = strconv.Atoi(*v.Value)
				if err != nil {
					endRow = 0
				}
			}

			rowsData := rows[startRow-1 : endRow]
			datas := []model.InvestasiTbkDetailEntity{}
			for _, row := range rowsData {
				if len(row) == 0 {
					continue
				}

				tmp := make([]float64, len(row))
				for i := 2; i < len(row); i++ {
					if len(row[i]) == 0 {
						tmp[i-3] = 0
						continue
					}
					tmpi, err := strconv.ParseFloat(row[i], 64)
					if err != nil {
						tmpi = 0
					}
					tmp[i-2] = tmpi
				}
				data := model.InvestasiTbkDetailEntity{
					Stock:          row[2],
					EndingShares:   &tmp[1],
					AvgPrice:       &tmp[2],
					AmountCost:     &tmp[3],
					ClosingPrice:   &tmp[4],
					AmountFv:       &tmp[5],
					UnrealizedGain: &tmp[6],
					RealizedGain:   &tmp[7],
					Fee:            &tmp[8],
					SortId:         counter,
				}
				datas = append(datas, data)
				counter++
			}
			allData = append(allData, datas)
		}
		var dataTB []model.InvestasiTbkEntityModel
		currentYear, currentMonth, _ := time.Now().Date()
		currentLocation := time.Now().Location()
		firstOfMonth := time.Date(currentYear, currentMonth, 1, 0, 0, 0, 0, currentLocation)
		lastOfMonth := firstOfMonth.AddDate(0, 1, -1)
		period := lastOfMonth.Format("2006-01-02")
		//cek company berdasarkan user
		//belum ada
		//skip
		companyID := 1
		criteriaTB := model.InvestasiTbkFilterModel{}
		criteriaTB.Period = &period
		criteriaTB.CompanyID = &companyID

		for i, datas := range allData {
			criteriaFormatter := model.FormatterFilterModel{}
			tmpConf := strings.Split(confRowStart[i], "-IMPORT-")
			tmpStr := tmpConf[0]
			criteriaFormatter.FormatterFor = &tmpStr

			getFormatter, err := s.FormatterRepository.Find(ctx, &criteriaFormatter)
			if err != nil {
				return err
			}
			getFormatterID, err := s.FormatterRepository.FindWithDetail(ctx, &criteriaFormatter)
			if err != nil {
				return err
			}

			formatterID := 0
			for _, tmpFormatter := range *getFormatter {
				formatterID = tmpFormatter.ID
			}
			if formatterID == 0 {
				fmt.Println("No Formatter found")
				return err
			}

			var dataFormatterBridgeds model.FormatterBridgesEntityModel

			criteriaFB := model.FormatterBridgesFilterModel{}
			criteriaFB.Source = &tmpStr
			criteriaFB.FormatterID = &formatterID
			criteriaFB.TrxRefID = &payload.IDInvestasiTbk
			temptTBK := "INVESTASI-TBK"
			dataFormatterBridgeds.Context = ctx
			dataFormatterBridgeds.FormatterBridgesEntity = model.FormatterBridgesEntity{
				TrxRefID:    payload.IDInvestasiTbk,
				FormatterID: getFormatterID.ID,
				Source:      temptTBK,
			}
			resultFormatterBridge, err := s.FormatterBridgesRepository.Create(ctx, &dataFormatterBridgeds)
			if err != nil {
				return err
			}

			for _, v := range datas {
				dataTBD := model.InvestasiTbkDetailEntityModel{
					Context:                  ctx,
					InvestasiTbkDetailEntity: v,
				}
				dataTBD.FormatterBridgesID = resultFormatterBridge.ID
				_, err := s.InvestasiTbkDetailRepository.Create(ctx, &dataTBD)
				if err != nil {
					return err
				}
			}
			dataTB = append(dataTB)
		}
		var data model.ImportedWorksheetDetailEntityModel
		_, err = s.ImportedWorksheetDetailRepository.FindByID(ctx, &payload.IDWorksheetDetailInvestasiTbk)
		if err != nil {
			return err
		}
		data.Context = ctx
		data.ImportedWorksheetDetailEntity = model.ImportedWorksheetDetailEntity{
			Status: 2,
		}

		_, err = s.ImportedWorksheetDetailRepository.Update(ctx, &payload.IDWorksheetDetailInvestasiTbk, &data)
		if err != nil {
			return err
		}
		return nil
	}); err != nil {
		return &dto.InvestasiTbkImportResponse{}, err
	}
	result := &dto.InvestasiTbkImportResponse{}
	return result, nil
}
func (s *service) ImportMutasiDta(ctx *abstraction.Context, payload *abstraction.JsonDataReUpload, datas []model.MutasiDtaDetailEntity) (*dto.MutasiDtaImportResponse, error) {

	tmpFolder := payload.MutasiDta
	_, err := os.Stat(tmpFolder)
	if err != nil {
		if os.IsNotExist(err) {
			if err := os.MkdirAll(tmpFolder, os.ModePerm); err != nil {
				log.Fatal(err)
			}
		} else {
			panic(err)
		}
	}

	f, err := excelize.OpenFile(tmpFolder)
	sheet := f.GetSheetName(f.GetActiveSheetIndex())
	rows, err := f.GetRows(sheet)
	if err != nil {
		panic(err)
	}
	confRowStart := []string{"MUTASI-DTA-IMPORT-ROW-START"}
	confRowEnd := []string{"MUTASI-DTA-IMPORT-ROW-END"}
	allData := [][]model.MutasiDtaDetailEntity{}
	counter := 1

	if err := trxmanager.New(s.Db).WithTrx(ctx, func(ctx *abstraction.Context) error {
		c2, err := f.GetCellValue(sheet, "B2")
		if c2 != "Mutasi DTA" {
			return errors.New("caannot")
		}
		for i, conf := range confRowStart {
			confCode := conf
			importConfigStartRow, err := s.ParameterRepository.FindWithCode(ctx, &confCode)
			if err != nil || len(*importConfigStartRow) == 0 {
				panic(err)
			}

			startRow := 0

			for _, v := range *importConfigStartRow {
				startRow, err = strconv.Atoi(*v.Value)
				if err != nil {
					startRow = 0
				}
			}

			confCode = confRowEnd[i]
			importConfigEndRow, err := s.ParameterRepository.FindWithCode(ctx, &confCode)
			if err != nil || len(*importConfigEndRow) == 0 {
				panic(err)
			}

			endRow := 0

			for _, v := range *importConfigEndRow {
				endRow, err = strconv.Atoi(*v.Value)
				if err != nil {
					endRow = 0
				}
			}

			rowsData := rows[startRow-1 : endRow]
			datas := []model.MutasiDtaDetailEntity{}
			for _, row := range rowsData {
				if len(row) == 0 {
					continue
				}

				tmp := make([]float64, len(row))
				for i := 2; i < len(row); i++ {
					if len(row[i]) == 0 {
						tmp[i-2] = 0
						continue
					}
					tmpi, err := strconv.ParseFloat(row[i], 64)
					if err != nil {
						tmpi = 0
					}
					tmp[i-2] = tmpi
				}

				data := model.MutasiDtaDetailEntity{
					Description:         row[2],
					Code:                strings.Replace(strings.ToUpper(row[2]), " ", "_", -1),
					SaldoAwal:           &tmp[1],
					ManfaatBebanPajak:   &tmp[2],
					Oci:                 &tmp[3],
					AkuisisiEntitasAnak: &tmp[4],
					DibebankanKeLr:      &tmp[5],
					DibebankanKeOci:     &tmp[6],
					SaldoAkhir:          &tmp[7],
					SortId:              counter,
				}
				datas = append(datas, data)
				counter++
			}
			allData = append(allData, datas)
		}

		var dataTB []model.MutasiDtaEntityModel
		currentYear, currentMonth, _ := time.Now().Date()
		currentLocation := time.Now().Location()
		firstOfMonth := time.Date(currentYear, currentMonth, 1, 0, 0, 0, 0, currentLocation)
		lastOfMonth := firstOfMonth.AddDate(0, 1, -1)
		period := lastOfMonth.Format("2006-01-02")
		//cek company berdasarkan user
		//belum ada
		//skip
		companyID := 1
		criteriaTB := model.MutasiDtaFilterModel{}
		criteriaTB.Period = &period
		criteriaTB.CompanyID = &companyID

		for i, datas := range allData {
			criteriaFormatter := model.FormatterFilterModel{}
			tmpConf := strings.Split(confRowStart[i], "-IMPORT-")
			tmpStr := tmpConf[0]
			criteriaFormatter.FormatterFor = &tmpStr

			getFormatter, err := s.FormatterRepository.Find(ctx, &criteriaFormatter)
			if err != nil {
				return err
			}
			getFormatterID, err := s.FormatterRepository.FindWithDetail(ctx, &criteriaFormatter)
			if err != nil {
				return err
			}
			formatterID := 0
			for _, tmpFormatter := range *getFormatter {
				formatterID = tmpFormatter.ID
			}
			if formatterID == 0 {
				fmt.Println("No Formatter found")
				return err
			}

			var dataFormatterBridgeds model.FormatterBridgesEntityModel

			criteriaFB := model.FormatterBridgesFilterModel{}
			criteriaFB.Source = &tmpStr
			criteriaFB.FormatterID = &formatterID
			criteriaFB.TrxRefID = &payload.IDMutasiDta

			temptMDTA := "MUTASI-DTA"
			dataFormatterBridgeds.Context = ctx
			dataFormatterBridgeds.FormatterBridgesEntity = model.FormatterBridgesEntity{
				TrxRefID:    payload.IDMutasiDta,
				FormatterID: getFormatterID.ID,
				Source:      temptMDTA,
			}
			resultFormatterBridge, err := s.FormatterBridgesRepository.Create(ctx, &dataFormatterBridgeds)
			if err != nil {
				return err
			}

			for _, v := range datas {
				dataTBD := model.MutasiDtaDetailEntityModel{
					Context:               ctx,
					MutasiDtaDetailEntity: v,
				}
				dataTBD.FormatterBridgesID = resultFormatterBridge.ID
				_, err := s.MutasiDtaDetailRepository.Create(ctx, &dataTBD)
				if err != nil {
					return err
				}
			}
			dataTB = append(dataTB)
		}
		var data model.ImportedWorksheetDetailEntityModel
		_, err = s.ImportedWorksheetDetailRepository.FindByID(ctx, &payload.IDWorksheetDetailMutasiDta)
		if err != nil {
			return err
		}
		data.Context = ctx
		data.ImportedWorksheetDetailEntity = model.ImportedWorksheetDetailEntity{
			Status: 2,
		}

		_, err = s.ImportedWorksheetDetailRepository.Update(ctx, &payload.IDWorksheetDetailMutasiDta, &data)
		if err != nil {
			return err
		}
		return nil
	}); err != nil {
		return &dto.MutasiDtaImportResponse{}, err
	}
	result := &dto.MutasiDtaImportResponse{}
	return result, nil
}
func (s *service) ImportMutasiIa(ctx *abstraction.Context, payload *abstraction.JsonDataReUpload, datas []model.MutasiIaDetailEntity) (*dto.MutasiIaImportResponse, error) {

	tmpFolder := payload.MutasiIa
	_, err := os.Stat(tmpFolder)
	if err != nil {
		if os.IsNotExist(err) {
			if err := os.MkdirAll(tmpFolder, os.ModePerm); err != nil {
				log.Fatal(err)
			}
		} else {
			panic(err)
		}
	}

	f, err := excelize.OpenFile(tmpFolder)
	sheet := f.GetSheetName(f.GetActiveSheetIndex())
	rows, err := f.GetRows(sheet)
	if err != nil {
		panic(err)
	}

	confRowStart := []string{"MUTASI-IA-COST-IMPORT-ROW-START", "MUTASI-IA-ACCUMULATED-DEPRECATION-IMPORT-ROW-START"}
	confRowEnd := []string{"MUTASI-IA-COST-IMPORT-ROW-END", "MUTASI-IA-ACCUMULATED-DEPRECATION-IMPORT-ROW-END"}
	allData := [][]model.MutasiIaDetailEntity{}
	counter := 1

	var dataTB []model.MutasiIaEntityModel
	currentYear, currentMonth, _ := time.Now().Date()
	currentLocation := time.Now().Location()
	firstOfMonth := time.Date(currentYear, currentMonth, 1, 0, 0, 0, 0, currentLocation)
	lastOfMonth := firstOfMonth.AddDate(0, 1, -1)
	period := lastOfMonth.Format("2006-01-02")

	if err := trxmanager.New(s.Db).WithTrx(ctx, func(ctx *abstraction.Context) error {
		c2, err := f.GetCellValue(sheet, "B2")
		if c2 != "Mutasi Intangible Assets (IA)" {
			return errors.New("caannot")
		}
		for i, conf := range confRowStart {
			confCode := conf
			importConfigStartRow, err := s.ParameterRepository.FindWithCode(ctx, &confCode)
			if err != nil || len(*importConfigStartRow) == 0 {
				panic(err)
			}

			startRow := 0

			for _, v := range *importConfigStartRow {
				startRow, err = strconv.Atoi(*v.Value)
				if err != nil {
					startRow = 0
				}
			}

			confCode = confRowEnd[i]
			importConfigEndRow, err := s.ParameterRepository.FindWithCode(ctx, &confCode)
			if err != nil || len(*importConfigEndRow) == 0 {
				panic(err)
			}

			endRow := 0

			for _, v := range *importConfigEndRow {
				endRow, err = strconv.Atoi(*v.Value)
				if err != nil {
					endRow = 0
				}
			}

			rowsData := rows[startRow:endRow]
			datas := []model.MutasiIaDetailEntity{}
			for _, row := range rowsData {
				if len(row) == 0 {
					continue
				}

				tmp := make([]float64, len(row))
				for i := 3; i < len(row); i++ {
					if len(row[i]) == 0 {
						tmp[i-3] = 0
						continue
					}
					tmpi, err := strconv.ParseFloat(row[i], 64)
					if err != nil {
						tmpi = 0
					}
					tmp[i-3] = tmpi
				}

				data := model.MutasiIaDetailEntity{
					Description:             row[1],
					Code:                    strings.Replace(strings.ToUpper(row[1]), " ", "_", -1),
					BeginningBalance:        &tmp[0],
					AcquisitionOfSubsidiary: &tmp[1],
					Additions:               &tmp[2],
					Deductions:              &tmp[3],
					Reclassification:        &tmp[4],
					Revaluation:             &tmp[5],
					EndingBalance:           &tmp[6],
					Control:                 &tmp[7],
					SortId:                  counter,
				}
				datas = append(datas, data)
				counter++
			}
			allData = append(allData, datas)
		}
		//cek company berdasarkan user
		//belum ada
		//skip
		companyID := 1
		criteriaTB := model.MutasiIaFilterModel{}
		criteriaTB.Period = &period
		criteriaTB.CompanyID = &companyID

		for i, datas := range allData {
			criteriaFormatter := model.FormatterFilterModel{}
			tmpConf := strings.Split(confRowStart[i], "-IMPORT-")
			tmpStr := tmpConf[0]
			criteriaFormatter.FormatterFor = &tmpStr

			getFormatter, err := s.FormatterRepository.Find(ctx, &criteriaFormatter)
			if err != nil {
				return err
			}
			getFormatterID, err := s.FormatterRepository.FindWithDetail(ctx, &criteriaFormatter)
			if err != nil {
				return err
			}
			formatterID := 0
			for _, tmpFormatter := range *getFormatter {
				formatterID = tmpFormatter.ID
			}
			if formatterID == 0 {
				fmt.Println("No Formatter found")
				return err
			}

			var dataFormatterBridgeds model.FormatterBridgesEntityModel
			temptMIA := "MUTASI-IA"
			criteriaFB := model.FormatterBridgesFilterModel{}
			criteriaFB.Source = &tmpStr
			criteriaFB.FormatterID = &formatterID
			criteriaFB.TrxRefID = &payload.IDMutasiIa

			dataFormatterBridgeds.Context = ctx
			dataFormatterBridgeds.FormatterBridgesEntity = model.FormatterBridgesEntity{
				TrxRefID:    payload.IDMutasiIa,
				FormatterID: getFormatterID.ID,
				Source:      temptMIA,
			}
			resultFormatterBridge, err := s.FormatterBridgesRepository.Create(ctx, &dataFormatterBridgeds)
			if err != nil {
				return err
			}
			for _, v := range datas {
				dataTBD := model.MutasiIaDetailEntityModel{
					Context:              ctx,
					MutasiIaDetailEntity: v,
				}
				dataTBD.FormatterBridgesID = resultFormatterBridge.ID
				_, err := s.MutasiIaDetailRepository.Create(ctx, &dataTBD)
				if err != nil {
					return err
				}
			}
			dataTB = append(dataTB)
		}
		var data model.ImportedWorksheetDetailEntityModel
		_, err = s.ImportedWorksheetDetailRepository.FindByID(ctx, &payload.IDWorksheetDetailMutasiIa)
		if err != nil {
			return err
		}
		data.Context = ctx
		data.ImportedWorksheetDetailEntity = model.ImportedWorksheetDetailEntity{
			Status: 2,
		}

		_, err = s.ImportedWorksheetDetailRepository.Update(ctx, &payload.IDWorksheetDetailMutasiIa, &data)
		if err != nil {
			return err
		}
		return nil
	}); err != nil {
		return &dto.MutasiIaImportResponse{}, err
	}
	result := &dto.MutasiIaImportResponse{}
	return result, nil
}
func (s *service) ImportMutasiRua(ctx *abstraction.Context, payload *abstraction.JsonDataReUpload, datas []model.MutasiRuaDetailEntity) (*dto.MutasiRuaImportResponse, error) {

	tmpFolder := payload.MutasiRua
	_, err := os.Stat(tmpFolder)
	if err != nil {
		if os.IsNotExist(err) {
			if err := os.MkdirAll(tmpFolder, os.ModePerm); err != nil {
				log.Fatal(err)
			}
		} else {
			panic(err)
		}
	}

	f, err := excelize.OpenFile(tmpFolder)
	sheet := f.GetSheetName(f.GetActiveSheetIndex())
	rows, err := f.GetRows(sheet)
	if err != nil {
		panic(err)
	}
	confRowStart := []string{"MUTASI-RUA-COST-IMPORT-ROW-START", "MUTASI-RUA-ACCUMULATED-DEPRECATION-IMPORT-ROW-START"}
	confRowEnd := []string{"MUTASI-RUA-COST-IMPORT-ROW-END", "MUTASI-RUA-ACCUMULATED-DEPRECATION-IMPORT-ROW-END"}
	allData := [][]model.MutasiRuaDetailEntity{}
	counter := 1

	var dataTB []model.MutasiRuaEntityModel
	currentYear, currentMonth, _ := time.Now().Date()
	currentLocation := time.Now().Location()
	firstOfMonth := time.Date(currentYear, currentMonth, 1, 0, 0, 0, 0, currentLocation)
	lastOfMonth := firstOfMonth.AddDate(0, 1, -1)
	period := lastOfMonth.Format("2006-01-02")

	if err := trxmanager.New(s.Db).WithTrx(ctx, func(ctx *abstraction.Context) error {
		c2, err := f.GetCellValue(sheet, "B2")
		if c2 != "Mutasi Right of Used Assets (RUA)" {
			return errors.New("caannot")
		}
		for i, conf := range confRowStart {
			confCode := conf
			importConfigStartRow, err := s.ParameterRepository.FindWithCode(ctx, &confCode)
			if err != nil || len(*importConfigStartRow) == 0 {
				panic(err)
			}

			startRow := 0

			for _, v := range *importConfigStartRow {
				startRow, err = strconv.Atoi(*v.Value)
				if err != nil {
					startRow = 0
				}
			}

			confCode = confRowEnd[i]
			importConfigEndRow, err := s.ParameterRepository.FindWithCode(ctx, &confCode)
			if err != nil || len(*importConfigEndRow) == 0 {
				panic(err)
			}

			endRow := 0

			for _, v := range *importConfigEndRow {
				endRow, err = strconv.Atoi(*v.Value)
				if err != nil {
					endRow = 0
				}
			}

			rowsData := rows[startRow-1 : endRow]
			datas := []model.MutasiRuaDetailEntity{}
			for _, row := range rowsData {
				if len(row) == 0 {
					continue
				}

				tmp := make([]float64, len(row))
				for i := 3; i < len(row); i++ {
					if len(row[i]) == 0 {
						tmp[i-3] = 0
						continue
					}
					tmpi, err := strconv.ParseFloat(row[i], 64)
					if err != nil {
						tmpi = 0
					}
					tmp[i-3] = tmpi
				}

				data := model.MutasiRuaDetailEntity{
					Description:             row[1],
					Code:                    strings.Replace(strings.ToUpper(row[1]), " ", "_", -1),
					BeginningBalance:        &tmp[0],
					AcquisitionOfSubsidiary: &tmp[1],
					Additions:               &tmp[2],
					Deductions:              &tmp[3],
					Reclassification:        &tmp[4],
					Remeasurement:           &tmp[5],
					EndingBalance:           &tmp[6],
					Control:                 &tmp[7],
					SortId:                  counter,
				}
				datas = append(datas, data)
				counter++
			}
			allData = append(allData, datas)
		}
		//cek company berdasarkan user
		//belum ada
		//skip
		companyID := 1
		criteriaTB := model.MutasiRuaFilterModel{}
		criteriaTB.Period = &period
		criteriaTB.CompanyID = &companyID

		if err != nil {
			return err
		}

		for i, datas := range allData {
			criteriaFormatter := model.FormatterFilterModel{}
			tmpConf := strings.Split(confRowStart[i], "-IMPORT-")
			tmpStr := tmpConf[0]
			criteriaFormatter.FormatterFor = &tmpStr

			getFormatter, err := s.FormatterRepository.Find(ctx, &criteriaFormatter)
			if err != nil {
				return err
			}
			getFormatterID, err := s.FormatterRepository.FindWithDetail(ctx, &criteriaFormatter)
			if err != nil {
				return err
			}
			formatterID := 0
			for _, tmpFormatter := range *getFormatter {
				formatterID = tmpFormatter.ID
			}
			if formatterID == 0 {
				fmt.Println("No Formatter found")
				return err
			}

			var dataFormatterBridgeds model.FormatterBridgesEntityModel

			criteriaFB := model.FormatterBridgesFilterModel{}
			criteriaFB.Source = &tmpStr
			criteriaFB.FormatterID = &formatterID
			criteriaFB.TrxRefID = &payload.IDMutasiRua
			temptRUA := "MUTASI-RUA"
			dataFormatterBridgeds.Context = ctx
			dataFormatterBridgeds.FormatterBridgesEntity = model.FormatterBridgesEntity{
				TrxRefID:    payload.IDMutasiRua,
				FormatterID: getFormatterID.ID,
				Source:      temptRUA,
			}
			resultFormatterBridge, err := s.FormatterBridgesRepository.Create(ctx, &dataFormatterBridgeds)
			if err != nil {
				return err
			}

			for _, v := range datas {
				dataTBD := model.MutasiRuaDetailEntityModel{
					Context:               ctx,
					MutasiRuaDetailEntity: v,
				}
				dataTBD.FormatterBridgesID = resultFormatterBridge.ID
				_, err := s.MutasiRuaDetailRepository.Create(ctx, &dataTBD)
				if err != nil {
					return err
				}
			}
			dataTB = append(dataTB)
		}
		var data model.ImportedWorksheetDetailEntityModel
		_, err = s.ImportedWorksheetDetailRepository.FindByID(ctx, &payload.IDWorksheetDetailMutasiRua)
		if err != nil {
			return err
		}
		data.Context = ctx
		data.ImportedWorksheetDetailEntity = model.ImportedWorksheetDetailEntity{
			Status: 2,
		}

		_, err = s.ImportedWorksheetDetailRepository.Update(ctx, &payload.IDWorksheetDetailMutasiRua, &data)
		if err != nil {
			return err
		}
		return nil
	}); err != nil {
		return &dto.MutasiRuaImportResponse{}, err
	}
	result := &dto.MutasiRuaImportResponse{}
	return result, nil
}
func (s *service) ImportMutasiPersediaan(ctx *abstraction.Context, payload *abstraction.JsonDataReUpload, datas []model.MutasiPersediaanDetailEntity) (*dto.MutasiPersediaanImportResponse, error) {

	tmpFolder := payload.MutasiPersediaan
	_, err := os.Stat(tmpFolder)
	if err != nil {
		if os.IsNotExist(err) {
			if err := os.MkdirAll(tmpFolder, os.ModePerm); err != nil {
				log.Fatal(err)
			}
		} else {
			panic(err)
		}
	}

	f, err := excelize.OpenFile(tmpFolder)
	sheet := f.GetSheetName(f.GetActiveSheetIndex())
	rows, err := f.GetRows(sheet)
	if err != nil {
		panic(err)
	}

	confRowStart := []string{"MUTASI-PERSEDIAAN-IMPORT-ROW-START", "MUTASI-CADANGAN-PENGHAPUSAN-PERSEDIAAN-IMPORT-ROW-START"}
	confRowEnd := []string{"MUTASI-PERSEDIAAN-IMPORT-ROW-END", "MUTASI-CADANGAN-PENGHAPUSAN-PERSEDIAAN-IMPORT-ROW-END"}
	allData := [][]model.MutasiPersediaanDetailEntity{}
	counter := 1

	var dataTB []model.MutasiPersediaanEntityModel
	currentYear, currentMonth, _ := time.Now().Date()
	currentLocation := time.Now().Location()
	firstOfMonth := time.Date(currentYear, currentMonth, 1, 0, 0, 0, 0, currentLocation)
	lastOfMonth := firstOfMonth.AddDate(0, 1, -1)
	period := lastOfMonth.Format("2006-01-02")

	if err := trxmanager.New(s.Db).WithTrx(ctx, func(ctx *abstraction.Context) error {
		c2, err := f.GetCellValue(sheet, "B2")
		if c2 != "Mutasi Persediaan" {
			return errors.New("caannot")
		}
		for i, conf := range confRowStart {
			confCode := conf
			importConfigStartRow, err := s.ParameterRepository.FindWithCode(ctx, &confCode)
			if err != nil || len(*importConfigStartRow) == 0 {
				panic(err)
			}

			startRow := 0

			for _, v := range *importConfigStartRow {
				startRow, err = strconv.Atoi(*v.Value)
				if err != nil {
					startRow = 0
				}
			}

			confCode = confRowEnd[i]
			importConfigEndRow, err := s.ParameterRepository.FindWithCode(ctx, &confCode)
			if err != nil || len(*importConfigEndRow) == 0 {
				panic(err)
			}

			endRow := 0

			for _, v := range *importConfigEndRow {
				endRow, err = strconv.Atoi(*v.Value)
				if err != nil {
					endRow = 0
				}
			}

			rowsData := rows[startRow-1 : endRow]
			datas := []model.MutasiPersediaanDetailEntity{}
			for _, row := range rowsData {
				if len(row) == 0 {
					continue
				}

				tmp, err := strconv.ParseFloat(row[2], 64)
				if err != nil {
					tmp = 0
				}

				data := model.MutasiPersediaanDetailEntity{
					Description: row[1],
					Code:        strings.Replace(strings.ToUpper(row[1]), " ", "_", -1),
					Amount:      &tmp,
					SortID:      counter,
				}
				datas = append(datas, data)
				counter++
			}
			allData = append(allData, datas)
		}

		//cek company berdasarkan user
		//belum ada
		//skip
		companyID := 1
		criteriaTB := model.MutasiFaFilterModel{}
		criteriaTB.Period = &period
		criteriaTB.CompanyID = &companyID

		for i, datas := range allData {
			criteriaFormatter := model.FormatterFilterModel{}
			tmpConf := strings.Split(confRowStart[i], "-IMPORT-")
			tmpStr := tmpConf[0]
			criteriaFormatter.FormatterFor = &tmpStr

			getFormatter, err := s.FormatterRepository.Find(ctx, &criteriaFormatter)
			if err != nil {
				errs = append(errs, err)
				log.Fatalln(err)
				return err
			}
			getFormatterID, err := s.FormatterRepository.FindWithDetail(ctx, &criteriaFormatter)
			if err != nil {
				return err
			}
			formatterID := 0
			for _, tmpFormatter := range *getFormatter {
				formatterID = tmpFormatter.ID
			}
			if formatterID == 0 {
				fmt.Println("No Formatter found")
				errs = append(errs, err)
				log.Fatalln(err)
				return err
			}

			var dataFormatterBridgeds model.FormatterBridgesEntityModel

			criteriaFB := model.FormatterBridgesFilterModel{}
			criteriaFB.Source = &tmpStr
			criteriaFB.FormatterID = &formatterID
			criteriaFB.TrxRefID = &payload.IDMutasiPersediaan
			temptP := "MUTASI-PERSEDIAAN"
			dataFormatterBridgeds.Context = ctx
			dataFormatterBridgeds.FormatterBridgesEntity = model.FormatterBridgesEntity{
				TrxRefID:    payload.IDMutasiPersediaan,
				FormatterID: getFormatterID.ID,
				Source:      temptP,
			}
			resultFormatterBridge, err := s.FormatterBridgesRepository.Create(ctx, &dataFormatterBridgeds)
			if err != nil {
				return err
			}
			for _, v := range datas {
				dataTBD := model.MutasiPersediaanDetailEntityModel{
					Context:                      ctx,
					MutasiPersediaanDetailEntity: v,
				}
				dataTBD.FormatterBridgesID = resultFormatterBridge.ID
				_, err := s.MutasiPersediaanDetailRepository.Create(ctx, &dataTBD)
				if err != nil {
					errs = append(errs, err)
					log.Fatalln(err)
					return err
				}
			}
			dataTB = append(dataTB)
		}
		var data model.ImportedWorksheetDetailEntityModel
		_, err = s.ImportedWorksheetDetailRepository.FindByID(ctx, &payload.IDWorksheetDetailMutasiPersediaan)
		if err != nil {
			return err
		}
		data.Context = ctx
		data.ImportedWorksheetDetailEntity = model.ImportedWorksheetDetailEntity{
			Status: 2,
		}

		_, err = s.ImportedWorksheetDetailRepository.Update(ctx, &payload.IDWorksheetDetailMutasiPersediaan, &data)
		if err != nil {
			return err
		}
		return nil
	}); err != nil {
		return &dto.MutasiPersediaanImportResponse{}, err
	}
	result := &dto.MutasiPersediaanImportResponse{}
	return result, nil
}
func (s *service) ImportEmployeeBenefit(ctx *abstraction.Context, payload *abstraction.JsonDataReUpload, datas []model.EmployeeBenefitDetailEntity) (*dto.EmployeeBenefitImportResponse, error) {

	tmpFolder := payload.EmployeeBenefit
	_, err := os.Stat(tmpFolder)
	if err != nil {
		if os.IsNotExist(err) {
			if err := os.MkdirAll(tmpFolder, os.ModePerm); err != nil {
				log.Fatal(err)
			}
		} else {
			panic(err)
		}
	}

	f, err := excelize.OpenFile(tmpFolder)
	sheet := f.GetSheetName(f.GetActiveSheetIndex())
	rows, err := f.GetRows(sheet)
	if err != nil {
		panic(err)
	}

	confRowStart := []string{"EMPLOYEE-BENEFIT-ASUMSI-IMPORT-ROW-START", "EMPLOYEE-BENEFIT-REKONSILIASI-IMPORT-ROW-START", "EMPLOYEE-BENEFIT-RINCIAN-LAPORAN-IMPORT-ROW-START", "EMPLOYEE-BENEFIT-RINCIAN-EKUITAS-IMPORT-ROW-START", "EMPLOYEE-BENEFIT-MUTASI-IMPORT-ROW-START", "EMPLOYEE-BENEFIT-INFORMASI-IMPORT-ROW-START", "EMPLOYEE-BENEFIT-ANALISIS-IMPORT-ROW-START"}
	confRowEnd := []string{"EMPLOYEE-BENEFIT-ASUMSI-IMPORT-ROW-END", "EMPLOYEE-BENEFIT-REKONSILIASI-IMPORT-ROW-END", "EMPLOYEE-BENEFIT-RINCIAN-LAPORAN-IMPORT-ROW-END", "EMPLOYEE-BENEFIT-RINCIAN-EKUITAS-IMPORT-ROW-END", "EMPLOYEE-BENEFIT-MUTASI-IMPORT-ROW-END", "EMPLOYEE-BENEFIT-INFORMASI-IMPORT-ROW-END", "EMPLOYEE-BENEFIT-ANALISIS-IMPORT-ROW-END"}
	allData := [][]model.EmployeeBenefitDetailEntity{}
	counter := 1

	var dataTB []model.EmployeeBenefitEntityModel
	currentYear, currentMonth, _ := time.Now().Date()
	currentLocation := time.Now().Location()
	firstOfMonth := time.Date(currentYear, currentMonth, 1, 0, 0, 0, 0, currentLocation)
	lastOfMonth := firstOfMonth.AddDate(0, 1, -1)
	period := lastOfMonth.Format("2006-01-02")

	if err := trxmanager.New(s.Db).WithTrx(ctx, func(ctx *abstraction.Context) error {
		c2, err := f.GetCellValue(sheet, "A2")
		if c2 != "MUTASI LIABILITAS IMBALAN KERJA" {
			return errors.New("caannot")
		}
		for i, conf := range confRowStart {
			confCode := conf
			importConfigStartRow, err := s.ParameterRepository.FindWithCode(ctx, &confCode)
			if err != nil || len(*importConfigStartRow) == 0 {
				panic(err)
			}

			startRow := 0

			for _, v := range *importConfigStartRow {
				startRow, err = strconv.Atoi(*v.Value)
				if err != nil {
					startRow = 0
				}
			}

			confCode = confRowEnd[i]
			importConfigEndRow, err := s.ParameterRepository.FindWithCode(ctx, &confCode)
			if err != nil || len(*importConfigEndRow) == 0 {
				panic(err)
			}

			endRow := 0

			for _, v := range *importConfigEndRow {
				endRow, err = strconv.Atoi(*v.Value)
				if err != nil {
					endRow = 0
				}
			}

			rowsData := rows[startRow-1 : endRow]
			datas := []model.EmployeeBenefitDetailEntity{}
			for _, row := range rowsData {
				if len(row) == 0 {
					continue
				}

				tmp := row[6]
				regex, _ := regexp.Compile(`[a-z]+`)

				isMatch := regex.MatchString(tmp)
				coass := strings.Replace(strings.ToUpper(row[1]), "+", "#", -1)
				coas := strings.Replace(strings.ToUpper(coass), " ", "_", -1)
				coa := strings.Replace(strings.ToUpper(coas), "-", "~", -1)
				if isMatch == true {
					data := model.EmployeeBenefitDetailEntity{
						Description: row[1],
						Code:        coa,
						SortID:      counter,
						Value:       row[6],
						IsValue:     true,
					}
					datas = append(datas, data)
					counter++
				}
				if isMatch == false {
					_, err := strconv.ParseFloat(row[6], 64)
					if err != nil {
						data := model.EmployeeBenefitDetailEntity{
							Description: row[1],
							Code:        coa,
							SortID:      counter,
							Value:       row[6],
							IsValue:     true,
						}
						datas = append(datas, data)
						counter++
					}

				}

				if isMatch == false {
					tmp, err := strconv.ParseFloat(row[6], 64)
					if err != nil {
						continue
					}
					data := model.EmployeeBenefitDetailEntity{
						Description: row[1],
						Code:        coa,
						SortID:      counter,
						Amount:      &tmp,
						IsValue:     false,
					}
					datas = append(datas, data)
					counter++
				}
			}

			allData = append(allData, datas)
		}

		//cek company berdasarkan user
		//belum ada
		//skip

		criteriaTB := model.EmployeeBenefitFilterModel{}
		criteriaTB.Period = &period

		if err != nil {
			errs = append(errs, err)
			log.Fatalln(err)
			return err
		}
		for i, datas := range allData {
			criteriaFormatter := model.FormatterFilterModel{}
			tmpConf := strings.Split(confRowStart[i], "-IMPORT-")
			tmpStr := tmpConf[0]
			criteriaFormatter.FormatterFor = &tmpStr

			getFormatter, err := s.FormatterRepository.Find(ctx, &criteriaFormatter)
			if err != nil {
				errs = append(errs, err)
				log.Fatalln(err)
				return err
			}
			getFormatterID, err := s.FormatterRepository.FindWithDetail(ctx, &criteriaFormatter)
			if err != nil {
				return err
			}
			formatterID := 0
			for _, tmpFormatter := range *getFormatter {
				formatterID = tmpFormatter.ID
			}
			if formatterID == 0 {
				fmt.Println("No Formatter found")
				errs = append(errs, err)
				log.Fatalln(err)
				return err
			}

			var dataFormatterBridgeds model.FormatterBridgesEntityModel

			criteriaFB := model.FormatterBridgesFilterModel{}
			criteriaFB.Source = &tmpStr
			criteriaFB.FormatterID = &formatterID
			criteriaFB.TrxRefID = &payload.IDEmployeeBenefit
			temptP := "EMPLOYEE-BENEFIT"
			dataFormatterBridgeds.Context = ctx
			dataFormatterBridgeds.FormatterBridgesEntity = model.FormatterBridgesEntity{
				TrxRefID:    payload.IDEmployeeBenefit,
				FormatterID: getFormatterID.ID,
				Source:      temptP,
			}
			trxRefId := &payload.IDEmployeeBenefit
			resultFormatterBridge, err := s.FormatterBridgesRepository.FindByIDTrx(ctx, &temptP, trxRefId)
			if err != nil {
				return err
			}
			for _, v := range datas {
				dataTBD := model.EmployeeBenefitDetailEntityModel{
					Context:                     ctx,
					EmployeeBenefitDetailEntity: v,
				}
				dataTBD.FormatterBridgesID = resultFormatterBridge.ID
				_, err := s.EmployeeBenefitDetailRepository.Create(ctx, &dataTBD)
				if err != nil {
					errs = append(errs, err)
					log.Fatalln(err)
					return err
				}
			}
			dataTB = append(dataTB)
		}
		var data model.ImportedWorksheetDetailEntityModel
		_, err = s.ImportedWorksheetDetailRepository.FindByID(ctx, &payload.IDWorksheetDetailEmployeeBenefit)
		if err != nil {
			return err
		}
		data.Context = ctx
		data.ImportedWorksheetDetailEntity = model.ImportedWorksheetDetailEntity{
			Status: 2,
		}

		_, err = s.ImportedWorksheetDetailRepository.Update(ctx, &payload.IDWorksheetDetailEmployeeBenefit, &data)
		if err != nil {
			return err
		}
		return nil
	}); err != nil {
		return &dto.EmployeeBenefitImportResponse{}, err
	}
	result := &dto.EmployeeBenefitImportResponse{}
	return result, nil
}

// non formatter
func (s *service) ImportInvestasiNonTbk(ctx *abstraction.Context, payload *abstraction.JsonDataReUpload, datas []model.InvestasiNonTbkDetailEntity) (*dto.InvestasiNonTbkImportResponse, error) {

	tmpFolder := payload.InvestasiNonTbk
	_, err := os.Stat(tmpFolder)
	if err != nil {
		if os.IsNotExist(err) {
			if err := os.MkdirAll(tmpFolder, os.ModePerm); err != nil {
				log.Fatal(err)
			}
		} else {
			panic(err)
		}
	}

	f, err := excelize.OpenFile(tmpFolder)
	sheet := f.GetSheetName(f.GetActiveSheetIndex())
	rows, err := f.GetRows(sheet)
	if err != nil {
		panic(err)
	}

	confRowStart := []string{"INVESTASI-NON-TBK-IMPORT-ROW-START"}
	confRowEnd := []string{"INVESTASI-NON-TBK-IMPORT-ROW-END"}
	allData := [][]model.InvestasiNonTbkDetailEntity{}
	counter := 1

	if err := trxmanager.New(s.Db).WithTrx(ctx, func(ctx *abstraction.Context) error {
		c2, err := f.GetCellValue(sheet, "B2")
		if err != nil {
			return err
		}
		if c2 != "Detail investasi anak usaha Non TBK" {
			return errors.New("caannot")
		}
		for i, conf := range confRowStart {
			confCode := conf
			importConfigStartRow, err := s.ParameterRepository.FindWithCode(ctx, &confCode)
			if err != nil || len(*importConfigStartRow) == 0 {
				panic(err)
			}

			startRow := 0

			for _, v := range *importConfigStartRow {
				startRow, err = strconv.Atoi(*v.Value)
				if err != nil {
					startRow = 0
				}
			}

			confCode = confRowEnd[i]
			importConfigEndRow, err := s.ParameterRepository.FindWithCode(ctx, &confCode)
			if err != nil {
				panic(err)
			}

			endRow := 0

			for _, v := range *importConfigEndRow {
				endRow, err = strconv.Atoi(*v.Value)
				if err != nil {
					endRow = 0
				}
			}

			rowsData := rows[startRow-1:]
			if endRow != 0 {
				rowsData = rows[startRow-1 : endRow]
			}

			datas := []model.InvestasiNonTbkDetailEntity{}
			for _, row := range rowsData {
				if len(row) == 0 {
					continue
				}

				tmp := make([]float64, len(row))
				for i := 3; i < len(row); i++ {
					if len(row[i]) == 0 {
						tmp[i-3] = 0
						continue
					}
					tmpi, err := strconv.ParseFloat(row[i], 64)
					if err != nil {
						tmpi = 0
					}
					tmp[i-3] = tmpi
				}

				data := model.InvestasiNonTbkDetailEntity{
					Code:                row[2],
					LbrSahamOwnership:   &tmp[0],
					TotalLbrSaham:       &tmp[1],
					PercentageOwnership: &tmp[2],
					HargaPar:            &tmp[3],
					TotalHargaPar:       &tmp[4],
					HargaBeli:           &tmp[5],
					TotalHargaBeli:      &tmp[6],
					SortId:              counter,
				}
				datas = append(datas, data)
				counter++
			}
			allData = append(allData, datas)
		}

		var dataTB []model.InvestasiNonTbkEntityModel
		currentYear, currentMonth, _ := time.Now().Date()
		currentLocation := time.Now().Location()
		firstOfMonth := time.Date(currentYear, currentMonth, 1, 0, 0, 0, 0, currentLocation)
		lastOfMonth := firstOfMonth.AddDate(0, 1, -1)
		period := lastOfMonth.Format("2006-01-02")
		//cek company berdasarkan user
		//belum ada
		//skip
		companyID := 1
		criteriaTB := model.InvestasiNonTbkFilterModel{}
		criteriaTB.Period = &period
		criteriaTB.CompanyID = &companyID

		for _, datas := range allData {

			for _, v := range datas {
				dataTBD := model.InvestasiNonTbkDetailEntityModel{
					Context:                     ctx,
					InvestasiNonTbkDetailEntity: v,
				}
				dataTBD.InvestasiNonTbkID = payload.IDInvestasiNonTbk
				_, err := s.InvestasiNonTbkDetailRepository.Create(ctx, &dataTBD)
				if err != nil {
					return err
				}
			}
			dataTB = append(dataTB)
		}
		var data model.ImportedWorksheetDetailEntityModel
		_, err = s.ImportedWorksheetDetailRepository.FindByID(ctx, &payload.IDWorksheetDetailInvestasiNonTbk)
		if err != nil {
			return err
		}
		data.Context = ctx
		data.ImportedWorksheetDetailEntity = model.ImportedWorksheetDetailEntity{
			Status: 2,
		}

		_, err = s.ImportedWorksheetDetailRepository.Update(ctx, &payload.IDWorksheetDetailInvestasiNonTbk, &data)
		if err != nil {
			return err
		}
		return nil
	}); err != nil {
		return &dto.InvestasiNonTbkImportResponse{}, err
	}
	result := &dto.InvestasiNonTbkImportResponse{}
	return result, nil
}
func (s *service) ImportPembelianPenjualanBerelasi(ctx *abstraction.Context, payload *abstraction.JsonDataReUpload, datas []model.PembelianPenjualanBerelasiDetailEntity) (*dto.PembelianPenjualanBerelasiImportResponse, error) {
	tmpFolder := payload.PembelianPenjualanBerelasi
	_, err := os.Stat(tmpFolder)
	if err != nil {
		if os.IsNotExist(err) {
			if err := os.MkdirAll(tmpFolder, os.ModePerm); err != nil {
				log.Fatal(err)
			}
		} else {
			panic(err)
		}
	}

	f, err := excelize.OpenFile(tmpFolder)
	sheet := f.GetSheetName(f.GetActiveSheetIndex())
	rows, err := f.GetRows(sheet)
	if err != nil {
		panic(err)
	}

	var dataTB []model.PembelianPenjualanBerelasiEntityModel
	currentYear, currentMonth, _ := time.Now().Date()
	currentLocation := time.Now().Location()
	firstOfMonth := time.Date(currentYear, currentMonth, 1, 0, 0, 0, 0, currentLocation)
	lastOfMonth := firstOfMonth.AddDate(0, 1, -1)
	period := lastOfMonth.Format("2006-01-02")

	if err := trxmanager.New(s.Db).WithTrx(ctx, func(ctx *abstraction.Context) error {
		c2, err := f.GetCellValue(sheet, "B2")
		if err != nil {
			return err
		}
		if c2 != "List pembelian dan penjualan berelasi" {
			return errors.New("caannot")
		}
		counter := 1
		confCode := "PEMBELIAN-PENJUALAN-BERELASI-IMPORT-ROW-START"
		importConfigStartRow, err := s.ParameterRepository.FindWithCode(ctx, &confCode)
		if err != nil || len(*importConfigStartRow) == 0 {
			panic(err)
		}

		startRow := 0

		for _, v := range *importConfigStartRow {
			startRow, err = strconv.Atoi(*v.Value)
			if err != nil {
				startRow = 0
			}
		}

		confCode = "PEMBELIAN-PENJUALAN-BERELASI-IMPORT-ROW-END"
		importConfigEndRow, err := s.ParameterRepository.FindWithCode(ctx, &confCode)
		if err != nil {
			panic(err)
		}

		endRow := 0

		for _, v := range *importConfigEndRow {
			endRow, err = strconv.Atoi(*v.Value)
			if err != nil {
				endRow = 0
			}
		}

		rowsData := rows[startRow-1:]
		if endRow != 0 {
			rowsData = rows[startRow-1 : endRow]
		}

		datas = []model.PembelianPenjualanBerelasiDetailEntity{}
		for _, row := range rowsData {
			if len(row) == 0 {
				continue
			}

			tmp, err := strconv.ParseFloat(row[4], 64)
			if err != nil {
				tmp = 0
			}

			tmp1, err := strconv.ParseFloat(row[5], 64)
			if err != nil {
				tmp1 = 0
			}

			data := model.PembelianPenjualanBerelasiDetailEntity{
				Name:         row[3],
				Code:         strings.Replace(strings.ToUpper(row[2]), " ", "_", -1),
				BoughtAmount: &tmp,
				SalesAmount:  &tmp1,
				SortID:       counter,
			}
			datas = append(datas, data)
			counter++
		}

		//cek company berdasarkan user
		//belum ada
		//skip
		companyID := 1
		criteriaTB := model.PembelianPenjualanBerelasiFilterModel{}
		criteriaTB.Period = &period
		criteriaTB.CompanyID = &companyID

		for _, v := range datas {
			dataTBD := model.PembelianPenjualanBerelasiDetailEntityModel{
				Context:                                ctx,
				PembelianPenjualanBerelasiDetailEntity: v,
			}
			dataTBD.PembelianPenjualanBerelasiID = payload.IDPembelianPenjualanBerelasi
			_, err := s.PPBerelasiDetailRepository.Create(ctx, &dataTBD)
			if err != nil {
				return err
			}
		}
		dataTB = append(dataTB)
		var data model.ImportedWorksheetDetailEntityModel
		_, err = s.ImportedWorksheetDetailRepository.FindByID(ctx, &payload.IDWorksheetDetailPembelianPenjualanBerelasi)
		if err != nil {
			return err
		}
		data.Context = ctx
		data.ImportedWorksheetDetailEntity = model.ImportedWorksheetDetailEntity{
			Status: 2,
		}

		_, err = s.ImportedWorksheetDetailRepository.Update(ctx, &payload.IDWorksheetDetailPembelianPenjualanBerelasi, &data)
		if err != nil {
			return err
		}
		return nil
	}); err != nil {
		return &dto.PembelianPenjualanBerelasiImportResponse{}, err
	}
	result := &dto.PembelianPenjualanBerelasiImportResponse{
		Data: dataTB,
	}

	return result, nil
}
