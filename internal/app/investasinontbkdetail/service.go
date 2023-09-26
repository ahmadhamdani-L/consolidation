package investasinontbkdetail

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

	"gorm.io/gorm"
)

type service struct {
	Repository                repository.InvestasiNonTbkDetail
	InvestasiNonTbkRepository repository.InvestasiNonTbk
	Db                        *gorm.DB
}

type Service interface {
	Find(ctx *abstraction.Context, payload *dto.InvestasiNonTbkDetailGetRequest) (*dto.InvestasiNonTbkDetailGetResponse, error)
	FindByID(ctx *abstraction.Context, payload *dto.InvestasiNonTbkDetailGetByIDRequest) (*dto.InvestasiNonTbkDetailGetByIDResponse, error)
	Create(ctx *abstraction.Context, payload *dto.InvestasiNonTbkDetailCreateRequest) (*dto.InvestasiNonTbkDetailCreateResponse, error)
	Update(ctx *abstraction.Context, payload *dto.InvestasiNonTbkDetailUpdateRequest) (*dto.InvestasiNonTbkDetailUpdateResponse, error)
	Delete(ctx *abstraction.Context, payload *dto.InvestasiNonTbkDetailDeleteRequest) (*dto.InvestasiNonTbkDetailDeleteResponse, error)
}

func NewService(f *factory.Factory) *service {
	repository := f.InvestasiNonTbkDetailRepository
	investasinontbk := f.InvestasiNonTbkRepository
	db := f.Db
	return &service{
		Repository:                repository,
		InvestasiNonTbkRepository: investasinontbk,
		Db:                        db,
	}
}

func (s *service) Find(ctx *abstraction.Context, payload *dto.InvestasiNonTbkDetailGetRequest) (*dto.InvestasiNonTbkDetailGetResponse, error) {
	investasinontbk, err := s.InvestasiNonTbkRepository.FindByID(ctx, payload.InvestasiNonTbkID)
	if err != nil {
		return &dto.InvestasiNonTbkDetailGetResponse{}, helper.ErrorHandler(err)
	}
	allowed := helper.CompanyValidation(ctx.Auth.ID, investasinontbk.CompanyID)
	if !allowed {
		return &dto.InvestasiNonTbkDetailGetResponse{}, response.ErrorBuilder(&response.ErrorConstant.UnauthorizedAccess, response.ErrorBuilder(&response.ErrorConstant.UnauthorizedAccess, errors.New("Not Allowed")))
	}
	data, info, err := s.Repository.Find(ctx, &payload.InvestasiNonTbkDetailFilterModel, &payload.Pagination)
	if err != nil {
		return &dto.InvestasiNonTbkDetailGetResponse{}, helper.ErrorHandler(err)
	}
	investasinontbk.InvestasiNonTbkDetail = *data
	result := &dto.InvestasiNonTbkDetailGetResponse{
		Datas:          *investasinontbk,
		PaginationInfo: *info,
	}
	return result, nil
}

func (s *service) FindByID(ctx *abstraction.Context, payload *dto.InvestasiNonTbkDetailGetByIDRequest) (*dto.InvestasiNonTbkDetailGetByIDResponse, error) {
	data, err := s.Repository.FindByID(ctx, &payload.ID)
	if err != nil {
		return &dto.InvestasiNonTbkDetailGetByIDResponse{}, helper.ErrorHandler(err)
	}

	int, err := s.InvestasiNonTbkRepository.FindByID(ctx, &data.InvestasiNonTbkID)
	if err != nil {
		return nil, err
	}

	allowed := helper.CompanyValidation(ctx.Auth.ID, int.CompanyID)
	if !allowed {
		return nil, response.ErrorBuilder(&response.ErrorConstant.UnauthorizedAccess, errors.New("Not Allowed"))
	}

	result := &dto.InvestasiNonTbkDetailGetByIDResponse{
		InvestasiNonTbkDetailEntityModel: *data,
	}
	return result, nil
}

func (s *service) Create(ctx *abstraction.Context, payload *dto.InvestasiNonTbkDetailCreateRequest) (*dto.InvestasiNonTbkDetailCreateResponse, error) {
	var data model.InvestasiNonTbkDetailEntityModel

	aup, err := s.InvestasiNonTbkRepository.FindByID(ctx, &payload.InvestasiNonTbkID)
	if err != nil {
		return nil, err
	}

	allowed := helper.CompanyValidation(ctx.Auth.ID, aup.CompanyID)
	if !allowed {
		return nil, response.ErrorBuilder(&response.ErrorConstant.UnauthorizedAccess, errors.New("Not Allowed"))
	}

	if err := trxmanager.New(s.Db).WithTrx(ctx, func(ctx *abstraction.Context) error {
		data.Context = ctx
		data.InvestasiNonTbkDetailEntity = payload.InvestasiNonTbkDetailEntity

		result, err := s.Repository.Create(ctx, &data)
		if err != nil {
			return helper.ErrorHandler(err)
		}
		data = *result
		return nil
	}); err != nil {
		return &dto.InvestasiNonTbkDetailCreateResponse{}, err
	}

	result := &dto.InvestasiNonTbkDetailCreateResponse{
		InvestasiNonTbkDetailEntityModel: data,
	}
	return result, nil
}

