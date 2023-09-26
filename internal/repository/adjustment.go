package repository

import (
	"fmt"
	"mcash-finance-console-core/internal/abstraction"
	"mcash-finance-console-core/internal/model"
	"mcash-finance-console-core/pkg/util/helper"
	"strings"

	"gorm.io/gorm"
)

type Adjustment interface {
	Find(ctx *abstraction.Context, m *model.AdjustmentFilterModel, p *abstraction.Pagination) (*[]model.AdjustmentEntityModel, *abstraction.PaginationInfo, error)
	FindByID(ctx *abstraction.Context, id *int) (*model.AdjustmentEntityModel, error)
	Create(ctx *abstraction.Context, e *model.AdjustmentEntityModel) (*model.AdjustmentEntityModel, error)
	Update(ctx *abstraction.Context, id *int, e *model.AdjustmentEntityModel) (*model.AdjustmentEntityModel, error)
	Delete(ctx *abstraction.Context, id *int, e *model.AdjustmentEntityModel) (*model.AdjustmentEntityModel, error)
	Get(ctx *abstraction.Context, m *model.AdjustmentFilterModel) (*[]model.AdjustmentEntityModel, error)
	GetCount(ctx *abstraction.Context, m *model.AdjustmentFilterModel) (*int64, error)
	GetVersion(ctx *abstraction.Context, m *model.AdjustmentFilterModel) (*model.GetVersionModel, error)
	UpdateTbd(ctx *abstraction.Context, id *int, e *model.TrialBalanceDetailEntityModel) (*model.TrialBalanceDetailEntityModel, error)
	FindByTbd(ctx *abstraction.Context, id *int, c *string) (*model.TrialBalanceDetailEntityModel, error)
	FindByFormatter(ctx *abstraction.Context, id *int, c *string) (*model.FormatterBridgesEntityModel, error)
	Export(ctx *abstraction.Context, id *int) (*model.AdjustmentEntityModel, error)
	FindByTbds(ctx *abstraction.Context, m *model.TrialBalanceDetailFilterModel) (*[]model.TrialBalanceDetailEntityModel, error)
	FindSummary(ctx *abstraction.Context, tbID *int) (*model.AdjustmentDetailEntityModel, error)
	ExportAll(ctx *abstraction.Context, trialBalanceID *int) (*[]model.AdjustmentDetailEntityModel, error)
}

type adjustment struct {
	abstraction.Repository
}

func NewAdjustment(db *gorm.DB) *adjustment {
	return &adjustment{
		abstraction.Repository{
			Db: db,
		},
	}
}

