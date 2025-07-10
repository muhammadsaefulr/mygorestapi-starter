package database

import (
	"fmt"
	"time"

	"github.com/muhammadsaefulr/NimeStreamAPI/config"

	"github.com/muhammadsaefulr/NimeStreamAPI/internal/domain/model"

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
	// db.Migrator().DropTable(&model.User{}, &model.Token{}, &model.Watchlist{}, &model.TrackEpisodeView{}, model.Comment{})
	db.AutoMigrate(
		&model.User{},
		&model.Token{},
		&model.Watchlist{},
		&model.TrackEpisodeView{},
		&model.Comment{},
		&model.History{},
		&model.RequestMovie{},
		&model.MovieDetails{},
		&model.MovieEpisode{},
		&model.ReportError{},
	)

	// Config connection pooling
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetConnMaxLifetime(60 * time.Minute)

	return db
}
