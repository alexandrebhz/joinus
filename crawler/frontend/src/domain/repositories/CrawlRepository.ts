import { CrawlResult } from '../entities/CrawlResult';

export interface CrawlRepository {
  execute(siteId: string): Promise<CrawlResult>;
}

