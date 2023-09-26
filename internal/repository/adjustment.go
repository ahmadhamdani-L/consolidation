package repository

import (
	"worker-consol/internal/abstraction"
	"worker-consol/internal/model"

	"gorm.io/gorm"
)

type Adjustment interface {
	Find(ctx *abstraction.Context, m *model.AdjustmentFilterModel) (*[]model.AdjustmentEntityModel, error)
	FindByID(ctx *abstraction.Context, id *int) (*model.AdjustmentEntityModel, error)
	Export(ctx *abstraction.Context, e *model.AdjustmentFilterModel) (*model.AdjustmentEntityModel, error)
}

type adjustment struct {
	abstraction.Repository
}

func NewAdjustment(db *gorm.DB) *adjustment {
	return &adjustment{
		abstraction.Repository{
			Db: db,
		},
	}
}

func (r *adjustment) Find(ctx *abstraction.Context, m *model.AdjustmentFilterModel) (*[]model.AdjustmentEntityModel, error) {
	conn := r.CheckTrx(ctx)

	var datas []model.AdjustmentEntityModel
	query := conn.Model(&model.AdjustmentEntityModel{})

	// filter
	query = r.Filter(ctx, query, *m)

	err := query.Preload("AdjustmentDetail").Find(&datas).Error
	if err != nil {
		return &datas, err
	}

	return &datas, nil
}

func (r *adjustment) FindByID(ctx *abstraction.Context, id *int) (*model.AdjustmentEntityModel, error) {
	conn := r.CheckTrx(ctx)

	var data model.AdjustmentEntityModel

	err := conn.Where("id = ?", id).Preload("AdjustmentDetail").First(&data).Error
	if err != nil {
		return nil, err
	}
	return &data, nil
}

func (r *adjustment) Export(ctx *abstraction.Context, e *model.AdjustmentFilterModel) (*model.AdjustmentEntityModel, error) {
	conn := r.CheckTrx(ctx)
	var data model.AdjustmentEntityModel
	query := conn.Model(&model.AdjustmentEntityModel{}).Preload("Company").Preload("AdjustmentDetail").Where("tb_id = ?", &e.TrialBalanceID).Where("company_id = ?", &e.CompanyID).Where("period = ?", &e.Period).Find(&data)
	if err := query.Error; err != nil {
		return nil, err
	}
	return &data, nil
}

func (r *adjustment) Update(ctx *abstraction.Context, id *int, e *model.AdjustmentEntityModel) (*model.AdjustmentEntityModel, error) {
	conn := r.CheckTrx(ctx)

	if err := conn.Model(e).Where("id = ?", &id).Updates(e).Error; err != nil {
		return nil, err
	}
	if err := conn.Model(e).Where("id = ?", &id).First(e).Error; err != nil {
		return nil, err
	}
	return e, nil
}
