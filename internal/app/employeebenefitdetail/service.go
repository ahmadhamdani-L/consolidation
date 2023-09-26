package employeebenefitdetail

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
	Repository                   repository.EmployeeBenefitDetail
	EmployeeBenefitRepository    repository.EmployeeBenefit
	FormatterBridgesRepository   repository.FormatterBridges
	FormatterDetailRepository    repository.FormatterDetail
	ValidationDetailRepository   repository.ValidationDetail
	FormatterRepository          repository.Formatter
	TrialBalanceDetailRepository repository.TrialBalanceDetail
	Db                           *gorm.DB
}

type Service interface {
	Find(ctx *abstraction.Context, payload *dto.EmployeeBenefitDetailGetRequest) (*dto.EmployeeBenefitDetailGetResponse, error)
	FindByID(ctx *abstraction.Context, payload *dto.EmployeeBenefitDetailGetByIDRequest) (*dto.EmployeeBenefitDetailGetByIDResponse, error)
	Create(ctx *abstraction.Context, payload *dto.EmployeeBenefitDetailCreateRequest) (*dto.EmployeeBenefitDetailCreateResponse, error)
	Update(ctx *abstraction.Context, payload *dto.EmployeeBenefitDetailUpdateRequest) (*dto.EmployeeBenefitDetailUpdateResponse, error)
	Delete(ctx *abstraction.Context, payload *dto.EmployeeBenefitDetailDeleteRequest) (*dto.EmployeeBenefitDetailDeleteResponse, error)
}

func NewService(f *factory.Factory) *service {
	repository := f.EmployeeBenefitDetailRepository
	employeeBenefitRepo := f.EmployeeBenefitRepository
	formatterBridgesRepo := f.FormatterBridgesRepository
	formatterDetailRepo := f.FormatterDetailRepository
	validationDetailRepo := f.ValidationDetailRepository
	formatterRepo := f.FormatterRepository
	trialBalanceDetailRepo := f.TrialBalanceDetailRepository
	db := f.Db
	return &service{
		Repository:                   repository,
		EmployeeBenefitRepository:    employeeBenefitRepo,
		FormatterBridgesRepository:   formatterBridgesRepo,
		FormatterDetailRepository:    formatterDetailRepo,
		ValidationDetailRepository:   validationDetailRepo,
		FormatterRepository:          formatterRepo,
		TrialBalanceDetailRepository: trialBalanceDetailRepo,
		Db:                           db,
	}
}

