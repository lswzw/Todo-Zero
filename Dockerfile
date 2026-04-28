# ===== Stage 1: 构建前端 =====
FROM node:20-alpine AS frontend-builder

WORKDIR /app
COPY web/package.json web/package-lock.json ./web/
RUN cd web && npm ci
COPY web/ ./web/
RUN cd web && npm run build

# ===== Stage 2: 编译后端 =====
FROM golang:1.25-alpine AS backend-builder

WORKDIR /app/server
COPY server/go.mod server/go.sum ./
RUN go mod download
COPY server/ ./
COPY --from=frontend-builder /app/server/dist ./dist/
RUN CGO_ENABLED=0 go build -ldflags="-s -w" -o todo-zero .

# ===== Stage 3: 运行 =====
FROM alpine:3.21

RUN apk add --no-cache tzdata ca-certificates \
    && cp /usr/share/zoneinfo/Asia/Shanghai /etc/localtime \
    && echo "Asia/Shanghai" > /etc/timezone \
    && apk del tzdata

WORKDIR /app
COPY --from=backend-builder /app/server/todo-zero .

ENV TZ=Asia/Shanghai
EXPOSE 8888
VOLUME ["/app/data"]

ENTRYPOINT ["./todo-zero"]
CMD ["-port", "8888"]
