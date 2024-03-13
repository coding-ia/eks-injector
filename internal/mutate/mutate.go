package mutate

import (
	"eks-injector/internal/parameter"
	"eks-injector/internal/policies"
	"eks-injector/internal/string_parser"
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

func ProcessAdmissionReview(content []byte, variables map[string]string, p policies.Policies) ([]byte, error) {
	adminReview := admissionv1.AdmissionReview{}

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
		err = mutateDeployment(adminReview.Request, &adminResponse, variables, p)
	case "ConfigMap":
		err = mutateConfigMap(adminReview.Request, &adminResponse, variables, p)
	case "DaemonSet":
		err = mutateDaemonSet(adminReview.Request, &adminResponse, variables, p)
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

func mutateDeployment(request *admissionv1.AdmissionRequest, response *admissionv1.AdmissionResponse, variables map[string]string, p policies.Policies) error {
	var deployment *appsv1.Deployment
	if err := json.Unmarshal(request.Object.Raw, &deployment); err != nil {
		return fmt.Errorf("unable unmarshal pod json object %v", err)
	}

	matchingPolicies := policies.FindDeploymentPolicy(p.Deployments, deployment.ObjectMeta.Namespace, deployment.ObjectMeta.Name, "env")
	patches, err := createEnvironmentPatches(matchingPolicies, deployment.Spec.Template.Spec.Containers, variables)
	if err != nil {
		return err
	}

	if len(patches) > 0 {
		err := setAdmissionResponsePatches(response, patches)
		if err != nil {
			return err
		}
	}

	return nil
}

func mutateDaemonSet(request *admissionv1.AdmissionRequest, response *admissionv1.AdmissionResponse, variables map[string]string, p policies.Policies) error {
	var daemonSet *appsv1.DaemonSet
	if err := json.Unmarshal(request.Object.Raw, &daemonSet); err != nil {
		return fmt.Errorf("unable unmarshal pod json object %v", err)
	}

	matchingPolicies := policies.FindDaemonSetPolicy(p.DaemonSets, daemonSet.ObjectMeta.Namespace, daemonSet.ObjectMeta.Name, "env")
	patches, err := createEnvironmentPatches(matchingPolicies, daemonSet.Spec.Template.Spec.Containers, variables)
	if err != nil {
		return err
	}

	if len(patches) > 0 {
		err := setAdmissionResponsePatches(response, patches)
		if err != nil {
			return err
		}
	}

	return nil
}

func mutateConfigMap(request *admissionv1.AdmissionRequest, response *admissionv1.AdmissionResponse, variables map[string]string, p policies.Policies) error {
	var configMap *corev1.ConfigMap
	if err := json.Unmarshal(request.Object.Raw, &configMap); err != nil {
		return fmt.Errorf("unable unmarshal pod json object %v", err)
	}

	matchingPolicies := policies.FindConfigMapPolicy(p.ConfigMaps, configMap.ObjectMeta.Namespace, configMap.ObjectMeta.Name)
	var patches []PatchOperation

	for _, policy := range matchingPolicies {
		value, err := getValue(policy, variables)
		if err != nil {
			return err
		}

		found := false
		configVar := 0
		for key := range configMap.Data {
			configVar++
			if key == policy.Key {
				patches = append(patches, PatchOperation{
					Op:    "replace",
					Path:  fmt.Sprintf("/data/%s", key),
					Value: value,
				})
				found = true
				break
			}
		}

		if !found {
			if configVar == 0 {
				patches = append(patches, PatchOperation{
					Op:   "add",
					Path: "/data",
					Value: map[string]string{
						policy.Key: value,
					},
				})
			} else {
				patches = append(patches, PatchOperation{
					Op:    "add",
					Path:  fmt.Sprintf("/data/%s", policy.Key),
					Value: value,
				})
			}
		}
	}

	if len(patches) > 0 {
		err := setAdmissionResponsePatches(response, patches)
		if err != nil {
			return err
		}
	}

	return nil
}

func createEnvironmentPatches(policies []policies.Policy, containers []corev1.Container, variables map[string]string) ([]PatchOperation, error) {
	var patches []PatchOperation

	for _, policy := range policies {
		value, err := getValue(policy, variables)
		if err != nil {
			return patches, err
		}

		found := false
		for cdx, container := range containers {
			envVar := 0

			for edx, env := range container.Env {
				envVar++

				if env.Name == policy.Key {
					patches = append(patches, PatchOperation{
						Op:    "replace",
						Path:  fmt.Sprintf("/spec/template/spec/containers/%d/env/%d/value", cdx, edx),
						Value: value,
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
							{Name: policy.Key, Value: value},
						},
					})
				} else {
					patches = append(patches, PatchOperation{
						Op:    "add",
						Path:  fmt.Sprintf("/spec/template/spec/containers/%d/env/-", cdx),
						Value: corev1.EnvVar{Name: policy.Key, Value: value},
					})
				}
			}
		}
	}

	return patches, nil
}

func setAdmissionResponsePatches(response *admissionv1.AdmissionResponse, patches []PatchOperation) error {
	patchBytes, err := json.Marshal(patches)
	if err != nil {
		return err
	}

	pT := admissionv1.PatchTypeJSONPatch
	response.PatchType = &pT
	response.Patch = patchBytes
	response.Result = &metav1.Status{
		Status: "Success",
	}

	return nil
}

func getValue(policy policies.Policy, variables map[string]string) (string, error) {
	if policy.Value == "" && policy.SSM.Name != "" {
		parameterName, err := string_parser.ParseString(policy.SSM.Name, variables)
		if err != nil {
			return "", err
		}
		value, err := parameter.GetParameter(policy.SSM.Region, parameterName, policy.SSM.Decrypt, policy.SSM.AssumeRole)
		return value, err
	}
	value, err := string_parser.ParseString(policy.Value, variables)
	if err != nil {
		return "", err
	}

	return value, nil
}
