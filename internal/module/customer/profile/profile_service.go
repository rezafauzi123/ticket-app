package profile

import (
	"errors"
	"ticket-app/internal/entity"
	"ticket-app/internal/module/customer/profile/request"
	"ticket-app/internal/repository"
	"ticket-app/pkg/constant"
)

type ProfileService interface {
	GetMe(userID string) (*entity.User, error)
	UpdateUser(userID string, req request.UpdateProfileRequest) (*entity.User, error)
	DeleteUser(userID string, deletedBy string) error
}

type profileService struct {
	userRepo repository.UserRepository
}

func NewProfileService(userRepo repository.UserRepository) ProfileService {
	return &profileService{userRepo: userRepo}
}

func (u *profileService) GetMe(userID string) (*entity.User, error) {
	data, err := u.userRepo.GetMe(userID)
	if data == nil {
		return nil, errors.New(constant.DATA_NOT_FOUND)
	}

	if err != nil {
		return nil, err
	}

	return data, nil
}

func (u *profileService) UpdateUser(userID string, req request.UpdateProfileRequest) (*entity.User, error) {
	existingUser, err := u.userRepo.GetMe(userID)
	if err != nil {
		return nil, err
	}

	if existingUser == nil {
		return nil, errors.New(constant.DATA_NOT_FOUND)
	}

	existingUser.Name = req.Name
	existingUser.Email = req.Email
	existingUser.Address = req.Address
	existingUser.Gender = req.Gender
	existingUser.MaritalStatus = req.MaritalStatus
	data, err := u.userRepo.UpdateUser(userID, *existingUser)
	if err != nil {
		return nil, err
	}

	return data, nil
}

func (u *profileService) DeleteUser(userID string, deletedBy string) error {
	err := u.userRepo.DeleteUser(userID, deletedBy)
	if err != nil {
		return err
	}

	return nil
}
