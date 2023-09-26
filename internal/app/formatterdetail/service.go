package formatterdetail

import (
	"errors"
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

	"github.com/xuri/excelize/v2"
	"gorm.io/gorm"
)

type service struct {
	Repository          repository.FormatterDetailDev
	FormatterRepository repository.Formatter
	CoaDevRepository    repository.CoaDev
	Db                  *gorm.DB
}

type Service interface {
	Find(ctx *abstraction.Context, payload *dto.FormatterDetailDevGetRequest) (*dto.FormatterDetailDevGetResponse, error)
}

func NewService(f *factory.Factory) *service {
	repository := f.FormatterDetailDevRepository
	coadevRepository := f.CoaDevRepository
	db := f.Db
	return &service{
		Repository:    repository,
		CoaDevRepository: coadevRepository,
		Db:            db,
	}
}

func (s *service) FindTree(ctx *abstraction.Context) (*dto.FormatterDetailDevDetailGetByParentRequests, error) {
	filter := model.FormatterDetailDevFilterModel{}

	pagination := abstraction.Pagination{}
	nolimit := 10000
	pagination.PageSize = &nolimit
	data, _, err := s.Repository.Finds(ctx, &filter, &pagination)
	if err != nil {
		return nil, helper.ErrorHandler(err)
	}

	// Check if data is empty
	if len(*data) == 0 {
		return nil, errors.New("no data found")
	}

	datas := makeTreeList1(*data, 0)
	result := &dto.FormatterDetailDevDetailGetByParentRequests{
		Datas: datas,
	}
	return result, nil
}

func makeTreeList1(dataTB []model.FormatterDetailDevFmtEntityModel, parent int) []model.FormatterDetailDevFmtEntityModel {
	tbData := []model.FormatterDetailDevFmtEntityModel{}
	for _, v := range dataTB {
		if v.ParentID == parent {
			v.Children = makeTreeList1(dataTB, v.ID)
			tbData = append(tbData, v)
		}
	}
	return tbData
}
func makeTreeList2(dataTB []model.FormatterDetailDevFmtEntityModel, parent int) []model.FormatterDetailDevFmtEntityModel {
	tbData := []model.FormatterDetailDevFmtEntityModel{}
	for _, v := range dataTB {
		if v.ParentID == parent {
			v.Children = makeTreeList2(dataTB, v.ID)
			tbData = append(tbData, v)

		}
	}
	return tbData
}
func (s *service) DragAndDrop(ctx *abstraction.Context, payload *dto.FormatterDragAndDropDevRequest) (*dto.FormatterDragAndDropDevResponse, error) {
	filter := model.FormatterDetailDevFilterModel{}
	pagination := abstraction.Pagination{}
	nolimit := 10000
	pagination.PageSize = &nolimit
	data, _, err := s.Repository.FindWithTotal(ctx, &filter, &pagination)
	if err != nil {
		return nil, helper.ErrorHandler(err)
	}
	if len(*data) == 0 {
		return nil, errors.New("tidak ditemukan data")
	}

	// Convert data to tree structure
	datas := makeTreeList1(*data, payload.ParentID)
	valueDatas := make(map[int]model.FormatterDetailDevFmtEntityModel)

	var datasToArray []model.FormatterDetailDevFmtEntityModel
	for _, d := range datas {
		valueDatas[d.ID] = d
	}
	for _, dc := range datas {
		datasToArray = append(datasToArray, dc)
	}
	var datas2 []model.FormatterDetailDevFmtEntityModel
	for _, p := range payload.Datas {
		if _, ok := valueDatas[p.ID]; ok {
			datas2 = append(datas2, p)
			for i, dp := range datasToArray {
				if dp.ID == p.ID {
					dp = datasToArray[i+1]
					datas2 = append(datas2, dp)
				}
			}
		}
	}
	var dataPayload []model.FormatterDetailDevFmtEntityModel
	for _, dP := range datas2 {
		da := makeTreeList2(*data, dP.ID)
		dP.Children = da
		dataPayload = append(dataPayload, dP)
	}

	var allSortIDs []float64
	for _, tree := range datas {
		gatherSortIDs(&tree, &allSortIDs)
	}
	if err := s.updateSortIDInDatabase(ctx, dataPayload, allSortIDs); err != nil {
		return nil, err
	}

	result := &dto.FormatterDragAndDropDevResponse{
		Datas: datas,
	}
	return result, nil
}
func gatherSortIDs(node *model.FormatterDetailDevFmtEntityModel, sortIDs *[]float64) {
	*sortIDs = append(*sortIDs, node.SortID)
	for _, child := range node.Children {
		gatherSortIDs(&child, sortIDs)
	}
}
func updateSortID(nodes []model.FormatterDetailDevFmtEntityModel, datas []model.FormatterDetailDevFmtEntityModel) {
	for i := 0; i < len(nodes) && i < len(datas); i++ {
		nodes[i].SortID = datas[i].SortID
		if len(nodes[i].Children) > 0 && len(datas[i].Children) > 0 {
			updateSortID(nodes[i].Children, datas[i].Children)
		}
	}
}
func (s *service) updateSortIDInDatabase(ctx *abstraction.Context, datas []model.FormatterDetailDevFmtEntityModel, sortIDs []float64) error {
	
	index := 0
	for _, data := range datas {
		// data := &datas[i]
		find, err := s.Repository.FindByID(ctx, &data.ID)
		if err != nil {
			return err
		}
		find.SortID = sortIDs[index]
		
		_, err = s.Repository.Update(ctx, &find.ID, find)
		if err != nil {
			return err
		}
		if len(data.Children) > 0 {
			
			index = s.child(ctx, data.Children, index+1, sortIDs)
		}else {
			index++
		}
	}

	return nil
}

