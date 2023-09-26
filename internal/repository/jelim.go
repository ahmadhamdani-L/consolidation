package repository

import (
	"gorm.io/gorm"
	// "gorm.io/gorm/clause"
	"worker-validation/internal/abstraction"
	"worker-validation/internal/model"
)

type Jelim interface {
	Find(ctx *abstraction.Context, m *model.JelimFilterModel) (*[]model.JelimEntityModel, error)
	FindByID(ctx *abstraction.Context, id *int) (*model.JelimEntityModel, error)
	Export(ctx *abstraction.Context, e *model.JelimFilterModel) (*model.JelimEntityModel, error)
}

type jelim struct {
	abstraction.Repository
}

func NewJelim(db *gorm.DB) *jelim {
	return &jelim{
		abstraction.Repository{
			Db: db,
		},
	}
}

func (r *jelim) Find(ctx *abstraction.Context, m *model.JelimFilterModel) (*[]model.JelimEntityModel, error) {
	conn := r.CheckTrx(ctx)

	var datas []model.JelimEntityModel
	query := conn.Model(&model.JelimEntityModel{})

	// filter
	query = r.Filter(ctx, query, *m)

	err := query.Preload("JelimDetail").Find(&datas).Error
	if err != nil {
		return &datas, err
	}

	return &datas, nil
}

func (r *jelim) FindByID(ctx *abstraction.Context, id *int) (*model.JelimEntityModel, error) {
	conn := r.CheckTrx(ctx)

	var data model.JelimEntityModel
	err := conn.Where("id = ?", id).First(&data).Error
	if err != nil {
		return nil, err
	}
	return &data, nil
}

func (r *jelim) Export(ctx *abstraction.Context, e *model.JelimFilterModel) (*model.JelimEntityModel, error) {
	conn := r.CheckTrx(ctx)
	var data model.JelimEntityModel
	query := conn.Model(&model.JelimEntityModel{}).Preload("Company").Preload("JelimDetail").Where("tb_id = ?", &e.TbID).Where("period = ?", &e.Period).Where("company_id = ?", &e.CompanyID).Find(&data)
	if err := query.Error; err != nil {
		return nil, err
	}
	return &data, nil
}
