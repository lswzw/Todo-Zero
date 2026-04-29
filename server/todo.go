package main

import (
	"embed"
	"flag"
	"fmt"
	"io/fs"
	"net/http"
	"os"
	"strings"

	"server/internal/config"
	"server/internal/db"
	"server/internal/handler"
	"server/internal/pkg/xerr"
	"server/internal/scheduler"
	"server/internal/svc"

	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/rest"
	"github.com/zeromicro/go-zero/rest/httpx"
)

// Version is set via -ldflags at build time.
var Version = "dev"

//go:embed dist
var staticFiles embed.FS

//go:embed etc/todo-api.yaml
var defaultConfig []byte

//go:embed docs/openapi.json
var openapiJSON []byte

var (
	host        = flag.String("host", "0.0.0.0", "listen host")
	port        = flag.Int("port", 8888, "listen port")
	dataDir     = flag.String("data-dir", "data", "data directory for SQLite database")
	dbFile      = flag.String("db-file", "todo.db", "SQLite database filename")
	jwtSecret   = flag.String("jwt-secret", "todo-app-jwt-secret-key-2024", "JWT signing secret")
	jwtExpire   = flag.Int64("jwt-expire", 86400, "JWT token expiration in seconds")
	debugMode   = flag.Bool("debug", false, "enable debug mode (API docs at /api-docs)")
	configFile  = flag.String("f", "", "config file path (overrides command-line flags)")
	showVersion = flag.Bool("version", false, "print version and exit")
)

func main() {
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Todo App - A standalone todo management application\n\n")
		fmt.Fprintf(os.Stderr, "Usage: %s [options]\n\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "Options:\n")
		flag.PrintDefaults()
		fmt.Fprintf(os.Stderr, "\nExamples:\n")
		fmt.Fprintf(os.Stderr, "  %s                          # Run with defaults\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "  %s -port 9090               # Run on port 9090\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "  %s -data-dir /var/todo      # Use custom data directory\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "  %s -f /etc/todo-api.yaml    # Use config file\n", os.Args[0])
	}
	flag.Parse()

	if *showVersion {
		fmt.Printf("Todo-Zero %s\n", Version)
		os.Exit(0)
	}

	fmt.Printf("[Main] Todo-Zero %s starting...\n", Version)

	var c config.Config

	// Load configuration: config file → embedded default, then flags override
	if *configFile != "" {
		conf.MustLoad(*configFile, &c)
	} else {
		if err := conf.LoadFromYamlBytes(defaultConfig, &c); err != nil {
			fmt.Printf("[Config] Failed to load embedded config: %v\n", err)
			os.Exit(1)
		}
	}

	// Command-line flags always override config file / defaults (only non-zero values)
	if *host != "0.0.0.0" {
		c.Host = *host
	}
	if *port != 8888 {
		c.Port = *port
	}
	if *dataDir != "data" {
		c.Database.DataDir = *dataDir
	}
	if *dbFile != "todo.db" {
		c.Database.DBFile = *dbFile
	}
	if *jwtSecret != "todo-app-jwt-secret-key-2024" {
		c.Auth.AccessSecret = *jwtSecret
	}
	if *jwtExpire != 86400 {
		c.Auth.AccessExpire = *jwtExpire
	}
	if *debugMode {
		c.Debug = true
	}

	// Initialize SQLite database
	sqliteDB, err := db.InitDB(c.Database.DataDir, c.Database.DBFile)
	if err != nil {
		fmt.Printf("[DB] Failed to initialize database: %v\n", err)
		os.Exit(1)
	}
	defer sqliteDB.Close()

	// JWT Secret security: auto-generate and persist if using default or empty
	const defaultJWTSecret = "todo-app-jwt-secret-key-2024"
	if c.Auth.AccessSecret == defaultJWTSecret || c.Auth.AccessSecret == "" {
		secret, generated, err := db.GetOrCreateJWTSecret(sqliteDB)
		if err != nil {
			fmt.Printf("[Security] WARNING: Failed to manage JWT secret: %v\n", err)
		} else {
			c.Auth.AccessSecret = secret
			if generated {
				fmt.Println("[Security] JWT secret auto-generated and persisted to database")
				fmt.Println("[Security] All existing tokens are now invalid, please re-login")
			} else {
				fmt.Println("[Security] JWT secret loaded from database")
			}
		}
	} else {
		fmt.Println("[Security] Using JWT secret from configuration")
	}
	// Create sub filesystem for embedded frontend (strip "dist/" prefix)
	distFS, err := fs.Sub(staticFiles, "dist")
	if err != nil {
		fmt.Printf("[Static] Failed to create sub filesystem: %v\n", err)
		os.Exit(1)
	}

	server := rest.MustNewServer(c.RestConf, rest.WithNotFoundHandler(
		staticFileHandler(distFS),
	))
	defer server.Stop()

	// Register custom error handler
	httpx.SetErrorHandlerCtx(xerr.ErrorResponse)

	ctx := svc.NewServiceContext(c, sqliteDB)

	fmt.Println("[Static] Frontend embedded and serving from /")

	handler.RegisterHandlers(server, ctx)

	// API documentation routes — only in debug mode and only when running in development environment
	if c.Debug {
		// Security: Only allow debug mode in development environment
		if os.Getenv("GO_ENV") == "production" {
			fmt.Println("[Security] WARNING: Debug mode disabled in production environment")
		} else {
			registerAPIDocRoutes(server)
			fmt.Println("[Debug] API documentation enabled: /api-docs, /openapi.json")
		}
	}

	scheduler.StartCleanupScheduler(ctx)
	scheduler.StartBackupScheduler(ctx)

	fmt.Printf("Starting server at %s:%d...\n", c.Host, c.Port)
	server.Start()
}

// staticFileHandler returns an http.Handler that serves embedded frontend files.
// This is used as the NotFoundHandler so that any non-API route serves the SPA.
// - Exact file matches are served directly
// - Other paths fall back to index.html for client-side routing
func staticFileHandler(content fs.FS) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Security headers for all responses
		w.Header().Set("X-Content-Type-Options", "nosniff")
		w.Header().Set("X-Frame-Options", "DENY")
		w.Header().Set("X-XSS-Protection", "1; mode=block")
		w.Header().Set("Referrer-Policy", "strict-origin-when-cross-origin")

		path := r.URL.Path

		// Root path or SPA route → serve index.html
		embedPath := strings.TrimPrefix(path, "/")

		if embedPath == "" {
			// Root path "/"
			data, err := fs.ReadFile(content, "index.html")
			if err != nil {
				http.NotFound(w, r)
				return
			}
			w.Header().Set("Content-Type", "text/html; charset=utf-8")
			w.Write(data)
			return
		}

		// Try to serve the exact file from embed FS
		data, err := fs.ReadFile(content, embedPath)
		if err == nil {
			contentType := getContentType(embedPath)
			w.Header().Set("Content-Type", contentType)
			if strings.HasPrefix(embedPath, "assets/") {
				w.Header().Set("Cache-Control", "public, max-age=31536000, immutable")
			}
			w.Write(data)
			return
		}

		// SPA fallback: serve index.html for client-side routing
		idxData, err := fs.ReadFile(content, "index.html")
		if err != nil {
			http.NotFound(w, r)
			return
		}
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		w.Write(idxData)
	})
}

