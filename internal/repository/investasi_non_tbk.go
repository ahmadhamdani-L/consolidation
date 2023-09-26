package repository

import (
	"worker-consol/internal/abstraction"
	"worker-consol/internal/model"

	"gorm.io/gorm"
)

type InvestasiNonTbk interface {
	Find(ctx *abstraction.Context, m *model.InvestasiNonTbkFilterModel) (*[]model.InvestasiNonTbkEntityModel, error)
	GetCount(ctx *abstraction.Context, m *model.InvestasiNonTbkFilterModel) (*int64, error)
	Update(ctx *abstraction.Context, id *int, e *model.InvestasiNonTbkEntityModel) (*model.InvestasiNonTbkEntityModel, error)
}

type investasinontbk struct {
	abstraction.Repository
}

func NewInvestasiNonTbk(db *gorm.DB) *investasinontbk {
	return &investasinontbk{
		abstraction.Repository{
			Db: db,
		},
	}
}

func (r *investasinontbk) Find(ctx *abstraction.Context, m *model.InvestasiNonTbkFilterModel) (*[]model.InvestasiNonTbkEntityModel, error) {
	conn := r.CheckTrx(ctx)

	var datas []model.InvestasiNonTbkEntityModel

	query := conn.Model(&model.InvestasiNonTbkEntityModel{})
	//filter
	query = r.Filter(ctx, query, *m)

	if err := query.Find(&datas).Error; err != nil {
		return &datas, err
	}

	return &datas, nil
}

func (r *investasinontbk) GetCount(ctx *abstraction.Context, m *model.InvestasiNonTbkFilterModel) (*int64, error) {
	var jmlData int64
	conn := r.CheckTrx(ctx)
	query := conn.Model(&model.InvestasiNonTbkEntityModel{})
	//filter
	query = r.Filter(ctx, query, *m)

	if err := query.Count(&jmlData).Error; err != nil {
		return &jmlData, err
	}

	return &jmlData, nil
}

func (r *investasinontbk) Update(ctx *abstraction.Context, id *int, e *model.InvestasiNonTbkEntityModel) (*model.InvestasiNonTbkEntityModel, error) {
	conn := r.CheckTrx(ctx)

	if err := conn.Model(e).Where("id = ?", &id).Updates(e).Error; err != nil {
		return nil, err
	}
	if err := conn.Model(e).Where("id = ?", &id).First(e).Error; err != nil {
		return nil, err
	}
	return e, nil
}
