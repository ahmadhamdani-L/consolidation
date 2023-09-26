package pembelianpenjualanberelasidetail

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
	Repository                 repository.PembelianPenjualanBerelasiDetail
	PPBRepo                    repository.PembelianPenjualanBerelasi
	ValidationDetailRepository repository.ValidationDetail
	Db                         *gorm.DB
}

type Service interface {
	Find(ctx *abstraction.Context, payload *dto.PembelianPenjualanBerelasiDetailGetRequest) (*dto.PembelianPenjualanBerelasiDetailGetResponse, error)
	FindByID(ctx *abstraction.Context, payload *dto.PembelianPenjualanBerelasiDetailGetByIDRequest) (*dto.PembelianPenjualanBerelasiDetailGetByIDResponse, error)
	Create(ctx *abstraction.Context, payload *dto.PembelianPenjualanBerelasiDetailCreateRequest) (*dto.PembelianPenjualanBerelasiDetailCreateResponse, error)
	Update(ctx *abstraction.Context, payload *dto.PembelianPenjualanBerelasiDetailUpdateRequest) (*dto.PembelianPenjualanBerelasiDetailUpdateResponse, error)
	Delete(ctx *abstraction.Context, payload *dto.PembelianPenjualanBerelasiDetailDeleteRequest) (*dto.PembelianPenjualanBerelasiDetailDeleteResponse, error)
}

func NewService(f *factory.Factory) *service {
	repository := f.PembelianPenjualanBerelasiDetailRepository
	ppbrepo := f.PembelianPenjualanBerelasiRepository
	validationDetailRepo := f.ValidationDetailRepository
	db := f.Db
	return &service{
		Repository:                 repository,
		Db:                         db,
		PPBRepo:                    ppbrepo,
		ValidationDetailRepository: validationDetailRepo,
	}
}

func (s *service) Find(ctx *abstraction.Context, payload *dto.PembelianPenjualanBerelasiDetailGetRequest) (*dto.PembelianPenjualanBerelasiDetailGetResponse, error) {
	pembelianpenjualanberelasi, err := s.PPBRepo.FindByID(ctx, payload.PembelianPenjualanBerelasiID)
	if err != nil {
		return &dto.PembelianPenjualanBerelasiDetailGetResponse{}, helper.ErrorHandler(err)
	}
	allowed := helper.CompanyValidation(ctx.Auth.ID, pembelianpenjualanberelasi.CompanyID)
	if !allowed {
		return &dto.PembelianPenjualanBerelasiDetailGetResponse{}, response.ErrorBuilder(&response.ErrorConstant.InternalServerError, errors.New("Not Allowed"))
	}
	data, info, err := s.Repository.Find(ctx, &payload.PembelianPenjualanBerelasiDetailFilterModel, &payload.Pagination)
	if err != nil {
		return &dto.PembelianPenjualanBerelasiDetailGetResponse{}, helper.ErrorHandler(err)
	}

	criteriaTotal := model.PembelianPenjualanBerelasiDetailEntityModel{}
	criteriaTotal.PembelianPenjualanBerelasiID = *payload.PembelianPenjualanBerelasiID
	pembelian, penjualan, err := s.Repository.GetTotal(ctx, &payload.PembelianPenjualanBerelasiDetailFilterModel)
	if err != nil {
		return &dto.PembelianPenjualanBerelasiDetailGetResponse{}, helper.ErrorHandler(err)
	}
	pembelianpenjualanberelasi.PembelianPenjualanBerelasiDetail = *data

	result := &dto.PembelianPenjualanBerelasiDetailGetResponse{
		PaginationInfo: *info,
	}
	result.Datas.Data = *pembelianpenjualanberelasi
	result.Datas.TotalPembelian = pembelian
	result.Datas.TotalPenjualan = penjualan
	return result, nil
}

func (s *service) FindByID(ctx *abstraction.Context, payload *dto.PembelianPenjualanBerelasiDetailGetByIDRequest) (*dto.PembelianPenjualanBerelasiDetailGetByIDResponse, error) {
	data, err := s.Repository.FindByID(ctx, &payload.ID)
	if err != nil {
		return &dto.PembelianPenjualanBerelasiDetailGetByIDResponse{}, helper.ErrorHandler(err)
	}

	ppbData, err := s.PPBRepo.FindByID(ctx, &data.PembelianPenjualanBerelasiID)
	if err != nil {
		return nil, helper.ErrorHandler(err)
	}

	allowed := helper.CompanyValidation(ctx.Auth.ID, ppbData.CompanyID)
	if !allowed {
		return nil, response.ErrorBuilder(&response.ErrorConstant.UnauthorizedAccess, errors.New("Not Allowed"))
	}

	result := &dto.PembelianPenjualanBerelasiDetailGetByIDResponse{
		PembelianPenjualanBerelasiDetailEntityModel: *data,
	}
	return result, nil
}

