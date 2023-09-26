package investasitbkdetail

import (
	"errors"
	"fmt"
	"mcash-finance-console-core/internal/abstraction"
	"mcash-finance-console-core/internal/dto"
	"mcash-finance-console-core/internal/factory"
	"mcash-finance-console-core/internal/model"
	"mcash-finance-console-core/internal/repository"
	"mcash-finance-console-core/pkg/constant"
	"mcash-finance-console-core/pkg/util/helper"
	"mcash-finance-console-core/pkg/util/response"
	"mcash-finance-console-core/pkg/util/trxmanager"
	"regexp"
	"strings"

	"github.com/Knetic/govaluate"
	"gorm.io/gorm"
)

type service struct {
	Repository                 repository.InvestasiTbkDetail
	InvestasiTBKRepository     repository.InvestasiTbk
	FormatterBridgesRepository repository.FormatterBridges
	FormatterDetailRepository  repository.FormatterDetail
	ValidationDetailRepository repository.ValidationDetail
	Db                         *gorm.DB
}

type Service interface {
	Find(ctx *abstraction.Context, payload *dto.InvestasiTbkDetailGetRequest) (*dto.InvestasiTbkDetailGetResponse, error)
	FindByID(ctx *abstraction.Context, payload *dto.InvestasiTbkDetailGetByIDRequest) (*dto.InvestasiTbkDetailGetByIDResponse, error)
	Create(ctx *abstraction.Context, payload *dto.InvestasiTbkDetailCreateRequest) (*dto.InvestasiTbkDetailCreateResponse, error)
	Update(ctx *abstraction.Context, payload *dto.InvestasiTbkDetailUpdateRequest) (*dto.InvestasiTbkDetailUpdateResponse, error)
	Delete(ctx *abstraction.Context, payload *dto.InvestasiTbkDetailDeleteRequest) (*dto.InvestasiTbkDetailDeleteResponse, error)
}

func NewService(f *factory.Factory) *service {
	repository := f.InvestasiTbkDetailRepository
	investasiTbkRepository := f.InvestasiTbkRepository
	formatterBridgesRepo := f.FormatterBridgesRepository
	validationDetailRepo := f.ValidationDetailRepository
	formatterDetailRepo := f.FormatterDetailRepository
	db := f.Db
	return &service{
		Repository:                 repository,
		InvestasiTBKRepository:     investasiTbkRepository,
		FormatterBridgesRepository: formatterBridgesRepo,
		ValidationDetailRepository: validationDetailRepo,
		FormatterDetailRepository:  formatterDetailRepo,
		Db:                         db,
	}
}

func (s *service) Find(ctx *abstraction.Context, payload *dto.InvestasiTbkDetailGetRequest) (*dto.InvestasiTbkDetailGetResponse, error) {
	investasi, err := s.InvestasiTBKRepository.FindByID(ctx, payload.InvestasiTbkID)
	if err != nil {
		return &dto.InvestasiTbkDetailGetResponse{}, helper.ErrorHandler(err)
	}

	allowed := helper.CompanyValidation(ctx.Auth.ID, investasi.CompanyID)
	if !allowed {
		return &dto.InvestasiTbkDetailGetResponse{}, response.ErrorBuilder(&response.ErrorConstant.InternalServerError, errors.New("not allowed"))
	}

	data, info, err := s.Repository.FindWithFormatter(ctx, &payload.InvestasiTbkDetailFilterModel, &payload.Pagination)
	if err != nil {
		return &dto.InvestasiTbkDetailGetResponse{}, helper.ErrorHandler(err)
	}

	TotalAmountCost, TotalAmountFv, TotalUnrealizedGain, TotalRealizedGain, TotalFee, err := s.Repository.GetTotal(ctx, &payload.InvestasiTbkDetailFilterModel)
	if err != nil {
		return nil, helper.ErrorHandler(err)
	}
	investasi.InvestasiTbkDetail = *data
	result := &dto.InvestasiTbkDetailGetResponse{
		PaginationInfo: *info,
	}
	result.Datas.Data = *investasi
	result.Datas.TotalAmountCost = TotalAmountCost
	result.Datas.TotalAmountFv = TotalAmountFv
	result.Datas.TotalUnrealizedGain = TotalUnrealizedGain
	result.Datas.TotalRealizedGain = TotalRealizedGain
	result.Datas.TotalFee = TotalFee
	return result, nil
}

