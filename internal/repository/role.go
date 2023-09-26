package repository

import (
	"fmt"
	"mcash-finance-console-core/internal/abstraction"
	"mcash-finance-console-core/internal/model"

	"gorm.io/gorm"
)

type Role interface {
	Find(ctx *abstraction.Context, m *model.RoleFilterModel, p *abstraction.Pagination) (*[]model.RoleEntityModel, *abstraction.PaginationInfo, error)
	FindByID(ctx *abstraction.Context, id *int) (*model.RoleEntityModel, error)
	Create(ctx *abstraction.Context, e *model.RoleEntityModel) (*model.RoleEntityModel, error)
	Update(ctx *abstraction.Context, id *int, e *model.RoleEntityModel) (*model.RoleEntityModel, error)
	Delete(ctx *abstraction.Context, id *int, e *model.RoleEntityModel) (*model.RoleEntityModel, error)
	DeleteRelatedRole(ctx *abstraction.Context, roleID *int) error
}

type role struct {
	abstraction.Repository
}

func NewRole(db *gorm.DB) *role {
	return &role{
		abstraction.Repository{
			Db: db,
		},
	}
}

func (r *role) Find(ctx *abstraction.Context, m *model.RoleFilterModel, p *abstraction.Pagination) (*[]model.RoleEntityModel, *abstraction.PaginationInfo, error) {
	conn := r.CheckTrx(ctx)

	var datas []model.RoleEntityModel
	var info abstraction.PaginationInfo

	query := conn.Model(&model.RoleEntityModel{})
	//filter
	query = r.Filter(ctx, query, *m)

	if m.Search != nil {
		query = query.Where("code ILIKE ? OR name ILIKE ?", "%"+*m.Search+"%", "%"+*m.Search+"%")
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

	if err := query.Preload("RolePermission.PermissionDef").Preload("UserCreated").Preload("UserModified").Find(&datas).WithContext(ctx.Request().Context()).Error; err != nil {
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

func (r *role) FindByID(ctx *abstraction.Context, id *int) (*model.RoleEntityModel, error) {
	conn := r.CheckTrx(ctx)

	var data model.RoleEntityModel
	if err := conn.Where("id = ?", &id).Preload("UserCreated").Preload("RolePermission").Preload("UserModified").First(&data).WithContext(ctx.Request().Context()).Error; err != nil {
		return &data, err
	}
	data.UserCreatedString = data.UserCreated.Name
	data.UserModifiedString = &data.UserModified.Name
	return &data, nil
}

func (r *role) Create(ctx *abstraction.Context, e *model.RoleEntityModel) (*model.RoleEntityModel, error) {
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

func (r *role) Update(ctx *abstraction.Context, id *int, e *model.RoleEntityModel) (*model.RoleEntityModel, error) {
	conn := r.CheckTrx(ctx)

	if err := conn.Model(e).Where("id = ?", &id).Updates(e).WithContext(ctx.Request().Context()).Error; err != nil {
		return nil, err
	}

	if err := conn.Model(e).Where("id = ?", &id).Preload("UserCreated").Preload("UserModified").First(e).WithContext(ctx.Request().Context()).Error; err != nil {
		return nil, err
	}
	e.UserCreatedString = e.UserCreated.Name
	e.UserModifiedString = &e.UserModified.Name

	return e, nil
}

func (r *role) Delete(ctx *abstraction.Context, id *int, e *model.RoleEntityModel) (*model.RoleEntityModel, error) {
	conn := r.CheckTrx(ctx)

	if err := conn.Model(e).Where("id =?", id).Delete(e).WithContext(ctx.Request().Context()).Error; err != nil {
		return nil, err
	}
	return e, nil
}

func (r *role) DeleteRelatedRole(ctx *abstraction.Context, roleID *int) error {
	conn := r.CheckTrx(ctx)
	if err := conn.Model(&model.RolePermissionApiEntityModel{}).Where("role_id = ?", roleID).Delete(&model.RolePermissionApiEntityModel{}).WithContext(ctx.Request().Context()).Error; err != nil {
		return err
	}

	if err := conn.Model(&model.RolePermissionEntityModel{}).Where("role_id = ?", roleID).Delete(&model.RolePermissionEntityModel{}).WithContext(ctx.Request().Context()).Error; err != nil {
		return err
	}

	return nil
}
