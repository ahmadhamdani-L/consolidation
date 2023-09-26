package notification

import (
	"encoding/json"
	"errors"
	"fmt"
	"mcash-finance-console-core/internal/abstraction"
	"mcash-finance-console-core/internal/dto"
	"mcash-finance-console-core/internal/factory"
	"mcash-finance-console-core/internal/model"
	"mcash-finance-console-core/internal/repository"
	"mcash-finance-console-core/pkg/kafka"
	"mcash-finance-console-core/pkg/util/helper"
	"mcash-finance-console-core/pkg/util/response"
	"mcash-finance-console-core/pkg/util/trxmanager"
	"time"

	"gorm.io/gorm"
)

type service struct {
	Repository repository.Notification
	Db         *gorm.DB
}

type Service interface {
	Find(ctx *abstraction.Context, payload *dto.NotificationGetRequest) (*dto.NotificationGetResponse, error)
	FindByID(ctx *abstraction.Context, payload *dto.NotificationGetByIDRequest) (*dto.NotificationGetByIDResponse, error)
	Update(ctx *abstraction.Context, payload *dto.NotificationMarkAsReadRequest) (*dto.NotificationMarkAsReadResponse, error)
	Test(ctx *abstraction.Context) error
}

func NewService(f *factory.Factory) *service {
	repository := f.NotificationRepository
	db := f.Db
	return &service{
		Repository: repository,
		Db:         db,
	}
}

func (s *service) Find(ctx *abstraction.Context, payload *dto.NotificationGetRequest) (*dto.NotificationGetResponse, error) {
	payload.CreatedBy = &ctx.Auth.ID
	data, info, err := s.Repository.Find(ctx, &payload.NotificationFilterModel, &payload.Pagination)
	if err != nil {
		return &dto.NotificationGetResponse{}, helper.ErrorHandler(err)
	}
	var datas dto.NotifCountData
	datas.Total, datas.Read, datas.Unread, err = s.Repository.Count(ctx, &payload.NotificationFilterModel)
	if err != nil {
		return nil, err
	}
	datas.NotificationData = *data
	result := &dto.NotificationGetResponse{
		Datas:          datas,
		PaginationInfo: *info,
	}
	return result, nil
}

func (s *service) FindByID(ctx *abstraction.Context, payload *dto.NotificationGetByIDRequest) (*dto.NotificationGetByIDResponse, error) {
	data, err := s.Repository.FindByID(ctx, &payload.ID)
	if err != nil {
		return &dto.NotificationGetByIDResponse{}, helper.ErrorHandler(err)
	}
	result := &dto.NotificationGetByIDResponse{
		NotificationEntityModel: *data,
	}
	return result, nil
}

func (s *service) MarkAsRead(ctx *abstraction.Context, payload *dto.NotificationMarkAsReadRequest) (*dto.NotificationMarkAsReadResponse, error) {
	var datas []model.NotificationEntityModel
	if err := trxmanager.New(s.Db).WithTrx(ctx, func(ctx *abstraction.Context) error {
		for _, id := range *payload.ArrID {
			var data model.NotificationEntityModel
			notifData, err := s.Repository.FindByID(ctx, &id)
			if err != nil {
				return helper.ErrorHandler(err)
			}
			if notifData.CreatedBy != ctx.Auth.ID {
				return response.ErrorBuilder(&response.ErrorConstant.UnauthorizedAccess, errors.New("Cannot update data"))
			}
			data.Context = ctx
			tru := true
			dataToUpdate := model.NotificationEntity{
				IsOpened: &tru,
			}
			data.NotificationEntity = dataToUpdate
			_, err = s.Repository.Update(ctx, &notifData.ID, &data)
			if err != nil {
				return helper.ErrorHandler(err)
			}
			datas = append(datas, data)
		}
		return nil
	}); err != nil {
		return &dto.NotificationMarkAsReadResponse{}, err
	}
	result := &dto.NotificationMarkAsReadResponse{
		Data: datas,
	}
	return result, nil
}

func (s *service) Test(ctx *abstraction.Context) error {
	waktu := time.Now()
	msg := kafka.JsonData{
		UserID:    ctx.Auth.ID,
		CompanyID: ctx.Auth.CompanyID,
		Timestamp: &waktu,
		Name:      "Lutfi",
		FileLoc:   "/filenya",
	}

	jsonStr, err := json.Marshal(msg)
	if err != nil {
		fmt.Printf("Error: %s", err.Error())
		return err
	}
	kafka.NewService("NOTIFICATION").SendMessage("NOTIFICATION", string(jsonStr))
	return nil
}
