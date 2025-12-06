# 🚀 WxPush Docker Go (企业微信消息推送服务)

![Go](https://img.shields.io/badge/Go-1.21-blue.svg)
![Docker](https://img.shields.io/badge/Docker-Lightweight-blue)
![License](https://img.shields.io/badge/License-MIT-green)

这是一个基于 **Go 语言** 重写的企业微信（WeCom）消息推送服务。

相比于传统的 Python/Java 实现，本项目主打 **极致轻量**、**高性能** 和 **部署简单**。支持 Docker 一键部署，内存占用仅 **10MB** 左右。

## ✨ 特性

- **极低资源占用**：基于 Alpine 和 Go 静态编译，镜像体积仅 ~15MB，运行时内存占用 < 10MB。
- **全消息类型支持**：采用透传（Pass-through）机制，支持企业微信所有消息类型（文本、Markdown、卡片、图文等），无需更新代码即可支持新类型。
- **安全可靠**：配置通过环境变量注入，敏感信息不落地；支持自定义推送密钥（Key）校验。
- **自动维护 Token**：内置 AccessToken 自动刷新与过期重试机制，无需人工干预。
- **多平台支持**：支持 Docker、Docker Compose、二进制直接运行以及 Koyeb/Render 等云平台。

---

## 🛠️ 准备工作

在部署之前，你需要获取企业微信应用的配置信息：

1. 注册 [企业微信](https://work.weixin.qq.com/)。
2. 进入后台 -> **应用管理** -> **创建应用**。
3. 获取以下信息：
   - **企业ID (CORP_ID)**: 我的企业 -> 企业信息 -> 最下方。
   - **应用ID (AGENT_ID)**: 应用管理 -> 点击应用 -> AgentId。
   - **应用密钥 (SECRET)**: 应用管理 -> 点击应用 -> Secret (需在手机端查看)。

---

## 🐳 部署方式 1：Docker (推荐)

这是最简单的方式，直接使用 GitHub Packages 提供的镜像。

### 1. 创建配置文件
在服务器上创建一个 `.env` 文件，填入你的配置：

```bash
# .env 文件内容 (请修改为真实值)
CORP_ID=ww49d7776235xxxxxx
AGENT_ID=1000002
SECRET=e5Nly7h_C_yl8lzw-poNTYNVpSx3f98b7OJztxxxxxx

# 自定义鉴权密钥 (调用接口时需带上 ?key=hexin123)
PUSH_KEY=hexin123

# 服务监听端口 (容器内部)
PORT=10008

# 时区
TZ=Asia/Shanghai
```

### 2. 启动容器
运行以下命令即可（映射到服务器的 10008 端口）：

```bash
docker run -d \
  --name wxpush \
  --restart always \
  -p 10008:10008 \
  --env-file .env \
  ghcr.io/x9xi99/wxpush-docker-go:latest
```

---

## 📂 部署方式 2：Docker Compose

如果你喜欢用 Compose 管理，创建一个 `docker-compose.yml`：

```yaml
version: '3'
services:
  wxpush:
    image: ghcr.io/x9xi99/wxpush-docker-go:latest
    container_name: wxpush
    restart: always
    ports:
      - "10008:10008"
    env_file:
      - .env
```

然后运行：
```bash
docker-compose up -d
```

---

## ☁️ 部署方式 3：云平台 (Koyeb 等)

本项目完美支持无服务器（Serverless/PaaS）平台部署。

1. **镜像地址**: `ghcr.io/x9xi99/wxpush-docker-go:latest`
2. **端口设置**: 在平台设置中暴露端口 `10008`。
3. **环境变量**: 在平台的 Settings -> Environment Variables 中添加 `CORP_ID`, `SECRET` 等变量（无需创建 .env 文件）。

---

## 🔌 API 使用说明

- **接口地址**: `http://你的IP:10008/send?key=你的PUSH_KEY`
- **请求方式**: `POST`
- **Content-Type**: `application/json`

### 1. 发送纯文本 (Text)
支持简写模式，直接传 `content` 即可。

```bash
curl -X POST "http://127.0.0.1:10008/send?key=hexin123" \
     -H "Content-Type: application/json" \
     -d '{
       "msgtype": "text",
       "content": "服务器 SSH 登录警告！\n来源IP: 192.168.1.100"
     }'
```

### 2. 发送 Markdown
```bash
curl -X POST "http://127.0.0.1:10008/send?key=hexin123" \
     -H "Content-Type: application/json" \
     -d '{
       "msgtype": "markdown",
       "markdown": {
         "content": "# 每日监控\n> CPU: <font color=\"warning\">80%</font>\n> 内存: <font color=\"info\">Normal</font>"
       }
     }'
```

### 3. 发送文本卡片 (TextCard)
适合做漂亮的通知跳转。

```bash
curl -X POST "http://127.0.0.1:10008/send?key=hexin123" \
     -H "Content-Type: application/json" \
     -d '{
       "msgtype": "textcard",
       "textcard": {
         "title": "GitHub Actions 构建成功",
         "description": "<div class=\"gray\">2025-12-06</div> <br>项目 wxpush-docker-go 构建完成。",
         "url": "https://github.com/x9xi99",
         "btntxt": "查看详情"
       }
     }'
```

### 4. 高级用法 (透传模式)
本服务支持企业微信的所有消息格式。你只需要参考 [企业微信官方文档 - 发送消息](https://developer.work.weixin.qq.com/document/path/90236)，构建对应的 JSON 结构即可。

例如发送图片：
```json
{
  "msgtype": "image",
  "image": {
    "media_id": "MEDIA_ID_xxxx"
  },
  "touser": "@all"
}
```

---

## ⚙️ 环境变量详解

| 变量名 | 必填 | 默认值 | 说明 |
| :--- | :--- | :--- | :--- |
| `CORP_ID` | ✅ | - | 企业微信 CorpID |
| `AGENT_ID` | ✅ | - | 企业微信应用 AgentID |
| `SECRET` | ✅ | - | 企业微信应用 Secret |
| `PUSH_KEY` | ✅ | - | 自定义鉴权密钥，防止未授权调用 |
| `PORT` | ❌ | `10001` | 服务监听端口 (建议改为 10008) |
| `TZ` | ❌ | `Asia/Shanghai` | 容器时区设置 |

---

## 🛡️ 安全提示

1. **不要上传 .env 文件**：请确保 `.env` 文件被包含在 `.gitignore` 中，绝对不要提交到 GitHub 公开仓库。
2. **保护 PUSH_KEY**：`PUSH_KEY` 是你的 API 密码，请设置复杂一点，并不要泄露给他人。

## 📝 License

MIT License.
