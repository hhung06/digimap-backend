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

var activateScheduledSurveysCmd = &cobra.Command{
	Use:   "activate-scheduled-surveys",
	Short: "Activate due surveys and close ended surveys",
	RunE:  runActivateScheduledSurveys,
}

func init() {
	rootCmd.AddCommand(activateScheduledSurveysCmd)
}

func runActivateScheduledSurveys(_ *cobra.Command, _ []string) error {
	cfg, err := config.Load()
	if err != nil {
		return fmt.Errorf("load config: %w", err)
	}

	logger := applog.NewLogger(cfg)
	pool, err := database.NewPool(context.Background(), cfg.Database)
	if err != nil {
		return fmt.Errorf("connect database: %w", err)
	}
	defer pool.Close()

	notificationRepo := postgresrepo.NewNotificationRepository(pool)
	surveyRepo := postgresrepo.NewSurveyRepository(pool)
	pusher := firebase.NewLogPusher(logger)
	notificationSvc := service.NewNotificationService(notificationRepo, pusher)
	surveySvc := service.NewSurveyService(surveyRepo, notificationRepo, notificationSvc)

	activated, closed, err := surveySvc.ProcessScheduledTransitions(context.Background(), time.Now())
	if err != nil {
		return fmt.Errorf("process scheduled survey transitions: %w", err)
	}

	fmt.Printf("activated %d survey(s), closed %d survey(s)\n", activated, closed)
	return nil
}