func (r *adjustment) Find(ctx *abstraction.Context, m *model.AdjustmentFilterModel, p *abstraction.Pagination) (*[]model.AdjustmentEntityModel, *abstraction.PaginationInfo, error) {
	conn := r.CheckTrx(ctx)

	var datas []model.AdjustmentEntityModel
	var info abstraction.PaginationInfo

	query := conn.Model(&model.AdjustmentEntityModel{})
	//filter
	tableName := model.AdjustmentEntityModel{}.TableName()
	query = r.FilterTable(ctx, query, *m, tableName)
	query = r.AllowedCompany(ctx, query, tableName)

	//filter custom
	tmp1 := m.CompanyCustomFilter
	queryCompany := conn.Model(&model.CompanyEntityModel{}).Select("id")
	query = r.FilterMultiCompany(ctx, query, queryCompany, tmp1, tableName)
	query = r.FilterMultiVersion(ctx, query, m.AdjustmentFilter)
	queryUser := conn.Model(&model.UserEntityModel{}).Select("id")
	query = r.FilterUser(ctx, query, queryUser, m.Filter, tableName)
	query = query.Where("status != 0")

	if m.Search != nil {
		query = query.Where("adjustment.created_by IN (?)", conn.Model(&model.UserEntityModel{}).Select("id").Where("name ILIKE ?", "%"+*m.Search+"%"))
	}

	//sort
	if p.Sort == nil {
		sort := "desc"
		p.Sort = &sort
	}
	if p.SortBy == nil {
		sortBy := "created_at"
		p.SortBy = &sortBy
	}
	tmpSortBy := p.SortBy
	if p.SortBy != nil && *p.SortBy == "company" {
		sortBy := "\"Company\".name"
		p.SortBy = &sortBy
	} else if p.SortBy != nil && *p.SortBy == "user" {
		sortBy := "\"UserCreated\".name"
		p.SortBy = &sortBy
	}
	if p.SortBy != nil && (tmpSortBy != nil && *tmpSortBy != "company" && *tmpSortBy != "user") {
		sortBy := fmt.Sprintf("adjustment.%s", *p.SortBy)
		p.SortBy = &sortBy
	}
	sort := fmt.Sprintf("%s %s", *p.SortBy, *p.Sort)
	query = query.Order(sort)
	p.SortBy = tmpSortBy
	//pagination
	if p.Page == nil {
		page := 1
		p.Page = &page
	}
	if p.PageSize == nil {
		pageSize := 10
		p.PageSize = &pageSize
	}
	info = abstraction.PaginationInfo{
		Pagination: p,
	}
	limit := *p.PageSize
	offset := limit * (*p.Page - 1)
	var totalData int64
	query = query.Count(&totalData).Limit(limit).Offset(offset)

	if err := query.Joins("Company").Joins("UserCreated", func(db *gorm.DB) *gorm.DB {
		db = db.Select("id, name")
		return db
	}).Preload("UserModified").Preload("TrialBalance").Find(&datas).WithContext(ctx.Request().Context()).Error; err != nil {
		return &datas, &info, err
	}

	for i, v := range datas {
		datas[i].UserCreatedString = v.UserCreated.Name
		datas[i].UserModifiedString = &v.UserModified.Name
	}

	info.Count = int(totalData)
	info.MoreRecords = false
	if int(totalData) > *p.PageSize {
		info.MoreRecords = true
		// info.Count -= 1
		// datas = datas[:len(datas)-1]
	}

	return &datas, &info, nil
}

func (r *adjustment) FindByID(ctx *abstraction.Context, id *int) (*model.AdjustmentEntityModel, error) {
	conn := r.CheckTrx(ctx)

	var data model.AdjustmentEntityModel
	if err := conn.Where("id = ?", &id).Preload("TrialBalance").Preload("Company").Preload("AdjustmentDetail").Preload("UserCreated").Preload("UserModified").First(&data).WithContext(ctx.Request().Context()).Error; err != nil {
		return &data, err
	}
	data.UserCreatedString = data.UserCreated.Name
	data.UserModifiedString = &data.UserModified.Name
	return &data, nil
}

func (r *adjustment) Create(ctx *abstraction.Context, e *model.AdjustmentEntityModel) (*model.AdjustmentEntityModel, error) {
	conn := r.CheckTrx(ctx)

	if err := conn.Create(e).WithContext(ctx.Request().Context()).Error; err != nil {
		return nil, err
	}
	if err := conn.Model(e).Preload("UserCreated").Preload("UserModified").First(e).WithContext(ctx.Request().Context()).Error; err != nil {
		return nil, err
	}
	e.UserCreatedString = e.UserCreated.Name
	e.UserModifiedString = &e.UserModified.Name
	return e, nil
}

func (r *adjustment) Update(ctx *abstraction.Context, id *int, e *model.AdjustmentEntityModel) (*model.AdjustmentEntityModel, error) {
	conn := r.CheckTrx(ctx)

	if err := conn.Model(e).Where("id = ?", &id).Updates(e).Preload("Company").WithContext(ctx.Request().Context()).Error; err != nil {
		return nil, err
	}
	if err := conn.Model(e).Where("id = ?", &id).Preload("Company").Preload("UserCreated").Preload("UserModified").First(e).WithContext(ctx.Request().Context()).Error; err != nil {
		return nil, err
	}
	e.UserCreatedString = e.UserCreated.Name
	e.UserModifiedString = &e.UserModified.Name
	return e, nil

}
func (r *adjustment) UpdateTbd(ctx *abstraction.Context, id *int, e *model.TrialBalanceDetailEntityModel) (*model.TrialBalanceDetailEntityModel, error) {
	conn := r.CheckTrx(ctx)

	if err := conn.Model(e).Where("id = ?", &id).Updates(e).WithContext(ctx.Request().Context()).Error; err != nil {
		return nil, err
	}
	if err := conn.Model(e).First(e).WithContext(ctx.Request().Context()).Error; err != nil {
		return nil, err
	}
	return e, nil

}

