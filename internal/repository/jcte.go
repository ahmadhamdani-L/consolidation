package repository

import (
	"fmt"
	"strings"

	"gorm.io/gorm"
	// "gorm.io/gorm/clause"
	"worker/internal/abstraction"
	"worker/internal/model"
)

type Jcte interface {
	Find(ctx *abstraction.Context, m *model.JcteFilterModel) (*[]model.JcteEntityModel, error)
	FindByID(ctx *abstraction.Context, id *int) (*model.JcteEntityModel, error)
	ExportOne(ctx *abstraction.Context, e *model.JcteFilterModel) (*model.JcteEntityModel, error)
	FindSummary(ctx *abstraction.Context, consolID *int) (*model.JcteDetailEntityModel, error)
	ExportAll(ctx *abstraction.Context, consolidationID *int) (*[]model.JcteDetailEntityModel, error)
}

type jcte struct {
	abstraction.Repository
}

func NewJcte(db *gorm.DB) *jcte {
	return &jcte{
		abstraction.Repository{
			Db: db,
		},
	}
}

func (r *jcte) Find(ctx *abstraction.Context, m *model.JcteFilterModel) (*[]model.JcteEntityModel, error) {
	conn := r.CheckTrx(ctx)

	var datas []model.JcteEntityModel

	query := conn.Model(&model.JcteEntityModel{})

	// filter
	query = r.Filter(ctx, query, *m)
	err := query.Preload("JcteDetail").Find(&datas).Error
	if err != nil {
		return &datas, err
	}

	return &datas, nil
}

func (r *jcte) FindByID(ctx *abstraction.Context, id *int) (*model.JcteEntityModel, error) {
	conn := r.CheckTrx(ctx)

	var data model.JcteEntityModel

	err := conn.Where("id = ?", id).Preload("JcteDetail").First(&data).Error
	if err != nil {
		return nil, err
	}
	return &data, nil
}

func (r *jcte) ExportOne(ctx *abstraction.Context, e *model.JcteFilterModel) (*model.JcteEntityModel, error) {
	conn := r.CheckTrx(ctx)
	var data model.JcteEntityModel
	query := conn.Model(&model.JcteEntityModel{}).Preload("Company").Preload("JcteDetail").Where("consolidation_id = ?", &e.ConsolidationID).Where("period = ?", &e.Period).Where("company_id = ?", &e.CompanyID).Find(&data)
	if err := query.Error; err != nil {
		return nil, err
	}
	return &data, nil
}

func (r *jcte) FindSummary(ctx *abstraction.Context, consolID *int) (*model.JcteDetailEntityModel, error) {
	conn := r.CheckTrx(ctx)
	var data []string
	if err := conn.Model(&model.JcteEntityModel{}).Where("status != 4").Where("consolidation_id = ?", &consolID).Pluck("id", &data).Error; err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	}

	if len(data) == 0 {
		return &model.JcteDetailEntityModel{}, nil
	}

	listID := strings.Join(data, ",")
	var sumData model.JcteDetailEntityModel
	if err := conn.Model(&model.JcteDetailEntityModel{}).Where(fmt.Sprintf("jcte_id IN (%s)", listID)).Select("SUM(balance_sheet_cr) balance_sheet_cr, SUM(balance_sheet_dr) balance_sheet_dr, SUM(income_statement_cr) income_statement_cr, SUM(income_statement_dr) income_statement_dr").Find(&sumData).Error; err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	}

	return &sumData, nil
}

func (r *jcte) ExportAll(ctx *abstraction.Context, consolidationID *int) (*[]model.JcteDetailEntityModel, error) {
	conn := r.CheckTrx(ctx)
	var data []string
	if err := conn.Model(&model.JcteEntityModel{}).Where("status != 4").Where("consolidation_id = ?", &consolidationID).Pluck("id", &data).Error; err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	}

	if len(data) == 0 {
		return &[]model.JcteDetailEntityModel{}, nil
	}

	listID := strings.Join(data, ",")
	var result []model.JcteDetailEntityModel
	if err := conn.Model(&model.JcteDetailEntityModel{}).Where(fmt.Sprintf("jcte_id IN (%s)", listID)).Order("id ASC").Find(&result).Error; err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	}

	return &result, nil
}
