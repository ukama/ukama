import bundleAnalyzer from '@next/bundle-analyzer';

const withBundleAnalyzer = bundleAnalyzer({
  enabled: process.env.ANALYZE === 'true',
});

// Collect backend origins from env so CSP connect-src stays accurate across environments.
const backendOrigins = [
  process.env.NEXT_PUBLIC_API_GW,
  process.env.API_GW_4SS,
  process.env.NEXT_PUBLIC_METRIC_URL,
  process.env.NEXT_PUBLIC_METRIC_WEBSOCKET_URL,
]
  .filter(Boolean)
  .map((u) => {
    try { return new URL(u).origin; } catch { return null; }
  })
  .filter(Boolean);

const uniqueOrigins = [...new Set(backendOrigins)];

const securityHeaders = [
  { key: 'X-Frame-Options', value: 'DENY' },
  { key: 'X-Content-Type-Options', value: 'nosniff' },
  { key: 'Referrer-Policy', value: 'strict-origin-when-cross-origin' },
  {
    key: 'Strict-Transport-Security',
    value: 'max-age=63072000; includeSubDomains; preload',
  },
  { key: 'X-XSS-Protection', value: '1; mode=block' },
  { key: 'Permissions-Policy', value: 'camera=(), microphone=(), geolocation=()' },
  {
    key: 'Content-Security-Policy',
    value: [
      "default-src 'self'",
      "script-src 'self' 'unsafe-eval' 'unsafe-inline'",
      "style-src 'self' 'unsafe-inline' https://fonts.googleapis.com",
      "font-src 'self' https://fonts.gstatic.com",
      "img-src 'self' data: blob: https://ukama-site-assets.s3.amazonaws.com https://*.mapbox.com",
      `connect-src 'self' ${uniqueOrigins.join(' ')} https://*.mapbox.com https://api.mapbox.com wss://*.mapbox.com`,
      "worker-src blob:",
      "frame-ancestors 'none'",
    ].join('; '),
  },
];

/** @type {import('next').NextConfig} */
const nextConfig = {
  output: 'standalone',
  async headers() {
    return [
      {
        source: '/(.*)',
        headers: securityHeaders,
      },
    ];
  },
  images: {
    remotePatterns: [
      {
        protocol: 'https',
        hostname: 'ukama-site-assets.s3.amazonaws.com',
        pathname: '**',
      },
    ],
    formats: ['image/avif', 'image/webp'],
    deviceSizes: [640, 768, 1024, 1280, 1536],
    imageSizes: [16, 32, 48, 64, 96, 128, 256],
  },
};

export default withBundleAnalyzer(nextConfig);
