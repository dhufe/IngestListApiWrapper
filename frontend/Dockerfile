# Stage 1: Build the Angular app
FROM node:alpine as build
WORKDIR /app

# Copy package.json and package-lock.json to install dependencies
COPY package*.json ./
RUN npm ci

# Copy the entire project and build it for production
COPY . .
RUN npm run build --prod

# Stage 2: Serve the app with Node.js and Express

FROM docker.io/nginx:alpine AS prod

RUN apk update --no-cache -U
COPY --from=build /app/dist/frontend /usr/share/nginx/html
