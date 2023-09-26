package consolidation

import (
	"encoding/json"
	"errors"
	"fmt"
	"mcash-finance-console-core/internal/abstraction"
	"mcash-finance-console-core/internal/dto"
	"mcash-finance-console-core/internal/factory"
	"mcash-finance-console-core/internal/model"
	"mcash-finance-console-core/internal/repository"
	"mcash-finance-console-core/pkg/kafka"
	"mcash-finance-console-core/pkg/util/helper"
	"mcash-finance-console-core/pkg/util/response"
	"mcash-finance-console-core/pkg/util/trxmanager"
	"net/http"
	"strconv"
	"time"

	"gorm.io/gorm"
)

type service struct {
	Repository             repository.Consolidation
	TrialBalanceRepository repository.TrialBalance
	AjeRepository          repository.Adjustment ``
	Db                     *gorm.DB
}

type Service interface {
	Find(ctx *abstraction.Context, payload *dto.ConsolidationGetRequest) (*dto.ConsolidationGetResponse, error)
	ListAvailable(ctx *abstraction.Context, payload *dto.ConsolidationGetListAvailable) (*dto.ConsolidationGetListAvaibleResponse, error)
	GetVersion(ctx *abstraction.Context, payload *dto.GetVersionRequest) (*dto.GetVersionResponse, error)
	GetControl(ctx *abstraction.Context, payload *dto.ConsolidationGetControlRequest) (*dto.ConsolidationGetControlResponse, error)
}

func NewService(f *factory.Factory) *service {
	repository := f.ConsolidationRepository
	trialBalanceRepository := f.TrialBalanceRepository
	db := f.Db
	return &service{
		Repository:             repository,
		TrialBalanceRepository: trialBalanceRepository,
		Db:                     db,
	}
}

type DataToConsolidate struct {
	ConsolidatedID     int
	MasterID           int
	ListDataID         []int
	ListConsolidatedID []int
}

func (s *service) Find(ctx *abstraction.Context, payload *dto.ConsolidationGetRequest) (*dto.ConsolidationGetResponse, error) {
	data, info, err := s.Repository.Find(ctx, &payload.ConsolidationFilterModel, &payload.Pagination)
	if err != nil {
		return &dto.ConsolidationGetResponse{}, response.ErrorBuilder(&response.ErrorConstant.InternalServerError, err)
	}
	result := &dto.ConsolidationGetResponse{
		Datas:          *data,
		PaginationInfo: *info,
	}
	return result, nil
}

