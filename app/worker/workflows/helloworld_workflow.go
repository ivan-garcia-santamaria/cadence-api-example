package workflows

import (
	"fmt"
	"time"

	"github.com/ivan-garcia-santamaria/cadence-api-example/app/worker-activity/activities"
	"go.uber.org/cadence/workflow"
	"go.uber.org/zap"
)

/**
 * This is the hello world workflow sample.
 */

// ApplicationName is the task list for this sample
const TaskListName = "helloWorldGroup"
const SignalName = "helloWorldSignal"

var activityOptions = workflow.ActivityOptions{
	ScheduleToStartTimeout: time.Minute,
	StartToCloseTimeout:    time.Minute,
	HeartbeatTimeout:       time.Second * 20,
	TaskList:               activities.TaskListName,
}

func Workflow(ctx workflow.Context, name, color string) (string, error) {
	ctx = workflow.WithActivityOptions(ctx, activityOptions)
	logger := workflow.GetLogger(ctx)
	logger.Info("helloworld workflow started " + name)

	var activityResult string

	var helloworldActivity *activities.HelloworldActivity
	err := workflow.ExecuteActivity(ctx, helloworldActivity.Hello, name, color).Get(ctx, &activityResult)

	if err != nil {
		logger.Error("Activity failed.", zap.Error(err))
		return "", err
	}

	// After saying hello, the workflow will wait for you to inform it of your age!
	signalName := SignalName
	selector := workflow.NewSelector(ctx)
	var ageResult int

	for {
		signalChan := workflow.GetSignalChannel(ctx, signalName)

		selector.AddReceive(signalChan, func(c workflow.Channel, more bool) {
			c.Receive(ctx, &ageResult)
			workflow.GetLogger(ctx).Info("Received age results from signal!", zap.String("signal", signalName), zap.Int("value", ageResult))
		})
		workflow.GetLogger(ctx).Info("Waiting for signal on channel.. " + signalName)
		// Wait for signal
		selector.Select(ctx)

		// We can check the age and return an appropriate response
		if ageResult > 0 && ageResult < 150 {
			logger.Info("Workflow completed.", zap.String("NameResult", activityResult), zap.Int("AgeResult", ageResult))

			return fmt.Sprintf("Hello "+name+"! You are %v years old!", ageResult), nil
		} else {
			return "You can't be that old!", nil
		}
	}
}
