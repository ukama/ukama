FROM node:lts-alpine AS base
WORKDIR /
COPY package*.json yarn.lock ./
RUN yarn install --frozen-lockfile --ignore-scripts

FROM base AS build
WORKDIR /
ENV NODE_ENV production
ENV NODE_ENV PLANNING_TOOL_DB=$PLANNING_TOOL_DB
COPY . .
RUN yarn build

FROM base AS release
WORKDIR /
ENV NODE_ENV production
ENV NODE_ENV PLANNING_TOOL_DB=$PLANNING_TOOL_DB
COPY --from=build . .
RUN apk update && apk add --no-cache supervisor
COPY supervisord.conf /etc/supervisor/conf.d/supervisord.conf
RUN mkdir -p /var/log/
EXPOSE 4000 4001

CMD ["/usr/bin/supervisord"]
