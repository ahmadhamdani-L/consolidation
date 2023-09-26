package jpm

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
	Repository          repository.Jpm
	JpmDetailRepository repository.JpmDetail
	Db                  *gorm.DB
}

type Service interface {
	Find(ctx *abstraction.Context, payload *dto.JpmGetRequest) (*dto.JpmGetResponse, error)
	FindByID(ctx *abstraction.Context, payload *dto.JpmGetByIDRequest) (*dto.JpmGetByIDResponse, error)
	Create(ctx *abstraction.Context, payload *dto.JpmCreateRequest) (*dto.JpmCreateResponse, error)
	Update(ctx *abstraction.Context, payload *dto.JpmUpdateRequest) (*dto.JpmUpdateResponse, error)
	Delete(ctx *abstraction.Context, payload *dto.JpmDeleteRequest) (*dto.JpmDeleteResponse, error)
	Export(ctx *abstraction.Context, payload *dto.JpmExportRequest) (*dto.JpmExportResponse, error)
	GetVersion(ctx *abstraction.Context, payload *dto.GetVersionRequest) (*dto.GetVersionResponse, error)
}

func NewService(f *factory.Factory) *service {
	repository := f.JpmRepository
	JpmDetailRepository := f.JpmDetailRepository

	db := f.Db
	return &service{
		Repository:          repository,
		JpmDetailRepository: JpmDetailRepository,
		Db:                  db,
	}
}

func (s *service) Find(ctx *abstraction.Context, payload *dto.JpmGetRequest) (*dto.JpmGetResponse, error) {
	data, info, err := s.Repository.Find(ctx, &payload.JpmFilterModel, &payload.Pagination)
	if err != nil {
		return &dto.JpmGetResponse{}, response.ErrorBuilder(&response.ErrorConstant.InternalServerError, err)
	}
	result := &dto.JpmGetResponse{
		Datas:          *data,
		PaginationInfo: *info,
	}
	return result, nil
}

