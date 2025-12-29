package repositories

import (
	"github.com/DeepSyyy/Spatium-Backend/models"
	"gorm.io/gorm"
)

type ReactionRepository interface {
	GetReactionSummary(postInternalID int64) (map[string]int, error)
	GetReactionTypeIDByEmoji(emoji string) (int64, error)
	FindByUserAndPost(userID int64, postID int64) (*models.Reaction, error)
	Create(reaction *models.Reaction) error
	Update(reaction *models.Reaction) error
	Delete(reaction *models.Reaction) error
	CountByPost(postID int64) (int64, error)
	HasUserReacted(userID int64, postID int64) (bool, error)
}

type reactionRepository struct {
	db *gorm.DB
}

func NewReactionRepository(db *gorm.DB) ReactionRepository {
	return &reactionRepository{db}
}

func (r *reactionRepository) FindByUserAndPost(userID int64, postID int64) (*models.Reaction, error) {
	var reaction models.Reaction
	err := r.db.Where("user_internal_id = ? AND post_internal_id = ?", userID, postID).
		First(&reaction).Error
	if err != nil {
		return nil, err
	}
	return &reaction, nil
}

func (r *reactionRepository) Create(reaction *models.Reaction) error {
	return r.db.Create(reaction).Error
}

func (r *reactionRepository) Update(reaction *models.Reaction) error {
	return r.db.Save(reaction).Error
}

func (r *reactionRepository) Delete(reaction *models.Reaction) error {
	return r.db.Delete(reaction).Error
}

func (r *reactionRepository) CountByPost(postID int64) (int64, error) {
	var count int64
	err := r.db.Model(&models.Reaction{}).Where("post_internal_id = ?", postID).Count(&count).Error
	return count, err
}

func (r *reactionRepository) HasUserReacted(userID int64, postID int64) (bool, error) {
	var count int64
	err := r.db.Model(&models.Reaction{}).Where("user_internal_id = ? AND post_internal_id = ?", userID, postID).Count(&count).Error
	return count > 0, err
}

func (r *reactionRepository) GetReactionSummary(postInternalID int64) (map[string]int, error) {
	var results []struct {
		Emoji string
		Count int64
	}

	err := r.db.Table("reaction_types rt").
		Select("rt.emoji, COUNT(r.internal_id) AS count").
		Joins("JOIN reactions r ON rt.internal_id = r.reaction_type_internal_id").
		Where("r.post_internal_id = ?", postInternalID).
		Group("rt.emoji").
		Scan(&results).Error

	if err != nil {
		return nil, err
	}
	stats := make(map[string]int)
	for _, res := range results {
		stats[res.Emoji] = int(res.Count)
	}
	return stats, nil
}

func (r *reactionRepository) GetReactionTypeIDByEmoji(emoji string) (int64, error) {
	var reactionType models.ReactionType
	err := r.db.Select("internal_id").Where("emoji = ?", emoji).First(&reactionType).Error
	if err != nil {
		return 0, err
	}
	return reactionType.InternalID, nil
}
