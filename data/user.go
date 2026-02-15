package data

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-auth/database"
)

func GetUsers(c *gin.Context) {
	var users []database.User
	database.Db.Find(&users)

	c.JSON(http.StatusOK, gin.H{
		"status": "ok",
		"users":  users,
	})

}
