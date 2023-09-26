package mutasiruadetail

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
	Repository                 repository.MutasiRuaDetail
	MutasiRuaRepository         repository.MutasiRua
	FormatterBridgesRepository repository.FormatterBridges
	FormatterDetailRepository  repository.FormatterDetail
	ValidationDetailRepository repository.ValidationDetail
	FormatterRepository          repository.Formatter
	TrialBalanceDetailRepository repository.TrialBalanceDetail
	Db                         *gorm.DB
}

type Service interface {
	Find(ctx *abstraction.Context, payload *dto.MutasiRuaDetailGetRequest) (*dto.MutasiRuaDetailGetResponse, error)
	FindByID(ctx *abstraction.Context, payload *dto.MutasiRuaDetailGetByIDRequest) (*dto.MutasiRuaDetailGetByIDResponse, error)
	Create(ctx *abstraction.Context, payload *dto.MutasiRuaDetailCreateRequest) (*dto.MutasiRuaDetailCreateResponse, error)
	Update(ctx *abstraction.Context, payload *dto.MutasiRuaDetailUpdateRequest) (*dto.MutasiRuaDetailUpdateResponse, error)
	Delete(ctx *abstraction.Context, payload *dto.MutasiRuaDetailDeleteRequest) (*dto.MutasiRuaDetailDeleteResponse, error)
}

func NewService(f *factory.Factory) *service {
	repository := f.MutasiRuaDetailRepository
	MutasiRuaRepository := f.MutasiRuaRepository
	formatterBridgesRepo := f.FormatterBridgesRepository
	formatterDetailRepo := f.FormatterDetailRepository
	validationDetailRepo := f.ValidationDetailRepository
	formatterRepo := f.FormatterRepository
	trialBalanceDetailRepo := f.TrialBalanceDetailRepository
	db := f.Db
	return &service{
		Repository:                 repository,
		MutasiRuaRepository:         MutasiRuaRepository,
		FormatterBridgesRepository: formatterBridgesRepo,
		FormatterDetailRepository:  formatterDetailRepo,
		ValidationDetailRepository: validationDetailRepo,
		FormatterRepository:          formatterRepo,
		TrialBalanceDetailRepository: trialBalanceDetailRepo,
		Db:                         db,
	}
}

