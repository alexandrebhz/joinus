import { SiteRepository } from '../../domain/repositories/SiteRepository';
import { CrawlSite } from '../../domain/entities/CrawlSite';

export class GetSiteUseCase {
  constructor(private siteRepository: SiteRepository) {}

  async execute(id: string): Promise<CrawlSite> {
    return this.siteRepository.findById(id);
  }
}

