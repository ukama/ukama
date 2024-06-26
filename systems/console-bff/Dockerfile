FROM node:lts-alpine

WORKDIR /

COPY package*.json yarn.lock ./

ENV NODE_ENV production
ENV NODE_ENV PLANNING_TOOL_DB=$PLANNING_TOOL_DB
RUN yarn install --frozen-lockfile --ignore-scripts

COPY . .

RUN apk update && apk add --no-cache supervisor
COPY supervisord.conf /etc/supervisor/conf.d/supervisord.conf
COPY supervisord.conf /etc/supervisor/supervisord.conf
RUN mkdir -p /var/log/
CMD ["/usr/bin/supervisord", "-c", "/etc/supervisor/supervisord.conf"]
