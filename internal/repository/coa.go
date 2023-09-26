package repository

import (
	"fmt"
	"mcash-finance-console-core/internal/abstraction"
	"mcash-finance-console-core/internal/model"
	"strings"

	"gorm.io/gorm"
)

type Coa interface {
	Find(ctx *abstraction.Context, m *model.CoaFilterModel, p *abstraction.Pagination) (*[]model.CoaEntityModel, *abstraction.PaginationInfo, error)
	FindByID(ctx *abstraction.Context, id *int) (*model.CoaEntityModel, error)
	Create(ctx *abstraction.Context, e *model.CoaEntityModel) (*model.CoaEntityModel, error)
	Update(ctx *abstraction.Context, id *int, e *model.CoaEntityModel) (*model.CoaEntityModel, error)
	Delete(ctx *abstraction.Context, id *int, e *model.CoaEntityModel) (*model.CoaEntityModel, error)
	FindWithCode(ctx *abstraction.Context, code *string) (*[]model.CoaEntityModel, error)
	FindWithCodes(ctx *abstraction.Context, code *string) (*model.CoaEntityModel, error)
	Export(ctx *abstraction.Context) (*[]model.CoaGroupEntityModel, error)
	CountDash(ctx *abstraction.Context) (int, error)
}

type coa struct {
	abstraction.Repository
}

func NewCoa(db *gorm.DB) *coa {
	return &coa{
		abstraction.Repository{
			Db: db,
		},
	}
}

func (r *coa) Find(ctx *abstraction.Context, m *model.CoaFilterModel, p *abstraction.Pagination) (*[]model.CoaEntityModel, *abstraction.PaginationInfo, error) {
	conn := r.CheckTrx(ctx)

	var datas []model.CoaEntityModel
	var info abstraction.PaginationInfo

	query := conn.Model(&model.CoaEntityModel{})
	//filter
	tableName := model.CoaEntityModel{}.TableName()
	query = r.FilterTable(ctx, query, *m, model.CoaEntityModel{}.TableName())
	queryUser := conn.Model(&model.UserEntityModel{}).Select("id")
	query = r.FilterUser(ctx, query, queryUser, m.Filter, tableName)
	if len(*m.ArrCoaGroupID) > 0 {
		listCoaGroup := strings.Trim(strings.Join(strings.Split(fmt.Sprint(*m.ArrCoaGroupID), " "), ","), "[]")
		query = query.Where("m_coa.coa_type_id IN (?)", conn.Model(&model.CoaTypeEntityModel{}).Select("id").Where(fmt.Sprintf("coa_group_id IN (%s)", listCoaGroup)))
	}

	if m.Search != nil {
		query = query.Where("(m_coa.name ILIKE ? OR m_coa.code ILIKE ?)", "%"+*m.Search+"%", "%"+*m.Search+"%")
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
	if p.SortBy != nil && *p.SortBy == "coa_type" {
		sortBy := "\"m_coa_type\".name"
		p.SortBy = &sortBy
	} else if p.SortBy != nil && *p.SortBy == "coa_group" {
		sortBy := "\"m_coa_group\".name"
		p.SortBy = &sortBy
	} else {
		sortBy := "m_coa." + *p.SortBy
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

	if err := query.Joins("JOIN m_coa_type ON m_coa_type.id = m_coa.coa_type_id").Joins("JOIN m_coa_group ON m_coa_type.coa_group_id = m_coa_group.id").Preload("CoaType.CoaGroup").Preload("UserCreated").Preload("UserModified").Find(&datas).WithContext(ctx.Request().Context()).Error; err != nil {
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

func (r *coa) FindByID(ctx *abstraction.Context, id *int) (*model.CoaEntityModel, error) {
	conn := r.CheckTrx(ctx)

	var data model.CoaEntityModel
	if err := conn.Where("id = ?", &id).Preload("CoaType.CoaGroup").Preload("UserCreated").Preload("UserModified").First(&data).WithContext(ctx.Request().Context()).Error; err != nil {
		return &data, err
	}
	return &data, nil
}

func (r *coa) FindWithCode(ctx *abstraction.Context, code *string) (*[]model.CoaEntityModel, error) {
	conn := r.CheckTrx(ctx)

	var data []model.CoaEntityModel
	if err := conn.Where("code LIKE ?", *code+"%").Order("code asc").Find(&data).WithContext(ctx.Request().Context()).Error; err != nil {
		return &data, err
	}
	return &data, nil
}
func (r *coa) FindWithCodes(ctx *abstraction.Context, code *string) (*model.CoaEntityModel, error) {
	conn := r.CheckTrx(ctx)

	var data model.CoaEntityModel
	if err := conn.Where("code LIKE ?", *code+"%").Order("code asc").Find(&data).WithContext(ctx.Request().Context()).Error; err != nil {
		return &data, err
	}
	return &data, nil
}
func (r *coa) Create(ctx *abstraction.Context, e *model.CoaEntityModel) (*model.CoaEntityModel, error) {
	conn := r.CheckTrx(ctx)

	if err := conn.Create(e).WithContext(ctx.Request().Context()).Error; err != nil {
		return nil, err
	}
	if err := conn.Model(e).Preload("UserCreated").Preload("UserModified").First(e).WithContext(ctx.Request().Context()).Error; err != nil {
		return nil, err
	}

	return e, nil
}

func (r *coa) Update(ctx *abstraction.Context, id *int, e *model.CoaEntityModel) (*model.CoaEntityModel, error) {
	conn := r.CheckTrx(ctx)

	if err := conn.Model(e).Where("id = ?", &id).Updates(e).WithContext(ctx.Request().Context()).Error; err != nil {
		return nil, err
	}
	if err := conn.Model(e).Where("id = ?", &id).Preload("UserCreated").Preload("UserModified").First(e).WithContext(ctx.Request().Context()).Error; err != nil {
		return nil, err
	}
	return e, nil

}

func (r *coa) Delete(ctx *abstraction.Context, id *int, e *model.CoaEntityModel) (*model.CoaEntityModel, error) {
	conn := r.CheckTrx(ctx)

	if err := conn.Model(e).Where("id =?", id).Delete(e).WithContext(ctx.Request().Context()).Error; err != nil {
		return nil, err
	}
	return e, nil
}

func (r *coa) Export(ctx *abstraction.Context) (*[]model.CoaGroupEntityModel, error) {
	conn := r.CheckTrx(ctx)
	var data []model.CoaGroupEntityModel
	if err := conn.Model(&model.CoaGroupEntityModel{}).Preload("CoaType.Coa").Order("code ASC").Find(&data).Error; err != nil {
		return nil, err
	}
	return &data, nil
}

func (r *coa) CountDash(ctx *abstraction.Context) (int, error) {
	conn := r.CheckTrx(ctx)
	var TotalDash int
	if err := conn.Model(&model.CoaEntityModel{}).Select("LENGTH(name) - LENGTH(REPLACE(name, '-', '')) as TotalDash").Order("TotalDash DESC").Limit(1).Scan(&TotalDash).Error; err != nil {
		return 0, err
	}
	return TotalDash, nil
}
