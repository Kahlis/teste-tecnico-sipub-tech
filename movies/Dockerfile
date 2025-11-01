FROM golang:1.25.3-alpine AS builder

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags='-w -s' -a -o movies ./cmd/server

# -----------------------
# Testing package
# -----------------------
FROM golang:1.25.3-alpine AS tester

WORKDIR /app

COPY --from=builder /app/go.mod /app/go.sum ./
COPY --from=builder /app ./

RUN go install gotest.tools/gotestsum@latest
CMD ["go", "test", "-v", "./core/integration/test/..."]

# -----------------------
# Runtime minimal package
# -----------------------
FROM alpine:3.20
RUN apk add --no-cache curl ca-certificates

COPY --from=builder /app/movies /app/movies
COPY --from=builder /app/seed /app/seed
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

EXPOSE 9090

CMD ["/app/movies"]
