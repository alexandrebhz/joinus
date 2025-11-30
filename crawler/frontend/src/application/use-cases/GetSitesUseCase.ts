import { SiteRepository } from '../../domain/repositories/SiteRepository';
import { CrawlSite } from '../../domain/entities/CrawlSite';

export class GetSitesUseCase {
  constructor(private siteRepository: SiteRepository) {}

  async execute(): Promise<CrawlSite[]> {
    return this.siteRepository.findAll();
  }
}

