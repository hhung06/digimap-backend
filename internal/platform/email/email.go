package email

import (
	"context"
	"fmt"

	applog "github.com/hhung06/digimap-backend/log"
)

// Sender defines the contract for sending transactional emails.
type Sender interface {
	SendPasswordReset(ctx context.Context, toEmail, resetToken string) error
	SendInvitation(ctx context.Context, toEmail, venueName, inviteToken string) error
}

// LogSender is a development stub that logs emails instead of sending them.
// Swap for an SES implementation in production.
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
