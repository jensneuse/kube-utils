package podtemplates

import (
	"k8s.io/api/core/v1"
	metaV1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
)

// Minio S3 server listening on 9000
func Minio(podName, accessKey, secretKey string) *PodTemplate {
	return &PodTemplate{
		LocalPort:  9000,
		RemotePort: 9000,
		Pod: &v1.Pod{
			ObjectMeta: metaV1.ObjectMeta{
				Name: podName,
			},
			Spec: v1.PodSpec{
				Containers: []v1.Container{
					{
						Name:  "minio",
						Image: "minio/minio",
						Args:  []string{"server", "/data"},
						Env: []v1.EnvVar{
							{
								Name:  "MINIO_ACCESS_KEY",
								Value: accessKey,
							},
							{
								Name:  "MINIO_SECRET_KEY",
								Value: secretKey,
							},
						},
						Ports: []v1.ContainerPort{
							{
								ContainerPort: 9000,
							},
						},
						ReadinessProbe: &v1.Probe{
							InitialDelaySeconds: 1,
							TimeoutSeconds:      15,
							Handler: v1.Handler{
								TCPSocket: &v1.TCPSocketAction{
									Port: intstr.FromInt(9000),
								},
							},
						},
					},
				},
			},
		},
	}
}
