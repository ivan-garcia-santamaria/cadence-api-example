package activities

import (
	"context"

	"go.uber.org/cadence/activity"
)

/**
 * This is the hello world workflow sample.
 */

// ApplicationName is the task list for this sample
const TaskListName = "helloWorldActivitiesGroup"

type HelloworldActivity struct {
}

func (h *HelloworldActivity) Hello(ctx context.Context, name, color string) (string, error) {
	logger := activity.GetLogger(ctx)
	logger.Info("helloworld activity started " + name + " " + color)
	return "Hello " + name + "! How old are you!", nil
}

func (h *HelloworldActivity) Hello2(ctx context.Context, name string) (string, error) {
	logger := activity.GetLogger(ctx)
	logger.Info("helloworld activity 2 started " + name)
	return "Hello2 " + name + "! How old are you!", nil
}
