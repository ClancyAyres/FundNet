package services

import (
	"FundNet/models"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/go-resty/resty/v2"
	"github.com/sirupsen/logrus"
)

type FundDataService struct {
	client *resty.Client
	db     *DatabaseService
}

type TianTianFundResponse struct {
	Data struct {
		KfData [][]string `json:"kfData"`
	} `json:"data"`
}

type FundInfoResponse struct {
	Data struct {
		Name       string  `json:"name"`
		NetWorth   float64 `json:"netWorth"`
		ChangeRate float64 `json:"changeRate"`
	} `json:"data"`
}

func NewFundDataService(db *DatabaseService) *FundDataService {
	client := resty.New()
	client.SetTimeout(10 * time.Second)
	client.SetHeaders(map[string]string{
		"User-Agent": "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.124 Safari/537.36",
		"Referer":    "https://fund.eastmoney.com/",
	})

	return &FundDataService{
		client: client,
		db:     db,
	}
}

// GetFundInfo 从天天基金获取基金信息
func (s *FundDataService) GetFundInfo(code string) (*models.Fund, error) {
	// 天天基金网的API接口
	url := fmt.Sprintf("https://fundgz.1234567.com.cn/js/%s.js", code)

	resp, err := s.client.R().
		SetHeader("Accept", "*/*").
		SetHeader("Accept-Encoding", "gzip, deflate, br").
		Get(url)

	if err != nil {
		return nil, fmt.Errorf("failed to fetch fund info: %w", err)
	}

	// 处理返回的JSONP格式数据
	body := string(resp.Body())
	// 提取JSON数据
	start := strings.Index(body, "(")
	end := strings.LastIndex(body, ")")
	if start == -1 || end == -1 {
		return nil, fmt.Errorf("invalid response format")
	}

	jsonStr := body[start+1 : end]
	var fundInfo FundInfoResponse

	err = json.Unmarshal([]byte(jsonStr), &fundInfo)
	if err != nil {
		return nil, fmt.Errorf("failed to parse fund info: %w", err)
	}

	fund := &models.Fund{
		Code:         code,
		Name:         fundInfo.Data.Name,
		CurrentPrice: fundInfo.Data.NetWorth,
		ChangeRate:   fundInfo.Data.ChangeRate,
		LastUpdated:  time.Now(),
	}

	// 保存到数据库
	err = s.db.AddFund(fund)
	if err != nil {
		return nil, fmt.Errorf("failed to save fund to database: %w", err)
	}

	return fund, nil
}

// GetFundHistory 获取基金历史净值数据
func (s *FundDataService) GetFundHistory(code string, days int) ([]*models.FundValue, error) {
	// 天天基金网的历史净值接口
	url := fmt.Sprintf("https://api.fund.eastmoney.com/f10/lsjz?fundCode=%s&pageIndex=1&pageSize=%d&startDate=&endDate=&callback=jQuery183003459761034742444_1625000000000&_=1625000000000", code, days)

	resp, err := s.client.R().
		SetHeader("Accept", "*/*").
		SetHeader("Accept-Encoding", "gzip, deflate, br").
		Get(url)

	if err != nil {
		return nil, fmt.Errorf("failed to fetch fund history: %w", err)
	}

	// 处理返回的JSONP格式数据
	body := string(resp.Body())
	start := strings.Index(body, "(")
	end := strings.LastIndex(body, ")")
	if start == -1 || end == -1 {
		return nil, fmt.Errorf("invalid response format")
	}

	jsonStr := body[start+1 : end]

	var historyResponse struct {
		Data struct {
			LSJZList []struct {
				FundCode    string  `json:"FSRQ"`
				NetWorth    float64 `json:"DWJZ"`
				ChangeRate  float64 `json:"JZZZL"`
				Accumulated float64 `json:"LJJZ"`
			} `json:"lsjzList"`
		} `json:"Data"`
	}

	err = json.Unmarshal([]byte(jsonStr), &historyResponse)
	if err != nil {
		return nil, fmt.Errorf("failed to parse fund history: %w", err)
	}

	var history []*models.FundValue
	for _, item := range historyResponse.Data.LSJZList {
		timestamp, err := time.Parse("2006-01-02", item.FundCode)
		if err != nil {
			continue
		}

		value := &models.FundValue{
			FundCode:       code,
			CurrentPrice:   item.NetWorth,
			ChangeRate:     item.ChangeRate,
			EstimateValue:  item.NetWorth,
			EstimateChange: item.ChangeRate,
			Timestamp:      timestamp,
		}

		history = append(history, value)
	}

	return history, nil
}

// CalculateEstimateValue 计算基金估值
func (s *FundDataService) CalculateEstimateValue(code string, shares, costPrice float64) (*models.FundValue, error) {
	// 获取基金最新信息
	fund, err := s.GetFundInfo(code)
	if err != nil {
		return nil, fmt.Errorf("failed to get fund info: %w", err)
	}

	// 计算估值
	estimateValue := shares * fund.CurrentPrice
	estimateChange := estimateValue - (shares * costPrice)

	value := &models.FundValue{
		FundCode:       code,
		CurrentPrice:   fund.CurrentPrice,
		ChangeRate:     fund.ChangeRate,
		EstimateValue:  estimateValue,
		EstimateChange: estimateChange,
		Timestamp:      time.Now(),
	}

	// 保存到历史记录
	err = s.db.AddValueHistory(value)
	if err != nil {
		return nil, fmt.Errorf("failed to save value history: %w", err)
	}

	return value, nil
}

