package repository

import (
	"database/sql"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

// SQLite SQLite 数据库
type SQLite struct {
	db *sql.DB
}

// NewSQLite 创建 SQLite 实例
func NewSQLite(dsn string) (*SQLite, error) {
	db, err := sql.Open("sqlite3", dsn)
	if err != nil {
		return nil, err
	}

	// 设置连接池
	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(5)
	db.SetConnMaxLifetime(5 * time.Minute)

	// 验证连接
	if err := db.Ping(); err != nil {
		return nil, err
	}

	return &SQLite{db: db}, nil
}

// Close 关闭数据库连接
func (s *SQLite) Close() error {
	return s.db.Close()
}

// DB 获取数据库实例
func (s *SQLite) DB() *sql.DB {
	return s.db
}

// InitSchema 初始化数据库表结构
func (s *SQLite) InitSchema() error {
	schema := `
	CREATE TABLE IF NOT EXISTS tasks (
		id TEXT PRIMARY KEY,
		name TEXT NOT NULL,
		description TEXT,
		status INTEGER NOT NULL DEFAULT 1,
		priority INTEGER NOT NULL DEFAULT 2,
		task_type TEXT,
		input_params TEXT,
		output_result TEXT,
		dependencies TEXT,
		retry_count INTEGER NOT NULL DEFAULT 0,
		max_retries INTEGER NOT NULL DEFAULT 0,
		error_message TEXT,
		created_at TEXT NOT NULL,
		updated_at TEXT NOT NULL,
		started_at TEXT,
		completed_at TEXT,
		created_by TEXT
	);

	CREATE INDEX IF NOT EXISTS idx_tasks_status ON tasks(status);
	CREATE INDEX IF NOT EXISTS idx_tasks_priority ON tasks(priority);
	CREATE INDEX IF NOT EXISTS idx_tasks_created_by ON tasks(created_by);
	CREATE INDEX IF NOT EXISTS idx_tasks_created_at ON tasks(created_at);

	CREATE TABLE IF NOT EXISTS task_events (
		id TEXT PRIMARY KEY,
		task_id TEXT NOT NULL,
		from_status INTEGER NOT NULL,
		to_status INTEGER NOT NULL,
		message TEXT,
		timestamp TEXT NOT NULL,
		operator TEXT,
		FOREIGN KEY (task_id) REFERENCES tasks(id) ON DELETE CASCADE
	);

	CREATE INDEX IF NOT EXISTS idx_task_events_task_id ON task_events(task_id);
	CREATE INDEX IF NOT EXISTS idx_task_events_timestamp ON task_events(timestamp);
	`

	_, err := s.db.Exec(schema)
	return err
}

// ExecTx 执行事务
func (s *SQLite) ExecTx(fn func(*sql.Tx) error) error {
	tx, err := s.db.Begin()
	if err != nil {
		return err
	}

	if err := fn(tx); err != nil {
		if rbErr := tx.Rollback(); rbErr != nil {
			return rbErr
		}
		return err
	}

	return tx.Commit()
}
