FROM node:21-alpine AS builder
WORKDIR /app

ENV GENERATE_SOURCEMAP=false

COPY package.json package.json
COPY yarn.lock yarn.lock
COPY tsconfig.json tsconfig.json
COPY public/ public/
COPY src/ src/


RUN yarn install --network-timeout 500000
RUN yarn build --production

FROM nginx:1.26.0-alpine
COPY --from=builder /app/build /usr/share/nginx/html
COPY ./nginx.conf /etc/nginx/conf.d/default.conf
EXPOSE 80
