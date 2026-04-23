package db

import (
	"crypto/rand"
	"database/sql"
	"embed"
	"encoding/hex"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	_ "modernc.org/sqlite"
)

//go:embed init.sql
var initSQL embed.FS

// InitDB initializes the SQLite database and returns the connection.
func InitDB(dataDir, dbFile string) (*sql.DB, error) {
	if err := os.MkdirAll(dataDir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create data directory: %w", err)
	}

	dbPath := filepath.Join(dataDir, dbFile)
	connStr := dbPath + "?_journal_mode=WAL&_busy_timeout=5000"

	sqliteDB, err := sql.Open("sqlite", connStr)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	// Configure connection pool for SQLite
	// SQLite only supports one writer at a time, limit connections
	sqliteDB.SetMaxOpenConns(1)
	sqliteDB.SetMaxIdleConns(1)

	if _, err := sqliteDB.Exec("PRAGMA foreign_keys = ON"); err != nil {
		return nil, fmt.Errorf("failed to enable foreign keys: %w", err)
	}

	tables, err := getTables(sqliteDB)
	if err != nil {
		return nil, fmt.Errorf("failed to check tables: %w", err)
	}

	if len(tables) == 0 {
		if err := runInitSQL(sqliteDB); err != nil {
			return nil, fmt.Errorf("failed to initialize database: %w", err)
		}
		fmt.Println("[DB] Database initialized successfully")
	} else {
		// Ensure indexes exist for existing databases
		if err := ensureIndexes(sqliteDB); err != nil {
			return nil, fmt.Errorf("failed to create indexes: %w", err)
		}
	}

	return sqliteDB, nil
}

var indexStatements = []string{
	`CREATE INDEX IF NOT EXISTS idx_tasks_user_id ON tasks (user_id, is_deleted)`,
	`CREATE INDEX IF NOT EXISTS idx_tasks_status ON tasks (user_id, status, is_deleted)`,
	`CREATE INDEX IF NOT EXISTS idx_tasks_category_id ON tasks (category_id)`,
	`CREATE INDEX IF NOT EXISTS idx_users_username ON users (username, is_deleted)`,
	`CREATE INDEX IF NOT EXISTS idx_login_log_username ON login_log (username)`,
}

func ensureIndexes(db *sql.DB) error {
	for _, stmt := range indexStatements {
		if _, err := db.Exec(stmt); err != nil {
			return fmt.Errorf("failed to create index: %w", err)
		}
	}
	return nil
}

func getTables(db *sql.DB) ([]string, error) {
	rows, err := db.Query("SELECT name FROM sqlite_master WHERE type='table'")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tables []string
	for rows.Next() {
		var name string
		if err := rows.Scan(&name); err != nil {
			return nil, err
		}
		tables = append(tables, name)
	}
	return tables, rows.Err()
}

// GetOrCreateJWTSecret retrieves the JWT secret from system_configs,
// or generates a new random one if it doesn't exist.
// Returns the secret and whether it was newly generated.
func GetOrCreateJWTSecret(db *sql.DB) (string, bool, error) {
	// Try to read existing secret from DB
	var secret string
	err := db.QueryRow("SELECT config_value FROM system_configs WHERE config_key = 'jwt_secret'").Scan(&secret)
	if err == nil && secret != "" {
		return secret, false, nil
	}

	// Generate a new random secret
	bytes := make([]byte, 32)
	if _, err := rand.Read(bytes); err != nil {
		return "", false, fmt.Errorf("failed to generate random secret: %w", err)
	}
	secret = hex.EncodeToString(bytes)

	// Persist to database
	_, err = db.Exec(
		"INSERT OR IGNORE INTO system_configs (config_key, config_value, group_name, description) VALUES (?, ?, ?, ?)",
		"jwt_secret", secret, "security", "JWT签名密钥（自动生成，请勿修改）",
	)
	if err != nil {
		return "", false, fmt.Errorf("failed to persist jwt_secret: %w", err)
	}

	return secret, true, nil
}

func runInitSQL(db *sql.DB) error {
	content, err := initSQL.ReadFile("init.sql")
	if err != nil {
		return fmt.Errorf("failed to read init.sql: %w", err)
	}

	// Parse SQL: remove comments carefully, handling strings
	lines := strings.Split(string(content), "\n")
	var cleanLines []string
	for _, line := range lines {
		line = removeLineComment(line)
		line = strings.TrimSpace(line)
		if line != "" {
			cleanLines = append(cleanLines, line)
		}
	}

	// Join and split by semicolon
	cleanSQL := strings.Join(cleanLines, "\n")
	statements := strings.Split(cleanSQL, ";")
	for _, stmt := range statements {
		stmt = strings.TrimSpace(stmt)
		if stmt == "" {
			continue
		}
		if _, err := db.Exec(stmt); err != nil {
			// Ignore UNIQUE constraint failures (INSERT OR IGNORE)
			if !strings.Contains(err.Error(), "UNIQUE constraint failed") {
				return fmt.Errorf("failed to execute: %s\nerror: %w", stmt, err)
			}
		}
	}
	return nil
}

// removeLineComment removes a line comment (-- ...) that is NOT inside a string literal.
// It scans the line tracking whether we are inside a single-quoted string.
func removeLineComment(line string) string {
	inString := false
	for i := 0; i < len(line); i++ {
		ch := line[i]
		if ch == '\'' {
			// Toggle string state (handles escaped '' as two separate toggles, which is correct)
			inString = !inString
		} else if !inString && i+1 < len(line) && line[i] == '-' && line[i+1] == '-' {
			return line[:i]
		}
	}
	return line
}
