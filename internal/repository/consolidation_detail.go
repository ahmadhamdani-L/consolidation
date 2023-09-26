package repository

import (
	"fmt"
	"worker-consol/internal/abstraction"
	"worker-consol/internal/model"

	"gorm.io/gorm"
)

type ConsolidationDetail interface {
	Find(ctx *abstraction.Context, m *model.ConsolidationDetailFilterModel) (*[]model.ConsolidationDetailEntityModel, error)
	FindWithCode(ctx *abstraction.Context, code *string) (*[]model.ConsolidationDetailEntityModel, error)
	Create(ctx *abstraction.Context, payload *model.ConsolidationDetailEntityModel) (*model.ConsolidationDetailEntityModel, error)
	DeleteByConsolID(ctx *abstraction.Context, consolID *int) error
}

type consolidationdetail struct {
	abstraction.Repository
}

func NewConsolidationDetail(db *gorm.DB) *consolidationdetail {
	return &consolidationdetail{
		abstraction.Repository{
			Db: db,
		},
	}
}

func (r *consolidationdetail) Find(ctx *abstraction.Context, m *model.ConsolidationDetailFilterModel) (*[]model.ConsolidationDetailEntityModel, error) {
	conn := r.CheckTrx(ctx)

	var datas []model.ConsolidationDetailEntityModel

	query := conn.Model(&model.ConsolidationDetailEntityModel{})
	//filter
	query = r.Filter(ctx, query, *m)
	if err := query.Find(&datas).Error; err != nil {
		return &datas, err
	}

	return &datas, nil
}

func (r *consolidationdetail) FindWithCode(ctx *abstraction.Context, code *string) (*[]model.ConsolidationDetailEntityModel, error) {
	conn := r.CheckTrx(ctx)

	var data []model.ConsolidationDetailEntityModel
	tmp := fmt.Sprintf("%s", *code)
	if err := conn.Where("code LIKE ?", tmp+"%").Find(&data).Error; err != nil {
		return &data, err
	}
	return &data, nil
}

func (r *consolidationdetail) Create(ctx *abstraction.Context, payload *model.ConsolidationDetailEntityModel) (*model.ConsolidationDetailEntityModel, error) {
	conn := r.CheckTrx(ctx)
	if err := conn.Create(&payload).Error; err != nil {
		return nil, err
	}

	if err := conn.Model(&payload).First(&payload).Error; err != nil {
		return nil, err
	}

	return payload, nil
}

func (r *consolidationdetail) DeleteByConsolID(ctx *abstraction.Context, consolID *int) error {
	conn := r.CheckTrx(ctx)

	if err := conn.Model(&model.ConsolidationDetailEntityModel{}).Where("consolidation_id = ?", consolID).Delete(&model.ConsolidationDetailEntityModel{}).Error; err != nil {
		return err
	}

	return nil
}
