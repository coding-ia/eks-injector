package mutate

import "testing"

func TestMutatesValidRequest(t *testing.T) {
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

	data, err := ProcessAdmissionReview([]byte(rawJSON), "test-cluster")
	_ = data
	_ = err
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

	data, err := ProcessAdmissionReview([]byte(rawJSON), "test-cluster")
	_ = data
	_ = err
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

	data, err := ProcessAdmissionReview([]byte(rawJSON), "test-cluster")
	_ = data
	_ = err
}
