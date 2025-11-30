import { useState, useCallback } from 'react';
import { container } from '@/src/infrastructure/di/container';
import { CrawlSite } from '@/src/domain/entities/CrawlSite';
import { CreateSiteInput } from '@/src/domain/value-objects/CreateSiteInput';
import { UpdateSiteInput } from '@/src/domain/value-objects/UpdateSiteInput';

export function useSites() {
  const [sites, setSites] = useState<CrawlSite[]>([]);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);

  const loadSites = useCallback(async () => {
    try {
      setLoading(true);
      setError(null);
      const data = await container.GetSitesUseCase.execute();
      setSites(data);
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Failed to load sites');
      throw err;
    } finally {
      setLoading(false);
    }
  }, []);

  const createSite = useCallback(async (input: CreateSiteInput) => {
    try {
      setLoading(true);
      setError(null);
      const newSite = await container.CreateSiteUseCase.execute(input);
      setSites((prev) => [...prev, newSite]);
      return newSite;
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Failed to create site');
      throw err;
    } finally {
      setLoading(false);
    }
  }, []);

  const updateSite = useCallback(async (id: string, input: UpdateSiteInput) => {
    try {
      setLoading(true);
      setError(null);
      const updatedSite = await container.UpdateSiteUseCase.execute(id, input);
      setSites((prev) =>
        prev.map((site) => (site.id === id ? updatedSite : site))
      );
      return updatedSite;
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Failed to update site');
      throw err;
    } finally {
      setLoading(false);
    }
  }, []);

  const deleteSite = useCallback(async (id: string) => {
    try {
      setLoading(true);
      setError(null);
      await container.DeleteSiteUseCase.execute(id);
      setSites((prev) => prev.filter((site) => site.id !== id));
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Failed to delete site');
      throw err;
    } finally {
      setLoading(false);
    }
  }, []);

  return {
    sites,
    loading,
    error,
    loadSites,
    createSite,
    updateSite,
    deleteSite,
  };
}

