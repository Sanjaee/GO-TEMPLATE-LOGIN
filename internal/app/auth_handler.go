package app

import (
	"net/http"
	"strings"

	"yourapp/internal/service"
	"yourapp/internal/util"

	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	authService service.AuthService
	jwtSecret   string
}

func NewAuthHandler(authService service.AuthService, jwtSecret string) *AuthHandler {
	return &AuthHandler{
		authService: authService,
		jwtSecret:   jwtSecret,
	}
}

// Register handles user registration
// POST /api/v1/auth/register
func (h *AuthHandler) Register(c *gin.Context) {
	var req service.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		util.BadRequest(c, err.Error())
		return
	}

	resp, err := h.authService.Register(req)
	if err != nil {
		util.ErrorResponse(c, http.StatusBadRequest, err.Error(), nil)
		return
	}

	util.SuccessResponse(c, http.StatusCreated, resp.Message, resp)
}

// Login handles user login
// POST /api/v1/auth/login
func (h *AuthHandler) Login(c *gin.Context) {
	var req service.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		util.BadRequest(c, err.Error())
		return
	}

	resp, err := h.authService.Login(req)
	if err != nil {
		if strings.Contains(err.Error(), "not verified") {
			// Return special response for unverified email with email in data
			util.ErrorResponse(c, http.StatusUnauthorized, err.Error(), gin.H{
				"email":                 req.Email,
				"requires_verification": true,
				"message":               "OTP telah dikirim ke email Anda. Silakan verifikasi email untuk melanjutkan.",
			})
			return
		}
		util.Unauthorized(c, err.Error())
		return
	}

	util.SuccessResponse(c, http.StatusOK, "Login successful", resp)
}

// VerifyOTP handles OTP verification
// POST /api/v1/auth/verify-otp
func (h *AuthHandler) VerifyOTP(c *gin.Context) {
	var req struct {
		Email   string `json:"email" binding:"required,email"`
		OTPCode string `json:"otp_code" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		util.BadRequest(c, err.Error())
		return
	}

	resp, err := h.authService.VerifyOTP(req.Email, req.OTPCode)
	if err != nil {
		util.ErrorResponse(c, http.StatusBadRequest, err.Error(), nil)
		return
	}

	util.SuccessResponse(c, http.StatusOK, "OTP verified successfully", resp)
}

// ResendOTP handles OTP resend
// POST /api/v1/auth/resend-otp
func (h *AuthHandler) ResendOTP(c *gin.Context) {
	var req struct {
		Email string `json:"email" binding:"required,email"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		util.BadRequest(c, err.Error())
		return
	}

	if err := h.authService.ResendOTP(req.Email); err != nil {
		util.ErrorResponse(c, http.StatusBadRequest, err.Error(), nil)
		return
	}

	util.SuccessResponse(c, http.StatusOK, "OTP sent successfully", nil)
}

// GoogleOAuth handles Google OAuth login
// POST /api/v1/auth/google-oauth
func (h *AuthHandler) GoogleOAuth(c *gin.Context) {
	var req service.GoogleOAuthRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		util.BadRequest(c, err.Error())
		return
	}

	resp, err := h.authService.GoogleOAuth(req)
	if err != nil {
		util.ErrorResponse(c, http.StatusBadRequest, err.Error(), nil)
		return
	}

	util.SuccessResponse(c, http.StatusOK, "Google OAuth successful", resp)
}

// RefreshToken handles token refresh
// POST /api/v1/auth/refresh-token
func (h *AuthHandler) RefreshToken(c *gin.Context) {
	var req struct {
		RefreshToken string `json:"refresh_token" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		util.BadRequest(c, err.Error())
		return
	}

	resp, err := h.authService.RefreshToken(req.RefreshToken)
	if err != nil {
		util.Unauthorized(c, err.Error())
		return
	}

	util.SuccessResponse(c, http.StatusOK, "Token refreshed successfully", resp)
}

// RequestResetPassword handles password reset request
// POST /api/v1/auth/forgot-password
func (h *AuthHandler) RequestResetPassword(c *gin.Context) {
	var req struct {
		Email string `json:"email" binding:"required,email"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		util.BadRequest(c, err.Error())
		return
	}

	if err := h.authService.RequestResetPassword(req.Email); err != nil {
		// Don't reveal if user exists or not
		util.SuccessResponse(c, http.StatusOK, "If email exists, reset link has been sent", nil)
		return
	}

	util.SuccessResponse(c, http.StatusOK, "If email exists, reset link has been sent", nil)
}

// VerifyResetPassword handles password reset with OTP verification
// POST /api/v1/auth/verify-reset-password
func (h *AuthHandler) VerifyResetPassword(c *gin.Context) {
	var req struct {
		Email       string `json:"email" binding:"required,email"`
		OTPCode     string `json:"otp_code" binding:"required"`
		NewPassword string `json:"new_password" binding:"required,min=8"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		util.BadRequest(c, err.Error())
		return
	}

	resp, err := h.authService.VerifyResetPassword(req.Email, req.OTPCode, req.NewPassword)
	if err != nil {
		util.ErrorResponse(c, http.StatusBadRequest, err.Error(), nil)
		return
	}

	util.SuccessResponse(c, http.StatusOK, "Password reset successfully", resp)
}

// ResetPassword handles password reset with token
// POST /api/v1/auth/reset-password
func (h *AuthHandler) ResetPassword(c *gin.Context) {
	var req struct {
		Token       string `json:"token" binding:"required"`
		NewPassword string `json:"newPassword" binding:"required,min=8"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		util.BadRequest(c, err.Error())
		return
	}

	resp, err := h.authService.ResetPassword(req.Token, req.NewPassword)
	if err != nil {
		util.ErrorResponse(c, http.StatusBadRequest, err.Error(), nil)
		return
	}

	util.SuccessResponse(c, http.StatusOK, "Password reset successfully", resp)
}

// VerifyEmail handles email verification
// POST /api/v1/auth/verify-email
func (h *AuthHandler) VerifyEmail(c *gin.Context) {
	var req struct {
		Token string `json:"token" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		util.BadRequest(c, err.Error())
		return
	}

	resp, err := h.authService.VerifyEmail(req.Token)
	if err != nil {
		util.ErrorResponse(c, http.StatusBadRequest, err.Error(), nil)
		return
	}

	util.SuccessResponse(c, http.StatusOK, "Email verified successfully", resp)
}

// GetMe handles getting current user info
// GET /api/v1/auth/me
func (h *AuthHandler) GetMe(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		util.Unauthorized(c, "User not authenticated")
		return
	}

	user, err := h.authService.GetMe(userID.(string))
	if err != nil {
		util.NotFound(c, "User not found")
		return
	}

	util.SuccessResponse(c, http.StatusOK, "User retrieved successfully", gin.H{"user": user})
}

// AuthMiddleware validates JWT token
func (h *AuthHandler) AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			util.Unauthorized(c, "Authorization header required")
			c.Abort()
			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			util.Unauthorized(c, "Invalid authorization header format")
			c.Abort()
			return
		}

		token := parts[1]
		claims, err := util.ValidateToken(token, h.jwtSecret)
		if err != nil {
			util.Unauthorized(c, "Invalid or expired token")
			c.Abort()
			return
		}

		c.Set("userID", claims.UserID)
		c.Set("email", claims.Email)
		c.Set("userType", claims.UserType)
		c.Next()
	}
}
