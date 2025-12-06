# 阶段1: 编译
FROM golang:1.21-alpine AS builder
WORKDIR /app
COPY main.go .
# 静态编译，去除调试符号，减小体积
RUN go build -ldflags="-s -w" -o server main.go

# 阶段2: 运行 (极简镜像)
FROM alpine:latest
# 安装证书(必须)和时区
RUN apk --no-cache add ca-certificates tzdata
WORKDIR /root/
COPY --from=builder /app/server .
EXPOSE 10001
CMD ["./server"]
