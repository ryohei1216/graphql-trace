# 1st stage
FROM golang:1.21 AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 go build -o server server.go

# 2nd stage
FROM alpine:latest

COPY --from=builder /app/server /app/server

CMD ["/app/server"]