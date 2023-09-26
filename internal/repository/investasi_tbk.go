package repository

import (
	"errors"
	"worker-validation/internal/abstraction"
	"worker-validation/internal/model"

	"gorm.io/gorm"
)

type InvestasiTbk interface {
	Find(ctx *abstraction.Context, m *model.InvestasiTbkFilterModel) (*[]model.InvestasiTbkEntityModel, error)
	GetCount(ctx *abstraction.Context, m *model.InvestasiTbkFilterModel) (*int64, error)
	FindByCriteria(ctx *abstraction.Context, filter *model.FilterData) (data *model.InvestasiTbkEntityModel, err error)
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

func (r *investasitbk) FindByCriteria(ctx *abstraction.Context, filter *model.FilterData) (data *model.InvestasiTbkEntityModel, err error) {
	conn := r.CheckTrx(ctx)
	query := conn.Model(&model.InvestasiTbkEntityModel{})
	if err = query.Where("company_id = ?", filter.CompanyID).Where("period = ?", filter.Period).Where("versions = ?", filter.Versions).First(&data).Error; err != nil {
		return
	}
	if data.ID == 0 {
		err = errors.New("Data Not Found")
	}
	return
}
