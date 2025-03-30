package service

import (
	"context"
	"errors"
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"

	"github.com/Ex6linz/BookSwap/backend/internal/auth/domain"
)

type AuthService struct {
	userRepo  UserRepository
	jwtSecret string
}

type UserRepository interface {
	Create(ctx context.Context, user *domain.User) error
	GetByEmail(ctx context.Context, email string) (*domain.User, error)
}

func NewAuthService(userRepo UserRepository, jwtSecret string) *AuthService {
	return &AuthService{
		userRepo:  userRepo,
		jwtSecret: jwtSecret,
	}
}

func (s *AuthService) Register(ctx context.Context, req *domain.UserRegister) (*domain.User, error) {
	existingUser, _ := s.userRepo.GetByEmail(ctx, req.Email)
	if existingUser != nil {
		return nil, domain.ErrEmailExists
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("failed to hash password: %w", err)
	}

	now := time.Now().UTC()
	user := &domain.User{
		ID:           uuid.New(),
		Name:         req.Name,
		Email:        req.Email,
		PasswordHash: string(hashedPassword),
		Location:     req.Location,
		CreatedAt:    now,
		UpdatedAt:    now,
	}

	if err := s.userRepo.Create(ctx, user); err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	// Nie zwracamy hasła
	user.PasswordHash = ""
	return user, nil
}

func (s *AuthService) Login(ctx context.Context, req *domain.UserLogin) (string, *domain.User, error) {
	user, err := s.userRepo.GetByEmail(ctx, req.Email)
	if err != nil {
		if errors.Is(err, domain.ErrUserNotFound) {
			return "", nil, domain.ErrInvalidCredentials
		}
		return "", nil, fmt.Errorf("failed to get user: %w", err)
	}

	if err := bcrypt.CompareHashAndPassword(
		[]byte(user.PasswordHash),
		[]byte(req.Password),
	); err != nil {
		return "", nil, domain.ErrInvalidCredentials
	}

	// Generowanie tokenu JWT
	token, err := generateJWT(user, s.jwtSecret)
	if err != nil {
		return "", nil, fmt.Errorf("failed to generate token: %w", err)
	}

	// Nie zwracamy hasła
	user.PasswordHash = ""
	return token, user, nil
}

func generateJWT(user *domain.User, secret string) (string, error) {
	claims := jwt.MapClaims{
		"sub":   user.ID.String(),
		"exp":   time.Now().Add(time.Hour * 24 * 7).Unix(), // 7 dni
		"name":  user.Name,
		"email": user.Email,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secret))
}
