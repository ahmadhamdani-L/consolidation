package repository

import (
	"gorm.io/gorm"

	// "gorm.io/gorm/clause"

	"worker-consol/internal/abstraction"
	"worker-consol/internal/model"
)

type Jpm interface {
	Find(ctx *abstraction.Context, m *model.JpmFilterModel) (*[]model.JpmEntityModel, error)
	FindByID(ctx *abstraction.Context, id *int) (*model.JpmEntityModel, error)
	Export(ctx *abstraction.Context, e *model.JpmFilterModel) (*model.JpmEntityModel, error)
	Update(ctx *abstraction.Context, id *int, e *model.JpmEntityModel) (*model.JpmEntityModel, error)
}

type jpm struct {
	abstraction.Repository
}

func NewJpm(db *gorm.DB) *jpm {
	return &jpm{
		abstraction.Repository{
			Db: db,
		},
	}
}

func (r *jpm) Find(ctx *abstraction.Context, m *model.JpmFilterModel) (*[]model.JpmEntityModel, error) {
	conn := r.CheckTrx(ctx)

	var datas []model.JpmEntityModel

	query := conn.Model(&model.JpmEntityModel{})

	// filter
	query = r.Filter(ctx, query, *m)

	err := query.Preload("JpmDetail").Find(&datas).Error
	if err != nil {
		return &datas, err
	}

	return &datas, nil
}

func (r *jpm) FindByID(ctx *abstraction.Context, id *int) (*model.JpmEntityModel, error) {
	conn := r.CheckTrx(ctx)

	var data model.JpmEntityModel

	err := conn.Where("id = ?", id).First(&data).Error
	if err != nil {
		return nil, err
	}
	return &data, nil
}

func (r *jpm) Export(ctx *abstraction.Context, e *model.JpmFilterModel) (*model.JpmEntityModel, error) {
	conn := r.CheckTrx(ctx)
	var data model.JpmEntityModel
	query := conn.Model(&model.JpmEntityModel{}).Where("period = ?", e.Period).Where("company_id = ?", e.CompanyID).Preload("Company").Preload("JpmDetail").Find(&data)
	if err := query.Error; err != nil {
		return nil, err
	}
	return &data, nil
}

func (r *jpm) Update(ctx *abstraction.Context, id *int, e *model.JpmEntityModel) (*model.JpmEntityModel, error) {
	conn := r.CheckTrx(ctx)

	if err := conn.Model(e).Where("id = ?", &id).Updates(e).Error; err != nil {
		return nil, err
	}
	if err := conn.Model(e).Where("id = ?", &id).First(e).Error; err != nil {
		return nil, err
	}
	return e, nil

}
