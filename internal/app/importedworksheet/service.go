package importedworksheet

import (
	"mcash-finance-console-core/internal/abstraction"
	"mcash-finance-console-core/internal/dto"
	"mcash-finance-console-core/internal/factory"
	"mcash-finance-console-core/internal/model"
	"net/http"
	"strconv"

	// "mcash-finance-console-core/internal/model"
	"mcash-finance-console-core/internal/repository"
	"mcash-finance-console-core/pkg/util/helper"
	"mcash-finance-console-core/pkg/util/response"
	"mcash-finance-console-core/pkg/util/trxmanager"

	"gorm.io/gorm"
)

type service struct {
	Repository                        repository.ImportedWorksheet
	ImportedWorksheetDetailRepository repository.ImportedWorksheetDetail
	Db                                *gorm.DB
	CompanyRepository                          repository.Company
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
	CoaRepository                              repository.Coa
}

type Service interface {
	Find(ctx *abstraction.Context, payload *dto.ImportedWorksheetGetRequest) (*dto.ImportedWorksheetGetResponse, error)
	FindByID(ctx *abstraction.Context, payload *dto.ImportedWorksheetGetByIDRequest) (*dto.ImportedWorksheetGetByIDResponse, error)
	GetVersion(ctx *abstraction.Context, payload *dto.GetVersionRequest) (*dto.GetVersionResponse, error)
}

func NewService(f *factory.Factory) *service {
	repository := f.ImportedWorksheetRepository
	ImportedWorksheetDetailRepository := f.ImportedWorksheetDetailRepository
	notifRepo := f.NotificationRepository
	importedWorksheetRepo := f.ImportedWorksheetRepository
	consolidationRepo := f.ConsolidationRepository
	agingUtangPiutangRepo := f.AgingUtangPiutangRepository
	agingUtangPiutangDetailRepo := f.AgingUtangPiutangDetailRepository
	formatterRepo := f.FormatterRepository
	formatterBridgesRepo := f.FormatterBridgesRepository
	employeeBenefitRepo := f.EmployeeBenefitRepository
	employeeBenefitDetailRepo := f.EmployeeBenefitDetailRepository
	adjustmentRepo := f.AdjustmentRepository
	adjustmentDetailRepo := f.AdjustmentDetailRepository
	investasiNonTbkRepo := f.InvestasiNonTbkRepository
	investasiNonTbkDetailRepo := f.InvestasiNonTbkDetailRepository
	investasiTbkRepo := f.InvestasiTbkRepository
	investasiTbkDetailRepo := f.InvestasiTbkDetailRepository
	mutasiDtaRepo := f.MutasiDtaRepository
	mutasiDtaDetailRepo := f.MutasiDtaDetailRepository
	mutasiFaRepo := f.MutasiFaRepository
	mutasiFaDetailRepo := f.MutasiFaDetailRepository
	mutasiIaRepo := f.MutasiIaRepository
	mutasiIaDetailRepo := f.MutasiIaDetailRepository
	mutasiPersediaanRepo := f.MutasiPersediaanRepository
	mutasiPersediaanDetailRepo := f.MutasiPersediaanDetailRepository
	mutasiRuaRepo := f.MutasiRuaRepository
	mutasiRuaDetailRepo := f.MutasiRuaDetailRepository
	pembelianPenjualanBerelasiRepo := f.PembelianPenjualanBerelasiRepository
	pembelianPenjualanBerelasiDetailRepo := f.PembelianPenjualanBerelasiDetailRepository
	trialBalanceRepo := f.TrialBalanceRepository
	trialBalanceDetailRepo := f.TrialBalanceDetailRepository
	formatterDetailRepo := f.FormatterDetailRepository
	companyRepo := f.CompanyRepository
	coaRepo := f.CoaRepository

	db := f.Db
	return &service{
		Repository:                        repository,
		ImportedWorksheetDetailRepository: ImportedWorksheetDetailRepository,
		NotificationRepository:                     notifRepo,
		ImportedWorksheetRepository:                importedWorksheetRepo,
		ConsolidationRepository:                    consolidationRepo,
		AgingUtangPiutangRepository:                agingUtangPiutangRepo,
		AgingUtangPiutangDetailRepository:          agingUtangPiutangDetailRepo,
		FormatterRepository:                        formatterRepo,
		FormatterBridgesRepository:                 formatterBridgesRepo,
		EmployeeBenefitRepository:                  employeeBenefitRepo,
		EmployeeBenefitDetailRepository:            employeeBenefitDetailRepo,
		AdjustmentRepository:                       adjustmentRepo,
		AdjustmentDetailRepository:                 adjustmentDetailRepo,
		InvestasiNonTbkRepository:                  investasiNonTbkRepo,
		InvestasiNonTbkDetailRepository:            investasiNonTbkDetailRepo,
		InvestasiTbkRepository:                     investasiTbkRepo,
		InvestasiTbkDetailRepository:               investasiTbkDetailRepo,
		MutasiDtaRepository:                        mutasiDtaRepo,
		MutasiDtaDetailRepository:                  mutasiDtaDetailRepo,
		MutasiFaRepository:                         mutasiFaRepo,
		MutasiFaDetailRepository:                   mutasiFaDetailRepo,
		MutasiIaRepository:                         mutasiIaRepo,
		MutasiIaDetailRepository:                   mutasiIaDetailRepo,
		MutasiPersediaanRepository:                 mutasiPersediaanRepo,
		MutasiPersediaanDetailRepository:           mutasiPersediaanDetailRepo,
		MutasiRuaRepository:                        mutasiRuaRepo,
		MutasiRuaDetailRepository:                  mutasiRuaDetailRepo,
		PembelianPenjualanBerelasiRepository:       pembelianPenjualanBerelasiRepo,
		PembelianPenjualanBerelasiDetailRepository: pembelianPenjualanBerelasiDetailRepo,
		TrialBalanceRepository:                     trialBalanceRepo,
		TrialBalanceDetailRepository:               trialBalanceDetailRepo,
		FormatterDetailRepository:                  formatterDetailRepo,
		CompanyRepository:                          companyRepo,
		CoaRepository:                              coaRepo,
		Db:                                db,
		
	}
}

