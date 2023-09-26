package repository

import (
	"fmt"
	"mcash-finance-console-core/internal/abstraction"
	"mcash-finance-console-core/internal/model"

	"gorm.io/gorm"
)

type JcteDetail interface {
	Find(ctx *abstraction.Context, m *model.JcteDetailFilterModel, p *abstraction.Pagination) (*[]model.JcteDetailEntityModel, *abstraction.PaginationInfo, error)
	FindByID(ctx *abstraction.Context, id *int) (*model.JcteDetailEntityModel, error)
	Create(ctx *abstraction.Context, e *model.JcteDetailEntityModel) (*model.JcteDetailEntityModel, error)
	Update(ctx *abstraction.Context, id *int, e *model.JcteDetailEntityModel) (*model.JcteDetailEntityModel, error)
	Delete(ctx *abstraction.Context, id *int, e *model.JcteDetailEntityModel) (*model.JcteDetailEntityModel, error)
	FindWithCode(ctx *abstraction.Context, code *string) (*[]model.JcteDetailEntityModel, error)
	FindWithAjeID(ctx *abstraction.Context, ajeID *int) (*[]model.JcteDetailEntityModel, error)
}

type jctedetail struct {
	abstraction.Repository
}

func NewJcteDetail(db *gorm.DB) *jctedetail {
	return &jctedetail{
		abstraction.Repository{
			Db: db,
		},
	}
}
func (r *jctedetail) FindWithAjeID(ctx *abstraction.Context, ajeID *int) (*[]model.JcteDetailEntityModel, error) {
	conn := r.CheckTrx(ctx)

	var data []model.JcteDetailEntityModel
	// tmp := fmt.Sprintf("%s", *code)
	if err := conn.Where("jcte_id = ?", ajeID).Find(&data).WithContext(ctx.Request().Context()).Error; err != nil {
		return &data, err
	}
	return &data, nil
}
func (r *jctedetail) Find(ctx *abstraction.Context, m *model.JcteDetailFilterModel, p *abstraction.Pagination) (*[]model.JcteDetailEntityModel, *abstraction.PaginationInfo, error) {
	conn := r.CheckTrx(ctx)

	var datas []model.JcteDetailEntityModel
	var info abstraction.PaginationInfo

	query := conn.Model(&model.JcteDetailEntityModel{})
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

	if err := query.Preload("Jcte").Find(&datas).WithContext(ctx.Request().Context()).Error; err != nil {
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
func (r *jctedetail) FindWithCode(ctx *abstraction.Context, code *string) (*[]model.JcteDetailEntityModel, error) {
	conn := r.CheckTrx(ctx)

	var data []model.JcteDetailEntityModel
	tmp := fmt.Sprintf("%s", *code)
	if err := conn.Where("code LIKE ?", tmp+"%").Find(&data).WithContext(ctx.Request().Context()).Error; err != nil {
		return &data, err
	}
	return &data, nil
}
func (r *jctedetail) FindByID(ctx *abstraction.Context, id *int) (*model.JcteDetailEntityModel, error) {
	conn := r.CheckTrx(ctx)

	var data model.JcteDetailEntityModel
	if err := conn.Where("id = ?", &id).Preload("Jcte").First(&data).WithContext(ctx.Request().Context()).Error; err != nil {
		return &data, err
	}
	return &data, nil
}
func (r *jctedetail) Create(ctx *abstraction.Context, e *model.JcteDetailEntityModel) (*model.JcteDetailEntityModel, error) {
	conn := r.CheckTrx(ctx)

	if err := conn.Create(e).WithContext(ctx.Request().Context()).Error; err != nil {
		return nil, err
	}
	if err := conn.Model(e).First(e).WithContext(ctx.Request().Context()).Error; err != nil {
		return nil, err
	}

	return e, nil
}
func (r *jctedetail) Update(ctx *abstraction.Context, id *int, e *model.JcteDetailEntityModel) (*model.JcteDetailEntityModel, error) {
	conn := r.CheckTrx(ctx)

	if err := conn.Model(e).Where("id = ?", &id).Updates(e).Preload("Jcte").WithContext(ctx.Request().Context()).Error; err != nil {
		return nil, err
	}
	if err := conn.Model(e).Where("id = ?", &id).Preload("Jcte").First(e).WithContext(ctx.Request().Context()).Error; err != nil {
		return nil, err
	}
	return e, nil
}
func (r *jctedetail) Delete(ctx *abstraction.Context, id *int, e *model.JcteDetailEntityModel) (*model.JcteDetailEntityModel, error) {
	conn := r.CheckTrx(ctx)

	if err := conn.Where("jcte_id =?", id).Delete(e).WithContext(ctx.Request().Context()).Error; err != nil {
		return nil, err
	}
	return e, nil
}
