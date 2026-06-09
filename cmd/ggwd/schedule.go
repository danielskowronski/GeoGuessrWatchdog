package main

import (
	"context"
	"errors"
	"time"

	"github.com/danielskowronski/geoguessrwatchdog/internal/config"
	"go.temporal.io/api/enums/v1"
	"go.temporal.io/api/serviceerror"
	"go.temporal.io/sdk/client"
)

func EnsureFetchDivisionsMapsSchedule(
	ctx context.Context,
	c client.Client,
	cfg config.WorkerConfig,
) error {
	scheduleID := TEMPORAL_SCHEDULE_FETCH_DIVISIONS_MAPS

	spec := client.ScheduleSpec{
		Intervals: []client.ScheduleIntervalSpec{
			{
				Every: time.Hour * time.Duration(cfg.Watchdogs.CompetitiveMaps.ScheduleFrequencyHours),
			},
		},
	}

	input := FetchDivisionsMapsWorkflowInput{
		TriggerApiUpdates:    true,
		TriggerMapUpdates:    true,
		TriggerNotifications: true,
		TemporalOptions:      cfg.Watchdogs.CompetitiveMaps.Temporal,
	}

	action := &client.ScheduleWorkflowAction{
		ID:        TEMPORAL_SCHEDULED_TASK_FETCH_DIVISIONS_MAPS,
		Workflow:  TEMPORAL_WORKFLOW_FETCH_DIVISIONS_MAPS,
		TaskQueue: cfg.Temporal.TaskQueue,
		Args:      []any{input},
	}

	handle := c.ScheduleClient().GetHandle(ctx, scheduleID)

	_, err := handle.Describe(ctx)
	if err != nil {
		var notFound *serviceerror.NotFound
		if errors.As(err, &notFound) {
			_, err = c.ScheduleClient().Create(ctx, client.ScheduleOptions{
				ID:      scheduleID,
				Spec:    spec,
				Action:  action,
				Overlap: enums.SCHEDULE_OVERLAP_POLICY_SKIP,
			})
			return err
		}

		return err
	}

	err = handle.Update(ctx, client.ScheduleUpdateOptions{
		DoUpdate: func(input client.ScheduleUpdateInput) (*client.ScheduleUpdate, error) {
			s := input.Description.Schedule

			s.Spec = &spec
			s.Action = action

			return &client.ScheduleUpdate{
				Schedule: &s,
			}, nil
		},
	})
	if err != nil {
		return err
	}

	if !cfg.Watchdogs.CompetitiveMaps.Enabled {
		err = handle.Pause(ctx, client.SchedulePauseOptions{
			Note: "disabled by config",
		})
		if err != nil {
			return err
		}
	} else {
		err = handle.Unpause(ctx, client.ScheduleUnpauseOptions{
			Note: "enabled by config",
		})
		if err != nil {
			return err
		}
	}

	return nil
}

func EnsureFetchUserStatsAndProgressSchedule(
	ctx context.Context,
	c client.Client,
	cfg config.WorkerConfig,
) error {
	scheduleID := TEMPORAL_SCHEDULE_FETCH_USER_STATS_AND_PROGRESS

	spec := client.ScheduleSpec{

		Calendars: []client.ScheduleCalendarSpec{
			{
				Hour: []client.ScheduleRange{
					{Start: cfg.Watchdogs.UserStats.ScheduleDaily.Hour},
				},
				Minute: []client.ScheduleRange{
					{Start: cfg.Watchdogs.UserStats.ScheduleDaily.Minute},
				},
			},
		},
	}

	input := FetchMuiltipleUsersStatsAndProgressWorkflowInput{
		TriggerUserStats:             true,
		TriggerUserProgress:          true,
		TriggerUserSingleplayerStats: true,
		UsersList:                    cfg.Watchdogs.UserStats.ObserveUsers,
		TemporalOptions:              cfg.Watchdogs.UserStats.Temporal,
	}

	action := &client.ScheduleWorkflowAction{
		ID:        TEMPORAL_SCHEDULED_TASK_FETCH_USER_STATS_AND_PROGRESS,
		Workflow:  TEMPORAL_WORKFLOW_FETCH_USER_STATS_AND_PROGRESS,
		TaskQueue: cfg.Temporal.TaskQueue,
		Args:      []any{input},
	}

	handle := c.ScheduleClient().GetHandle(ctx, scheduleID)

	_, err := handle.Describe(ctx)
	if err != nil {
		var notFound *serviceerror.NotFound
		if errors.As(err, &notFound) {
			_, err = c.ScheduleClient().Create(ctx, client.ScheduleOptions{
				ID:      scheduleID,
				Spec:    spec,
				Action:  action,
				Overlap: enums.SCHEDULE_OVERLAP_POLICY_SKIP,
			})
			return err
		}

		return err
	}

	err = handle.Update(ctx, client.ScheduleUpdateOptions{
		DoUpdate: func(input client.ScheduleUpdateInput) (*client.ScheduleUpdate, error) {
			s := input.Description.Schedule

			s.Spec = &spec
			s.Action = action

			return &client.ScheduleUpdate{
				Schedule: &s,
			}, nil
		},
	})
	if err != nil {
		return err
	}

	if !cfg.Watchdogs.UserStats.Enabled {
		err = handle.Pause(ctx, client.SchedulePauseOptions{
			Note: "disabled by config",
		})
		if err != nil {
			return err
		}
	} else {
		err = handle.Unpause(ctx, client.ScheduleUnpauseOptions{
			Note: "enabled by config",
		})
		if err != nil {
			return err
		}
	}

	return nil
}
