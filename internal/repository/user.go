package repository

import (
	"fmt"
	"mcash-finance-console-core/internal/abstraction"
	"mcash-finance-console-core/internal/model"

	"gorm.io/gorm"
)

type User interface {
	Find(ctx *abstraction.Context, m *model.UserFilterModel, p *abstraction.Pagination) (*[]model.UserEntityModel, *abstraction.PaginationInfo, error)
	FindByID(ctx *abstraction.Context, id *int) (*model.UserEntityModel, error)
	FindByUsername(ctx *abstraction.Context, username *string) (*model.UserEntityModel, error)
	Create(ctx *abstraction.Context, data *model.UserEntityModel) (*model.UserEntityModel, error)
	checkTrx(ctx *abstraction.Context) *gorm.DB
	Update(ctx *abstraction.Context, id *int, e *model.UserEntityModel) (*model.UserEntityModel, error)
	Delete(ctx *abstraction.Context, id *int, e *model.UserEntityModel) (*model.UserEntityModel, error)
	CountByCriteria(ctx *abstraction.Context, e *model.UserFilterModel) (int64, error)
	FindByEmail(ctx *abstraction.Context, email *string) (*model.UserEntityModel, error)
	UpdateUserPassword(ctx *abstraction.Context, token *string, hashPwd string) error
}

type user struct {
	abstraction.Repository
}

func NewUser(db *gorm.DB) *user {
	return &user{
		abstraction.Repository{
			Db: db,
		},
	}
}

func (r *user) FindByUsername(ctx *abstraction.Context, username *string) (*model.UserEntityModel, error) {
	conn := r.checkTrx(ctx)

	var data model.UserEntityModel
	err := conn.Where("username = ? OR email = ?", username, username).First(&data).WithContext(ctx.Request().Context()).Error
	if err != nil {
		return nil, err
	}
	return &data, nil
}

func (r *user) Create(ctx *abstraction.Context, e *model.UserEntityModel) (*model.UserEntityModel, error) {
	conn := r.checkTrx(ctx)

	err := conn.Create(&e).WithContext(ctx.Request().Context()).Error
	if err != nil {
		return nil, err
	}
	err = conn.Model(&e).First(&e).WithContext(ctx.Request().Context()).Error
	if err != nil {
		return nil, err
	}
	return e, nil
}

func (r *user) checkTrx(ctx *abstraction.Context) *gorm.DB {
	if ctx.Trx != nil {
		return ctx.Trx.Db
	}
	return r.Db
}

func (r *user) Find(ctx *abstraction.Context, m *model.UserFilterModel, p *abstraction.Pagination) (*[]model.UserEntityModel, *abstraction.PaginationInfo, error) {
	conn := r.CheckTrx(ctx)

	var datas []model.UserEntityModel
	var info abstraction.PaginationInfo

	if m.IsActive != nil && !*m.IsActive {
		tmpFalse := false
		m.IsActive = &tmpFalse
	}

	query := conn.Model(&model.UserEntityModel{})
	//filter
	query = r.FilterTable(ctx, query, *m, model.UserEntityModel{}.TableName())

	if m.Search != nil {
		query = query.Where("users.email ILIKE ? OR users.username ILIKE ? OR users.name ILIKE ?", "%"+*m.Search+"%", "%"+*m.Search+"%", "%"+*m.Search+"%")
	}

	//sort
	if p.Sort == nil {
		sort := "desc"
		p.Sort = &sort
	}
	if p.SortBy == nil {
		sortBy := "created_at"
		p.SortBy = &sortBy
	}

	tmpSortBy := p.SortBy
	if p.SortBy != nil && *p.SortBy == "company" {
		sortBy := "\"Company\".name"
		p.SortBy = &sortBy
	}
	if p.SortBy != nil && *p.SortBy == "role" {
		sortBy := "\"Role\".name"
		p.SortBy = &sortBy
	}

	sort := fmt.Sprintf("%s %s", *p.SortBy, *p.Sort)
	query = query.Order(sort)
	p.SortBy = tmpSortBy

	//pagination
	if p.Page == nil {
		page := 1
		p.Page = &page
	}
	if p.PageSize == nil {
		pageSize := 10
		p.PageSize = &pageSize
	}
	info = abstraction.PaginationInfo{
		Pagination: p,
	}
	limit := *p.PageSize
	offset := limit * (*p.Page - 1)
	var totalData int64
	query = query.Count(&totalData).Limit(limit).Offset(offset)

	if err := query.Joins("Company").Joins("Role").Find(&datas).WithContext(ctx.Request().Context()).Error; err != nil {
		return &datas, &info, err
	}

	info.Count = int(totalData)
	info.MoreRecords = false
	if len(datas) > *p.PageSize {
		info.MoreRecords = true
		// info.Count -= 1
		// datas = datas[:len(datas)-1]
	}

	return &datas, &info, nil
}

func (r *user) FindByID(ctx *abstraction.Context, id *int) (*model.UserEntityModel, error) {
	conn := r.CheckTrx(ctx)

	var data model.UserEntityModel
	if err := conn.Where("id = ?", &id).Preload("Company").Preload("UserCreated").First(&data).WithContext(ctx.Request().Context()).Error; err != nil {
		return &data, err
	}

	return &data, nil
}

func (r *user) Update(ctx *abstraction.Context, id *int, e *model.UserEntityModel) (*model.UserEntityModel, error) {
	conn := r.checkTrx(ctx)

	if err := conn.Model(&e).Where("id = ?", &id).Updates(&e).WithContext(ctx.Request().Context()).Error; err != nil {
		return nil, err
	}
	if err := conn.Model(&model.UserEntityModel{}).Where("id = ?", &id).First(&e).WithContext(ctx.Request().Context()).Error; err != nil {
		return nil, err
	}
	return e, nil

}

func (r *user) Delete(ctx *abstraction.Context, id *int, e *model.UserEntityModel) (*model.UserEntityModel, error) {
	conn := r.CheckTrx(ctx)

	if err := conn.Model(e).Where("id =?", id).Delete(e).WithContext(ctx.Request().Context()).Error; err != nil {
		return nil, err
	}
	return e, nil
}

func (r *user) CountByCriteria(ctx *abstraction.Context, e *model.UserFilterModel) (int64, error) {
	conn := r.CheckTrx(ctx)
	var total int64
	query := conn.Model(&model.UserEntityModel{})
	query = r.Filter(ctx, query, *e)

	if err := query.Count(&total).WithContext(ctx.Request().Context()).Error; err != nil {
		return 0, err
	}
	return total, nil
}

func (r *user) FindByEmail(ctx *abstraction.Context, email *string) (*model.UserEntityModel, error) {
	conn := r.CheckTrx(ctx)

	var data model.UserEntityModel
	if err := conn.Where("email = ?", &email).Preload("Company").First(&data).WithContext(ctx.Request().Context()).Error; err != nil {
		return nil, err
	}
	// data.UserCreatedString = data.UserCreated.Name
	// data.UserModifiedString = &data.UserModified.Name
	return &data, nil
}

func (r *user) UpdateUserPassword(ctx *abstraction.Context, token *string, hashPwd string) error {
	conn := r.checkTrx(ctx)

	var data model.UserEntityModel

	err := conn.Where("id = ?", ctx.Auth.ID).First(&data).
		WithContext(ctx.Request().Context()).Error
	if err != nil {
		return err
	}
	err = conn.Model(&data).UpdateColumn("password", hashPwd).
		WithContext(ctx.Request().Context()).Error
	if err != nil {
		return err
	}
	return nil
}
