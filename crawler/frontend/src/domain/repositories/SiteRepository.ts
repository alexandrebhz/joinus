import { CrawlSite } from '../entities/CrawlSite';
import { CreateSiteInput } from '../value-objects/CreateSiteInput';
import { UpdateSiteInput } from '../value-objects/UpdateSiteInput';

export interface SiteRepository {
  findAll(): Promise<CrawlSite[]>;
  findById(id: string): Promise<CrawlSite>;
  create(input: CreateSiteInput): Promise<CrawlSite>;
  update(id: string, input: UpdateSiteInput): Promise<CrawlSite>;
  delete(id: string): Promise<void>;
}

