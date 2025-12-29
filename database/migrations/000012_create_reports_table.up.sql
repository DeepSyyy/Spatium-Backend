-- Create reports table for user reports
CREATE TABLE IF NOT EXISTS reports (
    internal_id BIGSERIAL PRIMARY KEY,
    public_id UUID NOT NULL UNIQUE,
    reporter_id BIGINT NOT NULL REFERENCES users(internal_id) ON DELETE CASCADE,
    reported_user_id BIGINT REFERENCES users(internal_id) ON DELETE SET NULL,
    report_type VARCHAR(20) NOT NULL,
    target_id VARCHAR(36),
    reason VARCHAR(50) NOT NULL,
    description TEXT,
    status VARCHAR(20) DEFAULT 'pending',
    moderator_notes TEXT,
    resolved_at TIMESTAMP,
    resolved_by_id BIGINT REFERENCES users(internal_id) ON DELETE SET NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Create indexes for efficient querying
CREATE INDEX idx_reports_reporter_id ON reports(reporter_id);
CREATE INDEX idx_reports_reported_user_id ON reports(reported_user_id);
CREATE INDEX idx_reports_status ON reports(status);
CREATE INDEX idx_reports_report_type ON reports(report_type);
CREATE INDEX idx_reports_created_at ON reports(created_at);
