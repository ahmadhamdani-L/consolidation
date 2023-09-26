package adjustment

import (
	"fmt"
	"mcash-finance-console-core/internal/abstraction"
	"mcash-finance-console-core/internal/dto"
	"mcash-finance-console-core/internal/factory"
	"mcash-finance-console-core/internal/model"
	"mcash-finance-console-core/internal/repository"
	"mcash-finance-console-core/pkg/util/helper"
	"mcash-finance-console-core/pkg/util/response"
	"mcash-finance-console-core/pkg/util/trxmanager"
	"net/http"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/Knetic/govaluate"
	"github.com/xuri/excelize/v2"
	"gorm.io/gorm"
)

type service struct {
	Repository                   repository.Adjustment
	AdjustmentDetailRepository   repository.AdjustmentDetail
	JpmRepository                repository.Jpm
	FormatterBridgesRepository   repository.FormatterBridges
	FormatterDetailRepository    repository.FormatterDetail
	TrialBalanceDetailRepository repository.TrialBalanceDetail
	Db                           *gorm.DB
}

type Service interface {
	Find(ctx *abstraction.Context, payload *dto.AdjustmentGetRequest) (*dto.AdjustmentGetResponse, error)
	FindByID(ctx *abstraction.Context, payload *dto.AdjustmentGetByIDRequest) (*dto.AdjustmentGetByIDResponse, error)
	Create(ctx *abstraction.Context, payload *dto.AdjustmentCreateRequest) (*dto.AdjustmentCreateResponse, error)
	Update(ctx *abstraction.Context, payload *dto.AdjustmentUpdateRequest) (*dto.AdjustmentUpdateResponse, error)
	Delete(ctx *abstraction.Context, payload *dto.AdjustmentDeleteRequest) (*dto.AdjustmentDeleteResponse, error)
	Export(ctx *abstraction.Context, payload *dto.AdjustmentExportRequest) (*dto.AdjustmentExportResponse, error)
	GetVersion(ctx *abstraction.Context, payload *dto.GetVersionRequest) (*dto.GetVersionResponse, error)
}

func NewService(f *factory.Factory) *service {
	repository := f.AdjustmentRepository
	AdjustmentDetailRepository := f.AdjustmentDetailRepository
	JpmRepository := f.JpmRepository
	formatterDetailRepo := f.FormatterDetailRepository
	formatterBridgesRepo := f.FormatterBridgesRepository
	TrialBalanceDetailRepository := f.TrialBalanceDetailRepository

	db := f.Db
	return &service{
		Repository:                   repository,
		AdjustmentDetailRepository:   AdjustmentDetailRepository,
		JpmRepository:                JpmRepository,
		TrialBalanceDetailRepository: TrialBalanceDetailRepository,
		FormatterBridgesRepository:   formatterBridgesRepo,
		FormatterDetailRepository:    formatterDetailRepo,
		Db:                           db,
	}
}

func (s *service) Find(ctx *abstraction.Context, payload *dto.AdjustmentGetRequest) (*dto.AdjustmentGetResponse, error) {
	data, info, err := s.Repository.Find(ctx, &payload.AdjustmentFilterModel, &payload.Pagination)
	if err != nil {
		return &dto.AdjustmentGetResponse{}, response.ErrorBuilder(&response.ErrorConstant.InternalServerError, err)
	}
	result := &dto.AdjustmentGetResponse{
		Datas:          *data,
		PaginationInfo: *info,
	}
	return result, nil
}

