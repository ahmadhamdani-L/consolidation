package repository

import (
	"fmt"
	"mcash-finance-console-core/internal/abstraction"
	"mcash-finance-console-core/internal/model"
	"strconv"

	"gorm.io/gorm"
)

type Jpm interface {
	Find(ctx *abstraction.Context, m *model.JpmFilterModel, p *abstraction.Pagination) (*[]model.JpmEntityModel, *abstraction.PaginationInfo, error)
	FindByID(ctx *abstraction.Context, id *int) (*model.JpmEntityModel, error)
	Create(ctx *abstraction.Context, e *model.JpmEntityModel) (*model.JpmEntityModel, error)
	Update(ctx *abstraction.Context, id *int, e *model.JpmEntityModel) (*model.JpmEntityModel, error)
	Delete(ctx *abstraction.Context, id *int, e *model.JpmEntityModel) (*model.JpmEntityModel, error)
	Get(ctx *abstraction.Context, m *model.JpmFilterModel) (*[]model.JpmEntityModel, error)
	GetCount(ctx *abstraction.Context, m *model.JpmFilterModel) (*int64, error)
	GetVersion(ctx *abstraction.Context, m *model.JpmFilterModel) (*model.GetVersionModel, error)
	UpdateTbd(ctx *abstraction.Context, id *int, e *model.ConsolidationDetailEntityModel) (*model.ConsolidationDetailEntityModel, error)
	FindByTbd(ctx *abstraction.Context, id *int, c *string) (*model.ConsolidationDetailEntityModel, error)
	FindByFormatter(ctx *abstraction.Context, id *int) (*model.FormatterBridgesEntityModel, error)
	Export(ctx *abstraction.Context, id *int) (*model.JpmEntityModel, error)
	FindByCoa(ctx *abstraction.Context, c *string) (*model.CoaEntityModel, error)
	FindByCoaType(ctx *abstraction.Context, idType *int) (*model.CoaTypeEntityModel, error)
	FindByCoaGroup(ctx *abstraction.Context, id *int) (*model.CoaGroupEntityModel, error)
	FindByConsoleBridge(ctx *abstraction.Context, id *int) (*[]model.ConsolidationBridgeEntityModel, error)
	FindByConsoleBridgeDetail(ctx *abstraction.Context, id *int, code *string) (*model.ConsolidationBridgeDetailEntityModel, error)
	FindByCoaConsoleDetail(ctx *abstraction.Context, id int, code *string) (*[]model.ConsolidationDetailEntityModel, error)
	FindSummarys(ctx *abstraction.Context, codeCoa *string, formatterBridgesID *int, isCoa *bool) (*model.ConsolidationDetailEntityModel, error)
	FindSummary(ctx *abstraction.Context) (*[]model.FormatterDetailEntityModel, error)
	FindC(ctx *abstraction.Context, m *model.ConsolidationDetailFilterModel, p *abstraction.Pagination) (*[]model.ConsolidationDetailEntityModel, *abstraction.PaginationInfo, error)
	FindDetailConsole(ctx *abstraction.Context, m *model.ConsolidationDetailFilterModel) (*model.ConsolidationDetailEntityModel, error)
	Updates(ctx *abstraction.Context, id *int, e *model.ConsolidationDetailEntityModel) (*model.ConsolidationDetailEntityModel, error)
	SummaryByCodes(ctx *abstraction.Context, bridgeID *int, codes []string) (*model.ConsolidationDetailEntityModel, error)
	FindByCode(ctx *abstraction.Context, e *model.ConsolidationDetailFilterModel) (*model.ConsolidationDetailEntityModel, error)

}
 
type jpm struct {
	abstraction.Repository
}

func NewJpm(db *gorm.DB) *jpm {
	return &jpm{
		abstraction.Repository{
			Db: db,
		},
	}
}

