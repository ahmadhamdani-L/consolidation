package repository

import (
	"strings"
	"worker/internal/abstraction"
	"worker/internal/model"

	"gorm.io/gorm"
)

type Adjustment interface {
	Find(ctx *abstraction.Context, m *model.AdjustmentFilterModel) (*[]model.AdjustmentEntityModel, error)
	FindByID(ctx *abstraction.Context, id *int) (*model.AdjustmentEntityModel, error)
	Export(ctx *abstraction.Context, e *model.AdjustmentFilterModel) (*model.AdjustmentEntityModel, error)
	FindSummary(ctx *abstraction.Context, tbID *int) (*model.AdjustmentDetailEntityModel, error)
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

func (r *adjustment) FindSummary(ctx *abstraction.Context, tbID *int) (*model.AdjustmentDetailEntityModel, error) {
	conn := r.CheckTrx(ctx)
	var data []string
	if err := conn.Model(&model.AdjustmentEntityModel{}).Where("tb_id = ?", &tbID).Pluck("id", &data).Error; err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	}

	if len(data) == 0 {
		return &model.AdjustmentDetailEntityModel{}, nil
	}

	listID := strings.Join(data, ",")
	var sumData model.AdjustmentDetailEntityModel
	if err := conn.Model(&model.AdjustmentDetailEntityModel{}).Where("adjustment_id IN (?)", listID).Select("SUM(balance_sheet_cr) balance_sheet_cr, SUM(balance_sheet_dr) balance_sheet_dr, SUM(income_statement_cr) income_statement_cr, SUM(income_statement_dr) income_statement_dr").Find(&sumData).Error; err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	}

	return &sumData, nil
}
