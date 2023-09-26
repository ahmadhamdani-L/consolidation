package repository

import (
	"worker-consol/internal/abstraction"
	"worker-consol/internal/model"

	"gorm.io/gorm"
)

type MutasiRua interface {
	Find(ctx *abstraction.Context, m *model.MutasiRuaFilterModel) (*[]model.MutasiRuaEntityModel, error)
	GetCount(ctx *abstraction.Context, m *model.MutasiRuaFilterModel) (*int64, error)
}

type mutasirua struct {
	abstraction.Repository
}

func NewMutasiRua(db *gorm.DB) *mutasirua {
	return &mutasirua{
		abstraction.Repository{
			Db: db,
		},
	}
}

func (r *mutasirua) Find(ctx *abstraction.Context, m *model.MutasiRuaFilterModel) (*[]model.MutasiRuaEntityModel, error) {
	conn := r.CheckTrx(ctx)

	var datas []model.MutasiRuaEntityModel

	query := conn.Model(&model.MutasiRuaEntityModel{})
	//filter
	query = r.Filter(ctx, query, *m)

	if err := query.Find(&datas).Error; err != nil {
		return &datas, err
	}

	return &datas, nil
}

func (r *mutasirua) GetCount(ctx *abstraction.Context, m *model.MutasiRuaFilterModel) (*int64, error) {
	var jmlData int64
	conn := r.CheckTrx(ctx)
	query := conn.Model(&model.MutasiRuaEntityModel{})
	//filter
	query = r.Filter(ctx, query, *m)

	if err := query.Count(&jmlData).Error; err != nil {
		return &jmlData, err
	}

	return &jmlData, nil
}

func (r *mutasirua) Update(ctx *abstraction.Context, id *int, e *model.MutasiRuaEntityModel) (*model.MutasiRuaEntityModel, error) {
	conn := r.CheckTrx(ctx)

	if err := conn.Model(e).Where("id = ?", &id).Updates(e).Error; err != nil {
		return nil, err
	}
	if err := conn.Model(e).Where("id = ?", &id).First(e).Error; err != nil {
		return nil, err
	}

	return e, nil
}
