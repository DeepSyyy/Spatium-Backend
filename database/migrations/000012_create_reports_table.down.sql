-- Drop indexes
DROP INDEX IF EXISTS idx_reports_reporter_id;
DROP INDEX IF EXISTS idx_reports_reported_user_id;
DROP INDEX IF EXISTS idx_reports_status;
DROP INDEX IF EXISTS idx_reports_report_type;
DROP INDEX IF EXISTS idx_reports_created_at;

-- Drop reports table
DROP TABLE IF EXISTS reports;
