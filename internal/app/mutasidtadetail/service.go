package mutasidtadetail

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
	Repository                 repository.MutasiDtaDetail
	MutasiDtaRepository        repository.MutasiDta
	FormatterBridgesRepository repository.FormatterBridges
	FormatterDetailRepository  repository.FormatterDetail
	ValidationDetailRepository repository.ValidationDetail
	FormatterRepository         repository.Formatter
	TrialBalanceDetailRepository repository.TrialBalanceDetail
	Db                         *gorm.DB
}

type Service interface {
	Find(ctx *abstraction.Context, payload *dto.MutasiDtaDetailGetRequest) (*dto.MutasiDtaDetailGetResponse, error)
	FindByID(ctx *abstraction.Context, payload *dto.MutasiDtaDetailGetByIDRequest) (*dto.MutasiDtaDetailGetByIDResponse, error)
	Create(ctx *abstraction.Context, payload *dto.MutasiDtaDetailCreateRequest) (*dto.MutasiDtaDetailCreateResponse, error)
	Update(ctx *abstraction.Context, payload *dto.MutasiDtaDetailUpdateRequest) (*dto.MutasiDtaDetailUpdateResponse, error)
	Delete(ctx *abstraction.Context, payload *dto.MutasiDtaDetailDeleteRequest) (*dto.MutasiDtaDetailDeleteResponse, error)
}

func NewService(f *factory.Factory) *service {
	repository := f.MutasiDtaDetailRepository
	mutasiDtaRepo := f.MutasiDtaRepository
	formatterBridgesRepo := f.FormatterBridgesRepository
	formatterDetailRepo := f.FormatterDetailRepository
	validationDetailRepo := f.ValidationDetailRepository
	formatterRepo := f.FormatterRepository
	trialBalanceDetailRepo := f.TrialBalanceDetailRepository
	db := f.Db
	return &service{
		Repository:                 repository,
		Db:                         db,
		MutasiDtaRepository:        mutasiDtaRepo,
		FormatterBridgesRepository: formatterBridgesRepo,
		FormatterDetailRepository:  formatterDetailRepo,
		ValidationDetailRepository: validationDetailRepo,
		FormatterRepository:          formatterRepo,
		TrialBalanceDetailRepository: trialBalanceDetailRepo,
	}
}
func (s *service) Find(ctx *abstraction.Context, payload *dto.MutasiDtaDetailGetRequest) (*dto.MutasiDtaDetailGetResponse, error) {
	mutasi, err := s.MutasiDtaRepository.FindByID(ctx, payload.MutasiDtaID)
	if err != nil {
		return &dto.MutasiDtaDetailGetResponse{}, helper.ErrorHandler(err)
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
		return &dto.MutasiDtaDetailGetResponse{}, response.ErrorBuilder(&response.ErrorConstant.InternalServerError, errors.New("not allowed"))
	}

	criteriaFormatterBridges := model.FormatterBridgesFilterModel{}
	criteriaFormatterBridges.TrxRefID = &mutasi.ID
	src := "MUTASI-DTA"
	criteriaFormatterBridges.Source = &src
	id := "id"
	asc := "ASC"
	pagingFB := abstraction.Pagination{
		SortBy: &id,
		Sort:   &asc,
	}
	formatterBridges, _, err := s.FormatterBridgesRepository.Find(ctx, &criteriaFormatterBridges, &pagingFB)
	if err != nil {
		return &dto.MutasiDtaDetailGetResponse{}, helper.ErrorHandler(err)
	}
	for _, v := range *formatterBridges {
		payload.FormatterBridgesID = &v.ID
		data, err := s.Repository.FindWithFormatter(ctx, &payload.MutasiDtaDetailFilterModel)
		if err != nil {
			return &dto.MutasiDtaDetailGetResponse{}, helper.ErrorHandler(err)
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
		var dataControl model.MutasiDtaDetailEntityModel
		var dataControlDTA []model.MutasiDtaDetailEntityModel
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
				criteriaMFD := model.MutasiDtaDetailFilterModel{}
				criteriaMFD.Code = &vFmt.Code
				criteriaMFD.FormatterBridgesID = &v.ID
				dataMFD, err := s.Repository.FindByCriteria(ctx, &criteriaMFD) // data Aging Utang Piutang Detail
				if err != nil || dataMFD.ID == 0 {
					return nil, err
				}

				switch strings.ToLower(splitControllerCommand[1]) {
				case "manfaat_beban_pajak":
					dataControls := *dataMFD.ManfaatBebanPajak - *summaryCoa1.AmountBeforeAje
					dataControl.ManfaatBebanPajak = &dataControls
				case "oci":
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
						dataControls := *dataMFD.Oci - hasilCoa2
						dataControl.Oci = &dataControls
					} else {
						dataControls := *dataMFD.Oci - *summaryCoa1.AmountBeforeAje
						dataControl.Oci = &dataControls
					}

				case "saldo_akhir":
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
							hasilCoa2 += *summaryCoa2.AmountAfterAje

						}
						dataControls := *dataMFD.SaldoAkhir - hasilCoa2
						dataControl.SaldoAkhir = &dataControls
					} else {
						dataControls := *dataMFD.SaldoAkhir - *summaryCoa1.AmountBeforeAje
						dataControl.SaldoAkhir = &dataControls
					}
				}
				dataControlDTA = append(dataControlDTA, dataControl)

			}
		}
		switch v.Formatter.FormatterFor {
		case "MUTASI-DTA":
			mutasi.MutasiDtaDetail = *data
			mutasi.ControlDTA = dataControl
		}
	}
	result := &dto.MutasiDtaDetailGetResponse{
		Datas: *mutasi,
	}
	return result, nil
}

