package delivery

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

// HttpSuccessRespond message
type HttpSuccessRespond struct {
	Code   int         `json:"code" example:"200"`
	Status string      `json:"status" example:"OK"`
	Data   interface{} `json:"data"`
}

// HttpErrorRespond message
type HttpErrorRespond struct {
	Code   int    `json:"code" example:"400"`
	Status string `json:"status" example:"Bad Request"`
	Data   string `json:"data"`
}

// HttpValidationErrorRespond message
type HttpValidationErrorRespond struct {
	Code   int         `json:"code" example:"422"`
	Status string      `json:"status" example:"Unprocessable Entity"`
	Data   interface{} `json:"data"`
}

// HttpServerErrorRespond message
type HttpServerErrorRespond struct {
	Code   int    `json:"code" example:"500"`
	Status string `json:"status" example:"Internal Server Error"`
	Data   string `json:"data"`
}

// NewHttpRespond godoc
func NewHttpRespond(context *gin.Context, code int, data interface{}) {
	if code == http.StatusOK || code == http.StatusCreated {
		context.JSON(
			code,
			HttpSuccessRespond{
				Code:   code,
				Status: http.StatusText(code),
				Data:   data,
			},
		)

		return
	}

	if code == http.StatusUnprocessableEntity {
		context.JSON(
			code,
			HttpValidationErrorRespond{
				Code:   code,
				Status: http.StatusText(code),
				Data:   data,
			},
		)

		return
	}

	msg := func() string {
		switch {
		case data != nil:
			return data.(string)
		case code == http.StatusBadRequest:
			return "something went wrong with the request"
		default:
			return "something went wrong with the server"
		}
	}()

	context.JSON(
		code,
		HttpErrorRespond{
			Code:   code,
			Status: http.StatusText(code),
			Data:   msg,
		},
	)

	return
}
