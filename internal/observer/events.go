package observer


import (
	"fmt"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/client-go/tools/cache"
)

func (o *Observer) registerEventInformer() (informerRegistration, error) {
	informer := o.factory.Core().V1().Events().Informer()
	
	_, err := informer.AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc: func(obj interface{}) {
			event, ok := obj.(*corev1.Event)
			if ok {
				o.handleEventAdd(event)
			}
		},
		UpdateFunc: func(_, newObj interface{}) {
			event, ok := newObj.(*corev1.Event)
			if ok {
				o.handleEventUpdate(event)
			}
		},
		DeleteFunc: func(obj interface{}) {
			event, ok := eventFromDeletedObject(obj)
			if ok {
				o.handleEventDelete(event)
			}
		},

	})

	return informerRegistration{name: "Event", informer: informer}, err
}

func eventFromDeletedObject(obj interface{}) (*corev1.Event, bool) {
	if event, ok := obj.(*corev1.Event); ok {
		return event, true
	}

	tombstone, ok := obj.(cache.DeletedFinalStateUnknown)
	if !ok {
		return nil, false
	}
	event, ok := tombstone.Obj.(*corev1.Event)
	return event, ok
}

func (o *Observer) handleEventAdd(e *corev1.Event) {
	logResource("Added", "Event", describeEvent(e))
}

func (o *Observer) handleEventUpdate(e *corev1.Event) {
	logResource("Updated", "Event", describeEvent(e))
}

func (o *Observer) handleEventDelete(e *corev1.Event) {
	logResource("Deleted", "Event", describeEvent(e))
}

func describeEvent(e *corev1.Event) string {
	return fmt.Sprintf("%s/%s: %s - %s", e.Namespace, e.Name, e.Reason, e.Message)
}