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
	// Carrega a env
	err := godotenv.Load()
	if err != nil {
		fmt.Println("Arquivo .env não encontrado, usando variáveis de ambiente")
	}

	// Pega a config do SMTP
	smtpHost := os.Getenv("SMTP_HOST")
	smtpUser := os.Getenv("SMTP_USER")
	smtpPass := os.Getenv("SMTP_PASS")
	smtpPort, _ := strconv.Atoi(os.Getenv("SMTP_PORT"))

	if smtpHost == "" || smtpUser == "" || smtpPass == "" {
		panic("Configure SMTP_HOST, SMTP_USER, SMTP_PASS na .env ou variáveis de ambiente")
	}

	r := gin.Default()

	r.POST("/send-email", func(c *gin.Context) {
		var req struct {
			To string `json:"to" binding:"required"`
			Subject string `json:"subject" binding:"required"`
			Body string `json:"body" binding:"required"`
		}

		// Verifica se a requisição é válida
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// Envia o email
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