func (s *service) FindByID(ctx *abstraction.Context, payload *dto.AdjustmentGetByIDRequest) (*dto.AdjustmentGetByIDResponse, error) {
	data, err := s.Repository.FindByID(ctx, &payload.ID)
	if err != nil {
		return &dto.AdjustmentGetByIDResponse{}, response.ErrorBuilder(&response.ErrorConstant.InternalServerError, err)
	}
	result := &dto.AdjustmentGetByIDResponse{
		AdjustmentEntityModel: *data,
	}
	return result, nil
}
func (s *service) Create(ctx *abstraction.Context, payload *dto.AdjustmentCreateRequest) (*dto.AdjustmentCreateResponse, error) {
	var data model.AdjustmentEntityModel
	var datas []model.AdjustmentDetailEntity

	// var datas1 []model.TrialBalanceDetailEntity
	if payload.AdjustmentEntity.TbID == 0 {
		return nil, response.CustomErrorBuilder(http.StatusBadRequest, "Trial Balance ID Harus Di isi", "Trial Balance ID Harus Di isi")
	}
	if err := trxmanager.New(s.Db).WithTrx(ctx, func(ctx *abstraction.Context) error {
		fmt.Println(err)
		
		var getAje model.AdjustmentFilterModel
		getAje.CompanyID = &payload.CompanyID
		getAje.Period = &payload.Period
		getAje.TbID = &payload.TbID
		

		getAdjustment, err := s.Repository.Get(ctx, &getAje)
		if err != nil {
			fmt.Println(err)
			return err
		}

		data.Context = ctx
		data.AdjustmentEntity = payload.AdjustmentEntity
		data.AdjustmentEntity.Status = 1
		// data.AdjustmentEntity.Period = lastOfMonth.Format("2006-01-02")
		data.AdjustmentEntity.TrxNumber = "AJE" + "#" + strconv.Itoa(len(*getAdjustment)+1)
		result, err := s.Repository.Create(ctx, &data)

		if err != nil {
			fmt.Println(err)
			return response.ErrorBuilder(&response.ErrorConstant.UnprocessableEntity, err)
		}
		data = *result
		source := "TRIAL-BALANCE"
		fID, err := s.Repository.FindByFormatter(ctx, &data.TbID, &source)

		if err != nil {
			fmt.Println(err)
			return err
		}
		//create AdjustmentDetail
		datas = payload.AdjustmentDetail
		var arrAdjustmentDetail []model.AdjustmentDetailEntity
		for _, v := range datas {
			AdjustmentDetail := model.AdjustmentDetailEntityModel{
				Context:                ctx,
				AdjustmentDetailEntity: v,
			}
			AdjustmentDetail.AdjustmentID = data.ID
			AdjustmentDetail.ReffNumber = &data.TrxNumber

			_, err = s.AdjustmentDetailRepository.Create(ctx, &AdjustmentDetail)

			if err != nil {
				fmt.Println(err)
				return err
			}
			if v.CoaCode == "310401004" || v.CoaCode == "310501002" || v.CoaCode == "310502002" || v.CoaCode == "310503002" || v.CoaCode == "310402002" {
				return response.CustomErrorBuilder(http.StatusBadRequest, response.E_BAD_REQUEST, "Coa Tersebut Tidak Dapat Melalukan Jurnal "+AdjustmentDetail.CoaCode)
			}

			findCoa, err := s.JpmRepository.FindByCoa(ctx, &AdjustmentDetail.CoaCode)
			if err != nil {
				return response.CustomErrorBuilder(http.StatusBadRequest, "Tidak Dapat Menemukan Coa "+AdjustmentDetail.CoaCode, "Tidak Dapat Menemukan Coa "+AdjustmentDetail.CoaCode)
			}

			findCoaType, err := s.JpmRepository.FindByCoaType(ctx, &findCoa.CoaTypeID)
			if err != nil {
				return response.CustomErrorBuilder(http.StatusBadRequest, "Tidak Dapat Menemukan Coa di Type Coa "+AdjustmentDetail.CoaCode, "Tidak Dapat Menemukan Coa di Type Coa "+AdjustmentDetail.CoaCode)
			}
			findCoaGroup, err := s.JpmRepository.FindByCoaGroup(ctx, &findCoaType.CoaGroupID)
			if err != nil {
				return response.CustomErrorBuilder(http.StatusBadRequest, "Tidak Dapat Menemukan Coa di Group coa "+AdjustmentDetail.CoaCode, "Tidak Dapat Menemukan Coa di Group coa "+AdjustmentDetail.CoaCode)
			}
			if findCoaGroup.Name == "ASET" {
				// if *AdjustmentDetail.BalanceSheetDr == float64(0) {
				// 	return response.CustomErrorBuilder(http.StatusBadRequest, "Balance Sheet Debit Harus Di isi", "Balance Sheet Debit Harus Di isi")
				// }
				// if *AdjustmentDetail.BalanceSheetCr != float64(0) {
				// 	return response.CustomErrorBuilder(http.StatusBadRequest, "Tidak Dapat Masukan Balance Sheet Credit", "Tidak Dapat Masukan Balance Sheet Credit")
				// }
				// if *AdjustmentDetail.IncomeStatementCr != float64(0) {
				// 	return response.CustomErrorBuilder(http.StatusBadRequest, "Tidak Dapat Masukan Income Statement Credit", "Tidak Dapat Masukan Income Statement Credit")
				// }
				// if *AdjustmentDetail.IncomeStatementDr != float64(0) {
				// 	return response.CustomErrorBuilder(http.StatusBadRequest, "Tidak Dapat Masukan Income Statement Debit", "Tidak Dapat Masukan Income Statement Debit")
				// }

				tbID, err := s.Repository.FindByTbd(ctx, &fID.ID, &AdjustmentDetail.CoaCode)
				if err != nil {
					return response.CustomErrorBuilder(http.StatusBadRequest, "Tidak Dapat Menemukan Coa di Kombinasi Konsolidasi "+AdjustmentDetail.CoaCode, "Tidak Dapat Menemukan Coa "+AdjustmentDetail.CoaCode)
				}

				tbID.Context = ctx

				AmountAjeCr := *v.BalanceSheetCr + *tbID.TrialBalanceDetailEntity.AmountAjeCr
				AmountAjeDr := *v.BalanceSheetDr + *tbID.TrialBalanceDetailEntity.AmountAjeDr
				tbID.TrialBalanceDetailEntity.AmountAjeCr = &AmountAjeCr
				tbID.TrialBalanceDetailEntity.AmountAjeDr = &AmountAjeDr

				var AmountAfterAje float64 = *tbID.AmountBeforeAje + *tbID.AmountAjeDr - *tbID.AmountAjeCr

				tbID.AmountAfterAje = &AmountAfterAje

				_, err = s.Repository.UpdateTbd(ctx, &tbID.ID, tbID)
				if err != nil {
					return err
				}

				arrAdjustmentDetail = append(arrAdjustmentDetail, v)
			}
			if findCoaGroup.Name == "Liabilitas" {
				// if *AdjustmentDetail.BalanceSheetDr != float64(0) {
				// 	return response.CustomErrorBuilder(http.StatusBadRequest, "Tidak Dapat Masukan Balance Sheet Debit", "Tidak Dapat Masukan Balance Sheet Debit")
				// }
				// if *AdjustmentDetail.BalanceSheetCr == float64(0) {
				// 	return response.CustomErrorBuilder(http.StatusBadRequest, "Balance Sheet Credit Harus Di isi", "Balance Sheet Credit Harus Di isi")
				// }
				// if *AdjustmentDetail.IncomeStatementCr != float64(0) {
				// 	return response.CustomErrorBuilder(http.StatusBadRequest, "Tidak Dapat Masukan Income Statement Credit", "Tidak Dapat Masukan Income Statement Credit")
				// }
				// if *AdjustmentDetail.IncomeStatementDr != float64(0) {
				// 	return response.CustomErrorBuilder(http.StatusBadRequest, "Tidak Dapat Masukan Income Statement Debit", "Tidak Dapat Masukan Income Statement Debit")
				// }

				tbID, err := s.Repository.FindByTbd(ctx, &fID.ID, &AdjustmentDetail.CoaCode)
				if err != nil {
					return response.CustomErrorBuilder(http.StatusBadRequest, "Tidak Dapat Menemukan Coa di Kombinasi Konsolidasi "+AdjustmentDetail.CoaCode, "Tidak Dapat Menemukan Coa "+AdjustmentDetail.CoaCode)
				}
				tbID.Context = ctx

				AmountAjeCr := *v.BalanceSheetCr + *tbID.TrialBalanceDetailEntity.AmountAjeCr
				AmountAjeDr := *v.BalanceSheetDr + *tbID.TrialBalanceDetailEntity.AmountAjeDr
				tbID.TrialBalanceDetailEntity.AmountAjeCr = &AmountAjeCr
				tbID.TrialBalanceDetailEntity.AmountAjeDr = &AmountAjeDr

				var AmountAfterAje float64 = *tbID.AmountBeforeAje - *tbID.AmountAjeDr + *tbID.AmountAjeCr

				tbID.AmountAfterAje = &AmountAfterAje

				_, err = s.Repository.UpdateTbd(ctx, &tbID.ID, tbID)
				if err != nil {
					return err
				}

				arrAdjustmentDetail = append(arrAdjustmentDetail, v)
			}
			if findCoaGroup.Name == "EKUITAS" {
				// if *AdjustmentDetail.BalanceSheetDr != float64(0) {
				// 	return response.CustomErrorBuilder(http.StatusBadRequest, "Tidak Dapat Masukan Balance Sheet Debit", "Tidak Dapat Masukan Balance Sheet Debit")
				// }
				// if *AdjustmentDetail.BalanceSheetCr == float64(0) {
				// 	return response.CustomErrorBuilder(http.StatusBadRequest, "Balance Sheet Credit Harus Di isi", "Balance Sheet Credit Harus Di isi")
				// }
				// if *AdjustmentDetail.IncomeStatementCr != float64(0) {
				// 	return response.CustomErrorBuilder(http.StatusBadRequest, "Tidak Dapat Masukan Income Statement Credit", "Tidak Dapat Masukan Income Statement Credit")
				// }
				// if *AdjustmentDetail.IncomeStatementDr != float64(0) {
				// 	return response.CustomErrorBuilder(http.StatusBadRequest, "Tidak Dapat Masukan Income Statement Debit", "Tidak Dapat Masukan Income Statement Debit")
				// }

				tbID, err := s.Repository.FindByTbd(ctx, &fID.ID, &AdjustmentDetail.CoaCode)
				if err != nil {
					return response.CustomErrorBuilder(http.StatusBadRequest, "Tidak Dapat Menemukan Coa di Kombinasi Konsolidasi "+AdjustmentDetail.CoaCode, "Tidak Dapat Menemukan Coa "+AdjustmentDetail.CoaCode)
				}
				tbID.Context = ctx

				AmountAjeCr := *v.BalanceSheetCr + *tbID.TrialBalanceDetailEntity.AmountAjeCr
				AmountAjeDr := *v.BalanceSheetDr + *tbID.TrialBalanceDetailEntity.AmountAjeDr
				tbID.TrialBalanceDetailEntity.AmountAjeCr = &AmountAjeCr
				tbID.TrialBalanceDetailEntity.AmountAjeDr = &AmountAjeDr

				var AmountAfterAje float64 = *tbID.AmountBeforeAje - *tbID.AmountAjeDr + *tbID.AmountAjeCr

				tbID.AmountAfterAje = &AmountAfterAje

				_, err = s.Repository.UpdateTbd(ctx, &tbID.ID, tbID)
				if err != nil {
					return err
				}

				arrAdjustmentDetail = append(arrAdjustmentDetail, v)
			}
			if findCoaGroup.Name == "Pendapatan" {
				//    if *AdjustmentDetail.BalanceSheetDr != float64(0) {
				//    return response.CustomErrorBuilder(http.StatusBadRequest, "Tidak Dapat Masukan Balance Sheet Debit", "Tidak Dapat Masukan Balance Sheet Debit")
				//    }
				//    if *AdjustmentDetail.BalanceSheetCr != float64(0) {
				//    return response.CustomErrorBuilder(http.StatusBadRequest, "Tidak Dapat Masukan Balance Sheet Credit", "Tidak Dapat Masukan Balance Sheet Credit")
				//    }
				//    if *AdjustmentDetail.IncomeStatementCr == float64(0) {
				//    return response.CustomErrorBuilder(http.StatusBadRequest, "Income Statement Credit Harus Di isi", "Income Statement Credit Harus Di isi")
				//    }
				//    if *AdjustmentDetail.IncomeStatementDr != float64(0) {
				//    return response.CustomErrorBuilder(http.StatusBadRequest, "Tidak Dapat Masukan Income Statement Debit", "Tidak Dapat Masukan Income Statement Debit")
				//    }

				tbID, err := s.Repository.FindByTbd(ctx, &fID.ID, &AdjustmentDetail.CoaCode)
				if err != nil {
					return response.CustomErrorBuilder(http.StatusBadRequest, "Tidak Dapat Menemukan Coa di Kombinasi Konsolidasi "+AdjustmentDetail.CoaCode, "Tidak Dapat Menemukan Coa "+AdjustmentDetail.CoaCode)
				}
				tbID.Context = ctx

				AmountAjeCr := *v.IncomeStatementCr + *tbID.TrialBalanceDetailEntity.AmountAjeCr
				AmountAjeDr := *v.IncomeStatementDr + *tbID.TrialBalanceDetailEntity.AmountAjeDr
				tbID.TrialBalanceDetailEntity.AmountAjeCr = &AmountAjeCr
				tbID.TrialBalanceDetailEntity.AmountAjeDr = &AmountAjeDr

				var AmountAfterAje float64 = *tbID.AmountBeforeAje - *tbID.AmountAjeDr + *tbID.AmountAjeCr

				tbID.AmountAfterAje = &AmountAfterAje

				_, err = s.Repository.UpdateTbd(ctx, &tbID.ID, tbID)
				if err != nil {
					return err
				}

				arrAdjustmentDetail = append(arrAdjustmentDetail, v)
			}
			if findCoaGroup.Name == "HPP/COGS" {
				//    if *AdjustmentDetail.BalanceSheetDr != float64(0) {
				//    return response.CustomErrorBuilder(http.StatusBadRequest, "Tidak Dapat Masukan Balance Sheet Debit", "Tidak Dapat Masukan Balance Sheet Debit")
				//    }
				//    if *AdjustmentDetail.BalanceSheetCr != float64(0) {
				//    return response.CustomErrorBuilder(http.StatusBadRequest, "Tidak Dapat Masukan Balance Sheet Credit", "Tidak Dapat Masukan Balance Sheet Credit")
				//    }
				//    if *AdjustmentDetail.IncomeStatementCr != float64(0) {
				//    return response.CustomErrorBuilder(http.StatusBadRequest, "Tidak Dapat Masukan Income Statement Credit", "Tidak Dapat Masukan Income Statement Credit")
				//    }
				//    if *AdjustmentDetail.IncomeStatementDr == float64(0) {
				//    return response.CustomErrorBuilder(http.StatusBadRequest, "Income Statement Debit Harus Di isi", "Income Statement Debit Harus Di isi")
				//    }

				tbID, err := s.Repository.FindByTbd(ctx, &fID.ID, &AdjustmentDetail.CoaCode)
				if err != nil {
					return response.CustomErrorBuilder(http.StatusBadRequest, "Tidak Dapat Menemukan Coa di Kombinasi Konsolidasi "+AdjustmentDetail.CoaCode, "Tidak Dapat Menemukan Coa "+AdjustmentDetail.CoaCode)
				}
				tbID.Context = ctx

				AmountAjeCr := *v.IncomeStatementCr + *tbID.TrialBalanceDetailEntity.AmountAjeCr
				AmountAjeDr := *v.IncomeStatementDr + *tbID.TrialBalanceDetailEntity.AmountAjeDr
				tbID.TrialBalanceDetailEntity.AmountAjeCr = &AmountAjeCr
				tbID.TrialBalanceDetailEntity.AmountAjeDr = &AmountAjeDr

				var AmountAfterAje float64 = *tbID.AmountBeforeAje + *tbID.AmountAjeDr - *tbID.AmountAjeCr

				tbID.AmountAfterAje = &AmountAfterAje

				_, err = s.Repository.UpdateTbd(ctx, &tbID.ID, tbID)
				if err != nil {
					return err
				}

				arrAdjustmentDetail = append(arrAdjustmentDetail, v)
			}
			if findCoaGroup.Name == "Selling Expense" {
				//    if *AdjustmentDetail.BalanceSheetDr != float64(0) {
				//    return response.CustomErrorBuilder(http.StatusBadRequest, "Tidak Dapat Masukan Balance Sheet Debit", "Tidak Dapat Masukan Balance Sheet Debit")
				//    }
				//    if *AdjustmentDetail.BalanceSheetCr != float64(0) {
				//    return response.CustomErrorBuilder(http.StatusBadRequest, "Tidak Dapat Masukan Balance Sheet Credit", "Tidak Dapat Masukan Balance Sheet Credit")
				//    }
				//    if *AdjustmentDetail.IncomeStatementCr != float64(0) {
				//    return response.CustomErrorBuilder(http.StatusBadRequest, "Tidak Dapat Masukan Income Statement Credit", "Tidak Dapat Masukan Income Statement Credit")
				//    }
				//    if *AdjustmentDetail.IncomeStatementDr == float64(0) {
				//    return response.CustomErrorBuilder(http.StatusBadRequest, "Income Statement Debit Harus Di isi", "Income Statement Debit Harus Di isi")
				//    }

				tbID, err := s.Repository.FindByTbd(ctx, &fID.ID, &AdjustmentDetail.CoaCode)
				if err != nil {
					return response.CustomErrorBuilder(http.StatusBadRequest, "Tidak Dapat Menemukan Coa di Kombinasi Konsolidasi "+AdjustmentDetail.CoaCode, "Tidak Dapat Menemukan Coa "+AdjustmentDetail.CoaCode)
				}
				tbID.Context = ctx

				AmountAjeCr := *v.IncomeStatementCr + *tbID.TrialBalanceDetailEntity.AmountAjeCr
				AmountAjeDr := *v.IncomeStatementDr + *tbID.TrialBalanceDetailEntity.AmountAjeDr
				tbID.TrialBalanceDetailEntity.AmountAjeCr = &AmountAjeCr
				tbID.TrialBalanceDetailEntity.AmountAjeDr = &AmountAjeDr

				var AmountAfterAje float64 = *tbID.AmountBeforeAje + *tbID.AmountAjeDr - *tbID.AmountAjeCr

				tbID.AmountAfterAje = &AmountAfterAje

				_, err = s.Repository.UpdateTbd(ctx, &tbID.ID, tbID)
				if err != nil {
					return err
				}

				arrAdjustmentDetail = append(arrAdjustmentDetail, v)
			}
			if findCoaGroup.Name == "General & Admin Expense" {
				//    if *AdjustmentDetail.BalanceSheetDr != float64(0) {
				//    return response.CustomErrorBuilder(http.StatusBadRequest, "Tidak Dapat Masukan Balance Sheet Debit", "Tidak Dapat Masukan Balance Sheet Debit")
				//    }
				//    if *AdjustmentDetail.BalanceSheetCr != float64(0) {
				//    return response.CustomErrorBuilder(http.StatusBadRequest, "Tidak Dapat Masukan Balance Sheet Credit", "Tidak Dapat Masukan Balance Sheet Credit")
				//    }
				//    if *AdjustmentDetail.IncomeStatementCr != float64(0) {
				//    return response.CustomErrorBuilder(http.StatusBadRequest, "Tidak Dapat Masukan Income Statement Credit", "Tidak Dapat Masukan Income Statement Credit")
				//    }
				//    if *AdjustmentDetail.IncomeStatementDr == float64(0) {
				//    return response.CustomErrorBuilder(http.StatusBadRequest, "Income Statement Debit Harus Di isi", "Income Statement Debit Harus Di isi")
				//    }

				tbID, err := s.Repository.FindByTbd(ctx, &fID.ID, &AdjustmentDetail.CoaCode)
				if err != nil {
					return response.CustomErrorBuilder(http.StatusBadRequest, "Tidak Dapat Menemukan Coa di Kombinasi Konsolidasi "+AdjustmentDetail.CoaCode, "Tidak Dapat Menemukan Coa "+AdjustmentDetail.CoaCode)
				}
				tbID.Context = ctx

				AmountAjeCr := *v.IncomeStatementCr + *tbID.TrialBalanceDetailEntity.AmountAjeCr
				AmountAjeDr := *v.IncomeStatementDr + *tbID.TrialBalanceDetailEntity.AmountAjeDr
				tbID.TrialBalanceDetailEntity.AmountAjeCr = &AmountAjeCr
				tbID.TrialBalanceDetailEntity.AmountAjeDr = &AmountAjeDr

				var AmountAfterAje float64 = *tbID.AmountBeforeAje + *tbID.AmountAjeDr - *tbID.AmountAjeCr

				tbID.AmountAfterAje = &AmountAfterAje

				_, err = s.Repository.UpdateTbd(ctx, &tbID.ID, tbID)
				if err != nil {
					return err
				}

				arrAdjustmentDetail = append(arrAdjustmentDetail, v)
			}
			if findCoaGroup.Name == "Other Income" {
				//    if *AdjustmentDetail.BalanceSheetDr != float64(0) {
				//    return response.CustomErrorBuilder(http.StatusBadRequest, "Tidak Dapat Masukan Balance Sheet Debit", "Tidak Dapat Masukan Balance Sheet Debit")
				//    }
				//    if *AdjustmentDetail.BalanceSheetCr != float64(0) {
				//    return response.CustomErrorBuilder(http.StatusBadRequest, "Tidak Dapat Masukan Balance Sheet Credit", "Tidak Dapat Masukan Balance Sheet Credit")
				//    }
				//    if *AdjustmentDetail.IncomeStatementCr == float64(0) {
				//    return response.CustomErrorBuilder(http.StatusBadRequest, "Income Statement Credit Harus Di isi", "Income Statement Credit Harus Di isi")
				//    }
				//    if *AdjustmentDetail.IncomeStatementDr != float64(0) {
				//    return response.CustomErrorBuilder(http.StatusBadRequest, "Tidak Dapat Masukan Income Statement Debit", "Tidak Dapat Masukan Income Statement Debit")
				//    }

				tbID, err := s.Repository.FindByTbd(ctx, &fID.ID, &AdjustmentDetail.CoaCode)
				if err != nil {
					return response.CustomErrorBuilder(http.StatusBadRequest, "Tidak Dapat Menemukan Coa di Kombinasi Konsolidasi "+AdjustmentDetail.CoaCode, "Tidak Dapat Menemukan Coa "+AdjustmentDetail.CoaCode)
				}
				tbID.Context = ctx

				AmountAjeCr := *v.IncomeStatementCr + *tbID.TrialBalanceDetailEntity.AmountAjeCr
				AmountAjeDr := *v.IncomeStatementDr + *tbID.TrialBalanceDetailEntity.AmountAjeDr
				tbID.TrialBalanceDetailEntity.AmountAjeCr = &AmountAjeCr
				tbID.TrialBalanceDetailEntity.AmountAjeDr = &AmountAjeDr

				var AmountAfterAje float64 = *tbID.AmountBeforeAje - *tbID.AmountAjeDr + *tbID.AmountAjeCr

				tbID.AmountAfterAje = &AmountAfterAje

				_, err = s.Repository.UpdateTbd(ctx, &tbID.ID, tbID)
				if err != nil {
					return err
				}

				arrAdjustmentDetail = append(arrAdjustmentDetail, v)
			}
			if findCoaGroup.Name == "Other Expense" {
				//    if *AdjustmentDetail.BalanceSheetDr != float64(0) {
				//    return response.CustomErrorBuilder(http.StatusBadRequest, "Tidak Dapat Masukan Balance Sheet Debit", "Tidak Dapat Masukan Balance Sheet Debit")
				//    }
				//    if *AdjustmentDetail.BalanceSheetCr != float64(0) {
				//    return response.CustomErrorBuilder(http.StatusBadRequest, "Tidak Dapat Masukan Balance Sheet Credit", "Tidak Dapat Masukan Balance Sheet Credit")
				//    }
				//    if *AdjustmentDetail.IncomeStatementCr != float64(0) {
				//    return response.CustomErrorBuilder(http.StatusBadRequest, "Tidak Dapat Masukan Income Statement Credit", "Tidak Dapat Masukan Income Statement Credit")
				//    }
				//    if *AdjustmentDetail.IncomeStatementDr == float64(0) {
				//    return response.CustomErrorBuilder(http.StatusBadRequest, "Income Statement Debit Harus Di isi", "Income Statement Debit Harus Di isi")
				//    }

				tbID, err := s.Repository.FindByTbd(ctx, &fID.ID, &AdjustmentDetail.CoaCode)
				if err != nil {
					return response.CustomErrorBuilder(http.StatusBadRequest, "Tidak Dapat Menemukan Coa di Kombinasi Konsolidasi "+AdjustmentDetail.CoaCode, "Tidak Dapat Menemukan Coa "+AdjustmentDetail.CoaCode)
				}
				tbID.Context = ctx

				AmountAjeCr := *v.IncomeStatementCr + *tbID.TrialBalanceDetailEntity.AmountAjeCr
				AmountAjeDr := *v.IncomeStatementDr + *tbID.TrialBalanceDetailEntity.AmountAjeDr
				tbID.TrialBalanceDetailEntity.AmountAjeCr = &AmountAjeCr
				tbID.TrialBalanceDetailEntity.AmountAjeDr = &AmountAjeDr

				var AmountAfterAje float64 = *tbID.AmountBeforeAje + *tbID.AmountAjeDr - *tbID.AmountAjeCr

				tbID.AmountAfterAje = &AmountAfterAje

				_, err = s.Repository.UpdateTbd(ctx, &tbID.ID, tbID)
				if err != nil {
					return err
				}

				arrAdjustmentDetail = append(arrAdjustmentDetail, v)
			}
			if findCoaGroup.Name == "Tax Expense" {
				//    if *AdjustmentDetail.BalanceSheetDr != float64(0) {
				//    return response.CustomErrorBuilder(http.StatusBadRequest, "Tidak Dapat Masukan Balance Sheet Debit", "Tidak Dapat Masukan Balance Sheet Debit")
				//    }
				//    if *AdjustmentDetail.BalanceSheetCr != float64(0) {
				//    return response.CustomErrorBuilder(http.StatusBadRequest, "Tidak Dapat Masukan Balance Sheet Credit", "Tidak Dapat Masukan Balance Sheet Credit")
				//    }
				//    if *AdjustmentDetail.IncomeStatementCr != float64(0) {
				//    return response.CustomErrorBuilder(http.StatusBadRequest, "Tidak Dapat Masukan Income Statement Credit", "Tidak Dapat Masukan Income Statement Credit")
				//    }
				//    if *AdjustmentDetail.IncomeStatementDr == float64(0) {
				//    return response.CustomErrorBuilder(http.StatusBadRequest, "Income Statement Debit Harus Di isi", "Income Statement Debit Harus Di isi")
				//    }

				tbID, err := s.Repository.FindByTbd(ctx, &fID.ID, &AdjustmentDetail.CoaCode)
				if err != nil {
					return response.CustomErrorBuilder(http.StatusBadRequest, "Tidak Dapat Menemukan Coa di Kombinasi Konsolidasi "+AdjustmentDetail.CoaCode, "Tidak Dapat Menemukan Coa "+AdjustmentDetail.CoaCode)
				}
				tbID.Context = ctx

				AmountAjeCr := *v.IncomeStatementCr + *tbID.TrialBalanceDetailEntity.AmountAjeCr
				AmountAjeDr := *v.IncomeStatementDr + *tbID.TrialBalanceDetailEntity.AmountAjeDr
				tbID.TrialBalanceDetailEntity.AmountAjeCr = &AmountAjeCr
				tbID.TrialBalanceDetailEntity.AmountAjeDr = &AmountAjeDr

				var AmountAfterAje float64 = *tbID.AmountBeforeAje + *tbID.AmountAjeDr - *tbID.AmountAjeCr

				tbID.AmountAfterAje = &AmountAfterAje

				_, err = s.Repository.UpdateTbd(ctx, &tbID.ID, tbID)
				if err != nil {
					return err
				}

				arrAdjustmentDetail = append(arrAdjustmentDetail, v)
			}
			if findCoaGroup.Name == "Income (Loss) from subsidiary" {
				//    if *AdjustmentDetail.BalanceSheetDr != float64(0) {
				//    return response.CustomErrorBuilder(http.StatusBadRequest, "Tidak Dapat Masukan Balance Sheet Debit", "Tidak Dapat Masukan Balance Sheet Debit")
				//    }
				//    if *AdjustmentDetail.BalanceSheetCr != float64(0) {
				//    return response.CustomErrorBuilder(http.StatusBadRequest, "Tidak Dapat Masukan Balance Sheet Credit", "Tidak Dapat Masukan Balance Sheet Credit")
				//    }
				//    if *AdjustmentDetail.IncomeStatementCr == float64(0) {
				//    return response.CustomErrorBuilder(http.StatusBadRequest, "Income Statement Credit Harus Di isi", "Income Statement Credit Harus Di isi")
				//    }
				//    if *AdjustmentDetail.IncomeStatementDr != float64(0) {
				//    return response.CustomErrorBuilder(http.StatusBadRequest, "Tidak Dapat Masukan Income Statement Debit", "Tidak Dapat Masukan Income Statement Debit")
				//    }

				tbID, err := s.Repository.FindByTbd(ctx, &fID.ID, &AdjustmentDetail.CoaCode)
				if err != nil {
					return response.CustomErrorBuilder(http.StatusBadRequest, "Tidak Dapat Menemukan Coa di Kombinasi Konsolidasi "+AdjustmentDetail.CoaCode, "Tidak Dapat Menemukan Coa "+AdjustmentDetail.CoaCode)
				}
				tbID.Context = ctx

				AmountAjeCr := *v.IncomeStatementCr + *tbID.TrialBalanceDetailEntity.AmountAjeCr
				AmountAjeDr := *v.IncomeStatementDr + *tbID.TrialBalanceDetailEntity.AmountAjeDr
				tbID.TrialBalanceDetailEntity.AmountAjeCr = &AmountAjeCr
				tbID.TrialBalanceDetailEntity.AmountAjeDr = &AmountAjeDr

				var AmountAfterAje float64 = *tbID.AmountBeforeAje - *tbID.AmountAjeDr + *tbID.AmountAjeCr

				tbID.AmountAfterAje = &AmountAfterAje

				_, err = s.Repository.UpdateTbd(ctx, &tbID.ID, tbID)
				if err != nil {
					return err
				}

				arrAdjustmentDetail = append(arrAdjustmentDetail, v)
			}
			if findCoaGroup.Name == "MINORITY INTEREST IN NET INCOME (NCI)" {
				//    if *AdjustmentDetail.BalanceSheetDr != float64(0) {
				//    return response.CustomErrorBuilder(http.StatusBadRequest, "Tidak Dapat Masukan Balance Sheet Debit", "Tidak Dapat Masukan Balance Sheet Debit")
				//    }
				//    if *AdjustmentDetail.BalanceSheetCr != float64(0) {
				//    return response.CustomErrorBuilder(http.StatusBadRequest, "Tidak Dapat Masukan Balance Sheet Credit", "Tidak Dapat Masukan Balance Sheet Credit")
				//    }
				//    if *AdjustmentDetail.IncomeStatementCr == float64(0) {
				//    return response.CustomErrorBuilder(http.StatusBadRequest, "Income Statement Credit Harus Di isi", "Income Statement Credit Harus Di isi")
				//    }
				//    if *AdjustmentDetail.IncomeStatementDr != float64(0) {
				//    return response.CustomErrorBuilder(http.StatusBadRequest, "Tidak Dapat Masukan Income Statement Debit", "Tidak Dapat Masukan Income Statement Debit")
				//    }

				tbID, err := s.Repository.FindByTbd(ctx, &fID.ID, &AdjustmentDetail.CoaCode)
				if err != nil {
					return response.CustomErrorBuilder(http.StatusBadRequest, "Tidak Dapat Menemukan Coa di Kombinasi Konsolidasi "+AdjustmentDetail.CoaCode, "Tidak Dapat Menemukan Coa "+AdjustmentDetail.CoaCode)
				}
				tbID.Context = ctx

				AmountAjeCr := *v.IncomeStatementCr + *tbID.TrialBalanceDetailEntity.AmountAjeCr
				AmountAjeDr := *v.IncomeStatementDr + *tbID.TrialBalanceDetailEntity.AmountAjeDr
				tbID.TrialBalanceDetailEntity.AmountAjeCr = &AmountAjeCr
				tbID.TrialBalanceDetailEntity.AmountAjeDr = &AmountAjeDr

				var AmountAfterAje float64 = *tbID.AmountBeforeAje - *tbID.AmountAjeDr + *tbID.AmountAjeCr

				tbID.AmountAfterAje = &AmountAfterAje

				_, err = s.Repository.UpdateTbd(ctx, &tbID.ID, tbID)
				if err != nil {
					return err
				}

				arrAdjustmentDetail = append(arrAdjustmentDetail, v)
			}
			if findCoaGroup.Name == "Other Comprehensive Income" {

				//    if *AdjustmentDetail.BalanceSheetDr != float64(0) {
				//    return response.CustomErrorBuilder(http.StatusBadRequest, "Tidak Dapat Masukan Balance Sheet Debit", "Tidak Dapat Masukan Balance Sheet Debit")
				//    }
				//    if *AdjustmentDetail.BalanceSheetCr != float64(0) {
				//    return response.CustomErrorBuilder(http.StatusBadRequest, "Tidak Dapat Masukan Balance Sheet Credit", "Tidak Dapat Masukan Balance Sheet Credit")
				//    }
				//    if *AdjustmentDetail.IncomeStatementCr == float64(0) {
				//    return response.CustomErrorBuilder(http.StatusBadRequest, "Income Statement Credit Harus Di isi", "Income Statement Credit Harus Di isi")
				//    }
				//    if *AdjustmentDetail.IncomeStatementDr != float64(0) {
				//    return response.CustomErrorBuilder(http.StatusBadRequest, "Tidak Dapat Masukan Income Statement Debit", "Tidak Dapat Masukan Income Statement Debit")
				//    }

				tbID, err := s.Repository.FindByTbd(ctx, &fID.ID, &AdjustmentDetail.CoaCode)
				if err != nil {
					return response.CustomErrorBuilder(http.StatusBadRequest, "Tidak Dapat Menemukan Coa di Kombinasi Konsolidasi "+AdjustmentDetail.CoaCode, "Tidak Dapat Menemukan Coa "+AdjustmentDetail.CoaCode)
				}
				tbID.Context = ctx

				AmountAjeCr := *v.IncomeStatementCr + *tbID.TrialBalanceDetailEntity.AmountAjeCr
				AmountAjeDr := *v.IncomeStatementDr + *tbID.TrialBalanceDetailEntity.AmountAjeDr
				tbID.TrialBalanceDetailEntity.AmountAjeCr = &AmountAjeCr
				tbID.TrialBalanceDetailEntity.AmountAjeDr = &AmountAjeDr

				var AmountAfterAje float64 = *tbID.AmountBeforeAje - *tbID.AmountAjeDr + *tbID.AmountAjeCr

				tbID.AmountAfterAje = &AmountAfterAje

				_, err = s.Repository.UpdateTbd(ctx, &tbID.ID, tbID)
				if err != nil {
					return err
				}

				arrAdjustmentDetail = append(arrAdjustmentDetail, v)
			}
			if findCoaGroup.Name == "Dampak penyesuaian proforma  atas OCI Entitas anak" {
				//    if *AdjustmentDetail.BalanceSheetDr != float64(0) {
				//    return response.CustomErrorBuilder(http.StatusBadRequest, "Tidak Dapat Masukan Balance Sheet Debit", "Tidak Dapat Masukan Balance Sheet Debit")
				//    }
				//    if *AdjustmentDetail.BalanceSheetCr != float64(0) {
				//    return response.CustomErrorBuilder(http.StatusBadRequest, "Tidak Dapat Masukan Balance Sheet Credit", "Tidak Dapat Masukan Balance Sheet Credit")
				//    }
				//    if *AdjustmentDetail.IncomeStatementCr == float64(0) {
				//    return response.CustomErrorBuilder(http.StatusBadRequest, "Income Statement Credit Harus Di isi", "Income Statement Credit Harus Di isi")
				//    }
				//    if *AdjustmentDetail.IncomeStatementDr != float64(0) {
				//    return response.CustomErrorBuilder(http.StatusBadRequest, "Tidak Dapat Masukan Income Statement Debit", "Tidak Dapat Masukan Income Statement Debit")
				//    }

				tbID, err := s.Repository.FindByTbd(ctx, &fID.ID, &AdjustmentDetail.CoaCode)
				if err != nil {
					return response.CustomErrorBuilder(http.StatusBadRequest, "Tidak Dapat Menemukan Coa di Kombinasi Konsolidasi "+AdjustmentDetail.CoaCode, "Tidak Dapat Menemukan Coa "+AdjustmentDetail.CoaCode)
				}
				tbID.Context = ctx

				AmountAjeCr := *v.IncomeStatementCr + *tbID.TrialBalanceDetailEntity.AmountAjeCr
				AmountAjeDr := *v.IncomeStatementDr + *tbID.TrialBalanceDetailEntity.AmountAjeDr
				tbID.TrialBalanceDetailEntity.AmountAjeCr = &AmountAjeCr
				tbID.TrialBalanceDetailEntity.AmountAjeDr = &AmountAjeDr

				var AmountAfterAje float64 = *tbID.AmountBeforeAje - *tbID.AmountAjeDr + *tbID.AmountAjeCr

				tbID.AmountAfterAje = &AmountAfterAje

				_, err = s.Repository.UpdateTbd(ctx, &tbID.ID, tbID)
				if err != nil {
					return err
				}

				arrAdjustmentDetail = append(arrAdjustmentDetail, v)
			}
			if findCoaGroup.Name == "Non Controlling OCI" {
				//    if *AdjustmentDetail.BalanceSheetDr != float64(0) {
				//    return response.CustomErrorBuilder(http.StatusBadRequest, "Tidak Dapat Masukan Balance Sheet Debit", "Tidak Dapat Masukan Balance Sheet Debit")
				//    }
				//    if *AdjustmentDetail.BalanceSheetCr != float64(0) {
				//    return response.CustomErrorBuilder(http.StatusBadRequest, "Tidak Dapat Masukan Balance Sheet Credit", "Tidak Dapat Masukan Balance Sheet Credit")
				//    }
				//    if *AdjustmentDetail.IncomeStatementCr == float64(0) {
				//    return response.CustomErrorBuilder(http.StatusBadRequest, "Income Statement Credit Harus Di isi", "Income Statement Credit Harus Di isi")
				//    }
				//    if *AdjustmentDetail.IncomeStatementDr != float64(0) {
				//    return response.CustomErrorBuilder(http.StatusBadRequest, "Tidak Dapat Masukan Income Statement Debit", "Tidak Dapat Masukan Income Statement Debit")
				//    }

				tbID, err := s.Repository.FindByTbd(ctx, &fID.ID, &AdjustmentDetail.CoaCode)
				if err != nil {
					return response.CustomErrorBuilder(http.StatusBadRequest, "Tidak Dapat Menemukan Coa di Kombinasi Konsolidasi "+AdjustmentDetail.CoaCode, "Tidak Dapat Menemukan Coa "+AdjustmentDetail.CoaCode)
				}
				tbID.Context = ctx

				AmountAjeCr := *v.IncomeStatementCr + *tbID.TrialBalanceDetailEntity.AmountAjeCr
				AmountAjeDr := *v.IncomeStatementDr + *tbID.TrialBalanceDetailEntity.AmountAjeDr
				tbID.TrialBalanceDetailEntity.AmountAjeCr = &AmountAjeCr
				tbID.TrialBalanceDetailEntity.AmountAjeDr = &AmountAjeDr

				var AmountAfterAje float64 = *tbID.AmountBeforeAje - *tbID.AmountAjeDr + *tbID.AmountAjeCr

				tbID.AmountAfterAje = &AmountAfterAje

				_, err = s.Repository.UpdateTbd(ctx, &tbID.ID, tbID)
				if err != nil {
					return err
				}

				arrAdjustmentDetail = append(arrAdjustmentDetail, v)
			}
		}
		criteriaFormatterDetailSumData := model.FormatterBridgesFilterModel{}
		mfa := "TRIAL-BALANCE"
		criteriaFormatterDetailSumData.Source = &mfa
		criteriaFormatterDetailSumData.TrxRefID = &data.TbID

		formatterDetailSumData, err := s.FormatterBridgesRepository.FindSummaryTB(ctx)
		if err != nil {
			return helper.ErrorHandler(err)
		}

		for _, v := range *formatterDetailSumData {
			criteriaTBDetail := model.TrialBalanceDetailFilterModel{}
			criteriaTBDetail.TrialBalanceID = &fID.TrxRefID
			criteriaTBDetail.FormatterBridgesID = &fID.ID
			if v.IsTotal != nil && *v.IsTotal && v.FxSummary != "" {
				// tmpString := []string{"AmountBeforeAje"}
				tmpString := []string{"AmountBeforeAje", "AmountAjeDr", "AmountAjeCr", "AmountAfterAje"}
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
						sumTBD, err := s.TrialBalanceDetailRepository.FindSummary(ctx, &vMatch, &fID.ID, v.IsCoa)
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
				dataTB, err := s.TrialBalanceDetailRepository.FindByExactCode(ctx, &criteriaTBDetail)
				if err != nil {
					return helper.ErrorHandler(err)
				}

				updateSummary := model.TrialBalanceDetailEntityModel{
					TrialBalanceDetailEntity: model.TrialBalanceDetailEntity{
						AmountAjeDr:    tmpTotalFl["AmountAjeDr"],
						AmountAjeCr:    tmpTotalFl["AmountAjeCr"],
						AmountAfterAje: tmpTotalFl["AmountAfterAje"],
					},
				}
				_, err = s.TrialBalanceDetailRepository.Update(ctx, &dataTB.ID, &updateSummary)
				if err != nil {
					return helper.ErrorHandler(err)
				}
			}
			if v.Code == "LABA_KOMPREHENSIF" {
				//UPDATE CUSTOM ROW "310401004" "310402002"
				// COA 310501002 = Row 3712 --> ambil angka dari 4337
				// COA 310502002 = Row 3718 --> ambil angka dari 4342+4343
				// COA 310503002 = Row 3724 --> ambil angka dari 4345+4346
				code := "310401004"
				criteriaTBDetail := model.TrialBalanceDetailFilterModel{}
				criteriaTBDetail.FormatterBridgesID = &fID.ID
				criteriaTBDetail.Code = &code
				customRowOne, err := s.TrialBalanceDetailRepository.FindByCode(ctx, &criteriaTBDetail)
				if err != nil {
					return helper.ErrorHandler(err)
				}

				code = "310402002"
				criteriaTBDetail.Code = &code
				customRowTwo, err := s.TrialBalanceDetailRepository.FindByCode(ctx, &criteriaTBDetail)
				if err != nil {
					return helper.ErrorHandler(err)
				}

				code = "310501002"
				criteriaTBDetail.Code = &code
				customRowThree, err := s.TrialBalanceDetailRepository.FindByCode(ctx, &criteriaTBDetail)
				if err != nil {
					return helper.ErrorHandler(err)
				}

				code = "310502002"
				criteriaTBDetail.Code = &code
				customRowFour, err := s.TrialBalanceDetailRepository.FindByCode(ctx, &criteriaTBDetail)
				if err != nil {
					return helper.ErrorHandler(err)
				}

				code = "310503002"
				criteriaTBDetail.Code = &code
				customRowFive, err := s.TrialBalanceDetailRepository.FindByCode(ctx, &criteriaTBDetail)
				if err != nil {
					return helper.ErrorHandler(err)
				}

				code = "950101001" //REVALUATION FA (4337)
				criteriaTBDetail.Code = &code
				dataReFa, err := s.TrialBalanceDetailRepository.FindByCode(ctx, &criteriaTBDetail)
				if err != nil {
					return helper.ErrorHandler(err)
				}

				code = "950301001" //Financial Instrument (4342)
				criteriaTBDetail.Code = &code
				dataFinIn, err := s.TrialBalanceDetailRepository.FindByCode(ctx, &criteriaTBDetail)
				if err != nil {
					return helper.ErrorHandler(err)
				}

				code = "950301002" //Income tax relating to components of OCI (4343)
				criteriaTBDetail.Code = &code
				dataIncomeTax, err := s.TrialBalanceDetailRepository.FindByCode(ctx, &criteriaTBDetail)
				if err != nil {
					return helper.ErrorHandler(err)
				}

				code = "950401001" //Foreign Exchange (4345)
				criteriaTBDetail.Code = &code
				dataForex, err := s.TrialBalanceDetailRepository.FindByCode(ctx, &criteriaTBDetail)
				if err != nil {
					return helper.ErrorHandler(err)
				}

				code = "950401002" //Income tax relating to components of OCI (4346)
				criteriaTBDetail.Code = &code
				dataIncomeTax2, err := s.TrialBalanceDetailRepository.FindByCode(ctx, &criteriaTBDetail)
				if err != nil {
					return helper.ErrorHandler(err)
				}

				code = "LABA_BERSIH"
				criteriaTBDetail.Code = &code
				dataLaba, err := s.TrialBalanceDetailRepository.FindByCode(ctx, &criteriaTBDetail)
				if err != nil {
					return helper.ErrorHandler(err)
				}

				code = "TOTAL_PENGHASILAN_KOMPREHENSIF_LAIN~BS"
				criteriaTBDetail.Code = &code
				dataKomprehensif, err := s.TrialBalanceDetailRepository.FindByCode(ctx, &criteriaTBDetail)
				if err != nil {
					return helper.ErrorHandler(err)
				}

				summaryCodes, err := s.TrialBalanceDetailRepository.SummaryByCodes(ctx, &fID.ID, []string{"310501002", "310502002", "310503002"})
				if err != nil {
					return helper.ErrorHandler(err)
				}

				updatedCustomRowOne := model.TrialBalanceDetailEntityModel{}
				updatedCustomRowOne.Context = ctx
				updatedCustomRowOne.AmountBeforeAje = dataLaba.AmountBeforeAje
				updatedCustomRowOne.AmountAjeCr = dataLaba.AmountAjeCr
				updatedCustomRowOne.AmountAjeDr = dataLaba.AmountAjeDr
				updatedCustomRowOne.AmountAfterAje = dataLaba.AmountAfterAje

				_, err = s.TrialBalanceDetailRepository.Update(ctx, &customRowOne.ID, &updatedCustomRowOne)
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

				_, err = s.TrialBalanceDetailRepository.Update(ctx, &customRowTwo.ID, &updateCustomRowTwo)
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

				_, err = s.TrialBalanceDetailRepository.Update(ctx, &customRowThree.ID, &updatedCustomRowThree)
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

				_, err = s.TrialBalanceDetailRepository.Update(ctx, &customRowFour.ID, &updateCustomRowFour)
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

				_, err = s.TrialBalanceDetailRepository.Update(ctx, &customRowFive.ID, &updateCustomRowFive)
				if err != nil {
					return helper.ErrorHandler(err)
				}
			}
		}
		return nil

	}); err != nil {
		fmt.Println(err)
		return &dto.AdjustmentCreateResponse{}, err
	}
	result := &dto.AdjustmentCreateResponse{
		AdjustmentEntityModel: data,
	}
	return result, nil
}
func (s *service) Update(ctx *abstraction.Context, payload *dto.AdjustmentUpdateRequest) (*dto.AdjustmentUpdateResponse, error) {
	var data model.AdjustmentEntityModel

	if err := trxmanager.New(s.Db).WithTrx(ctx, func(ctx *abstraction.Context) error {
		if _, err := s.Repository.FindByID(ctx, &payload.ID); err != nil {
			return response.ErrorBuilder(&response.ErrorConstant.BadRequest, err)
		}
		data.Context = ctx
		data.AdjustmentEntity = payload.AdjustmentEntity
		result, err := s.Repository.Update(ctx, &payload.ID, &data)
		if err != nil {
			return response.ErrorBuilder(&response.ErrorConstant.InternalServerError, err)
		}
		data = *result
		return nil
	}); err != nil {
		return &dto.AdjustmentUpdateResponse{}, err
	}
	result := &dto.AdjustmentUpdateResponse{
		AdjustmentEntityModel: data,
	}
	return result, nil
}
func (s *service) Delete(ctx *abstraction.Context, payload *dto.AdjustmentDeleteRequest) (*dto.AdjustmentDeleteResponse, error) {
	var data model.AdjustmentEntityModel

	if err := trxmanager.New(s.Db).WithTrx(ctx, func(ctx *abstraction.Context) error {
		findById, err := s.Repository.FindByID(ctx, &payload.ID)
		if err != nil {
			return err
		}
		FindAjeDetail, err := s.AdjustmentDetailRepository.FindWithAjeID(ctx, &payload.ID)
		if err != nil {
			return err
		}
		source := "TRIAL-BALANCE"
		fID, err := s.Repository.FindByFormatter(ctx, &findById.TbID, &source)
		if err != nil {
			return err
		}
		var arrAdjustmentDetail []model.AdjustmentDetailEntityModel
		for _, v := range *FindAjeDetail {
			AdjustmentDetail := model.AdjustmentDetailEntityModel{
				Context: ctx,
			}

			AdjustmentID, err := s.Repository.FindByID(ctx, &payload.ID)
			if err != nil {
				return err
			}
			source := "TRIAL-BALANCE"
			fID, err := s.Repository.FindByFormatter(ctx, &AdjustmentID.TbID, &source)
			if err != nil {
				return err
			}

			// tbID, err := s.Repository.FindByTbd(ctx, &fID.ID, &v.CoaCode)
			// if err != nil {
			// 	return response.CustomErrorBuilder(http.StatusBadRequest, "Tidak Dapat Menemukan Coa "+v.CoaCode, "Tidak Dapat Menemukan Coa "+v.CoaCode)
			// }
			// tbID.Context = ctx

			// AmountAjeCr := *tbID.TrialBalanceDetailEntity.AmountAjeCr - *v.BalanceSheetCr
			// AmountAjeDr := *tbID.TrialBalanceDetailEntity.AmountAjeDr - *v.BalanceSheetDr
			// tbID.TrialBalanceDetailEntity.AmountAjeCr = &AmountAjeCr
			// tbID.TrialBalanceDetailEntity.AmountAjeDr = &AmountAjeDr
			// AmountAfterAje := *tbID.AmountBeforeAje + *tbID.AmountAjeDr - *tbID.AmountAjeCr
			// tbID.AmountAfterAje = &AmountAfterAje

			// replace := ""

			// if *v.BalanceSheetDr > float64(0) {
			// 	ReffAjeDr := strings.Replace(*tbID.ReffAjeDr, *v.ReffNumber+",", replace, 1)
			// 	tbID.ReffAjeCr = &ReffAjeDr
			// }

			// if *v.BalanceSheetCr > float64(0) {
			// 	ReffAjeCr := strings.Replace(*tbID.ReffAjeCr, *v.ReffNumber+",", replace, 1)
			// 	tbID.ReffAjeCr = &ReffAjeCr
			// }
			// _, err = s.Repository.UpdateTbd(ctx, &tbID.ID, tbID)
			// if err != nil {
			// 	return err
			// }
			// arrAdjustmentDetail = append(arrAdjustmentDetail, AdjustmentDetail)
			findCoa, err := s.JpmRepository.FindByCoa(ctx, &v.CoaCode)
			if err != nil {
				return response.CustomErrorBuilder(http.StatusBadRequest, "Tidak Dapat Menemukan Coa "+AdjustmentDetail.CoaCode, "Tidak Dapat Menemukan Coa "+AdjustmentDetail.CoaCode)
			}

			findCoaType, err := s.JpmRepository.FindByCoaType(ctx, &findCoa.CoaTypeID)
			if err != nil {
				return response.CustomErrorBuilder(http.StatusBadRequest, "Tidak Dapat Menemukan Coa di Type Coa "+AdjustmentDetail.CoaCode, "Tidak Dapat Menemukan Coa di Type Coa "+AdjustmentDetail.CoaCode)
			}
			findCoaGroup, err := s.JpmRepository.FindByCoaGroup(ctx, &findCoaType.CoaGroupID)
			if err != nil {
				return response.CustomErrorBuilder(http.StatusBadRequest, "Tidak Dapat Menemukan Coa di Group coa "+AdjustmentDetail.CoaCode, "Tidak Dapat Menemukan Coa di Group coa "+AdjustmentDetail.CoaCode)
			}
			if findCoaGroup.Name == "ASET" {
				// if *AdjustmentDetail.BalanceSheetDr == float64(0) {
				// 	return response.CustomErrorBuilder(http.StatusBadRequest, "Balance Sheet Debit Harus Di isi", "Balance Sheet Debit Harus Di isi")
				// }
				// if *AdjustmentDetail.BalanceSheetCr != float64(0) {
				// 	return response.CustomErrorBuilder(http.StatusBadRequest, "Tidak Dapat Masukan Balance Sheet Credit", "Tidak Dapat Masukan Balance Sheet Credit")
				// }
				// if *AdjustmentDetail.IncomeStatementCr != float64(0) {
				// 	return response.CustomErrorBuilder(http.StatusBadRequest, "Tidak Dapat Masukan Income Statement Credit", "Tidak Dapat Masukan Income Statement Credit")
				// }
				// if *AdjustmentDetail.IncomeStatementDr != float64(0) {
				// 	return response.CustomErrorBuilder(http.StatusBadRequest, "Tidak Dapat Masukan Income Statement Debit", "Tidak Dapat Masukan Income Statement Debit")
				// }

				tbID, err := s.Repository.FindByTbd(ctx, &fID.ID, &v.CoaCode)
				if err != nil {
					return response.CustomErrorBuilder(http.StatusBadRequest, "Tidak Dapat Menemukan Coa di Kombinasi Konsolidasi "+AdjustmentDetail.CoaCode, "Tidak Dapat Menemukan Coa "+AdjustmentDetail.CoaCode)
				}

				tbID.Context = ctx

				AmountAjeCr := *tbID.TrialBalanceDetailEntity.AmountAjeCr - *v.BalanceSheetCr
				AmountAjeDr := *tbID.TrialBalanceDetailEntity.AmountAjeDr - *v.BalanceSheetDr
				tbID.TrialBalanceDetailEntity.AmountAjeCr = &AmountAjeCr
				tbID.TrialBalanceDetailEntity.AmountAjeDr = &AmountAjeDr

				var AmountAfterAje float64 = *tbID.AmountBeforeAje + *tbID.AmountAjeDr - *tbID.AmountAjeCr

				tbID.AmountAfterAje = &AmountAfterAje

				_, err = s.Repository.UpdateTbd(ctx, &tbID.ID, tbID)
				if err != nil {
					return err
				}

				arrAdjustmentDetail = append(arrAdjustmentDetail, v)
			}
			if findCoaGroup.Name == "Liabilitas" {
				// if *AdjustmentDetail.BalanceSheetDr != float64(0) {
				// 	return response.CustomErrorBuilder(http.StatusBadRequest, "Tidak Dapat Masukan Balance Sheet Debit", "Tidak Dapat Masukan Balance Sheet Debit")
				// }
				// if *AdjustmentDetail.BalanceSheetCr == float64(0) {
				// 	return response.CustomErrorBuilder(http.StatusBadRequest, "Balance Sheet Credit Harus Di isi", "Balance Sheet Credit Harus Di isi")
				// }
				// if *AdjustmentDetail.IncomeStatementCr != float64(0) {
				// 	return response.CustomErrorBuilder(http.StatusBadRequest, "Tidak Dapat Masukan Income Statement Credit", "Tidak Dapat Masukan Income Statement Credit")
				// }
				// if *AdjustmentDetail.IncomeStatementDr != float64(0) {
				// 	return response.CustomErrorBuilder(http.StatusBadRequest, "Tidak Dapat Masukan Income Statement Debit", "Tidak Dapat Masukan Income Statement Debit")
				// }

				tbID, err := s.Repository.FindByTbd(ctx, &fID.ID, &v.CoaCode)
				if err != nil {
					return response.CustomErrorBuilder(http.StatusBadRequest, "Tidak Dapat Menemukan Coa di Kombinasi Konsolidasi "+v.CoaCode, "Tidak Dapat Menemukan Coa "+v.CoaCode)
				}
				tbID.Context = ctx

				AmountAjeCr := *tbID.TrialBalanceDetailEntity.AmountAjeCr - *v.BalanceSheetCr
				AmountAjeDr := *tbID.TrialBalanceDetailEntity.AmountAjeDr - *v.BalanceSheetDr
				tbID.TrialBalanceDetailEntity.AmountAjeCr = &AmountAjeCr
				tbID.TrialBalanceDetailEntity.AmountAjeDr = &AmountAjeDr

				var AmountAfterAje float64 = *tbID.AmountBeforeAje - *tbID.AmountAjeDr + *tbID.AmountAjeCr

				tbID.AmountAfterAje = &AmountAfterAje

				_, err = s.Repository.UpdateTbd(ctx, &tbID.ID, tbID)
				if err != nil {
					return err
				}

				arrAdjustmentDetail = append(arrAdjustmentDetail, v)
			}
			if findCoaGroup.Name == "EKUITAS" {
				// if *AdjustmentDetail.BalanceSheetDr != float64(0) {
				// 	return response.CustomErrorBuilder(http.StatusBadRequest, "Tidak Dapat Masukan Balance Sheet Debit", "Tidak Dapat Masukan Balance Sheet Debit")
				// }
				// if *AdjustmentDetail.BalanceSheetCr == float64(0) {
				// 	return response.CustomErrorBuilder(http.StatusBadRequest, "Balance Sheet Credit Harus Di isi", "Balance Sheet Credit Harus Di isi")
				// }
				// if *AdjustmentDetail.IncomeStatementCr != float64(0) {
				// 	return response.CustomErrorBuilder(http.StatusBadRequest, "Tidak Dapat Masukan Income Statement Credit", "Tidak Dapat Masukan Income Statement Credit")
				// }
				// if *AdjustmentDetail.IncomeStatementDr != float64(0) {
				// 	return response.CustomErrorBuilder(http.StatusBadRequest, "Tidak Dapat Masukan Income Statement Debit", "Tidak Dapat Masukan Income Statement Debit")
				// }

				tbID, err := s.Repository.FindByTbd(ctx, &fID.ID, &v.CoaCode)
				if err != nil {
					return response.CustomErrorBuilder(http.StatusBadRequest, "Tidak Dapat Menemukan Coa di Kombinasi Konsolidasi "+v.CoaCode, "Tidak Dapat Menemukan Coa "+v.CoaCode)
				}
				tbID.Context = ctx

				AmountAjeCr := *tbID.TrialBalanceDetailEntity.AmountAjeCr - *v.BalanceSheetCr
				AmountAjeDr := *tbID.TrialBalanceDetailEntity.AmountAjeDr - *v.BalanceSheetDr
				tbID.TrialBalanceDetailEntity.AmountAjeCr = &AmountAjeCr
				tbID.TrialBalanceDetailEntity.AmountAjeDr = &AmountAjeDr

				var AmountAfterAje float64 = *tbID.AmountBeforeAje - *tbID.AmountAjeDr + *tbID.AmountAjeCr

				tbID.AmountAfterAje = &AmountAfterAje

				_, err = s.Repository.UpdateTbd(ctx, &tbID.ID, tbID)
				if err != nil {
					return err
				}

				arrAdjustmentDetail = append(arrAdjustmentDetail, v)
			}
			if findCoaGroup.Name == "Pendapatan" {
				// if *AdjustmentDetail.BalanceSheetDr != float64(0) {
				// 	return response.CustomErrorBuilder(http.StatusBadRequest, "Tidak Dapat Masukan Balance Sheet Debit", "Tidak Dapat Masukan Balance Sheet Debit")
				// }
				// if *AdjustmentDetail.BalanceSheetCr != float64(0) {
				// 	return response.CustomErrorBuilder(http.StatusBadRequest, "Tidak Dapat Masukan Balance Sheet Credit", "Tidak Dapat Masukan Balance Sheet Credit")
				// }
				//    if *AdjustmentDetail.IncomeStatementCr == float64(0) {
				//    return response.CustomErrorBuilder(http.StatusBadRequest, "Income Statement Credit Harus Di isi", "Income Statement Credit Harus Di isi")
				//    }
				//    if *AdjustmentDetail.IncomeStatementDr != float64(0) {
				//    return response.CustomErrorBuilder(http.StatusBadRequest, "Tidak Dapat Masukan Income Statement Debit", "Tidak Dapat Masukan Income Statement Debit")
				//    }

				tbID, err := s.Repository.FindByTbd(ctx, &fID.ID, &v.CoaCode)
				if err != nil {
					return response.CustomErrorBuilder(http.StatusBadRequest, "Tidak Dapat Menemukan Coa di Kombinasi Konsolidasi "+v.CoaCode, "Tidak Dapat Menemukan Coa "+v.CoaCode)
				}
				tbID.Context = ctx

				AmountAjeCr := *tbID.TrialBalanceDetailEntity.AmountAjeCr - *v.IncomeStatementCr
				AmountAjeDr := *tbID.TrialBalanceDetailEntity.AmountAjeDr - *v.IncomeStatementDr
				tbID.TrialBalanceDetailEntity.AmountAjeCr = &AmountAjeCr
				tbID.TrialBalanceDetailEntity.AmountAjeDr = &AmountAjeDr

				var AmountAfterAje float64 = *tbID.AmountBeforeAje - *tbID.AmountAjeDr + *tbID.AmountAjeCr

				tbID.AmountAfterAje = &AmountAfterAje

				_, err = s.Repository.UpdateTbd(ctx, &tbID.ID, tbID)
				if err != nil {
					return err
				}

				arrAdjustmentDetail = append(arrAdjustmentDetail, v)
			}
			if findCoaGroup.Name == "HPP/COGS" {
				// if *AdjustmentDetail.BalanceSheetDr != float64(0) {
				// 	return response.CustomErrorBuilder(http.StatusBadRequest, "Tidak Dapat Masukan Balance Sheet Debit", "Tidak Dapat Masukan Balance Sheet Debit")
				// }
				// if *AdjustmentDetail.BalanceSheetCr != float64(0) {
				// 	return response.CustomErrorBuilder(http.StatusBadRequest, "Tidak Dapat Masukan Balance Sheet Credit", "Tidak Dapat Masukan Balance Sheet Credit")
				// }
				//    if *AdjustmentDetail.IncomeStatementCr != float64(0) {
				//    return response.CustomErrorBuilder(http.StatusBadRequest, "Tidak Dapat Masukan Income Statement Credit", "Tidak Dapat Masukan Income Statement Credit")
				//    }
				//    if *AdjustmentDetail.IncomeStatementDr == float64(0) {
				//    return response.CustomErrorBuilder(http.StatusBadRequest, "Income Statement Debit Harus Di isi", "Income Statement Debit Harus Di isi")
				//    }

				tbID, err := s.Repository.FindByTbd(ctx, &fID.ID, &v.CoaCode)
				if err != nil {
					return response.CustomErrorBuilder(http.StatusBadRequest, "Tidak Dapat Menemukan Coa di Kombinasi Konsolidasi "+v.CoaCode, "Tidak Dapat Menemukan Coa "+v.CoaCode)
				}
				tbID.Context = ctx

				AmountAjeCr := *tbID.TrialBalanceDetailEntity.AmountAjeCr - *v.IncomeStatementCr
				AmountAjeDr := *tbID.TrialBalanceDetailEntity.AmountAjeDr - *v.IncomeStatementDr
				tbID.TrialBalanceDetailEntity.AmountAjeCr = &AmountAjeCr
				tbID.TrialBalanceDetailEntity.AmountAjeDr = &AmountAjeDr

				var AmountAfterAje float64 = *tbID.AmountBeforeAje + *tbID.AmountAjeDr - *tbID.AmountAjeCr

				tbID.AmountAfterAje = &AmountAfterAje

				_, err = s.Repository.UpdateTbd(ctx, &tbID.ID, tbID)
				if err != nil {
					return err
				}

				arrAdjustmentDetail = append(arrAdjustmentDetail, v)
			}
			if findCoaGroup.Name == "Selling Expense" {
				// if *AdjustmentDetail.BalanceSheetDr != float64(0) {
				// 	return response.CustomErrorBuilder(http.StatusBadRequest, "Tidak Dapat Masukan Balance Sheet Debit", "Tidak Dapat Masukan Balance Sheet Debit")
				// }
				// if *AdjustmentDetail.BalanceSheetCr != float64(0) {
				// 	return response.CustomErrorBuilder(http.StatusBadRequest, "Tidak Dapat Masukan Balance Sheet Credit", "Tidak Dapat Masukan Balance Sheet Credit")
				// }
				//    if *AdjustmentDetail.IncomeStatementCr != float64(0) {
				//    return response.CustomErrorBuilder(http.StatusBadRequest, "Tidak Dapat Masukan Income Statement Credit", "Tidak Dapat Masukan Income Statement Credit")
				//    }
				//    if *AdjustmentDetail.IncomeStatementDr == float64(0) {
				//    return response.CustomErrorBuilder(http.StatusBadRequest, "Income Statement Debit Harus Di isi", "Income Statement Debit Harus Di isi")
				//    }

				tbID, err := s.Repository.FindByTbd(ctx, &fID.ID, &v.CoaCode)
				if err != nil {
					return response.CustomErrorBuilder(http.StatusBadRequest, "Tidak Dapat Menemukan Coa di Kombinasi Konsolidasi "+v.CoaCode, "Tidak Dapat Menemukan Coa "+v.CoaCode)
				}
				tbID.Context = ctx

				AmountAjeCr := *tbID.TrialBalanceDetailEntity.AmountAjeCr - *v.IncomeStatementCr
				AmountAjeDr := *tbID.TrialBalanceDetailEntity.AmountAjeDr - *v.IncomeStatementDr
				tbID.TrialBalanceDetailEntity.AmountAjeCr = &AmountAjeCr
				tbID.TrialBalanceDetailEntity.AmountAjeDr = &AmountAjeDr

				var AmountAfterAje float64 = *tbID.AmountBeforeAje + *tbID.AmountAjeDr - *tbID.AmountAjeCr

				tbID.AmountAfterAje = &AmountAfterAje

				_, err = s.Repository.UpdateTbd(ctx, &tbID.ID, tbID)
				if err != nil {
					return err
				}

				arrAdjustmentDetail = append(arrAdjustmentDetail, v)
			}
			if findCoaGroup.Name == "General & Admin Expense" {
				// if *AdjustmentDetail.BalanceSheetDr != float64(0) {
				// 	return response.CustomErrorBuilder(http.StatusBadRequest, "Tidak Dapat Masukan Balance Sheet Debit", "Tidak Dapat Masukan Balance Sheet Debit")
				// }
				// if *AdjustmentDetail.BalanceSheetCr != float64(0) {
				// 	return response.CustomErrorBuilder(http.StatusBadRequest, "Tidak Dapat Masukan Balance Sheet Credit", "Tidak Dapat Masukan Balance Sheet Credit")
				// }
				//    if *AdjustmentDetail.IncomeStatementCr != float64(0) {
				//    return response.CustomErrorBuilder(http.StatusBadRequest, "Tidak Dapat Masukan Income Statement Credit", "Tidak Dapat Masukan Income Statement Credit")
				//    }
				//    if *AdjustmentDetail.IncomeStatementDr == float64(0) {
				//    return response.CustomErrorBuilder(http.StatusBadRequest, "Income Statement Debit Harus Di isi", "Income Statement Debit Harus Di isi")
				//    }

				tbID, err := s.Repository.FindByTbd(ctx, &fID.ID, &v.CoaCode)
				if err != nil {
					return response.CustomErrorBuilder(http.StatusBadRequest, "Tidak Dapat Menemukan Coa di Kombinasi Konsolidasi "+v.CoaCode, "Tidak Dapat Menemukan Coa "+v.CoaCode)
				}
				tbID.Context = ctx

				AmountAjeCr := *tbID.TrialBalanceDetailEntity.AmountAjeCr - *v.IncomeStatementCr
				AmountAjeDr := *tbID.TrialBalanceDetailEntity.AmountAjeDr - *v.IncomeStatementDr
				tbID.TrialBalanceDetailEntity.AmountAjeCr = &AmountAjeCr
				tbID.TrialBalanceDetailEntity.AmountAjeDr = &AmountAjeDr

				var AmountAfterAje float64 = *tbID.AmountBeforeAje + *tbID.AmountAjeDr - *tbID.AmountAjeCr

				tbID.AmountAfterAje = &AmountAfterAje

				_, err = s.Repository.UpdateTbd(ctx, &tbID.ID, tbID)
				if err != nil {
					return err
				}

				arrAdjustmentDetail = append(arrAdjustmentDetail, v)
			}
			if findCoaGroup.Name == "Other Income" {
				// if *AdjustmentDetail.BalanceSheetDr != float64(0) {
				// 	return response.CustomErrorBuilder(http.StatusBadRequest, "Tidak Dapat Masukan Balance Sheet Debit", "Tidak Dapat Masukan Balance Sheet Debit")
				// }
				// if *AdjustmentDetail.BalanceSheetCr != float64(0) {
				// 	return response.CustomErrorBuilder(http.StatusBadRequest, "Tidak Dapat Masukan Balance Sheet Credit", "Tidak Dapat Masukan Balance Sheet Credit")
				// }
				//    if *AdjustmentDetail.IncomeStatementCr == float64(0) {
				//    return response.CustomErrorBuilder(http.StatusBadRequest, "Income Statement Credit Harus Di isi", "Income Statement Credit Harus Di isi")
				//    }
				//    if *AdjustmentDetail.IncomeStatementDr != float64(0) {
				//    return response.CustomErrorBuilder(http.StatusBadRequest, "Tidak Dapat Masukan Income Statement Debit", "Tidak Dapat Masukan Income Statement Debit")
				//    }

				tbID, err := s.Repository.FindByTbd(ctx, &fID.ID, &v.CoaCode)
				if err != nil {
					return response.CustomErrorBuilder(http.StatusBadRequest, "Tidak Dapat Menemukan Coa di Kombinasi Konsolidasi "+v.CoaCode, "Tidak Dapat Menemukan Coa "+v.CoaCode)
				}
				tbID.Context = ctx

				AmountAjeCr := *tbID.TrialBalanceDetailEntity.AmountAjeCr - *v.IncomeStatementCr
				AmountAjeDr := *tbID.TrialBalanceDetailEntity.AmountAjeDr - *v.IncomeStatementDr
				tbID.TrialBalanceDetailEntity.AmountAjeCr = &AmountAjeCr
				tbID.TrialBalanceDetailEntity.AmountAjeDr = &AmountAjeDr

				var AmountAfterAje float64 = *tbID.AmountBeforeAje - *tbID.AmountAjeDr + *tbID.AmountAjeCr

				tbID.AmountAfterAje = &AmountAfterAje

				_, err = s.Repository.UpdateTbd(ctx, &tbID.ID, tbID)
				if err != nil {
					return err
				}

				arrAdjustmentDetail = append(arrAdjustmentDetail, v)
			}
			if findCoaGroup.Name == "Other Expense" {
				// if *AdjustmentDetail.BalanceSheetDr != float64(0) {
				// 	return response.CustomErrorBuilder(http.StatusBadRequest, "Tidak Dapat Masukan Balance Sheet Debit", "Tidak Dapat Masukan Balance Sheet Debit")
				// }
				// if *AdjustmentDetail.BalanceSheetCr != float64(0) {
				// 	return response.CustomErrorBuilder(http.StatusBadRequest, "Tidak Dapat Masukan Balance Sheet Credit", "Tidak Dapat Masukan Balance Sheet Credit")
				// }
				//    if *AdjustmentDetail.IncomeStatementCr != float64(0) {
				//    return response.CustomErrorBuilder(http.StatusBadRequest, "Tidak Dapat Masukan Income Statement Credit", "Tidak Dapat Masukan Income Statement Credit")
				//    }
				//    if *AdjustmentDetail.IncomeStatementDr == float64(0) {
				//    return response.CustomErrorBuilder(http.StatusBadRequest, "Income Statement Debit Harus Di isi", "Income Statement Debit Harus Di isi")
				//    }

				tbID, err := s.Repository.FindByTbd(ctx, &fID.ID, &v.CoaCode)
				if err != nil {
					return response.CustomErrorBuilder(http.StatusBadRequest, "Tidak Dapat Menemukan Coa di Kombinasi Konsolidasi "+v.CoaCode, "Tidak Dapat Menemukan Coa "+v.CoaCode)
				}
				tbID.Context = ctx

				AmountAjeCr := *tbID.TrialBalanceDetailEntity.AmountAjeCr - *v.IncomeStatementCr
				AmountAjeDr := *tbID.TrialBalanceDetailEntity.AmountAjeDr - *v.IncomeStatementDr
				tbID.TrialBalanceDetailEntity.AmountAjeCr = &AmountAjeCr
				tbID.TrialBalanceDetailEntity.AmountAjeDr = &AmountAjeDr

				var AmountAfterAje float64 = *tbID.AmountBeforeAje + *tbID.AmountAjeDr - *tbID.AmountAjeCr

				tbID.AmountAfterAje = &AmountAfterAje

				_, err = s.Repository.UpdateTbd(ctx, &tbID.ID, tbID)
				if err != nil {
					return err
				}

				arrAdjustmentDetail = append(arrAdjustmentDetail, v)
			}
			if findCoaGroup.Name == "Tax Expense" {
				// if *AdjustmentDetail.BalanceSheetDr != float64(0) {
				// 	return response.CustomErrorBuilder(http.StatusBadRequest, "Tidak Dapat Masukan Balance Sheet Debit", "Tidak Dapat Masukan Balance Sheet Debit")
				// }
				// if *AdjustmentDetail.BalanceSheetCr != float64(0) {
				// 	return response.CustomErrorBuilder(http.StatusBadRequest, "Tidak Dapat Masukan Balance Sheet Credit", "Tidak Dapat Masukan Balance Sheet Credit")
				// }
				//    if *AdjustmentDetail.IncomeStatementCr != float64(0) {
				//    return response.CustomErrorBuilder(http.StatusBadRequest, "Tidak Dapat Masukan Income Statement Credit", "Tidak Dapat Masukan Income Statement Credit")
				//    }
				//    if *AdjustmentDetail.IncomeStatementDr == float64(0) {
				//    return response.CustomErrorBuilder(http.StatusBadRequest, "Income Statement Debit Harus Di isi", "Income Statement Debit Harus Di isi")
				//    }

				tbID, err := s.Repository.FindByTbd(ctx, &fID.ID, &v.CoaCode)
				if err != nil {
					return response.CustomErrorBuilder(http.StatusBadRequest, "Tidak Dapat Menemukan Coa di Kombinasi Konsolidasi "+v.CoaCode, "Tidak Dapat Menemukan Coa "+v.CoaCode)
				}
				tbID.Context = ctx

				AmountAjeCr := *tbID.TrialBalanceDetailEntity.AmountAjeCr - *v.IncomeStatementCr
				AmountAjeDr := *tbID.TrialBalanceDetailEntity.AmountAjeDr - *v.IncomeStatementDr
				tbID.TrialBalanceDetailEntity.AmountAjeCr = &AmountAjeCr
				tbID.TrialBalanceDetailEntity.AmountAjeDr = &AmountAjeDr

				var AmountAfterAje float64 = *tbID.AmountBeforeAje + *tbID.AmountAjeDr - *tbID.AmountAjeCr

				tbID.AmountAfterAje = &AmountAfterAje

				_, err = s.Repository.UpdateTbd(ctx, &tbID.ID, tbID)
				if err != nil {
					return err
				}

				arrAdjustmentDetail = append(arrAdjustmentDetail, v)
			}
			if findCoaGroup.Name == "Income (Loss) from subsidiary" {
				// if *AdjustmentDetail.BalanceSheetDr != float64(0) {
				// 	return response.CustomErrorBuilder(http.StatusBadRequest, "Tidak Dapat Masukan Balance Sheet Debit", "Tidak Dapat Masukan Balance Sheet Debit")
				// }
				// if *AdjustmentDetail.BalanceSheetCr != float64(0) {
				// 	return response.CustomErrorBuilder(http.StatusBadRequest, "Tidak Dapat Masukan Balance Sheet Credit", "Tidak Dapat Masukan Balance Sheet Credit")
				// }
				//    if *AdjustmentDetail.IncomeStatementCr == float64(0) {
				//    return response.CustomErrorBuilder(http.StatusBadRequest, "Income Statement Credit Harus Di isi", "Income Statement Credit Harus Di isi")
				//    }
				//    if *AdjustmentDetail.IncomeStatementDr != float64(0) {
				//    return response.CustomErrorBuilder(http.StatusBadRequest, "Tidak Dapat Masukan Income Statement Debit", "Tidak Dapat Masukan Income Statement Debit")
				//    }

				tbID, err := s.Repository.FindByTbd(ctx, &fID.ID, &v.CoaCode)
				if err != nil {
					return response.CustomErrorBuilder(http.StatusBadRequest, "Tidak Dapat Menemukan Coa di Kombinasi Konsolidasi "+v.CoaCode, "Tidak Dapat Menemukan Coa "+v.CoaCode)
				}
				tbID.Context = ctx

				AmountAjeCr := *tbID.TrialBalanceDetailEntity.AmountAjeCr - *v.IncomeStatementCr
				AmountAjeDr := *tbID.TrialBalanceDetailEntity.AmountAjeDr - *v.IncomeStatementDr
				tbID.TrialBalanceDetailEntity.AmountAjeCr = &AmountAjeCr
				tbID.TrialBalanceDetailEntity.AmountAjeDr = &AmountAjeDr

				var AmountAfterAje float64 = *tbID.AmountBeforeAje - *tbID.AmountAjeDr + *tbID.AmountAjeCr

				tbID.AmountAfterAje = &AmountAfterAje

				_, err = s.Repository.UpdateTbd(ctx, &tbID.ID, tbID)
				if err != nil {
					return err
				}

				arrAdjustmentDetail = append(arrAdjustmentDetail, v)
			}
			if findCoaGroup.Name == "MINORITY INTEREST IN NET INCOME (NCI)" {
				// if *AdjustmentDetail.BalanceSheetDr != float64(0) {
				// 	return response.CustomErrorBuilder(http.StatusBadRequest, "Tidak Dapat Masukan Balance Sheet Debit", "Tidak Dapat Masukan Balance Sheet Debit")
				// }
				// if *AdjustmentDetail.BalanceSheetCr != float64(0) {
				// 	return response.CustomErrorBuilder(http.StatusBadRequest, "Tidak Dapat Masukan Balance Sheet Credit", "Tidak Dapat Masukan Balance Sheet Credit")
				// }
				//    if *AdjustmentDetail.IncomeStatementCr == float64(0) {
				//    return response.CustomErrorBuilder(http.StatusBadRequest, "Income Statement Credit Harus Di isi", "Income Statement Credit Harus Di isi")
				//    }
				//    if *AdjustmentDetail.IncomeStatementDr != float64(0) {
				//    return response.CustomErrorBuilder(http.StatusBadRequest, "Tidak Dapat Masukan Income Statement Debit", "Tidak Dapat Masukan Income Statement Debit")
				//    }

				tbID, err := s.Repository.FindByTbd(ctx, &fID.ID, &v.CoaCode)
				if err != nil {
					return response.CustomErrorBuilder(http.StatusBadRequest, "Tidak Dapat Menemukan Coa di Kombinasi Konsolidasi "+v.CoaCode, "Tidak Dapat Menemukan Coa "+v.CoaCode)
				}
				tbID.Context = ctx

				AmountAjeCr := *tbID.TrialBalanceDetailEntity.AmountAjeCr - *v.IncomeStatementCr
				AmountAjeDr := *tbID.TrialBalanceDetailEntity.AmountAjeDr - *v.IncomeStatementDr
				tbID.TrialBalanceDetailEntity.AmountAjeCr = &AmountAjeCr
				tbID.TrialBalanceDetailEntity.AmountAjeDr = &AmountAjeDr

				var AmountAfterAje float64 = *tbID.AmountBeforeAje - *tbID.AmountAjeDr + *tbID.AmountAjeCr

				tbID.AmountAfterAje = &AmountAfterAje

				_, err = s.Repository.UpdateTbd(ctx, &tbID.ID, tbID)
				if err != nil {
					return err
				}

				arrAdjustmentDetail = append(arrAdjustmentDetail, v)
			}
			if findCoaGroup.Name == "Other Comprehensive Income" {

				//    if *AdjustmentDetail.BalanceSheetDr != float64(0) {
				//    return response.CustomErrorBuilder(http.StatusBadRequest, "Tidak Dapat Masukan Balance Sheet Debit", "Tidak Dapat Masukan Balance Sheet Debit")
				//    }
				//    if *AdjustmentDetail.BalanceSheetCr != float64(0) {
				//    return response.CustomErrorBuilder(http.StatusBadRequest, "Tidak Dapat Masukan Balance Sheet Credit", "Tidak Dapat Masukan Balance Sheet Credit")
				//    }
				//    if *AdjustmentDetail.IncomeStatementCr == float64(0) {
				//    return response.CustomErrorBuilder(http.StatusBadRequest, "Income Statement Credit Harus Di isi", "Income Statement Credit Harus Di isi")
				//    }
				//    if *AdjustmentDetail.IncomeStatementDr != float64(0) {
				//    return response.CustomErrorBuilder(http.StatusBadRequest, "Tidak Dapat Masukan Income Statement Debit", "Tidak Dapat Masukan Income Statement Debit")
				//    }

				tbID, err := s.Repository.FindByTbd(ctx, &fID.ID, &v.CoaCode)
				if err != nil {
					return response.CustomErrorBuilder(http.StatusBadRequest, "Tidak Dapat Menemukan Coa di Kombinasi Konsolidasi "+v.CoaCode, "Tidak Dapat Menemukan Coa "+v.CoaCode)
				}
				tbID.Context = ctx

				AmountAjeCr := *tbID.TrialBalanceDetailEntity.AmountAjeCr - *v.IncomeStatementCr
				AmountAjeDr := *tbID.TrialBalanceDetailEntity.AmountAjeDr - *v.IncomeStatementDr
				tbID.TrialBalanceDetailEntity.AmountAjeCr = &AmountAjeCr
				tbID.TrialBalanceDetailEntity.AmountAjeDr = &AmountAjeDr

				var AmountAfterAje float64 = *tbID.AmountBeforeAje - *tbID.AmountAjeDr + *tbID.AmountAjeCr

				tbID.AmountAfterAje = &AmountAfterAje

				_, err = s.Repository.UpdateTbd(ctx, &tbID.ID, tbID)
				if err != nil {
					return err
				}

				arrAdjustmentDetail = append(arrAdjustmentDetail, v)
			}
			if findCoaGroup.Name == "Dampak penyesuaian proforma  atas OCI Entitas anak" {
				//    if *AdjustmentDetail.BalanceSheetDr != float64(0) {
				//    return response.CustomErrorBuilder(http.StatusBadRequest, "Tidak Dapat Masukan Balance Sheet Debit", "Tidak Dapat Masukan Balance Sheet Debit")
				//    }
				//    if *AdjustmentDetail.BalanceSheetCr != float64(0) {
				//    return response.CustomErrorBuilder(http.StatusBadRequest, "Tidak Dapat Masukan Balance Sheet Credit", "Tidak Dapat Masukan Balance Sheet Credit")
				//    }
				//    if *AdjustmentDetail.IncomeStatementCr == float64(0) {
				//    return response.CustomErrorBuilder(http.StatusBadRequest, "Income Statement Credit Harus Di isi", "Income Statement Credit Harus Di isi")
				//    }
				//    if *AdjustmentDetail.IncomeStatementDr != float64(0) {
				//    return response.CustomErrorBuilder(http.StatusBadRequest, "Tidak Dapat Masukan Income Statement Debit", "Tidak Dapat Masukan Income Statement Debit")
				//    }

				tbID, err := s.Repository.FindByTbd(ctx, &fID.ID, &v.CoaCode)
				if err != nil {
					return response.CustomErrorBuilder(http.StatusBadRequest, "Tidak Dapat Menemukan Coa di Kombinasi Konsolidasi "+v.CoaCode, "Tidak Dapat Menemukan Coa "+v.CoaCode)
				}
				tbID.Context = ctx

				AmountAjeCr := *tbID.TrialBalanceDetailEntity.AmountAjeCr - *v.IncomeStatementCr
				AmountAjeDr := *tbID.TrialBalanceDetailEntity.AmountAjeDr - *v.IncomeStatementDr 
				tbID.TrialBalanceDetailEntity.AmountAjeCr = &AmountAjeCr
				tbID.TrialBalanceDetailEntity.AmountAjeDr = &AmountAjeDr

				var AmountAfterAje float64 = *tbID.AmountBeforeAje - *tbID.AmountAjeDr + *tbID.AmountAjeCr

				tbID.AmountAfterAje = &AmountAfterAje

				_, err = s.Repository.UpdateTbd(ctx, &tbID.ID, tbID)
				if err != nil {
					return err
				}

				arrAdjustmentDetail = append(arrAdjustmentDetail, v)
			}
			if findCoaGroup.Name == "Non Controlling OCI" {
				//    if *AdjustmentDetail.BalanceSheetDr != float64(0) {
				//    return response.CustomErrorBuilder(http.StatusBadRequest, "Tidak Dapat Masukan Balance Sheet Debit", "Tidak Dapat Masukan Balance Sheet Debit")
				//    }
				//    if *AdjustmentDetail.BalanceSheetCr != float64(0) {
				//    return response.CustomErrorBuilder(http.StatusBadRequest, "Tidak Dapat Masukan Balance Sheet Credit", "Tidak Dapat Masukan Balance Sheet Credit")
				//    }
				//    if *AdjustmentDetail.IncomeStatementCr == float64(0) {
				//    return response.CustomErrorBuilder(http.StatusBadRequest, "Income Statement Credit Harus Di isi", "Income Statement Credit Harus Di isi")
				//    }
				//    if *AdjustmentDetail.IncomeStatementDr != float64(0) {
				//    return response.CustomErrorBuilder(http.StatusBadRequest, "Tidak Dapat Masukan Income Statement Debit", "Tidak Dapat Masukan Income Statement Debit")
				//    }

				tbID, err := s.Repository.FindByTbd(ctx, &fID.ID, &v.CoaCode)
				if err != nil {
					return response.CustomErrorBuilder(http.StatusBadRequest, "Tidak Dapat Menemukan Coa di Kombinasi Konsolidasi "+v.CoaCode, "Tidak Dapat Menemukan Coa "+v.CoaCode)
				}
				tbID.Context = ctx

				AmountAjeCr := *tbID.TrialBalanceDetailEntity.AmountAjeCr - *v.IncomeStatementCr
				AmountAjeDr := *tbID.TrialBalanceDetailEntity.AmountAjeDr - *v.IncomeStatementDr
				tbID.TrialBalanceDetailEntity.AmountAjeCr = &AmountAjeCr
				tbID.TrialBalanceDetailEntity.AmountAjeDr = &AmountAjeDr

				var AmountAfterAje float64 = *tbID.AmountBeforeAje - *tbID.AmountAjeDr + *tbID.AmountAjeCr

				tbID.AmountAfterAje = &AmountAfterAje

				_, err = s.Repository.UpdateTbd(ctx, &tbID.ID, tbID)
				if err != nil {
					return err
				}

				arrAdjustmentDetail = append(arrAdjustmentDetail, v)
			}
		}

		data.Context = ctx
		result, err := s.Repository.Delete(ctx, &payload.ID, &data)
		if err != nil {
			return response.ErrorBuilder(&response.ErrorConstant.InternalServerError, err)
		}
		data = *result

		criteriaFormatterDetailSumData := model.FormatterBridgesFilterModel{}
		mfa := "TRIAL-BALANCE"
		criteriaFormatterDetailSumData.Source = &mfa
		criteriaFormatterDetailSumData.TrxRefID = &fID.TrxRefID

		formatterDetailSumData, err := s.FormatterBridgesRepository.FindSummaryTB(ctx)
		if err != nil {
			return helper.ErrorHandler(err)
		}

		for _, v := range *formatterDetailSumData {
			criteriaTBDetail := model.TrialBalanceDetailFilterModel{}
			criteriaTBDetail.TrialBalanceID = &fID.TrxRefID
			criteriaTBDetail.FormatterBridgesID = &fID.ID
			if v.IsTotal != nil && *v.IsTotal && v.FxSummary != "" {
				
				tmpString := []string{"AmountBeforeAje", "AmountAjeDr", "AmountAjeCr", "AmountAfterAje"}
				tmpTotalFl := make(map[string]*float64)
				
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
						sumTBD, err := s.TrialBalanceDetailRepository.FindSummary(ctx, &vMatch, &fID.ID, v.IsCoa)
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
				dataTB, err := s.TrialBalanceDetailRepository.FindByExactCode(ctx, &criteriaTBDetail)
				if err != nil {
					return helper.ErrorHandler(err)
				}

				updateSummary := model.TrialBalanceDetailEntityModel{
					TrialBalanceDetailEntity: model.TrialBalanceDetailEntity{
						AmountAjeDr:    tmpTotalFl["AmountAjeDr"],
						AmountAjeCr:    tmpTotalFl["AmountAjeCr"],
						AmountAfterAje: tmpTotalFl["AmountAfterAje"],
					},
				}
				_, err = s.TrialBalanceDetailRepository.Update(ctx, &dataTB.ID, &updateSummary)
				if err != nil {
					return helper.ErrorHandler(err)
				}
			}
			if v.Code == "LABA_BERSIH" {
				//UPDATE CUSTOM ROW "310401004" "310402002"
				// COA 310501002 = Row 3712 --> ambil angka dari 4337
				// COA 310502002 = Row 3718 --> ambil angka dari 4342+4343
				// COA 310503002 = Row 3724 --> ambil angka dari 4345+4346
				code := "310401004"
				criteriaTBDetail := model.TrialBalanceDetailFilterModel{}
				criteriaTBDetail.FormatterBridgesID = &fID.ID
				criteriaTBDetail.Code = &code
				customRowOne, err := s.TrialBalanceDetailRepository.FindByCode(ctx, &criteriaTBDetail)
				if err != nil {
					return helper.ErrorHandler(err)
				}

				code = "310402002"
				criteriaTBDetail.Code = &code
				customRowTwo, err := s.TrialBalanceDetailRepository.FindByCode(ctx, &criteriaTBDetail)
				if err != nil {
					return helper.ErrorHandler(err)
				}

				code = "310501002"
				criteriaTBDetail.Code = &code
				customRowThree, err := s.TrialBalanceDetailRepository.FindByCode(ctx, &criteriaTBDetail)
				if err != nil {
					return helper.ErrorHandler(err)
				}

				code = "310502002"
				criteriaTBDetail.Code = &code
				customRowFour, err := s.TrialBalanceDetailRepository.FindByCode(ctx, &criteriaTBDetail)
				if err != nil {
					return helper.ErrorHandler(err)
				}

				code = "310503002"
				criteriaTBDetail.Code = &code
				customRowFive, err := s.TrialBalanceDetailRepository.FindByCode(ctx, &criteriaTBDetail)
				if err != nil {
					return helper.ErrorHandler(err)
				}

				code = "950101001" //REVALUATION FA (4337)
				criteriaTBDetail.Code = &code
				dataReFa, err := s.TrialBalanceDetailRepository.FindByCode(ctx, &criteriaTBDetail)
				if err != nil {
					return helper.ErrorHandler(err)
				}

				code = "950301001" //Financial Instrument (4342)
				criteriaTBDetail.Code = &code
				dataFinIn, err := s.TrialBalanceDetailRepository.FindByCode(ctx, &criteriaTBDetail)
				if err != nil {
					return helper.ErrorHandler(err)
				}

				code = "950301002" //Income tax relating to components of OCI (4343)
				criteriaTBDetail.Code = &code
				dataIncomeTax, err := s.TrialBalanceDetailRepository.FindByCode(ctx, &criteriaTBDetail)
				if err != nil {
					return helper.ErrorHandler(err)
				}

				code = "950401001" //Foreign Exchange (4345)
				criteriaTBDetail.Code = &code
				dataForex, err := s.TrialBalanceDetailRepository.FindByCode(ctx, &criteriaTBDetail)
				if err != nil {
					return helper.ErrorHandler(err)
				}

				code = "950401002" //Income tax relating to components of OCI (4346)
				criteriaTBDetail.Code = &code
				dataIncomeTax2, err := s.TrialBalanceDetailRepository.FindByCode(ctx, &criteriaTBDetail)
				if err != nil {
					return helper.ErrorHandler(err)
				}

				code = "LABA_BERSIH"
				criteriaTBDetail.Code = &code
				dataLaba, err := s.TrialBalanceDetailRepository.FindByCode(ctx, &criteriaTBDetail)
				if err != nil {
					return helper.ErrorHandler(err)
				}

				code = "TOTAL_PENGHASILAN_KOMPREHENSIF_LAIN~BS"
				criteriaTBDetail.Code = &code
				dataKomprehensif, err := s.TrialBalanceDetailRepository.FindByCode(ctx, &criteriaTBDetail)
				if err != nil {
					return helper.ErrorHandler(err)
				}

				summaryCodes, err := s.TrialBalanceDetailRepository.SummaryByCodes(ctx, &fID.ID, []string{"310501002", "310502002", "310503002"})
				if err != nil {
					return helper.ErrorHandler(err)
				}

				updatedCustomRowOne := model.TrialBalanceDetailEntityModel{}
				updatedCustomRowOne.Context = ctx
				updatedCustomRowOne.AmountBeforeAje = dataLaba.AmountBeforeAje
				updatedCustomRowOne.AmountAjeCr = dataLaba.AmountAjeCr
				updatedCustomRowOne.AmountAjeDr = dataLaba.AmountAjeDr
				updatedCustomRowOne.AmountAfterAje = dataLaba.AmountAfterAje

				_, err = s.TrialBalanceDetailRepository.Update(ctx, &customRowOne.ID, &updatedCustomRowOne)
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

				_, err = s.TrialBalanceDetailRepository.Update(ctx, &customRowTwo.ID, &updateCustomRowTwo)
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

				_, err = s.TrialBalanceDetailRepository.Update(ctx, &customRowThree.ID, &updatedCustomRowThree)
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

				_, err = s.TrialBalanceDetailRepository.Update(ctx, &customRowFour.ID, &updateCustomRowFour)
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

				_, err = s.TrialBalanceDetailRepository.Update(ctx, &customRowFive.ID, &updateCustomRowFive)
				if err != nil {
					return helper.ErrorHandler(err)
				}
			}
		}
		{
			
			{
				//UPDATE CUSTOM ROW "310401004" "310402002"
				// COA 310501002 = Row 3712 --> ambil angka dari 4337
				// COA 310502002 = Row 3718 --> ambil angka dari 4342+4343
				// COA 310503002 = Row 3724 --> ambil angka dari 4345+4346
				code := "310401004"
				criteriaTBDetail := model.TrialBalanceDetailFilterModel{}
				criteriaTBDetail.FormatterBridgesID = &fID.ID
				criteriaTBDetail.Code = &code
				customRowOne, err := s.TrialBalanceDetailRepository.FindByCode(ctx, &criteriaTBDetail)
				if err != nil {
					return helper.ErrorHandler(err)
				}

				code = "310402002"
				criteriaTBDetail.Code = &code
				customRowTwo, err := s.TrialBalanceDetailRepository.FindByCode(ctx, &criteriaTBDetail)
				if err != nil {
					return helper.ErrorHandler(err)
				}

				code = "310501002"
				criteriaTBDetail.Code = &code
				customRowThree, err := s.TrialBalanceDetailRepository.FindByCode(ctx, &criteriaTBDetail)
				if err != nil {
					return helper.ErrorHandler(err)
				}

				code = "310502002"
				criteriaTBDetail.Code = &code
				customRowFour, err := s.TrialBalanceDetailRepository.FindByCode(ctx, &criteriaTBDetail)
				if err != nil {
					return helper.ErrorHandler(err)
				}

				code = "310503002"
				criteriaTBDetail.Code = &code
				customRowFive, err := s.TrialBalanceDetailRepository.FindByCode(ctx, &criteriaTBDetail)
				if err != nil {
					return helper.ErrorHandler(err)
				}

				code = "950101001" //REVALUATION FA (4337)
				criteriaTBDetail.Code = &code
				dataReFa, err := s.TrialBalanceDetailRepository.FindByCode(ctx, &criteriaTBDetail)
				if err != nil {
					return helper.ErrorHandler(err)
				}

				code = "950301001" //Financial Instrument (4342)
				criteriaTBDetail.Code = &code
				dataFinIn, err := s.TrialBalanceDetailRepository.FindByCode(ctx, &criteriaTBDetail)
				if err != nil {
					return helper.ErrorHandler(err)
				}

				code = "950301002" //Income tax relating to components of OCI (4343)
				criteriaTBDetail.Code = &code
				dataIncomeTax, err := s.TrialBalanceDetailRepository.FindByCode(ctx, &criteriaTBDetail)
				if err != nil {
					return helper.ErrorHandler(err)
				}

				code = "950401001" //Foreign Exchange (4345)
				criteriaTBDetail.Code = &code
				dataForex, err := s.TrialBalanceDetailRepository.FindByCode(ctx, &criteriaTBDetail)
				if err != nil {
					return helper.ErrorHandler(err)
				}

				code = "950401002" //Income tax relating to components of OCI (4346)
				criteriaTBDetail.Code = &code
				dataIncomeTax2, err := s.TrialBalanceDetailRepository.FindByCode(ctx, &criteriaTBDetail)
				if err != nil {
					return helper.ErrorHandler(err)
				}

				code = "LABA_BERSIH"
				criteriaTBDetail.Code = &code
				dataLaba, err := s.TrialBalanceDetailRepository.FindByCode(ctx, &criteriaTBDetail)
				if err != nil {
					return helper.ErrorHandler(err)
				}

				code = "TOTAL_PENGHASILAN_KOMPREHENSIF_LAIN~BS"
				criteriaTBDetail.Code = &code
				dataKomprehensif, err := s.TrialBalanceDetailRepository.FindByCode(ctx, &criteriaTBDetail)
				if err != nil {
					return helper.ErrorHandler(err)
				}

				summaryCodes, err := s.TrialBalanceDetailRepository.SummaryByCodes(ctx, &fID.ID, []string{"310501002", "310502002", "310503002"})
				if err != nil {
					return helper.ErrorHandler(err)
				}

				updatedCustomRowOne := model.TrialBalanceDetailEntityModel{}
				updatedCustomRowOne.Context = ctx
				updatedCustomRowOne.AmountBeforeAje = dataLaba.AmountBeforeAje
				updatedCustomRowOne.AmountAjeCr = dataLaba.AmountAjeCr
				updatedCustomRowOne.AmountAjeDr = dataLaba.AmountAjeDr
				updatedCustomRowOne.AmountAfterAje = dataLaba.AmountAfterAje

				_, err = s.TrialBalanceDetailRepository.Update(ctx, &customRowOne.ID, &updatedCustomRowOne)
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

				_, err = s.TrialBalanceDetailRepository.Update(ctx, &customRowTwo.ID, &updateCustomRowTwo)
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

				_, err = s.TrialBalanceDetailRepository.Update(ctx, &customRowThree.ID, &updatedCustomRowThree)
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

				_, err = s.TrialBalanceDetailRepository.Update(ctx, &customRowFour.ID, &updateCustomRowFour)
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

				_, err = s.TrialBalanceDetailRepository.Update(ctx, &customRowFive.ID, &updateCustomRowFive)
				if err != nil {
					return helper.ErrorHandler(err)
				}
			}
		}
		for _, v := range *formatterDetailSumData {
			criteriaTBDetail := model.TrialBalanceDetailFilterModel{}
			criteriaTBDetail.TrialBalanceID = &fID.TrxRefID
			criteriaTBDetail.FormatterBridgesID = &fID.ID

			if v.AutoSummary != nil && *v.AutoSummary {
				code := fmt.Sprintf("%s_Subtotal", v.Code)
				criteriaTBDetail.Code = &code
				mfadetailsum, _, err := s.TrialBalanceDetailRepository.Find(ctx, &criteriaTBDetail, &abstraction.Pagination{})
				if err != nil {
					return helper.ErrorHandler(err)
				}

				for _, a := range *mfadetailsum {
					sumTBD, err := s.TrialBalanceDetailRepository.FindSummary(ctx, &v.Code, &fID.ID, v.IsCoa)
					if err != nil {
						return helper.ErrorHandler(err)
					}
					updateSummary := model.TrialBalanceDetailEntityModel{
						TrialBalanceDetailEntity: model.TrialBalanceDetailEntity{
							AmountAjeDr:    sumTBD.AmountAjeDr,
							AmountAjeCr:    sumTBD.AmountAjeCr,
							AmountAfterAje: sumTBD.AmountAfterAje,
						},
					}
					_, err = s.TrialBalanceDetailRepository.Update(ctx, &a.ID, &updateSummary)
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
						sumTBD, err := s.TrialBalanceDetailRepository.FindSummary(ctx, &vMatch, &fID.ID, v.IsCoa)
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
				dataTB, err := s.TrialBalanceDetailRepository.FindByExactCode(ctx, &criteriaTBDetail)
				if err != nil {
					return helper.ErrorHandler(err)
				}

				updateSummary := model.TrialBalanceDetailEntityModel{
					TrialBalanceDetailEntity: model.TrialBalanceDetailEntity{
						AmountAjeDr:    tmpTotalFl["AmountAjeDr"],
						AmountAjeCr:    tmpTotalFl["AmountAjeCr"],
						AmountAfterAje: tmpTotalFl["AmountAfterAje"],
					},
				}
				_, err = s.TrialBalanceDetailRepository.Update(ctx, &dataTB.ID, &updateSummary)
				if err != nil {
					return helper.ErrorHandler(err)
				}
			}
		}
		return nil
	}); err != nil {
		return &dto.AdjustmentDeleteResponse{}, err
	}
	result := &dto.AdjustmentDeleteResponse{
		AdjustmentEntityModel: data,
	}
	return result, nil
}

func (s *service) Export(ctx *abstraction.Context, payload *dto.AdjustmentExportRequest) (*dto.AdjustmentExportResponse, error) {

	datas, err := s.Repository.Export(ctx, &payload.AdjustmentID)
	if err != nil {
		return &dto.AdjustmentExportResponse{}, response.ErrorBuilder(&response.ErrorConstant.InternalServerError, err)
	}

	datePeriod, err := time.Parse(time.RFC3339, datas.Period)
	if err != nil {
		return nil, err
	}

	f := excelize.NewFile()
	sheet := f.GetSheetName(f.GetActiveSheetIndex())
	arrStyleColWidth := []map[string]interface{}{
		{"COL": "A", "WIDTH": 6.43},
		{"COL": "B", "WIDTH": 15.38},
		{"COL": "C", "WIDTH": 7.25},
		{"COL": "D", "WIDTH": 2.14},
		{"COL": "E", "WIDTH": 57.45},
		{"COL": "F", "WIDTH": 7.25},
		{"COL": "G", "WIDTH": 15.38},
		{"COL": "H", "WIDTH": 15.38},
		{"COL": "I", "WIDTH": 15.38},
		{"COL": "J", "WIDTH": 15.38},
	}

	for _, v := range arrStyleColWidth {
		tmpColWidth := fmt.Sprintf("%f", v["WIDTH"])
		colWidth, _ := strconv.ParseFloat(tmpColWidth, 64)
		f.SetColWidth(sheet, fmt.Sprintf("%s", v["COL"]), fmt.Sprintf("%s", v["COL"]), colWidth)
	}

	styleDefault, err := f.NewStyle(&excelize.Style{
		Font: &excelize.Font{
			Family: "Arial",
			Size:   10,
		},
	})

	f.SetColStyle(sheet, "A:Z", styleDefault)

	// stylingBorderLeftOnly, err := f.NewStyle(&excelize.Style{
	// 	Border: []excelize.Border{
	// 		{Type: "left", Color: "000000", Style: 2},
	// 	},
	// })

	stylingBorderLROnly, err := f.NewStyle(&excelize.Style{
		Border: []excelize.Border{
			{Type: "right", Color: "000000", Style: 2},
			{Type: "left", Color: "000000", Style: 2},
		},
	})

	stylingBorderAll, err := f.NewStyle(&excelize.Style{
		Border: []excelize.Border{
			{Type: "top", Color: "000000", Style: 2},
			{Type: "right", Color: "000000", Style: 2},
			{Type: "left", Color: "000000", Style: 2},
			{Type: "bottom", Color: "000000", Style: 2},
		},
	})

	stylingHeader, err := f.NewStyle(&excelize.Style{
		Font: &excelize.Font{
			Bold: true,
		},
		Border: []excelize.Border{
			{Type: "top", Color: "000000", Style: 2},
			{Type: "bottom", Color: "000000", Style: 2},
			{Type: "left", Color: "000000", Style: 2},
			{Type: "right", Color: "000000", Style: 2},
		},
		Fill: excelize.Fill{
			Type:    "pattern",
			Pattern: 1,
			Color:   []string{"#fac090"},
		},
		Alignment: &excelize.Alignment{
			Horizontal: "center",
			Vertical:   "center",
		},
	})

	f.MergeCell(sheet, "A6", "A7")
	f.MergeCell(sheet, "B6", "B7")
	f.MergeCell(sheet, "C6", "C7")
	f.MergeCell(sheet, "D6", "E7")
	f.MergeCell(sheet, "F6", "F7")
	f.MergeCell(sheet, "G6", "H6")
	f.MergeCell(sheet, "I6", "j6")

	stylingCurrency, err := f.NewStyle(&excelize.Style{
		NumFmt: 7,
		Border: []excelize.Border{
			{Type: "left", Color: "000000", Style: 1},
			{Type: "right", Color: "000000", Style: 1},
		},
	})

	stylingSubTotal, err := f.NewStyle(&excelize.Style{
		Border: []excelize.Border{
			{Type: "top", Color: "000000", Style: 2},
			{Type: "bottom", Color: "000000", Style: 2},
			{Type: "left", Color: "000000", Style: 2},
			{Type: "right", Color: "000000", Style: 2},
		},

		Fill: excelize.Fill{
			Type:    "pattern",
			Pattern: 1,
			Color:   []string{"#00ff00"},
		},

		Font: &excelize.Font{
			Bold: true,
		},
	})

	stylingSubTotal2, err := f.NewStyle(&excelize.Style{
		Border: []excelize.Border{
			{Type: "top", Color: "000000", Style: 2},
			{Type: "bottom", Color: "000000", Style: 2},
			{Type: "left", Color: "000000", Style: 2},
			{Type: "right", Color: "000000", Style: 2},
		},

		Fill: excelize.Fill{
			Type:    "pattern",
			Pattern: 1,
			Color:   []string{"#99ff33"},
		},

		Font: &excelize.Font{
			Bold: true,
		},
	})
	f.SetCellValue(sheet, "A2", "Company")
	f.SetCellValue(sheet, "B2", ": "+datas.Company.Name)
	f.SetCellValue(sheet, "A3", "Date")
	f.SetCellValue(sheet, "B3", ": "+datePeriod.Format("02-Jan-06"))
	f.SetCellValue(sheet, "A4", "Subject")
	f.SetCellValue(sheet, "B4", ": AJE")

	f.SetCellStyle(sheet, "A6", "J7", stylingHeader)
	f.SetCellValue(sheet, "A6", "NO")
	f.SetCellValue(sheet, "B6", "COA")
	f.SetCellValue(sheet, "C6", "AJE")
	f.SetCellValue(sheet, "D6", "DESCRIPTION")
	f.SetCellValue(sheet, "F6", "WP")
	f.SetCellValue(sheet, "F7", "REFF")
	f.SetCellValue(sheet, "G6", "Balance Sheet")
	f.SetCellValue(sheet, "G7", "DR")
	f.SetCellValue(sheet, "H7", "CR")
	f.SetCellValue(sheet, "I6", "Income Stat")
	f.SetCellValue(sheet, "I7", "DR")
	f.SetCellValue(sheet, "J7", "CR")

	row := 8

	rowBefore := row

	for i, v := range datas.AdjustmentDetail {

		f.SetCellValue(sheet, fmt.Sprintf("A%d", row), i+1)
		f.SetCellValue(sheet, fmt.Sprintf("B%d", row), v.CoaCode)
		f.SetCellValue(sheet, fmt.Sprintf("E%d", row), *v.Description)
		f.SetCellValue(sheet, fmt.Sprintf("E%d", row+1), *v.Note)
		if v.BalanceSheetDr != nil && *v.BalanceSheetDr != 0 {
			f.SetCellValue(sheet, fmt.Sprintf("G%d", row), *v.BalanceSheetDr)
		}
		if v.BalanceSheetCr != nil && *v.BalanceSheetCr != 0 {
			f.SetCellValue(sheet, fmt.Sprintf("H%d", row), *v.BalanceSheetCr)
		}
		if v.IncomeStatementDr != nil && *v.IncomeStatementDr != 0 {
			f.SetCellValue(sheet, fmt.Sprintf("I%d", row), *v.IncomeStatementDr)
		}
		if v.IncomeStatementCr != nil && *v.IncomeStatementCr != 0 {
			f.SetCellValue(sheet, fmt.Sprintf("J%d", row), *v.IncomeStatementCr)
		}

		row += 2
	}

	f.SetCellStyle(sheet, "A8", fmt.Sprintf("B%d", row), stylingBorderLROnly)
	f.SetCellStyle(sheet, "C8", fmt.Sprintf("C%d", row), stylingBorderLROnly)
	f.SetCellStyle(sheet, "F8", fmt.Sprintf("F%d", row), stylingBorderLROnly)
	f.SetCellStyle(sheet, "G8", fmt.Sprintf("J%d", row), stylingCurrency)
	f.SetCellStyle(sheet, fmt.Sprintf("A%d", row), fmt.Sprintf("F%d", row+1), stylingBorderAll)
	f.MergeCell(sheet, fmt.Sprintf("D%d", row), fmt.Sprintf("E%d", row))
	f.MergeCell(sheet, fmt.Sprintf("D%d", row+1), fmt.Sprintf("E%d", row+1))

	f.SetSheetFormatPr(sheet, excelize.DefaultRowHeight(12.85))

	f.SetCellValue(sheet, fmt.Sprintf("D%d", row), "Total PM")
	f.SetCellFormula(sheet, fmt.Sprintf("G%d", row), fmt.Sprintf("=SUM(G%d:G%d)", rowBefore, row-1))
	f.SetCellFormula(sheet, fmt.Sprintf("H%d", row), fmt.Sprintf("=SUM(H%d:H%d)", rowBefore, row-1))
	f.SetCellFormula(sheet, fmt.Sprintf("I%d", row), fmt.Sprintf("=SUM(I%d:I%d)", rowBefore, row-1))
	f.SetCellFormula(sheet, fmt.Sprintf("J%d", row), fmt.Sprintf("=SUM(J%d:J%d)", rowBefore, row-1))
	f.SetCellFormula(sheet, fmt.Sprintf("I%d", row+1), fmt.Sprintf("=G%d+I%d", row, row))
	f.SetCellFormula(sheet, fmt.Sprintf("J%d", row+1), fmt.Sprintf("=H%d+J%d", row, row))
	f.SetCellFormula(sheet, fmt.Sprintf("J%d", row+2), fmt.Sprintf("=I%d-J%d", row+1, row+1))

	f.SetCellStyle(sheet, fmt.Sprintf("D%d", row), fmt.Sprintf("J%d", row), stylingSubTotal)
	f.SetCellStyle(sheet, fmt.Sprintf("G%d", row+1), fmt.Sprintf("J%d", row+1), stylingSubTotal2)

	f.SetDefaultFont("Arial")

	tmpFolder := fmt.Sprintf("assets/%d", ctx.Auth.ID)
	_, err = os.Stat(tmpFolder)
	if err != nil {
		if os.IsNotExist(err) {
			if err := os.MkdirAll(tmpFolder, os.ModePerm); err != nil {
				return nil, err
			}
		} else {
			return nil, err
		}
	}
	period := datePeriod.Format("2006-01-02")
	fileName := fmt.Sprintf("Adjustment_%s.xlsx", period)
	fileLoc := fmt.Sprintf("assets/%d/%s", ctx.Auth.ID, fileName)
	err = f.SaveAs(fileLoc)
	if err != nil {
		return nil, err
	}

	result := &dto.AdjustmentExportResponse{
		FileName: fileName,
		Path:     fileLoc,
	}
	return result, nil
}

func (s *service) GetVersion(ctx *abstraction.Context, payload *dto.GetVersionRequest) (*dto.GetVersionResponse, error) {
	filter := model.AdjustmentFilterModel{
		CompanyCustomFilter: model.CompanyCustomFilter{
			CompanyID:          payload.CompanyID,
			ArrCompanyID:       payload.ArrCompanyID,
			ArrCompanyString:   payload.ArrCompanyString,
			ArrCompanyOperator: payload.ArrCompanyOperator,
		},
	}
	filter.Period = payload.Period
	data, err := s.Repository.GetVersion(ctx, &filter)
	if err != nil {
		return &dto.GetVersionResponse{}, response.ErrorBuilder(&response.ErrorConstant.InternalServerError, err)
	}
	result := &dto.GetVersionResponse{
		Data: *data,
	}
	return result, nil
}
