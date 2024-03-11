package string_parser

import "testing"

func TestParser(t *testing.T) {
	values := map[string]string{
		"ClusterName": "test_cluster",
		"Version":     "1.27",
		"Environment": "sbx",
	}
	parsedString, err := ParseString("{{ .ClusterName }}/{{ .Version }}/{{ .Environment }}", values)
	if err != nil {
		t.Fail()
	}
	if parsedString != "test_cluster/1.27/sbx" {
		t.Fail()
	}
}
