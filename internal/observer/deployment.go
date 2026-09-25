package observer

import (
	"fmt"

	appsv1 "k8s.io/api/apps/v1"
)

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
