package main

import (
	"context"
	"fmt"
	"os"
	"time"

	"go.temporal.io/sdk/client"
	"go.temporal.io/sdk/worker"

	"github.com/danielskowronski/geoguessrwatchdog/internal/config"
)

func main() {
	mode := getenv("GGWD_MODE", "worker")

	configPath := getenv(CONF_PATH_ENV_VAR, DEFAULT_CONF_PATH)
	cfg, err := config.Load(configPath)
	if err != nil {
		panic(fmt.Errorf("failed to load config: %w", err))
	}

	switch mode {
	case "worker":
		runWorker(cfg)
	case "trigger": // FIXME: this will need to select which workflow
		triggerWorkflow(cfg)
	case "schedule": // FIXME: this will need to configure schedules as per config; also this should be moved to worker init
		createSchedule(cfg)
	default:
		panic("unknown GGWD_MODE: " + mode)
	}
}

func runWorker(cfg *config.Config) {
	c := mustTemporalClient(cfg.Temporal)
	defer c.Close()

	w := worker.New(c, taskQueue(), worker.Options{})

	acts := &Activities{
		DatabaseURL:      cfg.Database.URL,
		HttpProxyURL:     cfg.GeoguessrAPI.Proxy,
		GeoGuessrApiBase: cfg.GeoguessrAPI.BaseURL,
		GeoGuessrCookie:  cfg.GeoguessrAPI.Cookie,
	}
	fmt.Println(cfg.GeoguessrAPI.Cookie)
	w.RegisterWorkflow(FetchFanoutWorkflow)
	w.RegisterActivity(acts)

	if err := w.Run(worker.InterruptCh()); err != nil {
		panic(err)
	}
}

func triggerWorkflow(cfg *config.Config) {
	c := mustTemporalClient(cfg.Temporal)
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

func createSchedule(cfg *config.Config) {
	c := mustTemporalClient(cfg.Temporal)
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

func mustTemporalClient(cfg config.TemporalConfig) client.Client {
	c, err := client.Dial(client.Options{
		HostPort:  cfg.Address,
		Namespace: cfg.Namespace,
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