func (s *service) Update(ctx *abstraction.Context, payload *dto.InvestasiNonTbkDetailUpdateRequest) (*dto.InvestasiNonTbkDetailUpdateResponse, error) {
	var data model.InvestasiNonTbkDetailEntityModel
	if err := trxmanager.New(s.Db).WithTrx(ctx, func(ctx *abstraction.Context) error {
		existing, err := s.Repository.FindByID(ctx, &payload.ID)
		if err != nil {
			return helper.ErrorHandler(err)
		}

		intData, err := s.InvestasiNonTbkRepository.FindByID(ctx, &existing.InvestasiNonTbkID)
		if err != nil {
			return err
		}

		if intData.Status != 1 {
			return response.ErrorBuilder(&response.ErrorConstant.DataValidated, errors.New("Cannot update validated data"))
		}

		allowed := helper.CompanyValidation(ctx.Auth.ID, intData.CompanyID)
		if !allowed {
			return response.ErrorBuilder(&response.ErrorConstant.UnauthorizedAccess, errors.New("Not Allowed"))
		}

		data.Context = ctx
		persentaseOwner := (*payload.LbrSahamOwnership) / (*payload.TotalLbrSaham)
		totalHargaPar := (*payload.LbrSahamOwnership) * (*payload.HargaPar)
		totalHargaBeli := (*payload.LbrSahamOwnership) * (*payload.HargaBeli)
		data.InvestasiNonTbkDetailEntity = model.InvestasiNonTbkDetailEntity{
			LbrSahamOwnership:   payload.LbrSahamOwnership,
			TotalLbrSaham:       payload.TotalLbrSaham,
			PercentageOwnership: &persentaseOwner,
			HargaPar:            payload.HargaPar,
			TotalHargaPar:       &totalHargaPar,
			HargaBeli:           payload.HargaBeli,
			TotalHargaBeli:      &totalHargaBeli,
		}
		// data.InvestasiNonTbkDetailEntity = payload.InvestasiNonTbkDetailEntity

		result, err := s.Repository.Update(ctx, &payload.ID, &data)
		if err != nil {
			return helper.ErrorHandler(err)
		}

		// criteriaValidation := model.ValidationDetailFilterModel{}
		// criteriaValidation.CompanyID = &intData.CompanyID
		// criteriaValidation.Period = &intData.Period
		// criteriaValidation.Versions = &intData.Versions
		// criteriaValidation.Name = &formatterBridgesData.Source

		// updateValidation := model.ValidationDetailEntityModel{}
		// updateValidation.Status = constant.VALIDATION_STATUS_NOT_BALANCE
		// updateValidation.Note = "Terdapat perubahan pada data. Silakan validasi ulang."

		// _, err = s.ValidationDetailRepository.UpdateByCriteria(ctx, &criteriaValidation, &updateValidation)
		// if err != nil {
		// 	return helper.ErrorHandler(err)
		// }

		data = *result
		return nil
	}); err != nil {
		return &dto.InvestasiNonTbkDetailUpdateResponse{}, err
	}
	result := &dto.InvestasiNonTbkDetailUpdateResponse{
		InvestasiNonTbkDetailEntityModel: data,
	}
	return result, nil
}

func (s *service) Delete(ctx *abstraction.Context, payload *dto.InvestasiNonTbkDetailDeleteRequest) (*dto.InvestasiNonTbkDetailDeleteResponse, error) {
	var data model.InvestasiNonTbkDetailEntityModel
	if err := trxmanager.New(s.Db).WithTrx(ctx, func(ctx *abstraction.Context) error {
		existing, err := s.Repository.FindByID(ctx, &payload.ID)
		if err != nil {
			return helper.ErrorHandler(err)
		}

		int, err := s.InvestasiNonTbkRepository.FindByID(ctx, &existing.InvestasiNonTbkID)
		if err != nil {
			return helper.ErrorHandler(err)
		}

		allowed := helper.CompanyValidation(ctx.Auth.ID, int.CompanyID)
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
		return &dto.InvestasiNonTbkDetailDeleteResponse{}, helper.ErrorHandler(err)
	}
	result := &dto.InvestasiNonTbkDetailDeleteResponse{
		InvestasiNonTbkDetailEntityModel: data,
	}
	return result, nil
}
