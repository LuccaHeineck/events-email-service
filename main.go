package main

import (
	"fmt"
	"net/http"
	"os"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"gopkg.in/gomail.v2"
)

func main() {
	_ = godotenv.Load()

	smtpHost := os.Getenv("SMTP_HOST")
	smtpUser := os.Getenv("SMTP_USER")
	smtpPass := os.Getenv("SMTP_PASS")
	smtpPort, _ := strconv.Atoi(os.Getenv("SMTP_PORT"))
	gatewayKey := os.Getenv("GATEWAY_KEY") // chave secreta compartilhada

	if smtpHost == "" || smtpUser == "" || smtpPass == "" || gatewayKey == "" {
		panic("Configure SMTP_* e GATEWAY_KEY na .env")
	}

	r := gin.Default()

	// Middleware para verificar se a requisição veio do gateway
	r.Use(func(c *gin.Context) {
		key := c.GetHeader("X-Gateway-Key")
		if key == "" || key != gatewayKey {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "Acesso negado"})
			return
		}
		c.Next()
	})

	r.POST("/send-email", func(c *gin.Context) {
		var req struct {
			To      string `json:"to" binding:"required"`
			Subject string `json:"subject" binding:"required"`
			Body    string `json:"body" binding:"required"`
		}

		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		if err := sendEmail(smtpHost, smtpPort, smtpUser, smtpPass, req.To, req.Subject, req.Body); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "Email enviado com sucesso!"})
	})

	r.Run(":8082")
}

func sendEmail(host string, port int, user, pass, to, subject, body string) error {
	m := gomail.NewMessage()
	m.SetHeader("From", user)
	m.SetHeader("To", to)
	m.SetHeader("Subject", subject)
	m.SetBody("text/plain", body)

	d := gomail.NewDialer(host, port, user, pass)
	return d.DialAndSend(m)
}
