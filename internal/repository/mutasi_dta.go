package repository

import (
	"errors"
	"worker-validation/internal/abstraction"
	"worker-validation/internal/model"

	"gorm.io/gorm"
)

type MutasiDta interface {
	Find(ctx *abstraction.Context, m *model.MutasiDtaFilterModel) (*[]model.MutasiDtaEntityModel, error)
	GetCount(ctx *abstraction.Context, m *model.MutasiDtaFilterModel) (*int64, error)
	FindByCriteria(ctx *abstraction.Context, filter *model.FilterData) (data *model.MutasiDtaEntityModel, err error)
	Update(ctx *abstraction.Context, id *int, e *model.MutasiDtaEntityModel) (*model.MutasiDtaEntityModel, error)
}

type mutasidta struct {
	abstraction.Repository
}

func NewMutasiDta(db *gorm.DB) *mutasidta {
	return &mutasidta{
		abstraction.Repository{
			Db: db,
		},
	}
}

func (r *mutasidta) Find(ctx *abstraction.Context, m *model.MutasiDtaFilterModel) (*[]model.MutasiDtaEntityModel, error) {
	conn := r.CheckTrx(ctx)

	var datas []model.MutasiDtaEntityModel

	query := conn.Model(&model.MutasiDtaEntityModel{})
	//filter
	query = r.Filter(ctx, query, *m)

	if err := query.Find(&datas).Error; err != nil {
		return &datas, err
	}

	return &datas, nil
}

func (r *mutasidta) GetCount(ctx *abstraction.Context, m *model.MutasiDtaFilterModel) (*int64, error) {
	var jmlData int64
	conn := r.CheckTrx(ctx)
	query := conn.Model(&model.MutasiDtaEntityModel{})
	//filter
	query = r.Filter(ctx, query, *m)

	if err := query.Count(&jmlData).Error; err != nil {
		return &jmlData, err
	}

	return &jmlData, nil
}

func (r *mutasidta) FindByCriteria(ctx *abstraction.Context, filter *model.FilterData) (data *model.MutasiDtaEntityModel, err error) {
	conn := r.CheckTrx(ctx)
	query := conn.Model(&model.MutasiDtaEntityModel{})
	if err = query.Where("company_id = ?", filter.CompanyID).Where("period = ?", filter.Period).Where("versions = ?", filter.Versions).First(&data).Error; err != nil {
		return
	}
	if data.ID == 0 {
		err = errors.New("Data Not Found")
	}
	return
}

func (r *mutasidta) Update(ctx *abstraction.Context, id *int, e *model.MutasiDtaEntityModel) (*model.MutasiDtaEntityModel, error) {
	conn := r.CheckTrx(ctx)

	if err := conn.Model(e).Where("id = ?", id).Updates(e).Error; err != nil {
		return nil, err
	}
	if err := conn.Model(e).Where("id = ?", id).First(e).Error; err != nil {
		return nil, err
	}
	return e, nil
}
