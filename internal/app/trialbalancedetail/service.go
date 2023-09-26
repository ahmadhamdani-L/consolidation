package trialbalancedetail

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
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"sync"

	"github.com/Knetic/govaluate"
	"gorm.io/gorm"
)

type service struct {
	Repository                 repository.TrialBalanceDetail
	TrialBalanceRepository     repository.TrialBalance
	FormatterBridgesRepository repository.FormatterBridges
	FormatterDetailRepository  repository.FormatterDetail
	ValidationDetailRepository repository.ValidationDetail
	AjeRepository              repository.Adjustment
	JpmRepository              repository.Jpm
	JcteRepository             repository.Jcte
	JelimRepository            repository.Jelim
	Db                         *gorm.DB
}

type Service interface {
	Find(ctx *abstraction.Context, payload *dto.TrialBalanceDetailGetRequest) (*dto.TrialBalanceDetailGetResponse, error)
	FindByID(ctx *abstraction.Context, payload *dto.TrialBalanceDetailGetByIDRequest) (*dto.TrialBalanceDetailGetByIDResponse, error)
	Create(ctx *abstraction.Context, payload *dto.TrialBalanceDetailCreateRequest) (*dto.TrialBalanceDetailCreateResponse, error)
	Update(ctx *abstraction.Context, payload *dto.TrialBalanceDetailUpdateRequest) (*dto.TrialBalanceDetailUpdateResponse, error)
	Delete(ctx *abstraction.Context, payload *dto.TrialBalanceDetailDeleteRequest) (*dto.TrialBalanceDetailDeleteResponse, error)
	FindByParent(ctx *abstraction.Context, payload *dto.TrialBalanceDetailGetByParentRequest) (*dto.TrialBalanceDetailGetByParentResponse, error)
	FindAll(ctx *abstraction.Context, payload *dto.TrialBalanceDetailGetByParentRequest) (*dto.TrialBalanceDetailGetByParentResponse, error)
}

func NewService(f *factory.Factory) *service {
	repository := f.TrialBalanceDetailRepository
	tbrepository := f.TrialBalanceRepository
	formatterDetailRepo := f.FormatterDetailRepository
	formatterBridgesRepo := f.FormatterBridgesRepository
	validationDetailRepo := f.ValidationDetailRepository
	ajeRepo := f.AdjustmentRepository
	jpmRepo := f.JpmRepository
	jcteRepo := f.JcteRepository
	jelimRepo := f.JelimRepository

	db := f.Db
	return &service{
		Repository:                 repository,
		Db:                         db,
		TrialBalanceRepository:     tbrepository,
		FormatterBridgesRepository: formatterBridgesRepo,
		FormatterDetailRepository:  formatterDetailRepo,
		ValidationDetailRepository: validationDetailRepo,
		AjeRepository:              ajeRepo,
		JpmRepository:              jpmRepo,
		JcteRepository:             jcteRepo,
		JelimRepository:            jelimRepo,
	}
}

func (s *service) Find(ctx *abstraction.Context, payload *dto.TrialBalanceDetailGetRequest) (*dto.TrialBalanceDetailGetResponse, error) {
	tb, err := s.TrialBalanceRepository.FindByID(ctx, payload.TrialBalanceID)
	if err != nil {
		return &dto.TrialBalanceDetailGetResponse{}, helper.ErrorHandler(err)
	}

	allowed := helper.CompanyValidation(ctx.Auth.ID, tb.CompanyID)
	if !allowed {
		return &dto.TrialBalanceDetailGetResponse{}, response.ErrorBuilder(&response.ErrorConstant.InternalServerError, errors.New("Not Allowed"))
	}

	data, info, err := s.Repository.FindWithFormatter(ctx, &payload.TrialBalanceDetailFilterModel, &payload.Pagination)
	if err != nil {
		return &dto.TrialBalanceDetailGetResponse{}, helper.ErrorHandler(err)
	}
	tb.TrialBalanceDetail = *data
	result := &dto.TrialBalanceDetailGetResponse{
		Datas: *tb,
		// Datas:          *data,
		PaginationInfo: *info,
	}
	return result, nil
}

