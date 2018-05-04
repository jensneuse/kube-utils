package readiness

import (
	"context"
	"errors"
	"k8s.io/api/core/v1"
	metaV1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"time"
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

	err = blockUntilPodExists(client, context, opts)
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
			return errors.New("context cancelled/timeout")
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

type Ready bool

func isPodReady(pod *v1.Pod) Ready {
	for _, containerStatus := range pod.Status.ContainerStatuses {
		if !containerStatus.Ready {
			return false
		}
	}

	return true
}

func blockUntilPodExists(client *kubernetes.Clientset, context context.Context, opts Opts) error {

	exists := make(chan error)

	go func() {
		for {
			pod, err := client.CoreV1().Pods(opts.Namespace).Get(opts.PodName, metaV1.GetOptions{})
			if err != nil {
				exists <- err
				break
			}

			if pod != nil && pod.Status.Phase != v1.PodPending {
				close(exists)
				break
			}

			time.Sleep(time.Millisecond * time.Duration(200))
		}
	}()

	select {
	case <-context.Done():
		return errors.New("context cancelled/timeout")
	case err, ok := <-exists:
		if err != nil {
			return err
		}

		if ok {
			return errors.New("unintended")
		}

		return nil
	}
}
