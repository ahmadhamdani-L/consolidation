package repository

import (
	"worker-consol/internal/abstraction"
	"worker-consol/internal/model"

	"gorm.io/gorm"
)

type TrialBalanceDetail interface {
	Find(ctx *abstraction.Context, m *model.TrialBalanceDetailFilterModel) (*[]model.TrialBalanceDetailEntityModel, error)
	FindWithCode(ctx *abstraction.Context, id *int, code *string) (*[]model.TrialBalanceDetailEntityModel, error)
	FindByCriteria(ctx *abstraction.Context, m *model.TrialBalanceFilterModel, code *string) (*[]model.TrialBalanceDetailEntityModel, error)
	FindSummaryByCompanyCode(ctx *abstraction.Context, list *string, code *string) (*model.TrialBalanceDetailEntityModel, error)
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
	if m.TrialBalanceID != nil {
		query = query.Where("formatter_bridges_id IN (?)", conn.Model(&model.FormatterBridgesEntityModel{}).Select("id").Where("source = ?", "TRIAL-BALANCE").Where("trx_ref_id = ?", m.TrialBalanceID))
	}
	if err := query.Find(&datas).Error; err != nil {
		return &datas, err
	}

	return &datas, nil
}

func (r *trialbalancedetail) FindWithCode(ctx *abstraction.Context, id *int, code *string) (*[]model.TrialBalanceDetailEntityModel, error) {
	conn := r.CheckTrx(ctx)

	var data []model.TrialBalanceDetailEntityModel
	if err := conn.Where("formatter_bridges_id = ?", id).Where("code LIKE ?", *code+"%").Find(&data).Error; err != nil {
		return &data, err
	}
	return &data, nil
}

func (r *trialbalancedetail) FindByCriteria(ctx *abstraction.Context, m *model.TrialBalanceFilterModel, code *string) (*[]model.TrialBalanceDetailEntityModel, error) {
	conn := r.CheckTrx(ctx)

	var dataTB model.TrialBalanceEntityModel

	queryTB := conn.Model(&model.TrialBalanceEntityModel{})
	queryTB = queryTB.Where("versions = ?", m.Versions).Where("period = ?", m.Period).Where("company_id = ?", m.CompanyID)

	if err := queryTB.First(&dataTB).Error; err != nil {
		return nil, err
	}

	var datas []model.TrialBalanceDetailEntityModel

	query := conn.Model(&model.TrialBalanceDetailEntityModel{})
	//filter
	query = r.Filter(ctx, query, *m)
	query = query.Where("formatter_bridges_id = ?", conn.Model(&model.FormatterBridgesEntityModel{}).Select("id").Where("source = ?", "TRIAL-BALANCE").Where("trx_ref_id = ?", dataTB.ID)).Where("code = ?", code)
	if err := query.First(&datas).Error; err != nil {
		return &datas, err
	}

	return &datas, nil
}

func (r *trialbalancedetail) FindSummaryByCompanyCode(ctx *abstraction.Context, list *string, code *string) (*model.TrialBalanceDetailEntityModel, error) {
	conn := r.CheckTrx(ctx)

	var detailData model.TrialBalanceDetailEntityModel
	querySum := conn.Model(&model.TrialBalanceDetailEntityModel{}).Where("code LIKE ?", *code+"%").Where("consolidation_bridge_id IN (?)", list)
	querySum = querySum.Select("SUM(amount_after_aje) amount_after_aje")
	if err := querySum.First(&detailData).Error; err != nil {
		return nil, err
	}

	return &detailData, nil
}
