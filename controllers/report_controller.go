package controllers

import (
	"github.com/DeepSyyy/Spatium-Backend/models"
	"github.com/DeepSyyy/Spatium-Backend/services"
	"github.com/DeepSyyy/Spatium-Backend/utils"
	"github.com/gofiber/fiber/v2"
)

type ReportController struct {
	reportService services.ReportService
}

func NewReportController(reportService services.ReportService) *ReportController {
	return &ReportController{reportService: reportService}
}

// CreateReport handles POST /reports
// @Summary Create a new report
// @Description Report a post, comment, or user for violating community guidelines
func (c *ReportController) CreateReport(ctx *fiber.Ctx) error {
	var req models.ReportRequest

	if err := ctx.BodyParser(&req); err != nil {
		return utils.BadRequest(ctx, "Input tidak valid", err.Error())
	}

	// Get user ID from JWT token
	userID, ok := ctx.Locals("user_id").(int64)
	if !ok {
		return utils.BadRequest(ctx, "Unauthorized", "Invalid or missing user ID")
	}

	// Validate required fields
	if req.TargetID == "" {
		return utils.BadRequest(ctx, "Target ID wajib diisi", "")
	}
	if req.ReportType == "" {
		return utils.BadRequest(ctx, "Tipe laporan wajib diisi", "")
	}
	if req.Reason == "" {
		return utils.BadRequest(ctx, "Alasan laporan wajib diisi", "")
	}

	// Validate report type
	validTypes := map[models.ReportType]bool{
		models.ReportTypePost:    true,
		models.ReportTypeComment: true,
		models.ReportTypeUser:    true,
	}
	if !validTypes[req.ReportType] {
		return utils.BadRequest(ctx, "Tipe laporan tidak valid", "")
	}

	// Validate report reason
	validReasons := map[models.ReportReason]bool{
		models.ReportReasonHarassment:    true,
		models.ReportReasonHateSpeech:    true,
		models.ReportReasonSpam:          true,
		models.ReportReasonSelfHarm:      true,
		models.ReportReasonViolence:      true,
		models.ReportReasonInappropriate: true,
		models.ReportReasonOther:         true,
	}
	if !validReasons[req.Reason] {
		return utils.BadRequest(ctx, "Alasan laporan tidak valid", "")
	}

	report, err := c.reportService.Create(userID, &req)
	if err != nil {
		return utils.BadRequest(ctx, "Gagal membuat laporan", err.Error())
	}

	resp := models.ReportResponse{
		PublicID:    report.PublicID.String(),
		ReportType:  report.ReportType,
		TargetID:    report.TargetID,
		Reason:      report.Reason,
		Description: report.Description,
		Status:      report.Status,
		CreatedAt:   report.CreatedAt,
	}

	return utils.Success(ctx, "Laporan berhasil dikirim. Terima kasih telah membantu menjaga komunitas.", resp)
}

// GetMyReports handles GET /reports/me
// @Summary Get current user's reports
// @Description Get all reports submitted by the current user
func (c *ReportController) GetMyReports(ctx *fiber.Ctx) error {
	userID, ok := ctx.Locals("user_id").(int64)
	if !ok {
		return utils.BadRequest(ctx, "Unauthorized", "Invalid or missing user ID")
	}

	reports, err := c.reportService.GetByReporterID(userID)
	if err != nil {
		return utils.BadRequest(ctx, "Gagal mengambil laporan", err.Error())
	}

	var resp []models.ReportResponse
	for _, report := range reports {
		resp = append(resp, models.ReportResponse{
			PublicID:    report.PublicID.String(),
			ReportType:  report.ReportType,
			TargetID:    report.TargetID,
			Reason:      report.Reason,
			Description: report.Description,
			Status:      report.Status,
			CreatedAt:   report.CreatedAt,
		})
	}

	return utils.Success(ctx, "Laporan berhasil diambil", resp)
}

// GetReportReasons handles GET /reports/reasons
// @Summary Get available report reasons
// @Description Get list of all available report reasons with descriptions
func (c *ReportController) GetReportReasons(ctx *fiber.Ctx) error {
	reasons := c.reportService.GetReportReasons()
	return utils.Success(ctx, "Alasan laporan berhasil diambil", reasons)
}

// GetPendingReports handles GET /admin/reports/pending
// @Summary Get pending reports (Admin only)
// @Description Get all pending reports for moderation
func (c *ReportController) GetPendingReports(ctx *fiber.Ctx) error {
	// TODO: Add admin role check
	reports, err := c.reportService.GetPendingReports()
	if err != nil {
		return utils.BadRequest(ctx, "Gagal mengambil laporan", err.Error())
	}

	return utils.Success(ctx, "Laporan pending berhasil diambil", reports)
}

// UpdateReportStatus handles PUT /admin/reports/:report_id/status
// @Summary Update report status (Admin only)
// @Description Update the status of a report
func (c *ReportController) UpdateReportStatus(ctx *fiber.Ctx) error {
	reportID := ctx.Params("report_id")
	if reportID == "" {
		return utils.BadRequest(ctx, "Report ID wajib diisi", "")
	}

	var req struct {
		Status         string `json:"status"`
		ModeratorNotes string `json:"moderator_notes"`
	}

	if err := ctx.BodyParser(&req); err != nil {
		return utils.BadRequest(ctx, "Input tidak valid", err.Error())
	}

	userID, ok := ctx.Locals("user_id").(int64)
	if !ok {
		return utils.BadRequest(ctx, "Unauthorized", "Invalid or missing user ID")
	}

	status := models.ReportStatus(req.Status)
	validStatuses := map[models.ReportStatus]bool{
		models.ReportStatusPending:   true,
		models.ReportStatusReviewed:  true,
		models.ReportStatusResolved:  true,
		models.ReportStatusDismissed: true,
	}
	if !validStatuses[status] {
		return utils.BadRequest(ctx, "Status tidak valid", "")
	}

	err := c.reportService.UpdateStatus(reportID, status, req.ModeratorNotes, userID)
	if err != nil {
		return utils.BadRequest(ctx, "Gagal mengubah status laporan", err.Error())
	}

	return utils.Success(ctx, "Status laporan berhasil diubah", nil)
}
