FROM node:alpine as build
WORKDIR /app

COPY package*.json ./
RUN npm ci
COPY . .
RUN npm run build --production

FROM docker.io/nginx:alpine AS prod

RUN apk update --no-cache -U \
    && rm -rf /usr/share/nginx/html

COPY --from=build /app/dist/frontend/browser /usr/share/nginx/html
