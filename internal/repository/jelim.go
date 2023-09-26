package repository

import (
	"fmt"
	"strings"

	"gorm.io/gorm"
	// "gorm.io/gorm/clause"
	"worker/internal/abstraction"
	"worker/internal/model"
)

type Jelim interface {
	Find(ctx *abstraction.Context, m *model.JelimFilterModel) (*[]model.JelimEntityModel, error)
	FindByID(ctx *abstraction.Context, id *int) (*model.JelimEntityModel, error)
	ExportOne(ctx *abstraction.Context, e *model.JelimFilterModel) (*model.JelimEntityModel, error)
	FindSummary(ctx *abstraction.Context, consolID *int) (*model.JelimDetailEntityModel, error)
	ExportAll(ctx *abstraction.Context, consolidationID *int) (*[]model.JelimDetailEntityModel, error)
}

type jelim struct {
	abstraction.Repository
}

func NewJelim(db *gorm.DB) *jelim {
	return &jelim{
		abstraction.Repository{
			Db: db,
		},
	}
}

func (r *jelim) Find(ctx *abstraction.Context, m *model.JelimFilterModel) (*[]model.JelimEntityModel, error) {
	conn := r.CheckTrx(ctx)

	var datas []model.JelimEntityModel
	query := conn.Model(&model.JelimEntityModel{})

	// filter
	query = r.Filter(ctx, query, *m)

	err := query.Preload("JelimDetail").Find(&datas).Error
	if err != nil {
		return &datas, err
	}

	return &datas, nil
}

func (r *jelim) FindByID(ctx *abstraction.Context, id *int) (*model.JelimEntityModel, error) {
	conn := r.CheckTrx(ctx)

	var data model.JelimEntityModel
	err := conn.Where("id = ?", id).First(&data).Error
	if err != nil {
		return nil, err
	}
	return &data, nil
}

func (r *jelim) ExportOne(ctx *abstraction.Context, e *model.JelimFilterModel) (*model.JelimEntityModel, error) {
	conn := r.CheckTrx(ctx)
	var data model.JelimEntityModel
	query := conn.Model(&model.JelimEntityModel{}).Preload("Company").Preload("JelimDetail").Where("consolidation_id = ?", &e.ConsolidationID).Where("period = ?", &e.Period).Where("company_id = ?", &e.CompanyID).Find(&data)
	if err := query.Error; err != nil {
		return nil, err
	}
	return &data, nil
}

func (r *jelim) FindSummary(ctx *abstraction.Context, consolID *int) (*model.JelimDetailEntityModel, error) {
	conn := r.CheckTrx(ctx)
	var data []string
	if err := conn.Model(&model.JelimEntityModel{}).Where("status != 0").Where("consolidation_id = ?", &consolID).Pluck("id", &data).Error; err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	}

	if len(data) == 0 {
		return &model.JelimDetailEntityModel{}, nil
	}

	listID := strings.Join(data, ",")
	var sumData model.JelimDetailEntityModel
	if err := conn.Model(&model.JelimDetailEntityModel{}).Where(fmt.Sprintf("jelim_id IN (%s)", listID)).Select("SUM(balance_sheet_cr) balance_sheet_cr, SUM(balance_sheet_dr) balance_sheet_dr, SUM(income_statement_cr) income_statement_cr, SUM(income_statement_dr) income_statement_dr").Find(&sumData).Error; err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	}

	return &sumData, nil
}

func (r *jelim) ExportAll(ctx *abstraction.Context, consolidationID *int) (*[]model.JelimDetailEntityModel, error) {
	conn := r.CheckTrx(ctx)
	var data []string
	if err := conn.Model(&model.JelimEntityModel{}).Where("status != 4").Where("consolidation_id = ?", &consolidationID).Pluck("id", &data).Error; err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	}

	if len(data) == 0 {
		return &[]model.JelimDetailEntityModel{}, nil
	}

	listID := strings.Join(data, ",")
	var result []model.JelimDetailEntityModel
	if err := conn.Model(&model.JelimDetailEntityModel{}).Where(fmt.Sprintf("jelim_id IN (%s)", listID)).Order("id ASC").Find(&result).Error; err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	}

	return &result, nil
}
