package repository

import (
	"worker-validation/internal/abstraction"
	"worker-validation/internal/model"

	"gorm.io/gorm"
)

type Controller interface {
	Find(ctx *abstraction.Context, m *model.ControllerFilterModel) (*[]model.ControllerEntityModel, error)
	FindByCriteria(ctx *abstraction.Context, m *model.ControllerFilterModel) (*[]model.ControllerEntityModel, error)
}

type controller struct {
	abstraction.Repository
}

func NewController(db *gorm.DB) *controller {
	return &controller{
		abstraction.Repository{
			Db: db,
		},
	}
}

func (r *controller) Find(ctx *abstraction.Context, m *model.ControllerFilterModel) (*[]model.ControllerEntityModel, error) {
	conn := r.CheckTrx(ctx)

	var datas []model.ControllerEntityModel

	query := conn.Model(&model.ControllerEntityModel{})

	//filter
	query = r.Filter(ctx, query, *m)

	if err := query.Find(&datas).Error; err != nil {
		return &datas, err
	}

	return &datas, nil
}

func (r *controller) FindByCriteria(ctx *abstraction.Context, m *model.ControllerFilterModel) (*[]model.ControllerEntityModel, error) {
	conn := r.CheckTrx(ctx)

	var datas []model.ControllerEntityModel

	query := conn.Model(&model.ControllerEntityModel{})

	//filter
	query = r.Filter(ctx, query, *m)

	if err := query.Find(&datas).Error; err != nil {
		return &datas, err
	}

	return &datas, nil
}
