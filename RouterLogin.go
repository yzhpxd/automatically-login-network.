package main

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
)

// --- 配置区域 ---
const (
	LoginURL = "http://192.168.3.1/ac_portal/login.php"
	CheckURL = "https://www.baidu.com"
	
	UserName = "321081119"
	
	// 固定凭证 (基于你抓包成功的那个)
	FixedPwd     = "e0f67a0ec942"
	FixedAuthTag = "1770107462614"
)

func CheckNetwork() bool {
	client := http.Client{Timeout: 5 * time.Second}
	resp, err := client.Get(CheckURL)
	if err != nil {
		return false
	}
	defer resp.Body.Close()
	return resp.StatusCode == 200
}

func DoLogin() {
	fmt.Println(">>> [动作] 发送固定凭证登录...")
	data := url.Values{}
	data.Set("opr", "pwdLogin")
	data.Set("userName", UserName)
	data.Set("pwd", FixedPwd)
	data.Set("auth_tag", FixedAuthTag)
	data.Set("rememberPwd", "0")

	req, err := http.NewRequest("POST", LoginURL, strings.NewReader(data.Encode()))
	if err != nil {
		fmt.Println("构造请求失败:", err)
		return
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded; charset=UTF-8")
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/144.0.0.0 Safari/537.36")
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
		fmt.Println(">>> [成功] 登录成功！")
	} else {
		fmt.Println(">>> [服务器响应] ", result)
	}
}

func main() {
	fmt.Println("--- 路由器全天候保活程序 (24h/10min) ---")
	
	if CheckNetwork() {
		fmt.Println(">>> 当前网络正常，无需操作。")
	} else {
		fmt.Println(">>> 启动自检：断网，立即重连...")
		DoLogin()
	}

	// 每10分钟检查一次
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
