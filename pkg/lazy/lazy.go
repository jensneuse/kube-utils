package lazy

import (
	"github.com/jensneuse/kube-utils/pkg/cleanup"
	"github.com/jensneuse/kube-utils/pkg/clientset"
	"github.com/jensneuse/kube-utils/pkg/pods"
	"github.com/jensneuse/kube-utils/pkg/podtemplates"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"log"
)

type Lazy struct {
	config    *rest.Config
	clientSet *kubernetes.Clientset
	results   []pods.Result
}

func New() *Lazy {
	config, client, err := clientset.New("")
	if err != nil {
		panic(err)
	}

	return &Lazy{config, client, []pods.Result{}}
}

func first(results []pods.Result, err error) pods.Result {

	if err != nil {
		log.Panic(err)
	}

	if len(results) != 1 {
		log.Panicf("invalid results length, expected 1 but got %d", len(results))
	}

	return results[0]
}

func (l *Lazy) Cleanup() error {
	for _, result := range l.results {
		result.Tunnel.Close()
		err := cleanup.DeletePods(l.clientSet, result.Pod)
		if err != nil {
			return err
		}
	}

	return nil
}

func (l *Lazy) CreateAndForwardPods(namespace string, templates ...*podtemplates.PodTemplate) []pods.Result {
	results, err := pods.CreateAndForwardBlocking(l.clientSet, l.config, namespace, templates...)
	if err != nil {
		log.Panic(err)
	}
	l.results = append(l.results, results...)
	return results
}

func (l *Lazy) CreateAndForwardMinio(namespace, podName, accessKey, secretKey string) pods.Result {
	results, err := pods.CreateAndForwardBlocking(l.clientSet, l.config, namespace, podtemplates.Minio(podName, accessKey, secretKey))
	l.results = append(l.results, results...)
	return first(results, err)
}

func (l *Lazy) CreateAndForwardPostgres(namespace, podName string) pods.Result {
	results, err := pods.CreateAndForwardBlocking(l.clientSet, l.config, namespace, podtemplates.Postgresql(podName))
	l.results = append(l.results, results...)
	return first(results, err)
}
