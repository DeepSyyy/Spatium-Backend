package repositories

import (
	"github.com/DeepSyyy/Spatium-Backend/models"
	"gorm.io/gorm"
)

type BlockRepository interface {
	Create(block *models.Block) error
	GetByID(id int64) (*models.Block, error)
	GetByPublicID(publicID string) (*models.Block, error)
	GetByBlockerID(blockerID int64) ([]models.Block, error)
	GetBlockedUserIDs(blockerID int64) ([]int64, error)
	IsBlocked(blockerID, blockedUserID int64) (bool, error)
	IsBlockedEitherWay(userID1, userID2 int64) (bool, error)
	Delete(id int64) error
	DeleteByBlockerAndBlocked(blockerID, blockedUserID int64) error
}

type blockRepository struct {
	db *gorm.DB
}

func NewBlockRepository(db *gorm.DB) BlockRepository {
	return &blockRepository{db: db}
}

func (r *blockRepository) Create(block *models.Block) error {
	return r.db.Create(block).Error
}

func (r *blockRepository) GetByID(id int64) (*models.Block, error) {
	var block models.Block
	err := r.db.First(&block, id).Error
	if err != nil {
		return nil, err
	}
	return &block, nil
}

func (r *blockRepository) GetByPublicID(publicID string) (*models.Block, error) {
	var block models.Block
	err := r.db.Where("public_id = ?", publicID).First(&block).Error
	if err != nil {
		return nil, err
	}
	return &block, nil
}

func (r *blockRepository) GetByBlockerID(blockerID int64) ([]models.Block, error) {
	var blocks []models.Block
	err := r.db.Where("blocker_id = ?", blockerID).Order("created_at DESC").Find(&blocks).Error
	return blocks, err
}

func (r *blockRepository) GetBlockedUserIDs(blockerID int64) ([]int64, error) {
	var blockedIDs []int64
	err := r.db.Model(&models.Block{}).
		Where("blocker_id = ?", blockerID).
		Pluck("blocked_user_id", &blockedIDs).Error
	return blockedIDs, err
}

func (r *blockRepository) IsBlocked(blockerID, blockedUserID int64) (bool, error) {
	var count int64
	err := r.db.Model(&models.Block{}).
		Where("blocker_id = ? AND blocked_user_id = ?", blockerID, blockedUserID).
		Count(&count).Error
	return count > 0, err
}

// IsBlockedEitherWay checks if either user has blocked the other
func (r *blockRepository) IsBlockedEitherWay(userID1, userID2 int64) (bool, error) {
	var count int64
	err := r.db.Model(&models.Block{}).
		Where("(blocker_id = ? AND blocked_user_id = ?) OR (blocker_id = ? AND blocked_user_id = ?)",
			userID1, userID2, userID2, userID1).
		Count(&count).Error
	return count > 0, err
}

func (r *blockRepository) Delete(id int64) error {
	return r.db.Delete(&models.Block{}, id).Error
}

func (r *blockRepository) DeleteByBlockerAndBlocked(blockerID, blockedUserID int64) error {
	return r.db.Where("blocker_id = ? AND blocked_user_id = ?", blockerID, blockedUserID).
		Delete(&models.Block{}).Error
}
