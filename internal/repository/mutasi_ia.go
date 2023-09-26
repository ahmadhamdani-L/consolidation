package repository

import (
	"errors"
	"worker-validation/internal/abstraction"
	"worker-validation/internal/model"

	"gorm.io/gorm"
)

type MutasiIa interface {
	Find(ctx *abstraction.Context, m *model.MutasiIaFilterModel) (*[]model.MutasiIaEntityModel, error)
	GetCount(ctx *abstraction.Context, m *model.MutasiIaFilterModel) (*int64, error)
	FindByCriteria(ctx *abstraction.Context, filter *model.FilterData) (data *model.MutasiIaEntityModel, err error)
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

func (r *mutasiia) FindByCriteria(ctx *abstraction.Context, filter *model.FilterData) (data *model.MutasiIaEntityModel, err error) {
	conn := r.CheckTrx(ctx)
	query := conn.Model(&model.MutasiIaEntityModel{})
	if err = query.Where("company_id = ?", filter.CompanyID).Where("period = ?", filter.Period).Where("versions = ?", filter.Versions).First(&data).Error; err != nil {
		return
	}
	if data.ID == 0 {
		err = errors.New("Data Not Found")
	}
	return
}

func (r *mutasiia) Update(ctx *abstraction.Context, id *int, e *model.MutasiIaEntityModel) (*model.MutasiIaEntityModel, error) {
	conn := r.CheckTrx(ctx)

	if err := conn.Model(e).Where("id = ?", id).Updates(e).Error; err != nil {
		return nil, err
	}
	if err := conn.Model(e).Where("id = ?", id).First(e).Error; err != nil {
		return nil, err
	}
	return e, nil
}
