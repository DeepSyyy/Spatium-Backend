package services

import (
	"errors"
	"time"

	"github.com/DeepSyyy/Spatium-Backend/models"
	"github.com/DeepSyyy/Spatium-Backend/repositories"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ReactionService interface {
	ReactToPost(userID int64, postID int64, reactionTypeID int64, emoji string) error
	GetReactionSummary(postInternalID int64) (map[string]int, error)
	GetReactionTypeIDByEmoji(emoji string) (int64, error)
}

type reactionService struct {
	reactionRepo repositories.ReactionRepository
}

func NewReactionService(reactionRepo repositories.ReactionRepository) ReactionService {
	return &reactionService{reactionRepo}
}

func (s *reactionService) ReactToPost(userID int64, postID int64, reactionTypeID int64, emoji string) error {
	reaction, err := s.reactionRepo.FindByUserAndPost(userID, postID)

	// Jika belum ada reaction → create baru
	if errors.Is(err, gorm.ErrRecordNotFound) {
		newReaction := models.Reaction{
			PublicID:               uuid.New(),
			UserID:                 userID,
			PostID:                 postID,
			ReactionTypeInternalID: reactionTypeID,
			Emoji:                  emoji,
			CreatedAt:              time.Now(),
		}
		return s.reactionRepo.Create(&newReaction)
	}

	// Kalau error lain → langsung return
	if err != nil {
		return err
	}

	// Jika sudah ada → update reaction lama
	reaction.ReactionTypeInternalID = reactionTypeID
	reaction.Emoji = emoji
	return s.reactionRepo.Update(reaction)
}

func (s *reactionService) GetReactionSummary(postInternalID int64) (map[string]int, error) {
	return s.reactionRepo.GetReactionSummary(postInternalID)
}

func (s *reactionService) GetReactionTypeIDByEmoji(emoji string) (int64, error) {
	return s.reactionRepo.GetReactionTypeIDByEmoji(emoji)
}
