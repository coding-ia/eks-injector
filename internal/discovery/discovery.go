package discovery

import (
	"fmt"
	"golang.org/x/net/context"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"os"
)

func DiscoverClusterName() (string, error) {
	return getNodeLabel("awseks.coding-ia.com/cluster-name")
}

func DiscoverEnvironment() (string, error) {
	return getNodeLabel("awseks.coding-ia.com/environment")
}

func DiscoverKubernetesVersion() (string, error) {
	config, err := rest.InClusterConfig()
	if err != nil {
		return "", err
	}

	clientSet, err := kubernetes.NewForConfig(config)
	if err != nil {
		return "", err
	}

	serverVersion, err := clientSet.Discovery().ServerVersion()
	if err != nil {
		return "", err
	}

	versionNumber := fmt.Sprintf("%s.%s", serverVersion.Major, serverVersion.Minor)

	return versionNumber, nil
}

func getNodeLabel(label string) (string, error) {
	config, err := rest.InClusterConfig()
	if err != nil {
		return "", err
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return "", err
	}

	nodeName := os.Getenv("NODE_NAME")

	node, err := clientset.CoreV1().Nodes().Get(context.TODO(), nodeName, metav1.GetOptions{})
	if err != nil {
		return "", err
	}

	clusterName := node.Labels[label]
	return clusterName, nil
}
