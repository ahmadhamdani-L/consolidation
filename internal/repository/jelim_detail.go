package repository

import (
	"fmt"
	"strings"
	"worker-consol/internal/abstraction"
	"worker-consol/internal/model"

	"gorm.io/gorm"
)

type JelimDetail interface {
	Find(ctx *abstraction.Context, m *model.JelimDetailFilterModel) (*[]model.JelimDetailEntityModel, error)
	FindByID(ctx *abstraction.Context, id *int) (*model.JelimDetailEntityModel, error)
	FindWithCode(ctx *abstraction.Context, code *string) (*[]model.JelimDetailEntityModel, error)
	FindSummary(ctx *abstraction.Context, tbID *int, code *string) (*model.JelimDetailEntityModel, error)
}

type jelimdetail struct {
	abstraction.Repository
}

func NewJelimDetail(db *gorm.DB) *jelimdetail {
	return &jelimdetail{
		abstraction.Repository{
			Db: db,
		},
	}
}

func (r *jelimdetail) Find(ctx *abstraction.Context, m *model.JelimDetailFilterModel) (*[]model.JelimDetailEntityModel, error) {
	conn := r.CheckTrx(ctx)

	var datas []model.JelimDetailEntityModel

	query := conn.Model(&model.JelimDetailEntityModel{})
	//filter
	query = r.Filter(ctx, query, *m)
	if err := query.Preload("Jelim").Find(&datas).Error; err != nil {
		return &datas, err
	}

	return &datas, nil
}

func (r *jelimdetail) FindWithCode(ctx *abstraction.Context, code *string) (*[]model.JelimDetailEntityModel, error) {
	conn := r.CheckTrx(ctx)

	var data []model.JelimDetailEntityModel
	tmp := fmt.Sprintf("%s", *code)
	if err := conn.Where("code LIKE ?", tmp+"%").Find(&data).Error; err != nil {
		return &data, err
	}
	return &data, nil
}

func (r *jelimdetail) FindByID(ctx *abstraction.Context, id *int) (*model.JelimDetailEntityModel, error) {
	conn := r.CheckTrx(ctx)

	var data model.JelimDetailEntityModel
	if err := conn.Where("id = ?", &id).Preload("Jelim").First(&data).Error; err != nil {
		return &data, err
	}
	return &data, nil
}

func (r *jelimdetail) FindSummary(ctx *abstraction.Context, tbID *int, code *string) (*model.JelimDetailEntityModel, error) {
	conn := r.CheckTrx(ctx)
	statusDeleted := 4
	query := conn.Model(&model.JelimEntityModel{}).Where("consolidation_id = ? AND status != ?", tbID, statusDeleted)
	var tmp []string
	if err := query.Pluck("id", &tmp).Error; err != nil {
		if err.Error() == "record not found" {
			return &model.JelimDetailEntityModel{}, nil
		}
		return nil, err
	}

	if len(tmp) == 0 {
		return &model.JelimDetailEntityModel{}, nil
	}

	listID := strings.Join(tmp, ",")
	var detailData model.JelimDetailEntityModel
	querySum := conn.Model(&model.JelimDetailEntityModel{}).Where("coa_code LIKE ?", *code+"%").Where(fmt.Sprintf("jelim_id IN (%s)", listID))
	querySum = querySum.Select("SUM(balance_sheet_dr) balance_sheet_dr, SUM(balance_sheet_cr) balance_sheet_cr, SUM(income_statement_dr) income_statement_dr, SUM(income_statement_cr) income_statement_cr")
	if err := querySum.Find(&detailData).Error; err != nil {
		if err.Error() == "record not found" {
			return &model.JelimDetailEntityModel{}, nil
		}
		return nil, err
	}

	return &detailData, nil
}
