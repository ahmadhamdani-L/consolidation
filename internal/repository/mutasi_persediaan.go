package repository

import (
	"worker/internal/abstraction"
	"worker/internal/model"

	"gorm.io/gorm"
)

type MutasiPersediaan interface {
	Find(ctx *abstraction.Context, m *model.MutasiPersediaanFilterModel) (*[]model.MutasiPersediaanEntityModel, error)
	GetCount(ctx *abstraction.Context, m *model.MutasiPersediaanFilterModel) (*int64, error)
}

type mutasipersediaan struct {
	abstraction.Repository
}

func NewMutasiPersediaan(db *gorm.DB) *mutasipersediaan {
	return &mutasipersediaan{
		abstraction.Repository{
			Db: db,
		},
	}
}

func (r *mutasipersediaan) Find(ctx *abstraction.Context, m *model.MutasiPersediaanFilterModel) (*[]model.MutasiPersediaanEntityModel, error) {
	conn := r.CheckTrx(ctx)

	var datas []model.MutasiPersediaanEntityModel

	query := conn.Model(&model.MutasiPersediaanEntityModel{})
	//filter
	query = r.Filter(ctx, query, *m)

	if err := query.Find(&datas).Error; err != nil {
		return &datas, err
	}

	return &datas, nil
}

func (r *mutasipersediaan) GetCount(ctx *abstraction.Context, m *model.MutasiPersediaanFilterModel) (*int64, error) {
	var jmlData int64
	conn := r.CheckTrx(ctx)
	query := conn.Model(&model.MutasiPersediaanEntityModel{})
	//filter
	query = r.Filter(ctx, query, *m)

	if err := query.Count(&jmlData).Error; err != nil {
		return &jmlData, err
	}

	return &jmlData, nil
}
