package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ReportError struct {
	ReportId     uuid.UUID  `gorm:"type:uuid;primaryKey" json:"report_id"`
	ReportedBy   uuid.UUID  `gorm:"type:uuid;not null" json:"reported_by,omitempty"`
	HandledBy    *uuid.UUID `gorm:"type:uuid" json:"handle_by,omitempty"`
	ProblemDesc  string     `gorm:"type:text;not null" json:"problem_desc"`
	EpisodeId    string     `gorm:"type:varchar(100);not null" json:"episode_id"`
	TypeMovie    string     `gorm:"type:varchar(100);not null" json:"type_movie"`
	StatusReport string     `gorm:"type:varchar(50);default:'pending'" json:"status_report"`
	CreatedAt    time.Time  `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt    time.Time  `gorm:"autoUpdateTime" json:"updated_at"`

	User *User `gorm:"foreignKey:ReportedBy" json:"user_reported_by"`
}

func (r *ReportError) BeforeCreate(tx *gorm.DB) (err error) {
	if r.ReportId == uuid.Nil {
		r.ReportId = uuid.New()
	}
	return
}