func (s *service) Delete(ctx *abstraction.Context, payload *dto.ImportedWorksheetDeleteRequest) (*dto.ImportedWorksheetDeleteResponse, error) {
	var data model.ImportedWorksheetEntityModel
	if err := trxmanager.New(s.Db).WithTrx(ctx, func(ctx *abstraction.Context) error {
		
		imprtwrksht, err := s.Repository.FindByID(ctx, &payload.ID)
		if err != nil {
			return helper.ErrorHandler(err)
		}

		var tb model.TrialBalanceFilterModel
		tb.Period = &imprtwrksht.Period
		tb.Versions = &imprtwrksht.Versions
		tb.CompanyID = &imprtwrksht.CompanyID

		findTb, err := s.TrialBalanceRepository.FindByCriteria(ctx, &tb)
		if err != nil {
			return helper.ErrorHandler(err)
		}
		if findTb.Status == 2 || findTb.Status == 3 {
			return response.CustomErrorBuilder(http.StatusBadRequest, "Status Data Sudah Confirm / Console Tidak Bisa Di Hapus", "Status Data Sudah Confirm / Console Tidak Bisa Di Hapus")
		}
		deleted := 4
		var trialbalance model.TrialBalanceEntityModel
		trialbalance.Context = ctx
		trialbalance.TrialBalanceEntity = model.TrialBalanceEntity{
			Status: deleted,
		}
		_, err = s.Repository.DeleteTBDFNew(ctx, &findTb.ID, &trialbalance)
		if err != nil {
			return response.CustomErrorBuilder(http.StatusBadRequest, "deletedTb", "deletedTb")
		}

		findAJE, err := s.Repository.FindAjeWithTb(ctx, &findTb.ID)
		if err != nil {
			return response.CustomErrorBuilder(http.StatusBadRequest, "AJE", "AJE")
		}

		if len(*findAJE) > 0 {
			for _, v := range *findAJE {
				
				var aje model.AdjustmentEntityModel
				aje.Context = ctx
				aje.AdjustmentEntity = model.AdjustmentEntity{
					Status: deleted,
				}
				_, err = s.Repository.DeleteAjeWithTb(ctx, &v.ID, &aje)
				if err != nil {
					return response.CustomErrorBuilder(http.StatusBadRequest, "da", "da")
				}
			}
		}
		var ag model.AgingUtangPiutangFilterModel
		ag.Period = &imprtwrksht.Period
		ag.Versions = &imprtwrksht.Versions
		ag.CompanyID = &imprtwrksht.CompanyID

		findAg, err := s.AgingUtangPiutangRepository.FindByCriteria(ctx, &ag)
		if err != nil {
			return response.CustomErrorBuilder(http.StatusBadRequest, "ag", "ag")
		}
		var agingutangpiutang model.AgingUtangPiutangEntityModel
		agingutangpiutang.Context = ctx
		agingutangpiutang.AgingUtangPiutangEntity = model.AgingUtangPiutangEntity{
			Status: deleted,
		}
		_, err = s.Repository.DeleteAPDF(ctx, &findAg.ID, &agingutangpiutang)
		if err != nil {
			return response.CustomErrorBuilder(http.StatusBadRequest, "da", "da")
		}
		
		var fa model.MutasiFaFilterModel
		fa.Period = &imprtwrksht.Period
		fa.Versions = &imprtwrksht.Versions
		fa.CompanyID = &imprtwrksht.CompanyID

		findFa, err := s.MutasiFaRepository.FindByCriteria(ctx, &fa)
		if err != nil {
			return helper.ErrorHandler(err)
		}
		var mutasifa model.MutasiFaEntityModel
		mutasifa.Context = ctx
		mutasifa.MutasiFaEntity = model.MutasiFaEntity{
			Status: deleted,
		}
		_, err = s.Repository.DeleteMFDF(ctx, &findFa.ID, &mutasifa)
		if err != nil {
			return helper.ErrorHandler(err)
		}

		var rua model.MutasiRuaFilterModel
		rua.Period = &imprtwrksht.Period
		rua.Versions = &imprtwrksht.Versions
		rua.CompanyID = &imprtwrksht.CompanyID

		findRua, err := s.MutasiRuaRepository.FindByCriteria(ctx, &rua)
		if err != nil {
			return helper.ErrorHandler(err)
		}
		var mutasirua model.MutasiRuaEntityModel
		mutasirua.Context = ctx
		mutasirua.MutasiRuaEntity = model.MutasiRuaEntity{
			Status: deleted,
		}
		_, err = s.Repository.DeleteMRDF(ctx, &findRua.ID, &mutasirua)
		if err != nil {
			return helper.ErrorHandler(err)
		}

		var ia model.MutasiIaFilterModel
		ia.Period = &imprtwrksht.Period
		ia.Versions = &imprtwrksht.Versions
		ia.CompanyID = &imprtwrksht.CompanyID

		findia, err := s.MutasiIaRepository.FindByCriteria(ctx, &ia)
		if err != nil {
			return helper.ErrorHandler(err)
		}
		var mutasiia model.MutasiIaEntityModel
		mutasiia.Context = ctx
		mutasiia.MutasiIaEntity = model.MutasiIaEntity{
			Status: deleted,
		}
		_, err = s.Repository.DeleteMIDF(ctx, &findia.ID, &mutasiia)
		if err != nil {
			return helper.ErrorHandler(err)
		}

		var it model.InvestasiTbkFilterModel
		it.Period = &imprtwrksht.Period
		it.Versions = &imprtwrksht.Versions
		it.CompanyID = &imprtwrksht.CompanyID

		findit, err := s.InvestasiTbkRepository.FindByCriteria(ctx, &it)
		if err != nil {
			return helper.ErrorHandler(err)
		}
		var investasitbk model.InvestasiTbkEntityModel
		investasitbk.Context = ctx
		investasitbk.InvestasiTbkEntity = model.InvestasiTbkEntity{
			Status: deleted,
		}
		_, err = s.Repository.DeleteITDF(ctx, &findit.ID, &investasitbk)
		if err != nil {
			return helper.ErrorHandler(err)
		}
		data.Context = ctx
		result, err := s.Repository.DeleteImportedWorksheet(ctx, &payload.ID, &data)
		if err != nil {
			return helper.ErrorHandler(err)
		}

		var intb model.InvestasiNonTbkFilterModel
		intb.Period = &imprtwrksht.Period
		intb.Versions = &imprtwrksht.Versions
		intb.CompanyID = &imprtwrksht.CompanyID

		findintb, err := s.InvestasiNonTbkRepository.FindByCriteria(ctx, &intb)
		if err != nil {
			return helper.ErrorHandler(err)
		}
		var investasinontbk model.InvestasiNonTbkEntityModel
		investasinontbk.Context = ctx
		investasinontbk.InvestasiNonTbkEntity = model.InvestasiNonTbkEntity{
			Status: deleted,
		}
		_, err = s.Repository.DeleteINTDF(ctx, &findintb.ID, &investasinontbk)
		if err != nil {
			return helper.ErrorHandler(err)
		}

		var mp model.MutasiPersediaanFilterModel
		mp.Period = &imprtwrksht.Period
		mp.Versions = &imprtwrksht.Versions
		mp.CompanyID = &imprtwrksht.CompanyID

		findmp, err := s.MutasiPersediaanRepository.FindByCriteria(ctx, &mp)
		if err != nil {
			return helper.ErrorHandler(err)
		}
		var mutasipersediaan model.MutasiPersediaanEntityModel
		mutasipersediaan.Context = ctx
		mutasipersediaan.MutasiPersediaanEntity = model.MutasiPersediaanEntity{
			Status: deleted,
		}
		_, err = s.Repository.DeleteMPDF(ctx, &findmp.ID, &mutasipersediaan)
		if err != nil {
			return helper.ErrorHandler(err)
		}
		var em model.EmployeeBenefitFilterModel
		em.Period = &imprtwrksht.Period
		em.Versions = &imprtwrksht.Versions
		em.CompanyID = &imprtwrksht.CompanyID

		findem, err := s.EmployeeBenefitRepository.FindByCriteria(ctx, &em)
		if err != nil {
			return helper.ErrorHandler(err)
		}
		var employee model.EmployeeBenefitEntityModel
		employee.Context = ctx
		employee.EmployeeBenefitEntity = model.EmployeeBenefitEntity{
			Status: deleted,
		}
		_, err = s.Repository.DeleteEBDF(ctx, &findem.ID, &employee)
		if err != nil {
			return helper.ErrorHandler(err)
		}

		var pb model.PembelianPenjualanBerelasiFilterModel
		pb.Period = &imprtwrksht.Period
		pb.Versions = &imprtwrksht.Versions
		pb.CompanyID = &imprtwrksht.CompanyID

		findpb, err := s.PembelianPenjualanBerelasiRepository.FindByCriteria(ctx, &pb)
		if err != nil {
			return helper.ErrorHandler(err)
		}
		var pembelian model.PembelianPenjualanBerelasiEntityModel
		pembelian.Context = ctx
		pembelian.PembelianPenjualanBerelasiEntity = model.PembelianPenjualanBerelasiEntity{
			Status: deleted,
		}
		_, err = s.Repository.DeletePPBDF(ctx, &findpb.ID, &pembelian)
		if err != nil {
			return helper.ErrorHandler(err)
		}
		data.Context = ctx
		result, err = s.Repository.DeleteImportedWorksheet(ctx, &payload.ID, &data)
		if err != nil {
			return helper.ErrorHandler(err)
		}
		data = *result
		return nil
	}); err != nil {
		return &dto.ImportedWorksheetDeleteResponse{}, err
	}
	result := &dto.ImportedWorksheetDeleteResponse{
		// ImportedWorksheet: data,
	}
	return result, nil
}

