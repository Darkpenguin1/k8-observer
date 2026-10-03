package observer

import (
	"fmt"
	appsv1 "k8s.io/api/apps/v1"
	"k8s.io/client-go/tools/cache"
)

func (o *Observer) registerDeploymentInformer() (informerRegistration, error) {
	informer := o.factory.Apps().V1().Deployments().Informer()

	_, err := informer.AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc: func(obj interface{}) {
			deployment, ok := obj.(*appsv1.Deployment)
			if ok {
				o.handleDeploymentAdd(deployment)
			}
		},
		UpdateFunc: func(_, newObj interface{}) {
			deployment, ok := newObj.(*appsv1.Deployment)
			if ok {
				o.handleDeploymentUpdate(deployment)
			}
		},
		DeleteFunc: func(obj interface{}) {
			deployment, ok := deploymentFromDeletedObject(obj)
			if ok {
				o.handleDeploymentDelete(deployment)
			}
		},
	})
	if err != nil {
		return informerRegistration{}, fmt.Errorf("register Deployment event handler: %w", err)
	}

	return informerRegistration{name: "Deployment", informer: informer}, nil
}

func deploymentFromDeletedObject(obj interface{}) (*appsv1.Deployment, bool) {
	if deployment, ok := obj.(*appsv1.Deployment); ok {
		return deployment, true
	}

	tombstone, ok := obj.(cache.DeletedFinalStateUnknown)
	if !ok {
		return nil, false
	}
	deployment, ok := tombstone.Obj.(*appsv1.Deployment)
	return deployment, ok
}

func (o *Observer) handleDeploymentAdd(d *appsv1.Deployment) {
	logResource("Added", "Deployment", describeDeployment(d))
}

func (o *Observer) handleDeploymentUpdate(d *appsv1.Deployment) {
	logResource("Updated", "Deployment", describeDeployment(d))
}

func (o *Observer) handleDeploymentDelete(d *appsv1.Deployment) {
	logResource("Deleted", "Deployment", describeDeployment(d))
}

func describeDeployment(d *appsv1.Deployment) string {
	desired := int32(1) // Kubernetes defaults replicas to 1 if omitted
	if d.Spec.Replicas != nil {
		desired = *d.Spec.Replicas
	}

	return fmt.Sprintf(
		"%s/%s: %d/%d replicas available",
		d.Namespace,
		d.Name,
		d.Status.AvailableReplicas,
		desired,
	)
}
