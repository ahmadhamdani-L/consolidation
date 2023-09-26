package repository

import (
	"errors"
	"fmt"
	"worker-validation/internal/abstraction"
	"worker-validation/internal/model"

	"gorm.io/gorm"
)

type AgingUtangPiutangDetail interface {
	Find(ctx *abstraction.Context, m *model.AgingUtangPiutangDetailFilterModel) (*[]model.AgingUtangPiutangDetailEntityModel, error)
	FindWithCode(ctx *abstraction.Context, code *string) (*[]model.AgingUtangPiutangDetailEntityModel, error)
	FindByCriteria(ctx *abstraction.Context, filter *model.AgingUtangPiutangDetailFilterModel) (data *model.AgingUtangPiutangDetailEntityModel, err error)
}

type agingutangpiutangdetail struct {
	abstraction.Repository
}

func NewAgingUtangPiutangDetail(db *gorm.DB) *agingutangpiutangdetail {
	return &agingutangpiutangdetail{
		abstraction.Repository{
			Db: db,
		},
	}
}

func (r *agingutangpiutangdetail) Find(ctx *abstraction.Context, m *model.AgingUtangPiutangDetailFilterModel) (*[]model.AgingUtangPiutangDetailEntityModel, error) {
	conn := r.CheckTrx(ctx)

	var datas []model.AgingUtangPiutangDetailEntityModel

	query := conn.Model(&model.AgingUtangPiutangDetailEntityModel{})
	//filter
	query = r.Filter(ctx, query, *m)
	query = query.Where("formatter_bridges_id IN (?)", conn.Model(&model.FormatterBridgesEntityModel{}).Select("id").Where("source = ?", "AGING-UTANG-PIUTANG").Where("trx_ref_id = ?", m.AgingUtangPiutangID))
	if err := query.Find(&datas).Error; err != nil {
		return &datas, err
	}

	return &datas, nil
}

func (r *agingutangpiutangdetail) FindWithCode(ctx *abstraction.Context, code *string) (*[]model.AgingUtangPiutangDetailEntityModel, error) {
	conn := r.CheckTrx(ctx)

	var data []model.AgingUtangPiutangDetailEntityModel
	tmp := fmt.Sprintf("%s", *code)
	if err := conn.Where("code ILIKE ?", tmp+"%").Find(&data).Error; err != nil {
		return &data, err
	}
	return &data, nil
}

func (r *agingutangpiutangdetail) FindByCriteria(ctx *abstraction.Context, filter *model.AgingUtangPiutangDetailFilterModel) (data *model.AgingUtangPiutangDetailEntityModel, err error) {
	conn := r.CheckTrx(ctx)
	query := conn.Model(&model.AgingUtangPiutangDetailEntityModel{})
	query = r.Filter(ctx, query, *filter)
	if err = query.First(&data).Error; err != nil {
		return
	}
	if data.ID == 0 {
		err = errors.New("Data Not Found")
	}
	return
}
