package repository

import (
	"errors"
	"fmt"
	"worker-validation/internal/abstraction"
	"worker-validation/internal/model"

	"gorm.io/gorm"
)

type TrialBalanceDetail interface {
	Find(ctx *abstraction.Context, m *model.TrialBalanceDetailFilterModel) (*[]model.TrialBalanceDetailEntityModel, error)
	FindWithCode(ctx *abstraction.Context, code *string) (*[]model.TrialBalanceDetailEntityModel, error)
	FindSummary(ctx *abstraction.Context, m *model.TrialBalanceDetailFilterModel) (*model.TrialBalanceDetailEntityModel, error)
	FindByCriteria(ctx *abstraction.Context, filter *model.TrialBalanceDetailFilterModel) (data *model.TrialBalanceDetailEntityModel, err error)
	FindWithCodes(ctx *abstraction.Context, fmtID *int ,code *string) (*model.TrialBalanceDetailEntityModel, error)
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

func (r *trialbalancedetail) FindWithCodes(ctx *abstraction.Context, fmtID *int ,code *string) (*model.TrialBalanceDetailEntityModel, error) {
	conn := r.CheckTrx(ctx)

	var data model.TrialBalanceDetailEntityModel
	tmp := fmt.Sprintf("%s", *code)
	if err := conn.Where("formatter_bridges_id = ? AND code LIKE ?", fmtID, tmp+"%").First(&data).Error; err != nil {
		return &data, err
	}
	return &data, nil
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

func (r *trialbalancedetail) FindSummary(ctx *abstraction.Context, m *model.TrialBalanceDetailFilterModel) (*model.TrialBalanceDetailEntityModel, error) {
	conn := r.CheckTrx(ctx)

	var datas model.TrialBalanceDetailEntityModel

	query := conn.Model(&model.TrialBalanceDetailEntityModel{})
	//filter
	query = r.Filter(ctx, query, *m)
	query = query.Select("SUM(amount_before_aje) amount_before_aje, SUM(amount_after_aje) amount_after_aje, SUM(amount_aje_cr) amount_aje_cr, SUM(amount_aje_dr) amount_aje_dr")
	if err := query.Find(&datas).Error; err != nil {
		return &datas, err
	}

	return &datas, nil
}

func (r *trialbalancedetail) FindByCriteria(ctx *abstraction.Context, filter *model.TrialBalanceDetailFilterModel) (data *model.TrialBalanceDetailEntityModel, err error) {
	conn := r.CheckTrx(ctx)
	query := conn.Model(&model.TrialBalanceDetailEntityModel{})
	query = r.Filter(ctx, query, *filter)
	if err = query.First(&data).Error; err != nil {
		return
	}
	if data.ID == 0 {
		err = errors.New("Data Not Found")
	}
	return
}
