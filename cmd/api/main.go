package main

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/danielgtaylor/huma/v2"
	"github.com/danielgtaylor/huma/v2/adapters/humachi"
	"github.com/danielskowronski/geoguessrwatchdog/internal/config"

	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

const (
	API_VER           = "0.2.0"
	DEFAULT_CONF_PATH = "/etc/ggwd-api/config.yaml"
)

type App struct {
	db          *pgxpool.Pool
	userAliases config.UserAliasesConfig
}

type HealthOutput struct {
	Body struct {
		Status string `json:"status"`
	}
}

func mustLoad() config.ApiConfig {
	cfg, err := config.LoadConfig[config.ApiConfig](DEFAULT_CONF_PATH, config.ApiConfigDefaults())
	if err != nil {
		panic(fmt.Sprintf("failed to load config: %v", err))
	}
	return cfg
}

func linkToDocs(bind string) string {
	elements := strings.Split(bind, ":")
	if len(elements) == 2 {
		host := elements[0]
		if host == "" {
			host = "localhost"
		}
		return fmt.Sprintf("http://%s:%s/docs", host, elements[1])
	}

	return fmt.Sprintf("http://localhost%s/docs", bind)
}

func main() {
	cfg := mustLoad()

	dbPool, err := pgxpool.New(context.Background(), cfg.Database.URL)
	if err != nil {
		panic(fmt.Sprintf("failed to connect to database: %v", err))
	}
	defer dbPool.Close()
	app := &App{
		userAliases: cfg.UserAliases,
		db:          dbPool,
	}

	r := chi.NewRouter()
	api := humachi.New(r, huma.DefaultConfig("GGWD API", API_VER))

	huma.Get(api, "/health", func(ctx context.Context, input *struct{}) (*HealthOutput, error) {
		resp := &HealthOutput{}
		resp.Body.Status = "ok"
		return resp, nil
	})

	huma.Get(api, "/divisions", app.GetDivisions)

	// TODO: /maps
	huma.Get(api, "/map/{id}", app.GetMapHistory)

	huma.Get(api, "/users", app.GetUsers)
	huma.Get(api, "/user/{id}", app.GetUserStats)

	fmt.Printf("Starting server at %s\n", cfg.Server.Bind)
	fmt.Printf("API documentation available at %s\n", linkToDocs(cfg.Server.Bind))
	http.ListenAndServe(cfg.Server.Bind, r)
}
