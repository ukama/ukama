/** @type {import('next').NextConfig} */
const nextConfig = {
  output: 'standalone',
  webpack: (config) => {
    config.cache = false;
    return config;
  },
  images: {
    remotePatterns: [
      {
        protocol: 'https',
        hostname: 'ukama-site-assets.s3.amazonaws.com',
        pathname: '**',
      },
      {
        protocol: 'https',
        hostname: 'ukama-site-assets.s3.us-east-1.amazonaws.com',
        pathname: '**',
      },
    ],
  },
};

export default nextConfig;
