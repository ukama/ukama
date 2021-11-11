FROM nginx:1.17.8-alpine
FROM node

RUN mkdir -p /usr/node/app
WORKDIR /usr/node/app
COPY package*.json ./
RUN yarn install

COPY  build /usr/share/nginx/html
RUN rm /etc/nginx/conf.d/default.conf
COPY nginx/nginx.conf /etc/nginx/conf.d

COPY . .

# RUN Yarn build

COPY .env ./build

# WORKDIR ./build

EXPOSE 3000

CMD [ "nginx", "yarn", "start" ]