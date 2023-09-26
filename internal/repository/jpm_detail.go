package repository

import (
	"fmt"
	"strings"
	"worker-consol/internal/abstraction"
	"worker-consol/internal/model"

	"gorm.io/gorm"
)

type JpmDetail interface {
	Find(ctx *abstraction.Context, m *model.JpmDetailFilterModel) (*[]model.JpmDetailEntityModel, error)
	FindByID(ctx *abstraction.Context, id *int) (*model.JpmDetailEntityModel, error)
	FindWithCode(ctx *abstraction.Context, code *string) (*[]model.JpmDetailEntityModel, error)
	FindSummary(ctx *abstraction.Context, consolID *int, code *string) (*model.JpmDetailEntityModel, error)
}

type jpmdetail struct {
	abstraction.Repository
}

func NewJpmDetail(db *gorm.DB) *jpmdetail {
	return &jpmdetail{
		abstraction.Repository{
			Db: db,
		},
	}
}

func (r *jpmdetail) Find(ctx *abstraction.Context, m *model.JpmDetailFilterModel) (*[]model.JpmDetailEntityModel, error) {
	conn := r.CheckTrx(ctx)

	var datas []model.JpmDetailEntityModel

	query := conn.Model(&model.JpmDetailEntityModel{})
	//filter
	query = r.Filter(ctx, query, *m)

	if err := query.Find(&datas).Error; err != nil {
		return &datas, err
	}

	return &datas, nil
}

func (r *jpmdetail) FindWithCode(ctx *abstraction.Context, code *string) (*[]model.JpmDetailEntityModel, error) {
	conn := r.CheckTrx(ctx)

	var data []model.JpmDetailEntityModel
	tmp := fmt.Sprintf("%s", *code)
	if err := conn.Where("code LIKE ?", tmp+"%").Find(&data).Error; err != nil {
		return &data, err
	}
	return &data, nil
}

func (r *jpmdetail) FindByID(ctx *abstraction.Context, id *int) (*model.JpmDetailEntityModel, error) {
	conn := r.CheckTrx(ctx)

	var data model.JpmDetailEntityModel
	if err := conn.Where("id = ?", &id).First(&data).Error; err != nil {
		return &data, err
	}
	return &data, nil
}

func (r *jpmdetail) FindSummary(ctx *abstraction.Context, consolID *int, code *string) (*model.JpmDetailEntityModel, error) {
	conn := r.CheckTrx(ctx)
	statusDeleted := 4
	query := conn.Model(&model.JpmEntityModel{}).Where("consolidation_id = ? AND status != ?", consolID, statusDeleted)
	var tmp []string
	if err := query.Pluck("id", &tmp).Error; err != nil {
		if err.Error() == "record not found" {
			return &model.JpmDetailEntityModel{}, nil
		}
		return nil, err
	}

	if len(tmp) == 0 {
		return &model.JpmDetailEntityModel{}, nil
	}

	listID := strings.Join(tmp, ",")
	var detailData model.JpmDetailEntityModel
	querySum := conn.Model(&model.JpmDetailEntityModel{}).Where("coa_code LIKE ?", *code+"%").Where(fmt.Sprintf("jpm_id IN (%s)", listID))
	querySum = querySum.Select("SUM(balance_sheet_dr) balance_sheet_dr, SUM(balance_sheet_cr) balance_sheet_cr, SUM(income_statement_dr) income_statement_dr, SUM(income_statement_cr) income_statement_cr")
	if err := querySum.Find(&detailData).Error; err != nil {
		if err.Error() == "record not found" {
			return &model.JpmDetailEntityModel{}, nil
		}
		return nil, err
	}

	return &detailData, nil
}
