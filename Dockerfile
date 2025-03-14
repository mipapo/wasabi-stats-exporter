FROM golang:1.22-alpine AS builder

WORKDIR /app

COPY main.go go.mod go.sum .

RUN go build -o /app/wasabi_exporter main.go

FROM alpine:latest

WORKDIR /root/

COPY --from=builder /app/wasabi_exporter .

EXPOSE 8080

CMD ["./wasabi_exporter"]

