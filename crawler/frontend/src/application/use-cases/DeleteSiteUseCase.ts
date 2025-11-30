import { SiteRepository } from '../../domain/repositories/SiteRepository';

export class DeleteSiteUseCase {
  constructor(private siteRepository: SiteRepository) {}

  async execute(id: string): Promise<void> {
    // Ensure site exists
    await this.siteRepository.findById(id);
    return this.siteRepository.delete(id);
  }
}

