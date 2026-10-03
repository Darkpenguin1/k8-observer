package observer

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"k8s.io/client-go/informers"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/tools/clientcmd"
)

type Observer struct {
	factory informers.SharedInformerFactory
}

// informerRegistration keeps Run independent of any particular Kubernetes
// resource. New resource observers only need to register an informer and give
// Run a name to use in startup errors.
type informerRegistration struct {
	name     string
	informer cache.SharedIndexInformer
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
		factory: factory,
	}, nil
}

func (o *Observer) Run(ctx context.Context) error {

	registerFuncs := []func() (informerRegistration, error){
		o.registerDeploymentInformer,
		o.registerPodInformer,
		o.registerEventInformer,
	}

	registrations := make([]informerRegistration, 0, len(registerFuncs))
	for _, register := range registerFuncs {
		registration, err := register()
		if err != nil {
			return err
		}
		registrations = append(registrations, registration)
	}

	o.factory.Start(ctx.Done())

	for _, registration := range registrations {
		if !cache.WaitForCacheSync(ctx.Done(), registration.informer.HasSynced) {
			if err := ctx.Err(); err != nil {
				return fmt.Errorf("sync %s informer cache: %w", registration.name, err)
			}
			return fmt.Errorf("sync %s informer cache", registration.name)
		}
	}

	<-ctx.Done()
	return nil
}