// getContentType returns the MIME type based on file extension.
func getContentType(filepath string) string {
	switch {
	case strings.HasSuffix(filepath, ".html"):
		return "text/html; charset=utf-8"
	case strings.HasSuffix(filepath, ".js"):
		return "application/javascript"
	case strings.HasSuffix(filepath, ".css"):
		return "text/css"
	case strings.HasSuffix(filepath, ".json"):
		return "application/json"
	case strings.HasSuffix(filepath, ".png"):
		return "image/png"
	case strings.HasSuffix(filepath, ".svg"):
		return "image/svg+xml"
	case strings.HasSuffix(filepath, ".ico"):
		return "image/x-icon"
	case strings.HasSuffix(filepath, ".woff"):
		return "font/woff"
	case strings.HasSuffix(filepath, ".woff2"):
		return "font/woff2"
	case strings.HasSuffix(filepath, ".ttf"):
		return "font/ttf"
	case strings.HasSuffix(filepath, ".eot"):
		return "application/vnd.ms-fontobject"
	default:
		return "application/octet-stream"
	}
}

// registerAPIDocRoutes adds OpenAPI documentation endpoints.
func registerAPIDocRoutes(srv *rest.Server) {
	// Serve OpenAPI JSON spec
	srv.AddRoutes([]rest.Route{
		{
			Method:  http.MethodGet,
			Path:    "/openapi.json",
			Handler: openapiJSONHandler,
		},
		{
			Method:  http.MethodGet,
			Path:    "/api-docs",
			Handler: scalarDocHandler,
		},
	})
}

func openapiJSONHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Header().Set("Cache-Control", "no-cache")
	w.Write(openapiJSON)
}

const scalarHTML = `<!DOCTYPE html>
<html>
<head>
  <title>Todo-Zero API Docs</title>
  <meta charset="utf-8" />
  <meta name="viewport" content="width=device-width, initial-scale=1" />
</head>
<body>
  <script id="api-reference" data-url="/openapi.json"></script>
  <script src="https://cdn.jsdelivr.net/npm/@scalar/api-reference"></script>
</body>
</html>`

func scalarDocHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Header().Set("Cache-Control", "no-cache")
	w.Write([]byte(scalarHTML))
}
