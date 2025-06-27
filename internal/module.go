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

	authService "github.com/muhammadsaefulr/NimeStreamAPI/internal/service/auth_service"
	odService "github.com/muhammadsaefulr/NimeStreamAPI/internal/service/otakudesu_scrape"
	systemService "github.com/muhammadsaefulr/NimeStreamAPI/internal/service/system_service"
	"github.com/muhammadsaefulr/NimeStreamAPI/internal/shared/validation"

	"gorm.io/gorm"
)

func InitModule(app *fiber.App, db *gorm.DB) {
	validate := validation.Validator()

	// Init services
	userRepo := userRepo.NewUserRepositryImpl(db)
	userSvc := userService.NewUserService(userRepo, validate)

	trackEpsRepo := trackRepo.NewTrackEpisodeViewRepository(db)
	animeSvc := odService.NewAnimeService(trackEpsRepo)

	watchlistRepo := watchListRepo.NewWatchlistRepositoryImpl(db)
	watchListSvc := watchlistService.NewWatchlistService(watchlistRepo, validate, animeSvc)

	commentRepo := commentRepo.NewCommentRepository(db)
	commentSvc := commetService.NewCommentService(commentRepo)

	historyRepo := historyRepo.NewHistoryRepositoryImpl(db)
	historySvc := HistoryService.NewHistoryService(historyRepo, validate)

	requestMovieRepo := requestMovieRepo.NewRequestMovieRepositoryImpl(db)
	requestMovieSvc := requestMovieService.NewRequestMovieService(requestMovieRepo, validate)

	tokenSvc := systemService.NewTokenService(db, validate, userSvc)
	authSvc := authService.NewAuthService(db, validate, userSvc, tokenSvc)
	emailSvc := systemService.NewEmailService()
	healthSvc := systemService.NewHealthCheckService(db)

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

	if !config.IsProd {
		v1.Get("/docs", func(c *fiber.Ctx) error {
			return c.SendString("API Docs here")
		})
	}
}
