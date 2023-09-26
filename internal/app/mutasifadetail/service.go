package mutasifadetail

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
	Repository                 repository.MutasiFaDetail
	MutasiFaRepository         repository.MutasiFa
	FormatterBridgesRepository repository.FormatterBridges
	FormatterDetailRepository  repository.FormatterDetail
	ValidationDetailRepository repository.ValidationDetail
	FormatterRepository         repository.Formatter
	TrialBalanceDetailRepository			repository.TrialBalanceDetail
	Db                         *gorm.DB
}

type Service interface {
	Find(ctx *abstraction.Context, payload *dto.MutasiFaDetailGetRequest) (*dto.MutasiFaDetailGetResponse, error)
	FindByID(ctx *abstraction.Context, payload *dto.MutasiFaDetailGetByIDRequest) (*dto.MutasiFaDetailGetByIDResponse, error)
	Create(ctx *abstraction.Context, payload *dto.MutasiFaDetailCreateRequest) (*dto.MutasiFaDetailCreateResponse, error)
	Update(ctx *abstraction.Context, payload *dto.MutasiFaDetailUpdateRequest) (*dto.MutasiFaDetailUpdateResponse, error)
	Delete(ctx *abstraction.Context, payload *dto.MutasiFaDetailDeleteRequest) (*dto.MutasiFaDetailDeleteResponse, error)
}

func NewService(f *factory.Factory) *service {
	repository := f.MutasiFaDetailRepository
	mutasiFaRepo := f.MutasiFaRepository
	formatterBridgesRepo := f.FormatterBridgesRepository
	formatterDetailRepo := f.FormatterDetailRepository
	validationDetailRepo := f.ValidationDetailRepository
	formatterRepo := f.FormatterRepository
	trialBalanceDetailRepo := f.TrialBalanceDetailRepository
	db := f.Db
	return &service{
		Repository:                 repository,
		MutasiFaRepository:         mutasiFaRepo,
		FormatterBridgesRepository: formatterBridgesRepo,
		FormatterDetailRepository:  formatterDetailRepo,
		ValidationDetailRepository: validationDetailRepo,
		FormatterRepository :        formatterRepo,
		TrialBalanceDetailRepository:trialBalanceDetailRepo,
		Db:                         db,
	}
}

