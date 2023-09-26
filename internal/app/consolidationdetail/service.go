package consolidationdetail

import (
	// "fmt"
	// "errors"
	"fmt"
	// "log"
	"mcash-finance-console-core/internal/abstraction"
	"mcash-finance-console-core/internal/dto"
	"mcash-finance-console-core/internal/factory"
	"mcash-finance-console-core/internal/model"
	"mcash-finance-console-core/internal/repository"

	// "sort"
	"sync"

	// "runtime"
	// "sort"
	// "sync"

	// "runtime"
	// "mcash-finance-console-core/pkg/util/helper"
	"mcash-finance-console-core/pkg/util/response"
	"strings"

	// "sync"

	"gorm.io/gorm"
)

type service struct {
	Repository              repository.ConsolidationDetail
	ConsolidationRepository repository.Consolidation
	Db                      *gorm.DB
}

type Service interface {
	Find(ctx *abstraction.Context, payload *dto.ConsolidationDetailGetRequest) (*dto.ConsolidationDetailGetResponse, error)
}

func NewService(f *factory.Factory) *service {
	repository := f.ConsolidationDetailRepository
	tbrepository := f.ConsolidationRepository
	db := f.Db
	return &service{
		Repository:              repository,
		Db:                      db,
		ConsolidationRepository: tbrepository,
	}
}

// func (s *service) FindAll(ctx *abstraction.Context, payload *dto.ConsolidationDetailGetByParentRequests) (*dto.ConsolidationDetailGetByParentResponses, error) {
// 	tb, err := s.ConsolidationRepository.FindByID(ctx, &payload.ConsolidationID)
// 	if err != nil {
// 		return nil, err
// 	}

// 	detailData, err := s.Repository.FindAllDetail(ctx, &payload.ConsolidationID)
// 	if err != nil {
// 		return nil, err
// 	}

// 	company2, err := s.Repository.FindAnakUsahaOnlyss(ctx, payload.ConsolidationID)
// 	if err != nil {
// 		return nil, err
// 	}

// 	var wg sync.WaitGroup
// 	detailChan := make(chan model.ConsolidationDetailFmtEntityModel, len(*detailData))

// 	for _, p := range *detailData {
// 		wg.Add(1)

// 		go func(p model.ConsolidationDetailFmtEntityModel) {
// 			defer wg.Done()
// 			payload.ParentCode = p.Code
// 			p.ConsolidationBridge = company2

// 			var ArrConsolidationss []model.ConsolidationBridgeEntityModel
// 			for _, c := range p.ConsolidationBridge {
// 				amount, err := s.Repository.FindAllDetailAmounts(ctx, &c.ID , &p.Code)
// 				if err != nil {
// 					log.Println(err)
// 					continue // skip the failed item
// 				}
// 				c.ConsolidationBridgeDetail = *amount
// 				ArrConsolidationss = append(ArrConsolidationss, c)
// 			}
// 			p.ConsolidationBridge = ArrConsolidationss
// 			detailChan <- p
// 		}(p)
// 	}
// 	wg.Wait()
// 	close(detailChan)

// 	var ArrConsolidations []model.ConsolidationDetailFmtEntityModel
// 	for p := range detailChan {
// 		ArrConsolidations = append(ArrConsolidations, p)
// 	}

// 	sort.Slice(ArrConsolidations, func(i, j int)bool {
// 		return ArrConsolidations[i].SortID < ArrConsolidations[j].SortID
// 	})

// 	// handle the case when the tree list cannot be created
// 	if len(ArrConsolidations) == 0 {
// 		return nil, errors.New("no consolidation details found")
// 	}

// 	tb.ConsolidationDetails = makeTreeList(ArrConsolidations, 0)

// 	result := &dto.ConsolidationDetailGetByParentResponses{
// 		Data: *tb,
// 	}
// 	return result, nil
// }

