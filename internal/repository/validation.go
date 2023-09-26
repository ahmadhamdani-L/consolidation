package repository

import (
	"fmt"
	"mcash-finance-console-core/internal/abstraction"
	"mcash-finance-console-core/internal/model"
	"mcash-finance-console-core/pkg/constant"

	"gorm.io/gorm"
)

type Validation interface {
	Find(ctx *abstraction.Context, m *model.TrialBalanceFilterModel, p *abstraction.Pagination) (*[]model.TrialBalanceEntityModel, *abstraction.PaginationInfo, error)
	FindListCompany(ctx *abstraction.Context, m *model.TrialBalanceFilterModel, parentCompanyID *int, p *abstraction.Pagination) (*[]model.TrialBalanceEntityModel, *abstraction.PaginationInfo, error)
	MakeSure(ctx *abstraction.Context, m *model.ValidationFilterModel) (*model.TrialBalanceEntityModel, error)
	GetListCompany(ctx *abstraction.Context, parentCompanyID *int) (map[int]bool, error)
}

type validation struct {
	abstraction.Repository
}

func NewValidation(db *gorm.DB) *validation {
	return &validation{
		abstraction.Repository{
			Db: db,
		},
	}
}

func (r *validation) Find(ctx *abstraction.Context, m *model.TrialBalanceFilterModel, p *abstraction.Pagination) (*[]model.TrialBalanceEntityModel, *abstraction.PaginationInfo, error) {
	conn := r.CheckTrx(ctx)

	var datas []model.TrialBalanceEntityModel
	var info abstraction.PaginationInfo

	query := conn.Model(&model.TrialBalanceEntityModel{})
	//filter
	tableName := model.TrialBalanceEntityModel{}.TableName()
	query = r.Filter(ctx, query, *m)
	query = r.AllowedCompany(ctx, query, tableName)

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
	}
	if p.SortBy != nil && (tmpSortBy != nil && *tmpSortBy != "company") {
		sortBy := fmt.Sprintf("trial_balance.%s", *p.SortBy)
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
	//.Select("id, company_id, period, versions, status, CASE WHEN (SELECT COUNT(*) FROM tabel2 WHERE tabel1_id = t1.id AND status = 0) = 0 THEN 'Balance' ELSE 'Imbalance' END as description")
	if err := query.Joins("Company").Where("status IN (?,?)", constant.MODUL_STATUS_DRAFT, constant.MODUL_STATUS_CONFIRMED).Find(&datas).WithContext(ctx.Request().Context()).Error; err != nil {
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

func (r *validation) FindListCompany(ctx *abstraction.Context, m *model.TrialBalanceFilterModel, parentCompanyID *int, p *abstraction.Pagination) (*[]model.TrialBalanceEntityModel, *abstraction.PaginationInfo, error) {
	conn := r.CheckTrx(ctx)

	var datas []model.TrialBalanceEntityModel
	var info abstraction.PaginationInfo

	query := conn.Model(&model.TrialBalanceEntityModel{})
	//filter
	tableName := model.TrialBalanceEntityModel{}.TableName()
	query = r.Filter(ctx, query, *m)
	query = r.AllowedCompany(ctx, query, tableName)

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

	queryRecursive := "WITH RECURSIVE company_recursive AS (SELECT m_company.id, parent_company_id FROM m_company WHERE parent_company_id = ? OR m_company.id = ? UNION ALL SELECT t1.id, t1.parent_company_id FROM company_recursive t1r INNER JOIN m_company t1 ON t1r.id = t1.parent_company_id) SELECT * FROM company_recursive"

	query = query.Where("status = ?", constant.MODUL_STATUS_DRAFT).Joins(fmt.Sprintf("INNER JOIN (%s) tbl ON trial_balance.company_id = tbl.id", queryRecursive), *parentCompanyID, *parentCompanyID).Count(&totalData).Limit(limit).Offset(offset)
	query = query.Order("company_id asc").Order("versions asc")

	if err := query.Select("trial_balance.*").Group("trial_balance.id").Preload("Company").Find(&datas).WithContext(ctx.Request().Context()).Error; err != nil {
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

func (r *validation) GetListCompany(ctx *abstraction.Context, parentCompanyID *int) (map[int]bool, error) {
	conn := r.CheckTrx(ctx)
	var companyData model.CompanyEntityModel
	if err := conn.Model(&model.CompanyEntityModel{}).Where("id = ?", parentCompanyID).Find(&companyData).Error; err != nil {
		return nil, err
	}

	query := conn.Model(&model.CompanyEntityModel{})
	var data []model.CompanyEntityModel
	queryRecursive := "WITH RECURSIVE company_recursive AS (SELECT m_company.id, parent_company_id FROM m_company WHERE parent_company_id = ? OR id = ? UNION ALL SELECT t1.id, t1.parent_company_id FROM company_recursive t1r INNER JOIN m_company t1 ON t1r.id = t1.parent_company_id) SELECT * FROM company_recursive"

	if err := query.Joins(fmt.Sprintf("INNER JOIN (%s) tbl ON tbl.id = m_company.id", queryRecursive), companyData.ID, companyData.ID).Find(&data).Error; err != nil {
		return nil, err
	}

	res := make(map[int]bool)
	res[companyData.ID] = true
	for _, v := range data {
		res[v.ID] = true
	}

	return res, nil
}

func (r *validation) MakeSure(ctx *abstraction.Context, m *model.ValidationFilterModel) (*model.TrialBalanceEntityModel, error) {
	conn := r.CheckTrx(ctx)
	dataTB := model.TrialBalanceEntityModel{}
	queryTB := conn.Model(&dataTB)
	tableName := model.TrialBalanceEntityModel{}.TableName()
	queryTB = r.Filter(ctx, queryTB, *m)
	queryTB = r.AllowedCompany(ctx, queryTB, tableName)
	queryTB = queryTB.Where("status = ?", constant.MODUL_STATUS_DRAFT)

	if err := queryTB.Find(&dataTB).Error; err != nil {
		return nil, err
	}

	// dataValidation := model.ValidationDetailEntityModel{}
	// dataValidation.CompanyID = dataTB.CompanyID
	// dataValidation.Period = dataTB.Period
	// dataValidation.Status = dataTB.Status
	// queryValidation := conn.Model(&dataValidation)
	// queryValidation = queryValidation.Where("company_id = ?", dataTB.CompanyID).Where("period = ?", dataTB.Period).Where("versions = ?", dataTB.Versions)
	// if err := queryValidation.FirstOrCreate(&dataValidation).Error; err != nil {
	// 	return nil, err
	// }

	// if dataValidation.Status != constant.VALIDATION_STATUS_NOT_BALANCE {
	// 	return nil, errors.New("Data has been validated")
	// }

	return &dataTB, nil
}
