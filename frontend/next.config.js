/** @type {import('next').NextConfig} */
const nextConfig = {
  reactStrictMode: true,
  output: process.env.NODE_ENV === 'production' ? 'standalone' : undefined,
  
  // Memory optimizations
  compress: true,
  
  // Optimize production builds
  productionBrowserSourceMaps: false,
  
  // Reduce memory usage during build
  experimental: {
    optimizeCss: true,
  },
  
  images: {
    // Use remotePatterns instead of deprecated domains
    remotePatterns: [
      {
        protocol: 'http',
        hostname: 'localhost',
        port: '9000',
        pathname: '/**',
      },
    ],
    // Optimize image loading
    formats: ['image/avif', 'image/webp'],
    minimumCacheTTL: 60,
  },
  
  // Turbopack config (Next.js 16 uses Turbopack by default)
  // Empty config allows webpack to be used when --webpack flag is passed
  turbopack: {},
  
  // Webpack optimizations (used when --webpack flag is passed)
  webpack: (config, { isServer }) => {
    if (!isServer) {
      // Reduce client bundle size
      config.optimization = {
        ...config.optimization,
        moduleIds: 'deterministic',
        runtimeChunk: 'single',
        splitChunks: {
          chunks: 'all',
          cacheGroups: {
            default: false,
            vendors: false,
            // Vendor chunk
            vendor: {
              name: 'vendor',
              chunks: 'all',
              test: /node_modules/,
              priority: 20,
            },
            // Common chunk
            common: {
              name: 'common',
              minChunks: 2,
              chunks: 'all',
              priority: 10,
              reuseExistingChunk: true,
              enforce: true,
            },
          },
        },
      }
    }
    return config
  },
  
  // Note: NEXT_PUBLIC_* variables are automatically exposed by Next.js
  // No need to manually set them in the env section
}

module.exports = nextConfig

