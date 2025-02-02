package bizerror

import (
	"encoding/json"
	"errors"
	"fmt"
	"fundinsight/pkg/misc"
	"io"
	"net/http"
	"runtime/debug"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
	"github.com/sirupsen/logrus"
)

func ErrorHandling() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer handle(c)
		c.Next()
	}
}

func handle(c *gin.Context) {
	if ret := recover(); ret != nil {
		err, ok := ret.(error)
		if !ok {
			err = errors.New(fmt.Sprintf("%s", ret))
		}
		HandleError(c, err)
	} else {
		if err := c.Errors.Last(); err != nil {
			HandleError(c, err)
		}
	}
}

func HandleError(c *gin.Context, err error) {
	logrus.Error(err)
	debug.Stack()

	genericErr := err
	var ginErr *gin.Error
	if errors.As(err, &ginErr) {
		genericErr = ginErr.Err
	}

	if bizErr, ok := genericErr.(BizError); ok {
		respond := bizErr.Respond()
		c.JSON(respond.Status, &misc.ErrorBody{Code: respond.Code, Message: respond.Message, Data: respond.Data})
		c.Abort()
		return
	}

	// bad request:  io.EOF (no body).
	if errors.Is(genericErr, io.EOF) {
		c.JSON(http.StatusBadRequest, &misc.ErrorBody{Code: "bad_request.body_not_found", Message: "body not found"})
		c.Abort()
		return
	}
	// bad request: json syntax Error
	if syntaxErr, ok := genericErr.(*json.SyntaxError); ok {
		c.JSON(http.StatusBadRequest, &misc.ErrorBody{Code: "bad_request.invalid_body_format", Message: "invalid body format", Data: syntaxErr.Error()})
		c.Abort()
		return
	}
	// validation failed
	if validationErr, ok := genericErr.(validator.ValidationErrors); ok {
		c.JSON(http.StatusBadRequest, &misc.ErrorBody{Code: "bad_request.validation_failed", Message: "validation failed", Data: validationErr.Error()})
		c.Abort()
		return
	}

	if errors.Is(genericErr, ErrUnauthenticated) {
		c.JSON(http.StatusUnauthorized, &misc.ErrorBody{Code: "common.unauthenticated", Message: "unauthenticated"})
		c.Abort()
		return
	}
	if errors.Is(genericErr, ErrForbidden) {
		c.JSON(http.StatusForbidden, &misc.ErrorBody{Code: "security.forbidden", Message: "access forbidden"})
		c.Abort()
		return
	}
	if errors.Is(genericErr, ErrUnknownState) {
		c.JSON(http.StatusBadRequest, &misc.ErrorBody{Code: "workflow.unknown_state", Message: "unknown state"})
		c.Abort()
		return
	}
	if errors.Is(genericErr, ErrStateExisted) {
		c.JSON(http.StatusBadRequest, &misc.ErrorBody{Code: "workflow.state_existed", Message: "state existed"})
		c.Abort()
		return
	}
	if errors.Is(genericErr, gorm.ErrRecordNotFound) || errors.Is(genericErr, ErrNotFound) {
		c.JSON(http.StatusNotFound, &misc.ErrorBody{Code: "common.record_not_found", Message: "record not found"})
		c.Abort()
		return
	}
	if errors.Is(genericErr, mysql.ErrInvalidConn) {
		c.JSON(http.StatusServiceUnavailable, &misc.ErrorBody{Code: "common.internal_server_error", Message: err.Error()})
		c.Abort()
		return
	}

	c.JSON(500, &misc.ErrorBody{Code: "common.internal_server_error", Message: err.Error()})
	c.Abort()
}
