package main

import (
	"fmt"

	"github.com/ivan-garcia-santamaria/cadence-api-example/app/commons/grpc"
	"github.com/ivan-garcia-santamaria/cadence-api-example/app/config"
	"github.com/ivan-garcia-santamaria/cadence-api-example/app/worker/workflows"

	"go.uber.org/cadence/client"
	"go.uber.org/cadence/worker"
)

func startWorkers(h *grpc.CadenceAdapter, taskList string) {

	workerOptions := worker.Options{
		MetricsScope: h.WorkerMetricScope,
		Logger:       h.Logger,
		FeatureFlags: client.FeatureFlags{
			WorkflowExecutionAlreadyCompletedErrorEnabled: true,
		},
		DisableActivityWorker: true,
	}

	// Configure worker options.
	// workerOptions := worker.Options{
	// 	MetricsScope:          h.WorkerMetricScope,
	// 	Logger:                h.Logger,
	// 	DisableWorkflowWorker: true,
	// }

	h.StartWorkers(h.Config.DomainName, taskList, workerOptions)
	// cadenceWorker := worker.New(h.ServiceClient, h.Config.DomainName, taskList, workerOptions)
	// err := cadenceWorker.Start()
	// if err != nil {
	// 	h.Logger.Error("Failed to start workers.", zap.Error(err))
	// 	panic("Failed to start workers")
	// }
}

func main() {
	fmt.Println("Starting Worker..")
	var appConfig config.AppConfig
	appConfig.Setup()
	var cadenceClient grpc.CadenceAdapter
	cadenceClient.SetupServiceConfig(&appConfig.Cadence)

	// workflow.Register(workflows.Workflow)
	cadenceClient.RegisterWorkflow(workflows.Workflow)
	startWorkers(&cadenceClient, workflows.TaskListName)
	// The workers are supposed to be long running process that should not exit.
	select {}
}
