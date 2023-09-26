package repository

import (
	"fmt"
	"mcash-finance-console-core/internal/abstraction"
	"mcash-finance-console-core/internal/model"

	"gorm.io/gorm"
)

type PermissionDef interface {
	Find(ctx *abstraction.Context, m *model.PermissionDefFilterModel, p *abstraction.Pagination) (*[]model.PermissionDefEntityModel, *abstraction.PaginationInfo, error)
	FindByID(ctx *abstraction.Context, id *int) (*model.PermissionDefEntityModel, error)
	FindByFunctionalID(ctx *abstraction.Context, functionalID string) (*model.PermissionDefEntityModel, error)
	Create(ctx *abstraction.Context, e *model.PermissionDefEntityModel) (*model.PermissionDefEntityModel, error)
	Update(ctx *abstraction.Context, id *int, e *model.PermissionDefEntityModel) (*model.PermissionDefEntityModel, error)
	Delete(ctx *abstraction.Context, id *int, e *model.PermissionDefEntityModel) (*model.PermissionDefEntityModel, error)
}

type permissiondef struct {
	abstraction.Repository
}

func NewPermissionDef(db *gorm.DB) *permissiondef {
	return &permissiondef{
		abstraction.Repository{
			Db: db,
		},
	}
}

func (r *permissiondef) Find(ctx *abstraction.Context, m *model.PermissionDefFilterModel, p *abstraction.Pagination) (*[]model.PermissionDefEntityModel, *abstraction.PaginationInfo, error) {
	conn := r.CheckTrx(ctx)

	var datas []model.PermissionDefEntityModel
	var info abstraction.PaginationInfo

	query := conn.Model(&model.PermissionDefEntityModel{})
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

	if err := query.Preload("UserCreated").Preload("UserModified").Find(&datas).WithContext(ctx.Request().Context()).Error; err != nil {
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

func (r *permissiondef) FindByID(ctx *abstraction.Context, id *int) (*model.PermissionDefEntityModel, error) {
	conn := r.CheckTrx(ctx)

	var data model.PermissionDefEntityModel
	if err := conn.Where("id = ?", &id).Preload("UserCreated").Preload("UserModified").First(&data).WithContext(ctx.Request().Context()).Error; err != nil {
		return &data, err
	}
	data.UserCreatedString = data.UserCreated.Name
	data.UserModifiedString = &data.UserModified.Name
	return &data, nil
}

func (r *permissiondef) Create(ctx *abstraction.Context, e *model.PermissionDefEntityModel) (*model.PermissionDefEntityModel, error) {
	conn := r.CheckTrx(ctx)

	if err := conn.Create(e).WithContext(ctx.Request().Context()).Error; err != nil {
		return nil, err
	}

	return e, nil
}

func (r *permissiondef) Update(ctx *abstraction.Context, id *int, e *model.PermissionDefEntityModel) (*model.PermissionDefEntityModel, error) {
	conn := r.CheckTrx(ctx)

	if err := conn.Model(e).Where("id = ?", &id).Updates(e).WithContext(ctx.Request().Context()).Error; err != nil {
		return nil, err
	}
	return e, nil

}

func (r *permissiondef) Delete(ctx *abstraction.Context, id *int, e *model.PermissionDefEntityModel) (*model.PermissionDefEntityModel, error) {
	conn := r.CheckTrx(ctx)

	if err := conn.Model(e).Where("id =?", id).Delete(e).WithContext(ctx.Request().Context()).Error; err != nil {
		return nil, err
	}
	return e, nil
}

func (r *permissiondef) FindByFunctionalID(ctx *abstraction.Context, functionalID string) (*model.PermissionDefEntityModel, error) {
	conn := r.CheckTrx(ctx)

	var data model.PermissionDefEntityModel
	if err := conn.Where("functional_id = ?", &functionalID).Preload("UserCreated").Preload("UserModified").First(&data).WithContext(ctx.Request().Context()).Error; err != nil {
		return &data, err
	}
	data.UserCreatedString = data.UserCreated.Name
	data.UserModifiedString = &data.UserModified.Name
	return &data, nil
}
