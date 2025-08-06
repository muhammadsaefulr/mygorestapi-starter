package database

import (
	"fmt"
	"time"

	"github.com/muhammadsaefulr/NimeStreamAPI/config"

	// "log"
	// "github.com/muhammadsaefulr/NimeStreamAPI/internal/domain/model"
	// "github.com/muhammadsaefulr/NimeStreamAPI/internal/infrastructure/persistence/seed"

	"github.com/muhammadsaefulr/NimeStreamAPI/internal/shared/utils"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func Connect(dbHost, dbName string) *gorm.DB {
	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%d sslmode=disable TimeZone=Asia/Jakarta",
		dbHost, config.DBUser, config.DBPassword, dbName, config.DBPort,
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger:                 logger.Default.LogMode(logger.Info),
		SkipDefaultTransaction: true,
		PrepareStmt:            true,
		TranslateError:         true,
	})
	if err != nil {
		utils.Log.Errorf("Failed to connect to database: %+v", err)
	}

	sqlDB, errDB := db.DB()
	if errDB != nil {
		utils.Log.Errorf("Failed to connect to database: %+v", errDB)
	}

	//auto migrations
	// db.AutoMigrate(
	// 	&model.User{},
	// 	&model.Token{},
	// 	&model.Watchlist{},
	// 	&model.TrackEpisodeView{},
	// 	&model.Comment{},
	// 	&model.CommentLike{},
	// 	&model.History{},
	// 	&model.RequestMovie{},
	// 	&model.MovieDetails{},
	// 	&model.MovieEpisode{},
	// 	&model.ReportError{},
	// 	&model.RolePermissions{},
	// 	&model.UserRole{},
	// 	&model.UserSubscription{},
	// 	&model.UserPoints{},
	// 	&model.UserBadge{},
	// 	&model.UserBadgeInfo{},
	// 	&model.BannerApp{},
	// 	&model.RequestVip{},
	// )

	// implements seed

	// if err := seed.SeedRolesAndPermissions(db); err != nil {
	// 	log.Fatal("gagal seed role permission:", err)
	// }

	// if err := seed.SeedUserRoles(db); err != nil {
	// 	log.Fatal("gagal seed user role:", err)
	// }

	// if err := seed.SeedUsers(db); err != nil {
	// 	log.Fatal("gagal seed user:", err)
	// }

	// if err := seed.SeedSubscriptionPlans(db); err != nil {
	// 	log.Fatal("gagal seed subscription plan:", err)
	// }

	// if err := seed.SeedUserBadges(db); err != nil {
	// 	log.Fatal("gagal seed user badge:", err)
	// }

	// if err := seed.SeedDrakor(db); err != nil {
	// 	log.Fatal("gagal seed drakor:", err)
	// }

	// if err := seed.SeedDrakorEpisodes(db); err != nil {
	// 	log.Fatal("gagal seed drakor episodes:", err)
	// }

	// if err := seed.SeedMovie(db); err != nil {
	// 	log.Fatal("gagal seed movie:", err)
	// }

	// if err := seed.SeedMovieEpisodes(db); err != nil {
	// 	log.Fatal("gagal seed movie episodes:", err)
	// }

	// if err := seed.SeedBannerApp(db); err != nil {
	// 	log.Fatal("gagal seed banner app:", err)
	// }

	// Config connection pooling
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetConnMaxLifetime(60 * time.Minute)

	return db
}
