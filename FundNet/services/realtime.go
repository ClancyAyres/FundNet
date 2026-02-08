package services

import (
	"FundNet/models"
	"fmt"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

type RealTimeService struct {
	db              *DatabaseService
	fundData        *FundDataService
	portfolio       *PortfolioService
	clients         map[*websocket.Conn]bool
	broadcast       chan *models.FundValue
	mu              sync.Mutex
	isRunning       bool
	stopChan        chan bool
	refreshInterval time.Duration
}

type RealTimeMessage struct {
	Type      string      `json:"type"`
	Data      interface{} `json:"data"`
	Timestamp time.Time   `json:"timestamp"`
}

func NewRealTimeService(db *DatabaseService) *RealTimeService {
	return &RealTimeService{
		db:              db,
		fundData:        NewFundDataService(db),
		portfolio:       NewPortfolioService(db),
		clients:         make(map[*websocket.Conn]bool),
		broadcast:       make(chan *models.FundValue, 100),
		refreshInterval: 30 * time.Second,
	}
}

// Start 启动实时服务
func (s *RealTimeService) Start() {
	if s.isRunning {
		return
	}

	s.isRunning = true
	s.stopChan = make(chan bool)

	// 启动广播协程
	go s.broadcastLoop()

	// 启动数据更新协程
	go s.updateLoop()
}

// Stop 停止实时服务
func (s *RealTimeService) Stop() {
	if !s.isRunning {
		return
	}

	s.isRunning = false
	close(s.stopChan)

	// 关闭所有客户端连接
	s.mu.Lock()
	for client := range s.clients {
		client.Close()
	}
	s.clients = make(map[*websocket.Conn]bool)
	s.mu.Unlock()
}

// AddClient 添加WebSocket客户端
func (s *RealTimeService) AddClient(conn *websocket.Conn) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.clients[conn] = true
}

// RemoveClient 移除WebSocket客户端
func (s *RealTimeService) RemoveClient(conn *websocket.Conn) {
	s.mu.Lock()
	defer s.mu.Unlock()

	delete(s.clients, conn)
}

// broadcastLoop 广播循环
func (s *RealTimeService) broadcastLoop() {
	for {
		select {
		case <-s.stopChan:
			return
		case fundValue := <-s.broadcast:
			s.mu.Lock()
			for client := range s.clients {
				err := client.WriteJSON(&RealTimeMessage{
					Type:      "fund_update",
					Data:      fundValue,
					Timestamp: time.Now(),
				})
				if err != nil {
					// 客户端连接出错，移除客户端
					client.Close()
					delete(s.clients, client)
				}
			}
			s.mu.Unlock()
		}
	}
}

// updateLoop 数据更新循环
func (s *RealTimeService) updateLoop() {
	ticker := time.NewTicker(s.refreshInterval)
	defer ticker.Stop()

	for {
		select {
		case <-s.stopChan:
			return
		case <-ticker.C:
			if s.isRunning {
				s.updateAllFunds()
			}
		}
	}
}

// updateAllFunds 更新所有基金数据
func (s *RealTimeService) updateAllFunds() {
	// 获取所有持仓的基金代码
	positions, err := s.db.GetAllPositions()
	if err != nil {
		return
	}

	// 获取所有基金代码
	fundCodes := make(map[string]bool)
	for _, position := range positions {
		fundCodes[position.FundCode] = true
	}

	// 更新每个基金的信息
	for code := range fundCodes {
		fund, err := s.fundData.GetFundInfo(code)
		if err != nil {
			continue
		}

		// 计算估值
		position, err := s.db.GetPosition(code)
		if err != nil || position == nil {
			continue
		}

		estimateValue := position.Shares * fund.CurrentPrice
		estimateChange := estimateValue - (position.Shares * position.CostPrice)

		value := &models.FundValue{
			FundCode:       code,
			CurrentPrice:   fund.CurrentPrice,
			ChangeRate:     fund.ChangeRate,
			EstimateValue:  estimateValue,
			EstimateChange: estimateChange,
			Timestamp:      time.Now(),
		}

		// 广播更新
		s.broadcast <- value

		// 保存到历史记录
		s.db.AddValueHistory(value)
	}

	// 广播投资组合概览
	summary, err := s.portfolio.GetPortfolioSummary()
	if err == nil {
		s.broadcast <- &models.FundValue{
			FundCode:       "portfolio",
			CurrentPrice:   summary.TotalValue,
			ChangeRate:     summary.TotalGainRate,
			EstimateValue:  summary.TotalValue,
			EstimateChange: summary.TotalGain,
			Timestamp:      time.Now(),
		}
	}
}

