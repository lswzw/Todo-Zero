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
	"server/internal/svc"

	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/rest"
	"github.com/zeromicro/go-zero/rest/httpx"
)

//go:embed dist
var staticFiles embed.FS

//go:embed etc/todo-api.yaml
var defaultConfig []byte

var (
	host        = flag.String("host", "0.0.0.0", "listen host")
	port        = flag.Int("port", 8888, "listen port")
	dataDir     = flag.String("data-dir", "data", "data directory for SQLite database")
	dbFile      = flag.String("db-file", "todo.db", "SQLite database filename")
	jwtSecret   = flag.String("jwt-secret", "todo-app-jwt-secret-key-2024", "JWT signing secret")
	jwtExpire   = flag.Int64("jwt-expire", 86400, "JWT token expiration in seconds")
	configFile  = flag.String("f", "", "config file path (overrides command-line flags)")
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

	var c config.Config

	// If config file is specified, load it; otherwise build config from flags with defaults
	if *configFile != "" {
		conf.MustLoad(*configFile, &c)
	} else {
		// Try embedded default config first
		if err := conf.LoadFromYamlBytes(defaultConfig, &c); err != nil {
			fmt.Printf("[Config] Failed to load embedded config: %v\n", err)
			os.Exit(1)
		}
		// Override with command-line flags (flags always take precedence)
		c.Host = *host
		c.Port = *port
		c.Database.DataDir = *dataDir
		c.Database.DBFile = *dbFile
		c.Auth.AccessSecret = *jwtSecret
		c.Auth.AccessExpire = *jwtExpire
	}

	// Initialize SQLite database
	sqliteDB, err := db.InitDB(c.Database.DataDir, c.Database.DBFile)
	if err != nil {
		fmt.Printf("[DB] Failed to initialize database: %v\n", err)
		os.Exit(1)
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

	fmt.Printf("Starting server at %s:%d...\n", c.Host, c.Port)
	server.Start()
}

// staticFileHandler returns an http.Handler that serves embedded frontend files.
// This is used as the NotFoundHandler so that any non-API route serves the SPA.
// - Exact file matches are served directly
// - Other paths fall back to index.html for client-side routing
func staticFileHandler(content fs.FS) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
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
