package repository

import (
	"fmt"
	"worker-consol/internal/abstraction"
	"worker-consol/internal/model"

	"gorm.io/gorm"
)

type MutasiDtaDetail interface {
	Find(ctx *abstraction.Context, m *model.MutasiDtaDetailFilterModel) (*[]model.MutasiDtaDetailEntityModel, error)
	FindWithCode(ctx *abstraction.Context, code *string) (*[]model.MutasiDtaDetailEntityModel, error)
}

type mutasidtadetail struct {
	abstraction.Repository
}

func NewMutasiDtaDetail(db *gorm.DB) *mutasidtadetail {
	return &mutasidtadetail{
		abstraction.Repository{
			Db: db,
		},
	}
}

func (r *mutasidtadetail) Find(ctx *abstraction.Context, m *model.MutasiDtaDetailFilterModel) (*[]model.MutasiDtaDetailEntityModel, error) {
	conn := r.CheckTrx(ctx)

	var datas []model.MutasiDtaDetailEntityModel

	query := conn.Model(&model.MutasiDtaDetailEntityModel{})
	//filter
	query = r.Filter(ctx, query, *m)
	query = query.Where("formatter_bridges_id IN (?)", conn.Model(&model.FormatterBridgesEntityModel{}).Select("id").Where("source = ?", "MUTASI-DTA").Where("trx_ref_id = ?", m.MutasiDtaID))
	if err := query.Find(&datas).Error; err != nil {
		return &datas, err
	}

	return &datas, nil
}

func (r *mutasidtadetail) FindWithCode(ctx *abstraction.Context, code *string) (*[]model.MutasiDtaDetailEntityModel, error) {
	conn := r.CheckTrx(ctx)

	var data []model.MutasiDtaDetailEntityModel
	tmp := fmt.Sprintf("%s", *code)
	if err := conn.Where("code ILIKE ?", tmp+"%").Find(&data).Error; err != nil {
		return &data, err
	}
	return &data, nil
}
