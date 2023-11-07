/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */

/** @type {import('next').NextConfig} */
const nextConfig = {
  experimental: {
    appDir: false,
  },
  env: {
    NEXT_PUBLIC_METRIC_URL: '',
    NEXT_PUBLIC_METRIC_WEBSOCKET_URL: '',
  },
  output: 'node .next/standalone/server.js',
  images: {
    remotePatterns: [
      {
        protocol: 'https',
        hostname: 'ukama-site-assets.s3.amazonaws.com',
        port: '',
        pathname: '/images/**',
      },
    ],
  },
  async redirects() {
    return [
      {
        source: '/',
        destination: '/home',
        permanent: true,
      },
    ];
  },
};

module.exports = nextConfig;
