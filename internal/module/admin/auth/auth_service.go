package auth

import (
	"errors"
	"strconv"
	"ticket-app/config"
	"ticket-app/internal/entity"
	"ticket-app/internal/module/customer/auth/request"
	"ticket-app/internal/repository"
	"ticket-app/pkg/constant"
	"ticket-app/pkg/jwt"
)

type AuthService interface {
	Login(req request.LoginRequest) (*entity.User, *map[string]string, error)
	RefreshToken(userID string) (string, error)
}

type authService struct {
	config   config.AppConfig
	userRepo repository.UserRepository
}

func NewAuthService(userRepo repository.UserRepository, confing config.AppConfig) AuthService {
	return &authService{
		userRepo: userRepo,
		config:   confing,
	}
}

func (u *authService) Login(req request.LoginRequest) (*entity.User, *map[string]string, error) {
	role := strconv.Itoa(constant.ADMIN)
	user, err := u.userRepo.FindByEmail(req.Email)
	if err != nil {
		return nil, nil, errors.New(constant.DATA_NOT_FOUND)
	}

	if user.RoleID != role {
		return nil, nil, errors.New(constant.DATA_NOT_FOUND)
	}

	if err := user.CheckPassword(req.Password); err != nil {
		return nil, nil, errors.New(constant.INVALID_PASSWORD)
	}

	token, refreshToken, err := jwt.GenerateTokens(user.ID, user.RoleID)
	if err != nil {
		return nil, nil, err
	}

	return user, &map[string]string{
		"access_token":  token,
		"refresh_token": refreshToken,
	}, nil
}

func (u *authService) RefreshToken(userID string) (string, error) {
	user, err := u.userRepo.GetMe(userID)
	if err != nil {
		return "", errors.New(constant.DATA_NOT_FOUND)
	}

	token, err := jwt.GenerateAccessToken(user.ID, user.RoleID)
	if err != nil {
		return "", err
	}

	return token, nil
}
