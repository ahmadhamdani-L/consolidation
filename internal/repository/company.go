package repository

import (
	"notification/internal/abstraction"
	"notification/internal/model"

	"gorm.io/gorm"
)

type Company interface {
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

func (r *company) FindByID(ctx *abstraction.Context, id *int) (*model.CompanyEntityModel, error) {
	conn := r.CheckTrx(ctx)

	var data model.CompanyEntityModel

	err := conn.Where("id = ?", id).First(&data).Error
	if err != nil {
		return nil, err
	}
	return &data, nil
}
