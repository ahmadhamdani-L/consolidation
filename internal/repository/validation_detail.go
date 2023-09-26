package repository

import (
	"mcash-finance-console-core/internal/abstraction"
	"mcash-finance-console-core/internal/model"
	"mcash-finance-console-core/pkg/constant"

	"gorm.io/gorm"
)

type ValidationDetail interface {
	Find(ctx *abstraction.Context, m *model.ValidationDetailFilterModel) (*[]model.ValidationDetailEntityModel, error)
	FindByID(ctx *abstraction.Context, id *int) (*model.ValidationDetailEntityModel, error)
	GetStatus(ctx *abstraction.Context, m *model.ValidationDetailFilterModel) (statusDesc int, status int, err error)
	MakeSure(ctx *abstraction.Context, m *model.ValidationDetailFilterModel) (*model.ValidationDetailEntityModel, error)
	CheckExist(ctx *abstraction.Context, m *model.ValidationDetailFilterModel) (bool, error)
	UpdateByCriteria(ctx *abstraction.Context, m *model.ValidationDetailFilterModel, data *model.ValidationDetailEntityModel) (*model.ValidationDetailEntityModel, error)
}

type validationdetail struct {
	abstraction.Repository
}

func NewValidationDetail(db *gorm.DB) *validationdetail {
	return &validationdetail{
		abstraction.Repository{
			Db: db,
		},
	}
}

func (r *validationdetail) Find(ctx *abstraction.Context, m *model.ValidationDetailFilterModel) (*[]model.ValidationDetailEntityModel, error) {
	conn := r.CheckTrx(ctx)

	var datas []model.ValidationDetailEntityModel

	query := conn.Model(&model.ValidationDetailEntityModel{})
	//filter
	query = r.Filter(ctx, query, *m)

	if err := query.Preload("UserValidate").Find(&datas).WithContext(ctx.Request().Context()).Error; err != nil {
		return &datas, err
	}

	for i, v := range datas {
		datas[i].UserValidateBy = v.UserValidate.Name
	}

	return &datas, nil
}

func (r *validationdetail) FindByID(ctx *abstraction.Context, id *int) (*model.ValidationDetailEntityModel, error) {
	conn := r.CheckTrx(ctx)

	var data model.ValidationDetailEntityModel
	if err := conn.Where("id = ?", &id).First(&data).WithContext(ctx.Request().Context()).Error; err != nil {
		return &data, err
	}

	return &data, nil
}

func (r *validationdetail) GetStatus(ctx *abstraction.Context, m *model.ValidationDetailFilterModel) (statusDesc int, status int, err error) {
	conn := r.CheckTrx(ctx)

	var datas []model.ValidationDetailEntityModel

	query := conn.Model(&model.ValidationDetailEntityModel{})
	query = r.Filter(ctx, query, *m)

	query = query.Group("status")

	if err := query.Find(&datas).WithContext(ctx.Request().Context()).Error; err != nil {
		return 0, 0, err
	}

	for _, v := range datas {
		if v.Status == 0 {
			return 0, 0, nil
		}
	}
	return 1, 1, nil
}

func (r *validationdetail) Create(ctx *abstraction.Context, e *model.ValidationDetailEntityModel) (*model.ValidationDetailEntityModel, error) {
	conn := r.CheckTrx(ctx)

	err := conn.Create(&e).Error
	if err != nil {
		return nil, err
	}
	err = conn.Model(&e).First(&e).Error
	if err != nil {
		return nil, err
	}
	return e, nil
}

func (r *validationdetail) MakeSure(ctx *abstraction.Context, m *model.ValidationDetailFilterModel) (*model.ValidationDetailEntityModel, error) {
	conn := r.CheckTrx(ctx)
	data := model.ValidationDetailEntityModel{}
	query := conn.Model(&data)
	tableName := model.ValidationDetailEntityModel{}.TableName()
	query = r.Filter(ctx, query, *m)
	query = r.AllowedCompany(ctx, query, tableName)
	query = query.Where("status = ?", constant.VALIDATION_STATUS_NOT_BALANCE)

	if err := query.Find(&data).Error; err != nil {
		return nil, err
	}

	return &data, nil
}

func (r *validationdetail) CheckExist(ctx *abstraction.Context, m *model.ValidationDetailFilterModel) (bool, error) {
	conn := r.CheckTrx(ctx)
	query := conn.Model(&model.ValidationDetailEntityModel{})
	tableName := model.ValidationDetailEntityModel{}.TableName()
	query = r.Filter(ctx, query, *m)
	query = r.AllowedCompany(ctx, query, tableName)
	var tmp int64

	if err := query.Count(&tmp).Error; err != nil || tmp == 0 {
		return false, err
	}

	return true, nil
}

func (r *validationdetail) UpdateByCriteria(ctx *abstraction.Context, m *model.ValidationDetailFilterModel, data *model.ValidationDetailEntityModel) (*model.ValidationDetailEntityModel, error) {
	conn := r.CheckTrx(ctx)
	query := conn.Model(&model.ValidationDetailEntityModel{})
	tableName := model.ValidationDetailEntityModel{}.TableName()
	query = r.Filter(ctx, query, *m)
	query = r.AllowedCompany(ctx, query, tableName)
	var validationData model.ValidationDetailEntityModel

	if err := query.First(&validationData).Error; err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	}

	if validationData.ID == 0 {
		return nil, nil
	}

	if err := query.Updates(data).Error; err != nil {
		return nil, err
	}

	if err := query.First(&validationData).Error; err != nil {
		return nil, err
	}

	queryTB := conn.Table(model.TrialBalanceEntityModel{}.TableName())
	queryTB = queryTB.Where("company_id = ? AND period = ? AND versions = ? AND status = ?", m.CompanyID, m.Period, m.Versions, constant.MODUL_STATUS_DRAFT)

	if err := queryTB.Update("validation_note", constant.VALIDATION_NOTE_NOT_BALANCE).Error; err != nil {
		return nil, err
	}

	return &validationData, nil
}
