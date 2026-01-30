package models

import (
	"time"

	"gorm.io/gorm"
)

// Permission represents the access control mapping between a team and a Kubernetes namespace
// Each permission grants a team access to a specific namespace in a specific cluster
type Permission struct {
	ID                  uint           `gorm:"primaryKey" json:"id"`
	TeamID              uint           `gorm:"not null;index" json:"teamId"`
	Team                *Team          `gorm:"foreignKey:TeamID" json:"team,omitempty"`
	ClusterName         string         `gorm:"type:varchar(255);not null" json:"clusterName"`
	Namespace           string         `gorm:"type:varchar(255);not null" json:"namespace"`
	ServiceAccountToken string         `gorm:"type:text;not null" json:"-"` // Hidden from JSON for security
	CreatedAt           time.Time      `json:"createdAt"`
	UpdatedAt           time.Time      `json:"updatedAt"`
	DeletedAt           gorm.DeletedAt `gorm:"index" json:"-"`
}

// TableName specifies the table name for the Permission model
func (Permission) TableName() string {
	return "permissions"
}

// PermissionDTO is a data transfer object for permissions without sensitive data
type PermissionDTO struct {
	ID          uint   `json:"id"`
	TeamID      uint   `json:"teamId"`
	ClusterName string `json:"clusterName"`
	Namespace   string `json:"namespace"`
}

// ToDTO converts a Permission to a PermissionDTO (without sensitive fields)
func (p *Permission) ToDTO() PermissionDTO {
	return PermissionDTO{
		ID:          p.ID,
		TeamID:      p.TeamID,
		ClusterName: p.ClusterName,
		Namespace:   p.Namespace,
	}
}