func (s *service) Find(ctx *abstraction.Context, payload *dto.EmployeeBenefitDetailGetRequest) (*dto.EmployeeBenefitDetailGetResponse, error) {
	mutasi, err := s.EmployeeBenefitRepository.FindByID(ctx, payload.EmployeeBenefitID)
	if err != nil {
		return &dto.EmployeeBenefitDetailGetResponse{}, helper.ErrorHandler(err)
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
		return &dto.EmployeeBenefitDetailGetResponse{}, response.ErrorBuilder(&response.ErrorConstant.InternalServerError, errors.New("not allowed"))
	}

	criteriaFormatterBridges := model.FormatterBridgesFilterModel{}
	criteriaFormatterBridges.TrxRefID = &mutasi.ID
	src := "EMPLOYEE-BENEFIT"
	criteriaFormatterBridges.Source = &src
	id := "id"
	asc := "ASC"
	pagingFB := abstraction.Pagination{
		SortBy: &id,
		Sort:   &asc,
	}
	formatterBridges, _, err := s.FormatterBridgesRepository.Find(ctx, &criteriaFormatterBridges, &pagingFB)
	if err != nil {
		return &dto.EmployeeBenefitDetailGetResponse{}, helper.ErrorHandler(err)
	}
	idP1 := 0
	for _, v := range *formatterBridges {
		payload.FormatterBridgesID = &v.ID
		data, err := s.Repository.FindWithFormatter(ctx, &payload.EmployeeBenefitDetailFilterModel)
		if err != nil {
			return &dto.EmployeeBenefitDetailGetResponse{}, helper.ErrorHandler(err)
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
		var dataControl model.EmployeeBenefitDetailEntityModel
		var dataControlMIA []model.EmployeeBenefitDetailEntityModel
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
				criteriaMFD := model.EmployeeBenefitDetailFilterModel{}
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
						test := 3
						temp := idP1 + test
						splitCoa2 := strings.Split(cntrl.Coa2, ".")
						for _, coa2 := range splitCoa2 {
							if coa2 == "PENGHASILAN_KOMPREHENSIF_LAIN" {
								temp = temp-1
							}
							criteriaMFD := model.EmployeeBenefitDetailFilterModel{}
							criteriaMFD.Code = &coa2
							criteriaMFD.FormatterBridgesID = &temp
							dataMFDs, err := s.Repository.FindByCriteria(ctx, &criteriaMFD) // data Aging Utang Piutang Detail
							if err != nil || dataMFD.ID == 0 {
								return nil, err
							}
							hasilCoa2 = *dataMFDs.Amount

						}
						dataControls := *dataMFD.Amount - hasilCoa2
						dataControl.Amount = &dataControls
					} else {
						angka := 0.0
						if dataMFD.Amount == nil {
							dataMFD.Amount = &angka
						}
						if summaryCoa1.AmountBeforeAje == nil {
							summaryCoa1.AmountBeforeAje = &angka
						}
						dataControls := *dataMFD.Amount - *summaryCoa1.AmountBeforeAje
						dataControl.Amount = &dataControls
					}
				}
				dataControlMIA = append(dataControlMIA, dataControl)

			}
		}
		switch v.Formatter.FormatterFor {
		case "EMPLOYEE-BENEFIT-ASUMSI":
			mutasi.EmployeeBenefitDetailAsumsi = *data
		case "EMPLOYEE-BENEFIT-REKONSILIASI":
			mutasi.EmployeeBenefitDetailRekonsiliasi = *data
		case "EMPLOYEE-BENEFIT-RINCIAN-LAPORAN":
			mutasi.EmployeeBenefitDetailRincianLaporan = *data
			mutasi.ControllRincianLaporan = dataControlMIA
		case "EMPLOYEE-BENEFIT-RINCIAN-EKUITAS":
			mutasi.EmployeeBenefitDetailRincianEkuitas = *data
			mutasi.ControllRincianEkuitas  = dataControlMIA
		case "EMPLOYEE-BENEFIT-MUTASI":
			mutasi.EmployeeBenefitDetailMutasi = *data
			mutasi.ControllMutasi = dataControlMIA
		case "EMPLOYEE-BENEFIT-INFORMASI":
			mutasi.EmployeeBenefitDetailInformasi = *data
		case "EMPLOYEE-BENEFIT-ANALISIS":
			mutasi.EmployeeBenefitDetailAnalisis = *data
		}
		idP1 = v.ID
	}
	result := &dto.EmployeeBenefitDetailGetResponse{
		Datas: *mutasi,
	}
	return result, nil
}

func (s *service) FindByID(ctx *abstraction.Context, payload *dto.EmployeeBenefitDetailGetByIDRequest) (*dto.EmployeeBenefitDetailGetByIDResponse, error) {
	data, err := s.Repository.FindByID(ctx, &payload.ID)
	if err != nil {
		return &dto.EmployeeBenefitDetailGetByIDResponse{}, helper.ErrorHandler(err)
	}

	fmtBridges, err := s.FormatterBridgesRepository.FindByID(ctx, &data.FormatterBridgesID)
	if err != nil {
		return nil, helper.ErrorHandler(err)
	}

	aup, err := s.EmployeeBenefitRepository.FindByID(ctx, &fmtBridges.TrxRefID)
	if err != nil {
		return nil, helper.ErrorHandler(err)
	}

	allowed := helper.CompanyValidation(ctx.Auth.ID, aup.CompanyID)
	if !allowed {
		return nil, response.ErrorBuilder(&response.ErrorConstant.UnauthorizedAccess, errors.New("Not Allowed"))
	}

	result := &dto.EmployeeBenefitDetailGetByIDResponse{
		EmployeeBenefitDetailEntityModel: *data,
	}
	return result, nil
}

func (s *service) Create(ctx *abstraction.Context, payload *dto.EmployeeBenefitDetailCreateRequest) (*dto.EmployeeBenefitDetailCreateResponse, error) {
	var data model.EmployeeBenefitDetailEntityModel

	fmtBridges, err := s.FormatterBridgesRepository.FindByID(ctx, &payload.FormatterBridgesID)
	if err != nil {
		return nil, err
	}

	eb, err := s.EmployeeBenefitRepository.FindByID(ctx, &fmtBridges.TrxRefID)
	if err != nil {
		return nil, err
	}

	allowed := helper.CompanyValidation(ctx.Auth.ID, eb.CompanyID)
	if !allowed {
		return nil, response.ErrorBuilder(&response.ErrorConstant.UnauthorizedAccess, errors.New("Not Allowed"))
	}

	if err := trxmanager.New(s.Db).WithTrx(ctx, func(ctx *abstraction.Context) error {
		data.Context = ctx
		data.EmployeeBenefitDetailEntity = payload.EmployeeBenefitDetailEntity

		result, err := s.Repository.Create(ctx, &data)
		if err != nil {
			return err
		}
		data = *result
		return nil
	}); err != nil {
		return &dto.EmployeeBenefitDetailCreateResponse{}, err
	}

	result := &dto.EmployeeBenefitDetailCreateResponse{
		EmployeeBenefitDetailEntityModel: data,
	}
	return result, nil
}

