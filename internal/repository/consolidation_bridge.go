package repository

import (
	"errors"
	"worker/internal/abstraction"
	"worker/internal/model"

	"gorm.io/gorm"
)

type ConsolidationBridge interface {
	Find(ctx *abstraction.Context, m *model.ConsolidationBridgeFilterModel) (*[]model.ConsolidationBridgeEntityModel, error)
	FindByCriteria(ctx *abstraction.Context, filter *model.ConsolidationBridgeFilterModel) (data *model.ConsolidationBridgeEntityModel, err error)
}

type consolidationbridge struct {
	abstraction.Repository
}

func NewConsolidationBridge(db *gorm.DB) *consolidationbridge {
	return &consolidationbridge{
		abstraction.Repository{
			Db: db,
		},
	}
}

func (r *consolidationbridge) Find(ctx *abstraction.Context, m *model.ConsolidationBridgeFilterModel) (*[]model.ConsolidationBridgeEntityModel, error) {
	conn := r.CheckTrx(ctx)

	var datas []model.ConsolidationBridgeEntityModel

	query := conn.Model(&model.ConsolidationBridgeEntityModel{})
	//filter
	query = r.Filter(ctx, query, *m)

	if err := query.Find(&datas).Error; err != nil {
		return &datas, err
	}

	return &datas, nil
}

func (r *consolidationbridge) FindByCriteria(ctx *abstraction.Context, filter *model.ConsolidationBridgeFilterModel) (data *model.ConsolidationBridgeEntityModel, err error) {
	conn := r.CheckTrx(ctx)
	query := conn.Model(&model.ConsolidationBridgeEntityModel{})
	query = r.Filter(ctx, query, *filter)
	if err = query.First(&data).Error; err != nil {
		return
	}
	if data.ID == 0 {
		err = errors.New("Data Not Found")
	}
	return
}
