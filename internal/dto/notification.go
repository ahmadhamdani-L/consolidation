package dto

import (
	"mcash-finance-console-core/internal/abstraction"
	"mcash-finance-console-core/internal/model"
	res "mcash-finance-console-core/pkg/util/response"
)

// Get
type NotificationGetRequest struct {
	abstraction.Pagination
	model.NotificationFilterModel
}

type NotifCountData struct {
	Total            int64                           `json:"total"`
	Unread           int64                           `json:"total_unread"`
	Read             int64                           `json:"total_read"`
	NotificationData []model.NotificationEntityModel `json:"notification_data"`
}
type NotificationGetResponse struct {
	Datas          NotifCountData
	PaginationInfo abstraction.PaginationInfo
}
type NotificationGetResponseDoc struct {
	Body struct {
		Meta res.Meta       `json:"meta"`
		Data NotifCountData `json:"data"`
	} `json:"body"`
}

// GetByID
type NotificationGetByIDRequest struct {
	ID int `param:"id" validate:"required,numeric"`
}
type NotificationGetByIDResponse struct {
	model.NotificationEntityModel
}
type NotificationGetByIDResponseDoc struct {
	Body struct {
		Meta res.Meta                    `json:"meta"`
		Data NotificationGetByIDResponse `json:"data"`
	} `json:"body"`
}

// MarkAsRead
type NotificationMarkAsReadRequest struct {
	// ID int `query:"id" validate:"required,numeric"`
	ArrID *[]int `json:"id" validate:"required"`
}
type NotificationMarkAsReadResponse struct {
	Data []model.NotificationEntityModel
}
type NotificationMarkAsReadResponseDoc struct {
	Body struct {
		Meta res.Meta                         `json:"meta"`
		Data []NotificationMarkAsReadResponse `json:"data"`
	} `json:"body"`
}
