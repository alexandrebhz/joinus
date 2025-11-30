export interface CrawlSite {
  id: string;
  name: string;
  base_url: string;
  backend_startup_id: string;
  active: boolean;
  schedule: string;
  last_crawled_at?: string;
  next_crawl_at?: string;
  crawl_interval: CrawlInterval;
  pagination_config: PaginationConfig;
  extraction_rules: ExtractionRules;
  deduplication_key: DeduplicationKey;
  request_delay: number;
  user_agent: string;
  created_at: string;
  updated_at: string;
}

export type CrawlInterval = 'daily' | 'weekly' | 'custom';
export type DeduplicationKey = 'url' | 'composite' | 'external_id';

export interface PaginationConfig {
  type: PaginationType;
  param_name?: string;
  start_page?: number;
  increment?: number;
  max_pages?: number;
  next_page_selector?: string;
  url_pattern?: string;
  api_config?: APIPaginationConfig;
}

export type PaginationType = 'query_param' | 'url_pattern' | 'link_follow' | 'api_pagination';

export interface APIPaginationConfig {
  endpoint: string;
  page_param: string;
  page_size: number;
  max_pages: number;
}

export interface ExtractionRules {
  job_list_selector: string;
  job_detail_url: JobURLRule;
  fields: Record<string, FieldRule>;
}

export interface JobURLRule {
  type: 'relative' | 'absolute' | 'attribute';
  selector: string;
  attribute?: string;
  base_url?: string;
}

export interface FieldRule {
  selector: string;
  type: 'text' | 'html' | 'attribute' | 'regex';
  attribute?: string;
  regex_pattern?: string;
  required: boolean;
  default_value?: string;
  transformations?: string[];
}

