FROM golang:1.21-alpine AS builder

WORKDIR /app

COPY . .

RUN go mod tidy && \
    go build -o wasabi_exporter

FROM alpine:latest

WORKDIR /root/

COPY --from=builder /app/wasabi_exporter .

EXPOSE 8080

CMD ["./wasabi_exporter"]

