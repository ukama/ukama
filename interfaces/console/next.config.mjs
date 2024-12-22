/** @type {import('next').NextConfig} */
const nextConfig = {
  output: 'standalone',
  webpack: (config) => {
    config.cache = false;
    return config;
  },
  images: {
    domains: ['ukama-site-assets.s3.amazonaws.com'],
  },
};

export default nextConfig;
