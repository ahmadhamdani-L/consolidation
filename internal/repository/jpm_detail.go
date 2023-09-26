package repository

import (
	"fmt"
	"mcash-finance-console-core/internal/abstraction"
	"mcash-finance-console-core/internal/model"

	"gorm.io/gorm"
)

type JpmDetail interface {
	Find(ctx *abstraction.Context, m *model.JpmDetailFilterModel, p *abstraction.Pagination) (*[]model.JpmDetailEntityModel, *abstraction.PaginationInfo, error)
	FindByID(ctx *abstraction.Context, id *int) (*model.JpmDetailEntityModel, error)
	Create(ctx *abstraction.Context, e *model.JpmDetailEntityModel) (*model.JpmDetailEntityModel, error)
	Update(ctx *abstraction.Context, id *int, e *model.JpmDetailEntityModel) (*model.JpmDetailEntityModel, error)
	Delete(ctx *abstraction.Context, id *int, e *model.JpmDetailEntityModel) (*model.JpmDetailEntityModel, error)
	FindWithCode(ctx *abstraction.Context, code *string) (*[]model.JpmDetailEntityModel, error)
	FindWithAjeID(ctx *abstraction.Context, ajeID *int) (*[]model.JpmDetailEntityModel, error)
}

type jpmdetail struct {
	abstraction.Repository
}

func NewJpmDetail(db *gorm.DB) *jpmdetail {
	return &jpmdetail{
		abstraction.Repository{
			Db: db,
		},
	}
}
func (r *jpmdetail) FindWithAjeID(ctx *abstraction.Context, ajeID *int) (*[]model.JpmDetailEntityModel, error) {
	conn := r.CheckTrx(ctx)

	var data []model.JpmDetailEntityModel
	// tmp := fmt.Sprintf("%s", *code)
	if err := conn.Where("jpm_id = ?", ajeID).Find(&data).WithContext(ctx.Request().Context()).Error; err != nil {
		return &data, err
	}
	return &data, nil
}
func (r *jpmdetail) Find(ctx *abstraction.Context, m *model.JpmDetailFilterModel, p *abstraction.Pagination) (*[]model.JpmDetailEntityModel, *abstraction.PaginationInfo, error) {
	conn := r.CheckTrx(ctx)

	var datas []model.JpmDetailEntityModel
	var info abstraction.PaginationInfo

	query := conn.Model(&model.JpmDetailEntityModel{})
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

	if err := query.Preload("Jpm").Find(&datas).WithContext(ctx.Request().Context()).Error; err != nil {
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
func (r *jpmdetail) FindWithCode(ctx *abstraction.Context, code *string) (*[]model.JpmDetailEntityModel, error) {
	conn := r.CheckTrx(ctx)

	var data []model.JpmDetailEntityModel
	tmp := fmt.Sprintf("%s", *code)
	if err := conn.Where("code LIKE ?", tmp+"%").Find(&data).WithContext(ctx.Request().Context()).Error; err != nil {
		return &data, err
	}
	return &data, nil
}
func (r *jpmdetail) FindByID(ctx *abstraction.Context, id *int) (*model.JpmDetailEntityModel, error) {
	conn := r.CheckTrx(ctx)

	var data model.JpmDetailEntityModel
	if err := conn.Where("id = ?", &id).Preload("Jpm").First(&data).WithContext(ctx.Request().Context()).Error; err != nil {
		return &data, err
	}
	return &data, nil
}
func (r *jpmdetail) Create(ctx *abstraction.Context, e *model.JpmDetailEntityModel) (*model.JpmDetailEntityModel, error) {
	conn := r.CheckTrx(ctx)

	if err := conn.Create(e).WithContext(ctx.Request().Context()).Error; err != nil {
		return nil, err
	}
	if err := conn.Model(e).First(e).WithContext(ctx.Request().Context()).Error; err != nil {
		return nil, err
	}

	return e, nil
}
func (r *jpmdetail) Update(ctx *abstraction.Context, id *int, e *model.JpmDetailEntityModel) (*model.JpmDetailEntityModel, error) {
	conn := r.CheckTrx(ctx)

	if err := conn.Model(e).Where("id = ?", &id).Updates(e).Preload("Jpm").WithContext(ctx.Request().Context()).Error; err != nil {
		return nil, err
	}
	if err := conn.Model(e).Where("id = ?", &id).Preload("Jpm").First(e).WithContext(ctx.Request().Context()).Error; err != nil {
		return nil, err
	}
	return e, nil
}
func (r *jpmdetail) Delete(ctx *abstraction.Context, id *int, e *model.JpmDetailEntityModel) (*model.JpmDetailEntityModel, error) {
	conn := r.CheckTrx(ctx)

	if err := conn.Where("jpm_id =?", id).Delete(e).WithContext(ctx.Request().Context()).Error; err != nil {
		return nil, err
	}
	return e, nil
}
