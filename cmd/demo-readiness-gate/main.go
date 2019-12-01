package main

import (
	"fmt"
	"os"

	genericapiserver "k8s.io/apiserver/pkg/server"
	_ "k8s.io/client-go/plugin/pkg/client/auth" // Load all client auth plugins for GCP, Azure, Openstack, etc
	"k8s.io/component-base/logs"

	"github.com/howardshaw/demo-readiness-gate/cmd/demo-readiness-gate/app"
)

// Controller-manager main.
func main() {
	logs.InitLogs()
	defer logs.FlushLogs()

	stopChan := genericapiserver.SetupSignalHandler()

	if err := app.NewReadinessGateControllerCommand(stopChan).Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}
}
