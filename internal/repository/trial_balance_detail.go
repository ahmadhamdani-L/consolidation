package repository

import (
	"fmt"
	"worker/internal/abstraction"
	"worker/internal/model"

	"gorm.io/gorm"
)

type TrialBalanceDetail interface {
	Find(ctx *abstraction.Context, m *model.TrialBalanceDetailFilterModel) (*[]model.TrialBalanceDetailEntityModel, error)
	FindWithCode(ctx *abstraction.Context, code *string) (*[]model.TrialBalanceDetailEntityModel, error)
	FindToExport(ctx *abstraction.Context, code *string, fmtBridgeID *int) (*[]model.TrialBalanceDetailEntityModel, error)
}

type trialbalancedetail struct {
	abstraction.Repository
}

func NewTrialBalanceDetail(db *gorm.DB) *trialbalancedetail {
	return &trialbalancedetail{
		abstraction.Repository{
			Db: db,
		},
	}
}

func (r *trialbalancedetail) Find(ctx *abstraction.Context, m *model.TrialBalanceDetailFilterModel) (*[]model.TrialBalanceDetailEntityModel, error) {
	conn := r.CheckTrx(ctx)

	var datas []model.TrialBalanceDetailEntityModel

	query := conn.Model(&model.TrialBalanceDetailEntityModel{})
	//filter
	query = r.Filter(ctx, query, *m)
	query = query.Where("formatter_bridges_id IN (?)", conn.Model(&model.FormatterBridgesEntityModel{}).Select("id").Where("source = ?", "TRIAL-BALANCE").Where("trx_ref_id = ?", m.TrialBalanceID))
	if err := query.Find(&datas).Error; err != nil {
		return &datas, err
	}

	return &datas, nil
}

func (r *trialbalancedetail) FindWithCode(ctx *abstraction.Context, code *string) (*[]model.TrialBalanceDetailEntityModel, error) {
	conn := r.CheckTrx(ctx)

	var data []model.TrialBalanceDetailEntityModel
	tmp := fmt.Sprintf("%s", *code)
	if err := conn.Where("code LIKE ?", tmp+"%").Find(&data).Error; err != nil {
		return &data, err
	}
	return &data, nil
}

func (r *trialbalancedetail) FindToExport(ctx *abstraction.Context, code *string, fmtBridgeID *int) (*[]model.TrialBalanceDetailEntityModel, error) {
	conn := r.CheckTrx(ctx)

	var data []model.TrialBalanceDetailEntityModel
	if err := conn.Where("code ILIKE ?", *code+"%").Where("formatter_bridges_id = ?", *fmtBridgeID).Find(&data).Error; err != nil {
		return &data, err
	}
	return &data, nil
}
