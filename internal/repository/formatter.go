package repository

import (
	"fmt"
	"mcash-finance-console-core/internal/abstraction"
	"mcash-finance-console-core/internal/model"

	"gorm.io/gorm"
)

type Formatter interface {
	Find(ctx *abstraction.Context, m *model.FormatterFilterModel, p *abstraction.Pagination) (*[]model.FormatterEntityModel, *abstraction.PaginationInfo, error)
	FindByID(ctx *abstraction.Context, id *int) (*model.FormatterEntityModel, error)
	Create(ctx *abstraction.Context, e *model.FormatterEntityModel) (*model.FormatterEntityModel, error)
	Update(ctx *abstraction.Context, id *int, e *model.FormatterEntityModel) (*model.FormatterEntityModel, error)
	Delete(ctx *abstraction.Context, id *int, e *model.FormatterEntityModel) (*model.FormatterEntityModel, error)
	FindWithDetail(ctx *abstraction.Context, m *model.FormatterFilterModel) (*model.FormatterEntityModel, error)
}

type formatter struct {
	abstraction.Repository
}

func NewFormatter(db *gorm.DB) *formatter {
	return &formatter{
		abstraction.Repository{
			Db: db,
		},
	}
}

func (r *formatter) Find(ctx *abstraction.Context, m *model.FormatterFilterModel, p *abstraction.Pagination) (*[]model.FormatterEntityModel, *abstraction.PaginationInfo, error) {
	conn := r.CheckTrx(ctx)

	var datas []model.FormatterEntityModel
	var info abstraction.PaginationInfo

	query := conn.Model(&model.FormatterEntityModel{})
	//filter
	query = r.Filter(ctx, query, *m)
	queryUser := conn.Model(&model.UserEntityModel{}).Select("id")
	query = r.FilterUser(ctx, query, queryUser, m.Filter, "")

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

func (r *formatter) FindWithDetail(ctx *abstraction.Context, m *model.FormatterFilterModel) (*model.FormatterEntityModel, error) {
	conn := r.CheckTrx(ctx)
	var data model.FormatterEntityModel

	if err := conn.Where("formatter_for", &m.FormatterFor).Preload("FormatterDetail", func(db *gorm.DB) *gorm.DB {
		return db.Order("sort_id asc")
	}).Preload("UserCreated").Preload("UserModified").First(&data).Error; err != nil {
		return &data, err
	}
	return &data, nil
}

func (r *formatter) FindByID(ctx *abstraction.Context, id *int) (*model.FormatterEntityModel, error) {
	conn := r.CheckTrx(ctx)

	var data model.FormatterEntityModel
	if err := conn.Where("id = ?", &id).Preload("UserCreated").Preload("UserModified").First(&data).WithContext(ctx.Request().Context()).Error; err != nil {
		return &data, err
	}
	return &data, nil
}

func (r *formatter) Create(ctx *abstraction.Context, e *model.FormatterEntityModel) (*model.FormatterEntityModel, error) {
	conn := r.CheckTrx(ctx)

	if err := conn.Create(e).WithContext(ctx.Request().Context()).Error; err != nil {
		return nil, err
	}
	if err := conn.Model(e).First(e).WithContext(ctx.Request().Context()).Error; err != nil {
		return nil, err
	}

	return e, nil
}

func (r *formatter) Update(ctx *abstraction.Context, id *int, e *model.FormatterEntityModel) (*model.FormatterEntityModel, error) {
	conn := r.CheckTrx(ctx)

	if err := conn.Model(e).Where("id = ?", &id).Updates(e).WithContext(ctx.Request().Context()).Error; err != nil {
		return nil, err
	}
	if err := conn.Model(e).Where("id = ?", &id).First(e).WithContext(ctx.Request().Context()).Error; err != nil {
		return nil, err
	}
	return e, nil

}

func (r *formatter) Delete(ctx *abstraction.Context, id *int, e *model.FormatterEntityModel) (*model.FormatterEntityModel, error) {
	conn := r.CheckTrx(ctx)

	if err := conn.Model(e).Where("id =?", id).Delete(e).WithContext(ctx.Request().Context()).Error; err != nil {
		return nil, err
	}
	return e, nil
}
