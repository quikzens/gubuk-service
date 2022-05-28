package util

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type response struct {
	Code   int         `json:"code"`
	Status string      `json:"status"`
	Data   interface{} `json:"data,omitempty"`
	Error  string      `json:"error,omitempty"`
}

func SendSuccess(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, response{
		Code:   200,
		Status: "SUCCESS",
		Data:   data,
	})
}

func SendBadRequest(c *gin.Context, err error) {
	c.JSON(http.StatusBadRequest, response{
		Code:   400,
		Status: "BAD REQUEST",
		Error:  err.Error(),
	})
	c.Abort()
}

func SendServerError(c *gin.Context, err error) {
	c.JSON(http.StatusInternalServerError, response{
		Code:   500,
		Status: "SERVER ERROR",
		Error:  err.Error(),
	})
	c.Abort()
}

func SendUnauthorized(c *gin.Context, err error) {
	c.JSON(http.StatusUnauthorized, response{
		Code:   401,
		Status: "UNAUTHORIZED",
		Error:  err.Error(),
	})
	c.Abort()
}

func SendNotFound(c *gin.Context, err error) {
	c.JSON(http.StatusNotFound, response{
		Code:   404,
		Status: "NOT FOUND",
		Error:  err.Error(),
	})
	c.Abort()
}
