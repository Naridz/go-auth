package auth

import (
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-auth/database"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type RegisterRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// var jwtSecret = []byte(os.Getenv("JWT_SECRET"))

func Register(c *gin.Context) {
	var register RegisterRequest
	if err := c.ShouldBindJSON(&register); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	hashPassword, err := bcrypt.GenerateFromPassword([]byte(register.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to hash password",
		})
		return
	}
	user := database.User{Username: register.Username, Password: string(hashPassword)}
	if err := database.Db.Create(&user).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": "username already exists"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "ok", "message": "register successfully", "userID": user.ID, "username": user.Username})
}

func Login(c *gin.Context) {
	var user LoginRequest
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	loginUser := new(database.User)
	res := database.Db.Where("username = ?", user.Username).First(&loginUser)
	if res.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": "invalid username or password"})
		return
	}
	err := bcrypt.CompareHashAndPassword([]byte(loginUser.Password), []byte(user.Password))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": "invalid username or password"})
		return
	}
	claims := jwt.MapClaims{
		"user_id":  loginUser.ID,
		"username": loginUser.Username,
		"exp":      time.Now().Add(time.Minute * 5).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "error"})
		return
	}

	tokenString, err := token.SignedString([]byte(secret))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "could not generate token"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "ok", "message": "login success", "token": tokenString})
}