func (s *service) Find(ctx *abstraction.Context, payload *dto.ImportedWorksheetGetRequest) (*dto.ImportedWorksheetGetResponse, error) {
	data, info, err := s.Repository.Find(ctx, &payload.ImportedWorksheetFilterModel, &payload.Pagination)
	if err != nil {
		return &dto.ImportedWorksheetGetResponse{}, response.ErrorBuilder(&response.ErrorConstant.InternalServerError, err)
	}
	var arrImportedWorksheet []model.ImportedWorksheetEntityModel
	for _, v := range *data {
		var dataImportedWorksheetDetail model.ImportedWorksheetDetailEntityModel
		dataImportedWorksheetDetail.Context = ctx

		dataImportedWorksheetDetail.ImportedWorksheetDetailEntity = model.ImportedWorksheetDetailEntity{
			ImportedWorksheetID: v.ID,
			Status:              1,
		}
		failed, err := s.ImportedWorksheetDetailRepository.GetCountStatus(ctx, &dataImportedWorksheetDetail)
		if err != nil {
			return &dto.ImportedWorksheetGetResponse{}, response.ErrorBuilder(&response.ErrorConstant.InternalServerError, err)
		}

		dataImportedWorksheetDetail.ImportedWorksheetDetailEntity = model.ImportedWorksheetDetailEntity{
			ImportedWorksheetID: v.ID,
			Status:              2,
		}
		succes, err := s.ImportedWorksheetDetailRepository.GetCountStatus(ctx, &dataImportedWorksheetDetail)
		if err != nil {
			return &dto.ImportedWorksheetGetResponse{}, response.ErrorBuilder(&response.ErrorConstant.InternalServerError, err)
		}
		var tb model.TrialBalanceFilterModel
		tb.Period = &v.Period
		tb.Versions = &v.Versions
		tb.CompanyID = &v.CompanyID

		findTb, err := s.TrialBalanceRepository.FindByCriteria(ctx, &tb)
		if err != nil {
			return &dto.ImportedWorksheetGetResponse{}, response.ErrorBuilder(&response.ErrorConstant.InternalServerError, err)
		}
		v.TrialBalance = *findTb
		v.Note = "SUCCES" + " " + strconv.Itoa(len(*succes)) + " " + "FAILED" + " " + strconv.Itoa(len(*failed))
		// v.Succes = len(*succes)
		// v.Failed = len(*failed)
		arrImportedWorksheet = append(arrImportedWorksheet, v)

		*data = arrImportedWorksheet
	}
	*data = arrImportedWorksheet
	result := &dto.ImportedWorksheetGetResponse{
		Datas:          *data,
		PaginationInfo: *info,
	}
	
	return result, nil
}

func (s *service) FindByID(ctx *abstraction.Context, payload *dto.ImportedWorksheetGetByIDRequest) (*dto.ImportedWorksheetGetByIDResponse, error) {
	data, err := s.Repository.FindByID(ctx, &payload.ID)
	if err != nil {
		return &dto.ImportedWorksheetGetByIDResponse{}, response.ErrorBuilder(&response.ErrorConstant.InternalServerError, err)
	}
	result := &dto.ImportedWorksheetGetByIDResponse{
		ImportedWorksheetEntityModel: *data,
	}
	return result, nil
}
