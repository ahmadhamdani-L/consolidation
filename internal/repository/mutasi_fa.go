package repository

import (
	"worker/internal/abstraction"
	"worker/internal/model"

	"gorm.io/gorm"
)

type MutasiFa interface {
	Find(ctx *abstraction.Context, m *model.MutasiFaFilterModel) (*[]model.MutasiFaEntityModel, error)
	GetCount(ctx *abstraction.Context, m *model.MutasiFaFilterModel) (*int64, error)
	FindByCriteria(ctx *abstraction.Context, m *model.MutasiFaFilterModel) (*model.MutasiFaEntityModel, error)
}

type mutasifa struct {
	abstraction.Repository
}

func NewMutasiFa(db *gorm.DB) *mutasifa {
	return &mutasifa{
		abstraction.Repository{
			Db: db,
		},
	}
}
func (r *mutasifa) FindByCriteria(ctx *abstraction.Context, m *model.MutasiFaFilterModel) (*model.MutasiFaEntityModel, error) {
	conn := r.CheckTrx(ctx)

	var data model.MutasiFaEntityModel
	query := conn.Model(&model.MutasiFaEntityModel{})
	query = r.Filter(ctx, query, *m)

	if err := query.Preload("Company").First(&data).Error; err != nil {
		return &data, err
	}

	return &data, nil
}
func (r *mutasifa) Find(ctx *abstraction.Context, m *model.MutasiFaFilterModel) (*[]model.MutasiFaEntityModel, error) {
	conn := r.CheckTrx(ctx)

	var datas []model.MutasiFaEntityModel

	query := conn.Model(&model.MutasiFaEntityModel{})
	//filter
	query = r.Filter(ctx, query, *m)

	if err := query.Find(&datas).Error; err != nil {
		return &datas, err
	}

	return &datas, nil
}

func (r *mutasifa) GetCount(ctx *abstraction.Context, m *model.MutasiFaFilterModel) (*int64, error) {
	var jmlData int64
	conn := r.CheckTrx(ctx)
	query := conn.Model(&model.MutasiFaEntityModel{})
	//filter
	query = r.Filter(ctx, query, *m)

	if err := query.Count(&jmlData).Error; err != nil {
		return &jmlData, err
	}

	return &jmlData, nil
}
