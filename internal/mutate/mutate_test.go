package mutate

import (
	"crypto/md5"
	"eks-inject/internal/policies"
	"encoding/hex"
	"encoding/json"
	"fmt"
	admissionv1 "k8s.io/api/admission/v1"
	"testing"
)

func TestMutatesDeploymentRequest(t *testing.T) {
	rawJSON := `{
		"kind": "AdmissionReview",
		"apiVersion": "admission.k8s.io/v1",
		"request": {
			"uid": "803630c2-34c6-45af-8768-3cf416e9b1b1",
			"kind": {
				"group": "apps",
				"version": "v1",
				"kind": "Deployment"
			},
			"resource": {
				"group": "apps",
				"version": "v1",
				"resource": "deployments"
			},
			"requestKind": {
				"group": "apps",
				"version": "v1",
				"kind": "Deployment"
			},
			"requestResource": {
				"group": "apps",
				"version": "v1",
				"resource": "deployments"
			},
			"name": "nginx-deployment",
			"namespace": "default",
			"operation": "CREATE",
			"userInfo": {
				"username": "system:admin",
				"groups": [
					"system:masters",
					"system:authenticated"
				]
			},
			"object": {
				"kind": "Deployment",
				"apiVersion": "apps/v1",
				"metadata": {
					"name": "nginx-deployment",
					"namespace": "default",
					"creationTimestamp": null,
					"annotations": {
						"kubectl.kubernetes.io/last-applied-configuration": "{\"apiVersion\":\"apps/v1\",\"kind\":\"Deployment\",\"metadata\":{\"annotations\":{},\"name\":\"nginx-deployment\",\"namespace\":\"default\"},\"spec\":{\"replicas\":1,\"selector\":{\"matchLabels\":{\"app\":\"nginx\"}},\"template\":{\"metadata\":{\"labels\":{\"app\":\"nginx\"}},\"spec\":{\"containers\":[{\"image\":\"nginx:latest\",\"name\":\"nginx\",\"ports\":[{\"containerPort\":80}]}]}}}}\n"
					},
					"managedFields": [
						{
							"manager": "kubectl-client-side-apply",
							"operation": "Update",
							"apiVersion": "apps/v1",
							"time": "2024-03-07T00:02:47Z",
							"fieldsType": "FieldsV1",
							"fieldsV1": {
								"f:metadata": {
									"f:annotations": {
										".": {},
										"f:kubectl.kubernetes.io/last-applied-configuration": {}
									}
								},
								"f:spec": {
									"f:progressDeadlineSeconds": {},
									"f:replicas": {},
									"f:revisionHistoryLimit": {},
									"f:selector": {},
									"f:strategy": {
										"f:rollingUpdate": {
											".": {},
											"f:maxSurge": {},
											"f:maxUnavailable": {}
										},
										"f:type": {}
									},
									"f:template": {
										"f:metadata": {
											"f:labels": {
												".": {},
												"f:app": {}
											}
										},
										"f:spec": {
											"f:containers": {
												"k:{\"name\":\"nginx\"}": {
													".": {},
													"f:image": {},
													"f:imagePullPolicy": {},
													"f:name": {},
													"f:ports": {
														".": {},
														"k:{\"containerPort\":80,\"protocol\":\"TCP\"}": {
															".": {},
															"f:containerPort": {},
															"f:protocol": {}
														}
													},
													"f:resources": {},
													"f:terminationMessagePath": {},
													"f:terminationMessagePolicy": {}
												}
											},
											"f:dnsPolicy": {},
											"f:restartPolicy": {},
											"f:schedulerName": {},
											"f:securityContext": {},
											"f:terminationGracePeriodSeconds": {}
										}
									}
								}
							}
						}
					]
				},
				"spec": {
					"replicas": 1,
					"selector": {
						"matchLabels": {
							"app": "nginx"
						}
					},
					"template": {
						"metadata": {
							"creationTimestamp": null,
							"labels": {
								"app": "nginx"
							}
						},
						"spec": {
							"containers": [
								{
									"name": "nginx",
									"image": "nginx:latest",
									"ports": [
										{
											"containerPort": 80,
											"protocol": "TCP"
										}
									],
									"resources": {},
									"terminationMessagePath": "/dev/termination-log",
									"terminationMessagePolicy": "File",
									"imagePullPolicy": "Always"
								}
							],
							"restartPolicy": "Always",
							"terminationGracePeriodSeconds": 30,
							"dnsPolicy": "ClusterFirst",
							"securityContext": {},
							"schedulerName": "default-scheduler"
						}
					},
					"strategy": {
						"type": "RollingUpdate",
						"rollingUpdate": {
							"maxUnavailable": "25%",
							"maxSurge": "25%"
						}
					},
					"revisionHistoryLimit": 10,
					"progressDeadlineSeconds": 600
				},
				"status": {}
			},
			"oldObject": null,
			"dryRun": false,
			"options": {
				"kind": "CreateOptions",
				"apiVersion": "meta.k8s.io/v1",
				"fieldManager": "kubectl-client-side-apply",
				"fieldValidation": "Strict"
			}
		}
	}`

	values := map[string]string{
		"ClusterName": "test-cluster",
		"Version":     "1.27",
		"Environment": "sbx",
	}
	tp := testPolicies()
	data, err := ProcessAdmissionReview([]byte(rawJSON), values, tp)
	if err == nil {
		ar, err := getAdmissionReview(data)
		if err != nil {
			t.Fail()
		}
		if ar.Request.UID != ar.Response.UID {
			t.Fail()
		}
		hash := getMD5Hash(ar.Response.Patch)
		if hash != "f54e0ee9713c2142a45cc2a8d7e194aa" {
			t.Fail()
		}
	}
}