func (s *service) child(ctx *abstraction.Context, dataTB []model.FormatterDetailDevFmtEntityModel, parent int, sortIds []float64) int {
	i := parent

	for _, v := range dataTB {
		find, err := s.Repository.FindByID(ctx, &v.ID)
		if err != nil {
			continue
		}
		find.SortID = sortIds[i]
		_, err = s.Repository.Update(ctx, &find.ID, find)
		if err != nil {
			continue
		}

		if len(v.Children) > 0 {
			i = s.child(ctx, v.Children, i+1, sortIds)
		}
		if len(v.Children) < 1 {
			i++
		}
	}

	return i
}

func (s *service) Create(ctx *abstraction.Context, payload *dto.FormatterDetailDevCreateRequest) (*dto.FormatterDetailDevCreateResponse, error) {
	var data model.FormatterDetailDevEntityModel
	fmtID := 3
	criteriaF := model.FormatterDetailDevFilterModel{}
	criteriaF.ParentID = payload.ParentID
	criteriaF.FormatterDevID = &fmtID

	findByParent, err := s.Repository.FindWithCriteria(ctx, &criteriaF)
	if err != nil {
		return nil, err
	}
	max := 0.0
	for _, parent := range *findByParent {
		if parent.SortID > max {
			max = parent.SortID
		}
	}
	findByID, err := s.Repository.FindByID(ctx, payload.ParentID)
	if err != nil {
		return nil, err
	}
	if max == 0 {
		max = findByID.SortID + 0.01
	}
	var subCoaWithTotal []string
	// t := true
	totalSubCoa := "TOTAL "
	subCoaTotal := totalSubCoa + payload.Description
	//isTotal == true
	if *payload.IsTotal == true {
		subCoa := payload.Description
		subCoa = strings.Replace(strings.ToUpper(subCoa), "+", "#", -1)
		subCoa = strings.Replace(strings.ToUpper(subCoa), " ", "_", -1)
		subCoa = strings.Replace(strings.ToUpper(subCoa), "-", "~", -1)

		subCoaWithTotal = append(subCoaWithTotal, subCoa)

		subCoaTotal = strings.Replace(strings.ToUpper(subCoaTotal), "+", "#", -1)
		subCoaTotal = strings.Replace(strings.ToUpper(subCoaTotal), " ", "_", -1)
		subCoaTotal = strings.Replace(strings.ToUpper(subCoaTotal), "-", "~", -1)
		subCoaWithTotal = append(subCoaWithTotal, subCoaTotal)

		findByCode, err := s.Repository.FindByCode(ctx, &fmtID, &subCoaTotal)
		if err != nil {
			return nil, err
		}
		bilanganKosong := 0

		if findByCode.ID != bilanganKosong {
			return nil, response.CustomErrorBuilder(http.StatusBadRequest, "Coa Sudah Terdaftar di Formatter "+subCoa, "Coa Sudah Terdaftar di Formatter "+subCoa)
		}

	} else {
		subCoa := payload.Description
		subCoa = strings.Replace(strings.ToUpper(subCoa), "+", "#", -1)
		subCoa = strings.Replace(strings.ToUpper(subCoa), " ", "_", -1)
		subCoa = strings.Replace(strings.ToUpper(subCoa), "-", "~", -1)
		subCoaWithTotal = append(subCoaWithTotal, subCoa)

		findByCode, err := s.Repository.FindByCode(ctx, &fmtID, &subCoa)
		if err != nil {
			return nil, err
		}
		bilanganKosong := 0

		if findByCode.ID != bilanganKosong {
			return nil, response.CustomErrorBuilder(http.StatusBadRequest, "Coa Sudah Terdaftar di Formatter "+subCoa, "Coa Sudah Terdaftar di Formatter "+subCoa)
		}
	}

	if err := trxmanager.New(s.Db).WithTrx(ctx, func(ctx *abstraction.Context) error {
		for _, v := range subCoaWithTotal {
			f := false
			t := true
			plusSatu := 1 + *findByID.Level
			coa := v
			coa = strings.Replace(strings.ToUpper(coa), "#", "+", -1)
			coa = strings.Replace(strings.ToUpper(coa), "_", " ", -1)
			coa = strings.Replace(strings.ToUpper(coa), "~", "-", -1)
			
			
				var data model.FormatterDetailDevEntityModel
				stringKosong := ""
				// data.Context = ctx
				data.FormatterDetailDevEntity = payload.FormatterDetailDevEntity
				data.FormatterDetailDevEntity.Code = v
				data.FormatterDetailDevEntity.Description = coa
				data.FormatterDetailDevEntity.IsControl = &f
				data.FormatterDetailDevEntity.SortID = max + 0.1
				data.FormatterDetailDevEntity.IsCoa = &f
				data.FormatterDetailDevEntity.AutoSummary = &f
				data.FormatterDetailDevEntity.IsTotal = &f
				data.FormatterDetailDevEntity.IsControl = &f
				data.FormatterDetailDevEntity.IsLabel = &t
				data.FormatterDetailDevEntity.ControlFormula = stringKosong
				data.FormatterDetailDevEntity.FormatterDevID = 3
				data.FormatterDetailDevEntity.ShowGroupCoa = &f
				data.FormatterDetailDevEntity.ParentID = payload.ParentID
				data.FormatterDetailDevEntity.IsParent = &t
				data.FormatterDetailDevEntity.IsShowView = &t
				data.FormatterDetailDevEntity.IsShowExport = &t
				data.FormatterDetailDevEntity.IsRecalculate = &f
				data.FormatterDetailDevEntity.CoaGroupID = findByID.CoaGroupID
				data.FormatterDetailDevEntity.Level = &plusSatu

				if v == subCoaTotal {
					data.FormatterDetailDevEntity.SortID = max + 0.2
					data.FormatterDetailDevEntity.IsTotal = &t
					data.FormatterDetailDevEntity.IsLabel = &f
					data.FormatterDetailDevEntity.IsParent = &f
					data.FormatterDetailDevEntity.Level = &plusSatu
				}
				if len(subCoaWithTotal) == 1 {
					data.FormatterDetailDevEntity.SortID = max + 0.01
				}
				result, err := s.Repository.Create(ctx, &data)
				if err != nil {
					return helper.ErrorHandler(err)
				}
				data = *result

		}

		return nil
	}); err != nil {
		return &dto.FormatterDetailDevCreateResponse{}, err
	}

	result := &dto.FormatterDetailDevCreateResponse{
		FormatterDetailDevEntityModel: data,
	}
	return result, nil
}
func (s *service) GetCoa(ctx *abstraction.Context, payload *dto.FormatterDetailDevGetCoaRequest) (*dto.FormatterDetailDevGetCoaResponse, error) {
   
    filterModel := &model.FormatterDetailDevFilterModel{}
    filterModel.FormatterDevID = payload.FormatterDevID

    pagination := abstraction.Pagination{}
    nolimit := 10000
    pagination.PageSize = &nolimit

    coaList, err := s.Repository.FindCoaListing(ctx)
    if err != nil {
        return nil, err
    }

 
    var response dto.FormatterDetailDevGetCoaResponse
    for _, coa := range *coaList {
        coaDTO := model.CoaDevEntityModel{
            CoaDevEntity: model.CoaDevEntity{
                Code: coa.Code,
            },
            
        }
        response.Datas = append(response.Datas, coaDTO)
    }

  
    response.Pagination = abstraction.Pagination{
        // Page:        paginationInfo.Pagination.Page,
        // PageSize:    paginationInfo.Pagination.PageSize,
        // Sort:        paginationInfo.Pagination.Sort,
        // SortBy:      paginationInfo.Pagination.SortBy,
    }

    return &response, nil
}
func (s *service) CreateCoaFmt(ctx *abstraction.Context, payload *dto.FormatterDetailDevCreateRequest) (*dto.FormatterDetailDevCreateResponse, error) {

	var dataCoa model.CoaDevEntityModel

	findByCodeCoa, err := s.CoaDevRepository.FindWithCode(ctx, &payload.Code)
	if err != nil {
		return nil, err
	}

	var data model.FormatterDetailDevEntityModel

	fmtID := 3
	criteriaF := model.FormatterDetailDevFilterModel{}
	criteriaF.ParentID = payload.ParentID
	criteriaF.FormatterDevID = &fmtID

	findByParent, err := s.Repository.FindWithCriteria(ctx, &criteriaF)
	if err != nil {
		return nil, err
	}
	max := 0.0
	for _, parent := range *findByParent {
		if parent.SortID > max {
			max = parent.SortID
		}
	}
	findByID, err := s.Repository.FindByID(ctx, payload.ParentID)
	if err != nil {
		return nil, err
	}
	if max == 0 {
		max = findByID.SortID
	}
	str := payload.Code
	sixChars := str[:6]

	findByCode, err := s.Repository.FindByCode(ctx, &fmtID, &sixChars)
	if err != nil {
		return nil, err
	}
	bilanganKosong := 0
	// if findByCode.ID != bilanganKosong {
	// 	return nil, response.CustomErrorBuilder(http.StatusBadRequest, "Coa Sudah Terdaftar di Formatter "+sixChars, "Coa Sudah Terdaftar di Formatter "+sixChars)
	// }
	if err := trxmanager.New(s.Db).WithTrx(ctx, func(ctx *abstraction.Context) error {
		plusSatu := 1 + *findByID.Level
		if len(*findByCodeCoa) == 0 {
			dataCoa.Context = ctx
			dataCoa.Code = payload.Code
			dataCoa.Name = payload.Description
			dataCoa.CoaTypeID = *payload.CoaTypeID
			_, err := s.CoaDevRepository.Create(ctx, &dataCoa)
			if err != nil {
				return helper.ErrorHandler(err)
			}

		}
		if findByCode.ID == bilanganKosong {
			f := false
			t := true
			stringKosong := ""
			data.Context = ctx
			data.FormatterDetailDevEntity = payload.FormatterDetailDevEntity
			data.FormatterDetailDevEntity.Code = sixChars
			data.FormatterDetailDevEntity.Description = sixChars
			data.FormatterDetailDevEntity.IsControl = &f
			data.FormatterDetailDevEntity.SortID = max + 0.1
			data.FormatterDetailDevEntity.IsCoa = &t
			data.FormatterDetailDevEntity.AutoSummary = &t
			data.FormatterDetailDevEntity.FxSummary = stringKosong
			data.FormatterDetailDevEntity.IsTotal = &f
			data.FormatterDetailDevEntity.IsControl = &f
			data.FormatterDetailDevEntity.IsLabel = &f
			data.FormatterDetailDevEntity.ControlFormula = stringKosong
			data.FormatterDetailDevEntity.FormatterDevID = 3
			data.FormatterDetailDevEntity.ShowGroupCoa = &f
			data.FormatterDetailDevEntity.ParentID = payload.ParentID
			data.FormatterDetailDevEntity.SummaryCoa = &stringKosong
			data.FormatterDetailDevEntity.IsParent = &f
			data.FormatterDetailDevEntity.IsShowView = &t
			data.FormatterDetailDevEntity.IsShowExport = &t
			data.FormatterDetailDevEntity.IsRecalculate = &f
			data.FormatterDetailDevEntity.CoaGroupID = findByID.CoaGroupID
			data.FormatterDetailDevEntity.Level = &plusSatu

			result, err := s.Repository.Create(ctx, &data)
			if err != nil {
				return helper.ErrorHandler(err)
			}
			data = *result
		}

		return nil
	}); err != nil {
		return &dto.FormatterDetailDevCreateResponse{}, err
	}

	result := &dto.FormatterDetailDevCreateResponse{
		FormatterDetailDevEntityModel: data,
	}
	return result, nil
}
func collectIDs(t model.FormatterDetailDevFmtEntityModel, ids *[]int) {
	*ids = append(*ids, t.ID)
	for _, child := range t.Children {
		collectIDs(child, ids)
	}
}
func (s *service) Delete(ctx *abstraction.Context, payload *dto.FormatterDetailDevDeleteRequest) (*dto.FormatterDetailDevDeleteResponse, error) {
	var data *[]model.FormatterDetailDevFmtEntityModel
	filter := model.FormatterDetailDevFilterModel{}
	pagination := abstraction.Pagination{}
	nolimit := 10000
	pagination.PageSize = &nolimit
	data, _, err := s.Repository.Finds(ctx, &filter, &pagination)
	if err != nil {
		return nil, helper.ErrorHandler(err)
	}
	if len(*data) == 0 {
		return nil, errors.New("tidak ditemukan data")
	}
	datas := makeTreeList1(*data, payload.ID)
	var allID []int
	for _, d := range datas {
		collectIDs(d, &allID)
	}
	allID = append(allID, payload.ID)
	findByID, err := s.Repository.FindByID(ctx, &payload.ID)
	if err != nil {
		return nil, err
	}
	total := "TOTAL_"
	totals := total+findByID.Code
	fmt := 3
	findByCode, err := s.Repository.FindByCode(ctx,&fmt, &totals)
	if err != nil {
		return nil, err
	}
	allID = append(allID, findByCode.ID)
	if err := trxmanager.New(s.Db).WithTrx(ctx, func(ctx *abstraction.Context) error {

		for _, id := range allID {
			if id == 0 {
				continue
			}
			datas, err := s.Repository.FindByID(ctx, &id)
			if err != nil {
				return helper.ErrorHandler(err)
			}

			datas.Context = ctx
			_, err = s.Repository.Delete(ctx, &datas.ID, datas)
			if err != nil {
				return helper.ErrorHandler(err)
			}
		}

		return nil
	}); err != nil {
		return &dto.FormatterDetailDevDeleteResponse{}, err
	}
	result := &dto.FormatterDetailDevDeleteResponse{
		// MutasiFaEntityModel: data,
	}
	return result, nil
}

