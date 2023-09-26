package repository

import (
	"worker-consol/internal/abstraction"
	"worker-consol/internal/model"

	"gorm.io/gorm"
)

type MutasiIa interface {
	Find(ctx *abstraction.Context, m *model.MutasiIaFilterModel) (*[]model.MutasiIaEntityModel, error)
	GetCount(ctx *abstraction.Context, m *model.MutasiIaFilterModel) (*int64, error)
	Update(ctx *abstraction.Context, id *int, e *model.MutasiIaEntityModel) (*model.MutasiIaEntityModel, error)
}

type mutasiia struct {
	abstraction.Repository
}

func NewMutasiIa(db *gorm.DB) *mutasiia {
	return &mutasiia{
		abstraction.Repository{
			Db: db,
		},
	}
}

func (r *mutasiia) Find(ctx *abstraction.Context, m *model.MutasiIaFilterModel) (*[]model.MutasiIaEntityModel, error) {
	conn := r.CheckTrx(ctx)

	var datas []model.MutasiIaEntityModel

	query := conn.Model(&model.MutasiIaEntityModel{})
	//filter
	query = r.Filter(ctx, query, *m)

	if err := query.Find(&datas).Error; err != nil {
		return &datas, err
	}

	return &datas, nil
}

func (r *mutasiia) GetCount(ctx *abstraction.Context, m *model.MutasiIaFilterModel) (*int64, error) {
	var jmlData int64
	conn := r.CheckTrx(ctx)
	query := conn.Model(&model.MutasiIaEntityModel{})
	//filter
	query = r.Filter(ctx, query, *m)

	if err := query.Count(&jmlData).Error; err != nil {
		return &jmlData, err
	}

	return &jmlData, nil
}

func (r *mutasiia) Update(ctx *abstraction.Context, id *int, e *model.MutasiIaEntityModel) (*model.MutasiIaEntityModel, error) {
	conn := r.CheckTrx(ctx)

	if err := conn.Model(e).Where("id = ?", &id).Updates(e).Error; err != nil {
		return nil, err
	}
	if err := conn.Model(e).Where("id = ?", &id).First(e).Error; err != nil {
		return nil, err
	}
	return e, nil
}
