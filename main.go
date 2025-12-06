package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"
)

// 全局变量
var (
	CorpID, AgentID, Secret, PushKey, Port string
	AccessToken                            string
)

func init() {
	// 读取环境变量
	CorpID = os.Getenv("CORP_ID")
	AgentID = os.Getenv("AGENT_ID")
	Secret = os.Getenv("SECRET")
	PushKey = os.Getenv("PUSH_KEY")
	Port = os.Getenv("PORT")
	if Port == "" { Port = "10001" }

	if CorpID == "" || AgentID == "" || Secret == "" || PushKey == "" {
		log.Fatal("❌ 配置缺失，请检查 .env 文件或环境变量！")
	}
}

// 自动刷新 Token
func refreshToken() {
	for {
		url := fmt.Sprintf("https://qyapi.weixin.qq.com/cgi-bin/gettoken?corpid=%s&corpsecret=%s", CorpID, Secret)
		if resp, err := http.Get(url); err == nil {
			body, _ := io.ReadAll(resp.Body)
			resp.Body.Close()
			var res map[string]interface{}
			json.Unmarshal(body, &res)
			if token, ok := res["access_token"].(string); ok {
				AccessToken = token
				log.Println("✅ Token 刷新成功")
			} else {
				log.Printf("❌ Token 获取失败: %s", body)
			}
		}
		time.Sleep(7000 * time.Second)
	}
}

// 处理推送请求
func handleSend(w http.ResponseWriter, r *http.Request) {
	// 1. 验证 Key
	if r.URL.Query().Get("key") != PushKey {
		http.Error(w, `{"err":"无效密钥"}`, 401)
		return
	}
	if AccessToken == "" {
		http.Error(w, `{"err":"服务正在初始化，请稍后"}`, 503)
		return
	}

	// 2. 解析 JSON
	var body map[string]interface{}
	if json.NewDecoder(r.Body).Decode(&body) != nil {
		http.Error(w, `{"err":"JSON格式错误"}`, 400)
		return
	}

	// 3. 构造微信消息体
	msgType := "text"
	if t, ok := body["msgtype"].(string); ok { msgType = t }

	payload := map[string]interface{}{
		"touser":  body["touser"],
		"toparty": body["toparty"],
		"totag":   body["totag"],
		"msgtype": msgType,
		"agentid": AgentID,
	}
	// 默认 text 处理
	if msgType == "text" && body["text"] == nil {
		content := ""
		if c, ok := body["content"].(string); ok { content = c }
		payload["text"] = map[string]string{"content": content}
	} else {
		// 其他类型直接透传
		if val, ok := body[msgType]; ok { payload[msgType] = val }
	}
	// 如果 body 里没有 touser，默认为 @all
	if payload["touser"] == nil { payload["touser"] = "@all" }

	// 4. 发送给微信
	jsonBytes, _ := json.Marshal(payload)
	wxUrl := fmt.Sprintf("https://qyapi.weixin.qq.com/cgi-bin/message/send?access_token=%s", AccessToken)
	resp, err := http.Post(wxUrl, "application/json", bytes.NewBuffer(jsonBytes))
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	defer resp.Body.Close()
	io.Copy(w, resp.Body)
}

func main() {
	go refreshToken()
	http.HandleFunc("/send", handleSend)
	log.Printf("🚀 服务启动，监听端口: %s", Port)
	http.ListenAndServe(":"+Port, nil)
}
