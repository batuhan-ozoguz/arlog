package services

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"os"
	"time"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

// KubernetesService provides methods to interact with the Kubernetes API
type KubernetesService struct {
	clientset *kubernetes.Clientset
	config    *rest.Config
}

// PodInfo represents basic information about a Kubernetes pod
type PodInfo struct {
	Name      string
	Status    string
	Namespace string
	Ready     string
	Restarts  int32
	Age       string
}

// NewKubernetesService creates a new Kubernetes service instance
// For local development, it uses kubeconfig or kubectl proxy
// For production, it will use service account tokens
func NewKubernetesService() (*KubernetesService, error) {
	var config *rest.Config
	var err error

	// Try to use in-cluster configuration first (when running inside K8s)
	config, err = rest.InClusterConfig()
	if err != nil {
		// If in-cluster config fails, try kubeconfig file
		kubeconfig := os.Getenv("KUBECONFIG")
		if kubeconfig == "" {
			homeDir, _ := os.UserHomeDir()
			kubeconfig = homeDir + "/.kube/config"
		}

		config, err = clientcmd.BuildConfigFromFlags("", kubeconfig)
		if err != nil {
			// If kubeconfig fails, try kubectl proxy
			proxyURL := os.Getenv("KUBE_PROXY_URL")
			if proxyURL == "" {
				proxyURL = "http://localhost:8001"
			}

			config = &rest.Config{
				Host: proxyURL,
			}
		}
	}

	// Create the clientset
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, fmt.Errorf("failed to create Kubernetes client: %w", err)
	}

	return &KubernetesService{
		clientset: clientset,
		config:    config,
	}, nil
}

// NewKubernetesServiceWithToken creates a Kubernetes service with a specific service account token
// This is used when accessing namespaces with specific permissions
func NewKubernetesServiceWithToken(token string, apiServerURL string) (*KubernetesService, error) {
	config := &rest.Config{
		Host:        apiServerURL,
		BearerToken: token,
		TLSClientConfig: rest.TLSClientConfig{
			Insecure: true, // For MVP; in production, use proper TLS verification
		},
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, fmt.Errorf("failed to create Kubernetes client with token: %w", err)
	}

	return &KubernetesService{
		clientset: clientset,
		config:    config,
	}, nil
}

// ListPods returns a list of pods in the specified namespace
func (k *KubernetesService) ListPods(namespace string) ([]PodInfo, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	pods, err := k.clientset.CoreV1().Pods(namespace).List(ctx, metav1.ListOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed to list pods: %w", err)
	}

	podInfos := make([]PodInfo, 0, len(pods.Items))
	for _, pod := range pods.Items {
		podInfo := PodInfo{
			Name:      pod.Name,
			Namespace: pod.Namespace,
			Status:    string(pod.Status.Phase),
			Ready:     getPodReadyStatus(&pod),
			Restarts:  getPodRestartCount(&pod),
			Age:       calculateAge(pod.CreationTimestamp.Time),
		}
		podInfos = append(podInfos, podInfo)
	}

	return podInfos, nil
}

// StreamLogs streams logs from a pod to the provided writer
// This function follows the logs in real-time
func (k *KubernetesService) StreamLogs(namespace, podName, container string, writer io.Writer) error {
	ctx := context.Background()

	// Get pod to check if container name is needed
	pod, err := k.clientset.CoreV1().Pods(namespace).Get(ctx, podName, metav1.GetOptions{})
	if err != nil {
		return fmt.Errorf("failed to get pod: %w", err)
	}

	// If no container specified and pod has multiple containers, use the first one
	if container == "" {
		if len(pod.Spec.Containers) > 0 {
			container = pod.Spec.Containers[0].Name
		}
	}

	// Configure log options
	tailLines := int64(100)
	logOptions := &corev1.PodLogOptions{
		Container:  container,
		Follow:     true,
		Timestamps: true,
		TailLines:  &tailLines,
	}

	// Get log stream
	req := k.clientset.CoreV1().Pods(namespace).GetLogs(podName, logOptions)
	stream, err := req.Stream(ctx)
	if err != nil {
		return fmt.Errorf("failed to get log stream: %w", err)
	}
	defer stream.Close()

	// Stream logs to writer
	reader := bufio.NewReader(stream)
	for {
		line, err := reader.ReadBytes('\n')
		if err != nil {
			if err == io.EOF {
				break
			}
			return fmt.Errorf("error reading log stream: %w", err)
		}

		// Write log line to the provided writer
		_, writeErr := writer.Write(line)
		if writeErr != nil {
			return fmt.Errorf("error writing log line: %w", writeErr)
		}
	}

	return nil
}

// GetPodLogs retrieves logs from a pod (non-streaming, for historical logs)
func (k *KubernetesService) GetPodLogs(namespace, podName, container string, tailLines int64) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	logOptions := &corev1.PodLogOptions{
		Container:  container,
		Timestamps: true,
		TailLines:  &tailLines,
	}

	req := k.clientset.CoreV1().Pods(namespace).GetLogs(podName, logOptions)
	logs, err := req.DoRaw(ctx)
	if err != nil {
		return "", fmt.Errorf("failed to get pod logs: %w", err)
	}

	return string(logs), nil
}

// getPodReadyStatus returns the ready status of a pod (e.g., "2/3" meaning 2 out of 3 containers are ready)
func getPodReadyStatus(pod *corev1.Pod) string {
	totalContainers := len(pod.Spec.Containers)
	readyContainers := 0

	for _, containerStatus := range pod.Status.ContainerStatuses {
		if containerStatus.Ready {
			readyContainers++
		}
	}

	return fmt.Sprintf("%d/%d", readyContainers, totalContainers)
}

// getPodRestartCount returns the total number of restarts for all containers in a pod
func getPodRestartCount(pod *corev1.Pod) int32 {
	var restarts int32 = 0
	for _, containerStatus := range pod.Status.ContainerStatuses {
		restarts += containerStatus.RestartCount
	}
	return restarts
}

// calculateAge calculates the age of a resource from its creation timestamp
func calculateAge(creationTime time.Time) string {
	duration := time.Since(creationTime)

	// Format duration in a human-readable way
	days := int(duration.Hours() / 24)
	hours := int(duration.Hours()) % 24
	minutes := int(duration.Minutes()) % 60

	if days > 0 {
		return fmt.Sprintf("%dd%dh", days, hours)
	} else if hours > 0 {
		return fmt.Sprintf("%dh%dm", hours, minutes)
	} else if minutes > 0 {
		return fmt.Sprintf("%dm", minutes)
	} else {
		return fmt.Sprintf("%ds", int(duration.Seconds()))
	}
}