func (r *adjustment) FindByTbd(ctx *abstraction.Context, id *int, c *string) (*model.TrialBalanceDetailEntityModel, error) {
	conn := r.CheckTrx(ctx)

	var data model.TrialBalanceDetailEntityModel
	if err := conn.Where("formatter_bridges_id = ? and code = ?", &id, &c).First(&data).WithContext(ctx.Request().Context()).Error; err != nil {
		return &data, err
	}
	return &data, nil
}

func (r *adjustment) FindByTbds(ctx *abstraction.Context, m *model.TrialBalanceDetailFilterModel) (*[]model.TrialBalanceDetailEntityModel, error) {
	conn := r.CheckTrx(ctx)

	var data []model.TrialBalanceDetailEntityModel
	if err := conn.Where("formatter_bridges_id = ? and code = ?", &m.FormatterBridgesID, m.Code).Find(&data).WithContext(ctx.Request().Context()).Error; err != nil {
		return &data, err
	}
	return &data, nil
}

func (r *adjustment) FindByFormatter(ctx *abstraction.Context, id *int, c *string) (*model.FormatterBridgesEntityModel, error) {
	conn := r.CheckTrx(ctx)

	var data model.FormatterBridgesEntityModel
	if err := conn.Where("trx_ref_id = ? AND source = ? ", &id, &c).Find(&data).WithContext(ctx.Request().Context()).Error; err != nil {
		return &data, err
	}
	return &data, nil
}

func (r *adjustment) Destroy(ctx *abstraction.Context, id *int, e *model.AdjustmentEntityModel) (*model.AdjustmentEntityModel, error) {
	conn := r.CheckTrx(ctx)

	if err := conn.Model(e).Where("id =?", id).Delete(e).WithContext(ctx.Request().Context()).Error; err != nil {
		return nil, err
	}
	return e, nil
}

func (r *adjustment) Delete(ctx *abstraction.Context, id *int, e *model.AdjustmentEntityModel) (*model.AdjustmentEntityModel, error) {
	conn := r.CheckTrx(ctx)

	if err := conn.Model(e).Where("id =?", id).Update("status", 4).WithContext(ctx.Request().Context()).Error; err != nil {
		return nil, err
	}
	e.UserCreatedString = e.UserCreated.Name
	e.UserModifiedString = &e.UserModified.Name
	return e, nil
}

func (r *adjustment) Get(ctx *abstraction.Context, m *model.AdjustmentFilterModel) (*[]model.AdjustmentEntityModel, error) {
	var datas []model.AdjustmentEntityModel

	conn := r.CheckTrx(ctx)
	query := conn.Model(&model.AdjustmentEntityModel{})
	//filter
	query = r.Filter(ctx, query, *m)

	if err := query.Find(&datas).WithContext(ctx.Request().Context()).Error; err != nil {
		return &datas, err
	}

	// var formatterBridges []model.FormatterBridgesEntityModel

	// query = conn.Model(&model.FormatterBridgesEntityModel{}).Where("trx_ref_id = ?", datas.ID).Where("source = ?", "TRIAL-BALANCE")
	// if err := query.Find(&formatterBridges).Error; err != nil {
	// 	return &datas, err
	// }
	// datas.FormatterBridges = formatterBridges

	return &datas, nil
}

