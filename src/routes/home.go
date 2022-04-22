package routes

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func HomePage(c *gin.Context) {

	c.JSON(http.StatusOK, gin.H{
		"message": "Welcome to Account Management System.",
		//"data":    o,
	})
}
