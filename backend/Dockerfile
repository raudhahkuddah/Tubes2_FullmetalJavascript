FROM golang:1.24-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o api-server .

FROM alpine:latest

# Add ca certificates for HTTPS requests
RUN apk --no-cache add ca-certificates

WORKDIR /root/

COPY --from=builder /app/api-server .

RUN mkdir -p /root/data

EXPOSE 8080

CMD ["./api-server"]