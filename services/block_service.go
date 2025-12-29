package services

import (
	"errors"

	"github.com/DeepSyyy/Spatium-Backend/models"
	"github.com/DeepSyyy/Spatium-Backend/repositories"
	"github.com/google/uuid"
)

type BlockService interface {
	Block(blockerID int64, blockedUserPublicID string, reason string) (*models.Block, error)
	Unblock(blockerID int64, blockedUserPublicID string) error
	GetBlockedUsers(blockerID int64) ([]models.BlockedUserInfo, error)
	GetBlockedUserIDs(blockerID int64) ([]int64, error)
	IsBlocked(blockerID, blockedUserID int64) (bool, error)
	IsBlockedEitherWay(userID1, userID2 int64) (bool, error)
}

type blockService struct {
	repo     repositories.BlockRepository
	userRepo repositories.UserRepository
}

func NewBlockService(repo repositories.BlockRepository, userRepo repositories.UserRepository) BlockService {
	return &blockService{
		repo:     repo,
		userRepo: userRepo,
	}
}

func (s *blockService) Block(blockerID int64, blockedUserPublicID string, reason string) (*models.Block, error) {
	// Get the blocked user by public ID
	blockedUser, err := s.userRepo.GetByPublicID(blockedUserPublicID)
	if err != nil {
		return nil, errors.New("User tidak ditemukan")
	}

	// Prevent self-blocking
	if blockedUser.InternalID == blockerID {
		return nil, errors.New("Anda tidak dapat memblokir diri sendiri")
	}

	// Check if already blocked
	isBlocked, err := s.repo.IsBlocked(blockerID, blockedUser.InternalID)
	if err != nil {
		return nil, err
	}
	if isBlocked {
		return nil, errors.New("User ini sudah diblokir")
	}

	block := &models.Block{
		PublicID:      uuid.New(),
		BlockerID:     blockerID,
		BlockedUserID: blockedUser.InternalID,
		Reason:        reason,
	}

	if err := s.repo.Create(block); err != nil {
		return nil, err
	}

	return block, nil
}

func (s *blockService) Unblock(blockerID int64, blockedUserPublicID string) error {
	// Get the blocked user by public ID
	blockedUser, err := s.userRepo.GetByPublicID(blockedUserPublicID)
	if err != nil {
		return errors.New("User tidak ditemukan")
	}

	// Check if actually blocked
	isBlocked, err := s.repo.IsBlocked(blockerID, blockedUser.InternalID)
	if err != nil {
		return err
	}
	if !isBlocked {
		return errors.New("User ini tidak sedang diblokir")
	}

	return s.repo.DeleteByBlockerAndBlocked(blockerID, blockedUser.InternalID)
}

func (s *blockService) GetBlockedUsers(blockerID int64) ([]models.BlockedUserInfo, error) {
	blocks, err := s.repo.GetByBlockerID(blockerID)
	if err != nil {
		return nil, err
	}

	var blockedUsers []models.BlockedUserInfo
	for _, block := range blocks {
		user, err := s.userRepo.GetByID(block.BlockedUserID)
		if err != nil {
			continue // Skip if user not found
		}
		blockedUsers = append(blockedUsers, models.BlockedUserInfo{
			PublicID:  user.PublicID.String(),
			Alias:     user.Alias,
			BlockedAt: block.CreatedAt,
		})
	}

	return blockedUsers, nil
}

func (s *blockService) GetBlockedUserIDs(blockerID int64) ([]int64, error) {
	return s.repo.GetBlockedUserIDs(blockerID)
}

func (s *blockService) IsBlocked(blockerID, blockedUserID int64) (bool, error) {
	return s.repo.IsBlocked(blockerID, blockedUserID)
}

func (s *blockService) IsBlockedEitherWay(userID1, userID2 int64) (bool, error) {
	return s.repo.IsBlockedEitherWay(userID1, userID2)
}
