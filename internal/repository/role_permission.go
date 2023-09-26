package repository

import (
	"fmt"
	"mcash-finance-console-core/internal/abstraction"
	"mcash-finance-console-core/internal/model"

	"gorm.io/gorm"
)

type RolePermission interface {
	Find(ctx *abstraction.Context, m *model.RolePermissionFilterModel, p *abstraction.Pagination) (*[]model.RolePermissionEntityModel, *abstraction.PaginationInfo, error)
	FindByID(ctx *abstraction.Context, id *int) (*model.RolePermissionEntityModel, error)
	FindByFunctionalID(ctx *abstraction.Context, functionalID *string) (*model.RolePermissionEntityModel, error)
	FirstOrCreate(ctx *abstraction.Context, e *model.RolePermissionEntityModel) (*model.RolePermissionEntityModel, error)
	Create(ctx *abstraction.Context, e *model.RolePermissionEntityModel) (*model.RolePermissionEntityModel, error)
	Update(ctx *abstraction.Context, id *int, e *model.RolePermissionEntityModel) (*model.RolePermissionEntityModel, error)
	UpdateByRoleFunctionalID(ctx *abstraction.Context, roleID *int, functionalID *string, e *model.RolePermissionEntityModel) (*model.RolePermissionEntityModel, error)
	Delete(ctx *abstraction.Context, id *int, e *model.RolePermissionEntityModel) (*model.RolePermissionEntityModel, error)
	DeleteByRoleFunctionalID(ctx *abstraction.Context, e *model.RolePermissionEntityModel) (*model.RolePermissionEntityModel, error)
	FindsByCriteria(ctx *abstraction.Context, m *model.RolePermissionFilterModel) (*[]model.RolePermissionEntityModel, error)
}

type rolepermission struct {
	abstraction.Repository
}

func NewRolePermission(db *gorm.DB) *rolepermission {
	return &rolepermission{
		abstraction.Repository{
			Db: db,
		},
	}
}

func (r *rolepermission) Find(ctx *abstraction.Context, m *model.RolePermissionFilterModel, p *abstraction.Pagination) (*[]model.RolePermissionEntityModel, *abstraction.PaginationInfo, error) {
	conn := r.CheckTrx(ctx)

	var datas []model.RolePermissionEntityModel
	var info abstraction.PaginationInfo

	query := conn.Model(&model.RolePermissionEntityModel{})
	//filter
	query = r.Filter(ctx, query, *m)

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

func (r *rolepermission) FindByID(ctx *abstraction.Context, id *int) (*model.RolePermissionEntityModel, error) {
	conn := r.CheckTrx(ctx)

	var data model.RolePermissionEntityModel
	if err := conn.Where("id = ?", &id).Preload("Role").First(&data).WithContext(ctx.Request().Context()).Error; err != nil {
		return &data, err
	}
	return &data, nil
}

func (r *rolepermission) FindByFunctionalID(ctx *abstraction.Context, functionalID *string) (*model.RolePermissionEntityModel, error) {
	conn := r.CheckTrx(ctx)

	var data model.RolePermissionEntityModel
	if err := conn.Where("functional_id = ?", &functionalID).First(&data).WithContext(ctx.Request().Context()).Error; err != nil {
		return &data, err
	}
	return &data, nil
}

func (r *rolepermission) Create(ctx *abstraction.Context, e *model.RolePermissionEntityModel) (*model.RolePermissionEntityModel, error) {
	conn := r.CheckTrx(ctx)

	if err := conn.Table(e.TableName()).Create(e.RolePermissionEntity).Error; err != nil {
		return nil, err
	}

	return e, nil
}

func (r *rolepermission) Update(ctx *abstraction.Context, id *int, e *model.RolePermissionEntityModel) (*model.RolePermissionEntityModel, error) {
	conn := r.CheckTrx(ctx)

	if err := conn.Table(e.TableName()).Where("id = ?", &id).Updates(e.RolePermissionEntity).Error; err != nil {
		return nil, err
	}
	return e, nil

}

func (r *rolepermission) FirstOrCreate(ctx *abstraction.Context, e *model.RolePermissionEntityModel) (*model.RolePermissionEntityModel, error) {
	conn := r.CheckTrx(ctx)

	if err := conn.Where("functional_id = ?", e.FunctionalID).Where("role_id = ?", e.RoleID).FirstOrCreate(e).WithContext(ctx.Request().Context()).Error; err != nil {
		return nil, err
	}

	return e, nil
}

func (r *rolepermission) UpdateByRoleFunctionalID(ctx *abstraction.Context, roleID *int, functionalID *string, e *model.RolePermissionEntityModel) (*model.RolePermissionEntityModel, error) {
	conn := r.CheckTrx(ctx)

	if err := conn.Table(e.TableName()).Where("role_id = ?", &roleID).Where("functional_id = ?", &functionalID).WithContext(ctx.Request().Context()).Updates(e.RolePermissionEntity).Error; err != nil {
		return nil, err
	}
	return e, nil

}

func (r *rolepermission) Delete(ctx *abstraction.Context, id *int, e *model.RolePermissionEntityModel) (*model.RolePermissionEntityModel, error) {
	conn := r.CheckTrx(ctx)

	if err := conn.Model(e).Where("id = ?", id).Delete(e).WithContext(ctx.Request().Context()).Error; err != nil {
		return nil, err
	}
	return e, nil
}

func (r *rolepermission) DeleteByRoleFunctionalID(ctx *abstraction.Context, e *model.RolePermissionEntityModel) (*model.RolePermissionEntityModel, error) {
	conn := r.CheckTrx(ctx)

	if err := conn.Model(e).Where("role_id = ?", &e.RoleID).Where("functional_id = ?", &e.FunctionalID).Delete(e).WithContext(ctx.Request().Context()).Error; err != nil {
		return nil, err
	}
	return e, nil
}

// create function FindsByCriteria for searching data by role_id
func (r *rolepermission) FindsByCriteria(ctx *abstraction.Context, m *model.RolePermissionFilterModel) (*[]model.RolePermissionEntityModel, error) {
	conn := r.CheckTrx(ctx)
	query := conn.Model(&model.RolePermissionEntityModel{})
	//filter
	query = r.Filter(ctx, query, *m)
	var datas []model.RolePermissionEntityModel
	if err := query.Find(&datas).WithContext(ctx.Request().Context()).Error; err != nil {
		return nil, err
	}
	return &datas, nil
}
