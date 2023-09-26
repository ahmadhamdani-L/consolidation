package repository

import (
	"fmt"
	"mcash-finance-console-core/internal/abstraction"
	"mcash-finance-console-core/internal/model"
	"strings"

	"gorm.io/gorm"
)

type AccessScope interface {
	Find(ctx *abstraction.Context, m *model.AccessScopeFilterModel, p *abstraction.Pagination) (*[]model.AccessScopeEntityModel, *abstraction.PaginationInfo, error)
	FindByID(ctx *abstraction.Context, id *int) (*model.AccessScopeEntityModel, error)
	Create(ctx *abstraction.Context, e *model.AccessScopeEntityModel) (*model.AccessScopeEntityModel, error)
	Update(ctx *abstraction.Context, id *int, e *model.AccessScopeEntityModel) (*model.AccessScopeEntityModel, error)
	Delete(ctx *abstraction.Context, id *int, e *model.AccessScopeEntityModel) (*model.AccessScopeEntityModel, error)
	FindWithCompanyByID(ctx *abstraction.Context, id *int) (*[]model.AccessScopeDetailListEntityModel, error)
	FindByCriteria(ctx *abstraction.Context, m *model.AccessScopeFilterModel) (*model.AccessScopeEntityModel, error)
}

type accessscope struct {
	abstraction.Repository
}

func NewAccessScope(db *gorm.DB) *accessscope {
	return &accessscope{
		abstraction.Repository{
			Db: db,
		},
	}
}

func (r *accessscope) Find(ctx *abstraction.Context, m *model.AccessScopeFilterModel, p *abstraction.Pagination) (*[]model.AccessScopeEntityModel, *abstraction.PaginationInfo, error) {
	conn := r.CheckTrx(ctx)

	var datas []model.AccessScopeEntityModel
	var info abstraction.PaginationInfo

	query := conn.Model(&model.AccessScopeEntityModel{})
	//filter
	tmpCompanyID := m.CompanyID
	m.CompanyID = nil
	query = r.Filter(ctx, query, *m)
	m.CompanyID = tmpCompanyID
	tmp1 := m.CompanyCustomFilter
	queryCompany := conn.Model(&model.CompanyEntityModel{}).Select("id")
	query = r.FilterMultiCompany(ctx, query, queryCompany, tmp1, "")
	queryUser := conn.Model(&model.UserEntityModel{}).Select("id")

	if m.UserString != nil {
		query = query.Where("user_id IN (?)", queryUser.Where("username ILIKE ?", "%"+*m.UserString+"%"))
	}

	if m.UserIsActive != nil {
		query = query.Where("user_id IN (?)", queryUser.Where("is_active = ?", m.UserIsActive))
	}

	if m.CompanyID != nil {
		query = query.Where("(id IN (?) OR access_all = true)", conn.Model(&model.AccessScopeDetailEntityModel{}).Select("access_scope_id").Where("company_id = ?", m.CompanyID))
	}

	//sort
	if p.Sort == nil {
		sort := "desc"
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
	var totalData int64
	query = query.Count(&totalData).Limit(limit).Offset(offset)

	if err := query.Preload("AccessScopeDetail.Company").Preload("User").Find(&datas).WithContext(ctx.Request().Context()).Error; err != nil {
		return &datas, &info, err
	}
	for i, v := range datas {
		datas[i].UserIsActive = v.User.IsActive
		var arrCompanyList []string
		datas[i].UserString = v.User.Username
		for iDetail, vDetail := range v.AccessScopeDetail {
			arrCompanyList = append(arrCompanyList, vDetail.Company.Name)
			if iDetail == 2 {
				break
			}
		}
		datas[i].JmlCompany = len(v.AccessScopeDetail)
		if v.AccessAll != nil && *v.AccessAll {
			datas[i].CompanyList = "All Branch"
		} else if len(v.AccessScopeDetail) > 0 {
			strCompanyList := strings.Join(arrCompanyList, ", ")
			datas[i].CompanyList = strings.TrimSpace(strCompanyList) + "..."
		}
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

func (r *accessscope) FindWithCompanyByID(ctx *abstraction.Context, id *int) (*[]model.AccessScopeDetailListEntityModel, error) {
	conn := r.CheckTrx(ctx)

	var datasDetail []model.AccessScopeDetailListEntityModel
	queryDetail := conn.Model(&model.CompanyEntityModel{}).Joins("c LEFT JOIN access_scope_detail ON c.id = access_scope_detail.company_id AND access_scope_detail.access_scope_id = ?", id).Select("c.*, CASE WHEN ( access_scope_detail.id IS NOT NULL ) THEN TRUE ELSE FALSE END checked").Where("is_active = true").Order("checked DESC, name ASC")
	if err := queryDetail.Find(&datasDetail).WithContext(ctx.Request().Context()).Error; err != nil {
		return &datasDetail, err
	}

	return &datasDetail, nil
}

func (r *accessscope) FindByID(ctx *abstraction.Context, id *int) (*model.AccessScopeEntityModel, error) {
	conn := r.CheckTrx(ctx)

	var data model.AccessScopeEntityModel
	if err := conn.Where("id = ?", &id).Preload("AccessScopeDetail").Preload("User").First(&data).WithContext(ctx.Request().Context()).Error; err != nil {
		return &data, err
	}
	return &data, nil
}

func (r *accessscope) Create(ctx *abstraction.Context, e *model.AccessScopeEntityModel) (*model.AccessScopeEntityModel, error) {
	conn := r.CheckTrx(ctx)

	if err := conn.Create(e).WithContext(ctx.Request().Context()).Error; err != nil {
		return nil, err
	}

	return e, nil
}

func (r *accessscope) Update(ctx *abstraction.Context, id *int, e *model.AccessScopeEntityModel) (*model.AccessScopeEntityModel, error) {
	conn := r.CheckTrx(ctx)

	if err := conn.Model(e).Where("id = ?", &id).Updates(e).WithContext(ctx.Request().Context()).Error; err != nil {
		return nil, err
	}
	return e, nil

}

func (r *accessscope) Delete(ctx *abstraction.Context, id *int, e *model.AccessScopeEntityModel) (*model.AccessScopeEntityModel, error) {
	conn := r.CheckTrx(ctx)

	if err := conn.Model(e).Where("id =?", id).Delete(e).WithContext(ctx.Request().Context()).Error; err != nil {
		return nil, err
	}
	return e, nil
}

func (r *accessscope) FindByCriteria(ctx *abstraction.Context, m *model.AccessScopeFilterModel) (*model.AccessScopeEntityModel, error) {
	conn := r.CheckTrx(ctx)
	query := conn.Model(&model.AccessScopeEntityModel{})
	query = r.Filter(ctx, query, *m)
	var data model.AccessScopeEntityModel
	if err := query.Preload("AccessScopeDetail.Company").Preload("User").First(&data).WithContext(ctx.Request().Context()).Error; err != nil {
		return nil, err
	}
	data.UserString = data.User.Name

	return &data, nil
}
