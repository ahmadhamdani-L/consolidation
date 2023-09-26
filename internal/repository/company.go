package repository

import (
	"worker-validation/internal/abstraction"
	"worker-validation/internal/model"

	"gorm.io/gorm"
)

type Company interface {
	Find(ctx *abstraction.Context, m *model.CompanyFilterModel) (*[]model.CompanyEntityModel, error)
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
