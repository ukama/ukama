/** @type {import('next').NextConfig} */
const nextConfig = {
  output: 'standalone',
  webpack: (config) => {
    config.cache = false;
    return config;
  },
};

export default nextConfig;
