package company

import (
	"errors"
	"mcash-finance-console-core/internal/abstraction"
	"mcash-finance-console-core/internal/dto"
	"mcash-finance-console-core/internal/factory"
	"mcash-finance-console-core/internal/model"
	"mcash-finance-console-core/internal/repository"
	"mcash-finance-console-core/pkg/util/helper"
	"mcash-finance-console-core/pkg/util/response"
	"mcash-finance-console-core/pkg/util/trxmanager"
	"net/http"
	"regexp"
	"strings"

	"gorm.io/gorm"
)

type service struct {
	Repository             repository.Company
	TrialBalanceRepository repository.TrialBalance
	Db                     *gorm.DB
}

type Service interface {
	Find(ctx *abstraction.Context, payload *dto.CompanyGetRequest) (*dto.CompanyGetResponse, error)
	FindByID(ctx *abstraction.Context, payload *dto.CompanyGetByIDRequest) (*dto.CompanyGetByIDResponse, error)
	Create(ctx *abstraction.Context, payload *dto.CompanyCreateRequest) (*dto.CompanyCreateResponse, error)
	Update(ctx *abstraction.Context, payload *dto.CompanyUpdateRequest) (*dto.CompanyUpdateResponse, error)
	Delete(ctx *abstraction.Context, payload *dto.CompanyDeleteRequest) (*dto.CompanyDeleteResponse, error)
	FindTree(ctx *abstraction.Context) (*dto.CompanyGetTreeviewResponse, error)
	FindListFilter(ctx *abstraction.Context) (*dto.CompanyGetResponse, error)
}

func NewService(f *factory.Factory) *service {
	repository := f.CompanyRepository
	db := f.Db
	tbRepo := f.TrialBalanceRepository
	return &service{
		Repository:             repository,
		Db:                     db,
		TrialBalanceRepository: tbRepo,
	}
}

func (s *service) Find(ctx *abstraction.Context, payload *dto.CompanyGetRequest) (*dto.CompanyGetResponse, error) {
	filter := model.CompanyFilterModel{}
	filter = payload.CompanyFilterModel
	if payload.ChildCompany != nil && *payload.ChildCompany {
		filter.ParentCompanyID = nil
	}
	data, info, err := s.Repository.Find(ctx, &filter, &payload.Pagination)
	if err != nil {
		return &dto.CompanyGetResponse{}, helper.ErrorHandler(err)
	}
	result := &dto.CompanyGetResponse{
		Datas:          *data,
		PaginationInfo: *info,
	}
	return result, nil
}

func (s *service) FindByID(ctx *abstraction.Context, payload *dto.CompanyGetByIDRequest) (*dto.CompanyGetByIDResponse, error) {
	data, err := s.Repository.FindByID(ctx, &payload.ID)
	if err != nil {
		return &dto.CompanyGetByIDResponse{}, err
	}
	result := &dto.CompanyGetByIDResponse{
		CompanyEntityModel: *data,
	}
	return result, nil
}

func (s *service) Create(ctx *abstraction.Context, payload *dto.CompanyCreateRequest) (*dto.CompanyCreateResponse, error) {
	var data model.CompanyEntityModel

	if err := trxmanager.New(s.Db).WithTrx(ctx, func(ctx *abstraction.Context) error {

		regex := regexp.MustCompile("[^a-zA-Z0-9]+")
		payload.Code = regex.ReplaceAllString(payload.Code , "")
		payload.Code = strings.ToUpper(payload.Code)
		// payload.Name = strings.ToUpper(payload.Name)
		
		uniqCode, err := s.Repository.FindWithCode(ctx, &payload.Code)
		if err != nil {
			return err
		}
		uniqName, err := s.Repository.FindWithName(ctx, &payload.Name)
		if err != nil && err != gorm.ErrRecordNotFound {
			return err
		}
		
		if uniqCode.Code == payload.Code {
			return response.CustomErrorBuilder(http.StatusBadRequest, "Code Company "+payload.Code+ " Sudah Terpakai", "Code Company "+payload.Code+ " Sudah Terpakai")
		}
		if len(*uniqName) > 0 {
			return response.CustomErrorBuilder(http.StatusBadRequest, "Name Company "+payload.Name+ " Sudah Terpakai", "Name Company "+payload.Name+ " Sudah Terpakai")
		}
		

		data.Context = ctx
		data.CompanyEntity = payload.CompanyEntity
		result, err := s.Repository.Create(ctx, &data)
		if err != nil {
			return helper.ErrorHandler(err)
		}
		data = *result
		return nil
	}); err != nil {
		return &dto.CompanyCreateResponse{}, err
	}
	result := &dto.CompanyCreateResponse{
		CompanyEntityModel: data,
	}
	return result, nil
}

