package repository

import (
	"fmt"
	"worker-consol/internal/abstraction"
	"worker-consol/internal/model"

	"gorm.io/gorm"
)

type ConsolidationBridgeDetail interface {
	Find(ctx *abstraction.Context, m *model.ConsolidationBridgeDetailFilterModel) (*[]model.ConsolidationBridgeDetailEntityModel, error)
	FindWithCode(ctx *abstraction.Context, code *string) (*[]model.ConsolidationBridgeDetailEntityModel, error)
	Create(ctx *abstraction.Context, payload *model.ConsolidationBridgeDetailEntityModel) (*model.ConsolidationBridgeDetailEntityModel, error)
	FindSummary(ctx *abstraction.Context, list *string, code *string) (*model.ConsolidationBridgeDetailEntityModel, error)
	DeleteByListBridgeID(ctx *abstraction.Context, list *string) error
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

func (r *consolidationbridgedetail) FindWithCode(ctx *abstraction.Context, code *string) (*[]model.ConsolidationBridgeDetailEntityModel, error) {
	conn := r.CheckTrx(ctx)

	var data []model.ConsolidationBridgeDetailEntityModel
	tmp := fmt.Sprintf("%s", *code)
	if err := conn.Where("code LIKE ?", tmp+"%").Find(&data).Error; err != nil {
		return &data, err
	}
	return &data, nil
}

func (r *consolidationbridgedetail) Create(ctx *abstraction.Context, payload *model.ConsolidationBridgeDetailEntityModel) (*model.ConsolidationBridgeDetailEntityModel, error) {
	conn := r.CheckTrx(ctx)
	if err := conn.Create(&payload).Error; err != nil {
		return nil, err
	}

	if err := conn.Model(&payload).First(&payload).Error; err != nil {
		return nil, err
	}

	return payload, nil
}

func (r *consolidationbridgedetail) FindSummary(ctx *abstraction.Context, list *string, code *string) (*model.ConsolidationBridgeDetailEntityModel, error) {
	conn := r.CheckTrx(ctx)

	var detailData model.ConsolidationBridgeDetailEntityModel
	querySum := conn.Model(&model.ConsolidationBridgeDetailEntityModel{}).Where("code LIKE ?", *code+"%").Where(fmt.Sprintf("consolidation_bridge_id IN (%s)", *list))
	querySum = querySum.Select("SUM(amount) amount")
	if err := querySum.Find(&detailData).Error; err != nil {
		return nil, err
	}

	return &detailData, nil
}

func (r *consolidationbridgedetail) DeleteByListBridgeID(ctx *abstraction.Context, list *string) error {
	conn := r.CheckTrx(ctx)

	if err := conn.Model(&model.ConsolidationBridgeDetailEntityModel{}).Where(fmt.Sprintf("consolidation_bridge_id IN (%s)", *list)).Delete(&model.ConsolidationBridgeDetailEntityModel{}).Error; err != nil {
		return err
	}

	return nil
}
