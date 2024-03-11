package string_parser

import "testing"

func TestParser(t *testing.T) {
	values := map[string]string{
		"cluster_name": "test_cluster",
		"version":      "1.27",
		"environment":  "sbx",
	}
	parsedString, err := ParseString("{{ .cluster_name }}/{{ .version }}/{{ .environment }}", values)
	if err != nil {
		t.Fail()
	}
	_ = parsedString

}
