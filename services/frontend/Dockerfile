# этап сборки (build stage)
FROM node:8.10.0-alpine as build-stage
WORKDIR /app
COPY package*.json ./
RUN npm install
RUN npm install @vue/cli -g
COPY . .
RUN npm run build

# этап production (production-stage)
FROM nginx:1.16.0-alpine as production-stage
COPY --from=build-stage /app/dist /usr/share/nginx/html
RUN rm /etc/nginx/conf.d/default.conf
COPY nginx.conf /etc/nginx/conf.d