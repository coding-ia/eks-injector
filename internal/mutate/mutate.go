package mutate

import (
	"eks-inject/internal/policies"
	"encoding/json"
	"errors"
	"fmt"
	admissionv1 "k8s.io/api/admission/v1"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"log"
)

type PatchOperation struct {
	Op    string      `json:"op"`
	Path  string      `json:"path"`
	Value interface{} `json:"value,omitempty"`
}

func ProcessAdmissionReview(content []byte, clusterName string) ([]byte, error) {
	adminReview := admissionv1.AdmissionReview{}

	if clusterName == "" {
		return nil, errors.New("cluster name is not defined")
	}

	if err := json.Unmarshal(content, &adminReview); err != nil {
		return nil, fmt.Errorf("unmarshaling request failed with %s", err)
	}

	if adminReview.Request == nil {
		return nil, errors.New("no admin request available")
	}

	adminResponse := admissionv1.AdmissionResponse{
		Allowed: true,
		UID:     adminReview.Request.UID,
	}

	var err error
	switch adminReview.Request.Kind.Kind {
	case "Deployment":
		err = mutateDeployment(adminReview.Request, &adminResponse, clusterName)
	case "ConfigMap":
		err = mutateConfigMap(adminReview.Request, &adminResponse, clusterName)
	case "DaemonSet":
		err = mutateDaemonSet(adminReview.Request, &adminResponse, clusterName)
	default:
		log.Println("Unable to process resource.")
	}

	if err != nil {
		log.Printf("Error: %s", err)
	}

	adminReview.Response = &adminResponse
	responseBody, err := json.Marshal(adminReview)
	if err != nil {
		return nil, err
	}

	return responseBody, nil
}

func mutateDeployment(request *admissionv1.AdmissionRequest, response *admissionv1.AdmissionResponse, clusterName string) error {
	var deployment *appsv1.Deployment
	if err := json.Unmarshal(request.Object.Raw, &deployment); err != nil {
		return fmt.Errorf("unable unmarshal pod json object %v", err)
	}

	policy, _ := policies.FindDeploymentPolicy(deployment.ObjectMeta.Namespace, deployment.ObjectMeta.Name, "env")
	if policy == nil {
		return nil
	}

	pT := admissionv1.PatchTypeJSONPatch
	response.PatchType = &pT

	var patches []PatchOperation
	found := false
	for cdx, container := range deployment.Spec.Template.Spec.Containers {
		envVar := 0

		for edx, env := range container.Env {
			envVar++

			if env.Name == policy.Key {
				patches = append(patches, PatchOperation{
					Op:    "replace",
					Path:  fmt.Sprintf("/spec/template/spec/containers/%d/env/%d/value", cdx, edx),
					Value: clusterName,
				})
				found = true
				break
			}
		}

		if !found {
			if envVar == 0 {
				patches = append(patches, PatchOperation{
					Op:   "add",
					Path: fmt.Sprintf("/spec/template/spec/containers/%d/env", cdx),
					Value: []corev1.EnvVar{
						{
							Name:  policy.Key,
							Value: clusterName,
						},
					},
				})
			} else {
				patches = append(patches, PatchOperation{
					Op:    "add",
					Path:  fmt.Sprintf("/spec/template/spec/containers/%d/env/-", cdx),
					Value: corev1.EnvVar{Name: policy.Key, Value: clusterName},
				})
			}
		}
	}

	patchBytes, err := json.Marshal(patches)
	if err != nil {
		return err
	}

	response.Patch = patchBytes
	response.Result = &metav1.Status{
		Status: "Success",
	}

	return nil
}

func mutateDaemonSet(request *admissionv1.AdmissionRequest, response *admissionv1.AdmissionResponse, clusterName string) error {
	var daemonSet *appsv1.DaemonSet
	if err := json.Unmarshal(request.Object.Raw, &daemonSet); err != nil {
		return fmt.Errorf("unable unmarshal pod json object %v", err)
	}

	policy, _ := policies.FindDaemonSetPolicy(daemonSet.ObjectMeta.Namespace, daemonSet.ObjectMeta.Name, "env")
	if policy == nil {
		return nil
	}

	pT := admissionv1.PatchTypeJSONPatch
	response.PatchType = &pT

	var patches []PatchOperation
	found := false
	for cdx, container := range daemonSet.Spec.Template.Spec.Containers {
		envVar := 0

		for edx, env := range container.Env {
			envVar++

			if env.Name == policy.Key {
				patches = append(patches, PatchOperation{
					Op:    "replace",
					Path:  fmt.Sprintf("/spec/template/spec/containers/%d/env/%d/value", cdx, edx),
					Value: clusterName,
				})
				found = true
				break
			}
		}

		if !found {
			if envVar == 0 {
				patches = append(patches, PatchOperation{
					Op:   "add",
					Path: fmt.Sprintf("/spec/template/spec/containers/%d/env", cdx),
					Value: []corev1.EnvVar{
						{Name: policy.Key, Value: clusterName},
					},
				})
			} else {
				patches = append(patches, PatchOperation{
					Op:    "add",
					Path:  fmt.Sprintf("/spec/template/spec/containers/%d/env/-", cdx),
					Value: corev1.EnvVar{Name: policy.Key, Value: clusterName},
				})
			}
		}
	}

	patchBytes, err := json.Marshal(patches)
	if err != nil {
		return err
	}

	response.Patch = patchBytes
	response.Result = &metav1.Status{
		Status: "Success",
	}

	return nil
}

func mutateConfigMap(request *admissionv1.AdmissionRequest, response *admissionv1.AdmissionResponse, clusterName string) error {
	var configMap *corev1.ConfigMap
	if err := json.Unmarshal(request.Object.Raw, &configMap); err != nil {
		return fmt.Errorf("unable unmarshal pod json object %v", err)
	}

	policy, _ := policies.FindConfigMapPolicy(configMap.ObjectMeta.Namespace, configMap.ObjectMeta.Name, "")
	if policy == nil {
		return nil
	}

	pT := admissionv1.PatchTypeJSONPatch
	response.PatchType = &pT

	var patches []PatchOperation
	found := false
	for key, _ := range configMap.Data {
		if key == policy.Key {
			patches = append(patches, PatchOperation{
				Op:    "replace",
				Path:  fmt.Sprintf("/data/%s", key),
				Value: clusterName,
			})
			found = true
			break
		}
	}

	if !found {
		patches = append(patches, PatchOperation{
			Op:    "add",
			Path:  fmt.Sprintf("/data/%s", policy.Key),
			Value: clusterName,
		})
	}

	patchBytes, err := json.Marshal(patches)
	if err != nil {
		return err
	}

	response.Patch = patchBytes
	response.Result = &metav1.Status{
		Status: "Success",
	}

	return nil
}
