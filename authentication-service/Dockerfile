FROM golang:1.20-alpine as builder

WORKDIR /app

COPY . .

RUN go build -o authentication cmd/api/main.go cmd/api/grpc.go 

FROM alpine:latest
LABEL org.opencontainers.image.source https://github.com/nu12/go-microservices

WORKDIR /app

COPY --from=builder /app/authentication /app/authentication

EXPOSE 50001

ENTRYPOINT ["./authentication"]