func TestMutatesConfigMapRequest(t *testing.T) {
	rawJSON := `{
		"kind": "AdmissionReview",
		"apiVersion": "admission.k8s.io/v1",
		"request": {
			"uid": "42b2eea6-2458-421a-ad7d-cd49f90abec5",
			"kind": {
				"group": "",
				"version": "v1",
				"kind": "ConfigMap"
			},
			"resource": {
				"group": "",
				"version": "v1",
				"resource": "configmaps"
			},
			"requestKind": {
				"group": "",
				"version": "v1",
				"kind": "ConfigMap"
			},
			"requestResource": {
				"group": "",
				"version": "v1",
				"resource": "configmaps"
			},
			"name": "example-configmap",
			"namespace": "default",
			"operation": "CREATE",
			"userInfo": {
				"username": "system:admin",
				"groups": [
					"system:masters",
					"system:authenticated"
				]
			},
			"object": {
				"kind": "ConfigMap",
				"apiVersion": "v1",
				"metadata": {
					"name": "example-configmap",
					"namespace": "default",
					"creationTimestamp": null,
					"annotations": {
						"kubectl.kubernetes.io/last-applied-configuration": "{\"apiVersion\":\"v1\",\"data\":{\"key1\":\"value1\",\"key2\":\"value2\",\"key3\":\"value3\"},\"kind\":\"ConfigMap\",\"metadata\":{\"annotations\":{},\"name\":\"example-configmap\",\"namespace\":\"default\"}}\n"
					},
					"managedFields": [
						{
							"manager": "kubectl-client-side-apply",
							"operation": "Update",
							"apiVersion": "v1",
							"time": "2024-03-07T02:43:38Z",
							"fieldsType": "FieldsV1",
							"fieldsV1": {
								"f:data": {
									".": {},
									"f:key1": {},
									"f:key2": {},
									"f:key3": {}
								},
								"f:metadata": {
									"f:annotations": {
										".": {},
										"f:kubectl.kubernetes.io/last-applied-configuration": {}
									}
								}
							}
						}
					]
				},
				"data": {
					"logicalName": "test"
				}
			},
			"oldObject": null,
			"dryRun": false,
			"options": {
				"kind": "CreateOptions",
				"apiVersion": "meta.k8s.io/v1",
				"fieldManager": "kubectl-client-side-apply",
				"fieldValidation": "Strict"
			}
		}
	}`

	values := map[string]string{
		"ClusterName": "test-cluster",
		"Version":     "1.27",
		"Environment": "sbx",
	}
	tp := testPolicies()
	data, err := ProcessAdmissionReview([]byte(rawJSON), values, tp)
	if err == nil {
		ar, err := getAdmissionReview(data)
		if err != nil {
			t.Fail()
		}
		if ar.Request.UID != ar.Response.UID {
			t.Fail()
		}
		hash := getMD5Hash(ar.Response.Patch)
		if hash != "93a40493cb6cb61e4d7e409737099fe6" {
			t.Fail()
		}
	}
}

