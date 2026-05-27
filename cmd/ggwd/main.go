package main

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"time"

	"go.temporal.io/sdk/client"
	"go.temporal.io/sdk/worker"

	"github.com/danielskowronski/geoguessrwatchdog/internal/config"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/spf13/cobra"
	"go.temporal.io/sdk/log"
)

func mustLoad() config.WorkerConfig {
	cfg, err := config.LoadConfig[config.WorkerConfig](DEFAULT_CONF_PATH, config.WorkerConfigDefaults())
	if err != nil {
		slog.Error("failed to load config", "err", err)
		panic("error during initialization")
	}

	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: cfg.Logging.GetLevel(),
	}))
	slog.SetDefault(logger)

	return cfg
}

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	}))
	slog.SetDefault(logger)

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
	cfg := mustLoad()

	if cfg.Preflight.IpInfoCheckEnabled {
		// FUTURE: this should be reported with exporters
		slog.Info("checking public IP", "check_url", cfg.Preflight.IpInfoCheckURL)
		api, err := NewAPIClient(cfg.GeoguessrAPI.Proxy, cfg.GeoguessrAPI.BaseURL, cfg.GeoguessrAPI.Cookie, cfg.GeoguessrAPI.Cache,
			map[string]string{
				cfg.Preflight.IpInfoCheckURL: "public_ip.json",
			})
		if err != nil {
			slog.Error("failed to create API client for public IP check", "err", err)
			panic("error during initialization")
		}
		ip, err := api.GetPublicIP(context.Background(), cfg.Preflight.IpInfoCheckURL)
		if err != nil {
			slog.Error("failed to get public IP", "err", err)
			panic("error during initialization")
		}

		slog.Info("public IP verified", "ip", ip)
	}

	temporalClient := mustTemporalClient(cfg.Temporal, slog.Default())
	defer temporalClient.Close()

	healthCtx, healthCancel := context.WithCancel(context.Background())
	defer healthCancel()

	healthState := &HealthState{}

	dbpool, err := pgxpool.New(context.Background(), cfg.Database.URL)
	if err != nil {
		slog.Error("failed to create healthcheck db pool", "err", err)
		panic("error during initialization")
	}
	defer dbpool.Close()

	startHealthServer(healthCtx, cfg.HealthCheck.HealthBind, dbpool, temporalClient, healthState)

	slog.Info("ensuring schedules are registered in Temporal")
	err = EnsureFetchDivisionsMapsSchedule(context.Background(), temporalClient, cfg)
	if err != nil {
		slog.Error("failed to ensure schedule", "schedule", "FetchDivisionsMaps", "err", err)
		panic("error during initialization")
	}
	err = EnsureFetchUserStatsAndProgressSchedule(context.Background(), temporalClient, cfg)
	if err != nil {
		slog.Error("failed to ensure schedule", "schedule", "FetchUserStatsAndProgress", "err", err)
		panic("error during initialization")
	}

	temporalWorker := worker.New(temporalClient, cfg.Temporal.TaskQueue, worker.Options{})
	acts := &Activities{
		Config: cfg,
	}
	temporalWorker.RegisterWorkflow(FetchDivisionsMapsWorkflow)
	temporalWorker.RegisterWorkflow(FetchMuiltipleUsersStatsAndProgressWorkflow)
	temporalWorker.RegisterWorkflow(FetchSingleUserStatsAndProgressWorkflow)
	temporalWorker.RegisterActivity(acts)
	slog.Info("starting Temporal worker with task queue", "task_queue", cfg.Temporal.TaskQueue)

	healthState.started.Store(true)

	if err := temporalWorker.Run(worker.InterruptCh()); err != nil {
		healthState.started.Store(false)
		slog.Error("Temporal worker encountered an error", "err", err)
		panic("error during worker execution")
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
	cfg := mustLoad()
	c := mustTemporalClient(cfg.Temporal, slog.Default())
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
		slog.Error("failed to start workflow", "workflow", Workflow, "err", runErr)
		panic("error during workflow execution")
	}

	slog.Info("started workflow", "workflow_id", run.GetID(), "run_id", run.GetRunID())
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
			slog.Info("triggering workflow", "workflow", "FetchDivisionsMaps")
			return startWorkflow(WorkflowEnumFetchDivisionsMaps)
		},
	})

	cmd.AddCommand(&cobra.Command{
		Use:   "FetchUserStatsAndProgress",
		Short: "Trigger FetchUserStatsAndProgress workflow",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			slog.Info("triggering workflow", "workflow", "FetchUserStatsAndProgress")
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
			slog.Info("triggering activity", "activity", "FetchDivision")
			return nil
		},
	})

	cmd.AddCommand(&cobra.Command{
		Use:   "FetchMap <map-id>",
		Short: "Trigger FetchMap activity",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			mapID := args[0]
			slog.Info("triggering activity", "activity", "FetchMap", "map_id", mapID)
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
			slog.Info("triggering activity", "activity", "NotifyMap", "map_id", mapID, "message", message)
			return nil
		},
	})

	cmd.AddCommand(&cobra.Command{
		Use:   "FetchUserStats <user-id>",
		Short: "Trigger FetchUserStats activity",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			userID := args[0]
			slog.Info("triggering activity", "activity", "FetchUserStats", "user_id", userID)
			return nil
		},
	})

	cmd.AddCommand(&cobra.Command{
		Use:   "FetchUserHistory <user-id>",
		Short: "Trigger FetchUserHistory activity",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			userID := args[0]
			slog.Info("triggering activity", "activity", "FetchUserHistory", "user_id", userID)
			return nil
		},
	})

	return cmd
}

func mustTemporalClient(cfg config.TemporalConfig, logger *slog.Logger) client.Client {
	temporalLogger := log.NewStructuredLogger(logger)
	c, err := client.Dial(client.Options{
		HostPort:  cfg.Address,
		Namespace: cfg.Namespace,
		Logger:    temporalLogger,
	})
	if err != nil {
		slog.Error("failed to create Temporal client", "err", err)
		panic("error during initialization")
	}

	return c
}
