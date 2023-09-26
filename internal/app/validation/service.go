package validation

import (
	"encoding/json"
	"errors"
	"mcash-finance-console-core/internal/abstraction"
	"mcash-finance-console-core/internal/dto"
	"mcash-finance-console-core/internal/factory"
	"mcash-finance-console-core/internal/model"
	"mcash-finance-console-core/internal/repository"
	"mcash-finance-console-core/pkg/constant"
	"mcash-finance-console-core/pkg/kafka"
	"mcash-finance-console-core/pkg/util/helper"
	"mcash-finance-console-core/pkg/util/response"
	"mcash-finance-console-core/pkg/util/trxmanager"
	"net/http"
	"time"

	"gorm.io/gorm"
)

type service struct {
	FormatterRepository               repository.Formatter
	FormatterDetailRepository         repository.FormatterDetail
	TrialBalanceRepository            repository.TrialBalance
	TrialBalanceDetailRepository      repository.TrialBalanceDetail
	AgingUtangPiutangRepository       repository.AgingUtangPiutang
	AgingUtangPiutangDetailRepository repository.AgingUtangPiutangDetail
	MutasiPersediaanRepository        repository.MutasiPersediaan
	MutasiPersediaanDetailRepository  repository.MutasiPersediaanDetail
	MutasiFaRepository                repository.MutasiFa
	MutasiFaDetailRepository          repository.MutasiFaDetail
	MutasiIaRepository                repository.MutasiIa
	MutasiIaDetailRepository          repository.MutasiIaDetail
	MutasiRuaRepository               repository.MutasiRua
	MutasiRuaDetailRepository         repository.MutasiRuaDetail
	EmployeeBenefitRepository         repository.EmployeeBenefit
	EmployeeBenefitDetailRepository   repository.EmployeeBenefitDetail
	MutasiDtaRepository               repository.MutasiDta
	MutasiDtaDetailRepository         repository.MutasiDtaDetail
	ValidationRepository              repository.Validation
	ValidationDetailRepository        repository.ValidationDetail
	Db                                *gorm.DB
}

type Service interface {
	// Validation(ctx *abstraction.Context, payload *dto.UserDeleteRequest) (*dto.UserDeleteResponse, error)
	NewValidation(ctx *abstraction.Context) error
	Find(ctx *abstraction.Context, payload *dto.ValidationGetRequest) (*dto.ValidationGetResponse, error)
	FindByID(ctx *abstraction.Context, payload *dto.ValidationGetRequest) (*dto.ValidationGetResponse, error)
	ListAvailable(ctx *abstraction.Context, payload *dto.ValidationGetListAvailable) (*dto.ValidationGetResponse, error)
}

func NewService(f *factory.Factory) *service {
	formatterRepo := f.FormatterRepository
	formatterDetailRepo := f.FormatterDetailRepository
	trialBalanceRepo := f.TrialBalanceRepository
	trialBalanceDetailRepo := f.TrialBalanceDetailRepository
	agingUtangPiutangRepo := f.AgingUtangPiutangRepository
	agingUtangPiutangDetailRepo := f.AgingUtangPiutangDetailRepository
	mutasiPersediaanRepo := f.MutasiPersediaanRepository
	mutasiPersediaanDetailRepo := f.MutasiPersediaanDetailRepository
	mutasiFaRepo := f.MutasiFaRepository
	mutasiFaDetailRepo := f.MutasiFaDetailRepository
	mutasiIaRepo := f.MutasiIaRepository
	mutasiIaDetailRepo := f.MutasiIaDetailRepository
	mutasiRuaRepo := f.MutasiRuaRepository
	mutasiRuaDetailRepo := f.MutasiRuaDetailRepository
	employeeBenefitRepo := f.EmployeeBenefitRepository
	employeeBenefitDetailRepo := f.EmployeeBenefitDetailRepository
	mutasiDtaRepo := f.MutasiDtaRepository
	mutasiDtaDetailRepo := f.MutasiDtaDetailRepository
	validationRepo := f.ValidationRepository
	validationDetailRepo := f.ValidationDetailRepository
	db := f.Db
	return &service{
		FormatterRepository:               formatterRepo,
		FormatterDetailRepository:         formatterDetailRepo,
		TrialBalanceRepository:            trialBalanceRepo,
		TrialBalanceDetailRepository:      trialBalanceDetailRepo,
		AgingUtangPiutangRepository:       agingUtangPiutangRepo,
		AgingUtangPiutangDetailRepository: agingUtangPiutangDetailRepo,
		MutasiPersediaanRepository:        mutasiPersediaanRepo,
		MutasiPersediaanDetailRepository:  mutasiPersediaanDetailRepo,
		MutasiFaRepository:                mutasiFaRepo,
		MutasiFaDetailRepository:          mutasiFaDetailRepo,
		MutasiIaRepository:                mutasiIaRepo,
		MutasiIaDetailRepository:          mutasiIaDetailRepo,
		MutasiRuaRepository:               mutasiRuaRepo,
		MutasiRuaDetailRepository:         mutasiRuaDetailRepo,
		EmployeeBenefitRepository:         employeeBenefitRepo,
		EmployeeBenefitDetailRepository:   employeeBenefitDetailRepo,
		MutasiDtaRepository:               mutasiDtaRepo,
		MutasiDtaDetailRepository:         mutasiDtaDetailRepo,
		ValidationRepository:              validationRepo,
		ValidationDetailRepository:        validationDetailRepo,
		Db:                                db,
	}
}

