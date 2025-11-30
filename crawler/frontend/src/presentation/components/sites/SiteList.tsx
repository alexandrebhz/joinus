'use client';

import { useEffect, useState } from 'react';
import { container } from '@/src/infrastructure/di/container';
import { CrawlSite } from '@/src/domain/entities/CrawlSite';
import { SiteCard } from './SiteCard';
import Link from 'next/link';

export function SiteList() {
  const [sites, setSites] = useState<CrawlSite[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    loadSites();
  }, []);

  const loadSites = async () => {
    try {
      setLoading(true);
      const sites = await container.GetSitesUseCase.execute();
      setSites(sites);
      setError(null);
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Failed to load sites');
    } finally {
      setLoading(false);
    }
  };

  const handleDelete = async (id: string) => {
    if (!confirm('Are you sure you want to delete this site?')) {
      return;
    }

    try {
      await container.DeleteSiteUseCase.execute(id);
      await loadSites();
    } catch (err) {
      alert(err instanceof Error ? err.message : 'Failed to delete site');
    }
  };

  const handleToggleActive = async (site: CrawlSite) => {
    try {
      await container.UpdateSiteUseCase.execute(site.id, { 
        active: !site.active 
      });
      await loadSites();
    } catch (err) {
      alert(err instanceof Error ? err.message : 'Failed to update site');
    }
  };

  if (loading) {
    return (
      <div className="flex justify-center items-center p-8">
        <div className="animate-spin rounded-full h-8 w-8 border-b-2 border-blue-600"></div>
      </div>
    );
  }

  if (error) {
    return (
      <div className="bg-red-50 border border-red-200 text-red-700 px-4 py-3 rounded">
        {error}
      </div>
    );
  }

  return (
    <div className="space-y-4">
      <div className="flex justify-between items-center">
        <h2 className="text-2xl font-bold">Crawl Sites</h2>
        <Link
          href="/sites/new"
          className="bg-blue-600 text-white px-4 py-2 rounded hover:bg-blue-700"
        >
          Add New Site
        </Link>
      </div>

      {sites.length === 0 ? (
        <div className="text-center py-8 text-gray-500">
          No sites configured. Create your first crawl site to get started.
        </div>
      ) : (
        <div className="grid gap-4">
          {sites.map((site) => (
            <SiteCard
              key={site.id}
              site={site}
              onToggleActive={handleToggleActive}
              onDelete={handleDelete}
            />
          ))}
        </div>
      )}
    </div>
  );
}