func (s *service) ListAvailable(ctx *abstraction.Context, payload *dto.ConsolidationGetListAvailable) (*dto.ConsolidationGetListAvaibleResponse, error) {
	criteriaTB := model.TrialBalanceFilterModel{}
	criteriaTB.Period = &payload.Period
	criteriaTB.CompanyID = &payload.CompanyID

	paginationTB := abstraction.Pagination{}
	paginationTB.Page = payload.Page
	paginationTB.PageSize = payload.PageSize
	paginationTB.Sort = payload.Sort
	paginationTB.SortBy = payload.SortBy

	dataTBParent, _, err := s.Repository.FindListCompanyParents(ctx, &criteriaTB, &paginationTB)
	if err != nil {
		return nil, err
	}

	dataCompanyParentCompanyID, err := s.Repository.FindByParentCompanyID(ctx, criteriaTB.CompanyID)
	if err != nil {
		return nil, err
	}

	var TbParent []model.TrialBalanceEntityModel
	var ConsolidationParent []model.ConsolidationEntityModel
	var ConsolidationParents []model.ConsolidationEntityModel
	for _, v := range *dataCompanyParentCompanyID {

		dataCompanyParentCompanyIDS, err := s.Repository.FindByParentCompanyID(ctx, &v.ID)
		if err != nil && err != gorm.ErrRecordNotFound {
			fmt.Println(err)
			return nil, err
		}
		if len(*dataCompanyParentCompanyIDS) == 0 {
			criteriaTBS := model.TrialBalanceFilterModel{}
			criteriaTBS.Period = &payload.Period
			criteriaTBS.CompanyID = &v.ID
			dataTBChild, _, err := s.Repository.FindListCompanyChilds(ctx, &criteriaTBS, &paginationTB)
			if err != nil {
				return nil, err
			}

			TbParent = append(TbParent, *dataTBChild...)
		}
		if len(*dataCompanyParentCompanyIDS) != 0 {
			criteriaTBS := model.ConsolidationFilterModel{}
			criteriaTBS.Period = &payload.Period
			criteriaTBS.CompanyID = &v.ID

			dataTBChildConsole, _, err := s.Repository.FindListCompanyChildConsole(ctx, &criteriaTBS, &paginationTB)
			if err != nil {
				return nil, err
			}

			// for _, c := range *dataTBChildConsole {
			// 	c.StatusString = "Consolidation"
			// 	ConsolidationParent = append(ConsolidationParent,c)
			// }
			ConsolidationParent = append(ConsolidationParent, *dataTBChildConsole...)

			for _, c := range *dataTBChildConsole {
				c.StatusString = "Consolidation"
				ConsolidationParents = append(ConsolidationParents, c)
			}
		}

	}

	// dataTBChild, _, err := s.Repository.FindListCompanyChild(ctx, &criteriaTB, &paginationTB)
	// if err != nil {
	// 	return nil, err
	// }

	results := dto.ConsolidationGetListAvaibleResponse{
		Parent:      *dataTBParent,
		ChildOnly:   TbParent,
		ChildParent: ConsolidationParents,
	}

	return &results, nil
}
func (s *service) ListDuplicateAvailable(ctx *abstraction.Context, payload *dto.ConsolidationGetListDuplicateAvailable) (*dto.ConsolidationGetListDuplicateAvaibleResponse, error) {

	dataConsole, err := s.Repository.FindByID(ctx, &payload.ConsolidationID)
	if err != nil {
		return nil, err
	}
	datePeriod, err := time.Parse(time.RFC3339, dataConsole.Period)
	if err != nil {
		return nil, err
	}
	period := datePeriod.Format("2006-01-02")

	criteriaTB := model.TrialBalanceFilterModel{}
	criteriaTB.Period = &period
	criteriaTB.CompanyID = &dataConsole.CompanyID
	criteriaTB.Versions = &dataConsole.Versions

	dataTb, err := s.Repository.FindByTBID(ctx, &criteriaTB)
	if err != nil {
		return nil, err
	}
	criteriaCB := model.ConsolidationBridgeFilterModel{}
	criteriaCB.Period = &period
	criteriaCB.ConsolidationID = &payload.ConsolidationID

	paginationTB := abstraction.Pagination{}
	paginationTB.Page = payload.Page
	paginationTB.PageSize = payload.PageSize
	paginationTB.Sort = payload.Sort
	paginationTB.SortBy = payload.SortBy

	dataTBParent, _, err := s.Repository.FindListCompanyParent(ctx, &criteriaTB, dataTb.ID, &paginationTB)
	if err != nil {
		return nil, err
	}

	dataCompanyParentCompanyID, err := s.Repository.FindByParentCompanyID(ctx, criteriaTB.CompanyID)
	if err != nil {
		return nil, err
	}
	dataTbDuplicate, _, err := s.Repository.FindListCompanyConsole(ctx, &criteriaCB, &paginationTB)
	if err != nil {
		return nil, err
	}

	type DataToConsolidate struct {
		ListDataIDConfirm       []model.TrialBalanceEntityModel
		ListDataIDConsolidation []model.ConsolidationEntityModel
	}
	var idConfirm = []int{
		0,
	}
	var idConsole = []int{
		0,
	}
	var tmpData DataToConsolidate
	for _, vDataconsolidation := range *dataTbDuplicate {

		if vDataconsolidation.ConsolidationVersions == 0 {
			period := datePeriod.Format("2006-01-02")
			criteriaTB := model.TrialBalanceFilterModel{}
			criteriaTB.Period = &period
			criteriaTB.CompanyID = &vDataconsolidation.CompanyID
			criteriaTB.Versions = &vDataconsolidation.Versions
			dataTbIDconsole, err := s.Repository.FindByTBID(ctx, &criteriaTB)
			if err != nil {
				return nil, err
			}
			idConfirm = append(idConfirm, dataTbIDconsole.ID)
			tmpData.ListDataIDConfirm = append(tmpData.ListDataIDConfirm, *dataTbIDconsole)
		}

		if vDataconsolidation.ConsolidationVersions != 0 {
			period := datePeriod.Format("2006-01-02")
			criteriaC := model.ConsolidationFilterModel{}
			criteriaC.Period = &period
			criteriaC.CompanyID = &vDataconsolidation.CompanyID
			criteriaC.ConsolidationVersions = &vDataconsolidation.ConsolidationVersions
			dataConsolidationID, err := s.Repository.FindByConsolidationID(ctx, &criteriaC)
			if err != nil {
				return nil, err
			}
			idConsole = append(idConsole, dataConsolidationID.ID)
			tmpData.ListDataIDConsolidation = append(tmpData.ListDataIDConsolidation, *dataConsolidationID)
		}

	}
	var TbParent []model.TrialBalanceEntityModel
	var ConsolidationParent []model.ConsolidationEntityModel
	var ConsolidationParents []model.ConsolidationEntityModel

	for _, v := range *dataCompanyParentCompanyID {

		dataCompanyParentCompanyIDS, err := s.Repository.FindByParentCompanyID(ctx, &v.ID)
		if err != nil && err != gorm.ErrRecordNotFound {
			fmt.Println(err)
			return nil, err
		}
		if len(*dataCompanyParentCompanyIDS) == 0 {
			criteriaTBS := model.TrialBalanceFilterModel{}
			criteriaTBS.Period = &dataConsole.Period
			criteriaTBS.CompanyID = &v.ID
			dataTBChild, _, err := s.Repository.FindListCompanyChild(ctx, &criteriaTBS, idConfirm, &paginationTB)
			if err != nil {
				return nil, err
			}

			TbParent = append(TbParent, *dataTBChild...)
		}
		if len(*dataCompanyParentCompanyIDS) != 0 {
			// period := dataConsole.Format("2006-01-02")
			criteriaTBS := model.ConsolidationFilterModel{}
			criteriaTBS.Period = &dataConsole.Period
			criteriaTBS.CompanyID = &v.ID
			dataTBChildConsole, _, err := s.Repository.FindListCompanyChildConsoles(ctx, &criteriaTBS, idConsole, &paginationTB)
			if err != nil {
				return nil, err
			}

			// for _, c := range *dataTBChildConsole {
			// 	c.StatusString = "Consolidation"
			// 	ConsolidationParent = append(ConsolidationParent,c)
			// }
			ConsolidationParent = append(ConsolidationParent, *dataTBChildConsole...)

			for _, c := range *dataTBChildConsole {
				c.StatusString = "Consolidation"
				ConsolidationParents = append(ConsolidationParents, c)
			}
		}

	}
	results := dto.ConsolidationGetListDuplicateAvaibleResponse{
		Parent:                       *dataTBParent,
		ChildOnly:                    TbParent,
		ChildParent:                  ConsolidationParents,
		ConsolidationParent:          *dataTb,
		ConsolidationChildOnly:       tmpData.ListDataIDConfirm,
		ConsolidationChildParentOnly: tmpData.ListDataIDConsolidation,
		General:                      *dataConsole,
	}

	return &results, nil
}
func (s *service) RequestToCombaine(ctx *abstraction.Context, payload *dto.ConsolidationCombaineRequest) error {
	if len(payload.ListConsolidation) == 0 && len(payload.ListConsolidationParent) == 0 {
		return response.CustomErrorBuilder(http.StatusBadRequest, "Data Anak Usaha Harus di isi", "Data Anak Usaha Harus di isi")
	}
	var tmpData DataToConsolidate
	dataTB, err := s.TrialBalanceRepository.FindByID(ctx, &payload.ConsolidationMasterID)
	if err != nil {
		return err
	}
	tmpData.MasterID = dataTB.ID
	tmpData.ListDataID = payload.ListConsolidation
	// tmpDataConsolidationdID = payload.ListConsolidationParent
	tmpData.ListConsolidatedID = payload.ListConsolidationParent

	// for _, vDataconsolidation := range payload.ListConsolidation {
	// 	criteriaConsolidation := model.ConsolidationFilterModel{}
	// 	criteriaConsolidation.ConsolidationID = &vDataconsolidation

	// 	dataConsole, err := s.TrialBalanceRepository.FindByID(ctx, &vDataconsolidation)
	// 	if err != nil {
	// 		return err
	// 	}

	// 	if dataConsole.Status == 2 {
	// 		tmpData.ListDataID = append(tmpData.ListDataID, dataConsole.ID)

	// 	}
	// 	if dataConsole.Status == 3 {
	// 		tmpData.ConsolidatedID = append(tmpData.ConsolidatedID, dataConsole.ID)
	// 	}
	// }

	jsonTmpData, err := json.Marshal(tmpData)
	if err != nil {
		return err
	}

	data := kafka.JsonData{}
	data.Name = "COMBINE"
	data.CompanyID = dataTB.CompanyID
	data.UserID = ctx.Auth.ID
	data.Data = string(jsonTmpData)
	data.Filter.Period = dataTB.Period
	data.Filter.Versions = dataTB.Versions

	jsonStr, err := json.Marshal(data)
	if err != nil {
		return err
	}

	go kafka.NewService("CONSOLIDATION").SendMessage("CONSOLIDATE", string(jsonStr))
	return nil
}
func (s *service) RequestToConsolidation(ctx *abstraction.Context, payload *dto.ConsolidationConsolidateRequest) error {

	if len(payload.ListConsolidation) == 0 && len(payload.ListConsolidationParent) == 0 {
		return response.CustomErrorBuilder(http.StatusBadRequest, "Data Anak Usaha Harus di isi", "Data Anak Usaha Harus di isi")
	}
	var tmpData DataToConsolidate
	dataTB, err := s.TrialBalanceRepository.FindByID(ctx, &payload.ConsolidationMasterID)
	if err != nil {
		return err
	}

	if dataTB.Status != 2 {
		return errors.New("data belum tervalidasi atau sudah terkonsolidasi")
	}
	// datePeriod, err := time.Parse(time.RFC3339, dataTB.Period)
	// if err != nil {
	// 	return err
	// }
	// period := datePeriod.Format("2006-01-02")
	// criteriaTB := model.ConsolidationFilterModel{}
	// criteriaTB.Period = &period
	// criteriaTB.CompanyID = &dataTB.CompanyID
	// criteriaTB.Versions = &dataTB.Versions
	dataTbIDconsole, err := s.Repository.FindByID(ctx, &payload.ConsolidationID)
	if err != nil {
		return err
	}

	if dataTbIDconsole.Status == 0 {
		return errors.New("data sedang dalam proses consolidate ,tunggu hingga proses conolidation selesai")
	}

	tmpData.MasterID = dataTB.ID
	tmpData.ListDataID = payload.ListConsolidation
	tmpData.ConsolidatedID = dataTbIDconsole.ID
	tmpData.ListConsolidatedID = payload.ListConsolidationParent

	// for _, vDataconsolidation := range payload.ListConsolidation {
	// 	criteriaConsolidation := model.ConsolidationFilterModel{}
	// 	criteriaConsolidation.ConsolidationID = &vDataconsolidation

	// 	dataConsole, err := s.TrialBalanceRepository.FindByID(ctx, &vDataconsolidation)
	// 	if err != nil {
	// 		return err
	// 	}

	// 	if dataConsole.Status == 2 {
	// 		tmpData.ListDataID = append(tmpData.ListDataID, dataConsole.ID)

	// 	}
	// 	if dataConsole.Status == 3 {
	// 		tmpData.ConsolidatedID = append(tmpData.ConsolidatedID, dataConsole.ID)
	// 	}
	// }
	var cd model.ConsolidationEntityModel

	if err := trxmanager.New(s.Db).WithTrx(ctx, func(ctx *abstraction.Context) error {
		_, err := s.Repository.FindByID(ctx, &payload.ConsolidationID)
		if err != nil {
			return helper.ErrorHandler(err)
		}

		// if existing.Status != 1 {
		// 	return response.ErrorBuilder(&response.ErrorConstant.DataValidated, errors.New("Cannot Update Data"))
		// }

		// allowed := helper.CompanyValidation(ctx.Auth.ID, existing.CompanyID)
		// allowed2 := helper.CompanyValidation(ctx.Auth.ID, payload.CompanyID)
		// if !allowed || !allowed2 {
		// 	return response.ErrorBuilder(&response.ErrorConstant.UnauthorizedAccess, errors.New("Not Allowed"))
		// }

		cd.Context = ctx
		cd.ID = payload.ConsolidationID
		cd.Status = 0
		cd.ConsolidationEntity.Status = 0
		_, err = s.Repository.Update(ctx, &payload.ConsolidationID, &cd)
		if err != nil {
			return helper.ErrorHandler(err)
		}
		// data = *result
		return nil
	}); err != nil {
		return err
	}

	jsonTmpData, err := json.Marshal(tmpData)
	if err != nil {
		return err
	}

	data := kafka.JsonData{}
	data.Name = "CONSOLIDATION"
	data.CompanyID = dataTB.CompanyID
	data.UserID = ctx.Auth.ID
	data.Data = string(jsonTmpData)
	data.Filter.Period = dataTB.Period
	data.Filter.Versions = dataTB.Versions

	jsonStr, err := json.Marshal(data)
	if err != nil {
		return err
	}

	go kafka.NewService("CONSOLIDATION").SendMessage("CONSOLIDATE", string(jsonStr))
	return nil
}
func (s *service) RequestToDuplicate(ctx *abstraction.Context, payload *dto.ConsolidationConsolidateRequest) error {
	if len(payload.ListConsolidation) == 0 && len(payload.ListConsolidationParent) == 0 {
		return response.CustomErrorBuilder(http.StatusBadRequest, "Data Anak Usaha Harus di isi", "Data Anak Usaha Harus di isi")
	}
	var tmpData DataToConsolidate
	dataTB, err := s.TrialBalanceRepository.FindByID(ctx, &payload.ConsolidationMasterID)
	if err != nil {
		return err
	}

	// if dataTB.Status != 2 {
	// 	return errors.New(fmt.Sprintf("data belum tervalidasi atau sudah terkonsolidasi"))
	// }
	// datePeriod, err := time.Parse(time.RFC3339, dataTB.Period)
	// if err != nil {
	// 	return err
	// }
	// period := datePeriod.Format("2006-01-02")
	// criteriaTB := model.ConsolidationFilterModel{}
	// criteriaTB.Period = &period
	// criteriaTB.CompanyID = &dataTB.CompanyID
	// criteriaTB.Versions = &dataTB.Versions
	dataTbIDconsole, err := s.Repository.FindByID(ctx, &payload.ConsolidationID)
	if err != nil {
		return err
	}

	tmpData.MasterID = dataTB.ID
	tmpData.ListDataID = payload.ListConsolidation
	tmpData.ConsolidatedID = dataTbIDconsole.ID
	tmpData.ListConsolidatedID = payload.ListConsolidationParent

	// for _, vDataconsolidation := range payload.ListConsolidation {
	// 	criteriaConsolidation := model.ConsolidationFilterModel{}
	// 	criteriaConsolidation.ConsolidationID = &vDataconsolidation

	// 	dataConsole, err := s.TrialBalanceRepository.FindByID(ctx, &vDataconsolidation)
	// 	if err != nil {
	// 		return err
	// 	}

	// 	if dataConsole.Status == 2 {
	// 		tmpData.ListDataID = append(tmpData.ListDataID, dataConsole.ID)

	// 	}
	// 	if dataConsole.Status == 3 {
	// 		tmpData.ConsolidatedID = append(tmpData.ConsolidatedID, dataConsole.ID)
	// 	}
	// }

	jsonTmpData, err := json.Marshal(tmpData)
	if err != nil {
		return err
	}

	data := kafka.JsonData{}
	data.Name = "DUPLICATE"
	data.CompanyID = dataTB.CompanyID
	data.UserID = ctx.Auth.ID
	data.Data = string(jsonTmpData)
	data.Filter.Period = dataTB.Period
	data.Filter.Versions = dataTB.Versions

	jsonStr, err := json.Marshal(data)
	if err != nil {
		return err
	}

	go kafka.NewService("CONSOLIDATION").SendMessage("CONSOLIDATE", string(jsonStr))
	return nil
}
func (s *service) RequestToEditCombain(ctx *abstraction.Context, payload *dto.ConsolidationConsolidateRequest) error {
	if len(payload.ListConsolidation) == 0 && len(payload.ListConsolidationParent) == 0 {
		return response.CustomErrorBuilder(http.StatusBadRequest, "Data Anak Usaha Harus di isi", "Data Anak Usaha Harus di isi")
	}
	var tmpData DataToConsolidate
	dataTB, err := s.TrialBalanceRepository.FindByID(ctx, &payload.ConsolidationMasterID)
	if err != nil {
		return err
	}

	// if dataTB.Status != 2 {
	// 	return errors.New(fmt.Sprintf("data belum tervalidasi atau sudah terkonsolidasi"))
	// }
	// datePeriod, err := time.Parse(time.RFC3339, dataTB.Period)
	// if err != nil {
	// 	return err
	// }
	// period := datePeriod.Format("2006-01-02")
	// criteriaTB := model.ConsolidationFilterModel{}
	// criteriaTB.Period = &period
	// criteriaTB.CompanyID = &dataTB.CompanyID
	// criteriaTB.Versions = &dataTB.Versions
	dataTbIDconsole, err := s.Repository.FindByID(ctx, &payload.ConsolidationID)
	if err != nil {
		return err
	}

	tmpData.MasterID = dataTB.ID
	tmpData.ListDataID = payload.ListConsolidation
	tmpData.ConsolidatedID = dataTbIDconsole.ID
	tmpData.ListConsolidatedID = payload.ListConsolidationParent

	// for _, vDataconsolidation := range payload.ListConsolidation {
	// 	criteriaConsolidation := model.ConsolidationFilterModel{}
	// 	criteriaConsolidation.ConsolidationID = &vDataconsolidation

	// 	dataConsole, err := s.TrialBalanceRepository.FindByID(ctx, &vDataconsolidation)
	// 	if err != nil {
	// 		return err
	// 	}

	// 	if dataConsole.Status == 2 {
	// 		tmpData.ListDataID = append(tmpData.ListDataID, dataConsole.ID)

	// 	}
	// 	if dataConsole.Status == 3 {
	// 		tmpData.ConsolidatedID = append(tmpData.ConsolidatedID, dataConsole.ID)
	// 	}
	// }

	jsonTmpData, err := json.Marshal(tmpData)
	if err != nil {
		return err
	}
	data := kafka.JsonData{}
	data.Name = "EDIT_COMBINE"
	data.CompanyID = dataTB.CompanyID
	data.UserID = ctx.Auth.ID
	data.Data = string(jsonTmpData)
	data.Filter.Period = dataTB.Period
	data.Filter.Versions = dataTB.Versions

	jsonStr, err := json.Marshal(data)
	if err != nil {
		return err
	}

	go kafka.NewService("CONSOLIDATION").SendMessage("CONSOLIDATE", string(jsonStr))
	return nil
}
func (s *service) GetVersion(ctx *abstraction.Context, payload *dto.GetVersionRequest) (*dto.GetVersionResponse, error) {
	filter := model.ConsolidationFilterModel{
		CompanyCustomFilter: model.CompanyCustomFilter{
			CompanyID:          payload.CompanyID,
			ArrCompanyID:       payload.ArrCompanyID,
			ArrCompanyString:   payload.ArrCompanyString,
			ArrCompanyOperator: payload.ArrCompanyOperator,
		},
	}
	filter.ArrStatus = payload.ArrStatus
	filter.Period = payload.Period
	filter.Status = payload.Status
	data, err := s.Repository.GetVersion(ctx, &filter)
	if err != nil {
		return &dto.GetVersionResponse{}, response.ErrorBuilder(&response.ErrorConstant.InternalServerError, err)
	}
	result := &dto.GetVersionResponse{
		Data: *data,
	}
	return result, nil
}
func (s *service) FindListCompanyCreateNewCombine(ctx *abstraction.Context, payload *dto.FindListCompanyCreateNewCombineGetRequest) (*dto.FindListCompanyCreateNewCombineGetResponse, error) {
	filter := model.CompanyFilterModel{}
	filter = payload.CompanyFilterModel
	// if payload.ChildCompany != nil && *payload.ChildCompany == true {
	// 	filter.ParentCompanyID = nil
	// }
	data, err := s.Repository.FindListCompanyCreateNewCombine(ctx, &filter)
	if err != nil {
		return &dto.FindListCompanyCreateNewCombineGetResponse{}, helper.ErrorHandler(err)
	}
	result := &dto.FindListCompanyCreateNewCombineGetResponse{
		Datas: *data,
	}
	return result, nil
}
func (s *service) Delete(ctx *abstraction.Context, payload *dto.ConsolidationDeleteRequest) (*dto.ConsolidationDeleteResponse, error) {
	var data model.ConsolidationEntityModel
	var datas model.ConsolidationDetailEntityModel

	if err := trxmanager.New(s.Db).WithTrx(ctx, func(ctx *abstraction.Context) error {
		FindByConsolidationID, err := s.Repository.FindByID(ctx, &payload.ID)
		if err != nil {
			fmt.Println(err)
			return err
		}
		datePeriod, err := time.Parse(time.RFC3339, FindByConsolidationID.Period)
		if err != nil {
			return err
		}
		period := datePeriod.Format("2006-01-02")
		criteriaCB := model.ConsolidationBridgeFilterModel{}
		criteriaCB.CompanyID = &FindByConsolidationID.CompanyID
		criteriaCB.Versions = &FindByConsolidationID.Versions
		criteriaCB.ConsolidationVersions = &FindByConsolidationID.ConsolidationVersions
		criteriaCB.Period = &period

		cekConsolidationBridge, err := s.Repository.FindByConsolidationBridge(ctx, &criteriaCB)

		if err != nil && err != gorm.ErrRecordNotFound {
			fmt.Println(err)
			return err
		}

		if err == gorm.ErrRecordNotFound {
			data.Context = ctx
			_, err := s.Repository.DestroyDetail(ctx, &payload.ID, &datas)
			if err != nil {
				fmt.Println(err)
				return response.ErrorBuilder(&response.ErrorConstant.InternalServerError, err)
			}

			result, err := s.Repository.Destroy(ctx, &payload.ID, &data)
			if err != nil {
				fmt.Println(err)
				return response.ErrorBuilder(&response.ErrorConstant.InternalServerError, err)
			}
			data = *result
		}

		if cekConsolidationBridge.ID != 0 {
			consolidationVers := strconv.Itoa(cekConsolidationBridge.ConsolidationVersions)
			return response.CustomErrorBuilder(http.StatusBadRequest, "Tidak Dapat Menghapus Consolidation Karena Telah Di Consolidation Versions "+consolidationVers, "Tidak Dapat Menghapus Consolidation Karena Telah Di Consolidation Versions "+consolidationVers)
		}
		return nil
	}); err != nil {
		fmt.Println(err)
		return &dto.ConsolidationDeleteResponse{}, err
	}
	result := &dto.ConsolidationDeleteResponse{
		ConsolidationEntityModel: data,
	}
	return result, nil
}
func (s *service) FindByID(ctx *abstraction.Context, payload *dto.ConsolidationGetByIDRequest) (*dto.ConsolidationGetByIDResponse, error) {
	data, err := s.Repository.FindByID(ctx, &payload.ID)
	if err != nil {
		return &dto.ConsolidationGetByIDResponse{}, response.ErrorBuilder(&response.ErrorConstant.InternalServerError, err)
	}
	result := &dto.ConsolidationGetByIDResponse{
		ConsolidationEntityModel: *data,
	}
	return result, nil
}

