package models

import (
	"time"

	"gorm.io/gorm"
)

// Team represents a team/group that has access to specific Kubernetes namespaces
// Each team is mapped to an Okta group for authentication
type Team struct {
	ID           uint           `gorm:"primaryKey" json:"id"`
	TeamName     string         `gorm:"type:varchar(255);not null;unique" json:"teamName"`
	OktaGroupID  string         `gorm:"type:varchar(255);not null;unique" json:"oktaGroupId"`
	Permissions  []Permission   `gorm:"foreignKey:TeamID;constraint:OnDelete:CASCADE" json:"permissions,omitempty"`
	CreatedAt    time.Time      `json:"createdAt"`
	UpdatedAt    time.Time      `json:"updatedAt"`
	DeletedAt    gorm.DeletedAt `gorm:"index" json:"-"`
}

// TableName specifies the table name for the Team model
func (Team) TableName() string {
	return "teams"
}

