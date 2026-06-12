FROM golang:alpine AS builder

WORKDIR /app
COPY . .
RUN CGO_ENABLED=0 go build -mod=vendor -o server ./cmd/server/

FROM alpine:3.21

RUN apk add --no-cache tzdata ca-certificates
WORKDIR /app
COPY --from=builder /app/server .
EXPOSE 8080
CMD ["./server"]
