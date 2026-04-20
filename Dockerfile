# 阶段1: 编译
# --platform=$BUILDPLATFORM 确保编译器在 GitHub Actions 的原生 x86 环境下高速运行
FROM --platform=$BUILDPLATFORM golang:1.21-alpine AS builder

# 这三个变量由 Docker Buildx 自动注入，不需要你手动赋值
ARG TARGETOS
ARG TARGETARCH
ARG TARGETVARIANT

WORKDIR /app
COPY main.go .

# 关键改动：使用环境变量控制交叉编译
# 1. CGO_ENABLED=0 确保静态编译，不依赖宿主机 C 库
# 2. GOOS 和 GOARCH 使用自动注入的参数
# 3. 如果是 arm/v7，TARGETVARIANT 会是 v7，我们需要处理一下 GOARM
RUN export GOARM=$(echo $TARGETVARIANT | cut -c 2); \
    CGO_ENABLED=0 GOOS=$TARGETOS GOARCH=$TARGETARCH GOARM=$GOARM \
    go build -ldflags="-s -w" -o server main.go

# 阶段2: 运行 (极简镜像)
FROM alpine:latest
RUN apk --no-cache add ca-certificates tzdata
WORKDIR /root/
# 从 builder 阶段拷贝编译好的二进制文件
COPY --from=builder /app/server .
EXPOSE 10001
CMD ["./server"]
