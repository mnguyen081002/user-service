package controller

import (
	"erp/internal/api/request"
	"erp/internal/api/response"
	"erp/internal/api_errors"
	utils2 "erp/internal/utils"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type BaseController struct {
}

func Response(c *gin.Context, statusCode int, message string, data interface{}) {
	c.JSON(statusCode, response.SimpleResponse{
		Data:    data,
		Message: message,
	})
}

func ResponseList(c *gin.Context, message string, total *int64, data interface{}) {
	var o request.PageOptions
	if err := c.ShouldBindQuery(&o); err != nil {
		ResponseValidationError(c, err)
		return
	}

	if o.Limit == 0 {
		o.Limit = 10
	}

	if o.Page == 0 {
		o.Page = 1
	}

	pageCount := utils2.GetPageCount(*total, o.Limit)
	c.JSON(http.StatusOK, response.SimpleResponseList{
		Message: message,
		Data:    data,
		Meta: response.Meta{
			Total:       total,
			Page:        o.Page,
			Limit:       o.Limit,
			Sort:        o.Sort,
			PageCount:   pageCount,
			HasPrevPage: o.Page > 1,
			HasNextPage: o.Page < pageCount,
		},
	})
}

func ResponseError(c *gin.Context, err error) {

	mas, ok := api_errors.MapErrorCodeMessage[err.Error()]
	var status int
	ginType := gin.ErrorTypePublic
	errp := err
	if !ok {
		status = http.StatusInternalServerError
		ginType = gin.ErrorTypePrivate
		mas = api_errors.MapErrorCodeMessage[api_errors.ErrInternalServerError]
		errp = errors.New(api_errors.ErrInternalServerError)
	}

	c.Errors = append(c.Errors, &gin.Error{
		Err:  err,
		Type: ginType,
		Meta: status,
	})

	c.AbortWithStatusJSON(mas.Status, response.ResponseError{
		Code:    errp.Error(),
		Message: mas.Message,
	})
}

func ResponseValidationError(c *gin.Context, err error) {
	var ve validator.ValidationErrors
	if errors.As(err, &ve) {
		err = errors.New(utils2.StructPascalToSnakeCase(ve[0].Field()) + " is " + ve[0].Tag())
	}

	c.AbortWithStatusJSON(http.StatusUnprocessableEntity, response.ResponseError{
		Code:    api_errors.ErrValidation,
		Message: err.Error(),
	})
}
