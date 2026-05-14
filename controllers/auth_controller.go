package controllers

import (
	"net/http"

	"demodiqit_api/config"
	"demodiqit_api/helpers/crypt"
	"demodiqit_api/helpers/respond"
	"demodiqit_api/middleware"
	"demodiqit_api/models"
	"demodiqit_api/request"

	"github.com/gin-gonic/gin"
)

type AuthController struct {
	Config *config.Config
}

func NewAuthController(cfg *config.Config) *AuthController {
	return &AuthController{
		Config: cfg,
	}
}

func (ac *AuthController) Login(c *gin.Context) {
	var req request.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, respond.ErrorRespond{
			Message: "Invalid request payload",
			Code:    "AUTH-001",
		})
		return
	}

	var user models.User
	result := config.DB.Where("email = ? OR username = ?", req.Email, req.Email).First(&user)
	if result.Error != nil {
		c.JSON(http.StatusUnauthorized, respond.ErrorRespond{
			Message: "Invalid email/username or password",
			Code:    "AUTH-002",
		})
		return
	}

	// Compare password
	if !crypt.CheckPasswordHash(req.Password, user.Password) {
		c.JSON(http.StatusUnauthorized, respond.ErrorRespond{
			Message: "Invalid email/username or password",
			Code:    "AUTH-002",
		})
		return
	}

	// Generate JWT Token
	token, err := crypt.GenerateJWT(user.ID, user.Username, user.Email, user.Roles, ac.Config.JWTSecret, ac.Config.JWTExpirationDays)
	if err != nil {
		c.JSON(http.StatusInternalServerError, respond.ErrorRespond{
			Message: "Failed to generate token",
			Code:    "AUTH-003",
		})
		return
	}

	// Save token to user table
	user.UserToken = token
	if err := config.DB.Save(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, respond.ErrorRespond{
			Message: "Failed to save user token",
			Code:    "AUTH-004",
		})
		return
	}

	rsp := respond.SuccessRespond{
		Message: "Login successfully!",
		Data: request.LoginResponse{
			ID:       user.ID,
			Username: user.Username,
			Email:    user.Email,
			FullName: user.FullName,
			Roles:    user.Roles,
			Token:    token,
		},
	}

	c.JSON(http.StatusOK, rsp)
}

// Me returns the profile of the currently authenticated user.
func (ac *AuthController) Me(c *gin.Context) {
	val, exists := c.Get(middleware.CurrentUserKey)
	if !exists {
		c.JSON(http.StatusUnauthorized, respond.ErrorRespond{
			Code:    "AUTH-005",
			Message: "Unauthorized",
		})
		return
	}

	claims, ok := val.(*crypt.JWTClaim)
	if !ok {
		c.JSON(http.StatusUnauthorized, respond.ErrorRespond{
			Code:    "AUTH-006",
			Message: "Invalid token claims",
		})
		return
	}

	c.JSON(http.StatusOK, respond.SuccessRespond{
		Message: "OK",
		Data: gin.H{
			"user_id":  claims.UserID,
			"username": claims.Username,
			"email":    claims.Email,
			"roles":    claims.Roles,
		},
	})
}
