'use client';

import { useState, FormEvent } from 'react';
import { CreateSiteInput } from '@/src/domain/value-objects/CreateSiteInput';
import { UpdateSiteInput } from '@/src/domain/value-objects/UpdateSiteInput';
import { CrawlSite } from '@/src/domain/entities/CrawlSite';
import { DynamicFieldsEditor } from './DynamicFieldsEditor';

interface SiteFormProps {
  initialData?: CrawlSite;
  onSubmit: (data: CreateSiteInput | UpdateSiteInput) => Promise<void>;
  onCancel: () => void;
  isLoading?: boolean;
}

export function SiteForm({ 
  initialData, 
  onSubmit, 
  onCancel, 
  isLoading = false 
}: SiteFormProps) {
  const [formData, setFormData] = useState<CreateSiteInput>({
    name: initialData?.name || '',
    base_url: initialData?.base_url || '',
    backend_startup_id: initialData?.backend_startup_id || '',
    schedule: initialData?.schedule || '0 0 * * *',
    crawl_interval: initialData?.crawl_interval || 'daily',
    pagination_config: initialData?.pagination_config || {
      type: 'query_param',
      param_name: 'page',
      start_page: 1,
      increment: 1,
      max_pages: 10,
    },
    extraction_rules: initialData?.extraction_rules || {
      job_list_selector: '.job-item',
      job_detail_url: {
        type: 'relative',
        selector: 'a.job-link',
        attribute: 'href',
        base_url: '',
      },
      fields: {},
    },
    deduplication_key: initialData?.deduplication_key || 'url',
    request_delay: initialData?.request_delay || 2,
  });

  const handleSubmit = async (e: FormEvent) => {
    e.preventDefault();
    await onSubmit(formData);
  };

  const updateField = <K extends keyof CreateSiteInput>(
    field: K,
    value: CreateSiteInput[K]
  ) => {
    setFormData((prev) => ({ ...prev, [field]: value }));
  };

  return (
    <form onSubmit={handleSubmit} className="space-y-6">
      {/* Basic Information */}
      <div className="bg-white p-6 rounded-lg shadow">
        <h2 className="text-xl font-semibold mb-4">Basic Information</h2>
        <div className="space-y-4">
          <div>
            <label className="block text-sm font-medium mb-1">Site Name</label>
            <input
              type="text"
              required
              value={formData.name}
              onChange={(e) => updateField('name', e.target.value)}
              className="w-full px-3 py-2 border rounded-md"
            />
          </div>
          <div>
            <label className="block text-sm font-medium mb-1">Base URL</label>
            <input
              type="url"
              required
              value={formData.base_url}
              onChange={(e) => updateField('base_url', e.target.value)}
              className="w-full px-3 py-2 border rounded-md"
            />
          </div>
          <div>
            <label className="block text-sm font-medium mb-1">
              Backend Startup ID
            </label>
            <input
              type="text"
              required
              value={formData.backend_startup_id}
              onChange={(e) => updateField('backend_startup_id', e.target.value)}
              className="w-full px-3 py-2 border rounded-md"
            />
          </div>
        </div>
      </div>

      {/* Schedule */}
      <div className="bg-white p-6 rounded-lg shadow">
        <h2 className="text-xl font-semibold mb-4">Schedule</h2>
        <div className="space-y-4">
          <div>
            <label className="block text-sm font-medium mb-1">
              Cron Expression
            </label>
            <input
              type="text"
              required
              value={formData.schedule}
              onChange={(e) => updateField('schedule', e.target.value)}
              placeholder="0 0 * * *"
              className="w-full px-3 py-2 border rounded-md"
            />
            <p className="text-xs text-gray-500 mt-1">
              Example: 0 0 * * * (daily at midnight)
            </p>
          </div>
          <div>
            <label className="block text-sm font-medium mb-1">
              Crawl Interval
            </label>
            <select
              value={formData.crawl_interval}
              onChange={(e) =>
                updateField('crawl_interval', e.target.value as any)
              }
              className="w-full px-3 py-2 border rounded-md"
            >
              <option value="daily">Daily</option>
              <option value="weekly">Weekly</option>
              <option value="custom">Custom</option>
            </select>
          </div>
          <div>
            <label className="block text-sm font-medium mb-1">
              Request Delay (seconds)
            </label>
            <input
              type="number"
              min="0"
              max="60"
              value={formData.request_delay}
              onChange={(e) =>
                updateField('request_delay', parseInt(e.target.value))
              }
              className="w-full px-3 py-2 border rounded-md"
            />
          </div>
        </div>
      </div>

      {/* Pagination Config */}
      <div className="bg-white p-6 rounded-lg shadow">
        <h2 className="text-xl font-semibold mb-4">Pagination</h2>
        <div className="space-y-4">
          <div>
            <label className="block text-sm font-medium mb-1">Type</label>
            <select
              value={formData.pagination_config.type}
              onChange={(e) =>
                updateField('pagination_config', {
                  ...formData.pagination_config,
                  type: e.target.value as any,
                })
              }
              className="w-full px-3 py-2 border rounded-md"
            >
              <option value="query_param">Query Parameter</option>
              <option value="url_pattern">URL Pattern</option>
              <option value="link_follow">Link Following</option>
              <option value="api_pagination">API Pagination</option>
            </select>
          </div>
          {formData.pagination_config.type === 'query_param' && (
            <>
              <div>
                <label className="block text-sm font-medium mb-1">
                  Parameter Name
                </label>
                <input
                  type="text"
                  value={formData.pagination_config.param_name || ''}
                  onChange={(e) =>
                    updateField('pagination_config', {
                      ...formData.pagination_config,
                      param_name: e.target.value,
                    })
                  }
                  className="w-full px-3 py-2 border rounded-md"
                />
              </div>
              <div>
                <label className="block text-sm font-medium mb-1">
                  Max Pages
                </label>
                <input
                  type="number"
                  value={formData.pagination_config.max_pages || 10}
                  onChange={(e) =>
                    updateField('pagination_config', {
                      ...formData.pagination_config,
                      max_pages: parseInt(e.target.value),
                    })
                  }
                  className="w-full px-3 py-2 border rounded-md"
                />
              </div>
            </>
          )}
        </div>
      </div>

      {/* Extraction Rules */}
      <div className="bg-white p-6 rounded-lg shadow">
        <h2 className="text-xl font-semibold mb-4">Extraction Rules</h2>
        <div className="space-y-4">
          <div>
            <label className="block text-sm font-medium mb-1">
              Job List Selector
            </label>
            <input
              type="text"
              value={formData.extraction_rules.job_list_selector || ''}
              onChange={(e) =>
                updateField('extraction_rules', {
                  ...formData.extraction_rules,
                  job_list_selector: e.target.value,
                })
              }
              className="w-full px-3 py-2 border rounded-md"
              placeholder=".job-item"
            />
            <p className="text-xs text-gray-500 mt-1">
              CSS selector that matches each job listing container
            </p>
          </div>

          <div>
            <label className="block text-sm font-medium mb-1">
              Job Detail URL Selector
            </label>
            <div className="grid grid-cols-2 gap-2">
              <input
                type="text"
                value={formData.extraction_rules.job_detail_url?.selector || ''}
                onChange={(e) =>
                  updateField('extraction_rules', {
                    ...formData.extraction_rules,
                    job_detail_url: {
                      ...formData.extraction_rules.job_detail_url,
                      selector: e.target.value,
                      type: formData.extraction_rules.job_detail_url?.type || 'relative',
                      attribute: formData.extraction_rules.job_detail_url?.attribute || 'href',
                    },
                  })
                }
                className="px-3 py-2 border rounded-md"
                placeholder="a.job-link"
              />
              <select
                value={formData.extraction_rules.job_detail_url?.type || 'relative'}
                onChange={(e) =>
                  updateField('extraction_rules', {
                    ...formData.extraction_rules,
                    job_detail_url: {
                      ...formData.extraction_rules.job_detail_url,
                      type: e.target.value as any,
                    },
                  })
                }
                className="px-3 py-2 border rounded-md"
              >
                <option value="relative">Relative URL</option>
                <option value="absolute">Absolute URL</option>
                <option value="attribute">Attribute</option>
              </select>
            </div>
          </div>

          <DynamicFieldsEditor
            fields={formData.extraction_rules.fields || {}}
            onFieldsChange={(fields) =>
              updateField('extraction_rules', {
                ...formData.extraction_rules,
                fields,
              })
            }
          />
        </div>
      </div>

      <div className="flex gap-4">
        <button
          type="submit"
          disabled={isLoading}
          className="bg-blue-600 text-white px-6 py-2 rounded hover:bg-blue-700 disabled:opacity-50"
        >
          {isLoading ? 'Saving...' : initialData ? 'Update Site' : 'Create Site'}
        </button>
        <button
          type="button"
          onClick={onCancel}
          className="bg-gray-200 text-gray-800 px-6 py-2 rounded hover:bg-gray-300"
        >
          Cancel
        </button>
      </div>
    </form>
  );
}

