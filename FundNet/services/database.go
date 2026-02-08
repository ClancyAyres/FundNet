package services

import (
	"FundNet/models"
	"database/sql"
	"fmt"

	_ "github.com/mattn/go-sqlite3"
)

type DatabaseService struct {
	db *sql.DB
}

func NewDatabaseService(dbPath string) (*DatabaseService, error) {
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	service := &DatabaseService{db: db}
	if err := service.initDB(); err != nil {
		return nil, fmt.Errorf("failed to initialize database: %w", err)
	}

	return service, nil
}

func (s *DatabaseService) initDB() error {
	// 创建基金表
	createFundTable := `
	CREATE TABLE IF NOT EXISTS funds (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		code TEXT UNIQUE NOT NULL,
		name TEXT NOT NULL,
		current_price REAL DEFAULT 0,
		change_rate REAL DEFAULT 0,
		last_updated TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	);`

	// 创建持仓表
	createPositionTable := `
	CREATE TABLE IF NOT EXISTS fund_positions (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		fund_code TEXT NOT NULL,
		shares REAL NOT NULL,
		cost_price REAL NOT NULL,
		group_name TEXT DEFAULT '默认分组',
		last_updated TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		FOREIGN KEY (fund_code) REFERENCES funds(code)
	);`

	// 创建估值历史表
	createValueHistoryTable := `
	CREATE TABLE IF NOT EXISTS fund_value_history (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		fund_code TEXT NOT NULL,
		current_price REAL NOT NULL,
		change_rate REAL NOT NULL,
		estimate_value REAL NOT NULL,
		estimate_change REAL NOT NULL,
		timestamp TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		FOREIGN KEY (fund_code) REFERENCES funds(code)
	);`

	// 创建配置表
	createConfigTable := `
	CREATE TABLE IF NOT EXISTS app_config (
		key TEXT PRIMARY KEY,
		value TEXT NOT NULL
	);`

	// 插入默认配置
	insertDefaultConfig := `
	INSERT OR IGNORE INTO app_config (key, value) VALUES
	('refresh_interval', '30'),
	('auto_refresh', 'true'),
	('max_history_days', '30');`

	_, err := s.db.Exec(createFundTable)
	if err != nil {
		return fmt.Errorf("failed to create funds table: %w", err)
	}

	_, err = s.db.Exec(createPositionTable)
	if err != nil {
		return fmt.Errorf("failed to create fund_positions table: %w", err)
	}

	_, err = s.db.Exec(createValueHistoryTable)
	if err != nil {
		return fmt.Errorf("failed to create fund_value_history table: %w", err)
	}

	_, err = s.db.Exec(createConfigTable)
	if err != nil {
		return fmt.Errorf("failed to create app_config table: %w", err)
	}

	_, err = s.db.Exec(insertDefaultConfig)
	if err != nil {
		return fmt.Errorf("failed to insert default config: %w", err)
	}

	return nil
}

// Fund operations
func (s *DatabaseService) AddFund(fund *models.Fund) error {
	query := `
	INSERT OR REPLACE INTO funds (code, name, current_price, change_rate, last_updated)
	VALUES (?, ?, ?, ?, ?)`

	_, err := s.db.Exec(query, fund.Code, fund.Name, fund.CurrentPrice, fund.ChangeRate, fund.LastUpdated)
	if err != nil {
		return fmt.Errorf("failed to add fund: %w", err)
	}

	return nil
}

func (s *DatabaseService) GetFund(code string) (*models.Fund, error) {
	query := `SELECT id, code, name, current_price, change_rate, last_updated FROM funds WHERE code = ?`
	row := s.db.QueryRow(query, code)

	var fund models.Fund
	err := row.Scan(&fund.ID, &fund.Code, &fund.Name, &fund.CurrentPrice, &fund.ChangeRate, &fund.LastUpdated)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get fund: %w", err)
	}

	return &fund, nil
}

