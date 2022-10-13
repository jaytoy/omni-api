package controllers

import (
	"net/http"
	"strings"
	"time"

	"blitzomni.com/m/models"
	"blitzomni.com/m/utils"
	"github.com/gin-gonic/gin"
	"github.com/thanhpk/randstr"
	"gorm.io/gorm"
)

type SignUpInput struct {
	FirstName       string `json:"first_name" binding:"required"`
	LastName        string `json:"last_name" binding:"required"`
	Email           string `json:"email" binding:"required"`
	Password        string `json:"password" binding:"required,min=8"`
	PasswordConfirm string `json:"password_confirm" binding:"required"`
}

type LoginInput struct {
	Email    string `json:"email"  binding:"required"`
	Password string `json:"password"  binding:"required"`
}

type AuthController struct {
	DB *gorm.DB
}

func NewAuthController(DB *gorm.DB) AuthController {
	return AuthController{DB}
}

func (ac *AuthController) SignUp(c *gin.Context) {
	var payload SignUpInput

	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": err.Error()})
		return
	}

	if payload.Password != payload.PasswordConfirm {
		c.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": "Passwords do not match"})
		return
	}

	newUser := models.User{
		FirstName: payload.FirstName,
		LastName:  payload.LastName,
		Email:     strings.ToLower(payload.Email),
		Password:  utils.HashPassword(payload.Password),
		Verified:  false,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	result := ac.DB.Create(&newUser)

	if result.Error != nil && strings.Contains(result.Error.Error(), "duplicate key value violates unique") {
		c.JSON(http.StatusConflict, gin.H{"status": "fail", "message": "User with that email already exists"})
		return
	} else if result.Error != nil {
		c.JSON(http.StatusBadGateway, gin.H{"status": "error", "message": "Something bad happened"})
		return
	}

	config, _ := utils.LoadConfig(".")

	// Generate Verification Code
	code := randstr.String(20)

	verification_code := utils.Encode(code)

	// Update User in Database
	newUser.VerificationCode = verification_code
	ac.DB.Save(newUser)

	var firstName = newUser.FirstName

	if strings.Contains(firstName, " ") {
		firstName = strings.Split(firstName, " ")[1]
	}

	// Send Email
	emailData := utils.EmailData{
		URL:       config.ClientOrigin + "/verify/" + code,
		FirstName: firstName,
		Subject:   "Your account verification code",
	}

	utils.SendEmail(&newUser, &emailData)

	message := "We sent an email with a verification code to " + newUser.Email
	c.JSON(http.StatusCreated, gin.H{"status": "success", "message": message})
}

func (ac *AuthController) VerifyEmail(c *gin.Context) {

	code := c.Params.ByName("verificationCode")
	verification_code := utils.Encode(code)

	var updatedUser models.User
	result := ac.DB.First(&updatedUser, "verification_code = ?", verification_code)
	if result.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": "Invalid verification code or user doesn't exists"})
		return
	}

	if updatedUser.Verified {
		c.JSON(http.StatusConflict, gin.H{"status": "fail", "message": "User already verified"})
		return
	}

	updatedUser.Verified = true
	ac.DB.Save(&updatedUser)

	c.JSON(http.StatusOK, gin.H{"status": "success", "message": "Email verified successfully"})
}

func (ac *AuthController) Login(c *gin.Context) {
	var payload LoginInput

	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": err.Error()})
		return
	}

	var user models.User
	result := ac.DB.First(&user, "email = ?", strings.ToLower(payload.Email))
	if result.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": "Invalid email"})
		return
	}

	if !user.Verified {
		c.JSON(http.StatusForbidden, gin.H{"status": "fail", "message": "Please verify your email"})
		return
	}

	if err := utils.VerifyPassword(user.Password, payload.Password); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": "Invalid Password"})
		return
	}

	config, _ := utils.LoadConfig(".")

	// Generate Token
	token, err := utils.GenerateToken(config.TokenExpiresIn, user.ID, config.TokenSecret)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": err.Error()})
		return
	}

	c.SetCookie("token", token, config.TokenMaxAge*60, "/", "localhost", false, true)

	c.JSON(http.StatusOK, gin.H{"status": "success", "token": token})
}

func (ac *AuthController) Logout(c *gin.Context) {
	c.SetCookie("token", "", -1, "/", "localhost", false, true)
	c.JSON(http.StatusOK, gin.H{"status": "success"})
}
