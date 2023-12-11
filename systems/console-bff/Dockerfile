FROM node:18-alpine

WORKDIR /

COPY package*.json ./
ENV NODE_ENV production
ENV NODE_ENV PLANNING_TOOL_DB=$PLANNING_TOOL_DB
RUN yarn install
COPY . .
RUN apt-get update && apt-get install -y supervisor && apt-get clean
COPY supervisord.conf /etc/supervisor/conf.d/supervisord.conf
RUN mkdir -p /var/log/
EXPOSE 4000 4001

CMD ["/usr/bin/supervisord"]