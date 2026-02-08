package services

import (
	"FundNet/models"
	"fmt"
	"time"
)

type PortfolioService struct {
	db *DatabaseService
}

func NewPortfolioService(db *DatabaseService) *PortfolioService {
	return &PortfolioService{
		db: db,
	}
}

// AddFundPosition 添加基金持仓
func (s *PortfolioService) AddFundPosition(position *models.FundPosition) error {
	// 验证基金代码是否存在
	fund, err := s.db.GetFund(position.FundCode)
	if err != nil {
		return fmt.Errorf("failed to get fund: %w", err)
	}
	if fund == nil {
		return fmt.Errorf("fund with code %s not found", position.FundCode)
	}

	// 设置更新时间
	position.LastUpdated = time.Now()

	// 保存到数据库
	err = s.db.AddPosition(position)
	if err != nil {
		return fmt.Errorf("failed to add position: %w", err)
	}

	return nil
}

// UpdateFundPosition 更新基金持仓
func (s *PortfolioService) UpdateFundPosition(position *models.FundPosition) error {
	// 验证基金代码是否存在
	fund, err := s.db.GetFund(position.FundCode)
	if err != nil {
		return fmt.Errorf("failed to get fund: %w", err)
	}
	if fund == nil {
		return fmt.Errorf("fund with code %s not found", position.FundCode)
	}

	// 设置更新时间
	position.LastUpdated = time.Now()

	// 保存到数据库
	err = s.db.AddPosition(position)
	if err != nil {
		return fmt.Errorf("failed to update position: %w", err)
	}

	return nil
}

// DeleteFundPosition 删除基金持仓
func (s *PortfolioService) DeleteFundPosition(fundCode string) error {
	err := s.db.DeletePosition(fundCode)
	if err != nil {
		return fmt.Errorf("failed to delete position: %w", err)
	}

	return nil
}

// GetPortfolioSummary 获取投资组合概览
func (s *PortfolioService) GetPortfolioSummary() (*models.PortfolioSummary, error) {
	summary, err := s.db.GetPortfolioSummary()
	if err != nil {
		return nil, fmt.Errorf("failed to get portfolio summary: %w", err)
	}

	return summary, nil
}

// GetGroupSummary 获取分组概览
func (s *PortfolioService) GetGroupSummary(groupName string) (*models.GroupSummary, error) {
	summary, err := s.db.GetGroupSummary(groupName)
	if err != nil {
		return nil, fmt.Errorf("failed to get group summary: %w", err)
	}

	return summary, nil
}

// GetAllGroups 获取所有分组
func (s *PortfolioService) GetAllGroups() ([]string, error) {
	groups, err := s.db.GetGroupNames()
	if err != nil {
		return nil, fmt.Errorf("failed to get groups: %w", err)
	}

	return groups, nil
}

// GetGroupPositions 获取分组下的所有持仓
func (s *PortfolioService) GetGroupPositions(groupName string) ([]*models.FundPosition, error) {
	positions, err := s.db.GetAllPositions()
	if err != nil {
		return nil, fmt.Errorf("failed to get positions: %w", err)
	}

	var groupPositions []*models.FundPosition
	for _, position := range positions {
		if position.GroupName == groupName {
			groupPositions = append(groupPositions, position)
		}
	}

	return groupPositions, nil
}

// CalculatePositionValue 计算持仓价值
func (s *PortfolioService) CalculatePositionValue(fundCode string) (*models.FundValue, error) {
	// 获取持仓信息
	position, err := s.db.GetPosition(fundCode)
	if err != nil {
		return nil, fmt.Errorf("failed to get position: %w", err)
	}
	if position == nil {
		return nil, fmt.Errorf("position not found for fund %s", fundCode)
	}

	// 获取基金最新价格
	fund, err := s.db.GetFund(fundCode)
	if err != nil {
		return nil, fmt.Errorf("failed to get fund: %w", err)
	}
	if fund == nil {
		return nil, fmt.Errorf("fund not found: %s", fundCode)
	}

	// 计算持仓价值
	positionValue := position.Shares * fund.CurrentPrice
	positionCost := position.Shares * position.CostPrice
	positionGain := positionValue - positionCost

	value := &models.FundValue{
		FundCode:       fundCode,
		CurrentPrice:   fund.CurrentPrice,
		ChangeRate:     fund.ChangeRate,
		EstimateValue:  positionValue,
		EstimateChange: positionGain,
		Timestamp:      time.Now(),
	}

	return value, nil
}

// CalculateGroupValue 计算分组总价值
func (s *PortfolioService) CalculateGroupValue(groupName string) (*models.GroupSummary, error) {
	positions, err := s.GetGroupPositions(groupName)
	if err != nil {
		return nil, fmt.Errorf("failed to get group positions: %w", err)
	}

	var totalValue, totalCost, totalGain, dailyGain float64

	for _, position := range positions {
		// 获取基金最新价格
		fund, err := s.db.GetFund(position.FundCode)
		if err != nil {
			return nil, fmt.Errorf("failed to get fund %s: %w", position.FundCode, err)
		}
		if fund == nil {
			continue
		}

		// 计算持仓价值
		positionValue := position.Shares * fund.CurrentPrice
		positionCost := position.Shares * position.CostPrice
		positionGain := positionValue - positionCost

		totalValue += positionValue
		totalCost += positionCost
		totalGain += positionGain

		// 计算当日收益
		dailyGain += position.Shares * position.CostPrice * fund.ChangeRate / 100
	}

	summary := &models.GroupSummary{
		GroupName: groupName,
		Value:     totalValue,
		Cost:      totalCost,
		Gain:      totalGain,
		DailyGain: dailyGain,
	}

	if totalCost > 0 {
		summary.GainRate = totalGain / totalCost * 100
	}

	return summary, nil
}