func (s *DatabaseService) GetAllFunds() ([]*models.Fund, error) {
	query := `SELECT id, code, name, current_price, change_rate, last_updated FROM funds ORDER BY code`
	rows, err := s.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("failed to get all funds: %w", err)
	}
	defer rows.Close()

	var funds []*models.Fund
	for rows.Next() {
		var fund models.Fund
		err := rows.Scan(&fund.ID, &fund.Code, &fund.Name, &fund.CurrentPrice, &fund.ChangeRate, &fund.LastUpdated)
		if err != nil {
			return nil, fmt.Errorf("failed to scan fund: %w", err)
		}
		funds = append(funds, &fund)
	}

	return funds, nil
}

// Position operations
func (s *DatabaseService) AddPosition(position *models.FundPosition) error {
	query := `
	INSERT OR REPLACE INTO fund_positions (fund_code, shares, cost_price, group_name, last_updated)
	VALUES (?, ?, ?, ?, ?)`

	_, err := s.db.Exec(query, position.FundCode, position.Shares, position.CostPrice, position.GroupName, position.LastUpdated)
	if err != nil {
		return fmt.Errorf("failed to add position: %w", err)
	}

	return nil
}

func (s *DatabaseService) GetPosition(fundCode string) (*models.FundPosition, error) {
	query := `SELECT id, fund_code, shares, cost_price, group_name, last_updated FROM fund_positions WHERE fund_code = ?`
	row := s.db.QueryRow(query, fundCode)

	var position models.FundPosition
	err := row.Scan(&position.ID, &position.FundCode, &position.Shares, &position.CostPrice, &position.GroupName, &position.LastUpdated)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get position: %w", err)
	}

	return &position, nil
}

func (s *DatabaseService) GetAllPositions() ([]*models.FundPosition, error) {
	query := `SELECT id, fund_code, shares, cost_price, group_name, last_updated FROM fund_positions ORDER BY group_name, fund_code`
	rows, err := s.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("failed to get all positions: %w", err)
	}
	defer rows.Close()

	var positions []*models.FundPosition
	for rows.Next() {
		var position models.FundPosition
		err := rows.Scan(&position.ID, &position.FundCode, &position.Shares, &position.CostPrice, &position.GroupName, &position.LastUpdated)
		if err != nil {
			return nil, fmt.Errorf("failed to scan position: %w", err)
		}
		positions = append(positions, &position)
	}

	return positions, nil
}

func (s *DatabaseService) DeletePosition(fundCode string) error {
	query := `DELETE FROM fund_positions WHERE fund_code = ?`
	_, err := s.db.Exec(query, fundCode)
	if err != nil {
		return fmt.Errorf("failed to delete position: %w", err)
	}

	return nil
}

// Value history operations
func (s *DatabaseService) AddValueHistory(value *models.FundValue) error {
	query := `
	INSERT INTO fund_value_history (fund_code, current_price, change_rate, estimate_value, estimate_change, timestamp)
	VALUES (?, ?, ?, ?, ?, ?)`

	_, err := s.db.Exec(query, value.FundCode, value.CurrentPrice, value.ChangeRate, value.EstimateValue, value.EstimateChange, value.Timestamp)
	if err != nil {
		return fmt.Errorf("failed to add value history: %w", err)
	}

	return nil
}

func (s *DatabaseService) GetValueHistory(fundCode string, days int) ([]*models.FundValue, error) {
	query := `
	SELECT fund_code, current_price, change_rate, estimate_value, estimate_change, timestamp
	FROM fund_value_history
	WHERE fund_code = ? AND timestamp >= datetime('now', '-? days')
	ORDER BY timestamp DESC`

	rows, err := s.db.Query(query, fundCode, days)
	if err != nil {
		return nil, fmt.Errorf("failed to get value history: %w", err)
	}
	defer rows.Close()

	var history []*models.FundValue
	for rows.Next() {
		var value models.FundValue
		err := rows.Scan(&value.FundCode, &value.CurrentPrice, &value.ChangeRate, &value.EstimateValue, &value.EstimateChange, &value.Timestamp)
		if err != nil {
			return nil, fmt.Errorf("failed to scan value history: %w", err)
		}
		history = append(history, &value)
	}

	return history, nil
}