func (s *service) Find(ctx *abstraction.Context, payload *dto.MutasiRuaDetailGetRequest) (*dto.MutasiRuaDetailGetResponse, error) {
	mutasi, err := s.MutasiRuaRepository.FindByID(ctx, payload.MutasiRuaID)
	if err != nil {
		return &dto.MutasiRuaDetailGetResponse{}, helper.ErrorHandler(err)
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
		return &dto.MutasiRuaDetailGetResponse{}, response.ErrorBuilder(&response.ErrorConstant.InternalServerError, errors.New("not allowed"))
	}

	criteriaFormatterBridges := model.FormatterBridgesFilterModel{}
	criteriaFormatterBridges.TrxRefID = &mutasi.ID
	src := "MUTASI-RUA"
	criteriaFormatterBridges.Source = &src
	id := "id"
	asc := "ASC"
	pagingFB := abstraction.Pagination{
		SortBy: &id,
		Sort:   &asc,
	}
	formatterBridges, _, err := s.FormatterBridgesRepository.Find(ctx, &criteriaFormatterBridges, &pagingFB)
	if err != nil {
		return &dto.MutasiRuaDetailGetResponse{}, helper.ErrorHandler(err)
	}
	for _, v := range *formatterBridges {
		payload.FormatterBridgesID = &v.ID
		data, err := s.Repository.FindWithFormatter(ctx, &payload.MutasiRuaDetailFilterModel)
		if err != nil {
			return &dto.MutasiRuaDetailGetResponse{}, helper.ErrorHandler(err)
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
		var dataControl model.MutasiRuaDetailEntityModel
		var dataControlMIA []model.MutasiRuaDetailEntityModel
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
				criteriaMFD := model.MutasiRuaDetailFilterModel{}
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
					if cntrl.Coa2 != "" {
						hasilCoa2 := 0.0
						splitCoa2 := strings.Split(cntrl.Coa2, ".")
						for _, coa2 := range splitCoa2 {
							tmpCodeCoa1 := coa2
							criteriaSummaryCoa1 := model.TrialBalanceDetailFilterModel{}
							criteriaSummaryCoa1.Code = &tmpCodeCoa1
							criteriaSummaryCoa1.FormatterBridgesID = &fmtBridgesTB.ID

							summaryCoa2, err := s.TrialBalanceDetailRepository.FindSummaryTb(ctx, &criteriaSummaryCoa1)
							if err != nil {
								return nil, err
							}
							hasilCoa2 -= *summaryCoa2.AmountAfterAje

						}
						dataControls := *dataMFD.Additions - hasilCoa2
						dataControl.Control = &dataControls
					} else {
						dataControls := *dataMFD.Additions - *summaryCoa1.AmountBeforeAje
						dataControl.Control = &dataControls
					}

				case "deductions":
					dataControls := *dataMFD.Deductions - *summaryCoa1.AmountBeforeAje
					dataControl.Control = &dataControls
				}

				dataControlMIA = append(dataControlMIA, dataControl)

			}
		}
		switch v.Formatter.FormatterFor {
		case "MUTASI-RUA-COST":
			mutasi.MutasiRuaCostDetail = *data
			mutasi.ControlMIACost = dataControlMIA
		case "MUTASI-RUA-ACCUMULATED-DEPRECATION":
			mutasi.MutasiRuaADDetail = *data
			mutasi.ControlMIAD = dataControlMIA
		case "MUTASI-DETAIL-PENGURANGAN":
			mutasi.MutasiDetailPengurangan = *data
			// mutasi.ControlMIAD = dataControlMIA
		}
	}
	result := &dto.MutasiRuaDetailGetResponse{
		Datas: *mutasi,
	}
	return result, nil
}

func (s *service) FindByID(ctx *abstraction.Context, payload *dto.MutasiRuaDetailGetByIDRequest) (*dto.MutasiRuaDetailGetByIDResponse, error) {
	data, err := s.Repository.FindByID(ctx, &payload.ID)
	if err != nil {
		return &dto.MutasiRuaDetailGetByIDResponse{}, helper.ErrorHandler(err)
	}

	fmtBridges, err := s.FormatterBridgesRepository.FindByID(ctx, &data.FormatterBridgesID)
	if err != nil {
		return nil, helper.ErrorHandler(err)
	}

	mruaData, err := s.MutasiRuaRepository.FindByID(ctx, &fmtBridges.TrxRefID)
	if err != nil {
		return nil, helper.ErrorHandler(err)
	}

	allowed := helper.CompanyValidation(ctx.Auth.ID, mruaData.CompanyID)
	if !allowed {
		return nil, response.ErrorBuilder(&response.ErrorConstant.UnauthorizedAccess, errors.New("not allowed"))
	}

	result := &dto.MutasiRuaDetailGetByIDResponse{
		MutasiRuaDetailEntityModel: *data,
	}
	return result, nil
}

func (s *service) Create(ctx *abstraction.Context, payload *dto.MutasiRuaDetailCreateRequest) (*dto.MutasiRuaDetailCreateResponse, error) {
	var data model.MutasiRuaDetailEntityModel

	fmtBridges, err := s.FormatterBridgesRepository.FindByID(ctx, &payload.FormatterBridgesID)
	if err != nil {
		return nil, helper.ErrorHandler(err)
	}

	mruaData, err := s.MutasiRuaRepository.FindByID(ctx, &fmtBridges.TrxRefID)
	if err != nil {
		return nil, helper.ErrorHandler(err)
	}

	if mruaData.Status != 1 {
		return nil, response.ErrorBuilder(&response.ErrorConstant.DataValidated, err)
	}

	allowed := helper.CompanyValidation(ctx.Auth.ID, mruaData.CompanyID)
	if !allowed {
		return nil, response.ErrorBuilder(&response.ErrorConstant.UnauthorizedAccess, errors.New("not allowed"))
	}

	if err := trxmanager.New(s.Db).WithTrx(ctx, func(ctx *abstraction.Context) error {
		data.Context = ctx
		data.MutasiRuaDetailEntity = payload.MutasiRuaDetailEntity

		result, err := s.Repository.Create(ctx, &data)
		if err != nil {
			return helper.ErrorHandler(err)
		}
		data = *result
		return nil
	}); err != nil {
		return &dto.MutasiRuaDetailCreateResponse{}, err
	}

	result := &dto.MutasiRuaDetailCreateResponse{
		MutasiRuaDetailEntityModel: data,
	}
	return result, nil
}

