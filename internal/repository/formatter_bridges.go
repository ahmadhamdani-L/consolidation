package repository

import (
	"fmt"
	"mcash-finance-console-core/internal/abstraction"
	"mcash-finance-console-core/internal/model"

	"gorm.io/gorm"
)

type FormatterBridges interface {
	Find(ctx *abstraction.Context, m *model.FormatterBridgesFilterModel, p *abstraction.Pagination) (*[]model.FormatterBridgesEntityModel, *abstraction.PaginationInfo, error)
	FindByID(ctx *abstraction.Context, id *int) (*model.FormatterBridgesEntityModel, error)
	Create(ctx *abstraction.Context, e *model.FormatterBridgesEntityModel) (*model.FormatterBridgesEntityModel, error)
	Update(ctx *abstraction.Context, id *int, e *model.FormatterBridgesEntityModel) (*model.FormatterBridgesEntityModel, error)
	Delete(ctx *abstraction.Context, id *int, e *model.FormatterBridgesEntityModel) (*model.FormatterBridgesEntityModel, error)
	FindWithCriteria(ctx *abstraction.Context, m *model.FormatterBridgesFilterModel) (*model.FormatterBridgesEntityModel, error)
	FindSummary(ctx *abstraction.Context, e *model.FormatterBridgesFilterModel) (*[]model.FormatterDetailEntityModel, error)
	FindSummaryTB(ctx *abstraction.Context) (*[]model.FormatterDetailEntityModel, error)
}

type formatterbridges struct {
	abstraction.Repository
}

func NewFormatterBridges(db *gorm.DB) *formatterbridges {
	return &formatterbridges{
		abstraction.Repository{
			Db: db,
		},
	}
}

func (r *formatterbridges) Find(ctx *abstraction.Context, m *model.FormatterBridgesFilterModel, p *abstraction.Pagination) (*[]model.FormatterBridgesEntityModel, *abstraction.PaginationInfo, error) {
	conn := r.CheckTrx(ctx)

	var datas []model.FormatterBridgesEntityModel
	var info abstraction.PaginationInfo

	query := conn.Model(&model.FormatterBridgesEntityModel{})
	//filter
	query = r.Filter(ctx, query, *m)
	queryUser := conn.Model(&model.UserEntityModel{}).Select("id")
	query = r.FilterUser(ctx, query, queryUser, m.FormatterBridgesFilter, "")

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

	if err := query.Preload("Formatter").Find(&datas).WithContext(ctx.Request().Context()).Error; err != nil {
		return &datas, &info, err
	}

	for i, v := range datas {
		datas[i].UserCreatedString = v.UserCreated.Name
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

func (r *formatterbridges) FindByID(ctx *abstraction.Context, id *int) (*model.FormatterBridgesEntityModel, error) {
	conn := r.CheckTrx(ctx)

	var data model.FormatterBridgesEntityModel
	if err := conn.Where("id = ?", &id).First(&data).WithContext(ctx.Request().Context()).Error; err != nil {
		return &data, err
	}
	return &data, nil
}

func (r *formatterbridges) FindWithCriteria(ctx *abstraction.Context, m *model.FormatterBridgesFilterModel) (*model.FormatterBridgesEntityModel, error) {
	conn := r.CheckTrx(ctx)

	var datas model.FormatterBridgesEntityModel

	query := conn.Model(&model.FormatterBridgesEntityModel{})
	//filter
	query = r.Filter(ctx, query, *m)

	if err := query.Preload("Formatter").Order("created_at ASC").Limit(1).Find(&datas).WithContext(ctx.Request().Context()).Error; err != nil {
		return &datas, err
	}

	return &datas, nil
}

func (r *formatterbridges) Create(ctx *abstraction.Context, e *model.FormatterBridgesEntityModel) (*model.FormatterBridgesEntityModel, error) {
	conn := r.CheckTrx(ctx)

	if err := conn.Create(e).WithContext(ctx.Request().Context()).Error; err != nil {
		return nil, err
	}
	if err := conn.Model(e).First(e).WithContext(ctx.Request().Context()).Error; err != nil {
		return nil, err
	}

	return e, nil
}

func (r *formatterbridges) Update(ctx *abstraction.Context, id *int, e *model.FormatterBridgesEntityModel) (*model.FormatterBridgesEntityModel, error) {
	conn := r.CheckTrx(ctx)

	if err := conn.Model(e).Where("id = ?", &id).Updates(e).WithContext(ctx.Request().Context()).Error; err != nil {
		return nil, err
	}
	if err := conn.Model(e).Where("id = ?", &id).First(e).WithContext(ctx.Request().Context()).Error; err != nil {
		return nil, err
	}
	return e, nil

}

func (r *formatterbridges) Delete(ctx *abstraction.Context, id *int, e *model.FormatterBridgesEntityModel) (*model.FormatterBridgesEntityModel, error) {
	conn := r.CheckTrx(ctx)

	if err := conn.Model(e).Where("id =?", id).Delete(e).WithContext(ctx.Request().Context()).Error; err != nil {
		return nil, err
	}
	return e, nil
}

func (r *formatterbridges) FindSummary(ctx *abstraction.Context, e *model.FormatterBridgesFilterModel) (*[]model.FormatterDetailEntityModel, error) {
	conn := r.CheckTrx(ctx)

	data := []model.FormatterDetailEntityModel{}
	err := conn.Model(&model.FormatterBridgesEntityModel{}).
		Joins("INNER JOIN m_formatter f ON formatter_bridges.formatter_id = f.id").
		Joins("INNER JOIN m_formatter_detail fd ON f.id = fd.formatter_id").
		Where("trx_ref_id = ? AND source = ? AND (auto_summary = true OR (is_total = true AND fx_summary IS NOT NULL))", e.TrxRefID, e.Source).
		Order("fd.formatter_id ASC").
		Order("fd.sort_id ASC").
		Select("fd.*").Find(&data).Error
	if err != nil {
		return nil, err
	}
	return &data, nil
}

func (r *formatterbridges) FindSummaryTB(ctx *abstraction.Context) (*[]model.FormatterDetailEntityModel, error) {
	conn := r.CheckTrx(ctx)

	data := []model.FormatterDetailEntityModel{}

	coa4dst := conn.Model(&model.FormatterDetailEntityModel{}).Where("code = 'PENDAPATAN'").Where("formatter_id = ?", 3).Select("sort_id").Limit(1)

	query1 := conn.Model(&model.FormatterDetailEntityModel{}).
		Where("formatter_id = ? AND (auto_summary = true OR (is_total = true AND fx_summary IS NOT NULL)) AND is_recalculate = true", 3).
		Order("sort_id ASC").Where("sort_id >= (?)", coa4dst)
	query2 := conn.Model(&model.FormatterDetailEntityModel{}).
		Where("formatter_id = ? AND (auto_summary = true OR (is_total = true AND fx_summary IS NOT NULL)) AND is_recalculate = true", 3).
		Order("sort_id ASC").Where("sort_id < (?)", coa4dst)

	err := query1.Find(&data).Error
	if err != nil {
		return nil, err
	}
	tmp := []model.FormatterDetailEntityModel{}
	err = query2.Find(&tmp).Error
	if err != nil {
		return nil, err
	}

	data = append(data, tmp...)

	return &data, nil
}

// -64309649
