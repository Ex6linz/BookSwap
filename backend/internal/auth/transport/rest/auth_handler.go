package rest

import (
	"errors"
	"net/http"
	_ "time"

	"github.com/gin-gonic/gin"
	_ "github.com/google/uuid"

	"github.com/Ex6linz/BookSwap/backend/internal/auth/domain"
	"github.com/Ex6linz/BookSwap/backend/internal/auth/service"
)

type AuthHandler struct {
	authService *service.AuthService
}

func NewAuthHandler(authService *service.AuthService) *AuthHandler {
	return &AuthHandler{authService: authService}
}

// @Summary Rejestracja nowego użytkownika
// @Accept json
// @Produce json
// @Param input body domain.UserRegister true "Dane rejestracji"
// @Success 201 {object} domain.User
// @Failure 400 {object} ErrorResponse
// @Failure 409 {object} ErrorResponse
// @Router /auth/register [post]
func (h *AuthHandler) Register(c *gin.Context) {
	var req domain.UserRegister
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Code:    "invalid-request",
			Message: "Nieprawidłowy format żądania",
		})
		return
	}

	user, err := h.authService.Register(c.Request.Context(), &req)
	if err != nil {
		handleAuthError(c, err)
		return
	}

	c.JSON(http.StatusCreated, user)
}

// @Summary Logowanie użytkownika
// @Accept json
// @Produce json
// @Param input body domain.UserLogin true "Dane logowania"
// @Success 200 {object} LoginResponse
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Router /auth/login [post]
func (h *AuthHandler) Login(c *gin.Context) {
	var req domain.UserLogin
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Code:    "invalid-request",
			Message: "Nieprawidłowy format żądania",
		})
		return
	}

	token, user, err := h.authService.Login(c.Request.Context(), &req)
	if err != nil {
		handleAuthError(c, err)
		return
	}

	c.JSON(http.StatusOK, LoginResponse{
		User:  user,
		Token: token,
	})
}

func handleAuthError(c *gin.Context, err error) {
	switch {
	case errors.Is(err, domain.ErrEmailExists):
		c.JSON(http.StatusConflict, ErrorResponse{
			Code:    "email-exists",
			Message: "Użytkownik z tym adresem email już istnieje",
		})
	case errors.Is(err, domain.ErrInvalidCredentials):
		c.JSON(http.StatusUnauthorized, ErrorResponse{
			Code:    "invalid-credentials",
			Message: "Nieprawidłowy email lub hasło",
		})
	default:
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Code:    "internal-error",
			Message: "Wystąpił błąd wewnętrzny",
		})
	}
}

type ErrorResponse struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

type LoginResponse struct {
	User  *domain.User `json:"user"`
	Token string       `json:"token"`
}
