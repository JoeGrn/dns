FROM golang:1.23.5-alpine AS builder

WORKDIR /app

COPY go.mod ./

RUN go mod download

COPY . .

RUN go build -o /app/dns cmd/dns/main.go

FROM alpine:latest

RUN apk --no-cache add ca-certificates

WORKDIR /root/

COPY --from=builder /app/dns .

EXPOSE 2053/udp

CMD ["./dns"]
