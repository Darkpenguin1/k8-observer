package observer

import (
	"fmt"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/client-go/tools/cache"
)

func (o *Observer) registerPodInformer() (informerRegistration, error) {
	informer := o.factory.Core().V1().Pods().Informer()

	// whole lotta type assertions
	_, err := informer.AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc: func(obj interface{}) {
			pod, ok := obj.(*corev1.Pod)

			if !ok || pod == nil {
				return
			}

			logCrashLoopDetections(nil, pod)

			o.handlePodAdd(pod)

		},
		UpdateFunc: func(oldObj, newObj interface{}) {

			oldPod, oldOk := oldObj.(*corev1.Pod)
			pod, ok := newObj.(*corev1.Pod)

			if !oldOk || oldPod == nil || !ok || pod == nil {
				return
			}

			logCrashLoopDetections(oldPod, pod)

			o.handlePodUpdate(pod)

		},
		DeleteFunc: func(obj interface{}) {
			pod, ok := podFromDeletedObject(obj)
			if ok {
				o.handlePodDelete(pod)
			}
		},
	})

	return informerRegistration{name: "Pod", informer: informer}, err
}

func podFromDeletedObject(obj interface{}) (*corev1.Pod, bool) {
	if pod, ok := obj.(*corev1.Pod); ok {
		return pod, pod != nil
	}

	tombstone, ok := obj.(cache.DeletedFinalStateUnknown)
	if !ok {
		return nil, false
	}
	pod, ok := tombstone.Obj.(*corev1.Pod)
	return pod, ok && pod != nil
}

func (o *Observer) handlePodAdd(p *corev1.Pod) {
	logResource("Added", "Pod", describePod(p))
}

func (o *Observer) handlePodUpdate(p *corev1.Pod) {
	logResource("Updated", "Pod", describePod(p))
}

func (o *Observer) handlePodDelete(p *corev1.Pod) {
	logResource("Deleted", "Pod", describePod(p))
}

func describePod(p *corev1.Pod) string {
	return fmt.Sprintf("%s/%s (Phase: %s)", p.Namespace, p.Name, p.Status.Phase)
}

func logCrashLoopDetections(oldPod, pod *corev1.Pod) {
	for _, containerName := range crashBackLoopDetected(oldPod, pod) {
		description := fmt.Sprintf("%s/%s container=%s", pod.Namespace, pod.Name, containerName)
		logResource("Detected", "CrashLoopBackOff", description)
	}
}