func (s *service) FindByID(ctx *abstraction.Context, payload *dto.MutasiDtaDetailGetByIDRequest) (*dto.MutasiDtaDetailGetByIDResponse, error) {
	data, err := s.Repository.FindByID(ctx, &payload.ID)
	if err != nil {
		return &dto.MutasiDtaDetailGetByIDResponse{}, response.ErrorBuilder(&response.ErrorConstant.InternalServerError, err)
	}

	fmtBridges, err := s.FormatterBridgesRepository.FindByID(ctx, &data.FormatterBridgesID)
	if err != nil {
		return nil, helper.ErrorHandler(err)
	}

	aup, err := s.MutasiDtaRepository.FindByID(ctx, &fmtBridges.TrxRefID)
	if err != nil {
		return nil, helper.ErrorHandler(err)
	}

	allowed := helper.CompanyValidation(ctx.Auth.ID, aup.CompanyID)
	if !allowed {
		return nil, response.ErrorBuilder(&response.ErrorConstant.UnauthorizedAccess, errors.New("Not Allowed"))
	}

	result := &dto.MutasiDtaDetailGetByIDResponse{
		MutasiDtaDetailEntityModel: *data,
	}
	return result, nil
}

func (s *service) Create(ctx *abstraction.Context, payload *dto.MutasiDtaDetailCreateRequest) (*dto.MutasiDtaDetailCreateResponse, error) {
	var data model.MutasiDtaDetailEntityModel

	fmtBridges, err := s.FormatterBridgesRepository.FindByID(ctx, &payload.FormatterBridgesID)
	if err != nil {
		return nil, helper.ErrorHandler(err)
	}

	mdtaData, err := s.MutasiDtaRepository.FindByID(ctx, &fmtBridges.TrxRefID)
	if err != nil {
		return nil, helper.ErrorHandler(err)
	}

	allowed := helper.CompanyValidation(ctx.Auth.ID, mdtaData.CompanyID)
	if !allowed {
		return nil, response.ErrorBuilder(&response.ErrorConstant.UnauthorizedAccess, errors.New("Not Allowed"))
	}

	if err := trxmanager.New(s.Db).WithTrx(ctx, func(ctx *abstraction.Context) error {
		data.Context = ctx
		data.MutasiDtaDetailEntity = payload.MutasiDtaDetailEntity

		result, err := s.Repository.Create(ctx, &data)
		if err != nil {
			return response.ErrorBuilder(&response.ErrorConstant.UnprocessableEntity, err)
		}
		data = *result
		return nil
	}); err != nil {
		return &dto.MutasiDtaDetailCreateResponse{}, err
	}

	result := &dto.MutasiDtaDetailCreateResponse{
		MutasiDtaDetailEntityModel: data,
	}
	return result, nil
}

