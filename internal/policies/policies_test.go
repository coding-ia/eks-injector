package policies

import (
	"os"
	"testing"
)

func TestPolicies(t *testing.T) {
	rawJSON := `
{
	"DaemonSets": [
		{
			"namespace": "aqua",
			"name": "aqua-enforcer-ds",
			"key": "AQUA_LOGICAL_NAME",
			"value": "{{ .ClusterName }}",
			"keyType": "env"
		}
	],
	"Deployments": [
		{
			"namespace": "sentinelone",
			"name": "sentinelone-helper",
			"key": "CLUSTER_NAME",
			"value": "{{ .ClusterName }}",
			"keyType": "env"
		}
	],
	"ConfigMaps": [
		{
			"namespace": "aqua",
			"name": "aqua-csp-kube-enforcer",
			"key": "AQUA_LOGICAL_NAME",
			"value": "{{ .ClusterName }}"
		}
	]
}
`
	tempFile, err := os.CreateTemp("", "example")
	if err != nil {
		t.Fatalf("Error creating temporary file: %v", err)
	}
	defer os.Remove(tempFile.Name()) // Clean up the temporary file
	data := []byte(rawJSON)
	if _, err := tempFile.Write(data); err != nil {
		tempFile.Close()
		t.Fatalf("Error writing JSON to temporary file: %v", err)
	}
	filePath := tempFile.Name()
	policies, _ := LoadPolicies(filePath)
	_ = policies
}
