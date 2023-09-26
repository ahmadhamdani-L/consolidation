package service

import (
	"encoding/json"
	"notification/internal/abstraction"
	"notification/internal/centrifugo"
	"notification/internal/factory"
	"notification/internal/model"
	"notification/internal/repository"

	"gorm.io/gorm"
)

type service struct {
	Db                *gorm.DB
	CompanyRepository repository.Company
	UserRepository    repository.User
}

type Service interface {
	SendNotif(ctx *abstraction.Context, payload *abstraction.JsonData)
}

func NewService(f *factory.Factory) *service {
	companyRepository := f.CompanyRepository
	userRepository := f.UserRepository
	db := f.Db

	return &service{
		Db:                db,
		CompanyRepository: companyRepository,
		UserRepository:    userRepository,
	}
}

func (s *service) SendNotif(data model.JsonData) {
	ctx := new(abstraction.Context)
	user, err := s.UserRepository.FindByID(ctx, &data.UserID)
	if err != nil {
		return
	}
	company, err := s.CompanyRepository.FindByID(ctx, &data.CompanyID)
	if err != nil {
		return
	}

	message := model.Message{
		NotificationID: data.ID,
		Name:           data.Name,
		User:           user.Name,
		Company:        company.Name,
		Period:         data.Filter.Period,
		Versions:       data.Filter.Versions,
		Message:        data.Data,
		FileUrl:        data.FileLoc,
		Timestamp:      data.Timestamp,
	}
	jsonStr, err := json.Marshal(message)
	if err != nil {
		panic(err)
	}
	go centrifugo.NewService("NOTIFICATION", "NOTIFICATION", user.ID).Subs().BroadcastMessage(jsonStr)
}
