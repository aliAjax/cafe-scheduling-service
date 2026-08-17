FROM golang:1.22-alpine AS builder

WORKDIR /src
COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -trimpath -o /out/cafe-scheduling-api ./cmd/api

FROM alpine:3.20

RUN apk add --no-cache ca-certificates tzdata
WORKDIR /app
COPY --from=builder /out/cafe-scheduling-api /app/cafe-scheduling-api
COPY migrations /app/migrations

EXPOSE 8080
ENTRYPOINT ["/app/cafe-scheduling-api"]
