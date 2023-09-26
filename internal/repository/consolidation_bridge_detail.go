package repository

import (
	"worker/internal/abstraction"
	"worker/internal/model"

	"gorm.io/gorm"
)

type ConsolidationBridgeDetail interface {
	Find(ctx *abstraction.Context, m *model.ConsolidationBridgeDetailFilterModel) (*[]model.ConsolidationBridgeDetailEntityModel, error)
	FindWithCode(ctx *abstraction.Context, consolBridgeID *int, code *string) (*[]model.ConsolidationBridgeDetailEntityModel, error)
	GetWithCode(ctx *abstraction.Context, consolBridgeID *int, code *string) (*model.ConsolidationBridgeDetailEntityModel, error)
}

type consolidationbridgedetail struct {
	abstraction.Repository
}

func NewConsolidationBridgeDetail(db *gorm.DB) *consolidationbridgedetail {
	return &consolidationbridgedetail{
		abstraction.Repository{
			Db: db,
		},
	}
}

func (r *consolidationbridgedetail) Find(ctx *abstraction.Context, m *model.ConsolidationBridgeDetailFilterModel) (*[]model.ConsolidationBridgeDetailEntityModel, error) {
	conn := r.CheckTrx(ctx)

	var datas []model.ConsolidationBridgeDetailEntityModel

	query := conn.Model(&model.ConsolidationBridgeDetailEntityModel{})
	//filter
	query = r.Filter(ctx, query, *m)
	if err := query.Find(&datas).Error; err != nil {
		return &datas, err
	}

	return &datas, nil
}

func (r *consolidationbridgedetail) FindWithCode(ctx *abstraction.Context, consolBridgeID *int, code *string) (*[]model.ConsolidationBridgeDetailEntityModel, error) {
	conn := r.CheckTrx(ctx)

	var data []model.ConsolidationBridgeDetailEntityModel
	if err := conn.Where("consolidation_bridge_id = ?", consolBridgeID).Where("code LIKE ?", *code+"%").Find(&data).Error; err != nil {
		return &data, err
	}
	return &data, nil
}

func (r *consolidationbridgedetail) GetWithCode(ctx *abstraction.Context, consolBridgeID *int, code *string) (*model.ConsolidationBridgeDetailEntityModel, error) {
	conn := r.CheckTrx(ctx)

	var data model.ConsolidationBridgeDetailEntityModel
	if err := conn.Where("consolidation_bridge_id = ?", consolBridgeID).Where("code = ?", *code).First(&data).Error; err != nil {
		return &data, err
	}
	return &data, nil
}
