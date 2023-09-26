package agingutangpiutangdetail

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
	Repository                  repository.AgingUtangPiutangDetail
	AgingUtangPiutangRepository repository.AgingUtangPiutang
	FormatterBridgesRepository  repository.FormatterBridges
	FormatterDetailRepository   repository.FormatterDetail
	ValidationDetailRepository  repository.ValidationDetail
	FormatterRepository         repository.Formatter
	TrialBalanceDetailRepository			repository.TrialBalanceDetail
	Db                          *gorm.DB
}

type Service interface {
	Find(ctx *abstraction.Context, payload *dto.AgingUtangPiutangDetailGetRequest) (*dto.AgingUtangPiutangDetailGetResponse, error)
	FindByID(ctx *abstraction.Context, payload *dto.AgingUtangPiutangDetailGetByIDRequest) (*dto.AgingUtangPiutangDetailGetByIDResponse, error)
	Create(ctx *abstraction.Context, payload *dto.AgingUtangPiutangDetailCreateRequest) (*dto.AgingUtangPiutangDetailCreateResponse, error)
	Update(ctx *abstraction.Context, payload *dto.AgingUtangPiutangDetailUpdateRequest) (*dto.AgingUtangPiutangDetailUpdateResponse, error)
	Delete(ctx *abstraction.Context, payload *dto.AgingUtangPiutangDetailDeleteRequest) (*dto.AgingUtangPiutangDetailDeleteResponse, error)
}

func NewService(f *factory.Factory) *service {
	repository := f.AgingUtangPiutangDetailRepository
	agingUtangPiutangRepo := f.AgingUtangPiutangRepository
	formatterBridgesRepo := f.FormatterBridgesRepository
	formatterDetailRepo := f.FormatterDetailRepository
	formatterRepo := f.FormatterRepository
	validationDetailRepo := f.ValidationDetailRepository
	trialBalanceDetailRepo := f.TrialBalanceDetailRepository
	db := f.Db
	return &service{
		Repository:                  repository,
		AgingUtangPiutangRepository: agingUtangPiutangRepo,
		FormatterBridgesRepository:  formatterBridgesRepo,
		FormatterDetailRepository:   formatterDetailRepo,
		ValidationDetailRepository:  validationDetailRepo,
		FormatterRepository :        formatterRepo,
		TrialBalanceDetailRepository:trialBalanceDetailRepo,
		Db:                          db,
	}
}

