package policies

import "testing"

func TestDaemonSetPolicy(t *testing.T) {
	policy, err := FindDaemonSetPolicy("test", "nginx-daemonset", "env")
	_ = policy
	_ = err
}
