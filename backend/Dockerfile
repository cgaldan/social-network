FROM golang:1.24-alpine AS builder

RUN apk add --no-cache build-base sqlite-dev

WORKDIR /src

COPY go.mod go.sum ./
RUN go mod download

COPY . .

ENV CGO_ENABLED=1 \
    GOOS=linux

RUN go build -ldflags="-s -w" -o /out/server ./cmd/server

FROM alpine:3.20 AS runtime

RUN apk add --no-cache sqlite-libs \
    && addgroup -S app && adduser -S -G app app

WORKDIR /app

COPY --from=builder /out/server /app/server

RUN touch /app/.env \
    && mkdir -p /app/data/database /app/data/uploads \
    && chown -R app:app /app

USER app

EXPOSE 8000

ENV SERVER_PORT=8000 \
    DATABASE_PATH=/app/data/database/social.db \
    UPLOAD_PATH=/app/data/uploads
    
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
  CMD wget --no-verbose --tries=1 --spider http://localhost:8000/health || exit 1

CMD ["/app/server"]
