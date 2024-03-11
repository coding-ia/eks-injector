package policies

type Policies struct {
	DaemonSets  []*Policy `json:"DaemonSets"`
	Deployments []*Policy `json:"Deployments"`
	ConfigMaps  []*Policy `json:"ConfigMaps"`
}

type Policy struct {
	Namespace   string `json:"namespace"`
	Name        string `json:"name"`
	Key         string `json:"key"`
	Type        string `json:"keyType"`
	Value       string `json:"value,omitempty"`
	SkipAdd     bool   `json:"skipAdd"`
	SkipReplace bool   `json:"skipReplace"`
}

type ConfigPolicy struct {
	Namespace   string       `json:"namespace"`
	Name        string       `json:"name"`
	Key         string       `json:"key"`
	Type        string       `json:"keyType"`
	Value       string       `json:"value,omitempty"`
	SkipAdd     bool         `json:"skipAdd"`
	SkipReplace bool         `json:"skipReplace"`
	SSM         SSMParameter `json:"ssm,omitempty"`
}

type SSMParameter struct {
	Region  string `json:"region"`
	Name    string `json:"name"`
	Decrypt bool   `json:"decrypt"`
}

func FindDeploymentPolicy(namespace string, name string, keyType string) (*Policy, error) {
	var policy *Policy

	if namespace == "sentinelone" && name == "sentinelone-helper" && keyType == "env" {
		policy = &Policy{
			Namespace: namespace,
			Name:      name,
			Key:       "CLUSTER_NAME",
			Value:     "{{ .ClusterName }}",
			Type:      keyType,
		}
	}
	if namespace == "default" && name == "nginx-deployment" && keyType == "env" {
		policy = &Policy{
			Namespace: namespace,
			Name:      name,
			Key:       "CLUSTER_NAME",
			Value:     "{{ .ClusterName }}",
			Type:      keyType,
		}
	}

	return policy, nil
}

func FindDaemonSetPolicy(namespace string, name string, keyType string) (*Policy, error) {
	var policy *Policy

	if namespace == "aqua" && name == "aqua-enforcer-ds" && keyType == "env" {
		policy = &Policy{
			Namespace: namespace,
			Name:      name,
			Key:       "AQUA_LOGICAL_NAME",
			Value:     "{{ .ClusterName }}",
			Type:      keyType,
		}
	}
	if namespace == "test" && name == "nginx-daemonset" && keyType == "env" {
		policy = &Policy{
			Namespace: namespace,
			Name:      name,
			Key:       "AQUA_LOGICAL_NAME",
			Value:     "{{ .ClusterName }}",
			Type:      keyType,
		}
	}

	return policy, nil
}

func FindConfigMapPolicy(namespace string, name string, keyType string) (*Policy, error) {
	var policy *Policy

	if namespace == "aqua" && name == "aqua-csp-kube-enforcer" && keyType == "" {
		policy = &Policy{
			Namespace: namespace,
			Name:      name,
			Key:       "AQUA_LOGICAL_NAME",
			Value:     "{{ .ClusterName }}",
			Type:      "",
		}
	}
	if namespace == "default" && name == "example-configmap" && keyType == "" {
		policy = &Policy{
			Namespace: namespace,
			Name:      name,
			Key:       "logicalName",
			Value:     "{{ .ClusterName }}",
			Type:      "",
		}
	}

	return policy, nil
}
