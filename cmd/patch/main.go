package main

import (
	"flag"
	"fmt"
	"github.com/howardshaw/demo-readiness-gate/pkg/webhook"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"time"
)

var (
	kubeconfig, masterURL string
)

func main() {
	addFlags()
	flag.Parse()

	var config *rest.Config
	var err error

	if kubeconfig != "" || masterURL != "" {
		config, err = clientcmd.BuildConfigFromFlags(masterURL, kubeconfig)
		if err != nil {
			panic(err.Error())
		}
	} else {
		config, err = rest.InClusterConfig()
		if err != nil {
			panic(err.Error())
		}
	}

	// create the clientset
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err.Error())
	}

	_, err = clientset.ServerVersion()
	if err != nil {
		panic(err.Error())
	}
	pods, err := clientset.CoreV1().Pods(corev1.NamespaceDefault).List(metav1.ListOptions{})
	if err != nil {
		panic(err.Error())
	}
	gatePodCondition := corev1.PodCondition{
		Type:               webhook.DemoReadinessGateName,
		Status:             corev1.ConditionTrue,
		LastProbeTime:      metav1.Time{Time: time.Now()},
		LastTransitionTime: metav1.Time{Time: time.Now()},
		Reason:             "ready",
		Message:            "success",
	}

	for _, pod := range pods.Items {
		pod.Status.Conditions = append(pod.Status.Conditions, gatePodCondition)
		_, err := clientset.CoreV1().Pods(pod.GetNamespace()).UpdateStatus(&pod)
		if errors.IsNotFound(err) {
			fmt.Printf("Pod %s in namespace %s not found\n", pod.GetName(), pod.GetNamespace())
		} else if statusError, isStatus := err.(*errors.StatusError); isStatus {
			fmt.Printf("Error getting pod %s in namespace %s: %v\n",
				pod.GetName(), pod.GetNamespace(), statusError.ErrStatus.Message)
		} else if err != nil {
			panic(err.Error())
		} else {
			fmt.Printf("Patched readiness gate for pod %s in namespace %s\n", pod.GetName(), pod.GetNamespace())
		}
	}
}

func addFlags() {
	flag.StringVar(&kubeconfig, "kubeconfig", "", "Path to a kubeconfig. Only required if out-of-cluster.")
	flag.StringVar(&masterURL, "master", "", "The address of the Kubernetes API server. Overrides any value in kubeconfig. Only required if out-of-cluster.")
}
