package handlers

import (
	"encoding/json"
	"log"
	"net/http"

	"arlog/backend/services"
)

// PodInfo represents basic pod information
type PodInfo struct {
	Name      string `json:"name"`
	Status    string `json:"status"`
	Namespace string `json:"namespace"`
	Ready     string `json:"ready"`
	Restarts  int32  `json:"restarts"`
	Age       string `json:"age"`
}

// PodsResponse represents the response for listing pods
type PodsResponse struct {
	Success   bool      `json:"success"`
	Pods      []PodInfo `json:"pods"`
	Namespace string    `json:"namespace"`
	Message   string    `json:"message,omitempty"`
}

// GetPods lists all pods in the specified namespace
// Query parameters:
//   - namespace: The Kubernetes namespace to query (required)
func GetPods(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// Get namespace from query parameters
	namespace := r.URL.Query().Get("namespace")
	if namespace == "" {
		response := PodsResponse{
			Success: false,
			Message: "namespace query parameter is required",
		}
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(response)
		return
	}

	// TODO: In Task 4.5, validate that the user has permission to access this namespace

	// Get Kubernetes service
	k8sService, err := services.NewKubernetesService()
	if err != nil {
		log.Printf("Error creating Kubernetes service: %v", err)
		response := PodsResponse{
			Success: false,
			Message: "Failed to connect to Kubernetes cluster",
		}
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(response)
		return
	}

	// List pods in the namespace
	pods, err := k8sService.ListPods(namespace)
	if err != nil {
		log.Printf("Error listing pods in namespace %s: %v", namespace, err)
		response := PodsResponse{
			Success:   false,
			Namespace: namespace,
			Message:   "Failed to list pods: " + err.Error(),
		}
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(response)
		return
	}

	// Convert to PodInfo structs
	podInfos := make([]PodInfo, len(pods))
	for i, pod := range pods {
		podInfos[i] = PodInfo{
			Name:      pod.Name,
			Status:    pod.Status,
			Namespace: pod.Namespace,
			Ready:     pod.Ready,
			Restarts:  pod.Restarts,
			Age:       pod.Age,
		}
	}

	response := PodsResponse{
		Success:   true,
		Pods:      podInfos,
		Namespace: namespace,
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

