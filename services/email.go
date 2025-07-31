package services

import (
	"fmt"
	"log"
	"net/smtp"
	"olhourbano2/config"
)

// EmailTemplate represents an email template
type EmailTemplate struct {
	Subject string
	Body    string
}

// GetConfirmationEmailTemplate returns the email template for report confirmation
func GetConfirmationEmailTemplate(reportID int, categoryName string) EmailTemplate {
	subject := fmt.Sprintf("Olho Urbano - Denúncia #%d Recebida", reportID)

	body := fmt.Sprintf(`
Olá,

Sua denúncia foi recebida com sucesso!

Detalhes da Denúncia:
- Número: #%d
- Categoria: %s
- Status: Pendente de Análise

Sua denúncia será analisada pela nossa equipe e você receberá atualizações sobre o andamento.

Para acompanhar o status da sua denúncia, acesse:
https://olhourbano.com.br/report/%d

Obrigado por contribuir para uma cidade melhor!

--
Equipe Olho Urbano
olhourbano.contato@gmail.com
`, reportID, categoryName, reportID)

	return EmailTemplate{
		Subject: subject,
		Body:    body,
	}
}

// SendEmail sends an email using SMTP configuration
func SendEmail(to string, template EmailTemplate) error {
	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		return fmt.Errorf("erro ao carregar configuração: %v", err)
	}

	// Email configuration
	from := cfg.SMTPUsername
	password := cfg.SMTPPassword
	smtpHost := cfg.SMTPHost
	smtpPort := fmt.Sprintf("%d", cfg.SMTPPort)

	// Authentication
	auth := smtp.PlainAuth("", from, password, smtpHost)

	// Message
	message := fmt.Sprintf("From: %s\r\nTo: %s\r\nSubject: %s\r\n\r\n%s", from, to, template.Subject, template.Body)

	// Send email
	err = smtp.SendMail(smtpHost+":"+smtpPort, auth, from, []string{to}, []byte(message))
	if err != nil {
		return fmt.Errorf("erro ao enviar email: %v", err)
	}

	log.Printf("Email de confirmação enviado para: %s", to)
	return nil
}

// SendConfirmationEmail sends a confirmation email for a report
func SendConfirmationEmail(email string, reportID int, categoryName string) {
	template := GetConfirmationEmailTemplate(reportID, categoryName)

	err := SendEmail(email, template)
	if err != nil {
		log.Printf("Erro ao enviar email de confirmação para %s: %v", email, err)
	}
}