func (s *service) Find(ctx *abstraction.Context, payload *dto.ValidationGetRequest) (*dto.ValidationGetResponse, error) {
	criteriaTB := model.TrialBalanceFilterModel{}
	criteriaTB.Period = payload.Period
	if payload.ArrCompanyID != nil && len(*payload.ArrCompanyID) > 0 {
		criteriaTB.ArrCompanyID = payload.ArrCompanyID
	}
	if payload.CompanyID != nil {
		criteriaTB.CompanyID = payload.CompanyID
	}
	criteriaTB.Status = payload.Status

	paginationTB := abstraction.Pagination{}
	paginationTB.Page = payload.Page
	paginationTB.PageSize = payload.PageSize
	paginationTB.Sort = payload.Sort
	paginationTB.SortBy = payload.SortBy

	dataTB, pagination, err := s.ValidationRepository.Find(ctx, &criteriaTB, &paginationTB)
	if err != nil {
		return nil, helper.ErrorHandler(err)
	}

	results := dto.ValidationGetResponse{
		Datas:          *dataTB,
		PaginationInfo: *pagination,
	}

	return &results, nil
}

func (s *service) FindByID(ctx *abstraction.Context, payload *dto.ValidationGetByIDRequest) (*dto.ValidationGetByIDResponse, error) {
	result := model.ValidationEntityModel{}
	if err := trxmanager.New(s.Db).WithTrx(ctx, func(ctx *abstraction.Context) error {
		dataTB, err := s.TrialBalanceRepository.FindByID(ctx, &payload.ID)
		if err != nil {
			return helper.ErrorHandler(err)
		}

		criteriaValidationDetail := model.ValidationDetailFilterModel{}
		criteriaValidationDetail.CompanyID = &dataTB.CompanyID
		criteriaValidationDetail.Period = &dataTB.Period
		criteriaValidationDetail.Versions = &dataTB.Versions

		dataValidationDetail, err := s.ValidationDetailRepository.Find(ctx, &criteriaValidationDetail)
		if err != nil {
			return helper.ErrorHandler(err)
		}

		result = model.ValidationEntityModel{
			ID:      dataTB.ID,
			Company: dataTB.Company,
			ValidationEntity: model.ValidationEntity{
				CompanyID:  dataTB.CompanyID,
				Period:     dataTB.Period,
				Versions:   dataTB.Versions,
				CreatedAt:  dataTB.CreatedAt,
				ModifiedAt: dataTB.ModifiedAt,
			},
			ValidationDetail: *dataValidationDetail,
		}

		return nil
	}); err != nil {
		return nil, err
	}

	results := dto.ValidationGetByIDResponse{
		Data: result,
	}

	return &results, nil
}

func (s *service) RequestToValidate(ctx *abstraction.Context, payload *dto.ValidationValidateRequest) error {
	type DataToValidate struct {
		MasterID   int
		ListDataID []int
		ModulID    map[int]string
	}
	listCompany, err := s.ValidationRepository.GetListCompany(ctx, &payload.CompanyID)
	if err != nil {
		return helper.ErrorHandler(err)
	}

	if len(listCompany) == 0 {
		// return response.CustomErrorBuilder(http.StatusBadRequest, "No Data with that criteria or Data has been validated", "No Data with that criteria or Data has been validated")
		return response.ErrorBuilder(&response.ErrorConstant.NotFound, err)
	}

	checkVal := make(map[int]bool)
	processTB := model.TrialBalanceEntityModel{}
	processTB.Context = ctx
	processTB.ValidationNote = "Validation on process"

	listData := []int{}
	for _, vDatavalidation := range payload.ListValidation {
		criteriaValidation := model.ValidationFilterModel{}
		criteriaValidation.TrialBalanceID = &vDatavalidation
		if checkVal[*criteriaValidation.TrialBalanceID] {
			continue
		}
		dataTBValidation, err := s.ValidationRepository.MakeSure(ctx, &criteriaValidation)
		if err != nil {
			return helper.ErrorHandler(err)
		}

		if dataTBValidation.ID == 0 {
			return response.ErrorBuilder(&response.ErrorConstant.NotFound, err)
		}

		tmp, err := time.Parse(time.RFC3339, dataTBValidation.Period)
		if err != nil {
			return err
		}
		period := tmp.Format("2006-01-02")

		if payload.Period != period {
			return response.CustomErrorBuilder(http.StatusBadRequest, "Data tidak dalam periode yang sama!", "Data tidak dalam periode yang sama!")
		}

		if listCompany[dataTBValidation.CompanyID] != true {
			return response.CustomErrorBuilder(http.StatusBadRequest, "Company tidak dalam hirarki yang sama", "Company tidak dalam hirarki yang sama")
		}

		_, err = s.TrialBalanceRepository.Update(ctx, &dataTBValidation.ID, &processTB)
		if err != nil {
			return helper.ErrorHandler(err)
		}

		listData = append(listData, dataTBValidation.ID)
		checkVal[*criteriaValidation.TrialBalanceID] = true
	}

	for _, vData := range listData {
		var tmpData DataToValidate
		tmpData.ListDataID = append(tmpData.ListDataID, vData)
		jsonTmpData, err := json.Marshal(tmpData)
		if err != nil {
			return err
		}

		data := kafka.JsonData{}
		// data.CompanyID = dataTB.CompanyID
		data.UserID = ctx.Auth.ID
		data.Data = string(jsonTmpData) //tb_id
		// data.Filter.Period = dataTB.Period
		// data.Filter.Versions = dataTB.Versions

		jsonStr, err := json.Marshal(data)
		if err != nil {
			return err
		}

		go kafka.NewService("VALIDATION").SendMessage("VALIDATE", string(jsonStr))
	}
	return nil
}

