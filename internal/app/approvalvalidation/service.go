package approvalvalidation

import (
	"errors"
	"mcash-finance-console-core/internal/abstraction"
	"mcash-finance-console-core/internal/dto"
	"mcash-finance-console-core/internal/factory"
	"mcash-finance-console-core/internal/model"
	"mcash-finance-console-core/internal/repository"
	"mcash-finance-console-core/pkg/constant"
	"mcash-finance-console-core/pkg/util/helper"
	"mcash-finance-console-core/pkg/util/response"
	"mcash-finance-console-core/pkg/util/trxmanager"
	"net/http"

	"gorm.io/gorm"
)

type service struct {
	ApprovalValidationRepository repository.ApprovalValidation
	TrialBalanceRepository       repository.TrialBalance
	Db                           *gorm.DB
}

type Service interface {
	Find(ctx *abstraction.Context, payload *dto.ApprovalValidationGetRequest) (*dto.ApprovalValidationGetResponse, error)
	Approve(ctx *abstraction.Context, payload *dto.ApproveValidationRequest) (*dto.ApproveValidationResponse, error)
}

func NewService(f *factory.Factory) *service {
	approveRepo := f.ApprovalValidationRepository
	tbRepo := f.TrialBalanceRepository
	db := f.Db
	return &service{
		ApprovalValidationRepository: approveRepo,
		TrialBalanceRepository:       tbRepo,
		Db:                           db,
	}
}

func (s *service) Find(ctx *abstraction.Context, payload *dto.ApprovalValidationGetRequest) (*dto.ApprovalValidationGetResponse, error) {

	criteria := model.TrialBalanceFilterModel{}
	criteria.TrialBalanceFilter = payload.TrialBalanceFilter
	data, pagination, err := s.ApprovalValidationRepository.Find(ctx, &criteria, &payload.Pagination)
	if err != nil {
		return nil, helper.ErrorHandler(err)
	}

	return &dto.ApprovalValidationGetResponse{
		Datas:          *data,
		PaginationInfo: *pagination,
	}, nil
}

func (s *service) Approve(ctx *abstraction.Context, payload *dto.ApproveValidationRequest) (*dto.ApproveValidationResponse, error) {
	data, err := s.TrialBalanceRepository.FindByID(ctx, &payload.TrialBalanceID)
	if err != nil {
		return nil, helper.ErrorHandler(err)
	}

	if data.Status != constant.MODUL_STATUS_DRAFT {
		return nil, response.CustomErrorBuilder(http.StatusBadRequest, response.E_BAD_REQUEST, "data not found/has been validated")
	}

	allowed := helper.CompanyValidation(ctx.Auth.ID, data.CompanyID)
	if !allowed {
		return nil, response.ErrorBuilder(&response.ErrorConstant.UnauthorizedAccess, errors.New("not allowed"))
	}

	if err := trxmanager.New(s.Db).WithTrx(ctx, func(ctx *abstraction.Context) error {
		criteriaApprove := model.ValidationDetailFilterModel{}
		criteriaApprove.CompanyID = &data.CompanyID
		criteriaApprove.Period = &data.Period
		criteriaApprove.Versions = &data.Versions
		err = s.ApprovalValidationRepository.Approve(ctx, &criteriaApprove)
		if err != nil {
			return helper.ErrorHandler(err)
		}

		return nil
	}); err != nil {
		return nil, helper.ErrorHandler(err)
	}

	return &dto.ApproveValidationResponse{
		Success: true,
	}, nil
}
