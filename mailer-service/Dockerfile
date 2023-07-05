FROM golang:1.20-alpine as builder

WORKDIR /app

COPY . .

RUN go build -o mailer cmd/api/main.go cmd/api/routes.go  cmd/api/handlers.go cmd/api/helpers.go cmd/api/mailer.go 

FROM alpine:latest 
LABEL org.opencontainers.image.source https://github.com/nu12/go-microservices

WORKDIR /app

COPY --from=builder /app/mailer /app/mailer
COPY --from=builder /app/templates/ /app/templates/

EXPOSE 8080

ENTRYPOINT ["./mailer"]