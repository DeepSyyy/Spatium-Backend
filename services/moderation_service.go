package services

import (
	"github.com/DeepSyyy/Spatium-Backend/utils"
)

type ModerationService interface {
	ModerateContent(content string) (*utils.ModerationResult, error)
	SanitizeContent(content string) string
	ContainsCrisisIndicators(content string) bool
	GetCrisisSupportMessage() string
}

type moderationService struct{}

func NewModerationService() ModerationService {
	return &moderationService{}
}

func (s *moderationService) ModerateContent(content string) (*utils.ModerationResult, error) {
	return utils.ModerateContent(content)
}

func (s *moderationService) SanitizeContent(content string) string {
	return utils.SanitizeContent(content)
}

func (s *moderationService) ContainsCrisisIndicators(content string) bool {
	return utils.ContainsCrisisIndicators(content)
}

func (s *moderationService) GetCrisisSupportMessage() string {
	return utils.GetCrisisSupportMessage()
}
