package module

import (
	"github.com/gofiber/fiber/v2"
	"github.com/muhammadsaefulr/NimeStreamAPI/config"
	"github.com/muhammadsaefulr/NimeStreamAPI/internal/delivery/http/router"
	"github.com/muhammadsaefulr/NimeStreamAPI/internal/delivery/middleware"
	userRepo "github.com/muhammadsaefulr/NimeStreamAPI/internal/repository/user"
	watchListRepo "github.com/muhammadsaefulr/NimeStreamAPI/internal/repository/watchlist"
	userService "github.com/muhammadsaefulr/NimeStreamAPI/internal/service/user_service"
	watchlistService "github.com/muhammadsaefulr/NimeStreamAPI/internal/service/watchlist_service"

	trackRepo "github.com/muhammadsaefulr/NimeStreamAPI/internal/repository/track_episode_view"

	commentRepo "github.com/muhammadsaefulr/NimeStreamAPI/internal/repository/comment"
	commetService "github.com/muhammadsaefulr/NimeStreamAPI/internal/service/comment_service"

	historyRepo "github.com/muhammadsaefulr/NimeStreamAPI/internal/repository/history"
	HistoryService "github.com/muhammadsaefulr/NimeStreamAPI/internal/service/history_service"

	requestMovieRepo "github.com/muhammadsaefulr/NimeStreamAPI/internal/repository/request_movie"
	requestMovieService "github.com/muhammadsaefulr/NimeStreamAPI/internal/service/request_movie_service"

	movieDetailRepo "github.com/muhammadsaefulr/NimeStreamAPI/internal/repository/movie_details"
	movieDetailService "github.com/muhammadsaefulr/NimeStreamAPI/internal/service/movie_details_service"

	movieUploaderRepo "github.com/muhammadsaefulr/NimeStreamAPI/internal/repository/movie_episode"
	movieUploaderSvc "github.com/muhammadsaefulr/NimeStreamAPI/internal/service/movie_episode_service"

	authService "github.com/muhammadsaefulr/NimeStreamAPI/internal/service/auth_service"
	odService "github.com/muhammadsaefulr/NimeStreamAPI/internal/service/otakudesu_scrape"
	systemService "github.com/muhammadsaefulr/NimeStreamAPI/internal/service/system_service"

	AnilistService "github.com/muhammadsaefulr/NimeStreamAPI/internal/service/anilist_service"
	mdlService "github.com/muhammadsaefulr/NimeStreamAPI/internal/service/mdl_service"
	tmdbService "github.com/muhammadsaefulr/NimeStreamAPI/internal/service/tmdb_service"

	discoverySvc "github.com/muhammadsaefulr/NimeStreamAPI/internal/service/discovery_service"

	// User Auth And Role

	userRoleRepo "github.com/muhammadsaefulr/NimeStreamAPI/internal/repository/user_role"
	userRoleService "github.com/muhammadsaefulr/NimeStreamAPI/internal/service/user_role_service"

	rolePermissionRepo "github.com/muhammadsaefulr/NimeStreamAPI/internal/repository/role_permissions"
	rolePermissionService "github.com/muhammadsaefulr/NimeStreamAPI/internal/service/role_permissions_service"

	// reporting

	reportErrRepo "github.com/muhammadsaefulr/NimeStreamAPI/internal/repository/report_error"
	reportErrService "github.com/muhammadsaefulr/NimeStreamAPI/internal/service/report_error_service"

	// other services

	userPointsRepo "github.com/muhammadsaefulr/NimeStreamAPI/internal/repository/user_points"
	userPointsService "github.com/muhammadsaefulr/NimeStreamAPI/internal/service/user_points_service"

	bannerAppRepo "github.com/muhammadsaefulr/NimeStreamAPI/internal/repository/banner_app"
	bannerAppService "github.com/muhammadsaefulr/NimeStreamAPI/internal/service/banner_app_service"

	userBadgeRepo "github.com/muhammadsaefulr/NimeStreamAPI/internal/repository/user_badge"
	userBadgeService "github.com/muhammadsaefulr/NimeStreamAPI/internal/service/user_badge_service"

	"github.com/muhammadsaefulr/NimeStreamAPI/internal/shared/utils"
	"github.com/muhammadsaefulr/NimeStreamAPI/internal/shared/validation"

	"gorm.io/gorm"
)

