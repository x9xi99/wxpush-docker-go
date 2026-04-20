# --- 阶段 1: 编译阶段 ---
# 使用 golang:1-alpine 保持最新稳定版，使用 BUILDPLATFORM 加速
FROM --platform=$BUILDPLATFORM golang:1-alpine AS builder

# 声明 Buildx 自动注入的变量
ARG TARGETOS
ARG TARGETARCH
ARG TARGETVARIANT

WORKDIR /app

# 拷贝源代码
COPY main.go .

# 解决 "cannot find main module" 并处理跨平台编译
RUN if [ "$TARGETARCH" = "arm" ]; then \
      export GOARM=$(echo $TARGETVARIANT | cut -c 2); \
    fi; \
    CGO_ENABLED=0 GOOS=$TARGETOS GOARCH=$TARGETARCH \
    go build -ldflags="-s -w" -o server main.go

# --- 阶段 2: 运行阶段 ---
FROM alpine:latest
# 安装 HTTPS 请求所需的证书和时区数据（生产环境必备）
RUN apk --no-cache add ca-certificates tzdata
WORKDIR /root/
# 仅从编译阶段拷贝最终产物
COPY --from=builder /app/server .

# 暴露程序端口
EXPOSE 10001

# 启动程序
CMD ["./server"]
