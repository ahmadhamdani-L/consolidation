package repository

import (
	"fmt"
	"mcash-finance-console-core/internal/abstraction"
	"mcash-finance-console-core/internal/model"

	"gorm.io/gorm"
)

type InvestasiNonTbkDetail interface {
	Find(ctx *abstraction.Context, m *model.InvestasiNonTbkDetailFilterModel, p *abstraction.Pagination) (*[]model.InvestasiNonTbkDetailEntityModel, *abstraction.PaginationInfo, error)
	FindByID(ctx *abstraction.Context, id *int) (*model.InvestasiNonTbkDetailEntityModel, error)
	Create(ctx *abstraction.Context, e *model.InvestasiNonTbkDetailEntityModel) (*model.InvestasiNonTbkDetailEntityModel, error)
	Update(ctx *abstraction.Context, id *int, e *model.InvestasiNonTbkDetailEntityModel) (*model.InvestasiNonTbkDetailEntityModel, error)
	Delete(ctx *abstraction.Context, id *int, e *model.InvestasiNonTbkDetailEntityModel) (*model.InvestasiNonTbkDetailEntityModel, error)
	FindWithCode(ctx *abstraction.Context, code *string) (*[]model.InvestasiNonTbkDetailEntityModel, error)
}

type investasinontbkdetail struct {
	abstraction.Repository
}

func NewInvestasiNonTbkDetail(db *gorm.DB) *investasinontbkdetail {
	return &investasinontbkdetail{
		abstraction.Repository{
			Db: db,
		},
	}
}

func (r *investasinontbkdetail) Find(ctx *abstraction.Context, m *model.InvestasiNonTbkDetailFilterModel, p *abstraction.Pagination) (*[]model.InvestasiNonTbkDetailEntityModel, *abstraction.PaginationInfo, error) {
	conn := r.CheckTrx(ctx)

	var datas []model.InvestasiNonTbkDetailEntityModel
	var info abstraction.PaginationInfo

	query := conn.Model(&model.InvestasiNonTbkDetailEntityModel{}).Preload("Company")
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

	if err := query.Find(&datas).WithContext(ctx.Request().Context()).Error; err != nil {
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

func (r *investasinontbkdetail) FindByID(ctx *abstraction.Context, id *int) (*model.InvestasiNonTbkDetailEntityModel, error) {
	conn := r.CheckTrx(ctx)

	var data model.InvestasiNonTbkDetailEntityModel
	if err := conn.Where("id = ?", &id).First(&data).WithContext(ctx.Request().Context()).Error; err != nil {
		return &data, err
	}
	return &data, nil
}

func (r *investasinontbkdetail) Create(ctx *abstraction.Context, e *model.InvestasiNonTbkDetailEntityModel) (*model.InvestasiNonTbkDetailEntityModel, error) {
	conn := r.CheckTrx(ctx)

	if err := conn.Create(e).WithContext(ctx.Request().Context()).Error; err != nil {
		return nil, err
	}
	if err := conn.Model(e).First(e).WithContext(ctx.Request().Context()).Error; err != nil {
		return nil, err
	}

	return e, nil
}

func (r *investasinontbkdetail) Update(ctx *abstraction.Context, id *int, e *model.InvestasiNonTbkDetailEntityModel) (*model.InvestasiNonTbkDetailEntityModel, error) {
	conn := r.CheckTrx(ctx)

	if err := conn.Model(e).Where("id = ?", &id).Updates(e).WithContext(ctx.Request().Context()).Error; err != nil {
		return nil, err
	}
	if err := conn.Model(e).Where("id = ?", &id).First(e).WithContext(ctx.Request().Context()).Error; err != nil {
		return nil, err
	}
	return e, nil

}

func (r *investasinontbkdetail) Delete(ctx *abstraction.Context, id *int, e *model.InvestasiNonTbkDetailEntityModel) (*model.InvestasiNonTbkDetailEntityModel, error) {
	conn := r.CheckTrx(ctx)

	if err := conn.Model(e).Where("id =?", id).Delete(e).WithContext(ctx.Request().Context()).Error; err != nil {
		return nil, err
	}
	return e, nil
}

func (r *investasinontbkdetail) FindWithCode(ctx *abstraction.Context, code *string) (*[]model.InvestasiNonTbkDetailEntityModel, error) {
	conn := r.CheckTrx(ctx)

	var data []model.InvestasiNonTbkDetailEntityModel
	if err := conn.Where("code ILIKE ?", *code+"%").Find(&data).WithContext(ctx.Request().Context()).Error; err != nil {
		return &data, err
	}
	return &data, nil
}
