FROM golang:1.20-alpine as builder

WORKDIR /app

COPY . .

RUN go build -o listener main.go

FROM alpine:latest 
LABEL org.opencontainers.image.source https://github.com/nu12/go-microservices

WORKDIR /app

COPY --from=builder /app/listener /app/listener

EXPOSE 8080

ENTRYPOINT ["./listener"]