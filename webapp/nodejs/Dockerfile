FROM node:latest

WORKDIR /home/isucon/webapp/nodejs
COPY package.json package.json
COPY package-lock.json package-lock.json

RUN npm ci
