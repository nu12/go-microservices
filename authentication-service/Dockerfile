FROM golang:1.20-alpine as builder

WORKDIR /app

COPY . .

RUN go build -o authentication cmd/api/main.go cmd/api/routes.go cmd/api/handlers.go cmd/api/helpers.go

FROM alpine:latest

WORKDIR /app

COPY --from=builder /app/authentication /app/authentication

EXPOSE 8080

ENTRYPOINT ["./authentication"]