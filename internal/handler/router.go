package handler

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	"github.com/hhung06/digimap-backend/config"
	_ "github.com/hhung06/digimap-backend/docs"
	"github.com/hhung06/digimap-backend/internal/dto"
	"github.com/hhung06/digimap-backend/internal/enricher"
	"github.com/hhung06/digimap-backend/internal/handler/middleware"
	"github.com/hhung06/digimap-backend/internal/repository"
	"github.com/hhung06/digimap-backend/internal/service"
	applog "github.com/hhung06/digimap-backend/log"
	"github.com/hhung06/digimap-backend/version"
)

// Dependencies holds all service and repository instances needed by handlers.
type Dependencies struct {
	AuthService             service.AuthService
	CustomerService         service.CustomerService
	VenueService            service.VenueService
	LevelService            service.LevelService
	LocationCategoryService service.LocationCategoryService
	LocationService         service.LocationService
	ProductService          service.ProductService
	StorageService          service.StorageService
	EventService            service.EventService
	UserService             service.UserService
	NotificationService     service.NotificationService
	SurveyService           service.SurveyService
	BeaconService           service.BeaconService
	ConnectionService       service.ConnectionService
	AdvertisementService    service.AdvertisementService
	ArticleService          service.ArticleService
	CouponService           service.CouponService
	VideoService            service.VideoService
	TagService              service.TagService
	AnalyticsService        service.AnalyticsService
	SnapshotService         service.SnapshotService
	LevelBundleService      service.LevelBundleService
	AssetService            service.AssetService
	LevelTypeService        service.LevelTypeService
	ThemeService            service.ThemeService
	ProductPlazaService     service.ProductPlazaService
	LanguageService         service.LanguageService
	EnricherRegistry        *enricher.Registry
	UserRepo                repository.UserRepository
	VenueRepo               repository.VenueRepository
	LocationRepo            repository.LocationRepository
	ProductRepo             repository.ProductRepository
	AppUserRepo             repository.AppUserRepository
	RedisClient             *redis.Client
	DB                      *pgxpool.Pool
}

