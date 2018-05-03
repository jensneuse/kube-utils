package portforwarder

import (
	"context"
	"github.com/jensneuse/kube-utils/pkg/clientset"
	"github.com/jensneuse/kube-utils/pkg/readiness"
	"k8s.io/api/core/v1"
	metaV1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
	"math/rand"
	"strconv"
	"testing"
)

func TestPortForwarder(t *testing.T) {

	namespace := "default"

	config, client, err := clientset.New("")
	if err != nil {
		t.Fatal(err)
	}

	pod := &v1.Pod{
		ObjectMeta: metaV1.ObjectMeta{
			Name: "integration-test-database-" + strconv.FormatInt(rand.Int63(), 10),
			Labels: map[string]string{
				"integration": "test",
			},
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
	}

	createdPod, err := client.CoreV1().Pods(namespace).Create(pod)
	if err != nil {
		t.Fatal(err)
	}

	err = readiness.BlockUntilPodReady(client, context.Background(), readiness.Opts{
		Namespace: namespace,
		PodName:   createdPod.Name,
	})

	if err != nil {
		t.Fatal(err)
	}

	tunnel, err := New(client, config, namespace, createdPod.Name, 5432, 15432)
	if err != nil {
		t.Fatal(err)
	}

	tunnel.Close()

	err = client.CoreV1().Pods(namespace).Delete(pod.Name, nil)
	if err != nil {
		t.Fatal(err)
	}
}
