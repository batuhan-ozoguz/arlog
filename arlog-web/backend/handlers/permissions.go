package handlers

import (
	"encoding/json"
	"log"
	"net/http"

	"arlog/backend/database"
	"arlog/backend/models"
)

// PermissionResponse represents the response for user permissions
type PermissionResponse struct {
	Success     bool                    `json:"success"`
	Permissions []models.PermissionDTO  `json:"permissions"`
	Message     string                  `json:"message,omitempty"`
}

// GetUserPermissions returns the namespaces the authenticated user can access
func GetUserPermissions(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// Get user from context (set by auth middleware)
	// In dev mode, this will be a dummy user; in okta mode, it's from the JWT
	teamName := "Cosmos Team" // Default for dev mode
	
	var team models.Team
	result := database.DB.Preload("Permissions").Where("team_name = ?", teamName).First(&team)
	
	if result.Error != nil {
		log.Printf("Error fetching team: %v", result.Error)
		response := PermissionResponse{
			Success: false,
			Message: "Failed to fetch user permissions",
		}
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(response)
		return
	}

	// Convert permissions to DTOs (without sensitive data)
	permissionDTOs := make([]models.PermissionDTO, len(team.Permissions))
	for i, perm := range team.Permissions {
		permissionDTOs[i] = perm.ToDTO()
	}

	response := PermissionResponse{
		Success:     true,
		Permissions: permissionDTOs,
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