func (s *service) GetControl(ctx *abstraction.Context, payload *dto.ConsolidationGetControlRequest) (*dto.ConsolidationGetControlResponse, error) {
	data, err := s.Repository.FindByID(ctx, &payload.ConsolidationID)
	if err != nil {
		return nil, helper.ErrorHandler(err)
	}

	allowed := helper.CompanyValidation(ctx.Auth.ID, data.CompanyID)
	if !allowed {
		return nil, response.ErrorBuilder(&response.ErrorConstant.UnauthorizedAccess, errors.New("not allowed"))
	}

	summaryJournal, err := s.Repository.FindSummaryJournal(ctx, &data.ID)
	if err != nil {
		return nil, helper.ErrorHandler(err)
	}

	control2WBS1, err := s.Repository.FindControl2Wbs1(ctx, &data.ID)
	if err != nil {
		return nil, helper.ErrorHandler(err)
	}
	tmp1 := *control2WBS1.AmountJpmCr - *summaryJournal.AmountJpmCr
	control2WBS1.AmountJpmCr = &tmp1
	tmp2 := *control2WBS1.AmountJpmDr - *summaryJournal.AmountJpmDr
	control2WBS1.AmountJpmDr = &tmp2
	tmp3 := *control2WBS1.AmountJcteCr - *summaryJournal.AmountJcteCr
	control2WBS1.AmountJcteCr = &tmp3
	tmp4 := *control2WBS1.AmountJcteDr - *summaryJournal.AmountJcteDr
	control2WBS1.AmountJcteDr = &tmp4
	tmp5 := *control2WBS1.AmountJelimCr - *summaryJournal.AmountJelimCr
	control2WBS1.AmountJelimCr = &tmp5
	tmp6 := *control2WBS1.AmountJelimDr - *summaryJournal.AmountJelimDr
	control2WBS1.AmountJelimDr = &tmp6

	return &dto.ConsolidationGetControlResponse{
		Datas: *control2WBS1,
	}, nil
}
