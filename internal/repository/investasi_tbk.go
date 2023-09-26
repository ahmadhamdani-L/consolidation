package repository

import (
	"worker-consol/internal/abstraction"
	"worker-consol/internal/model"

	"gorm.io/gorm"
)

type InvestasiTbk interface {
	Find(ctx *abstraction.Context, m *model.InvestasiTbkFilterModel) (*[]model.InvestasiTbkEntityModel, error)
	GetCount(ctx *abstraction.Context, m *model.InvestasiTbkFilterModel) (*int64, error)
	Update(ctx *abstraction.Context, id *int, e *model.InvestasiTbkEntityModel) (*model.InvestasiTbkEntityModel, error)
}

type investasitbk struct {
	abstraction.Repository
}

func NewInvestasiTbk(db *gorm.DB) *investasitbk {
	return &investasitbk{
		abstraction.Repository{
			Db: db,
		},
	}
}

func (r *investasitbk) Find(ctx *abstraction.Context, m *model.InvestasiTbkFilterModel) (*[]model.InvestasiTbkEntityModel, error) {
	conn := r.CheckTrx(ctx)

	var datas []model.InvestasiTbkEntityModel

	query := conn.Model(&model.InvestasiTbkEntityModel{})
	//filter
	query = r.Filter(ctx, query, *m)

	if err := query.Find(&datas).Error; err != nil {
		return &datas, err
	}

	return &datas, nil
}

func (r *investasitbk) GetCount(ctx *abstraction.Context, m *model.InvestasiTbkFilterModel) (*int64, error) {
	var jmlData int64
	conn := r.CheckTrx(ctx)
	query := conn.Model(&model.InvestasiTbkEntityModel{})
	//filter
	query = r.Filter(ctx, query, *m)

	if err := query.Count(&jmlData).Error; err != nil {
		return &jmlData, err
	}

	return &jmlData, nil
}

func (r *investasitbk) Update(ctx *abstraction.Context, id *int, e *model.InvestasiTbkEntityModel) (*model.InvestasiTbkEntityModel, error) {
	conn := r.CheckTrx(ctx)

	if err := conn.Model(e).Where("id = ?", &id).Updates(e).Error; err != nil {
		return nil, err
	}
	if err := conn.Model(e).Where("id = ?", &id).First(e).Error; err != nil {
		return nil, err
	}
	return e, nil
}
