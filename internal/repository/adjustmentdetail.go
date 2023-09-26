package repository

import (
	"fmt"
	"mcash-finance-console-core/internal/abstraction"
	"mcash-finance-console-core/internal/model"

	"gorm.io/gorm"
)

type AdjustmentDetail interface {
	Find(ctx *abstraction.Context, m *model.AdjustmentDetailFilterModel, p *abstraction.Pagination) (*[]model.AdjustmentDetailEntityModel, *abstraction.PaginationInfo, error)
	FindByID(ctx *abstraction.Context, id *int) (*model.AdjustmentDetailEntityModel, error)
	Create(ctx *abstraction.Context, e *model.AdjustmentDetailEntityModel) (*model.AdjustmentDetailEntityModel, error)
	Update(ctx *abstraction.Context, id *int, e *model.AdjustmentDetailEntityModel) (*model.AdjustmentDetailEntityModel, error)
	Delete(ctx *abstraction.Context, id *int, e *model.AdjustmentDetailEntityModel) (*model.AdjustmentDetailEntityModel, error)
	FindWithCode(ctx *abstraction.Context, code *string) (*[]model.AdjustmentDetailEntityModel, error)
	FindWithAjeID(ctx *abstraction.Context, ajeID *int) (*[]model.AdjustmentDetailEntityModel, error)
	FindWithAjeIDNoArray(ctx *abstraction.Context, ajeID *int) (*model.AdjustmentDetailEntityModel, error)
}

type adjustmentdetail struct {
	abstraction.Repository
}

func NewAdjustmentDetail(db *gorm.DB) *adjustmentdetail {
	return &adjustmentdetail{
		abstraction.Repository{
			Db: db,
		},
	}
}
func (r *adjustmentdetail) FindWithAjeID(ctx *abstraction.Context, ajeID *int) (*[]model.AdjustmentDetailEntityModel, error) {
	conn := r.CheckTrx(ctx)

	var data []model.AdjustmentDetailEntityModel
	// tmp := fmt.Sprintf("%s", *code)
	if err := conn.Where("adjustment_id = ?", ajeID).Find(&data).WithContext(ctx.Request().Context()).Error; err != nil {
		return &data, err
	}
	return &data, nil
}
func (r *adjustmentdetail) FindWithAjeIDNoArray(ctx *abstraction.Context, ajeID *int) (*model.AdjustmentDetailEntityModel, error) {
	conn := r.CheckTrx(ctx)

	var data model.AdjustmentDetailEntityModel
	// tmp := fmt.Sprintf("%s", *code)
	if err := conn.Where("adjustment_id = ?", ajeID).Find(&data).WithContext(ctx.Request().Context()).Error; err != nil {
		return &data, err
	}
	return &data, nil
}
func (r *adjustmentdetail) Find(ctx *abstraction.Context, m *model.AdjustmentDetailFilterModel, p *abstraction.Pagination) (*[]model.AdjustmentDetailEntityModel, *abstraction.PaginationInfo, error) {
	conn := r.CheckTrx(ctx)

	var datas []model.AdjustmentDetailEntityModel
	var info abstraction.PaginationInfo

	query := conn.Model(&model.AdjustmentDetailEntityModel{})
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

	if err := query.Preload("Adjustment").Find(&datas).WithContext(ctx.Request().Context()).Error; err != nil {
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
func (r *adjustmentdetail) FindWithCode(ctx *abstraction.Context, code *string) (*[]model.AdjustmentDetailEntityModel, error) {
	conn := r.CheckTrx(ctx)

	var data []model.AdjustmentDetailEntityModel
	if err := conn.Where("code LIKE ?", *code+"%").Find(&data).WithContext(ctx.Request().Context()).Error; err != nil {
		return &data, err
	}
	return &data, nil
}
func (r *adjustmentdetail) FindByID(ctx *abstraction.Context, id *int) (*model.AdjustmentDetailEntityModel, error) {
	conn := r.CheckTrx(ctx)

	var data model.AdjustmentDetailEntityModel
	if err := conn.Where("id = ?", &id).Preload("Adjustment").Find(&data).WithContext(ctx.Request().Context()).Error; err != nil {
		return &data, err
	}
	return &data, nil
}
func (r *adjustmentdetail) Create(ctx *abstraction.Context, e *model.AdjustmentDetailEntityModel) (*model.AdjustmentDetailEntityModel, error) {
	conn := r.CheckTrx(ctx)

	if err := conn.Create(e).WithContext(ctx.Request().Context()).Error; err != nil {
		return nil, err
	}
	if err := conn.Model(e).First(e).WithContext(ctx.Request().Context()).Error; err != nil {
		return nil, err
	}

	return e, nil
}
func (r *adjustmentdetail) Update(ctx *abstraction.Context, id *int, e *model.AdjustmentDetailEntityModel) (*model.AdjustmentDetailEntityModel, error) {
	conn := r.CheckTrx(ctx)

	if err := conn.Model(e).Where("id = ?", &id).Updates(e).Preload("Adjustment").WithContext(ctx.Request().Context()).Error; err != nil {
		return nil, err
	}
	if err := conn.Model(e).Where("id = ?", &id).Preload("Adjustment").First(e).WithContext(ctx.Request().Context()).Error; err != nil {
		return nil, err
	}
	return e, nil
}
func (r *adjustmentdetail) Delete(ctx *abstraction.Context, id *int, e *model.AdjustmentDetailEntityModel) (*model.AdjustmentDetailEntityModel, error) {
	conn := r.CheckTrx(ctx)

	if err := conn.Where("adjustment_id =?", id).Delete(e).WithContext(ctx.Request().Context()).Error; err != nil {
		return nil, err
	}
	return e, nil
}
