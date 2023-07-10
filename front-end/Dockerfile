FROM golang:1.20-alpine as builder

WORKDIR /app

COPY . .

RUN go build -o front cmd/web/main.go

FROM alpine:latest 
LABEL org.opencontainers.image.source https://github.com/nu12/go-microservices

ARG commit
ENV COMMIT=${commit}

WORKDIR /app

COPY --from=builder /app/front /app/front
COPY --from=builder /app/cmd/web/templates/ /app/cmd/web/templates/

EXPOSE 8080

ENTRYPOINT ["./front"]