func (s *service) FindByID(ctx *abstraction.Context, payload *dto.JpmGetByIDRequest) (*dto.JpmGetByIDResponse, error) {
	data, err := s.Repository.FindByID(ctx, &payload.ID)
	if err != nil {
		return &dto.JpmGetByIDResponse{}, response.ErrorBuilder(&response.ErrorConstant.InternalServerError, err)
	}
	result := &dto.JpmGetByIDResponse{
		JpmEntityModel: *data,
	}
	return result, nil
}
func (s *service) Create(ctx *abstraction.Context, payload *dto.JpmCreateRequest) (*dto.JpmCreateResponse, error) {
	var data model.JpmEntityModel
	var datas []model.JpmDetailEntity

	if err := trxmanager.New(s.Db).WithTrx(ctx, func(ctx *abstraction.Context) error {
		var getAje model.JpmFilterModel
		getAje.CompanyID = &payload.CompanyID
		getAje.Period = &payload.Period
		getAje.ConsolidationID = &payload.ConsolidationID
		

		getJpm, err := s.Repository.Get(ctx, &getAje)
		if err != nil {
			fmt.Println(err)
			return err
		}

		data.Context = ctx
		data.JpmEntity = payload.JpmEntity
		data.JpmEntity.Status = 1
		// data.JpmEntity.Period = lastOfMonth.Format("2006-01-02")
		data.JpmEntity.TrxNumber = "JPM" + "#" + strconv.Itoa(len(*getJpm)+1)
		result, err := s.Repository.Create(ctx, &data)
		if err != nil {
			return response.ErrorBuilder(&response.ErrorConstant.UnprocessableEntity, err)
		}
		data = *result

		//create JpmDetail
		datas = payload.JpmDetail
		var arrJpmDetail []model.JpmDetailEntity
		for _, v := range datas {
			JpmDetail := model.JpmDetailEntityModel{
				Context:         ctx,
				JpmDetailEntity: v,
			}
			JpmDetail.JpmID = data.ID
			JpmDetail.ReffNumber = &data.TrxNumber
			_, err := s.JpmDetailRepository.Create(ctx, &JpmDetail)
			if err != nil {
				return err
			}
			if v.CoaCode == "310401004" || v.CoaCode == "310501002" || v.CoaCode == "310502002" || v.CoaCode == "310503002" || v.CoaCode == "310402002" {
				return response.CustomErrorBuilder(http.StatusBadRequest, response.E_BAD_REQUEST, "Coa Tersebut Tidak Dapat Melalukan Jurnal "+JpmDetail.CoaCode)
			}
			findCoa, err := s.Repository.FindByCoa(ctx, &JpmDetail.CoaCode)
			if err != nil {
				return response.CustomErrorBuilder(http.StatusBadRequest, "Tidak Dapat Menemukan Coa "+JpmDetail.CoaCode, "Tidak Dapat Menemukan Coa "+JpmDetail.CoaCode)
			}

			findCoaType, err := s.Repository.FindByCoaType(ctx, &findCoa.CoaTypeID)
			if err != nil {
				return response.CustomErrorBuilder(http.StatusBadRequest, "Tidak Dapat Menemukan Coa di Type Coa "+JpmDetail.CoaCode, "Tidak Dapat Menemukan Coa di Type Coa "+JpmDetail.CoaCode)
			}
			findCoaGroup, err := s.Repository.FindByCoaGroup(ctx, &findCoaType.CoaGroupID)
			if err != nil {
				return response.CustomErrorBuilder(http.StatusBadRequest, "Tidak Dapat Menemukan Coa di Group coa "+JpmDetail.CoaCode, "Tidak Dapat Menemukan Coa di Group coa "+JpmDetail.CoaCode)
			}

			if findCoaGroup.Name == "ASET" {
				// if *JpmDetail.BalanceSheetDr == float64(0) {
				// 	return response.CustomErrorBuilder(http.StatusBadRequest, "Balance Sheet Debit Harus Di isi", "Balance Sheet Debit Harus Di isi")
				// }
				// if *JpmDetail.BalanceSheetCr > float64(0) {
				// 	return response.CustomErrorBuilder(http.StatusBadRequest, "Tidak Dapat Masukan Balance Sheet Credit", "Tidak Dapat Masukan Balance Sheet Credit")
				// }
				// if *JpmDetail.IncomeStatementCr > float64(0) {
				// 	return response.CustomErrorBuilder(http.StatusBadRequest, "Tidak Dapat Masukan Income Statement Credit", "Tidak Dapat Masukan Income Statement Credit")
				// }
				// if *JpmDetail.IncomeStatementDr > float64(0) {
				// 	return response.CustomErrorBuilder(http.StatusBadRequest, "Tidak Dapat Masukan Income Statement Debit", "Tidak Dapat Masukan Income Statement Debit")
				// }

				tbID, err := s.Repository.FindByTbd(ctx, &data.ConsolidationID, &JpmDetail.CoaCode)
				if err != nil {
					return response.CustomErrorBuilder(http.StatusBadRequest, "Tidak Dapat Menemukan Coa di Kombinasi Konsolidasi "+JpmDetail.CoaCode, "Tidak Dapat Menemukan Coa "+JpmDetail.CoaCode)
				}
				tbID.Context = ctx

				// AmountAjeCr := *v.BalanceSheetCr + *tbID.ConsolidationDetailEntity.AmountJcteCr
				AmountJpmDr := *v.BalanceSheetDr + *tbID.ConsolidationDetailEntity.AmountJpmDr
				AmountJpmCr := *v.BalanceSheetCr + *tbID.ConsolidationDetailEntity.AmountJpmCr
				tbID.ConsolidationDetailEntity.AmountJpmCr = &AmountJpmCr
				tbID.ConsolidationDetailEntity.AmountJpmDr = &AmountJpmDr
				// amount After Jpm
				AmountAfterJpm := *tbID.AmountBeforeJpm + *tbID.AmountJpmDr - *tbID.AmountJpmCr
				tbID.AmountAfterJpm = &AmountAfterJpm
				// amount After Jcte
				AmountAfterJcte := *tbID.AmountAfterJpm + *tbID.AmountJcteDr - *tbID.AmountJcteCr
				tbID.AmountAfterJcte = &AmountAfterJcte

				findConsoleBridge, err := s.Repository.FindByConsoleBridge(ctx, &result.ConsolidationID)
				if err != nil {
					return response.CustomErrorBuilder(http.StatusBadRequest, "Tidak Dapat Menemukan Anak Usaha ", "Tidak Dapat Menemukan Anak Usaha ")
				}

				var dataAmountAnakUsahaa []float64
				for _, cb := range *findConsoleBridge {
					findConsoleBridgeDetail, err := s.Repository.FindByConsoleBridgeDetail(ctx, &cb.ID, &v.CoaCode)
					if err != nil {
						return response.CustomErrorBuilder(http.StatusBadRequest, "Tidak Dapat Menemukan Amount Anak Usaha ", "Tidak Dapat Menemukan Amount Anak Usaha ")
					}
					dataAmountAnakUsahaa = append(dataAmountAnakUsahaa, findConsoleBridgeDetail.Amount)
				}
				var sumAmountAnakUsaha float64 = 0.0

				for i := 0; i < len(dataAmountAnakUsahaa); i++ {
					sumAmountAnakUsaha += dataAmountAnakUsahaa[i]
				}

				var combineUnaudited float64 = *tbID.AmountAfterJcte + sumAmountAnakUsaha
				tbID.AmountCombineSubsidiary = &combineUnaudited

				var amountConsole float64 = *tbID.AmountCombineSubsidiary + *tbID.AmountJelimDr - *tbID.AmountJelimCr

				tbID.AmountConsole = &amountConsole

				_, err = s.Repository.UpdateTbd(ctx, &tbID.ID, tbID)
				if err != nil {
					return err
				}

				arrJpmDetail = append(arrJpmDetail, v)
			}
			if findCoaGroup.Name == "Liabilitas" {
				// if *JpmDetail.BalanceSheetDr > float64(0) {
				// 	return response.CustomErrorBuilder(http.StatusBadRequest, "Tidak Dapat Masukan Balance Sheet Debit", "Tidak Dapat Masukan Balance Sheet Debit")
				// }
				// if *JpmDetail.BalanceSheetCr == float64(0) {
				// 	return response.CustomErrorBuilder(http.StatusBadRequest, "Balance Sheet Credit Harus Di isi", "Balance Sheet Credit Harus Di isi")
				// }
				// if *JpmDetail.IncomeStatementCr > float64(0) {
				// 	return response.CustomErrorBuilder(http.StatusBadRequest, "Tidak Dapat Masukan Income Statement Credit", "Tidak Dapat Masukan Income Statement Credit")
				// }
				// if *JpmDetail.IncomeStatementDr > float64(0) {
				// 	return response.CustomErrorBuilder(http.StatusBadRequest, "Tidak Dapat Masukan Income Statement Debit", "Tidak Dapat Masukan Income Statement Debit")
				// }

				tbID, err := s.Repository.FindByTbd(ctx, &data.ConsolidationID, &JpmDetail.CoaCode)
				if err != nil {
					return response.CustomErrorBuilder(http.StatusBadRequest, "Tidak Dapat Menemukan Coa di Kombinasi Konsolidasi "+JpmDetail.CoaCode, "Tidak Dapat Menemukan Coa "+JpmDetail.CoaCode)
				}
				tbID.Context = ctx

				AmountJpmCr := *v.BalanceSheetCr + *tbID.ConsolidationDetailEntity.AmountJpmCr
				AmountAjeDr := *v.BalanceSheetDr + *tbID.ConsolidationDetailEntity.AmountJpmDr
				tbID.ConsolidationDetailEntity.AmountJpmCr = &AmountJpmCr
				tbID.ConsolidationDetailEntity.AmountJpmDr = &AmountAjeDr
				// amount After Jpm
				AmountAfterJpm := *tbID.AmountBeforeJpm - *tbID.AmountJpmDr + *tbID.AmountJpmCr
				tbID.AmountAfterJpm = &AmountAfterJpm
				// amount After Jcte
				AmountAfterJcte := *tbID.AmountAfterJpm - *tbID.AmountJcteDr + *tbID.AmountJcteCr
				tbID.AmountAfterJcte = &AmountAfterJcte

				findConsoleBridge, err := s.Repository.FindByConsoleBridge(ctx, &result.ConsolidationID)
				if err != nil {
					return response.CustomErrorBuilder(http.StatusBadRequest, "Tidak Dapat Menemukan Anak Usaha ", "Tidak Dapat Menemukan Anak Usaha ")
				}

				var dataAmountAnakUsahaa []float64
				for _, cb := range *findConsoleBridge {
					findConsoleBridgeDetail, err := s.Repository.FindByConsoleBridgeDetail(ctx, &cb.ID, &v.CoaCode)
					_, err = s.Repository.UpdateTbd(ctx, &tbID.ID, tbID)
					if err != nil {
						return response.CustomErrorBuilder(http.StatusBadRequest, "Tidak Dapat Menemukan Amount Anak Usaha ", "Tidak Dapat Menemukan Amount Anak Usaha ")
					}
					dataAmountAnakUsahaa = append(dataAmountAnakUsahaa, findConsoleBridgeDetail.Amount)
				}
				var sumAmountAnakUsaha float64 = 0.0

				for i := 0; i < len(dataAmountAnakUsahaa); i++ {
					sumAmountAnakUsaha += dataAmountAnakUsahaa[i]
				}

				var combineUnaudited float64 = *tbID.AmountAfterJcte + sumAmountAnakUsaha
				tbID.AmountCombineSubsidiary = &combineUnaudited

				var amountConsole float64 = *tbID.AmountCombineSubsidiary - *tbID.AmountJelimDr + *tbID.AmountJelimCr

				tbID.AmountConsole = &amountConsole

				_, err = s.Repository.UpdateTbd(ctx, &tbID.ID, tbID)
				if err != nil {
					return err
				}
				arrJpmDetail = append(arrJpmDetail, v)
			}
			if findCoaGroup.Name == "EKUITAS" {
				// if *JpmDetail.BalanceSheetDr > float64(0) {
				// 	return response.CustomErrorBuilder(http.StatusBadRequest, "Tidak Dapat Masukan Balance Sheet Debit", "Tidak Dapat Masukan Balance Sheet Debit")
				// }
				// if *JpmDetail.BalanceSheetCr == float64(0) {
				// 	return response.CustomErrorBuilder(http.StatusBadRequest, "Balance Sheet Credit Harus Di isi", "Balance Sheet Credit Harus Di isi")
				// }
				// if *JpmDetail.IncomeStatementCr > float64(0) {
				// 	return response.CustomErrorBuilder(http.StatusBadRequest, "Tidak Dapat Masukan Income Statement Credit", "Tidak Dapat Masukan Income Statement Credit")
				// }
				// if *JpmDetail.IncomeStatementDr > float64(0) {
				// 	return response.CustomErrorBuilder(http.StatusBadRequest, "Tidak Dapat Masukan Income Statement Debit", "Tidak Dapat Masukan Income Statement Debit")
				// }

				tbID, err := s.Repository.FindByTbd(ctx, &data.ConsolidationID, &JpmDetail.CoaCode)
				if err != nil {
					return response.CustomErrorBuilder(http.StatusBadRequest, "Tidak Dapat Menemukan Coa di Kombinasi Konsolidasi "+JpmDetail.CoaCode, "Tidak Dapat Menemukan Coa "+JpmDetail.CoaCode)
				}
				tbID.Context = ctx

				AmountJpmCr := *v.BalanceSheetCr + *tbID.ConsolidationDetailEntity.AmountJpmCr
				AmountAjeDr := *v.BalanceSheetDr + *tbID.ConsolidationDetailEntity.AmountJpmDr
				tbID.ConsolidationDetailEntity.AmountJpmCr = &AmountJpmCr
				tbID.ConsolidationDetailEntity.AmountJpmDr = &AmountAjeDr
				// amount After Jpm
				AmountAfterJpm := *tbID.AmountBeforeJpm - *tbID.AmountJpmDr + *tbID.AmountJpmCr
				tbID.AmountAfterJpm = &AmountAfterJpm
				// amount After Jcte
				AmountAfterJcte := *tbID.AmountAfterJpm - *tbID.AmountJcteDr + *tbID.AmountJcteCr
				tbID.AmountAfterJcte = &AmountAfterJcte

				findConsoleBridge, err := s.Repository.FindByConsoleBridge(ctx, &result.ConsolidationID)
				if err != nil {
					return response.CustomErrorBuilder(http.StatusBadRequest, "Tidak Dapat Menemukan Anak Usaha ", "Tidak Dapat Menemukan Anak Usaha ")
				}

				var dataAmountAnakUsahaa []float64
				for _, cb := range *findConsoleBridge {
					findConsoleBridgeDetail, err := s.Repository.FindByConsoleBridgeDetail(ctx, &cb.ID, &v.CoaCode)
					_, err = s.Repository.UpdateTbd(ctx, &tbID.ID, tbID)
					if err != nil {
						return response.CustomErrorBuilder(http.StatusBadRequest, "Tidak Dapat Menemukan Amount Anak Usaha ", "Tidak Dapat Menemukan Amount Anak Usaha ")
					}
					dataAmountAnakUsahaa = append(dataAmountAnakUsahaa, findConsoleBridgeDetail.Amount)
				}
				var sumAmountAnakUsaha float64 = 0.0

				for i := 0; i < len(dataAmountAnakUsahaa); i++ {
					sumAmountAnakUsaha += dataAmountAnakUsahaa[i]
				}

				var combineUnaudited float64 = *tbID.AmountAfterJcte + sumAmountAnakUsaha
				tbID.AmountCombineSubsidiary = &combineUnaudited

				var amountConsole float64 = *tbID.AmountCombineSubsidiary - *tbID.AmountJelimDr + *tbID.AmountJelimCr

				tbID.AmountConsole = &amountConsole

				_, err = s.Repository.UpdateTbd(ctx, &tbID.ID, tbID)
				if err != nil {
					return err
				}
				arrJpmDetail = append(arrJpmDetail, v)
			}
			if findCoaGroup.Name == "Pendapatan" {
				// if *JpmDetail.BalanceSheetDr != float64(0) {
				// 	return response.CustomErrorBuilder(http.StatusBadRequest, "Tidak Dapat Masukan Balance Sheet Debit", "Tidak Dapat Masukan Balance Sheet Debit")
				// }
				// if *JpmDetail.BalanceSheetCr != float64(0) {
				// 	return response.CustomErrorBuilder(http.StatusBadRequest, "Tidak Dapat Masukan Balance Sheet Credit", "Tidak Dapat Masukan Balance Sheet Credit")
				// }
				// if *JpmDetail.IncomeStatementCr == float64(0) {
				// 	return response.CustomErrorBuilder(http.StatusBadRequest, "Income Statement Credit Harus Di isi", "Income Statement Credit Harus Di isi")
				// }
				// if *JpmDetail.IncomeStatementDr != float64(0) {
				// 	return response.CustomErrorBuilder(http.StatusBadRequest, "Tidak Dapat Masukan Income Statement Debit", "Tidak Dapat Masukan Income Statement Debit")
				// }

				tbID, err := s.Repository.FindByTbd(ctx, &data.ConsolidationID, &JpmDetail.CoaCode)
				if err != nil {
					return response.CustomErrorBuilder(http.StatusBadRequest, "Tidak Dapat Menemukan Coa di Kombinasi Konsolidasi "+JpmDetail.CoaCode, "Tidak Dapat Menemukan Coa "+JpmDetail.CoaCode)
				}
				tbID.Context = ctx

				AmountJpmCr := *v.IncomeStatementCr + *tbID.ConsolidationDetailEntity.AmountJpmCr
				AmountAjeDr := *v.IncomeStatementDr + *tbID.ConsolidationDetailEntity.AmountJpmDr
				tbID.ConsolidationDetailEntity.AmountJpmCr = &AmountJpmCr
				tbID.ConsolidationDetailEntity.AmountJpmDr = &AmountAjeDr
				// amount After Jpm
				AmountAfterJpm := *tbID.AmountBeforeJpm - *tbID.AmountJpmDr + *tbID.AmountJpmCr
				tbID.AmountAfterJpm = &AmountAfterJpm
				// amount After Jcte
				AmountAfterJcte := *tbID.AmountAfterJpm - *tbID.AmountJcteDr + *tbID.AmountJcteCr
				tbID.AmountAfterJcte = &AmountAfterJcte

				findConsoleBridge, err := s.Repository.FindByConsoleBridge(ctx, &result.ConsolidationID)
				if err != nil {
					return response.CustomErrorBuilder(http.StatusBadRequest, "Tidak Dapat Menemukan Anak Usaha ", "Tidak Dapat Menemukan Anak Usaha ")
				}

				var dataAmountAnakUsahaa []float64
				for _, cb := range *findConsoleBridge {
					findConsoleBridgeDetail, err := s.Repository.FindByConsoleBridgeDetail(ctx, &cb.ID, &v.CoaCode)
					_, err = s.Repository.UpdateTbd(ctx, &tbID.ID, tbID)
					if err != nil {
						return response.CustomErrorBuilder(http.StatusBadRequest, "Tidak Dapat Menemukan Amount Anak Usaha ", "Tidak Dapat Menemukan Amount Anak Usaha ")
					}
					dataAmountAnakUsahaa = append(dataAmountAnakUsahaa, findConsoleBridgeDetail.Amount)
				}
				var sumAmountAnakUsaha float64 = 0.0

				for i := 0; i < len(dataAmountAnakUsahaa); i++ {
					sumAmountAnakUsaha += dataAmountAnakUsahaa[i]
				}

				var combineUnaudited float64 = *tbID.AmountAfterJcte + sumAmountAnakUsaha
				tbID.AmountCombineSubsidiary = &combineUnaudited

				var amountConsole float64 = *tbID.AmountCombineSubsidiary - *tbID.AmountJelimDr + *tbID.AmountJelimCr

				tbID.AmountConsole = &amountConsole

				_, err = s.Repository.UpdateTbd(ctx, &tbID.ID, tbID)
				if err != nil {
					return err
				}
				arrJpmDetail = append(arrJpmDetail, v)
			}
			if findCoaGroup.Name == "HPP/COGS" {
				// if *JpmDetail.BalanceSheetDr != float64(0) {
				// 	return response.CustomErrorBuilder(http.StatusBadRequest, "Tidak Dapat Masukan Balance Sheet Debit", "Tidak Dapat Masukan Balance Sheet Debit")
				// }
				// if *JpmDetail.BalanceSheetCr != float64(0) {
				// 	return response.CustomErrorBuilder(http.StatusBadRequest, "Tidak Dapat Masukan Balance Sheet Credit", "Tidak Dapat Masukan Balance Sheet Credit")
				// }
				// if *JpmDetail.IncomeStatementCr != float64(0) {
				// 	return response.CustomErrorBuilder(http.StatusBadRequest, "Tidak Dapat Masukan Income Statement Credit", "Tidak Dapat Masukan Income Statement Credit")
				// }
				// if *JpmDetail.IncomeStatementDr == float64(0) {
				// 	return response.CustomErrorBuilder(http.StatusBadRequest, "Income Statement Debit Harus Di isi", "Income Statement Debit Harus Di isi")
				// }

				tbID, err := s.Repository.FindByTbd(ctx, &data.ConsolidationID, &JpmDetail.CoaCode)
				if err != nil {
					return response.CustomErrorBuilder(http.StatusBadRequest, "Tidak Dapat Menemukan Coa di Kombinasi Konsolidasi "+JpmDetail.CoaCode, "Tidak Dapat Menemukan Coa "+JpmDetail.CoaCode)
				}
				tbID.Context = ctx

				AmountAjeCr := *v.IncomeStatementCr + *tbID.ConsolidationDetailEntity.AmountJpmCr
				AmountJpmDr := *v.IncomeStatementDr + *tbID.ConsolidationDetailEntity.AmountJpmDr
				tbID.ConsolidationDetailEntity.AmountJpmCr = &AmountAjeCr
				tbID.ConsolidationDetailEntity.AmountJpmDr = &AmountJpmDr
				AmountAfterJpm := *tbID.AmountBeforeJpm + *tbID.AmountJpmDr - *tbID.AmountJpmCr
				tbID.AmountAfterJpm = &AmountAfterJpm
				AmountAfterJcte := *tbID.AmountAfterJpm + *tbID.AmountJcteDr - *tbID.AmountJcteCr
				tbID.AmountAfterJcte = &AmountAfterJcte

				findConsoleBridge, err := s.Repository.FindByConsoleBridge(ctx, &result.ConsolidationID)
				if err != nil {
					return response.CustomErrorBuilder(http.StatusBadRequest, "Tidak Dapat Menemukan Anak Usaha ", "Tidak Dapat Menemukan Anak Usaha ")
				}

				var dataAmountAnakUsahaa []float64
				for _, cb := range *findConsoleBridge {
					findConsoleBridgeDetail, err := s.Repository.FindByConsoleBridgeDetail(ctx, &cb.ID, &v.CoaCode)
					_, err = s.Repository.UpdateTbd(ctx, &tbID.ID, tbID)
					if err != nil {
						return response.CustomErrorBuilder(http.StatusBadRequest, "Tidak Dapat Menemukan Amount Anak Usaha ", "Tidak Dapat Menemukan Amount Anak Usaha ")
					}
					dataAmountAnakUsahaa = append(dataAmountAnakUsahaa, findConsoleBridgeDetail.Amount)
				}
				var sumAmountAnakUsaha float64 = 0.0

				for i := 0; i < len(dataAmountAnakUsahaa); i++ {
					sumAmountAnakUsaha += dataAmountAnakUsahaa[i]
				}

				var combineUnaudited float64 = *tbID.AmountAfterJcte + sumAmountAnakUsaha
				tbID.AmountCombineSubsidiary = &combineUnaudited

				var amountConsole float64 = *tbID.AmountCombineSubsidiary + *tbID.AmountJelimDr - *tbID.AmountJelimCr

				tbID.AmountConsole = &amountConsole

				_, err = s.Repository.UpdateTbd(ctx, &tbID.ID, tbID)
				if err != nil {
					return err
				}
				arrJpmDetail = append(arrJpmDetail, v)
			}
			if findCoaGroup.Name == "Selling Expense" {
				// if *JpmDetail.BalanceSheetDr != float64(0) {
				// 	return response.CustomErrorBuilder(http.StatusBadRequest, "Tidak Dapat Masukan Balance Sheet Debit", "Tidak Dapat Masukan Balance Sheet Debit")
				// }
				// if *JpmDetail.BalanceSheetCr != float64(0) {
				// 	return response.CustomErrorBuilder(http.StatusBadRequest, "Tidak Dapat Masukan Balance Sheet Credit", "Tidak Dapat Masukan Balance Sheet Credit")
				// }
				// if *JpmDetail.IncomeStatementCr != float64(0) {
				// 	return response.CustomErrorBuilder(http.StatusBadRequest, "Tidak Dapat Masukan Income Statement Credit", "Tidak Dapat Masukan Income Statement Credit")
				// }
				// if *JpmDetail.IncomeStatementDr == float64(0) {
				// 	return response.CustomErrorBuilder(http.StatusBadRequest, "Income Statement Debit Harus Di isi", "Income Statement Debit Harus Di isi")
				// }

				tbID, err := s.Repository.FindByTbd(ctx, &data.ConsolidationID, &JpmDetail.CoaCode)
				if err != nil {
					return response.CustomErrorBuilder(http.StatusBadRequest, "Tidak Dapat Menemukan Coa di Kombinasi Konsolidasi "+JpmDetail.CoaCode, "Tidak Dapat Menemukan Coa "+JpmDetail.CoaCode)
				}
				tbID.Context = ctx

				AmountAjeCr := *v.IncomeStatementCr + *tbID.ConsolidationDetailEntity.AmountJpmCr
				AmountJpmDr := *v.IncomeStatementDr + *tbID.ConsolidationDetailEntity.AmountJpmDr
				tbID.ConsolidationDetailEntity.AmountJpmCr = &AmountAjeCr
				tbID.ConsolidationDetailEntity.AmountJpmDr = &AmountJpmDr
				AmountAfterJpm := *tbID.AmountBeforeJpm + *tbID.AmountJpmDr - *tbID.AmountJpmCr
				tbID.AmountAfterJpm = &AmountAfterJpm
				AmountAfterJcte := *tbID.AmountAfterJpm + *tbID.AmountJcteDr - *tbID.AmountJcteCr
				tbID.AmountAfterJcte = &AmountAfterJcte

				findConsoleBridge, err := s.Repository.FindByConsoleBridge(ctx, &result.ConsolidationID)
				if err != nil {
					return response.CustomErrorBuilder(http.StatusBadRequest, "Tidak Dapat Menemukan Anak Usaha ", "Tidak Dapat Menemukan Anak Usaha ")
				}

				var dataAmountAnakUsahaa []float64
				for _, cb := range *findConsoleBridge {
					findConsoleBridgeDetail, err := s.Repository.FindByConsoleBridgeDetail(ctx, &cb.ID, &v.CoaCode)
					_, err = s.Repository.UpdateTbd(ctx, &tbID.ID, tbID)
					if err != nil {
						return response.CustomErrorBuilder(http.StatusBadRequest, "Tidak Dapat Menemukan Amount Anak Usaha ", "Tidak Dapat Menemukan Amount Anak Usaha ")
					}
					dataAmountAnakUsahaa = append(dataAmountAnakUsahaa, findConsoleBridgeDetail.Amount)
				}
				var sumAmountAnakUsaha float64 = 0.0

				for i := 0; i < len(dataAmountAnakUsahaa); i++ {
					sumAmountAnakUsaha += dataAmountAnakUsahaa[i]
				}

				var combineUnaudited float64 = *tbID.AmountAfterJcte + sumAmountAnakUsaha
				tbID.AmountCombineSubsidiary = &combineUnaudited

				var amountConsole float64 = *tbID.AmountCombineSubsidiary + *tbID.AmountJelimDr - *tbID.AmountJelimCr

				tbID.AmountConsole = &amountConsole

				_, err = s.Repository.UpdateTbd(ctx, &tbID.ID, tbID)
				if err != nil {
					return err
				}
				arrJpmDetail = append(arrJpmDetail, v)
			}
			if findCoaGroup.Name == "General & Admin Expense" {
				// if *JpmDetail.BalanceSheetDr != float64(0) {
				// 	return response.CustomErrorBuilder(http.StatusBadRequest, "Tidak Dapat Masukan Balance Sheet Debit", "Tidak Dapat Masukan Balance Sheet Debit")
				// }
				// if *JpmDetail.BalanceSheetCr != float64(0) {
				// 	return response.CustomErrorBuilder(http.StatusBadRequest, "Tidak Dapat Masukan Balance Sheet Credit", "Tidak Dapat Masukan Balance Sheet Credit")
				// }
				// if *JpmDetail.IncomeStatementCr != float64(0) {
				// 	return response.CustomErrorBuilder(http.StatusBadRequest, "Tidak Dapat Masukan Income Statement Credit", "Tidak Dapat Masukan Income Statement Credit")
				// }
				// if *JpmDetail.IncomeStatementDr == float64(0) {
				// 	return response.CustomErrorBuilder(http.StatusBadRequest, "Income Statement Debit Harus Di isi", "Income Statement Debit Harus Di isi")
				// }

				tbID, err := s.Repository.FindByTbd(ctx, &data.ConsolidationID, &JpmDetail.CoaCode)
				if err != nil {
					return response.CustomErrorBuilder(http.StatusBadRequest, "Tidak Dapat Menemukan Coa di Kombinasi Konsolidasi "+JpmDetail.CoaCode, "Tidak Dapat Menemukan Coa "+JpmDetail.CoaCode)
				}
				tbID.Context = ctx

				AmountAjeCr := *v.IncomeStatementCr + *tbID.ConsolidationDetailEntity.AmountJpmCr
				AmountJpmDr := *v.IncomeStatementDr + *tbID.ConsolidationDetailEntity.AmountJpmDr
				tbID.ConsolidationDetailEntity.AmountJpmCr = &AmountAjeCr
				tbID.ConsolidationDetailEntity.AmountJpmDr = &AmountJpmDr
				AmountAfterJpm := *tbID.AmountBeforeJpm + *tbID.AmountJpmDr - *tbID.AmountJpmCr
				tbID.AmountAfterJpm = &AmountAfterJpm
				AmountAfterJcte := *tbID.AmountAfterJpm + *tbID.AmountJcteDr - *tbID.AmountJcteCr
				tbID.AmountAfterJcte = &AmountAfterJcte

				findConsoleBridge, err := s.Repository.FindByConsoleBridge(ctx, &result.ConsolidationID)
				if err != nil {
					return response.CustomErrorBuilder(http.StatusBadRequest, "Tidak Dapat Menemukan Anak Usaha ", "Tidak Dapat Menemukan Anak Usaha ")
				}

				var dataAmountAnakUsahaa []float64
				for _, cb := range *findConsoleBridge {
					findConsoleBridgeDetail, err := s.Repository.FindByConsoleBridgeDetail(ctx, &cb.ID, &v.CoaCode)
					_, err = s.Repository.UpdateTbd(ctx, &tbID.ID, tbID)
					if err != nil {
						return response.CustomErrorBuilder(http.StatusBadRequest, "Tidak Dapat Menemukan Amount Anak Usaha ", "Tidak Dapat Menemukan Amount Anak Usaha ")
					}
					dataAmountAnakUsahaa = append(dataAmountAnakUsahaa, findConsoleBridgeDetail.Amount)
				}
				var sumAmountAnakUsaha float64 = 0.0

				for i := 0; i < len(dataAmountAnakUsahaa); i++ {
					sumAmountAnakUsaha += dataAmountAnakUsahaa[i]
				}

				var combineUnaudited float64 = *tbID.AmountAfterJcte + sumAmountAnakUsaha
				tbID.AmountCombineSubsidiary = &combineUnaudited

				var amountConsole float64 = *tbID.AmountCombineSubsidiary + *tbID.AmountJelimDr - *tbID.AmountJelimCr

				tbID.AmountConsole = &amountConsole

				_, err = s.Repository.UpdateTbd(ctx, &tbID.ID, tbID)
				if err != nil {
					return err
				}
				arrJpmDetail = append(arrJpmDetail, v)
			}
			if findCoaGroup.Name == "Other Income" {
				// if *JpmDetail.BalanceSheetDr != float64(0) {
				// 	return response.CustomErrorBuilder(http.StatusBadRequest, "Tidak Dapat Masukan Balance Sheet Debit", "Tidak Dapat Masukan Balance Sheet Debit")
				// }
				// if *JpmDetail.BalanceSheetCr != float64(0) {
				// 	return response.CustomErrorBuilder(http.StatusBadRequest, "Tidak Dapat Masukan Balance Sheet Credit", "Tidak Dapat Masukan Balance Sheet Credit")
				// }
				// if *JpmDetail.IncomeStatementCr == float64(0) {
				// 	return response.CustomErrorBuilder(http.StatusBadRequest, "Income Statement Credit Harus Di isi", "Income Statement Credit Harus Di isi")
				// }
				// if *JpmDetail.IncomeStatementDr != float64(0) {
				// 	return response.CustomErrorBuilder(http.StatusBadRequest, "Tidak Dapat Masukan Income Statement Debit", "Tidak Dapat Masukan Income Statement Debit")
				// }

				tbID, err := s.Repository.FindByTbd(ctx, &data.ConsolidationID, &JpmDetail.CoaCode)
				if err != nil {
					return response.CustomErrorBuilder(http.StatusBadRequest, "Tidak Dapat Menemukan Coa di Kombinasi Konsolidasi "+JpmDetail.CoaCode, "Tidak Dapat Menemukan Coa "+JpmDetail.CoaCode)
				}
				tbID.Context = ctx

				AmountJpmCr := *v.IncomeStatementCr + *tbID.ConsolidationDetailEntity.AmountJpmCr
				AmountAjeDr := *v.IncomeStatementDr + *tbID.ConsolidationDetailEntity.AmountJpmDr
				tbID.ConsolidationDetailEntity.AmountJpmCr = &AmountJpmCr
				tbID.ConsolidationDetailEntity.AmountJcteDr = &AmountAjeDr
				AmountAfterJpm := *tbID.AmountBeforeJpm - *tbID.AmountJpmDr + *tbID.AmountJpmCr
				tbID.AmountAfterJpm = &AmountAfterJpm
				AmountAfterJcte := *tbID.AmountAfterJpm - *tbID.AmountJcteDr + *tbID.AmountJcteCr
				tbID.AmountAfterJcte = &AmountAfterJcte

				findConsoleBridge, err := s.Repository.FindByConsoleBridge(ctx, &result.ConsolidationID)
				if err != nil {
					return response.CustomErrorBuilder(http.StatusBadRequest, "Tidak Dapat Menemukan Anak Usaha ", "Tidak Dapat Menemukan Anak Usaha ")
				}

				var dataAmountAnakUsahaa []float64
				for _, cb := range *findConsoleBridge {
					findConsoleBridgeDetail, err := s.Repository.FindByConsoleBridgeDetail(ctx, &cb.ID, &v.CoaCode)
					_, err = s.Repository.UpdateTbd(ctx, &tbID.ID, tbID)
					if err != nil {
						return response.CustomErrorBuilder(http.StatusBadRequest, "Tidak Dapat Menemukan Amount Anak Usaha ", "Tidak Dapat Menemukan Amount Anak Usaha ")
					}
					dataAmountAnakUsahaa = append(dataAmountAnakUsahaa, findConsoleBridgeDetail.Amount)
				}
				var sumAmountAnakUsaha float64 = 0.0

				for i := 0; i < len(dataAmountAnakUsahaa); i++ {
					sumAmountAnakUsaha += dataAmountAnakUsahaa[i]
				}

				var totalAmount float64 = *tbID.AmountAfterJcte + sumAmountAnakUsaha
				tbID.AmountCombineSubsidiary = &totalAmount

				var amountConsole float64 = *tbID.AmountCombineSubsidiary - *tbID.AmountJelimDr + *tbID.AmountJelimCr

				tbID.AmountConsole = &amountConsole

				_, err = s.Repository.UpdateTbd(ctx, &tbID.ID, tbID)
				if err != nil {
					return err
				}
				arrJpmDetail = append(arrJpmDetail, v)
			}
			if findCoaGroup.Name == "Other Expense" {
				// if *JpmDetail.BalanceSheetDr != float64(0) {
				// 	return response.CustomErrorBuilder(http.StatusBadRequest, "Tidak Dapat Masukan Balance Sheet Debit", "Tidak Dapat Masukan Balance Sheet Debit")
				// }
				// if *JpmDetail.BalanceSheetCr != float64(0) {
				// 	return response.CustomErrorBuilder(http.StatusBadRequest, "Tidak Dapat Masukan Balance Sheet Credit", "Tidak Dapat Masukan Balance Sheet Credit")
				// }
				// if *JpmDetail.IncomeStatementCr != float64(0) {
				// 	return response.CustomErrorBuilder(http.StatusBadRequest, "Tidak Dapat Masukan Income Statement Credit", "Tidak Dapat Masukan Income Statement Credit")
				// }
				// if *JpmDetail.IncomeStatementDr == float64(0) {
				// 	return response.CustomErrorBuilder(http.StatusBadRequest, "Income Statement Debit Harus Di isi", "Income Statement Debit Harus Di isi")
				// }

				tbID, err := s.Repository.FindByTbd(ctx, &data.ConsolidationID, &JpmDetail.CoaCode)
				if err != nil {
					return response.CustomErrorBuilder(http.StatusBadRequest, "Tidak Dapat Menemukan Coa di Kombinasi Konsolidasi "+JpmDetail.CoaCode, "Tidak Dapat Menemukan Coa "+JpmDetail.CoaCode)
				}
				tbID.Context = ctx

				AmountAjeCr := *v.IncomeStatementCr + *tbID.ConsolidationDetailEntity.AmountJpmCr
				AmountJpmDr := *v.IncomeStatementDr + *tbID.ConsolidationDetailEntity.AmountJpmDr
				tbID.ConsolidationDetailEntity.AmountJpmCr = &AmountAjeCr
				tbID.ConsolidationDetailEntity.AmountJpmDr = &AmountJpmDr
				AmountAfterJpm := *tbID.AmountBeforeJpm + *tbID.AmountJpmDr - *tbID.AmountJpmCr
				tbID.AmountAfterJpm = &AmountAfterJpm
				AmountAfterJcte := *tbID.AmountAfterJpm + *tbID.AmountJcteDr - *tbID.AmountJcteCr
				tbID.AmountAfterJcte = &AmountAfterJcte

				findConsoleBridge, err := s.Repository.FindByConsoleBridge(ctx, &result.ConsolidationID)
				if err != nil {
					return response.CustomErrorBuilder(http.StatusBadRequest, "Tidak Dapat Menemukan Anak Usaha ", "Tidak Dapat Menemukan Anak Usaha ")
				}

				var dataAmountAnakUsahaa []float64
				for _, cb := range *findConsoleBridge {
					findConsoleBridgeDetail, err := s.Repository.FindByConsoleBridgeDetail(ctx, &cb.ID, &v.CoaCode)
					_, err = s.Repository.UpdateTbd(ctx, &tbID.ID, tbID)
					if err != nil {
						return response.CustomErrorBuilder(http.StatusBadRequest, "Tidak Dapat Menemukan Amount Anak Usaha ", "Tidak Dapat Menemukan Amount Anak Usaha ")
					}
					dataAmountAnakUsahaa = append(dataAmountAnakUsahaa, findConsoleBridgeDetail.Amount)
				}
				var sumAmountAnakUsaha float64 = 0.0

				for i := 0; i < len(dataAmountAnakUsahaa); i++ {
					sumAmountAnakUsaha += dataAmountAnakUsahaa[i]
				}

				var combineUnaudited float64 = *tbID.AmountAfterJcte + sumAmountAnakUsaha
				tbID.AmountCombineSubsidiary = &combineUnaudited

				var amountConsole float64 = *tbID.AmountCombineSubsidiary + *tbID.AmountJelimDr - *tbID.AmountJelimCr

				tbID.AmountConsole = &amountConsole

				_, err = s.Repository.UpdateTbd(ctx, &tbID.ID, tbID)
				if err != nil {
					return err
				}
				arrJpmDetail = append(arrJpmDetail, v)
			}
			if findCoaGroup.Name == "Tax Expense" {
				// if *JpmDetail.BalanceSheetDr != float64(0) {
				// 	return response.CustomErrorBuilder(http.StatusBadRequest, "Tidak Dapat Masukan Balance Sheet Debit", "Tidak Dapat Masukan Balance Sheet Debit")
				// }
				// if *JpmDetail.BalanceSheetCr != float64(0) {
				// 	return response.CustomErrorBuilder(http.StatusBadRequest, "Tidak Dapat Masukan Balance Sheet Credit", "Tidak Dapat Masukan Balance Sheet Credit")
				// }
				// if *JpmDetail.IncomeStatementCr != float64(0) {
				// 	return response.CustomErrorBuilder(http.StatusBadRequest, "Tidak Dapat Masukan Income Statement Credit", "Tidak Dapat Masukan Income Statement Credit")
				// }
				// if *JpmDetail.IncomeStatementDr == float64(0) {
				// 	return response.CustomErrorBuilder(http.StatusBadRequest, "Income Statement Debit Harus Di isi", "Income Statement Debit Harus Di isi")
				// }

				tbID, err := s.Repository.FindByTbd(ctx, &data.ConsolidationID, &JpmDetail.CoaCode)
				if err != nil {
					return response.CustomErrorBuilder(http.StatusBadRequest, "Tidak Dapat Menemukan Coa di Kombinasi Konsolidasi "+JpmDetail.CoaCode, "Tidak Dapat Menemukan Coa "+JpmDetail.CoaCode)
				}
				tbID.Context = ctx

				AmountAjeCr := *v.IncomeStatementDr + *tbID.ConsolidationDetailEntity.AmountJpmCr
				AmountJpmDr := *v.IncomeStatementDr + *tbID.ConsolidationDetailEntity.AmountJpmDr
				tbID.ConsolidationDetailEntity.AmountJpmCr = &AmountAjeCr
				tbID.ConsolidationDetailEntity.AmountJpmDr = &AmountJpmDr
				AmountAfterJpm := *tbID.AmountBeforeJpm + *tbID.AmountJpmDr - *tbID.AmountJpmCr
				tbID.AmountAfterJpm = &AmountAfterJpm
				AmountAfterJcte := *tbID.AmountAfterJpm + *tbID.AmountJcteDr - *tbID.AmountJcteCr
				tbID.AmountAfterJcte = &AmountAfterJcte

				findConsoleBridge, err := s.Repository.FindByConsoleBridge(ctx, &result.ConsolidationID)
				if err != nil {
					return response.CustomErrorBuilder(http.StatusBadRequest, "Tidak Dapat Menemukan Anak Usaha ", "Tidak Dapat Menemukan Anak Usaha ")
				}

				var dataAmountAnakUsahaa []float64
				for _, cb := range *findConsoleBridge {
					findConsoleBridgeDetail, err := s.Repository.FindByConsoleBridgeDetail(ctx, &cb.ID, &v.CoaCode)
					_, err = s.Repository.UpdateTbd(ctx, &tbID.ID, tbID)
					if err != nil {
						return response.CustomErrorBuilder(http.StatusBadRequest, "Tidak Dapat Menemukan Amount Anak Usaha ", "Tidak Dapat Menemukan Amount Anak Usaha ")
					}
					dataAmountAnakUsahaa = append(dataAmountAnakUsahaa, findConsoleBridgeDetail.Amount)
				}
				var sumAmountAnakUsaha float64 = 0.0

				for i := 0; i < len(dataAmountAnakUsahaa); i++ {
					sumAmountAnakUsaha += dataAmountAnakUsahaa[i]
				}

				var combineUnaudited float64 = *tbID.AmountAfterJcte + sumAmountAnakUsaha
				tbID.AmountCombineSubsidiary = &combineUnaudited

				var amountConsole float64 = *tbID.AmountCombineSubsidiary + *tbID.AmountJelimDr - *tbID.AmountJelimCr

				tbID.AmountConsole = &amountConsole

				_, err = s.Repository.UpdateTbd(ctx, &tbID.ID, tbID)
				if err != nil {
					return err
				}
				arrJpmDetail = append(arrJpmDetail, v)
			}
			if findCoaGroup.Name == "Income (Loss) from subsidiary" {
				// if *JpmDetail.BalanceSheetDr != float64(0) {
				// 	return response.CustomErrorBuilder(http.StatusBadRequest, "Tidak Dapat Masukan Balance Sheet Debit", "Tidak Dapat Masukan Balance Sheet Debit")
				// }
				// if *JpmDetail.BalanceSheetCr != float64(0) {
				// 	return response.CustomErrorBuilder(http.StatusBadRequest, "Tidak Dapat Masukan Balance Sheet Credit", "Tidak Dapat Masukan Balance Sheet Credit")
				// }
				// if *JpmDetail.IncomeStatementCr == float64(0) {
				// 	return response.CustomErrorBuilder(http.StatusBadRequest, "Income Statement Credit Harus Di isi", "Income Statement Credit Harus Di isi")
				// }
				// if *JpmDetail.IncomeStatementDr != float64(0) {
				// 	return response.CustomErrorBuilder(http.StatusBadRequest, "Tidak Dapat Masukan Income Statement Debit", "Tidak Dapat Masukan Income Statement Debit")
				// }

				tbID, err := s.Repository.FindByTbd(ctx, &data.ConsolidationID, &JpmDetail.CoaCode)
				if err != nil {
					return response.CustomErrorBuilder(http.StatusBadRequest, "Tidak Dapat Menemukan Coa di Kombinasi Konsolidasi "+JpmDetail.CoaCode, "Tidak Dapat Menemukan Coa "+JpmDetail.CoaCode)
				}
				tbID.Context = ctx

				AmountJpmCr := *v.IncomeStatementCr + *tbID.ConsolidationDetailEntity.AmountJpmCr
				AmountAjeDr := *v.IncomeStatementDr + *tbID.ConsolidationDetailEntity.AmountJpmDr
				tbID.ConsolidationDetailEntity.AmountJpmCr = &AmountJpmCr
				tbID.ConsolidationDetailEntity.AmountJcteDr = &AmountAjeDr
				AmountAfterJpm := *tbID.AmountBeforeJpm - *tbID.AmountJpmDr + *tbID.AmountJpmCr
				tbID.AmountAfterJpm = &AmountAfterJpm
				AmountAfterJcte := *tbID.AmountAfterJpm - *tbID.AmountJcteDr + *tbID.AmountJcteCr
				tbID.AmountAfterJcte = &AmountAfterJcte

				findConsoleBridge, err := s.Repository.FindByConsoleBridge(ctx, &result.ConsolidationID)
				if err != nil {
					return response.CustomErrorBuilder(http.StatusBadRequest, "Tidak Dapat Menemukan Anak Usaha ", "Tidak Dapat Menemukan Anak Usaha ")
				}

				var dataAmountAnakUsahaa []float64
				for _, cb := range *findConsoleBridge {
					findConsoleBridgeDetail, err := s.Repository.FindByConsoleBridgeDetail(ctx, &cb.ID, &v.CoaCode)
					_, err = s.Repository.UpdateTbd(ctx, &tbID.ID, tbID)
					if err != nil {
						return response.CustomErrorBuilder(http.StatusBadRequest, "Tidak Dapat Menemukan Amount Anak Usaha ", "Tidak Dapat Menemukan Amount Anak Usaha ")
					}
					dataAmountAnakUsahaa = append(dataAmountAnakUsahaa, findConsoleBridgeDetail.Amount)
				}
				var sumAmountAnakUsaha float64 = 0.0

				for i := 0; i < len(dataAmountAnakUsahaa); i++ {
					sumAmountAnakUsaha += dataAmountAnakUsahaa[i]
				}

				var totalAmount float64 = *tbID.AmountAfterJcte + sumAmountAnakUsaha
				tbID.AmountCombineSubsidiary = &totalAmount

				var amountConsole float64 = *tbID.AmountCombineSubsidiary - *tbID.AmountJelimDr + *tbID.AmountJelimCr

				tbID.AmountConsole = &amountConsole

				_, err = s.Repository.UpdateTbd(ctx, &tbID.ID, tbID)
				if err != nil {
					return err
				}
				arrJpmDetail = append(arrJpmDetail, v)
			}
			if findCoaGroup.Name == "MINORITY INTEREST IN NET INCOME (NCI)" {
				// if *JpmDetail.BalanceSheetDr != float64(0) {
				// 	return response.CustomErrorBuilder(http.StatusBadRequest, "Tidak Dapat Masukan Balance Sheet Debit", "Tidak Dapat Masukan Balance Sheet Debit")
				// }
				// if *JpmDetail.BalanceSheetCr != float64(0) {
				// 	return response.CustomErrorBuilder(http.StatusBadRequest, "Tidak Dapat Masukan Balance Sheet Credit", "Tidak Dapat Masukan Balance Sheet Credit")
				// }
				// if *JpmDetail.IncomeStatementCr == float64(0) {
				// 	return response.CustomErrorBuilder(http.StatusBadRequest, "Income Statement Credit Harus Di isi", "Income Statement Credit Harus Di isi")
				// }
				// if *JpmDetail.IncomeStatementDr != float64(0) {
				// 	return response.CustomErrorBuilder(http.StatusBadRequest, "Tidak Dapat Masukan Income Statement Debit", "Tidak Dapat Masukan Income Statement Debit")
				// }

				tbID, err := s.Repository.FindByTbd(ctx, &data.ConsolidationID, &JpmDetail.CoaCode)
				if err != nil {
					return response.CustomErrorBuilder(http.StatusBadRequest, "Tidak Dapat Menemukan Coa di Kombinasi Konsolidasi "+JpmDetail.CoaCode, "Tidak Dapat Menemukan Coa "+JpmDetail.CoaCode)
				}
				tbID.Context = ctx

				AmountJpmCr := *v.IncomeStatementCr + *tbID.ConsolidationDetailEntity.AmountJpmCr
				AmountAjeDr := *v.IncomeStatementDr + *tbID.ConsolidationDetailEntity.AmountJpmDr
				tbID.ConsolidationDetailEntity.AmountJpmCr = &AmountJpmCr
				tbID.ConsolidationDetailEntity.AmountJpmDr = &AmountAjeDr
				AmountAfterJpm := *tbID.AmountBeforeJpm - *tbID.AmountJpmDr + *tbID.AmountJpmCr
				tbID.AmountAfterJpm = &AmountAfterJpm
				AmountAfterJcte := *tbID.AmountAfterJpm - *tbID.AmountJcteDr + *tbID.AmountJcteCr
				tbID.AmountAfterJcte = &AmountAfterJcte

				findConsoleBridge, err := s.Repository.FindByConsoleBridge(ctx, &result.ConsolidationID)
				if err != nil {
					return response.CustomErrorBuilder(http.StatusBadRequest, "Tidak Dapat Menemukan Anak Usaha ", "Tidak Dapat Menemukan Anak Usaha ")
				}

				var dataAmountAnakUsahaa []float64
				for _, cb := range *findConsoleBridge {
					findConsoleBridgeDetail, err := s.Repository.FindByConsoleBridgeDetail(ctx, &cb.ID, &v.CoaCode)
					_, err = s.Repository.UpdateTbd(ctx, &tbID.ID, tbID)
					if err != nil {
						return response.CustomErrorBuilder(http.StatusBadRequest, "Tidak Dapat Menemukan Amount Anak Usaha ", "Tidak Dapat Menemukan Amount Anak Usaha ")
					}
					dataAmountAnakUsahaa = append(dataAmountAnakUsahaa, findConsoleBridgeDetail.Amount)
				}
				var sumAmountAnakUsaha float64 = 0.0

				for i := 0; i < len(dataAmountAnakUsahaa); i++ {
					sumAmountAnakUsaha += dataAmountAnakUsahaa[i]
				}

				var totalAmount float64 = *tbID.AmountAfterJcte + sumAmountAnakUsaha
				tbID.AmountCombineSubsidiary = &totalAmount

				var amountConsole float64 = *tbID.AmountCombineSubsidiary - *tbID.AmountJelimDr + *tbID.AmountJelimCr

				tbID.AmountConsole = &amountConsole

				_, err = s.Repository.UpdateTbd(ctx, &tbID.ID, tbID)
				if err != nil {
					return err
				}
				arrJpmDetail = append(arrJpmDetail, v)
			}
			if findCoaGroup.Name == "Other Comprehensive Income" {

				// if *JpmDetail.BalanceSheetDr != float64(0) {
				// 	return response.CustomErrorBuilder(http.StatusBadRequest, "Tidak Dapat Masukan Balance Sheet Debit", "Tidak Dapat Masukan Balance Sheet Debit")
				// }
				// if *JpmDetail.BalanceSheetCr != float64(0) {
				// 	return response.CustomErrorBuilder(http.StatusBadRequest, "Tidak Dapat Masukan Balance Sheet Credit", "Tidak Dapat Masukan Balance Sheet Credit")
				// }
				// if *JpmDetail.IncomeStatementCr == float64(0) {
				// 	return response.CustomErrorBuilder(http.StatusBadRequest, "Income Statement Credit Harus Di isi", "Income Statement Credit Harus Di isi")
				// }
				// if *JpmDetail.IncomeStatementDr != float64(0) {
				// 	return response.CustomErrorBuilder(http.StatusBadRequest, "Tidak Dapat Masukan Income Statement Debit", "Tidak Dapat Masukan Income Statement Debit")
				// }

				tbID, err := s.Repository.FindByTbd(ctx, &data.ConsolidationID, &JpmDetail.CoaCode)
				if err != nil {
					return response.CustomErrorBuilder(http.StatusBadRequest, "Tidak Dapat Menemukan Coa di Kombinasi Konsolidasi "+JpmDetail.CoaCode, "Tidak Dapat Menemukan Coa "+JpmDetail.CoaCode)
				}
				tbID.Context = ctx

				AmountJpmCr := *v.IncomeStatementCr + *tbID.ConsolidationDetailEntity.AmountJpmCr
				AmountAjeDr := *v.IncomeStatementDr + *tbID.ConsolidationDetailEntity.AmountJpmDr
				tbID.ConsolidationDetailEntity.AmountJpmCr = &AmountJpmCr
				tbID.ConsolidationDetailEntity.AmountJpmDr = &AmountAjeDr
				AmountAfterJpm := *tbID.AmountBeforeJpm - *tbID.AmountJpmDr + *tbID.AmountJpmCr
				tbID.AmountAfterJpm = &AmountAfterJpm
				AmountAfterJcte := *tbID.AmountAfterJpm - *tbID.AmountJcteDr + *tbID.AmountJcteCr
				tbID.AmountAfterJcte = &AmountAfterJcte

				findConsoleBridge, err := s.Repository.FindByConsoleBridge(ctx, &result.ConsolidationID)
				if err != nil {
					return response.CustomErrorBuilder(http.StatusBadRequest, "Tidak Dapat Menemukan Anak Usaha ", "Tidak Dapat Menemukan Anak Usaha ")
				}

				var dataAmountAnakUsahaa []float64
				for _, cb := range *findConsoleBridge {
					findConsoleBridgeDetail, err := s.Repository.FindByConsoleBridgeDetail(ctx, &cb.ID, &v.CoaCode)
					_, err = s.Repository.UpdateTbd(ctx, &tbID.ID, tbID)
					if err != nil {
						return response.CustomErrorBuilder(http.StatusBadRequest, "Tidak Dapat Menemukan Amount Anak Usaha ", "Tidak Dapat Menemukan Amount Anak Usaha ")
					}
					dataAmountAnakUsahaa = append(dataAmountAnakUsahaa, findConsoleBridgeDetail.Amount)
				}
				var sumAmountAnakUsaha float64 = 0.0

				for i := 0; i < len(dataAmountAnakUsahaa); i++ {
					sumAmountAnakUsaha += dataAmountAnakUsahaa[i]
				}

				var totalAmount float64 = *tbID.AmountAfterJcte + sumAmountAnakUsaha
				tbID.AmountCombineSubsidiary = &totalAmount

				var amountConsole float64 = *tbID.AmountCombineSubsidiary - *tbID.AmountJelimDr + *tbID.AmountJelimCr

				tbID.AmountConsole = &amountConsole

				_, err = s.Repository.UpdateTbd(ctx, &tbID.ID, tbID)
				if err != nil {
					return err
				}
				arrJpmDetail = append(arrJpmDetail, v)
			}
			if findCoaGroup.Name == "Dampak penyesuaian proforma  atas OCI Entitas anak" {
				// if *JpmDetail.BalanceSheetDr != float64(0) {
				// 	return response.CustomErrorBuilder(http.StatusBadRequest, "Tidak Dapat Masukan Balance Sheet Debit", "Tidak Dapat Masukan Balance Sheet Debit")
				// }
				// if *JpmDetail.BalanceSheetCr != float64(0) {
				// 	return response.CustomErrorBuilder(http.StatusBadRequest, "Tidak Dapat Masukan Balance Sheet Credit", "Tidak Dapat Masukan Balance Sheet Credit")
				// }
				// if *JpmDetail.IncomeStatementCr == float64(0) {
				// 	return response.CustomErrorBuilder(http.StatusBadRequest, "Income Statement Credit Harus Di isi", "Income Statement Credit Harus Di isi")
				// }
				// if *JpmDetail.IncomeStatementDr != float64(0) {
				// 	return response.CustomErrorBuilder(http.StatusBadRequest, "Tidak Dapat Masukan Income Statement Debit", "Tidak Dapat Masukan Income Statement Debit")
				// }

				tbID, err := s.Repository.FindByTbd(ctx, &data.ConsolidationID, &JpmDetail.CoaCode)
				if err != nil {
					return response.CustomErrorBuilder(http.StatusBadRequest, "Tidak Dapat Menemukan Coa di Kombinasi Konsolidasi "+JpmDetail.CoaCode, "Tidak Dapat Menemukan Coa "+JpmDetail.CoaCode)
				}
				tbID.Context = ctx
				AmountJpmCr := *v.IncomeStatementCr + *tbID.ConsolidationDetailEntity.AmountJpmCr
				AmountAjeDr := *v.IncomeStatementDr + *tbID.ConsolidationDetailEntity.AmountJpmDr
				tbID.ConsolidationDetailEntity.AmountJpmCr = &AmountJpmCr
				tbID.ConsolidationDetailEntity.AmountJpmDr = &AmountAjeDr
				AmountAfterJpm := *tbID.AmountBeforeJpm - *tbID.AmountJpmDr + *tbID.AmountJpmCr
				tbID.AmountAfterJpm = &AmountAfterJpm
				AmountAfterJcte := *tbID.AmountAfterJpm - *tbID.AmountJcteDr + *tbID.AmountJcteCr
				tbID.AmountAfterJcte = &AmountAfterJcte

				findConsoleBridge, err := s.Repository.FindByConsoleBridge(ctx, &result.ConsolidationID)
				if err != nil {
					return response.CustomErrorBuilder(http.StatusBadRequest, "Tidak Dapat Menemukan Anak Usaha ", "Tidak Dapat Menemukan Anak Usaha ")
				}

				var dataAmountAnakUsahaa []float64
				for _, cb := range *findConsoleBridge {
					findConsoleBridgeDetail, err := s.Repository.FindByConsoleBridgeDetail(ctx, &cb.ID, &v.CoaCode)
					_, err = s.Repository.UpdateTbd(ctx, &tbID.ID, tbID)
					if err != nil {
						return response.CustomErrorBuilder(http.StatusBadRequest, "Tidak Dapat Menemukan Amount Anak Usaha ", "Tidak Dapat Menemukan Amount Anak Usaha ")
					}
					dataAmountAnakUsahaa = append(dataAmountAnakUsahaa, findConsoleBridgeDetail.Amount)
				}
				var sumAmountAnakUsaha float64 = 0.0

				for i := 0; i < len(dataAmountAnakUsahaa); i++ {
					sumAmountAnakUsaha += dataAmountAnakUsahaa[i]
				}

				var totalAmount float64 = *tbID.AmountAfterJcte + sumAmountAnakUsaha
				tbID.AmountCombineSubsidiary = &totalAmount

				var amountConsole float64 = *tbID.AmountCombineSubsidiary - *tbID.AmountJelimDr + *tbID.AmountJelimCr

				tbID.AmountConsole = &amountConsole

				_, err = s.Repository.UpdateTbd(ctx, &tbID.ID, tbID)
				if err != nil {
					return err
				}
				arrJpmDetail = append(arrJpmDetail, v)
			}
			if findCoaGroup.Name == "Non Controlling OCI" {
				// if *JpmDetail.BalanceSheetDr != float64(0) {
				// 	return response.CustomErrorBuilder(http.StatusBadRequest, "Tidak Dapat Masukan Balance Sheet Debit", "Tidak Dapat Masukan Balance Sheet Debit")
				// }
				// if *JpmDetail.BalanceSheetCr != float64(0) {
				// 	return response.CustomErrorBuilder(http.StatusBadRequest, "Tidak Dapat Masukan Balance Sheet Credit", "Tidak Dapat Masukan Balance Sheet Credit")
				// }
				// if *JpmDetail.IncomeStatementCr == float64(0) {
				// 	return response.CustomErrorBuilder(http.StatusBadRequest, "Income Statement Credit Harus Di isi", "Income Statement Credit Harus Di isi")
				// }
				// if *JpmDetail.IncomeStatementDr != float64(0) {
				// 	return response.CustomErrorBuilder(http.StatusBadRequest, "Tidak Dapat Masukan Income Statement Debit", "Tidak Dapat Masukan Income Statement Debit")
				// }

				tbID, err := s.Repository.FindByTbd(ctx, &data.ConsolidationID, &JpmDetail.CoaCode)
				if err != nil {
					return response.CustomErrorBuilder(http.StatusBadRequest, "Tidak Dapat Menemukan Coa di Kombinasi Konsolidasi "+JpmDetail.CoaCode, "Tidak Dapat Menemukan Coa "+JpmDetail.CoaCode)
				}
				tbID.Context = ctx

				AmountJpmCr := *v.IncomeStatementCr + *tbID.ConsolidationDetailEntity.AmountJpmCr
				AmountAjeDr := *v.IncomeStatementDr + *tbID.ConsolidationDetailEntity.AmountJpmDr
				tbID.ConsolidationDetailEntity.AmountJpmCr = &AmountJpmCr
				tbID.ConsolidationDetailEntity.AmountJpmDr = &AmountAjeDr
				AmountAfterJpm := *tbID.AmountBeforeJpm - *tbID.AmountJpmDr + *tbID.AmountJpmCr
				tbID.AmountAfterJpm = &AmountAfterJpm
				AmountAfterJcte := *tbID.AmountAfterJpm - *tbID.AmountJcteDr + *tbID.AmountJcteCr
				tbID.AmountAfterJcte = &AmountAfterJcte

				findConsoleBridge, err := s.Repository.FindByConsoleBridge(ctx, &result.ConsolidationID)
				if err != nil {
					return response.CustomErrorBuilder(http.StatusBadRequest, "Tidak Dapat Menemukan Anak Usaha ", "Tidak Dapat Menemukan Anak Usaha ")
				}

				var dataAmountAnakUsahaa []float64
				for _, cb := range *findConsoleBridge {
					findConsoleBridgeDetail, err := s.Repository.FindByConsoleBridgeDetail(ctx, &cb.ID, &v.CoaCode)
					_, err = s.Repository.UpdateTbd(ctx, &tbID.ID, tbID)
					if err != nil {
						return response.CustomErrorBuilder(http.StatusBadRequest, "Tidak Dapat Menemukan Amount Anak Usaha ", "Tidak Dapat Menemukan Amount Anak Usaha ")
					}
					dataAmountAnakUsahaa = append(dataAmountAnakUsahaa, findConsoleBridgeDetail.Amount)
				}
				var sumAmountAnakUsaha float64 = 0.0

				for i := 0; i < len(dataAmountAnakUsahaa); i++ {
					sumAmountAnakUsaha += dataAmountAnakUsahaa[i]
				}

				var totalAmount float64 = *tbID.AmountAfterJcte + sumAmountAnakUsaha
				tbID.AmountCombineSubsidiary = &totalAmount

				var amountConsole float64 = *tbID.AmountCombineSubsidiary - *tbID.AmountJelimDr + *tbID.AmountJelimCr

				tbID.AmountConsole = &amountConsole

				_, err = s.Repository.UpdateTbd(ctx, &tbID.ID, tbID)
				if err != nil {
					return err
				}
				arrJpmDetail = append(arrJpmDetail, v)
			}
		}

		// criteriaFormatterDetailSumData.TrxRefID = &ConsolidationData.TrxRefID

		formatterDetailSumData, err := s.Repository.FindSummary(ctx)
		if err != nil {
			return response.CustomErrorBuilder(http.StatusBadRequest, "Gagal Dapatkan sumary Coa", "Gagal Dapatkan sumary Coa")
		}
		for _, v := range *formatterDetailSumData {
			criteriaTBDetail := model.ConsolidationDetailFilterModel{}
			criteriaTBDetail.ConsolidationID = &result.ConsolidationID
			if v.AutoSummary != nil && *v.AutoSummary {
				code := fmt.Sprintf("%s_Subtotal", v.Code)
				criteriaTBDetail.Code = &code
				criteriaTBDetail.ConsolidationID = &result.ConsolidationID
				mfadetailsum, _, err := s.Repository.FindC(ctx, &criteriaTBDetail, &abstraction.Pagination{})
				if err != nil {
					return response.CustomErrorBuilder(http.StatusBadRequest, "Gagal Find Detail Consolidation", "Gagal Find Detail Consolidation")
				}

				for _, a := range *mfadetailsum {
					sumTBD, err := s.Repository.FindSummarys(ctx, &v.Code, &result.ConsolidationID, v.IsCoa)
					if err != nil {
						return response.CustomErrorBuilder(http.StatusBadRequest, "Gagal Menjumlahkan Summary Coa", "Gagal Menjumlahkan Summary Coa")
					}

					updateSummary := model.ConsolidationDetailEntityModel{
						ConsolidationDetailEntity: model.ConsolidationDetailEntity{
							AmountJpmCr:             sumTBD.AmountJpmCr,
							AmountJpmDr:             sumTBD.AmountJpmDr,
							AmountAfterJpm:          sumTBD.AmountAfterJpm,
							AmountJcteCr:            sumTBD.AmountJcteCr,
							AmountJcteDr:            sumTBD.AmountJcteDr,
							AmountAfterJcte:         sumTBD.AmountAfterJcte,
							AmountCombineSubsidiary: sumTBD.AmountCombineSubsidiary,
							AmountJelimCr:           sumTBD.AmountJelimCr,
							AmountJelimDr:           sumTBD.AmountJelimDr,
							AmountConsole:           sumTBD.AmountConsole,
						},
					}

					_, err = s.Repository.Updates(ctx, &a.ID, &updateSummary)
					if err != nil {
						return response.CustomErrorBuilder(http.StatusBadRequest, "Gagal Update Consolidation Detail", "Gagal Update Consolidation Detail")
					}
				}
			}
			if v.IsTotal != nil && *v.IsTotal == true && v.FxSummary != "" {

				tmpString := []string{"AmountBeforeJpm", "AmountJpmCr", "AmountJpmDr", "AmountAfterJpm", "AmountJcteCr", "AmountJcteDr", "AmountAfterJcte", "AmountCombineSubsidiary", "AmountJelimCr", "AmountJelimDr", "AmountConsole"}
				tmpTotalFl := make(map[string]*float64)
				// reg := regexp.MustCompile(`[0-9]+\d{3,}`)
				reg := regexp.MustCompile(`[A-Za-z0-9_~:()]+|[0-9]+\d{3,}`)

				for _, tipe := range tmpString {
					formula := strings.ToUpper(v.FxSummary)
					match := reg.FindAllString(formula, -1)
					amountBeforeJpm := make(map[string]interface{}, 0)
					for _, vMatch := range match {
						//cari jml berdasarkan code

						if len(vMatch) < 3 {
							continue
						}
						sumTBD, err := s.Repository.FindSummarys(ctx, &vMatch, &result.ConsolidationID, v.IsCoa)
						if err != nil {
							return response.CustomErrorBuilder(http.StatusBadRequest, "Gagal Mendapatkan Jumlah", "Gagal Mendapatkan Jumlah")
						}
						angka := 0.0

						if tipe == "AmountBeforeJpm" && sumTBD.AmountBeforeJpm != nil {
							angka = *sumTBD.AmountBeforeJpm
						} else if tipe == "AmountJpmCr" && sumTBD.AmountJpmCr != nil {
							angka = *sumTBD.AmountJpmCr
						} else if tipe == "AmountJpmDr" && sumTBD.AmountJpmDr != nil {
							angka = *sumTBD.AmountJpmDr
						} else if tipe == "AmountAfterJpm" && sumTBD.AmountAfterJpm != nil {
							angka = *sumTBD.AmountAfterJpm
						} else if tipe == "AmountJcteCr" && sumTBD.AmountJcteCr != nil {
							angka = *sumTBD.AmountJcteCr
						} else if tipe == "AmountJcteDr" && sumTBD.AmountJcteDr != nil {
							angka = *sumTBD.AmountJcteDr
						} else if tipe == "AmountAfterJcte" && sumTBD.AmountAfterJcte != nil {
							angka = *sumTBD.AmountAfterJcte
						} else if tipe == "AmountCombineSubsidiary" && sumTBD.AmountCombineSubsidiary != nil {
							angka = *sumTBD.AmountCombineSubsidiary
						} else if tipe == "AmountJelimCr" && sumTBD.AmountJelimCr != nil {
							angka = *sumTBD.AmountJelimCr
						} else if tipe == "AmountJelimDr" && sumTBD.AmountJelimDr != nil {
							angka = *sumTBD.AmountJelimDr
						} else if tipe == "AmountConsole" && sumTBD.AmountConsole != nil {
							angka = *sumTBD.AmountConsole
						}

						formula = helper.ReplaceWholeWord(formula, vMatch, fmt.Sprintf("(%2.f)", angka))
						// parameters[vMatch] = angka

					}

					expressionFormula, err := govaluate.NewEvaluableExpression(formula)
					if err != nil {
						fmt.Println(err)
						return err
					}
					result, err := expressionFormula.Evaluate(amountBeforeJpm)
					if err != nil {
						return response.CustomErrorBuilder(http.StatusBadRequest, "Tidak Dapat Masukan Balance Sheet Debit", "Tidak Dapat Masukan Balance Sheet Debit")
					}
					tmp := result.(float64)
					tmpTotalFl[tipe] = &tmp

				}
				criteriaTBDetail.ConsolidationID = &result.ConsolidationID
				criteriaTBDetail.Code = &v.Code
				mfadetailsum, err := s.Repository.FindDetailConsole(ctx, &criteriaTBDetail)
				if err != nil {
					fmt.Sprintln(err)
					return err
				}

				updateSummary := model.ConsolidationDetailEntityModel{
					ConsolidationDetailEntity: model.ConsolidationDetailEntity{
						AmountJpmCr:             tmpTotalFl["AmountJpmCr"],
						AmountJpmDr:             tmpTotalFl["AmountJpmDr"],
						AmountAfterJpm:          tmpTotalFl["AmountAfterJpm"],
						AmountJcteCr:            tmpTotalFl["AmountJcteCr"],
						AmountJcteDr:            tmpTotalFl["AmountJcteDr"],
						AmountAfterJcte:         tmpTotalFl["AmountAfterJcte"],
						AmountCombineSubsidiary: tmpTotalFl["AmountCombineSubsidiary"],
						AmountJelimCr:           tmpTotalFl["AmountJelimCr"],
						AmountJelimDr:           tmpTotalFl["AmountJelimDr"],
						AmountConsole:           tmpTotalFl["AmountConsole"],
					},
				}

				_, err = s.Repository.Updates(ctx, &mfadetailsum.ID, &updateSummary)
				if err != nil {
					return response.CustomErrorBuilder(http.StatusBadRequest, "Tidak Dapat Masukan Balance Sheet Debit", "Tidak Dapat Masukan Balance Sheet Debit")
				}
			}
			if v.Code == "LABA_BERSIH" {
				code := "310401004"
				criteriaTBDetail := model.ConsolidationDetailFilterModel{}
				criteriaTBDetail.ConsolidationID = &result.ConsolidationID
				criteriaTBDetail.Code = &code
				customRowOne, err := s.Repository.FindByTbd(ctx, &result.ConsolidationID, &code)
				if err != nil {
					return helper.ErrorHandler(err)
				}

				code = "310402002"
				criteriaTBDetail.Code = &code
				customRowTwo, err := s.Repository.FindByTbd(ctx, &result.ConsolidationID, &code)
				if err != nil {
					return helper.ErrorHandler(err)
				}

				code = "310501002"
				criteriaTBDetail.Code = &code
				customRowThree, err := s.Repository.FindByTbd(ctx, &result.ConsolidationID, &code)
				if err != nil {
					return helper.ErrorHandler(err)
				}

				code = "310502002"
				criteriaTBDetail.Code = &code
				customRowFour, err := s.Repository.FindByTbd(ctx, &result.ConsolidationID, &code)
				if err != nil {
					return helper.ErrorHandler(err)
				}

				code = "310503002"
				criteriaTBDetail.Code = &code
				customRowFive, err := s.Repository.FindByTbd(ctx, &result.ConsolidationID, &code)
				if err != nil {
					return helper.ErrorHandler(err)
				}

				code = "950101001" //REVALUATION FA (4337)
				criteriaTBDetail.Code = &code
				dataReFa, err := s.Repository.FindByTbd(ctx, &result.ConsolidationID, &code)
				if err != nil {
					return helper.ErrorHandler(err)
				}

				code = "950301001" //Financial Instrument (4342)
				criteriaTBDetail.Code = &code
				dataFinIn, err := s.Repository.FindByTbd(ctx, &result.ConsolidationID, &code)
				if err != nil {
					return helper.ErrorHandler(err)
				}

				code = "950301002" //Income tax relating to components of OCI (4343)
				criteriaTBDetail.Code = &code
				dataIncomeTax, err := s.Repository.FindByTbd(ctx, &result.ConsolidationID, &code)
				if err != nil {
					return helper.ErrorHandler(err)
				}

				code = "950401001" //Foreign Exchange (4345)
				criteriaTBDetail.Code = &code
				dataForex, err := s.Repository.FindByTbd(ctx, &result.ConsolidationID, &code)
				if err != nil {
					return helper.ErrorHandler(err)
				}

				code = "950401002" //Income tax relating to components of OCI (4346)
				criteriaTBDetail.Code = &code
				dataIncomeTax2, err := s.Repository.FindByTbd(ctx, &result.ConsolidationID, &code)
				if err != nil {
					return helper.ErrorHandler(err)
				}

				code = "LABA_BERSIH"
				criteriaTBDetail.Code = &code
				dataLaba, err := s.Repository.FindByTbd(ctx, &result.ConsolidationID, &code)
				if err != nil {
					return helper.ErrorHandler(err)
				}

				code = "TOTAL_PENGHASILAN_KOMPREHENSIF_LAIN~BS"
				criteriaTBDetail.Code = &code
				dataKomprehensif, err := s.Repository.FindByTbd(ctx, &result.ConsolidationID, &code)
				if err != nil {
					return helper.ErrorHandler(err)
				}

				summaryCodes, err := s.Repository.SummaryByCodes(ctx, &result.ConsolidationID, []string{"310501002", "310502002", "310503002"})
				if err != nil {
					return helper.ErrorHandler(err)
				}

				updatedCustomRowOne := model.ConsolidationDetailEntityModel{}
				updatedCustomRowOne.Context = ctx
				updatedCustomRowOne.AmountBeforeJpm = dataLaba.AmountBeforeJpm
				updatedCustomRowOne.AmountJpmCr = dataLaba.AmountJpmCr
				updatedCustomRowOne.AmountJpmDr = dataLaba.AmountJpmDr
				updatedCustomRowOne.AmountAfterJpm = dataLaba.AmountAfterJpm

				updatedCustomRowOne.AmountJcteCr = dataLaba.AmountJcteCr
				updatedCustomRowOne.AmountJcteDr = dataLaba.AmountJcteDr
				updatedCustomRowOne.AmountAfterJcte = dataLaba.AmountAfterJcte

				updatedCustomRowOne.AmountCombineSubsidiary = dataLaba.AmountCombineSubsidiary

				updatedCustomRowOne.AmountJelimCr = dataLaba.AmountJelimCr
				updatedCustomRowOne.AmountJelimDr = dataLaba.AmountJelimDr
				updatedCustomRowOne.AmountConsole = dataLaba.AmountConsole

				_, err = s.Repository.Updates(ctx, &customRowOne.ID, &updatedCustomRowOne)
				if err != nil {
					return helper.ErrorHandler(err)
				}

				//

				updatedCustomRowThree := model.ConsolidationDetailEntityModel{}
				updatedCustomRowThree.Context = ctx
				updatedCustomRowThree.AmountBeforeJpm = dataReFa.AmountBeforeJpm
				updatedCustomRowThree.AmountJpmCr = dataReFa.AmountJpmCr
				updatedCustomRowThree.AmountJpmDr = dataReFa.AmountJpmDr
				updatedCustomRowThree.AmountAfterJpm = dataReFa.AmountAfterJpm

				updatedCustomRowThree.AmountJcteCr = dataReFa.AmountJcteCr
				updatedCustomRowThree.AmountJcteDr = dataReFa.AmountJcteDr
				updatedCustomRowThree.AmountAfterJcte = dataReFa.AmountAfterJcte

				updatedCustomRowThree.AmountCombineSubsidiary = dataReFa.AmountCombineSubsidiary

				updatedCustomRowThree.AmountJelimCr = dataReFa.AmountJelimCr
				updatedCustomRowThree.AmountJelimDr = dataReFa.AmountJelimDr
				updatedCustomRowThree.AmountConsole = dataReFa.AmountConsole

				_, err = s.Repository.Updates(ctx, &customRowThree.ID, &updatedCustomRowThree)
				if err != nil {
					return helper.ErrorHandler(err)
				}

				updateCustomRowFour := model.ConsolidationDetailEntityModel{}
				updateCustomRowFour.Context = ctx

				tmp1 := 0.0
				if dataFinIn.AmountBeforeJpm != nil && *dataFinIn.AmountBeforeJpm != 0 {
					tmp1 = *dataFinIn.AmountBeforeJpm
				}
				if dataIncomeTax.AmountBeforeJpm != nil && *dataIncomeTax.AmountBeforeJpm != 0 {
					tmp1 = tmp1 + *dataIncomeTax.AmountBeforeJpm
				}
				updateCustomRowFour.AmountBeforeJpm = &tmp1

				tmp2 := 0.0
				if dataFinIn.AmountJpmCr != nil && *dataFinIn.AmountJpmCr != 0 {
					tmp2 = *dataFinIn.AmountJpmCr
				}
				if dataIncomeTax.AmountJpmCr != nil && *dataIncomeTax.AmountJpmCr != 0 {
					tmp2 = tmp2 + *dataIncomeTax.AmountJpmCr
				}
				updateCustomRowFour.AmountJpmCr = &tmp2

				tmp3 := 0.0
				if dataFinIn.AmountJpmDr != nil && *dataFinIn.AmountJpmDr != 0 {
					tmp3 = *dataFinIn.AmountJpmDr
				}
				if dataIncomeTax.AmountJpmDr != nil && *dataIncomeTax.AmountJpmDr != 0 {
					tmp3 = tmp3 + *dataIncomeTax.AmountJpmDr
				}
				updateCustomRowFour.AmountJpmDr = &tmp3

				tmp4 := 0.0
				if dataFinIn.AmountAfterJpm != nil && *dataFinIn.AmountAfterJpm != 0 {
					tmp4 = *dataFinIn.AmountAfterJpm
				}
				if dataIncomeTax.AmountAfterJpm != nil && *dataIncomeTax.AmountAfterJpm != 0 {
					tmp4 = tmp4 + *dataIncomeTax.AmountAfterJpm
				}
				updateCustomRowFour.AmountAfterJpm = &tmp4

				//jcte

				tmp5 := 0.0
				if dataFinIn.AmountJcteCr != nil && *dataFinIn.AmountJcteCr != 0 {
					tmp5 = *dataFinIn.AmountJcteCr
				}
				if dataIncomeTax.AmountJcteCr != nil && *dataIncomeTax.AmountJcteCr != 0 {
					tmp5 = tmp5 + *dataIncomeTax.AmountJcteCr
				}
				updateCustomRowFour.AmountJcteCr = &tmp5

				tmp6 := 0.0
				if dataFinIn.AmountJcteDr != nil && *dataFinIn.AmountJcteDr != 0 {
					tmp6 = *dataFinIn.AmountJcteDr
				}
				if dataIncomeTax.AmountJcteDr != nil && *dataIncomeTax.AmountJcteDr != 0 {
					tmp6 = tmp6 + *dataIncomeTax.AmountJcteDr
				}
				updateCustomRowFour.AmountJcteDr = &tmp6

				tmp7 := 0.0
				if dataFinIn.AmountAfterJcte != nil && *dataFinIn.AmountAfterJcte != 0 {
					tmp7 = *dataFinIn.AmountAfterJcte
				}
				if dataIncomeTax.AmountAfterJcte != nil && *dataIncomeTax.AmountAfterJcte != 0 {
					tmp7 = tmp7 + *dataIncomeTax.AmountAfterJcte
				}
				updateCustomRowFour.AmountAfterJcte = &tmp7

				// acs
				tmp8 := 0.0
				if dataFinIn.AmountCombineSubsidiary != nil && *dataFinIn.AmountCombineSubsidiary != 0 {
					tmp8 = *dataFinIn.AmountCombineSubsidiary
				}
				if dataIncomeTax.AmountCombineSubsidiary != nil && *dataIncomeTax.AmountCombineSubsidiary != 0 {
					tmp8 = tmp8 + *dataIncomeTax.AmountCombineSubsidiary
				}
				updateCustomRowFour.AmountCombineSubsidiary = &tmp8

				//jelim
				tmp9 := 0.0
				if dataFinIn.AmountJelimCr != nil && *dataFinIn.AmountJelimCr != 0 {
					tmp9 = *dataFinIn.AmountJelimCr
				}
				if dataIncomeTax.AmountJelimCr != nil && *dataIncomeTax.AmountJelimCr != 0 {
					tmp9 = tmp9 + *dataIncomeTax.AmountJelimCr
				}
				updateCustomRowFour.AmountJelimCr = &tmp9

				tmp10 := 0.0
				if dataFinIn.AmountJelimDr != nil && *dataFinIn.AmountJelimDr != 0 {
					tmp10 = *dataFinIn.AmountJelimDr
				}
				if dataIncomeTax.AmountJelimDr != nil && *dataIncomeTax.AmountJelimDr != 0 {
					tmp10 = tmp10 + *dataIncomeTax.AmountJelimDr
				}
				updateCustomRowFour.AmountJelimDr = &tmp10

				tmp11 := 0.0
				if dataFinIn.AmountConsole != nil && *dataFinIn.AmountConsole != 0 {
					tmp11 = *dataFinIn.AmountConsole
				}
				if dataIncomeTax.AmountConsole != nil && *dataIncomeTax.AmountConsole != 0 {
					tmp11 = tmp11 + *dataIncomeTax.AmountConsole
				}
				updateCustomRowFour.AmountConsole = &tmp11

				_, err = s.Repository.Updates(ctx, &customRowFour.ID, &updateCustomRowFour)
				if err != nil {
					return helper.ErrorHandler(err)
				}

				updateCustomRowFive := model.ConsolidationDetailEntityModel{}
				updateCustomRowFive.Context = ctx

				tmpa := 0.0
				if dataForex.AmountBeforeJpm != nil && *dataForex.AmountBeforeJpm != 0 {
					tmpa = *dataForex.AmountBeforeJpm
				}
				if dataIncomeTax2.AmountBeforeJpm != nil && *dataIncomeTax2.AmountBeforeJpm != 0 {
					tmpa = tmpa + *dataIncomeTax2.AmountBeforeJpm
				}
				updateCustomRowFive.AmountBeforeJpm = &tmpa

				tmpb := 0.0
				if dataForex.AmountJpmCr != nil && *dataForex.AmountJpmCr != 0 {
					tmpb = *dataForex.AmountJpmCr
				}
				if dataIncomeTax2.AmountJpmCr != nil && *dataIncomeTax2.AmountJpmCr != 0 {
					tmpb = tmpb + *dataIncomeTax2.AmountJpmCr
				}
				updateCustomRowFive.AmountJpmCr = &tmpb

				tmpc := 0.0
				if dataForex.AmountJpmDr != nil && *dataForex.AmountJpmDr != 0 {
					tmpc = *dataForex.AmountJpmDr
				}
				if dataIncomeTax2.AmountJpmDr != nil && *dataIncomeTax2.AmountJpmDr != 0 {
					tmpc = tmpc + *dataIncomeTax2.AmountJpmDr
				}
				updateCustomRowFive.AmountJpmDr = &tmpc

				tmpd := 0.0
				if dataForex.AmountAfterJpm != nil && *dataForex.AmountAfterJpm != 0 {
					tmpd = *dataForex.AmountAfterJpm
				}
				if dataIncomeTax2.AmountAfterJpm != nil && *dataIncomeTax2.AmountAfterJpm != 0 {
					tmpd = tmpd + *dataIncomeTax2.AmountAfterJpm
				}
				updateCustomRowFive.AmountAfterJpm = &tmpd

				//jcte
				tmpf := 0.0
				if dataForex.AmountJcteCr != nil && *dataForex.AmountJcteCr != 0 {
					tmpf = *dataForex.AmountJcteCr
				}
				if dataIncomeTax2.AmountJcteCr != nil && *dataIncomeTax2.AmountJcteCr != 0 {
					tmpf = tmpf + *dataIncomeTax2.AmountJcteCr
				}
				updateCustomRowFive.AmountJcteCr = &tmpf

				tmpe := 0.0
				if dataForex.AmountJcteDr != nil && *dataForex.AmountJcteDr != 0 {
					tmpe = *dataForex.AmountJcteDr
				}
				if dataIncomeTax2.AmountJcteDr != nil && *dataIncomeTax2.AmountJcteDr != 0 {
					tmpe = tmpe + *dataIncomeTax2.AmountJcteDr
				}
				updateCustomRowFive.AmountJcteDr = &tmpe

				tmph := 0.0
				if dataForex.AmountAfterJcte != nil && *dataForex.AmountAfterJcte != 0 {
					tmph = *dataForex.AmountAfterJcte
				}
				if dataIncomeTax2.AmountAfterJcte != nil && *dataIncomeTax2.AmountAfterJcte != 0 {
					tmph = tmph + *dataIncomeTax2.AmountAfterJcte
				}
				updateCustomRowFive.AmountAfterJcte = &tmph

				// acs
				tmpg := 0.0
				if dataForex.AmountCombineSubsidiary != nil && *dataForex.AmountCombineSubsidiary != 0 {
					tmpg = *dataForex.AmountCombineSubsidiary
				}
				if dataIncomeTax2.AmountCombineSubsidiary != nil && *dataIncomeTax2.AmountCombineSubsidiary != 0 {
					tmpg = tmpg + *dataIncomeTax2.AmountCombineSubsidiary
				}
				updateCustomRowFive.AmountCombineSubsidiary = &tmpg

				//jelim
				tmpi := 0.0
				if dataForex.AmountJelimCr != nil && *dataForex.AmountJelimCr != 0 {
					tmpi = *dataForex.AmountJelimCr
				}
				if dataIncomeTax2.AmountJelimCr != nil && *dataIncomeTax2.AmountJelimCr != 0 {
					tmpi = tmpi + *dataIncomeTax2.AmountJelimCr
				}
				updateCustomRowFive.AmountJelimCr = &tmpi

				tmpj := 0.0
				if dataForex.AmountJelimDr != nil && *dataForex.AmountJelimDr != 0 {
					tmpj = *dataForex.AmountJelimDr
				}
				if dataIncomeTax2.AmountJelimDr != nil && *dataIncomeTax2.AmountJelimDr != 0 {
					tmpj = tmpj + *dataIncomeTax2.AmountJelimDr
				}
				updateCustomRowFive.AmountJelimDr = &tmpj

				tmpk := 0.0
				if dataForex.AmountConsole != nil && *dataForex.AmountConsole != 0 {
					tmpk = *dataForex.AmountConsole
				}
				if dataIncomeTax2.AmountConsole != nil && *dataIncomeTax2.AmountConsole != 0 {
					tmpk = tmpk + *dataIncomeTax2.AmountConsole
				}
				updateCustomRowFive.AmountConsole = &tmpk
				_, err = s.Repository.Updates(ctx, &customRowFive.ID, &updateCustomRowFive)
				if err != nil {
					return helper.ErrorHandler(err)
				}
				updateCustomRowTwo := model.ConsolidationDetailEntityModel{}
				updateCustomRowTwo.Context = ctx

				tmp12 := 0.0
				if dataKomprehensif.AmountBeforeJpm != nil && *dataKomprehensif.AmountBeforeJpm != 0 {
					tmp12 = *dataKomprehensif.AmountBeforeJpm
				}
				if summaryCodes.AmountBeforeJpm != nil && *summaryCodes.AmountBeforeJpm != 0 {
					tmp12 = tmp12 - *summaryCodes.AmountBeforeJpm
				}
				updateCustomRowTwo.AmountBeforeJpm = &tmp12

				tmp13 := 0.0
				if dataKomprehensif.AmountJpmCr != nil && *dataKomprehensif.AmountJpmCr != 0 {
					tmp13 = *dataKomprehensif.AmountJpmCr
				}
				if summaryCodes.AmountJpmCr != nil && *summaryCodes.AmountJpmCr != 0 {
					tmp13 = tmp13 - *summaryCodes.AmountJpmCr
				}
				updateCustomRowTwo.AmountJpmCr = &tmp13

				tmp14 := 0.0
				if dataKomprehensif.AmountJpmDr != nil && *dataKomprehensif.AmountJpmDr != 0 {
					tmp14 = *dataKomprehensif.AmountJpmDr
				}
				if summaryCodes.AmountJpmDr != nil && *summaryCodes.AmountJpmDr != 0 {
					tmp14 = tmp14 - *summaryCodes.AmountJpmDr
				}
				updateCustomRowTwo.AmountJpmDr = &tmp14

				tmp15 := 0.0
				if dataKomprehensif.AmountAfterJpm != nil && *dataKomprehensif.AmountAfterJpm != 0 {
					tmp15 = *dataKomprehensif.AmountAfterJpm
				}
				if summaryCodes.AmountAfterJpm != nil && *summaryCodes.AmountAfterJpm != 0 {
					tmp15 = tmp15 - *summaryCodes.AmountAfterJpm
				}
				updateCustomRowTwo.AmountAfterJpm = &tmp15

				//jcte

				tmp16 := 0.0
				if dataKomprehensif.AmountJcteCr != nil && *dataKomprehensif.AmountJcteCr != 0 {
					tmp16 = *dataKomprehensif.AmountJcteCr
				}
				if summaryCodes.AmountJcteCr != nil && *summaryCodes.AmountJcteCr != 0 {
					tmp16 = tmp16 - *summaryCodes.AmountJcteCr
				}
				updateCustomRowTwo.AmountJcteCr = &tmp16

				tmp17 := 0.0
				if dataKomprehensif.AmountJcteDr != nil && *dataKomprehensif.AmountJcteDr != 0 {
					tmp17 = *dataKomprehensif.AmountJcteDr
				}
				if summaryCodes.AmountJcteDr != nil && *summaryCodes.AmountJcteDr != 0 {
					tmp17 = tmp17 - *summaryCodes.AmountJcteDr
				}
				updateCustomRowTwo.AmountJcteDr = &tmp17

				tmp18 := 0.0
				if dataKomprehensif.AmountAfterJcte != nil && *dataKomprehensif.AmountAfterJcte != 0 {
					tmp18 = *dataKomprehensif.AmountAfterJcte
				}
				if summaryCodes.AmountAfterJcte != nil && *summaryCodes.AmountAfterJcte != 0 {
					tmp18 = tmp18 - *summaryCodes.AmountAfterJcte
				}
				updateCustomRowTwo.AmountAfterJcte = &tmp18

				// acs
				tmp19 := 0.0
				if dataKomprehensif.AmountCombineSubsidiary != nil && *dataKomprehensif.AmountCombineSubsidiary != 0 {
					tmp19 = *dataKomprehensif.AmountCombineSubsidiary
				}
				if summaryCodes.AmountCombineSubsidiary != nil && *summaryCodes.AmountCombineSubsidiary != 0 {
					tmp19 = tmp19 - *summaryCodes.AmountCombineSubsidiary
				}
				updateCustomRowTwo.AmountCombineSubsidiary = &tmp19

				//jelim
				tmp20 := 0.0
				if dataKomprehensif.AmountJelimCr != nil && *dataKomprehensif.AmountJelimCr != 0 {
					tmp20 = *dataKomprehensif.AmountJelimCr
				}
				if summaryCodes.AmountJelimCr != nil && *summaryCodes.AmountJelimCr != 0 {
					tmp20 = tmp20 - *summaryCodes.AmountJelimCr
				}
				updateCustomRowTwo.AmountJelimCr = &tmp20

				tmp21 := 0.0
				if dataKomprehensif.AmountJelimDr != nil && *dataKomprehensif.AmountJelimDr != 0 {
					tmp21 = *dataKomprehensif.AmountJelimDr
				}
				if summaryCodes.AmountJelimDr != nil && *summaryCodes.AmountJelimDr != 0 {
					tmp21 = tmp21 - *summaryCodes.AmountJelimDr
				}
				updateCustomRowTwo.AmountJelimDr = &tmp21

				tmp22 := 0.0
				if dataKomprehensif.AmountConsole != nil && *dataKomprehensif.AmountConsole != 0 {
					tmp22 = *dataKomprehensif.AmountConsole
				}
				if summaryCodes.AmountConsole != nil && *summaryCodes.AmountConsole != 0 {
					tmp22 = tmp22 - *summaryCodes.AmountConsole
				}
				updateCustomRowTwo.AmountConsole = &tmp22
				_, err = s.Repository.Updates(ctx, &customRowTwo.ID, &updateCustomRowTwo)
				if err != nil {
					return helper.ErrorHandler(err)
				}
			}
		}
		{
			{
				code := "310401004"
				criteriaTBDetail := model.ConsolidationDetailFilterModel{}
				criteriaTBDetail.ConsolidationID = &result.ConsolidationID
				criteriaTBDetail.Code = &code
				customRowOne, err := s.Repository.FindByTbd(ctx, &result.ConsolidationID, &code)
				if err != nil {
					return helper.ErrorHandler(err)
				}

				code = "310402002"
				criteriaTBDetail.Code = &code
				customRowTwo, err := s.Repository.FindByTbd(ctx, &result.ConsolidationID, &code)
				if err != nil {
					return helper.ErrorHandler(err)
				}

				code = "310501002"
				criteriaTBDetail.Code = &code
				customRowThree, err := s.Repository.FindByTbd(ctx, &result.ConsolidationID, &code)
				if err != nil {
					return helper.ErrorHandler(err)
				}

				code = "310502002"
				criteriaTBDetail.Code = &code
				customRowFour, err := s.Repository.FindByTbd(ctx, &result.ConsolidationID, &code)
				if err != nil {
					return helper.ErrorHandler(err)
				}

				code = "310503002"
				criteriaTBDetail.Code = &code
				customRowFive, err := s.Repository.FindByTbd(ctx, &result.ConsolidationID, &code)
				if err != nil {
					return helper.ErrorHandler(err)
				}

				code = "950101001" //REVALUATION FA (4337)
				criteriaTBDetail.Code = &code
				dataReFa, err := s.Repository.FindByTbd(ctx, &result.ConsolidationID, &code)
				if err != nil {
					return helper.ErrorHandler(err)
				}

				code = "950301001" //Financial Instrument (4342)
				criteriaTBDetail.Code = &code
				dataFinIn, err := s.Repository.FindByTbd(ctx, &result.ConsolidationID, &code)
				if err != nil {
					return helper.ErrorHandler(err)
				}

				code = "950301002" //Income tax relating to components of OCI (4343)
				criteriaTBDetail.Code = &code
				dataIncomeTax, err := s.Repository.FindByTbd(ctx, &result.ConsolidationID, &code)
				if err != nil {
					return helper.ErrorHandler(err)
				}

				code = "950401001" //Foreign Exchange (4345)
				criteriaTBDetail.Code = &code
				dataForex, err := s.Repository.FindByTbd(ctx, &result.ConsolidationID, &code)
				if err != nil {
					return helper.ErrorHandler(err)
				}

				code = "950401002" //Income tax relating to components of OCI (4346)
				criteriaTBDetail.Code = &code
				dataIncomeTax2, err := s.Repository.FindByTbd(ctx, &result.ConsolidationID, &code)
				if err != nil {
					return helper.ErrorHandler(err)
				}

				code = "LABA_BERSIH"
				criteriaTBDetail.Code = &code
				dataLaba, err := s.Repository.FindByTbd(ctx, &result.ConsolidationID, &code)
				if err != nil {
					return helper.ErrorHandler(err)
				}

				code = "TOTAL_PENGHASILAN_KOMPREHENSIF_LAIN~BS"
				criteriaTBDetail.Code = &code
				dataKomprehensif, err := s.Repository.FindByTbd(ctx, &result.ConsolidationID, &code)
				if err != nil {
					return helper.ErrorHandler(err)
				}

				summaryCodes, err := s.Repository.SummaryByCodes(ctx, &result.ConsolidationID, []string{"310501002", "310502002", "310503002"})
				if err != nil {
					return helper.ErrorHandler(err)
				}

				updatedCustomRowOne := model.ConsolidationDetailEntityModel{}
				updatedCustomRowOne.Context = ctx
				updatedCustomRowOne.AmountBeforeJpm = dataLaba.AmountBeforeJpm
				updatedCustomRowOne.AmountJpmCr = dataLaba.AmountJpmCr
				updatedCustomRowOne.AmountJpmDr = dataLaba.AmountJpmDr
				updatedCustomRowOne.AmountAfterJpm = dataLaba.AmountAfterJpm

				updatedCustomRowOne.AmountJcteCr = dataLaba.AmountJcteCr
				updatedCustomRowOne.AmountJcteDr = dataLaba.AmountJcteDr
				updatedCustomRowOne.AmountAfterJcte = dataLaba.AmountAfterJcte

				updatedCustomRowOne.AmountCombineSubsidiary = dataLaba.AmountCombineSubsidiary

				updatedCustomRowOne.AmountJelimCr = dataLaba.AmountJelimCr
				updatedCustomRowOne.AmountJelimDr = dataLaba.AmountJelimDr
				updatedCustomRowOne.AmountConsole = dataLaba.AmountConsole

				_, err = s.Repository.Updates(ctx, &customRowOne.ID, &updatedCustomRowOne)
				if err != nil {
					return helper.ErrorHandler(err)
				}

				//

				updatedCustomRowThree := model.ConsolidationDetailEntityModel{}
				updatedCustomRowThree.Context = ctx
				updatedCustomRowThree.AmountBeforeJpm = dataReFa.AmountBeforeJpm
				updatedCustomRowThree.AmountJpmCr = dataReFa.AmountJpmCr
				updatedCustomRowThree.AmountJpmDr = dataReFa.AmountJpmDr
				updatedCustomRowThree.AmountAfterJpm = dataReFa.AmountAfterJpm

				updatedCustomRowThree.AmountJcteCr = dataReFa.AmountJcteCr
				updatedCustomRowThree.AmountJcteDr = dataReFa.AmountJcteDr
				updatedCustomRowThree.AmountAfterJcte = dataReFa.AmountAfterJcte

				updatedCustomRowThree.AmountCombineSubsidiary = dataReFa.AmountCombineSubsidiary

				updatedCustomRowThree.AmountJelimCr = dataReFa.AmountJelimCr
				updatedCustomRowThree.AmountJelimDr = dataReFa.AmountJelimDr
				updatedCustomRowThree.AmountConsole = dataReFa.AmountConsole

				_, err = s.Repository.Updates(ctx, &customRowThree.ID, &updatedCustomRowThree)
				if err != nil {
					return helper.ErrorHandler(err)
				}

				updateCustomRowFour := model.ConsolidationDetailEntityModel{}
				updateCustomRowFour.Context = ctx

				tmp1 := 0.0
				if dataFinIn.AmountBeforeJpm != nil && *dataFinIn.AmountBeforeJpm != 0 {
					tmp1 = *dataFinIn.AmountBeforeJpm
				}
				if dataIncomeTax.AmountBeforeJpm != nil && *dataIncomeTax.AmountBeforeJpm != 0 {
					tmp1 = tmp1 + *dataIncomeTax.AmountBeforeJpm
				}
				updateCustomRowFour.AmountBeforeJpm = &tmp1

				tmp2 := 0.0
				if dataFinIn.AmountJpmCr != nil && *dataFinIn.AmountJpmCr != 0 {
					tmp2 = *dataFinIn.AmountJpmCr
				}
				if dataIncomeTax.AmountJpmCr != nil && *dataIncomeTax.AmountJpmCr != 0 {
					tmp2 = tmp2 + *dataIncomeTax.AmountJpmCr
				}
				updateCustomRowFour.AmountJpmCr = &tmp2

				tmp3 := 0.0
				if dataFinIn.AmountJpmDr != nil && *dataFinIn.AmountJpmDr != 0 {
					tmp3 = *dataFinIn.AmountJpmDr
				}
				if dataIncomeTax.AmountJpmDr != nil && *dataIncomeTax.AmountJpmDr != 0 {
					tmp3 = tmp3 + *dataIncomeTax.AmountJpmDr
				}
				updateCustomRowFour.AmountJpmDr = &tmp3

				tmp4 := 0.0
				if dataFinIn.AmountAfterJpm != nil && *dataFinIn.AmountAfterJpm != 0 {
					tmp4 = *dataFinIn.AmountAfterJpm
				}
				if dataIncomeTax.AmountAfterJpm != nil && *dataIncomeTax.AmountAfterJpm != 0 {
					tmp4 = tmp4 + *dataIncomeTax.AmountAfterJpm
				}
				updateCustomRowFour.AmountAfterJpm = &tmp4

				//jcte

				tmp5 := 0.0
				if dataFinIn.AmountJcteCr != nil && *dataFinIn.AmountJcteCr != 0 {
					tmp5 = *dataFinIn.AmountJcteCr
				}
				if dataIncomeTax.AmountJcteCr != nil && *dataIncomeTax.AmountJcteCr != 0 {
					tmp5 = tmp5 + *dataIncomeTax.AmountJcteCr
				}
				updateCustomRowFour.AmountJcteCr = &tmp5

				tmp6 := 0.0
				if dataFinIn.AmountJcteDr != nil && *dataFinIn.AmountJcteDr != 0 {
					tmp6 = *dataFinIn.AmountJcteDr
				}
				if dataIncomeTax.AmountJcteDr != nil && *dataIncomeTax.AmountJcteDr != 0 {
					tmp6 = tmp6 + *dataIncomeTax.AmountJcteDr
				}
				updateCustomRowFour.AmountJcteDr = &tmp6

				tmp7 := 0.0
				if dataFinIn.AmountAfterJcte != nil && *dataFinIn.AmountAfterJcte != 0 {
					tmp7 = *dataFinIn.AmountAfterJcte
				}
				if dataIncomeTax.AmountAfterJcte != nil && *dataIncomeTax.AmountAfterJcte != 0 {
					tmp7 = tmp7 + *dataIncomeTax.AmountAfterJcte
				}
				updateCustomRowFour.AmountAfterJcte = &tmp7

				// acs
				tmp8 := 0.0
				if dataFinIn.AmountCombineSubsidiary != nil && *dataFinIn.AmountCombineSubsidiary != 0 {
					tmp8 = *dataFinIn.AmountCombineSubsidiary
				}
				if dataIncomeTax.AmountCombineSubsidiary != nil && *dataIncomeTax.AmountCombineSubsidiary != 0 {
					tmp8 = tmp8 + *dataIncomeTax.AmountCombineSubsidiary
				}
				updateCustomRowFour.AmountCombineSubsidiary = &tmp8

				//jelim
				tmp9 := 0.0
				if dataFinIn.AmountJelimCr != nil && *dataFinIn.AmountJelimCr != 0 {
					tmp9 = *dataFinIn.AmountJelimCr
				}
				if dataIncomeTax.AmountJelimCr != nil && *dataIncomeTax.AmountJelimCr != 0 {
					tmp9 = tmp9 + *dataIncomeTax.AmountJelimCr
				}
				updateCustomRowFour.AmountJelimCr = &tmp9

				tmp10 := 0.0
				if dataFinIn.AmountJelimDr != nil && *dataFinIn.AmountJelimDr != 0 {
					tmp10 = *dataFinIn.AmountJelimDr
				}
				if dataIncomeTax.AmountJelimDr != nil && *dataIncomeTax.AmountJelimDr != 0 {
					tmp10 = tmp10 + *dataIncomeTax.AmountJelimDr
				}
				updateCustomRowFour.AmountJelimDr = &tmp10

				tmp11 := 0.0
				if dataFinIn.AmountConsole != nil && *dataFinIn.AmountConsole != 0 {
					tmp11 = *dataFinIn.AmountConsole
				}
				if dataIncomeTax.AmountConsole != nil && *dataIncomeTax.AmountConsole != 0 {
					tmp11 = tmp11 + *dataIncomeTax.AmountConsole
				}
				updateCustomRowFour.AmountConsole = &tmp11

				_, err = s.Repository.Updates(ctx, &customRowFour.ID, &updateCustomRowFour)
				if err != nil {
					return helper.ErrorHandler(err)
				}

				updateCustomRowFive := model.ConsolidationDetailEntityModel{}
				updateCustomRowFive.Context = ctx

				tmpa := 0.0
				if dataForex.AmountBeforeJpm != nil && *dataForex.AmountBeforeJpm != 0 {
					tmpa = *dataForex.AmountBeforeJpm
				}
				if dataIncomeTax2.AmountBeforeJpm != nil && *dataIncomeTax2.AmountBeforeJpm != 0 {
					tmpa = tmpa + *dataIncomeTax2.AmountBeforeJpm
				}
				updateCustomRowFive.AmountBeforeJpm = &tmpa

				tmpb := 0.0
				if dataForex.AmountJpmCr != nil && *dataForex.AmountJpmCr != 0 {
					tmpb = *dataForex.AmountJpmCr
				}
				if dataIncomeTax2.AmountJpmCr != nil && *dataIncomeTax2.AmountJpmCr != 0 {
					tmpb = tmpb + *dataIncomeTax2.AmountJpmCr
				}
				updateCustomRowFive.AmountJpmCr = &tmpb

				tmpc := 0.0
				if dataForex.AmountJpmDr != nil && *dataForex.AmountJpmDr != 0 {
					tmpc = *dataForex.AmountJpmDr
				}
				if dataIncomeTax2.AmountJpmDr != nil && *dataIncomeTax2.AmountJpmDr != 0 {
					tmpc = tmpc + *dataIncomeTax2.AmountJpmDr
				}
				updateCustomRowFive.AmountJpmDr = &tmpc

				tmpd := 0.0
				if dataForex.AmountAfterJpm != nil && *dataForex.AmountAfterJpm != 0 {
					tmpd = *dataForex.AmountAfterJpm
				}
				if dataIncomeTax2.AmountAfterJpm != nil && *dataIncomeTax2.AmountAfterJpm != 0 {
					tmpd = tmpd + *dataIncomeTax2.AmountAfterJpm
				}
				updateCustomRowFive.AmountAfterJpm = &tmpd

				//jcte
				tmpf := 0.0
				if dataForex.AmountJcteCr != nil && *dataForex.AmountJcteCr != 0 {
					tmpf = *dataForex.AmountJcteCr
				}
				if dataIncomeTax2.AmountJcteCr != nil && *dataIncomeTax2.AmountJcteCr != 0 {
					tmpf = tmpf + *dataIncomeTax2.AmountJcteCr
				}
				updateCustomRowFive.AmountJcteCr = &tmpf

				tmpe := 0.0
				if dataForex.AmountJcteDr != nil && *dataForex.AmountJcteDr != 0 {
					tmpe = *dataForex.AmountJcteDr
				}
				if dataIncomeTax2.AmountJcteDr != nil && *dataIncomeTax2.AmountJcteDr != 0 {
					tmpe = tmpe + *dataIncomeTax2.AmountJcteDr
				}
				updateCustomRowFive.AmountJcteDr = &tmpe

				tmph := 0.0
				if dataForex.AmountAfterJcte != nil && *dataForex.AmountAfterJcte != 0 {
					tmph = *dataForex.AmountAfterJcte
				}
				if dataIncomeTax2.AmountAfterJcte != nil && *dataIncomeTax2.AmountAfterJcte != 0 {
					tmph = tmph + *dataIncomeTax2.AmountAfterJcte
				}
				updateCustomRowFive.AmountAfterJcte = &tmph

				// acs
				tmpg := 0.0
				if dataForex.AmountCombineSubsidiary != nil && *dataForex.AmountCombineSubsidiary != 0 {
					tmpg = *dataForex.AmountCombineSubsidiary
				}
				if dataIncomeTax2.AmountCombineSubsidiary != nil && *dataIncomeTax2.AmountCombineSubsidiary != 0 {
					tmpg = tmpg + *dataIncomeTax2.AmountCombineSubsidiary
				}
				updateCustomRowFive.AmountCombineSubsidiary = &tmpg

				//jelim
				tmpi := 0.0
				if dataForex.AmountJelimCr != nil && *dataForex.AmountJelimCr != 0 {
					tmpi = *dataForex.AmountJelimCr
				}
				if dataIncomeTax2.AmountJelimCr != nil && *dataIncomeTax2.AmountJelimCr != 0 {
					tmpi = tmpi + *dataIncomeTax2.AmountJelimCr
				}
				updateCustomRowFive.AmountJelimCr = &tmpi

				tmpj := 0.0
				if dataForex.AmountJelimDr != nil && *dataForex.AmountJelimDr != 0 {
					tmpj = *dataForex.AmountJelimDr
				}
				if dataIncomeTax2.AmountJelimDr != nil && *dataIncomeTax2.AmountJelimDr != 0 {
					tmpj = tmpj + *dataIncomeTax2.AmountJelimDr
				}
				updateCustomRowFive.AmountJelimDr = &tmpj

				tmpk := 0.0
				if dataForex.AmountConsole != nil && *dataForex.AmountConsole != 0 {
					tmpk = *dataForex.AmountConsole
				}
				if dataIncomeTax2.AmountConsole != nil && *dataIncomeTax2.AmountConsole != 0 {
					tmpk = tmpk + *dataIncomeTax2.AmountConsole
				}
				updateCustomRowFive.AmountConsole = &tmpk
				_, err = s.Repository.Updates(ctx, &customRowFive.ID, &updateCustomRowFive)
				if err != nil {
					return helper.ErrorHandler(err)
				}
				updateCustomRowTwo := model.ConsolidationDetailEntityModel{}
				updateCustomRowTwo.Context = ctx

				tmp12 := 0.0
				if dataKomprehensif.AmountBeforeJpm != nil && *dataKomprehensif.AmountBeforeJpm != 0 {
					tmp12 = *dataKomprehensif.AmountBeforeJpm
				}
				if summaryCodes.AmountBeforeJpm != nil && *summaryCodes.AmountBeforeJpm != 0 {
					tmp12 = tmp12 - *summaryCodes.AmountBeforeJpm
				}
				updateCustomRowTwo.AmountBeforeJpm = &tmp12

				tmp13 := 0.0
				if dataKomprehensif.AmountJpmCr != nil && *dataKomprehensif.AmountJpmCr != 0 {
					tmp13 = *dataKomprehensif.AmountJpmCr
				}
				if summaryCodes.AmountJpmCr != nil && *summaryCodes.AmountJpmCr != 0 {
					tmp13 = tmp13 - *summaryCodes.AmountJpmCr
				}
				updateCustomRowTwo.AmountJpmCr = &tmp13

				tmp14 := 0.0
				if dataKomprehensif.AmountJpmDr != nil && *dataKomprehensif.AmountJpmDr != 0 {
					tmp14 = *dataKomprehensif.AmountJpmDr
				}
				if summaryCodes.AmountJpmDr != nil && *summaryCodes.AmountJpmDr != 0 {
					tmp14 = tmp14 - *summaryCodes.AmountJpmDr
				}
				updateCustomRowTwo.AmountJpmDr = &tmp14

				tmp15 := 0.0
				if dataKomprehensif.AmountAfterJpm != nil && *dataKomprehensif.AmountAfterJpm != 0 {
					tmp15 = *dataKomprehensif.AmountAfterJpm
				}
				if summaryCodes.AmountAfterJpm != nil && *summaryCodes.AmountAfterJpm != 0 {
					tmp15 = tmp15 - *summaryCodes.AmountAfterJpm
				}
				updateCustomRowTwo.AmountAfterJpm = &tmp15

				//jcte

				tmp16 := 0.0
				if dataKomprehensif.AmountJcteCr != nil && *dataKomprehensif.AmountJcteCr != 0 {
					tmp16 = *dataKomprehensif.AmountJcteCr
				}
				if summaryCodes.AmountJcteCr != nil && *summaryCodes.AmountJcteCr != 0 {
					tmp16 = tmp16 - *summaryCodes.AmountJcteCr
				}
				updateCustomRowTwo.AmountJcteCr = &tmp16

				tmp17 := 0.0
				if dataKomprehensif.AmountJcteDr != nil && *dataKomprehensif.AmountJcteDr != 0 {
					tmp17 = *dataKomprehensif.AmountJcteDr
				}
				if summaryCodes.AmountJcteDr != nil && *summaryCodes.AmountJcteDr != 0 {
					tmp17 = tmp17 - *summaryCodes.AmountJcteDr
				}
				updateCustomRowTwo.AmountJcteDr = &tmp17

				tmp18 := 0.0
				if dataKomprehensif.AmountAfterJcte != nil && *dataKomprehensif.AmountAfterJcte != 0 {
					tmp18 = *dataKomprehensif.AmountAfterJcte
				}
				if summaryCodes.AmountAfterJcte != nil && *summaryCodes.AmountAfterJcte != 0 {
					tmp18 = tmp18 - *summaryCodes.AmountAfterJcte
				}
				updateCustomRowTwo.AmountAfterJcte = &tmp18

				// acs
				tmp19 := 0.0
				if dataKomprehensif.AmountCombineSubsidiary != nil && *dataKomprehensif.AmountCombineSubsidiary != 0 {
					tmp19 = *dataKomprehensif.AmountCombineSubsidiary
				}
				if summaryCodes.AmountCombineSubsidiary != nil && *summaryCodes.AmountCombineSubsidiary != 0 {
					tmp19 = tmp19 - *summaryCodes.AmountCombineSubsidiary
				}
				updateCustomRowTwo.AmountCombineSubsidiary = &tmp19

				//jelim
				tmp20 := 0.0
				if dataKomprehensif.AmountJelimCr != nil && *dataKomprehensif.AmountJelimCr != 0 {
					tmp20 = *dataKomprehensif.AmountJelimCr
				}
				if summaryCodes.AmountJelimCr != nil && *summaryCodes.AmountJelimCr != 0 {
					tmp20 = tmp20 - *summaryCodes.AmountJelimCr
				}
				updateCustomRowTwo.AmountJelimCr = &tmp20

				tmp21 := 0.0
				if dataKomprehensif.AmountJelimDr != nil && *dataKomprehensif.AmountJelimDr != 0 {
					tmp21 = *dataKomprehensif.AmountJelimDr
				}
				if summaryCodes.AmountJelimDr != nil && *summaryCodes.AmountJelimDr != 0 {
					tmp21 = tmp21 - *summaryCodes.AmountJelimDr
				}
				updateCustomRowTwo.AmountJelimDr = &tmp21

				tmp22 := 0.0
				if dataKomprehensif.AmountConsole != nil && *dataKomprehensif.AmountConsole != 0 {
					tmp22 = *dataKomprehensif.AmountConsole
				}
				if summaryCodes.AmountConsole != nil && *summaryCodes.AmountConsole != 0 {
					tmp22 = tmp22 - *summaryCodes.AmountConsole
				}
				updateCustomRowTwo.AmountConsole = &tmp22
				_, err = s.Repository.Updates(ctx, &customRowTwo.ID, &updateCustomRowTwo)
				if err != nil {
					return helper.ErrorHandler(err)
				}
			}
		}
		for _, v := range *formatterDetailSumData {
			criteriaTBDetail := model.ConsolidationDetailFilterModel{}
			criteriaTBDetail.ConsolidationID = &result.ConsolidationID

			if v.AutoSummary != nil && *v.AutoSummary {
				code := fmt.Sprintf("%s_Subtotal", v.Code)
				criteriaTBDetail.Code = &code
				criteriaTBDetail.ConsolidationID = &result.ConsolidationID
				mfadetailsum, _, err := s.Repository.FindC(ctx, &criteriaTBDetail, &abstraction.Pagination{})
				if err != nil {
					return response.CustomErrorBuilder(http.StatusBadRequest, "Gagal Find Detail Consolidation", "Gagal Find Detail Consolidation")
				}

				for _, a := range *mfadetailsum {
					sumTBD, err := s.Repository.FindSummarys(ctx, &v.Code, &result.ConsolidationID, v.IsCoa)
					if err != nil {
						return response.CustomErrorBuilder(http.StatusBadRequest, "Gagal Menjumlahkan Summary Coa", "Gagal Menjumlahkan Summary Coa")
					}

					updateSummary := model.ConsolidationDetailEntityModel{
						ConsolidationDetailEntity: model.ConsolidationDetailEntity{
							AmountJpmCr:             sumTBD.AmountJpmCr,
							AmountJpmDr:             sumTBD.AmountJpmDr,
							AmountAfterJpm:          sumTBD.AmountAfterJpm,
							AmountJcteCr:            sumTBD.AmountJcteCr,
							AmountJcteDr:            sumTBD.AmountJcteDr,
							AmountAfterJcte:         sumTBD.AmountAfterJcte,
							AmountCombineSubsidiary: sumTBD.AmountCombineSubsidiary,
							AmountJelimCr:           sumTBD.AmountJelimCr,
							AmountJelimDr:           sumTBD.AmountJelimDr,
							AmountConsole:           sumTBD.AmountConsole,
						},
					}

					_, err = s.Repository.Updates(ctx, &a.ID, &updateSummary)
					if err != nil {
						return response.CustomErrorBuilder(http.StatusBadRequest, "Gagal Update Consolidation Detail", "Gagal Update Consolidation Detail")
					}
				}
			}
			if v.IsTotal != nil && *v.IsTotal == true && v.FxSummary != "" {

				tmpString := []string{"AmountBeforeJpm", "AmountJpmCr", "AmountJpmDr", "AmountAfterJpm", "AmountJcteCr", "AmountJcteDr", "AmountAfterJcte", "AmountCombineSubsidiary", "AmountJelimCr", "AmountJelimDr", "AmountConsole"}
				tmpTotalFl := make(map[string]*float64)
				// reg := regexp.MustCompile(`[0-9]+\d{3,}`)
				reg := regexp.MustCompile(`[A-Za-z0-9_~:()]+|[0-9]+\d{3,}`)

				for _, tipe := range tmpString {
					formula := strings.ToUpper(v.FxSummary)
					match := reg.FindAllString(formula, -1)
					amountBeforeJpm := make(map[string]interface{}, 0)
					for _, vMatch := range match {
						//cari jml berdasarkan code

						if len(vMatch) < 3 {
							continue
						}
						sumTBD, err := s.Repository.FindSummarys(ctx, &vMatch, &result.ConsolidationID, v.IsCoa)
						if err != nil {
							return response.CustomErrorBuilder(http.StatusBadRequest, "Gagal Mendapatkan Jumlah", "Gagal Mendapatkan Jumlah")
						}
						angka := 0.0

						if tipe == "AmountBeforeJpm" && sumTBD.AmountBeforeJpm != nil {
							angka = *sumTBD.AmountBeforeJpm
						} else if tipe == "AmountJpmCr" && sumTBD.AmountJpmCr != nil {
							angka = *sumTBD.AmountJpmCr
						} else if tipe == "AmountJpmDr" && sumTBD.AmountJpmDr != nil {
							angka = *sumTBD.AmountJpmDr
						} else if tipe == "AmountAfterJpm" && sumTBD.AmountAfterJpm != nil {
							angka = *sumTBD.AmountAfterJpm
						} else if tipe == "AmountJcteCr" && sumTBD.AmountJcteCr != nil {
							angka = *sumTBD.AmountJcteCr
						} else if tipe == "AmountJcteDr" && sumTBD.AmountJcteDr != nil {
							angka = *sumTBD.AmountJcteDr
						} else if tipe == "AmountAfterJcte" && sumTBD.AmountAfterJcte != nil {
							angka = *sumTBD.AmountAfterJcte
						} else if tipe == "AmountCombineSubsidiary" && sumTBD.AmountCombineSubsidiary != nil {
							angka = *sumTBD.AmountCombineSubsidiary
						} else if tipe == "AmountJelimCr" && sumTBD.AmountJelimCr != nil {
							angka = *sumTBD.AmountJelimCr
						} else if tipe == "AmountJelimDr" && sumTBD.AmountJelimDr != nil {
							angka = *sumTBD.AmountJelimDr
						} else if tipe == "AmountConsole" && sumTBD.AmountConsole != nil {
							angka = *sumTBD.AmountConsole
						}

						formula = helper.ReplaceWholeWord(formula, vMatch, fmt.Sprintf("(%2.f)", angka))
						// parameters[vMatch] = angka

					}

					expressionFormula, err := govaluate.NewEvaluableExpression(formula)
					if err != nil {
						fmt.Println(err)
						return err
					}
					result, err := expressionFormula.Evaluate(amountBeforeJpm)
					if err != nil {
						return response.CustomErrorBuilder(http.StatusBadRequest, "Tidak Dapat Masukan Balance Sheet Debit", "Tidak Dapat Masukan Balance Sheet Debit")
					}
					tmp := result.(float64)
					tmpTotalFl[tipe] = &tmp

				}
				criteriaTBDetail.ConsolidationID = &result.ConsolidationID
				criteriaTBDetail.Code = &v.Code
				mfadetailsum, err := s.Repository.FindDetailConsole(ctx, &criteriaTBDetail)
				if err != nil {
					fmt.Sprintln(err)
					return err
				}

				updateSummary := model.ConsolidationDetailEntityModel{
					ConsolidationDetailEntity: model.ConsolidationDetailEntity{
						AmountJpmCr:             tmpTotalFl["AmountJpmCr"],
						AmountJpmDr:             tmpTotalFl["AmountJpmDr"],
						AmountAfterJpm:          tmpTotalFl["AmountAfterJpm"],
						AmountJcteCr:            tmpTotalFl["AmountJcteCr"],
						AmountJcteDr:            tmpTotalFl["AmountJcteDr"],
						AmountAfterJcte:         tmpTotalFl["AmountAfterJcte"],
						AmountCombineSubsidiary: tmpTotalFl["AmountCombineSubsidiary"],
						AmountJelimCr:           tmpTotalFl["AmountJelimCr"],
						AmountJelimDr:           tmpTotalFl["AmountJelimDr"],
						AmountConsole:           tmpTotalFl["AmountConsole"],
					},
				}

				_, err = s.Repository.Updates(ctx, &mfadetailsum.ID, &updateSummary)
				if err != nil {
					return response.CustomErrorBuilder(http.StatusBadRequest, "Tidak Dapat Masukan Balance Sheet Debit", "Tidak Dapat Masukan Balance Sheet Debit")
				}
			}
		}

		return nil

	}); err != nil {
		return &dto.JpmCreateResponse{}, err
	}
	result := &dto.JpmCreateResponse{
		JpmEntityModel: data,
	}
	return result, nil
}
func (s *service) Update(ctx *abstraction.Context, payload *dto.JpmUpdateRequest) (*dto.JpmUpdateResponse, error) {
	var data model.JpmEntityModel

	if err := trxmanager.New(s.Db).WithTrx(ctx, func(ctx *abstraction.Context) error {
		if _, err := s.Repository.FindByID(ctx, &payload.ID); err != nil {
			return response.ErrorBuilder(&response.ErrorConstant.BadRequest, err)
		}
		data.Context = ctx
		data.JpmEntity = payload.JpmEntity
		result, err := s.Repository.Update(ctx, &payload.ID, &data)
		if err != nil {
			return response.ErrorBuilder(&response.ErrorConstant.InternalServerError, err)
		}
		data = *result
		return nil
	}); err != nil {
		return &dto.JpmUpdateResponse{}, err
	}
	result := &dto.JpmUpdateResponse{
		JpmEntityModel: data,
	}
	return result, nil
}
func (s *service) Delete(ctx *abstraction.Context, payload *dto.JpmDeleteRequest) (*dto.JpmDeleteResponse, error) {
	var data model.JpmEntityModel
	if err := trxmanager.New(s.Db).WithTrx(ctx, func(ctx *abstraction.Context) error {
		if _, err := s.Repository.FindByID(ctx, &payload.ID); err != nil {
			return response.ErrorBuilder(&response.ErrorConstant.BadRequest, err)
		}
		FindAjeDetail, err := s.JpmDetailRepository.FindWithAjeID(ctx, &payload.ID)
		if err != nil {
			return err
		}
		JpmID, err := s.Repository.FindByID(ctx, &payload.ID)
		if err != nil {
			return err
		}
		var arrJpmDetails []model.JpmDetailEntityModel

		for _, v := range *FindAjeDetail {
			JpmDetail := model.JpmDetailEntityModel{
				Context: ctx,
			}
			JpmDetail.ID = v.ID
			JpmID, err := s.Repository.FindByID(ctx, &payload.ID)
			if err != nil {
				return err
			}

			// tbID, err := s.Repository.FindByTbd(ctx, &JpmID.ConsolidationID, &v.CoaCode )
			// if err != nil {
			// 	return response.CustomErrorBuilder(http.StatusBadRequest, "Tidak Dapat Menemukan Coa "  + v.CoaCode, "Tidak Dapat Menemukan Coa "  + v.CoaCode)
			// }

			// tbID.Context = ctx

			// AmountAjeCr := *v.BalanceSheetCr - *tbID.ConsolidationDetailEntity.AmountJpmCr
			// AmountAjeDr := *v.BalanceSheetDr - *tbID.ConsolidationDetailEntity.AmountJpmDr
			// tbID.ConsolidationDetailEntity.AmountJpmCr = &AmountAjeCr
			// tbID.ConsolidationDetailEntity.AmountJpmDr = &AmountAjeDr
			// AmountAfterJpm := +*tbID.AmountCombineSubsidiary + *tbID.AmountJpmDr - *tbID.AmountJpmDr
			// tbID.AmountAfterJpm = &AmountAfterJpm

			// _, err = s.Repository.UpdateTbd(ctx, &tbID.ID, tbID)
			// if err != nil {
			// 	return err
			// }
			// arrJpmDetails = append(arrJpmDetails, JpmDetail)
			if v.CoaCode == "310401004" || v.CoaCode == "310501002" || v.CoaCode == "310502002" || v.CoaCode == "310503002" || v.CoaCode == "310402002" {
				return response.CustomErrorBuilder(http.StatusBadRequest, response.E_BAD_REQUEST, "Coa Tersebut Tidak Dapat Melalukan Jurnal "+v.CoaCode)
			}
			findCoa, err := s.Repository.FindByCoa(ctx, &v.CoaCode)
			if err != nil {
				return response.CustomErrorBuilder(http.StatusBadRequest, "Tidak Dapat Menemukan Coa "+v.CoaCode, "Tidak Dapat Menemukan Coa "+v.CoaCode)
			}

			findCoaType, err := s.Repository.FindByCoaType(ctx, &findCoa.CoaTypeID)
			if err != nil {
				return response.CustomErrorBuilder(http.StatusBadRequest, "Tidak Dapat Menemukan Coa di Type Coa "+v.CoaCode, "Tidak Dapat Menemukan Coa di Type Coa "+v.CoaCode)
			}
			findCoaGroup, err := s.Repository.FindByCoaGroup(ctx, &findCoaType.CoaGroupID)
			if err != nil {
				return response.CustomErrorBuilder(http.StatusBadRequest, "Tidak Dapat Menemukan Coa di Group coa "+v.CoaCode, "Tidak Dapat Menemukan Coa di Group coa "+v.CoaCode)
			}

			if findCoaGroup.Name == "ASET" {
				// if *JpmDetail.BalanceSheetDr == float64(0) {
				// 	return response.CustomErrorBuilder(http.StatusBadRequest, "Balance Sheet Debit Harus Di isi", "Balance Sheet Debit Harus Di isi")
				// }
				// if *JpmDetail.BalanceSheetCr > float64(0) {
				// 	return response.CustomErrorBuilder(http.StatusBadRequest, "Tidak Dapat Masukan Balance Sheet Credit", "Tidak Dapat Masukan Balance Sheet Credit")
				// }
				// if *JpmDetail.IncomeStatementCr > float64(0) {
				// 	return response.CustomErrorBuilder(http.StatusBadRequest, "Tidak Dapat Masukan Income Statement Credit", "Tidak Dapat Masukan Income Statement Credit")
				// }
				// if *JpmDetail.IncomeStatementDr > float64(0) {
				// 	return response.CustomErrorBuilder(http.StatusBadRequest, "Tidak Dapat Masukan Income Statement Debit", "Tidak Dapat Masukan Income Statement Debit")
				// }

				tbID, err := s.Repository.FindByTbd(ctx, &JpmID.ConsolidationID, &v.CoaCode)
				if err != nil {
					return response.CustomErrorBuilder(http.StatusBadRequest, "Tidak Dapat Menemukan Coa di Kombinasi Konsolidasi "+v.CoaCode, "Tidak Dapat Menemukan Coa "+v.CoaCode)
				}
				tbID.Context = ctx

				// AmountAjeCr := *v.BalanceSheetCr + *tbID.ConsolidationDetailEntity.AmountJcteCr
				AmountJpmDr := *tbID.ConsolidationDetailEntity.AmountJpmDr - *v.BalanceSheetDr
				AmountJpmCr := +*tbID.ConsolidationDetailEntity.AmountJpmCr - *v.BalanceSheetCr
				tbID.ConsolidationDetailEntity.AmountJpmCr = &AmountJpmCr
				tbID.ConsolidationDetailEntity.AmountJpmDr = &AmountJpmDr
				// amount After Jpm
				AmountAfterJpm := *tbID.AmountBeforeJpm + *tbID.AmountJpmDr - *tbID.AmountJpmCr
				tbID.AmountAfterJpm = &AmountAfterJpm
				// amount After Jcte
				AmountAfterJcte := *tbID.AmountAfterJpm + *tbID.AmountJcteDr - *tbID.AmountJcteCr
				tbID.AmountAfterJcte = &AmountAfterJcte

				findConsoleBridge, err := s.Repository.FindByConsoleBridge(ctx, &JpmID.ConsolidationID)
				if err != nil {
					return response.CustomErrorBuilder(http.StatusBadRequest, "Tidak Dapat Menemukan Anak Usaha ", "Tidak Dapat Menemukan Anak Usaha ")
				}

				var dataAmountAnakUsahaa []float64
				for _, cb := range *findConsoleBridge {
					findConsoleBridgeDetail, err := s.Repository.FindByConsoleBridgeDetail(ctx, &cb.ID, &v.CoaCode)
					if err != nil {
						return response.CustomErrorBuilder(http.StatusBadRequest, "Tidak Dapat Menemukan Amount Anak Usaha ", "Tidak Dapat Menemukan Amount Anak Usaha ")
					}
					dataAmountAnakUsahaa = append(dataAmountAnakUsahaa, findConsoleBridgeDetail.Amount)
				}
				var sumAmountAnakUsaha float64 = 0.0

				for i := 0; i < len(dataAmountAnakUsahaa); i++ {
					sumAmountAnakUsaha += dataAmountAnakUsahaa[i]
				}

				var combineUnaudited float64 = *tbID.AmountAfterJcte + sumAmountAnakUsaha
				tbID.AmountCombineSubsidiary = &combineUnaudited

				var amountConsole float64 = *tbID.AmountCombineSubsidiary + *tbID.AmountJelimDr - *tbID.AmountJelimCr

				tbID.AmountConsole = &amountConsole

				_, err = s.Repository.UpdateTbd(ctx, &tbID.ID, tbID)
				if err != nil {
					return err
				}

				arrJpmDetails = append(arrJpmDetails, v)
			}
			if findCoaGroup.Name == "Liabilitas" {
				// if *JpmDetail.BalanceSheetDr > float64(0) {
				// 	return response.CustomErrorBuilder(http.StatusBadRequest, "Tidak Dapat Masukan Balance Sheet Debit", "Tidak Dapat Masukan Balance Sheet Debit")
				// }
				// if *JpmDetail.BalanceSheetCr == float64(0) {
				// 	return response.CustomErrorBuilder(http.StatusBadRequest, "Balance Sheet Credit Harus Di isi", "Balance Sheet Credit Harus Di isi")
				// }
				// if *JpmDetail.IncomeStatementCr > float64(0) {
				// 	return response.CustomErrorBuilder(http.StatusBadRequest, "Tidak Dapat Masukan Income Statement Credit", "Tidak Dapat Masukan Income Statement Credit")
				// }
				// if *JpmDetail.IncomeStatementDr > float64(0) {
				// 	return response.CustomErrorBuilder(http.StatusBadRequest, "Tidak Dapat Masukan Income Statement Debit", "Tidak Dapat Masukan Income Statement Debit")
				// }

				tbID, err := s.Repository.FindByTbd(ctx, &JpmID.ConsolidationID, &v.CoaCode)
				if err != nil {
					return response.CustomErrorBuilder(http.StatusBadRequest, "Tidak Dapat Menemukan Coa di Kombinasi Konsolidasi "+v.CoaCode, "Tidak Dapat Menemukan Coa "+v.CoaCode)
				}
				tbID.Context = ctx

				AmountJpmCr := *tbID.ConsolidationDetailEntity.AmountJpmCr - *v.BalanceSheetCr
				AmountAjeDr := *tbID.ConsolidationDetailEntity.AmountJpmDr - *v.BalanceSheetDr
				tbID.ConsolidationDetailEntity.AmountJpmCr = &AmountJpmCr
				tbID.ConsolidationDetailEntity.AmountJpmDr = &AmountAjeDr
				// amount After Jpm
				AmountAfterJpm := *tbID.AmountBeforeJpm - *tbID.AmountJpmDr + *tbID.AmountJpmCr
				tbID.AmountAfterJpm = &AmountAfterJpm
				// amount After Jcte
				AmountAfterJcte := *tbID.AmountAfterJpm - *tbID.AmountJcteDr + *tbID.AmountJcteCr
				tbID.AmountAfterJcte = &AmountAfterJcte

				findConsoleBridge, err := s.Repository.FindByConsoleBridge(ctx, &JpmID.ConsolidationID)
				if err != nil {
					return response.CustomErrorBuilder(http.StatusBadRequest, "Tidak Dapat Menemukan Anak Usaha ", "Tidak Dapat Menemukan Anak Usaha ")
				}

				var dataAmountAnakUsahaa []float64
				for _, cb := range *findConsoleBridge {
					findConsoleBridgeDetail, err := s.Repository.FindByConsoleBridgeDetail(ctx, &cb.ID, &v.CoaCode)
					_, err = s.Repository.UpdateTbd(ctx, &tbID.ID, tbID)
					if err != nil {
						return response.CustomErrorBuilder(http.StatusBadRequest, "Tidak Dapat Menemukan Amount Anak Usaha ", "Tidak Dapat Menemukan Amount Anak Usaha ")
					}
					dataAmountAnakUsahaa = append(dataAmountAnakUsahaa, findConsoleBridgeDetail.Amount)
				}
				var sumAmountAnakUsaha float64 = 0.0

				for i := 0; i < len(dataAmountAnakUsahaa); i++ {
					sumAmountAnakUsaha += dataAmountAnakUsahaa[i]
				}

				var combineUnaudited float64 = *tbID.AmountAfterJcte + sumAmountAnakUsaha
				tbID.AmountCombineSubsidiary = &combineUnaudited

				var amountConsole float64 = *tbID.AmountCombineSubsidiary - *tbID.AmountJelimDr + *tbID.AmountJelimCr

				tbID.AmountConsole = &amountConsole

				_, err = s.Repository.UpdateTbd(ctx, &tbID.ID, tbID)
				if err != nil {
					return err
				}
				arrJpmDetails = append(arrJpmDetails, v)
			}
			if findCoaGroup.Name == "EKUITAS" {
				// if *JpmDetail.BalanceSheetDr > float64(0) {
				// 	return response.CustomErrorBuilder(http.StatusBadRequest, "Tidak Dapat Masukan Balance Sheet Debit", "Tidak Dapat Masukan Balance Sheet Debit")
				// }
				// if *JpmDetail.BalanceSheetCr == float64(0) {
				// 	return response.CustomErrorBuilder(http.StatusBadRequest, "Balance Sheet Credit Harus Di isi", "Balance Sheet Credit Harus Di isi")
				// }
				// if *JpmDetail.IncomeStatementCr > float64(0) {
				// 	return response.CustomErrorBuilder(http.StatusBadRequest, "Tidak Dapat Masukan Income Statement Credit", "Tidak Dapat Masukan Income Statement Credit")
				// }
				// if *JpmDetail.IncomeStatementDr > float64(0) {
				// 	return response.CustomErrorBuilder(http.StatusBadRequest, "Tidak Dapat Masukan Income Statement Debit", "Tidak Dapat Masukan Income Statement Debit")
				// }

				tbID, err := s.Repository.FindByTbd(ctx, &JpmID.ConsolidationID, &v.CoaCode)
				if err != nil {
					return response.CustomErrorBuilder(http.StatusBadRequest, "Tidak Dapat Menemukan Coa di Kombinasi Konsolidasi "+v.CoaCode, "Tidak Dapat Menemukan Coa "+v.CoaCode)
				}
				tbID.Context = ctx

				AmountJpmCr := *tbID.ConsolidationDetailEntity.AmountJpmCr - *v.BalanceSheetCr
				AmountAjeDr := *tbID.ConsolidationDetailEntity.AmountJpmDr - *v.BalanceSheetDr
				tbID.ConsolidationDetailEntity.AmountJpmCr = &AmountJpmCr
				tbID.ConsolidationDetailEntity.AmountJpmDr = &AmountAjeDr
				// amount After Jpm
				AmountAfterJpm := *tbID.AmountBeforeJpm - *tbID.AmountJpmDr + *tbID.AmountJpmCr
				tbID.AmountAfterJpm = &AmountAfterJpm
				// amount After Jcte
				AmountAfterJcte := *tbID.AmountAfterJpm - *tbID.AmountJcteDr + *tbID.AmountJcteCr
				tbID.AmountAfterJcte = &AmountAfterJcte

				findConsoleBridge, err := s.Repository.FindByConsoleBridge(ctx, &JpmID.ConsolidationID)
				if err != nil {
					return response.CustomErrorBuilder(http.StatusBadRequest, "Tidak Dapat Menemukan Anak Usaha ", "Tidak Dapat Menemukan Anak Usaha ")
				}

				var dataAmountAnakUsahaa []float64
				for _, cb := range *findConsoleBridge {
					findConsoleBridgeDetail, err := s.Repository.FindByConsoleBridgeDetail(ctx, &cb.ID, &v.CoaCode)
					_, err = s.Repository.UpdateTbd(ctx, &tbID.ID, tbID)
					if err != nil {
						return response.CustomErrorBuilder(http.StatusBadRequest, "Tidak Dapat Menemukan Amount Anak Usaha ", "Tidak Dapat Menemukan Amount Anak Usaha ")
					}
					dataAmountAnakUsahaa = append(dataAmountAnakUsahaa, findConsoleBridgeDetail.Amount)
				}
				var sumAmountAnakUsaha float64 = 0.0

				for i := 0; i < len(dataAmountAnakUsahaa); i++ {
					sumAmountAnakUsaha += dataAmountAnakUsahaa[i]
				}

				var combineUnaudited float64 = *tbID.AmountAfterJcte + sumAmountAnakUsaha
				tbID.AmountCombineSubsidiary = &combineUnaudited

				var amountConsole float64 = *tbID.AmountCombineSubsidiary - *tbID.AmountJelimDr + *tbID.AmountJelimCr

				tbID.AmountConsole = &amountConsole

				_, err = s.Repository.UpdateTbd(ctx, &tbID.ID, tbID)
				if err != nil {
					return err
				}
				arrJpmDetails = append(arrJpmDetails, v)
			}
			if findCoaGroup.Name == "Pendapatan" {
				// if *JpmDetail.BalanceSheetDr != float64(0) {
				// 	return response.CustomErrorBuilder(http.StatusBadRequest, "Tidak Dapat Masukan Balance Sheet Debit", "Tidak Dapat Masukan Balance Sheet Debit")
				// }
				// if *JpmDetail.BalanceSheetCr != float64(0) {
				// 	return response.CustomErrorBuilder(http.StatusBadRequest, "Tidak Dapat Masukan Balance Sheet Credit", "Tidak Dapat Masukan Balance Sheet Credit")
				// }
				// if *JpmDetail.IncomeStatementCr == float64(0) {
				// 	return response.CustomErrorBuilder(http.StatusBadRequest, "Income Statement Credit Harus Di isi", "Income Statement Credit Harus Di isi")
				// }
				// if *JpmDetail.IncomeStatementDr != float64(0) {
				// 	return response.CustomErrorBuilder(http.StatusBadRequest, "Tidak Dapat Masukan Income Statement Debit", "Tidak Dapat Masukan Income Statement Debit")
				// }

				tbID, err := s.Repository.FindByTbd(ctx, &JpmID.ConsolidationID, &v.CoaCode)
				if err != nil {
					return response.CustomErrorBuilder(http.StatusBadRequest, "Tidak Dapat Menemukan Coa di Kombinasi Konsolidasi "+v.CoaCode, "Tidak Dapat Menemukan Coa "+v.CoaCode)
				}
				tbID.Context = ctx

				AmountJpmCr := *tbID.ConsolidationDetailEntity.AmountJpmCr - *v.IncomeStatementCr
				AmountAjeDr := *tbID.ConsolidationDetailEntity.AmountJpmDr - *v.IncomeStatementDr
				tbID.ConsolidationDetailEntity.AmountJpmCr = &AmountJpmCr
				tbID.ConsolidationDetailEntity.AmountJpmDr = &AmountAjeDr
				// amount After Jpm
				AmountAfterJpm := *tbID.AmountBeforeJpm - *tbID.AmountJpmDr + *tbID.AmountJpmCr
				tbID.AmountAfterJpm = &AmountAfterJpm
				// amount After Jcte
				AmountAfterJcte := *tbID.AmountAfterJpm - *tbID.AmountJcteDr + *tbID.AmountJcteCr
				tbID.AmountAfterJcte = &AmountAfterJcte

				findConsoleBridge, err := s.Repository.FindByConsoleBridge(ctx, &JpmID.ConsolidationID)
				if err != nil {
					return response.CustomErrorBuilder(http.StatusBadRequest, "Tidak Dapat Menemukan Anak Usaha ", "Tidak Dapat Menemukan Anak Usaha ")
				}

				var dataAmountAnakUsahaa []float64
				for _, cb := range *findConsoleBridge {
					findConsoleBridgeDetail, err := s.Repository.FindByConsoleBridgeDetail(ctx, &cb.ID, &v.CoaCode)
					_, err = s.Repository.UpdateTbd(ctx, &tbID.ID, tbID)
					if err != nil {
						return response.CustomErrorBuilder(http.StatusBadRequest, "Tidak Dapat Menemukan Amount Anak Usaha ", "Tidak Dapat Menemukan Amount Anak Usaha ")
					}
					dataAmountAnakUsahaa = append(dataAmountAnakUsahaa, findConsoleBridgeDetail.Amount)
				}
				var sumAmountAnakUsaha float64 = 0.0

				for i := 0; i < len(dataAmountAnakUsahaa); i++ {
					sumAmountAnakUsaha += dataAmountAnakUsahaa[i]
				}

				var combineUnaudited float64 = *tbID.AmountAfterJcte + sumAmountAnakUsaha
				tbID.AmountCombineSubsidiary = &combineUnaudited

				var amountConsole float64 = *tbID.AmountCombineSubsidiary - *tbID.AmountJelimDr + *tbID.AmountJelimCr

				tbID.AmountConsole = &amountConsole

				_, err = s.Repository.UpdateTbd(ctx, &tbID.ID, tbID)
				if err != nil {
					return err
				}
				arrJpmDetails = append(arrJpmDetails, v)
			}
			if findCoaGroup.Name == "HPP/COGS" {
				// if *JpmDetail.BalanceSheetDr != float64(0) {
				// 	return response.CustomErrorBuilder(http.StatusBadRequest, "Tidak Dapat Masukan Balance Sheet Debit", "Tidak Dapat Masukan Balance Sheet Debit")
				// }
				// if *JpmDetail.BalanceSheetCr != float64(0) {
				// 	return response.CustomErrorBuilder(http.StatusBadRequest, "Tidak Dapat Masukan Balance Sheet Credit", "Tidak Dapat Masukan Balance Sheet Credit")
				// }
				// if *JpmDetail.IncomeStatementCr != float64(0) {
				// 	return response.CustomErrorBuilder(http.StatusBadRequest, "Tidak Dapat Masukan Income Statement Credit", "Tidak Dapat Masukan Income Statement Credit")
				// }
				// if *JpmDetail.IncomeStatementDr == float64(0) {
				// 	return response.CustomErrorBuilder(http.StatusBadRequest, "Income Statement Debit Harus Di isi", "Income Statement Debit Harus Di isi")
				// }

				tbID, err := s.Repository.FindByTbd(ctx, &JpmID.ConsolidationID, &v.CoaCode)
				if err != nil {
					return response.CustomErrorBuilder(http.StatusBadRequest, "Tidak Dapat Menemukan Coa di Kombinasi Konsolidasi "+v.CoaCode, "Tidak Dapat Menemukan Coa "+v.CoaCode)
				}
				tbID.Context = ctx

				AmountAjeCr := *tbID.ConsolidationDetailEntity.AmountJpmCr - *v.IncomeStatementCr
				AmountJpmDr := *tbID.ConsolidationDetailEntity.AmountJpmDr - *v.IncomeStatementDr
				tbID.ConsolidationDetailEntity.AmountJpmCr = &AmountAjeCr
				tbID.ConsolidationDetailEntity.AmountJpmDr = &AmountJpmDr
				AmountAfterJpm := *tbID.AmountBeforeJpm + *tbID.AmountJpmDr - *tbID.AmountJpmCr
				tbID.AmountAfterJpm = &AmountAfterJpm
				AmountAfterJcte := *tbID.AmountAfterJpm + *tbID.AmountJcteDr - *tbID.AmountJcteCr
				tbID.AmountAfterJcte = &AmountAfterJcte

				findConsoleBridge, err := s.Repository.FindByConsoleBridge(ctx, &JpmID.ConsolidationID)
				if err != nil {
					return response.CustomErrorBuilder(http.StatusBadRequest, "Tidak Dapat Menemukan Anak Usaha ", "Tidak Dapat Menemukan Anak Usaha ")
				}

				var dataAmountAnakUsahaa []float64
				for _, cb := range *findConsoleBridge {
					findConsoleBridgeDetail, err := s.Repository.FindByConsoleBridgeDetail(ctx, &cb.ID, &v.CoaCode)
					_, err = s.Repository.UpdateTbd(ctx, &tbID.ID, tbID)
					if err != nil {
						return response.CustomErrorBuilder(http.StatusBadRequest, "Tidak Dapat Menemukan Amount Anak Usaha ", "Tidak Dapat Menemukan Amount Anak Usaha ")
					}
					dataAmountAnakUsahaa = append(dataAmountAnakUsahaa, findConsoleBridgeDetail.Amount)
				}
				var sumAmountAnakUsaha float64 = 0.0

				for i := 0; i < len(dataAmountAnakUsahaa); i++ {
					sumAmountAnakUsaha += dataAmountAnakUsahaa[i]
				}

				var combineUnaudited float64 = *tbID.AmountAfterJcte + sumAmountAnakUsaha
				tbID.AmountCombineSubsidiary = &combineUnaudited

				var amountConsole float64 = *tbID.AmountCombineSubsidiary + *tbID.AmountJelimDr - *tbID.AmountJelimCr

				tbID.AmountConsole = &amountConsole

				_, err = s.Repository.UpdateTbd(ctx, &tbID.ID, tbID)
				if err != nil {
					return err
				}
				arrJpmDetails = append(arrJpmDetails, v)
			}
			if findCoaGroup.Name == "Selling Expense" {
				// if *JpmDetail.BalanceSheetDr != float64(0) {
				// 	return response.CustomErrorBuilder(http.StatusBadRequest, "Tidak Dapat Masukan Balance Sheet Debit", "Tidak Dapat Masukan Balance Sheet Debit")
				// }
				// if *JpmDetail.BalanceSheetCr != float64(0) {
				// 	return response.CustomErrorBuilder(http.StatusBadRequest, "Tidak Dapat Masukan Balance Sheet Credit", "Tidak Dapat Masukan Balance Sheet Credit")
				// }
				// if *JpmDetail.IncomeStatementCr != float64(0) {
				// 	return response.CustomErrorBuilder(http.StatusBadRequest, "Tidak Dapat Masukan Income Statement Credit", "Tidak Dapat Masukan Income Statement Credit")
				// }
				// if *JpmDetail.IncomeStatementDr == float64(0) {
				// 	return response.CustomErrorBuilder(http.StatusBadRequest, "Income Statement Debit Harus Di isi", "Income Statement Debit Harus Di isi")
				// }

				tbID, err := s.Repository.FindByTbd(ctx, &JpmID.ConsolidationID, &v.CoaCode)
				if err != nil {
					return response.CustomErrorBuilder(http.StatusBadRequest, "Tidak Dapat Menemukan Coa di Kombinasi Konsolidasi "+v.CoaCode, "Tidak Dapat Menemukan Coa "+v.CoaCode)
				}
				tbID.Context = ctx

				AmountAjeCr := *tbID.ConsolidationDetailEntity.AmountJpmCr - *v.IncomeStatementCr
				AmountJpmDr := *tbID.ConsolidationDetailEntity.AmountJpmDr - *v.IncomeStatementDr
				tbID.ConsolidationDetailEntity.AmountJpmCr = &AmountAjeCr
				tbID.ConsolidationDetailEntity.AmountJpmDr = &AmountJpmDr
				AmountAfterJpm := *tbID.AmountBeforeJpm + *tbID.AmountJpmDr - *tbID.AmountJpmCr
				tbID.AmountAfterJpm = &AmountAfterJpm
				AmountAfterJcte := *tbID.AmountAfterJpm + *tbID.AmountJcteDr - *tbID.AmountJcteCr
				tbID.AmountAfterJcte = &AmountAfterJcte

				findConsoleBridge, err := s.Repository.FindByConsoleBridge(ctx, &JpmID.ConsolidationID)
				if err != nil {
					return response.CustomErrorBuilder(http.StatusBadRequest, "Tidak Dapat Menemukan Anak Usaha ", "Tidak Dapat Menemukan Anak Usaha ")
				}

				var dataAmountAnakUsahaa []float64
				for _, cb := range *findConsoleBridge {
					findConsoleBridgeDetail, err := s.Repository.FindByConsoleBridgeDetail(ctx, &cb.ID, &v.CoaCode)
					_, err = s.Repository.UpdateTbd(ctx, &tbID.ID, tbID)
					if err != nil {
						return response.CustomErrorBuilder(http.StatusBadRequest, "Tidak Dapat Menemukan Amount Anak Usaha ", "Tidak Dapat Menemukan Amount Anak Usaha ")
					}
					dataAmountAnakUsahaa = append(dataAmountAnakUsahaa, findConsoleBridgeDetail.Amount)
				}
				var sumAmountAnakUsaha float64 = 0.0

				for i := 0; i < len(dataAmountAnakUsahaa); i++ {
					sumAmountAnakUsaha += dataAmountAnakUsahaa[i]
				}

				var combineUnaudited float64 = *tbID.AmountAfterJcte + sumAmountAnakUsaha
				tbID.AmountCombineSubsidiary = &combineUnaudited

				var amountConsole float64 = *tbID.AmountCombineSubsidiary + *tbID.AmountJelimDr - *tbID.AmountJelimCr

				tbID.AmountConsole = &amountConsole

				_, err = s.Repository.UpdateTbd(ctx, &tbID.ID, tbID)
				if err != nil {
					return err
				}
				arrJpmDetails = append(arrJpmDetails, v)
			}
			if findCoaGroup.Name == "General & Admin Expense" {
				// if *JpmDetail.BalanceSheetDr != float64(0) {
				// 	return response.CustomErrorBuilder(http.StatusBadRequest, "Tidak Dapat Masukan Balance Sheet Debit", "Tidak Dapat Masukan Balance Sheet Debit")
				// }
				// if *JpmDetail.BalanceSheetCr != float64(0) {
				// 	return response.CustomErrorBuilder(http.StatusBadRequest, "Tidak Dapat Masukan Balance Sheet Credit", "Tidak Dapat Masukan Balance Sheet Credit")
				// }
				// if *JpmDetail.IncomeStatementCr != float64(0) {
				// 	return response.CustomErrorBuilder(http.StatusBadRequest, "Tidak Dapat Masukan Income Statement Credit", "Tidak Dapat Masukan Income Statement Credit")
				// }
				// if *JpmDetail.IncomeStatementDr == float64(0) {
				// 	return response.CustomErrorBuilder(http.StatusBadRequest, "Income Statement Debit Harus Di isi", "Income Statement Debit Harus Di isi")
				// }

				tbID, err := s.Repository.FindByTbd(ctx, &JpmID.ConsolidationID, &v.CoaCode)
				if err != nil {
					return response.CustomErrorBuilder(http.StatusBadRequest, "Tidak Dapat Menemukan Coa di Kombinasi Konsolidasi "+v.CoaCode, "Tidak Dapat Menemukan Coa "+v.CoaCode)
				}
				tbID.Context = ctx

				AmountAjeCr := *tbID.ConsolidationDetailEntity.AmountJpmCr - *v.IncomeStatementCr
				AmountJpmDr := *tbID.ConsolidationDetailEntity.AmountJpmDr - *v.IncomeStatementDr
				tbID.ConsolidationDetailEntity.AmountJpmCr = &AmountAjeCr
				tbID.ConsolidationDetailEntity.AmountJpmDr = &AmountJpmDr
				AmountAfterJpm := *tbID.AmountBeforeJpm + *tbID.AmountJpmDr - *tbID.AmountJpmCr
				tbID.AmountAfterJpm = &AmountAfterJpm
				AmountAfterJcte := *tbID.AmountAfterJpm + *tbID.AmountJcteDr - *tbID.AmountJcteCr
				tbID.AmountAfterJcte = &AmountAfterJcte

				findConsoleBridge, err := s.Repository.FindByConsoleBridge(ctx, &JpmID.ConsolidationID)
				if err != nil {
					return response.CustomErrorBuilder(http.StatusBadRequest, "Tidak Dapat Menemukan Anak Usaha ", "Tidak Dapat Menemukan Anak Usaha ")
				}

				var dataAmountAnakUsahaa []float64
				for _, cb := range *findConsoleBridge {
					findConsoleBridgeDetail, err := s.Repository.FindByConsoleBridgeDetail(ctx, &cb.ID, &v.CoaCode)
					_, err = s.Repository.UpdateTbd(ctx, &tbID.ID, tbID)
					if err != nil {
						return response.CustomErrorBuilder(http.StatusBadRequest, "Tidak Dapat Menemukan Amount Anak Usaha ", "Tidak Dapat Menemukan Amount Anak Usaha ")
					}
					dataAmountAnakUsahaa = append(dataAmountAnakUsahaa, findConsoleBridgeDetail.Amount)
				}
				var sumAmountAnakUsaha float64 = 0.0

				for i := 0; i < len(dataAmountAnakUsahaa); i++ {
					sumAmountAnakUsaha += dataAmountAnakUsahaa[i]
				}

				var combineUnaudited float64 = *tbID.AmountAfterJcte + sumAmountAnakUsaha
				tbID.AmountCombineSubsidiary = &combineUnaudited

				var amountConsole float64 = *tbID.AmountCombineSubsidiary + *tbID.AmountJelimDr - *tbID.AmountJelimCr

				tbID.AmountConsole = &amountConsole

				_, err = s.Repository.UpdateTbd(ctx, &tbID.ID, tbID)
				if err != nil {
					return err
				}
				arrJpmDetails = append(arrJpmDetails, v)
			}
			if findCoaGroup.Name == "Other Income" {
				// if *JpmDetail.BalanceSheetDr != float64(0) {
				// 	return response.CustomErrorBuilder(http.StatusBadRequest, "Tidak Dapat Masukan Balance Sheet Debit", "Tidak Dapat Masukan Balance Sheet Debit")
				// }
				// if *JpmDetail.BalanceSheetCr != float64(0) {
				// 	return response.CustomErrorBuilder(http.StatusBadRequest, "Tidak Dapat Masukan Balance Sheet Credit", "Tidak Dapat Masukan Balance Sheet Credit")
				// }
				// if *JpmDetail.IncomeStatementCr == float64(0) {
				// 	return response.CustomErrorBuilder(http.StatusBadRequest, "Income Statement Credit Harus Di isi", "Income Statement Credit Harus Di isi")
				// }
				// if *JpmDetail.IncomeStatementDr != float64(0) {
				// 	return response.CustomErrorBuilder(http.StatusBadRequest, "Tidak Dapat Masukan Income Statement Debit", "Tidak Dapat Masukan Income Statement Debit")
				// }

				tbID, err := s.Repository.FindByTbd(ctx, &JpmID.ConsolidationID, &v.CoaCode)
				if err != nil {
					return response.CustomErrorBuilder(http.StatusBadRequest, "Tidak Dapat Menemukan Coa di Kombinasi Konsolidasi "+v.CoaCode, "Tidak Dapat Menemukan Coa "+v.CoaCode)
				}
				tbID.Context = ctx

				AmountJpmCr := *tbID.ConsolidationDetailEntity.AmountJpmCr - *v.IncomeStatementCr
				AmountAjeDr := *tbID.ConsolidationDetailEntity.AmountJpmDr - *v.IncomeStatementDr
				tbID.ConsolidationDetailEntity.AmountJpmCr = &AmountJpmCr
				tbID.ConsolidationDetailEntity.AmountJcteDr = &AmountAjeDr
				AmountAfterJpm := *tbID.AmountBeforeJpm - *tbID.AmountJpmDr + *tbID.AmountJpmCr
				tbID.AmountAfterJpm = &AmountAfterJpm
				AmountAfterJcte := *tbID.AmountAfterJpm - *tbID.AmountJcteDr + *tbID.AmountJcteCr
				tbID.AmountAfterJcte = &AmountAfterJcte

				findConsoleBridge, err := s.Repository.FindByConsoleBridge(ctx, &JpmID.ConsolidationID)
				if err != nil {
					return response.CustomErrorBuilder(http.StatusBadRequest, "Tidak Dapat Menemukan Anak Usaha ", "Tidak Dapat Menemukan Anak Usaha ")
				}

				var dataAmountAnakUsahaa []float64
				for _, cb := range *findConsoleBridge {
					findConsoleBridgeDetail, err := s.Repository.FindByConsoleBridgeDetail(ctx, &cb.ID, &v.CoaCode)
					_, err = s.Repository.UpdateTbd(ctx, &tbID.ID, tbID)
					if err != nil {
						return response.CustomErrorBuilder(http.StatusBadRequest, "Tidak Dapat Menemukan Amount Anak Usaha ", "Tidak Dapat Menemukan Amount Anak Usaha ")
					}
					dataAmountAnakUsahaa = append(dataAmountAnakUsahaa, findConsoleBridgeDetail.Amount)
				}
				var sumAmountAnakUsaha float64 = 0.0

				for i := 0; i < len(dataAmountAnakUsahaa); i++ {
					sumAmountAnakUsaha += dataAmountAnakUsahaa[i]
				}

				var totalAmount float64 = *tbID.AmountAfterJcte + sumAmountAnakUsaha
				tbID.AmountCombineSubsidiary = &totalAmount

				var amountConsole float64 = *tbID.AmountCombineSubsidiary - *tbID.AmountJelimDr + *tbID.AmountJelimCr

				tbID.AmountConsole = &amountConsole

				_, err = s.Repository.UpdateTbd(ctx, &tbID.ID, tbID)
				if err != nil {
					return err
				}
				arrJpmDetails = append(arrJpmDetails, v)
			}
			if findCoaGroup.Name == "Other Expense" {
				// if *JpmDetail.BalanceSheetDr != float64(0) {
				// 	return response.CustomErrorBuilder(http.StatusBadRequest, "Tidak Dapat Masukan Balance Sheet Debit", "Tidak Dapat Masukan Balance Sheet Debit")
				// }
				// if *JpmDetail.BalanceSheetCr != float64(0) {
				// 	return response.CustomErrorBuilder(http.StatusBadRequest, "Tidak Dapat Masukan Balance Sheet Credit", "Tidak Dapat Masukan Balance Sheet Credit")
				// }
				// if *JpmDetail.IncomeStatementCr != float64(0) {
				// 	return response.CustomErrorBuilder(http.StatusBadRequest, "Tidak Dapat Masukan Income Statement Credit", "Tidak Dapat Masukan Income Statement Credit")
				// }
				// if *JpmDetail.IncomeStatementDr == float64(0) {
				// 	return response.CustomErrorBuilder(http.StatusBadRequest, "Income Statement Debit Harus Di isi", "Income Statement Debit Harus Di isi")
				// }

				tbID, err := s.Repository.FindByTbd(ctx, &JpmID.ConsolidationID, &v.CoaCode)
				if err != nil {
					return response.CustomErrorBuilder(http.StatusBadRequest, "Tidak Dapat Menemukan Coa di Kombinasi Konsolidasi "+v.CoaCode, "Tidak Dapat Menemukan Coa "+v.CoaCode)
				}
				tbID.Context = ctx

				AmountAjeCr := *tbID.ConsolidationDetailEntity.AmountJpmCr - *v.IncomeStatementCr
				AmountJpmDr := *tbID.ConsolidationDetailEntity.AmountJpmDr - *v.IncomeStatementDr
				tbID.ConsolidationDetailEntity.AmountJpmCr = &AmountAjeCr
				tbID.ConsolidationDetailEntity.AmountJpmDr = &AmountJpmDr
				AmountAfterJpm := *tbID.AmountBeforeJpm + *tbID.AmountJpmDr - *tbID.AmountJpmCr
				tbID.AmountAfterJpm = &AmountAfterJpm
				AmountAfterJcte := *tbID.AmountAfterJpm + *tbID.AmountJcteDr - *tbID.AmountJcteCr
				tbID.AmountAfterJcte = &AmountAfterJcte

				findConsoleBridge, err := s.Repository.FindByConsoleBridge(ctx, &JpmID.ConsolidationID)
				if err != nil {
					return response.CustomErrorBuilder(http.StatusBadRequest, "Tidak Dapat Menemukan Anak Usaha ", "Tidak Dapat Menemukan Anak Usaha ")
				}

				var dataAmountAnakUsahaa []float64
				for _, cb := range *findConsoleBridge {
					findConsoleBridgeDetail, err := s.Repository.FindByConsoleBridgeDetail(ctx, &cb.ID, &v.CoaCode)
					_, err = s.Repository.UpdateTbd(ctx, &tbID.ID, tbID)
					if err != nil {
						return response.CustomErrorBuilder(http.StatusBadRequest, "Tidak Dapat Menemukan Amount Anak Usaha ", "Tidak Dapat Menemukan Amount Anak Usaha ")
					}
					dataAmountAnakUsahaa = append(dataAmountAnakUsahaa, findConsoleBridgeDetail.Amount)
				}
				var sumAmountAnakUsaha float64 = 0.0

				for i := 0; i < len(dataAmountAnakUsahaa); i++ {
					sumAmountAnakUsaha += dataAmountAnakUsahaa[i]
				}

				var combineUnaudited float64 = *tbID.AmountAfterJcte + sumAmountAnakUsaha
				tbID.AmountCombineSubsidiary = &combineUnaudited

				var amountConsole float64 = *tbID.AmountCombineSubsidiary + *tbID.AmountJelimDr - *tbID.AmountJelimCr

				tbID.AmountConsole = &amountConsole

				_, err = s.Repository.UpdateTbd(ctx, &tbID.ID, tbID)
				if err != nil {
					return err
				}
				arrJpmDetails = append(arrJpmDetails, v)
			}
			if findCoaGroup.Name == "Tax Expense" {
				// if *JpmDetail.BalanceSheetDr != float64(0) {
				// 	return response.CustomErrorBuilder(http.StatusBadRequest, "Tidak Dapat Masukan Balance Sheet Debit", "Tidak Dapat Masukan Balance Sheet Debit")
				// }
				// if *JpmDetail.BalanceSheetCr != float64(0) {
				// 	return response.CustomErrorBuilder(http.StatusBadRequest, "Tidak Dapat Masukan Balance Sheet Credit", "Tidak Dapat Masukan Balance Sheet Credit")
				// }
				// if *JpmDetail.IncomeStatementCr != float64(0) {
				// 	return response.CustomErrorBuilder(http.StatusBadRequest, "Tidak Dapat Masukan Income Statement Credit", "Tidak Dapat Masukan Income Statement Credit")
				// }
				// if *JpmDetail.IncomeStatementDr == float64(0) {
				// 	return response.CustomErrorBuilder(http.StatusBadRequest, "Income Statement Debit Harus Di isi", "Income Statement Debit Harus Di isi")
				// }

				tbID, err := s.Repository.FindByTbd(ctx, &JpmID.ConsolidationID, &v.CoaCode)
				if err != nil {
					return response.CustomErrorBuilder(http.StatusBadRequest, "Tidak Dapat Menemukan Coa di Kombinasi Konsolidasi "+v.CoaCode, "Tidak Dapat Menemukan Coa "+v.CoaCode)
				}
				tbID.Context = ctx

				AmountAjeCr := *tbID.ConsolidationDetailEntity.AmountJpmCr - *v.IncomeStatementCr
				AmountJpmDr := *tbID.ConsolidationDetailEntity.AmountJpmDr - *v.IncomeStatementDr
				tbID.ConsolidationDetailEntity.AmountJpmCr = &AmountAjeCr
				tbID.ConsolidationDetailEntity.AmountJpmDr = &AmountJpmDr
				AmountAfterJpm := *tbID.AmountBeforeJpm + *tbID.AmountJpmDr - *tbID.AmountJpmCr
				tbID.AmountAfterJpm = &AmountAfterJpm
				AmountAfterJcte := *tbID.AmountAfterJpm + *tbID.AmountJcteDr - *tbID.AmountJcteCr
				tbID.AmountAfterJcte = &AmountAfterJcte

				findConsoleBridge, err := s.Repository.FindByConsoleBridge(ctx, &JpmID.ConsolidationID)
				if err != nil {
					return response.CustomErrorBuilder(http.StatusBadRequest, "Tidak Dapat Menemukan Anak Usaha ", "Tidak Dapat Menemukan Anak Usaha ")
				}

				var dataAmountAnakUsahaa []float64
				for _, cb := range *findConsoleBridge {
					findConsoleBridgeDetail, err := s.Repository.FindByConsoleBridgeDetail(ctx, &cb.ID, &v.CoaCode)
					_, err = s.Repository.UpdateTbd(ctx, &tbID.ID, tbID)
					if err != nil {
						return response.CustomErrorBuilder(http.StatusBadRequest, "Tidak Dapat Menemukan Amount Anak Usaha ", "Tidak Dapat Menemukan Amount Anak Usaha ")
					}
					dataAmountAnakUsahaa = append(dataAmountAnakUsahaa, findConsoleBridgeDetail.Amount)
				}
				var sumAmountAnakUsaha float64 = 0.0

				for i := 0; i < len(dataAmountAnakUsahaa); i++ {
					sumAmountAnakUsaha += dataAmountAnakUsahaa[i]
				}

				var combineUnaudited float64 = *tbID.AmountAfterJcte + sumAmountAnakUsaha
				tbID.AmountCombineSubsidiary = &combineUnaudited

				var amountConsole float64 = *tbID.AmountCombineSubsidiary + *tbID.AmountJelimDr - *tbID.AmountJelimCr

				tbID.AmountConsole = &amountConsole

				_, err = s.Repository.UpdateTbd(ctx, &tbID.ID, tbID)
				if err != nil {
					return err
				}
				arrJpmDetails = append(arrJpmDetails, v)
			}
			if findCoaGroup.Name == "Income (Loss) from subsidiary" {
				// if *JpmDetail.BalanceSheetDr != float64(0) {
				// 	return response.CustomErrorBuilder(http.StatusBadRequest, "Tidak Dapat Masukan Balance Sheet Debit", "Tidak Dapat Masukan Balance Sheet Debit")
				// }
				// if *JpmDetail.BalanceSheetCr != float64(0) {
				// 	return response.CustomErrorBuilder(http.StatusBadRequest, "Tidak Dapat Masukan Balance Sheet Credit", "Tidak Dapat Masukan Balance Sheet Credit")
				// }
				// if *JpmDetail.IncomeStatementCr == float64(0) {
				// 	return response.CustomErrorBuilder(http.StatusBadRequest, "Income Statement Credit Harus Di isi", "Income Statement Credit Harus Di isi")
				// }
				// if *JpmDetail.IncomeStatementDr != float64(0) {
				// 	return response.CustomErrorBuilder(http.StatusBadRequest, "Tidak Dapat Masukan Income Statement Debit", "Tidak Dapat Masukan Income Statement Debit")
				// }

				tbID, err := s.Repository.FindByTbd(ctx, &JpmID.ConsolidationID, &v.CoaCode)
				if err != nil {
					return response.CustomErrorBuilder(http.StatusBadRequest, "Tidak Dapat Menemukan Coa di Kombinasi Konsolidasi "+v.CoaCode, "Tidak Dapat Menemukan Coa "+v.CoaCode)
				}
				tbID.Context = ctx

				AmountJpmCr := *tbID.ConsolidationDetailEntity.AmountJpmCr - *v.IncomeStatementCr
				AmountAjeDr := *tbID.ConsolidationDetailEntity.AmountJpmDr - *v.IncomeStatementDr
				tbID.ConsolidationDetailEntity.AmountJpmCr = &AmountJpmCr
				tbID.ConsolidationDetailEntity.AmountJcteDr = &AmountAjeDr
				AmountAfterJpm := *tbID.AmountBeforeJpm - *tbID.AmountJpmDr + *tbID.AmountJpmCr
				tbID.AmountAfterJpm = &AmountAfterJpm
				AmountAfterJcte := *tbID.AmountAfterJpm - *tbID.AmountJcteDr + *tbID.AmountJcteCr
				tbID.AmountAfterJcte = &AmountAfterJcte

				findConsoleBridge, err := s.Repository.FindByConsoleBridge(ctx, &JpmID.ConsolidationID)
				if err != nil {
					return response.CustomErrorBuilder(http.StatusBadRequest, "Tidak Dapat Menemukan Anak Usaha ", "Tidak Dapat Menemukan Anak Usaha ")
				}

				var dataAmountAnakUsahaa []float64
				for _, cb := range *findConsoleBridge {
					findConsoleBridgeDetail, err := s.Repository.FindByConsoleBridgeDetail(ctx, &cb.ID, &v.CoaCode)
					_, err = s.Repository.UpdateTbd(ctx, &tbID.ID, tbID)
					if err != nil {
						return response.CustomErrorBuilder(http.StatusBadRequest, "Tidak Dapat Menemukan Amount Anak Usaha ", "Tidak Dapat Menemukan Amount Anak Usaha ")
					}
					dataAmountAnakUsahaa = append(dataAmountAnakUsahaa, findConsoleBridgeDetail.Amount)
				}
				var sumAmountAnakUsaha float64 = 0.0

				for i := 0; i < len(dataAmountAnakUsahaa); i++ {
					sumAmountAnakUsaha += dataAmountAnakUsahaa[i]
				}

				var totalAmount float64 = *tbID.AmountAfterJcte + sumAmountAnakUsaha
				tbID.AmountCombineSubsidiary = &totalAmount

				var amountConsole float64 = *tbID.AmountCombineSubsidiary - *tbID.AmountJelimDr + *tbID.AmountJelimCr

				tbID.AmountConsole = &amountConsole

				_, err = s.Repository.UpdateTbd(ctx, &tbID.ID, tbID)
				if err != nil {
					return err
				}
				arrJpmDetails = append(arrJpmDetails, v)
			}
			if findCoaGroup.Name == "MINORITY INTEREST IN NET INCOME (NCI)" {
				// if *JpmDetail.BalanceSheetDr != float64(0) {
				// 	return response.CustomErrorBuilder(http.StatusBadRequest, "Tidak Dapat Masukan Balance Sheet Debit", "Tidak Dapat Masukan Balance Sheet Debit")
				// }
				// if *JpmDetail.BalanceSheetCr != float64(0) {
				// 	return response.CustomErrorBuilder(http.StatusBadRequest, "Tidak Dapat Masukan Balance Sheet Credit", "Tidak Dapat Masukan Balance Sheet Credit")
				// }
				// if *JpmDetail.IncomeStatementCr == float64(0) {
				// 	return response.CustomErrorBuilder(http.StatusBadRequest, "Income Statement Credit Harus Di isi", "Income Statement Credit Harus Di isi")
				// }
				// if *JpmDetail.IncomeStatementDr != float64(0) {
				// 	return response.CustomErrorBuilder(http.StatusBadRequest, "Tidak Dapat Masukan Income Statement Debit", "Tidak Dapat Masukan Income Statement Debit")
				// }

				tbID, err := s.Repository.FindByTbd(ctx, &JpmID.ConsolidationID, &v.CoaCode)
				if err != nil {
					return response.CustomErrorBuilder(http.StatusBadRequest, "Tidak Dapat Menemukan Coa di Kombinasi Konsolidasi "+v.CoaCode, "Tidak Dapat Menemukan Coa "+v.CoaCode)
				}
				tbID.Context = ctx

				AmountJpmCr := *tbID.ConsolidationDetailEntity.AmountJpmCr - *v.IncomeStatementCr
				AmountAjeDr := *tbID.ConsolidationDetailEntity.AmountJpmDr - *v.IncomeStatementDr
				tbID.ConsolidationDetailEntity.AmountJpmCr = &AmountJpmCr
				tbID.ConsolidationDetailEntity.AmountJpmDr = &AmountAjeDr
				AmountAfterJpm := *tbID.AmountBeforeJpm - *tbID.AmountJpmDr + *tbID.AmountJpmCr
				tbID.AmountAfterJpm = &AmountAfterJpm
				AmountAfterJcte := *tbID.AmountAfterJpm - *tbID.AmountJcteDr + *tbID.AmountJcteCr
				tbID.AmountAfterJcte = &AmountAfterJcte

				findConsoleBridge, err := s.Repository.FindByConsoleBridge(ctx, &JpmID.ConsolidationID)
				if err != nil {
					return response.CustomErrorBuilder(http.StatusBadRequest, "Tidak Dapat Menemukan Anak Usaha ", "Tidak Dapat Menemukan Anak Usaha ")
				}

				var dataAmountAnakUsahaa []float64
				for _, cb := range *findConsoleBridge {
					findConsoleBridgeDetail, err := s.Repository.FindByConsoleBridgeDetail(ctx, &cb.ID, &v.CoaCode)
					_, err = s.Repository.UpdateTbd(ctx, &tbID.ID, tbID)
					if err != nil {
						return response.CustomErrorBuilder(http.StatusBadRequest, "Tidak Dapat Menemukan Amount Anak Usaha ", "Tidak Dapat Menemukan Amount Anak Usaha ")
					}
					dataAmountAnakUsahaa = append(dataAmountAnakUsahaa, findConsoleBridgeDetail.Amount)
				}
				var sumAmountAnakUsaha float64 = 0.0

				for i := 0; i < len(dataAmountAnakUsahaa); i++ {
					sumAmountAnakUsaha += dataAmountAnakUsahaa[i]
				}

				var totalAmount float64 = *tbID.AmountAfterJcte + sumAmountAnakUsaha
				tbID.AmountCombineSubsidiary = &totalAmount

				var amountConsole float64 = *tbID.AmountCombineSubsidiary - *tbID.AmountJelimDr + *tbID.AmountJelimCr

				tbID.AmountConsole = &amountConsole

				_, err = s.Repository.UpdateTbd(ctx, &tbID.ID, tbID)
				if err != nil {
					return err
				}
				arrJpmDetails = append(arrJpmDetails, v)
			}
			if findCoaGroup.Name == "Other Comprehensive Income" {

				// if *JpmDetail.BalanceSheetDr != float64(0) {
				// 	return response.CustomErrorBuilder(http.StatusBadRequest, "Tidak Dapat Masukan Balance Sheet Debit", "Tidak Dapat Masukan Balance Sheet Debit")
				// }
				// if *JpmDetail.BalanceSheetCr != float64(0) {
				// 	return response.CustomErrorBuilder(http.StatusBadRequest, "Tidak Dapat Masukan Balance Sheet Credit", "Tidak Dapat Masukan Balance Sheet Credit")
				// }
				// if *JpmDetail.IncomeStatementCr == float64(0) {
				// 	return response.CustomErrorBuilder(http.StatusBadRequest, "Income Statement Credit Harus Di isi", "Income Statement Credit Harus Di isi")
				// }
				// if *JpmDetail.IncomeStatementDr != float64(0) {
				// 	return response.CustomErrorBuilder(http.StatusBadRequest, "Tidak Dapat Masukan Income Statement Debit", "Tidak Dapat Masukan Income Statement Debit")
				// }

				tbID, err := s.Repository.FindByTbd(ctx, &JpmID.ConsolidationID, &v.CoaCode)
				if err != nil {
					return response.CustomErrorBuilder(http.StatusBadRequest, "Tidak Dapat Menemukan Coa di Kombinasi Konsolidasi "+v.CoaCode, "Tidak Dapat Menemukan Coa "+v.CoaCode)
				}
				tbID.Context = ctx

				AmountJpmCr := *tbID.ConsolidationDetailEntity.AmountJpmCr - *v.IncomeStatementCr
				AmountAjeDr := *tbID.ConsolidationDetailEntity.AmountJpmDr - *v.IncomeStatementDr
				tbID.ConsolidationDetailEntity.AmountJpmCr = &AmountJpmCr
				tbID.ConsolidationDetailEntity.AmountJpmDr = &AmountAjeDr
				AmountAfterJpm := *tbID.AmountBeforeJpm - *tbID.AmountJpmDr + *tbID.AmountJpmCr
				tbID.AmountAfterJpm = &AmountAfterJpm
				AmountAfterJcte := *tbID.AmountAfterJpm - *tbID.AmountJcteDr + *tbID.AmountJcteCr
				tbID.AmountAfterJcte = &AmountAfterJcte

				findConsoleBridge, err := s.Repository.FindByConsoleBridge(ctx, &JpmID.ConsolidationID)
				if err != nil {
					return response.CustomErrorBuilder(http.StatusBadRequest, "Tidak Dapat Menemukan Anak Usaha ", "Tidak Dapat Menemukan Anak Usaha ")
				}

				var dataAmountAnakUsahaa []float64
				for _, cb := range *findConsoleBridge {
					findConsoleBridgeDetail, err := s.Repository.FindByConsoleBridgeDetail(ctx, &cb.ID, &v.CoaCode)
					_, err = s.Repository.UpdateTbd(ctx, &tbID.ID, tbID)
					if err != nil {
						return response.CustomErrorBuilder(http.StatusBadRequest, "Tidak Dapat Menemukan Amount Anak Usaha ", "Tidak Dapat Menemukan Amount Anak Usaha ")
					}
					dataAmountAnakUsahaa = append(dataAmountAnakUsahaa, findConsoleBridgeDetail.Amount)
				}
				var sumAmountAnakUsaha float64 = 0.0

				for i := 0; i < len(dataAmountAnakUsahaa); i++ {
					sumAmountAnakUsaha += dataAmountAnakUsahaa[i]
				}

				var totalAmount float64 = *tbID.AmountAfterJcte + sumAmountAnakUsaha
				tbID.AmountCombineSubsidiary = &totalAmount

				var amountConsole float64 = *tbID.AmountCombineSubsidiary - *tbID.AmountJelimDr + *tbID.AmountJelimCr

				tbID.AmountConsole = &amountConsole

				_, err = s.Repository.UpdateTbd(ctx, &tbID.ID, tbID)
				if err != nil {
					return err
				}
				arrJpmDetails = append(arrJpmDetails, v)
			}
			if findCoaGroup.Name == "Dampak penyesuaian proforma  atas OCI Entitas anak" {
				// if *JpmDetail.BalanceSheetDr != float64(0) {
				// 	return response.CustomErrorBuilder(http.StatusBadRequest, "Tidak Dapat Masukan Balance Sheet Debit", "Tidak Dapat Masukan Balance Sheet Debit")
				// }
				// if *JpmDetail.BalanceSheetCr != float64(0) {
				// 	return response.CustomErrorBuilder(http.StatusBadRequest, "Tidak Dapat Masukan Balance Sheet Credit", "Tidak Dapat Masukan Balance Sheet Credit")
				// }
				// if *JpmDetail.IncomeStatementCr == float64(0) {
				// 	return response.CustomErrorBuilder(http.StatusBadRequest, "Income Statement Credit Harus Di isi", "Income Statement Credit Harus Di isi")
				// }
				// if *JpmDetail.IncomeStatementDr != float64(0) {
				// 	return response.CustomErrorBuilder(http.StatusBadRequest, "Tidak Dapat Masukan Income Statement Debit", "Tidak Dapat Masukan Income Statement Debit")
				// }

				tbID, err := s.Repository.FindByTbd(ctx, &JpmID.ConsolidationID, &v.CoaCode)
				if err != nil {
					return response.CustomErrorBuilder(http.StatusBadRequest, "Tidak Dapat Menemukan Coa di Kombinasi Konsolidasi "+v.CoaCode, "Tidak Dapat Menemukan Coa "+v.CoaCode)
				}
				tbID.Context = ctx
				AmountJpmCr := *tbID.ConsolidationDetailEntity.AmountJpmCr - *v.IncomeStatementCr
				AmountAjeDr := *tbID.ConsolidationDetailEntity.AmountJpmDr - *v.IncomeStatementDr
				tbID.ConsolidationDetailEntity.AmountJpmCr = &AmountJpmCr
				tbID.ConsolidationDetailEntity.AmountJpmDr = &AmountAjeDr
				AmountAfterJpm := *tbID.AmountBeforeJpm - *tbID.AmountJpmDr + *tbID.AmountJpmCr
				tbID.AmountAfterJpm = &AmountAfterJpm
				AmountAfterJcte := *tbID.AmountAfterJpm - *tbID.AmountJcteDr + *tbID.AmountJcteCr
				tbID.AmountAfterJcte = &AmountAfterJcte

				findConsoleBridge, err := s.Repository.FindByConsoleBridge(ctx, &JpmID.ConsolidationID)
				if err != nil {
					return response.CustomErrorBuilder(http.StatusBadRequest, "Tidak Dapat Menemukan Anak Usaha ", "Tidak Dapat Menemukan Anak Usaha ")
				}

				var dataAmountAnakUsahaa []float64
				for _, cb := range *findConsoleBridge {
					findConsoleBridgeDetail, err := s.Repository.FindByConsoleBridgeDetail(ctx, &cb.ID, &v.CoaCode)
					_, err = s.Repository.UpdateTbd(ctx, &tbID.ID, tbID)
					if err != nil {
						return response.CustomErrorBuilder(http.StatusBadRequest, "Tidak Dapat Menemukan Amount Anak Usaha ", "Tidak Dapat Menemukan Amount Anak Usaha ")
					}
					dataAmountAnakUsahaa = append(dataAmountAnakUsahaa, findConsoleBridgeDetail.Amount)
				}
				var sumAmountAnakUsaha float64 = 0.0

				for i := 0; i < len(dataAmountAnakUsahaa); i++ {
					sumAmountAnakUsaha += dataAmountAnakUsahaa[i]
				}

				var totalAmount float64 = *tbID.AmountAfterJcte + sumAmountAnakUsaha
				tbID.AmountCombineSubsidiary = &totalAmount

				var amountConsole float64 = *tbID.AmountCombineSubsidiary - *tbID.AmountJelimDr + *tbID.AmountJelimCr

				tbID.AmountConsole = &amountConsole

				_, err = s.Repository.UpdateTbd(ctx, &tbID.ID, tbID)
				if err != nil {
					return err
				}
				arrJpmDetails = append(arrJpmDetails, v)
			}
			if findCoaGroup.Name == "Non Controlling OCI" {
				// if *JpmDetail.BalanceSheetDr != float64(0) {
				// 	return response.CustomErrorBuilder(http.StatusBadRequest, "Tidak Dapat Masukan Balance Sheet Debit", "Tidak Dapat Masukan Balance Sheet Debit")
				// }
				// if *JpmDetail.BalanceSheetCr != float64(0) {
				// 	return response.CustomErrorBuilder(http.StatusBadRequest, "Tidak Dapat Masukan Balance Sheet Credit", "Tidak Dapat Masukan Balance Sheet Credit")
				// }
				// if *JpmDetail.IncomeStatementCr == float64(0) {
				// 	return response.CustomErrorBuilder(http.StatusBadRequest, "Income Statement Credit Harus Di isi", "Income Statement Credit Harus Di isi")
				// }
				// if *JpmDetail.IncomeStatementDr != float64(0) {
				// 	return response.CustomErrorBuilder(http.StatusBadRequest, "Tidak Dapat Masukan Income Statement Debit", "Tidak Dapat Masukan Income Statement Debit")
				// }

				tbID, err := s.Repository.FindByTbd(ctx, &JpmID.ConsolidationID, &v.CoaCode)
				if err != nil {
					return response.CustomErrorBuilder(http.StatusBadRequest, "Tidak Dapat Menemukan Coa di Kombinasi Konsolidasi "+v.CoaCode, "Tidak Dapat Menemukan Coa "+v.CoaCode)
				}
				tbID.Context = ctx

				AmountJpmCr := *tbID.ConsolidationDetailEntity.AmountJpmCr - *v.IncomeStatementCr
				AmountAjeDr := *tbID.ConsolidationDetailEntity.AmountJpmDr - *v.IncomeStatementDr
				tbID.ConsolidationDetailEntity.AmountJpmCr = &AmountJpmCr
				tbID.ConsolidationDetailEntity.AmountJpmDr = &AmountAjeDr
				AmountAfterJpm := *tbID.AmountBeforeJpm - *tbID.AmountJpmDr + *tbID.AmountJpmCr
				tbID.AmountAfterJpm = &AmountAfterJpm
				AmountAfterJcte := *tbID.AmountAfterJpm - *tbID.AmountJcteDr + *tbID.AmountJcteCr
				tbID.AmountAfterJcte = &AmountAfterJcte

				findConsoleBridge, err := s.Repository.FindByConsoleBridge(ctx, &JpmID.ConsolidationID)
				if err != nil {
					return response.CustomErrorBuilder(http.StatusBadRequest, "Tidak Dapat Menemukan Anak Usaha ", "Tidak Dapat Menemukan Anak Usaha ")
				}

				var dataAmountAnakUsahaa []float64
				for _, cb := range *findConsoleBridge {
					findConsoleBridgeDetail, err := s.Repository.FindByConsoleBridgeDetail(ctx, &cb.ID, &v.CoaCode)
					_, err = s.Repository.UpdateTbd(ctx, &tbID.ID, tbID)
					if err != nil {
						return response.CustomErrorBuilder(http.StatusBadRequest, "Tidak Dapat Menemukan Amount Anak Usaha ", "Tidak Dapat Menemukan Amount Anak Usaha ")
					}
					dataAmountAnakUsahaa = append(dataAmountAnakUsahaa, findConsoleBridgeDetail.Amount)
				}
				var sumAmountAnakUsaha float64 = 0.0

				for i := 0; i < len(dataAmountAnakUsahaa); i++ {
					sumAmountAnakUsaha += dataAmountAnakUsahaa[i]
				}

				var totalAmount float64 = *tbID.AmountAfterJcte + sumAmountAnakUsaha
				tbID.AmountCombineSubsidiary = &totalAmount

				var amountConsole float64 = *tbID.AmountCombineSubsidiary - *tbID.AmountJelimDr + *tbID.AmountJelimCr

				tbID.AmountConsole = &amountConsole

				_, err = s.Repository.UpdateTbd(ctx, &tbID.ID, tbID)
				if err != nil {
					return err
				}
				arrJpmDetails = append(arrJpmDetails, v)
			}

		}
		data.Context = ctx
		result, err := s.Repository.Delete(ctx, &payload.ID, &data)
		if err != nil {
			return response.ErrorBuilder(&response.ErrorConstant.InternalServerError, err)
		}
		data = *result
		formatterDetailSumData, err := s.Repository.FindSummary(ctx)
		if err != nil {
			return response.CustomErrorBuilder(http.StatusBadRequest, "Gagal Dapatkan sumary Coa", "Gagal Dapatkan sumary Coa")
		}
		for _, v := range *formatterDetailSumData {
			criteriaTBDetail := model.ConsolidationDetailFilterModel{}
			criteriaTBDetail.ConsolidationID = &JpmID.ConsolidationID
			if v.AutoSummary != nil && *v.AutoSummary {
				code := fmt.Sprintf("%s_Subtotal", v.Code)
				criteriaTBDetail.Code = &code
				criteriaTBDetail.ConsolidationID = &JpmID.ConsolidationID
				mfadetailsum, _, err := s.Repository.FindC(ctx, &criteriaTBDetail, &abstraction.Pagination{})
				if err != nil {
					return response.CustomErrorBuilder(http.StatusBadRequest, "Gagal Find Detail Consolidation", "Gagal Find Detail Consolidation")
				}

				for _, a := range *mfadetailsum {
					sumTBD, err := s.Repository.FindSummarys(ctx, &v.Code, &JpmID.ConsolidationID, v.IsCoa)
					if err != nil {
						return response.CustomErrorBuilder(http.StatusBadRequest, "Gagal Menjumlahkan Summary Coa", "Gagal Menjumlahkan Summary Coa")
					}

					updateSummary := model.ConsolidationDetailEntityModel{
						ConsolidationDetailEntity: model.ConsolidationDetailEntity{
							AmountJpmCr:             sumTBD.AmountJpmCr,
							AmountJpmDr:             sumTBD.AmountJpmDr,
							AmountAfterJpm:          sumTBD.AmountAfterJpm,
							AmountJcteCr:            sumTBD.AmountJcteCr,
							AmountJcteDr:            sumTBD.AmountJcteDr,
							AmountAfterJcte:         sumTBD.AmountAfterJcte,
							AmountCombineSubsidiary: sumTBD.AmountCombineSubsidiary,
							AmountJelimCr:           sumTBD.AmountJelimCr,
							AmountJelimDr:           sumTBD.AmountJelimDr,
							AmountConsole:           sumTBD.AmountConsole,
						},
					}

					_, err = s.Repository.Updates(ctx, &a.ID, &updateSummary)
					if err != nil {
						return response.CustomErrorBuilder(http.StatusBadRequest, "Gagal Update Consolidation Detail", "Gagal Update Consolidation Detail")
					}
				}
			}
			if v.IsTotal != nil && *v.IsTotal == true && v.FxSummary != "" {

				tmpString := []string{"AmountBeforeJpm", "AmountJpmCr", "AmountJpmDr", "AmountAfterJpm", "AmountJcteCr", "AmountJcteDr", "AmountAfterJcte", "AmountCombineSubsidiary", "AmountJelimCr", "AmountJelimDr", "AmountConsole"}
				tmpTotalFl := make(map[string]*float64)
				// reg := regexp.MustCompile(`[0-9]+\d{3,}`)
				reg := regexp.MustCompile(`[A-Za-z0-9_~:()]+|[0-9]+\d{3,}`)

				for _, tipe := range tmpString {
					formula := strings.ToUpper(v.FxSummary)
					if tipe == "AmountJpmCr" || tipe == "AmountJpmDr" || tipe == "AmountJcteCr"|| tipe == "AmountJcteDr"|| tipe == "AmountJelimDr"|| tipe == "AmountJelimCr"{
						newStr := strings.ReplaceAll(v.FxSummary, "-", "+")
						formula = strings.ToUpper(newStr)
					}
					match := reg.FindAllString(formula, -1)
					amountBeforeJpm := make(map[string]interface{}, 0)
					for _, vMatch := range match {
						//cari jml berdasarkan code

						if len(vMatch) < 3 {
							continue
						}
						sumTBD, err := s.Repository.FindSummarys(ctx, &vMatch, &JpmID.ConsolidationID, v.IsCoa)
						if err != nil {
							return response.CustomErrorBuilder(http.StatusBadRequest, "Gagal Mendapatkan Jumlah", "Gagal Mendapatkan Jumlah")
						}
						angka := 0.0

						if tipe == "AmountBeforeJpm" && sumTBD.AmountBeforeJpm != nil {
							angka = *sumTBD.AmountBeforeJpm
						} else if tipe == "AmountJpmCr" && sumTBD.AmountJpmCr != nil {
							angka = *sumTBD.AmountJpmCr
						} else if tipe == "AmountJpmDr" && sumTBD.AmountJpmDr != nil {
							angka = *sumTBD.AmountJpmDr
						} else if tipe == "AmountAfterJpm" && sumTBD.AmountAfterJpm != nil {
							angka = *sumTBD.AmountAfterJpm
						} else if tipe == "AmountJcteCr" && sumTBD.AmountJcteCr != nil {
							angka = *sumTBD.AmountJcteCr
						} else if tipe == "AmountJcteDr" && sumTBD.AmountJcteDr != nil {
							angka = *sumTBD.AmountJcteDr
						} else if tipe == "AmountAfterJcte" && sumTBD.AmountAfterJcte != nil {
							angka = *sumTBD.AmountAfterJcte
						} else if tipe == "AmountCombineSubsidiary" && sumTBD.AmountCombineSubsidiary != nil {
							angka = *sumTBD.AmountCombineSubsidiary
						} else if tipe == "AmountJelimCr" && sumTBD.AmountJelimCr != nil {
							angka = *sumTBD.AmountJelimCr
						} else if tipe == "AmountJelimDr" && sumTBD.AmountJelimDr != nil {
							angka = *sumTBD.AmountJelimDr
						} else if tipe == "AmountConsole" && sumTBD.AmountConsole != nil {
							angka = *sumTBD.AmountConsole
						}

						formula = helper.ReplaceWholeWord(formula, vMatch, fmt.Sprintf("(%2.f)", angka))
						// parameters[vMatch] = angka

					}

					expressionFormula, err := govaluate.NewEvaluableExpression(formula)
					if err != nil {
						fmt.Println(err)
						return err
					}
					JpmID, err := expressionFormula.Evaluate(amountBeforeJpm)
					if err != nil {
						return response.CustomErrorBuilder(http.StatusBadRequest, "Tidak Dapat Masukan Balance Sheet Debit", "Tidak Dapat Masukan Balance Sheet Debit")
					}
					tmp := JpmID.(float64)
					tmpTotalFl[tipe] = &tmp

				}
				criteriaTBDetail.ConsolidationID = &JpmID.ConsolidationID
				criteriaTBDetail.Code = &v.Code
				mfadetailsum, err := s.Repository.FindDetailConsole(ctx, &criteriaTBDetail)
				if err != nil {
					fmt.Sprintln(err)
					return err
				}

				updateSummary := model.ConsolidationDetailEntityModel{
					ConsolidationDetailEntity: model.ConsolidationDetailEntity{
						AmountJpmCr:             tmpTotalFl["AmountJpmCr"],
						AmountJpmDr:             tmpTotalFl["AmountJpmDr"],
						AmountAfterJpm:          tmpTotalFl["AmountAfterJpm"],
						AmountJcteCr:            tmpTotalFl["AmountJcteCr"],
						AmountJcteDr:            tmpTotalFl["AmountJcteDr"],
						AmountAfterJcte:         tmpTotalFl["AmountAfterJcte"],
						AmountCombineSubsidiary: tmpTotalFl["AmountCombineSubsidiary"],
						AmountJelimCr:           tmpTotalFl["AmountJelimCr"],
						AmountJelimDr:           tmpTotalFl["AmountJelimDr"],
						AmountConsole:           tmpTotalFl["AmountConsole"],
					},
				}

				_, err = s.Repository.Updates(ctx, &mfadetailsum.ID, &updateSummary)
				if err != nil {
					return response.CustomErrorBuilder(http.StatusBadRequest, "Tidak Dapat Masukan Balance Sheet Debit", "Tidak Dapat Masukan Balance Sheet Debit")
				}

			}
			if v.Code == "LABA_BERSIH" {
				code := "310401004"
				criteriaTBDetail := model.ConsolidationDetailFilterModel{}
				criteriaTBDetail.ConsolidationID = &JpmID.ConsolidationID
				criteriaTBDetail.Code = &code
				customRowOne, err := s.Repository.FindByTbd(ctx, &JpmID.ConsolidationID, &code)
				if err != nil {
					return helper.ErrorHandler(err)
				}

				code = "310402002"
				criteriaTBDetail.Code = &code
				customRowTwo, err := s.Repository.FindByTbd(ctx, &JpmID.ConsolidationID, &code)
				if err != nil {
					return helper.ErrorHandler(err)
				}

				code = "310501002"
				criteriaTBDetail.Code = &code
				customRowThree, err := s.Repository.FindByTbd(ctx, &JpmID.ConsolidationID, &code)
				if err != nil {
					return helper.ErrorHandler(err)
				}

				code = "310502002"
				criteriaTBDetail.Code = &code
				customRowFour, err := s.Repository.FindByTbd(ctx, &JpmID.ConsolidationID, &code)
				if err != nil {
					return helper.ErrorHandler(err)
				}

				code = "310503002"
				criteriaTBDetail.Code = &code
				customRowFive, err := s.Repository.FindByTbd(ctx, &JpmID.ConsolidationID, &code)
				if err != nil {
					return helper.ErrorHandler(err)
				}

				code = "950101001" //REVALUATION FA (4337)
				criteriaTBDetail.Code = &code
				dataReFa, err := s.Repository.FindByTbd(ctx, &JpmID.ConsolidationID, &code)
				if err != nil {
					return helper.ErrorHandler(err)
				}

				code = "950301001" //Financial Instrument (4342)
				criteriaTBDetail.Code = &code
				dataFinIn, err := s.Repository.FindByTbd(ctx, &JpmID.ConsolidationID, &code)
				if err != nil {
					return helper.ErrorHandler(err)
				}

				code = "950301002" //Income tax relating to components of OCI (4343)
				criteriaTBDetail.Code = &code
				dataIncomeTax, err := s.Repository.FindByTbd(ctx, &JpmID.ConsolidationID, &code)
				if err != nil {
					return helper.ErrorHandler(err)
				}

				code = "950401001" //Foreign Exchange (4345)
				criteriaTBDetail.Code = &code
				dataForex, err := s.Repository.FindByTbd(ctx, &JpmID.ConsolidationID, &code)
				if err != nil {
					return helper.ErrorHandler(err)
				}

				code = "950401002" //Income tax relating to components of OCI (4346)
				criteriaTBDetail.Code = &code
				dataIncomeTax2, err := s.Repository.FindByTbd(ctx, &JpmID.ConsolidationID, &code)
				if err != nil {
					return helper.ErrorHandler(err)
				}

				code = "LABA_BERSIH"
				criteriaTBDetail.Code = &code
				dataLaba, err := s.Repository.FindByTbd(ctx, &JpmID.ConsolidationID, &code)
				if err != nil {
					return helper.ErrorHandler(err)
				}

				code = "TOTAL_PENGHASILAN_KOMPREHENSIF_LAIN~BS"
				criteriaTBDetail.Code = &code
				dataKomprehensif, err := s.Repository.FindByTbd(ctx, &JpmID.ConsolidationID, &code)
				if err != nil {
					return helper.ErrorHandler(err)
				}

				summaryCodes, err := s.Repository.SummaryByCodes(ctx, &JpmID.ConsolidationID, []string{"310501002", "310502002", "310503002"})
				if err != nil {
					return helper.ErrorHandler(err)
				}

				updatedCustomRowOne := model.ConsolidationDetailEntityModel{}
				updatedCustomRowOne.Context = ctx
				updatedCustomRowOne.AmountBeforeJpm = dataLaba.AmountBeforeJpm
				updatedCustomRowOne.AmountJpmCr = dataLaba.AmountJpmCr
				updatedCustomRowOne.AmountJpmDr = dataLaba.AmountJpmDr
				updatedCustomRowOne.AmountAfterJpm = dataLaba.AmountAfterJpm

				updatedCustomRowOne.AmountJcteCr = dataLaba.AmountJcteCr
				updatedCustomRowOne.AmountJcteDr = dataLaba.AmountJcteDr
				updatedCustomRowOne.AmountAfterJcte = dataLaba.AmountAfterJcte

				updatedCustomRowOne.AmountCombineSubsidiary = dataLaba.AmountCombineSubsidiary

				updatedCustomRowOne.AmountJelimCr = dataLaba.AmountJelimCr
				updatedCustomRowOne.AmountJelimDr = dataLaba.AmountJelimDr
				updatedCustomRowOne.AmountConsole = dataLaba.AmountConsole

				_, err = s.Repository.Updates(ctx, &customRowOne.ID, &updatedCustomRowOne)
				if err != nil {
					return helper.ErrorHandler(err)
				}

				//

				updatedCustomRowThree := model.ConsolidationDetailEntityModel{}
				updatedCustomRowThree.Context = ctx
				updatedCustomRowThree.AmountBeforeJpm = dataReFa.AmountBeforeJpm
				updatedCustomRowThree.AmountJpmCr = dataReFa.AmountJpmCr
				updatedCustomRowThree.AmountJpmDr = dataReFa.AmountJpmDr
				updatedCustomRowThree.AmountAfterJpm = dataReFa.AmountAfterJpm

				updatedCustomRowThree.AmountJcteCr = dataReFa.AmountJcteCr
				updatedCustomRowThree.AmountJcteDr = dataReFa.AmountJcteDr
				updatedCustomRowThree.AmountAfterJcte = dataReFa.AmountAfterJcte

				updatedCustomRowThree.AmountCombineSubsidiary = dataReFa.AmountCombineSubsidiary

				updatedCustomRowThree.AmountJelimCr = dataReFa.AmountJelimCr
				updatedCustomRowThree.AmountJelimDr = dataReFa.AmountJelimDr
				updatedCustomRowThree.AmountConsole = dataReFa.AmountConsole

				_, err = s.Repository.Updates(ctx, &customRowThree.ID, &updatedCustomRowThree)
				if err != nil {
					return helper.ErrorHandler(err)
				}

				updateCustomRowFour := model.ConsolidationDetailEntityModel{}
				updateCustomRowFour.Context = ctx

				tmp1 := 0.0
				if dataFinIn.AmountBeforeJpm != nil && *dataFinIn.AmountBeforeJpm != 0 {
					tmp1 = *dataFinIn.AmountBeforeJpm
				}
				if dataIncomeTax.AmountBeforeJpm != nil && *dataIncomeTax.AmountBeforeJpm != 0 {
					tmp1 = tmp1 + *dataIncomeTax.AmountBeforeJpm
				}
				updateCustomRowFour.AmountBeforeJpm = &tmp1

				tmp2 := 0.0
				if dataFinIn.AmountJpmCr != nil && *dataFinIn.AmountJpmCr != 0 {
					tmp2 = *dataFinIn.AmountJpmCr
				}
				if dataIncomeTax.AmountJpmCr != nil && *dataIncomeTax.AmountJpmCr != 0 {
					tmp2 = tmp2 + *dataIncomeTax.AmountJpmCr
				}
				updateCustomRowFour.AmountJpmCr = &tmp2

				tmp3 := 0.0
				if dataFinIn.AmountJpmDr != nil && *dataFinIn.AmountJpmDr != 0 {
					tmp3 = *dataFinIn.AmountJpmDr
				}
				if dataIncomeTax.AmountJpmDr != nil && *dataIncomeTax.AmountJpmDr != 0 {
					tmp3 = tmp3 + *dataIncomeTax.AmountJpmDr
				}
				updateCustomRowFour.AmountJpmDr = &tmp3

				tmp4 := 0.0
				if dataFinIn.AmountAfterJpm != nil && *dataFinIn.AmountAfterJpm != 0 {
					tmp4 = *dataFinIn.AmountAfterJpm
				}
				if dataIncomeTax.AmountAfterJpm != nil && *dataIncomeTax.AmountAfterJpm != 0 {
					tmp4 = tmp4 + *dataIncomeTax.AmountAfterJpm
				}
				updateCustomRowFour.AmountAfterJpm = &tmp4

				//jcte

				tmp5 := 0.0
				if dataFinIn.AmountJcteCr != nil && *dataFinIn.AmountJcteCr != 0 {
					tmp5 = *dataFinIn.AmountJcteCr
				}
				if dataIncomeTax.AmountJcteCr != nil && *dataIncomeTax.AmountJcteCr != 0 {
					tmp5 = tmp5 + *dataIncomeTax.AmountJcteCr
				}
				updateCustomRowFour.AmountJcteCr = &tmp5

				tmp6 := 0.0
				if dataFinIn.AmountJcteDr != nil && *dataFinIn.AmountJcteDr != 0 {
					tmp6 = *dataFinIn.AmountJcteDr
				}
				if dataIncomeTax.AmountJcteDr != nil && *dataIncomeTax.AmountJcteDr != 0 {
					tmp6 = tmp6 + *dataIncomeTax.AmountJcteDr
				}
				updateCustomRowFour.AmountJcteDr = &tmp6

				tmp7 := 0.0
				if dataFinIn.AmountAfterJcte != nil && *dataFinIn.AmountAfterJcte != 0 {
					tmp7 = *dataFinIn.AmountAfterJcte
				}
				if dataIncomeTax.AmountAfterJcte != nil && *dataIncomeTax.AmountAfterJcte != 0 {
					tmp7 = tmp7 + *dataIncomeTax.AmountAfterJcte
				}
				updateCustomRowFour.AmountAfterJcte = &tmp7

				// acs
				tmp8 := 0.0
				if dataFinIn.AmountCombineSubsidiary != nil && *dataFinIn.AmountCombineSubsidiary != 0 {
					tmp8 = *dataFinIn.AmountCombineSubsidiary
				}
				if dataIncomeTax.AmountCombineSubsidiary != nil && *dataIncomeTax.AmountCombineSubsidiary != 0 {
					tmp8 = tmp8 + *dataIncomeTax.AmountCombineSubsidiary
				}
				updateCustomRowFour.AmountCombineSubsidiary = &tmp8

				//jelim
				tmp9 := 0.0
				if dataFinIn.AmountJelimCr != nil && *dataFinIn.AmountJelimCr != 0 {
					tmp9 = *dataFinIn.AmountJelimCr
				}
				if dataIncomeTax.AmountJelimCr != nil && *dataIncomeTax.AmountJelimCr != 0 {
					tmp9 = tmp9 + *dataIncomeTax.AmountJelimCr
				}
				updateCustomRowFour.AmountJelimCr = &tmp9

				tmp10 := 0.0
				if dataFinIn.AmountJelimDr != nil && *dataFinIn.AmountJelimDr != 0 {
					tmp10 = *dataFinIn.AmountJelimDr
				}
				if dataIncomeTax.AmountJelimDr != nil && *dataIncomeTax.AmountJelimDr != 0 {
					tmp10 = tmp10 + *dataIncomeTax.AmountJelimDr
				}
				updateCustomRowFour.AmountJelimDr = &tmp10

				tmp11 := 0.0
				if dataFinIn.AmountConsole != nil && *dataFinIn.AmountConsole != 0 {
					tmp11 = *dataFinIn.AmountConsole
				}
				if dataIncomeTax.AmountConsole != nil && *dataIncomeTax.AmountConsole != 0 {
					tmp11 = tmp11 + *dataIncomeTax.AmountConsole
				}
				updateCustomRowFour.AmountConsole = &tmp11

				_, err = s.Repository.Updates(ctx, &customRowFour.ID, &updateCustomRowFour)
				if err != nil {
					return helper.ErrorHandler(err)
				}

				updateCustomRowFive := model.ConsolidationDetailEntityModel{}
				updateCustomRowFive.Context = ctx

				tmpa := 0.0
				if dataForex.AmountBeforeJpm != nil && *dataForex.AmountBeforeJpm != 0 {
					tmpa = *dataForex.AmountBeforeJpm
				}
				if dataIncomeTax2.AmountBeforeJpm != nil && *dataIncomeTax2.AmountBeforeJpm != 0 {
					tmpa = tmpa + *dataIncomeTax2.AmountBeforeJpm
				}
				updateCustomRowFive.AmountBeforeJpm = &tmpa

				tmpb := 0.0
				if dataForex.AmountJpmCr != nil && *dataForex.AmountJpmCr != 0 {
					tmpb = *dataForex.AmountJpmCr
				}
				if dataIncomeTax2.AmountJpmCr != nil && *dataIncomeTax2.AmountJpmCr != 0 {
					tmpb = tmpb + *dataIncomeTax2.AmountJpmCr
				}
				updateCustomRowFive.AmountJpmCr = &tmpb

				tmpc := 0.0
				if dataForex.AmountJpmDr != nil && *dataForex.AmountJpmDr != 0 {
					tmpc = *dataForex.AmountJpmDr
				}
				if dataIncomeTax2.AmountJpmDr != nil && *dataIncomeTax2.AmountJpmDr != 0 {
					tmpc = tmpc + *dataIncomeTax2.AmountJpmDr
				}
				updateCustomRowFive.AmountJpmDr = &tmpc

				tmpd := 0.0
				if dataForex.AmountAfterJpm != nil && *dataForex.AmountAfterJpm != 0 {
					tmpd = *dataForex.AmountAfterJpm
				}
				if dataIncomeTax2.AmountAfterJpm != nil && *dataIncomeTax2.AmountAfterJpm != 0 {
					tmpd = tmpd + *dataIncomeTax2.AmountAfterJpm
				}
				updateCustomRowFive.AmountAfterJpm = &tmpd

				//jcte
				tmpf := 0.0
				if dataForex.AmountJcteCr != nil && *dataForex.AmountJcteCr != 0 {
					tmpf = *dataForex.AmountJcteCr
				}
				if dataIncomeTax2.AmountJcteCr != nil && *dataIncomeTax2.AmountJcteCr != 0 {
					tmpf = tmpf + *dataIncomeTax2.AmountJcteCr
				}
				updateCustomRowFive.AmountJcteCr = &tmpf

				tmpe := 0.0
				if dataForex.AmountJcteDr != nil && *dataForex.AmountJcteDr != 0 {
					tmpe = *dataForex.AmountJcteDr
				}
				if dataIncomeTax2.AmountJcteDr != nil && *dataIncomeTax2.AmountJcteDr != 0 {
					tmpe = tmpe + *dataIncomeTax2.AmountJcteDr
				}
				updateCustomRowFive.AmountJcteDr = &tmpe

				tmph := 0.0
				if dataForex.AmountAfterJcte != nil && *dataForex.AmountAfterJcte != 0 {
					tmph = *dataForex.AmountAfterJcte
				}
				if dataIncomeTax2.AmountAfterJcte != nil && *dataIncomeTax2.AmountAfterJcte != 0 {
					tmph = tmph + *dataIncomeTax2.AmountAfterJcte
				}
				updateCustomRowFive.AmountAfterJcte = &tmph

				// acs
				tmpg := 0.0
				if dataForex.AmountCombineSubsidiary != nil && *dataForex.AmountCombineSubsidiary != 0 {
					tmpg = *dataForex.AmountCombineSubsidiary
				}
				if dataIncomeTax2.AmountCombineSubsidiary != nil && *dataIncomeTax2.AmountCombineSubsidiary != 0 {
					tmpg = tmpg + *dataIncomeTax2.AmountCombineSubsidiary
				}
				updateCustomRowFive.AmountCombineSubsidiary = &tmpg

				//jelim
				tmpi := 0.0
				if dataForex.AmountJelimCr != nil && *dataForex.AmountJelimCr != 0 {
					tmpi = *dataForex.AmountJelimCr
				}
				if dataIncomeTax2.AmountJelimCr != nil && *dataIncomeTax2.AmountJelimCr != 0 {
					tmpi = tmpi + *dataIncomeTax2.AmountJelimCr
				}
				updateCustomRowFive.AmountJelimCr = &tmpi

				tmpj := 0.0
				if dataForex.AmountJelimDr != nil && *dataForex.AmountJelimDr != 0 {
					tmpj = *dataForex.AmountJelimDr
				}
				if dataIncomeTax2.AmountJelimDr != nil && *dataIncomeTax2.AmountJelimDr != 0 {
					tmpj = tmpj + *dataIncomeTax2.AmountJelimDr
				}
				updateCustomRowFive.AmountJelimDr = &tmpj

				tmpk := 0.0
				if dataForex.AmountConsole != nil && *dataForex.AmountConsole != 0 {
					tmpk = *dataForex.AmountConsole
				}
				if dataIncomeTax2.AmountConsole != nil && *dataIncomeTax2.AmountConsole != 0 {
					tmpk = tmpk + *dataIncomeTax2.AmountConsole
				}
				updateCustomRowFive.AmountConsole = &tmpk
				_, err = s.Repository.Updates(ctx, &customRowFive.ID, &updateCustomRowFive)
				if err != nil {
					return helper.ErrorHandler(err)
				}
				updateCustomRowTwo := model.ConsolidationDetailEntityModel{}
				updateCustomRowTwo.Context = ctx

				tmp12 := 0.0
				if dataKomprehensif.AmountBeforeJpm != nil && *dataKomprehensif.AmountBeforeJpm != 0 {
					tmp12 = *dataKomprehensif.AmountBeforeJpm
				}
				if summaryCodes.AmountBeforeJpm != nil && *summaryCodes.AmountBeforeJpm != 0 {
					tmp12 = tmp12 - *summaryCodes.AmountBeforeJpm
				}
				updateCustomRowTwo.AmountBeforeJpm = &tmp12

				tmp13 := 0.0
				if dataKomprehensif.AmountJpmCr != nil && *dataKomprehensif.AmountJpmCr != 0 {
					tmp13 = *dataKomprehensif.AmountJpmCr
				}
				if summaryCodes.AmountJpmCr != nil && *summaryCodes.AmountJpmCr != 0 {
					tmp13 = tmp13 - *summaryCodes.AmountJpmCr
				}
				updateCustomRowTwo.AmountJpmCr = &tmp13

				tmp14 := 0.0
				if dataKomprehensif.AmountJpmDr != nil && *dataKomprehensif.AmountJpmDr != 0 {
					tmp14 = *dataKomprehensif.AmountJpmDr
				}
				if summaryCodes.AmountJpmDr != nil && *summaryCodes.AmountJpmDr != 0 {
					tmp14 = tmp14 - *summaryCodes.AmountJpmDr
				}
				updateCustomRowTwo.AmountJpmDr = &tmp14

				tmp15 := 0.0
				if dataKomprehensif.AmountAfterJpm != nil && *dataKomprehensif.AmountAfterJpm != 0 {
					tmp15 = *dataKomprehensif.AmountAfterJpm
				}
				if summaryCodes.AmountAfterJpm != nil && *summaryCodes.AmountAfterJpm != 0 {
					tmp15 = tmp15 - *summaryCodes.AmountAfterJpm
				}
				updateCustomRowTwo.AmountAfterJpm = &tmp15

				//jcte

				tmp16 := 0.0
				if dataKomprehensif.AmountJcteCr != nil && *dataKomprehensif.AmountJcteCr != 0 {
					tmp16 = *dataKomprehensif.AmountJcteCr
				}
				if summaryCodes.AmountJcteCr != nil && *summaryCodes.AmountJcteCr != 0 {
					tmp16 = tmp16 - *summaryCodes.AmountJcteCr
				}
				updateCustomRowTwo.AmountJcteCr = &tmp16

				tmp17 := 0.0
				if dataKomprehensif.AmountJcteDr != nil && *dataKomprehensif.AmountJcteDr != 0 {
					tmp17 = *dataKomprehensif.AmountJcteDr
				}
				if summaryCodes.AmountJcteDr != nil && *summaryCodes.AmountJcteDr != 0 {
					tmp17 = tmp17 - *summaryCodes.AmountJcteDr
				}
				updateCustomRowTwo.AmountJcteDr = &tmp17

				tmp18 := 0.0
				if dataKomprehensif.AmountAfterJcte != nil && *dataKomprehensif.AmountAfterJcte != 0 {
					tmp18 = *dataKomprehensif.AmountAfterJcte
				}
				if summaryCodes.AmountAfterJcte != nil && *summaryCodes.AmountAfterJcte != 0 {
					tmp18 = tmp18 - *summaryCodes.AmountAfterJcte
				}
				updateCustomRowTwo.AmountAfterJcte = &tmp18

				// acs
				tmp19 := 0.0
				if dataKomprehensif.AmountCombineSubsidiary != nil && *dataKomprehensif.AmountCombineSubsidiary != 0 {
					tmp19 = *dataKomprehensif.AmountCombineSubsidiary
				}
				if summaryCodes.AmountCombineSubsidiary != nil && *summaryCodes.AmountCombineSubsidiary != 0 {
					tmp19 = tmp19 - *summaryCodes.AmountCombineSubsidiary
				}
				updateCustomRowTwo.AmountCombineSubsidiary = &tmp19

				//jelim
				tmp20 := 0.0
				if dataKomprehensif.AmountJelimCr != nil && *dataKomprehensif.AmountJelimCr != 0 {
					tmp20 = *dataKomprehensif.AmountJelimCr
				}
				if summaryCodes.AmountJelimCr != nil && *summaryCodes.AmountJelimCr != 0 {
					tmp20 = tmp20 - *summaryCodes.AmountJelimCr
				}
				updateCustomRowTwo.AmountJelimCr = &tmp20

				tmp21 := 0.0
				if dataKomprehensif.AmountJelimDr != nil && *dataKomprehensif.AmountJelimDr != 0 {
					tmp21 = *dataKomprehensif.AmountJelimDr
				}
				if summaryCodes.AmountJelimDr != nil && *summaryCodes.AmountJelimDr != 0 {
					tmp21 = tmp21 - *summaryCodes.AmountJelimDr
				}
				updateCustomRowTwo.AmountJelimDr = &tmp21

				tmp22 := 0.0
				if dataKomprehensif.AmountConsole != nil && *dataKomprehensif.AmountConsole != 0 {
					tmp22 = *dataKomprehensif.AmountConsole
				}
				if summaryCodes.AmountConsole != nil && *summaryCodes.AmountConsole != 0 {
					tmp22 = tmp22 - *summaryCodes.AmountConsole
				}
				updateCustomRowTwo.AmountConsole = &tmp22
				_, err = s.Repository.Updates(ctx, &customRowTwo.ID, &updateCustomRowTwo)
				if err != nil {
					return helper.ErrorHandler(err)
				}
			}
		}
		{
			{
				code := "310401004"
				criteriaTBDetail := model.ConsolidationDetailFilterModel{}
				criteriaTBDetail.ConsolidationID = &JpmID.ConsolidationID
				criteriaTBDetail.Code = &code
				customRowOne, err := s.Repository.FindByTbd(ctx, &JpmID.ConsolidationID, &code)
				if err != nil {
					return helper.ErrorHandler(err)
				}

				code = "310402002"
				criteriaTBDetail.Code = &code
				customRowTwo, err := s.Repository.FindByTbd(ctx, &JpmID.ConsolidationID, &code)
				if err != nil {
					return helper.ErrorHandler(err)
				}

				code = "310501002"
				criteriaTBDetail.Code = &code
				customRowThree, err := s.Repository.FindByTbd(ctx, &JpmID.ConsolidationID, &code)
				if err != nil {
					return helper.ErrorHandler(err)
				}

				code = "310502002"
				criteriaTBDetail.Code = &code
				customRowFour, err := s.Repository.FindByTbd(ctx, &JpmID.ConsolidationID, &code)
				if err != nil {
					return helper.ErrorHandler(err)
				}

				code = "310503002"
				criteriaTBDetail.Code = &code
				customRowFive, err := s.Repository.FindByTbd(ctx, &JpmID.ConsolidationID, &code)
				if err != nil {
					return helper.ErrorHandler(err)
				}

				code = "950101001" //REVALUATION FA (4337)
				criteriaTBDetail.Code = &code
				dataReFa, err := s.Repository.FindByTbd(ctx, &JpmID.ConsolidationID, &code)
				if err != nil {
					return helper.ErrorHandler(err)
				}

				code = "950301001" //Financial Instrument (4342)
				criteriaTBDetail.Code = &code
				dataFinIn, err := s.Repository.FindByTbd(ctx, &JpmID.ConsolidationID, &code)
				if err != nil {
					return helper.ErrorHandler(err)
				}

				code = "950301002" //Income tax relating to components of OCI (4343)
				criteriaTBDetail.Code = &code
				dataIncomeTax, err := s.Repository.FindByTbd(ctx, &JpmID.ConsolidationID, &code)
				if err != nil {
					return helper.ErrorHandler(err)
				}

				code = "950401001" //Foreign Exchange (4345)
				criteriaTBDetail.Code = &code
				dataForex, err := s.Repository.FindByTbd(ctx, &JpmID.ConsolidationID, &code)
				if err != nil {
					return helper.ErrorHandler(err)
				}

				code = "950401002" //Income tax relating to components of OCI (4346)
				criteriaTBDetail.Code = &code
				dataIncomeTax2, err := s.Repository.FindByTbd(ctx, &JpmID.ConsolidationID, &code)
				if err != nil {
					return helper.ErrorHandler(err)
				}

				code = "LABA_BERSIH"
				criteriaTBDetail.Code = &code
				dataLaba, err := s.Repository.FindByTbd(ctx, &JpmID.ConsolidationID, &code)
				if err != nil {
					return helper.ErrorHandler(err)
				}

				code = "TOTAL_PENGHASILAN_KOMPREHENSIF_LAIN~BS"
				criteriaTBDetail.Code = &code
				dataKomprehensif, err := s.Repository.FindByTbd(ctx, &JpmID.ConsolidationID, &code)
				if err != nil {
					return helper.ErrorHandler(err)
				}

				summaryCodes, err := s.Repository.SummaryByCodes(ctx, &JpmID.ConsolidationID, []string{"310501002", "310502002", "310503002"})
				if err != nil {
					return helper.ErrorHandler(err)
				}

				updatedCustomRowOne := model.ConsolidationDetailEntityModel{}
				updatedCustomRowOne.Context = ctx
				updatedCustomRowOne.AmountBeforeJpm = dataLaba.AmountBeforeJpm
				updatedCustomRowOne.AmountJpmCr = dataLaba.AmountJpmCr
				updatedCustomRowOne.AmountJpmDr = dataLaba.AmountJpmDr
				updatedCustomRowOne.AmountAfterJpm = dataLaba.AmountAfterJpm

				updatedCustomRowOne.AmountJcteCr = dataLaba.AmountJcteCr
				updatedCustomRowOne.AmountJcteDr = dataLaba.AmountJcteDr
				updatedCustomRowOne.AmountAfterJcte = dataLaba.AmountAfterJcte

				updatedCustomRowOne.AmountCombineSubsidiary = dataLaba.AmountCombineSubsidiary

				updatedCustomRowOne.AmountJelimCr = dataLaba.AmountJelimCr
				updatedCustomRowOne.AmountJelimDr = dataLaba.AmountJelimDr
				updatedCustomRowOne.AmountConsole = dataLaba.AmountConsole

				_, err = s.Repository.Updates(ctx, &customRowOne.ID, &updatedCustomRowOne)
				if err != nil {
					return helper.ErrorHandler(err)
				}

				//

				updatedCustomRowThree := model.ConsolidationDetailEntityModel{}
				updatedCustomRowThree.Context = ctx
				updatedCustomRowThree.AmountBeforeJpm = dataReFa.AmountBeforeJpm
				updatedCustomRowThree.AmountJpmCr = dataReFa.AmountJpmCr
				updatedCustomRowThree.AmountJpmDr = dataReFa.AmountJpmDr
				updatedCustomRowThree.AmountAfterJpm = dataReFa.AmountAfterJpm

				updatedCustomRowThree.AmountJcteCr = dataReFa.AmountJcteCr
				updatedCustomRowThree.AmountJcteDr = dataReFa.AmountJcteDr
				updatedCustomRowThree.AmountAfterJcte = dataReFa.AmountAfterJcte

				updatedCustomRowThree.AmountCombineSubsidiary = dataReFa.AmountCombineSubsidiary

				updatedCustomRowThree.AmountJelimCr = dataReFa.AmountJelimCr
				updatedCustomRowThree.AmountJelimDr = dataReFa.AmountJelimDr
				updatedCustomRowThree.AmountConsole = dataReFa.AmountConsole

				_, err = s.Repository.Updates(ctx, &customRowThree.ID, &updatedCustomRowThree)
				if err != nil {
					return helper.ErrorHandler(err)
				}

				updateCustomRowFour := model.ConsolidationDetailEntityModel{}
				updateCustomRowFour.Context = ctx

				tmp1 := 0.0
				if dataFinIn.AmountBeforeJpm != nil && *dataFinIn.AmountBeforeJpm != 0 {
					tmp1 = *dataFinIn.AmountBeforeJpm
				}
				if dataIncomeTax.AmountBeforeJpm != nil && *dataIncomeTax.AmountBeforeJpm != 0 {
					tmp1 = tmp1 + *dataIncomeTax.AmountBeforeJpm
				}
				updateCustomRowFour.AmountBeforeJpm = &tmp1

				tmp2 := 0.0
				if dataFinIn.AmountJpmCr != nil && *dataFinIn.AmountJpmCr != 0 {
					tmp2 = *dataFinIn.AmountJpmCr
				}
				if dataIncomeTax.AmountJpmCr != nil && *dataIncomeTax.AmountJpmCr != 0 {
					tmp2 = tmp2 + *dataIncomeTax.AmountJpmCr
				}
				updateCustomRowFour.AmountJpmCr = &tmp2

				tmp3 := 0.0
				if dataFinIn.AmountJpmDr != nil && *dataFinIn.AmountJpmDr != 0 {
					tmp3 = *dataFinIn.AmountJpmDr
				}
				if dataIncomeTax.AmountJpmDr != nil && *dataIncomeTax.AmountJpmDr != 0 {
					tmp3 = tmp3 + *dataIncomeTax.AmountJpmDr
				}
				updateCustomRowFour.AmountJpmDr = &tmp3

				tmp4 := 0.0
				if dataFinIn.AmountAfterJpm != nil && *dataFinIn.AmountAfterJpm != 0 {
					tmp4 = *dataFinIn.AmountAfterJpm
				}
				if dataIncomeTax.AmountAfterJpm != nil && *dataIncomeTax.AmountAfterJpm != 0 {
					tmp4 = tmp4 + *dataIncomeTax.AmountAfterJpm
				}
				updateCustomRowFour.AmountAfterJpm = &tmp4

				//jcte

				tmp5 := 0.0
				if dataFinIn.AmountJcteCr != nil && *dataFinIn.AmountJcteCr != 0 {
					tmp5 = *dataFinIn.AmountJcteCr
				}
				if dataIncomeTax.AmountJcteCr != nil && *dataIncomeTax.AmountJcteCr != 0 {
					tmp5 = tmp5 + *dataIncomeTax.AmountJcteCr
				}
				updateCustomRowFour.AmountJcteCr = &tmp5

				tmp6 := 0.0
				if dataFinIn.AmountJcteDr != nil && *dataFinIn.AmountJcteDr != 0 {
					tmp6 = *dataFinIn.AmountJcteDr
				}
				if dataIncomeTax.AmountJcteDr != nil && *dataIncomeTax.AmountJcteDr != 0 {
					tmp6 = tmp6 + *dataIncomeTax.AmountJcteDr
				}
				updateCustomRowFour.AmountJcteDr = &tmp6

				tmp7 := 0.0
				if dataFinIn.AmountAfterJcte != nil && *dataFinIn.AmountAfterJcte != 0 {
					tmp7 = *dataFinIn.AmountAfterJcte
				}
				if dataIncomeTax.AmountAfterJcte != nil && *dataIncomeTax.AmountAfterJcte != 0 {
					tmp7 = tmp7 + *dataIncomeTax.AmountAfterJcte
				}
				updateCustomRowFour.AmountAfterJcte = &tmp7

				// acs
				tmp8 := 0.0
				if dataFinIn.AmountCombineSubsidiary != nil && *dataFinIn.AmountCombineSubsidiary != 0 {
					tmp8 = *dataFinIn.AmountCombineSubsidiary
				}
				if dataIncomeTax.AmountCombineSubsidiary != nil && *dataIncomeTax.AmountCombineSubsidiary != 0 {
					tmp8 = tmp8 + *dataIncomeTax.AmountCombineSubsidiary
				}
				updateCustomRowFour.AmountCombineSubsidiary = &tmp8

				//jelim
				tmp9 := 0.0
				if dataFinIn.AmountJelimCr != nil && *dataFinIn.AmountJelimCr != 0 {
					tmp9 = *dataFinIn.AmountJelimCr
				}
				if dataIncomeTax.AmountJelimCr != nil && *dataIncomeTax.AmountJelimCr != 0 {
					tmp9 = tmp9 + *dataIncomeTax.AmountJelimCr
				}
				updateCustomRowFour.AmountJelimCr = &tmp9

				tmp10 := 0.0
				if dataFinIn.AmountJelimDr != nil && *dataFinIn.AmountJelimDr != 0 {
					tmp10 = *dataFinIn.AmountJelimDr
				}
				if dataIncomeTax.AmountJelimDr != nil && *dataIncomeTax.AmountJelimDr != 0 {
					tmp10 = tmp10 + *dataIncomeTax.AmountJelimDr
				}
				updateCustomRowFour.AmountJelimDr = &tmp10

				tmp11 := 0.0
				if dataFinIn.AmountConsole != nil && *dataFinIn.AmountConsole != 0 {
					tmp11 = *dataFinIn.AmountConsole
				}
				if dataIncomeTax.AmountConsole != nil && *dataIncomeTax.AmountConsole != 0 {
					tmp11 = tmp11 + *dataIncomeTax.AmountConsole
				}
				updateCustomRowFour.AmountConsole = &tmp11

				_, err = s.Repository.Updates(ctx, &customRowFour.ID, &updateCustomRowFour)
				if err != nil {
					return helper.ErrorHandler(err)
				}

				updateCustomRowFive := model.ConsolidationDetailEntityModel{}
				updateCustomRowFive.Context = ctx

				tmpa := 0.0
				if dataForex.AmountBeforeJpm != nil && *dataForex.AmountBeforeJpm != 0 {
					tmpa = *dataForex.AmountBeforeJpm
				}
				if dataIncomeTax2.AmountBeforeJpm != nil && *dataIncomeTax2.AmountBeforeJpm != 0 {
					tmpa = tmpa + *dataIncomeTax2.AmountBeforeJpm
				}
				updateCustomRowFive.AmountBeforeJpm = &tmpa

				tmpb := 0.0
				if dataForex.AmountJpmCr != nil && *dataForex.AmountJpmCr != 0 {
					tmpb = *dataForex.AmountJpmCr
				}
				if dataIncomeTax2.AmountJpmCr != nil && *dataIncomeTax2.AmountJpmCr != 0 {
					tmpb = tmpb + *dataIncomeTax2.AmountJpmCr
				}
				updateCustomRowFive.AmountJpmCr = &tmpb

				tmpc := 0.0
				if dataForex.AmountJpmDr != nil && *dataForex.AmountJpmDr != 0 {
					tmpc = *dataForex.AmountJpmDr
				}
				if dataIncomeTax2.AmountJpmDr != nil && *dataIncomeTax2.AmountJpmDr != 0 {
					tmpc = tmpc + *dataIncomeTax2.AmountJpmDr
				}
				updateCustomRowFive.AmountJpmDr = &tmpc

				tmpd := 0.0
				if dataForex.AmountAfterJpm != nil && *dataForex.AmountAfterJpm != 0 {
					tmpd = *dataForex.AmountAfterJpm
				}
				if dataIncomeTax2.AmountAfterJpm != nil && *dataIncomeTax2.AmountAfterJpm != 0 {
					tmpd = tmpd + *dataIncomeTax2.AmountAfterJpm
				}
				updateCustomRowFive.AmountAfterJpm = &tmpd

				//jcte
				tmpf := 0.0
				if dataForex.AmountJcteCr != nil && *dataForex.AmountJcteCr != 0 {
					tmpf = *dataForex.AmountJcteCr
				}
				if dataIncomeTax2.AmountJcteCr != nil && *dataIncomeTax2.AmountJcteCr != 0 {
					tmpf = tmpf + *dataIncomeTax2.AmountJcteCr
				}
				updateCustomRowFive.AmountJcteCr = &tmpf

				tmpe := 0.0
				if dataForex.AmountJcteDr != nil && *dataForex.AmountJcteDr != 0 {
					tmpe = *dataForex.AmountJcteDr
				}
				if dataIncomeTax2.AmountJcteDr != nil && *dataIncomeTax2.AmountJcteDr != 0 {
					tmpe = tmpe + *dataIncomeTax2.AmountJcteDr
				}
				updateCustomRowFive.AmountJcteDr = &tmpe

				tmph := 0.0
				if dataForex.AmountAfterJcte != nil && *dataForex.AmountAfterJcte != 0 {
					tmph = *dataForex.AmountAfterJcte
				}
				if dataIncomeTax2.AmountAfterJcte != nil && *dataIncomeTax2.AmountAfterJcte != 0 {
					tmph = tmph + *dataIncomeTax2.AmountAfterJcte
				}
				updateCustomRowFive.AmountAfterJcte = &tmph

				// acs
				tmpg := 0.0
				if dataForex.AmountCombineSubsidiary != nil && *dataForex.AmountCombineSubsidiary != 0 {
					tmpg = *dataForex.AmountCombineSubsidiary
				}
				if dataIncomeTax2.AmountCombineSubsidiary != nil && *dataIncomeTax2.AmountCombineSubsidiary != 0 {
					tmpg = tmpg + *dataIncomeTax2.AmountCombineSubsidiary
				}
				updateCustomRowFive.AmountCombineSubsidiary = &tmpg

				//jelim
				tmpi := 0.0
				if dataForex.AmountJelimCr != nil && *dataForex.AmountJelimCr != 0 {
					tmpi = *dataForex.AmountJelimCr
				}
				if dataIncomeTax2.AmountJelimCr != nil && *dataIncomeTax2.AmountJelimCr != 0 {
					tmpi = tmpi + *dataIncomeTax2.AmountJelimCr
				}
				updateCustomRowFive.AmountJelimCr = &tmpi

				tmpj := 0.0
				if dataForex.AmountJelimDr != nil && *dataForex.AmountJelimDr != 0 {
					tmpj = *dataForex.AmountJelimDr
				}
				if dataIncomeTax2.AmountJelimDr != nil && *dataIncomeTax2.AmountJelimDr != 0 {
					tmpj = tmpj + *dataIncomeTax2.AmountJelimDr
				}
				updateCustomRowFive.AmountJelimDr = &tmpj

				tmpk := 0.0
				if dataForex.AmountConsole != nil && *dataForex.AmountConsole != 0 {
					tmpk = *dataForex.AmountConsole
				}
				if dataIncomeTax2.AmountConsole != nil && *dataIncomeTax2.AmountConsole != 0 {
					tmpk = tmpk + *dataIncomeTax2.AmountConsole
				}
				updateCustomRowFive.AmountConsole = &tmpk
				_, err = s.Repository.Updates(ctx, &customRowFive.ID, &updateCustomRowFive)
				if err != nil {
					return helper.ErrorHandler(err)
				}
				updateCustomRowTwo := model.ConsolidationDetailEntityModel{}
				updateCustomRowTwo.Context = ctx

				tmp12 := 0.0
				if dataKomprehensif.AmountBeforeJpm != nil && *dataKomprehensif.AmountBeforeJpm != 0 {
					tmp12 = *dataKomprehensif.AmountBeforeJpm
				}
				if summaryCodes.AmountBeforeJpm != nil && *summaryCodes.AmountBeforeJpm != 0 {
					tmp12 = tmp12 - *summaryCodes.AmountBeforeJpm
				}
				updateCustomRowTwo.AmountBeforeJpm = &tmp12

				tmp13 := 0.0
				if dataKomprehensif.AmountJpmCr != nil && *dataKomprehensif.AmountJpmCr != 0 {
					tmp13 = *dataKomprehensif.AmountJpmCr
				}
				if summaryCodes.AmountJpmCr != nil && *summaryCodes.AmountJpmCr != 0 {
					tmp13 = tmp13 - *summaryCodes.AmountJpmCr
				}
				updateCustomRowTwo.AmountJpmCr = &tmp13

				tmp14 := 0.0
				if dataKomprehensif.AmountJpmDr != nil && *dataKomprehensif.AmountJpmDr != 0 {
					tmp14 = *dataKomprehensif.AmountJpmDr
				}
				if summaryCodes.AmountJpmDr != nil && *summaryCodes.AmountJpmDr != 0 {
					tmp14 = tmp14 - *summaryCodes.AmountJpmDr
				}
				updateCustomRowTwo.AmountJpmDr = &tmp14

				tmp15 := 0.0
				if dataKomprehensif.AmountAfterJpm != nil && *dataKomprehensif.AmountAfterJpm != 0 {
					tmp15 = *dataKomprehensif.AmountAfterJpm
				}
				if summaryCodes.AmountAfterJpm != nil && *summaryCodes.AmountAfterJpm != 0 {
					tmp15 = tmp15 - *summaryCodes.AmountAfterJpm
				}
				updateCustomRowTwo.AmountAfterJpm = &tmp15

				//jcte

				tmp16 := 0.0
				if dataKomprehensif.AmountJcteCr != nil && *dataKomprehensif.AmountJcteCr != 0 {
					tmp16 = *dataKomprehensif.AmountJcteCr
				}
				if summaryCodes.AmountJcteCr != nil && *summaryCodes.AmountJcteCr != 0 {
					tmp16 = tmp16 - *summaryCodes.AmountJcteCr
				}
				updateCustomRowTwo.AmountJcteCr = &tmp16

				tmp17 := 0.0
				if dataKomprehensif.AmountJcteDr != nil && *dataKomprehensif.AmountJcteDr != 0 {
					tmp17 = *dataKomprehensif.AmountJcteDr
				}
				if summaryCodes.AmountJcteDr != nil && *summaryCodes.AmountJcteDr != 0 {
					tmp17 = tmp17 - *summaryCodes.AmountJcteDr
				}
				updateCustomRowTwo.AmountJcteDr = &tmp17

				tmp18 := 0.0
				if dataKomprehensif.AmountAfterJcte != nil && *dataKomprehensif.AmountAfterJcte != 0 {
					tmp18 = *dataKomprehensif.AmountAfterJcte
				}
				if summaryCodes.AmountAfterJcte != nil && *summaryCodes.AmountAfterJcte != 0 {
					tmp18 = tmp18 - *summaryCodes.AmountAfterJcte
				}
				updateCustomRowTwo.AmountAfterJcte = &tmp18

				// acs
				tmp19 := 0.0
				if dataKomprehensif.AmountCombineSubsidiary != nil && *dataKomprehensif.AmountCombineSubsidiary != 0 {
					tmp19 = *dataKomprehensif.AmountCombineSubsidiary
				}
				if summaryCodes.AmountCombineSubsidiary != nil && *summaryCodes.AmountCombineSubsidiary != 0 {
					tmp19 = tmp19 - *summaryCodes.AmountCombineSubsidiary
				}
				updateCustomRowTwo.AmountCombineSubsidiary = &tmp19

				//jelim
				tmp20 := 0.0
				if dataKomprehensif.AmountJelimCr != nil && *dataKomprehensif.AmountJelimCr != 0 {
					tmp20 = *dataKomprehensif.AmountJelimCr
				}
				if summaryCodes.AmountJelimCr != nil && *summaryCodes.AmountJelimCr != 0 {
					tmp20 = tmp20 - *summaryCodes.AmountJelimCr
				}
				updateCustomRowTwo.AmountJelimCr = &tmp20

				tmp21 := 0.0
				if dataKomprehensif.AmountJelimDr != nil && *dataKomprehensif.AmountJelimDr != 0 {
					tmp21 = *dataKomprehensif.AmountJelimDr
				}
				if summaryCodes.AmountJelimDr != nil && *summaryCodes.AmountJelimDr != 0 {
					tmp21 = tmp21 - *summaryCodes.AmountJelimDr
				}
				updateCustomRowTwo.AmountJelimDr = &tmp21

				tmp22 := 0.0
				if dataKomprehensif.AmountConsole != nil && *dataKomprehensif.AmountConsole != 0 {
					tmp22 = *dataKomprehensif.AmountConsole
				}
				if summaryCodes.AmountConsole != nil && *summaryCodes.AmountConsole != 0 {
					tmp22 = tmp22 - *summaryCodes.AmountConsole
				}
				updateCustomRowTwo.AmountConsole = &tmp22
				_, err = s.Repository.Updates(ctx, &customRowTwo.ID, &updateCustomRowTwo)
				if err != nil {
					return helper.ErrorHandler(err)
				}
			}
		}
		for _, v := range *formatterDetailSumData {
			criteriaTBDetail := model.ConsolidationDetailFilterModel{}
			criteriaTBDetail.ConsolidationID = &JpmID.ConsolidationID
			if v.AutoSummary != nil && *v.AutoSummary {
				code := fmt.Sprintf("%s_Subtotal", v.Code)
				criteriaTBDetail.Code = &code
				criteriaTBDetail.ConsolidationID = &JpmID.ConsolidationID
				mfadetailsum, _, err := s.Repository.FindC(ctx, &criteriaTBDetail, &abstraction.Pagination{})
				if err != nil {
					return response.CustomErrorBuilder(http.StatusBadRequest, "Gagal Find Detail Consolidation", "Gagal Find Detail Consolidation")
				}

				for _, a := range *mfadetailsum {
					sumTBD, err := s.Repository.FindSummarys(ctx, &v.Code, &JpmID.ConsolidationID, v.IsCoa)
					if err != nil {
						return response.CustomErrorBuilder(http.StatusBadRequest, "Gagal Menjumlahkan Summary Coa", "Gagal Menjumlahkan Summary Coa")
					}

					updateSummary := model.ConsolidationDetailEntityModel{
						ConsolidationDetailEntity: model.ConsolidationDetailEntity{
							AmountJpmCr:             sumTBD.AmountJpmCr,
							AmountJpmDr:             sumTBD.AmountJpmDr,
							AmountAfterJpm:          sumTBD.AmountAfterJpm,
							AmountJcteCr:            sumTBD.AmountJcteCr,
							AmountJcteDr:            sumTBD.AmountJcteDr,
							AmountAfterJcte:         sumTBD.AmountAfterJcte,
							AmountCombineSubsidiary: sumTBD.AmountCombineSubsidiary,
							AmountJelimCr:           sumTBD.AmountJelimCr,
							AmountJelimDr:           sumTBD.AmountJelimDr,
							AmountConsole:           sumTBD.AmountConsole,
						},
					}

					_, err = s.Repository.Updates(ctx, &a.ID, &updateSummary)
					if err != nil {
						return response.CustomErrorBuilder(http.StatusBadRequest, "Gagal Update Consolidation Detail", "Gagal Update Consolidation Detail")
					}
				}
			}
			if v.IsTotal != nil && *v.IsTotal == true && v.FxSummary != "" {

				tmpString := []string{"AmountBeforeJpm", "AmountJpmCr", "AmountJpmDr", "AmountAfterJpm", "AmountJcteCr", "AmountJcteDr", "AmountAfterJcte", "AmountCombineSubsidiary", "AmountJelimCr", "AmountJelimDr", "AmountConsole"}
				tmpTotalFl := make(map[string]*float64)
				// reg := regexp.MustCompile(`[0-9]+\d{3,}`)
				reg := regexp.MustCompile(`[A-Za-z0-9_~:()]+|[0-9]+\d{3,}`)

				for _, tipe := range tmpString {
					formula := strings.ToUpper(v.FxSummary)
					if tipe == "AmountJpmCr" || tipe == "AmountJpmDr" || tipe == "AmountJcteCr"|| tipe == "AmountJcteDr"|| tipe == "AmountJelimDr"|| tipe == "AmountJelimCr"{
						newStr := strings.ReplaceAll(v.FxSummary, "-", "+")
						formula = strings.ToUpper(newStr)
					}
					match := reg.FindAllString(formula, -1)
					amountBeforeJpm := make(map[string]interface{}, 0)
					for _, vMatch := range match {
						//cari jml berdasarkan code

						if len(vMatch) < 3 {
							continue
						}
						sumTBD, err := s.Repository.FindSummarys(ctx, &vMatch, &JpmID.ConsolidationID, v.IsCoa)
						if err != nil {
							return response.CustomErrorBuilder(http.StatusBadRequest, "Gagal Mendapatkan Jumlah", "Gagal Mendapatkan Jumlah")
						}
						angka := 0.0

						if tipe == "AmountBeforeJpm" && sumTBD.AmountBeforeJpm != nil {
							angka = *sumTBD.AmountBeforeJpm
						} else if tipe == "AmountJpmCr" && sumTBD.AmountJpmCr != nil {
							angka = *sumTBD.AmountJpmCr
						} else if tipe == "AmountJpmDr" && sumTBD.AmountJpmDr != nil {
							angka = *sumTBD.AmountJpmDr
						} else if tipe == "AmountAfterJpm" && sumTBD.AmountAfterJpm != nil {
							angka = *sumTBD.AmountAfterJpm
						} else if tipe == "AmountJcteCr" && sumTBD.AmountJcteCr != nil {
							angka = *sumTBD.AmountJcteCr
						} else if tipe == "AmountJcteDr" && sumTBD.AmountJcteDr != nil {
							angka = *sumTBD.AmountJcteDr
						} else if tipe == "AmountAfterJcte" && sumTBD.AmountAfterJcte != nil {
							angka = *sumTBD.AmountAfterJcte
						} else if tipe == "AmountCombineSubsidiary" && sumTBD.AmountCombineSubsidiary != nil {
							angka = *sumTBD.AmountCombineSubsidiary
						} else if tipe == "AmountJelimCr" && sumTBD.AmountJelimCr != nil {
							angka = *sumTBD.AmountJelimCr
						} else if tipe == "AmountJelimDr" && sumTBD.AmountJelimDr != nil {
							angka = *sumTBD.AmountJelimDr
						} else if tipe == "AmountConsole" && sumTBD.AmountConsole != nil {
							angka = *sumTBD.AmountConsole
						}

						formula = helper.ReplaceWholeWord(formula, vMatch, fmt.Sprintf("(%2.f)", angka))
						// parameters[vMatch] = angka

					}

					expressionFormula, err := govaluate.NewEvaluableExpression(formula)
					if err != nil {
						fmt.Println(err)
						return err
					}
					JpmID, err := expressionFormula.Evaluate(amountBeforeJpm)
					if err != nil {
						return response.CustomErrorBuilder(http.StatusBadRequest, "Tidak Dapat Masukan Balance Sheet Debit", "Tidak Dapat Masukan Balance Sheet Debit")
					}
					tmp := JpmID.(float64)
					tmpTotalFl[tipe] = &tmp

				}
				criteriaTBDetail.ConsolidationID = &JpmID.ConsolidationID
				criteriaTBDetail.Code = &v.Code
				mfadetailsum, err := s.Repository.FindDetailConsole(ctx, &criteriaTBDetail)
				if err != nil {
					fmt.Sprintln(err)
					return err
				}

				updateSummary := model.ConsolidationDetailEntityModel{
					ConsolidationDetailEntity: model.ConsolidationDetailEntity{
						AmountJpmCr:             tmpTotalFl["AmountJpmCr"],
						AmountJpmDr:             tmpTotalFl["AmountJpmDr"],
						AmountAfterJpm:          tmpTotalFl["AmountAfterJpm"],
						AmountJcteCr:            tmpTotalFl["AmountJcteCr"],
						AmountJcteDr:            tmpTotalFl["AmountJcteDr"],
						AmountAfterJcte:         tmpTotalFl["AmountAfterJcte"],
						AmountCombineSubsidiary: tmpTotalFl["AmountCombineSubsidiary"],
						AmountJelimCr:           tmpTotalFl["AmountJelimCr"],
						AmountJelimDr:           tmpTotalFl["AmountJelimDr"],
						AmountConsole:           tmpTotalFl["AmountConsole"],
					},
				}

				_, err = s.Repository.Updates(ctx, &mfadetailsum.ID, &updateSummary)
				if err != nil {
					return response.CustomErrorBuilder(http.StatusBadRequest, "Tidak Dapat Masukan Balance Sheet Debit", "Tidak Dapat Masukan Balance Sheet Debit")
				}

			}
		}

		for _, v := range *formatterDetailSumData {
			criteriaTBDetail := model.ConsolidationDetailFilterModel{}
			criteriaTBDetail.ConsolidationID = &JpmID.ConsolidationID

			if v.AutoSummary != nil && *v.AutoSummary {
				code := fmt.Sprintf("%s_Subtotal", v.Code)
				criteriaTBDetail.Code = &code
				criteriaTBDetail.ConsolidationID = &JpmID.ConsolidationID
				mfadetailsum, _, err := s.Repository.FindC(ctx, &criteriaTBDetail, &abstraction.Pagination{})
				if err != nil {
					return response.CustomErrorBuilder(http.StatusBadRequest, "Gagal Find Detail Consolidation", "Gagal Find Detail Consolidation")
				}

				for _, a := range *mfadetailsum {
					sumTBD, err := s.Repository.FindSummarys(ctx, &v.Code, &JpmID.ConsolidationID, v.IsCoa)
					if err != nil {
						return response.CustomErrorBuilder(http.StatusBadRequest, "Gagal Menjumlahkan Summary Coa", "Gagal Menjumlahkan Summary Coa")
					}

					updateSummary := model.ConsolidationDetailEntityModel{
						ConsolidationDetailEntity: model.ConsolidationDetailEntity{
							AmountJpmCr:             sumTBD.AmountJpmCr,
							AmountJpmDr:             sumTBD.AmountJpmDr,
							AmountAfterJpm:          sumTBD.AmountAfterJpm,
							AmountJcteCr:            sumTBD.AmountJcteCr,
							AmountJcteDr:            sumTBD.AmountJcteDr,
							AmountAfterJcte:         sumTBD.AmountAfterJcte,
							AmountCombineSubsidiary: sumTBD.AmountCombineSubsidiary,
							AmountJelimCr:           sumTBD.AmountJelimCr,
							AmountJelimDr:           sumTBD.AmountJelimDr,
							AmountConsole:           sumTBD.AmountConsole,
						},
					}

					_, err = s.Repository.Updates(ctx, &a.ID, &updateSummary)
					if err != nil {
						return response.CustomErrorBuilder(http.StatusBadRequest, "Gagal Update Consolidation Detail", "Gagal Update Consolidation Detail")
					}
				}
			}

			if v.IsTotal != nil && *v.IsTotal == true && v.FxSummary != "" {

				tmpString := []string{"AmountBeforeJpm", "AmountJpmCr", "AmountJpmDr", "AmountAfterJpm", "AmountJcteCr", "AmountJcteDr", "AmountAfterJcte", "AmountCombineSubsidiary", "AmountJelimCr", "AmountJelimDr", "AmountConsole"}
				tmpTotalFl := make(map[string]*float64)
				// reg := regexp.MustCompile(`[0-9]+\d{3,}`)
				reg := regexp.MustCompile(`[A-Za-z0-9_~:()]+|[0-9]+\d{3,}`)

				for _, tipe := range tmpString {
					formula := strings.ToUpper(v.FxSummary)
					if tipe == "AmountJpmCr" || tipe == "AmountJpmDr" || tipe == "AmountJcteCr"|| tipe == "AmountJcteDr"|| tipe == "AmountJelimDr"|| tipe == "AmountJelimCr"{
						newStr := strings.ReplaceAll(v.FxSummary, "-", "+")
						formula = strings.ToUpper(newStr)
					}
					match := reg.FindAllString(formula, -1)
					amountBeforeJpm := make(map[string]interface{}, 0)
					for _, vMatch := range match {
						//cari jml berdasarkan code

						if len(vMatch) < 3 {
							continue
						}
						sumTBD, err := s.Repository.FindSummarys(ctx, &vMatch, &JpmID.ConsolidationID, v.IsCoa)
						if err != nil {
							return response.CustomErrorBuilder(http.StatusBadRequest, "Gagal Mendapatkan Jumlah", "Gagal Mendapatkan Jumlah")
						}
						angka := 0.0

						if tipe == "AmountBeforeJpm" && sumTBD.AmountBeforeJpm != nil {
							angka = *sumTBD.AmountBeforeJpm
						} else if tipe == "AmountJpmCr" && sumTBD.AmountJpmCr != nil {
							angka = *sumTBD.AmountJpmCr
						} else if tipe == "AmountJpmDr" && sumTBD.AmountJpmDr != nil {
							angka = *sumTBD.AmountJpmDr
						} else if tipe == "AmountAfterJpm" && sumTBD.AmountAfterJpm != nil {
							angka = *sumTBD.AmountAfterJpm
						} else if tipe == "AmountJcteCr" && sumTBD.AmountJcteCr != nil {
							angka = *sumTBD.AmountJcteCr
						} else if tipe == "AmountJcteDr" && sumTBD.AmountJcteDr != nil {
							angka = *sumTBD.AmountJcteDr
						} else if tipe == "AmountAfterJcte" && sumTBD.AmountAfterJcte != nil {
							angka = *sumTBD.AmountAfterJcte
						} else if tipe == "AmountCombineSubsidiary" && sumTBD.AmountCombineSubsidiary != nil {
							angka = *sumTBD.AmountCombineSubsidiary
						} else if tipe == "AmountJelimCr" && sumTBD.AmountJelimCr != nil {
							angka = *sumTBD.AmountJelimCr
						} else if tipe == "AmountJelimDr" && sumTBD.AmountJelimDr != nil {
							angka = *sumTBD.AmountJelimDr
						} else if tipe == "AmountConsole" && sumTBD.AmountConsole != nil {
							angka = *sumTBD.AmountConsole
						}

						formula = helper.ReplaceWholeWord(formula, vMatch, fmt.Sprintf("(%2.f)", angka))
						// parameters[vMatch] = angka

					}

					expressionFormula, err := govaluate.NewEvaluableExpression(formula)
					if err != nil {
						fmt.Println(err)
						return err
					}
					JpmID, err := expressionFormula.Evaluate(amountBeforeJpm)
					if err != nil {
						return response.CustomErrorBuilder(http.StatusBadRequest, "Tidak Dapat Masukan Balance Sheet Debit", "Tidak Dapat Masukan Balance Sheet Debit")
					}
					tmp := JpmID.(float64)
					tmpTotalFl[tipe] = &tmp

				}
				criteriaTBDetail.ConsolidationID = &JpmID.ConsolidationID
				criteriaTBDetail.Code = &v.Code
				mfadetailsum, err := s.Repository.FindDetailConsole(ctx, &criteriaTBDetail)
				if err != nil {
					fmt.Sprintln(err)
					return err
				}

				updateSummary := model.ConsolidationDetailEntityModel{
					ConsolidationDetailEntity: model.ConsolidationDetailEntity{
						AmountJpmCr:             tmpTotalFl["AmountJpmCr"],
						AmountJpmDr:             tmpTotalFl["AmountJpmDr"],
						AmountAfterJpm:          tmpTotalFl["AmountAfterJpm"],
						AmountJcteCr:            tmpTotalFl["AmountJcteCr"],
						AmountJcteDr:            tmpTotalFl["AmountJcteDr"],
						AmountAfterJcte:         tmpTotalFl["AmountAfterJcte"],
						AmountCombineSubsidiary: tmpTotalFl["AmountCombineSubsidiary"],
						AmountJelimCr:           tmpTotalFl["AmountJelimCr"],
						AmountJelimDr:           tmpTotalFl["AmountJelimDr"],
						AmountConsole:           tmpTotalFl["AmountConsole"],
					},
				}

				_, err = s.Repository.Updates(ctx, &mfadetailsum.ID, &updateSummary)
				if err != nil {
					return response.CustomErrorBuilder(http.StatusBadRequest, "Tidak Dapat Masukan Balance Sheet Debit", "Tidak Dapat Masukan Balance Sheet Debit")
				}

			}
		}

		return nil
	}); err != nil {
		return &dto.JpmDeleteResponse{}, err
	}
	result := &dto.JpmDeleteResponse{
		JpmEntityModel: data,
	}
	return result, nil
}

