package coa

import (
	"fmt"
	"mcash-finance-console-core/internal/abstraction"
	"mcash-finance-console-core/internal/dto"
	"mcash-finance-console-core/internal/factory"
	"mcash-finance-console-core/internal/model"
	"mcash-finance-console-core/internal/repository"
	"mcash-finance-console-core/pkg/util/helper"
	"mcash-finance-console-core/pkg/util/response"
	"mcash-finance-console-core/pkg/util/trxmanager"
	"os"
	"regexp"
	"strings"

	"github.com/xuri/excelize/v2"
	"gorm.io/gorm"
)

type service struct {
	Repository repository.Coa
	Db         *gorm.DB
}

type Service interface {
	Find(ctx *abstraction.Context, payload *dto.CoaGetRequest) (*dto.CoaGetResponse, error)
	FindByID(ctx *abstraction.Context, payload *dto.CoaGetByIDRequest) (*dto.CoaGetByIDResponse, error)
	Create(ctx *abstraction.Context, payload *dto.CoaCreateRequest) (*dto.CoaCreateResponse, error)
	Update(ctx *abstraction.Context, payload *dto.CoaUpdateRequest) (*dto.CoaUpdateResponse, error)
	Delete(ctx *abstraction.Context, payload *dto.CoaDeleteRequest) (*dto.CoaDeleteResponse, error)
}

func NewService(f *factory.Factory) *service {
	repository := f.CoaRepository
	db := f.Db
	return &service{
		Repository: repository,
		Db:         db,
	}
}

func (s *service) Find(ctx *abstraction.Context, payload *dto.CoaGetRequest) (*dto.CoaGetResponse, error) {
	data, info, err := s.Repository.Find(ctx, &payload.CoaFilterModel, &payload.Pagination)
	if err != nil {
		return &dto.CoaGetResponse{}, helper.ErrorHandler(err)
	}

	result := &dto.CoaGetResponse{
		Datas:          *data,
		PaginationInfo: *info,
	}
	return result, nil
}

func (s *service) FindByID(ctx *abstraction.Context, payload *dto.CoaGetByIDRequest) (*dto.CoaGetByIDResponse, error) {
	data, err := s.Repository.FindByID(ctx, &payload.ID)
	if err != nil {
		return &dto.CoaGetByIDResponse{}, helper.ErrorHandler(err)
	}
	result := &dto.CoaGetByIDResponse{
		CoaEntityModel: *data,
	}
	return result, nil
}

func (s *service) Create(ctx *abstraction.Context, payload *dto.CoaCreateRequest) (*dto.CoaCreateResponse, error) {
	var data model.CoaEntityModel

	if err := trxmanager.New(s.Db).WithTrx(ctx, func(ctx *abstraction.Context) error {
		data.Context = ctx
		data.CoaEntity = payload.CoaEntity

		result, err := s.Repository.Create(ctx, &data)
		if err != nil {
			return response.ErrorBuilder(&response.ErrorConstant.UnprocessableEntity, err)
		}
		data = *result
		return nil
	}); err != nil {
		return &dto.CoaCreateResponse{}, err
	}

	result := &dto.CoaCreateResponse{
		CoaEntityModel: data,
	}
	return result, nil
}

func (s *service) Update(ctx *abstraction.Context, payload *dto.CoaUpdateRequest) (*dto.CoaUpdateResponse, error) {
	var data model.CoaEntityModel
	if err := trxmanager.New(s.Db).WithTrx(ctx, func(ctx *abstraction.Context) error {
		if _, err := s.Repository.FindByID(ctx, &payload.ID); err != nil {
			return helper.ErrorHandler(err)
		}
		data.Context = ctx
		data.CoaEntity = payload.CoaEntity

		result, err := s.Repository.Update(ctx, &payload.ID, &data)
		if err != nil {
			return helper.ErrorHandler(err)
		}
		data = *result
		return nil
	}); err != nil {
		return &dto.CoaUpdateResponse{}, err
	}
	result := &dto.CoaUpdateResponse{
		CoaEntityModel: data,
	}
	return result, nil
}

func (s *service) Delete(ctx *abstraction.Context, payload *dto.CoaDeleteRequest) (*dto.CoaDeleteResponse, error) {
	var data model.CoaEntityModel
	if err := trxmanager.New(s.Db).WithTrx(ctx, func(ctx *abstraction.Context) error {
		if _, err := s.Repository.FindByID(ctx, &payload.ID); err != nil {
			return helper.ErrorHandler(err)
		}
		data.Context = ctx
		result, err := s.Repository.Delete(ctx, &payload.ID, &data)
		if err != nil {
			return helper.ErrorHandler(err)
		}
		data = *result
		return nil
	}); err != nil {
		return &dto.CoaDeleteResponse{}, err
	}
	result := &dto.CoaDeleteResponse{
		CoaEntityModel: data,
	}
	return result, nil
}

