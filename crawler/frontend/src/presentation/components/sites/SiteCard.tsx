import { CrawlSite } from '@/src/domain/entities/CrawlSite';
import Link from 'next/link';

interface SiteCardProps {
  site: CrawlSite;
  onToggleActive: (site: CrawlSite) => void;
  onDelete: (id: string) => void;
}

export function SiteCard({ site, onToggleActive, onDelete }: SiteCardProps) {
  return (
    <div className="border rounded-lg p-4 hover:shadow-md transition-shadow">
      <div className="flex justify-between items-start">
        <div className="flex-1">
          <div className="flex items-center gap-2 mb-2">
            <h3 className="text-xl font-semibold">{site.name}</h3>
            <span
              className={`px-2 py-1 text-xs rounded ${
                site.active
                  ? 'bg-green-100 text-green-800'
                  : 'bg-gray-100 text-gray-800'
              }`}
            >
              {site.active ? 'Active' : 'Inactive'}
            </span>
          </div>
          <p className="text-gray-600 mb-2">
            <a
              href={site.base_url}
              target="_blank"
              rel="noopener noreferrer"
              className="text-blue-600 hover:underline"
            >
              {site.base_url}
            </a>
          </p>
          <div className="text-sm text-gray-500 space-y-1">
            <p>Schedule: {site.schedule}</p>
            <p>Interval: {site.crawl_interval}</p>
            {site.last_crawled_at && (
              <p>
                Last crawled:{' '}
                {new Date(site.last_crawled_at).toLocaleString()}
              </p>
            )}
          </div>
        </div>
        <div className="flex gap-2 ml-4">
          <button
            onClick={() => onToggleActive(site)}
            className={`px-3 py-1 text-sm rounded ${
              site.active
                ? 'bg-yellow-100 text-yellow-800 hover:bg-yellow-200'
                : 'bg-green-100 text-green-800 hover:bg-green-200'
            }`}
          >
            {site.active ? 'Deactivate' : 'Activate'}
          </button>
          <Link
            href={`/sites/${site.id}/edit`}
            className="px-3 py-1 text-sm bg-blue-100 text-blue-800 rounded hover:bg-blue-200"
          >
            Edit
          </Link>
          <button
            onClick={() => onDelete(site.id)}
            className="px-3 py-1 text-sm bg-red-100 text-red-800 rounded hover:bg-red-200"
          >
            Delete
          </button>
        </div>
      </div>
    </div>
  );
}

