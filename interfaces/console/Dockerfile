# FROM nginx
# COPY  build /usr/share/nginx/html
# RUN rm /etc/nginx/conf.d/default.conf
# COPY nginx/nginx.conf /etc/nginx/conf.d
# EXPOSE 8080
# CMD ["nginx", "-g", "daemon off;"]

FROM node:18-alpine AS deps
RUN apk add --no-cache libc6-compat
WORKDIR /

COPY package.json yarn.lock ./
RUN  yarn install --production

FROM node:18-alpine AS builder
WORKDIR /
COPY --from=deps /node_modules ./node_modules
COPY . .

ENV NEXT_TELEMETRY_DISABLED 1

RUN yarn build

FROM node:18-alpine AS runner
WORKDIR /

ENV NODE_ENV production
ENV NEXT_TELEMETRY_DISABLED 1

RUN addgroup --system --gid 1001 nodejs
RUN adduser --system --uid 1001 nextjs

COPY --from=builder --chown=nextjs:nodejs /.next ./.next
COPY --from=builder /node_modules ./node_modules
COPY --from=builder /package.json ./package.json

USER nextjs

EXPOSE 3000

ENV PORT 3000

CMD ["yarn", "start"]
