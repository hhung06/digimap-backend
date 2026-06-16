package firebase

import (
	"context"
	"fmt"

	gofirebase "firebase.google.com/go/v4"
	"firebase.google.com/go/v4/messaging"
	"google.golang.org/api/option"

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

// FCMPusher sends real push notifications via the Firebase Admin SDK.
type FCMPusher struct {
	client *messaging.Client
}

// NewFCMPusher initialises a real FCM sender from a service-account credentials file.
func NewFCMPusher(ctx context.Context, credentialsPath string) (Pusher, error) {
	app, err := gofirebase.NewApp(ctx, nil, option.WithCredentialsFile(credentialsPath))
	if err != nil {
		return nil, fmt.Errorf("init firebase app: %w", err)
	}
	client, err := app.Messaging(ctx)
	if err != nil {
		return nil, fmt.Errorf("init firebase messaging: %w", err)
	}
	return &FCMPusher{client: client}, nil
}

func (p *FCMPusher) Send(ctx context.Context, msg Message) (string, error) {
	m := &messaging.Message{
		Notification: &messaging.Notification{Title: msg.Title, Body: msg.Body},
		Data:         msg.Data,
	}
	if msg.Token != "" {
		m.Token = msg.Token
	} else {
		m.Topic = msg.Topic
	}
	return p.client.Send(ctx, m)
}

func (p *FCMPusher) SendMulticast(ctx context.Context, tokens []string, title, body string, data map[string]string) (int, int, error) {
	resp, err := p.client.SendEachForMulticast(ctx, &messaging.MulticastMessage{
		Notification: &messaging.Notification{Title: title, Body: body},
		Data:         data,
		Tokens:       tokens,
	})
	if err != nil {
		return 0, 0, err
	}
	return resp.SuccessCount, resp.FailureCount, nil
}