func (s *service) Update(ctx *abstraction.Context, payload *dto.EmployeeBenefitDetailUpdateRequest) (*dto.EmployeeBenefitDetailUpdateResponse, error) {
	var data model.EmployeeBenefitDetailEntityModel
	if err := trxmanager.New(s.Db).WithTrx(ctx, func(ctx *abstraction.Context) error {
		employeebenefitdetail, err := s.Repository.FindByID(ctx, &payload.ID)
		if err != nil {
			return err
		}

		formatterBridgesData, err := s.FormatterBridgesRepository.FindByID(ctx, &employeebenefitdetail.FormatterBridgesID)
		if err != nil {
			return err
		}

		ebData, err := s.EmployeeBenefitRepository.FindByID(ctx, &formatterBridgesData.TrxRefID)
		if err != nil {
			return err
		}

		if ebData.Status != 1 {
			return response.ErrorBuilder(&response.ErrorConstant.DataValidated, errors.New("Cannot update data"))
		}

		allowed := helper.CompanyValidation(ctx.Auth.ID, ebData.CompanyID)
		if !allowed {
			return response.ErrorBuilder(&response.ErrorConstant.UnauthorizedAccess, errors.New("Not Allowed"))
		}

		//validasi data bukan isTotal atau autosummary
		criteriaFormatterValidasi := model.FormatterDetailFilterModel{}
		criteriaFormatterValidasi.FormatterID = &formatterBridgesData.FormatterID
		criteriaFormatterValidasi.Code = &employeebenefitdetail.Code

		formatterValidasi, jmlData, err := s.FormatterDetailRepository.Find(ctx, &criteriaFormatterValidasi, &abstraction.Pagination{})
		if err != nil {
			return helper.ErrorHandler(err)
		}

		if jmlData.Count == 0 {
			return response.ErrorBuilder(&response.ErrorConstant.NotFound, errors.New("Data Not Found"))
		}

		for _, v := range *formatterValidasi {
			if (v.IsTotal != nil && *v.IsTotal) || (v.AutoSummary != nil && *v.AutoSummary) || (v.IsLabel != nil && *v.IsLabel) {
				return helper.ErrorHandler(err)
			}
		}

		if *employeebenefitdetail.IsValue == true && *payload.Value == "" {
			return response.ErrorBuilder(&response.ErrorConstant.BadRequest, errors.New("terdapat data yang harus diisi"))
		}

		data.Context = ctx
		data.EmployeeBenefitDetailEntity = model.EmployeeBenefitDetailEntity{}
		if *employeebenefitdetail.IsValue == false && payload.Amount != nil && *payload.Amount != 0 {
			data.EmployeeBenefitDetailEntity.Amount = payload.Amount
		}
		if payload.Value != nil && *payload.Value != "" {
			data.EmployeeBenefitDetailEntity.Value = *payload.Value
		}
		// data.EmployeeBenefitDetailEntity = payload.EmployeeBenefitDetailEntity

		result, err := s.Repository.Update(ctx, &payload.ID, &data)
		if err != nil {
			return err
		}

		criteriaFormatterDetailSumData := model.FormatterBridgesFilterModel{}
		eb := "EMPLOYEE-BENEFIT"
		criteriaFormatterDetailSumData.Source = &eb
		criteriaFormatterDetailSumData.TrxRefID = &formatterBridgesData.TrxRefID

		formatterDetailSumData, err := s.FormatterBridgesRepository.FindSummary(ctx, &criteriaFormatterDetailSumData)
		if err != nil {
			return helper.ErrorHandler(err)
		}

		criteriaFormatterDetail := model.EmployeeBenefitDetailFilterModel{}
		// criteriaFormatterDetail.FormatterBridgesID = &formatterBridgesData.ID
		criteriaFormatterDetail.EmployeeBenefitID = &formatterBridgesData.TrxRefID

		formatterDetailData, err := s.Repository.FindWithFormatter(ctx, &criteriaFormatterDetail)
		if err != nil {
			return helper.ErrorHandler(err)
		}

		for _, v := range *formatterDetailSumData {
			criteriaEmployeeBenefitDetail := model.EmployeeBenefitDetailFilterModel{}
			// criteriaEmployeeBenefitDetail.FormatterBridgesID = &employeebenefitdetail.FormatterBridgesID
			criteriaEmployeeBenefitDetail.Code = &v.Code
			criteriaEmployeeBenefitDetail.EmployeeBenefitID = &formatterBridgesData.TrxRefID

			if v.AutoSummary != nil && *v.AutoSummary == true {
				filterHelperFormatterBridges := model.FormatterBridgesFilterModel{}
				filterHelperFormatterBridges.FormatterID = &v.FormatterID
				filterHelperFormatterBridges.Source = &eb
				filterHelperFormatterBridges.TrxRefID = &formatterBridgesData.TrxRefID
				helperFormatterBridges, _, err := s.FormatterBridgesRepository.Find(ctx, &filterHelperFormatterBridges, &abstraction.Pagination{})
				if err != nil {
					return helper.ErrorHandler(err)
				}
				for _, hv := range *helperFormatterBridges {
					criteriaEmployeeBenefitDetail.FormatterBridgesID = &hv.ID
				}

				ebdetailsum, _, err := s.Repository.Find(ctx, &criteriaEmployeeBenefitDetail, &abstraction.Pagination{})
				if err != nil {
					return helper.ErrorHandler(err)
				}

				for _, a := range *ebdetailsum {
					sumAUP, err := s.Repository.FindSummary(ctx, &v.FormatterID, &a.FormatterBridgesID, &v.SortID)
					if err != nil {
						return helper.ErrorHandler(err)
					}
					updateSummary := model.EmployeeBenefitDetailEntityModel{
						EmployeeBenefitDetailEntity: model.EmployeeBenefitDetailEntity{
							Amount: sumAUP.Amount,
						},
					}
					_, err = s.Repository.Update(ctx, &a.ID, &updateSummary)
					if err != nil {
						return helper.ErrorHandler(err)
					}
				}
			}

			if v.IsTotal != nil && *v.IsTotal == true {
				if v.Code == "CONTROL_1" || v.Code == "CONTROL_2" {
					continue
				}
				formula := v.FxSummary
				// parameterFormula := make(map[string]interface{})
				tmpString := []string{"Amount"}
				tmpTotalStr := make(map[string]string)
				tmpTotalFl := make(map[string]*float64)
				for _, vFormatterDetail := range *formatterDetailData {
					if strings.Contains(formula, strings.Trim(strings.ToUpper(vFormatterDetail.Code), " ")) {
						tmpTotalStr["Amount"] = strings.Replace(formula, vFormatterDetail.Code, fmt.Sprintf("%f", *vFormatterDetail.Amount), -1)
					}
				}

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

				updateSummary := model.EmployeeBenefitDetailEntityModel{
					EmployeeBenefitDetailEntity: model.EmployeeBenefitDetailEntity{
						Amount: tmpTotalFl["Amount"],
					},
				}
				ebsum, _, err := s.Repository.Find(ctx, &criteriaEmployeeBenefitDetail, &abstraction.Pagination{})
				if err != nil {
					return helper.ErrorHandler(err)
				}
				for _, vEmployee := range *ebsum {
					_, err = s.Repository.Update(ctx, &vEmployee.ID, &updateSummary)
					if err != nil {
						return helper.ErrorHandler(err)
					}
				}
			}
		}

		criteriaValidation := model.ValidationDetailFilterModel{}
		criteriaValidation.CompanyID = &ebData.CompanyID
		criteriaValidation.Period = &ebData.Period
		criteriaValidation.Versions = &ebData.Versions
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
		return &dto.EmployeeBenefitDetailUpdateResponse{}, err
	}
	result := &dto.EmployeeBenefitDetailUpdateResponse{
		EmployeeBenefitDetailEntityModel: data,
	}
	return result, nil
}

func (s *service) Delete(ctx *abstraction.Context, payload *dto.EmployeeBenefitDetailDeleteRequest) (*dto.EmployeeBenefitDetailDeleteResponse, error) {
	var data model.EmployeeBenefitDetailEntityModel
	if err := trxmanager.New(s.Db).WithTrx(ctx, func(ctx *abstraction.Context) error {
		existing, err := s.Repository.FindByID(ctx, &payload.ID)
		if err != nil {
			return helper.ErrorHandler(err)
		}

		fmtBridges, err := s.FormatterBridgesRepository.FindByID(ctx, &existing.FormatterBridgesID)
		if err != nil {
			return helper.ErrorHandler(err)
		}

		aup, err := s.EmployeeBenefitRepository.FindByID(ctx, &fmtBridges.TrxRefID)
		if err != nil {
			return helper.ErrorHandler(err)
		}

		allowed := helper.CompanyValidation(ctx.Auth.ID, aup.CompanyID)
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
		return &dto.EmployeeBenefitDetailDeleteResponse{}, helper.ErrorHandler(err)
	}
	result := &dto.EmployeeBenefitDetailDeleteResponse{
		// EmployeeBenefitDetailEntityModel: data,
	}
	return result, nil
}
