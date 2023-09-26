package repository

import (
	"fmt"
	"worker/internal/abstraction"
	"worker/internal/model"

	"gorm.io/gorm"
)

type AdjustmentDetail interface {
	Find(ctx *abstraction.Context, m *model.AdjustmentDetailFilterModel, p *abstraction.Pagination) (*[]model.AdjustmentDetailEntityModel, *abstraction.PaginationInfo, error)
	FindByID(ctx *abstraction.Context, id *int) (*model.AdjustmentDetailEntityModel, error)
	FindWithCode(ctx *abstraction.Context, code *string) (*[]model.AdjustmentDetailEntityModel, error)
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
	limit := *p.PageSize + 1
	offset := limit * (*p.Page - 1)
	query = query.Limit(limit).Offset(offset)

	if err := query.Preload("Adjustment").Find(&datas).Error; err != nil {
		return &datas, &info, err
	}

	info.Count = len(datas)
	info.MoreRecords = false
	if len(datas) > *p.PageSize {
		info.MoreRecords = true
		info.Count -= 1
		datas = datas[:len(datas)-1]
	}

	return &datas, &info, nil
}

func (r *adjustmentdetail) FindWithCode(ctx *abstraction.Context, code *string) (*[]model.AdjustmentDetailEntityModel, error) {
	conn := r.CheckTrx(ctx)

	var data []model.AdjustmentDetailEntityModel
	tmp := fmt.Sprintf("%s", *code)
	if err := conn.Where("code LIKE ?", tmp+"%").Find(&data).Error; err != nil {
		return &data, err
	}
	return &data, nil
}

func (r *adjustmentdetail) FindByID(ctx *abstraction.Context, id *int) (*model.AdjustmentDetailEntityModel, error) {
	conn := r.CheckTrx(ctx)

	var data model.AdjustmentDetailEntityModel
	if err := conn.Where("id = ?", &id).Preload("Adjustment").First(&data).Error; err != nil {
		return &data, err
	}
	return &data, nil
}
