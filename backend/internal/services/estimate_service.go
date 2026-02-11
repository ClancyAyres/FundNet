package services

import (
	"database/sql"
	"fundnet/backend/internal/models"
	"math"
	"strconv"
	"time"
)

// EstimateResult 估算结果
type EstimateResult struct {
	Code         string    `json:"code"`
	Name         string    `json:"name"`
	Nav          float64   `json:"nav"`
	EstimateNav  float64   `json:"estimate_nav"`
	DailyGrowth  float64   `json:"daily_growth"`
	EstimateTime time.Time `json:"estimate_time"`
}

// HistoryPoint 历史数据点
type HistoryPoint struct {
	Time        time.Time `json:"time"`
	EstimateNav float64   `json:"estimate_nav"`
	DailyGrowth float64   `json:"daily_growth"`
}

// EstimateService 估算服务
type EstimateService struct {
	db *sql.DB
}

// NewEstimateService 创建估算服务
func NewEstimateService() *EstimateService {
	return &EstimateService{
		db: models.GetDB(),
	}
}

// GetEstimate 获取基金估算
func (s *EstimateService) GetEstimate(code string) (*EstimateResult, error) {
	fund, err := s.GetFundFromDB(code)
	if err != nil {
		return nil, err
	}

	return &EstimateResult{
		Code:         fund.Code,
		Name:         fund.Name,
		Nav:          fund.Nav,
		EstimateNav:  fund.EstimateNav,
		DailyGrowth:  fund.DailyGrowth,
		EstimateTime: fund.EstimateTime,
	}, nil
}

// CalculateEstimate 计算估算净值
func (s *EstimateService) CalculateEstimate(fundCode string, shares, cost float64) (float64, float64, error) {
	fund, err := s.GetFundFromDB(fundCode)
	if err != nil {
		return 0, 0, err
	}

	var currentNav float64
	if fund.EstimateNav > 0 {
		currentNav = fund.EstimateNav
	} else if fund.Nav > 0 {
		currentNav = fund.Nav
	} else {
		currentNav = cost
	}

	currentValue := shares * currentNav
	profitRate := float64(0)
	if cost > 0 {
		profitRate = ((currentNav - cost) / cost) * 100
	}

	return currentValue, profitRate, nil
}

// RefreshAllEstimates 刷新所有估算
func (s *EstimateService) RefreshAllEstimates() {
	funds, _ := s.GetAllSubscribedFunds()
	for _, fund := range funds {
		s.RefreshEstimate(fund.Code)
	}
}

// RefreshEstimate 刷新单个基金估算
func (s *EstimateService) RefreshEstimate(code string) error {
	estimate, err := s.GetEstimate(code)
	if err != nil {
		return err
	}

	dailyGrowth := float64(0)
	if estimate.Nav > 0 {
		dailyGrowth = ((estimate.EstimateNav - estimate.Nav) / estimate.Nav) * 100
	}

	s.SaveEstimateHistory(code, estimate.EstimateNav, dailyGrowth)
	return nil
}

// SaveEstimateHistory 保存估算历史
func (s *EstimateService) SaveEstimateHistory(fundCode string, estimateNav, dailyGrowth float64) error {
	now := time.Now()
	_, err := s.db.Exec(`
		INSERT INTO estimate_history (fund_code, estimate_nav, daily_growth, recorded_at, created_at)
		VALUES (?, ?, ?, ?, ?)
	`, fundCode, estimateNav, dailyGrowth, now, now)
	return err
}

// GetHistory 获取历史估算数据
func (s *EstimateService) GetHistory(code string, days int) ([]HistoryPoint, error) {
	query := `
		SELECT estimate_nav, daily_growth, recorded_at
		FROM estimate_history
		WHERE fund_code = ?
		ORDER BY recorded_at DESC
		LIMIT ?
	`

	rows, err := s.db.Query(query, code, days)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var history []HistoryPoint
	for rows.Next() {
		var point HistoryPoint
		var dailyGrowth sql.NullString

		err := rows.Scan(&point.EstimateNav, &dailyGrowth, &point.Time)
		if err != nil {
			return nil, err
		}

		if dailyGrowth.Valid {
			point.DailyGrowth, _ = strconv.ParseFloat(dailyGrowth.String, 64)
		}

		history = append(history, point)
	}

	for i, j := 0, len(history)-1; i < j; i, j = i+1, j-1 {
		history[i], history[j] = history[j], history[i]
	}

	return history, nil
}

// GetAllSubscribedFunds 获取所有订阅的基金
func (s *EstimateService) GetAllSubscribedFunds() ([]models.Fund, error) {
	rows, err := s.db.Query(`
		SELECT id, code, name, sector, nav, nav_date, estimate_nav, estimate_time,
		       daily_growth, subscribed, subscribe_time, created_at, updated_at
		FROM funds
		WHERE subscribed = 1
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

// GetFundFromDB 从数据库获取基金信息
func (s *EstimateService) GetFundFromDB(code string) (*models.Fund, error) {
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

// CalculatePortfolioEstimate 计算组合估算
func (s *EstimateService) CalculatePortfolioEstimate(positions []models.Position) (map[string]interface{}, error) {
	var totalCost, totalValue, totalProfitLoss float64

	for _, pos := range positions {
		totalCost += pos.CostBasis
		totalValue += pos.CurrentValue
		totalProfitLoss += pos.ProfitLoss
	}

	profitRate := float64(0)
	if totalCost > 0 {
		profitRate = (totalProfitLoss / totalCost) * 100
	}

	totalCost = math.Round(totalCost*100) / 100
	totalValue = math.Round(totalValue*100) / 100
	totalProfitLoss = math.Round(totalProfitLoss*100) / 100
	profitRate = math.Round(profitRate*100) / 100

	return map[string]interface{}{
		"total_cost":     totalCost,
		"total_value":    totalValue,
		"total_profit":   totalProfitLoss,
		"profit_rate":    profitRate,
		"position_count": len(positions),
	}, nil
}
