package policies

import (
	"encoding/json"
	"os"
)

type Policies struct {
	DaemonSets  []*Policy `json:"DaemonSets"`
	Deployments []*Policy `json:"Deployments"`
	ConfigMaps  []*Policy `json:"ConfigMaps"`
}

type Policy struct {
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

func LoadPolicies(path string) (Policies, error) {
	var policies Policies
	data, err := os.ReadFile(path)
	if err != nil {
		return policies, err
	}

	err = json.Unmarshal(data, &policies)
	if err != nil {
		return policies, err
	}

	return policies, nil
}

func FindDeploymentPolicy(policies []*Policy, namespace string, name string, keyType string) (*Policy, error) {
	for _, v := range policies {
		if v.Namespace == namespace && v.Name == name && v.Type == keyType {
			return v, nil
		}
	}

	return nil, nil
}

func FindDaemonSetPolicy(policies []*Policy, namespace string, name string, keyType string) (*Policy, error) {
	for _, v := range policies {
		if v.Namespace == namespace && v.Name == name && v.Type == keyType {
			return v, nil
		}
	}

	return nil, nil
}

func FindConfigMapPolicy(policies []*Policy, namespace string, name string, keyType string) (*Policy, error) {
	for _, v := range policies {
		if v.Namespace == namespace && v.Name == name {
			return v, nil
		}
	}

	return nil, nil
}
