FROM nginx:1.17.8-alpine
FROM node

RUN mkdir -p /usr/node/app
WORKDIR /usr/node/app
COPY package*.json ./
RUN yarn install

COPY . .
COPY .env ./build

COPY  build /usr/share/nginx/html
RUN rm /etc/nginx/conf.d/default.conf
COPY nginx/nginx.conf /etc/nginx/conf.d

EXPOSE 8081

CMD [ "nginx", "yarn", "start" ]