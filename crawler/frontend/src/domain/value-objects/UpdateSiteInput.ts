import { CrawlSite, CrawlInterval, DeduplicationKey } from '../entities/CrawlSite';

export interface UpdateSiteInput {
  name?: string;
  base_url?: string;
  backend_startup_id?: string;
  active?: boolean;
  schedule?: string;
  crawl_interval?: CrawlInterval;
  pagination_config?: CrawlSite['pagination_config'];
  extraction_rules?: CrawlSite['extraction_rules'];
  deduplication_key?: DeduplicationKey;
  request_delay?: number;
  user_agent?: string;
}

