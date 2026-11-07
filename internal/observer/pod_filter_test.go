package observer

import (
	"reflect"
	"testing"

	corev1 "k8s.io/api/core/v1"
)

func TestCrashBackLoopDetected(t *testing.T) {
	running := func(name string) corev1.ContainerStatus {
		return corev1.ContainerStatus{
			Name: name,
			State: corev1.ContainerState{
				Running: &corev1.ContainerStateRunning{},
			},
		}
	}
	crashLoopBackOff := func(name string) corev1.ContainerStatus {
		return corev1.ContainerStatus{
			Name: name,
			State: corev1.ContainerState{
				Waiting: &corev1.ContainerStateWaiting{Reason: "CrashLoopBackOff"},
			},
		}
	}

	tests := []struct {
		name string
		old  []corev1.ContainerStatus
		new  []corev1.ContainerStatus
		want []string
	}{
		{
			name: "container enters CrashLoopBackOff",
			old:  []corev1.ContainerStatus{running("app")},
			new:  []corev1.ContainerStatus{crashLoopBackOff("app")},
			want: []string{"app"},
		},
		{
			name: "container stays in CrashLoopBackOff",
			old:  []corev1.ContainerStatus{crashLoopBackOff("app")},
			new:  []corev1.ContainerStatus{crashLoopBackOff("app")},
			want: nil,
		},
		{
			name: "matches containers by name when status order changes",
			old: []corev1.ContainerStatus{
				crashLoopBackOff("sidecar"),
				running("app"),
			},
			new: []corev1.ContainerStatus{
				crashLoopBackOff("app"),
				crashLoopBackOff("sidecar"),
			},
			want: []string{"app"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			oldPod := &corev1.Pod{Status: corev1.PodStatus{ContainerStatuses: tt.old}}
			newPod := &corev1.Pod{Status: corev1.PodStatus{ContainerStatuses: tt.new}}

			got := crashBackLoopDetected(oldPod, newPod)
			if !reflect.DeepEqual(got, tt.want) {
				t.Fatalf("crashBackLoopDetected() = %v, want %v", got, tt.want)
			}
		})
	}
}
