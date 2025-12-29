package repositories

import (
	"github.com/DeepSyyy/Spatium-Backend/models"
	"gorm.io/gorm"
)

type ReportRepository interface {
	Create(report *models.Report) error
	GetByID(id int64) (*models.Report, error)
	GetByPublicID(publicID string) (*models.Report, error)
	GetByReporterID(reporterID int64) ([]models.Report, error)
	GetByStatus(status models.ReportStatus) ([]models.Report, error)
	GetByTargetID(targetID string) ([]models.Report, error)
	GetAll() ([]models.Report, error)
	Update(report *models.Report) error
	Delete(id int64) error
	HasUserReported(reporterID int64, targetID string, reportType models.ReportType) (bool, error)
}

type reportRepository struct {
	db *gorm.DB
}

func NewReportRepository(db *gorm.DB) ReportRepository {
	return &reportRepository{db: db}
}

func (r *reportRepository) Create(report *models.Report) error {
	return r.db.Create(report).Error
}

func (r *reportRepository) GetByID(id int64) (*models.Report, error) {
	var report models.Report
	err := r.db.First(&report, id).Error
	if err != nil {
		return nil, err
	}
	return &report, nil
}

func (r *reportRepository) GetByPublicID(publicID string) (*models.Report, error) {
	var report models.Report
	err := r.db.Where("public_id = ?", publicID).First(&report).Error
	if err != nil {
		return nil, err
	}
	return &report, nil
}

func (r *reportRepository) GetByReporterID(reporterID int64) ([]models.Report, error) {
	var reports []models.Report
	err := r.db.Where("reporter_id = ?", reporterID).Order("created_at DESC").Find(&reports).Error
	return reports, err
}

func (r *reportRepository) GetByStatus(status models.ReportStatus) ([]models.Report, error) {
	var reports []models.Report
	err := r.db.Where("status = ?", status).Order("created_at ASC").Find(&reports).Error
	return reports, err
}

func (r *reportRepository) GetByTargetID(targetID string) ([]models.Report, error) {
	var reports []models.Report
	err := r.db.Where("target_id = ?", targetID).Find(&reports).Error
	return reports, err
}

func (r *reportRepository) GetAll() ([]models.Report, error) {
	var reports []models.Report
	err := r.db.Order("created_at DESC").Find(&reports).Error
	return reports, err
}

func (r *reportRepository) Update(report *models.Report) error {
	return r.db.Save(report).Error
}

func (r *reportRepository) Delete(id int64) error {
	return r.db.Delete(&models.Report{}, id).Error
}

func (r *reportRepository) HasUserReported(reporterID int64, targetID string, reportType models.ReportType) (bool, error) {
	var count int64
	err := r.db.Model(&models.Report{}).
		Where("reporter_id = ? AND target_id = ? AND report_type = ?", reporterID, targetID, reportType).
		Count(&count).Error
	return count > 0, err
}
