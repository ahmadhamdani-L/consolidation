package repository

import (
	"errors"
	"worker-consol/internal/abstraction"
	"worker-consol/internal/model"
	"worker-consol/pkg/constant"

	"gorm.io/gorm"
)

type Consolidation interface {
	Find(ctx *abstraction.Context, m *model.ConsolidationFilterModel) (*[]model.ConsolidationEntityModel, error)
	FindByCriteria(ctx *abstraction.Context, filter *model.ConsolidationFilterModel) (data *model.ConsolidationEntityModel, err error)
	Create(ctx *abstraction.Context, data *model.ConsolidationEntityModel) (*model.ConsolidationEntityModel, error)
	FindByID(ctx *abstraction.Context, id *int) (*model.ConsolidationEntityModel, error)
	Update(ctx *abstraction.Context, id *int, e *model.ConsolidationEntityModel) error
	UpdateStatusModul(ctx *abstraction.Context, companyID *int, period *string, version *int) error
	Count(ctx *abstraction.Context, m *model.ConsolidationFilterModel) (int64, error)
	UpdateStatusJurnal(ctx *abstraction.Context, consolID *int) error
}

type consolidation struct {
	abstraction.Repository
}

func NewConsolidation(db *gorm.DB) *consolidation {
	return &consolidation{
		abstraction.Repository{
			Db: db,
		},
	}
}

func (r *consolidation) Find(ctx *abstraction.Context, m *model.ConsolidationFilterModel) (*[]model.ConsolidationEntityModel, error) {
	conn := r.CheckTrx(ctx)

	var datas []model.ConsolidationEntityModel

	query := conn.Model(&model.ConsolidationEntityModel{})
	//filter
	query = r.Filter(ctx, query, *m)

	if err := query.Find(&datas).Error; err != nil {
		return &datas, err
	}

	return &datas, nil
}

func (r *consolidation) FindByCriteria(ctx *abstraction.Context, filter *model.ConsolidationFilterModel) (data *model.ConsolidationEntityModel, err error) {
	conn := r.CheckTrx(ctx)
	query := conn.Model(&model.ConsolidationEntityModel{})
	query = r.Filter(ctx, query, *filter)
	if err = query.Preload("ConsolidationBridge.Company").First(&data).Error; err != nil {
		return
	}
	if data.ID == 0 {
		err = errors.New("Data Not Found")
	}
	return
}

func (r *consolidation) Create(ctx *abstraction.Context, data *model.ConsolidationEntityModel) (*model.ConsolidationEntityModel, error) {
	conn := r.CheckTrx(ctx)
	var lenExistingData int64
	if err := conn.Model(&model.ConsolidationEntityModel{}).Where("period = ?", data.Period).Where("company_id", data.CompanyID).Count(&lenExistingData).Error; err != nil {
		return nil, err
	}

	data.ConsolidationVersions = int(lenExistingData) + 1
	if err := conn.Create(&data).Error; err != nil {
		return nil, err
	}

	err := conn.Model(&data).First(&data).Error
	if err != nil {
		return nil, err
	}

	return data, err

}

func (r *consolidation) FindByID(ctx *abstraction.Context, id *int) (data *model.ConsolidationEntityModel, err error) {
	conn := r.CheckTrx(ctx)
	if err = conn.Model(&model.ConsolidationEntityModel{}).Where("id = ?", id).First(&data).Error; err != nil {
		return
	}
	return
}

func (r *consolidation) Update(ctx *abstraction.Context, id *int, e *model.ConsolidationEntityModel) error {
	conn := r.CheckTrx(ctx)
	if err := conn.Model(&model.ConsolidationEntityModel{}).Where("id = ?", id).Updates(e).Error; err != nil {
		return err
	}
	return nil
}

func (r *consolidation) UpdateStatusModul(ctx *abstraction.Context, companyID *int, period *string, version *int) error {
	conn := r.CheckTrx(ctx)
	listModul := []string{
		model.AgingUtangPiutangEntityModel{}.TableName(),
		model.EmployeeBenefitEntityModel{}.TableName(),
		model.InvestasiNonTbkEntityModel{}.TableName(),
		model.InvestasiTbkEntityModel{}.TableName(),
		model.MutasiDtaEntityModel{}.TableName(),
		model.MutasiFaEntityModel{}.TableName(),
		model.MutasiIaEntityModel{}.TableName(),
		model.MutasiPersediaanEntityModel{}.TableName(),
		model.MutasiRuaEntityModel{}.TableName(),
		model.PembelianPenjualanBerelasiEntityModel{}.TableName(),
		model.TrialBalanceEntityModel{}.TableName(),
	}
	var trialBalanceData model.TrialBalanceEntityModel
	if err := conn.Model(&model.TrialBalanceEntityModel{}).Where("company_id = ?", companyID).Where("period = DATE(?)", period).Where("versions = ?", version).First(&trialBalanceData).Error; err != nil {
		return err
	}

	for _, modul := range listModul {
		if err := conn.Table(modul).Where("company_id = ?", companyID).Where("period = DATE(?)", period).Where("versions = ?", version).Update("status", constant.MODUL_STATUS_CONSOLIDATE).Error; err != nil {
			return err
		}
	}

	if err := conn.Table(model.AdjustmentEntityModel{}.TableName()).Where("company_id = ?", companyID).Where("period = DATE(?)", period).Where("tb_id = ?", trialBalanceData.ID).Update("status", constant.MODUL_STATUS_CONSOLIDATE).Error; err != nil {
		return err
	}
	return nil
}

func (r *consolidation) UpdateStatusJurnal(ctx *abstraction.Context, consolID *int) error {
	conn := r.CheckTrx(ctx)
	listModul := []string{
		model.JcteEntityModel{}.TableName(),
		model.JelimEntityModel{}.TableName(),
		model.JpmEntityModel{}.TableName(),
	}

	for _, modul := range listModul {
		if err := conn.Table(modul).Where("consolidation_id = ?", consolID).Update("status", constant.MODUL_STATUS_CONSOLIDATE).Error; err != nil {
			return err
		}
	}

	return nil
}

func (r *consolidation) Count(ctx *abstraction.Context, m *model.ConsolidationFilterModel) (int64, error) {
	conn := r.CheckTrx(ctx)
	var count int64
	query := conn.Model(&model.ConsolidationEntityModel{})
	query = r.Filter(ctx, query, *m)
	if err := query.Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}
