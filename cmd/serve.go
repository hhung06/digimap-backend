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

	"github.com/hhung06/digimap-backend/cmd/tenants"
	"github.com/hhung06/digimap-backend/config"
	"github.com/hhung06/digimap-backend/internal/enricher"
	"github.com/hhung06/digimap-backend/internal/handler"
	"github.com/hhung06/digimap-backend/internal/platform/cache"
	"github.com/hhung06/digimap-backend/internal/platform/cdn"
	"github.com/hhung06/digimap-backend/internal/platform/crypto"
	"github.com/hhung06/digimap-backend/internal/platform/database"
	"github.com/hhung06/digimap-backend/internal/platform/email"
	"github.com/hhung06/digimap-backend/internal/platform/firebase"
	"github.com/hhung06/digimap-backend/internal/platform/search"
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
	defer applog.Sync(logger)
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
	locationRepo := postgresrepo.NewLocationRepository(pool)
	productRepo := postgresrepo.NewProductRepository(pool)
	eventRepo := postgresrepo.NewEventRepository(pool)
	notificationRepo := postgresrepo.NewNotificationRepository(pool)
	surveyRepo := postgresrepo.NewSurveyRepository(pool)
	beaconRepo := postgresrepo.NewBeaconRepository(pool)
	connectionRepo := postgresrepo.NewConnectionRepository(pool)
	adRepo := postgresrepo.NewAdRepository(pool)
	articleRepo := postgresrepo.NewArticleRepository(pool)
	couponRepo := postgresrepo.NewCouponRepository(pool)
	couponUserRepo := postgresrepo.NewCouponUserRepository(pool)
	videoRepo := postgresrepo.NewVideoRepository(pool)
	tagRepo := postgresrepo.NewTagRepository(pool)
	eventLogRepo := postgresrepo.NewEventLogRepository(pool)
	searchQueryRepo := postgresrepo.NewSearchQueryRepository(pool)
	snapshotRepo := postgresrepo.NewSnapshotRepository(pool)
	levelBundleRepo := postgresrepo.NewLevelBundleRepository(pool)
	assetRepo := postgresrepo.NewAssetRepository(pool)
	levelTypeRepo := postgresrepo.NewLevelTypeRepository(pool)
	themeRepo := postgresrepo.NewThemeRepository(pool)
	productPlazaRepo := postgresrepo.NewProductPlazaRepository(pool)
	appUserRepo := postgresrepo.NewAppUserRepository(pool)
	languageRepo := postgresrepo.NewLanguageRepository(pool)
	appVersionRepo := postgresrepo.NewAppVersionRepository(pool)

	// ── Platform services ─────────────────────────────────────────────────
	var mailer email.Sender
	if cfg.SMTP.Host != "" {
		mailer = email.NewSMTPSender(cfg.SMTP)
		logger.Infof("SMTP mailer initialised host=%s port=%d", cfg.SMTP.Host, cfg.SMTP.Port)
	} else {
		mailer = email.NewLogSender(logger)
		logger.Info("SMTP mailer: using log stub (EMAIL_HOST not set)")
	}

	var assetStorer, snapshotStorer storage.Storer
	var invalidator cdn.Invalidator
	if cfg.AWS.S3AssetsBucket != "" {
		assetStorer, err = storage.NewS3Storer(cfg.AWS, cfg.AWS.S3AssetsBucket)
		if err != nil {
			return fmt.Errorf("init s3 assets storer: %w", err)
		}
		snapshotStorer, err = storage.NewS3Storer(cfg.AWS, cfg.AWS.S3SnapshotBucket)
		if err != nil {
			return fmt.Errorf("init s3 snapshot storer: %w", err)
		}
		invalidator, err = cdn.NewCloudFrontInvalidator(cfg.AWS)
		if err != nil {
			return fmt.Errorf("init cloudfront invalidator: %w", err)
		}
		logger.Info("S3 storers + CloudFront invalidator initialised")
	} else {
		assetStorer = storage.NewLogStorer()
		snapshotStorer = storage.NewLogStorer()
		invalidator = cdn.NewLogInvalidator()
		logger.Info("storage/CDN: using log stubs (AWS_S3_ASSETS_BUCKET not set)")
	}

	var pusher firebase.Pusher
	if cfg.Firebase.CredentialsPath != "" {
		pusher, err = firebase.NewFCMPusher(ctx, cfg.Firebase.CredentialsPath)
		if err != nil {
			return fmt.Errorf("init firebase pusher: %w", err)
		}
		logger.Info("firebase FCM pusher initialised")
	} else {
		pusher = firebase.NewLogPusher(logger)
		logger.Info("firebase FCM pusher: using log stub (FIREBASE_CREDENTIALS_PATH not set)")
	}

	var phoneEncryptor service.PhoneEncryptor
	if cfg.App.DataEncryptionKey != "" {
		phoneEncryptor, err = crypto.NewAESGCMEncryption(cfg.App.DataEncryptionKey)
		if err != nil {
			return fmt.Errorf("init phone encryptor: %w", err)
		}
		logger.Info("visitor phone encryption enabled")
	} else {
		logger.Info("visitor phone encryption disabled (DATA_ENCRYPTION_KEY not set)")
	}

	var searcher search.Searcher
	if cfg.OpenSearch.Endpoint != "" {
		searcher, err = search.NewSearcher(cfg.OpenSearch.Endpoint, cfg.OpenSearch.User, cfg.OpenSearch.Password)
		if err != nil {
			return fmt.Errorf("init opensearch: %w", err)
		}
		logger.Infof("opensearch searcher initialised endpoint=%s", cfg.OpenSearch.Endpoint)
	} else {
		searcher = search.NewLogSearcher()
		logger.Info("opensearch: using log stub (OPENSEARCH_ENDPOINT not set)")
	}

	// ── Application services ──────────────────────────────────────────────
	authSvc := service.NewAuthService(userRepo, tokenRepo, mailer, cfg.JWT)
	customerSvc := service.NewCustomerService(customerRepo)
	venueSvc := service.NewVenueService(venueRepo, levelRepo, pool)
	levelSvc := service.NewLevelService(levelRepo)
	locationCategorySvc := service.NewLocationCategoryService(locationCategoryRepo)
	appVersionSvc := service.NewAppVersionService(appVersionRepo, assetStorer, invalidator, cfg.App.Environment)
	snapshotSvc := service.NewSnapshotService(snapshotRepo, venueRepo, languageRepo, locationRepo, locationCategoryRepo, productRepo, themeRepo, snapshotStorer, assetStorer, invalidator, appVersionSvc, cfg.App.Environment, searcher, logger)
	locationSvc := service.NewLocationService(locationRepo, venueRepo, assetStorer, invalidator, appVersionSvc, cfg.App.Environment)
	productSvc := service.NewProductService(productRepo)
	mediaSvc := service.NewMediaService(assetStorer, cfg.App.Environment)
	storageSvc := service.NewStorageService(assetStorer)
	eventSvc := service.NewEventService(eventRepo)
	userSvc := service.NewUserService(userRepo, mailer)
	notificationSvc := service.NewNotificationService(notificationRepo, pusher)
	surveySvc := service.NewSurveyService(surveyRepo, notificationRepo, notificationSvc)
	beaconSvc := service.NewBeaconService(beaconRepo)
	connectionSvc := service.NewConnectionService(connectionRepo)
	adSvc := service.NewAdvertisementService(adRepo)
	articleSvc := service.NewArticleService(articleRepo, mediaSvc)
	couponSvc := service.NewCouponService(couponRepo, couponUserRepo)
	videoSvc := service.NewVideoService(videoRepo)
	tagSvc := service.NewTagService(tagRepo)
	analyticsSvc := service.NewAnalyticsService(eventLogRepo, searchQueryRepo, venueRepo, redisClient)
	levelBundleSvc := service.NewLevelBundleService(levelBundleRepo, snapshotRepo, assetStorer, cfg.App.Environment)
	assetSvc := service.NewAssetService(assetRepo)
	levelTypeSvc := service.NewLevelTypeService(levelTypeRepo)
	themeSvc := service.NewThemeService(themeRepo, venueRepo, assetStorer, invalidator, cfg.App.Environment)
	productPlazaSvc := service.NewProductPlazaService(productPlazaRepo)
	languageSvc := service.NewLanguageService(languageRepo)
	memoSvc := service.NewMemoService(locationRepo, venueRepo, assetStorer, invalidator, appVersionSvc, cfg.App.Environment)

	enricherRegistry := enricher.NewRegistry(venueRepo)
	tenants.RegisterAll(enricherRegistry)

	// ── HTTP server ───────────────────────────────────────────────────────
	deps := handler.Dependencies{
		VisitorPhoneEncryptor:   phoneEncryptor,
		AuthService:             authSvc,
		CustomerService:         customerSvc,
		VenueService:            venueSvc,
		LevelService:            levelSvc,
		LocationCategoryService: locationCategorySvc,
		LocationService:         locationSvc,
		ProductService:          productSvc,
		MediaService:            mediaSvc,
		StorageService:          storageSvc,
		EventService:            eventSvc,
		UserService:             userSvc,
		NotificationService:     notificationSvc,
		SurveyService:           surveySvc,
		BeaconService:           beaconSvc,
		ConnectionService:       connectionSvc,
		AdvertisementService:    adSvc,
		ArticleService:          articleSvc,
		CouponService:           couponSvc,
		VideoService:            videoSvc,
		TagService:              tagSvc,
		AnalyticsService:        analyticsSvc,
		SnapshotService:         snapshotSvc,
		LevelBundleService:      levelBundleSvc,
		AssetService:            assetSvc,
		LevelTypeService:        levelTypeSvc,
		ThemeService:            themeSvc,
		ProductPlazaService:     productPlazaSvc,
		LanguageService:         languageSvc,
		AppVersionService:       appVersionSvc,
		MemoService:             memoSvc,
		EnricherRegistry:        enricherRegistry,
		UserRepo:                userRepo,
		VenueRepo:               venueRepo,
		LocationRepo:            locationRepo,
		ProductRepo:             productRepo,
		LocationCategoryRepo:    locationCategoryRepo,
		AppUserRepo:             appUserRepo,
		RedisClient:             redisClient,
		DB:                      pool,
		Searcher:                searcher,
		Environment:             cfg.App.Environment,
	}
	router := handler.NewRouter(cfg, logger, deps)

	srv := &http.Server{
		Addr:         fmt.Sprintf(":%d", cfg.Server.Port),
		Handler:      router,
		ReadTimeout:  cfg.Server.ReadTimeout,
		WriteTimeout: cfg.Server.WriteTimeout,
		IdleTimeout:  cfg.Server.IdleTimeout,
	}

	schedulerCtx, cancelScheduler := context.WithCancel(context.Background())
	defer cancelScheduler()

	go func() {
		logger.Info("notification scheduler started")
		ticker := time.NewTicker(time.Minute)
		defer ticker.Stop()
		for {
			select {
			case <-schedulerCtx.Done():
				return
			case t := <-ticker.C:
				sent, err := notificationSvc.SendDueScheduled(schedulerCtx, t)
				if err != nil {
					logger.Errorf("notification scheduler: %v", err)
				} else if sent > 0 {
					logger.Infof("notification scheduler: sent %d notification(s)", sent)
				}
			}
		}
	}()

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
	cancelScheduler()

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	if err := srv.Shutdown(shutdownCtx); err != nil {
		return fmt.Errorf("server shutdown: %w", err)
	}

	logger.Info("server stopped")
	return nil
}
