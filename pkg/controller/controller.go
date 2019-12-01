package controller

import (
	"context"
	"github.com/howardshaw/demo-readiness-gate/pkg/webhook"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/klog"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

// ReconcilePod reconciles Pod
type ReconcilePod struct {
	// client can be used to retrieve objects from the APIServer.
	client client.Client
}

// Implement reconcile.Reconciler so the controller can reconcile objects
var _ reconcile.Reconciler = &ReconcilePod{}

func NewReconcilePod(client client.Client) *ReconcilePod {
	return &ReconcilePod{client: client}
}

func (r *ReconcilePod) Reconcile(request reconcile.Request) (reconcile.Result, error) {
	// Fetch the ReplicaSet from the cache
	//rs := &appsv1.ReplicaSet{}
	pod := &corev1.Pod{}
	err := r.client.Get(context.TODO(), request.NamespacedName, pod)
	if errors.IsNotFound(err) {
		klog.Error(nil, "Could not find pod")
		return reconcile.Result{}, nil
	}

	if err != nil {
		klog.Error(err, "Could not fetch pod")
		return reconcile.Result{}, err
	}

	// do something

	// update condition
	for _, readiness := range pod.Spec.ReadinessGates {
		if readiness.ConditionType == webhook.DemoReadinessGateName {
			// Print the Pod
			klog.Info("Reconciling Pod ", request.NamespacedName)
			demoReadinessCondition := false
			persistStatus := false
			for i, condition := range pod.Status.Conditions {
				if condition.Type == webhook.DemoReadinessGateName {
					demoReadinessCondition = true
					if condition.Status != corev1.ConditionTrue {
						pod.Status.Conditions[i] = webhook.GenerateReadinessGateCondition()
						persistStatus = true
					}
				}
			}

			// Update the condition
			if !demoReadinessCondition {
				if pod.Status.Conditions == nil {
					pod.Status.Conditions = make([]corev1.PodCondition, 0)
				}
				pod.Status.Conditions = append(pod.Status.Conditions, webhook.GenerateReadinessGateCondition())
				persistStatus = true
			}
			// need to update status
			if persistStatus {
				klog.Info("Need to update status ", request.NamespacedName)
				if err := r.client.Status().Update(context.TODO(), pod); err != nil {
					klog.Error(err, "Could not update readiness gate condition")
					return reconcile.Result{}, err
				}
			}
		}
	}

	return reconcile.Result{}, nil
}
