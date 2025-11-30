import { SiteRepository } from '../../domain/repositories/SiteRepository';
import { CreateSiteInput } from '../../domain/value-objects/CreateSiteInput';
import { CrawlSite } from '../../domain/entities/CrawlSite';

export class CreateSiteUseCase {
  constructor(private siteRepository: SiteRepository) {}

  async execute(input: CreateSiteInput): Promise<CrawlSite> {
    // Domain validation could go here
    this.validateInput(input);
    return this.siteRepository.create(input);
  }

  private validateInput(input: CreateSiteInput): void {
    if (!input.name || input.name.trim().length === 0) {
      throw new Error('Site name is required');
    }
    if (!input.base_url || !this.isValidUrl(input.base_url)) {
      throw new Error('Valid base URL is required');
    }
    if (!input.backend_startup_id || input.backend_startup_id.trim().length === 0) {
      throw new Error('Backend startup ID is required');
    }
    if (!input.schedule || input.schedule.trim().length === 0) {
      throw new Error('Schedule is required');
    }
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