func (s *service) RequestToValidateModul(ctx *abstraction.Context, payload *dto.ValidationValidateModulRequest) error {
	type DataToValidate struct {
		MasterID   int
		ListDataID []int
		ModulID    map[int]string
	}
	var tmpData DataToValidate
	dataTB, err := s.TrialBalanceRepository.FindByID(ctx, &payload.ValidationMasterID)
	if err != nil {
		return helper.ErrorHandler(err)
	}

	if dataTB.Status != constant.MODUL_STATUS_DRAFT {
		return response.CustomErrorBuilder(http.StatusBadRequest, "No Data with that criteria or Data has been validated", "No Data with that criteria or Data has been validated")
	}

	tmpData.MasterID = dataTB.ID
	tmpData.ListDataID = append(tmpData.ListDataID, dataTB.ID)
	tmp := make(map[int]string)
	for _, vDatavalidation := range payload.ListValidation {
		criteriaValidation := model.ValidationDetailFilterModel{}
		criteriaValidation.CompanyID = &dataTB.CompanyID
		criteriaValidation.Period = &dataTB.Period
		criteriaValidation.Versions = &dataTB.Versions
		criteriaValidation.ID = &vDatavalidation

		dataValidation, err := s.ValidationDetailRepository.MakeSure(ctx, &criteriaValidation)
		if err != nil {
			return helper.ErrorHandler(err)
		}
		if dataValidation.ID == 0 {
			return response.ErrorBuilder(&response.ErrorConstant.NotFound, errors.New("No Data with that criteria or Data has been validated"))
		}
		tmp[dataValidation.ID] = dataValidation.Name
	}
	tmpData.ModulID = tmp

	jsonTmpData, err := json.Marshal(tmpData)
	if err != nil {
		return err
	}

	data := kafka.JsonData{}
	data.CompanyID = dataTB.CompanyID
	data.UserID = ctx.Auth.ID
	data.Data = string(jsonTmpData) //tb_id
	data.Filter.Period = dataTB.Period
	data.Filter.Versions = dataTB.Versions

	jsonStr, err := json.Marshal(data)
	if err != nil {
		return err
	}

	go kafka.NewService("VALIDATION").SendMessage("VALIDATE", string(jsonStr))
	return nil
}

func (s *service) ListAvailable(ctx *abstraction.Context, payload *dto.ValidationGetListAvailable) (*dto.ValidationGetResponse, error) {

	allowed := helper.CompanyValidation(ctx.Auth.ID, payload.CompanyID)
	if !allowed {
		return nil, response.ErrorBuilder(&response.ErrorConstant.UnauthorizedAccess, errors.New("Not Allowed"))
	}

	criteriaTB := model.TrialBalanceFilterModel{}
	criteriaTB.Period = &payload.Period

	paginationTB := abstraction.Pagination{}
	paginationTB.Page = payload.Page
	paginationTB.PageSize = payload.PageSize
	paginationTB.Sort = payload.Sort
	paginationTB.SortBy = payload.SortBy

	dataTB, pagination, err := s.ValidationRepository.FindListCompany(ctx, &criteriaTB, &payload.CompanyID, &paginationTB)
	if err != nil {
		return nil, helper.ErrorHandler(err)
	}

	results := dto.ValidationGetResponse{
		Datas:          *dataTB,
		PaginationInfo: *pagination,
	}

	return &results, nil
}
