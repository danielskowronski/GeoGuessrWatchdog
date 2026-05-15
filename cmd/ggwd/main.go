package main

import (
	"context"
	"fmt"
	"os"
	"time"

	"go.temporal.io/sdk/client"
	"go.temporal.io/sdk/worker"

	"github.com/danielskowronski/geoguessrwatchdog/internal/config"
	"github.com/spf13/cobra"
)

func getenv(name string, fallback string) string {
	value := os.Getenv(name)
	if value == "" {
		return fallback
	}
	return value
}

func loadConfig() (config.Config, error) {
	configPath := getenv(CONF_PATH_ENV_VAR, DEFAULT_CONF_PATH)
	cfg, err := config.Load(configPath)
	if err != nil {
		return config.Config{}, fmt.Errorf("failed to load config: %w", err)
	}
	return *cfg, nil
}

func main() {
	rootCmd := &cobra.Command{
		Use:   "ggwd",
		Short: "GeoGuessrWatchdog - Temporal worker and CLI for GeoGuessr ranked stats monitoring",
	}

	rootCmd.AddCommand(workerCmd())
	rootCmd.AddCommand(triggerWorkflowCmd())
	rootCmd.AddCommand(triggerActivityCmd())

	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func startWorker() error {
	cfg, err := loadConfig()
	if err != nil {
		panic(fmt.Sprintf("failed to load config: %v", err))
	}

	c := mustTemporalClient(cfg.Temporal)
	defer c.Close()

	fmt.Println("ensuring schedules...")
	err = EnsureFetchDivisionsMapsSchedule(context.Background(), c, cfg)
	if err != nil {
		panic(fmt.Sprintf("failed to ensure schedules: %v", err))
	}
	err = EnsureFetchUserStatsAndProgressSchedule(context.Background(), c, cfg)
	if err != nil {
		panic(fmt.Sprintf("failed to ensure schedules: %v", err))
	}

	w := worker.New(c, cfg.Temporal.TaskQueue, worker.Options{})
	acts := &Activities{
		Config: cfg,
	}
	w.RegisterWorkflow(FetchDivisionsMapsWorkflow)
	w.RegisterWorkflow(FetchMuiltipleUsersStatsAndProgressWorkflow)
	w.RegisterWorkflow(FetchSingleUserStatsAndProgressWorkflow)
	w.RegisterActivity(acts)
	fmt.Println("starting Temporal worker...")

	if err := w.Run(worker.InterruptCh()); err != nil {
		panic(err)
	}

	return nil
}

func workerCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "worker",
		Short: "Start Temporal worker",
		RunE: func(cmd *cobra.Command, args []string) error {
			return startWorker()
		},
	}
}

func startWorkflow(Workflow WorkflowEnum) error {
	cfg, err := loadConfig()
	if err != nil {
		panic(fmt.Sprintf("failed to load config: %v", err))
	}
	c := mustTemporalClient(cfg.Temporal)
	defer c.Close()

	workflowIdElement := "unknown"
	switch Workflow {
	case WorkflowEnumFetchDivisionsMaps:
		workflowIdElement = "FetchDivisionsMaps"
	case WorkflowEnumFetchUserStatsAndProgress:
		workflowIdElement = "FetchUserStatsAndProgress"
	}
	opts := client.StartWorkflowOptions{
		ID:        fmt.Sprintf("manual-%s-%d", workflowIdElement, time.Now().Unix()),
		TaskQueue: cfg.Temporal.TaskQueue,
	}
	var run client.WorkflowRun
	var runErr error
	switch Workflow {
	case WorkflowEnumFetchDivisionsMaps:
		input := FetchDivisionsMapsWorkflowInput{
			TriggerApiUpdates:    true,
			TriggerMapUpdates:    true,
			TriggerNotifications: true,
			TemporalOptions:      cfg.Watchdogs.CompetitiveMaps.Temporal,
		}
		run, runErr = c.ExecuteWorkflow(context.Background(), opts, FetchDivisionsMapsWorkflow, input)
	case WorkflowEnumFetchUserStatsAndProgress:
		input := FetchMuiltipleUsersStatsAndProgressWorkflowInput{
			TriggerUserStats:    true,
			TriggerUserProgress: true,
			UsersList:           cfg.Watchdogs.UserStats.ObserveUsers,
			TemporalOptions:     cfg.Watchdogs.UserStats.Temporal,
		}
		run, runErr = c.ExecuteWorkflow(context.Background(), opts, FetchMuiltipleUsersStatsAndProgressWorkflow, input)
	default:
		return fmt.Errorf("unknown workflow: %v", Workflow)
	}

	if runErr != nil {
		panic(runErr)
	}

	fmt.Printf("started workflow_id=%s run_id=%s\n", run.GetID(), run.GetRunID())
	return nil
}

func triggerWorkflowCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "trigger-workflow",
		Short: "Trigger Temporal workflows",
	}

	cmd.AddCommand(&cobra.Command{
		Use:   "FetchDivisionsMaps",
		Short: "Trigger FetchDivisionsMaps workflow",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			fmt.Println("trigger workflow: FetchDivisionsMaps")
			return startWorkflow(WorkflowEnumFetchDivisionsMaps)
		},
	})

	cmd.AddCommand(&cobra.Command{
		Use:   "FetchUserStatsAndProgress",
		Short: "Trigger FetchUserStatsAndProgress workflow",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			fmt.Println("trigger workflow: FetchUserStatsAndProgress")
			return startWorkflow(WorkflowEnumFetchUserStatsAndProgress)
		},
	})

	return cmd
}

// TODO: runActivity with params

func triggerActivityCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "trigger-activity",
		Short: "Trigger Temporal activities directly",
	}

	cmd.AddCommand(&cobra.Command{
		Use:   "FetchDivision",
		Short: "Trigger FetchDivision activity",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			fmt.Println("trigger activity: FetchDivision")
			return nil
		},
	})

	cmd.AddCommand(&cobra.Command{
		Use:   "FetchMap <map-id>",
		Short: "Trigger FetchMap activity",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			mapID := args[0]
			fmt.Println("trigger activity: FetchMap", mapID)
			return nil
		},
	})

	cmd.AddCommand(&cobra.Command{
		Use:   "NotifyMap <map-id> <message>",
		Short: "Trigger NotifyMap activity",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			mapID := args[0]
			message := args[1]
			fmt.Println("trigger activity: NotifyMap", mapID, message)
			return nil
		},
	})

	cmd.AddCommand(&cobra.Command{
		Use:   "FetchUserStats <user-id>",
		Short: "Trigger FetchUserStats activity",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			userID := args[0]
			fmt.Println("trigger activity: FetchUserStats", userID)
			return nil
		},
	})

	cmd.AddCommand(&cobra.Command{
		Use:   "FetchUserHistory <user-id>",
		Short: "Trigger FetchUserHistory activity",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			userID := args[0]
			fmt.Println("trigger activity: FetchUserHistory", userID)
			return nil
		},
	})

	return cmd
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
