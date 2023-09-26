package mutasipersediaandetail

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
	Repository                   repository.MutasiPersediaanDetail
	MutasiPersediaanRepository   repository.MutasiPersediaan
	FormatterBridgesRepository   repository.FormatterBridges
	FormatterDetailRepository    repository.FormatterDetail
	ValidationDetailRepository   repository.ValidationDetail
	FormatterRepository          repository.Formatter
	TrialBalanceDetailRepository repository.TrialBalanceDetail
	Db                           *gorm.DB
}

type Service interface {
	Find(ctx *abstraction.Context, payload *dto.MutasiPersediaanDetailGetRequest) (*dto.MutasiPersediaanDetailGetResponse, error)
	FindByID(ctx *abstraction.Context, payload *dto.MutasiPersediaanDetailGetByIDRequest) (*dto.MutasiPersediaanDetailGetByIDResponse, error)
	Create(ctx *abstraction.Context, payload *dto.MutasiPersediaanDetailCreateRequest) (*dto.MutasiPersediaanDetailCreateResponse, error)
	Update(ctx *abstraction.Context, payload *dto.MutasiPersediaanDetailUpdateRequest) (*dto.MutasiPersediaanDetailUpdateResponse, error)
	Delete(ctx *abstraction.Context, payload *dto.MutasiPersediaanDetailDeleteRequest) (*dto.MutasiPersediaanDetailDeleteResponse, error)
}

func NewService(f *factory.Factory) *service {
	repository := f.MutasiPersediaanDetailRepository
	mutasiPersediaanRepo := f.MutasiPersediaanRepository
	formatterBridgesRepo := f.FormatterBridgesRepository
	formatterDetailRepo := f.FormatterDetailRepository
	validationDetailRepo := f.ValidationDetailRepository
	formatterRepo := f.FormatterRepository
	trialBalanceDetailRepo := f.TrialBalanceDetailRepository
	db := f.Db
	return &service{
		Repository:                   repository,
		MutasiPersediaanRepository:   mutasiPersediaanRepo,
		FormatterBridgesRepository:   formatterBridgesRepo,
		FormatterDetailRepository:    formatterDetailRepo,
		ValidationDetailRepository:   validationDetailRepo,
		FormatterRepository:          formatterRepo,
		TrialBalanceDetailRepository: trialBalanceDetailRepo,
		Db:                           db,
	}
}

func (s *service) Find(ctx *abstraction.Context, payload *dto.MutasiPersediaanDetailGetRequest) (*dto.MutasiPersediaanDetailGetResponse, error) {
	mutasi, err := s.MutasiPersediaanRepository.FindByID(ctx, payload.MutasiPersediaanID)
	if err != nil {
		return &dto.MutasiPersediaanDetailGetResponse{}, helper.ErrorHandler(err)
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
		return &dto.MutasiPersediaanDetailGetResponse{}, response.ErrorBuilder(&response.ErrorConstant.InternalServerError, errors.New("not allowed"))
	}

	criteriaFormatterBridges := model.FormatterBridgesFilterModel{}
	criteriaFormatterBridges.TrxRefID = &mutasi.ID
	src := "MUTASI-PERSEDIAAN"
	criteriaFormatterBridges.Source = &src
	id := "id"
	asc := "ASC"
	pagingFB := abstraction.Pagination{
		SortBy: &id,
		Sort:   &asc,
	}
	formatterBridges, _, err := s.FormatterBridgesRepository.Find(ctx, &criteriaFormatterBridges, &pagingFB)
	if err != nil {
		return &dto.MutasiPersediaanDetailGetResponse{}, helper.ErrorHandler(err)
	}
	idP1 := 0
	for _, v := range *formatterBridges {
		payload.FormatterBridgesID = &v.ID
		data, err := s.Repository.FindWithFormatter(ctx, &payload.MutasiPersediaanDetailFilterModel)
		if err != nil {
			return &dto.MutasiPersediaanDetailGetResponse{}, helper.ErrorHandler(err)
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
		var dataControl model.MutasiPersediaanDetailEntityModel
		var dataControlPERSEDIAAN []model.MutasiPersediaanDetailEntityModel
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
				criteriaMFD := model.MutasiPersediaanDetailFilterModel{}
				criteriaMFD.Code = &vFmt.Code
				criteriaMFD.FormatterBridgesID = &v.ID
				dataMFD, err := s.Repository.FindByCriteria(ctx, &criteriaMFD) // data Aging Utang Piutang Detail
				if err != nil || dataMFD.ID == 0 {
					return nil, err
				}

				switch strings.ToLower(splitControllerCommand[1]) {
				case "amount":
					if cntrl.Coa2 != "" {
						hasilCoa2 := 0.0
						splitCoa2 := strings.Split(cntrl.Coa2, ".")
						for _, coa2 := range splitCoa2 {
							criteriaMFD := model.MutasiPersediaanDetailFilterModel{}
							criteriaMFD.Code = &coa2
							criteriaMFD.FormatterBridgesID = &idP1
							dataMFDs, err := s.Repository.FindByCriteria(ctx, &criteriaMFD) // data Aging Utang Piutang Detail
							if err != nil || dataMFD.ID == 0 {
								return nil, err
							}
							hasilCoa2 = *dataMFDs.Amount

						}
						dataControls := *dataMFD.Amount - hasilCoa2
						dataControl.Amount = &dataControls
					} else {
						dataControls := *dataMFD.Amount - *summaryCoa1.AmountBeforeAje
						dataControl.Amount = &dataControls
					}
				}

				dataControlPERSEDIAAN = append(dataControlPERSEDIAAN, dataControl)

			}
		}
		idP1 = v.ID
		switch v.Formatter.FormatterFor {
		case "MUTASI-PERSEDIAAN":
			mutasi.MutasiPersediaanDetail = *data
			mutasi.ControlPersediaan = dataControlPERSEDIAAN
		case "MUTASI-CADANGAN-PENGHAPUSAN-PERSEDIAAN":
			mutasi.MutasiCadanganPenghapusanpersediaan = *data
			mutasi.ControlPersediaanPenghapusan = dataControlPERSEDIAAN
		}
	}
	result := &dto.MutasiPersediaanDetailGetResponse{
		Datas: *mutasi,
	}
	return result, nil
}

