package handlers

import (
	"log"
	"net/http"

	"arlog/backend/services"

	"github.com/gorilla/websocket"
)

// WebSocket upgrader configuration
var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	// Allow all origins for development (should be restricted in production)
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

// StreamLogs handles WebSocket connections for streaming pod logs
// Query parameters:
//   - namespace: The Kubernetes namespace (required)
//   - podName: The name of the pod (required)
//   - container: The container name (optional, uses first container if not specified)
//   - follow: Whether to follow logs (default: true)
//   - tailLines: Number of lines to show from the end (default: 100)
func StreamLogs(w http.ResponseWriter, r *http.Request) {
	// Get query parameters
	namespace := r.URL.Query().Get("namespace")
	podName := r.URL.Query().Get("podName")
	container := r.URL.Query().Get("container")

	// Validate required parameters
	if namespace == "" || podName == "" {
		http.Error(w, "namespace and podName query parameters are required", http.StatusBadRequest)
		return
	}

	// TODO: In Task 4.5, validate that the user has permission to access this namespace

	// Upgrade HTTP connection to WebSocket
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("Error upgrading connection to WebSocket: %v", err)
		return
	}
	defer conn.Close()

	log.Printf("WebSocket connection established for pod: %s/%s", namespace, podName)

	// Get Kubernetes service
	k8sService, err := services.NewKubernetesService()
	if err != nil {
		log.Printf("Error creating Kubernetes service: %v", err)
		conn.WriteMessage(websocket.TextMessage, []byte("Error: Failed to connect to Kubernetes cluster"))
		return
	}

	// Create a custom writer that sends data to the WebSocket
	wsWriter := &WebSocketWriter{
		conn: conn,
	}

	// Stream logs to the WebSocket
	err = k8sService.StreamLogs(namespace, podName, container, wsWriter)
	if err != nil {
		log.Printf("Error streaming logs for pod %s/%s: %v", namespace, podName, err)
		conn.WriteMessage(websocket.TextMessage, []byte("Error: "+err.Error()))
		return
	}

	log.Printf("WebSocket connection closed for pod: %s/%s", namespace, podName)
}

// WebSocketWriter is an io.Writer that writes to a WebSocket connection
type WebSocketWriter struct {
	conn *websocket.Conn
}

// Write implements the io.Writer interface for WebSocket
func (w *WebSocketWriter) Write(p []byte) (n int, err error) {
	err = w.conn.WriteMessage(websocket.TextMessage, p)
	if err != nil {
		return 0, err
	}
	return len(p), nil
}

