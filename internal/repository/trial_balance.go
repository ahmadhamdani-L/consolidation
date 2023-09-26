package repository

import (
	"worker/internal/abstraction"
	"worker/internal/model"

	"gorm.io/gorm"
)

type TrialBalance interface {
	Find(ctx *abstraction.Context, m *model.TrialBalanceFilterModel) (*[]model.TrialBalanceEntityModel, error)
	Get(ctx *abstraction.Context, m *model.TrialBalanceFilterModel) (*model.TrialBalanceEntityModel, error)
	GetCount(ctx *abstraction.Context, m *model.TrialBalanceFilterModel) (*int64, error)
}

type trialbalance struct {
	abstraction.Repository
}

func NewTrialBalance(db *gorm.DB) *trialbalance {
	return &trialbalance{
		abstraction.Repository{
			Db: db,
		},
	}
}

func (r *trialbalance) Find(ctx *abstraction.Context, m *model.TrialBalanceFilterModel) (*[]model.TrialBalanceEntityModel, error) {
	conn := r.CheckTrx(ctx)

	var datas []model.TrialBalanceEntityModel

	query := conn.Model(&model.TrialBalanceEntityModel{})
	//filter
	query = r.Filter(ctx, query, *m)

	if err := query.Find(&datas).Error; err != nil {
		return &datas, err
	}

	return &datas, nil
}

func (r *trialbalance) Get(ctx *abstraction.Context, m *model.TrialBalanceFilterModel) (*model.TrialBalanceEntityModel, error) {
	var datas model.TrialBalanceEntityModel

	conn := r.CheckTrx(ctx)
	query := conn.Model(&model.TrialBalanceEntityModel{})
	//filter
	query = r.Filter(ctx, query, *m)

	if err := query.Preload("Company").Find(&datas).Error; err != nil {
		return &datas, err
	}

	var formatterBridges []model.FormatterBridgesEntityModel

	query = conn.Model(&model.FormatterBridgesEntityModel{}).Where("trx_ref_id = ?", &datas.ID).Where("source = ?", "TRIAL-BALANCE")
	if err := query.Find(&formatterBridges).Error; err != nil {
		return &datas, err
	}
	datas.FormatterBridges = formatterBridges

	return &datas, nil
}

func (r *trialbalance) GetCount(ctx *abstraction.Context, m *model.TrialBalanceFilterModel) (*int64, error) {
	var jmlData int64
	conn := r.CheckTrx(ctx)
	query := conn.Model(&model.TrialBalanceEntityModel{})
	//filter
	query = r.Filter(ctx, query, *m)

	if err := query.Count(&jmlData).Error; err != nil {
		return &jmlData, err
	}

	return &jmlData, nil
}
