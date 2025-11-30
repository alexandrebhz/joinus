'use client';

import { useRouter, useParams } from 'next/navigation';
import { useEffect, useState } from 'react';
import { container } from '@/src/infrastructure/di/container';
import { SiteForm } from '@/src/presentation/components/forms/SiteForm';
import { CrawlSite } from '@/src/domain/entities/CrawlSite';
import { CreateSiteInput } from '@/src/domain/value-objects/CreateSiteInput';
import { UpdateSiteInput } from '@/src/domain/value-objects/UpdateSiteInput';

export default function EditSitePage() {
  const router = useRouter();
  const params = useParams();
  const siteId = params.id as string;
  const [site, setSite] = useState<CrawlSite | null>(null);
  const [loading, setLoading] = useState(true);
  const [saving, setSaving] = useState(false);

  useEffect(() => {
    loadSite();
  }, [siteId]);

  const loadSite = async () => {
    try {
      setLoading(true);
      const data = await container.GetSiteUseCase.execute(siteId);
      setSite(data);
    } catch (err) {
      alert(err instanceof Error ? err.message : 'Failed to load site');
      router.push('/sites');
    } finally {
      setLoading(false);
    }
  };

  const handleSubmit = async (data: CreateSiteInput | UpdateSiteInput) => {
    setSaving(true);
    try {
      // Convert CreateSiteInput to UpdateSiteInput (all fields become optional)
      const updateInput: UpdateSiteInput = {
        name: 'name' in data ? data.name : undefined,
        base_url: 'base_url' in data ? data.base_url : undefined,
        backend_startup_id: 'backend_startup_id' in data ? data.backend_startup_id : undefined,
        schedule: 'schedule' in data ? data.schedule : undefined,
        crawl_interval: 'crawl_interval' in data ? data.crawl_interval : undefined,
        pagination_config: 'pagination_config' in data ? data.pagination_config : undefined,
        extraction_rules: 'extraction_rules' in data ? data.extraction_rules : undefined,
        deduplication_key: 'deduplication_key' in data ? data.deduplication_key : undefined,
        request_delay: 'request_delay' in data ? data.request_delay : undefined,
      };
      await container.UpdateSiteUseCase.execute(siteId, updateInput);
      router.push('/sites');
    } catch (err) {
      alert(err instanceof Error ? err.message : 'Failed to update site');
      throw err;
    } finally {
      setSaving(false);
    }
  };

  if (loading) {
    return (
      <div className="flex justify-center items-center p-8">
        <div className="animate-spin rounded-full h-8 w-8 border-b-2 border-blue-600"></div>
      </div>
    );
  }

  if (!site) {
    return null;
  }

  return (
    <div className="container mx-auto px-4 py-8 max-w-4xl">
      <h1 className="text-3xl font-bold mb-6">Edit Crawl Site</h1>
      <SiteForm
        initialData={site}
        onSubmit={handleSubmit}
        onCancel={() => router.back()}
        isLoading={saving}
      />
    </div>
  );
}

