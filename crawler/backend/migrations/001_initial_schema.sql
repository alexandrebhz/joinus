-- Create crawl_sites table
CREATE TABLE IF NOT EXISTS crawl_sites (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(255) NOT NULL,
    base_url TEXT NOT NULL,
    backend_startup_id UUID NOT NULL,
    active BOOLEAN DEFAULT TRUE,
    schedule VARCHAR(100) NOT NULL,
    last_crawled_at TIMESTAMP,
    next_crawl_at TIMESTAMP,
    crawl_interval VARCHAR(50) NOT NULL,
    pagination_config JSONB NOT NULL,
    extraction_rules JSONB NOT NULL,
    deduplication_key VARCHAR(50) DEFAULT 'url',
    request_delay INTEGER DEFAULT 2,
    user_agent TEXT,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMP
);

-- Create crawled_jobs table
CREATE TABLE IF NOT EXISTS crawled_jobs (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    site_id UUID NOT NULL REFERENCES crawl_sites(id) ON DELETE CASCADE,
    external_id VARCHAR(255),
    detail_url TEXT NOT NULL,
    title VARCHAR(500) NOT NULL,
    description TEXT,
    requirements TEXT,
    company VARCHAR(255),
    location VARCHAR(255),
    city VARCHAR(100),
    country VARCHAR(100),
    job_type VARCHAR(50),
    location_type VARCHAR(50),
    salary_min INTEGER,
    salary_max INTEGER,
    currency VARCHAR(3),
    application_url TEXT,
    application_email VARCHAR(255),
    expires_at TIMESTAMP,
    raw_html TEXT,
    deduplication_hash VARCHAR(255),
    synced BOOLEAN DEFAULT FALSE,
    synced_at TIMESTAMP,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMP
);

-- Create indexes
CREATE INDEX IF NOT EXISTS idx_crawl_sites_active ON crawl_sites(active);
CREATE INDEX IF NOT EXISTS idx_crawled_jobs_site_id ON crawled_jobs(site_id);
CREATE INDEX IF NOT EXISTS idx_crawled_jobs_detail_url ON crawled_jobs(detail_url);
CREATE INDEX IF NOT EXISTS idx_crawled_jobs_external_id ON crawled_jobs(external_id);
CREATE INDEX IF NOT EXISTS idx_crawled_jobs_deduplication_hash ON crawled_jobs(deduplication_hash);
CREATE INDEX IF NOT EXISTS idx_crawled_jobs_synced ON crawled_jobs(synced);

