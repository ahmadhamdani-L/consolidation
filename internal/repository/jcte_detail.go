package repository

import (
	"fmt"
	"strings"
	"worker-consol/internal/abstraction"
	"worker-consol/internal/model"

	"gorm.io/gorm"
)

type JcteDetail interface {
	Find(ctx *abstraction.Context, m *model.JcteDetailFilterModel) (*[]model.JcteDetailEntityModel, error)
	FindByID(ctx *abstraction.Context, id *int) (*model.JcteDetailEntityModel, error)
	FindWithCode(ctx *abstraction.Context, code *string) (*[]model.JcteDetailEntityModel, error)
	FindSummary(ctx *abstraction.Context, tbID *int, code *string) (*model.JcteDetailEntityModel, error)
}

type jctedetail struct {
	abstraction.Repository
}

func NewJcteDetail(db *gorm.DB) *jctedetail {
	return &jctedetail{
		abstraction.Repository{
			Db: db,
		},
	}
}

func (r *jctedetail) Find(ctx *abstraction.Context, m *model.JcteDetailFilterModel) (*[]model.JcteDetailEntityModel, error) {
	conn := r.CheckTrx(ctx)

	var datas []model.JcteDetailEntityModel

	query := conn.Model(&model.JcteDetailEntityModel{})
	//filter
	query = r.Filter(ctx, query, *m)

	if err := query.Preload("Jcte").Find(&datas).Error; err != nil {
		return &datas, err
	}

	return &datas, nil
}

func (r *jctedetail) FindWithCode(ctx *abstraction.Context, code *string) (*[]model.JcteDetailEntityModel, error) {
	conn := r.CheckTrx(ctx)

	var data []model.JcteDetailEntityModel
	tmp := fmt.Sprintf("%s", *code)
	if err := conn.Where("code LIKE ?", tmp+"%").Find(&data).Error; err != nil {
		return &data, err
	}
	return &data, nil
}

func (r *jctedetail) FindByID(ctx *abstraction.Context, id *int) (*model.JcteDetailEntityModel, error) {
	conn := r.CheckTrx(ctx)

	var data model.JcteDetailEntityModel
	if err := conn.Where("id = ?", &id).Preload("Jcte").First(&data).Error; err != nil {
		return &data, err
	}
	return &data, nil
}

func (r *jctedetail) FindSummary(ctx *abstraction.Context, tbID *int, code *string) (*model.JcteDetailEntityModel, error) {
	conn := r.CheckTrx(ctx)
	statusDeleted := 4
	var tmp []string
	query := conn.Model(&model.JcteEntityModel{}).Where("consolidation_id = ? AND status != ?", tbID, statusDeleted)
	if err := query.Pluck("id", &tmp).Error; err != nil {
		if err.Error() == "record not found" {
			return &model.JcteDetailEntityModel{}, nil
		}
		return nil, err
	}

	if len(tmp) == 0 {
		return &model.JcteDetailEntityModel{}, nil
	}

	listID := strings.Join(tmp, ",")
	var detailData model.JcteDetailEntityModel
	querySum := conn.Model(&model.JcteDetailEntityModel{}).Where("coa_code LIKE ?", *code+"%").Where(fmt.Sprintf("jcte_id IN (%s)", listID))
	querySum = querySum.Select("SUM(balance_sheet_dr) balance_sheet_dr, SUM(balance_sheet_cr) balance_sheet_cr, SUM(income_statement_dr) income_statement_dr, SUM(income_statement_cr) income_statement_cr")
	if err := querySum.Find(&detailData).Error; err != nil {
		if err.Error() == "record not found" {
			return &model.JcteDetailEntityModel{}, nil
		}
		return nil, err
	}

	return &detailData, nil
}
