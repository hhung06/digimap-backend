package cmd

import (
	"context"
	"fmt"
	"time"

	"github.com/spf13/cobra"

	"github.com/hhung06/digimap-backend/config"
	"github.com/hhung06/digimap-backend/internal/platform/database"
	"github.com/hhung06/digimap-backend/internal/platform/firebase"
	postgresrepo "github.com/hhung06/digimap-backend/internal/repository/postgres"
	"github.com/hhung06/digimap-backend/internal/service"
	applog "github.com/hhung06/digimap-backend/log"
)

var sendScheduledNotificationsCmd = &cobra.Command{
	Use:   "send-scheduled-notifications",
	Short: "Send due scheduled notifications",
	RunE:  runSendScheduledNotifications,
}

func init() {
	rootCmd.AddCommand(sendScheduledNotificationsCmd)
}

func runSendScheduledNotifications(_ *cobra.Command, _ []string) error {
	cfg, err := config.Load()
	if err != nil {
		return fmt.Errorf("load config: %w", err)
	}

	logger := applog.NewLogger(cfg)
	defer applog.Sync(logger)
	pool, err := database.NewPool(context.Background(), cfg.Database)
	if err != nil {
		return fmt.Errorf("connect database: %w", err)
	}
	defer pool.Close()

	notificationRepo := postgresrepo.NewNotificationRepository(pool)
	pusher := firebase.NewLogPusher(logger)
	notificationSvc := service.NewNotificationService(notificationRepo, pusher)

	sent, err := notificationSvc.SendDueScheduled(context.Background(), time.Now())
	if err != nil {
		return fmt.Errorf("send due scheduled notifications: %w", err)
	}

	fmt.Printf("sent %d scheduled notification(s)\n", sent)
	return nil
}
