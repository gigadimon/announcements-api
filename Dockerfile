FROM golang:alpine as builder

WORKDIR /app

COPY . .

RUN go mod download

RUN go build -o /app/main main.go

FROM alpine:3

COPY --from=builder /app/main main