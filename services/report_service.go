package services

import (
	"errors"
	"time"

	"github.com/DeepSyyy/Spatium-Backend/models"
	"github.com/DeepSyyy/Spatium-Backend/repositories"
	"github.com/google/uuid"
)

type ReportService interface {
	Create(reporterID int64, req *models.ReportRequest) (*models.Report, error)
	GetByPublicID(publicID string) (*models.Report, error)
	GetByReporterID(reporterID int64) ([]models.Report, error)
	GetPendingReports() ([]models.Report, error)
	GetAllReports() ([]models.Report, error)
	UpdateStatus(publicID string, status models.ReportStatus, moderatorNotes string, resolvedByID int64) error
	HasUserReported(reporterID int64, targetID string, reportType models.ReportType) (bool, error)
	GetReportReasons() []models.ReportReasonOption
}

type reportService struct {
	repo     repositories.ReportRepository
	userRepo repositories.UserRepository
	postRepo repositories.PostRepository
}

func NewReportService(repo repositories.ReportRepository, userRepo repositories.UserRepository, postRepo repositories.PostRepository) ReportService {
	return &reportService{
		repo:     repo,
		userRepo: userRepo,
		postRepo: postRepo,
	}
}

func (s *reportService) Create(reporterID int64, req *models.ReportRequest) (*models.Report, error) {
	// Check if user already reported this target
	alreadyReported, err := s.repo.HasUserReported(reporterID, req.TargetID, req.ReportType)
	if err != nil {
		return nil, err
	}
	if alreadyReported {
		return nil, errors.New("Anda sudah melaporkan konten ini sebelumnya")
	}

	// Validate target exists based on report type
	var reportedUserID int64
	switch req.ReportType {
	case models.ReportTypePost:
		post, err := s.postRepo.GetPostDetail(req.TargetID)
		if err != nil {
			return nil, errors.New("Post tidak ditemukan")
		}
		reportedUserID = post.UserID
	case models.ReportTypeComment:
		// For comments, we'd need to look up the comment's user
		// For now, we'll leave reportedUserID as 0
	case models.ReportTypeUser:
		user, err := s.userRepo.GetByPublicID(req.TargetID)
		if err != nil {
			return nil, errors.New("User tidak ditemukan")
		}
		reportedUserID = user.InternalID
	default:
		return nil, errors.New("Tipe laporan tidak valid")
	}

	// Prevent self-reporting
	if reportedUserID == reporterID {
		return nil, errors.New("Anda tidak dapat melaporkan diri sendiri")
	}

	report := &models.Report{
		PublicID:       uuid.New(),
		ReporterID:     reporterID,
		ReportedUserID: reportedUserID,
		ReportType:     req.ReportType,
		TargetID:       req.TargetID,
		Reason:         req.Reason,
		Description:    req.Description,
		Status:         models.ReportStatusPending,
	}

	if err := s.repo.Create(report); err != nil {
		return nil, err
	}

	return report, nil
}

func (s *reportService) GetByPublicID(publicID string) (*models.Report, error) {
	return s.repo.GetByPublicID(publicID)
}

func (s *reportService) GetByReporterID(reporterID int64) ([]models.Report, error) {
	return s.repo.GetByReporterID(reporterID)
}

func (s *reportService) GetPendingReports() ([]models.Report, error) {
	return s.repo.GetByStatus(models.ReportStatusPending)
}

func (s *reportService) GetAllReports() ([]models.Report, error) {
	return s.repo.GetAll()
}

func (s *reportService) UpdateStatus(publicID string, status models.ReportStatus, moderatorNotes string, resolvedByID int64) error {
	report, err := s.repo.GetByPublicID(publicID)
	if err != nil {
		return err
	}

	report.Status = status
	report.ModeratorNotes = moderatorNotes

	if status == models.ReportStatusResolved || status == models.ReportStatusDismissed {
		now := time.Now()
		report.ResolvedAt = &now
		report.ResolvedByID = &resolvedByID
	}

	return s.repo.Update(report)
}

func (s *reportService) HasUserReported(reporterID int64, targetID string, reportType models.ReportType) (bool, error) {
	return s.repo.HasUserReported(reporterID, targetID, reportType)
}

func (s *reportService) GetReportReasons() []models.ReportReasonOption {
	return models.GetReportReasonOptions()
}