func (s *service) Update(ctx *abstraction.Context, payload *dto.CompanyUpdateRequest) (*dto.CompanyUpdateResponse, error) {
	var data model.CompanyEntityModel
	if err := trxmanager.New(s.Db).WithTrx(ctx, func(ctx *abstraction.Context) error {
		if _, err := s.Repository.FindByID(ctx, &payload.ID); err != nil {
			return helper.ErrorHandler(err)
		}
		findByOne, err := s.Repository.FindByID(ctx, &payload.ID)
		if err != nil {
			return err
		}
		regex := regexp.MustCompile("[^a-zA-Z0-9]+")
		payload.Code = regex.ReplaceAllString(payload.Code , "")
		payload.Code = strings.ToUpper(payload.Code )
		
		uniq, err := s.Repository.FindWithCode(ctx, &payload.Code)
		if err != nil {
			return err
		}
		if payload.Code != findByOne.Code  {
			if uniq.Code == payload.Code {
				return response.CustomErrorBuilder(http.StatusBadRequest, "Code Company "+payload.Code+ " Sudah Terpakai", "Code Company "+payload.Code+ " Sudah Terpakai")
			}
		}
		
		data.Context = ctx
		data.CompanyEntity = payload.CompanyEntity
		result, err := s.Repository.Update(ctx, &payload.ID, &data)
		if err != nil {
			return helper.ErrorHandler(err)
		}
		data = *result
		return nil
	}); err != nil {
		return &dto.CompanyUpdateResponse{}, err
	}
	result := &dto.CompanyUpdateResponse{
		CompanyEntityModel: data,
	}
	return result, nil
}

func (s *service) Delete(ctx *abstraction.Context, payload *dto.CompanyDeleteRequest) (*dto.CompanyDeleteResponse, error) {
	var data model.CompanyEntityModel
	if err := trxmanager.New(s.Db).WithTrx(ctx, func(ctx *abstraction.Context) error {
		if _, err := s.Repository.FindByID(ctx, &payload.ID); err != nil {
			return helper.ErrorHandler(err)
		}
		tbCriteria := model.TrialBalanceFilterModel{}
		tbCriteria.CompanyID = &payload.ID
		jmlData, err := s.TrialBalanceRepository.GetCount(ctx, &tbCriteria)
		if err != nil {
			return helper.ErrorHandler(err)
		}
		if jmlData != nil && *jmlData > 0 {
			return helper.ErrorHandler(errors.New("Tidak bisa menghapus data karena sudah terdapat data trial balance"))
		}
		data.Context = ctx
		result, err := s.Repository.Delete(ctx, &payload.ID, &data)
		if err != nil {
			return helper.ErrorHandler(err)
		}
		data = *result
		return nil
	}); err != nil {
		return &dto.CompanyDeleteResponse{}, err
	}
	result := &dto.CompanyDeleteResponse{
		// CompanyEntityModel: data,
	}
	return result, nil
}

func (s *service) FindTree(ctx *abstraction.Context) (*dto.CompanyGetTreeviewResponse, error) {
	filter := model.CompanyFilterModel{}
	pagination := abstraction.Pagination{}
	nolimit := 10000
	pagination.PageSize = &nolimit
	data, _, err := s.Repository.Find(ctx, &filter, &pagination)
	if err != nil {
		return nil, helper.ErrorHandler(err)
	}
	datas := makeTreeList(*data, 0)
	result := &dto.CompanyGetTreeviewResponse{
		Datas: datas,
	}
	return result, nil
}

func makeTreeList(dataTB []model.CompanyEntityModel, parent int) []model.CompanyEntityModel {
	tbData := []model.CompanyEntityModel{}
	for _, v := range dataTB {
		if v.ParentCompanyID != nil && *v.ParentCompanyID == parent {
			v.ChildCompany = makeTreeList(dataTB, v.ID)
			tbData = append(tbData, v)
		} else {
			tbData = append(tbData, v)
		}
	}
	return tbData
}

func (s *service) FindListFilter(ctx *abstraction.Context) (*dto.CompanyGetResponse, error) {
	data, err := s.Repository.FindFilterList(ctx)
	if err != nil {
		return nil, helper.ErrorHandler(err)
	}
	result := &dto.CompanyGetResponse{
		Datas: *data,
	}
	return result, nil
}
