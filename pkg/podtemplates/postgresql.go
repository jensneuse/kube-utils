package podtemplates

import (
	"k8s.io/api/core/v1"
	metaV1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
)

// Postgresql Database on default port (5432)
func Postgresql() *PodTemplate {
	return &PodTemplate{
		LocalPort:  15432,
		RemotePort: 5432,
		Pod: &v1.Pod{
			ObjectMeta: metaV1.ObjectMeta{
				Name: "integration-test-database",
			},
			Spec: v1.PodSpec{
				Containers: []v1.Container{
					{
						Name:  "postgres",
						Image: "postgres",
						Ports: []v1.ContainerPort{
							{
								ContainerPort: 5432,
							},
						},
						ReadinessProbe: &v1.Probe{
							InitialDelaySeconds: 1,
							TimeoutSeconds:      15,
							Handler: v1.Handler{
								TCPSocket: &v1.TCPSocketAction{
									Port: intstr.FromInt(5432),
								},
							},
						},
					},
				},
			},
		},
	}
}
