package observer

import (
	"fmt"
	"time"

	// corev1 "k8s.io/api/core/v1"
	// metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/informers"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"os"
	"path/filepath"
)

type Observer struct {
	// informer and anything it needs to run
	factory   informers.SharedInformerFactory
	clientset *kubernetes.Clientset
}

func New(inCluster bool) (*Observer, error) {
	var (
		config *rest.Config
		err    error
	)

	if inCluster {
		config, err = rest.InClusterConfig()
	} else {
		kubeconfig := os.Getenv("KUBECONFIG")
		if kubeconfig == "" {
			home, homeErr := os.UserHomeDir()
			if homeErr != nil {
				return nil, fmt.Errorf("find home directory: %w", homeErr)
			}
			kubeconfig = filepath.Join(home, ".kube", "config")
		}
		config, err = clientcmd.BuildConfigFromFlags("", kubeconfig)
	}
	if err != nil {
		return nil, fmt.Errorf("load Kubernetes config: %w", err)
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, fmt.Errorf("create Kubernetes client: %w", err)
	}

	factory := informers.NewSharedInformerFactory(clientset, 10*time.Minute)
	return &Observer{
		factory:   factory,
		clientset: clientset,
	}, nil
}

func (o *Observer) Run() {
	// register handlers and start the informer
}

func (o *Observer) handleDeploymentAdd() {
	// call describeDeployment, then log
}

func (o *Observer) handleDeploymentUpdate() {
	// call describeDeployment on the new object, then log
}
