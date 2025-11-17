FROM golang:1.23-alpine AS builder

WORKDIR /build

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /build/svc ./cmd/server/main.go

FROM alpine:3.20

RUN apk add --no-cache ca-certificates

WORKDIR /app

COPY --from=builder /build/svc /app/svc

ENV HTTP_HOST=0.0.0.0 HTTP_PORT=8080

EXPOSE 8080

CMD ["./svc"]
