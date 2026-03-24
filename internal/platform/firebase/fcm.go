package firebase

import (
	"context"
	"fmt"

	applog "github.com/hhung06/digimap-backend/log"
)

// Message is the payload to push to FCM.
type Message struct {
	// Either Token (single device) or Topic (subscribed group)
	Token string
	Topic string

	Title string
	Body  string
	Data  map[string]string
}

// Pusher defines the contract for sending push notifications.
type Pusher interface {
	// Send dispatches a single FCM message. Returns the FCM message ID.
	Send(ctx context.Context, msg Message) (string, error)
	// SendMulticast sends to up to 500 device tokens at once.
	SendMulticast(ctx context.Context, tokens []string, title, body string, data map[string]string) (successCount int, failureCount int, err error)
}

// LogPusher is a development stub that logs instead of actually sending to FCM.
type LogPusher struct {
	logger applog.Logger
}

// NewLogPusher returns a Pusher that logs all outgoing messages.
func NewLogPusher(logger applog.Logger) Pusher {
	return &LogPusher{logger: logger}
}

func (p *LogPusher) Send(_ context.Context, msg Message) (string, error) {
	target := msg.Token
	if target == "" {
		target = "topic:" + msg.Topic
	}
	p.logger.Infof("[FCM] send to=%s title=%q body=%q", target, msg.Title, msg.Body)
	fmt.Printf("\n[DEV FCM] to=%s title=%q body=%q\n", target, msg.Title, msg.Body)
	return "dev-message-id", nil
}

func (p *LogPusher) SendMulticast(_ context.Context, tokens []string, title, body string, _ map[string]string) (int, int, error) {
	p.logger.Infof("[FCM] multicast tokens=%d title=%q", len(tokens), title)
	fmt.Printf("\n[DEV FCM] multicast to %d tokens title=%q body=%q\n", len(tokens), title, body)
	return len(tokens), 0, nil
}
