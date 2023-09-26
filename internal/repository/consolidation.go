package repository

import (
	"errors"
	"strconv"
	"worker/internal/abstraction"
	"worker/internal/model"

	"gorm.io/gorm"
)

type Consolidation interface {
	Find(ctx *abstraction.Context, m *model.ConsolidationFilterModel) (*[]model.ConsolidationEntityModel, error)
	FindByCriteria(ctx *abstraction.Context, filter *model.ConsolidationFilterModel) (data *model.ConsolidationEntityModel, err error)
	FindDetailByCode(ctx *abstraction.Context, consolidationID *int, code *string) (*[]model.ConsolidationDetailEntityModel, error)
	FindByID(ctx *abstraction.Context, id *int) (*model.ConsolidationEntityModel, error)
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
	if err = query.Preload("Company").Preload("ConsolidationBridge.Company").First(&data).Error; err != nil {
		return
	}
	if data.ID == 0 {
		err = errors.New("data not found")
	}
	return
}

func (r *consolidation) FindDetailByCode(ctx *abstraction.Context, consolidationID *int, code *string) (*[]model.ConsolidationDetailEntityModel, error) {
	tmpStr := *code
	if _, err := strconv.Atoi(*code); err == nil {
		tmpStr += "%"
	}
	conn := r.CheckTrx(ctx)
	var data *[]model.ConsolidationDetailEntityModel
	query := conn.Model(&model.ConsolidationDetailEntityModel{})
	query = query.Where("consolidation_id = ?", consolidationID).Where("code LIKE ?", tmpStr)
	query = query.Order("id asc")
	if err := query.Find(&data).Error; err != nil {
		return nil, err
	}
	return data, nil
}

func (r *consolidation) FindByID(ctx *abstraction.Context, id *int) (*model.ConsolidationEntityModel, error) {
	conn := r.CheckTrx(ctx)
	var data *model.ConsolidationEntityModel
	query := conn.Model(&model.ConsolidationEntityModel{})
	query = query.Where("id = ?", id)
	if err := query.First(&data).Error; err != nil {
		return nil, err
	}
	return data, nil
}
