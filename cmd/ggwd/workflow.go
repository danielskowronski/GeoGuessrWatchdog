package main

import (
	"time"

	"go.temporal.io/sdk/temporal"
	"go.temporal.io/sdk/workflow"
)

func FetchFanoutWorkflow(ctx workflow.Context, input WorkflowInput) error {
	if input.MaxParallel <= 0 {
		input.MaxParallel = 2
	}

	activityOptions := workflow.ActivityOptions{ // TODO: parametrize these timeouts and retry policies
		StartToCloseTimeout: 1 * time.Minute,
		RetryPolicy: &temporal.RetryPolicy{
			InitialInterval:    5 * time.Second,
			BackoffCoefficient: 2,
			MaximumInterval:    1 * time.Minute,
			MaximumAttempts:    3,
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
			for start := 0; start < len(divisionFetchResult.UniqueMaps); start += input.MaxParallel {
				end := start + input.MaxParallel
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
