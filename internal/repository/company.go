package repository

import (
	"worker/internal/abstraction"
	"worker/internal/model"

	"gorm.io/gorm"
)

type Company interface {
	Find(ctx *abstraction.Context, m *model.CompanyFilterModel) (*[]model.CompanyEntityModel, error)
	FindByID(ctx *abstraction.Context, id *int) (*model.CompanyEntityModel, error)
}

type company struct {
	abstraction.Repository
}

func NewCompany(db *gorm.DB) *company {
	return &company{
		abstraction.Repository{
			Db: db,
		},
	}
}

func (r *company) Find(ctx *abstraction.Context, m *model.CompanyFilterModel) (*[]model.CompanyEntityModel, error) {
	conn := r.CheckTrx(ctx)

	var datas []model.CompanyEntityModel

	query := conn.Model(&model.CompanyEntityModel{})

	//filter
	query = r.Filter(ctx, query, *m)

	if err := query.Find(&datas).Error; err != nil {
		return &datas, err
	}

	return &datas, nil
}

func (r *company) FindByID(ctx *abstraction.Context, id *int) (*model.CompanyEntityModel, error) {
	conn := r.CheckTrx(ctx)

	var data model.CompanyEntityModel

	err := conn.Where("id = ?", id).First(&data).Error
	if err != nil {
		return nil, err
	}
	return &data, nil
}
