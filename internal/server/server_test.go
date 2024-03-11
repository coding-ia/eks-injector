package server

import (
	"os"
	"testing"
)

func TestServer_BuildConfig(t *testing.T) {
	os.Setenv("CLUSTER_NAME", "test-cluster")
	os.Setenv("CLUSTER_ENVIRONMENT", "sbx")
	os.Setenv("CLUSTER_VERSION", "1.27")

	variables := BuildConfig()
	if len(variables) == 3 {
		if variables["ClusterName"] != "test-cluster" {
			t.Fail()
		}
		if variables["Environment"] != "sbx" {
			t.Fail()
		}
		if variables["Version"] != "1.27" {
			t.Fail()
		}
	} else {
		t.Fail()
	}
}