func (s *service) Find(ctx *abstraction.Context, payload *dto.AgingUtangPiutangDetailGetRequest) (*dto.AgingUtangPiutangDetailGetResponse, error) {
	aup, err := s.AgingUtangPiutangRepository.FindByID(ctx, payload.AgingUtangPiutangID)
	if err != nil {
		return &dto.AgingUtangPiutangDetailGetResponse{}, helper.ErrorHandler(err)
	}
	criteriaTb := model.TrialBalanceFilterModel{}
	criteriaTb.Period = &aup.Period
	criteriaTb.CompanyID = &aup.CompanyID
	criteriaTb.Status = &aup.Status
	criteriaTb.Versions = &aup.Versions

	// ==== Aging Utang Piutang =====
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

	allowed := helper.CompanyValidation(ctx.Auth.ID, aup.CompanyID)
	if !allowed {
		return nil, response.ErrorBuilder(&response.ErrorConstant.UnauthorizedAccess, errors.New("not allowed"))
	}

	criteriaFormatterBridges := model.FormatterBridgesFilterModel{}
	criteriaFormatterBridges.TrxRefID = &aup.ID
	src := "AGING-UTANG-PIUTANG"
	criteriaFormatterBridges.Source = &src
	id := "id"
	asc := "ASC"
	pagingFB := abstraction.Pagination{
		SortBy: &id,
		Sort:   &asc,
	}

	formatterBridges, _, err := s.FormatterBridgesRepository.Find(ctx, &criteriaFormatterBridges, &pagingFB)
	if err != nil {
		return &dto.AgingUtangPiutangDetailGetResponse{}, helper.ErrorHandler(err)
	}

	for _, v := range *formatterBridges {
		payload.FormatterBridgesID = &v.ID
		data, err := s.Repository.FindWithFormatter(ctx, &payload.AgingUtangPiutangDetailFilterModel)
		if err != nil {
			return &dto.AgingUtangPiutangDetailGetResponse{}, helper.ErrorHandler(err)
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
			var dataControl model.AgingUtangPiutangDetailEntityModel
			var dataControlEmc2 []model.AgingUtangPiutangDetailEntityModel
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
					criteriaAUPD := model.AgingUtangPiutangDetailFilterModel{}
					criteriaAUPD.Code = &vFmt.Code
					criteriaAUPD.FormatterBridgesID = &v.ID
					dataAUPD, err := s.Repository.FindByCriteria(ctx, &criteriaAUPD) // data Aging Utang Piutang Detail
					if err != nil || dataAUPD.ID == 0 {
						return nil, err
					}
					
					switch strings.ToLower(splitControllerCommand[1]) {
					case "piutangusaha_3rdparty":
						dataControls := *dataAUPD.Piutangusaha3rdparty - *summaryCoa1.AmountBeforeAje
						dataControl.Piutangusaha3rdparty = &dataControls
					case "piutangusaha_berelasi":
						dataControls1 := *dataAUPD.PiutangusahaBerelasi - *summaryCoa1.AmountBeforeAje
						dataControl.PiutangusahaBerelasi = &dataControls1
					case "piutanglainshortterm_3rdparty":
						dataControl2 := *dataAUPD.Piutanglainshortterm3rdparty - *summaryCoa1.AmountBeforeAje
						dataControl.Piutanglainshortterm3rdparty = &dataControl2
					case "piutanglainshortterm_berelasi":
						dataControl3 := *dataAUPD.PiutanglainshorttermBerelasi - *summaryCoa1.AmountBeforeAje
						dataControl.PiutanglainshorttermBerelasi = &dataControl3
					case "piutangberelasishortterm":
						dataControl4 := *dataAUPD.Piutangberelasishortterm - *summaryCoa1.AmountBeforeAje
						dataControl.Piutangberelasishortterm = &dataControl4
					case "piutanglainlongterm_3rdparty":
						dataControl5 := *dataAUPD.Piutanglainlongterm3rdparty - *summaryCoa1.AmountBeforeAje
						dataControl.Piutanglainlongterm3rdparty = &dataControl5
					case "piutanglainlongterm_berelasi":
						dataControl6 := *dataAUPD.PiutanglainlongtermBerelasi - *summaryCoa1.AmountBeforeAje
						dataControl.PiutanglainlongtermBerelasi = &dataControl6
					case "piutangberelasilongterm":
						dataControl7 := *dataAUPD.Piutangberelasilongterm - *summaryCoa1.AmountBeforeAje
						dataControl.Piutangberelasilongterm = &dataControl7
					case "utangusaha_3rdparty":
						dataControl8 := *dataAUPD.Utangusaha3rdparty - *summaryCoa1.AmountBeforeAje
						dataControl.Utangusaha3rdparty = &dataControl8
					case "utangusaha_berelasi":
						dataControl9 := *dataAUPD.UtangusahaBerelasi - *summaryCoa1.AmountBeforeAje
						dataControl.UtangusahaBerelasi = &dataControl9
					case "utanglainshortterm_3rdparty":
						dataControl10 := *dataAUPD.Utanglainshortterm3rdparty - *summaryCoa1.AmountBeforeAje
						dataControl.Utanglainshortterm3rdparty = &dataControl10
					case "utanglainshortterm_berelasi":
						dataControl11 := *dataAUPD.UtanglainshorttermBerelasi - *summaryCoa1.AmountBeforeAje
						dataControl.UtanglainshorttermBerelasi = &dataControl11
					case "utangberelasishortterm":
						dataControl12 := *dataAUPD.Utangberelasishortterm - *summaryCoa1.AmountBeforeAje
						dataControl.Utangberelasishortterm = &dataControl12
					case "utanglainlongterm_3rdparty":
						dataControl13 := *dataAUPD.Utanglainlongterm3rdparty - *summaryCoa1.AmountBeforeAje
						dataControl.Utanglainlongterm3rdparty = &dataControl13
					case "utanglainlongterm_berelasi":
						dataControl14 := *dataAUPD.UtanglainlongtermBerelasi - *summaryCoa1.AmountBeforeAje
						dataControl.UtanglainlongtermBerelasi = &dataControl14
					case "utangberelasilongterm":
						dataControl15 := *dataAUPD.Utangberelasilongterm - *summaryCoa1.AmountBeforeAje
						dataControl.Utangberelasilongterm = &dataControl15
					}
					if cntrl.FormatterID == 2 {
						dataControlEmc2 = append(dataControlEmc2, dataControl)
					}
				}
			}
		switch v.Formatter.FormatterFor {
		case "AGING-UTANG-PIUTANG":
			aup.AgingUtangPiutangDetail = *data
			aup.Control = dataControl
		case "AGING-UTANG-PIUTANG-MUTASI-ECL":
			aup.AgingUtangPiutangMEcl = *data
			aup.ControlMEcl = dataControlEmc2
		}

	}
	result := &dto.AgingUtangPiutangDetailGetResponse{
		Datas: *aup,
	}
	return result, nil
}