func (s *service) FindByID(ctx *abstraction.Context, payload *dto.InvestasiTbkDetailGetByIDRequest) (*dto.InvestasiTbkDetailGetByIDResponse, error) {
	data, err := s.Repository.FindByID(ctx, &payload.ID)
	if err != nil {
		return &dto.InvestasiTbkDetailGetByIDResponse{}, response.ErrorBuilder(&response.ErrorConstant.InternalServerError, err)
	}

	fmtBridges, err := s.FormatterBridgesRepository.FindByID(ctx, &data.FormatterBridgesID)
	if err != nil {
		return nil, helper.ErrorHandler(err)
	}

	aup, err := s.InvestasiTBKRepository.FindByID(ctx, &fmtBridges.TrxRefID)
	if err != nil {
		return nil, helper.ErrorHandler(err)
	}

	allowed := helper.CompanyValidation(ctx.Auth.ID, aup.CompanyID)
	if !allowed {
		return nil, response.ErrorBuilder(&response.ErrorConstant.UnauthorizedAccess, errors.New("Not Allowed"))
	}

	result := &dto.InvestasiTbkDetailGetByIDResponse{
		InvestasiTbkDetailEntityModel: *data,
	}
	return result, nil
}

func (s *service) Create(ctx *abstraction.Context, payload *dto.InvestasiTbkDetailCreateRequest) (*dto.InvestasiTbkDetailCreateResponse, error) {
	var data model.InvestasiTbkDetailEntityModel

	if err := trxmanager.New(s.Db).WithTrx(ctx, func(ctx *abstraction.Context) error {
		data.Context = ctx
		data.InvestasiTbkDetailEntity = payload.InvestasiTbkDetailEntity

		result, err := s.Repository.Create(ctx, &data)
		if err != nil {
			return response.ErrorBuilder(&response.ErrorConstant.UnprocessableEntity, err)
		}
		data = *result
		return nil
	}); err != nil {
		return &dto.InvestasiTbkDetailCreateResponse{}, err
	}

	result := &dto.InvestasiTbkDetailCreateResponse{
		InvestasiTbkDetailEntityModel: data,
	}
	return result, nil
}

