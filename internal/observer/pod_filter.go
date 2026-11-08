package observer

import (
	corev1 "k8s.io/api/core/v1"
)

func crashBackLoopDetected(oldPod, newPod *corev1.Pod) []string {
	var detected []string
	oldCrashLoopingSet := make(map[string]bool)

	// on each pod update we would do a for loop on ALL of its statuses?

	if oldPod != nil {
		for _, status := range oldPod.Status.ContainerStatuses {
			waiting := status.State.Waiting
			inCrashLoopBackOff := waiting != nil && waiting.Reason == "CrashLoopBackOff"

			if inCrashLoopBackOff {
				oldCrashLoopingSet[status.Name] = true
			}

		}
	}

	for _, status := range newPod.Status.ContainerStatuses {
		waiting := status.State.Waiting
		inCrashLoopBackOff := waiting != nil && waiting.Reason == "CrashLoopBackOff"

		if inCrashLoopBackOff && !oldCrashLoopingSet[status.Name] {
			detected = append(detected, status.Name)

		}

	}

	return detected
}
