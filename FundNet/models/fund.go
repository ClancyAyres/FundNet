package models

import (
	"time"
)

// Fund 基金信息
type Fund struct {
	ID           int       `json:"id" db:"id"`
	Code         string    `json:"code" db:"code"`
	Name         string    `json:"name" db:"name"`
	CurrentPrice float64   `json:"current_price" db:"current_price"`
	ChangeRate   float64   `json:"change_rate" db:"change_rate"`
	LastUpdated  time.Time `json:"last_updated" db:"last_updated"`
}

// FundPosition 基金持仓信息
type FundPosition struct {
	ID          int       `json:"id" db:"id"`
	FundCode    string    `json:"fund_code" db:"fund_code"`
	Shares      float64   `json:"shares" db:"shares"`
	CostPrice   float64   `json:"cost_price" db:"cost_price"`
	GroupName   string    `json:"group_name" db:"group_name"`
	LastUpdated time.Time `json:"last_updated" db:"last_updated"`
}

// FundValue 基金估值信息
type FundValue struct {
	FundCode       string    `json:"fund_code"`
	CurrentPrice   float64   `json:"current_price"`
	ChangeRate     float64   `json:"change_rate"`
	EstimateValue  float64   `json:"estimate_value"`
	EstimateChange float64   `json:"estimate_change"`
	Timestamp      time.Time `json:"timestamp"`
}

// PortfolioSummary 投资组合概览
type PortfolioSummary struct {
	TotalValue    float64 `json:"total_value"`
	TotalCost     float64 `json:"total_cost"`
	TotalGain     float64 `json:"total_gain"`
	TotalGainRate float64 `json:"total_gain_rate"`
	DailyGain     float64 `json:"daily_gain"`
	DailyGainRate float64 `json:"daily_gain_rate"`
}

// GroupSummary 分组概览
type GroupSummary struct {
	GroupName     string  `json:"group_name"`
	Value         float64 `json:"value"`
	Cost          float64 `json:"cost"`
	Gain          float64 `json:"gain"`
	GainRate      float64 `json:"gain_rate"`
	DailyGain     float64 `json:"daily_gain"`
	DailyGainRate float64 `json:"daily_gain_rate"`
}