// BatchUpdateFunds 批量更新基金信息
func (s *FundDataService) BatchUpdateFunds() error {
	// 获取所有持仓的基金代码
	positions, err := s.db.GetAllPositions()
	if err != nil {
		return fmt.Errorf("failed to get positions: %w", err)
	}

	// 获取所有基金代码
	fundCodes := make(map[string]bool)
	for _, position := range positions {
		fundCodes[position.FundCode] = true
	}

	// 更新每个基金的信息
	for code := range fundCodes {
		_, err := s.GetFundInfo(code)
		if err != nil {
			logrus.Warnf("Failed to update fund %s: %v", code, err)
			continue
		}
	}

	return nil
}

// GetRealTimeData 获取实时数据（模拟）
func (s *FundDataService) GetRealTimeData(code string) (*models.Fund, error) {
	// 这里可以实现更复杂的实时数据获取逻辑
	// 目前简单地调用 GetFundInfo
	return s.GetFundInfo(code)
}

// GetFundList 获取热门基金列表
func (s *FundDataService) GetFundList() ([]*models.Fund, error) {
	// 天天基金网的热门基金接口
	url := "https://fund.eastmoney.com/data/rankhandler.aspx?op=ph&dt=kf&ft=all&rs=&gs=0&sc=6pn=50&pi=1&lx=0&v=0.2639394374594143"

	_, err := s.client.R().
		SetHeader("Accept", "application/json, text/javascript, */*, q=0.01").
		SetHeader("X-Requested-With", "XMLHttpRequest").
		Get(url)

	if err != nil {
		return nil, fmt.Errorf("failed to fetch fund list: %w", err)
	}

	// 处理返回的数据
	// body := string(resp.Body())

	// 这里需要解析返回的文本格式数据
	// 由于格式比较复杂，这里简化处理
	funds := []*models.Fund{
		{
			Code:         "000001",
			Name:         "华夏成长混合",
			CurrentPrice: 1.2345,
			ChangeRate:   0.5,
			LastUpdated:  time.Now(),
		},
		{
			Code:         "000002",
			Name:         "易方达消费行业股票",
			CurrentPrice: 2.3456,
			ChangeRate:   -0.2,
			LastUpdated:  time.Now(),
		},
	}

	return funds, nil
}

// DownloadFundData 下载基金数据到本地
func (s *FundDataService) DownloadFundData(code string, days int) error {
	history, err := s.GetFundHistory(code, days)
	if err != nil {
		return fmt.Errorf("failed to download fund data: %w", err)
	}

	// 保存历史数据到数据库
	for _, value := range history {
		err := s.db.AddValueHistory(value)
		if err != nil {
			logrus.Warnf("Failed to save history for fund %s: %v", code, err)
		}
	}

	return nil
}

// GetFundPerformance 获取基金业绩表现
func (s *FundDataService) GetFundPerformance(code string) (map[string]float64, error) {
	// 获取历史数据
	history, err := s.GetFundHistory(code, 365)
	if err != nil {
		return nil, fmt.Errorf("failed to get fund performance: %w", err)
	}

	if len(history) == 0 {
		return nil, fmt.Errorf("no history data available")
	}

	// 计算各种时间段的收益率
	performance := make(map[string]float64)

	// 计算最新净值
	latestValue := history[0].CurrentPrice

	// 计算各时间段的收益率
	timePeriods := map[string]int{
		"1_week":   7,
		"1_month":  30,
		"3_months": 90,
		"6_months": 180,
		"1_year":   365,
	}

	for period, days := range timePeriods {
		if len(history) > days {
			periodValue := history[days].CurrentPrice
			if periodValue > 0 {
				returnRate := (latestValue - periodValue) / periodValue * 100
				performance[period] = returnRate
			}
		}
	}

	return performance, nil
}

// GetFundDetails 获取基金详细信息
func (s *FundDataService) GetFundDetails(code string) (map[string]interface{}, error) {
	details := make(map[string]interface{})

	// 获取基本信息
	fund, err := s.GetFundInfo(code)
	if err != nil {
		return nil, fmt.Errorf("failed to get fund details: %w", err)
	}

	// 获取业绩表现
	performance, err := s.GetFundPerformance(code)
	if err != nil {
		logrus.Warnf("Failed to get performance for fund %s: %v", code, err)
	}

	// 获取持仓信息（如果有）
	position, err := s.db.GetPosition(code)
	if err != nil && err.Error() != "sql: no rows in result set" {
		return nil, fmt.Errorf("failed to get position: %w", err)
	}

	details["fund_info"] = fund
	details["performance"] = performance
	if position != nil {
		details["position"] = position
	}

	return details, nil
}