func (s *service) FindByID(ctx *abstraction.Context, payload *dto.MutasiPersediaanDetailGetByIDRequest) (*dto.MutasiPersediaanDetailGetByIDResponse, error) {
	data, err := s.Repository.FindByID(ctx, &payload.ID)
	if err != nil {
		return &dto.MutasiPersediaanDetailGetByIDResponse{}, helper.ErrorHandler(err)
	}

	fmtBridges, err := s.FormatterBridgesRepository.FindByID(ctx, &data.FormatterBridgesID)
	if err != nil {
		return nil, helper.ErrorHandler(err)
	}

	mpData, err := s.MutasiPersediaanRepository.FindByID(ctx, &fmtBridges.TrxRefID)
	if err != nil {
		return nil, helper.ErrorHandler(err)
	}

	allowed := helper.CompanyValidation(ctx.Auth.ID, mpData.CompanyID)
	if !allowed {
		return nil, response.ErrorBuilder(&response.ErrorConstant.UnauthorizedAccess, errors.New("Not Allowed"))
	}

	result := &dto.MutasiPersediaanDetailGetByIDResponse{
		MutasiPersediaanDetailEntityModel: *data,
	}
	return result, nil
}

func (s *service) Create(ctx *abstraction.Context, payload *dto.MutasiPersediaanDetailCreateRequest) (*dto.MutasiPersediaanDetailCreateResponse, error) {
	var data model.MutasiPersediaanDetailEntityModel

	fmtBridges, err := s.FormatterBridgesRepository.FindByID(ctx, &payload.FormatterBridgesID)
	if err != nil {
		return nil, helper.ErrorHandler(err)
	}

	mpData, err := s.MutasiPersediaanRepository.FindByID(ctx, &fmtBridges.TrxRefID)
	if err != nil {
		return nil, helper.ErrorHandler(err)
	}

	if mpData.Status != 1 {
		return nil, response.ErrorBuilder(&response.ErrorConstant.DataValidated, err)
	}

	allowed := helper.CompanyValidation(ctx.Auth.ID, mpData.CompanyID)
	if !allowed {
		return nil, response.ErrorBuilder(&response.ErrorConstant.UnauthorizedAccess, errors.New("Not Allowed"))
	}

	if err := trxmanager.New(s.Db).WithTrx(ctx, func(ctx *abstraction.Context) error {
		data.Context = ctx
		data.MutasiPersediaanDetailEntity = payload.MutasiPersediaanDetailEntity

		result, err := s.Repository.Create(ctx, &data)
		if err != nil {
			return helper.ErrorHandler(err)
		}
		data = *result
		return nil
	}); err != nil {
		return &dto.MutasiPersediaanDetailCreateResponse{}, err
	}

	result := &dto.MutasiPersediaanDetailCreateResponse{
		MutasiPersediaanDetailEntityModel: data,
	}
	return result, nil
}

