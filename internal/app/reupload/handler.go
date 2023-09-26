package reupload

import (
	"worker/internal/abstraction"
	"worker/internal/factory"
)

type handler struct {
	Service *service
}

func NewHandler(f *factory.Factory) *handler {
	return &handler{
		Service: NewService(f),
	}
}

func (h *handler) ReUpload(dataReUpload abstraction.JsonDataReUpload) {
	cc := abstraction.Context{
		Auth: &abstraction.AuthContext{
			ID:        dataReUpload.UserID,
			CompanyID: dataReUpload.CompanyID,
		},
	}
	h.Service.ReUpload(&cc, &dataReUpload)
	return
}