# Build stage
FROM golang:1.19-alpine3.16 AS builder

RUN apk add --no-cache gcc musl-dev

WORKDIR /app

COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY . ./

RUN CGO_ENABLED=1 GOOS=linux go build -ldflags="-w -s" -o /zerotier-webhook ./cmd/server

# Run stage
FROM alpine:3.16

RUN apk add --no-cache ca-certificates

WORKDIR /app

COPY --from=builder /zerotier-webhook /app/zerotier-webhook

EXPOSE 8080

CMD ["/app/zerotier-webhook"]
