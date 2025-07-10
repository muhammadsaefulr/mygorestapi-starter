package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ReportError struct {
	ReportId     uuid.UUID `gorm:"type:uuid;primaryKey"`
	ReportedBy   uuid.UUID `gorm:"type:uuid;not null" json:"reported_by,omitempty"`
	HandledBy    uuid.UUID `gorm:"type:uuid"`
	ProblemDesc  string    `gorm:"type:text;not null"`
	EpisodeId    string    `gorm:"type:varchar(100);not null"`
	StatusReport string    `gorm:"type:varchar(50);default:'pending'"`
	CreatedAt    time.Time `gorm:"autoCreateTime"`
	UpdatedAt    time.Time `gorm:"autoUpdateTime"`
}

func (r *ReportError) BeforeCreate(tx *gorm.DB) (err error) {
	if r.ReportId == uuid.Nil {
		r.ReportId = uuid.New()
	}
	return
}
