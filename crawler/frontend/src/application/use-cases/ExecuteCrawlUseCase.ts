import { CrawlRepository } from '../../domain/repositories/CrawlRepository';
import { CrawlResult } from '../../domain/entities/CrawlResult';

export class ExecuteCrawlUseCase {
  constructor(private crawlRepository: CrawlRepository) {}

  async execute(siteId: string): Promise<CrawlResult> {
    if (!siteId || siteId.trim().length === 0) {
      throw new Error('Site ID is required');
    }
    return this.crawlRepository.execute(siteId);
  }
}