func (s *service) Update(ctx *abstraction.Context, payload *dto.MutasiPersediaanDetailUpdateRequest) (*dto.MutasiPersediaanDetailUpdateResponse, error) {
	var data model.MutasiPersediaanDetailEntityModel
	if err := trxmanager.New(s.Db).WithTrx(ctx, func(ctx *abstraction.Context) error {
		mpDetail, err := s.Repository.FindByID(ctx, &payload.ID)
		if err != nil {
			return helper.ErrorHandler(err)
		}

		formatterBridgesData, err := s.FormatterBridgesRepository.FindByID(ctx, &mpDetail.FormatterBridgesID)
		if err != nil {
			return helper.ErrorHandler(err)
		}

		mpData, err := s.MutasiPersediaanRepository.FindByID(ctx, &formatterBridgesData.TrxRefID)
		if err != nil {
			return helper.ErrorHandler(err)
		}

		allowed := helper.CompanyValidation(ctx.Auth.ID, mpData.CompanyID)
		if !allowed {
			return response.ErrorBuilder(&response.ErrorConstant.UnauthorizedAccess, errors.New("Not Allowed"))
		}

		//validasi data bukan isTotal atau autosummary
		criteriaFormatterValidasi := model.FormatterDetailFilterModel{}
		criteriaFormatterValidasi.FormatterID = &formatterBridgesData.FormatterID
		criteriaFormatterValidasi.Code = &mpDetail.Code

		formatterValidasi, jmlData, err := s.FormatterDetailRepository.Find(ctx, &criteriaFormatterValidasi, &abstraction.Pagination{})
		if err != nil {
			return helper.ErrorHandler(err)
		}

		if jmlData.Count == 0 {
			return response.ErrorBuilder(&response.ErrorConstant.NotFound, errors.New("Data not found"))
		}

		for _, v := range *formatterValidasi {
			if (v.IsTotal != nil && *v.IsTotal == true) || (v.AutoSummary != nil && *v.AutoSummary == true) {
				return response.ErrorBuilder(&response.ErrorConstant.BadRequest, errors.New("Cannot update data"))
			}
		}

		data.Context = ctx
		// data.MutasiPersediaanDetailEntity = payload.MutasiPersediaanDetailEntity
		data.MutasiPersediaanDetailEntity = model.MutasiPersediaanDetailEntity{
			Amount: payload.Amount,
		}

		result, err := s.Repository.Update(ctx, &payload.ID, &data)
		if err != nil {
			return helper.ErrorHandler(err)
		}

		criteriaFormatterDetailSumData := model.FormatterBridgesFilterModel{}
		mp := "MUTASI-PERSEDIAAN"
		criteriaFormatterDetailSumData.Source = &mp
		criteriaFormatterDetailSumData.TrxRefID = &formatterBridgesData.TrxRefID

		formatterDetailSumData, err := s.FormatterBridgesRepository.FindSummary(ctx, &criteriaFormatterDetailSumData)
		if err != nil {
			return helper.ErrorHandler(err)
		}

		criteriaFormatterDetail := model.MutasiPersediaanDetailFilterModel{}
		criteriaFormatterDetail.MutasiPersediaanID = &formatterBridgesData.TrxRefID

		formatterDetailData, err := s.Repository.FindWithFormatter(ctx, &criteriaFormatterDetail)
		if err != nil {
			return helper.ErrorHandler(err)
		}

		for _, v := range *formatterDetailSumData {
			criteriaMutasiDetail := model.MutasiPersediaanDetailFilterModel{}
			criteriaMutasiDetail.Code = &v.Code
			criteriaMutasiDetail.MutasiPersediaanID = &formatterBridgesData.TrxRefID

			if v.AutoSummary != nil && *v.AutoSummary == true {
				filterHelperFormatterBridges := model.FormatterBridgesFilterModel{}
				filterHelperFormatterBridges.FormatterID = &v.FormatterID
				filterHelperFormatterBridges.Source = &mp
				filterHelperFormatterBridges.TrxRefID = &formatterBridgesData.TrxRefID
				helperFormatterBridges, _, err := s.FormatterBridgesRepository.Find(ctx, &filterHelperFormatterBridges, &abstraction.Pagination{})
				if err != nil {
					return helper.ErrorHandler(err)
				}
				for _, hv := range *helperFormatterBridges {
					criteriaMutasiDetail.FormatterBridgesID = &hv.ID
				}

				mpdetailsum, _, err := s.Repository.Find(ctx, &criteriaMutasiDetail, &abstraction.Pagination{})
				if err != nil {
					return helper.ErrorHandler(err)
				}

				for _, a := range *mpdetailsum {
					sumMP, err := s.Repository.FindSummary(ctx, &v.FormatterID, &a.FormatterBridgesID, &v.SortID)
					if err != nil {
						return helper.ErrorHandler(err)
					}
					updateSummary := model.MutasiPersediaanDetailEntityModel{
						MutasiPersediaanDetailEntity: model.MutasiPersediaanDetailEntity{
							Amount: sumMP.Amount,
						},
					}
					_, err = s.Repository.Update(ctx, &a.ID, &updateSummary)
					if err != nil {
						return helper.ErrorHandler(err)
					}
				}
			}

			if v.IsTotal != nil && *v.IsTotal == true && v.FxSummary != "" {
				if v.Code == "CONTROL_1" || v.Code == "CONTROL_2" || v.Code == "CONTROL_3" {
					continue
				}
				formula := strings.ToUpper(v.FxSummary)
				// parameterFormula := make(map[string]interface{})
				tmpString := []string{"Amount"}
				tmpTotalStr := make(map[string]string)
				tmpTotalFl := make(map[string]*float64)
				tmpTotalStr["Amount"] = formula
				for _, vFormatterDetail := range *formatterDetailData {
					// parameterFormula[vFormatterDetail.Code] =
					if strings.Contains(formula, strings.Trim(strings.ToUpper(vFormatterDetail.Code), " ")) {
						tmpTotalStr["Amount"] = strings.Replace(tmpTotalStr["Amount"], vFormatterDetail.Code, fmt.Sprintf("%.2f", *vFormatterDetail.Amount), -1)
					}
				}
				// reg := regexp.MustCompile(`[^1234567890\.+\-\/\(\)\*\^]+`)
				reg := regexp.MustCompile(`[A-Za-z_]+|[0-9]+\d{3,}`)
				for _, str := range tmpString {
					formulaFinal := reg.ReplaceAllString(tmpTotalStr[str], "0")
					expressionFormula, err := govaluate.NewEvaluableExpression(formulaFinal)
					if err != nil {
						return err
					}
					parameters := make(map[string]interface{}, 0)
					result, err := expressionFormula.Evaluate(parameters)
					if err != nil {
						return err
					}
					tmp := result.(float64)
					tmpTotalFl[str] = &tmp
				}

				mpdetailsum, _, err := s.Repository.Find(ctx, &criteriaMutasiDetail, &abstraction.Pagination{})
				if err != nil {
					return helper.ErrorHandler(err)
				}

				for _, a := range *mpdetailsum {
					updateSummary := model.MutasiPersediaanDetailEntityModel{
						MutasiPersediaanDetailEntity: model.MutasiPersediaanDetailEntity{
							Amount: tmpTotalFl["Amount"],
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
		criteriaValidation.CompanyID = &mpData.CompanyID
		criteriaValidation.Period = &mpData.Period
		criteriaValidation.Versions = &mpData.Versions
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
		return &dto.MutasiPersediaanDetailUpdateResponse{}, err
	}
	result := &dto.MutasiPersediaanDetailUpdateResponse{
		MutasiPersediaanDetailEntityModel: data,
	}
	return result, nil
}

func (s *service) Delete(ctx *abstraction.Context, payload *dto.MutasiPersediaanDetailDeleteRequest) (*dto.MutasiPersediaanDetailDeleteResponse, error) {
	var data model.MutasiPersediaanDetailEntityModel
	if err := trxmanager.New(s.Db).WithTrx(ctx, func(ctx *abstraction.Context) error {
		existing, err := s.Repository.FindByID(ctx, &payload.ID)
		if err != nil {
			return helper.ErrorHandler(err)
		}

		fmtBridges, err := s.FormatterBridgesRepository.FindByID(ctx, &existing.FormatterBridgesID)
		if err != nil {
			return helper.ErrorHandler(err)
		}

		mpData, err := s.MutasiPersediaanRepository.FindByID(ctx, &fmtBridges.TrxRefID)
		if err != nil {
			return helper.ErrorHandler(err)
		}

		allowed := helper.CompanyValidation(ctx.Auth.ID, mpData.CompanyID)
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
		return &dto.MutasiPersediaanDetailDeleteResponse{}, err
	}
	result := &dto.MutasiPersediaanDetailDeleteResponse{
		// MutasiPersediaanDetailEntityModel: data,
	}
	return result, nil
}
