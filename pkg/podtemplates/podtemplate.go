package podtemplates

import "k8s.io/api/core/v1"

type PodTemplate struct {
	Pod        *v1.Pod
	LocalPort  int
	RemotePort int
}
