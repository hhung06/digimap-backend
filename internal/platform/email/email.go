package email

import (
	"context"
	"crypto/tls"
	"fmt"
	"net/smtp"
	"strings"

	"github.com/hhung06/digimap-backend/config"
	applog "github.com/hhung06/digimap-backend/log"
)

// Sender defines the contract for sending transactional emails.
type Sender interface {
	SendPasswordReset(ctx context.Context, toEmail, resetToken string) error
	SendInvitation(ctx context.Context, toEmail, venueName, inviteToken string) error
}

// SMTPSender sends email via SMTP (AWS SES SMTP endpoint or any SMTP provider).
type SMTPSender struct {
	cfg config.SMTPConfig
}

// NewSMTPSender creates a Sender backed by an SMTP server.
func NewSMTPSender(cfg config.SMTPConfig) Sender {
	return &SMTPSender{cfg: cfg}
}

func (s *SMTPSender) SendPasswordReset(_ context.Context, toEmail, resetToken string) error {
	subject := "Password Reset Request"
	body := fmt.Sprintf("You requested a password reset.\n\nYour reset token: %s\n\nIf you did not request this, please ignore this email.", resetToken)
	return s.send(toEmail, subject, body)
}

func (s *SMTPSender) SendInvitation(_ context.Context, toEmail, venueName, inviteToken string) error {
	subject := fmt.Sprintf("You've been invited to %s on Digimap", venueName)
	body := fmt.Sprintf("You have been invited to join %s on Digimap.\n\nYour invitation token: %s", venueName, inviteToken)
	return s.send(toEmail, subject, body)
}

func (s *SMTPSender) send(toEmail, subject, body string) error {
	addr := fmt.Sprintf("%s:%d", s.cfg.Host, s.cfg.Port)
	auth := smtp.PlainAuth("", s.cfg.User, s.cfg.Password, s.cfg.Host)
	msg, err := buildMessage(s.cfg.FromEmail, toEmail, subject, body)
	if err != nil {
		return err
	}

	if s.cfg.UseTLS {
		// Implicit TLS (port 465)
		return sendWithTLS(addr, s.cfg.Host, auth, s.cfg.FromEmail, toEmail, msg)
	}
	// STARTTLS (port 587) — smtp.SendMail upgrades automatically
	return smtp.SendMail(addr, auth, s.cfg.FromEmail, []string{toEmail}, msg)
}

// sendWithTLS dials with implicit TLS (port 465).
func sendWithTLS(addr, host string, auth smtp.Auth, from, to string, msg []byte) error {
	conn, err := tls.Dial("tcp", addr, &tls.Config{ServerName: host})
	if err != nil {
		return fmt.Errorf("tls dial: %w", err)
	}
	defer conn.Close()

	c, err := smtp.NewClient(conn, host)
	if err != nil {
		return fmt.Errorf("smtp client: %w", err)
	}
	defer c.Close()

	if err := c.Auth(auth); err != nil {
		return fmt.Errorf("smtp auth: %w", err)
	}
	if err := c.Mail(from); err != nil {
		return fmt.Errorf("smtp mail from: %w", err)
	}
	if err := c.Rcpt(to); err != nil {
		return fmt.Errorf("smtp rcpt to: %w", err)
	}
	w, err := c.Data()
	if err != nil {
		return fmt.Errorf("smtp data: %w", err)
	}
	if _, err := w.Write(msg); err != nil {
		return fmt.Errorf("smtp write: %w", err)
	}
	return w.Close()
}

func buildMessage(from, to, subject, body string) ([]byte, error) {
	for name, val := range map[string]string{"from": from, "to": to, "subject": subject} {
		if strings.ContainsAny(val, "\r\n") {
			return nil, fmt.Errorf("email header injection in %s", name)
		}
	}
	var sb strings.Builder
	sb.WriteString("From: " + from + "\r\n")
	sb.WriteString("To: " + to + "\r\n")
	sb.WriteString("Subject: " + subject + "\r\n")
	sb.WriteString("MIME-Version: 1.0\r\n")
	sb.WriteString("Content-Type: text/plain; charset=UTF-8\r\n")
	sb.WriteString("\r\n")
	sb.WriteString(body)
	return []byte(sb.String()), nil
}

// LogSender is a development stub that logs emails instead of sending them.
type LogSender struct {
	logger applog.Logger
}

// NewLogSender returns a Sender that logs all outgoing emails.
func NewLogSender(logger applog.Logger) Sender {
	return &LogSender{logger: logger}
}

func (s *LogSender) SendPasswordReset(_ context.Context, toEmail, resetToken string) error {
	s.logger.Infof("[EMAIL] password-reset to=%s token=%s", toEmail, resetToken)
	fmt.Printf("\n[DEV EMAIL] Password reset for %s\n  Token: %s\n\n", toEmail, resetToken)
	return nil
}

func (s *LogSender) SendInvitation(_ context.Context, toEmail, venueName, inviteToken string) error {
	s.logger.Infof("[EMAIL] invitation to=%s venue=%s token=%s", toEmail, venueName, inviteToken)
	fmt.Printf("\n[DEV EMAIL] Invitation for %s to venue %s\n  Token: %s\n\n", toEmail, venueName, inviteToken)
	return nil
}
