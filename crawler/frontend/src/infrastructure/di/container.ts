import { ApiClient } from '../http/ApiClient';
import { SiteRepository } from '../../domain/repositories/SiteRepository';
import { CrawlRepository } from '../../domain/repositories/CrawlRepository';
import { GetSitesUseCase } from '../../application/use-cases/GetSitesUseCase';
import { GetSiteUseCase } from '../../application/use-cases/GetSiteUseCase';
import { CreateSiteUseCase } from '../../application/use-cases/CreateSiteUseCase';
import { UpdateSiteUseCase } from '../../application/use-cases/UpdateSiteUseCase';
import { DeleteSiteUseCase } from '../../application/use-cases/DeleteSiteUseCase';
import { ExecuteCrawlUseCase } from '../../application/use-cases/ExecuteCrawlUseCase';

class Container {
  private apiClient: ApiClient;
  private siteRepository: SiteRepository;
  private crawlRepository: CrawlRepository;

  // Use cases
  private getSitesUseCase: GetSitesUseCase;
  private getSiteUseCase: GetSiteUseCase;
  private createSiteUseCase: CreateSiteUseCase;
  private updateSiteUseCase: UpdateSiteUseCase;
  private deleteSiteUseCase: DeleteSiteUseCase;
  private executeCrawlUseCase: ExecuteCrawlUseCase;

  constructor() {
    // Infrastructure layer
    this.apiClient = new ApiClient();
    this.siteRepository = this.apiClient;
    this.crawlRepository = this.apiClient;

    // Application layer - inject dependencies
    this.getSitesUseCase = new GetSitesUseCase(this.siteRepository);
    this.getSiteUseCase = new GetSiteUseCase(this.siteRepository);
    this.createSiteUseCase = new CreateSiteUseCase(this.siteRepository);
    this.updateSiteUseCase = new UpdateSiteUseCase(this.siteRepository);
    this.deleteSiteUseCase = new DeleteSiteUseCase(this.siteRepository);
    this.executeCrawlUseCase = new ExecuteCrawlUseCase(this.crawlRepository);
  }

  // Getters for use cases
  get GetSitesUseCase(): GetSitesUseCase {
    return this.getSitesUseCase;
  }

  get GetSiteUseCase(): GetSiteUseCase {
    return this.getSiteUseCase;
  }

  get CreateSiteUseCase(): CreateSiteUseCase {
    return this.createSiteUseCase;
  }

  get UpdateSiteUseCase(): UpdateSiteUseCase {
    return this.updateSiteUseCase;
  }

  get DeleteSiteUseCase(): DeleteSiteUseCase {
    return this.deleteSiteUseCase;
  }

  get ExecuteCrawlUseCase(): ExecuteCrawlUseCase {
    return this.executeCrawlUseCase;
  }
}

// Singleton instance
export const container = new Container();

