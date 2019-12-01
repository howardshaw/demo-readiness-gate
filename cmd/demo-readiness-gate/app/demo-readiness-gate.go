package app

import (
	"flag"
	"fmt"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/component-base/logs"
	"k8s.io/klog"
	"net/http"
	"os"
	"sigs.k8s.io/controller-runtime/pkg/controller"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/source"
	"sigs.k8s.io/controller-runtime/pkg/webhook"

	"github.com/howardshaw/demo-readiness-gate/cmd/demo-readiness-gate/app/options"
	readinessController "github.com/howardshaw/demo-readiness-gate/pkg/controller"
	"github.com/howardshaw/demo-readiness-gate/pkg/version"
	corev1 "k8s.io/api/core/v1"
	"sigs.k8s.io/controller-runtime/pkg/manager"

	readinessWebhook "github.com/howardshaw/demo-readiness-gate/pkg/webhook"
)

var (
	kubeconfig, kubeFedConfig, masterURL string
)

// NewControllerManagerCommand creates a *cobra.Command object with default parameters
func NewReadinessGateControllerCommand(stopChan <-chan struct{}) *cobra.Command {
	verFlag := false
	opts := options.NewOptions()

	cmd := &cobra.Command{
		Use:  "demo-readiness-gate",
		Long: "demo readiness gate for pod.",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Fprintf(os.Stdout, "Demo readiness gate version: %s\n", fmt.Sprintf("%#v", version.Get()))
			if verFlag {
				os.Exit(0)
			}
			PrintFlags(cmd.Flags())

			if err := Run(opts, stopChan); err != nil {
				fmt.Fprintf(os.Stderr, "%v\n", err)
				os.Exit(1)
			}
		},
	}

	// Add the command line flags from other dependencies(klog, kubebuilder, etc.)
	cmd.Flags().AddGoFlagSet(flag.CommandLine)

	opts.AddFlags(cmd.Flags())
	cmd.Flags().BoolVar(&verFlag, "version", false, "Prints the Version info of controller-manager")
	cmd.Flags().StringVar(&kubeFedConfig, "kubefed-config", "", "Path to a KubeFedConfig yaml file. Test only.")
	cmd.Flags().StringVar(&kubeconfig, "kubeconfig", "", "Path to a kubeconfig. Only required if out-of-cluster.")
	cmd.Flags().StringVar(&masterURL, "master", "", "The address of the Kubernetes API server. Overrides any value in kubeconfig. Only required if out-of-cluster.")

	return cmd
}

// Run runs the controller-manager with options. This should never exit.
func Run(opts *options.Options, stopChan <-chan struct{}) error {
	logs.InitLogs()
	defer logs.FlushLogs()

	// TODO: Make healthz endpoint configurable
	go serveHealthz(":18080")

	restConfig, err := clientcmd.BuildConfigFromFlags(masterURL, kubeconfig)
	if err != nil {
		panic(err)
	}

	klog.Info("setting up a manager")
	mgr, err := manager.New(restConfig, manager.Options{})
	if err != nil {
		klog.Error(err, "unable to set up overall controller manager")
		os.Exit(1)
	}
	klog.Info("Setting up controller")
	c, err := controller.New("demo-readiness-gate-controller", mgr, controller.Options{
		Reconciler: readinessController.NewReconcilePod(mgr.GetClient()),
	})
	if err != nil {
		klog.Error(err, "unable to set up controller")
		os.Exit(1)
	}

	// Watch Pod and enqueue Pod object key
	if err := c.Watch(&source.Kind{Type: &corev1.Pod{}}, &handler.EnqueueRequestForObject{}); err != nil {
		klog.Error(err, "unable to watch pod")
		os.Exit(1)
	}

	// Setup webhooks
	klog.Info("setting up webhook server")
	hookServer := mgr.GetWebhookServer()

	klog.Info("registering webhooks to the webhook server")
	hookServer.Register("/mutate-pod-readiness-gate", &webhook.Admission{Handler: &readinessWebhook.PodReadinessGate{}})

	klog.Info("starting manager")
	if err := mgr.Start(stopChan); err != nil {
		klog.Error(err, "unable to run manager")
		return err
	}
	return nil
}

// PrintFlags logs the flags in the flagset
func PrintFlags(flags *pflag.FlagSet) {
	flags.VisitAll(func(flag *pflag.Flag) {
		klog.Infof("FLAG: --%s=%q", flag.Name, flag.Value)
	})
}

func serveHealthz(address string) {
	http.HandleFunc("/healthz", func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("OK"))
	})

	klog.Fatal(http.ListenAndServe(address, nil))
}
