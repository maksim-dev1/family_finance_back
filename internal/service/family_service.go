package service

import (
	"time"

	"github.com/google/uuid"
	"myapp/internal/models"
	"myapp/internal/repository"
)

type FamilyService interface {
	CreateFamily(name, description, creatorUserID string) (*models.Family, error)
	GetUserFamilies(userID string) ([]*models.Family, error)
	JoinFamily(familyID, userID, role string) error
}

type familyService struct {
	familyRepo       repository.FamilyRepository
	familyMemberRepo repository.FamilyMemberRepository
}

func NewFamilyService(familyRepo repository.FamilyRepository, familyMemberRepo repository.FamilyMemberRepository) FamilyService {
	return &familyService{
		familyRepo:       familyRepo,
		familyMemberRepo: familyMemberRepo,
	}
}

func (s *familyService) CreateFamily(name, description, creatorUserID string) (*models.Family, error) {
	family := &models.Family{
		ID:          uuid.New().String(),
		Name:        name,
		Description: description,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}
	err := s.familyRepo.CreateFamily(family)
	if err != nil {
		return nil, err
	}
	member := &models.FamilyMember{
		ID:       uuid.New().String(),
		FamilyID: family.ID,
		UserID:   creatorUserID,
		Role:     "admin",
		JoinedAt: time.Now(),
	}
	err = s.familyMemberRepo.AddMember(member)
	if err != nil {
		return nil, err
	}
	return family, nil
}

func (s *familyService) GetUserFamilies(userID string) ([]*models.Family, error) {
	return s.familyRepo.GetFamiliesByUser(userID)
}

func (s *familyService) JoinFamily(familyID, userID, role string) error {
	member := &models.FamilyMember{
		ID:       uuid.New().String(),
		FamilyID: familyID,
		UserID:   userID,
		Role:     role,
		JoinedAt: time.Now(),
	}
	return s.familyMemberRepo.AddMember(member)
}
