package repository

import (
	"fmt"
	"mcash-finance-console-core/internal/abstraction"
	"mcash-finance-console-core/internal/model"
	"strings"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type Company interface {
	Find(ctx *abstraction.Context, m *model.CompanyFilterModel, p *abstraction.Pagination) (*[]model.CompanyEntityModel, *abstraction.PaginationInfo, error)
	FindByID(ctx *abstraction.Context, id *int) (*model.CompanyEntityModel, error)
	Create(ctx *abstraction.Context, e *model.CompanyEntityModel) (*model.CompanyEntityModel, error)
	Update(ctx *abstraction.Context, id *int, e *model.CompanyEntityModel) (*model.CompanyEntityModel, error)
	Delete(ctx *abstraction.Context, id *int, e *model.CompanyEntityModel) (*model.CompanyEntityModel, error)
	FindFilterList(ctx *abstraction.Context) (*[]model.CompanyEntityModel, error)
	FindWithCode(ctx *abstraction.Context, code *string) (*model.CompanyEntityModel, error)
	FindWithName(ctx *abstraction.Context, name *string) (*[]model.CompanyEntityModel, error)
	FindIsActive(ctx *abstraction.Context, t *bool) (*[]model.CompanyEntityModel, error)
}

type company struct {
	abstraction.Repository
}

func NewCompany(db *gorm.DB) *company {
	return &company{
		abstraction.Repository{
			Db: db,
		},
	}
}

func (r *company) FindIsActive(ctx *abstraction.Context, t *bool) (*[]model.CompanyEntityModel, error) {
	conn := r.CheckTrx(ctx)

	var data []model.CompanyEntityModel
	if err := conn.Where("is_active = ?", &t).Preload("UserCreated").Preload("UserModified").Find(&data).WithContext(ctx.Request().Context()).Error; err != nil {
		return &data, err
	}
	return &data, nil
}

func (r *company) FindWithCode(ctx *abstraction.Context, code *string) (*model.CompanyEntityModel, error) {
	conn := r.CheckTrx(ctx)

	var data model.CompanyEntityModel
	if err := conn.Where("code = ?", *code).Find(&data).WithContext(ctx.Request().Context()).Error; err != nil {
		return &data, err
	}
	return &data, nil
}
func (r *company) FindWithName(ctx *abstraction.Context, name *string) (*[]model.CompanyEntityModel, error) {
	conn := r.CheckTrx(ctx)

	var data []model.CompanyEntityModel
	if err := conn.Where("name ILIKE ?", "%"+*name).Find(&data).WithContext(ctx.Request().Context()).Error; err != nil {
		return &data, err
	}
	return &data, nil
}

func (r *company) Find(ctx *abstraction.Context, m *model.CompanyFilterModel, p *abstraction.Pagination) (*[]model.CompanyEntityModel, *abstraction.PaginationInfo, error) {
	conn := r.CheckTrx(ctx)

	var datas []model.CompanyEntityModel
	var info abstraction.PaginationInfo

	query := conn.Model(&model.CompanyEntityModel{})

	//filter
	tableName := model.CompanyEntityModel{}.TableName()
	query = r.Filter(ctx, query, *m)
	queryUser := conn.Model(&model.UserEntityModel{}).Select("id")
	query = r.FilterUser(ctx, query, queryUser, m.Filter, tableName)

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

	if m.ParentCompanyID == nil && m.WithChild != nil && *m.WithChild {
		query = query.Where("parent_company_id IS NULL").Preload(clause.Associations)
	}

	if err := query.Preload("ParentCompany").Preload("UserCreated").Preload("UserModified").Find(&datas).WithContext(ctx.Request().Context()).Error; err != nil {
		return &datas, &info, err
	}

	for i, v := range datas {
		datas[i].UserCreatedString = v.UserCreated.Name
		datas[i].UserModifiedString = &v.UserModified.Name
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

func (r *company) FindByID(ctx *abstraction.Context, id *int) (*model.CompanyEntityModel, error) {
	conn := r.CheckTrx(ctx)

	var data model.CompanyEntityModel
	if err := conn.Where("id = ?", &id).Preload("UserCreated").Preload("UserModified").First(&data).WithContext(ctx.Request().Context()).Error; err != nil {
		return &data, err
	}
	return &data, nil
}

func (r *company) Create(ctx *abstraction.Context, e *model.CompanyEntityModel) (*model.CompanyEntityModel, error) {
	conn := r.CheckTrx(ctx)

	if err := conn.Create(e).WithContext(ctx.Request().Context()).Error; err != nil {
		return nil, err
	}
	if err := conn.Model(e).Preload("UserCreated").Preload("UserModified").First(e).WithContext(ctx.Request().Context()).Error; err != nil {
		return nil, err
	}

	return e, nil
}

func (r *company) Update(ctx *abstraction.Context, id *int, e *model.CompanyEntityModel) (*model.CompanyEntityModel, error) {
	conn := r.CheckTrx(ctx)

	if err := conn.Model(e).Where("id = ?", &id).Updates(e).WithContext(ctx.Request().Context()).Error; err != nil {
		return nil, err
	}
	if err := conn.Where("id = ?", &id).Preload("ParentCompany").First(e).WithContext(ctx.Request().Context()).Error; err != nil {
		return nil, err
	}
	return e, nil 
}

func (r *company) Delete(ctx *abstraction.Context, id *int, e *model.CompanyEntityModel) (*model.CompanyEntityModel, error) {
	conn := r.CheckTrx(ctx)

	if err := conn.Model(e).Where("id =?", id).Delete(e).WithContext(ctx.Request().Context()).Error; err != nil {
		return nil, err
	}
	return e, nil
}

func FilterHasAccessedCompany(ctx *abstraction.Context, db *gorm.DB) {
	r := NewCompany(db)
	conn := r.CheckTrx(ctx)
	_ = conn.Model(&model.AccessScopeEntityModel{}).Where("user_id = ?", ctx.Auth.ID)
}

func (r *company) FindFilterList(ctx *abstraction.Context) (*[]model.CompanyEntityModel, error) {
	conn := r.CheckTrx(ctx)
	var accessScope model.AccessScopeEntityModel
	if err := conn.Model(&model.AccessScopeEntityModel{}).Where("user_id = ?", ctx.Auth.ID).First(&accessScope).Error; err != nil {
		return nil, err
	}

	if accessScope.AccessAll != nil && *accessScope.AccessAll {
		var datas []model.CompanyEntityModel
		if err := conn.Find(&datas).WithContext(ctx.Request().Context()).Error; err != nil {
			return nil, err
		}

		return &datas, nil
	}

	var data []string
	if err := conn.Model(&model.TrialBalanceEntityModel{}).Where("status != 0").Where("created_by = ?", ctx.Auth.ID).Group("company_id").Pluck("company_id", &data).Error; err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	}

	if len(data) == 0 {
		return &[]model.CompanyEntityModel{}, nil
	}

	var data2 []string
	if err := conn.Model(&model.AccessScopeDetailEntityModel{}).Where("access_scope_id = (?)").Pluck("company_id", &data2).Error; err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	}

	if len(data) == 0 {
		return &[]model.CompanyEntityModel{}, nil
	}

	data = append(data, data2...)
	listID := strings.Join(data, ",")

	var datas []model.CompanyEntityModel
	if err := conn.Where(fmt.Sprintf("id IN (%s)", listID)).Find(&datas).WithContext(ctx.Request().Context()).Error; err != nil {
		return nil, err
	}
	return &datas, nil
}
