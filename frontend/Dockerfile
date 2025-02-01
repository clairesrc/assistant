FROM --platform=amd64 node:latest as builder

WORKDIR /app

COPY package.json yarn.lock ./

RUN yarn install --check-files --non-interactive --ignore-optional --frozen-lockfile

COPY . .

RUN yarn build

FROM --platform=amd64 nginx:alpine

COPY --from=builder /app/dist /usr/share/nginx/html

EXPOSE 80

CMD ["nginx", "-g", "daemon off;"]
