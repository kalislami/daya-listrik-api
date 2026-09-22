FROM golang:1.24-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 go build -o /app/server ./cmd/server

FROM alpine:3.21
RUN apk add --no-cache ca-certificates
WORKDIR /app
COPY --from=builder /app/server ./server
COPY --from=builder /app/docs/openapi.yaml ./docs/openapi.yaml
EXPOSE 8080
CMD ["./server"]
