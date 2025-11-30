import type { Metadata } from "next";
import "./globals.css";
import Link from "next/link";

export const metadata: Metadata = {
  title: "Job Crawler System",
  description: "Manage and monitor web crawling operations",
};

export default function RootLayout({
  children,
}: Readonly<{
  children: React.ReactNode;
}>) {
  return (
    <html lang="en">
      <body className="bg-gray-50">
        <nav className="bg-white shadow-sm border-b">
          <div className="container mx-auto px-4">
            <div className="flex justify-between items-center h-16">
              <Link href="/" className="text-xl font-bold text-blue-600">
                Job Crawler
              </Link>
              <div className="flex gap-4">
                <Link
                  href="/"
                  className="text-gray-600 hover:text-gray-900"
                >
                  Home
                </Link>
                <Link
                  href="/sites"
                  className="text-gray-600 hover:text-gray-900"
                >
                  Sites
                </Link>
              </div>
            </div>
          </div>
        </nav>
        {children}
      </body>
    </html>
  );
}
