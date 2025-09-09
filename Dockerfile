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

VOLUME ["/app/"]

ENV DB_URL=/app/test.sqlite3
ENV PORT = 8000
ENV HOST = 0.0.0.0


CMD ["./server"]