func (s *service) Update(ctx *abstraction.Context, payload *dto.MutasiDtaDetailUpdateRequest) (*dto.MutasiDtaDetailUpdateResponse, error) {
	var data model.MutasiDtaDetailEntityModel
	if err := trxmanager.New(s.Db).WithTrx(ctx, func(ctx *abstraction.Context) error {
		mdDetail, err := s.Repository.FindByID(ctx, &payload.ID)
		if err != nil {
			return helper.ErrorHandler(err)
		}

		formatterBridgesData, err := s.FormatterBridgesRepository.FindByID(ctx, &mdDetail.FormatterBridgesID)
		if err != nil {
			return helper.ErrorHandler(err)
		}

		mdtaData, err := s.MutasiDtaRepository.FindByID(ctx, &formatterBridgesData.TrxRefID)
		if err != nil {
			return helper.ErrorHandler(err)
		}

		if mdtaData.Status != 1 {
			return response.ErrorBuilder(&response.ErrorConstant.DataValidated, errors.New("Cannot Update Data"))
		}

		allowed := helper.CompanyValidation(ctx.Auth.ID, mdtaData.CompanyID)
		if !allowed {
			return response.ErrorBuilder(&response.ErrorConstant.UnauthorizedAccess, errors.New("Not Allowed"))
		}

		//validasi data bukan isTotal atau autosummary
		criteriaFormatterValidasi := model.FormatterDetailFilterModel{}
		criteriaFormatterValidasi.FormatterID = &formatterBridgesData.FormatterID
		criteriaFormatterValidasi.Code = &mdDetail.Code

		formatterValidasi, jmlData, err := s.FormatterDetailRepository.Find(ctx, &criteriaFormatterValidasi, &abstraction.Pagination{})
		if err != nil {
			return helper.ErrorHandler(err)
		}

		if jmlData.Count == 0 {
			return response.ErrorBuilder(&response.ErrorConstant.NotFound, errors.New("Data Not Found"))
		}

		for _, v := range *formatterValidasi {
			if (v.IsTotal != nil && *v.IsTotal == true) || (v.AutoSummary != nil && *v.AutoSummary == true) {
				return response.ErrorBuilder(&response.ErrorConstant.BadRequest, errors.New("Cannot update data"))
			}
		}

		data.Context = ctx
		// data.MutasiDtaDetailEntity = payload.MutasiDtaDetailEntity
		saldoAkhir := (*payload.SaldoAwal) + (*payload.ManfaatBebanPajak) + (*payload.Oci) + (*payload.AkuisisiEntitasAnak) + (*payload.DibebankanKeLr) + (*payload.DibebankanKeOci)
		data.MutasiDtaDetailEntity = model.MutasiDtaDetailEntity{
			SaldoAwal:           payload.SaldoAwal,
			ManfaatBebanPajak:   payload.ManfaatBebanPajak,
			Oci:                 payload.Oci,
			AkuisisiEntitasAnak: payload.AkuisisiEntitasAnak,
			DibebankanKeLr:      payload.DibebankanKeLr,
			DibebankanKeOci:     payload.DibebankanKeOci,
			SaldoAkhir:          &saldoAkhir,
		}

		result, err := s.Repository.Update(ctx, &payload.ID, &data)
		if err != nil {
			return helper.ErrorHandler(err)
		}

		criteriaFormatterDetailSumData := model.FormatterBridgesFilterModel{}
		mdta := "MUTASI-DTA"
		criteriaFormatterDetailSumData.Source = &mdta
		criteriaFormatterDetailSumData.TrxRefID = &formatterBridgesData.TrxRefID

		formatterDetailSumData, err := s.FormatterBridgesRepository.FindSummary(ctx, &criteriaFormatterDetailSumData)
		if err != nil {
			return helper.ErrorHandler(err)
		}

		criteriaFormatterDetail := model.MutasiDtaDetailFilterModel{}
		criteriaFormatterDetail.MutasiDtaID = &formatterBridgesData.TrxRefID

		for _, v := range *formatterDetailSumData {
			criteriaMutasiDetail := model.MutasiDtaDetailFilterModel{}
			// criteriaMutasiDetail.FormatterBridgesID = &formatterBridgesData.ID
			criteriaMutasiDetail.Code = &v.Code
			criteriaMutasiDetail.MutasiDtaID = &formatterBridgesData.TrxRefID

			if v.AutoSummary != nil && *v.AutoSummary {
				filterHelperFormatterBridges := model.FormatterBridgesFilterModel{}
				filterHelperFormatterBridges.FormatterID = &v.FormatterID
				filterHelperFormatterBridges.Source = &mdta
				filterHelperFormatterBridges.TrxRefID = &formatterBridgesData.TrxRefID
				helperFormatterBridges, _, err := s.FormatterBridgesRepository.Find(ctx, &filterHelperFormatterBridges, &abstraction.Pagination{})
				if err != nil {
					return helper.ErrorHandler(err)
				}
				for _, hv := range *helperFormatterBridges {
					criteriaMutasiDetail.FormatterBridgesID = &hv.ID
				}
				mdtadetailsum, _, err := s.Repository.Find(ctx, &criteriaMutasiDetail, &abstraction.Pagination{})
				if err != nil {
					return helper.ErrorHandler(err)
				}

				for _, a := range *mdtadetailsum {
					sumMFA, err := s.Repository.FindSummary(ctx, &v.FormatterID, &a.FormatterBridgesID, &v.SortID)
					if err != nil {
						return helper.ErrorHandler(err)
					}
					updateSummary := model.MutasiDtaDetailEntityModel{
						MutasiDtaDetailEntity: model.MutasiDtaDetailEntity{
							SaldoAwal:           sumMFA.SaldoAwal,
							ManfaatBebanPajak:   sumMFA.ManfaatBebanPajak,
							Oci:                 sumMFA.Oci,
							AkuisisiEntitasAnak: sumMFA.AkuisisiEntitasAnak,
							DibebankanKeLr:      sumMFA.DibebankanKeLr,
							DibebankanKeOci:     sumMFA.DibebankanKeOci,
							SaldoAkhir:          sumMFA.SaldoAkhir,
						},
					}
					_, err = s.Repository.Update(ctx, &a.ID, &updateSummary)
					if err != nil {
						return helper.ErrorHandler(err)
					}
				}
			}

			if v.IsTotal != nil && *v.IsTotal && v.FxSummary != "" {
				if v.Code == "CONTROL_1" || v.Code == "CONTROL_2" || v.Code == "CONTROL_3" {
					continue
				}
				// tmpString := []string{"AmountBeforeAje"}
				tmpString := []string{"SaldoAwal", "ManfaatBebanPajak", "Oci", "AkuisisiEntitasAnak", "DibebankanKeLr", "DibebankanKeOci", "SaldoAkhir"}
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
						criteriaSumDTA := model.MutasiDtaDetailFilterModel{}
						criteriaSumDTA.Code = &vMatch
						criteriaSumDTA.FormatterBridgesID = &formatterBridgesData.ID
						sumDTA, err := s.Repository.FindTotal(ctx, &criteriaSumDTA)
						if err != nil {
							return helper.ErrorHandler(err)
						}
						angka := 0.0
						if tipe == "SaldoAwal" && sumDTA.SaldoAwal != nil {
							angka = *sumDTA.SaldoAwal
						}
						if tipe == "ManfaatBebanPajak" && sumDTA.ManfaatBebanPajak != nil {
							angka = *sumDTA.ManfaatBebanPajak
						}
						if tipe == "Oci" && sumDTA.Oci != nil {
							angka = *sumDTA.Oci
						}
						if tipe == "AkuisisiEntitasAnak" && sumDTA.AkuisisiEntitasAnak != nil {
							angka = *sumDTA.AkuisisiEntitasAnak
						}
						if tipe == "DibebankanKeLr" && sumDTA.DibebankanKeLr != nil {
							angka = *sumDTA.DibebankanKeLr
						}
						if tipe == "DibebankanKeOci" && sumDTA.DibebankanKeOci != nil {
							angka = *sumDTA.DibebankanKeOci
						}
						if tipe == "SaldoAkhir" && sumDTA.SaldoAkhir != nil {
							angka = *sumDTA.SaldoAkhir
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
					updateSummary := model.MutasiDtaDetailEntityModel{
						MutasiDtaDetailEntity: model.MutasiDtaDetailEntity{
							SaldoAwal:           tmpTotalFl["SaldoAwal"],
							ManfaatBebanPajak:   tmpTotalFl["ManfaatBebanPajak"],
							Oci:                 tmpTotalFl["Oci"],
							AkuisisiEntitasAnak: tmpTotalFl["AkuisisiEntitasAnak"],
							DibebankanKeLr:      tmpTotalFl["DibebankanKeLr"],
							DibebankanKeOci:     tmpTotalFl["DibebankanKeOci"],
							SaldoAkhir:          tmpTotalFl["SaldoAkhir"],
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
		criteriaValidation.CompanyID = &mdtaData.CompanyID
		criteriaValidation.Period = &mdtaData.Period
		criteriaValidation.Versions = &mdtaData.Versions
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
		return &dto.MutasiDtaDetailUpdateResponse{}, err
	}
	result := &dto.MutasiDtaDetailUpdateResponse{
		MutasiDtaDetailEntityModel: data,
	}
	return result, nil
}

func (s *service) Delete(ctx *abstraction.Context, payload *dto.MutasiDtaDetailDeleteRequest) (*dto.MutasiDtaDetailDeleteResponse, error) {
	var data model.MutasiDtaDetailEntityModel
	if err := trxmanager.New(s.Db).WithTrx(ctx, func(ctx *abstraction.Context) error {
		existing, err := s.Repository.FindByID(ctx, &payload.ID)
		if err != nil {
			return helper.ErrorHandler(err)
		}

		fmtBridges, err := s.FormatterBridgesRepository.FindByID(ctx, &existing.FormatterBridgesID)
		if err != nil {
			return helper.ErrorHandler(err)
		}

		mdtaData, err := s.MutasiDtaRepository.FindByID(ctx, &fmtBridges.TrxRefID)
		if err != nil {
			return helper.ErrorHandler(err)
		}

		allowed := helper.CompanyValidation(ctx.Auth.ID, mdtaData.CompanyID)
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
		return &dto.MutasiDtaDetailDeleteResponse{}, err
	}
	result := &dto.MutasiDtaDetailDeleteResponse{
		MutasiDtaDetailEntityModel: data,
	}
	return result, nil
}