func InitModule(app *fiber.App, db *gorm.DB) {
	validate := validation.Validator()

	uploader, err := utils.NewS3Uploader(
		"http://minio:9000", // Endpoint (MinIO or AWS)
		"admin", "4dm1n3rs", // Access key
		"https://dev.msaepul.my.id/minio", // Endpoint (MinIO or AWS)
	)

	if err != nil {
		utils.Log.Errorf("Failed to create S3 uploader: %v", err)
		return
	}
	// Init services
	userRepo := userRepo.NewUserRepositryImpl(db)
	userSvc := userService.NewUserService(userRepo, validate)

	// User Auth And Role

	userRoleRepo := userRoleRepo.NewUserRoleRepositoryImpl(db)
	userRoleSvc := userRoleService.NewUserRoleService(userRoleRepo, validate)

	rolePermissionRepo := rolePermissionRepo.NewRolePermissionsRepositoryImpl(db)
	rolePermissionSvc := rolePermissionService.NewRolePermissionsService(rolePermissionRepo, validate)

	trackEpsRepo := trackRepo.NewTrackEpisodeViewRepository(db)
	animeSvc := odService.NewAnimeService(trackEpsRepo)

	watchlistRepo := watchListRepo.NewWatchlistRepositoryImpl(db)
	watchListSvc := watchlistService.NewWatchlistService(watchlistRepo, validate, animeSvc)

	commentRepo := commentRepo.NewCommentRepository(db)
	commentSvc := commetService.NewCommentService(commentRepo)

	historyRepo := historyRepo.NewHistoryRepositoryImpl(db)
	historySvc := HistoryService.NewHistoryService(historyRepo, validate, animeSvc)

	requestMovieRepo := requestMovieRepo.NewRequestMovieRepositoryImpl(db)
	requestMovieSvc := requestMovieService.NewRequestMovieService(requestMovieRepo, validate)

	tokenSvc := systemService.NewTokenService(db, validate, userSvc)
	authSvc := authService.NewAuthService(db, validate, userSvc, tokenSvc)
	emailSvc := systemService.NewEmailService()
	healthSvc := systemService.NewHealthCheckService(db)

	// Other Sources Services

	anilistSvc := AnilistService.NewAnilistService(validate)
	tmdbSvc := tmdbService.NewTMDbService(validate)
	mdlSvc := mdlService.NewMdlService(validate)

	// Native Upload Data Manual

	movieDetailRepo := movieDetailRepo.NewMovieDetailsRepositoryImpl(db)
	movieDetailSvc := movieDetailService.NewMovieDetailsService(movieDetailRepo, validate, anilistSvc, tmdbSvc, mdlSvc, animeSvc)

	movieUploaderRepo := movieUploaderRepo.NewMovieEpisodeRepositoryImpl(db)
	movieUploaderSvc := movieUploaderSvc.NewMovieEpisodeService(movieUploaderRepo, validate, uploader, movieDetailSvc)

	// Dynamic Orchestrator Source

	discoverySvc := discoverySvc.NewDiscoveryService(validate, anilistSvc, tmdbSvc, mdlSvc, animeSvc, movieDetailSvc)

	// reporting

	reportErrRepo := reportErrRepo.NewReportErrorRepositoryImpl(db)
	reportErrSvc := reportErrService.NewReportErrorService(reportErrRepo, validate)

	// other services

	userPointsRepo := userPointsRepo.NewUserPointsRepositoryImpl(db)
	userPointsSvc := userPointsService.NewUserPointsService(userPointsRepo, validate)

	bannerAppRepo := bannerAppRepo.NewBannerAppRepositoryImpl(db)
	bannerAppSvc := bannerAppService.NewBannerAppService(bannerAppRepo, validate)

	userBadgeRepo := userBadgeRepo.NewUserBadgeRepositoryImpl(db)
	userBadgeSvc := userBadgeService.NewUserBadgeService(userBadgeRepo, validate)

	middleware.InitAuthMiddleware(userSvc)

	v1 := app.Group("/api/v1")

	router.AuthRoutes(v1, authSvc, userSvc, tokenSvc, emailSvc)
	router.UserRoutes(v1, userSvc, tokenSvc)
	router.OdRoutes(v1, animeSvc)
	router.HealthCheckRoutes(v1, healthSvc)
	router.DocsRoutes(v1)
	router.WatchlistRoutes(v1, watchListSvc)
	router.CommentsRoutes(v1, commentSvc)
	router.HistoryRoutes(v1, historySvc)
	router.RequestMovieRoutes(v1, requestMovieSvc)
	router.ReportErrorRoutes(v1, reportErrSvc)
	router.UserRoleRoutes(v1, userRoleSvc)
	router.RolePermissionsRoutes(v1, rolePermissionSvc)
	router.UserBadgeRoutes(v1, userBadgeSvc)

	// Other Sources

	router.TmdbRoutes(v1, tmdbSvc)
	router.AnilistRoutes(v1, anilistSvc)
	router.MdlRoutes(v1, mdlSvc)

	// Dynamic Orchestrator Source

	router.DiscoveryRoutes(v1, discoverySvc)

	// Native Upload Data Manual

	router.MovieDetailsRoutes(v1, movieDetailSvc)
	router.MovieEpisodeRoutes(v1, movieUploaderSvc)

	// other services

	router.UserPointsRoutes(v1, userPointsSvc)
	router.BannerAppRoutes(v1, bannerAppSvc)

	if !config.IsProd {
		v1.Get("/docs", func(c *fiber.Ctx) error {
			return c.SendString("API Docs here")
		})
	}
}
