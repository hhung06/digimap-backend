package cmd

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/spf13/cobra"

	"github.com/hhung06/digimap-backend/config"
	"github.com/hhung06/digimap-backend/internal/handler"
	"github.com/hhung06/digimap-backend/internal/platform/cache"
	"github.com/hhung06/digimap-backend/internal/platform/database"
	"github.com/hhung06/digimap-backend/internal/platform/email"
	"github.com/hhung06/digimap-backend/internal/platform/firebase"
	"github.com/hhung06/digimap-backend/internal/platform/storage"
	postgresrepo "github.com/hhung06/digimap-backend/internal/repository/postgres"
	"github.com/hhung06/digimap-backend/internal/service"
	applog "github.com/hhung06/digimap-backend/log"
)

var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Start the HTTP server",
	Long:  `Starts the Digimap REST API server on the configured port.`,
	RunE:  runServe,
}

func runServe(_ *cobra.Command, _ []string) error {
	cfg, err := config.Load()
	if err != nil {
		return fmt.Errorf("load config: %w", err)
	}

	logger := applog.NewLogger(cfg)
	logger.Infof("starting digimap-backend env=%s port=%d", cfg.App.Environment, cfg.Server.Port)

	ctx := context.Background()

	// ── Infrastructure ────────────────────────────────────────────────────
	pool, err := database.NewPool(ctx, cfg.Database)
	if err != nil {
		return fmt.Errorf("connect postgres: %w", err)
	}
	defer pool.Close()
	logger.Info("postgres connected")

	redisClient, err := cache.NewRedisClient(cfg.Redis)
	if err != nil {
		return fmt.Errorf("connect redis: %w", err)
	}
	defer redisClient.Close()
	logger.Info("redis connected")

	// ── Repositories ──────────────────────────────────────────────────────
	userRepo := postgresrepo.NewUserRepository(pool)
	tokenRepo := postgresrepo.NewTokenRepository(pool)
	customerRepo := postgresrepo.NewCustomerRepository(pool)
	venueRepo := postgresrepo.NewVenueRepository(pool)
	levelRepo := postgresrepo.NewLevelRepository(pool)
	locationCategoryRepo := postgresrepo.NewLocationCategoryRepository(pool)
	amenityRepo := postgresrepo.NewAmenityRepository(pool)
	locationRepo := postgresrepo.NewLocationRepository(pool)
	productRepo := postgresrepo.NewProductRepository(pool)
	eventRepo := postgresrepo.NewEventRepository(pool)
	notificationRepo := postgresrepo.NewNotificationRepository(pool)
	surveyRepo := postgresrepo.NewSurveyRepository(pool)
	beaconRepo := postgresrepo.NewBeaconRepository(pool)

	// ── Platform services ─────────────────────────────────────────────────
	mailer := email.NewLogSender(logger)
	storer := storage.NewLogStorer()
	pusher := firebase.NewLogPusher(logger)

	// ── Application services ──────────────────────────────────────────────
	authSvc := service.NewAuthService(userRepo, tokenRepo, mailer, cfg.JWT)
	customerSvc := service.NewCustomerService(customerRepo)
	venueSvc := service.NewVenueService(venueRepo, levelRepo, pool)
	levelSvc := service.NewLevelService(levelRepo)
	locationCategorySvc := service.NewLocationCategoryService(locationCategoryRepo)
	amenitySvc := service.NewAmenityService(amenityRepo)
	locationSvc := service.NewLocationService(locationRepo)
	productSvc := service.NewProductService(productRepo)
	storageSvc := service.NewStorageService(storer)
	eventSvc := service.NewEventService(eventRepo)
	userSvc := service.NewUserService(userRepo, mailer)
	notificationSvc := service.NewNotificationService(notificationRepo, pusher)
	surveySvc := service.NewSurveyService(surveyRepo)
	beaconSvc := service.NewBeaconService(beaconRepo)

	// ── HTTP server ───────────────────────────────────────────────────────
	deps := handler.Dependencies{
		AuthService:             authSvc,
		CustomerService:         customerSvc,
		VenueService:            venueSvc,
		LevelService:            levelSvc,
		LocationCategoryService: locationCategorySvc,
		AmenityService:          amenitySvc,
		LocationService:         locationSvc,
		ProductService:          productSvc,
		StorageService:          storageSvc,
		EventService:            eventSvc,
		UserService:             userSvc,
		NotificationService:     notificationSvc,
		SurveyService:           surveySvc,
		BeaconService:           beaconSvc,
		UserRepo:                userRepo,
	}
	router := handler.NewRouter(cfg, logger, deps)

	srv := &http.Server{
		Addr:         fmt.Sprintf(":%d", cfg.Server.Port),
		Handler:      router,
		ReadTimeout:  cfg.Server.ReadTimeout,
		WriteTimeout: cfg.Server.WriteTimeout,
		IdleTimeout:  cfg.Server.IdleTimeout,
	}

	go func() {
		logger.Infof("listening on :%d", cfg.Server.Port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Fatalf("server error: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	logger.Info("shutting down server...")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	if err := srv.Shutdown(shutdownCtx); err != nil {
		return fmt.Errorf("server shutdown: %w", err)
	}

	logger.Info("server stopped")
	return nil
}
