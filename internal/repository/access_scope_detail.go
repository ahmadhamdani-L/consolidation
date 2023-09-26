package repository

import (
	"fmt"
	"mcash-finance-console-core/internal/abstraction"
	"mcash-finance-console-core/internal/model"

	"gorm.io/gorm"
)

type AccessScopeDetail interface {
	Find(ctx *abstraction.Context, m *model.AccessScopeDetailFilterModel, p *abstraction.Pagination) (*[]model.AccessScopeDetailEntityModel, *abstraction.PaginationInfo, error)
	FindByID(ctx *abstraction.Context, id *int) (*model.AccessScopeDetailEntityModel, error)
	Create(ctx *abstraction.Context, e *model.AccessScopeDetailEntityModel) (*model.AccessScopeDetailEntityModel, error)
	Update(ctx *abstraction.Context, id *int, e *model.AccessScopeDetailEntityModel) (*model.AccessScopeDetailEntityModel, error)
	Delete(ctx *abstraction.Context, id *int, e *model.AccessScopeDetailEntityModel) (*model.AccessScopeDetailEntityModel, error)
	DeleteByParent(ctx *abstraction.Context, id *int, e *model.AccessScopeDetailEntityModel) error
}

type accessscopedetail struct {
	abstraction.Repository
}

func NewAccessScopeDetail(db *gorm.DB) *accessscopedetail {
	return &accessscopedetail{
		abstraction.Repository{
			Db: db,
		},
	}
}

func (r *accessscopedetail) Find(ctx *abstraction.Context, m *model.AccessScopeDetailFilterModel, p *abstraction.Pagination) (*[]model.AccessScopeDetailEntityModel, *abstraction.PaginationInfo, error) {
	conn := r.CheckTrx(ctx)

	var datas []model.AccessScopeDetailEntityModel
	var info abstraction.PaginationInfo

	query := conn.Model(&model.AccessScopeDetailEntityModel{})
	//filter
	query = r.Filter(ctx, query, *m)
	queryCompany := conn.Model(&model.CompanyEntityModel{}).Select("id")

	if m.CompanyString != nil {
		query = query.Where("company_id IN (?)", queryCompany.Where("name ILIKE ?", &m.CompanyString))
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

	if err := query.Preload("AccessScope").Preload("Company").Find(&datas).WithContext(ctx.Request().Context()).Error; err != nil {
		return &datas, &info, err
	}

	for i, v := range datas {
		datas[i].CompanyString = &v.Company.Name
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

func (r *accessscopedetail) FindByID(ctx *abstraction.Context, id *int) (*model.AccessScopeDetailEntityModel, error) {
	conn := r.CheckTrx(ctx)

	var data model.AccessScopeDetailEntityModel
	if err := conn.Where("id = ?", &id).Preload("AccessScope").Preload("Company").First(&data).WithContext(ctx.Request().Context()).Error; err != nil {
		return &data, err
	}
	return &data, nil
}

func (r *accessscopedetail) Create(ctx *abstraction.Context, e *model.AccessScopeDetailEntityModel) (*model.AccessScopeDetailEntityModel, error) {
	conn := r.CheckTrx(ctx)

	if err := conn.Create(e).WithContext(ctx.Request().Context()).Error; err != nil {
		return nil, err
	}

	return e, nil
}

func (r *accessscopedetail) Update(ctx *abstraction.Context, id *int, e *model.AccessScopeDetailEntityModel) (*model.AccessScopeDetailEntityModel, error) {
	conn := r.CheckTrx(ctx)

	if err := conn.Model(e).Where("id = ?", &id).Updates(e).WithContext(ctx.Request().Context()).Error; err != nil {
		return nil, err
	}
	return e, nil

}

func (r *accessscopedetail) Delete(ctx *abstraction.Context, id *int, e *model.AccessScopeDetailEntityModel) (*model.AccessScopeDetailEntityModel, error) {
	conn := r.CheckTrx(ctx)

	if err := conn.Model(e).Where("id = ?", id).Delete(e).WithContext(ctx.Request().Context()).Error; err != nil {
		return nil, err
	}
	return e, nil
}

func (r *accessscopedetail) DeleteByParent(ctx *abstraction.Context, accessScopeID *int, e *model.AccessScopeDetailEntityModel) error {
	conn := r.CheckTrx(ctx)

	if err := conn.Model(e).Where("access_scope_id = ?", accessScopeID).Where("company_id != ?", e.CompanyID).Delete(e).WithContext(ctx.Request().Context()).Error; err != nil {
		return err
	}
	return nil
}
