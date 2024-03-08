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
	var responseBody []byte
	admReview := admissionv1.AdmissionReview{}

	if clusterName == "" {
		return nil, errors.New("cluster name is not defined")
	}

	if err := json.Unmarshal(content, &admReview); err != nil {
		return nil, fmt.Errorf("unmarshaling request failed with %s", err)
	}

	switch admReview.Request.Kind.Kind {
	case "Deployment":
		responseBody, err := mutateDeployment(admReview, clusterName)
		return responseBody, err
	case "ConfigMap":
		responseBody, err := mutateConfigMap(admReview, clusterName)
		return responseBody, err
	case "DaemonSet":
		responseBody, err := mutateDaemonSet(admReview, clusterName)
		return responseBody, err
	default:
		log.Println("Unable to process resource.")
	}

	return responseBody, nil
}

func mutateDeployment(admReview admissionv1.AdmissionReview, clusterName string) ([]byte, error) {
	var responseBody []byte
	ar := admReview.Request

	if ar != nil {
		resp := admissionv1.AdmissionResponse{}
		var deployment *appsv1.Deployment
		if err := json.Unmarshal(ar.Object.Raw, &deployment); err != nil {
			return nil, fmt.Errorf("unable unmarshal pod json object %v", err)
		}

		policy, _ := policies.FindDeploymentPolicy(deployment.ObjectMeta.Namespace, deployment.ObjectMeta.Name, "env")
		if policy == nil {
			return responseBody, nil
		}

		resp.Allowed = true
		resp.UID = ar.UID
		pT := admissionv1.PatchTypeJSONPatch
		resp.PatchType = &pT

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
			return nil, err
		}

		resp.Patch = patchBytes
		resp.Result = &metav1.Status{
			Status: "Success",
		}

		admReview.Response = &resp
		response, err := json.Marshal(admReview)
		if err != nil {
			return nil, err
		}

		responseBody = response
	}

	return responseBody, nil
}

func mutateDaemonSet(admReview admissionv1.AdmissionReview, clusterName string) ([]byte, error) {
	var responseBody []byte
	ar := admReview.Request

	if ar != nil {
		resp := admissionv1.AdmissionResponse{}
		var daemonSet *appsv1.DaemonSet
		if err := json.Unmarshal(ar.Object.Raw, &daemonSet); err != nil {
			return nil, fmt.Errorf("unable unmarshal pod json object %v", err)
		}

		resp.Allowed = true
		resp.UID = ar.UID

		policy, _ := policies.FindDaemonSetPolicy(daemonSet.ObjectMeta.Namespace, daemonSet.ObjectMeta.Name, "env")
		if policy == nil {
			return responseBody, nil
		}

		pT := admissionv1.PatchTypeJSONPatch
		resp.PatchType = &pT

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
			return nil, err
		}

		resp.Patch = patchBytes
		resp.Result = &metav1.Status{
			Status: "Success",
		}

		admReview.Response = &resp
		response, err := json.Marshal(admReview)
		if err != nil {
			return nil, err
		}

		responseBody = response
	}

	return responseBody, nil
}

func mutateConfigMap(admReview admissionv1.AdmissionReview, clusterName string) ([]byte, error) {
	var responseBody []byte
	ar := admReview.Request

	if ar != nil {
		resp := admissionv1.AdmissionResponse{}
		var configMap *corev1.ConfigMap
		if err := json.Unmarshal(ar.Object.Raw, &configMap); err != nil {
			return nil, fmt.Errorf("unable unmarshal pod json object %v", err)
		}

		policy, _ := policies.FindConfigMapPolicy(configMap.ObjectMeta.Namespace, configMap.ObjectMeta.Name, "")
		if policy == nil {
			return responseBody, nil
		}

		resp.Allowed = true
		resp.UID = ar.UID
		pT := admissionv1.PatchTypeJSONPatch
		resp.PatchType = &pT

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
			return nil, err
		}

		resp.Patch = patchBytes
		resp.Result = &metav1.Status{
			Status: "Success",
		}

		admReview.Response = &resp
		response, err := json.Marshal(admReview)
		if err != nil {
			return nil, err
		}

		responseBody = response
	}

	return responseBody, nil
}
