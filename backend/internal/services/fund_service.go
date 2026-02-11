package services

import (
	"database/sql"
	"fmt"
	"time"

	"fundnet/backend/internal/models"
)

// FundService 基金服务
type FundService struct {
	db *sql.DB
}

// NewFundService 创建基金服务
func NewFundService() *FundService {
	return &FundService{
		db: models.GetDB(),
	}
}

// GetAllFunds 获取所有订阅的基金
func (s *FundService) GetAllFunds() ([]models.Fund, error) {
	rows, err := s.db.Query(`
		SELECT id, code, name, sector, nav, nav_date, estimate_nav, estimate_time,
		       daily_growth, subscribed, subscribe_time, created_at, updated_at
		FROM funds
		WHERE subscribed = 1
		ORDER BY updated_at DESC
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var funds []models.Fund
	for rows.Next() {
		var fund models.Fund
		var estimateTime, navDate sql.NullString

		err := rows.Scan(
			&fund.ID, &fund.Code, &fund.Name, &fund.Sector,
			&fund.Nav, &navDate, &fund.EstimateNav, &estimateTime,
			&fund.DailyGrowth, &fund.Subscribed, &fund.SubscribeTime,
			&fund.CreatedAt, &fund.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}

		if estimateTime.Valid {
			fund.EstimateTime, _ = time.Parse("2006-01-02 15:04:05", estimateTime.String)
		}
		if navDate.Valid {
			fund.NavDate, _ = time.Parse("2006-01-02", navDate.String)
		}

		funds = append(funds, fund)
	}

	return funds, nil
}

// GetFundByCode 根据代码获取基金
func (s *FundService) GetFundByCode(code string) (*models.Fund, error) {
	var fund models.Fund
	var estimateTime, navDate sql.NullString

	err := s.db.QueryRow(`
		SELECT id, code, name, sector, nav, nav_date, estimate_nav, estimate_time,
		       daily_growth, subscribed, subscribe_time, created_at, updated_at
		FROM funds
		WHERE code = ?
	`, code).Scan(
		&fund.ID, &fund.Code, &fund.Name, &fund.Sector,
		&fund.Nav, &navDate, &fund.EstimateNav, &estimateTime,
		&fund.DailyGrowth, &fund.Subscribed, &fund.SubscribeTime,
		&fund.CreatedAt, &fund.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	if estimateTime.Valid {
		fund.EstimateTime, _ = time.Parse("2006-01-02 15:04:05", estimateTime.String)
	}
	if navDate.Valid {
		fund.NavDate, _ = time.Parse("2006-01-02", navDate.String)
	}

	return &fund, nil
}

// AddFund 添加基金订阅
func (s *FundService) AddFund(code, name, sector string) (*models.Fund, error) {
	now := time.Now()
	result, err := s.db.Exec(`
		INSERT OR REPLACE INTO funds (code, name, sector, subscribed, subscribe_time, created_at, updated_at)
		VALUES (?, ?, ?, 1, ?, ?, ?)
	`, code, name, sector, now, now, now)
	if err != nil {
		return nil, err
	}

	id, _ := result.LastInsertId()
	return &models.Fund{
		ID:            id,
		Code:          code,
		Name:          name,
		Sector:        sector,
		Subscribed:    true,
		SubscribeTime: now,
		CreatedAt:     now,
		UpdatedAt:     now,
	}, nil
}

// RemoveFund 取消基金订阅
func (s *FundService) RemoveFund(code string) error {
	_, err := s.db.Exec(`DELETE FROM funds WHERE code = ?`, code)
	return err
}

// UpdateFund 更新基金信息
func (s *FundService) UpdateFund(code, name, sector string) (*models.Fund, error) {
	now := time.Now()
	_, err := s.db.Exec(`
		UPDATE funds SET name = ?, sector = ?, updated_at = ?
		WHERE code = ?
	`, name, sector, now, code)
	if err != nil {
		return nil, err
	}

	return s.GetFundByCode(code)
}

// UpdateFundData 更新基金数据
func (s *FundService) UpdateFundData(code string, nav, estimateNav float64, navDate, dailyGrowth string) error {
	now := time.Now()
	_, err := s.db.Exec(`
		UPDATE funds SET nav = ?, nav_date = ?, estimate_nav = ?,
		       estimate_time = ?, daily_growth = ?, updated_at = ?
		WHERE code = ?
	`, nav, navDate, estimateNav, now, dailyGrowth, now, code)
	return err
}

// UpdateAllFundData 更新所有基金数据
func (s *FundService) UpdateAllFundData() {
	funds, _ := s.GetAllFunds()
	for _, fund := range funds {
		// 这里可以调用爬虫获取最新数据
		_ = fund
	}
}

// GetAllSectors 获取所有板块
func (s *FundService) GetAllSectors() ([]models.Sector, error) {
	rows, err := s.db.Query(`
		SELECT id, name, color, sort_order, created_at, updated_at
		FROM sectors
		ORDER BY sort_order, name
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var sectors []models.Sector
	for rows.Next() {
		var sector models.Sector
		var createdAt, updatedAt string
		err := rows.Scan(&sector.ID, &sector.Name, &sector.Color, &sector.SortOrder,
			&createdAt, &updatedAt)
		if err != nil {
			return nil, err
		}
		sector.CreatedAt, _ = time.Parse("2006-01-02 15:04:05", createdAt)
		sector.UpdatedAt, _ = time.Parse("2006-01-02 15:04:05", updatedAt)
		sectors = append(sectors, sector)
	}

	return sectors, nil
}

// CreateSector 创建板块
func (s *FundService) CreateSector(name, color string, sortOrder int) (*models.Sector, error) {
	now := time.Now()
	if color == "" {
		color = "#1890ff"
	}
	result, err := s.db.Exec(`
		INSERT OR IGNORE INTO sectors (name, color, sort_order, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?)
	`, name, color, sortOrder, now, now)
	if err != nil {
		return nil, err
	}

	id, _ := result.LastInsertId()
	return &models.Sector{
		ID:        id,
		Name:      name,
		Color:     color,
		SortOrder: sortOrder,
		CreatedAt: now,
		UpdatedAt: now,
	}, nil
}

// UpdateSector 更新板块
func (s *FundService) UpdateSector(id int64, name, color string, sortOrder int) (*models.Sector, error) {
	now := time.Now()
	_, err := s.db.Exec(`
		UPDATE sectors SET name = ?, color = ?, sort_order = ?, updated_at = ?
		WHERE id = ?
	`, name, color, sortOrder, now, id)
	if err != nil {
		return nil, err
	}

	return s.GetSectorByID(id)
}

// DeleteSector 删除板块
func (s *FundService) DeleteSector(id int64) error {
	_, err := s.db.Exec(`DELETE FROM sectors WHERE id = ?`, id)
	return err
}

// GetSectorByID 根据ID获取板块
func (s *FundService) GetSectorByID(id int64) (*models.Sector, error) {
	var sector models.Sector
	err := s.db.QueryRow(`
		SELECT id, name, color, sort_order, created_at, updated_at
		FROM sectors WHERE id = ?
	`, id).Scan(&sector.ID, &sector.Name, &sector.Color, &sector.SortOrder,
		&sector.CreatedAt, &sector.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return &sector, nil
}

// GetAllPositions 获取所有持仓
func (s *FundService) GetAllPositions() ([]models.Position, error) {
	rows, err := s.db.Query(`
		SELECT id, fund_code, fund_name, shares, cost, cost_basis, current_value,
		       profit_loss, profit_rate, daily_growth, sector, created_at, updated_at
		FROM positions
		ORDER BY created_at DESC
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var positions []models.Position
	for rows.Next() {
		var position models.Position
		err := rows.Scan(
			&position.ID, &position.FundCode, &position.FundName,
			&position.Shares, &position.Cost, &position.CostBasis,
			&position.CurrentValue, &position.ProfitLoss, &position.ProfitRate,
			&position.DailyGrowth, &position.Sector,
			&position.CreatedAt, &position.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		positions = append(positions, position)
	}

	return positions, nil
}

// AddPosition 添加持仓
func (s *FundService) AddPosition(fundCode, fundName string, shares, cost float64, sector string) (*models.Position, error) {
	now := time.Now()
	costBasis := shares * cost
	result, err := s.db.Exec(`
		INSERT INTO positions (fund_code, fund_name, shares, cost, cost_basis, sector, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?)
	`, fundCode, fundName, shares, cost, costBasis, sector, now, now)
	if err != nil {
		return nil, err
	}

	id, _ := result.LastInsertId()
	return &models.Position{
		ID:        id,
		FundCode:  fundCode,
		FundName:  fundName,
		Shares:    shares,
		Cost:      cost,
		CostBasis: costBasis,
		Sector:    sector,
		CreatedAt: now,
		UpdatedAt: now,
	}, nil
}

// UpdatePosition 更新持仓
func (s *FundService) UpdatePosition(id int64, shares, cost float64, sector string) (*models.Position, error) {
	now := time.Now()
	costBasis := shares * cost
	_, err := s.db.Exec(`
		UPDATE positions SET shares = ?, cost = ?, cost_basis = ?, sector = ?, updated_at = ?
		WHERE id = ?
	`, shares, cost, costBasis, sector, now, id)
	if err != nil {
		return nil, err
	}

	return s.GetPositionByID(id)
}

// DeletePosition 删除持仓
func (s *FundService) DeletePosition(id int64) error {
	_, err := s.db.Exec(`DELETE FROM positions WHERE id = ?`, id)
	return err
}

// GetPositionByID 根据ID获取持仓
func (s *FundService) GetPositionByID(id int64) (*models.Position, error) {
	var position models.Position
	err := s.db.QueryRow(`
		SELECT id, fund_code, fund_name, shares, cost, cost_basis, current_value,
		       profit_loss, profit_rate, daily_growth, sector, created_at, updated_at
		FROM positions WHERE id = ?
	`, id).Scan(
		&position.ID, &position.FundCode, &position.FundName,
		&position.Shares, &position.Cost, &position.CostBasis,
		&position.CurrentValue, &position.ProfitLoss, &position.ProfitRate,
		&position.DailyGrowth, &position.Sector,
		&position.CreatedAt, &position.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return &position, nil
}

// GetAssetStats 获取资产统计
func (s *FundService) GetAssetStats() (map[string]interface{}, error) {
	positions, err := s.GetAllPositions()
	if err != nil {
		return nil, err
	}

	var totalCostBasis, totalCurrentValue, totalProfitLoss float64
	for _, pos := range positions {
		totalCostBasis += pos.CostBasis
		totalCurrentValue += pos.CurrentValue
		totalProfitLoss += pos.ProfitLoss
	}

	profitRate := float64(0)
	if totalCostBasis > 0 {
		profitRate = (totalProfitLoss / totalCostBasis) * 100
	}

	return map[string]interface{}{
		"total_cost_basis":    totalCostBasis,
		"total_current_value": totalCurrentValue,
		"total_profit_loss":   totalProfitLoss,
		"profit_rate":         profitRate,
		"position_count":      len(positions),
	}, nil
}

// GetAssetSummary 获取资产摘要（按板块分组）
func (s *FundService) GetAssetSummary() (map[string]interface{}, error) {
	positions, err := s.GetAllPositions()
	if err != nil {
		return nil, err
	}

	sectorStats := make(map[string]map[string]float64)
	for _, pos := range positions {
		if _, ok := sectorStats[pos.Sector]; !ok {
			sectorStats[pos.Sector] = map[string]float64{
				"cost_basis":    0,
				"current_value": 0,
				"profit_loss":   0,
			}
		}
		sectorStats[pos.Sector]["cost_basis"] += pos.CostBasis
		sectorStats[pos.Sector]["current_value"] += pos.CurrentValue
		sectorStats[pos.Sector]["profit_loss"] += pos.ProfitLoss
	}

	// 计算各板块占比
	var totalValue float64
	for _, stats := range sectorStats {
		totalValue += stats["current_value"]
	}

	sectors := make([]map[string]interface{}, 0)
	for name, stats := range sectorStats {
		weight := float64(0)
		if totalValue > 0 {
			weight = (stats["current_value"] / totalValue) * 100
		}
		profitRate := float64(0)
		if stats["cost_basis"] > 0 {
			profitRate = (stats["profit_loss"] / stats["cost_basis"]) * 100
		}
		sectors = append(sectors, map[string]interface{}{
			"name":          name,
			"cost_basis":    stats["cost_basis"],
			"current_value": stats["current_value"],
			"profit_loss":   stats["profit_loss"],
			"profit_rate":   profitRate,
			"weight":        weight,
		})
	}

	// 总计
	totalCostBasis, totalCurrentValue, totalProfitLoss := 0.0, 0.0, 0.0
	for _, s := range sectors {
		totalCostBasis += s["cost_basis"].(float64)
		totalCurrentValue += s["current_value"].(float64)
		totalProfitLoss += s["profit_loss"].(float64)
	}

	totalProfitRate := float64(0)
	if totalCostBasis > 0 {
		totalProfitRate = (totalProfitLoss / totalCostBasis) * 100
	}

	return map[string]interface{}{
		"sectors":             sectors,
		"total_cost_basis":    totalCostBasis,
		"total_current_value": totalCurrentValue,
		"total_profit_loss":   totalProfitLoss,
		"total_profit_rate":   totalProfitRate,
	}, nil
}

// GetConfig 获取配置
func (s *FundService) GetConfig() map[string]interface{} {
	refreshInterval := 60
	logLevel := "info"

	row := s.db.QueryRow(`SELECT value FROM config WHERE key = 'refresh_interval'`)
	row.Scan(&refreshInterval)

	row = s.db.QueryRow(`SELECT value FROM config WHERE key = 'log_level'`)
	row.Scan(&logLevel)

	return map[string]interface{}{
		"refresh_interval": refreshInterval,
		"log_level":        logLevel,
	}
}

// UpdateConfig 更新配置
func (s *FundService) UpdateConfig(refreshInterval int, logLevel string) error {
	now := time.Now()
	_, err := s.db.Exec(`
		INSERT OR REPLACE INTO config (key, value, updated_at) VALUES (?, ?, ?)
	`, "refresh_interval", fmt.Sprintf("%d", refreshInterval), now)
	if err != nil {
		return err
	}

	_, err = s.db.Exec(`
		INSERT OR REPLACE INTO config (key, value, updated_at) VALUES (?, ?, ?)
	`, "log_level", logLevel, now)
	return err
}
