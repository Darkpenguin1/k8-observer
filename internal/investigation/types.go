package investigation

import (
	"time"

	"k8s.io/apimachinery/pkg/types"
)

type Key struct {
	PodUID        types.UID
	ContainerName string
	Reason        string
}

type Investigation struct {
	Key       Key
	Namespace string
	PodName   string
	FirstSeen time.Time
	LastSeen  time.Time
}

type Detection struct {
	Key        Key
	Namespace  string
	PodName    string
	ObservedAt time.Time
}
