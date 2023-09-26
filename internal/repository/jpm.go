package repository

import (
	"fmt"
	"strings"

	"gorm.io/gorm"

	// "gorm.io/gorm/clause"

	"worker/internal/abstraction"
	"worker/internal/model"
)

type Jpm interface {
	Find(ctx *abstraction.Context, m *model.JpmFilterModel) (*[]model.JpmEntityModel, error)
	FindByID(ctx *abstraction.Context, id *int) (*model.JpmEntityModel, error)
	ExportOne(ctx *abstraction.Context, e *model.JpmFilterModel) (*model.JpmEntityModel, error)
	ExportAll(ctx *abstraction.Context, consolidationID *int) (*[]model.JpmDetailEntityModel, error)
	FindSummary(ctx *abstraction.Context, consolID *int) (*model.JpmDetailEntityModel, error)
}

type jpm struct {
	abstraction.Repository
}

func NewJpm(db *gorm.DB) *jpm {
	return &jpm{
		abstraction.Repository{
			Db: db,
		},
	}
}

func (r *jpm) Find(ctx *abstraction.Context, m *model.JpmFilterModel) (*[]model.JpmEntityModel, error) {
	conn := r.CheckTrx(ctx)

	var datas []model.JpmEntityModel

	query := conn.Model(&model.JpmEntityModel{})

	// filter
	query = r.Filter(ctx, query, *m)

	err := query.Preload("JpmDetail").Find(&datas).Error
	if err != nil {
		return &datas, err
	}

	return &datas, nil
}

func (r *jpm) FindByID(ctx *abstraction.Context, id *int) (*model.JpmEntityModel, error) {
	conn := r.CheckTrx(ctx)

	var data model.JpmEntityModel

	err := conn.Where("id = ?", id).First(&data).Error
	if err != nil {
		return nil, err
	}
	return &data, nil
}

func (r *jpm) ExportOne(ctx *abstraction.Context, e *model.JpmFilterModel) (*model.JpmEntityModel, error) {
	conn := r.CheckTrx(ctx)
	var data model.JpmEntityModel
	query := conn.Model(&model.JpmEntityModel{}).Where("period = ?", e.Period).Where("company_id = ?", e.CompanyID).Preload("Company").Where("consolidation_id = ?", e.ConsolidationID).Preload("JpmDetail").Find(&data)
	if err := query.Error; err != nil {
		return nil, err
	}
	return &data, nil
}

func (r *jpm) FindSummary(ctx *abstraction.Context, consolID *int) (*model.JpmDetailEntityModel, error) {
	conn := r.CheckTrx(ctx)
	var data []string
	if err := conn.Model(&model.JpmEntityModel{}).Where("status != 0").Where("consolidation_id = ?", &consolID).Pluck("id", &data).Error; err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	}

	if len(data) == 0 {
		return &model.JpmDetailEntityModel{}, nil
	}

	listID := strings.Join(data, ",")
	var sumData model.JpmDetailEntityModel
	if err := conn.Model(&model.JpmDetailEntityModel{}).Where(fmt.Sprintf("jpm_id IN (%s)", listID)).Select("SUM(balance_sheet_cr) balance_sheet_cr, SUM(balance_sheet_dr) balance_sheet_dr, SUM(income_statement_cr) income_statement_cr, SUM(income_statement_dr) income_statement_dr").Find(&sumData).Error; err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	}

	return &sumData, nil
}

func (r *jpm) ExportAll(ctx *abstraction.Context, consolidationID *int) (*[]model.JpmDetailEntityModel, error) {
	conn := r.CheckTrx(ctx)
	var data []string
	if err := conn.Model(&model.JpmEntityModel{}).Where("status != 4").Where("consolidation_id = ?", &consolidationID).Pluck("id", &data).Error; err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	}

	if len(data) == 0 {
		return &[]model.JpmDetailEntityModel{}, nil
	}

	listID := strings.Join(data, ",")
	var result []model.JpmDetailEntityModel
	if err := conn.Model(&model.JpmDetailEntityModel{}).Where(fmt.Sprintf("jpm_id IN (%s)", listID)).Order("id ASC").Find(&result).Error; err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	}

	return &result, nil
}
