# ===== Stage 1: Build Vue 前端 =====
FROM node:20-alpine AS frontend
WORKDIR /app/frontend
COPY frontend/package*.json ./
RUN npm ci
COPY frontend/ .
RUN npm run build

# ===== Stage 2: Build Go 后端 =====
FROM golang:alpine AS backend
WORKDIR /app
COPY go.mod go.sum ./
COPY vendor/ ./vendor/
COPY . .
RUN CGO_ENABLED=0 go build -mod=vendor -o server ./cmd/server/

# ===== Stage 3: 最终轻量镜像 =====
FROM alpine:3.21
RUN apk add --no-cache tzdata ca-certificates
WORKDIR /app
COPY --from=backend /app/server .
COPY --from=frontend /app/frontend/dist ./frontend/dist
EXPOSE 8080
CMD ["./server"]
