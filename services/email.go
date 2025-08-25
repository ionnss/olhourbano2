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

// GetStatusEmailTemplate returns the email template for report status
func GetStatusEmailTemplate(reportID int, status string) EmailTemplate {
	subject := fmt.Sprintf("Olho Urbano - Denúncia #%d Atualizada", reportID)

	body := fmt.Sprintf(`
Olá,

Sua denúncia foi atualizada para o status: %s

Para acompanhar o status da sua denúncia, acesse:
https://olhourbano.com.br/report/%d

Obrigado por contribuir para uma cidade melhor!

--
Equipe Olho Urbano
olhourbano.contato@gmail.com
`, status, reportID)

	return EmailTemplate{
		Subject: subject,
		Body:    body,
	}
}

// GetCommentNotificationEmailTemplate returns the email template for comment notifications
func GetCommentNotificationEmailTemplate(reportID int, commenterName, commentContent string) EmailTemplate {
	subject := fmt.Sprintf("Olho Urbano - Novo Comentário na Denúncia #%d", reportID)

	body := fmt.Sprintf(`
Olá,

Sua denúncia recebeu um novo comentário!

Detalhes:
- Denúncia: #%d
- Comentário de: %s
- Conteúdo: "%s"

Para visualizar o comentário e responder, acesse:
https://olhourbano.com.br/report/%d

Obrigado por contribuir para uma cidade melhor!

--
Equipe Olho Urbano
olhourbano.contato@gmail.com
`, reportID, commenterName, commentContent, reportID)

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

// SendStatusEmail sends a status email for a report
func SendStatusEmail(email string, reportID int, status string) {
	template := GetStatusEmailTemplate(reportID, status)

	err := SendEmail(email, template)
	if err != nil {
		log.Printf("Erro ao enviar email de status para %s: %v", email, err)
	}
}

// SendCommentNotificationEmail sends a notification email when a comment is posted
func SendCommentNotificationEmail(email string, reportID int, commenterName, commentContent string) {
	template := GetCommentNotificationEmailTemplate(reportID, commenterName, commentContent)

	err := SendEmail(email, template)
	if err != nil {
		log.Printf("Erro ao enviar email de notificação de comentário para %s: %v", email, err)
	}
}
