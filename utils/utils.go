/*
* @Author: thepoy
* @Email:  thepoy@aliyun.com
* @File Name: utils.go
* @Created:   2021-01-17 12:27:14
* @Modified:  2021-01-20 18:23:47
 */

package utils

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"path/filepath"
	"strings"
	"time"

	"gopkg.in/gomail.v2"
)

// 发送email的类型
const (
	RegisterMail uint = 0
)

// UserVisitInfo 用户访问时的信息
type UserVisitInfo struct {
	UA      UserAgent
	IP      string    `json:"ip"`
	Time    time.Time `json:"time"`
	Address string    `json:"address"`
}

// UserAgent ua信息
type UserAgent struct {
	Browser, Platform string
}

func parseUserAgent(ua string) UserAgent {
	var u UserAgent
	if strings.Contains(ua, "Firefox") {
		u.Browser = "Firefox"
	} else if strings.Contains(ua, "Edge") {
		u.Browser = "Edge"
	} else if strings.Contains(ua, "Opera") || strings.Contains(ua, "OPR") {
		u.Browser = "Opera"
	} else if strings.Contains(ua, "Chrome") {
		u.Browser = "Chrome"
	} else if strings.Contains(ua, "Safari") {
		u.Browser = "Safari"
	} else if strings.Contains(ua, "MSIE") || strings.Contains(ua, "Trident") {
		u.Browser = "IE"
	} else {
		u.Browser = "Unkown"
	}

	if strings.Contains(ua, "Windows") {
		u.Platform = "Windows"
	} else if strings.Contains(ua, "Linux") {
		u.Platform = "Linux"
	} else if strings.Contains(ua, "Android") {
		u.Platform = "Android"
	} else if strings.Contains(ua, "iPhone") {
		u.Platform = "iPhone"
	} else if strings.Contains(ua, "Mac OS") {
		u.Platform = "Mac OS"
	} else {
		u.Platform = "Unkown"
	}

	return u
}

func parseVisitInfo(ua, ip string) (*UserVisitInfo, error) {

	url := fmt.Sprintf("https://www.ip.cn/api/index?ip=%s&type=1", ip)
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var info UserVisitInfo
	err = json.Unmarshal(body, &info)
	if err != nil {
		return nil, err
	}

	info.UA = parseUserAgent(ua)

	return &info, nil
}

func loadHTML(t uint) (string, error) {
	switch t {
	case RegisterMail:
		path := filepath.Join(basicDir, "utils/email/register.html")
		b, err := ioutil.ReadFile(path)
		if err != nil {
			return "", err
		}
		html := strings.Replace(string(b), "{{ type }}", "注册", 1)
		return html, nil
	}
	return "", nil
}

// SendValidateCodeViaMail 发送邮件验证码
func SendValidateCodeViaMail(address, ua, ip string) error {
	// 生成6位随机验证码
	rnd := rand.New(rand.NewSource(time.Now().UnixNano()))
	vcode := fmt.Sprintf("%06v", rnd.Int31n(1000000))

	info, err := parseVisitInfo(ua, ip)
	if err != nil {
		return err
	}

	html, _ := loadHTML(0)
	html = strings.Replace(html, "{{ code }}", fmt.Sprint(vcode), 1)
	html = strings.Replace(html, "{{ browser }}", info.UA.Browser, 1)
	html = strings.Replace(html, "{{ platform }}", info.UA.Platform, 1)
	html = strings.Replace(html, "{{ address }}", info.Address, 1)
	html = strings.Replace(html, "{{ ip }}", ip, 1)
	now := time.Now()
	html = strings.Replace(html, "{{ time }}", now.Local().Format("15:04:02, 2006年01月02日(中国标准时间)"), 1)

	message := gomail.NewMessage()
	message.SetAddressHeader("From", email.Email, "无空的博客")
	message.SetHeader("To", address)
	message.SetHeader("Subject", "无空博客的验证码")
	message.SetBody("text/html", html)

	dial := gomail.NewDialer(email.Host, email.Port, email.Email, email.Password)
	if err := dial.DialAndSend(message); err != nil {
		fmt.Println("发送邮件时出错：", err)
		return err
	}

	if err := SetVerifyCodeOnRedis(address, vcode); err != nil {
		return err
	}

	log.Println("Mail sent successfully.")

	return nil
}