func (s *service) GetVersion(ctx *abstraction.Context, payload *dto.GetVersionRequest) (*dto.GetVersionResponse, error) {
	filter := model.JpmFilterModel{
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

func (s *service) Export(ctx *abstraction.Context, payload *dto.JpmExportRequest) (*dto.JpmExportResponse, error) {

	datas, err := s.Repository.Export(ctx, &payload.JpmID)
	if err != nil {
		return &dto.JpmExportResponse{}, response.ErrorBuilder(&response.ErrorConstant.InternalServerError, err)
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
	f.SetCellValue(sheet, "B4", ": JPM")

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

	for i, v := range datas.JpmDetail {

		f.SetCellValue(sheet, fmt.Sprintf("A%d", row), i+1)
		f.SetCellValue(sheet, fmt.Sprintf("B%d", row), v.CoaCode)
		f.SetCellValue(sheet, fmt.Sprintf("E%d", row), *v.Description)
		f.SetCellValue(sheet, fmt.Sprintf("E%d", row+1), *v.Note)
		f.SetCellValue(sheet, fmt.Sprintf("G%d", row), *v.BalanceSheetDr)
		f.SetCellValue(sheet, fmt.Sprintf("H%d", row), *v.BalanceSheetCr)
		f.SetCellValue(sheet, fmt.Sprintf("I%d", row), *v.IncomeStatementDr)
		f.SetCellValue(sheet, fmt.Sprintf("J%d", row), *v.IncomeStatementCr)

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
	fileName := fmt.Sprintf("Jurnal_PM_%s.xlsx", period)
	fileLoc := fmt.Sprintf("assets/%d/%s", ctx.Auth.ID, fileName)
	err = f.SaveAs(fileLoc)
	if err != nil {
		return nil, err
	}

	result := &dto.JpmExportResponse{
		FileName: fileName,
		Path:     fileLoc,
	}
	return result, nil
}