func (s *service) Create(ctx *abstraction.Context, payload *dto.PembelianPenjualanBerelasiDetailCreateRequest) (*dto.PembelianPenjualanBerelasiDetailCreateResponse, error) {
	var data model.PembelianPenjualanBerelasiDetailEntityModel

	ppbData, err := s.PPBRepo.FindByID(ctx, &payload.PembelianPenjualanBerelasiID)
	if err != nil {
		return nil, helper.ErrorHandler(err)
	}

	allowed := helper.CompanyValidation(ctx.Auth.ID, ppbData.CompanyID)
	if !allowed {
		return nil, response.ErrorBuilder(&response.ErrorConstant.UnauthorizedAccess, errors.New("Not Allowed"))
	}

	if err := trxmanager.New(s.Db).WithTrx(ctx, func(ctx *abstraction.Context) error {
		data.Context = ctx
		data.PembelianPenjualanBerelasiDetailEntity = payload.PembelianPenjualanBerelasiDetailEntity

		result, err := s.Repository.Create(ctx, &data)
		if err != nil {
			return response.ErrorBuilder(&response.ErrorConstant.UnprocessableEntity, err)
		}
		data = *result
		return nil
	}); err != nil {
		return &dto.PembelianPenjualanBerelasiDetailCreateResponse{}, err
	}

	result := &dto.PembelianPenjualanBerelasiDetailCreateResponse{
		PembelianPenjualanBerelasiDetailEntityModel: data,
	}
	return result, nil
}

func (s *service) Update(ctx *abstraction.Context, payload *dto.PembelianPenjualanBerelasiDetailUpdateRequest) (*dto.PembelianPenjualanBerelasiDetailUpdateResponse, error) {
	var data model.PembelianPenjualanBerelasiDetailEntityModel
	if err := trxmanager.New(s.Db).WithTrx(ctx, func(ctx *abstraction.Context) error {
		ppbDetail, err := s.Repository.FindByID(ctx, &payload.ID)
		if err != nil {
			return helper.ErrorHandler(err)
		}

		ppb, err := s.PPBRepo.FindByID(ctx, &ppbDetail.PembelianPenjualanBerelasiID)
		if err != nil {
			return helper.ErrorHandler(err)
		}

		if ppb.Status != 1 {
			return response.ErrorBuilder(&response.ErrorConstant.DataValidated, errors.New("Cannot Update Data"))
		}

		allowed := helper.CompanyValidation(ctx.Auth.ID, ppb.CompanyID)
		if !allowed {
			return response.ErrorBuilder(&response.ErrorConstant.UnauthorizedAccess, errors.New("Not Allowed"))
		}

		data.Context = ctx
		// data.PembelianPenjualanBerelasiDetailEntity = payload.PembelianPenjualanBerelasiDetailEntity
		data.PembelianPenjualanBerelasiDetailEntity = model.PembelianPenjualanBerelasiDetailEntity{
			BoughtAmount: payload.BoughtAmount,
			SalesAmount:  payload.SalesAmount,
		}

		result, err := s.Repository.Update(ctx, &payload.ID, &data)
		if err != nil {
			return helper.ErrorHandler(err)
		}

		/* criteriaValidation := model.ValidationDetailFilterModel{}
		criteriaValidation.CompanyID = &ppb.CompanyID
		criteriaValidation.Period = &ppb.Period
		criteriaValidation.Versions = &ppb.Versions
		criteriaValidation.Name = &formatterBridgesData.Source

		updateValidation := model.ValidationDetailEntityModel{}
		updateValidation.Status = constant.VALIDATION_STATUS_NOT_BALANCE
		updateValidation.Note = "Terdapat perubahan pada data. Silakan validasi ulang."

		_, err = s.ValidationDetailRepository.UpdateByCriteria(ctx, &criteriaValidation, &updateValidation)
		if err != nil {
			return helper.ErrorHandler(err)
		} */

		data = *result
		return nil
	}); err != nil {
		return &dto.PembelianPenjualanBerelasiDetailUpdateResponse{}, err
	}
	result := &dto.PembelianPenjualanBerelasiDetailUpdateResponse{
		PembelianPenjualanBerelasiDetailEntityModel: data,
	}
	return result, nil
}

func (s *service) Delete(ctx *abstraction.Context, payload *dto.PembelianPenjualanBerelasiDetailDeleteRequest) (*dto.PembelianPenjualanBerelasiDetailDeleteResponse, error) {
	var data model.PembelianPenjualanBerelasiDetailEntityModel
	if err := trxmanager.New(s.Db).WithTrx(ctx, func(ctx *abstraction.Context) error {
		existing, err := s.Repository.FindByID(ctx, &payload.ID)
		if err != nil {
			return helper.ErrorHandler(err)
		}
		ppbData, err := s.PPBRepo.FindByID(ctx, &existing.PembelianPenjualanBerelasiID)
		if err != nil {
			return helper.ErrorHandler(err)
		}

		allowed := helper.CompanyValidation(ctx.Auth.ID, ppbData.CompanyID)
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
		return &dto.PembelianPenjualanBerelasiDetailDeleteResponse{}, helper.ErrorHandler(err)
	}
	result := &dto.PembelianPenjualanBerelasiDetailDeleteResponse{
		// PembelianPenjualanBerelasiDetailEntityModel: data,
	}
	return result, nil
}
