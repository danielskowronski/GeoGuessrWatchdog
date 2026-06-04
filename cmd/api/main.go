package main

import (
	"context"
	"embed"
	"fmt"
	"html/template"
	"io/fs"
	"log/slog"
	"net/http"
	"os"
	"strings"

	"github.com/danielgtaylor/huma/v2"
	"github.com/danielgtaylor/huma/v2/adapters/humachi"
	"github.com/danielskowronski/geoguessrwatchdog/internal/config"

	"github.com/danielskowronski/geoguessrwatchdog/internal/buildinfo"
	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

const (
	API_VER           = "0.2.0"
	DEFAULT_CONF_PATH = "/etc/ggwd-api/config.yaml"
)

//go:embed web/html/*.html web/static/style/*.css web/static/*.js web/static/pages/*.js
var embeddedWebFS embed.FS

type App struct {
	db          *pgxpool.Pool
	userAliases config.UserAliasesConfig
	state       *HealthState
	logger      *slog.Logger
	metrics     *Metrics
}

type PageData struct {
	Title   string
	Heading string
	ID      string
}

func mustLoad() config.ApiConfig {
	cfg, err := config.LoadConfig[config.ApiConfig](DEFAULT_CONF_PATH, config.ApiConfigDefaults())
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

func getWebFS(serverCfg config.ApiServerConfig) fs.FS {
	if serverCfg.ServeLocally {
		localPath := serverCfg.LocalServePath
		if localPath == "" {
			localPath = "web"
		}

		if _, err := os.Stat(localPath); err != nil {
			slog.Error("failed to access local web path", "path", localPath, "err", err)
			panic("error during initialization")
		}

		slog.Info("serving web files from local path", "path", localPath)
		return os.DirFS(localPath)
	}

	webRoot, err := fs.Sub(embeddedWebFS, "web")
	if err != nil {
		slog.Error("failed to create sub filesystem for embedded web files", "err", err)
		panic("error during initialization")
	}

	slog.Info("serving web files from embedded filesystem")
	return webRoot
}

func renderPage(webRoot fs.FS, w http.ResponseWriter, page string, data PageData) {
	tpl, err := template.ParseFS(
		webRoot,
		"html/layout.html",
		"html/"+page+".html",
	)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	if err := tpl.ExecuteTemplate(w, "layout", data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func registerWebRoutes(r chi.Router, webRoot fs.FS) {
	staticFS, err := fs.Sub(webRoot, "static")
	if err != nil {
		slog.Error("failed to create sub filesystem for static files", "err", err)
		panic("error during initialization")
	}

	r.Handle("/static/*",
		http.StripPrefix("/static/",
			http.FileServer(http.FS(staticFS)),
		),
	)

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		renderPage(webRoot, w, "fetch_status", PageData{
			Title:   "GGWD Fetch Status",
			Heading: "Fetch Status",
		})
	})

	r.Get("/divisions", func(w http.ResponseWriter, r *http.Request) {
		renderPage(webRoot, w, "divisions", PageData{
			Title:   "GGWD Divisions",
			Heading: "Divisions",
		})
	})

	r.Get("/users", func(w http.ResponseWriter, r *http.Request) {
		renderPage(webRoot, w, "users", PageData{
			Title:   "GGWD Users",
			Heading: "Users",
		})
	})

	r.Get("/map/{id}", func(w http.ResponseWriter, r *http.Request) {
		id := chi.URLParam(r, "id")

		renderPage(webRoot, w, "map", PageData{
			Title:   "GGWD Map",
			Heading: "Map: " + id,
			ID:      id,
		})
	})

	r.Get("/user/{id}", func(w http.ResponseWriter, r *http.Request) {
		id := chi.URLParam(r, "id")

		renderPage(webRoot, w, "user", PageData{
			Title:   "GGWD User",
			Heading: "User: " + id,
			ID:      id,
		})
	})
}

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	}))
	slog.SetDefault(logger)
	cfg := mustLoad()

	dbPool, err := pgxpool.New(context.Background(), cfg.Database.URL)
	if err != nil {
		logger.Error("failed to connect to database", "err", err)
		os.Exit(1)
	}
	defer dbPool.Close()
	app := &App{
		logger:      logger.With("component", "api"),
		userAliases: cfg.UserAliases,
		db:          dbPool,
		state:       &HealthState{},
	}
	app.metrics = NewMetrics(app)
	r := chi.NewRouter()
	webRoot := getWebFS(cfg.Server)
	registerWebRoutes(r, webRoot)
	r.Handle("/metrics", app.metrics.Handler())
	api := humachi.New(r, huma.DefaultConfig("GGWD API", API_VER))
	huma.Get(api, "/api/divisions", app.GetDivisions)
	huma.Get(api, "/api/users", app.GetUsers)
	huma.Get(api, "/api/user/{id}", app.GetUserStats)
	huma.Get(api, "/api/map/{id}", app.GetMapHistory)
	huma.Get(api, "/api/fetch_statuses", app.GetFetchStatuses)
	huma.Get(api, "/version", app.Version)
	huma.Get(api, "/livez", app.Livez)
	huma.Get(api, "/readyz", app.Readyz)
	logger.Info("build info", "version", buildinfo.Version, "date", buildinfo.BuildDate)
	logger.Info("starting server",
		"bind", cfg.Server.Bind,
		"docs", linkToDocs(cfg.Server.Bind),
	)
	app.state.started.Store(true)
	if err := http.ListenAndServe(cfg.Server.Bind, r); err != nil {
		logger.Error("server stopped", "err", err)
		os.Exit(1)
	}
}
