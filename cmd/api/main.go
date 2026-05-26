package main

import (
	"context"
	"embed"
	"fmt"
	"html/template"
	"io/fs"
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
}

type HealthOutput struct {
	Body struct {
		Status string `json:"status"`
	}
}

type PageData struct {
	Title   string
	Heading string
	ID      string
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

func getWebFS(serverCfg config.ApiServerConfig) fs.FS {
	if serverCfg.ServeLocally {
		localPath := serverCfg.LocalServePath
		if localPath == "" {
			localPath = "web"
		}

		if _, err := os.Stat(localPath); err != nil {
			panic(fmt.Sprintf("failed to access local web path %q: %v", localPath, err))
		}

		fmt.Printf("Serving web files from local path: %s\n", localPath)
		return os.DirFS(localPath)
	}

	webRoot, err := fs.Sub(embeddedWebFS, "web")
	if err != nil {
		panic(fmt.Sprintf("failed to create embedded web filesystem: %v", err))
	}

	fmt.Println("Serving web files from embedded filesystem")
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
		panic(fmt.Sprintf("failed to create static filesystem: %v", err))
	}

	r.Handle("/static/*",
		http.StripPrefix("/static/",
			http.FileServer(http.FS(staticFS)),
		),
	)

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/divisions", http.StatusFound)
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
	cfg := mustLoad()
	fmt.Println(cfg.Server.LocalServePath)

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

	webRoot := getWebFS(cfg.Server)
	registerWebRoutes(r, webRoot)

	api := humachi.New(r, huma.DefaultConfig("GGWD API", API_VER))

	huma.Get(api, "/health", func(ctx context.Context, input *struct{}) (*HealthOutput, error) {
		resp := &HealthOutput{}
		resp.Body.Status = "ok"
		return resp, nil
	})

	huma.Get(api, "/api/divisions", app.GetDivisions)
	huma.Get(api, "/api/users", app.GetUsers)
	huma.Get(api, "/api/user/{id}", app.GetUserStats)
	huma.Get(api, "/api/map/{id}", app.GetMapHistory)

	fmt.Println(buildinfo.GetBuildInfo())

	fmt.Printf("Starting server at %s\n", cfg.Server.Bind)
	fmt.Printf("API documentation available at %s\n", linkToDocs(cfg.Server.Bind))

	if err := http.ListenAndServe(cfg.Server.Bind, r); err != nil {
		panic(err)
	}
}
