# Build stage
FROM golang:1.26-alpine AS builder

WORKDIR /app

RUN apk add --no-cache git

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s" -o /app/recodations-api ./cmd/api

# Production stage
FROM alpine:3.20

WORKDIR /app

RUN apk --no-cache add ca-certificates tzdata

COPY --from=builder /app/recodations-api /app/recodations-api
COPY --from=builder /app/docs /app/docs

EXPOSE 3000

ENTRYPOINT ["/app/recodations-api"]
