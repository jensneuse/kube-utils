package cleanup

import (
	"k8s.io/api/core/v1"
	"k8s.io/client-go/kubernetes"
)

func DeletePods(clientSet *kubernetes.Clientset, pods ...*v1.Pod) error {
	for _, pod := range pods {
		err := clientSet.CoreV1().Pods(pod.Namespace).Delete(pod.Name, nil)
		if err != nil {
			return err
		}
	}

	return nil
}