func (s *service) FindByID(ctx *abstraction.Context, payload *dto.TrialBalanceDetailGetByIDRequest) (*dto.TrialBalanceDetailGetByIDResponse, error) {
	data, err := s.Repository.FindByID(ctx, &payload.ID)
	if err != nil {
		return &dto.TrialBalanceDetailGetByIDResponse{}, helper.ErrorHandler(err)
	}
	fmtBridges, err := s.FormatterBridgesRepository.FindByID(ctx, &data.FormatterBridgesID)
	if err != nil {
		return nil, helper.ErrorHandler(err)
	}

	tbData, err := s.TrialBalanceRepository.FindByID(ctx, &fmtBridges.TrxRefID)
	if err != nil {
		return nil, helper.ErrorHandler(err)
	}

	allowed := helper.CompanyValidation(ctx.Auth.ID, tbData.CompanyID)
	if !allowed {
		return nil, response.ErrorBuilder(&response.ErrorConstant.UnauthorizedAccess, errors.New("Not Allowed"))
	}
	result := &dto.TrialBalanceDetailGetByIDResponse{
		TrialBalanceDetailEntityModel: *data,
	}
	return result, nil
}

func (s *service) Create(ctx *abstraction.Context, payload *dto.TrialBalanceDetailCreateRequest) (*dto.TrialBalanceDetailCreateResponse, error) {
	var data model.TrialBalanceDetailEntityModel

	if err := trxmanager.New(s.Db).WithTrx(ctx, func(ctx *abstraction.Context) error {
		data.Context = ctx
		data.TrialBalanceDetailEntity = payload.TrialBalanceDetailEntity

		fmtBridges, err := s.FormatterBridgesRepository.FindByID(ctx, &data.FormatterBridgesID)
		if err != nil {
			return helper.ErrorHandler(err)
		}

		tbData, err := s.TrialBalanceRepository.FindByID(ctx, &fmtBridges.TrxRefID)
		if err != nil {
			return helper.ErrorHandler(err)
		}

		allowed := helper.CompanyValidation(ctx.Auth.ID, tbData.CompanyID)
		if !allowed {
			return response.ErrorBuilder(&response.ErrorConstant.UnauthorizedAccess, errors.New("Not Allowed"))
		}

		result, err := s.Repository.Create(ctx, &data)
		if err != nil {
			return helper.ErrorHandler(err)
		}
		data = *result
		return nil
	}); err != nil {
		return &dto.TrialBalanceDetailCreateResponse{}, err
	}

	result := &dto.TrialBalanceDetailCreateResponse{
		TrialBalanceDetailEntityModel: data,
	}
	return result, nil
}