func TestMutatesConfigMapRequestSSM(t *testing.T) {
	rawJSON := `{
		"kind": "AdmissionReview",
		"apiVersion": "admission.k8s.io/v1",
		"request": {
			"uid": "42b2eea6-2458-421a-ad7d-cd49f90abec5",
			"kind": {
				"group": "",
				"version": "v1",
				"kind": "ConfigMap"
			},
			"resource": {
				"group": "",
				"version": "v1",
				"resource": "configmaps"
			},
			"requestKind": {
				"group": "",
				"version": "v1",
				"kind": "ConfigMap"
			},
			"requestResource": {
				"group": "",
				"version": "v1",
				"resource": "configmaps"
			},
			"name": "example-configmap-ssm",
			"namespace": "default",
			"operation": "CREATE",
			"userInfo": {
				"username": "system:admin",
				"groups": [
					"system:masters",
					"system:authenticated"
				]
			},
			"object": {
				"kind": "ConfigMap",
				"apiVersion": "v1",
				"metadata": {
					"name": "example-configmap-ssm",
					"namespace": "default",
					"creationTimestamp": null,
					"annotations": {
						"kubectl.kubernetes.io/last-applied-configuration": "{\"apiVersion\":\"v1\",\"data\":{\"key1\":\"value1\",\"key2\":\"value2\",\"key3\":\"value3\"},\"kind\":\"ConfigMap\",\"metadata\":{\"annotations\":{},\"name\":\"example-configmap\",\"namespace\":\"default\"}}\n"
					},
					"managedFields": [
						{
							"manager": "kubectl-client-side-apply",
							"operation": "Update",
							"apiVersion": "v1",
							"time": "2024-03-07T02:43:38Z",
							"fieldsType": "FieldsV1",
							"fieldsV1": {
								"f:data": {
									".": {},
									"f:key1": {},
									"f:key2": {},
									"f:key3": {}
								},
								"f:metadata": {
									"f:annotations": {
										".": {},
										"f:kubectl.kubernetes.io/last-applied-configuration": {}
									}
								}
							}
						}
					]
				},
				"data": {
					"logicalName": "test"
				}
			},
			"oldObject": null,
			"dryRun": false,
			"options": {
				"kind": "CreateOptions",
				"apiVersion": "meta.k8s.io/v1",
				"fieldManager": "kubectl-client-side-apply",
				"fieldValidation": "Strict"
			}
		}
	}`

	values := map[string]string{
		"ClusterName": "test-cluster",
		"Version":     "1.27",
		"Environment": "sbx",
	}
	tp := testPolicies()
	data, err := ProcessAdmissionReview([]byte(rawJSON), values, tp)
	if err == nil {
		ar, err := getAdmissionReview(data)
		if err != nil {
			t.Fail()
		}
		if ar.Request.UID != ar.Response.UID {
			t.Fail()
		}
		hash := getMD5Hash(ar.Response.Patch)
		if hash != "b62c19a25078503a5e26105ca5552256" {
			t.Fail()
		}
	}
}

