package main

import (
	"fmt"

	"github.com/ivan-garcia-santamaria/cadence-api-example/app/commons"
	"github.com/ivan-garcia-santamaria/cadence-api-example/app/config"
	"github.com/ivan-garcia-santamaria/cadence-api-example/app/worker-activity/activities"
	"go.uber.org/cadence/activity"
	"go.uber.org/cadence/worker"
	"go.uber.org/zap"
)

func startWorkers(h *commons.CadenceAdapter, taskList string) {
	// Configure worker options.
	workerOptions := worker.Options{
		MetricsScope:          h.Scope,
		Logger:                h.Logger,
		DisableWorkflowWorker: true,
	}

	cadenceWorker := worker.New(h.ServiceClient, h.Config.Domain, taskList, workerOptions)
	err := cadenceWorker.Start()
	if err != nil {
		h.Logger.Error("Failed to start workers.", zap.Error(err))
		panic("Failed to start workers")
	}
}

func main() {
	fmt.Println("Starting Worker..")
	var appConfig config.AppConfig
	appConfig.Setup()
	var cadenceClient commons.CadenceAdapter
	cadenceClient.Setup(&appConfig.Cadence)

	helloworldActivity := &activities.HelloworldActivity{}
	activity.Register(helloworldActivity.Hello)
	activity.Register(helloworldActivity.Hello2)

	startWorkers(&cadenceClient, activities.TaskListName)
	// The workers are supposed to be long running process that should not exit.
	select {}
}