func (s *service) Export(ctx *abstraction.Context) (*string, error) {
	
	f := excelize.NewFile()
	currentSheet := f.GetSheetName(f.GetActiveSheetIndex())
	f.DeleteSheet(currentSheet)

	errorMessages := []string{}
	
		a, err := s.ExportTrialBalance(ctx, f)
		if err != nil {
			return nil, err
		}
		f = a
	
		
	if len(errorMessages) > 0 {
		return nil, response.ErrorBuilder(&response.ErrorConstant.InternalServerError, errors.New(strings.Join(errorMessages, ", ")))
	}

	tmpFolder := fmt.Sprintf("assets/%d",1)
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
	
	fileName := "testFormatterDev.xlsx"
	fileLoc := fmt.Sprintf("assets/%d/%s", 1, fileName)
	err = f.SaveAs(fileLoc)
	if err != nil {
		return nil, err
	}

	return &fileLoc, nil
}
var tbRowCode = make(map[string]int)
func (s *service) ExportTrialBalance(ctx *abstraction.Context, f *excelize.File) (*excelize.File, error) {
	sheet := "TRIAL_BALANCE"
	indexSheet := f.NewSheet(sheet)
	f.SetActiveSheet(indexSheet)
	var (
		criteriaFormatter model.FormatterDetailDevFilterModel
	)
	formatterID := 3
	t := true
	criteriaFormatter.IsShowExport = &t
	criteriaFormatter.FormatterDevID = &formatterID
	pagesize := 100000
	tmpStr := "sort_id"
	tmpStr1 := "ASC"
	paginationTB := abstraction.Pagination{
		PageSize: &pagesize,
		SortBy:   &tmpStr,
		Sort:     &tmpStr1,
	}

	data, _, err := s.Repository.Find(ctx, &criteriaFormatter, &paginationTB)
	if err != nil {
		return nil, err
	}
	f.SetCellValue(sheet, "B2", "Company")
	f.SetCellValue(sheet, "D2", ": ")
	f.SetCellValue(sheet, "B3", "Date")
	f.SetCellValue(sheet, "D3", ": ")
	f.SetCellValue(sheet, "B4", "Subject")
	f.SetCellValue(sheet, "D4", ": DETAIL ASET, LIABILITAS & EKUITAS")

	arrStyleColWidth := []map[string]interface{}{
		{"COL": "A", "WIDTH": 0.83},
		{"COL": "B", "WIDTH": 15.38},
		{"COL": "C", "WIDTH": 2.14},
		{"COL": "D", "WIDTH": 2.14},
		{"COL": "E", "WIDTH": 57.45},
		{"COL": "F", "WIDTH": 6.43},
		{"COL": "G", "WIDTH": 17.65},
		{"COL": "H", "WIDTH": 10.71},
		{"COL": "I", "WIDTH": 16.83},
		{"COL": "J", "WIDTH": 10.10},
		{"COL": "K", "WIDTH": 17.65},
		{"COL": "L", "WIDTH": 22.14},
	}
	for _, v := range arrStyleColWidth {
		tmpColWidth := fmt.Sprintf("%f", v["WIDTH"])
		colWidth, _ := strconv.ParseFloat(tmpColWidth, 64)
		err = f.SetColWidth(sheet, fmt.Sprintf("%s", v["COL"]), fmt.Sprintf("%s", v["COL"]), colWidth)
		if err != nil {
			return nil, err
		}
	}

	styleDefault, err := f.NewStyle(&excelize.Style{
		Font: &excelize.Font{
			Family: "Arial",
			Size:   10,
		},
	})
	if err != nil {
		return nil, err
	}
	err = f.SetColStyle(sheet, "A:Z", styleDefault)
	if err != nil {
		return nil, err
	}
	numberFormat := "#,##"
	stylingBorderRightOnly, _ := f.NewStyle(&excelize.Style{
		Border: []excelize.Border{
			{Type: "right", Color: "000000", Style: 1},
		},
	})

	stylingBorderLeftOnly, _ := f.NewStyle(&excelize.Style{
		Border: []excelize.Border{
			{Type: "left", Color: "000000", Style: 1},
		},
	})

	stylingBorderLROnly, _ := f.NewStyle(&excelize.Style{
		Border: []excelize.Border{
			{Type: "right", Color: "000000", Style: 1},
			{Type: "left", Color: "000000", Style: 1},
		},
	})

	stylingBorderTopOnly, err := f.NewStyle(&excelize.Style{
		Border: []excelize.Border{
			{Type: "top", Color: "000000", Style: 1},
		},
	})
	if err != nil {
		return nil, err
	}

	stylingHeader, err := f.NewStyle(&excelize.Style{
		Font: &excelize.Font{
			Bold: true,
		},
		Border: []excelize.Border{
			{Type: "top", Color: "000000", Style: 1},
			{Type: "bottom", Color: "000000", Style: 1},
			{Type: "left", Color: "000000", Style: 1},
			{Type: "right", Color: "000000", Style: 1},
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
	if err != nil {
		return nil, err
	}

	
	stylingSubTotal, err := f.NewStyle(&excelize.Style{
		Border: []excelize.Border{
			{Type: "top", Color: "000000", Style: 1},
			{Type: "bottom", Color: "000000", Style: 1},
		},
		Font: &excelize.Font{
			Bold: true,
		},
	})
	if err != nil {
		return nil, err
	}

	stylingSubTotalCurrency, err := f.NewStyle(&excelize.Style{
		Border: []excelize.Border{
			{Type: "top", Color: "000000", Style: 1},
			{Type: "bottom", Color: "000000", Style: 1},
		},
		Font: &excelize.Font{
			Bold: true,
		},
		NumFmt: 41,
	})
	if err != nil {
		return nil, err
	}

	stylingTotal, err := f.NewStyle(&excelize.Style{
		Border: []excelize.Border{
			{Type: "top", Color: "000000", Style: 1},
			{Type: "bottom", Color: "000000", Style: 1},
			{Type: "left", Color: "000000", Style: 1},
			{Type: "right", Color: "000000", Style: 1},
		},
		Font: &excelize.Font{
			Bold: true,
		},
		Fill: excelize.Fill{
			Type:    "pattern",
			Pattern: 1,
			Color:   []string{"#ccff33"},
		},
		NumFmt: 41,
	})
	if err != nil {
		return nil, err
	}

	stylingTotalControl, err := f.NewStyle(&excelize.Style{
		Border: []excelize.Border{
			{Type: "top", Color: "000000", Style: 1},
			{Type: "left", Color: "000000", Style: 1},
			{Type: "right", Color: "000000", Style: 1},
			{Type: "bottom", Color: "000000", Style: 1},
		},
		Font: &excelize.Font{
			Bold: true,
		},
		Fill: excelize.Fill{
			Type:    "pattern",
			Pattern: 1,
			Color:   []string{"#3ada24"},
		},
		NumFmt: 41,
	})
	if err != nil {
		return nil, err
	}

	err = f.MergeCell(sheet, "B6", "B8")
	if err != nil {
		return nil, err
	}
	err = f.MergeCell(sheet, "C6", "E8")
	if err != nil {
		return nil, err
	}
	err = f.MergeCell(sheet, "F6", "F8")
	if err != nil {
		return nil, err
	}
	err = f.MergeCell(sheet, "H6", "K7")
	if err != nil {
		return nil, err
	}

	stylingCurrency, err := f.NewStyle(&excelize.Style{
		NumFmt: 41,
		Border: []excelize.Border{
			{Type: "left", Color: "000000", Style: 1},
			{Type: "right", Color: "000000", Style: 1},
		},
	})
	if err != nil {
		return nil, err
	}

	stylingCurrency2, err := f.NewStyle(&excelize.Style{
		NumFmt: 7,
		Border: []excelize.Border{
			{Type: "left", Color: "000000", Style: 1},
			{Type: "right", Color: "000000", Style: 1},
		},
		CustomNumFmt: &numberFormat,
	})
	if err != nil {
		return nil, err
	}
	err = f.SetCellStyle(sheet, "B6", "L8", stylingHeader)
	if err != nil {
		return nil, err
	}

	err = f.SetCellValue(sheet, "B6", "No Akun")
	if err != nil {
		return nil, err
	}
	err = f.SetCellValue(sheet, "C6", "Keterangan")
	if err != nil {
		return nil, err
	}
	err = f.SetCellValue(sheet, "F6", "WP Reff")
	if err != nil {
		return nil, err
	}
	err = f.SetCellValue(sheet, "G6", "")
	if err != nil {
		return nil, err
	}
	err = f.SetCellValue(sheet, "G7", "Unaudited")
	if err != nil {
		return nil, err
	}
	err = f.SetCellValue(sheet, "G8", "")
	if err != nil {
		return nil, err
	}
	err = f.SetCellValue(sheet, "H6", "Adjustment Journal Entry")
	if err != nil {
		return nil, err
	}
	err = f.SetCellValue(sheet, "I8", "Debet")
	if err != nil {
		return nil, err
	}
	err = f.SetCellValue(sheet, "K8", "Kredit")
	if err != nil {
		return nil, err
	}
	err = f.SetCellFormula(sheet, "L6", "=G6")
	if err != nil {
		return nil, err
	}
	err = f.SetCellFormula(sheet, "L7", "=G7")
	if err != nil {
		return nil, err
	}
	err = f.SetCellFormula(sheet, "L8", "=G8")
	if err != nil {
		return nil, err
	}

	row := 9
	r := 9
	rowCode := make(map[string]int)
	isAutoSum := make(map[string]bool)
	customRow := make(map[string]string)
	for _, a := range *data {
		rowCode[a.Code] = r
		if a.AutoSummary != nil && *a.AutoSummary {
			isAutoSum[a.Code] = true
		}
		if a.IsCoa != nil && *a.IsCoa {
			// rowBefore := row
			tbdetails, err := s.CoaDevRepository.FindWithCode(ctx, &a.Code)
			if err != nil {
				return nil, err
			}
			for _, vTbDetail := range *tbdetails {
				

				tbRowCode[vTbDetail.Code] = r
				
				
			}
			// rowAfter := row - 1
			rowTB := len(*tbdetails)
			if a.AutoSummary != nil && *a.AutoSummary && rowTB > 1 {
				
				rowCode[fmt.Sprintf("%s_SUBTOTAL", a.Code)] = r
				tbRowCode[fmt.Sprintf("%s_SUBTOTAL", a.Code)] = r
				r++
				
			}
		}

	}
	for _, v := range *data {
		rowCode[v.Code] = row
	
		if err = f.SetCellStyle(sheet, fmt.Sprintf("G%d", row), fmt.Sprintf("L%d", row), stylingCurrency); err != nil {
			return nil, err
		}
		if err = f.SetCellStyle(sheet, fmt.Sprintf("F%d", row), fmt.Sprintf("F%d", row), stylingBorderLROnly); err != nil {
			return nil, err
		}
		if v.AutoSummary != nil && *v.AutoSummary {
			isAutoSum[v.Code] = true
		}
		if !(v.IsTotal != nil && *v.IsTotal) && v.IsLabel != nil && *v.IsLabel {
			if err = f.SetCellStyle(sheet, fmt.Sprintf("G%d", row), fmt.Sprintf("L%d", row), stylingCurrency); err != nil {
				return nil, err
			}
			if err = f.SetCellStyle(sheet, fmt.Sprintf("F%d", row), fmt.Sprintf("F%d", row), stylingBorderLROnly); err != nil {
				return nil, err
			}
			f.SetCellValue(sheet, fmt.Sprintf("C%d", row), v.Description)
		}

		if v.IsCoa != nil && *v.IsCoa {
			rowBefore := row
			tbdetails, err := s.CoaDevRepository.FindWithCode(ctx, &v.Code)
			if err != nil {
				return nil, err
			}
			for _, vTbDetail := range *tbdetails {
				if err = f.SetCellStyle(sheet, fmt.Sprintf("G%d", row), fmt.Sprintf("L%d", row), stylingCurrency); err != nil {
					return nil, err
				}
				if err = f.SetCellStyle(sheet, fmt.Sprintf("F%d", row), fmt.Sprintf("F%d", row), stylingBorderLROnly); err != nil {
					return nil, err
				}
				if strings.Contains(strings.ToLower(vTbDetail.Code), "subtotal") {
					continue
				}

				tbRowCode[vTbDetail.Code] = row
				valueKosong := 0.0
				f.SetCellValue(sheet, fmt.Sprintf("B%d", row), vTbDetail.Code)
				f.SetCellValue(sheet, fmt.Sprintf("E%d", row), vTbDetail.Name)
				f.SetCellValue(sheet, fmt.Sprintf("G%d", row), valueKosong)
				f.SetCellValue(sheet, fmt.Sprintf("I%d", row), valueKosong)
				f.SetCellValue(sheet, fmt.Sprintf("K%d", row), valueKosong)
				if err = f.SetCellStyle(sheet, fmt.Sprintf("L%d", row), fmt.Sprintf("L%d", row), stylingCurrency); err != nil {
					return nil, err
				}
				if err = f.SetCellStyle(sheet, fmt.Sprintf("G%d", row), fmt.Sprintf("G%d", row), stylingCurrency); err != nil {
					return nil, err
				}
			
				if err = f.SetCellStyle(sheet, fmt.Sprintf("I%d", row), fmt.Sprintf("I%d", row), stylingCurrency2); err != nil {
					return nil, err
				}
				if err = f.SetCellStyle(sheet, fmt.Sprintf("K%d", row), fmt.Sprintf("K%d", row), stylingCurrency2); err != nil {
					return nil, err
				}
				tmpHeadCoa := fmt.Sprintf("%c", vTbDetail.Code[0])
				if tmpHeadCoa == "9" {
					tmpHeadCoa = vTbDetail.Code[:1]
				}
				switch tmpHeadCoa {
				case "1", "5", "6", "7", "91", "92":
					f.SetCellFormula(sheet, fmt.Sprintf("L%d", row), fmt.Sprintf("=G%d+I%d-K%d", row, row, row))
				default:
					f.SetCellFormula(sheet, fmt.Sprintf("L%d", row), fmt.Sprintf("=G%d-I%d+K%d", row, row, row))
				}
				row++
			}
			rowAfter := row - 1
			rowTB := len(*tbdetails)
			if v.AutoSummary != nil && *v.AutoSummary && rowTB > 1 {
				f.SetCellValue(sheet, fmt.Sprintf("D%d", row), "Subtotal")
				f.SetCellFormula(sheet, fmt.Sprintf("G%d", row), fmt.Sprintf("=SUM(G%d:G%d)", rowBefore, rowAfter))
				f.SetCellFormula(sheet, fmt.Sprintf("I%d", row), fmt.Sprintf("=SUM(I%d:I%d)", rowBefore, rowAfter))
				f.SetCellFormula(sheet, fmt.Sprintf("K%d", row), fmt.Sprintf("=SUM(K%d:K%d)", rowBefore, rowAfter))
				f.SetCellFormula(sheet, fmt.Sprintf("L%d", row), fmt.Sprintf("=SUM(L%d:L%d)", rowBefore, rowAfter))
				f.SetCellStyle(sheet, fmt.Sprintf("D%d", row), fmt.Sprintf("F%d", row), stylingSubTotal)
				if err = f.SetCellStyle(sheet, fmt.Sprintf("G%d", row), fmt.Sprintf("L%d", row), stylingSubTotalCurrency); err != nil {
					return nil, err
				}
				rowCode[fmt.Sprintf("%s_SUBTOTAL", v.Code)] = row
				tbRowCode[fmt.Sprintf("%s_SUBTOTAL", v.Code)] = row
				row++
				if err = f.SetCellStyle(sheet, fmt.Sprintf("F%d", row), fmt.Sprintf("L%d", row), stylingBorderLROnly); err != nil {
					return nil, err
				}
			}
		}

		if v.IsTotal != nil && *v.IsTotal {
			tbRowCode[v.Code] = row
			f.SetCellValue(sheet, fmt.Sprintf("C%d", row), v.Description)
			if v.Code == "TOTAL_LIABILITAS_DAN_EKUITAS" {
				rowAset := row
				if _, ok := rowCode["TOTAL_ASET"]; ok {
					rowAset = rowCode["TOTAL_ASET"]
				}
				f.SetCellFormula(sheet, "G5", fmt.Sprintf("=G%d-G%d", rowAset, row))
				f.SetCellFormula(sheet, "L5", fmt.Sprintf("=L%d-L%d", rowAset, row))
			}

			//show control aje
			if v.Code == "CONTROL" {
				f.SetCellFormula(sheet, "I5", fmt.Sprintf("=I%d", row))
				f.SetCellFormula(sheet, "K5", fmt.Sprintf("=K%d", row))
			}
			if v.Code == "CONTROL_TO_ADJUSTMENT_SHEET" {
				f.SetCellValue(sheet, fmt.Sprintf("C%d", row), v.Description)
				dbt := 0.0

				cdt := 0.0

				f.SetCellValue(sheet, fmt.Sprintf("I%d", row), dbt)
				f.SetCellValue(sheet, fmt.Sprintf("K%d", row), cdt)
				f.SetCellFormula(sheet, fmt.Sprintf("J%d", row), fmt.Sprintf("=I%d-K%d", row, row))
			}
			if v.Code == "TOTAL_JOURNAL_IN_WP" {
				f.SetCellFormula(sheet, fmt.Sprintf("J%d", row), fmt.Sprintf("=I%d-K%d", row, row))
			}

			if v.Code != "TOTAL_JOURNAL_IN_WP" && v.Code != "CONTROL_TO_ADJUSTMENT_SHEET" && v.Code != "CONTROL" && v.Code != "CONTROL_TO_WBS_1" {
				if v.IsControl != nil && *v.IsControl {
					f.SetCellStyle(sheet, fmt.Sprintf("B%d", row), fmt.Sprintf("L%d", row), stylingTotalControl)
				} else {
					f.SetCellStyle(sheet, fmt.Sprintf("B%d", row), fmt.Sprintf("L%d", row), stylingTotal)
				}
			}
			if v.FxSummary == "" {
				row++
				continue
			}
			for chr := 'G'; chr <= 'L'; chr++ {
				if fmt.Sprintf("%c", chr) == "H" || fmt.Sprintf("%c", chr) == "J" || ((fmt.Sprintf("%c", chr) == "G" || fmt.Sprintf("%c", chr) == "L") && (v.Code == "TOTAL_JOURNAL_IN_WP" || v.Code == "CONTROL_TO_ADJUSTMENT_SHEET" || v.Code == "CONTROL")) {
					continue
				}
				formula := v.FxSummary
				// reg := regexp.MustCompile(`[A-Za-z_]+|[0-9]+\d{3,}`)
				reg := regexp.MustCompile(`[A-Za-z0-9_~:()]+|[0-9]+\d{2,}`)
				match := reg.FindAllString(formula, -1)
				for _, vMatch := range match {
					if len(vMatch) < 3 {
						continue
					}
					if isAutoSum[vMatch] {
						if rowCode[fmt.Sprintf("%s_SUBTOTAL", vMatch)] != 0 {
							formula = strings.ReplaceAll(formula, vMatch, fmt.Sprintf("%c%d", chr, rowCode[fmt.Sprintf("%s_SUBTOTAL", vMatch)]))
						} else {
							formula = strings.ReplaceAll(formula, vMatch, fmt.Sprintf("%c%d", chr, rowCode[vMatch]))
						}
					} else {
						if rowCode[vMatch] != 0 {
							formula = strings.ReplaceAll(formula, vMatch, fmt.Sprintf("%c%d", chr, rowCode[vMatch]))
						}
					}
				}
				f.SetCellFormula(sheet, fmt.Sprintf("%c%d", chr, row), fmt.Sprintf("=%s", formula))
			}
			row++
			continue
		}
		row++
	}

	customRow["310401004"] = "=LABA_BERSIH"
	customRow["310402002"] = "=TOTAL_PENGHASILAN_KOMPREHENSIF_LAIN~BS-SUM(310501002,310502002,310503002)"
	customRow["310501002"] = "=950101001"
	customRow["310502002"] = "=950301001+950301002"
	customRow["310503002"] = "=950401001+950401002"
	for key, nRow := range tbRowCode {
		if strings.Contains(customRow["310401004"], key) {
			customRow["310401004"] = strings.ReplaceAll(customRow["310401004"], key, fmt.Sprintf("@%d", nRow))
		}
		if strings.Contains(customRow["310402002"], key) && key != "RE" {
			customRow["310402002"] = strings.ReplaceAll(customRow["310402002"], key, fmt.Sprintf("@%d", nRow))
		}
		if strings.Contains(customRow["310501002"], key) {
			customRow["310501002"] = strings.ReplaceAll(customRow["310501002"], key, fmt.Sprintf("@%d", nRow))
		}
		if strings.Contains(customRow["310502002"], key) {
			customRow["310502002"] = strings.ReplaceAll(customRow["310502002"], key, fmt.Sprintf("@%d", nRow))
		}
		if strings.Contains(customRow["310503002"], key) {
			customRow["310503002"] = strings.ReplaceAll(customRow["310503002"], key, fmt.Sprintf("@%d", nRow))
		}
	}

	for key, vCustomRow := range customRow {
		if val, ok := tbRowCode[key]; ok {
			f.SetCellFormula(sheet, fmt.Sprintf("G%d", val), strings.ReplaceAll(vCustomRow, "@", "G"))
			f.SetCellFormula(sheet, fmt.Sprintf("I%d", val), strings.ReplaceAll(vCustomRow, "@", "I"))
			f.SetCellFormula(sheet, fmt.Sprintf("K%d", val), strings.ReplaceAll(vCustomRow, "@", "K"))
			// f.SetCellFormula(sheet, fmt.Sprintf("L%d", val), strings.ReplaceAll(vCustomRow, "@", "L"))
		}
	}

	if err = f.SetCellStyle(sheet, fmt.Sprintf("B%d", row), fmt.Sprintf("L%d", row), stylingBorderTopOnly); err != nil {
		return nil, err
	}

	if err = f.SetCellStyle(sheet, "A9", fmt.Sprintf("A%d", row-1), stylingBorderRightOnly); err != nil {
		return nil, err
	}
	if err = f.SetCellStyle(sheet, "M9", fmt.Sprintf("M%d", row-1), stylingBorderLeftOnly); err != nil {
		return nil, err
	}

	if err = f.SetSheetFormatPr(sheet, excelize.DefaultRowHeight(12.85)); err != nil {
		return nil, err
	}
	f.SetDefaultFont("Arial")

	return f, nil
}