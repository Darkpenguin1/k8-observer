package observer

import (
	"fmt"
	"time"

	// corev1 "k8s.io/api/core/v1"
	// metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"context"
	appsv1 "k8s.io/api/apps/v1"
	"k8s.io/client-go/informers"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/tools/clientcmd"
	"log"
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

func (o *Observer) Run(ctx context.Context) error {
	deploymentInformer := o.factory.Apps().V1().Deployments().Informer()

	_, err := deploymentInformer.AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc: func(obj interface{}) {
			deployment, ok := obj.(*appsv1.Deployment)
			if !ok {
				return
			}
			o.handleDeploymentAdd(deployment)
		},
		UpdateFunc: func(oldObj, newObj interface{}) {
			deployment, ok := newObj.(*appsv1.Deployment)
			if !ok {
				return
			}
			o.handleDeploymentUpdate(deployment)
		},
	})
	if err != nil {
		return fmt.Errorf("register Deployment event handler: %w", err)
	}

	o.factory.Start(ctx.Done())

	if !cache.WaitForCacheSync(ctx.Done(), deploymentInformer.HasSynced) {
		return fmt.Errorf("Deployment cache did not sync: %w", ctx.Err())
	}

	<-ctx.Done()
	return nil
}

func (o *Observer) handleDeploymentAdd(d *appsv1.Deployment) {
	// call describeDeployment, then log
	log.Printf("Deployment added: %s", describeDeployment(d))
}

func (o *Observer) handleDeploymentUpdate(d *appsv1.Deployment) {
	// call describeDeployment on the new object, then log
	log.Printf("Deployment added: %s", describeDeployment(d))
}
