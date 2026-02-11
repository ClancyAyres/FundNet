package models

import (
	"database/sql"
	"fmt"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

var db *sql.DB

type Fund struct {
	ID            int64     `json:"id"`
	Code          string    `json:"code"`
	Name          string    `json:"name"`
	Sector        string    `json:"sector"`
	Nav           float64   `json:"nav"`
	NavDate       time.Time `json:"nav_date"`
	EstimateNav   float64   `json:"estimate_nav"`
	EstimateTime  time.Time `json:"estimate_time"`
	DailyGrowth   float64   `json:"daily_growth"`
	Subscribed    bool      `json:"subscribed"`
	SubscribeTime time.Time `json:"subscribe_time"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

type Position struct {
	ID           int64     `json:"id"`
	FundCode     string    `json:"fund_code"`
	FundName     string    `json:"fund_name"`
	Shares       float64   `json:"shares"`
	Cost         float64   `json:"cost"`
	CostBasis    float64   `json:"cost_basis"`
	CurrentValue float64   `json:"current_value"`
	ProfitLoss   float64   `json:"profit_loss"`
	ProfitRate   float64   `json:"profit_rate"`
	DailyGrowth  float64   `json:"daily_growth"`
	Sector       string    `json:"sector"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

type Sector struct {
	ID        int64     `json:"id"`
	Name      string    `json:"name"`
	Color     string    `json:"color"`
	SortOrder int       `json:"sort_order"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type EstimateHistory struct {
	ID          int64     `json:"id"`
	FundCode    string    `json:"fund_code"`
	EstimateNav float64   `json:"estimate_nav"`
	DailyGrowth float64   `json:"daily_growth"`
	RecordedAt  time.Time `json:"recorded_at"`
	CreatedAt   time.Time `json:"created_at"`
}

func InitDB(path string) error {
	var err error
	db, err = sql.Open("sqlite3", path)
	if err != nil {
		return fmt.Errorf("failed to open database: %w", err)
	}

	if err := createTables(); err != nil {
		return fmt.Errorf("failed to create tables: %w", err)
	}

	return nil
}

func createTables() error {
	// 先删除旧表（如果存在），重新创建以应用新的 DATETIME 类型
	dropTables := []string{
		"DROP TABLE IF EXISTS funds",
		"DROP TABLE IF EXISTS positions",
		"DROP TABLE IF EXISTS sectors",
		"DROP TABLE IF EXISTS estimate_history",
		"DROP TABLE IF EXISTS config",
	}
	for _, stmt := range dropTables {
		if _, err := db.Exec(stmt); err != nil {
			return err
		}
	}

	tables := []string{
		`CREATE TABLE funds (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			code TEXT UNIQUE NOT NULL,
			name TEXT,
			sector TEXT,
			nav REAL DEFAULT 0,
			nav_date DATETIME,
			estimate_nav REAL DEFAULT 0,
			estimate_time DATETIME,
			daily_growth REAL DEFAULT 0,
			subscribed INTEGER DEFAULT 0,
			subscribe_time DATETIME,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
		)`,
		`CREATE TABLE positions (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			fund_code TEXT NOT NULL,
			fund_name TEXT,
			shares REAL DEFAULT 0,
			cost REAL DEFAULT 0,
			cost_basis REAL DEFAULT 0,
			current_value REAL DEFAULT 0,
			profit_loss REAL DEFAULT 0,
			profit_rate REAL DEFAULT 0,
			daily_growth REAL DEFAULT 0,
			sector TEXT,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
		)`,
		`CREATE TABLE sectors (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			name TEXT UNIQUE NOT NULL,
			color TEXT DEFAULT '#1890ff',
			sort_order INTEGER DEFAULT 0,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
		)`,
		`CREATE TABLE estimate_history (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			fund_code TEXT NOT NULL,
			estimate_nav REAL DEFAULT 0,
			daily_growth REAL DEFAULT 0,
			recorded_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP
		)`,
		`CREATE TABLE config (
			key TEXT PRIMARY KEY,
			value TEXT,
			updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
		)`,
	}

	for _, table := range tables {
		if _, err := db.Exec(table); err != nil {
			return err
		}
	}

	defaultSectors := []string{
		"科技", "医疗", "新能源", "QDII", "消费", "金融", "军工", "半导体", "互联网", "房地产",
	}
	for _, sector := range defaultSectors {
		if _, err := db.Exec("INSERT OR IGNORE INTO sectors (name) VALUES (?)", sector); err != nil {
			return err
		}
	}

	return nil
}

func CloseDB() {
	if db != nil {
		db.Close()
	}
}

func GetDB() *sql.DB {
	return db
}
