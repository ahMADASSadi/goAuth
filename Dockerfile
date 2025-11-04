# syntax=docker/dockerfile:1

# --- Builder Stage ---
FROM golang:1.24-alpine AS builder

WORKDIR /app

RUN apk add --no-cache build-base sqlite-dev git

COPY src/go.mod src/go.sum ./

RUN go mod download

COPY src/ .

ENV CGO_ENABLED=1

ENV GOOS=linux

RUN go build -ldflags '-extldflags "-static"' -tags osusergo,netgo -o /app/server ./cmd/main.go


# --- Final Image Stage ---
FROM alpine:latest

WORKDIR /app


RUN apk add --no-cache sqlite

COPY --from=builder /app/server .

COPY --from=builder /app/.env .

EXPOSE 8000

VOLUME ["/app/data"]

# Environment variables (can be overridden at runtime)
ENV DB_URL=/app/data/test.sqlite3
ENV PORT=8000
ENV HOST=0.0.0.0

# Create data directory for SQLite database
RUN mkdir -p /app/data

# Health check
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
    CMD wget --no-verbose --tries=1 --spider http://localhost:8000/health || exit 1

CMD ["./server"]

