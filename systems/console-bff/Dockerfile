FROM node

RUN mkdir -p /console/bff
WORKDIR /console/bff
ENV NODE_ENV production
ENV NODE_ENV DATABASE_URL="postgresql://postgres:admin@127.0.0.1/planning"
COPY package*.json ./
RUN yarn install

COPY . .

EXPOSE 8080

CMD [ "yarn", "start", "metrics-start" ]