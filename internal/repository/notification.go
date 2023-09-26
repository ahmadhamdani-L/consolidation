package repository

import (
	"fmt"
	"mcash-finance-console-core/internal/abstraction"
	"mcash-finance-console-core/internal/model"

	"gorm.io/gorm"
)

type Notification interface {
	Find(ctx *abstraction.Context, m *model.NotificationFilterModel, p *abstraction.Pagination) (*[]model.NotificationEntityModel, *abstraction.PaginationInfo, error)
	FindByID(ctx *abstraction.Context, id *int) (*model.NotificationEntityModel, error)
	Update(ctx *abstraction.Context, id *int, e *model.NotificationEntityModel) (*model.NotificationEntityModel, error)
	Count(ctx *abstraction.Context, m *model.NotificationFilterModel) (int64, int64, int64, error)
}

type notification struct {
	abstraction.Repository
}

func NewNotification(db *gorm.DB) *notification {
	return &notification{
		abstraction.Repository{
			Db: db,
		},
	}
}

func (r *notification) Find(ctx *abstraction.Context, m *model.NotificationFilterModel, p *abstraction.Pagination) (*[]model.NotificationEntityModel, *abstraction.PaginationInfo, error) {
	conn := r.CheckTrx(ctx)

	var datas []model.NotificationEntityModel
	var info abstraction.PaginationInfo

	query := conn.Model(&model.NotificationEntityModel{})
	//filter
	query = r.Filter(ctx, query, *m)
	queryUser := conn.Model(&model.UserEntityModel{}).Select("id")
	query = r.FilterUser(ctx, query, queryUser, m.Filter, "")
	//sort
	if p.Sort == nil {
		sort := "desc"
		p.Sort = &sort
	}
	if p.SortBy == nil {
		sortBy := "created_at"
		p.SortBy = &sortBy
	}

	sort := fmt.Sprintf("%s %s", *p.SortBy, *p.Sort)
	query = query.Order(sort)

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

	if err := query.Find(&datas).WithContext(ctx.Request().Context()).Error; err != nil {
		return &datas, &info, err
	}

	for i, v := range datas {
		datas[i].UserCreatedString = v.UserCreated.Name
		datas[i].UserModifiedString = &v.UserModified.Name
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

func (r *notification) FindByID(ctx *abstraction.Context, id *int) (*model.NotificationEntityModel, error) {
	conn := r.CheckTrx(ctx)

	var data model.NotificationEntityModel
	if err := conn.Where("id = ?", &id).First(&data).WithContext(ctx.Request().Context()).Error; err != nil {
		return &data, err
	}
	return &data, nil
}

func (r *notification) Update(ctx *abstraction.Context, id *int, e *model.NotificationEntityModel) (*model.NotificationEntityModel, error) {
	conn := r.CheckTrx(ctx)

	if err := conn.Model(e).Where("id = ?", &id).Updates(e).WithContext(ctx.Request().Context()).Error; err != nil {
		return nil, err
	}
	if err := conn.Model(e).Where("id = ?", &id).First(e).WithContext(ctx.Request().Context()).Error; err != nil {
		return nil, err
	}
	return e, nil

}

func (r *notification) Count(ctx *abstraction.Context, m *model.NotificationFilterModel) (int64, int64, int64, error) {
	conn := r.CheckTrx(ctx)

	query := conn.Model(&model.NotificationEntityModel{}).Where("created_by = ?", ctx.Auth.ID)
	queryRead := conn.Model(&model.NotificationEntityModel{}).Where("is_opened = ?", true).Where("created_by = ?", ctx.Auth.ID)
	queryUnread := conn.Model(&model.NotificationEntityModel{}).Where("is_opened != ?", true).Where("created_by = ?", ctx.Auth.ID)
	queryUser := conn.Model(&model.UserEntityModel{}).Select("id")
	query = r.FilterUser(ctx, query, queryUser, m.Filter, "")
	queryRead = r.FilterUser(ctx, queryRead, queryUser, m.Filter, "")
	queryUnread = r.FilterUser(ctx, queryUnread, queryUser, m.Filter, "")

	var totalData int64
	var totalRead int64
	var totalUnread int64
	if err := query.Count(&totalData).Error; err != nil {
		return 0, 0, 0, err
	}
	if err := queryRead.Count(&totalRead).Error; err != nil {
		return 0, 0, 0, err
	}
	if err := queryUnread.Count(&totalUnread).Error; err != nil {
		return 0, 0, 0, err
	}

	return totalData, totalRead, totalUnread, nil
}
