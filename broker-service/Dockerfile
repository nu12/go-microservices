FROM golang:1.20-alpine as builder

WORKDIR /app

COPY . .

RUN go build -o broker cmd/api/main.go cmd/api/routes.go  cmd/api/handlers.go cmd/api/helpers.go cmd/api/amqp.go

FROM alpine:latest 
LABEL org.opencontainers.image.source https://github.com/nu12/go-microservices

WORKDIR /app

COPY --from=builder /app/broker /app/broker

EXPOSE 8080

ENTRYPOINT ["./broker"]