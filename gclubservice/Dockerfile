FROM node:14 as base

WORKDIR /home/node/app

ENV TZ="Asia/Bangkok"

COPY package*.json ./

RUN npm i

COPY . .

FROM base as production

ENV NODE_PATH=./build

RUN npm run build