package service

import (
	"errors"
	"runtime"

	"github.com/muhammadsaefulr/mygorestapi-starter/internal/shared/utils"

	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type HealthCheckService interface {
	GormCheck() error
	MemoryHeapCheck() error
	S3Check() error
}

type healthCheckService struct {
	Log *logrus.Logger
	DB  *gorm.DB
	S3  *utils.S3Uploader
}

func NewHealthCheckService(db *gorm.DB, S3 *utils.S3Uploader) HealthCheckService {
	return &healthCheckService{
		Log: utils.Log,
		DB:  db,
		S3:  S3,
	}
}

func (s *healthCheckService) GormCheck() error {
	sqlDB, errDB := s.DB.DB()
	if errDB != nil {
		s.Log.Errorf("failed to access the database connection pool: %v", errDB)
		return errDB
	}

	if err := sqlDB.Ping(); err != nil {
		s.Log.Errorf("failed to ping the database: %v", err)
		return err
	}

	return nil
}

// MemoryHeapCheck checks if heap memory usage exceeds a threshold
func (s *healthCheckService) MemoryHeapCheck() error {
	var memStats runtime.MemStats
	runtime.ReadMemStats(&memStats) // Collect memory statistics

	heapAlloc := memStats.HeapAlloc            // Heap memory currently allocated
	heapThreshold := uint64(300 * 1024 * 1024) // Example threshold: 300 MB

	s.Log.Infof("Heap Memory Allocation: %v bytes", heapAlloc)

	// If the heap allocation exceeds the threshold, return an error
	if heapAlloc > heapThreshold {
		s.Log.Errorf("Heap memory usage exceeds threshold: %v bytes", heapAlloc)
		return errors.New("heap memory usage too high")
	}

	return nil
}

func (s *healthCheckService) S3Check() error {
	err := s.S3.Ping()
	if err != nil {
		s.Log.Errorf("S3 ping failed: %v", err)
		return err
	}
	return nil
}
