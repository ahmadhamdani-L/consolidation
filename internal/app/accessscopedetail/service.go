package accessscopedetail

import (
	"mcash-finance-console-core/internal/abstraction"
	"mcash-finance-console-core/internal/dto"
	"mcash-finance-console-core/internal/factory"
	"mcash-finance-console-core/internal/model"
	"mcash-finance-console-core/internal/repository"
	"mcash-finance-console-core/pkg/util/helper"
	"mcash-finance-console-core/pkg/util/response"
	"mcash-finance-console-core/pkg/util/trxmanager"
	"net/http"

	"gorm.io/gorm"
)

type Service interface {
	Find(ctx *abstraction.Context, payload *dto.AccessScopeDetailGetRequest) (*dto.AccessScopeDetailGetResponse, error)
	FindByID(ctx *abstraction.Context, payload *dto.AccessScopeDetailGetByIDRequest) (*dto.AccessScopeDetailGetByIDResponse, error)
	Create(ctx *abstraction.Context, payload *dto.AccessScopeDetailCreateRequest) (*dto.AccessScopeDetailCreateResponse, error)
	Update(ctx *abstraction.Context, payload *dto.AccessScopeDetailUpdateRequest) (*dto.AccessScopeDetailUpdateResponse, error)
	Delete(ctx *abstraction.Context, payload *dto.AccessScopeDetailDeleteRequest) (*dto.AccessScopeDetailDeleteResponse, error)
	FindWithCompanyByID(ctx *abstraction.Context, payload *dto.AccessScopeDetailGetCompanyListRequest) (*dto.AccessScopeDetailGetCompanyListResponse, error)
}

type service struct {
	AccessScopeRepository repository.AccessScope
	Repository            repository.AccessScopeDetail
	UserRepository        repository.User
	Db                    *gorm.DB
}

func NewService(f *factory.Factory) *service {
	accessScopeRepo := f.AccessScopeRepository
	repository := f.AccessScopeDetailRepository
	userRepo := f.UserRepository
	db := f.Db
	return &service{
		AccessScopeRepository: accessScopeRepo,
		Repository:            repository,
		Db:                    db,
		UserRepository:        userRepo,
	}
}

func (s *service) Find(ctx *abstraction.Context, payload *dto.AccessScopeDetailGetRequest) (*dto.AccessScopeDetailGetResponse, error) {
	data, info, err := s.Repository.Find(ctx, &payload.AccessScopeDetailFilterModel, &payload.Pagination)
	if err != nil {
		return &dto.AccessScopeDetailGetResponse{}, helper.ErrorHandler(err)
	}

	result := &dto.AccessScopeDetailGetResponse{
		Datas:          *data,
		PaginationInfo: *info,
	}
	return result, nil
}

func (s *service) FindByID(ctx *abstraction.Context, payload *dto.AccessScopeDetailGetByIDRequest) (*dto.AccessScopeDetailGetByIDResponse, error) {
	data, err := s.Repository.FindByID(ctx, &payload.ID)
	if err != nil {
		return &dto.AccessScopeDetailGetByIDResponse{}, helper.ErrorHandler(err)
	}
	result := &dto.AccessScopeDetailGetByIDResponse{
		AccessScopeDetailEntityModel: *data,
	}
	return result, nil
}

func (s *service) Create(ctx *abstraction.Context, payload *dto.AccessScopeDetailCreateRequest) (*dto.AccessScopeDetailCreateResponse, error) {
	var data model.AccessScopeDetailEntityModel

	if err := trxmanager.New(s.Db).WithTrx(ctx, func(ctx *abstraction.Context) error {
		data.Context = ctx
		data.AccessScopeDetailEntity = payload.AccessScopeDetailEntity

		result, err := s.Repository.Create(ctx, &data)
		if err != nil {
			return helper.ErrorHandler(err)
		}
		data = *result
		return nil
	}); err != nil {
		return &dto.AccessScopeDetailCreateResponse{}, err
	}

	result := &dto.AccessScopeDetailCreateResponse{
		AccessScopeDetailEntityModel: data,
	}
	return result, nil
}

