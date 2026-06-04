# Consolidated console-bff image (CONSOLIDATION-DESIGN Phase C).
# One image, two containers: the API server (server/) and the subscriptions
# service (subscriptions/) run as separate containers (see docker-compose),
# each with its own command. No supervisord.
FROM node:18-alpine

WORKDIR /app

# pnpm, pinned via the package.json "packageManager" field (corepack).
RUN corepack enable

COPY package.json pnpm-lock.yaml ./
# Install before NODE_ENV=production so pnpm keeps devDeps (tsc) for the build.
RUN pnpm install --frozen-lockfile --ignore-scripts

COPY . .
RUN pnpm build

ENV NODE_ENV=production
EXPOSE 8080 8081

# Default command: the consolidated API server. The subscriptions container
# overrides this with `node dist/subscriptions/index.js` (see compose).
CMD ["node", "dist/server/index.js"]
