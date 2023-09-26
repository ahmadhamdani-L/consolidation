package repository

import (
	"fmt"
	"mcash-finance-console-core/internal/abstraction"
	"mcash-finance-console-core/internal/model"

	"gorm.io/gorm"
)

type ImportedWorksheetDetail interface {
	Create(ctx *abstraction.Context, e *model.ImportedWorksheetDetailEntityModel) (*model.ImportedWorksheetDetailEntityModel, error)
	GetCountStatus(ctx *abstraction.Context, e *model.ImportedWorksheetDetailEntityModel) (*[]model.ImportedWorksheetDetailEntityModel, error)
	Find(ctx *abstraction.Context, m *model.ImportedWorksheetDetailFilterModel, p *abstraction.Pagination) (*[]model.ImportedWorksheetDetailEntityModel, *abstraction.PaginationInfo, error)
	FindByID(ctx *abstraction.Context, id *int) (*model.ImportedWorksheetDetailEntityModel, error)
}

type importedworksheetdetail struct {
	abstraction.Repository
}

func NewImportedWorksheetDetail(db *gorm.DB) *importedworksheetdetail {
	return &importedworksheetdetail{
		abstraction.Repository{
			Db: db,
		},
	}
}

func (r *importedworksheetdetail) Create(ctx *abstraction.Context, e *model.ImportedWorksheetDetailEntityModel) (*model.ImportedWorksheetDetailEntityModel, error) {
	conn := r.CheckTrx(ctx)

	if err := conn.Create(e).Error; err != nil {
		return nil, err
	}
	if err := conn.Model(e).First(e).Error; err != nil {
		return nil, err
	}

	return e, nil
}

func (r *importedworksheetdetail) GetCountStatus(ctx *abstraction.Context, e *model.ImportedWorksheetDetailEntityModel) (*[]model.ImportedWorksheetDetailEntityModel, error) {
	conn := r.CheckTrx(ctx)

	var data []model.ImportedWorksheetDetailEntityModel

	if err := conn.Where("imported_worksheet_id =? AND status =? ", e.ImportedWorksheetID, e.Status).Find(&data).Error; err != nil {
		return &data, err
	}
	return &data, nil
}
func (r *importedworksheetdetail) Find(ctx *abstraction.Context, m *model.ImportedWorksheetDetailFilterModel, p *abstraction.Pagination) (*[]model.ImportedWorksheetDetailEntityModel, *abstraction.PaginationInfo, error) {
	conn := r.CheckTrx(ctx)

	var datas []model.ImportedWorksheetDetailEntityModel
	var info abstraction.PaginationInfo

	query := conn.Model(&model.ImportedWorksheetDetailEntityModel{})
	//filter
	query = r.Filter(ctx, query, *m)

	//sort
	if p.Sort == nil {
		sort := "asc"
		p.Sort = &sort
	}
	if p.SortBy == nil {
		sortBy := "id"
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

	info.Count = int(totalData)
	info.MoreRecords = false
	if len(datas) > *p.PageSize {
		info.MoreRecords = true
		// info.Count -= 1
		// datas = datas[:len(datas)-1]
	}

	return &datas, &info, nil
}
func (r *importedworksheetdetail) FindWithCode(ctx *abstraction.Context, code *string) (*[]model.ImportedWorksheetDetailEntityModel, error) {
	conn := r.CheckTrx(ctx)

	var data []model.ImportedWorksheetDetailEntityModel
	tmp := fmt.Sprintf("%s", *code)
	if err := conn.Where("code LIKE ?", tmp+"%").Find(&data).WithContext(ctx.Request().Context()).Error; err != nil {
		return &data, err
	}
	return &data, nil
}
func (r *importedworksheetdetail) FindByID(ctx *abstraction.Context, id *int) (*model.ImportedWorksheetDetailEntityModel, error) {
	conn := r.CheckTrx(ctx)

	var data model.ImportedWorksheetDetailEntityModel
	if err := conn.Where("id = ?", &id).Find(&data).WithContext(ctx.Request().Context()).Error; err != nil {
		return &data, err
	}
	return &data, nil
}