func (s *service) Update(ctx *abstraction.Context, payload *dto.InvestasiTbkDetailUpdateRequest) (*dto.InvestasiTbkDetailUpdateResponse, error) {
	var data model.InvestasiTbkDetailEntityModel
	if err := trxmanager.New(s.Db).WithTrx(ctx, func(ctx *abstraction.Context) error {
		existing, err := s.Repository.FindByID(ctx, &payload.ID)
		if err != nil {
			return helper.ErrorHandler(err)
		}

		fmtBridges, err := s.FormatterBridgesRepository.FindByID(ctx, &existing.FormatterBridgesID)
		if err != nil {
			return helper.ErrorHandler(err)
		}

		it, err := s.InvestasiTBKRepository.FindByID(ctx, &fmtBridges.TrxRefID)
		if err != nil {
			return helper.ErrorHandler(err)
		}

		//validasi data bukan isTotal atau autosummary
		criteriaFormatterValidasi := model.FormatterDetailFilterModel{}
		criteriaFormatterValidasi.FormatterID = &fmtBridges.FormatterID
		criteriaFormatterValidasi.Code = &existing.Stock

		formatterValidasi, jmlData, err := s.FormatterDetailRepository.Find(ctx, &criteriaFormatterValidasi, &abstraction.Pagination{})
		if err != nil {
			return helper.ErrorHandler(err)
		}

		if jmlData.Count == 0 {
			return response.ErrorBuilder(&response.ErrorConstant.NotFound, errors.New("data not found"))
		}

		for _, v := range *formatterValidasi {
			if (v.IsTotal != nil && *v.IsTotal) || (v.AutoSummary != nil && *v.AutoSummary) || (v.IsLabel != nil && *v.IsLabel) {
				return response.ErrorBuilder(&response.ErrorConstant.BadRequest, errors.New("cannot update data"))
			}
		}

		if it.Status != 1 {
			return response.ErrorBuilder(&response.ErrorConstant.DataValidated, errors.New("cannot Update Data"))
		}

		allowed := helper.CompanyValidation(ctx.Auth.ID, it.CompanyID)
		if !allowed {
			return response.ErrorBuilder(&response.ErrorConstant.UnauthorizedAccess, errors.New("not allowed"))
		}

		data.Context = ctx
		amountCost := (*payload.EndingShares) * (*payload.AvgPrice)
		amountFv := (*payload.EndingShares) * (*payload.ClosingPrice)
		unrealizedGain := amountCost - amountFv
		data.InvestasiTbkDetailEntity = model.InvestasiTbkDetailEntity{
			EndingShares:   payload.EndingShares,
			AvgPrice:       payload.AvgPrice,
			AmountCost:     &amountCost,
			ClosingPrice:   payload.ClosingPrice,
			AmountFv:       &amountFv,
			UnrealizedGain: &unrealizedGain,
			RealizedGain:   payload.RealizedGain,
			Fee:            payload.Fee,
		}
		// data.InvestasiTbkDetailEntity = payload.InvestasiTbkDetailEntity

		result, err := s.Repository.Update(ctx, &payload.ID, &data)
		if err != nil {
			return helper.ErrorHandler(err)
		}

		criteriaValidation := model.ValidationDetailFilterModel{}
		criteriaValidation.CompanyID = &it.CompanyID
		criteriaValidation.Period = &it.Period
		criteriaValidation.Versions = &it.Versions
		criteriaValidation.Name = &fmtBridges.Source

		updateValidation := model.ValidationDetailEntityModel{}
		updateValidation.Status = constant.VALIDATION_STATUS_NOT_BALANCE
		updateValidation.Note = "Terdapat perubahan pada data. Silakan validasi ulang."

		_, err = s.ValidationDetailRepository.UpdateByCriteria(ctx, &criteriaValidation, &updateValidation)
		if err != nil {
			return helper.ErrorHandler(err)
		}

		criteriaData := model.InvestasiTbkDetailFilterModel{}
		criteriaData.FormatterBridgesID = &fmtBridges.ID
		criteriaData.InvestasiTbkID = &it.ID
		nolimit := 1000
		paginationData := abstraction.Pagination{}
		paginationData.PageSize = &nolimit
		datas, _, err := s.Repository.FindWithFormatter(ctx, &criteriaData, &paginationData)
		if err != nil {
			return helper.ErrorHandler(err)
		}
		vSumEndingShare := 0.
		vSumAVGPrice := 0.
		vSumAmountCost := 0.
		vSumClosingPrice := 0.
		vSumAmountFv := 0.
		vSumUnrealizedGain := 0.
		vSumRealizedGain := 0.
		vSumFee := 0.

		for _, v := range *datas {
			if !(v.AutoSummary != nil && *v.AutoSummary) {
				vSumEndingShare += *v.EndingShares
				vSumAVGPrice += *v.AvgPrice
				vSumAmountCost += *v.AmountCost
				vSumClosingPrice += *v.ClosingPrice
				vSumAmountFv += *v.AmountFv
				vSumUnrealizedGain += *v.UnrealizedGain
				vSumRealizedGain += *v.RealizedGain
				vSumFee += *v.Fee
			}

			if v.AutoSummary != nil && *v.AutoSummary {
				updatedData := model.InvestasiTbkDetailEntityModel{}
				// updatedData.EndingShares = &vSumEndingShare
				// updatedData.AvgPrice = &vSumAVGPrice
				updatedData.AmountCost = &vSumAmountCost
				// updatedData.ClosingPrice = &vSumClosingPrice
				updatedData.AmountFv = &vSumAmountFv
				updatedData.UnrealizedGain = &vSumUnrealizedGain
				updatedData.RealizedGain = &vSumRealizedGain
				updatedData.Fee = &vSumFee

				_, err = s.Repository.Update(ctx, &v.ID, &updatedData)
				if err != nil {
					return helper.ErrorHandler(err)
				}

				vSumEndingShare = 0
				vSumAVGPrice = 0
				vSumAmountCost = 0
				vSumClosingPrice = 0
				vSumAmountFv = 0
				vSumUnrealizedGain = 0
				vSumRealizedGain = 0
				vSumFee = 0
			}

			if v.IsTotal != nil && *v.IsTotal && v.FxSummary != nil && *v.FxSummary != "" {
				// tmpString := []string{"AmountBeforeAje"}
				tmpString := []string{"EndingShare", "AVGPrice", "AmountCost", "ClosingPrice", "AmountFV", "UnrealizedGain", "RealizedGain", "Fee"}
				tmpTotalFl := make(map[string]*float64)
				// reg := regexp.MustCompile(`[0-9]+\d{3,}`)
				reg := regexp.MustCompile(`[A-Za-z0-9_~:()]+|[0-9]+\d{3,}`)

				for _, tipe := range tmpString {
					formula := strings.ToUpper(*v.FxSummary)
					match := reg.FindAllString(formula, -1)
					parameters := make(map[string]interface{}, 0)
					for _, vMatch := range match {
						if len(vMatch) < 3 {
							continue
						}
						//cari jml berdasarkan code
						sumIT, err := s.Repository.FindByStock(ctx, &fmtBridges.ID, &vMatch)
						if err != nil {
							return helper.ErrorHandler(err)
						}
						angka := 0.0
						if tipe == "EndingShare" && sumIT.EndingShares != nil {
							angka = *sumIT.EndingShares
						} else if tipe == "AVGPrice" && sumIT.AvgPrice != nil {
							angka = *sumIT.AvgPrice
						} else if tipe == "AmountAjeCr" && sumIT.AmountCost != nil {
							angka = *sumIT.AmountCost
						} else if tipe == "AmountAfterAje" && sumIT.ClosingPrice != nil {
							angka = *sumIT.ClosingPrice
						} else if tipe == "AmountFV" && sumIT.AmountFv != nil {
							angka = *sumIT.AmountFv
						} else if tipe == "UnrealizedGain" && sumIT.UnrealizedGain != nil {
							angka = *sumIT.UnrealizedGain
						} else if tipe == "RealizedGain" && sumIT.RealizedGain != nil {
							angka = *sumIT.RealizedGain
						} else if tipe == "Fee" && sumIT.Fee != nil {
							angka = *sumIT.Fee
						}
						formula = helper.ReplaceWholeWord(formula, vMatch, fmt.Sprintf("(%2.f)", angka))
						// parameters[vMatch] = angka

					}
					expressionFormula, err := govaluate.NewEvaluableExpression(formula)
					if err != nil {
						return err
					}
					result, err := expressionFormula.Evaluate(parameters)
					if err != nil {
						return helper.ErrorHandler(err)
					}
					tmp := result.(float64)
					tmpTotalFl[tipe] = &tmp
				}

				updateSummary := model.InvestasiTbkDetailEntityModel{}
				updateSummary.EndingShares = tmpTotalFl["EndingShare"]
				updateSummary.AvgPrice = tmpTotalFl["AVGPrice"]
				updateSummary.AmountCost = tmpTotalFl["AmountCost"]
				updateSummary.ClosingPrice = tmpTotalFl["ClosingPrice"]
				updateSummary.AmountFv = tmpTotalFl["AmountFV"]
				updateSummary.UnrealizedGain = tmpTotalFl["UnrealizedGain"]
				updateSummary.RealizedGain = tmpTotalFl["RealizedGain"]
				updateSummary.Fee = tmpTotalFl["Fee"]

				_, err = s.Repository.Update(ctx, &v.ID, &updateSummary)
				if err != nil {
					return helper.ErrorHandler(err)
				}
			}
		}

		data = *result
		return nil
	}); err != nil {
		return &dto.InvestasiTbkDetailUpdateResponse{}, err
	}
	result := &dto.InvestasiTbkDetailUpdateResponse{
		InvestasiTbkDetailEntityModel: data,
	}
	return result, nil
}

