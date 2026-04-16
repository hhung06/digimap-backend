package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"

	"github.com/hhung06/digimap-backend/config"
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
	EnricherRegistry        *enricher.Registry
	UserRepo                repository.UserRepository
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
		auth.POST("/register", authH.Register)
		auth.POST("/refresh", authH.RefreshToken)
		auth.POST("/password-reset", authH.RequestPasswordReset)
		auth.POST("/password-reset/confirm", authH.ConfirmPasswordReset)
	}

	// ── JWT-protected routes ──────────────────────────────────────────────────
	jwtAuth := middleware.AuthRequired(deps.AuthService)
	protected := v1.Group("/", jwtAuth)
	{
		protected.POST("/auth/logout", authH.Logout)
		protected.PUT("/auth/password-change", authH.ChangePassword)
	}

	// ── System admin — customers ──────────────────────────────────────────────
	custH := newCustomerHandler(deps.CustomerService)
	adminOnly := v1.Group("/", jwtAuth, middleware.SystemAdminRequired(), apiRL)
	{
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
		adminOnly.GET("/venues/:id/snapshots/:snapshotID", snapshotH.Get)
		adminOnly.DELETE("/venues/:id/snapshots/:snapshotID", snapshotH.Delete)

		adminOnly.POST("/venues/:id/snapshots/:snapshotID/bundles", bundleH.Create)
		adminOnly.GET("/venues/:id/snapshots/:snapshotID/bundles", bundleH.List)
		adminOnly.DELETE("/venues/:id/snapshots/:snapshotID/bundles/:bundleID", bundleH.Delete)
	}

	// ── Venue routes (JWT required; RBAC applied per action) ──────────────────
	venueH := newVenueHandler(deps.VenueService)
	venues := v1.Group("/venues", jwtAuth, apiRL)
	{
		venues.GET("", venueH.List)
		venues.POST("", venueH.Create)
		venues.GET("/:id", venueH.Get)
		venues.PUT("/:id", middleware.VenueAccess(deps.UserRepo, service.RoleEditor), venueH.Update)
		venues.DELETE("/:id", middleware.VenueAccess(deps.UserRepo, service.RoleOwner), venueH.Delete)
		venues.POST("/:id/publish", middleware.VenueAccess(deps.UserRepo, service.RoleOwner), venueH.Publish)
		venues.POST("/:id/clone", middleware.VenueAccess(deps.UserRepo, service.RoleOwner), venueH.Clone)
		venues.GET("/:id/key", middleware.VenueAccess(deps.UserRepo, service.RoleOwner), venueH.GetKey)
		venues.PUT("/:id/key", middleware.VenueAccess(deps.UserRepo, service.RoleOwner), venueH.RegenerateKey)

		// Level sub-resources
		levelH := newLevelHandler(deps.LevelService)
		editorAccess := middleware.VenueAccess(deps.UserRepo, service.RoleEditor)
		viewerAccess := middleware.VenueAccess(deps.UserRepo, service.RoleViewer)

		venues.GET("/:id/map-groups", viewerAccess, levelH.ListMapGroups)
		venues.POST("/:id/map-groups", editorAccess, levelH.CreateMapGroup)
		venues.PUT("/:id/map-groups/:mgID", editorAccess, levelH.UpdateMapGroup)
		venues.DELETE("/:id/map-groups/:mgID", editorAccess, levelH.DeleteMapGroup)

		venues.GET("/:id/levels", viewerAccess, levelH.ListLevels)
		venues.POST("/:id/levels", editorAccess, levelH.CreateLevel)
		venues.GET("/:id/levels/:levelID", viewerAccess, levelH.GetLevel)
		venues.PUT("/:id/levels/:levelID", editorAccess, levelH.UpdateLevel)
		venues.DELETE("/:id/levels/:levelID", editorAccess, levelH.DeleteLevel)
		venues.PUT("/:id/levels/:levelID/perspective", editorAccess, levelH.UpsertPerspective)
		venues.GET("/:id/levels/:levelID/geo-references", viewerAccess, levelH.ListGeoReferences)
		venues.POST("/:id/levels/:levelID/geo-references", editorAccess, levelH.CreateGeoReference)
		venues.DELETE("/:id/levels/:levelID/geo-references/:refID", editorAccess, levelH.DeleteGeoReference)

		// Location sub-resources
		locH := newLocationHandler(deps.LocationCategoryService, deps.LocationService, deps.EnricherRegistry)

		venues.GET("/:id/categories", viewerAccess, locH.ListCategories)
		venues.POST("/:id/categories", editorAccess, locH.CreateCategory)
		venues.GET("/:id/categories/:catID", viewerAccess, locH.GetCategory)
		venues.PUT("/:id/categories/:catID", editorAccess, locH.UpdateCategory)
		venues.DELETE("/:id/categories/:catID", editorAccess, locH.DeleteCategory)

		venues.GET("/:id/locations", viewerAccess, locH.ListLocations)
		venues.POST("/:id/locations", editorAccess, locH.CreateLocation)
		venues.GET("/:id/locations/:locationID", viewerAccess, locH.GetLocation)
		venues.PUT("/:id/locations/:locationID", editorAccess, locH.UpdateLocation)
		venues.DELETE("/:id/locations/:locationID", editorAccess, locH.DeleteLocation)
		venues.POST("/:id/locations/:locationID/duplicate", editorAccess, locH.DuplicateLocation)
		venues.PUT("/:id/locations/:locationID/set-top", editorAccess, locH.SetTop)
		venues.POST("/:id/locations/:locationID/images", editorAccess, locH.CreateImage)
		venues.DELETE("/:id/locations/:locationID/images/:imageID", editorAccess, locH.DeleteImage)

		// Product sub-resources
		prodH := newProductHandler(deps.ProductService, deps.StorageService, deps.EnricherRegistry)

		venues.GET("/:id/product-categories", viewerAccess, prodH.ListCategories)
		venues.POST("/:id/product-categories", editorAccess, prodH.CreateCategory)
		venues.PUT("/:id/product-categories/:catID", editorAccess, prodH.UpdateCategory)
		venues.DELETE("/:id/product-categories/:catID", editorAccess, prodH.DeleteCategory)

		venues.GET("/:id/products", viewerAccess, prodH.List)
		venues.POST("/:id/products", editorAccess, prodH.Create)
		venues.GET("/:id/products/:productID", viewerAccess, prodH.Get)
		venues.PUT("/:id/products/:productID", editorAccess, prodH.Update)
		venues.DELETE("/:id/products/:productID", editorAccess, prodH.Delete)
		venues.POST("/:id/products/:productID/attachments", editorAccess, prodH.CreateAttachment)
		venues.DELETE("/:id/products/:productID/attachments/:attID", editorAccess, prodH.DeleteAttachment)

		// Event sub-resources
		eventH := newEventHandler(deps.EventService)

		venues.GET("/:id/event-types", viewerAccess, eventH.ListEventTypes)
		venues.POST("/:id/event-types", editorAccess, eventH.CreateEventType)
		venues.PUT("/:id/event-types/:typeID", editorAccess, eventH.UpdateEventType)
		venues.DELETE("/:id/event-types/:typeID", editorAccess, eventH.DeleteEventType)

		venues.GET("/:id/events", viewerAccess, eventH.ListEvents)
		venues.POST("/:id/events", editorAccess, eventH.CreateEvent)
		venues.GET("/:id/events/:eventID", viewerAccess, eventH.GetEvent)
		venues.PUT("/:id/events/:eventID", editorAccess, eventH.UpdateEvent)
		venues.DELETE("/:id/events/:eventID", editorAccess, eventH.DeleteEvent)
		venues.POST("/:id/events/:eventID/images", editorAccess, eventH.CreateEventImage)
		venues.DELETE("/:id/events/:eventID/images/:imageID", editorAccess, eventH.DeleteEventImage)

		// Venue users + invitations (owner only for mutations)
		userH := newUserHandler(deps.UserService)
		ownerAccess := middleware.VenueAccess(deps.UserRepo, service.RoleOwner)

		venues.GET("/:id/users", viewerAccess, userH.ListVenueUsers)
		venues.POST("/:id/users/invite", ownerAccess, userH.InviteUser)
		venues.PUT("/:id/users/:userID/role", ownerAccess, userH.ChangeRole)
		venues.DELETE("/:id/users/:userID", ownerAccess, userH.RemoveFromVenue)
		venues.GET("/:id/invitations", ownerAccess, userH.ListInvitations)
		venues.POST("/:id/invitations/:invitationID/cancel", ownerAccess, userH.CancelInvitation)

		// Notification sub-resources
		notifH := newNotificationHandler(deps.NotificationService, deps.EnricherRegistry)

		venues.GET("/:id/notifications", viewerAccess, notifH.List)
		venues.POST("/:id/notifications", editorAccess, notifH.Create)
		venues.GET("/:id/notifications/:notifID", viewerAccess, notifH.Get)
		venues.PUT("/:id/notifications/:notifID", editorAccess, notifH.Update)
		venues.DELETE("/:id/notifications/:notifID", editorAccess, notifH.Delete)
		venues.POST("/:id/notifications/:notifID/send", editorAccess, notifH.Send)

		// Survey sub-resources
		surveyH := newSurveyHandler(deps.SurveyService)

		venues.GET("/:id/surveys", viewerAccess, surveyH.List)
		venues.POST("/:id/surveys", editorAccess, surveyH.Create)
		venues.GET("/:id/surveys/:surveyID", viewerAccess, surveyH.Get)
		venues.PUT("/:id/surveys/:surveyID", editorAccess, surveyH.Update)
		venues.DELETE("/:id/surveys/:surveyID", editorAccess, surveyH.Delete)
		venues.POST("/:id/surveys/:surveyID/questions", editorAccess, surveyH.CreateQuestion)
		venues.PUT("/:id/surveys/:surveyID/questions/:questionID", editorAccess, surveyH.UpdateQuestion)
		venues.DELETE("/:id/surveys/:surveyID/questions/:questionID", editorAccess, surveyH.DeleteQuestion)
		venues.POST("/:id/surveys/:surveyID/questions/:questionID/options", editorAccess, surveyH.CreateOption)
		venues.PUT("/:id/surveys/:surveyID/questions/:questionID/options/:optionID", editorAccess, surveyH.UpdateOption)
		venues.DELETE("/:id/surveys/:surveyID/questions/:questionID/options/:optionID", editorAccess, surveyH.DeleteOption)
		venues.GET("/:id/surveys/:surveyID/responses", viewerAccess, surveyH.ListResponses)
		venues.POST("/:id/surveys/:surveyID/responses", surveyH.SubmitResponse)

		// Beacon sub-resources
		beaconH := newBeaconHandler(deps.BeaconService)

		venues.GET("/:id/beacons", viewerAccess, beaconH.List)
		venues.POST("/:id/beacons", editorAccess, beaconH.Create)
		venues.GET("/:id/beacons/:beaconID", viewerAccess, beaconH.Get)
		venues.PUT("/:id/beacons/:beaconID", editorAccess, beaconH.Update)
		venues.DELETE("/:id/beacons/:beaconID", editorAccess, beaconH.Delete)

		// Connection sub-resources
		connH := newConnectionHandler(deps.ConnectionService)

		venues.GET("/:id/connections", viewerAccess, connH.List)
		venues.POST("/:id/connections", editorAccess, connH.Create)
		venues.GET("/:id/connections/:connectionID", viewerAccess, connH.Get)
		venues.PUT("/:id/connections/:connectionID", editorAccess, connH.Update)
		venues.DELETE("/:id/connections/:connectionID", editorAccess, connH.Delete)
		venues.GET("/:id/connections/:connectionID/levels", viewerAccess, connH.ListLevels)
		venues.POST("/:id/connections/:connectionID/levels", editorAccess, connH.AddLevel)
		venues.DELETE("/:id/connections/:connectionID/levels/:clID", editorAccess, connH.RemoveLevel)

		// Advertisement sub-resources
		adH := newAdHandler(deps.AdvertisementService)

		venues.GET("/:id/ads", viewerAccess, adH.List)
		venues.POST("/:id/ads", editorAccess, adH.Create)
		venues.GET("/:id/ads/:adID", viewerAccess, adH.Get)
		venues.PUT("/:id/ads/:adID", editorAccess, adH.Update)
		venues.DELETE("/:id/ads/:adID", editorAccess, adH.Delete)
		venues.POST("/:id/ads/:adID/publish", editorAccess, adH.Publish)

		// Article sub-resources
		articleH := newArticleHandler(deps.ArticleService)

		venues.GET("/:id/articles", viewerAccess, articleH.List)
		venues.POST("/:id/articles", editorAccess, articleH.Create)
		venues.GET("/:id/articles/:articleID", viewerAccess, articleH.Get)
		venues.PUT("/:id/articles/:articleID", editorAccess, articleH.Update)
		venues.DELETE("/:id/articles/:articleID", editorAccess, articleH.Delete)
		venues.POST("/:id/articles/:articleID/images", editorAccess, articleH.CreateImage)
		venues.DELETE("/:id/articles/:articleID/images/:imageID", editorAccess, articleH.DeleteImage)

		// Coupon sub-resources
		couponH := newCouponHandler(deps.CouponService)

		venues.GET("/:id/coupons", viewerAccess, couponH.List)
		venues.POST("/:id/coupons", editorAccess, couponH.Create)
		venues.GET("/:id/coupons/:couponID", viewerAccess, couponH.Get)
		venues.PUT("/:id/coupons/:couponID", editorAccess, couponH.Update)
		venues.DELETE("/:id/coupons/:couponID", editorAccess, couponH.Delete)

		// Video sub-resources
		videoH := newVideoHandler(deps.VideoService)

		venues.GET("/:id/videos", viewerAccess, videoH.List)
		venues.POST("/:id/videos", editorAccess, videoH.Create)
		venues.GET("/:id/videos/:videoID", viewerAccess, videoH.Get)
		venues.PUT("/:id/videos/:videoID", editorAccess, videoH.Update)
		venues.DELETE("/:id/videos/:videoID", editorAccess, videoH.Delete)

		// Analytics (protected — viewer access)
		analyticsH := newAnalyticsHandler(deps.AnalyticsService)
		venues.GET("/:id/analytics/events", viewerAccess, analyticsH.ListEventLogs)
		venues.GET("/:id/analytics/searches", viewerAccess, analyticsH.ListSearchQueries)
	}

	// ── Analytics: public (no auth required) ─────────────────────────────────
	analyticsH := newAnalyticsHandler(deps.AnalyticsService)
	public := v1.Group("/public", publicRL)
	{
		public.POST("/venues/:id/events", analyticsH.TrackEvent)
		public.POST("/venues/:id/searches", analyticsH.TrackSearch)
	}

	// ── Storage (pre-signed uploads) ──────────────────────────────────────────
	prodH := newProductHandler(deps.ProductService, deps.StorageService, deps.EnricherRegistry)
	protected.POST("/storage/presign-upload", prodH.PresignUpload)

	// ── Profile + invitation accept ────────────────────────────────────────────
	userH := newUserHandler(deps.UserService)
	protected.GET("/profile", userH.GetProfile)
	protected.PUT("/profile", userH.UpdateProfile)
	protected.POST("/invitations/accept", userH.AcceptInvitation)

	// ── Global event tags ──────────────────────────────────────────────────────
	eventH := newEventHandler(deps.EventService)
	v1.GET("/event-tags", jwtAuth, eventH.ListTags)
	v1.POST("/event-tags", jwtAuth, eventH.CreateTag)
	v1.PUT("/event-tags/:tagID", jwtAuth, eventH.UpdateTag)
	v1.DELETE("/event-tags/:tagID", jwtAuth, eventH.DeleteTag)

	// ── Global tags (polymorphic) ──────────────────────────────────────────────
	tagH := newTagHandler(deps.TagService)
	v1.GET("/tags", jwtAuth, tagH.List)
	v1.POST("/tags", jwtAuth, tagH.Create)
	v1.GET("/tags/:tagID", jwtAuth, tagH.Get)
	v1.PUT("/tags/:tagID", jwtAuth, tagH.Update)
	v1.DELETE("/tags/:tagID", jwtAuth, tagH.Delete)
	v1.POST("/tags/attach", jwtAuth, tagH.AttachTag)
	v1.DELETE("/tags/:tagID/detach", jwtAuth, tagH.DetachTag)
	v1.GET("/tags/entity", jwtAuth, tagH.ListEntityTags)

	return r
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
