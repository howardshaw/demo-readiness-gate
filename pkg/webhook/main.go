package webhook

import (
	"context"
	"encoding/json"
	"k8s.io/klog"
	"net/http"

	corev1 "k8s.io/api/core/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/webhook/admission"
)

// PodReadinessGate add demo readiness gate for Pods
type PodReadinessGate struct {
	client  client.Client
	decoder *admission.Decoder
}

// PodReadinessGate adds an demo readiness gate to every incoming pods.
func (a *PodReadinessGate) Handle(ctx context.Context, req admission.Request) admission.Response {
	pod := &corev1.Pod{}

	err := a.decoder.Decode(req, pod)
	if err != nil {
		return admission.Errored(http.StatusBadRequest, err)
	}

	if len(pod.Spec.ReadinessGates) == 0 {
		pod.Spec.ReadinessGates = make([]corev1.PodReadinessGate, 0)
	}
	pod.Spec.ReadinessGates = append(pod.Spec.ReadinessGates, corev1.PodReadinessGate{ConditionType: DemoReadinessGateName})

	marshaledPod, err := json.Marshal(pod)
	klog.Infof("webhook pod: %s", marshaledPod)
	if err != nil {
		return admission.Errored(http.StatusInternalServerError, err)
	}

	return admission.PatchResponseFromRaw(req.Object.Raw, marshaledPod)
}

// PodReadinessGate implements inject.Client.
// A client will be automatically injected.

// InjectClient injects the client.
func (a *PodReadinessGate) InjectClient(c client.Client) error {
	a.client = c
	return nil
}

// PodReadinessGate implements admission.DecoderInjector.
// A decoder will be automatically injected.

// InjectDecoder injects the decoder.
func (a *PodReadinessGate) InjectDecoder(d *admission.Decoder) error {
	a.decoder = d
	return nil
}
