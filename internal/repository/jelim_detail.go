package repository

import (
	"fmt"
	"mcash-finance-console-core/internal/abstraction"
	"mcash-finance-console-core/internal/model"

	"gorm.io/gorm"
)

type JelimDetail interface {
	Find(ctx *abstraction.Context, m *model.JelimDetailFilterModel, p *abstraction.Pagination) (*[]model.JelimDetailEntityModel, *abstraction.PaginationInfo, error)
	FindByID(ctx *abstraction.Context, id *int) (*model.JelimDetailEntityModel, error)
	Create(ctx *abstraction.Context, e *model.JelimDetailEntityModel) (*model.JelimDetailEntityModel, error)
	Update(ctx *abstraction.Context, id *int, e *model.JelimDetailEntityModel) (*model.JelimDetailEntityModel, error)
	Delete(ctx *abstraction.Context, id *int, e *model.JelimDetailEntityModel) (*model.JelimDetailEntityModel, error)
	FindWithCode(ctx *abstraction.Context, code *string) (*[]model.JelimDetailEntityModel, error)
	FindWithAjeID(ctx *abstraction.Context, ajeID *int) (*[]model.JelimDetailEntityModel, error)
}

type jelimdetail struct {
	abstraction.Repository
}

func NewJelimDetail(db *gorm.DB) *jelimdetail {
	return &jelimdetail{
		abstraction.Repository{
			Db: db,
		},
	}
}

func (r *jelimdetail) FindWithAjeID(ctx *abstraction.Context, ajeID *int) (*[]model.JelimDetailEntityModel, error) {
	conn := r.CheckTrx(ctx)

	var data []model.JelimDetailEntityModel
	// tmp := fmt.Sprintf("%s", *code)
	if err := conn.Where("jelim_id = ?", ajeID).Find(&data).WithContext(ctx.Request().Context()).Error; err != nil {
		return &data, err
	}
	return &data, nil
}
func (r *jelimdetail) Find(ctx *abstraction.Context, m *model.JelimDetailFilterModel, p *abstraction.Pagination) (*[]model.JelimDetailEntityModel, *abstraction.PaginationInfo, error) {
	conn := r.CheckTrx(ctx)

	var datas []model.JelimDetailEntityModel
	var info abstraction.PaginationInfo

	query := conn.Model(&model.JelimDetailEntityModel{})
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
	var totalData int64
	query = query.Count(&totalData).Limit(limit).Offset(offset)

	if err := query.Preload("Jelim").Find(&datas).WithContext(ctx.Request().Context()).Error; err != nil {
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
func (r *jelimdetail) FindWithCode(ctx *abstraction.Context, code *string) (*[]model.JelimDetailEntityModel, error) {
	conn := r.CheckTrx(ctx)

	var data []model.JelimDetailEntityModel
	tmp := fmt.Sprintf("%s", *code)
	if err := conn.Where("code LIKE ?", tmp+"%").Find(&data).WithContext(ctx.Request().Context()).Error; err != nil {
		return &data, err
	}
	return &data, nil
}
func (r *jelimdetail) FindByID(ctx *abstraction.Context, id *int) (*model.JelimDetailEntityModel, error) {
	conn := r.CheckTrx(ctx)

	var data model.JelimDetailEntityModel
	if err := conn.Where("id = ?", &id).Preload("Jelim").First(&data).WithContext(ctx.Request().Context()).Error; err != nil {
		return &data, err
	}
	return &data, nil
}
func (r *jelimdetail) Create(ctx *abstraction.Context, e *model.JelimDetailEntityModel) (*model.JelimDetailEntityModel, error) {
	conn := r.CheckTrx(ctx)

	if err := conn.Create(e).WithContext(ctx.Request().Context()).Error; err != nil {
		return nil, err
	}
	if err := conn.Model(e).First(e).WithContext(ctx.Request().Context()).Error; err != nil {
		return nil, err
	}

	return e, nil
}
func (r *jelimdetail) Update(ctx *abstraction.Context, id *int, e *model.JelimDetailEntityModel) (*model.JelimDetailEntityModel, error) {
	conn := r.CheckTrx(ctx)

	if err := conn.Model(e).Where("id = ?", &id).Updates(e).Preload("Jelim").WithContext(ctx.Request().Context()).Error; err != nil {
		return nil, err
	}
	if err := conn.Model(e).Where("id = ?", &id).Preload("Jelim").First(e).WithContext(ctx.Request().Context()).Error; err != nil {
		return nil, err
	}
	return e, nil
}
func (r *jelimdetail) Delete(ctx *abstraction.Context, id *int, e *model.JelimDetailEntityModel) (*model.JelimDetailEntityModel, error) {
	conn := r.CheckTrx(ctx)

	if err := conn.Where("jelim_id =?", id).Delete(e).WithContext(ctx.Request().Context()).Error; err != nil {
		return nil, err
	}
	return e, nil
}
