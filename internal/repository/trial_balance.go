package repository

import (
	"errors"
	"worker-validation/internal/abstraction"
	"worker-validation/internal/model"

	"gorm.io/gorm"
)

type TrialBalance interface {
	Find(ctx *abstraction.Context, m *model.TrialBalanceFilterModel) (*[]model.TrialBalanceEntityModel, error)
	Get(ctx *abstraction.Context, m *model.TrialBalanceFilterModel) (*model.TrialBalanceEntityModel, error)
	GetCount(ctx *abstraction.Context, m *model.TrialBalanceFilterModel) (*int64, error)
	FindByID(ctx *abstraction.Context, id *int) (data *model.TrialBalanceEntityModel, err error)
	FindByCriteria(ctx *abstraction.Context, filter *model.FilterData) (data *model.TrialBalanceEntityModel, err error)
	Update(ctx *abstraction.Context, id *int, e *model.TrialBalanceEntityModel) (*model.TrialBalanceEntityModel, error)
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

func (r *trialbalance) FindByID(ctx *abstraction.Context, id *int) (data *model.TrialBalanceEntityModel, err error) {
	conn := r.CheckTrx(ctx)
	query := conn.Model(&model.TrialBalanceEntityModel{})
	// err = query.Where("id = ?", id).First(&data).Error
	if err = query.Where("id = ?", id).Preload("Company").First(&data).Error; err != nil {
		return
	}
	if data.ID == 0 {
		err = errors.New("Data Not Found")
	}
	return
}

func (r *trialbalance) FindByCriteria(ctx *abstraction.Context, filter *model.FilterData) (data *model.TrialBalanceEntityModel, err error) {
	conn := r.CheckTrx(ctx)
	query := conn.Model(&model.TrialBalanceEntityModel{})
	if &filter.Status != nil && filter.Status != 0 {
		query = query.Where("status = ?", filter.Status)
	}
	if err = query.Where("company_id = ?", filter.CompanyID).Where("period = ?", filter.Period).Where("versions = ?", filter.Versions).First(&data).Error; err != nil {
		return
	}
	if data.ID == 0 {
		err = errors.New("Data Not Found")
	}
	return
}

func (r *trialbalance) Update(ctx *abstraction.Context, id *int, e *model.TrialBalanceEntityModel) (*model.TrialBalanceEntityModel, error) {
	conn := r.CheckTrx(ctx)
	if err := conn.Model(e).Where("id = ?", id).Updates(e).Error; err != nil {
		return nil, err
	}
	if err := conn.Where("id = ?", id).First(&e).Error; err != nil {
		return nil, err
	}

	return e, nil
}
