import Link from 'next/link';

export default function Home() {
  return (
    <div className="container mx-auto px-4 py-16">
      <div className="text-center">
        <h1 className="text-4xl font-bold mb-4">Job Crawler System</h1>
        <p className="text-xl text-gray-600 mb-8">
          Manage and monitor your web crawling operations
        </p>
        <div className="flex gap-4 justify-center">
          <Link
            href="/sites"
            className="bg-blue-600 text-white px-6 py-3 rounded-lg hover:bg-blue-700"
          >
            View Sites
          </Link>
          <Link
            href="/sites/new"
            className="bg-green-600 text-white px-6 py-3 rounded-lg hover:bg-green-700"
          >
            Add New Site
          </Link>
        </div>
      </div>

      <div className="mt-16 grid md:grid-cols-3 gap-6">
        <div className="bg-white p-6 rounded-lg shadow">
          <h3 className="text-xl font-semibold mb-2">Manage Sites</h3>
          <p className="text-gray-600">
            Create, update, and configure crawl sites with custom extraction
            rules and pagination strategies.
          </p>
        </div>
        <div className="bg-white p-6 rounded-lg shadow">
          <h3 className="text-xl font-semibold mb-2">Schedule Crawls</h3>
          <p className="text-gray-600">
            Set up automated crawling schedules using cron expressions for
            each site.
          </p>
        </div>
        <div className="bg-white p-6 rounded-lg shadow">
          <h3 className="text-xl font-semibold mb-2">Monitor Jobs</h3>
          <p className="text-gray-600">
            Track crawled jobs and monitor synchronization status with the
            backend API.
          </p>
        </div>
      </div>
    </div>
  );
}
