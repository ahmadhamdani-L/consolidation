package repository

import (
	"worker-consol/internal/abstraction"
	"worker-consol/internal/model"

	"gorm.io/gorm"
)

type TrialBalance interface {
	Find(ctx *abstraction.Context, m *model.TrialBalanceFilterModel) (*[]model.TrialBalanceEntityModel, error)
	FindByID(ctx *abstraction.Context, id *int) (*model.TrialBalanceEntityModel, error)
	FindByCriteria(ctx *abstraction.Context, m *model.TrialBalanceFilterModel) (*model.TrialBalanceEntityModel, error)
	Get(ctx *abstraction.Context, m *model.TrialBalanceFilterModel) (*model.TrialBalanceEntityModel, error)
	GetCount(ctx *abstraction.Context, m *model.TrialBalanceFilterModel) (*int64, error)
	Update(ctx *abstraction.Context, id *int, e *model.TrialBalanceEntityModel) (*model.TrialBalanceEntityModel, error)
	SetConsolID(ctx *abstraction.Context, id *int, consolID *int) (*model.TrialBalanceEntityModel, error)
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

func (r *trialbalance) FindByID(ctx *abstraction.Context, id *int) (*model.TrialBalanceEntityModel, error) {
	conn := r.CheckTrx(ctx)
	var data model.TrialBalanceEntityModel
	query := conn.Model(data)
	if err := query.Where("id = ?", id).Preload("Company").First(&data).Error; err != nil {
		return nil, err
	}
	return &data, nil
}

func (r *trialbalance) FindByCriteria(ctx *abstraction.Context, m *model.TrialBalanceFilterModel) (*model.TrialBalanceEntityModel, error) {
	conn := r.CheckTrx(ctx)
	var data model.TrialBalanceEntityModel
	query := conn.Model(data)
	query = r.Filter(ctx, query, *m)
	if err := query.First(&data).Error; err != nil {
		return nil, err
	}
	return &data, nil
}

func (r *trialbalance) Update(ctx *abstraction.Context, id *int, e *model.TrialBalanceEntityModel) (*model.TrialBalanceEntityModel, error) {
	conn := r.CheckTrx(ctx)

	if err := conn.Model(e).Where("id = ?", &id).Updates(e).Error; err != nil {
		return nil, err
	}
	if err := conn.Model(e).Where("id = ?", &id).First(e).Error; err != nil {
		return nil, err
	}
	return e, nil
}

func (r *trialbalance) SetConsolID(ctx *abstraction.Context, id *int, consolID *int) (*model.TrialBalanceEntityModel, error) {
	conn := r.CheckTrx(ctx)

	query := conn.Table(model.TrialBalanceEntityModel{}.TableName()).Where("id = ?", &id)
	if consolID == nil {
		query.Update("consolidation_id", nil)
	} else {
		query.Update("consolidation_id", *consolID)
	}

	if err := query.Error; err != nil {
		return nil, err
	}
	data := model.TrialBalanceEntityModel{}
	if err := conn.Model(&model.TrialBalanceEntityModel{}).Where("id = ?", &id).First(&data).Error; err != nil {
		return nil, err
	}
	return &data, nil
}
