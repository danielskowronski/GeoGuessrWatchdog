package main

const (
	GG_API_PATH_DIVISIONS_LIST = "/v4/ranked-system/divisions"
	GG_API_PATH_MAP_INFO       = "/maps/%s"
	GG_API_PATH_USER_STATS     = "/v4/stats/users/%s"
	GG_API_PATH_USER_PROGRESS  = "/v4/ranked-system/progress/%s"
	GG_COOKIE_NAME             = "_ncfa"

	// TODO: split API names to workflows and verbs (scheduled/schedule/manual), also handle manual- with consts
	TEMPORAL_SCHEDULE_FETCH_DIVISIONS_MAPS                = "fetch-divisions-maps-schedule"
	TEMPORAL_SCHEDULED_TASK_FETCH_DIVISIONS_MAPS          = "fetch-divisions-maps-scheduled"
	TEMPORAL_WORKFLOW_FETCH_DIVISIONS_MAPS                = "FetchDivisionsMapsWorkflow"
	TEMPORAL_SCHEDULE_FETCH_USER_STATS_AND_PROGRESS       = "fetch-user-stats-and-progress-schedule"
	TEMPORAL_SCHEDULED_TASK_FETCH_USER_STATS_AND_PROGRESS = "fetch-user-stats-and-progress-scheduled"
	TEMPORAL_WORKFLOW_FETCH_USER_STATS_AND_PROGRESS       = "FetchMuiltipleUsersStatsAndProgressWorkflow"

	CONF_PATH_ENV_VAR = "GGWD_CONFIG_PATH"
	DEFAULT_CONF_PATH = "/etc/ggwd/config.yaml"
)
