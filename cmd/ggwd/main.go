package main

import (
	"context"
	"fmt"
	"os"
	"time"

	"go.temporal.io/sdk/client"
	"go.temporal.io/sdk/worker"

	"github.com/danielskowronski/geoguessrwatchdog/internal/buildinfo"
	"github.com/danielskowronski/geoguessrwatchdog/internal/config"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/spf13/cobra"
)

func mustLoad() config.WorkerConfig {
	cfg, err := config.LoadConfig[config.WorkerConfig](DEFAULT_CONF_PATH, config.WorkerConfigDefaults())
	if err != nil {
		panic(fmt.Sprintf("failed to load config: %v", err))
	}
	return cfg
}

func main() {
	fmt.Println(buildinfo.GetBuildInfo())

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
		fmt.Println("checking public IP...")
		api, err := NewAPIClient(cfg.GeoguessrAPI.Proxy, cfg.GeoguessrAPI.BaseURL, cfg.GeoguessrAPI.Cookie, cfg.GeoguessrAPI.Cache,
			map[string]string{
				cfg.Preflight.IpInfoCheckURL: "public_ip.json",
			})
		if err != nil {
			panic(fmt.Sprintf("failed to create API client: %v", err))
		}
		ip, err := api.GetPublicIP(context.Background(), cfg.Preflight.IpInfoCheckURL)
		if err != nil {
			panic(fmt.Sprintf("failed to get public IP: %v", err))
		}
		fmt.Printf("public IP: %s\n", ip)
	}

	temporalClient := mustTemporalClient(cfg.Temporal)
	defer temporalClient.Close()

	healthCtx, healthCancel := context.WithCancel(context.Background())
	defer healthCancel()

	healthState := &HealthState{}

	dbpool, err := pgxpool.New(context.Background(), cfg.Database.URL)
	if err != nil {
		panic(fmt.Sprintf("failed to create healthcheck db pool: %v", err))
	}
	defer dbpool.Close()

	startHealthServer(healthCtx, cfg.HealthCheck.HealthBind, dbpool, temporalClient, healthState)

	fmt.Println("ensuring schedules...")
	err = EnsureFetchDivisionsMapsSchedule(context.Background(), temporalClient, cfg)
	if err != nil {
		panic(fmt.Sprintf("failed to ensure schedules: %v", err))
	}
	err = EnsureFetchUserStatsAndProgressSchedule(context.Background(), temporalClient, cfg)
	if err != nil {
		panic(fmt.Sprintf("failed to ensure schedules: %v", err))
	}

	temporalWorker := worker.New(temporalClient, cfg.Temporal.TaskQueue, worker.Options{})
	acts := &Activities{
		Config: cfg,
	}
	temporalWorker.RegisterWorkflow(FetchDivisionsMapsWorkflow)
	temporalWorker.RegisterWorkflow(FetchMuiltipleUsersStatsAndProgressWorkflow)
	temporalWorker.RegisterWorkflow(FetchSingleUserStatsAndProgressWorkflow)
	temporalWorker.RegisterActivity(acts)
	fmt.Println("starting Temporal worker...")

	healthState.started.Store(true)

	if err := temporalWorker.Run(worker.InterruptCh()); err != nil {
		healthState.started.Store(false)
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
	cfg := mustLoad()
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
