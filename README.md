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
# .env 文件内容
CORP_ID=ww49d7776235xxxxxx
AGENT_ID=1000002
SECRET=e5Nly7h_C_yl8lzw-poNTYNVpSx3f98b7OJztxxxxxx

# 自定义鉴权密钥 (调用接口时需带上 ?key=hexin123)
PUSH_KEY=hexin123

# 服务监听端口 (容器内部)
PORT=10008

# 时区
TZ=Asia/Shanghai
