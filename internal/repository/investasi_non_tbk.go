package repository

import (
	"errors"
	"worker-validation/internal/abstraction"
	"worker-validation/internal/model"

	"gorm.io/gorm"
)

type InvestasiNonTbk interface {
	Find(ctx *abstraction.Context, m *model.InvestasiNonTbkFilterModel) (*[]model.InvestasiNonTbkEntityModel, error)
	GetCount(ctx *abstraction.Context, m *model.InvestasiNonTbkFilterModel) (*int64, error)
	FindByCriteria(ctx *abstraction.Context, filter *model.FilterData) (data *model.InvestasiNonTbkEntityModel, err error)
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

func (r *investasinontbk) FindByCriteria(ctx *abstraction.Context, filter *model.FilterData) (data *model.InvestasiNonTbkEntityModel, err error) {
	conn := r.CheckTrx(ctx)
	query := conn.Model(&model.InvestasiNonTbkEntityModel{})
	if err = query.Where("company_id = ?", filter.CompanyID).Where("period = ?", filter.Period).Where("versions = ?", filter.Versions).First(&data).Error; err != nil {
		return
	}
	if data.ID == 0 {
		err = errors.New("Data Not Found")
	}
	return
}
