# --- 阶段 1: 编译阶段 ---
# 使用 --platform=$BUILDPLATFORM 强制在 GitHub 原生 x86 环境下运行编译器，避免模拟运行导致的极慢速度
FROM --platform=$BUILDPLATFORM golang:1.21-alpine AS builder

# 声明 Buildx 自动注入的变量
ARG TARGETOS
ARG TARGETARCH
ARG TARGETVARIANT

WORKDIR /app
COPY . .

# 关键改动：
# 1. CGO_ENABLED=0 确保静态链接（跨平台运行不报错）
# 2. 动态获取 GOARM（针对 arm/v7）
RUN if [ "$TARGETARCH" = "arm" ]; then \
      export GOARM=$(echo $TARGETVARIANT | cut -c 2); \
    fi; \
    CGO_ENABLED=0 GOOS=$TARGETOS GOARCH=$TARGETARCH \
    go build -ldflags="-s -w" -o server .

# --- 阶段 2: 运行阶段 ---
FROM alpine:latest
# 安装基础证书和时区（对于微信推送等 HTTPS 请求是必须的）
RUN apk --no-cache add ca-certificates tzdata
WORKDIR /root/
# 仅从编译阶段拷贝最终产物
COPY --from=builder /app/server .
EXPOSE 10001
CMD ["./server"]
