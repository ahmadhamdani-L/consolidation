package repository

import (
	"fmt"
	"worker/internal/abstraction"
	"worker/internal/model"

	"gorm.io/gorm"
)

type ConsolidationDetail interface {
	Find(ctx *abstraction.Context, m *model.ConsolidationDetailFilterModel) (*[]model.ConsolidationDetailEntityModel, error)
	FindWithCode(ctx *abstraction.Context, code *string) (*[]model.ConsolidationDetailEntityModel, error)
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