func (s *service) FindAll(ctx *abstraction.Context, payload *dto.ConsolidationDetailGetByParentRequests) (*dto.ConsolidationDetailGetByParentResponses, error) {
	tb, err := s.ConsolidationRepository.FindByID(ctx, &payload.ConsolidationID)
	if err != nil {
		return nil, err
	}

	// allowed := helper.CompanyValidation(ctx.Auth.ID, tb.CompanyID)
	// if !allowed {
	// 	return nil, response.ErrorBuilder(&response.ErrorConstant.UnauthorizedAccess, errors.New("Not Allowed"))
	// }

	detailData, err := s.Repository.FindAllDetail(ctx, &payload.ConsolidationID)
	if err != nil {
		return nil, err
	}

	var wg sync.WaitGroup

	var ArrConsolidations []model.ConsolidationDetailFmtEntityModel
	wg.Add(len(*detailData))
	for _, c := range *detailData {
			payload.ParentCode = c.Code
			// payload.ConsolidationDetailFilterModel.VersionConsolidation = &consolidation.ConsolidationVersions
			company, err := s.Repository.FindAnakUsahaOnlys(ctx, payload.ConsolidationID, payload.ParentCode)
			if err != nil {
				return nil, err
			}
			c.ConsolidationBridge = company
			ArrConsolidations = append(ArrConsolidations, c)
	}
	
	
	fmt.Println("Semua goroutine selesai dieksekusi")
	
	tb.ConsolidationDetails = makeTreeList(ArrConsolidations, 0)
	result := &dto.ConsolidationDetailGetByParentResponses{
		Data: *tb,
	}
	return result, nil
}
func makeTreeList(dataTB []model.ConsolidationDetailFmtEntityModel, parent int) []model.ConsolidationDetailFmtEntityModel {
	tbData := []model.ConsolidationDetailFmtEntityModel{}
	for _, v := range dataTB {
		if v.ParentID == parent {
			v.Children = makeTreeList(dataTB, v.FormatterDetailID)
			tbData = append(tbData, v)
		}
	}
	return tbData
}
func (s *service) Find(ctx *abstraction.Context, payload *dto.ConsolidationDetailGetRequest) (*dto.ConsolidationDetailGetResponse, error) {

	consolidation, err := s.ConsolidationRepository.FindByID(ctx, payload.ConsolidationID)
	if err != nil {
		return &dto.ConsolidationDetailGetResponse{}, response.ErrorBuilder(&response.ErrorConstant.InternalServerError, err)
	}
	source := "TB-CONSOLIDATION"

	mf, err := s.Repository.FindByFormatter(ctx, &source)
	if err != nil {
		return &dto.ConsolidationDetailGetResponse{}, response.ErrorBuilder(&response.ErrorConstant.InternalServerError, err)
	}

	mfd, err := s.Repository.FindByFormatterDetail(ctx, &mf.ID, payload.Parent)

	if err != nil {

		source := "ASET"
		payload.Parent = &source
		codeFormatter, err := s.Repository.FindByFormatterDetail(ctx, &mf.ID, payload.Parent)
		if err != nil {
			return &dto.ConsolidationDetailGetResponse{}, response.ErrorBuilder(&response.ErrorConstant.InternalServerError, err)

		}

		sr := "TOTAL_"
		sr2 := sr + codeFormatter.Code
		payload.Sourcea = &sr2

		data, info, err := s.Repository.FindisNull(ctx, &payload.ConsolidationDetailFilterModel, &payload.Pagination)
		if err != nil {
			return &dto.ConsolidationDetailGetResponse{}, response.ErrorBuilder(&response.ErrorConstant.InternalServerError, err)
		}

		// var ArrConsolidations []model.ConsolidationDetailEntityModel
		// payload.ConsolidationDetailFilterModel.VersionConsolidation = &consolidation.ConsolidationVersions

		// _, err = s.Repository.FindAnakUsahaOnly(ctx, &payload.ConsolidationDetailFilterModel)
		// 	if err != nil {
		// 		return nil, err
		// 	}

		// for _, c := range data {
		// 	// var ss int
		// 	ConsolidationDetail := model.ConsolidationDetailEntityModel{
		// 		Context: ctx,
		// 	}
		// 	ConsolidationDetail.ID = c.ID
		// 	ConsolidationDetail.ConsolidationDetailEntity = c.ConsolidationDetailEntity

		// 	payload.ConsolidationDetailFilterModel.Code = &c.Code
		// 	payload.ConsolidationDetailFilterModel.VersionConsolidation = &consolidation.ConsolidationVersions
		// 	// ConsolidationDetail.ConsolidationBridge = company
		// 	// for _, d := range company {
		// 	// 	// var ss int
		// 	// 	ConsolidationBridge := model.ConsolidationBridgeEntityModel{
		// 	// 		Context: ctx,
		// 	// 	}
		// 	// 	ConsolidationBridge.ID = d.ID
		// 	// 	ConsolidationBridge.ConsolidationBridgeEntity = d.ConsolidationBridgeEntity

		// 	// 	payload.ConsolidationDetailFilterModel.Code = &c.Code
		// 	// 	payload.ConsolidationDetailFilterModel.ConsolidationBridgeID = &d.ID

		// 	// 	amount, err := s.Repository.FindByAmount(ctx, &d.ID, &c.Code)
		// 	// 	if err != nil {
		// 	// 		continue
		// 	// 	}
		// 	// 	ConsolidationBridge.ConsolidationBridgeDetail = *amount
		// 	// 	ConsolidationDetail.ConsolidationBridge = append(ConsolidationDetail.ConsolidationBridge, ConsolidationBridge)
		// 	// 	// ArrConsolidations = append(ArrConsolidations, ConsolidationDetail)
		// 	// }
		// 	ArrConsolidations = append(ArrConsolidations, ConsolidationDetail)

		// }
		// consolidation.ConsolidationDetail = ArrConsolidations

		// result := &dto.ConsolidationDetailGetResponse{
		// 	Datas:          *consolidation,
		// 	// ChildCompany: company,
		// 	PaginationInfo: *info,
		// }
		// return result, nil
		var ArrConsolidations []model.ConsolidationDetailEntityModel

		payload.ConsolidationDetailFilterModel.VersionConsolidation = &consolidation.ConsolidationVersions

		for _, c := range data {
			// var ss int
			ConsolidationDetail := model.ConsolidationDetailEntityModel{
				Context: ctx,
			}
			ConsolidationDetail.ID = c.ID
			ConsolidationDetail.ConsolidationDetailEntity = c.ConsolidationDetailEntity

			payload.ConsolidationDetailFilterModel.Code = &c.Code
			payload.ConsolidationDetailFilterModel.VersionConsolidation = &consolidation.ConsolidationVersions

			company, err := s.Repository.FindAnakUsahaOnlyNull(ctx, &payload.ConsolidationDetailFilterModel)
			if err != nil {
				return nil, err
			}

			ConsolidationDetail.ConsolidationBridge = company

			// for _, d := range company {
			// 	// var ss int
			// 	ConsolidationBridges := model.ConsolidationBridgeEntityModel{
			// 		Context: ctx,
			// 	}
			// 	ConsolidationBridges.ID = d.ID
			// 	ConsolidationBridges.ConsolidationBridgeEntity = d.ConsolidationBridgeEntity

			// 	payload.ConsolidationDetailFilterModel.Code = &c.Code
			// 	// payload.ConsolidationDetailFilterModel.ConsolidationBridgeID = &d.ID

			// 	amount, err := s.Repository.FindByAmount(ctx, &d.ID, &c.Code)
			// 	if err != nil {
			// 		continue
			// 	}
			// 	ConsolidationBridges.ConsolidationBridgeDetail = *amount
			// 	ConsolidationDetail.ConsolidationBridge = append(ConsolidationDetail.ConsolidationBridge, ConsolidationBridges)
			// 	// ArrConsolidations = append(ArrConsolidations, ConsolidationDetail)
			// }
			ArrConsolidations = append(ArrConsolidations, ConsolidationDetail)

		}
		consolidation.ConsolidationDetail = ArrConsolidations
		company, err := s.Repository.FindAnakUsahaOnly(ctx, &payload.ConsolidationDetailFilterModel)
		if err != nil {
			return nil, err
		}
		result := &dto.ConsolidationDetailGetResponse{
			Datas:          *consolidation,
			ChildCompany:   company,
			PaginationInfo: *info,
		}
		return result, nil
	}

	payload.ParentID = &mfd.ID
	code, _, err := s.Repository.FindCode(ctx, &payload.ConsolidationDetailFilterModel, &payload.Pagination)
	if err != nil {
		return &dto.ConsolidationDetailGetResponse{}, response.ErrorBuilder(&response.ErrorConstant.InternalServerError, err)
	}
	var codeSource []string

	// var ArrConsolidation []model.ConsolidationDetailEntityModel

	for _, v := range code {
		var ss string

		sr := "TOTAL_"
		sr2 := sr + v.Code
		ss = sr2
		codeSource = append(codeSource, ss)
		// ArrConsolidation = append(ArrConsolidation, AdjustmentDetail)
	}
	// sr := "TOTAL_"
	replace := "'"
	payload.Source = &codeSource
	payload.ParentID = &mfd.ID
	justString := strings.Join(codeSource, " , ")
	testString := strings.ReplaceAll(justString, " ", replace)
	payload.Sourcea = &testString

	data, info, err := s.Repository.Find(ctx, &payload.ConsolidationDetailFilterModel, &payload.Pagination)
	if err != nil {
		return &dto.ConsolidationDetailGetResponse{}, response.ErrorBuilder(&response.ErrorConstant.InternalServerError, err)
	}

	payload.ConsolidationDetailFilterModel.VersionConsolidation = &consolidation.ConsolidationVersions

	var ArrConsolidations []model.ConsolidationDetailEntityModel
	for _, c := range data {
		// var ss int
		ConsolidationDetail := model.ConsolidationDetailEntityModel{
			Context: ctx,
		}
		ConsolidationDetail.ID = c.ID
		ConsolidationDetail.ConsolidationDetailEntity = c.ConsolidationDetailEntity

		payload.ConsolidationDetailFilterModel.Code = &c.Code
		payload.ConsolidationDetailFilterModel.VersionConsolidation = &consolidation.ConsolidationVersions

		company, err := s.Repository.FindAnakUsahaOnly(ctx, &payload.ConsolidationDetailFilterModel)
		if err != nil {
			return nil, err
		}

		ConsolidationDetail.ConsolidationBridge = company

		// for _, d := range company {
		// 	// var ss int
		// 	ConsolidationBridges := model.ConsolidationBridgeEntityModel{
		// 		Context: ctx,
		// 	}
		// 	ConsolidationBridges.ID = d.ID
		// 	ConsolidationBridges.ConsolidationBridgeEntity = d.ConsolidationBridgeEntity

		// 	payload.ConsolidationDetailFilterModel.Code = &c.Code
		// 	// payload.ConsolidationDetailFilterModel.ConsolidationBridgeID = &d.ID

		// 	amount, err := s.Repository.FindByAmount(ctx, &d.ID, &c.Code)
		// 	if err != nil {
		// 		continue
		// 	}
		// 	ConsolidationBridges.ConsolidationBridgeDetail = *amount
		// 	ConsolidationDetail.ConsolidationBridge = append(ConsolidationDetail.ConsolidationBridge, ConsolidationBridges)
		// 	// ArrConsolidations = append(ArrConsolidations, ConsolidationDetail)
		// }
		ArrConsolidations = append(ArrConsolidations, ConsolidationDetail)

	}
	consolidation.ConsolidationDetail = ArrConsolidations
	company, err := s.Repository.FindAnakUsahaOnly(ctx, &payload.ConsolidationDetailFilterModel)
	if err != nil {
		return nil, err
	}
	result := &dto.ConsolidationDetailGetResponse{
		Datas:          *consolidation,
		ChildCompany:   company,
		PaginationInfo: *info,
	}
	return result, nil
}

