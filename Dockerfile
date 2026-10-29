# Single image: builds the React UI, builds the Go server, and serves both on $PORT.
FROM node:22-alpine AS web
WORKDIR /web
COPY apps/web/package*.json ./
RUN npm install --no-audit --no-fund
COPY apps/web/ ./
RUN npm run build

FROM golang:1.23-alpine AS server
WORKDIR /src
COPY services/game/ ./
RUN CGO_ENABLED=0 go build -o /belot-server ./cmd/server

FROM alpine:3.20
RUN adduser -D app
COPY --from=server /belot-server /belot-server
COPY --from=web /web/dist /static
ENV STATIC_DIR=/static PORT=3000
USER app
EXPOSE 3000
HEALTHCHECK CMD wget -qO- http://localhost:3000/healthz || exit 1
CMD ["/belot-server"]