func (s *service) Delete(ctx *abstraction.Context, payload *dto.InvestasiTbkDetailDeleteRequest) (*dto.InvestasiTbkDetailDeleteResponse, error) {
	var data model.InvestasiTbkDetailEntityModel
	if err := trxmanager.New(s.Db).WithTrx(ctx, func(ctx *abstraction.Context) error {
		existing, err := s.Repository.FindByID(ctx, &payload.ID)
		if err != nil {
			return helper.ErrorHandler(err)
		}

		fmtBridges, err := s.FormatterBridgesRepository.FindByID(ctx, &existing.FormatterBridgesID)
		if err != nil {
			return helper.ErrorHandler(err)
		}

		aup, err := s.InvestasiTBKRepository.FindByID(ctx, &fmtBridges.TrxRefID)
		if err != nil {
			return helper.ErrorHandler(err)
		}

		allowed := helper.CompanyValidation(ctx.Auth.ID, aup.CompanyID)
		if !allowed {
			return response.ErrorBuilder(&response.ErrorConstant.UnauthorizedAccess, errors.New("not allowed"))
		}

		data.Context = ctx
		result, err := s.Repository.Delete(ctx, &payload.ID, &data)
		if err != nil {
			return helper.ErrorHandler(err)
		}
		data = *result
		return nil
	}); err != nil {
		return &dto.InvestasiTbkDetailDeleteResponse{}, response.ErrorBuilder(&response.ErrorConstant.InternalServerError, err)
	}
	result := &dto.InvestasiTbkDetailDeleteResponse{
		InvestasiTbkDetailEntityModel: data,
	}
	return result, nil
}
