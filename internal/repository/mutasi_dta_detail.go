package repository

import (
	"errors"
	"fmt"
	"worker-validation/internal/abstraction"
	"worker-validation/internal/model"

	"gorm.io/gorm"
)

type MutasiDtaDetail interface {
	Find(ctx *abstraction.Context, m *model.MutasiDtaDetailFilterModel) (*[]model.MutasiDtaDetailEntityModel, error)
	FindWithCode(ctx *abstraction.Context, code *string) (*[]model.MutasiDtaDetailEntityModel, error)
	FindByCriteria(ctx *abstraction.Context, filter *model.MutasiDtaDetailFilterModel) (data *model.MutasiDtaDetailEntityModel, err error)
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

func (r *mutasidtadetail) FindByCriteria(ctx *abstraction.Context, filter *model.MutasiDtaDetailFilterModel) (data *model.MutasiDtaDetailEntityModel, err error) {
	conn := r.CheckTrx(ctx)
	query := conn.Model(&model.MutasiDtaDetailEntityModel{})
	query = r.Filter(ctx, query, *filter)
	if err = query.First(&data).Error; err != nil {
		return
	}
	if data.ID == 0 {
		err = errors.New("Data Not Found")
	}
	return
}