func (s *service) Update(ctx *abstraction.Context, payload *dto.TrialBalanceDetailUpdateRequest) (*dto.TrialBalanceDetailUpdateResponse, error) {
	var data model.TrialBalanceDetailEntityModel
	if err := trxmanager.New(s.Db).WithTrx(ctx, func(ctx *abstraction.Context) error {
		dataTB, err := s.Repository.FindByID(ctx, &payload.ID)
		if err != nil {
			return helper.ErrorHandler(err)
		}

		if dataTB.Code == "310401004" || dataTB.Code == "310402002" || dataTB.Code == "310501002" || dataTB.Code == "310502002" || dataTB.Code == "310503002" || strings.Contains(strings.ToLower(dataTB.Code), "_subtotal") {
			return response.CustomErrorBuilder(http.StatusBadRequest, response.E_BAD_REQUEST, "cannot update the data because the data has its own count.")
		}

		formatterBridgesData, err := s.FormatterBridgesRepository.FindByID(ctx, &dataTB.FormatterBridgesID)
		if err != nil {
			return helper.ErrorHandler(err)
		}

		tbData, err := s.TrialBalanceRepository.FindByID(ctx, &formatterBridgesData.TrxRefID)
		if err != nil {
			return helper.ErrorHandler(err)
		}

		if tbData.Status != 1 {
			return response.ErrorBuilder(&response.ErrorConstant.DataValidated, errors.New("Cannot Update Data"))
		}

		allowed := helper.CompanyValidation(ctx.Auth.ID, tbData.CompanyID)
		if !allowed {
			return response.ErrorBuilder(&response.ErrorConstant.UnauthorizedAccess, errors.New("Not Allowed"))
		}

		_, err = strconv.Atoi(dataTB.Code)
		if err != nil {
			// jika diconvert ke angka error, maka kemungkinan itu total atau label, bukan coa
			criteriaFormatterValidasi := model.FormatterDetailFilterModel{}
			criteriaFormatterValidasi.FormatterID = &formatterBridgesData.FormatterID
			criteriaFormatterValidasi.Code = &data.Code

			formatterValidasi, jmlData, err := s.FormatterDetailRepository.Find(ctx, &criteriaFormatterValidasi, &abstraction.Pagination{})
			if err != nil {
				return err
			}

			if jmlData.Count == 0 {
				return errors.New("No Data Found in Formatter")
			}

			for _, v := range *formatterValidasi {
				if (v.IsTotal != nil && *v.IsTotal) || (v.IsLabel != nil && *v.IsLabel) {
					return response.ErrorBuilder(&response.ErrorConstant.BadRequest, errors.New("Cannot update data"))
				}
			}
		}

		tmpHeadCoa := fmt.Sprintf("%c", dataTB.Code[0])
		afterAje := 0.0
		if tmpHeadCoa == "9" {
			tmpHeadCoa = dataTB.Code[:1]
		}
		switch tmpHeadCoa {
		case "1", "5", "6", "7", "91", "92":
			afterAje = *payload.AmountBeforeAje + *dataTB.AmountAjeDr - *dataTB.AmountAjeCr
		default:
			afterAje = *payload.AmountBeforeAje - *dataTB.AmountAjeDr + *dataTB.AmountAjeCr
		}

		data.Context = ctx
		data.TrialBalanceDetailEntity.AmountBeforeAje = payload.AmountBeforeAje

		data.TrialBalanceDetailEntity.AmountAfterAje = &afterAje

		result, err := s.Repository.Update(ctx, &payload.ID, &data)
		if err != nil {
			return helper.ErrorHandler(err)
		}

		criteriaFormatterDetailSumData := model.FormatterBridgesFilterModel{}
		mfa := "TRIAL-BALANCE"
		criteriaFormatterDetailSumData.Source = &mfa
		criteriaFormatterDetailSumData.TrxRefID = &formatterBridgesData.TrxRefID

		formatterDetailSumData, err := s.FormatterBridgesRepository.FindSummaryTB(ctx)
		if err != nil {
			return helper.ErrorHandler(err)
		}

		for _, v := range *formatterDetailSumData {
			criteriaTBDetail := model.TrialBalanceDetailFilterModel{}
			criteriaTBDetail.TrialBalanceID = &formatterBridgesData.TrxRefID
			criteriaTBDetail.FormatterBridgesID = &formatterBridgesData.ID

			if v.AutoSummary != nil && *v.AutoSummary {
				code := fmt.Sprintf("%s_Subtotal", v.Code)
				criteriaTBDetail.Code = &code
				mfadetailsum, _, err := s.Repository.Find(ctx, &criteriaTBDetail, &abstraction.Pagination{})
				if err != nil {
					return helper.ErrorHandler(err)
				}

				for _, a := range *mfadetailsum {
					sumTBD, err := s.Repository.FindSummary(ctx, &v.Code, &formatterBridgesData.ID, v.IsCoa)
					if err != nil {
						return helper.ErrorHandler(err)
					}
					updateSummary := model.TrialBalanceDetailEntityModel{
						TrialBalanceDetailEntity: model.TrialBalanceDetailEntity{
							AmountBeforeAje: sumTBD.AmountBeforeAje,
							AmountAjeDr:     sumTBD.AmountAjeDr,
							AmountAjeCr:     sumTBD.AmountAjeCr,
							AmountAfterAje:  sumTBD.AmountAfterAje,
						},
					}
					_, err = s.Repository.Update(ctx, &a.ID, &updateSummary)
					if err != nil {
						return helper.ErrorHandler(err)
					}
				}

			}

			if v.IsTotal != nil && *v.IsTotal && v.FxSummary != "" {
				// tmpString := []string{"AmountBeforeAje"}
				tmpString := []string{"AmountBeforeAje", "AmountAjeDr", "AmountAjeCr", "AmountAfterAje"}
				tmpTotalFl := make(map[string]*float64)
				// reg := regexp.MustCompile(`[0-9]+\d{3,}`)
				reg := regexp.MustCompile(`[A-Za-z0-9_~:()]+|[0-9]+\d{2,}`)

				for _, tipe := range tmpString {
					formula := strings.ToUpper(v.FxSummary)
					match := reg.FindAllString(formula, -1)
					parameters := make(map[string]interface{}, 0)
					for _, vMatch := range match {
						if len(vMatch) < 3 {
							continue
						}
						//cari jml berdasarkan code
						sumTBD, err := s.Repository.FindSummary(ctx, &vMatch, &formatterBridgesData.ID, v.IsCoa)
						if err != nil {
							return helper.ErrorHandler(err)
						}
						angka := 0.0
						if tipe == "AmountBeforeAje" && sumTBD.AmountBeforeAje != nil {
							angka = *sumTBD.AmountBeforeAje
						} else if tipe == "AmountAjeDr" && sumTBD.AmountAjeDr != nil {
							angka = *sumTBD.AmountAjeDr
						} else if tipe == "AmountAjeCr" && sumTBD.AmountAjeCr != nil {
							angka = *sumTBD.AmountAjeCr
						} else if tipe == "AmountAfterAje" && sumTBD.AmountAfterAje != nil {
							angka = *sumTBD.AmountAfterAje
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
				criteriaTBDetail.Code = &v.Code
				dataTB, err := s.Repository.FindByExactCode(ctx, &criteriaTBDetail)
				if err != nil {
					return helper.ErrorHandler(err)
				}

				updateSummary := model.TrialBalanceDetailEntityModel{
					TrialBalanceDetailEntity: model.TrialBalanceDetailEntity{
						AmountBeforeAje: tmpTotalFl["AmountBeforeAje"],
						AmountAjeDr:     tmpTotalFl["AmountAjeDr"],
						AmountAjeCr:     tmpTotalFl["AmountAjeCr"],
						AmountAfterAje:  tmpTotalFl["AmountAfterAje"],
					},
				}
				_, err = s.Repository.Update(ctx, &dataTB.ID, &updateSummary)
				if err != nil {
					return helper.ErrorHandler(err)
				}
			}

			if v.Code == "LABA_KOMPREHENSIF" {
				{
					//UPDATE CUSTOM ROW "310401004" "310402002"
					// COA 310501002 = Row 3712 --> ambil angka dari 4337
					// COA 310502002 = Row 3718 --> ambil angka dari 4342+4343
					// COA 310503002 = Row 3724 --> ambil angka dari 4345+4346
					code := "310401004"
					criteriaTBDetail := model.TrialBalanceDetailFilterModel{}
					criteriaTBDetail.FormatterBridgesID = &formatterBridgesData.ID
					criteriaTBDetail.Code = &code
					customRowOne, err := s.Repository.FindByCode(ctx, &criteriaTBDetail)
					if err != nil {
						return helper.ErrorHandler(err)
					}

					code = "310402002"
					criteriaTBDetail.Code = &code
					customRowTwo, err := s.Repository.FindByCode(ctx, &criteriaTBDetail)
					if err != nil {
						return helper.ErrorHandler(err)
					}

					code = "310501002"
					criteriaTBDetail.Code = &code
					customRowThree, err := s.Repository.FindByCode(ctx, &criteriaTBDetail)
					if err != nil {
						return helper.ErrorHandler(err)
					}

					code = "310502002"
					criteriaTBDetail.Code = &code
					customRowFour, err := s.Repository.FindByCode(ctx, &criteriaTBDetail)
					if err != nil {
						return helper.ErrorHandler(err)
					}

					code = "310503002"
					criteriaTBDetail.Code = &code
					customRowFive, err := s.Repository.FindByCode(ctx, &criteriaTBDetail)
					if err != nil {
						return helper.ErrorHandler(err)
					}

					code = "950101001" //REVALUATION FA (4337)
					criteriaTBDetail.Code = &code
					dataReFa, err := s.Repository.FindByCode(ctx, &criteriaTBDetail)
					if err != nil {
						return helper.ErrorHandler(err)
					}

					code = "950301001" //Financial Instrument (4342)
					criteriaTBDetail.Code = &code
					dataFinIn, err := s.Repository.FindByCode(ctx, &criteriaTBDetail)
					if err != nil {
						return helper.ErrorHandler(err)
					}

					code = "950301002" //Income tax relating to components of OCI (4343)
					criteriaTBDetail.Code = &code
					dataIncomeTax, err := s.Repository.FindByCode(ctx, &criteriaTBDetail)
					if err != nil {
						return helper.ErrorHandler(err)
					}

					code = "950401001" //Foreign Exchange (4345)
					criteriaTBDetail.Code = &code
					dataForex, err := s.Repository.FindByCode(ctx, &criteriaTBDetail)
					if err != nil {
						return helper.ErrorHandler(err)
					}

					code = "950401002" //Income tax relating to components of OCI (4346)
					criteriaTBDetail.Code = &code
					dataIncomeTax2, err := s.Repository.FindByCode(ctx, &criteriaTBDetail)
					if err != nil {
						return helper.ErrorHandler(err)
					}

					code = "LABA_BERSIH"
					criteriaTBDetail.Code = &code
					dataLaba, err := s.Repository.FindByCode(ctx, &criteriaTBDetail)
					if err != nil {
						return helper.ErrorHandler(err)
					}

					code = "TOTAL_PENGHASILAN_KOMPREHENSIF_LAIN~BS"
					criteriaTBDetail.Code = &code
					dataKomprehensif, err := s.Repository.FindByCode(ctx, &criteriaTBDetail)
					if err != nil {
						return helper.ErrorHandler(err)
					}

					summaryCodes, err := s.Repository.SummaryByCodes(ctx, &formatterBridgesData.ID, []string{"310501002", "310502002", "310503002"})
					if err != nil {
						return helper.ErrorHandler(err)
					}

					updatedCustomRowOne := model.TrialBalanceDetailEntityModel{}
					updatedCustomRowOne.Context = ctx
					updatedCustomRowOne.AmountBeforeAje = dataLaba.AmountBeforeAje
					updatedCustomRowOne.AmountAjeCr = dataLaba.AmountAjeCr
					updatedCustomRowOne.AmountAjeDr = dataLaba.AmountAjeDr
					updatedCustomRowOne.AmountAfterAje = dataLaba.AmountAfterAje

					_, err = s.Repository.Update(ctx, &customRowOne.ID, &updatedCustomRowOne)
					if err != nil {
						return helper.ErrorHandler(err)
					}

					updateCustomRowTwo := model.TrialBalanceDetailEntityModel{}
					updateCustomRowTwo.Context = ctx

					tmp1 := 0.0
					if dataKomprehensif.AmountBeforeAje != nil && *dataKomprehensif.AmountBeforeAje != 0 {
						tmp1 = *dataKomprehensif.AmountBeforeAje
					}
					if summaryCodes.AmountBeforeAje != nil && *summaryCodes.AmountBeforeAje != 0 {
						tmp1 = tmp1 - *summaryCodes.AmountBeforeAje
					}
					updateCustomRowTwo.AmountBeforeAje = &tmp1

					tmp2 := 0.0
					if dataKomprehensif.AmountAjeCr != nil && *dataKomprehensif.AmountAjeCr != 0 {
						tmp2 = *dataKomprehensif.AmountAjeCr
					}
					if summaryCodes.AmountAjeCr != nil && *summaryCodes.AmountAjeCr != 0 {
						tmp2 = tmp2 - *summaryCodes.AmountAjeCr
					}
					updateCustomRowTwo.AmountAjeCr = &tmp2

					tmp3 := 0.0
					if dataKomprehensif.AmountAjeDr != nil && *dataKomprehensif.AmountAjeDr != 0 {
						tmp3 = *dataKomprehensif.AmountAjeDr
					}
					if summaryCodes.AmountAjeDr != nil && *summaryCodes.AmountAjeDr != 0 {
						tmp3 = tmp3 - *summaryCodes.AmountAjeDr
					}
					updateCustomRowTwo.AmountAjeDr = &tmp3

					tmp4 := 0.0
					if dataKomprehensif.AmountAfterAje != nil && *dataKomprehensif.AmountAfterAje != 0 {
						tmp4 = *dataKomprehensif.AmountAfterAje
					}
					if summaryCodes.AmountAfterAje != nil && *summaryCodes.AmountAfterAje != 0 {
						tmp4 = tmp4 - *summaryCodes.AmountAfterAje
					}
					updateCustomRowTwo.AmountAfterAje = &tmp4

					_, err = s.Repository.Update(ctx, &customRowTwo.ID, &updateCustomRowTwo)
					if err != nil {
						return helper.ErrorHandler(err)
					}

					//

					updatedCustomRowThree := model.TrialBalanceDetailEntityModel{}
					updatedCustomRowThree.Context = ctx
					updatedCustomRowThree.AmountBeforeAje = dataReFa.AmountBeforeAje
					updatedCustomRowThree.AmountAjeCr = dataReFa.AmountAjeCr
					updatedCustomRowThree.AmountAjeDr = dataReFa.AmountAjeDr
					updatedCustomRowThree.AmountAfterAje = dataReFa.AmountAfterAje

					_, err = s.Repository.Update(ctx, &customRowThree.ID, &updatedCustomRowThree)
					if err != nil {
						return helper.ErrorHandler(err)
					}

					updateCustomRowFour := model.TrialBalanceDetailEntityModel{}
					updateCustomRowFour.Context = ctx

					tmp5 := 0.0
					if dataFinIn.AmountBeforeAje != nil && *dataFinIn.AmountBeforeAje != 0 {
						tmp5 = *dataFinIn.AmountBeforeAje
					}
					if dataIncomeTax.AmountBeforeAje != nil && *dataIncomeTax.AmountBeforeAje != 0 {
						tmp5 = tmp5 + *dataIncomeTax.AmountBeforeAje
					}
					updateCustomRowFour.AmountBeforeAje = &tmp5

					tmp6 := 0.0
					if dataFinIn.AmountAjeCr != nil && *dataFinIn.AmountAjeCr != 0 {
						tmp6 = *dataFinIn.AmountAjeCr
					}
					if dataIncomeTax.AmountAjeCr != nil && *dataIncomeTax.AmountAjeCr != 0 {
						tmp6 = tmp6 + *dataIncomeTax.AmountAjeCr
					}
					updateCustomRowFour.AmountAjeCr = &tmp6

					tmp7 := 0.0
					if dataFinIn.AmountAjeDr != nil && *dataFinIn.AmountAjeDr != 0 {
						tmp7 = *dataFinIn.AmountAjeDr
					}
					if dataIncomeTax.AmountAjeDr != nil && *dataIncomeTax.AmountAjeDr != 0 {
						tmp7 = tmp7 + *dataIncomeTax.AmountAjeDr
					}
					updateCustomRowFour.AmountAjeDr = &tmp7

					tmp8 := 0.0
					if dataFinIn.AmountAfterAje != nil && *dataFinIn.AmountAfterAje != 0 {
						tmp8 = *dataFinIn.AmountAfterAje
					}
					if dataIncomeTax.AmountAfterAje != nil && *dataIncomeTax.AmountAfterAje != 0 {
						tmp8 = tmp8 + *dataIncomeTax.AmountAfterAje
					}
					updateCustomRowFour.AmountAfterAje = &tmp8

					_, err = s.Repository.Update(ctx, &customRowFour.ID, &updateCustomRowFour)
					if err != nil {
						return helper.ErrorHandler(err)
					}

					updateCustomRowFive := model.TrialBalanceDetailEntityModel{}
					updateCustomRowFive.Context = ctx

					tmp9 := 0.0
					if dataForex.AmountBeforeAje != nil && *dataForex.AmountBeforeAje != 0 {
						tmp9 = *dataForex.AmountBeforeAje
					}
					if dataIncomeTax2.AmountBeforeAje != nil && *dataIncomeTax2.AmountBeforeAje != 0 {
						tmp9 = tmp9 + *dataIncomeTax2.AmountBeforeAje
					}
					updateCustomRowFive.AmountBeforeAje = &tmp9

					tmp10 := 0.0
					if dataForex.AmountAjeCr != nil && *dataForex.AmountAjeCr != 0 {
						tmp10 = *dataForex.AmountAjeCr
					}
					if dataIncomeTax2.AmountAjeCr != nil && *dataIncomeTax2.AmountAjeCr != 0 {
						tmp10 = tmp10 + *dataIncomeTax2.AmountAjeCr
					}
					updateCustomRowFive.AmountAjeCr = &tmp10

					tmp11 := 0.0
					if dataForex.AmountAjeDr != nil && *dataForex.AmountAjeDr != 0 {
						tmp11 = *dataForex.AmountAjeDr
					}
					if dataIncomeTax2.AmountAjeDr != nil && *dataIncomeTax2.AmountAjeDr != 0 {
						tmp11 = tmp11 + *dataIncomeTax2.AmountAjeDr
					}
					updateCustomRowFive.AmountAjeDr = &tmp11

					tmp12 := 0.0
					if dataForex.AmountAfterAje != nil && *dataForex.AmountAfterAje != 0 {
						tmp12 = *dataForex.AmountAfterAje
					}
					if dataIncomeTax2.AmountAfterAje != nil && *dataIncomeTax2.AmountAfterAje != 0 {
						tmp12 = tmp12 + *dataIncomeTax2.AmountAfterAje
					}
					updateCustomRowFive.AmountAfterAje = &tmp12

					_, err = s.Repository.Update(ctx, &customRowFive.ID, &updateCustomRowFive)
					if err != nil {
						return helper.ErrorHandler(err)
					}
				}
			}
		}

		criteriaValidation := model.ValidationDetailFilterModel{}
		criteriaValidation.CompanyID = &tbData.CompanyID
		criteriaValidation.Period = &tbData.Period
		criteriaValidation.Versions = &tbData.Versions
		criteriaValidation.Name = &formatterBridgesData.Source

		updateValidation := model.ValidationDetailEntityModel{}
		updateValidation.Status = constant.VALIDATION_STATUS_NOT_BALANCE
		updateValidation.Note = "Terdapat perubahan pada data. Silakan validasi ulang."

		_, err = s.ValidationDetailRepository.UpdateByCriteria(ctx, &criteriaValidation, &updateValidation)
		if err != nil {
			return helper.ErrorHandler(err)
		}

		data = *result
		return nil
	}); err != nil {
		return &dto.TrialBalanceDetailUpdateResponse{}, err
	}
	result := &dto.TrialBalanceDetailUpdateResponse{
		TrialBalanceDetailEntityModel: data,
	}
	return result, nil
}

func (s *service) Delete(ctx *abstraction.Context, payload *dto.TrialBalanceDetailDeleteRequest) (*dto.TrialBalanceDetailDeleteResponse, error) {
	var data model.TrialBalanceDetailEntityModel
	if err := trxmanager.New(s.Db).WithTrx(ctx, func(ctx *abstraction.Context) error {
		existing, err := s.Repository.FindByID(ctx, &payload.ID)
		if err != nil {
			return helper.ErrorHandler(err)
		}

		fmtBridges, err := s.FormatterBridgesRepository.FindByID(ctx, &existing.FormatterBridgesID)
		if err != nil {
			return helper.ErrorHandler(err)
		}

		tbData, err := s.TrialBalanceRepository.FindByID(ctx, &fmtBridges.TrxRefID)
		if err != nil {
			return helper.ErrorHandler(err)
		}

		allowed := helper.CompanyValidation(ctx.Auth.ID, tbData.CompanyID)
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
		return &dto.TrialBalanceDetailDeleteResponse{}, err
	}
	result := &dto.TrialBalanceDetailDeleteResponse{
		// TrialBalanceDetailEntityModel: data,
	}
	return result, nil
}

func (s *service) FindByParent(ctx *abstraction.Context, payload *dto.TrialBalanceDetailGetByParentRequest) (*dto.TrialBalanceDetailGetByParentResponse, error) {
	tb, err := s.TrialBalanceRepository.FindByID(ctx, &payload.TrialBalanceID)
	if err != nil {
		return nil, err
	}

	allowed := helper.CompanyValidation(ctx.Auth.ID, tb.CompanyID)
	if !allowed {
		return nil, response.ErrorBuilder(&response.ErrorConstant.UnauthorizedAccess, errors.New("Not Allowed"))
	}

	detailData, err := s.Repository.FindDetail(ctx, &payload.TrialBalanceID, &payload.ParentID)
	if err != nil {
		return nil, helper.ErrorHandler(err)
	}
	tb.TrialBalanceDetail = *detailData
	result := &dto.TrialBalanceDetailGetByParentResponse{
		Data: *tb,
	}
	return result, nil
}

func (s *service) FindAll(ctx *abstraction.Context, payload *dto.TrialBalanceDetailGetByParentRequest) (*dto.TrialBalanceDetailGetByParentResponse, error) {
	tb, err := s.TrialBalanceRepository.FindByID(ctx, &payload.TrialBalanceID)
	if err != nil {
		return nil, err
	}

	allowed := helper.CompanyValidation(ctx.Auth.ID, tb.CompanyID)
	if !allowed {
		return nil, response.ErrorBuilder(&response.ErrorConstant.UnauthorizedAccess, errors.New("Not Allowed"))
	}

	var tmpErr error
	wg := sync.WaitGroup{}
	wg.Add(1)
	go func() {
		defer wg.Done()
		detailData, err := s.Repository.FindAllDetail(ctx, &payload.TrialBalanceID)
		if err != nil {
			tmpErr = err
		}

		tb.TrialBalanceDetail = makeTreeList(*detailData, 0)
	}()

	if tmpErr != nil {
		return nil, helper.ErrorHandler(tmpErr)
	}

	sumAJE, err := s.AjeRepository.FindSummary(ctx, &tb.ID)
	if err != nil {
		return nil, helper.ErrorHandler(err)
	}

	control2WBS1, err := s.Repository.FindControlWbs1(ctx, &tb.ID)
	if err != nil {
		return nil, helper.ErrorHandler(err)
	}

	tmpcr := *helper.AssignAmount(control2WBS1.AmountAjeCr) - (*helper.AssignAmount(sumAJE.IncomeStatementCr) + *helper.AssignAmount(sumAJE.BalanceSheetCr))
	tmpdr := *helper.AssignAmount(control2WBS1.AmountAjeDr) - (*helper.AssignAmount(sumAJE.IncomeStatementDr) + *helper.AssignAmount(sumAJE.BalanceSheetDr))

	wg.Wait()

	result := &dto.TrialBalanceDetailGetByParentResponse{
		AmountBeforeAje: control2WBS1.AmountBeforeAje,
		AmountAjeCr:     &tmpcr,
		AmountAjeDr:     &tmpdr,
		AmountAfterAje:  control2WBS1.AmountAfterAje,
		Data:            *tb,
	}
	return result, nil
}

func makeTreeList(dataTB []model.TrialBalanceDetailFmtEntityModel, parent int) []model.TrialBalanceDetailFmtEntityModel {
	tbData := []model.TrialBalanceDetailFmtEntityModel{}
	for _, v := range dataTB {
		if v.Code == "310401004" || v.Code == "310402002" || v.Code == "310501002" || v.Code == "310502002" || v.Code == "310503002" {
			v.IsLabel = true
		}
		if v.ParentID == parent {
			v.Children = makeTreeList(dataTB, v.FormatterDetailID)
			tbData = append(tbData, v)
		}
	}
	return tbData
}
