FROM golang:1.23.5-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
WORKDIR /app/cmd
RUN go build -o go-api .

FROM alpine:latest

RUN apk --no-cache add ca-certificates

WORKDIR /app/cmd

COPY --from=builder app/cmd/go-api .
COPY [/certs/ca.pem] /certs/ca.pem

EXPOSE 8080

CMD ["./go-api"]