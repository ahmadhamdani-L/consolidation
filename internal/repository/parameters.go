package repository

import (
	"fmt"
	"mcash-finance-console-core/internal/abstraction"
	"mcash-finance-console-core/internal/model"

	"gorm.io/gorm"
)

type Parameter interface {
	Find(ctx *abstraction.Context, m *model.ParameterFilterModel, p *abstraction.Pagination) (*[]model.ParameterEntityModel, *abstraction.PaginationInfo, error)
	FindByID(ctx *abstraction.Context, id *int) (*model.ParameterEntityModel, error)
	Create(ctx *abstraction.Context, e *model.ParameterEntityModel) (*model.ParameterEntityModel, error)
	Update(ctx *abstraction.Context, id *int, e *model.ParameterEntityModel) (*model.ParameterEntityModel, error)
	Delete(ctx *abstraction.Context, id *int, e *model.ParameterEntityModel) (*model.ParameterEntityModel, error)
	FindWithCode(ctx *abstraction.Context, code *string) (*[]model.ParameterEntityModel, error)
}

type parameter struct {
	abstraction.Repository
}

func NewParameter(db *gorm.DB) *parameter {
	return &parameter{
		abstraction.Repository{
			Db: db,
		},
	}
}

func (r *parameter) Find(ctx *abstraction.Context, m *model.ParameterFilterModel, p *abstraction.Pagination) (*[]model.ParameterEntityModel, *abstraction.PaginationInfo, error) {
	conn := r.CheckTrx(ctx)

	var datas []model.ParameterEntityModel
	var info abstraction.PaginationInfo

	query := conn.Model(&model.ParameterEntityModel{})
	//filter
	query = r.Filter(ctx, query, *m)
	queryUser := conn.Model(&model.UserEntityModel{}).Select("id")
	query = r.FilterUser(ctx, query, queryUser, m.Filter, "")
	if m.Search != nil {
		query = query.Where("(parameters.code ILIKE ? OR parameters.value ILIKE ?)", "%"+*m.Search+"%", "%"+*m.Search+"%")
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

	if err := query.Find(&datas).WithContext(ctx.Request().Context()).Error; err != nil {
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

func (r *parameter) FindByID(ctx *abstraction.Context, id *int) (*model.ParameterEntityModel, error) {
	conn := r.CheckTrx(ctx)

	var data model.ParameterEntityModel
	if err := conn.Where("id = ?", &id).First(&data).WithContext(ctx.Request().Context()).Error; err != nil {
		return &data, err
	}
	return &data, nil
}

func (r *parameter) FindWithCode(ctx *abstraction.Context, code *string) (*[]model.ParameterEntityModel, error) {
	conn := r.CheckTrx(ctx)

	var data []model.ParameterEntityModel
	if err := conn.Where("code LIKE ?", *code+"%").Find(&data).WithContext(ctx.Request().Context()).Error; err != nil {
		return &data, err
	}
	return &data, nil
}

func (r *parameter) Create(ctx *abstraction.Context, e *model.ParameterEntityModel) (*model.ParameterEntityModel, error) {
	conn := r.CheckTrx(ctx)

	if err := conn.Create(e).WithContext(ctx.Request().Context()).Error; err != nil {
		return nil, err
	}
	if err := conn.Model(e).First(e).WithContext(ctx.Request().Context()).Error; err != nil {
		return nil, err
	}

	return e, nil
}

func (r *parameter) Update(ctx *abstraction.Context, id *int, e *model.ParameterEntityModel) (*model.ParameterEntityModel, error) {
	conn := r.CheckTrx(ctx)

	if err := conn.Model(e).Where("id = ?", &id).Updates(e).WithContext(ctx.Request().Context()).Error; err != nil {
		return nil, err
	}
	if err := conn.Model(e).Where("id = ?", &id).First(e).WithContext(ctx.Request().Context()).Error; err != nil {
		return nil, err
	}
	return e, nil

}

func (r *parameter) Delete(ctx *abstraction.Context, id *int, e *model.ParameterEntityModel) (*model.ParameterEntityModel, error) {
	conn := r.CheckTrx(ctx)

	if err := conn.Model(e).Where("id =?", id).Delete(e).WithContext(ctx.Request().Context()).Error; err != nil {
		return nil, err
	}
	return e, nil
}