func (r *jpm) SummaryByCodes(ctx *abstraction.Context, bridgeID *int, codes []string) (*model.ConsolidationDetailEntityModel, error) {
	conn := r.CheckTrx(ctx)
	var data model.ConsolidationDetailEntityModel
	query := conn.Model(&model.ConsolidationDetailEntityModel{}).Where("consolidation_id = ?", bridgeID).Where("code IN (?)", codes).Select("SUM( amount_before_jpm ) amount_before_jpm, SUM( amount_jpm_dr ) amount_jpm_dr, SUM( amount_jpm_cr ) amount_jpm_cr, SUM( amount_after_jpm ) amount_after_jpm ,SUM( amount_jcte_dr ) amount_jcte_dr,SUM( amount_jcte_cr ) amount_jcte_cr,SUM( amount_after_jcte ) amount_after_jcte,SUM( amount_combine_subsidiary ) amount_combine_subsidiary,SUM( amount_jelim_dr ) amount_jelim_dr,SUM( amount_jelim_cr ) amount_jelim_cr,SUM( amount_console ) amount_console")
	if err := query.Take(&data).Error; err != nil {
		return nil, err
	}
	return &data, nil
}
func (r *jpm) FindByCode(ctx *abstraction.Context, e *model.ConsolidationDetailFilterModel) (*model.ConsolidationDetailEntityModel, error) {
	conn := r.CheckTrx(ctx)
	var data model.ConsolidationDetailEntityModel
	query := conn.Model(&model.ConsolidationDetailEntityModel{})
	query = r.Filter(ctx, query, *e)
	if err := query.First(&data).Error; err != nil {
		return nil, err
	}
	return &data, nil
}
func (r *jpm) Updates(ctx *abstraction.Context, id *int, e *model.ConsolidationDetailEntityModel) (*model.ConsolidationDetailEntityModel, error) {
	conn := r.CheckTrx(ctx)

	if err := conn.Model(e).Where("id = ?", &id).Updates(e).WithContext(ctx.Request().Context()).Error; err != nil {
		return nil, err
	}
	if err := conn.Model(e).Where("id = ?", &id).First(e).WithContext(ctx.Request().Context()).Error; err != nil {
		return nil, err
	}
	return e, nil
}
func (r *jpm) FindSummarys(ctx *abstraction.Context, codeCoa *string, formatterBridgesID *int, isCoa *bool) (*model.ConsolidationDetailEntityModel, error) {
	conn := r.CheckTrx(ctx)
	var data model.ConsolidationDetailEntityModel

	tmpStr := *codeCoa
	if _, err := strconv.Atoi(*codeCoa); err == nil {
		tmpStr += "%"
	}
	err := conn.Raw("SELECT SUM( amount_before_jpm ) amount_before_jpm, SUM( amount_jpm_dr ) amount_jpm_dr, SUM( amount_jpm_cr ) amount_jpm_cr, SUM( amount_after_jpm ) amount_after_jpm ,SUM( amount_jcte_dr ) amount_jcte_dr,SUM( amount_jcte_cr ) amount_jcte_cr,SUM( amount_after_jcte ) amount_after_jcte,SUM( amount_combine_subsidiary ) amount_combine_subsidiary,SUM( amount_jelim_dr ) amount_jelim_dr,SUM( amount_jelim_cr ) amount_jelim_cr,SUM( amount_console ) amount_console FROM consolidation_detail WHERE consolidation_id = ? AND LOWER(description) != 'sub total' AND code LIKE '"+tmpStr+"'", &formatterBridgesID).Find(&data).Error
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	return &data, nil
}
func (r *jpm) FindDetailConsole(ctx *abstraction.Context, m *model.ConsolidationDetailFilterModel) (*model.ConsolidationDetailEntityModel, error) {
	conn := r.CheckTrx(ctx)

	var data model.ConsolidationDetailEntityModel
	if err := conn.Where("consolidation_id = ? AND code = ?", &m.ConsolidationID, &m.Code).Find(&data).WithContext(ctx.Request().Context()).Error; err != nil {
		return &data, err
	}
	return &data, nil
}
func (r *jpm) FindC(ctx *abstraction.Context, m *model.ConsolidationDetailFilterModel, p *abstraction.Pagination) (*[]model.ConsolidationDetailEntityModel, *abstraction.PaginationInfo, error) {
	conn := r.CheckTrx(ctx)

	var datas []model.ConsolidationDetailEntityModel
	var info abstraction.PaginationInfo

	query := conn.Model(&model.ConsolidationDetailEntityModel{})
	//filter
	query = r.Filter(ctx, query, *m)

	//sort
	if p.Sort == nil {
		sort := "asc"
		p.Sort = &sort
	}
	if p.SortBy == nil {
		sortBy := "id"
		p.SortBy = &sortBy
	}

	sort := fmt.Sprintf("%s %s", *p.SortBy, *p.Sort)
	query = query.Order(sort)

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
	if m.ConsolidationID != nil {
		query = query.Where("consolidation_id = ?", m.ConsolidationID).Limit(1)
	}

	var totalData int64
	if *p.PageSize != -1 {
		query = query.Count(&totalData).Limit(limit).Offset(offset)
	} else {
		query = query.Count(&totalData)
	}
	if err := query.Find(&datas).WithContext(ctx.Request().Context()).Error; err != nil {
		return &datas, &info, err
	}

	info.Count = int(totalData)
	info.MoreRecords = false
	if len(datas) > *p.PageSize {
		info.MoreRecords = true
		// info.Count -= 1
		// datas = datas[:len(datas)-1]
	}

	return &datas, &info, nil
}
func (r *jpm) FindSummary(ctx *abstraction.Context) (*[]model.FormatterDetailEntityModel, error) {
	conn := r.CheckTrx(ctx)

	data := []model.FormatterDetailEntityModel{}

	coa4dst := conn.Model(&model.FormatterDetailEntityModel{}).Where("code = 'PENDAPATAN'").Where("formatter_id = ?", 3).Select("sort_id").Limit(1)

	query1 := conn.Model(&model.FormatterDetailEntityModel{}).
		Where("formatter_id = ? AND (auto_summary = true OR (is_total = true AND fx_summary IS NOT NULL)) AND is_recalculate = true", 3).
		Order("sort_id ASC").Where("sort_id >= (?)", coa4dst)
	query2 := conn.Model(&model.FormatterDetailEntityModel{}).
		Where("formatter_id = ? AND (auto_summary = true OR (is_total = true AND fx_summary IS NOT NULL)) AND is_recalculate = true", 3).
		Order("sort_id ASC").Where("sort_id < (?)", coa4dst)

	err := query1.Find(&data).Error
	if err != nil {
		return nil, err
	}
	tmp := []model.FormatterDetailEntityModel{}
	err = query2.Find(&tmp).Error
	if err != nil {
		return nil, err
	}

	data = append(data, tmp...)

	return &data, nil
}
func (r *jpm) FindByCoaConsoleDetail(ctx *abstraction.Context, id int, code *string) (*[]model.ConsolidationDetailEntityModel, error) {
	conn := r.CheckTrx(ctx)

	var data []model.ConsolidationDetailEntityModel
	if err := conn.Where("consolidation_id = ? AND code ILIKE ?", id, *code+"%").Find(&data).WithContext(ctx.Request().Context()).Error; err != nil {
		return &data, err
	}
	return &data, nil
}
func (r *jpm) FindByConsoleBridge(ctx *abstraction.Context, id *int) (*[]model.ConsolidationBridgeEntityModel, error) {
	conn := r.CheckTrx(ctx)

	var data []model.ConsolidationBridgeEntityModel
	if err := conn.Where("consolidation_id = ?", &id).Find(&data).WithContext(ctx.Request().Context()).Error; err != nil {
		return &data, err
	}
	return &data, nil
}
func (r *jpm) FindByConsoleBridgeDetail(ctx *abstraction.Context, id *int, code *string) (*model.ConsolidationBridgeDetailEntityModel, error) {
	conn := r.CheckTrx(ctx)

	var data model.ConsolidationBridgeDetailEntityModel
	if err := conn.Where("consolidation_bridge_id = ? AND code = ?", &id, &code).Find(&data).WithContext(ctx.Request().Context()).Error; err != nil {
		return &data, err
	}
	return &data, nil
}
func (r *jpm) FindByCoaGroup(ctx *abstraction.Context, id *int) (*model.CoaGroupEntityModel, error) {
	conn := r.CheckTrx(ctx)

	var data model.CoaGroupEntityModel
	if err := conn.Where("id = ?", &id).First(&data).WithContext(ctx.Request().Context()).Error; err != nil {
		return &data, err
	}
	return &data, nil
}
func (r *jpm) FindByCoa(ctx *abstraction.Context, c *string) (*model.CoaEntityModel, error) {
	conn := r.CheckTrx(ctx)

	var data model.CoaEntityModel
	if err := conn.Where("code = ? ", &c).First(&data).WithContext(ctx.Request().Context()).Error; err != nil {
		return &data, err
	}
	return &data, nil
}
func (r *jpm) FindByCoaType(ctx *abstraction.Context, idType *int) (*model.CoaTypeEntityModel, error) {
	conn := r.CheckTrx(ctx)

	var data model.CoaTypeEntityModel
	if err := conn.Where("id = ?", &idType).First(&data).WithContext(ctx.Request().Context()).Error; err != nil {
		return &data, err
	}
	return &data, nil
}
func (r *jpm) Find(ctx *abstraction.Context, m *model.JpmFilterModel, p *abstraction.Pagination) (*[]model.JpmEntityModel, *abstraction.PaginationInfo, error) {
	conn := r.CheckTrx(ctx)

	var datas []model.JpmEntityModel
	var info abstraction.PaginationInfo

	query := conn.Model(&model.JpmEntityModel{})
	//filter
	tableName := model.JpmEntityModel{}.TableName()
	query = r.FilterTable(ctx, query, *m, tableName)
	query = r.AllowedCompany(ctx, query, tableName)

	//filter custom
	tmp1 := m.CompanyCustomFilter
	queryCompany := conn.Model(&model.CompanyEntityModel{}).Select("id")
	query = r.FilterMultiCompany(ctx, query, queryCompany, tmp1, "")
	query = r.FilterMultiVersion(ctx, query, m.JpmFilter)
	queryUser := conn.Model(&model.UserEntityModel{}).Select("id")
	query = r.FilterUser(ctx, query, queryUser, m.Filter, "")

	if m.Search != nil {
		query = query.Where("created_by IN (?)", conn.Model(&model.UserEntityModel{}).Select("id").Where("name ILIKE ?", "%"+*m.Search+"%"))
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
		sortBy := fmt.Sprintf("jpm.%s", *p.SortBy)
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

	if err := query.Preload("Consolidation").Preload("JpmDetail").Joins("Company").Joins("UserCreated", func(db *gorm.DB) *gorm.DB {
		db = db.Select("id, name")
		return db
	}).Preload("UserModified").Find(&datas).WithContext(ctx.Request().Context()).Error; err != nil {
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
func (r *jpm) FindByID(ctx *abstraction.Context, id *int) (*model.JpmEntityModel, error) {
	conn := r.CheckTrx(ctx)

	var data model.JpmEntityModel
	if err := conn.Where("id = ?", &id).Preload("Company").Preload("Consolidation").Preload("JpmDetail").Preload("UserCreated").Preload("UserModified").First(&data).WithContext(ctx.Request().Context()).Error; err != nil {
		return &data, err
	}
	data.UserCreatedString = data.UserCreated.Name
	data.UserModifiedString = &data.UserModified.Name
	return &data, nil
}
func (r *jpm) Create(ctx *abstraction.Context, e *model.JpmEntityModel) (*model.JpmEntityModel, error) {
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
func (r *jpm) Update(ctx *abstraction.Context, id *int, e *model.JpmEntityModel) (*model.JpmEntityModel, error) {
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
func (r *jpm) UpdateTbd(ctx *abstraction.Context, id *int, e *model.ConsolidationDetailEntityModel) (*model.ConsolidationDetailEntityModel, error) {
	conn := r.CheckTrx(ctx)

	if err := conn.Model(e).Where("id = ?", &id).Updates(e).WithContext(ctx.Request().Context()).Error; err != nil {
		return nil, err
	}
	if err := conn.Model(e).First(e).WithContext(ctx.Request().Context()).Error; err != nil {
		return nil, err
	}
	return e, nil
}
func (r *jpm) FindByTbd(ctx *abstraction.Context, id *int, c *string) (*model.ConsolidationDetailEntityModel, error) {
	conn := r.CheckTrx(ctx)

	var data model.ConsolidationDetailEntityModel
	if err := conn.Where("consolidation_id = ? and code = ?", &id, &c).First(&data).WithContext(ctx.Request().Context()).Error; err != nil {
		return &data, err
	}
	return &data, nil
}
func (r *jpm) FindByFormatter(ctx *abstraction.Context, id *int) (*model.FormatterBridgesEntityModel, error) {
	conn := r.CheckTrx(ctx)

	var data model.FormatterBridgesEntityModel
	if err := conn.Where("trx_ref_id = ? ", &id).First(&data).WithContext(ctx.Request().Context()).Error; err != nil {
		return &data, err
	}
	return &data, nil
}
func (r *jpm) Destroy(ctx *abstraction.Context, id *int, e *model.JpmEntityModel) (*model.JpmEntityModel, error) {
	conn := r.CheckTrx(ctx)

	if err := conn.Model(e).Where("id =?", id).Delete(e).WithContext(ctx.Request().Context()).Error; err != nil {
		return nil, err
	}
	return e, nil
}
func (r *jpm) Delete(ctx *abstraction.Context, id *int, e *model.JpmEntityModel) (*model.JpmEntityModel, error) {
	conn := r.CheckTrx(ctx)

	if err := conn.Model(e).Where("id =?", id).Update("status", 4).WithContext(ctx.Request().Context()).Error; err != nil {
		return nil, err
	}
	e.UserCreatedString = e.UserCreated.Name
	e.UserModifiedString = &e.UserModified.Name
	return e, nil
}
func (r *jpm) Get(ctx *abstraction.Context, m *model.JpmFilterModel) (*[]model.JpmEntityModel, error) {
	var datas []model.JpmEntityModel

	conn := r.CheckTrx(ctx)
	query := conn.Model(&model.JpmEntityModel{})
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
func (r *jpm) GetCount(ctx *abstraction.Context, m *model.JpmFilterModel) (*int64, error) {
	var jmlData int64
	conn := r.CheckTrx(ctx)
	query := conn.Model(&model.JpmEntityModel{})
	//filter
	query = r.Filter(ctx, query, *m)

	if err := query.Count(&jmlData).WithContext(ctx.Request().Context()).Error; err != nil {
		return &jmlData, err
	}

	return &jmlData, nil
}
func (r *jpm) GetVersion(ctx *abstraction.Context, m *model.JpmFilterModel) (*model.GetVersionModel, error) {
	var data []model.JpmEntityModel
	conn := r.CheckTrx(ctx)
	query := conn.Model(&model.JpmEntityModel{})
	query = r.Filter(ctx, query, *m)

	//filter custom
	tmp1 := m.CompanyCustomFilter
	queryCompany := conn.Model(&model.CompanyEntityModel{}).Select("id")
	query = r.FilterMultiCompany(ctx, query, queryCompany, tmp1, "")

	query = query.Select("consolidation_versions").Group("consolidation_versions").Order("consolidation_versions ASC")

	if err := query.Find(&data).Error; err != nil {
		return &model.GetVersionModel{}, err
	}

	var result model.GetVersionModel
	tmp := []map[string]string{}
	for _, v := range data {
		tmp1 := map[string]string{
			"value": fmt.Sprintf("%d", v.ConsolidationID),
			"label": fmt.Sprintf("Version %d", v.ConsolidationID),
		}
		tmp = append(tmp, tmp1)
	}
	result.Version = tmp
	return &result, nil
}
func (r *jpm) Export(ctx *abstraction.Context, id *int) (*model.JpmEntityModel, error) {
	conn := r.CheckTrx(ctx)
	var data model.JpmEntityModel
	query := conn.Model(&model.JpmEntityModel{}).Preload("Company").Preload("JpmDetail").Where("id = ?", &id).Find(&data)
	if err := query.Error; err != nil {
		return nil, err
	}
	return &data, nil
}