func (s *service) Update(ctx *abstraction.Context, payload *dto.AccessScopeDetailUpdateRequest) (*dto.AccessScopeDetailUpdateResponse, error) {
	var data model.AccessScopeDetailEntityModel
	if err := trxmanager.New(s.Db).WithTrx(ctx, func(ctx *abstraction.Context) error {
		if _, err := s.Repository.FindByID(ctx, &payload.ID); err != nil {
			return helper.ErrorHandler(err)
		}
		data.Context = ctx
		data.AccessScopeDetailEntity = payload.AccessScopeDetailEntity

		result, err := s.Repository.Update(ctx, &payload.ID, &data)
		if err != nil {
			return helper.ErrorHandler(err)
		}

		data = *result
		return nil
	}); err != nil {
		return &dto.AccessScopeDetailUpdateResponse{}, err
	}
	result := &dto.AccessScopeDetailUpdateResponse{
		AccessScopeDetailEntityModel: data,
	}
	return result, nil
}

func (s *service) Delete(ctx *abstraction.Context, payload *dto.AccessScopeDetailDeleteRequest) (*dto.AccessScopeDetailDeleteResponse, error) {
	var data model.AccessScopeDetailEntityModel
	if err := trxmanager.New(s.Db).WithTrx(ctx, func(ctx *abstraction.Context) error {
		acessScope, err := s.Repository.FindByID(ctx, &payload.ID)
		if err != nil {
			return helper.ErrorHandler(err)
		}
		user, err := s.UserRepository.FindByID(ctx, &acessScope.AccessScope.UserID)
		if err != nil {
			return helper.ErrorHandler(err)
		}
		if user.CompanyID == acessScope.CompanyID {
			return response.CustomErrorBuilder(http.StatusBadRequest, response.E_BAD_REQUEST, "Cannot delete data because its default company of the user")
		}
		data.Context = ctx
		result, err := s.Repository.Delete(ctx, &payload.ID, &data)
		if err != nil {
			return helper.ErrorHandler(err)
		}
		data = *result
		return nil
	}); err != nil {
		return &dto.AccessScopeDetailDeleteResponse{}, helper.ErrorHandler(err)
	}
	result := &dto.AccessScopeDetailDeleteResponse{
		AccessScopeDetailEntityModel: data,
	}
	return result, nil
}

func (s *service) FindWithCompanyByID(ctx *abstraction.Context, payload *dto.AccessScopeDetailGetCompanyListRequest) (*dto.AccessScopeDetailGetCompanyListResponse, error) {
	data, err := s.AccessScopeRepository.FindByID(ctx, &payload.ID)
	if err != nil {
		return &dto.AccessScopeDetailGetCompanyListResponse{}, helper.ErrorHandler(err)
	}
	datas, err := s.AccessScopeRepository.FindWithCompanyByID(ctx, &data.ID)
	if err != nil {
		return &dto.AccessScopeDetailGetCompanyListResponse{}, helper.ErrorHandler(err)
	}

	dataResult := makeTreeList(*datas, nil)
	result := &dto.AccessScopeDetailGetCompanyListResponse{
		Data: dataResult,
	}
	return result, nil
}

func makeTreeList(data []model.AccessScopeDetailListEntityModel, parent *int) []model.AccessScopeDetailListEntityModel {
	result := []model.AccessScopeDetailListEntityModel{}
	for _, v := range data {
		if v.ParentCompanyID == parent || (v.ParentCompanyID != nil && parent != nil && *v.ParentCompanyID == *parent) {
			v.Child = makeTreeList(data, &v.ID)
			if len(v.Child) > 0 {
				v.IsParent = true
			}
			result = append(result, v)
		}
	}
	return result
}
