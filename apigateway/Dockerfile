FROM golang:1.25.3-alpine AS builder

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags='-w -s' -a -o apigateway ./cmd/api

FROM alpine:3.20
RUN apk add --no-cache curl ca-certificates

COPY --from=builder /app/apigateway /app/apigateway
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

EXPOSE 8080

CMD ["/app/apigateway"]