func TestMutatesDaemonSetRequest(t *testing.T) {
	rawJSON := `{
		"kind": "AdmissionReview",
		"apiVersion": "admission.k8s.io/v1",
		"request": {
			"uid": "d6b26140-ace4-4f05-9bc6-ce356cfd9f7d",
			"kind": {
				"group": "apps",
				"version": "v1",
				"kind": "DaemonSet"
			},
			"resource": {
				"group": "apps",
				"version": "v1",
				"resource": "daemonsets"
			},
			"requestKind": {
				"group": "apps",
				"version": "v1",
				"kind": "DaemonSet"
			},
			"requestResource": {
				"group": "apps",
				"version": "v1",
				"resource": "daemonsets"
			},
			"name": "nginx-daemonset",
			"namespace": "test",
			"operation": "CREATE",
			"userInfo": {
				"username": "system:admin",
				"groups": [
					"system:masters",
					"system:authenticated"
				]
			},
			"object": {
				"kind": "DaemonSet",
				"apiVersion": "apps/v1",
				"metadata": {
					"name": "nginx-daemonset",
					"namespace": "test",
					"creationTimestamp": null,
					"labels": {
						"app": "nginx"
					},
					"annotations": {
						"deprecated.daemonset.template.generation": "0",
						"kubectl.kubernetes.io/last-applied-configuration": "{\"apiVersion\":\"apps/v1\",\"kind\":\"DaemonSet\",\"metadata\":{\"annotations\":{},\"labels\":{\"app\":\"nginx\"},\"name\":\"nginx-daemonset\",\"namespace\":\"test\"},\"spec\":{\"selector\":{\"matchLabels\":{\"app\":\"nginx\"}},\"template\":{\"metadata\":{\"labels\":{\"app\":\"nginx\"}},\"spec\":{\"containers\":[{\"image\":\"nginx:latest\",\"name\":\"nginx\",\"ports\":[{\"containerPort\":80}]}]}}}}\n"
					},
					"managedFields": [
						{
							"manager": "kubectl-client-side-apply",
							"operation": "Update",
							"apiVersion": "apps/v1",
							"time": "2024-03-07T04:31:33Z",
							"fieldsType": "FieldsV1",
							"fieldsV1": {
								"f:metadata": {
									"f:annotations": {
										".": {},
										"f:deprecated.daemonset.template.generation": {},
										"f:kubectl.kubernetes.io/last-applied-configuration": {}
									},
									"f:labels": {
										".": {},
										"f:app": {}
									}
								},
								"f:spec": {
									"f:revisionHistoryLimit": {},
									"f:selector": {},
									"f:template": {
										"f:metadata": {
											"f:labels": {
												".": {},
												"f:app": {}
											}
										},
										"f:spec": {
											"f:containers": {
												"k:{\"name\":\"nginx\"}": {
													".": {},
													"f:image": {},
													"f:imagePullPolicy": {},
													"f:name": {},
													"f:ports": {
														".": {},
														"k:{\"containerPort\":80,\"protocol\":\"TCP\"}": {
															".": {},
															"f:containerPort": {},
															"f:protocol": {}
														}
													},
													"f:resources": {},
													"f:terminationMessagePath": {},
													"f:terminationMessagePolicy": {}
												}
											},
											"f:dnsPolicy": {},
											"f:restartPolicy": {},
											"f:schedulerName": {},
											"f:securityContext": {},
											"f:terminationGracePeriodSeconds": {}
										}
									},
									"f:updateStrategy": {
										"f:rollingUpdate": {
											".": {},
											"f:maxSurge": {},
											"f:maxUnavailable": {}
										},
										"f:type": {}
									}
								}
							}
						}
					]
				},
				"spec": {
					"selector": {
						"matchLabels": {
							"app": "nginx"
						}
					},
					"template": {
						"metadata": {
							"creationTimestamp": null,
							"labels": {
								"app": "nginx"
							}
						},
						"spec": {
							"containers": [
								{
									"name": "nginx",
									"image": "nginx:latest",
									"ports": [
										{
											"containerPort": 80,
											"protocol": "TCP"
										}
									],
									"resources": {},
									"terminationMessagePath": "/dev/termination-log",
									"terminationMessagePolicy": "File",
									"imagePullPolicy": "Always"
								}
							],
							"restartPolicy": "Always",
							"terminationGracePeriodSeconds": 30,
							"dnsPolicy": "ClusterFirst",
							"securityContext": {},
							"schedulerName": "default-scheduler"
						}
					},
					"updateStrategy": {
						"type": "RollingUpdate",
						"rollingUpdate": {
							"maxUnavailable": 1,
							"maxSurge": 0
						}
					},
					"revisionHistoryLimit": 10
				},
				"status": {
					"currentNumberScheduled": 0,
					"numberMisscheduled": 0,
					"desiredNumberScheduled": 0,
					"numberReady": 0
				}
			},
			"oldObject": null,
			"dryRun": false,
			"options": {
				"kind": "CreateOptions",
				"apiVersion": "meta.k8s.io/v1",
				"fieldManager": "kubectl-client-side-apply",
				"fieldValidation": "Strict"
			}
		}
	}`

	values := map[string]string{
		"ClusterName": "test-cluster",
		"Version":     "1.27",
		"Environment": "sbx",
	}
	tp := testPolicies()
	data, err := ProcessAdmissionReview([]byte(rawJSON), values, tp)
	if err == nil {
		ar, err := getAdmissionReview(data)
		if err != nil {
			t.Fail()
		}
		if ar.Request.UID != ar.Response.UID {
			t.Fail()
		}
		hash := getMD5Hash(ar.Response.Patch)
		if hash != "78d237c223420aefb9f540ce20650e53" {
			t.Fail()
		}
	}
}

func getAdmissionReview(data []byte) (*admissionv1.AdmissionReview, error) {
	var adminReview *admissionv1.AdmissionReview
	if err := json.Unmarshal(data, &adminReview); err != nil {
		return nil, fmt.Errorf("unable unmarshal pod json object %v", err)
	}
	return adminReview, nil
}

func getMD5Hash(data []byte) string {
	hash := md5.Sum(data)
	return hex.EncodeToString(hash[:])
}

func testPolicies() policies.Policies {
	return policies.Policies{
		Deployments: []policies.Policy{
			{
				Namespace: "default",
				Name:      "nginx-deployment",
				Key:       "CLUSTER_NAME",
				Value:     "{{ .ClusterName }}",
				Type:      "env",
			},
		},
		DaemonSets: []policies.Policy{
			{
				Namespace: "test",
				Name:      "nginx-daemonset",
				Key:       "AQUA_LOGICAL_NAME",
				Value:     "{{ .ClusterName }}",
				Type:      "env",
			},
		},
		ConfigMaps: []policies.Policy{
			{
				Namespace: "default",
				Name:      "example-configmap",
				Key:       "logicalName",
				Value:     "{{ .ClusterName }}",
			},
			{
				Namespace: "default",
				Name:      "example-configmap-ssm",
				Key:       "logicalName",
				Value:     "",
				SSM: policies.SSMParameter{
					Region:  "us-east-2",
					Name:    "/cluster/{{ .Version }}/license",
					Decrypt: false,
				},
			},
		},
	}
}
