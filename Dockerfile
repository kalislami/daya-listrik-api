# Stage pertama: Build Go application
FROM golang:1.23.1-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod tidy

COPY . ./

RUN go build -o app cmd/server/main.go

# Stage kedua: Menjalankan aplikasi Go
FROM alpine:3.17

# Menyalin aplikasi Go yang telah dibangun dari stage builder
COPY --from=builder /app/app /app/
COPY --from=builder /app/migrations /app/migrations
COPY .env /app/

WORKDIR /app

# Menambahkan dependensi runtime
RUN apk --no-cache add ca-certificates

EXPOSE 8080

CMD ["./app"]