func (r *adjustment) GetCount(ctx *abstraction.Context, m *model.AdjustmentFilterModel) (*int64, error) {
	var jmlData int64
	conn := r.CheckTrx(ctx)
	query := conn.Model(&model.AdjustmentEntityModel{})
	//filter
	query = r.Filter(ctx, query, *m)

	if err := query.Count(&jmlData).WithContext(ctx.Request().Context()).Error; err != nil {
		return &jmlData, err
	}

	return &jmlData, nil
}

func (r *adjustment) GetVersion(ctx *abstraction.Context, m *model.AdjustmentFilterModel) (*model.GetVersionModel, error) {
	var data []model.AdjustmentEntityModel
	conn := r.CheckTrx(ctx)
	query := conn.Model(&model.AdjustmentEntityModel{})
	query = r.Filter(ctx, query, *m)

	//filter custom
	tmp1 := m.CompanyCustomFilter
	queryCompany := conn.Model(&model.CompanyEntityModel{}).Select("id")
	query = r.FilterMultiCompany(ctx, query, queryCompany, tmp1, "")

	query = query.Select("tb_id").Group("tb_id").Order("tb_id ASC")

	if err := query.Find(&data).Error; err != nil {
		return &model.GetVersionModel{}, err
	}

	var result model.GetVersionModel
	tmp := []map[string]string{}
	for _, v := range data {
		tmp1 := map[string]string{
			"value": fmt.Sprintf("%d", v.TbID),
			"label": fmt.Sprintf("Version %d", v.TbID),
		}
		tmp = append(tmp, tmp1)
	}
	result.Version = tmp
	return &result, nil
}

func (r *adjustment) Export(ctx *abstraction.Context, id *int) (*model.AdjustmentEntityModel, error) {
	conn := r.CheckTrx(ctx)
	var data model.AdjustmentEntityModel
	query := conn.Model(&model.AdjustmentEntityModel{}).Preload("Company").Preload("AdjustmentDetail").Where("id = ?", &id).Find(&data)
	if err := query.Error; err != nil {
		return nil, err
	}
	return &data, nil
}

func (r *adjustment) FindSummary(ctx *abstraction.Context, tbID *int) (*model.AdjustmentDetailEntityModel, error) {
	conn := r.CheckTrx(ctx)
	var data []string
	if err := conn.Model(&model.AdjustmentEntityModel{}).Where("tb_id = ?", &tbID).Pluck("id", &data).Error; err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	}

	if len(data) == 0 {
		return &model.AdjustmentDetailEntityModel{}, nil
	}

	listID := strings.Join(data, ",")
	var sumData model.AdjustmentDetailEntityModel
	if err := conn.Model(&model.AdjustmentDetailEntityModel{}).Where(fmt.Sprintf("adjustment_id IN (%s)", listID)).Select("SUM(balance_sheet_cr) balance_sheet_cr, SUM(balance_sheet_dr) balance_sheet_dr, SUM(income_statement_cr) income_statement_cr, SUM(income_statement_dr) income_statement_dr").Find(&sumData).Error; err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	}

	sumData.BalanceSheetCr = helper.AssignAmount(sumData.BalanceSheetCr)
	sumData.BalanceSheetDr = helper.AssignAmount(sumData.BalanceSheetDr)
	sumData.IncomeStatementCr = helper.AssignAmount(sumData.IncomeStatementCr)
	sumData.IncomeStatementDr = helper.AssignAmount(sumData.IncomeStatementDr)

	return &sumData, nil
}

func (r *adjustment) ExportAll(ctx *abstraction.Context, trialBalanceID *int) (*[]model.AdjustmentDetailEntityModel, error) {
	conn := r.CheckTrx(ctx)
	var data []model.AdjustmentDetailEntityModel
	queryAJE := conn.Model(&model.AdjustmentEntityModel{}).Where("tb_id = ?", &trialBalanceID).Select("id")
	query := conn.Model(&model.AdjustmentDetailEntityModel{}).Where("adjustment_id IN (?)", queryAJE).Order("id ASC").Find(&data)
	if err := query.Error; err != nil {
		return nil, err
	}
	return &data, nil
}