func (s *service) Find(ctx *abstraction.Context, payload *dto.MutasiFaDetailGetRequest) (*dto.MutasiFaDetailGetResponse, error) {
	mutasi, err := s.MutasiFaRepository.FindByID(ctx, payload.MutasiFaID)
	if err != nil {
		return &dto.MutasiFaDetailGetResponse{}, helper.ErrorHandler(err)
	}
	criteriaTb := model.TrialBalanceFilterModel{}
	criteriaTb.Period = &mutasi.Period
	criteriaTb.CompanyID = &mutasi.CompanyID
	criteriaTb.Status = &mutasi.Status
	criteriaTb.Versions = &mutasi.Versions

	trialBalance, err := s.TrialBalanceDetailRepository.FindByCriteriaTb(ctx, &criteriaTb)
	if err != nil {
		return nil, err
	}
	sourceTB := "TRIAL-BALANCE"
	criteriaFmtB1 := model.FormatterBridgesFilterModel{}
	criteriaFmtB1.Source = &sourceTB
	criteriaFmtB1.TrxRefID = &trialBalance.ID

	fmtBridgesTB, err := s.FormatterBridgesRepository.FindWithCriteria(ctx, &criteriaFmtB1)
	if err != nil {
		return nil, err
	}
	allowed := helper.CompanyValidation(ctx.Auth.ID, mutasi.CompanyID)
	if !allowed {
		return &dto.MutasiFaDetailGetResponse{}, response.ErrorBuilder(&response.ErrorConstant.InternalServerError, errors.New("not allowed"))
	}

	criteriaFormatterBridges := model.FormatterBridgesFilterModel{}
	criteriaFormatterBridges.TrxRefID = &mutasi.ID
	src := "MUTASI-FA"
	criteriaFormatterBridges.Source = &src
	id := "id"
	asc := "ASC"
	pagingFB := abstraction.Pagination{
		SortBy: &id,
		Sort:   &asc,
	}
	formatterBridges, _, err := s.FormatterBridgesRepository.Find(ctx, &criteriaFormatterBridges, &pagingFB)
	if err != nil {
		return &dto.MutasiFaDetailGetResponse{}, helper.ErrorHandler(err)
	}
	for _, v := range *formatterBridges {
		payload.FormatterBridgesID = &v.ID
		data, err := s.Repository.FindWithFormatter(ctx, &payload.MutasiFaDetailFilterModel)
		if err != nil {
			return &dto.MutasiFaDetailGetResponse{}, helper.ErrorHandler(err)
		}
		fmtData, err := s.FormatterRepository.FindByID(ctx, &v.FormatterID)
		if err != nil {
				return nil, err
			}
		tmpT := true
		criteriaFmtDetail := model.FormatterDetailFilterModel{}
		criteriaFmtDetail.IsControl = &tmpT
		criteriaFmtDetail.FormatterID = &fmtData.ID
		fmtDetail, err := s.FormatterDetailRepository.FindWithCriteria(ctx, &criteriaFmtDetail)
		if err != nil {
			continue
		}
			var dataControl model.MutasiFaDetailEntityModel
			var dataControlMFA []model.MutasiFaDetailEntityModel
			for _, vFmt := range *fmtDetail {
				criteriaController := model.ControllerFilterModel{}
				criteriaController.FormatterID = &vFmt.FormatterID
				criteriaController.Code = &vFmt.Code

				controller, err := s.FormatterDetailRepository.FindByCriteriaControl(ctx, &criteriaController)
				if err != nil {
					return nil, err
				}
				for _, cntrl := range *controller {
					tmpCodeCoa1 := cntrl.Coa1
					criteriaSummaryCoa1 := model.TrialBalanceDetailFilterModel{}
					criteriaSummaryCoa1.Code = &tmpCodeCoa1
					criteriaSummaryCoa1.FormatterBridgesID = &fmtBridgesTB.ID

					summaryCoa1, err := s.TrialBalanceDetailRepository.FindSummaryTb(ctx, &criteriaSummaryCoa1)
					if err != nil {
						return nil, err
					}

					splitControllerCommand := strings.Split(cntrl.ControllerCommand, ".")
					if len(splitControllerCommand) < 2 || len(splitControllerCommand) > 4 {
						return nil, err
					}
					// find table
					criteriaMFD := model.MutasiFaDetailFilterModel{}
					criteriaMFD.Code = &vFmt.Code
					criteriaMFD.FormatterBridgesID = &v.ID
					dataMFD, err := s.Repository.FindByCriteria(ctx, &criteriaMFD) // data Aging Utang Piutang Detail
					if err != nil || dataMFD.ID == 0 {
						return nil, err
					}
					
					switch strings.ToLower(splitControllerCommand[1]) {
					case "ending_balance":
						dataControls := *dataMFD.BeginningBalance - *summaryCoa1.AmountBeforeAje
						dataControl.Control = &dataControls
					case "additions":
						dataControls := *dataMFD.Additions - *summaryCoa1.AmountBeforeAje
						dataControl.Additions = &dataControls
					case "deductions":
					dataControls := *dataMFD.Deductions - *summaryCoa1.AmountBeforeAje
					dataControl.Reclassification = &dataControls
					}
					
					dataControlMFA = append(dataControlMFA, dataControl)
					
				}
			}
		switch v.Formatter.FormatterFor {
		case "MUTASI-FA-COST":
			mutasi.MutasiFaCostDetail = *data
			mutasi.ControlMFACost = dataControlMFA
		case "MUTASI-FA-ACCUMULATED-DEPRECATION":
			mutasi.MutasiFaADDetail = *data
			mutasi.ControlMFAD = dataControlMFA
		case "MUTASI-DETAIL-PENGURANGAN":
			mutasi.MutasiDetailPengurangan = *data
			mutasi.ControlMFAPengurangan = dataControlMFA
		}
	}
	result := &dto.MutasiFaDetailGetResponse{
		Datas: *mutasi,
	}
	return result, nil
}

