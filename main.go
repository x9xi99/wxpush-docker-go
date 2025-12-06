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

var (
	CorpID, AgentID, Secret, PushKey, Port string
	AccessToken                            string
)

func init() {
	CorpID = os.Getenv("CORP_ID")
	AgentID = os.Getenv("AGENT_ID")
	Secret = os.Getenv("SECRET")
	PushKey = os.Getenv("PUSH_KEY")
	Port = os.Getenv("PORT")
	if Port == "" { Port = "10001" }

	if CorpID == "" || AgentID == "" || Secret == "" || PushKey == "" {
		log.Fatal("❌ 配置缺失，请检查环境变量！")
	}
}

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

	body := make(map[string]interface{})

	// 2. 智能解析：如果是 GET 请求，从 URL 获取参数；如果是 POST，解析 JSON
	if r.Method == http.MethodGet {
		query := r.URL.Query()
		for k, v := range query {
			body[k] = v[0]
		}
		// 🌟 增强功能：如果 URL 里传了 card 相关的字段，自动组装成 textcard 结构
		if body["msgtype"] == "textcard" {
			body["textcard"] = map[string]string{
				"title":       getString(body, "title", "新通知"),
				"description": getString(body, "description", "请查看详情"),
				"url":         getString(body, "url", "https://work.weixin.qq.com"),
				"btntxt":      getString(body, "btntxt", "查看详情"),
			}
		}
	} else {
		if json.NewDecoder(r.Body).Decode(&body) != nil {
			http.Error(w, `{"err":"JSON格式错误"}`, 400)
			return
		}
	}

	// 3. 构造 Payload
	msgType := getString(body, "msgtype", "text")
	payload := map[string]interface{}{
		"touser":  getString(body, "touser", "@all"),
		"agentid": AgentID,
		"msgtype": msgType,
	}

	// 兼容处理
	if msgType == "text" {
		if _, ok := body["text"]; !ok {
			payload["text"] = map[string]string{"content": getString(body, "content", "")}
		} else {
			payload["text"] = body["text"]
		}
	} else if val, ok := body[msgType]; ok {
		payload[msgType] = val
	}

	// 4. 发送
	jsonBytes, _ := json.Marshal(payload)
	wxUrl := fmt.Sprintf("https://qyapi.weixin.qq.com/cgi-bin/message/send?access_token=%s", AccessToken)
	resp, err := http.Post(wxUrl, "application/json", bytes.NewBuffer(jsonBytes))
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	defer resp.Body.Close()
	w.Header().Set("Content-Type", "application/json")
	io.Copy(w, resp.Body)
}

func getString(m map[string]interface{}, key, defaultVal string) string {
	if v, ok := m[key]; ok {
		if s, ok := v.(string); ok {
			return s
		}
	}
	return defaultVal
}

func main() {
	go refreshToken()
	http.HandleFunc("/send", handleSend)
	log.Printf("🚀 WxPush 服务启动，端口: %s", Port)
	http.ListenAndServe(":"+Port, nil)
}
