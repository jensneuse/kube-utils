package pods

import (
	"context"
	"github.com/jensneuse/kube-utils/pkg/podtemplates"
	"github.com/jensneuse/kube-utils/pkg/portforwarder"
	"github.com/jensneuse/kube-utils/pkg/readiness"
	"k8s.io/api/core/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"sync"
)

type Result struct {
	Tunnel portforwarder.Tunnel
	Pod    *v1.Pod
}

func CreateAndForwardBlocking(client *kubernetes.Clientset, config *rest.Config, namespace string, templates ...*podtemplates.PodTemplate) ([]Result, error) {

	var wg sync.WaitGroup
	errChan := make(chan error)
	doneChan := make(chan struct{})
	resultChan := make(chan Result, len(templates))

	for _, template := range templates {

		wg.Add(1)
		go func(template *podtemplates.PodTemplate) {

			var err error
			template.Pod, err = client.CoreV1().Pods(namespace).Create(template.Pod)
			if err != nil {
				errChan <- err
				return
			}

			// wait 1 second for resource to become available
			//time.Sleep(time.Second)

			opts := readiness.Opts{
				Namespace: template.Pod.Namespace,
				PodName:   template.Pod.Name,
			}

			err = readiness.BlockUntilPodReady(client, context.Background(), opts)
			if err != nil {
				errChan <- err
				return
			}

			tunnel, err := portforwarder.New(client, config, template.Pod, template.RemotePort, template.LocalPort)
			if err != nil {
				errChan <- err
				return
			}

			result := Result{
				Tunnel: tunnel,
				Pod:    template.Pod,
			}

			resultChan <- result
			wg.Done()

		}(template)
	}

	go func() {
		wg.Wait()
		close(doneChan)
	}()

	select {
	case err := <-errChan:
		return nil, err
	case <-doneChan:
		results := make([]Result, len(templates))
		for i := 0; i < len(templates); i++ {
			results[i] = <-resultChan
		}

		return results, nil
	}
}
