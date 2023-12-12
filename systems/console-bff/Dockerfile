FROM node:18-alpine

# Create app directory
WORKDIR /usr/src/app

# Install app dependencies
COPY package.json yarn.lock ./
RUN yarn install --frozen-lockfile

# Bundle app source
COPY . .

# Create a new user 'nodeuser' and switch to it
RUN adduser -D nodeuser
USER nodeuser

# Install supervisor
USER root
RUN apk update && apk add supervisor
COPY supervisord.conf /etc/supervisor/conf.d/supervisord.conf
RUN mkdir -p /var/log/
USER nodeuser

# Expose ports
EXPOSE 4000 4001

# Start the app
CMD ["/usr/bin/supervisord"]