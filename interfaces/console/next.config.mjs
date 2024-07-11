/** @type {import('next').NextConfig} */
const nextConfig = {
  output: 'standalone',
  env: {
    NEXT_PUBLIC_NODE_ENV: process.env.NEXT_PUBLIC_NODE_ENV,
    NEXT_PUBLIC_API_GW: process.env.NEXT_PUBLIC_API_GW,
    NEXT_PUBLIC_APP_URL: process.env.NEXT_PUBLIC_APP_URL,
    NEXT_PUBLIC_AUTH_APP_URL: process.env.NEXT_PUBLIC_AUTH_APP_URL,
    NEXT_PUBLIC_METRIC_URL: process.env.NEXT_PUBLIC_METRIC_URL,
    NEXT_PUBLIC_METRIC_WEBSOCKET_URL: process.env.NEXT_PUBLIC_METRIC_WEBSOCKET_URL,
    NEXT_PUBLIC_MAP_BOX_TOKEN: process.env.NEXT_PUBLIC_MAP_BOX_TOKEN,
    NEXT_PUBLIC_API_GW_4SS: process.env.NEXT_PUBLIC_API_GW_4SS,
  }
};

export default nextConfig;
