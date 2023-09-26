package repository

import (
	"worker-validation/internal/abstraction"
	"worker-validation/internal/model"

	"gorm.io/gorm"
)

type FormatterBridges interface {
	FindWithCriteria(ctx *abstraction.Context, m *model.FormatterBridgesFilterModel) (*[]model.FormatterBridgesEntityModel, error)
	FindWithCriterias(ctx *abstraction.Context, m *model.FormatterBridgesFilterModel) (*model.FormatterBridgesEntityModel, error)
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

func (r *formatterbridges) FindWithCriteria(ctx *abstraction.Context, m *model.FormatterBridgesFilterModel) (*[]model.FormatterBridgesEntityModel, error) {
	conn := r.CheckTrx(ctx)

	var datas []model.FormatterBridgesEntityModel

	query := conn.Model(&model.FormatterBridgesEntityModel{})
	//filter
	query = r.Filter(ctx, query, *m)

	if err := query.Order("created_at ASC").Find(&datas).Error; err != nil {
		return &datas, err
	}

	return &datas, nil
}
func (r *formatterbridges) FindWithCriterias(ctx *abstraction.Context, m *model.FormatterBridgesFilterModel) (*model.FormatterBridgesEntityModel, error) {
	conn := r.CheckTrx(ctx)

	var datas model.FormatterBridgesEntityModel

	query := conn.Model(&model.FormatterBridgesEntityModel{})
	//filter
	query = r.Filter(ctx, query, *m)

	if err := query.Order("created_at ASC").First(&datas).Error; err != nil {
		return &datas, err
	}

	return &datas, nil
}

