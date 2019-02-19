package controller

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

type Controller interface {
	Validate(c *gin.Context, validator interface{}) bool
}

type BaseController struct {}

func (bc *BaseController) Validate (c *gin.Context, validator interface{}) bool {
	if err := c.ShouldBindJSON(&validator); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return false
	}

	return true
}