// Config operations
func (s *DatabaseService) GetConfig(key string) (string, error) {
	query := `SELECT value FROM app_config WHERE key = ?`
	row := s.db.QueryRow(query, key)

	var value string
	err := row.Scan(&value)
	if err != nil {
		if err == sql.ErrNoRows {
			return "", nil
		}
		return "", fmt.Errorf("failed to get config: %w", err)
	}

	return value, nil
}

func (s *DatabaseService) SetConfig(key, value string) error {
	query := `INSERT OR REPLACE INTO app_config (key, value) VALUES (?, ?)`
	_, err := s.db.Exec(query, key, value)
	if err != nil {
		return fmt.Errorf("failed to set config: %w", err)
	}

	return nil
}

func (s *DatabaseService) Close() error {
	return s.db.Close()
}

// Utility functions
func (s *DatabaseService) GetGroupNames() ([]string, error) {
	query := `SELECT DISTINCT group_name FROM fund_positions ORDER BY group_name`
	rows, err := s.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("failed to get group names: %w", err)
	}
	defer rows.Close()

	var groups []string
	for rows.Next() {
		var groupName string
		err := rows.Scan(&groupName)
		if err != nil {
			return nil, fmt.Errorf("failed to scan group name: %w", err)
		}
		groups = append(groups, groupName)
	}

	return groups, nil
}

func (s *DatabaseService) GetPortfolioSummary() (*models.PortfolioSummary, error) {
	// 获取所有持仓
	positions, err := s.GetAllPositions()
	if err != nil {
		return nil, fmt.Errorf("failed to get positions: %w", err)
	}

	// 获取所有基金的最新价格
	funds, err := s.GetAllFunds()
	if err != nil {
		return nil, fmt.Errorf("failed to get funds: %w", err)
	}

	// 创建基金价格映射
	priceMap := make(map[string]float64)
	for _, fund := range funds {
		priceMap[fund.Code] = fund.CurrentPrice
	}

	// 计算投资组合概览
	var totalValue, totalCost, totalGain, dailyGain float64

	for _, position := range positions {
		currentPrice, exists := priceMap[position.FundCode]
		if !exists {
			continue
		}

		positionValue := currentPrice * position.Shares
		positionCost := position.CostPrice * position.Shares
		positionGain := positionValue - positionCost

		totalValue += positionValue
		totalCost += positionCost
		totalGain += positionGain

		// 计算当日收益（基于涨跌幅）
		fund, _ := s.GetFund(position.FundCode)
		if fund != nil {
			dailyGain += position.Shares * position.CostPrice * fund.ChangeRate / 100
		}
	}

	summary := &models.PortfolioSummary{
		TotalValue: totalValue,
		TotalCost:  totalCost,
		TotalGain:  totalGain,
		DailyGain:  dailyGain,
	}

	if totalCost > 0 {
		summary.TotalGainRate = totalGain / totalCost * 100
	}

	return summary, nil
}

func (s *DatabaseService) GetGroupSummary(groupName string) (*models.GroupSummary, error) {
	query := `
	SELECT fp.fund_code, fp.shares, fp.cost_price, f.current_price, f.change_rate
	FROM fund_positions fp
	JOIN funds f ON fp.fund_code = f.code
	WHERE fp.group_name = ?`

	rows, err := s.db.Query(query, groupName)
	if err != nil {
		return nil, fmt.Errorf("failed to get group summary: %w", err)
	}
	defer rows.Close()

	var value, cost, gain, dailyGain float64

	for rows.Next() {
		var fundCode string
		var shares, costPrice, currentPrice, changeRate float64
		err := rows.Scan(&fundCode, &shares, &costPrice, &currentPrice, &changeRate)
		if err != nil {
			return nil, fmt.Errorf("failed to scan group summary: %w", err)
		}

		positionValue := currentPrice * shares
		positionCost := costPrice * shares
		positionGain := positionValue - positionCost

		value += positionValue
		cost += positionCost
		gain += positionGain
		dailyGain += shares * costPrice * changeRate / 100
	}

	summary := &models.GroupSummary{
		GroupName: groupName,
		Value:     value,
		Cost:      cost,
		Gain:      gain,
		DailyGain: dailyGain,
	}

	if cost > 0 {
		summary.GainRate = gain / cost * 100
	}

	return summary, nil
}
