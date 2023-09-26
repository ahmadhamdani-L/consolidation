package repository

import (
	"fmt"
	"worker/internal/abstraction"
	"worker/internal/model"

	"gorm.io/gorm"
)

type InvestasiTbkDetail interface {
	Find(ctx *abstraction.Context, m *model.InvestasiTbkDetailFilterModel) (*[]model.InvestasiTbkDetailEntityModel, error)
	FindWithCode(ctx *abstraction.Context, code *string) (*[]model.InvestasiTbkDetailEntityModel, error)
}

type investasitbkdetail struct {
	abstraction.Repository
}

func NewInvestasiTbkDetail(db *gorm.DB) *investasitbkdetail {
	return &investasitbkdetail{
		abstraction.Repository{
			Db: db,
		},
	}
}

func (r *investasitbkdetail) Find(ctx *abstraction.Context, m *model.InvestasiTbkDetailFilterModel) (*[]model.InvestasiTbkDetailEntityModel, error) {
	conn := r.CheckTrx(ctx)

	var datas []model.InvestasiTbkDetailEntityModel

	query := conn.Model(&model.InvestasiTbkDetailEntityModel{})
	//filter
	query = r.Filter(ctx, query, *m)
	// query = query.Where("formatter_bridges_id = ?", m.FormatterBridgesID)
	query = query.Where("formatter_bridges_id IN (?)", conn.Model(&model.FormatterBridgesEntityModel{}).Select("id").Where("source = ?", "INVESTASI-TBK").Where("trx_ref_id = ?", m.InvestasiTbkID))
	if err := query.Find(&datas).Error; err != nil {
		return &datas, err
	}

	return &datas, nil
}

func (r *investasitbkdetail) FindWithCode(ctx *abstraction.Context, code *string) (*[]model.InvestasiTbkDetailEntityModel, error) {
	conn := r.CheckTrx(ctx)

	var data []model.InvestasiTbkDetailEntityModel
	tmp := fmt.Sprintf("%s", *code)
	if err := conn.Where("code ILIKE ?", tmp+"%").Find(&data).Error; err != nil {
		return &data, err
	}
	return &data, nil
}