func (s *service) Update(ctx *abstraction.Context, payload *dto.MutasiRuaDetailUpdateRequest) (*dto.MutasiRuaDetailUpdateResponse, error) {
	var data model.MutasiRuaDetailEntityModel
	if err := trxmanager.New(s.Db).WithTrx(ctx, func(ctx *abstraction.Context) error {
		mruaDetail, err := s.Repository.FindByID(ctx, &payload.ID)
		if err != nil {
			return helper.ErrorHandler(err)
		}

		formatterBridgesData, err := s.FormatterBridgesRepository.FindByID(ctx, &mruaDetail.FormatterBridgesID)
		if err != nil {
			return helper.ErrorHandler(err)
		}

		mruaData, err := s.MutasiRuaRepository.FindByID(ctx, &formatterBridgesData.TrxRefID)
		if err != nil {
			return helper.ErrorHandler(err)
		}

		allowed := helper.CompanyValidation(ctx.Auth.ID, mruaData.CompanyID)
		if !allowed {
			return response.ErrorBuilder(&response.ErrorConstant.UnauthorizedAccess, errors.New("not allowed"))
		}

		//validasi data bukan isTotal atau autosummary
		criteriaFormatterValidasi := model.FormatterDetailFilterModel{}
		criteriaFormatterValidasi.FormatterID = &formatterBridgesData.FormatterID
		criteriaFormatterValidasi.Code = &mruaDetail.Code

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
		endingBalance := (*payload.BeginningBalance) + (*payload.AcquisitionOfSubsidiary) + (*payload.Additions) - (*payload.Deductions) + (*payload.Reclassification) + (*payload.Remeasurement)
		data.MutasiRuaDetailEntity = model.MutasiRuaDetailEntity{
			BeginningBalance:        payload.BeginningBalance,
			AcquisitionOfSubsidiary: payload.AcquisitionOfSubsidiary,
			Additions:               payload.Additions,
			Deductions:              payload.Deductions,
			Reclassification:        payload.Reclassification,
			Remeasurement:           payload.Remeasurement,
			EndingBalance:           &endingBalance,
			// Control:                 payload.Control,
		}

		result, err := s.Repository.Update(ctx, &payload.ID, &data)
		if err != nil {
			return helper.ErrorHandler(err)
		}

		criteriaFormatterDetailSumData := model.FormatterBridgesFilterModel{}
		mrua := "MUTASI-RUA"
		criteriaFormatterDetailSumData.Source = &mrua
		criteriaFormatterDetailSumData.TrxRefID = &formatterBridgesData.TrxRefID

		formatterDetailSumData, err := s.FormatterBridgesRepository.FindSummary(ctx, &criteriaFormatterDetailSumData)
		if err != nil {
			return helper.ErrorHandler(err)
		}

		criteriaFormatterDetail := model.MutasiRuaDetailFilterModel{}
		criteriaFormatterDetail.MutasiRuaID = &formatterBridgesData.TrxRefID

		for _, v := range *formatterDetailSumData {
			criteriaMutasiDetail := model.MutasiRuaDetailFilterModel{}
			// criteriaMutasiDetail.FormatterBridgesID = &formatterBridgesData.ID
			criteriaMutasiDetail.Code = &v.Code
			criteriaMutasiDetail.MutasiRuaID = &formatterBridgesData.TrxRefID

			if v.AutoSummary != nil && *v.AutoSummary {
				mruadetailsum, _, err := s.Repository.Find(ctx, &criteriaMutasiDetail, &abstraction.Pagination{})
				if err != nil {
					return helper.ErrorHandler(err)
				}

				for _, a := range *mruadetailsum {
					sumMFA, err := s.Repository.FindSummary(ctx, &v.FormatterID, &a.FormatterBridgesID, &v.SortID)
					if err != nil {
						return helper.ErrorHandler(err)
					}
					updateSummary := model.MutasiRuaDetailEntityModel{
						MutasiRuaDetailEntity: model.MutasiRuaDetailEntity{
							BeginningBalance:        sumMFA.BeginningBalance,
							AcquisitionOfSubsidiary: sumMFA.AcquisitionOfSubsidiary,
							Additions:               sumMFA.Additions,
							Deductions:              sumMFA.Deductions,
							Reclassification:        sumMFA.Reclassification,
							Remeasurement:           sumMFA.Remeasurement,
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
				tmpString := []string{"BeginningBalance", "AcquisitionOfSubsidiary", "Additions", "Deductions", "Reclassification", "Remeasurement", "EndingBalance"}
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
						criteriaSumRUA := model.MutasiRuaDetailFilterModel{}
						criteriaSumRUA.Code = &vMatch
						criteriaSumRUA.MutasiRuaID = &mruaData.ID
						sumFA, err := s.Repository.FindTotal(ctx, &criteriaSumRUA)
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
						if tipe == "Remeasurement" && sumFA.Remeasurement != nil {
							angka = *sumFA.Remeasurement
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

				for _, a := range *mdtadetailsum {
					updateSummary := model.MutasiRuaDetailEntityModel{
						MutasiRuaDetailEntity: model.MutasiRuaDetailEntity{
							BeginningBalance:        tmpTotalFl["BeginningBalance"],
							AcquisitionOfSubsidiary: tmpTotalFl["AcquisitionOfSubsidiary"],
							Additions:               tmpTotalFl["Additions"],
							Deductions:              tmpTotalFl["Deductions"],
							Reclassification:        tmpTotalFl["Reclassification"],
							Remeasurement:           tmpTotalFl["Remeasurement"],
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
		criteriaValidation.CompanyID = &mruaData.CompanyID
		criteriaValidation.Period = &mruaData.Period
		criteriaValidation.Versions = &mruaData.Versions
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
		return &dto.MutasiRuaDetailUpdateResponse{}, err
	}
	result := &dto.MutasiRuaDetailUpdateResponse{
		MutasiRuaDetailEntityModel: data,
	}
	return result, nil
}

func (s *service) Delete(ctx *abstraction.Context, payload *dto.MutasiRuaDetailDeleteRequest) (*dto.MutasiRuaDetailDeleteResponse, error) {
	var data model.MutasiRuaDetailEntityModel
	if err := trxmanager.New(s.Db).WithTrx(ctx, func(ctx *abstraction.Context) error {
		existing, err := s.Repository.FindByID(ctx, &payload.ID)
		if err != nil {
			return helper.ErrorHandler(err)
		}

		fmtBridges, err := s.FormatterBridgesRepository.FindByID(ctx, &existing.FormatterBridgesID)
		if err != nil {
			return helper.ErrorHandler(err)
		}

		mruaData, err := s.MutasiRuaRepository.FindByID(ctx, &fmtBridges.TrxRefID)
		if err != nil {
			return helper.ErrorHandler(err)
		}

		allowed := helper.CompanyValidation(ctx.Auth.ID, mruaData.CompanyID)
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
		return &dto.MutasiRuaDetailDeleteResponse{}, err
	}
	result := &dto.MutasiRuaDetailDeleteResponse{
		// MutasiRuaDetailEntityModel: data,
	}
	return result, nil
}
