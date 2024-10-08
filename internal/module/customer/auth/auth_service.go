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
	"time"
)

type AuthService interface {
	Register(req request.RegisterRequest) (*entity.User, *map[string]string, error)
	Login(req request.LoginRequest) (*entity.User, *map[string]string, error)
	RefreshToken(userID string) (string, error)
}

type authService struct {
	userRepo repository.UserRepository
	config   config.AppConfig
}

func NewAuthService(userRepo repository.UserRepository, config config.AppConfig) AuthService {
	return &authService{
		userRepo: userRepo,
		config:   config,
	}
}

func (u *authService) Register(req request.RegisterRequest) (*entity.User, *map[string]string, error) {
	existingUser, _ := u.userRepo.FindByEmail(req.Email)
	if existingUser != nil {
		return nil, nil, errors.New(constant.EMAIL_EXIST)
	}

	now := time.Now()
	var user = &entity.User{
		Name:          req.Name,
		Email:         req.Email,
		Address:       req.Address,
		Gender:        req.Gender,
		MaritalStatus: req.MaritalStatus,
		RoleID:        strconv.Itoa(constant.CUSTOMER),
	}

	err := user.HashPassword(req.Password)
	if err != nil {
		return nil, nil, err
	}

	user.CreatedAt = now
	user.UpdatedAt = &now
	user.CreatedBy = &req.Email
	user.UpdatedBy = &req.Email

	err = u.userRepo.Create(user)
	if err != nil {
		return nil, nil, err
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

func (u *authService) Login(req request.LoginRequest) (*entity.User, *map[string]string, error) {
	role := strconv.Itoa(constant.CUSTOMER)
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
