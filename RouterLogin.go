package main

import (
	"crypto/rc4"
	"encoding/hex"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)

// --- 配置区域 ---
const (
	LoginURL = "http://192.168.3.1/ac_portal/login.php"
	CheckURL = "http://connect.rom.miui.com/generate_204" // 强连通性测试接口（防校园网网关劫持）
	
	UserName = "账号"
	// !!! 请在这里填入你平时登录用的真实明文密码 !!!
	Password = "密码" 
)

// CheckNetwork 检测是否真正连通外网
func CheckNetwork() bool {
	client := http.Client{Timeout: 5 * time.Second}
	resp, err := client.Get(CheckURL)
	if err != nil {
		return false
	}
	defer resp.Body.Close()
	// 被强制跳转登录页时通常返回200或302，只有彻底畅通才会返回 204 No Content
	return resp.StatusCode == 204
}

// encryptRC4 实现 ac_portal 的核心密码加密逻辑
func encryptRC4(password, key string) string {
	cipher, err := rc4.NewCipher([]byte(key))
	if err != nil {
		return ""
	}
	dst := make([]byte, len(password))
	cipher.XORKeyStream(dst, []byte(password))
	return hex.EncodeToString(dst) // 转换为你在抓包里看到的 12 位十六进制字符串
}

func DoLogin() {
	// 1. 生成当前毫秒级时间戳作为 auth_tag 和 密钥
	authTag := strconv.FormatInt(time.Now().UnixMilli(), 10)
	
	// 2. 动态生成加密密码 (这行代码彻底解决了你之前失效的问题)
	dynamicPwd := encryptRC4(Password, authTag)

	fmt.Printf(">>> [动作] 生成动态凭证...\n   - auth_tag: %s\n   - pwd(加密后): %s\n", authTag, dynamicPwd)
	
	data := url.Values{}
	data.Set("opr", "pwdLogin")
	data.Set("userName", UserName)
	data.Set("pwd", dynamicPwd)
	data.Set("auth_tag", authTag)
	data.Set("rememberPwd", "0")

	req, err := http.NewRequest("POST", LoginURL, strings.NewReader(data.Encode()))
	if err != nil {
		fmt.Println("构造请求失败:", err)
		return
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded; charset=UTF-8")
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36")
	req.Header.Set("X-Requested-With", "XMLHttpRequest")
	req.Header.Set("Referer", "http://192.168.3.1/ac_portal/default/pc.html")

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("发送失败:", err)
		return
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	result := string(body)
	if strings.Contains(result, "success") || strings.Contains(result, "true") {
		fmt.Println(">>> [成功] 登录成功！网络已恢复。")
	} else {
		fmt.Println(">>> [失败] 认证未通过，服务器响应:", result)
	}
}

func main() {
	fmt.Println("--- 路由器全天候保活程序 (24h/10min) ---")
	
	if CheckNetwork() {
		fmt.Println(">>> 当前网络正常，无需操作。")
	} else {
		fmt.Println(">>> 启动自检：断网，立即尝试登录...")
		DoLogin()
	}

	ticker := time.NewTicker(10 * time.Minute)
	defer ticker.Stop()

	fmt.Println(">>> 监控进程已启动 (策略：全天候，每10分钟检查)...")

	for range ticker.C {
		if !CheckNetwork() {
			now := time.Now().Format("15:04:05")
			fmt.Printf("[%s] 检测到断网，正在重连...\n", now)
			DoLogin()
		}
	}
}