// SetRefreshInterval 设置刷新间隔
func (s *RealTimeService) SetRefreshInterval(interval time.Duration) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.refreshInterval = interval
}

// GetRefreshInterval 获取刷新间隔
func (s *RealTimeService) GetRefreshInterval() time.Duration {
	s.mu.Lock()
	defer s.mu.Unlock()

	return s.refreshInterval
}

// GetRealTimeSummary 获取实时概览
func (s *RealTimeService) GetRealTimeSummary() (*models.PortfolioSummary, error) {
	return s.portfolio.GetPortfolioSummary()
}

// GetRealTimeGroupSummary 获取实时分组概览
func (s *RealTimeService) GetRealTimeGroupSummary(groupName string) (*models.GroupSummary, error) {
	return s.portfolio.GetGroupSummary(groupName)
}

// GetRealTimeFundValue 获取实时基金估值
func (s *RealTimeService) GetRealTimeFundValue(fundCode string) (*models.FundValue, error) {
	// 获取基金最新信息
	fund, err := s.fundData.GetFundInfo(fundCode)
	if err != nil {
		return nil, fmt.Errorf("failed to get fund info: %w", err)
	}

	// 获取持仓信息
	position, err := s.db.GetPosition(fundCode)
	if err != nil || position == nil {
		return nil, fmt.Errorf("position not found for fund %s", fundCode)
	}

	// 计算估值
	estimateValue := position.Shares * fund.CurrentPrice
	estimateChange := estimateValue - (position.Shares * position.CostPrice)

	value := &models.FundValue{
		FundCode:       fundCode,
		CurrentPrice:   fund.CurrentPrice,
		ChangeRate:     fund.ChangeRate,
		EstimateValue:  estimateValue,
		EstimateChange: estimateChange,
		Timestamp:      time.Now(),
	}

	return value, nil
}

// GetRealTimeGroupValues 获取实时分组价值
func (s *RealTimeService) GetRealTimeGroupValues() (map[string]*models.GroupSummary, error) {
	return s.portfolio.CalculateAllGroupsValue()
}

// GetRealTimeFundHistory 获取实时基金历史数据
func (s *RealTimeService) GetRealTimeFundHistory(fundCode string, days int) ([]*models.FundValue, error) {
	return s.db.GetValueHistory(fundCode, days)
}

// GetClientCount 获取客户端连接数
func (s *RealTimeService) GetClientCount() int {
	s.mu.Lock()
	defer s.mu.Unlock()

	return len(s.clients)
}

// IsRunning 检查服务是否正在运行
func (s *RealTimeService) IsRunning() bool {
	s.mu.Lock()
	defer s.mu.Unlock()

	return s.isRunning
}

// SendCustomMessage 发送自定义消息给所有客户端
func (s *RealTimeService) SendCustomMessage(messageType string, data interface{}) {
	s.mu.Lock()
	defer s.mu.Unlock()

	for client := range s.clients {
		err := client.WriteJSON(&RealTimeMessage{
			Type:      messageType,
			Data:      data,
			Timestamp: time.Now(),
		})
		if err != nil {
			// 客户端连接出错，移除客户端
			client.Close()
			delete(s.clients, client)
		}
	}
}

// SendFundAlert 发送基金提醒
func (s *RealTimeService) SendFundAlert(fundCode string, alertType string, value float64) {
	alert := map[string]interface{}{
		"fundCode":  fundCode,
		"alertType": alertType,
		"value":     value,
		"timestamp": time.Now(),
	}

	s.SendCustomMessage("fund_alert", alert)
}

// SendPortfolioAlert 发送投资组合提醒
func (s *RealTimeService) SendPortfolioAlert(alertType string, value float64) {
	alert := map[string]interface{}{
		"alertType": alertType,
		"value":     value,
		"timestamp": time.Now(),
	}

	s.SendCustomMessage("portfolio_alert", alert)
}
