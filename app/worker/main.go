package main

import (
	"fmt"
	"strings"

	"github.com/ivan-garcia-santamaria/cadence-api-example/app/commons"
	"github.com/ivan-garcia-santamaria/cadence-api-example/app/commons/grpc"
	"github.com/ivan-garcia-santamaria/cadence-api-example/app/config"
	"github.com/ivan-garcia-santamaria/cadence-api-example/app/worker/workflows"

	"go.uber.org/cadence/client"
	"go.uber.org/cadence/worker"
	"go.uber.org/cadence/workflow"
	"go.uber.org/zap"
)

func startWorkersThrift(h *commons.CadenceAdapter, taskList string) {
	// Configure worker options.
	workerOptions := worker.Options{
		MetricsScope:          h.Scope,
		Logger:                h.Logger,
		DisableActivityWorker: true,
	}

	cadenceWorker := worker.New(h.ServiceClient, h.Config.Domain, taskList, workerOptions)
	err := cadenceWorker.Start()
	if err != nil {
		h.Logger.Error("Failed to start workers.", zap.Error(err))
		panic("Failed to start workers")
	}
}

func startWorkers(h *grpc.CadenceAdapter, taskList string) {

	workerOptions := worker.Options{
		MetricsScope: h.WorkerMetricScope,
		Logger:       h.Logger,
		FeatureFlags: client.FeatureFlags{
			WorkflowExecutionAlreadyCompletedErrorEnabled: true,
		},
		DisableActivityWorker: true,
	}

	h.StartWorkers(h.Config.DomainName, taskList, workerOptions)
}

func main() {
	fmt.Println("Starting Worker..")
	var appConfig config.AppConfig
	appConfig.Setup()

	if strings.Contains(appConfig.Cadence.HostPort, ":7833") {

		var cadenceClient grpc.CadenceAdapter
		cadenceClient.SetupServiceConfig(&appConfig.Cadence)

		cadenceClient.RegisterWorkflow(workflows.Workflow)
		startWorkers(&cadenceClient, workflows.TaskListName)
	} else {
		var cadenceClient commons.CadenceAdapter
		cadenceClient.Setup(&appConfig.Cadence)

		workflow.Register(workflows.Workflow)
		startWorkersThrift(&cadenceClient, workflows.TaskListName)

	}

	// The workers are supposed to be long running process that should not exit.
	select {}
}