func (s *service) FindByID(ctx *abstraction.Context, payload *dto.AgingUtangPiutangDetailGetByIDRequest) (*dto.AgingUtangPiutangDetailGetByIDResponse, error) {
	data, err := s.Repository.FindByID(ctx, &payload.ID)
	if err != nil {
		return &dto.AgingUtangPiutangDetailGetByIDResponse{}, helper.ErrorHandler(err)
	}

	fmtBridges, err := s.FormatterBridgesRepository.FindByID(ctx, &data.FormatterBridgesID)
	if err != nil {
		return nil, helper.ErrorHandler(err)
	}

	aup, err := s.AgingUtangPiutangRepository.FindByID(ctx, &fmtBridges.TrxRefID)
	if err != nil {
		return nil, helper.ErrorHandler(err)
	}

	allowed := helper.CompanyValidation(ctx.Auth.ID, aup.CompanyID)
	if !allowed {
		return nil, response.ErrorBuilder(&response.ErrorConstant.UnauthorizedAccess, errors.New("not allowed"))
	}

	result := &dto.AgingUtangPiutangDetailGetByIDResponse{
		AgingUtangPiutangDetailEntityModel: *data,
	}
	return result, nil
}

func (s *service) Create(ctx *abstraction.Context, payload *dto.AgingUtangPiutangDetailCreateRequest) (*dto.AgingUtangPiutangDetailCreateResponse, error) {
	var data model.AgingUtangPiutangDetailEntityModel

	fmtBridges, err := s.FormatterBridgesRepository.FindByID(ctx, &payload.FormatterBridgesID)
	if err != nil {
		return nil, helper.ErrorHandler(err)
	}

	aup, err := s.AgingUtangPiutangRepository.FindByID(ctx, &fmtBridges.TrxRefID)
	if err != nil {
		return nil, helper.ErrorHandler(err)
	}

	allowed := helper.CompanyValidation(ctx.Auth.ID, aup.CompanyID)
	if !allowed {
		return nil, response.ErrorBuilder(&response.ErrorConstant.UnauthorizedAccess, errors.New("not allowed"))
	}

	if err := trxmanager.New(s.Db).WithTrx(ctx, func(ctx *abstraction.Context) error {
		data.Context = ctx
		data.AgingUtangPiutangDetailEntity = payload.AgingUtangPiutangDetailEntity

		result, err := s.Repository.Create(ctx, &data)
		if err != nil {
			return helper.ErrorHandler(err)
		}
		data = *result
		return nil
	}); err != nil {
		return &dto.AgingUtangPiutangDetailCreateResponse{}, err
	}

	result := &dto.AgingUtangPiutangDetailCreateResponse{
		AgingUtangPiutangDetailEntityModel: data,
	}
	return result, nil
}

