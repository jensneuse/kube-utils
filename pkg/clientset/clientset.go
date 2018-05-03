package clientset

import (
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"os"
	"path/filepath"
)

// Make a new clientset + rest.Config
// optional kubeconfigPath
func New(kubeConfigPath string) (*rest.Config, *kubernetes.Clientset, error) {

	if kubeConfigPath == "" {
		kubeConfigPath = filepath.Join(
			os.Getenv("HOME"), ".kube", "config",
		)
	}

	config, err := localConfig(kubeConfigPath)
	if err != nil {
		config, err = inClusterConfig()
		if err != nil {
			return nil, nil, err
		}
	}

	client, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, nil, err
	}

	return config, client, nil
}

func localConfig(kubeConfigPath string) (*rest.Config, error) {
	return clientcmd.BuildConfigFromFlags("", kubeConfigPath)
}

func inClusterConfig() (*rest.Config, error) {
	return clientcmd.BuildConfigFromFlags("", "")
}