func (s *service) FindByID(ctx *abstraction.Context, payload *dto.MutasiFaDetailGetByIDRequest) (*dto.MutasiFaDetailGetByIDResponse, error) {
	data, err := s.Repository.FindByID(ctx, &payload.ID)
	if err != nil {
		return &dto.MutasiFaDetailGetByIDResponse{}, helper.ErrorHandler(err)
	}

	fmtBridges, err := s.FormatterBridgesRepository.FindByID(ctx, &data.FormatterBridgesID)
	if err != nil {
		return nil, err
	}

	mfaData, err := s.MutasiFaRepository.FindByID(ctx, &fmtBridges.TrxRefID)
	if err != nil {
		return nil, err
	}

	allowed := helper.CompanyValidation(ctx.Auth.ID, mfaData.CompanyID)
	if !allowed {
		return nil, response.ErrorBuilder(&response.ErrorConstant.UnauthorizedAccess, errors.New("not allowed"))
	}

	result := &dto.MutasiFaDetailGetByIDResponse{
		MutasiFaDetailEntityModel: *data,
	}
	return result, nil
}

func (s *service) Create(ctx *abstraction.Context, payload *dto.MutasiFaDetailCreateRequest) (*dto.MutasiFaDetailCreateResponse, error) {
	var data model.MutasiFaDetailEntityModel

	fmtBridges, err := s.FormatterBridgesRepository.FindByID(ctx, &payload.FormatterBridgesID)
	if err != nil {
		return nil, helper.ErrorHandler(err)
	}

	mfaData, err := s.MutasiFaRepository.FindByID(ctx, &fmtBridges.TrxRefID)
	if err != nil {
		return nil, helper.ErrorHandler(err)
	}

	allowed := helper.CompanyValidation(ctx.Auth.ID, mfaData.CompanyID)
	if !allowed {
		return nil, response.ErrorBuilder(&response.ErrorConstant.UnauthorizedAccess, errors.New("not allowed"))
	}

	if err := trxmanager.New(s.Db).WithTrx(ctx, func(ctx *abstraction.Context) error {
		data.Context = ctx
		data.MutasiFaDetailEntity = payload.MutasiFaDetailEntity

		result, err := s.Repository.Create(ctx, &data)
		if err != nil {
			return helper.ErrorHandler(err)
		}
		data = *result
		return nil
	}); err != nil {
		return &dto.MutasiFaDetailCreateResponse{}, err
	}

	result := &dto.MutasiFaDetailCreateResponse{
		MutasiFaDetailEntityModel: data,
	}
	return result, nil
}

