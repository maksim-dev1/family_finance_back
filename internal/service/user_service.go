package service

import (
	"family_finance_back/internal/models"
	"family_finance_back/internal/repository"
)

type UserService struct {
	userRepo repository.UserRepository
}

func NewUserService(userRepo repository.UserRepository) *UserService {
	return &UserService{userRepo: userRepo}
}

func (s *UserService) GetAllUsers() ([]*models.User, error) {
	return s.userRepo.GetAllUsers()
}

func (s *UserService) GetUserByEmail(email string) (*models.User, error) {
	return s.userRepo.GetUserByEmail(email)
}

func (s *UserService) GetUserByID(id string) (*models.User, error) {
	return s.userRepo.GetUserByID(id)
}

func (s *UserService) DeleteUserByEmail(email string) error {
	return s.userRepo.DeleteUserByEmail(email)
}