// NewRouter builds and returns the configured Gin engine with all routes registered.
func NewRouter(cfg *config.Config, logger applog.Logger, deps Dependencies) *gin.Engine {
	if cfg.App.IsProduction() {
		gin.SetMode(gin.ReleaseMode)
	}

	r := gin.New()

	r.Use(middleware.Recovery(logger))
	r.Use(middleware.RequestID())
	r.Use(middleware.Logger(logger))
	r.Use(middleware.CORS(cfg.App))
	r.Use(middleware.SecurityHeaders(cfg.App))

	// Infrastructure endpoints
	r.GET("/health", newHealthHandler(deps.DB, deps.RedisClient))
	r.GET("/version", versionHandler)
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	v1 := r.Group("/api/v1")

	// Rate limiters (Redis-backed, shared across replicas)
	publicRL := middleware.RateLimitByIP(deps.RedisClient, "120-M") // 120 req/min per IP for public endpoints
	authRL := middleware.RateLimitByIP(deps.RedisClient, "10-M")    // 10 req/min per IP for auth endpoints
	apiRL := middleware.RateLimitByUser(deps.RedisClient, "600-M")  // 600 req/min per user for API endpoints

	// ── Public auth routes ────────────────────────────────────────────────────
	authH := newAuthHandler(deps.AuthService, authConfig{
		accessExpirySeconds:  int(cfg.JWT.AccessExpiry.Seconds()),
		refreshExpirySeconds: int(cfg.JWT.RefreshExpiry.Seconds()),
	})

	auth := v1.Group("/auth", authRL)
	{
		auth.POST("/login", authH.Login)
		auth.POST("/refresh", authH.RefreshToken)
		auth.POST("/password-reset", authH.RequestPasswordReset)
		auth.POST("/password-reset/confirm", authH.ConfirmPasswordReset)
	}

	// ── JWT-protected routes ──────────────────────────────────────────────────
	jwtAuth := middleware.AuthRequired(deps.AuthService)
	adminJWT := v1.Group("/", jwtAuth, middleware.SystemAdminRequired(), apiRL)
	{
		adminJWT.POST("/auth/logout", authH.Logout)
		adminJWT.PUT("/auth/password-change", authH.ChangePassword)
	}

	// ── System admin — customers ──────────────────────────────────────────────
	custH := newCustomerHandler(deps.CustomerService)
	{
		adminOnly := adminJWT
		// adminOnly.POST("/auth/register", authH.Register)

		adminOnly.GET("/customers", custH.List)
		adminOnly.POST("/customers", custH.Create)
		adminOnly.GET("/customers/:id", custH.Get)
		adminOnly.PUT("/customers/:id", custH.Update)
		adminOnly.DELETE("/customers/:id", custH.Delete)

		// Snapshot sub-resources (system admin only)
		snapshotH := newSnapshotHandler(deps.SnapshotService)
		bundleH := newLevelBundleHandler(deps.LevelBundleService)

		adminOnly.GET("/venues/:id/snapshots", snapshotH.List)
		adminOnly.POST("/venues/:id/snapshots", snapshotH.CreateDraft)
		adminOnly.GET("/venues/:id/snapshots/recent", snapshotH.LatestPublished)
		adminOnly.POST("/venues/:id/snapshots/auto-publish", snapshotH.AutoPublish)
		adminOnly.GET("/venues/:id/snapshots/:snapshotID", snapshotH.Get)
		adminOnly.DELETE("/venues/:id/snapshots/:snapshotID", snapshotH.Delete)
		adminOnly.POST("/venues/:id/snapshots/:snapshotID/publish", snapshotH.Publish)
		adminOnly.POST("/venues/:id/snapshots/:snapshotID/revert", snapshotH.Revert)

		adminOnly.GET("/venues/:id/snapshots/:snapshotID/bundles", bundleH.List)
		adminOnly.POST("/venues/:id/snapshots/:snapshotID/bundles", bundleH.Create)
		adminOnly.DELETE("/venues/:id/snapshots/:snapshotID/bundles/:bundleID", bundleH.Delete)

		assetH := newAssetHandler(deps.AssetService)
		adminOnly.GET("/venues/:id/assets", assetH.List)
		adminOnly.POST("/venues/:id/assets", assetH.Create)
		adminOnly.GET("/venues/:id/assets/:assetID", assetH.Get)
		adminOnly.PUT("/venues/:id/assets/:assetID", assetH.Update)
		adminOnly.DELETE("/venues/:id/assets/:assetID", assetH.Delete)

		ltH := newLevelTypeHandler(deps.LevelTypeService)
		adminOnly.GET("/level-types", ltH.List)
		adminOnly.POST("/level-types", ltH.Create)
		adminOnly.PUT("/level-types/:typeID", ltH.Update)
		adminOnly.DELETE("/level-types/:typeID", ltH.Delete)

		themeH := newThemeHandler(deps.ThemeService)
		adminOnly.GET("/venues/:id/themes", themeH.List)
		adminOnly.POST("/venues/:id/themes", themeH.Create)
		adminOnly.PUT("/venues/:id/themes/:themeID", themeH.Update)
		adminOnly.DELETE("/venues/:id/themes/:themeID", themeH.Delete)

		plazaH := newProductPlazaHandler(deps.ProductPlazaService)
		adminOnly.GET("/venues/:id/product-plazas", plazaH.List)
		adminOnly.POST("/venues/:id/product-plazas", plazaH.Create)
		adminOnly.PUT("/venues/:id/product-plazas/:plazaID", plazaH.Update)
		adminOnly.DELETE("/venues/:id/product-plazas/:plazaID", plazaH.Delete)

		geoLocH := newLocationHandler(deps.LocationCategoryService, deps.LocationService, deps.EnricherRegistry)
		adminOnly.GET("/geo-search", geoLocH.GeoSearch)
	}

	// ── Venue routes (JWT required; RBAC applied per action) ──────────────────
	venueH := newVenueHandler(deps.VenueService)
	venues := adminJWT.Group("/venues")
	{
		venues.GET("", venueH.List)
		venues.POST("", venueH.Create)
		venues.GET("/:id", venueH.Get)
		venues.PUT("/:id", venueH.Update)
		venues.DELETE("/:id", venueH.Delete)
		// venues.POST("/:id/publish", venueH.Publish)
		// venues.POST("/:id/clone", venueH.Clone)
		venues.GET("/:id/key", venueH.GetKey)
		venues.PUT("/:id/key", venueH.RegenerateKey)

		// Level sub-resources
		levelH := newLevelHandler(deps.LevelService)

		venues.GET("/:id/map-groups", levelH.ListMapGroups)
		venues.POST("/:id/map-groups", levelH.CreateMapGroup)
		venues.PUT("/:id/map-groups/:mgID", levelH.UpdateMapGroup)
		venues.DELETE("/:id/map-groups/:mgID", levelH.DeleteMapGroup)

		venues.GET("/:id/levels", levelH.ListLevels)
		venues.POST("/:id/levels", levelH.CreateLevel)
		venues.GET("/:id/levels/:levelID", levelH.GetLevel)
		venues.PUT("/:id/levels/:levelID", levelH.UpdateLevel)
		venues.DELETE("/:id/levels/:levelID", levelH.DeleteLevel)
		venues.PUT("/:id/levels/:levelID/perspective", levelH.UpsertPerspective)
		venues.GET("/:id/levels/:levelID/geo-references", levelH.ListGeoReferences)
		venues.POST("/:id/levels/:levelID/geo-references", levelH.CreateGeoReference)
		venues.DELETE("/:id/levels/:levelID/geo-references/:refID", levelH.DeleteGeoReference)

		// Location sub-resources
		locH := newLocationHandler(deps.LocationCategoryService, deps.LocationService, deps.EnricherRegistry)

		venues.GET("/:id/categories", locH.ListCategories)
		venues.POST("/:id/categories", locH.CreateCategory)
		venues.GET("/:id/categories/:catID", locH.GetCategory)
		venues.PUT("/:id/categories/:catID", locH.UpdateCategory)
		venues.DELETE("/:id/categories/:catID", locH.DeleteCategory)

		venues.GET("/:id/locations", locH.ListLocations)
		venues.POST("/:id/locations", locH.CreateLocation)
		venues.GET("/:id/locations/:locationID", locH.GetLocation)
		venues.PUT("/:id/locations/:locationID", locH.UpdateLocation)
		venues.DELETE("/:id/locations/:locationID", locH.DeleteLocation)
		venues.PUT("/:id/locations/:locationID/set-top", locH.SetTop)
		venues.POST("/:id/locations/:locationID/images", locH.CreateImage)
		venues.DELETE("/:id/locations/:locationID/images/:imageID", locH.DeleteImage)

		// Product sub-resources
		prodH := newProductHandler(deps.ProductService, deps.StorageService, deps.EnricherRegistry)

		venues.GET("/:id/product-categories", prodH.ListCategories)
		venues.POST("/:id/product-categories", prodH.CreateCategory)
		venues.PUT("/:id/product-categories/:catID", prodH.UpdateCategory)
		venues.DELETE("/:id/product-categories/:catID", prodH.DeleteCategory)

		venues.GET("/:id/products", prodH.List)
		venues.POST("/:id/products", prodH.Create)
		venues.GET("/:id/products/:productID", prodH.Get)
		venues.PUT("/:id/products/:productID", prodH.Update)
		venues.DELETE("/:id/products/:productID", prodH.Delete)
		venues.POST("/:id/products/:productID/attachments", prodH.CreateAttachment)
		venues.DELETE("/:id/products/:productID/attachments/:attID", prodH.DeleteAttachment)

		// Event sub-resources
		eventH := newEventHandler(deps.EventService)

		venues.GET("/:id/event-types", eventH.ListEventTypes)
		venues.POST("/:id/event-types", eventH.CreateEventType)
		venues.PUT("/:id/event-types/:typeID", eventH.UpdateEventType)
		venues.DELETE("/:id/event-types/:typeID", eventH.DeleteEventType)

		venues.GET("/:id/events", eventH.ListEvents)
		venues.POST("/:id/events", eventH.CreateEvent)
		venues.GET("/:id/events/:eventID", eventH.GetEvent)
		venues.PUT("/:id/events/:eventID", eventH.UpdateEvent)
		venues.DELETE("/:id/events/:eventID", eventH.DeleteEvent)
		venues.POST("/:id/events/:eventID/images", eventH.CreateEventImage)
		venues.DELETE("/:id/events/:eventID/images/:imageID", eventH.DeleteEventImage)

		// Venue users + invitations (owner only for mutations)
		userH := newUserHandler(deps.UserService)
		venues.GET("/:id/users", userH.ListVenueUsers)
		venues.POST("/:id/users/invite", userH.InviteUser)
		venues.PUT("/:id/users/:userID/role", userH.ChangeRole)
		venues.DELETE("/:id/users/:userID", userH.RemoveFromVenue)
		venues.GET("/:id/invitations", userH.ListInvitations)
		venues.POST("/:id/invitations/:invitationID/cancel", userH.CancelInvitation)

		// Notification sub-resources
		notifH := newNotificationHandler(deps.NotificationService, deps.EnricherRegistry)

		venues.GET("/:id/notifications", notifH.List)
		venues.POST("/:id/notifications", notifH.Create)
		venues.GET("/:id/notifications/:notifID", notifH.Get)
		venues.PUT("/:id/notifications/:notifID", notifH.Update)
		venues.DELETE("/:id/notifications/:notifID", notifH.Delete)
		venues.POST("/:id/notifications/:notifID/send", notifH.Send)

		// Survey sub-resources
		surveyH := newSurveyHandler(deps.SurveyService, deps.EnricherRegistry)

		venues.GET("/:id/surveys", surveyH.List)
		venues.POST("/:id/surveys", surveyH.Create)
		venues.GET("/:id/surveys/:surveyID", surveyH.Get)
		venues.PUT("/:id/surveys/:surveyID", surveyH.Update)
		venues.DELETE("/:id/surveys/:surveyID", surveyH.Delete)
		venues.POST("/:id/surveys/:surveyID/questions", surveyH.CreateQuestion)
		venues.PUT("/:id/surveys/:surveyID/questions/:questionID", surveyH.UpdateQuestion)
		venues.DELETE("/:id/surveys/:surveyID/questions/:questionID", surveyH.DeleteQuestion)
		venues.POST("/:id/surveys/:surveyID/questions/:questionID/options", surveyH.CreateOption)
		venues.PUT("/:id/surveys/:surveyID/questions/:questionID/options/:optionID", surveyH.UpdateOption)
		venues.DELETE("/:id/surveys/:surveyID/questions/:questionID/options/:optionID", surveyH.DeleteOption)
		venues.GET("/:id/surveys/:surveyID/responses", surveyH.ListResponses)
		venues.POST("/:id/surveys/:surveyID/responses", surveyH.SubmitResponse)

		// Beacon sub-resources
		beaconH := newBeaconHandler(deps.BeaconService)

		venues.GET("/:id/beacons", beaconH.List)
		venues.POST("/:id/beacons", beaconH.Create)
		venues.GET("/:id/beacons/:beaconID", beaconH.Get)
		venues.PUT("/:id/beacons/:beaconID", beaconH.Update)
		venues.DELETE("/:id/beacons/:beaconID", beaconH.Delete)

		// Connection sub-resources
		connH := newConnectionHandler(deps.ConnectionService)

		venues.GET("/:id/connections", connH.List)
		venues.POST("/:id/connections", connH.Create)
		venues.GET("/:id/connections/:connectionID", connH.Get)
		venues.PUT("/:id/connections/:connectionID", connH.Update)
		venues.DELETE("/:id/connections/:connectionID", connH.Delete)
		venues.GET("/:id/connections/:connectionID/levels", connH.ListLevels)
		venues.POST("/:id/connections/:connectionID/levels", connH.AddLevel)
		venues.DELETE("/:id/connections/:connectionID/levels/:clID", connH.RemoveLevel)

		// Advertisement sub-resources
		adH := newAdHandler(deps.AdvertisementService, deps.EnricherRegistry)

		venues.GET("/:id/ads", adH.List)
		venues.POST("/:id/ads", adH.Create)
		venues.GET("/:id/ads/:adID", adH.Get)
		venues.PUT("/:id/ads/:adID", adH.Update)
		venues.DELETE("/:id/ads/:adID", adH.Delete)
		venues.POST("/:id/ads/:adID/publish", adH.Publish)

		// Article sub-resources
		articleH := newArticleHandler(deps.ArticleService)

		venues.GET("/:id/articles", articleH.List)
		venues.POST("/:id/articles", articleH.Create)
		venues.GET("/:id/articles/:articleID", articleH.Get)
		venues.PUT("/:id/articles/:articleID", articleH.Update)
		venues.DELETE("/:id/articles/:articleID", articleH.Delete)
		venues.POST("/:id/articles/:articleID/images", articleH.CreateImage)
		venues.DELETE("/:id/articles/:articleID/images/:imageID", articleH.DeleteImage)

		// Coupon sub-resources
		couponH := newCouponHandler(deps.CouponService)

		venues.GET("/:id/coupons", couponH.List)
		venues.POST("/:id/coupons", couponH.Create)
		venues.GET("/:id/coupons/:couponID", couponH.Get)
		venues.PUT("/:id/coupons/:couponID", couponH.Update)
		venues.DELETE("/:id/coupons/:couponID", couponH.Delete)

		// Video sub-resources
		videoH := newVideoHandler(deps.VideoService)

		venues.GET("/:id/videos", videoH.List)
		venues.POST("/:id/videos", videoH.Create)
		venues.GET("/:id/videos/:videoID", videoH.Get)
		venues.PUT("/:id/videos/:videoID", videoH.Update)
		venues.DELETE("/:id/videos/:videoID", videoH.Delete)

		// Language sub-resources
		langH := newLanguageHandler(deps.LanguageService)

		venues.GET("/:id/languages", langH.List)
		venues.POST("/:id/languages", langH.Create)
		venues.PUT("/:id/languages/:langID", langH.Update)
		venues.DELETE("/:id/languages/:langID", langH.Delete)

		// Analytics
		analyticsH := newAnalyticsHandler(deps.AnalyticsService)
		// venues.GET("/:id/analytics/events", analyticsH.ListEventLogs)
		venues.GET("/:id/analytics/searches", analyticsH.ListSearchQueries)
	}

	// ── Analytics: public (no auth required) ─────────────────────────────────
	analyticsH := newAnalyticsHandler(deps.AnalyticsService)
	public := v1.Group("/public", publicRL)
	{
		public.POST("/venues/:id/events", analyticsH.TrackEvent)
		public.POST("/venues/:id/searches", analyticsH.TrackSearch)
	}

	// ── JMA Webhooks (no auth — external callback) ─────────────────────────────
	webhookH := newWebhookHandler(deps.VenueRepo, deps.LocationRepo, deps.ProductRepo)
	webhooks := v1.Group("/webhooks")
	{
		jma := webhooks.Group("/jma")
		jma.POST("/exhibitors", webhookH.JMAExhibitorUpdate)
		jma.POST("/products", webhookH.JMAProductUpdate)
		jma.POST("/push", webhookH.JMAPushNotification)
	}

	// ── Storage (pre-signed uploads) ──────────────────────────────────────────
	prodH := newProductHandler(deps.ProductService, deps.StorageService, deps.EnricherRegistry)
	adminJWT.POST("/storage/presign-upload", prodH.PresignUpload)

	// ── Profile + invitation accept ────────────────────────────────────────────
	userH := newUserHandler(deps.UserService)
	adminJWT.GET("/profile", userH.GetProfile)
	adminJWT.PUT("/profile", userH.UpdateProfile)
	adminJWT.POST("/invitations/accept", userH.AcceptInvitation)

	// ── Global event tags ──────────────────────────────────────────────────────
	eventH := newEventHandler(deps.EventService)
	adminJWT.GET("/event-tags", eventH.ListTags)
	adminJWT.POST("/event-tags", eventH.CreateTag)
	adminJWT.PUT("/event-tags/:tagID", eventH.UpdateTag)
	adminJWT.DELETE("/event-tags/:tagID", eventH.DeleteTag)

	// ── Global tags (polymorphic) ──────────────────────────────────────────────
	tagH := newTagHandler(deps.TagService)
	adminJWT.GET("/tags", tagH.List)
	adminJWT.POST("/tags", tagH.Create)
	adminJWT.GET("/tags/:tagID", tagH.Get)
	adminJWT.PUT("/tags/:tagID", tagH.Update)
	adminJWT.DELETE("/tags/:tagID", tagH.Delete)
	adminJWT.POST("/tags/attach", tagH.AttachTag)
	adminJWT.DELETE("/tags/:tagID/detach", tagH.DetachTag)
	adminJWT.GET("/tags/entity", tagH.ListEntityTags)

	// ── App layer (API Key auth) ───────────────────────────────────────────────
	appH := newAppHandler(
		deps.LocationService,
		deps.EventService,
		deps.ProductService,
		deps.ProductPlazaService,
		deps.ArticleService,
		deps.NotificationService,
		deps.CouponService,
		deps.AdvertisementService,
		deps.SurveyService,
		deps.LocationCategoryService,
		deps.AnalyticsService,
	)
	appKey := r.Group("/app/v1", middleware.APIKeyAuth(venueKeyLookup{deps.VenueRepo}))
	{
		appKey.GET("/locations", appH.ListLocations)
		appKey.GET("/locations/:locationID", appH.GetLocation)
		appKey.GET("/events", appH.ListEvents)
		appKey.GET("/products", appH.ListProducts)
		appKey.GET("/products/:productID", appH.GetProduct)
		appKey.GET("/product-categories", appH.ListProductCategories)
		appKey.GET("/product-plazas", appH.ListProductPlazas)
		appKey.GET("/product-plazas/:plazaID", appH.GetProductPlaza)
		appKey.GET("/articles", appH.ListArticles)
		appKey.GET("/articles/:articleID", appH.GetArticle)
		appKey.GET("/featured-zones", appH.ListFeaturedZones)
		appKey.GET("/featured-zones/:zoneID", appH.GetFeaturedZone)
		appKey.GET("/notifications", appH.ListNotifications)
		appKey.GET("/notifications/:notifID", appH.GetNotification)
		appKey.GET("/coupons", appH.ListCoupons)
		appKey.GET("/coupons/:couponID", appH.GetCoupon)
		appKey.PUT("/coupons/:couponID/redeem", appH.RedeemCoupon)
		appKey.GET("/ads", appH.ListAds)
		appKey.GET("/surveys/:surveyID", appH.GetSurvey)
		appKey.POST("/surveys/:surveyID/submit-response", appH.SubmitSurveyResponse)
		appKey.GET("/top-search", appH.TopSearch)
		appKey.GET("/search-options", appH.SearchOptions)
		appKey.GET("/promotions", appH.Promotions)
	}

	// ── Public API (no auth) ──────────────────────────────────────────────────
	visitorSurveySvc := service.NewVisitorSurveySubmissionService(deps.VenueRepo, deps.AppUserRepo)
	publicH := newPublicHandler(deps.VenueService, deps.SurveyService, deps.ProductPlazaService, deps.AppUserRepo, visitorSurveySvc)
	publicAPI := r.Group("/public/v1")
	{
		publicAPI.GET("/venues/:id/information", publicH.VenueInformation)
		publicAPI.GET("/surveys/:id", publicH.GetSurvey)
		publicAPI.POST("/surveys/:id/submit-response", publicH.SubmitSurveyResponse)
		publicAPI.GET("/venue-info", publicH.VenueInfo)
		publicAPI.POST("/visitor-surveys", publicH.SubmitVisitorSurvey)
	}

	return r
}

// venueKeyLookup adapts repository.VenueRepository to the middleware.VenueByPublicKeyLookup interface.
type venueKeyLookup struct {
	repo repository.VenueRepository
}

func (v venueKeyLookup) FindByPublicKey(ctx context.Context, publicKey string) (uuid.UUID, string, error) {
	venue, err := v.repo.FindByPublicKey(ctx, publicKey)
	if err != nil {
		return uuid.Nil, "", err
	}
	return venue.ID, venue.PrivateKey, nil
}

func versionHandler(c *gin.Context) {
	c.JSON(http.StatusOK, dto.OK(gin.H{
		"version":    version.Version,
		"git_commit": version.GitCommit,
		"build_date": version.BuildDate,
		"go_version": version.GoVersion,
		"os_arch":    version.OsArch,
	}))
}
