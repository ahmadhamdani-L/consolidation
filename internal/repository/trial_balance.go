package repository

import (
	"worker/internal/abstraction"
	"worker/internal/model"

	"gorm.io/gorm"
)

type TrialBalance interface {
	Find(ctx *abstraction.Context, m *model.TrialBalanceFilterModel) (*[]model.TrialBalanceEntityModel, error)
	Get(ctx *abstraction.Context, m *model.TrialBalanceFilterModel, opts []string) (*model.TrialBalanceEntityModel, error)
	GetCount(ctx *abstraction.Context, m *model.TrialBalanceFilterModel) (*int64, error)
	Create(ctx *abstraction.Context, e *model.TrialBalanceEntityModel) (*model.TrialBalanceEntityModel, error)
	FindByID(ctx *abstraction.Context, version *int, company *int, period *string) (*model.TrialBalanceEntityModel, error)
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

func (r *trialbalance) Create(ctx *abstraction.Context, e *model.TrialBalanceEntityModel) (*model.TrialBalanceEntityModel, error) {
	conn := r.CheckTrx(ctx)

	if err := conn.Create(e).Error; err != nil {
		return nil, err
	}
	if err := conn.Model(e).Preload("UserCreated").Preload("UserModified").First(e).Error; err != nil {
		return nil, err
	}
	e.UserCreatedString = e.UserCreated.Name
	e.UserModifiedString = &e.UserModified.Name

	return e, nil
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

func (r *trialbalance) Get(ctx *abstraction.Context, m *model.TrialBalanceFilterModel, opts []string) (*model.TrialBalanceEntityModel, error) {
	var datas model.TrialBalanceEntityModel

	conn := r.CheckTrx(ctx)
	query := conn.Model(&model.TrialBalanceEntityModel{})
	//filter
	query = r.Filter(ctx, query, *m)

	for _, val := range opts {
		if val == "Company" || val == "Formatter" || val == "TrialBalanceDetail" {
			query = query.Preload(val)
		}
	}

	if err := query.Find(&datas).Error; err != nil {
		return &datas, err
	}

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

func (r *trialbalance) FindByID(ctx *abstraction.Context, version *int, company *int, period *string) (*model.TrialBalanceEntityModel, error) {
	conn := r.CheckTrx(ctx)

	var data model.TrialBalanceEntityModel
	err := conn.Where("versions = ? AND company_id = ? AND period = ?", version, company, period).Find(&data).Error
	if err != nil {
		return nil, err
	}
	return &data, nil
}

func (r *trialbalance) Update(ctx *abstraction.Context, id *int, e *model.TrialBalanceEntityModel) (*model.TrialBalanceEntityModel, error) {
	conn := r.CheckTrx(ctx)

	if err := conn.Model(e).Where("id = ?", &id).Updates(e).Error; err != nil {
		return nil, err
	}

	return e, nil

}
