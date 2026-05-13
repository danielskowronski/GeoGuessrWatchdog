package main

import (
	"time"

	"go.temporal.io/sdk/temporal"
	"go.temporal.io/sdk/workflow"
)

// TODO: DRY!

type WorkflowEnum int

const (
	WorkflowEnumFetchDivisionsMaps WorkflowEnum = iota
	WorkflowEnumFetchUserStatsAndProgress
)

func FetchDivisionsMapsWorkflow(ctx workflow.Context, input FetchDivisionsMapsWorkflowInput) error {

	activityOptions := workflow.ActivityOptions{
		StartToCloseTimeout: time.Duration(input.TemporalOptions.FanoutActivity.StartToCloseTimeoutSeconds) * time.Second,
		RetryPolicy: &temporal.RetryPolicy{
			InitialInterval:    time.Duration(input.TemporalOptions.FanoutActivity.RetryInitialIntervalSeconds) * time.Second,
			BackoffCoefficient: input.TemporalOptions.FanoutActivity.RetryBackoffCoefficient,
			MaximumInterval:    time.Duration(input.TemporalOptions.FanoutActivity.RetryMaximumIntervalSeconds) * time.Second,
			MaximumAttempts:    int32(input.TemporalOptions.FanoutActivity.RetryMaximumAttempts),
		},
	}

	ctx = workflow.WithActivityOptions(ctx, activityOptions)

	if input.TriggerApiUpdates {
		var divisionFetchResult DivisionFetchResult
		if err := workflow.ExecuteActivity(ctx, (*Activities).DivisionFetchActivity, input).Get(ctx, &divisionFetchResult); err != nil {
			return err
		}

		mapFetchResults := make([]MapFetchResult, 0, len(divisionFetchResult.UniqueMaps))
		if input.TriggerMapUpdates {
			for start := 0; start < len(divisionFetchResult.UniqueMaps); start += input.TemporalOptions.Parallelism {
				end := start + input.TemporalOptions.Parallelism
				if end > len(divisionFetchResult.UniqueMaps) {
					end = len(divisionFetchResult.UniqueMaps)
				}
				batch := divisionFetchResult.UniqueMaps[start:end]
				futures := make([]workflow.Future, 0, len(batch))
				for _, child := range batch {
					f := workflow.ExecuteActivity(ctx, (*Activities).MapFetchActivity, MapFetchInput{
						MapID: child,
					})
					futures = append(futures, f)
				}
				for _, f := range futures {
					var result MapFetchResult
					if err := f.Get(ctx, &result); err != nil {
						return err
					}
					mapFetchResults = append(mapFetchResults, result)
				}
			}
		}
	}

	if input.TriggerNotifications {
		// TODO: analyze results - we may skip notifier queries to db
		if err := workflow.ExecuteActivity(ctx, (*Activities).NotifyAboutMapChangeActivity, NotifierInput{}).Get(ctx, nil); err != nil {
			return err
		}
	}

	return nil
}

func FetchUserStatsAndProgressWorkflow(ctx workflow.Context, input FetchUserStatsAndProgressWorkflowInput) error {

	activityOptions := workflow.ActivityOptions{
		StartToCloseTimeout: time.Duration(input.TemporalOptions.FanoutActivity.StartToCloseTimeoutSeconds) * time.Second,
		RetryPolicy: &temporal.RetryPolicy{
			InitialInterval:    time.Duration(input.TemporalOptions.FanoutActivity.RetryInitialIntervalSeconds) * time.Second,
			BackoffCoefficient: input.TemporalOptions.FanoutActivity.RetryBackoffCoefficient,
			MaximumInterval:    time.Duration(input.TemporalOptions.FanoutActivity.RetryMaximumIntervalSeconds) * time.Second,
			MaximumAttempts:    int32(input.TemporalOptions.FanoutActivity.RetryMaximumAttempts),
		},
	}

	ctx = workflow.WithActivityOptions(ctx, activityOptions)

	if input.TriggerUserStats {
		for start := 0; start < len(input.UsersList); start += input.TemporalOptions.Parallelism {
			end := start + input.TemporalOptions.Parallelism
			if end > len(input.UsersList) {
				end = len(input.UsersList)
			}
			batch := input.UsersList[start:end]
			futures := make([]workflow.Future, 0, len(batch))
			for _, userID := range batch {
				f := workflow.ExecuteActivity(ctx, (*Activities).UserStatsFetchActivity, UserStatsFetchInput{
					UserID: userID,
				})
				futures = append(futures, f)
			}
			for _, f := range futures {
				if err := f.Get(ctx, nil); err != nil {
					return err
				}
			}
		}
	}
	if input.TriggerUserProgress {
		for start := 0; start < len(input.UsersList); start += input.TemporalOptions.Parallelism {
			end := start + input.TemporalOptions.Parallelism
			if end > len(input.UsersList) {
				end = len(input.UsersList)
			}
			batch := input.UsersList[start:end]
			futures := make([]workflow.Future, 0, len(batch))
			for _, userID := range batch {
				f := workflow.ExecuteActivity(ctx, (*Activities).UserProgressFetchActivity, UserProgressFetchInput{
					UserID: userID,
				})
				futures = append(futures, f)
			}
			for _, f := range futures {
				if err := f.Get(ctx, nil); err != nil {
					return err
				}
			}
		}
	}
	return nil
}