func (s *service) Update(ctx *abstraction.Context, payload *dto.AgingUtangPiutangDetailUpdateRequest) (*dto.AgingUtangPiutangDetailUpdateResponse, error) {
	var data model.AgingUtangPiutangDetailEntityModel
	if err := trxmanager.New(s.Db).WithTrx(ctx, func(ctx *abstraction.Context) error {
		agingupdetail, err := s.Repository.FindByID(ctx, &payload.ID)
		if err != nil {
			return helper.ErrorHandler(err)
		}

		formatterBridgesData, err := s.FormatterBridgesRepository.FindByID(ctx, &agingupdetail.FormatterBridgesID)
		if err != nil {
			return helper.ErrorHandler(err)
		}

		aupData, err := s.AgingUtangPiutangRepository.FindByID(ctx, &formatterBridgesData.TrxRefID)
		if err != nil {
			return helper.ErrorHandler(err)
		}

		if aupData.Status != 1 {
			return response.ErrorBuilder(&response.ErrorConstant.DataValidated, errors.New("cannot update data"))
		}

		allowed := helper.CompanyValidation(ctx.Auth.ID, aupData.CompanyID)
		if !allowed {
			return response.ErrorBuilder(&response.ErrorConstant.UnauthorizedAccess, errors.New("not allowed"))
		}

		//validasi data bukan isTotal atau autosummary
		criteriaFormatterValidasi := model.FormatterDetailFilterModel{}
		criteriaFormatterValidasi.FormatterID = &formatterBridgesData.FormatterID
		criteriaFormatterValidasi.Code = &agingupdetail.Code

		formatterValidasi, jmlData, err := s.FormatterDetailRepository.Find(ctx, &criteriaFormatterValidasi, &abstraction.Pagination{})
		if err != nil {
			return helper.ErrorHandler(err)
		}

		if jmlData.Count == 0 {
			return response.ErrorBuilder(&response.ErrorConstant.NotFound, errors.New("not allowed"))
		}

		for _, v := range *formatterValidasi {
			if (v.IsTotal != nil && *v.IsTotal) || (v.AutoSummary != nil && *v.AutoSummary) || (v.IsLabel != nil && *v.IsLabel) {
				return response.ErrorBuilder(&response.ErrorConstant.BadRequest, errors.New("cannot update data"))
			}
		}

		data.Context = ctx
		data.AgingUtangPiutangDetailEntity = model.AgingUtangPiutangDetailEntity{
			Piutangusaha3rdparty:         payload.Piutangusaha3rdparty,
			PiutangusahaBerelasi:         payload.PiutangusahaBerelasi,
			Piutanglainshortterm3rdparty: payload.Piutanglainshortterm3rdparty,
			PiutanglainshorttermBerelasi: payload.PiutanglainshorttermBerelasi,
			Piutangberelasishortterm:     payload.Piutangberelasishortterm,
			Piutanglainlongterm3rdparty:  payload.Piutanglainlongterm3rdparty,
			PiutanglainlongtermBerelasi:  payload.PiutanglainlongtermBerelasi,
			Piutangberelasilongterm:      payload.Piutangberelasilongterm,
			Utangusaha3rdparty:           payload.Utangusaha3rdparty,
			UtangusahaBerelasi:           payload.UtangusahaBerelasi,
			Utanglainshortterm3rdparty:   payload.Utanglainshortterm3rdparty,
			UtanglainshorttermBerelasi:   payload.UtanglainshorttermBerelasi,
			Utangberelasishortterm:       payload.Utangberelasishortterm,
			Utanglainlongterm3rdparty:    payload.Utanglainlongterm3rdparty,
			UtanglainlongtermBerelasi:    payload.UtanglainlongtermBerelasi,
			Utangberelasilongterm:        payload.Utangberelasilongterm,
		}
		// data.AgingUtangPiutangDetailEntity = payload.AgingUtangPiutangDetailEntity

		result, err := s.Repository.Update(ctx, &payload.ID, &data)
		if err != nil {
			return helper.ErrorHandler(err)
		}

		criteriaFormatterDetailSumData := model.FormatterBridgesFilterModel{}
		aup := "AGING-UTANG-PIUTANG"
		criteriaFormatterDetailSumData.Source = &aup
		criteriaFormatterDetailSumData.TrxRefID = &formatterBridgesData.TrxRefID

		formatterDetailSumData, err := s.FormatterBridgesRepository.FindSummary(ctx, &criteriaFormatterDetailSumData)
		if err != nil {
			return helper.ErrorHandler(err)
		}

		criteriaFormatterDetail := model.AgingUtangPiutangDetailFilterModel{}
		// criteriaFormatterDetail.FormatterBridgesID = &formatterBridgesData.ID
		criteriaFormatterDetail.AgingUtangPiutangID = &formatterBridgesData.TrxRefID

		formatterDetailData, err := s.Repository.FindWithFormatter(ctx, &criteriaFormatterDetail)
		if err != nil {
			return helper.ErrorHandler(err)
		}

		for _, v := range *formatterDetailSumData {
			criteriaAgingDetail := model.AgingUtangPiutangDetailFilterModel{}
			// criteriaAgingDetail.FormatterBridgesID = &agingupdetail.FormatterBridgesID
			criteriaAgingDetail.Code = &v.Code
			criteriaAgingDetail.AgingUtangPiutangID = &formatterBridgesData.TrxRefID

			if v.AutoSummary != nil && *v.AutoSummary {
				filterHelperFormatterBridges := model.FormatterBridgesFilterModel{}
				filterHelperFormatterBridges.FormatterID = &v.FormatterID
				filterHelperFormatterBridges.Source = &aup
				filterHelperFormatterBridges.TrxRefID = &formatterBridgesData.TrxRefID
				helperFormatterBridges, _, err := s.FormatterBridgesRepository.Find(ctx, &filterHelperFormatterBridges, &abstraction.Pagination{})
				if err != nil {
					return helper.ErrorHandler(err)
				}
				for _, hv := range *helperFormatterBridges {
					criteriaAgingDetail.FormatterBridgesID = &hv.ID
				}
				aupdetailsum, _, err := s.Repository.Find(ctx, &criteriaAgingDetail, &abstraction.Pagination{})
				if err != nil {
					return helper.ErrorHandler(err)
				}

				for _, a := range *aupdetailsum {
					sumAUP, err := s.Repository.FindSummary(ctx, &v.FormatterID, &a.FormatterBridgesID, &v.SortID)
					if err != nil {
						return helper.ErrorHandler(err)
					}
					updateSummary := model.AgingUtangPiutangDetailEntityModel{
						AgingUtangPiutangDetailEntity: model.AgingUtangPiutangDetailEntity{
							Piutangusaha3rdparty:         sumAUP.Piutangusaha3rdparty,
							PiutangusahaBerelasi:         sumAUP.PiutangusahaBerelasi,
							Piutanglainshortterm3rdparty: sumAUP.Piutanglainshortterm3rdparty,
							PiutanglainshorttermBerelasi: sumAUP.PiutanglainshorttermBerelasi,
							Piutangberelasishortterm:     sumAUP.Piutangberelasishortterm,
							Piutanglainlongterm3rdparty:  sumAUP.Piutanglainlongterm3rdparty,
							PiutanglainlongtermBerelasi:  sumAUP.PiutanglainlongtermBerelasi,
							Piutangberelasilongterm:      sumAUP.Piutangberelasilongterm,
							Utangusaha3rdparty:           sumAUP.Utangusaha3rdparty,
							UtangusahaBerelasi:           sumAUP.UtangusahaBerelasi,
							Utanglainshortterm3rdparty:   sumAUP.Utanglainshortterm3rdparty,
							UtanglainshorttermBerelasi:   sumAUP.UtanglainshorttermBerelasi,
							Utangberelasishortterm:       sumAUP.Utangberelasishortterm,
							Utanglainlongterm3rdparty:    sumAUP.Utanglainlongterm3rdparty,
							UtanglainlongtermBerelasi:    sumAUP.UtanglainlongtermBerelasi,
							Utangberelasilongterm:        sumAUP.Utangberelasilongterm,
						},
					}
					_, err = s.Repository.Update(ctx, &a.ID, &updateSummary)
					if err != nil {
						return helper.ErrorHandler(err)
					}
				}
			}

			if v.IsTotal != nil && *v.IsTotal {
				if v.ControlFormula != "" {
					continue
				}
				formula := v.FxSummary
				// parameterFormula := make(map[string]interface{})
				tmpString := []string{"Piutangusaha3rdparty", "PiutangusahaBerelasi", "Piutanglainshortterm3rdparty", "PiutanglainshorttermBerelasi", "Piutangberelasishortterm", "Piutanglainlongterm3rdparty", "PiutanglainlongtermBerelasi", "Piutangberelasilongterm", "Utangusaha3rdparty", "UtangusahaBerelasi", "Utanglainshortterm3rdparty", "UtanglainshorttermBerelasi", "Utangberelasishortterm", "Utanglainlongterm3rdparty", "UtanglainlongtermBerelasi", "Utangberelasilongterm"}
				tmpTotalStr := make(map[string]string)
				tmpTotalFl := make(map[string]*float64)
				for _, vFormatterDetail := range *formatterDetailData {
					if strings.Contains(formula, strings.Trim(strings.ToUpper(vFormatterDetail.Code), " ")) {
						tmpTotalStr["Piutangusaha3rdparty"] = strings.Replace(formula, vFormatterDetail.Code, fmt.Sprintf("%f", *vFormatterDetail.Piutangusaha3rdparty), -1)
						tmpTotalStr["PiutangusahaBerelasi"] = strings.Replace(formula, vFormatterDetail.Code, fmt.Sprintf("%f", *vFormatterDetail.PiutangusahaBerelasi), -1)
						tmpTotalStr["Piutanglainshortterm3rdparty"] = strings.Replace(formula, vFormatterDetail.Code, fmt.Sprintf("%f", *vFormatterDetail.Piutanglainshortterm3rdparty), -1)
						tmpTotalStr["PiutanglainshorttermBerelasi"] = strings.Replace(formula, vFormatterDetail.Code, fmt.Sprintf("%f", *vFormatterDetail.PiutanglainshorttermBerelasi), -1)
						tmpTotalStr["Piutangberelasishortterm"] = strings.Replace(formula, vFormatterDetail.Code, fmt.Sprintf("%f", *vFormatterDetail.Piutangberelasishortterm), -1)
						tmpTotalStr["Piutanglainlongterm3rdparty"] = strings.Replace(formula, vFormatterDetail.Code, fmt.Sprintf("%f", *vFormatterDetail.Piutanglainlongterm3rdparty), -1)
						tmpTotalStr["PiutanglainlongtermBerelasi"] = strings.Replace(formula, vFormatterDetail.Code, fmt.Sprintf("%f", *vFormatterDetail.PiutanglainlongtermBerelasi), -1)
						tmpTotalStr["Piutangberelasilongterm"] = strings.Replace(formula, vFormatterDetail.Code, fmt.Sprintf("%f", *vFormatterDetail.Piutangberelasilongterm), -1)
						tmpTotalStr["Utangusaha3rdparty"] = strings.Replace(formula, vFormatterDetail.Code, fmt.Sprintf("%f", *vFormatterDetail.Utangusaha3rdparty), -1)
						tmpTotalStr["UtangusahaBerelasi"] = strings.Replace(formula, vFormatterDetail.Code, fmt.Sprintf("%f", *vFormatterDetail.UtangusahaBerelasi), -1)
						tmpTotalStr["Utanglainshortterm3rdparty"] = strings.Replace(formula, vFormatterDetail.Code, fmt.Sprintf("%f", *vFormatterDetail.Utanglainshortterm3rdparty), -1)
						tmpTotalStr["UtanglainshorttermBerelasi"] = strings.Replace(formula, vFormatterDetail.Code, fmt.Sprintf("%f", *vFormatterDetail.UtanglainshorttermBerelasi), -1)
						tmpTotalStr["Utangberelasishortterm"] = strings.Replace(formula, vFormatterDetail.Code, fmt.Sprintf("%f", *vFormatterDetail.Utangberelasishortterm), -1)
						tmpTotalStr["Utanglainlongterm3rdparty"] = strings.Replace(formula, vFormatterDetail.Code, fmt.Sprintf("%f", *vFormatterDetail.Utanglainlongterm3rdparty), -1)
						tmpTotalStr["UtanglainlongtermBerelasi"] = strings.Replace(formula, vFormatterDetail.Code, fmt.Sprintf("%f", *vFormatterDetail.UtanglainlongtermBerelasi), -1)
						tmpTotalStr["Utangberelasilongterm"] = strings.Replace(formula, vFormatterDetail.Code, fmt.Sprintf("%f", *vFormatterDetail.Utangberelasilongterm), -1)
					}
				}

				// reg := regexp.MustCompile(`[^1234567890+\-\/\(\)\*\^]+`)
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

				updateSummary := model.AgingUtangPiutangDetailEntityModel{
					AgingUtangPiutangDetailEntity: model.AgingUtangPiutangDetailEntity{
						Piutangusaha3rdparty:         tmpTotalFl["Piutangusaha3rdparty"],
						PiutangusahaBerelasi:         tmpTotalFl["PiutangusahaBerelasi"],
						Piutanglainshortterm3rdparty: tmpTotalFl["Piutanglainshortterm3rdparty"],
						PiutanglainshorttermBerelasi: tmpTotalFl["PiutanglainshorttermBerelasi"],
						Piutangberelasishortterm:     tmpTotalFl["Piutangberelasishortterm"],
						Piutanglainlongterm3rdparty:  tmpTotalFl["Piutanglainlongterm3rdparty"],
						PiutanglainlongtermBerelasi:  tmpTotalFl["PiutanglainlongtermBerelasi"],
						Piutangberelasilongterm:      tmpTotalFl["Piutangberelasilongterm"],
						Utangusaha3rdparty:           tmpTotalFl["Utangusaha3rdparty"],
						UtangusahaBerelasi:           tmpTotalFl["UtangusahaBerelasi"],
						Utanglainshortterm3rdparty:   tmpTotalFl["Utanglainshortterm3rdparty"],
						UtanglainshorttermBerelasi:   tmpTotalFl["UtanglainshorttermBerelasi"],
						Utangberelasishortterm:       tmpTotalFl["Utangberelasishortterm"],
						Utanglainlongterm3rdparty:    tmpTotalFl["Utanglainlongterm3rdparty"],
						UtanglainlongtermBerelasi:    tmpTotalFl["UtanglainlongtermBerelasi"],
						Utangberelasilongterm:        tmpTotalFl["Utangberelasilongterm"],
					},
				}

				ebsum, _, err := s.Repository.Find(ctx, &criteriaAgingDetail, &abstraction.Pagination{})
				if err != nil {
					return helper.ErrorHandler(err)
				}
				for _, vAUP := range *ebsum {
					_, err = s.Repository.Update(ctx, &vAUP.ID, &updateSummary)
					if err != nil {
						return helper.ErrorHandler(err)
					}
				}
			}
		}

		criteriaValidation := model.ValidationDetailFilterModel{}
		criteriaValidation.CompanyID = &aupData.CompanyID
		criteriaValidation.Period = &aupData.Period
		criteriaValidation.Versions = &aupData.Versions
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
		return &dto.AgingUtangPiutangDetailUpdateResponse{}, err
	}
	result := &dto.AgingUtangPiutangDetailUpdateResponse{
		AgingUtangPiutangDetailEntityModel: data,
	}
	return result, nil
}

func (s *service) Delete(ctx *abstraction.Context, payload *dto.AgingUtangPiutangDetailDeleteRequest) (*dto.AgingUtangPiutangDetailDeleteResponse, error) {
	var data model.AgingUtangPiutangDetailEntityModel
	if err := trxmanager.New(s.Db).WithTrx(ctx, func(ctx *abstraction.Context) error {
		existing, err := s.Repository.FindByID(ctx, &payload.ID)
		if err != nil {
			return helper.ErrorHandler(err)
		}

		fmtBridges, err := s.FormatterBridgesRepository.FindByID(ctx, &existing.FormatterBridgesID)
		if err != nil {
			return helper.ErrorHandler(err)
		}

		aup, err := s.AgingUtangPiutangRepository.FindByID(ctx, &fmtBridges.TrxRefID)
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
		return &dto.AgingUtangPiutangDetailDeleteResponse{}, err
	}
	result := &dto.AgingUtangPiutangDetailDeleteResponse{
		// AgingUtangPiutangDetailEntityModel: data,
	}
	return result, nil
}
