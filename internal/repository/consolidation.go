package repository

import (
	"fmt"
	"mcash-finance-console-core/internal/abstraction"
	"mcash-finance-console-core/internal/model"
	"mcash-finance-console-core/pkg/util/helper"
	utilHelper "mcash-finance-console-core/pkg/util/helper"

	"gorm.io/gorm"
)

type Consolidation interface {
	FindListCompanyParent(ctx *abstraction.Context, m *model.TrialBalanceFilterModel, id int, p *abstraction.Pagination) (*[]model.TrialBalanceEntityModel, *abstraction.PaginationInfo, error)
	FindListCompanyParents(ctx *abstraction.Context, m *model.TrialBalanceFilterModel, p *abstraction.Pagination) (*[]model.TrialBalanceEntityModel, *abstraction.PaginationInfo, error)
	FindListCompanyChild(ctx *abstraction.Context, m *model.TrialBalanceFilterModel, id []int, p *abstraction.Pagination) (*[]model.TrialBalanceEntityModel, *abstraction.PaginationInfo, error)
	FindListCompanyChilds(ctx *abstraction.Context, m *model.TrialBalanceFilterModel, p *abstraction.Pagination) (*[]model.TrialBalanceEntityModel, *abstraction.PaginationInfo, error)
	Find(ctx *abstraction.Context, m *model.ConsolidationFilterModel, p *abstraction.Pagination) (*[]model.ConsolidationEntityModel, *abstraction.PaginationInfo, error)
	FindListCompanyConsole(ctx *abstraction.Context, m *model.ConsolidationBridgeFilterModel, p *abstraction.Pagination) (*[]model.ConsolidationBridgeEntityModel, *abstraction.PaginationInfo, error)
	FindByID(ctx *abstraction.Context, id *int) (*model.ConsolidationEntityModel, error)
	FindByTBID(ctx *abstraction.Context, m *model.TrialBalanceFilterModel) (*model.TrialBalanceEntityModel, error)
	FindByConsolidationID(ctx *abstraction.Context, m *model.ConsolidationFilterModel) (*model.ConsolidationEntityModel, error)
	FindByConsolidationIDs(ctx *abstraction.Context, m *model.ConsolidationFilterModel) (*model.ConsolidationEntityModel, error)
	GetVersion(ctx *abstraction.Context, m *model.ConsolidationFilterModel) (*model.GetVersionModel, error)
	FindListCompanyCreateNewCombine(ctx *abstraction.Context, m *model.CompanyFilterModel) (*[]model.CompanyEntityModel, error)
	FindByParentCompanyID(ctx *abstraction.Context, id *int) (*[]model.CompanyEntityModel, error)
	FindByParentCompanyIDS(ctx *abstraction.Context, id *int) (*[]model.CompanyEntityModel, error)
	FindByConsolidationBridge(ctx *abstraction.Context, m *model.ConsolidationBridgeFilterModel) (*model.ConsolidationBridgeEntityModel, error)
	FindListCompanyChildConsole(ctx *abstraction.Context, m *model.ConsolidationFilterModel, p *abstraction.Pagination) (*[]model.ConsolidationEntityModel, *abstraction.PaginationInfo, error)
	FindListCompanyChildConsoles(ctx *abstraction.Context, m *model.ConsolidationFilterModel, id []int, p *abstraction.Pagination) (*[]model.ConsolidationEntityModel, *abstraction.PaginationInfo, error)
	Destroy(ctx *abstraction.Context, id *int, e *model.ConsolidationEntityModel) (*model.ConsolidationEntityModel, error)
	DestroyDetail(ctx *abstraction.Context, id *int, e *model.ConsolidationDetailEntityModel) (*model.ConsolidationDetailEntityModel, error)
	Update(ctx *abstraction.Context, id *int, e *model.ConsolidationEntityModel) (*model.ConsolidationEntityModel, error)
	FindSummaryJournal(ctx *abstraction.Context, consolidationID *int) (*model.ConsolidationDetailEntityModel, error)
	FindControl2Wbs1(ctx *abstraction.Context, consolidationID *int) (*model.ConsolidationDetailEntityModel, error)
	FindAnakUsaha(ctx *abstraction.Context, consolidationID int) ([]model.ConsolidationBridgeEntityModel, error)
	CompanyChild(ctx *abstraction.Context, m *model.CompanyEntityModel) (*[]model.CompanyEntityModel, error)
}

