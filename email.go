package main

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"gopkg.in/gomail.v2"
)

func EmailVerification(server *gin.Engine) {
	server.POST("/send_email", func(ctx *gin.Context) {
		sendEmailHandler(ctx.Writer, ctx.Request)
	})

	server.POST("/verify_code", func(ctx *gin.Context) {
		verifyCodeHandler(ctx.Writer, ctx.Request)
	})
}

func generateVerificationCode() (string, error) {
	// 生成一个随机的验证代码
	b := make([]byte, 32)
	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(b), nil
}

func sendVerificationEmail(to string, code string) error {
	m := gomail.NewMessage()
	m.SetHeader("From", "youremail@example.com")
	m.SetHeader("To", to)
	m.SetHeader("Subject", "Email Verification")
	m.SetBody("text/plain", fmt.Sprintf("Your verification code is: %s", code))

	d := gomail.NewDialer("smtp.example.com", 587, "yourusername", "yourpassword")

	// 发送邮件
	if err := d.DialAndSend(m); err != nil {
		return err
	}
	return nil
}

func sendEmailHandler(w http.ResponseWriter, r *http.Request) {
	email := r.FormValue("email")

	// 生成验证代码
	code, err := generateVerificationCode()
	if err != nil {
		http.Error(w, "Error generating verification code", http.StatusInternalServerError)
		return
	}

	// 保存验证代码与邮箱的关联 (通常是保存到数据库)
	// saveVerificationCode(email, code)

	// 发送验证邮件
	if err := sendVerificationEmail(email, code); err != nil {
		http.Error(w, "Error sending email", http.StatusInternalServerError)
		return
	}

	fmt.Fprintf(w, "Verification email sent to %s", email)
}

func verifyCodeHandler(w http.ResponseWriter, r *http.Request) {
	// email := r.FormValue("email")
	// code := r.FormValue("code")

	// 从数据库中获取存储的验证代码
	// storedCode, err := getStoredCode(email)
	// if err != nil {
	//     http.Error(w, "Error retrieving stored code", http.StatusInternalServerError)
	//     return
	// }

	// 比较用户输入的验证码与存储的验证代码
	// if code == storedCode {
	//     fmt.Fprintf(w, "Email verified successfully!")
	// } else {
	//     http.Error(w, "Invalid verification code", http.StatusUnauthorized)
	// }
}