func (s *service) FindByParent(ctx *abstraction.Context, payload *dto.ConsolidationDetailGetByParentRequest) (*dto.ConsolidationDetailGetByParentResponse, error) {

	detailData, err := s.Repository.FindDetail(ctx, &payload.TrialBalanceID, &payload.ParentID)
	if err != nil {
		return nil, err
	}
	var ArrConsolidations []model.ConsolidationDetailFmtEntityModel

	// payload.ConsolidationDetailFilterModel.VersionConsolidation = &consolidation.ConsolidationVersions
	consolidation, err := s.ConsolidationRepository.FindByID(ctx, payload.ConsolidationID)
	if err != nil {
		return &dto.ConsolidationDetailGetByParentResponse{}, response.ErrorBuilder(&response.ErrorConstant.InternalServerError, err)
	}
	for _, c := range *detailData {

		payload.ConsolidationDetailFilterModel.Code = &c.Code
		payload.ConsolidationDetailFilterModel.VersionConsolidation = &consolidation.ConsolidationVersions
		company, err := s.Repository.FindAnakUsahaOnly(ctx, &payload.ConsolidationDetailFilterModel)
		if err != nil {
			return nil, err
		}

		c.ConsolidationBridge = company
		ArrConsolidations = append(ArrConsolidations, c)

	}
	// consolidation.ConsolidationDetail = ArrConsolidations
	result := &dto.ConsolidationDetailGetByParentResponse{
		Data: &ArrConsolidations,
	}
	return result, nil
}
