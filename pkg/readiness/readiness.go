package readiness

import (
	"context"
	"errors"
	"fmt"
	"k8s.io/api/core/v1"
	metaV1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

type Opts struct {
	// required
	Namespace string
	// required
	PodName string
}

func (o Opts) validate() error {
	if o.PodName == "" {
		return errors.New("variable PodName must not be empty")
	}

	if o.Namespace == "" {
		return errors.New("variable Namespace must not be empty")
	}

	return nil
}

func BlockUntilPodReady(client *kubernetes.Clientset, context context.Context, opts Opts) error {

	err := opts.validate()
	if err != nil {
		return err
	}

	watcher, err := client.
		CoreV1().
		Pods(opts.Namespace).
		Watch(metaV1.SingleObject(metaV1.ObjectMeta{Name: opts.PodName}))

	if err != nil {
		return err
	}

	for {
		select {
		case <-context.Done():
			return errors.New("timeout reached")
		case event, ok := <-watcher.ResultChan():

			if !ok {
				return nil
			}

			pod, ok := event.Object.(*v1.Pod)
			if ok {
				if isPodReady(pod) {
					watcher.Stop()
					return nil
				}
			}
		}
	}
}

func SignalPodReady(client *kubernetes.Clientset, context context.Context, opts Opts) chan error {

	ch := make(chan error)

	go func() {
		err := BlockUntilPodReady(client, context, opts)
		if err != nil {
			ch <- err
		}

		close(ch)
	}()

	return ch
}

type Ready bool

func isPodReady(pod *v1.Pod) Ready {
	for _, containerStatus := range pod.Status.ContainerStatuses {
		if !containerStatus.Ready {
			fmt.Printf("ContainerState: %s\n", containerStatus.String())
			return false
		}
	}

	return true
}