// CalculateAllGroupsValue 计算所有分组的价值
func (s *PortfolioService) CalculateAllGroupsValue() (map[string]*models.GroupSummary, error) {
	groups, err := s.GetAllGroups()
	if err != nil {
		return nil, fmt.Errorf("failed to get groups: %w", err)
	}

	groupValues := make(map[string]*models.GroupSummary)

	for _, group := range groups {
		summary, err := s.CalculateGroupValue(group)
		if err != nil {
			return nil, fmt.Errorf("failed to calculate group value for %s: %w", group, err)
		}
		groupValues[group] = summary
	}

	return groupValues, nil
}

// GetPositionDetails 获取持仓详细信息
func (s *PortfolioService) GetPositionDetails(fundCode string) (map[string]interface{}, error) {
	details := make(map[string]interface{})

	// 获取持仓信息
	position, err := s.db.GetPosition(fundCode)
	if err != nil {
		return nil, fmt.Errorf("failed to get position: %w", err)
	}
	if position == nil {
		return nil, fmt.Errorf("position not found for fund %s", fundCode)
	}

	// 获取基金信息
	fund, err := s.db.GetFund(fundCode)
	if err != nil {
		return nil, fmt.Errorf("failed to get fund: %w", err)
	}

	// 计算持仓价值
	value, err := s.CalculatePositionValue(fundCode)
	if err != nil {
		return nil, fmt.Errorf("failed to calculate position value: %w", err)
	}

	// 获取历史数据
	history, err := s.db.GetValueHistory(fundCode, 30)
	if err != nil {
		return nil, fmt.Errorf("failed to get value history: %w", err)
	}

	details["position"] = position
	details["fund"] = fund
	details["value"] = value
	details["history"] = history

	return details, nil
}

// GetPortfolioPerformance 获取投资组合业绩表现
func (s *PortfolioService) GetPortfolioPerformance() (map[string]interface{}, error) {
	performance := make(map[string]interface{})

	// 获取投资组合概览
	summary, err := s.GetPortfolioSummary()
	if err != nil {
		return nil, fmt.Errorf("failed to get portfolio summary: %w", err)
	}

	// 获取所有分组
	groups, err := s.GetAllGroups()
	if err != nil {
		return nil, fmt.Errorf("failed to get groups: %w", err)
	}

	// 计算各分组表现
	groupPerformances := make(map[string]*models.GroupSummary)
	for _, group := range groups {
		groupSummary, err := s.GetGroupSummary(group)
		if err != nil {
			return nil, fmt.Errorf("failed to get group summary for %s: %w", group, err)
		}
		groupPerformances[group] = groupSummary
	}

	performance["summary"] = summary
	performance["groups"] = groupPerformances

	return performance, nil
}

// UpdatePositionShares 更新持仓份额
func (s *PortfolioService) UpdatePositionShares(fundCode string, shares float64) error {
	position, err := s.db.GetPosition(fundCode)
	if err != nil {
		return fmt.Errorf("failed to get position: %w", err)
	}
	if position == nil {
		return fmt.Errorf("position not found for fund %s", fundCode)
	}

	position.Shares = shares
	position.LastUpdated = time.Now()

	err = s.db.AddPosition(position)
	if err != nil {
		return fmt.Errorf("failed to update position shares: %w", err)
	}

	return nil
}

// UpdatePositionCostPrice 更新持仓成本价
func (s *PortfolioService) UpdatePositionCostPrice(fundCode string, costPrice float64) error {
	position, err := s.db.GetPosition(fundCode)
	if err != nil {
		return fmt.Errorf("failed to get position: %w", err)
	}
	if position == nil {
		return fmt.Errorf("position not found for fund %s", fundCode)
	}

	position.CostPrice = costPrice
	position.LastUpdated = time.Now()

	err = s.db.AddPosition(position)
	if err != nil {
		return fmt.Errorf("failed to update position cost price: %w", err)
	}

	return nil
}

// UpdatePositionGroup 更新持仓分组
func (s *PortfolioService) UpdatePositionGroup(fundCode string, groupName string) error {
	position, err := s.db.GetPosition(fundCode)
	if err != nil {
		return fmt.Errorf("failed to get position: %w", err)
	}
	if position == nil {
		return fmt.Errorf("position not found for fund %s", fundCode)
	}

	position.GroupName = groupName
	position.LastUpdated = time.Now()

	err = s.db.AddPosition(position)
	if err != nil {
		return fmt.Errorf("failed to update position group: %w", err)
	}

	return nil
}
