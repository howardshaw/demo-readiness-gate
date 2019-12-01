package webhook

import (
	"encoding/json"
	"fmt"
	"k8s.io/api/admission/v1beta1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/klog"
	"net/http"
	"time"
)

const DemoReadinessGateName = "www.example.com/feature-demo"

type patchReadness struct {
	Op    string                  `json:"op"`
	Path  string                  `json:"path"`
	Value corev1.PodReadinessGate `json:"value"`
}
type patchReadnesses struct {
	Op    string                    `json:"op"`
	Path  string                    `json:"path"`
	Value []corev1.PodReadinessGate `json:"value"`
}

func demoReadinessGate(ar v1beta1.AdmissionReview) *v1beta1.AdmissionResponse {
	klog.V(2).Info("calling add demo readiness gate")
	var pod *corev1.Pod
	raw := ar.Request.Object.Raw
	err := json.Unmarshal(raw, &pod)
	if err != nil {
		klog.Error(err)
		return toAdmissionResponse(err)
	}

	reviewResponse := v1beta1.AdmissionResponse{}
	reviewResponse.Allowed = true

	if len(pod.Spec.ReadinessGates) == 0 {
		patchReadness := []patchReadnesses{{
			Op:   "add",
			Path: "/spec/readinessGates",
			Value: []corev1.PodReadinessGate{
				{
					ConditionType: DemoReadinessGateName,
				},
			},
		}}
		reviewResponse.Patch, err = json.Marshal(patchReadness)
		if err != nil {
			klog.Error(err)
			return toAdmissionResponse(err)
		}
	} else {
		patchReadness := []patchReadness{{
			Op:   "add",
			Path: fmt.Sprintf("/spec/readinessGates/%d", len(pod.Spec.ReadinessGates)),
			Value: corev1.PodReadinessGate{

				ConditionType: DemoReadinessGateName,
			},
		}}
		reviewResponse.Patch, err = json.Marshal(patchReadness)
		if err != nil {
			klog.Error(err)
			return toAdmissionResponse(err)
		}
	}

	pt := v1beta1.PatchTypeJSONPatch
	reviewResponse.PatchType = &pt
	reviewResponse.Result = &metav1.Status{Status: "Success"}
	reviewResponse.AuditAnnotations = map[string]string{"demo-webhook": "add-label"}
	return &reviewResponse
}

func ServeDemoReadinessGate(w http.ResponseWriter, r *http.Request) {
	serve(w, r, demoReadinessGate)
}

func GenerateReadinessGateCondition() corev1.PodCondition {
	return corev1.PodCondition{
		Type:               DemoReadinessGateName,
		Status:             corev1.ConditionTrue,
		LastProbeTime:      metav1.Time{Time: time.Now()},
		LastTransitionTime: metav1.Time{Time: time.Now()},
		Reason:             "ready",
		Message:            "success",
	}
}
