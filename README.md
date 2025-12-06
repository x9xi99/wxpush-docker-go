# 🚀 WxPush Docker Go (企业微信消息推送服务)

![Go](https://img.shields.io/badge/Go-1.21-blue.svg)
![Docker](https://img.shields.io/badge/Docker-Lightweight-blue)
![License](https://img.shields.io/badge/License-MIT-green)

这是一个基于 **Go 语言** 重写的高性能企业微信（WeCom）消息推送服务。

相比于传统的 Python/Java 实现，本项目主打 **极致轻量**、**安全** 和 **零依赖**。基于 Alpine 镜像静态编译，Docker 镜像体积仅约 **15MB**，运行时内存占用仅 **5~10MB**。

---

## ✨ 核心特性

- **🚀 极致轻量**：告别 Python/Node.js 的庞大体积，Go 语言原生编译，瞬间启动，几乎不占资源。
- **📩 全类型支持**：采用“透传模式”，完美支持企业微信所有消息类型（文本、Markdown、图文、卡片等）。
- **🛡️ 安全设计**：敏感配置（Secret/Key）通过环境变量注入，不在代码中留痕；支持自定义鉴权密钥防止恶意调用。
- **🔄 自动维护**：内置 AccessToken 自动刷新与过期重试机制，7x24 小时稳定运行。
- **🐳 多样化部署**：支持 Docker、Docker Compose、Koyeb/Render 云平台以及二进制直接运行。

---

## 🛠️ 部署前准备

在部署之前，请确保你已获取企业微信的应用配置：
1. **企业ID (CORP_ID)**: [企业微信后台](https://work.weixin.qq.com/) -> 我的企业 -> 企业信息 -> 最下方。
2. **应用ID (AGENT_ID)**: 应用管理 -> 创建/选择应用 -> AgentId。
3. **应用密钥 (SECRET)**: 应用管理 -> 点击应用 -> Secret (需在手机端查看)。

---

## 📂 推荐部署流程 (Docker CLI)

为了确保 Docker 能正确读取配置文件，建议按照以下标准步骤操作：

### 1. 创建项目目录
在服务器上创建一个文件夹，用来存放配置文件，防止文件散乱。

```bash
# 新建文件夹
mkdir wxpush

# 进入文件夹 (⚠️ 重要：后续操作都在这个目录下进行)
cd wxpush
```

### 2. 创建配置文件 (.env)
在 `wxpush` 目录下创建一个名为 `.env` 的文件。

```bash
nano .env
```

**粘贴以下内容（请务必修改为你的真实信息）：**

```ini
# --- 企业微信配置 (必填) ---
CORP_ID=ww49d7776235xxxxxx
AGENT_ID=1000002
SECRET=e5Nly7h_C_yl8lzw-poNTYNVpSx3f98b7OJztxxxxxx

# --- 安全配置 (必填) ---
# 自定义鉴权密钥，调用接口时 URL 必须带上 ?key=这个值
PUSH_KEY=您的密钥

# --- 服务配置 ---
# 容器内部监听端口 (建议保持与外部映射一致)
# 默认PORT=10001，如不需要自定义可以注释掉此项
#以下为自定义端口
PORT=10008
# 时区设置 (保证日志时间正确)
TZ=Asia/Shanghai
```
*(保存退出：按 Ctrl+O -> 回车 -> Ctrl+X)*

### 3. 一键启动
确保你仍然在 `wxpush` 目录下，运行以下命令：
这里自定义端口-p参数10008:10008，默认就应该为10001:10001

```bash
docker run -d \
  --name wxpush \
  --restart always \
  -p 10008:10008 \
  --env-file .env \
  ghcr.io/x9xi99/wxpush-docker-go:latest
```

> **命令解析**：
> - `-p 10008:10008`: 将服务器的 10008 端口映射到容器的 10008 端口。
> - `--env-file .env`: 读取当前目录下的 `.env` 文件作为环境变量。

### 4. 验证部署
查看容器日志，确认 Token 获取成功：
```bash
docker logs -f wxpush
```
*如果看到 `✅ Token 刷新成功` 字样，说明服务已正常运行。*

---

## 🐳 备选部署：Docker Compose

如果你更习惯使用 Compose 管理，请在 `wxpush` 目录下创建 `docker-compose.yml`：

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

启动命令：
```bash
docker-compose up -d
```

---

## 🔌 API 使用说明

- **接口地址**: `http://你的服务器IP:10008/send?key=你的PUSH_KEY`
- **请求方式**: `POST`
- **Content-Type**: `application/json`

本服务支持企业微信官方文档中的**所有**消息格式。你只需要构建对应的 JSON 即可。

### 📨 示例 1：发送纯文本 (Text)
支持简写模式，直接传 `content` 即可。

```bash
curl -X POST "[http://127.0.0.1:10008/send?key=你的PUSH_KEY](http://127.0.0.1:10008/send?key=你的PUSH_KEY)" \
     -H "Content-Type: application/json" \
     -d '{
       "msgtype": "text",
       "content": "🔔 服务器 SSH 登录警告！\n来源IP: 192.168.1.100"
     }'
```

### 📝 示例 2：发送 Markdown
Markdown 支持颜色高亮，适合做监控告警。

```bash
curl -X POST "[http://127.0.0.1:10008/send?key=你的PUSH_KEY](http://127.0.0.1:10008/send?key=你的PUSH_KEY)" \
     -H "Content-Type: application/json" \
     -d '{
       "msgtype": "markdown",
       "markdown": {
         "content": "# 每日监控日报\n> CPU: <font color=\"warning\">80%</font>\n> 内存: <font color=\"info\">正常</font>\n> [查看详情](http://example.com)"
       }
     }'
```

### 🃏 示例 3：发送文本卡片 (TextCard)
适合做漂亮的通知跳转。

```bash
curl -X POST "[http://127.0.0.1:10008/send?key=你的PUSH_KEY](http://127.0.0.1:10008/send?key=你的PUSH_KEY)" \
     -H "Content-Type: application/json" \
     -d '{
       "msgtype": "textcard",
       "textcard": {
         "title": "GitHub Actions 构建成功",
         "description": "<div class=\"gray\">2025-12-06</div> <br>项目构建完成，耗时 35s。",
         "url": "[https://github.com/x9xi99](https://github.com/x9xi99)",
         "btntxt": "立即查看"
       }
     }'
```

### 🖼️ 示例 4：其他类型 (图片/文件等)
完全遵循 [企业微信官方 API 文档](https://developer.work.weixin.qq.com/document/path/90236) 的 JSON 结构。例如发送图片：

```json
{
  "msgtype": "image",
  "image": {
    "media_id": "MEDIA_ID_xxxxxx"
  },
  "touser": "@all"
}
```

---

## ⚙️ 环境变量详解

| 变量名 | 必填 | 示例 | 说明 |
| :--- | :--- | :--- | :--- |
| `CORP_ID` | ✅ | `ww49...` | 企业微信 CorpID |
| `AGENT_ID` | ✅ | `1000002` | 企业微信应用 AgentID |
| `SECRET` | ✅ | `e5Nl...` | 企业微信应用 Secret |
| `PUSH_KEY` | ✅ | `你的PUSH_KEY` | **自定义鉴权密钥**，防止接口被他人滥用 |
| `PORT` | ❌ | `10008` | 容器内部监听端口，默认为 10001 |
| `TZ` | ❌ | `Asia/Shanghai` | 容器时区，建议设置以保证日志时间正确 |

---

## ☁️ 云平台部署 (Koyeb/Render)

如果不想使用自己的服务器，可以使用免费的 PaaS 平台。

1. **镜像**: `ghcr.io/x9xi99/wxpush-docker-go:latest`
2. **端口**: 在平台设置中暴露端口 `10008`。
3. **环境变量**: 直接在平台的控制台 (Environment Variables) 中填入 `CORP_ID`, `SECRET` 等键值对（无需 .env 文件）。
4. **健康检查**: 如果平台需要，HTTP 路径为 `/` (需鉴权) 或直接调用接口测试。

---

## 🤝 贡献与构建

如果你想自己修改代码并构建镜像：

```bash
# 1. 克隆代码
git clone [https://github.com/x9xi99/wxpush-docker-go.git](https://github.com/x9xi99/wxpush-docker-go.git)

# 2. 构建镜像
docker build -t my-wxpush .

# 3. 运行
docker run -d -p 10008:10008 --env-file .env my-wxpush
```

## 📝 License

MIT License.
