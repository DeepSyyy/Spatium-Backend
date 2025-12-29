package models

import (
	"time"

	"github.com/google/uuid"
)

// ReportType defines what is being reported
type ReportType string

const (
	ReportTypePost    ReportType = "post"
	ReportTypeComment ReportType = "comment"
	ReportTypeUser    ReportType = "user"
)

// ReportStatus defines the status of a report
type ReportStatus string

const (
	ReportStatusPending   ReportStatus = "pending"
	ReportStatusReviewed  ReportStatus = "reviewed"
	ReportStatusResolved  ReportStatus = "resolved"
	ReportStatusDismissed ReportStatus = "dismissed"
)

// ReportReason defines predefined report reasons
type ReportReason string

const (
	ReportReasonHarassment    ReportReason = "harassment"
	ReportReasonHateSpeech    ReportReason = "hate_speech"
	ReportReasonSpam          ReportReason = "spam"
	ReportReasonSelfHarm      ReportReason = "self_harm"
	ReportReasonViolence      ReportReason = "violence"
	ReportReasonInappropriate ReportReason = "inappropriate"
	ReportReasonOther         ReportReason = "other"
)

// Report represents a user report against content or another user
type Report struct {
	InternalID     int64        `json:"internal_id" gorm:"primaryKey;autoIncrement"`
	PublicID       uuid.UUID    `json:"public_id" gorm:"uniqueIndex;not null"`
	ReporterID     int64        `json:"reporter_id" gorm:"not null;index"` // User who made the report
	ReportedUserID int64        `json:"reported_user_id" gorm:"index"`     // User being reported (if applicable)
	ReportType     ReportType   `json:"report_type" gorm:"type:varchar(20);not null"`
	TargetID       string       `json:"target_id" gorm:"type:varchar(36)"` // PublicID of post/comment/user being reported
	Reason         ReportReason `json:"reason" gorm:"type:varchar(50);not null"`
	Description    string       `json:"description" gorm:"type:text"`
	Status         ReportStatus `json:"status" gorm:"type:varchar(20);default:'pending'"`
	ModeratorNotes string       `json:"moderator_notes" gorm:"type:text"`
	ResolvedAt     *time.Time   `json:"resolved_at"`
	ResolvedByID   *int64       `json:"resolved_by_id"`
	CreatedAt      time.Time    `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt      time.Time    `json:"updated_at" gorm:"autoUpdateTime"`
}

// ReportRequest is the request body for creating a report
type ReportRequest struct {
	ReportType  ReportType   `json:"report_type" validate:"required"`
	TargetID    string       `json:"target_id" validate:"required"`
	Reason      ReportReason `json:"reason" validate:"required"`
	Description string       `json:"description"`
}

// ReportResponse is the response for a report
type ReportResponse struct {
	PublicID    string       `json:"public_id"`
	ReportType  ReportType   `json:"report_type"`
	TargetID    string       `json:"target_id"`
	Reason      ReportReason `json:"reason"`
	Description string       `json:"description,omitempty"`
	Status      ReportStatus `json:"status"`
	CreatedAt   time.Time    `json:"created_at"`
}

// ReportReasonOption represents a selectable report reason for the UI
type ReportReasonOption struct {
	Value       ReportReason `json:"value"`
	Label       string       `json:"label"`
	Description string       `json:"description"`
}

// GetReportReasonOptions returns all available report reasons with descriptions
func GetReportReasonOptions() []ReportReasonOption {
	return []ReportReasonOption{
		{
			Value:       ReportReasonHarassment,
			Label:       "Pelecehan",
			Description: "Perilaku yang menargetkan individu secara tidak pantas",
		},
		{
			Value:       ReportReasonHateSpeech,
			Label:       "Ujaran Kebencian",
			Description: "Konten yang mempromosikan kebencian terhadap kelompok tertentu",
		},
		{
			Value:       ReportReasonSpam,
			Label:       "Spam",
			Description: "Konten yang tidak relevan atau promosi berulang",
		},
		{
			Value:       ReportReasonSelfHarm,
			Label:       "Menyakiti Diri Sendiri",
			Description: "Konten yang mempromosikan atau mendorong tindakan menyakiti diri sendiri",
		},
		{
			Value:       ReportReasonViolence,
			Label:       "Kekerasan",
			Description: "Konten yang mengancam atau mempromosikan kekerasan",
		},
		{
			Value:       ReportReasonInappropriate,
			Label:       "Konten Tidak Pantas",
			Description: "Konten yang tidak sesuai untuk komunitas ini",
		},
		{
			Value:       ReportReasonOther,
			Label:       "Lainnya",
			Description: "Alasan lain yang tidak tercantum di atas",
		},
	}
}
