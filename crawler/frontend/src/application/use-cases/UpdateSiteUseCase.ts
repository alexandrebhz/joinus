import { SiteRepository } from '../../domain/repositories/SiteRepository';
import { UpdateSiteInput } from '../../domain/value-objects/UpdateSiteInput';
import { CrawlSite } from '../../domain/entities/CrawlSite';

export class UpdateSiteUseCase {
  constructor(private siteRepository: SiteRepository) {}

  async execute(id: string, input: UpdateSiteInput): Promise<CrawlSite> {
    // Ensure site exists
    await this.siteRepository.findById(id);
    
    // Validate if URL is being updated
    if (input.base_url && !this.isValidUrl(input.base_url)) {
      throw new Error('Invalid base URL');
    }

    return this.siteRepository.update(id, input);
  }

  private isValidUrl(url: string): boolean {
    try {
      new URL(url);
      return true;
    } catch {
      return false;
    }
  }
}

