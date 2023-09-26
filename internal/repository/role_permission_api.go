package repository

import (
	"fmt"
	"mcash-finance-console-core/internal/abstraction"
	"mcash-finance-console-core/internal/model"

	"gorm.io/gorm"
)

type RolePermissionApi interface {
	Find(ctx *abstraction.Context, m *model.RolePermissionApiFilterModel, p *abstraction.Pagination) (*[]model.RolePermissionApiEntityModel, *abstraction.PaginationInfo, error)
	FindByID(ctx *abstraction.Context, id *int) (*model.RolePermissionApiEntityModel, error)
	Create(ctx *abstraction.Context, e *model.RolePermissionApiEntityModel) (*model.RolePermissionApiEntityModel, error)
	FirstOrCreate(ctx *abstraction.Context, e *model.RolePermissionApiEntityModel) (*model.RolePermissionApiEntityModel, error)
	Update(ctx *abstraction.Context, id *int, e *model.RolePermissionApiEntityModel) (*model.RolePermissionApiEntityModel, error)
	Delete(ctx *abstraction.Context, id *int, e *model.RolePermissionApiEntityModel) (*model.RolePermissionApiEntityModel, error)
	DeleteByCriteria(ctx *abstraction.Context, e *model.RolePermissionApiEntityModel) (*model.RolePermissionApiEntityModel, error)
}

type rolepermissionapi struct {
	abstraction.Repository
}

func NewRolePermissionApi(db *gorm.DB) *rolepermissionapi {
	return &rolepermissionapi{
		abstraction.Repository{
			Db: db,
		},
	}
}

func (r *rolepermissionapi) Find(ctx *abstraction.Context, m *model.RolePermissionApiFilterModel, p *abstraction.Pagination) (*[]model.RolePermissionApiEntityModel, *abstraction.PaginationInfo, error) {
	conn := r.CheckTrx(ctx)

	var datas []model.RolePermissionApiEntityModel
	var info abstraction.PaginationInfo

	query := conn.Model(&model.RolePermissionApiEntityModel{})
	//filter
	query = r.Filter(ctx, query, *m)

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

	if err := query.Preload("Role").Find(&datas).WithContext(ctx.Request().Context()).Error; err != nil {
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

func (r *rolepermissionapi) FindByID(ctx *abstraction.Context, id *int) (*model.RolePermissionApiEntityModel, error) {
	conn := r.CheckTrx(ctx)

	var data model.RolePermissionApiEntityModel
	if err := conn.Where("id = ?", &id).Preload("RolePermissionApiPermissionApi").First(&data).WithContext(ctx.Request().Context()).Error; err != nil {
		return &data, err
	}
	return &data, nil
}

func (r *rolepermissionapi) Create(ctx *abstraction.Context, e *model.RolePermissionApiEntityModel) (*model.RolePermissionApiEntityModel, error) {
	conn := r.CheckTrx(ctx)

	if err := conn.Create(e).WithContext(ctx.Request().Context()).Error; err != nil {
		return nil, err
	}

	return e, nil
}

func (r *rolepermissionapi) Update(ctx *abstraction.Context, id *int, e *model.RolePermissionApiEntityModel) (*model.RolePermissionApiEntityModel, error) {
	conn := r.CheckTrx(ctx)

	if err := conn.Model(e).Where("id = ?", &id).Updates(e).WithContext(ctx.Request().Context()).Error; err != nil {
		return nil, err
	}
	return e, nil

}

func (r *rolepermissionapi) Delete(ctx *abstraction.Context, id *int, e *model.RolePermissionApiEntityModel) (*model.RolePermissionApiEntityModel, error) {
	conn := r.CheckTrx(ctx)

	if err := conn.Model(e).Where("id =?", id).Delete(e).WithContext(ctx.Request().Context()).Error; err != nil {
		return nil, err
	}
	return e, nil
}

func (r *rolepermissionapi) FirstOrCreate(ctx *abstraction.Context, e *model.RolePermissionApiEntityModel) (*model.RolePermissionApiEntityModel, error) {
	conn := r.CheckTrx(ctx)

	if err := conn.Where("role_id = ?", e.RoleID).Where("api_path = ?", e.ApiPath).Where("api_method = ?", e.ApiMethod).FirstOrCreate(e).WithContext(ctx.Request().Context()).Error; err != nil {
		return nil, err
	}

	return e, nil
}

func (r *rolepermissionapi) DeleteByCriteria(ctx *abstraction.Context, e *model.RolePermissionApiEntityModel) (*model.RolePermissionApiEntityModel, error) {
	conn := r.CheckTrx(ctx)

	if err := conn.Model(e).Where("role_id = ?", e.RoleID).Where("api_path = ?", e.ApiPath).Where("api_method = ?", e.ApiMethod).Delete(e).WithContext(ctx.Request().Context()).Error; err != nil {
		return nil, err
	}
	return e, nil
}
