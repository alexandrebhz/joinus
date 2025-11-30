-- Add custom_fields column to crawled_jobs table
ALTER TABLE crawled_jobs ADD COLUMN IF NOT EXISTS custom_fields JSONB DEFAULT '{}'::jsonb;

-- Create index for custom fields queries (optional, useful if you want to query by custom fields)
CREATE INDEX IF NOT EXISTS idx_crawled_jobs_custom_fields ON crawled_jobs USING GIN (custom_fields);

