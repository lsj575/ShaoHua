package api

import (
	"api/user-web/forms"
	"api/user-web/global"
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"go.uber.org/zap"
	"gopkg.in/gomail.v2"
	"math/rand"
	"net/http"
	"strings"
	"time"
)

const EmailName = "BBQ"
const EmailSubject = "BBQ注册验证码"

type Mail struct {
	senderId string
	toIds    []string
	subject  string
	body     string
}

// 验证redis中字符串类型数据是否匹配
func VerifyString(key, value string) bool {
	rdb := redis.NewClient(&redis.Options{
		Addr: fmt.Sprintf("%s:%d", global.ServiceConfig.RedisConfig.Host, global.ServiceConfig.RedisConfig.Port),
	})
	code, err := rdb.Get(context.Background(), key).Result()
	if err == redis.Nil || code != value {
		return false
	}
	return true
}

func GenValidateCode(width int) string {
	numeric := [10]byte{0, 1, 2, 3, 4, 5, 6, 7, 8, 9}
	r := len(numeric)
	rand.Seed(time.Now().UnixNano())

	var sb strings.Builder
	for i := 0; i < width; i++ {
		_, _ = fmt.Fprintf(&sb, "%d", numeric[rand.Intn(r)])
	}
	return sb.String()
}

func SendEmailVerificationCode(ctx *gin.Context) {
	emailVerificationForm := forms.EmailVerificationForm{}
	if err := ctx.ShouldBind(&emailVerificationForm); err != nil {
		HandleValidatorError(ctx, err)
		return
	}

	code := GenValidateCode(6)

	if emailVerificationForm.Username != "" {
		emailVerificationForm.Username = emailVerificationForm.Username + "，"
	}

	body := `
        <html>
        <body>
        <h3>
        ` + emailVerificationForm.Username + `您好：
        </h3>
        非常感谢您使用` + EmailName + `，您的邮箱验证码为：<br/>
        <b>` + code + `</b><br/>
        此验证码有效期30分钟，请妥善保存。<br/>
        如果这不是您本人的操作，请忽略本邮件。<br/>
        </body>
        </html>
        `

	if err := SendToMail(emailVerificationForm.Email, EmailSubject, body); err != nil {
		ctx.JSON(http.StatusInternalServerError, global.JsonError(global.ServerAPIInternalError, "邮件发送失败"))
		return
	}

	rdb := redis.NewClient(&redis.Options{
		Addr: fmt.Sprintf("%s:%d", global.ServiceConfig.RedisConfig.Host, global.ServiceConfig.RedisConfig.Port),
	})
	rdb.Set(context.Background(), emailVerificationForm.Email, code, time.Second*300)

	ctx.JSON(http.StatusOK, global.JsonSuccess("邮件发送成功", []int{}))
}

func (mail *Mail) BuildMessage() string {
	message := ""
	message += fmt.Sprintf("From: %s<%s>\r\n", EmailName, mail.senderId)
	if len(mail.toIds) > 0 {
		message += fmt.Sprintf("To: %s\r\n", strings.Join(mail.toIds, ";"))
	}

	message += fmt.Sprintf("Subject: %s\r\n", mail.subject)

	message += "Content-Type: text/html; charset=UTF-8"

	message += "\r\n\r\n" + mail.body

	return message
}

func SendToMail(to, subject, body string) error {
	m := gomail.NewMessage()
	m.SetHeader("From", global.ServiceConfig.SMTPServerConfig.User)
	m.SetHeader("To", strings.Split(to, ";")...)
	m.SetHeader("Subject", subject)
	m.SetBody("text/html", body)

	d := gomail.NewDialer("smtp.qq.com", 587,
		global.ServiceConfig.SMTPServerConfig.User, global.ServiceConfig.SMTPServerConfig.Password)

	// 发送邮件
	if err := d.DialAndSend(m); err != nil {
		zap.S().Errorf("DialAndSend err %v:", err)
		return err
	}
	zap.S().Info("Mail sent successfully")

	return nil
}
