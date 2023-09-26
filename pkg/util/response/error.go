package response

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

type errorResponse struct {
	Meta  Meta   `json:"meta"`
	Error string `json:"error"`
}

type Error struct {
	Response     errorResponse `json:"response"`
	Code         int           `json:"code"`
	Codes		string			`json:"codes"`
	ErrorMessage error
}

const (
	E_DUPLICATE            = "duplicate"
	E_NOT_FOUND            = "not_found"
	E_UNPROCESSABLE_ENTITY = "unprocessable_entity"
	E_UNAUTHORIZED         = "unauthorized"
	E_BAD_REQUEST          = "bad_request"
	E_SERVER_ERROR         = "server_error"
	E_COA_NOT_FOUND 	   = "coa tidak ditemukan"
)

type errorConstant struct {
	Duplicate           Error
	NotFound            Error
	RouteNotFound       Error
	UnprocessableEntity Error
	Unauthorized        Error
	BadRequest          Error
	Validation          Error
	InternalServerError Error
	UnauthorizedAccess  Error
	DataValidated       Error
	HasRelatedData      Error
	CoaNotFound 		Error
}

var ErrorConstant errorConstant = errorConstant{
	Duplicate: Error{
		Response: errorResponse{
			Meta: Meta{
				Success: false,
				Message: "Created value already exists",
			},
			Error: E_DUPLICATE,
		},
		Code: http.StatusConflict,
	},
	NotFound: Error{
		Response: errorResponse{
			Meta: Meta{
				Success: false,
				Message: "Data not found",
			},
			Error: E_NOT_FOUND,
		},
		Code: http.StatusNotFound,
	},
	RouteNotFound: Error{
		Response: errorResponse{
			Meta: Meta{
				Success: false,
				Message: "Route not found",
			},
			Error: E_NOT_FOUND,
		},
		Code: http.StatusNotFound,
	},
	UnprocessableEntity: Error{
		Response: errorResponse{
			Meta: Meta{
				Success: false,
				Message: "Invalid parameters or payload",
			},
			Error: E_UNPROCESSABLE_ENTITY,
		},
		Code: http.StatusUnprocessableEntity,
	},
	Unauthorized: Error{
		Response: errorResponse{
			Meta: Meta{
				Success: false,
				Message: "Unauthorized, please login",
			},
			Error: E_UNAUTHORIZED,
		},
		Code: http.StatusUnauthorized,
	},
	BadRequest: Error{
		Response: errorResponse{
			Meta: Meta{
				Success: false,
				Message: "Bad Request",
			},
			Error: E_BAD_REQUEST,
		},
		Code: http.StatusBadRequest,
	},
	Validation: Error{
		Response: errorResponse{
			Meta: Meta{
				Success: false,
				Message: "Invalid parameters or payload",
			},
			Error: E_BAD_REQUEST,
		},
		Code: http.StatusBadRequest,
	},
	InternalServerError: Error{
		Response: errorResponse{
			Meta: Meta{
				Success: false,
				Message: "Something bad happened",
			},
			Error: E_SERVER_ERROR,
		},
		Code: http.StatusInternalServerError,
	},
	UnauthorizedAccess: Error{
		Response: errorResponse{
			Meta: Meta{
				Success: false,
				Message: "Unauthorized Access to this resource",
			},
			Error: E_UNAUTHORIZED,
		},
		Code: http.StatusUnauthorized,
	},
	DataValidated: Error{
		Response: errorResponse{
			Meta: Meta{
				Success: false,
				Message: "Cannot Modify because data status has been validated",
			},
			Error: E_UNPROCESSABLE_ENTITY,
		},
		Code: http.StatusUnprocessableEntity,
	},
	HasRelatedData: Error{
		Response: errorResponse{
			Meta: Meta{
				Success: false,
				Message: "Cannot Delete because data has related data",
			},
			Error: E_UNPROCESSABLE_ENTITY,
		},
		Code: http.StatusUnprocessableEntity,
	},
	CoaNotFound: Error{
		Response: errorResponse{
			Meta: Meta{
				Success: false,
				Message: "Coa Tidak Ditemukan",
			},
			Error: E_COA_NOT_FOUND,
		},
		Code: http.StatusConflict,
	},
}

func ErrorBuilder(res *Error, message error) *Error {
	res.ErrorMessage = message
	return res
}
func CustomErrorBuilders(codes string, err string, message string) *Error {
	return &Error{
		Response: errorResponse{
			Meta: Meta{
				Success: false,
				Message: message,
			},
			Error: err,
		},
		Codes: codes,
	}
}
func CustomErrorBuilder(code int, err string, message string) *Error {
	return &Error{
		Response: errorResponse{
			Meta: Meta{
				Success: false,
				Message: message,
			},
			Error: err,
		},
		Code: code,
	}
}
func ErrorBuilders(err *Error) *Error {
	return err
}
func ErrorResponse(err error) *Error {
	re, ok := err.(*Error)
	if ok {
		return re
	} else {
		return ErrorBuilder(&ErrorConstant.InternalServerError, err)
	}
}

func (e *Error) Error() string {
	return fmt.Sprintf("error code %d", e.Code)
}

func (e *Error) ParseToError() error {
	return e
}

func (e *Error) Send(c echo.Context) error {
	var errorMessage string
	if e.ErrorMessage != nil {
		errorMessage = fmt.Sprintf("%+v", errors.WithStack(e.ErrorMessage))
	}
	logrus.Error(errorMessage)

	_, err := ioutil.ReadAll(c.Request().Body)
	if err != nil {
		logrus.Warn("error read body, message : ", e.Error())
	}

	_, err = json.Marshal(c.Request().Header)
	if err != nil {
		logrus.Warn("error read header, message : ", e.Error())
	}

	// go func() {
	// 	retries := 3
	// 	logError := log.LogError{
	// 		ID:           shortid.MustGenerate(),
	// 		Header:       string(bHeader),
	// 		Body:         string(body),
	// 		URL:          c.Request().URL.Path,
	// 		HttpMethod:   c.Request().Method,
	// 		ErrorMessage: errorMessage,
	// 		Level:        "Error",
	// 		AppName:      os.Getenv("APP"),
	// 		Version:      os.Getenv("VERSION"),
	// 		Env:          os.Getenv("ENV"),
	// 		CreatedAt:    *date.DateTodayLocal(),
	// 	}
	// 	for i := 0; i < retries; i++ {
	// 		err := log.InsertErrorLog(context.Background(), &logError)
	// 		if err == nil {
	// 			break
	// 		}
	// 	}
	// }()

	return c.JSON(e.Code, e.Response)
}
