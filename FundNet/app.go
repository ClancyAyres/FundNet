package main

import (
	"context"
	"fmt"
	"time"

	"FundNet/models"
	"FundNet/services"
)

// App struct
type App struct {
	ctx       context.Context
	db        *services.DatabaseService
	fundData  *services.FundDataService
	portfolio *services.PortfolioService
	realTime  *services.RealTimeService
}

// NewApp creates a new App application struct
func NewApp() *App {
	return &App{}
}

// startup is called when the app starts. The context is saved
// so we can call the runtime methods
func (a *App) startup(ctx context.Context) {
	a.ctx = ctx

	// 初始化服务
	var err error
	a.db, err = services.NewDatabaseService("fundnet.db")
	if err != nil {
		fmt.Printf("Failed to initialize database: %v\n", err)
		return
	}

	a.fundData = services.NewFundDataService(a.db)
	a.portfolio = services.NewPortfolioService(a.db)
	a.realTime = services.NewRealTimeService(a.db)

	// 启动实时服务
	a.realTime.Start()
}

// shutdown is called when the app closes
func (a *App) shutdown(ctx context.Context) {
	// 停止实时服务
	if a.realTime != nil {
		a.realTime.Stop()
	}

	// 关闭数据库连接
	if a.db != nil {
		a.db.Close()
	}
}

// Greet returns a greeting for the given name
func (a *App) Greet(name string) string {
	return fmt.Sprintf("Hello %s, It's show time!", name)
}

// Database operations
func (a *App) AddFund(fund *models.Fund) error {
	return a.db.AddFund(fund)
}

func (a *App) GetFund(code string) (*models.Fund, error) {
	return a.db.GetFund(code)
}

func (a *App) GetAllFunds() ([]*models.Fund, error) {
	return a.db.GetAllFunds()
}

func (a *App) AddPosition(position *models.FundPosition) error {
	return a.portfolio.AddFundPosition(position)
}

func (a *App) GetPosition(fundCode string) (*models.FundPosition, error) {
	return a.db.GetPosition(fundCode)
}

func (a *App) GetAllPositions() ([]*models.FundPosition, error) {
	return a.db.GetAllPositions()
}

func (a *App) DeletePosition(fundCode string) error {
	return a.portfolio.DeleteFundPosition(fundCode)
}

// Fund data operations
func (a *App) GetFundInfo(code string) (*models.Fund, error) {
	return a.fundData.GetFundInfo(code)
}

func (a *App) CalculateEstimateValue(code string, shares, costPrice float64) (*models.FundValue, error) {
	return a.fundData.CalculateEstimateValue(code, shares, costPrice)
}

func (a *App) GetFundHistory(code string, days int) ([]*models.FundValue, error) {
	return a.fundData.GetFundHistory(code, days)
}

func (a *App) BatchUpdateFunds() error {
	return a.fundData.BatchUpdateFunds()
}

// Portfolio operations
func (a *App) GetPortfolioSummary() (*models.PortfolioSummary, error) {
	return a.portfolio.GetPortfolioSummary()
}

func (a *App) GetGroupSummary(groupName string) (*models.GroupSummary, error) {
	return a.portfolio.GetGroupSummary(groupName)
}

func (a *App) GetAllGroups() ([]string, error) {
	return a.portfolio.GetAllGroups()
}

func (a *App) GetGroupPositions(groupName string) ([]*models.FundPosition, error) {
	return a.portfolio.GetGroupPositions(groupName)
}

func (a *App) CalculatePositionValue(fundCode string) (*models.FundValue, error) {
	return a.portfolio.CalculatePositionValue(fundCode)
}

func (a *App) CalculateGroupValue(groupName string) (*models.GroupSummary, error) {
	return a.portfolio.CalculateGroupValue(groupName)
}

func (a *App) CalculateAllGroupsValue() (map[string]*models.GroupSummary, error) {
	return a.portfolio.CalculateAllGroupsValue()
}

func (a *App) GetPositionDetails(fundCode string) (map[string]interface{}, error) {
	return a.portfolio.GetPositionDetails(fundCode)
}

func (a *App) GetPortfolioPerformance() (map[string]interface{}, error) {
	return a.portfolio.GetPortfolioPerformance()
}

// Real-time operations
func (a *App) StartRealTimeService() {
	a.realTime.Start()
}

func (a *App) StopRealTimeService() {
	a.realTime.Stop()
}

func (a *App) GetRealTimeSummary() (*models.PortfolioSummary, error) {
	return a.realTime.GetRealTimeSummary()
}

func (a *App) GetRealTimeGroupSummary(groupName string) (*models.GroupSummary, error) {
	return a.realTime.GetRealTimeGroupSummary(groupName)
}

func (a *App) GetRealTimeFundValue(fundCode string) (*models.FundValue, error) {
	return a.realTime.GetRealTimeFundValue(fundCode)
}

func (a *App) GetRealTimeGroupValues() (map[string]*models.GroupSummary, error) {
	return a.realTime.GetRealTimeGroupValues()
}

func (a *App) GetRealTimeFundHistory(fundCode string, days int) ([]*models.FundValue, error) {
	return a.realTime.GetRealTimeFundHistory(fundCode, days)
}

func (a *App) GetClientCount() int {
	return a.realTime.GetClientCount()
}

func (a *App) IsRealTimeRunning() bool {
	return a.realTime.IsRunning()
}

func (a *App) SetRefreshInterval(interval int) {
	a.realTime.SetRefreshInterval(time.Duration(interval) * time.Second)
}

func (a *App) GetRefreshInterval() int {
	return int(a.realTime.GetRefreshInterval().Seconds())
}

// Configuration operations
func (a *App) GetConfig(key string) (string, error) {
	return a.db.GetConfig(key)
}

func (a *App) SetConfig(key, value string) error {
	return a.db.SetConfig(key, value)
}

// Utility operations
func (a *App) UpdatePositionShares(fundCode string, shares float64) error {
	return a.portfolio.UpdatePositionShares(fundCode, shares)
}

func (a *App) UpdatePositionCostPrice(fundCode string, costPrice float64) error {
	return a.portfolio.UpdatePositionCostPrice(fundCode, costPrice)
}

func (a *App) UpdatePositionGroup(fundCode string, groupName string) error {
	return a.portfolio.UpdatePositionGroup(fundCode, groupName)
}
