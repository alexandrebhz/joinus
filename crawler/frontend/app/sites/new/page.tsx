'use client';

import { useRouter } from 'next/navigation';
import { useState } from 'react';
import { container } from '@/src/infrastructure/di/container';
import { SiteForm } from '@/src/presentation/components/forms/SiteForm';
import { CreateSiteInput } from '@/src/domain/value-objects/CreateSiteInput';
import { UpdateSiteInput } from '@/src/domain/value-objects/UpdateSiteInput';

export default function NewSitePage() {
  const router = useRouter();
  const [loading, setLoading] = useState(false);

  const handleSubmit = async (data: CreateSiteInput | UpdateSiteInput) => {
    setLoading(true);
    try {
      // Type guard to ensure it's CreateSiteInput
      if (!('name' in data && data.name) || !('base_url' in data && data.base_url)) {
        throw new Error('Invalid form data');
      }
      const createInput = data as CreateSiteInput;
      await container.CreateSiteUseCase.execute(createInput);
      router.push('/sites');
    } catch (err) {
      alert(err instanceof Error ? err.message : 'Failed to create site');
      throw err;
    } finally {
      setLoading(false);
    }
  };

  return (
    <div className="container mx-auto px-4 py-8 max-w-4xl">
      <h1 className="text-3xl font-bold mb-6">Create New Crawl Site</h1>
      <SiteForm
        onSubmit={handleSubmit}
        onCancel={() => router.back()}
        isLoading={loading}
      />
    </div>
  );
}
