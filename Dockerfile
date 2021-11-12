FROM node

RUN mkdir -p /usr/node/app
WORKDIR /usr/node/app
COPY package*.json ./
RUN yarn install

COPY . .

EXPOSE 8081

CMD [ "yarn", "start" ]