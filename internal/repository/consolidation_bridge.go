package repository

import (
	"errors"
	"fmt"
	"strings"
	"worker-consol/internal/abstraction"
	"worker-consol/internal/model"

	"gorm.io/gorm"
)

type ConsolidationBridge interface {
	Find(ctx *abstraction.Context, m *model.ConsolidationBridgeFilterModel) (*[]model.ConsolidationBridgeEntityModel, error)
	FindByCriteria(ctx *abstraction.Context, filter *model.ConsolidationBridgeFilterModel) (data *model.ConsolidationBridgeEntityModel, err error)
	Create(ctx *abstraction.Context, payload *model.ConsolidationBridgeEntityModel) (*model.ConsolidationBridgeEntityModel, error)
	FindListConsolBridge(ctx *abstraction.Context, m *model.ConsolidationBridgeFilterModel) (string, error)
	DeleteByConsolID(ctx *abstraction.Context, id *int) error
	FindTBByConsolBridgeID(ctx *abstraction.Context, listid *string) (*[]model.TrialBalanceEntityModel, error)
}

type consolidationbridge struct {
	abstraction.Repository
}

func NewConsolidationBridge(db *gorm.DB) *consolidationbridge {
	return &consolidationbridge{
		abstraction.Repository{
			Db: db,
		},
	}
}

func (r *consolidationbridge) Find(ctx *abstraction.Context, m *model.ConsolidationBridgeFilterModel) (*[]model.ConsolidationBridgeEntityModel, error) {
	conn := r.CheckTrx(ctx)

	var datas []model.ConsolidationBridgeEntityModel

	query := conn.Model(&model.ConsolidationBridgeEntityModel{})
	//filter
	query = r.Filter(ctx, query, *m)

	if err := query.Find(&datas).Error; err != nil {
		return &datas, err
	}

	return &datas, nil
}

func (r *consolidationbridge) FindByCriteria(ctx *abstraction.Context, filter *model.ConsolidationBridgeFilterModel) (data *model.ConsolidationBridgeEntityModel, err error) {
	conn := r.CheckTrx(ctx)
	query := conn.Model(&model.ConsolidationBridgeEntityModel{})
	query = r.Filter(ctx, query, *filter)
	if err = query.First(&data).Error; err != nil {
		return
	}
	if data.ID == 0 {
		err = errors.New("Data Not Found")
	}
	return
}

func (r *consolidationbridge) Create(ctx *abstraction.Context, payload *model.ConsolidationBridgeEntityModel) (*model.ConsolidationBridgeEntityModel, error) {
	conn := r.CheckTrx(ctx)
	if err := conn.Create(&payload).Error; err != nil {
		return nil, err
	}

	if err := conn.Model(&payload).First(&payload).Error; err != nil {
		return nil, err
	}

	return payload, nil
}

func (r *consolidationbridge) FindListConsolBridge(ctx *abstraction.Context, m *model.ConsolidationBridgeFilterModel) (string, error) {
	conn := r.CheckTrx(ctx)

	var datas []model.ConsolidationBridgeEntityModel

	query := conn.Model(&model.ConsolidationBridgeEntityModel{})
	//filter
	query = r.Filter(ctx, query, *m)

	if err := query.Find(&datas).Error; err != nil {
		return "", err
	}
	tmpList := []string{}
	for _, vData := range datas {
		tmpList = append(tmpList, fmt.Sprintf("%d", vData.ID))
	}
	listConsolBridge := strings.Join(tmpList, ",")
	return listConsolBridge, nil
}

func (r *consolidationbridge) DeleteByConsolID(ctx *abstraction.Context, id *int) error {
	conn := r.CheckTrx(ctx)
	query := conn.Model(&model.ConsolidationBridgeEntityModel{})
	if err := query.Where("consolidation_id = ?", id).Delete(&model.ConsolidationBridgeEntityModel{}).Error; err != nil {
		return err
	}
	return nil
}

func (r *consolidationbridge) FindTBByConsolBridgeID(ctx *abstraction.Context, listid *string) (*[]model.TrialBalanceEntityModel, error) {
	conn := r.CheckTrx(ctx)
	tmp := []model.ConsolidationBridgeEntityModel{}
	query := conn.Model(&model.ConsolidationBridgeEntityModel{})
	if err := query.Where(fmt.Sprintf("id IN (%s)", *listid)).Find(&tmp).Error; err != nil {
		return nil, err
	}
	var data []model.TrialBalanceEntityModel
	for _, vData := range tmp {
		tmpData := model.TrialBalanceEntityModel{}
		if vData.ConsolidationVersions == 0 {
			if err := conn.Model(&model.TrialBalanceEntityModel{}).Where("company_id = ?", vData.CompanyID).Where("versions = ?", vData.Versions).Where("period = ?", vData.Period).First(&tmpData).Error; err != nil {
				return nil, err
			}
		} else {
			consolData := model.ConsolidationEntityModel{}
			if err := conn.Model(&model.ConsolidationEntityModel{}).Where("company_id = ?", vData.CompanyID).Where("versions = ?", vData.Versions).Where("period = ?", vData.Period).Where("consolidation_versions = ?", vData.ConsolidationVersions).First(&consolData).Error; err != nil {
				return nil, err
			}

			if err := conn.Model(&model.TrialBalanceEntityModel{}).Where("company_id = ?", consolData.CompanyID).Where("versions = ?", consolData.Versions).Where("period = ?", consolData.Period).First(&tmpData).Error; err != nil {
				return nil, err
			}
		}
		data = append(data, tmpData)
	}

	return &data, nil
}