func (s *service) Export(ctx *abstraction.Context, payload *dto.CoaExportRequest) (*dto.CoaExportResponse, error) {
	coaPaging := abstraction.Pagination{}
	nolimit := 10000
	coaPaging.PageSize = &nolimit

	coaGroup, err := s.Repository.Export(ctx)
	if err != nil {
		return nil, err
	}

	f := excelize.NewFile()
	sheetCoa := "COA"
	f.SetSheetName(f.GetSheetName(0), sheetCoa)

	hitungstrip, err := s.Repository.CountDash(ctx)
	if err != nil {
		return nil, err
	}
	hitungstrip += 5

	tmplastColGroup, err := excelize.CoordinatesToCellName(hitungstrip, 1)
	if err != nil {
		return nil, err
	}
	tmplastColGroup = strings.ReplaceAll(tmplastColGroup, "1", "")

	colCoaType, err := excelize.CoordinatesToCellName(hitungstrip+1, 1)
	if err != nil {
		return nil, err
	}
	colCoaType = strings.ReplaceAll(colCoaType, "1", "")

	colCoaGroup, err := excelize.CoordinatesToCellName(hitungstrip+2, 1)
	if err != nil {
		return nil, err
	}
	colCoaGroup = strings.ReplaceAll(colCoaGroup, "1", "")

	colCoaName, err := excelize.ColumnNameToNumber("F")
	if err != nil {
		return nil, err
	}

	f.SetCellValue(sheetCoa, "B3", "COA")
	f.SetCellValue(sheetCoa, "C3", "Account Name")
	f.SetCellValue(sheetCoa, fmt.Sprintf("%s3", colCoaType), "Type")
	f.SetCellValue(sheetCoa, fmt.Sprintf("%s3", colCoaGroup), "Group")

	f.SetColWidth(sheetCoa, "B", "B", 10)
	f.SetColWidth(sheetCoa, "C", "C", 50)
	f.SetColWidth(sheetCoa, "D", "D", 0)
	f.SetColWidth(sheetCoa, "E", "E", 2)
	f.SetColWidth(sheetCoa, "F", "F", 8)
	f.SetColWidth(sheetCoa, "G", "G", 34)
	f.SetColWidth(sheetCoa, "H", "H", 17)
	for i := 8; i <= hitungstrip; i++ {
		tmpCoaGroupping, err := excelize.CoordinatesToCellName(hitungstrip+1, 1)
		if err != nil {
			return nil, err
		}
		tmpCoaGroupping = strings.ReplaceAll(tmpCoaGroupping, "1", "")
		f.SetColWidth(sheetCoa, tmpCoaGroupping, tmpCoaGroupping, 17)
	}
	f.SetColWidth(sheetCoa, colCoaType, colCoaType, 25)
	f.SetColWidth(sheetCoa, colCoaGroup, colCoaGroup, 13)

	err = f.MergeCell(sheetCoa, "D2", fmt.Sprintf("%s2", tmplastColGroup))
	if err != nil {
		return nil, err
	}
	f.SetCellValue(sheetCoa, "D2", "Grouping")

	styleHeader, err := f.NewStyle(&excelize.Style{
		Font: &excelize.Font{
			Bold: true,
		},
		Fill: excelize.Fill{
			Type:    "pattern",
			Pattern: 1,
			Color:   []string{"#f8cbad"},
		},
		Alignment: &excelize.Alignment{
			WrapText:   true,
			Horizontal: "center",
			Vertical:   "center",
		},
	})
	if err != nil {
		return nil, err
	}
	f.SetCellStyle(sheetCoa, "B2", fmt.Sprintf("%s3", colCoaGroup), styleHeader)

	row := 4
	for _, vCGroup := range *coaGroup {
		for _, vCType := range vCGroup.CoaType {
			for _, vCoa := range vCType.Coa {
				tmpColCoaName := colCoaName
				re := regexp.MustCompile(`\(-`)
				newCoa := re.ReplaceAllString(vCoa.Name, "(")
				coaSplit := strings.Split(newCoa, " - ")
				f.SetCellValue(sheetCoa, fmt.Sprintf("B%d", row), vCoa.Code)
				f.SetCellValue(sheetCoa, fmt.Sprintf("C%d", row), vCoa.Name)
				f.SetCellValue(sheetCoa, fmt.Sprintf("D%d", row), nil)
				f.SetCellValue(sheetCoa, fmt.Sprintf("E%d", row), "-")
				for _, v := range coaSplit {
					colCoaName, err := excelize.CoordinatesToCellName(tmpColCoaName, 1)
					if err != nil {
						return nil, err
					}
					colCoaName = strings.ReplaceAll(colCoaName, "1", "")
					f.SetCellValue(sheetCoa, fmt.Sprintf("%s%d", colCoaName, row), v)
					tmpColCoaName++
				}
				f.SetCellValue(sheetCoa, fmt.Sprintf("%s%d", colCoaType, row), vCType.Name)
				f.SetCellValue(sheetCoa, fmt.Sprintf("%s%d", colCoaGroup, row), vCGroup.Name)

				row++
			}
		}
	}

	err = f.AutoFilter(sheetCoa, "B3", fmt.Sprintf("%s3", colCoaGroup), "")
	if err != nil {
		return nil, err
	}
	tmpFolder := fmt.Sprintf("assets/%d", ctx.Auth.ID)
	_, err = os.Stat(tmpFolder)
	if err != nil {
		if os.IsNotExist(err) {
			if err := os.MkdirAll(tmpFolder, os.ModePerm); err != nil {
				return nil, err
			}
		} else {
			return nil, err
		}
	}
	fileName := "COA.xlsx"
	fileLoc := fmt.Sprintf("assets/%d/%s", ctx.Auth.ID, fileName)
	err = f.SaveAs(fileLoc)
	if err != nil {
		return nil, err
	}

	result := &dto.CoaExportResponse{
		FileName: fileName,
		Path:     fileLoc,
	}
	return result, nil
}
