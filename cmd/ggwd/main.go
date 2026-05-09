package main

import (
	"context"
	"fmt"
	"os"
	"time"

	"go.temporal.io/sdk/client"
	"go.temporal.io/sdk/worker"
)

func main() {
	mode := getenv("GGWD_MODE", "worker")

	switch mode {
	case "worker":
		runWorker()
	case "trigger":
		triggerWorkflow()
	case "schedule":
		createSchedule()
	default:
		panic("unknown GGWD_MODE: " + mode)
	}
}

func genActivities() *Activities {
	// TODO: maybe this should be also configurable via config file?
	return &Activities{
		DatabaseURL:      mustEnv("GGWD_DB_URL"),
		HttpProxyURL:     getenv("HTTP_PROXY_URL", ""),
		GeoGuessrApiBase: getenv("GGWD_GG_API_BASE", GG_API_DEFAULT_BASE),
		GeoGuessrCookie:  mustEnv("GGWD_GG_COOKIE"),
	}
}

func runWorker() {
	c := mustTemporalClient()
	defer c.Close()

	w := worker.New(c, taskQueue(), worker.Options{})

	acts := genActivities()
	w.RegisterWorkflow(FetchFanoutWorkflow)
	w.RegisterActivity(acts)

	if err := w.Run(worker.InterruptCh()); err != nil {
		panic(err)
	}
}

func triggerWorkflow() {
	c := mustTemporalClient()
	defer c.Close()

	input := WorkflowInput{
		TriggerApiUpdates:    true,
		TriggerMapUpdates:    true,
		TriggerNotifications: true,
	}

	opts := client.StartWorkflowOptions{
		ID:        fmt.Sprintf("manual-%d", time.Now().Unix()),
		TaskQueue: taskQueue(),
	}

	run, err := c.ExecuteWorkflow(context.Background(), opts, FetchFanoutWorkflow, input)
	if err != nil {
		panic(err)
	}

	fmt.Printf("started workflow_id=%s run_id=%s\n", run.GetID(), run.GetRunID())
}

func createSchedule() {
	c := mustTemporalClient()
	defer c.Close()

	scheduleID := getenv("SCHEDULE_ID", "api-fetch-every-6h")

	input := WorkflowInput{
		TriggerApiUpdates:    true,
		TriggerMapUpdates:    true,
		TriggerNotifications: true,
	}

	handle, err := c.ScheduleClient().Create(context.Background(), client.ScheduleOptions{
		ID: scheduleID,
		Spec: client.ScheduleSpec{
			Intervals: []client.ScheduleIntervalSpec{
				{Every: 6 * time.Hour},
			},
		},
		Action: &client.ScheduleWorkflowAction{
			ID:        "scheduled-api-fetch",
			Workflow:  FetchFanoutWorkflow,
			TaskQueue: taskQueue(),
			Args:      []any{input},
		},
	})
	if err != nil {
		panic(err)
	}

	fmt.Printf("created schedule_id=%s\n", handle.GetID())
}

func mustTemporalClient() client.Client {
	c, err := client.Dial(client.Options{
		HostPort:  mustEnv("TEMPORAL_ADDRESS"),
		Namespace: getenv("TEMPORAL_NAMESPACE", "default"),
	})
	if err != nil {
		panic(err)
	}

	return c
}

func taskQueue() string {
	return getenv("TEMPORAL_TASK_QUEUE", "ggwd-task-queue")
}
func mustEnv(name string) string {
	value := os.Getenv(name)
	if value == "" {
		panic("missing env var: " + name)
	}
	return value
}

func getenv(name string, fallback string) string {
	value := os.Getenv(name)
	if value == "" {
		return fallback
	}
	return value
}