type consolidation struct {
	abstraction.Repository
}

func NewConsolidation(db *gorm.DB) *consolidation {
	return &consolidation{
		abstraction.Repository{
			Db: db,
		},
	}
}

func (r *consolidation) CompanyChild(ctx *abstraction.Context, m *model.CompanyEntityModel) (*[]model.CompanyEntityModel, error) {
	conn := r.CheckTrx(ctx)

	var datas []model.CompanyEntityModel
	
	query := conn.Model(&model.CompanyEntityModel{})



	query = query.Where("parent_company_id = ?", m.ID)
	

	if err := query.Preload("UserCreated").Preload("UserModified").Find(&datas).WithContext(ctx.Request().Context()).Error; err != nil {
		return &datas , nil
	}


	// info.Count = int(totalData)
	// info.MoreRecords = false
	// if int(totalData) > *p.PageSize {
	// 	info.MoreRecords = true
	// 	info.Count -= 1
	// 	datas = datas[:len(datas)-1]
	// }

	return &datas, nil
}

func (r *consolidation) Update(ctx *abstraction.Context, id *int, e *model.ConsolidationEntityModel) (*model.ConsolidationEntityModel, error) {
	conn := r.CheckTrx(ctx)
	err := conn.Model(e).Where("id = ?", id).Update("status", 0).WithContext(ctx.Request().Context()).Error
	if err != nil {
		return nil, err
	}

	e.UserCreatedString = e.UserCreated.Name
	e.UserModifiedString = &e.UserModified.Name
	return e, nil

}
func (r *consolidation) FindListCompanyChildConsoles(ctx *abstraction.Context, m *model.ConsolidationFilterModel, id []int, p *abstraction.Pagination) (*[]model.ConsolidationEntityModel, *abstraction.PaginationInfo, error) {
	conn := r.CheckTrx(ctx)

	var datas []model.ConsolidationEntityModel
	var info abstraction.PaginationInfo

	query := conn.Model(&model.ConsolidationEntityModel{})

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

	if err := query.Where("DATE(period) = ? AND company_id = ? AND status IN (2) AND consolidation.id NOT IN (?)", m.Period, m.CompanyID, id).Preload("Company").Find(&datas).WithContext(ctx.Request().Context()).Error; err != nil {
		return &datas, &info, err
	}

	info.Count = int(totalData)
	info.MoreRecords = false
	if int(totalData) > *p.PageSize {
		info.MoreRecords = true
		info.Count -= 1
		datas = datas[:len(datas)-1]
	}

	return &datas, &info, nil
}
func (r *consolidation) Destroy(ctx *abstraction.Context, id *int, e *model.ConsolidationEntityModel) (*model.ConsolidationEntityModel, error) {
	conn := r.CheckTrx(ctx)

	if err := conn.Model(e).Where("id =?", id).Delete(e).WithContext(ctx.Request().Context()).Error; err != nil {
		return nil, err
	}
	return e, nil
}
func (r *consolidation) DestroyDetail(ctx *abstraction.Context, id *int, e *model.ConsolidationDetailEntityModel) (*model.ConsolidationDetailEntityModel, error) {
	conn := r.CheckTrx(ctx)

	if err := conn.Model(e).Where("consolidation_id =?", id).Delete(e).WithContext(ctx.Request().Context()).Error; err != nil {
		return nil, err
	}
	return e, nil
}
func (r *consolidation) FindListCompanyCreateNewCombine(ctx *abstraction.Context, m *model.CompanyFilterModel) (*[]model.CompanyEntityModel, error) {
	conn := r.CheckTrx(ctx)
	var data []model.CompanyEntityModel
	query := "SELECT distinct mc.id, mc.name FROM m_company mc INNER JOIN m_company mc2 ON mc.id = mc2.parent_company_id WHERE mc.id = mc2.parent_company_id"
	if err := conn.Raw(query).Find(&data).WithContext(ctx.Request().Context()).Error; err != nil {
		return &data, err
	}
	// Joins("INNER JOIN m_company mc ON m_company.id = mc.parent_company_id")
	// query = query.Where("m_company.id = mc.parent_company_id")

	return &data, nil
}
func (r *consolidation) FindListCompanyParent(ctx *abstraction.Context, m *model.TrialBalanceFilterModel, id int, p *abstraction.Pagination) (*[]model.TrialBalanceEntityModel, *abstraction.PaginationInfo, error) {
	conn := r.CheckTrx(ctx)

	var datas []model.TrialBalanceEntityModel
	var info abstraction.PaginationInfo

	query := conn.Model(&model.TrialBalanceEntityModel{})

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

	if err := query.Where("DATE(period) = ? AND company_id = ? AND status IN (2) AND trial_balance.id != ? AND consolidation_id IS NULL", m.Period, m.CompanyID, id).Joins(fmt.Sprintf("INNER JOIN m_company ON m_company.id = trial_balance.company_id")).Group("trial_balance.id").Preload("Company").Find(&datas).WithContext(ctx.Request().Context()).Error; err != nil {
		return &datas, &info, err
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
func (r *consolidation) FindListCompanyParents(ctx *abstraction.Context, m *model.TrialBalanceFilterModel, p *abstraction.Pagination) (*[]model.TrialBalanceEntityModel, *abstraction.PaginationInfo, error) {
	conn := r.CheckTrx(ctx)

	var datas []model.TrialBalanceEntityModel
	var info abstraction.PaginationInfo

	query := conn.Model(&model.TrialBalanceEntityModel{})

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

	if err := query.Where("DATE(period) = ? AND company_id = ? AND status IN (2) AND consolidation_id IS NULL ", m.Period, m.CompanyID).Joins(fmt.Sprintf("INNER JOIN m_company ON m_company.id = trial_balance.company_id")).Group("trial_balance.id").Preload("Company").Find(&datas).WithContext(ctx.Request().Context()).Error; err != nil {
		return &datas, &info, err
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
func (r *consolidation) FindListCompanyChild(ctx *abstraction.Context, m *model.TrialBalanceFilterModel, id []int, p *abstraction.Pagination) (*[]model.TrialBalanceEntityModel, *abstraction.PaginationInfo, error) {
	conn := r.CheckTrx(ctx)

	var datas []model.TrialBalanceEntityModel
	var info abstraction.PaginationInfo

	query := conn.Model(&model.TrialBalanceEntityModel{})
	//filter
	// query = r.Filter(ctx, query, *m)
	// query = r.AllowedCompany(ctx, query, tableName)

	//sort
	// if p.Sort == nil {
	// 	sort := "desc"
	// 	p.Sort = &sort
	// }
	// if p.SortBy == nil {
	// 	sortBy := "created_at"
	// 	p.SortBy = &sortBy
	// }

	// sort := fmt.Sprintf("%s %s", *p.SortBy, *p.Sort)
	// query = query.Order(sort)

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

	if err := query.Where("DATE(period) = ? AND mc.id = ? AND status IN (2) AND trial_balance.id NOT IN (?) AND consolidation_id IS NULL", m.Period, m.CompanyID, id).Joins(fmt.Sprintf("INNER JOIN m_company mc ON mc.id = trial_balance.company_id")).Group("trial_balance.id").Preload("Company").Find(&datas).WithContext(ctx.Request().Context()).Error; err != nil {
		return &datas, &info, err
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
func (r *consolidation) FindListCompanyChilds(ctx *abstraction.Context, m *model.TrialBalanceFilterModel, p *abstraction.Pagination) (*[]model.TrialBalanceEntityModel, *abstraction.PaginationInfo, error) {
	conn := r.CheckTrx(ctx)

	var datas []model.TrialBalanceEntityModel
	var info abstraction.PaginationInfo

	query := conn.Model(&model.TrialBalanceEntityModel{})
	//filter
	// query = r.Filter(ctx, query, *m)
	// query = r.AllowedCompany(ctx, query, tableName)

	//sort
	// if p.Sort == nil {
	// 	sort := "desc"
	// 	p.Sort = &sort
	// }
	// if p.SortBy == nil {
	// 	sortBy := "created_at"
	// 	p.SortBy = &sortBy
	// }

	// sort := fmt.Sprintf("%s %s", *p.SortBy, *p.Sort)
	// query = query.Order(sort)

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

	if err := query.Where("DATE(period) = ? AND mc.id = ? AND status IN (2) AND consolidation_id IS NULL", m.Period, m.CompanyID).Joins(fmt.Sprintf("INNER JOIN m_company mc ON mc.id = trial_balance.company_id")).Group("trial_balance.id").Preload("Company").Find(&datas).WithContext(ctx.Request().Context()).Error; err != nil {
		return &datas, &info, err
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
func (r *consolidation) FindListCompanyConsole(ctx *abstraction.Context, m *model.ConsolidationBridgeFilterModel, p *abstraction.Pagination) (*[]model.ConsolidationBridgeEntityModel, *abstraction.PaginationInfo, error) {
	conn := r.CheckTrx(ctx)

	var datas []model.ConsolidationBridgeEntityModel
	var info abstraction.PaginationInfo

	query := conn.Model(&model.ConsolidationBridgeEntityModel{})

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

	if err := query.Where("DATE(period) = ? AND consolidation_id = ? ", m.Period, m.ConsolidationID).Group("consolidation_bridge.id").Preload("Company").Find(&datas).WithContext(ctx.Request().Context()).Error; err != nil {
		return &datas, &info, err
	}

	info.Count = int(totalData)
	info.MoreRecords = false
	if int(totalData) > *p.PageSize {
		info.MoreRecords = true
		info.Count -= 1
		datas = datas[:len(datas)-1]
	}

	return &datas, &info, nil
}
func (r *consolidation) FindListCompanyConsoles(ctx *abstraction.Context, m *model.ConsolidationBridgeFilterModel, id []int, p *abstraction.Pagination) (*[]model.ConsolidationBridgeEntityModel, *abstraction.PaginationInfo, error) {
	conn := r.CheckTrx(ctx)

	var datas []model.ConsolidationBridgeEntityModel
	var info abstraction.PaginationInfo

	query := conn.Model(&model.ConsolidationBridgeEntityModel{})

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

	if err := query.Where("DATE(period) = ? AND consolidation_id = ? AND consolidation_detail.id NOT IN (?)", m.Period, m.ConsolidationID, id).Group("consolidation_bridge.id").Preload("Company").Find(&datas).WithContext(ctx.Request().Context()).Error; err != nil {
		return &datas, &info, err
	}

	info.Count = int(totalData)
	info.MoreRecords = false
	if int(totalData) > *p.PageSize {
		info.MoreRecords = true
		info.Count -= 1
		datas = datas[:len(datas)-1]
	}

	return &datas, &info, nil
}
func (r *consolidation) Find(ctx *abstraction.Context, m *model.ConsolidationFilterModel, p *abstraction.Pagination) (*[]model.ConsolidationEntityModel, *abstraction.PaginationInfo, error) {
	conn := r.CheckTrx(ctx)

	var datas []model.ConsolidationEntityModel
	var info abstraction.PaginationInfo

	query := conn.Model(&model.ConsolidationEntityModel{})
	//filter
	tableName := model.ConsolidationEntityModel{}.TableName()
	query = r.Filter(ctx, query, *m)
	query = r.AllowedCompany(ctx, query, tableName)

	//filter custom
	tmp1 := m.CompanyCustomFilter
	queryCompany := conn.Model(&model.CompanyEntityModel{}).Select("id")
	query = r.FilterMultiCompany(ctx, query, queryCompany, tmp1, "")
	query = r.FilterMultiVersion(ctx, query, m.ConsolidationFilter)
	queryUser := conn.Model(&model.UserEntityModel{}).Select("id")
	query = r.FilterUser(ctx, query, queryUser, m.Filter, "")
	query = query.Where("status != 5")

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
	var totalData int64
	query = query.Count(&totalData).Limit(limit).Offset(offset)

	if err := query.Preload("Company").Preload("UserCreated").Preload("UserModified").Find(&datas).WithContext(ctx.Request().Context()).Error; err != nil {
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
func (r *consolidation) FindByID(ctx *abstraction.Context, id *int) (*model.ConsolidationEntityModel, error) {
	conn := r.CheckTrx(ctx)

	var data model.ConsolidationEntityModel
	if err := conn.Where("id = ?", &id).Preload("Company").Preload("UserCreated").Preload("UserModified").First(&data).WithContext(ctx.Request().Context()).Error; err != nil {
		return &data, err
	}
	data.UserCreatedString = data.UserCreated.Name
	data.UserModifiedString = &data.UserModified.Name
	return &data, nil
}
func (r *consolidation) FindByTBID(ctx *abstraction.Context, m *model.TrialBalanceFilterModel) (*model.TrialBalanceEntityModel, error) {
	conn := r.CheckTrx(ctx)

	var data model.TrialBalanceEntityModel
	if err := conn.Where("DATE(period) = ? AND company_id = ? AND versions = ?", m.Period, m.CompanyID, m.Versions).Preload("Company").Preload("UserCreated").Preload("UserModified").First(&data).WithContext(ctx.Request().Context()).Error; err != nil {
		return &data, err
	}
	data.UserCreatedString = data.UserCreated.Name
	data.UserModifiedString = &data.UserModified.Name
	return &data, nil
}
func (r *consolidation) FindByConsolidationID(ctx *abstraction.Context, m *model.ConsolidationFilterModel) (*model.ConsolidationEntityModel, error) {
	conn := r.CheckTrx(ctx)

	var data model.ConsolidationEntityModel
	if err := conn.Where("DATE(period) = ? AND company_id = ? AND consolidation_versions = ?", m.Period, m.CompanyID, m.ConsolidationVersions).Preload("Company").Preload("UserCreated").Preload("UserModified").First(&data).WithContext(ctx.Request().Context()).Error; err != nil {
		return &data, err
	}
	data.UserCreatedString = data.UserCreated.Name
	data.UserModifiedString = &data.UserModified.Name
	return &data, nil
}
func (r *consolidation) FindByConsolidationIDs(ctx *abstraction.Context, m *model.ConsolidationFilterModel) (*model.ConsolidationEntityModel, error) {
	conn := r.CheckTrx(ctx)

	var data model.ConsolidationEntityModel
	if err := conn.Where("DATE(period) = ? AND company_id = ? AND versions = ?", m.Period, m.CompanyID, m.Versions).Preload("Company").Preload("UserCreated").Preload("UserModified").First(&data).WithContext(ctx.Request().Context()).Error; err != nil {
		return &data, err
	}
	data.UserCreatedString = data.UserCreated.Name
	data.UserModifiedString = &data.UserModified.Name
	return &data, nil
}
func (r *consolidation) GetVersion(ctx *abstraction.Context, m *model.ConsolidationFilterModel) (*model.GetVersionModel, error) {
	var data []model.ConsolidationEntityModel
	conn := r.CheckTrx(ctx)
	query := conn.Model(&model.ConsolidationEntityModel{})
	query = r.Filter(ctx, query, *m)
	query = query.Where("status != 0")

	//filter custom
	tmp1 := m.CompanyCustomFilter
	queryCompany := conn.Model(&model.CompanyEntityModel{}).Select("id")
	query = r.FilterMultiCompany(ctx, query, queryCompany, tmp1, "")
	query = r.FilterMultiStatus(ctx, query, *m)

	query = query.Select("consolidation_versions").Group("consolidation_versions").Order("consolidation_versions ASC")

	if err := query.Find(&data).Error; err != nil {
		return &model.GetVersionModel{}, err
	}

	var result model.GetVersionModel
	tmp := []map[string]string{}
	for _, v := range data {
		tmp1 := map[string]string{
			"value": fmt.Sprintf("%d", v.ConsolidationVersions),
			"label": fmt.Sprintf("Consolidation-Version %d", v.ConsolidationVersions),
		}
		tmp = append(tmp, tmp1)
	}
	result.Version = tmp
	return &result, nil
}
func (r *consolidation) FindByParentCompanyID(ctx *abstraction.Context, id *int) (*[]model.CompanyEntityModel, error) {
	conn := r.CheckTrx(ctx)

	var data []model.CompanyEntityModel
	if err := conn.Where("parent_company_id = ?", &id).Preload("UserCreated").Preload("UserModified").Find(&data).WithContext(ctx.Request().Context()).Error; err != nil {
		return &data, err
	}
	return &data, nil
}
func (r *consolidation) FindByParentCompanyIDS(ctx *abstraction.Context, id *int) (*[]model.CompanyEntityModel, error) {
	conn := r.CheckTrx(ctx)

	var data []model.CompanyEntityModel
	if err := conn.Where("parent_company_id = ?", &id).Preload("UserCreated").Preload("UserModified").Find(&data).WithContext(ctx.Request().Context()).Error; err != nil {
		return &data, err
	}
	return &data, nil
}
func (r *consolidation) FindListCompanyChildConsole(ctx *abstraction.Context, m *model.ConsolidationFilterModel, p *abstraction.Pagination) (*[]model.ConsolidationEntityModel, *abstraction.PaginationInfo, error) {
	conn := r.CheckTrx(ctx)

	var datas []model.ConsolidationEntityModel
	var info abstraction.PaginationInfo

	query := conn.Model(&model.ConsolidationEntityModel{})
	//filter
	// query = r.Filter(ctx, query, *m)
	// query = r.AllowedCompany(ctx, query, tableName)

	//sort
	// if p.Sort == nil {
	// 	sort := "desc"
	// 	p.Sort = &sort
	// }
	// if p.SortBy == nil {
	// 	sortBy := "created_at"
	// 	p.SortBy = &sortBy
	// }

	// sort := fmt.Sprintf("%s %s", *p.SortBy, *p.Sort)
	// query = query.Order(sort)

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

	if err := query.Where("DATE(period) = ? AND mc.id = ? AND status IN (2)", m.Period, m.CompanyID).Joins(fmt.Sprintf("INNER JOIN m_company mc ON mc.id = consolidation.company_id")).Group("consolidation.id").Preload("Company").Find(&datas).WithContext(ctx.Request().Context()).Error; err != nil {
		return &datas, &info, err
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
func (r *consolidation) FindByConsolidationBridge(ctx *abstraction.Context, m *model.ConsolidationBridgeFilterModel) (*model.ConsolidationBridgeEntityModel, error) {
	conn := r.CheckTrx(ctx)

	var data model.ConsolidationBridgeEntityModel
	if err := conn.Where("DATE(period) = ? AND company_id = ? AND versions = ?  AND consolidation_versions = ? ", m.Period, m.CompanyID, m.Versions, m.ConsolidationVersions).First(&data).WithContext(ctx.Request().Context()).Error; err != nil {
		return &data, err
	}
	return &data, nil
}

func (r *consolidation) FindSummaryJournal(ctx *abstraction.Context, consolidationID *int) (*model.ConsolidationDetailEntityModel, error) {
	conn := r.CheckTrx(ctx)
	var result model.ConsolidationDetailEntityModel
	zero := 0.0
	result.AmountJpmCr = &zero
	result.AmountJpmDr = &zero
	result.AmountJcteCr = &zero
	result.AmountJcteDr = &zero
	result.AmountJelimCr = &zero
	result.AmountJelimDr = &zero

	// Query dan perhitungan untuk data JpmEntityModel
	{
		var data []int
		if err := conn.Model(&model.JpmEntityModel{}).
			Where("status != 4").
			Where("consolidation_id = ?", consolidationID).
			Pluck("id", &data).
			Error; err != nil && err != gorm.ErrRecordNotFound {
			return nil, err
		}

		if len(data) != 0 {
			var sumDataJpm model.JpmDetailEntityModel
			if err := conn.Model(&model.JpmDetailEntityModel{}).
				Where("jpm_id IN (?)", data).
				Select("(SUM(balance_sheet_cr)+SUM(income_statement_cr)) balance_sheet_cr, (SUM(balance_sheet_dr) + SUM(income_statement_dr)) balance_sheet_dr").
				Find(&sumDataJpm).
				Error; err != nil && err != gorm.ErrRecordNotFound {
				return nil, err
			}

			result.AmountJpmCr = utilHelper.AssignAmount(sumDataJpm.BalanceSheetCr)
			result.AmountJpmDr = utilHelper.AssignAmount(sumDataJpm.BalanceSheetDr)
		}
	}

	// Query dan perhitungan untuk data JcteEntityModel
	{
		var data []int
		if err := conn.Model(&model.JcteEntityModel{}).
			Where("status != 4").
			Where("consolidation_id = ?", consolidationID).
			Pluck("id", &data).
			Error; err != nil && err != gorm.ErrRecordNotFound {
			return nil, err
		}

		if len(data) != 0 {
			var sumDataJcte model.JcteDetailEntityModel
			if err := conn.Model(&model.JcteDetailEntityModel{}).
				Where("jcte_id IN (?)", data).
				Select("(SUM(balance_sheet_cr)+SUM(income_statement_cr)) balance_sheet_cr, (SUM(balance_sheet_dr) + SUM(income_statement_dr)) balance_sheet_dr").
				Find(&sumDataJcte).
				Error; err != nil && err != gorm.ErrRecordNotFound {
				return nil, err
			}
			result.AmountJcteCr = utilHelper.AssignAmount(sumDataJcte.BalanceSheetCr)
			result.AmountJcteDr = utilHelper.AssignAmount(sumDataJcte.BalanceSheetDr)
		}
	}

	// Query dan perhitungan untuk data JelimEntityModel
	{
		var data []int
		if err := conn.Model(&model.JelimEntityModel{}).
			Where("status != 0").
			Where("consolidation_id = ?", consolidationID).
			Pluck("id", &data).
			Error; err != nil && err != gorm.ErrRecordNotFound {
			return nil, err
		}

		if len(data) != 0 {
			var sumDataJelim model.JelimDetailEntityModel
			if err := conn.Model(&model.JelimDetailEntityModel{}).
				Where("jelim_id IN (?)", data).
				Select("(SUM(balance_sheet_cr)+SUM(income_statement_cr)) balance_sheet_cr, (SUM(balance_sheet_dr) + SUM(income_statement_dr)) balance_sheet_dr").
				Find(&sumDataJelim).
				Error; err != nil && err != gorm.ErrRecordNotFound {
				return nil, err
			}

			result.AmountJelimCr = utilHelper.AssignAmount(sumDataJelim.BalanceSheetCr)
			result.AmountJelimDr = utilHelper.AssignAmount(sumDataJelim.BalanceSheetDr)
		}
	}

	return &result, nil
}
func (r *consolidation) FindControl2Wbs1(ctx *abstraction.Context, consolidationID *int) (*model.ConsolidationDetailEntityModel, error) {
	conn := r.CheckTrx(ctx)

	var dataFmtBridges []model.ConsolidationBridgeEntityModel
	err := conn.Model(&model.ConsolidationBridgeEntityModel{}).Where("consolidation_id = ?", consolidationID).Preload("Company").Preload("ConsolidationBridgeDetail", func(db *gorm.DB) *gorm.DB {
		query := db.Where("code = 'TOTAL_ASET'")
		strquery := "SELECT SUM(amount) FROM consolidation_bridge_detail tmp WHERE tmp.consolidation_bridge_id = consolidation_bridge_detail.consolidation_bridge_id AND code = 'TOTAL_LIABILITAS_DAN_EKUITAS'"
		query = query.Select(fmt.Sprintf("SUM(amount - (%s)) amount, consolidation_bridge_id ", strquery)).Group("consolidation_bridge_id")
		return query
	}).Find(&dataFmtBridges).Error
	if err != nil {
		return nil, err
	}

	var result model.ConsolidationDetailEntityModel
	queryFrom := "( SELECT consolidation_id, COALESCE ( amount_before_jpm, 0 ) amount_before_jpm, COALESCE ( amount_jpm_dr, 0 ) amount_jpm_dr, COALESCE ( amount_jpm_cr, 0 ) amount_jpm_cr, COALESCE ( amount_after_jpm, 0 ) amount_after_jpm, COALESCE ( amount_jcte_dr, 0 ) amount_jcte_dr, COALESCE ( amount_jcte_cr, 0 ) amount_jcte_cr, COALESCE ( amount_after_jcte, 0 ) amount_after_jcte, COALESCE ( amount_combine_subsidiary, 0 ) amount_combine_subsidiary, COALESCE ( amount_jelim_dr, 0 ) amount_jelim_dr, COALESCE ( amount_jelim_cr, 0 ) amount_jelim_cr, COALESCE ( amount_console, 0 ) amount_console FROM consolidation_detail WHERE code = ? AND consolidation_id = ? )"
	query := fmt.Sprintf("SELECT ( aset.amount_before_jpm - liaeku.amount_before_jpm ) amount_before_jpm, ( aset.amount_jpm_cr + liaeku.amount_jpm_cr ) amount_jpm_cr, ( aset.amount_jpm_dr + liaeku.amount_jpm_dr ) amount_jpm_dr, ( aset.amount_after_jpm - liaeku.amount_after_jpm ) amount_after_jpm, ( aset.amount_jcte_cr + liaeku.amount_jcte_cr ) amount_jcte_cr, ( aset.amount_jcte_dr + liaeku.amount_jcte_dr ) amount_jcte_dr, ( aset.amount_after_jcte - liaeku.amount_after_jcte ) amount_after_jcte, ( aset.amount_combine_subsidiary - liaeku.amount_combine_subsidiary ) amount_combine_subsidiary, ( aset.amount_jelim_cr + liaeku.amount_jelim_cr ) amount_jelim_cr, ( aset.amount_jelim_dr + liaeku.amount_jelim_dr ) amount_jelim_dr, ( aset.amount_console - liaeku.amount_console ) amount_console FROM %s aset JOIN %s liaeku ON aset.consolidation_id = liaeku.consolidation_id", queryFrom, queryFrom)
	if err := conn.Raw(query, "TOTAL_ASET", consolidationID, "TOTAL_LIABILITAS_DAN_EKUITAS", consolidationID).Find(&result).Error; err != nil {
		return nil, err
	}

	result.AmountBeforeJpm = helper.AssignAmount(result.AmountBeforeJpm)
	result.AmountJpmDr = helper.AssignAmount(result.AmountJpmDr)
	result.AmountJpmCr = helper.AssignAmount(result.AmountJpmCr)
	result.AmountAfterJpm = helper.AssignAmount(result.AmountAfterJpm)
	result.AmountJcteDr = helper.AssignAmount(result.AmountJcteDr)
	result.AmountJcteCr = helper.AssignAmount(result.AmountJcteCr)
	result.AmountAfterJcte = helper.AssignAmount(result.AmountAfterJcte)
	result.AmountCombineSubsidiary = helper.AssignAmount(result.AmountCombineSubsidiary)
	result.AmountJelimDr = helper.AssignAmount(result.AmountJelimDr)
	result.AmountJelimCr = helper.AssignAmount(result.AmountJelimCr)
	result.AmountConsole = helper.AssignAmount(result.AmountConsole)
	result.ConsolidationBridge = dataFmtBridges

	return &result, nil
}

func (r *consolidation) FindAnakUsaha(ctx *abstraction.Context, consolidationID int) ([]model.ConsolidationBridgeEntityModel, error) {
	conn := r.CheckTrx(ctx)

	var datas []model.ConsolidationBridgeEntityModel

	query := conn.Model(&model.ConsolidationBridgeEntityModel{})
	if err := query.Where("consolidation_id = ? ", consolidationID).Preload("Company").
	Find(&datas).WithContext(ctx.Request().Context()).Error; err != nil {
		return datas, nil
	}
	
	// querycount := query
	// var totalData int64
	// if err := querycount.Count(&totalData).Error; err != nil {
	// 	return &datas, &info, err
	// }
	// if *p.PageSize != -1 {
	// 	query = query.Limit(limit).Offset(offset)
	// }
	// if err := query.Find(&datas).WithContext(ctx.Request().Context()).Error; err != nil {
	// 	return &datas, &info, err
	// }

	// info.Count = int(totalData)
	// info.MoreRecords = false
	// if len(datas) > *p.PageSize {
	// 	info.MoreRecords = true
	// 	info.Count -= 1
	// 	datas = datas[:len(datas)-1]
	// }

	return datas, nil
}
