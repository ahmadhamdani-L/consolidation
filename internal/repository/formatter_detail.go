package repository

import (
	"worker-validation/internal/abstraction"
	"worker-validation/internal/model"

	"gorm.io/gorm"
)

type FormatterDetail interface {
	Find(ctx *abstraction.Context, m *model.FormatterDetailFilterModel) (*[]model.FormatterDetailEntityModel, error)
}

type formatterdetail struct {
	abstraction.Repository
}

func NewFormatterDetail(db *gorm.DB) *formatterdetail {
	return &formatterdetail{
		abstraction.Repository{
			Db: db,
		},
	}
}

func (r *formatterdetail) Find(ctx *abstraction.Context, m *model.FormatterDetailFilterModel) (*[]model.FormatterDetailEntityModel, error) {
	conn := r.CheckTrx(ctx)

	var datas []model.FormatterDetailEntityModel

	query := conn.Model(&model.FormatterDetailEntityModel{})
	//filter
	query = r.Filter(ctx, query, *m)

	query = query.Order("sort_id asc")

	if err := query.Find(&datas).Error; err != nil {
		return &datas, err
	}

	return &datas, nil
}