func (s *service) Update(ctx *abstraction.Context, payload *dto.MutasiFaDetailUpdateRequest) (*dto.MutasiFaDetailUpdateResponse, error) {
	var data model.MutasiFaDetailEntityModel
	if err := trxmanager.New(s.Db).WithTrx(ctx, func(ctx *abstraction.Context) error {
		mfDetail, err := s.Repository.FindByID(ctx, &payload.ID)
		if err != nil {
			return helper.ErrorHandler(err)
		}

		formatterBridgesData, err := s.FormatterBridgesRepository.FindByID(ctx, &mfDetail.FormatterBridgesID)
		if err != nil {
			return helper.ErrorHandler(err)
		}

		mfaData, err := s.MutasiFaRepository.FindByID(ctx, &formatterBridgesData.TrxRefID)
		if err != nil {
			return helper.ErrorHandler(err)
		}

		if mfaData.Status != 1 {
			return response.ErrorBuilder(&response.ErrorConstant.DataValidated, errors.New("cannot update data"))
		}

		allowed := helper.CompanyValidation(ctx.Auth.ID, mfaData.CompanyID)
		if !allowed {
			return response.ErrorBuilder(&response.ErrorConstant.UnauthorizedAccess, errors.New("not allowed"))
		}

		//validasi data bukan isTotal atau autosummary
		criteriaFormatterValidasi := model.FormatterDetailFilterModel{}
		criteriaFormatterValidasi.FormatterID = &formatterBridgesData.FormatterID
		criteriaFormatterValidasi.Code = &mfDetail.Code

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

		data.Context = ctx
		// data.MutasiFaDetailEntity = payload.MutasiFaDetailEntity
		endingBalance := (*payload.BeginningBalance) + (*payload.AcquisitionOfSubsidiary) + (*payload.Additions) - (*payload.Deductions) + (*payload.Reclassification) + (*payload.Revaluation)
		data.MutasiFaDetailEntity = model.MutasiFaDetailEntity{
			BeginningBalance:        payload.BeginningBalance,
			AcquisitionOfSubsidiary: payload.AcquisitionOfSubsidiary,
			Additions:               payload.Additions,
			Deductions:              payload.Deductions,
			Reclassification:        payload.Reclassification,
			Revaluation:             payload.Revaluation,
			EndingBalance:           &endingBalance,
		}

		result, err := s.Repository.Update(ctx, &payload.ID, &data)
		if err != nil {
			return helper.ErrorHandler(err)
		}

		criteriaFormatterDetailSumData := model.FormatterBridgesFilterModel{}
		mfa := "MUTASI-FA"
		criteriaFormatterDetailSumData.Source = &mfa
		criteriaFormatterDetailSumData.TrxRefID = &formatterBridgesData.TrxRefID

		formatterDetailSumData, err := s.FormatterBridgesRepository.FindSummary(ctx, &criteriaFormatterDetailSumData)
		if err != nil {
			return helper.ErrorHandler(err)
		}

		criteriaFormatterDetail := model.MutasiFaDetailFilterModel{}
		criteriaFormatterDetail.MutasiFaID = &formatterBridgesData.TrxRefID

		for _, v := range *formatterDetailSumData {
			criteriaMutasiDetail := model.MutasiFaDetailFilterModel{}
			// criteriaMutasiDetail.FormatterBridgesID = &formatterBridgesData.ID
			criteriaMutasiDetail.Code = &v.Code
			criteriaMutasiDetail.MutasiFaID = &formatterBridgesData.TrxRefID

			if v.AutoSummary != nil && *v.AutoSummary {
				filterHelperFormatterBridges := model.FormatterBridgesFilterModel{}
				filterHelperFormatterBridges.FormatterID = &v.FormatterID
				filterHelperFormatterBridges.Source = &mfa
				filterHelperFormatterBridges.TrxRefID = &formatterBridgesData.TrxRefID
				helperFormatterBridges, _, err := s.FormatterBridgesRepository.Find(ctx, &filterHelperFormatterBridges, &abstraction.Pagination{})
				if err != nil {
					return helper.ErrorHandler(err)
				}
				for _, hv := range *helperFormatterBridges {
					criteriaMutasiDetail.FormatterBridgesID = &hv.ID
				}

				mfadetailsum, _, err := s.Repository.Find(ctx, &criteriaMutasiDetail, &abstraction.Pagination{})
				if err != nil {
					return helper.ErrorHandler(err)
				}

				for _, a := range *mfadetailsum {
					sumMFA, err := s.Repository.FindSummary(ctx, &v.FormatterID, &a.FormatterBridgesID, &v.SortID)
					if err != nil {
						return helper.ErrorHandler(err)
					}
					updateSummary := model.MutasiFaDetailEntityModel{
						MutasiFaDetailEntity: model.MutasiFaDetailEntity{
							BeginningBalance:        sumMFA.BeginningBalance,
							AcquisitionOfSubsidiary: sumMFA.AcquisitionOfSubsidiary,
							Additions:               sumMFA.Additions,
							Deductions:              sumMFA.Deductions,
							Reclassification:        sumMFA.Reclassification,
							Revaluation:             sumMFA.Revaluation,
							EndingBalance:           sumMFA.EndingBalance,
						},
					}
					_, err = s.Repository.Update(ctx, &a.ID, &updateSummary)
					if err != nil {
						return helper.ErrorHandler(err)
					}
				}
			}

			if v.IsTotal != nil && *v.IsTotal && v.FxSummary != "" {
				if v.Code == "CONTROL_1" || v.Code == "CONTROL_2" {
					continue
				}
				// tmpString := []string{"AmountBeforeAje"}
				tmpString := []string{"BeginningBalance", "AcquisitionOfSubsidiary", "Additions", "Deductions", "Reclassification", "Revaluation", "EndingBalance"}
				tmpTotalFl := make(map[string]*float64)
				// reg := regexp.MustCompile(`[0-9]+\d{3,}`)
				reg := regexp.MustCompile(`[A-Za-z0-9_~:()]+|[0-9]+\d{3,}`)
				for _, tipe := range tmpString {
					formula := strings.ToUpper(v.FxSummary)
					match := reg.FindAllString(formula, -1)
					parameters := make(map[string]interface{}, 0)
					for _, vMatch := range match {
						if len(vMatch) < 3 {
							continue
						}
						//cari jml berdasarkan code
						criteriaSumFA := model.MutasiFaDetailFilterModel{}
						criteriaSumFA.Code = &vMatch
						criteriaSumFA.MutasiFaID = &mfaData.ID
						sumFA, err := s.Repository.FindTotal(ctx, &criteriaSumFA)
						if err != nil {
							return helper.ErrorHandler(err)
						}
						angka := 0.0
						if tipe == "BeginningBalance" && sumFA.BeginningBalance != nil {
							angka = *sumFA.BeginningBalance
						}
						if tipe == "AcquisitionOfSubsidiary" && sumFA.AcquisitionOfSubsidiary != nil {
							angka = *sumFA.AcquisitionOfSubsidiary
						}
						if tipe == "Additions" && sumFA.Additions != nil {
							angka = *sumFA.Additions
						}
						if tipe == "Deductions" && sumFA.Deductions != nil {
							angka = *sumFA.Deductions
						}
						if tipe == "Reclassification" && sumFA.Reclassification != nil {
							angka = *sumFA.Reclassification
						}
						if tipe == "Revaluation" && sumFA.Revaluation != nil {
							angka = *sumFA.Revaluation
						}
						if tipe == "EndingBalance" && sumFA.EndingBalance != nil {
							angka = *sumFA.EndingBalance
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

				mdtadetailsum, _, err := s.Repository.Find(ctx, &criteriaMutasiDetail, &abstraction.Pagination{})
				if err != nil {
					return helper.ErrorHandler(err)
				}
				// defaultFloat := 0.0
				for _, a := range *mdtadetailsum {
					updateSummary := model.MutasiFaDetailEntityModel{
						
						MutasiFaDetailEntity: model.MutasiFaDetailEntity{
							BeginningBalance:        tmpTotalFl["BeginningBalance"],
							AcquisitionOfSubsidiary: tmpTotalFl["AcquisitionOfSubsidiary"],
							Additions:               tmpTotalFl["Additions"],
							Deductions:              tmpTotalFl["Deductions"],
							Reclassification:        tmpTotalFl["Reclassification"],
							Revaluation:             tmpTotalFl["Revaluation"],
							EndingBalance:           tmpTotalFl["EndingBalance"],
						},
					}
					_, err = s.Repository.Update(ctx, &a.ID, &updateSummary)
					if err != nil {
						return helper.ErrorHandler(err)
					}
				}
			}
		}

		criteriaValidation := model.ValidationDetailFilterModel{}
		criteriaValidation.CompanyID = &mfaData.CompanyID
		criteriaValidation.Period = &mfaData.Period
		criteriaValidation.Versions = &mfaData.Versions
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
		return &dto.MutasiFaDetailUpdateResponse{}, err
	}
	result := &dto.MutasiFaDetailUpdateResponse{
		MutasiFaDetailEntityModel: data,
	}
	return result, nil
}

func (s *service) Delete(ctx *abstraction.Context, payload *dto.MutasiFaDetailDeleteRequest) (*dto.MutasiFaDetailDeleteResponse, error) {
	var data model.MutasiFaDetailEntityModel
	if err := trxmanager.New(s.Db).WithTrx(ctx, func(ctx *abstraction.Context) error {
		if _, err := s.Repository.FindByID(ctx, &payload.ID); err != nil {
			return helper.ErrorHandler(err)
		}
		data.Context = ctx
		result, err := s.Repository.Delete(ctx, &payload.ID, &data)
		if err != nil {
			return helper.ErrorHandler(err)
		}
		data = *result
		return nil
	}); err != nil {
		return &dto.MutasiFaDetailDeleteResponse{}, err
	}
	result := &dto.MutasiFaDetailDeleteResponse{
		MutasiFaDetailEntityModel: data,
	}
	return result, nil
}
