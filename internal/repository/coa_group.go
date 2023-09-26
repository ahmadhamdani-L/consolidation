package repository

import (
	"worker-consol/internal/abstraction"
	"worker-consol/internal/model"

	"gorm.io/gorm"
)

type CoaGroup interface {
	Find(ctx *abstraction.Context, m *model.CoaGroupFilterModel) (*[]model.CoaGroupEntityModel, error)
	Create(ctx *abstraction.Context, e *model.CoaGroupEntityModel) (*model.CoaGroupEntityModel, error)
}

type coagroup struct {
	abstraction.Repository
}

func NewCoaGroup(db *gorm.DB) *coagroup {
	return &coagroup{
		abstraction.Repository{
			Db: db,
		},
	}
}

func (r *coagroup) Find(ctx *abstraction.Context, m *model.CoaGroupFilterModel) (*[]model.CoaGroupEntityModel, error) {
	conn := r.CheckTrx(ctx)

	var datas []model.CoaGroupEntityModel

	query := conn.Model(&model.CoaGroupEntityModel{})

	//filter
	query = r.Filter(ctx, query, *m)

	if err := query.Find(&datas).Error; err != nil {
		return &datas, err
	}

	return &datas, nil
}

func (r *coagroup) Create(ctx *abstraction.Context, e *model.CoaGroupEntityModel) (*model.CoaGroupEntityModel, error) {
	conn := r.CheckTrx(ctx)

	if err := conn.Create(e).Error; err != nil {
		return nil, err
	}
	if err := conn.Model(e).First(e).Error; err != nil {
		return nil, err
	}

	return e, nil
}
