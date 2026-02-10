package handlers

import (
	"net/http"
	"strconv"

	"fundnet/backend/internal/services"

	"github.com/gin-gonic/gin"
)

// FundHandler 基金处理器
type FundHandler struct {
	fundService     *services.FundService
	estimateService *services.EstimateService
}

// NewFundHandler 创建基金处理器
func NewFundHandler(fundService *services.FundService, estimateService *services.EstimateService) *FundHandler {
	return &FundHandler{
		fundService:     fundService,
		estimateService: estimateService,
	}
}

// RegisterRoutes 注册路由
func RegisterRoutes(router *gin.Engine, fundService *services.FundService, estimateService *services.EstimateService) {
	handler := NewFundHandler(fundService, estimateService)

	api := router.Group("/api")
	{
		// 基金相关接口
		funds := api.Group("/funds")
		{
			funds.GET("", handler.GetFunds)
			funds.GET("/:code", handler.GetFund)
			funds.GET("/:code/estimate", handler.GetFundEstimate)
			funds.POST("", handler.AddFund)
			funds.DELETE("/:code", handler.RemoveFund)
			funds.PUT("/:code", handler.UpdateFund)
		}

		// 板块相关接口
		sectors := api.Group("/sectors")
		{
			sectors.GET("", handler.GetSectors)
			sectors.POST("", handler.CreateSector)
			sectors.PUT("/:id", handler.UpdateSector)
			sectors.DELETE("/:id", handler.DeleteSector)
		}

		// 持仓相关接口
		positions := api.Group("/positions")
		{
			positions.GET("", handler.GetPositions)
			positions.POST("", handler.AddPosition)
			positions.PUT("/:id", handler.UpdatePosition)
			positions.DELETE("/:id", handler.DeletePosition)
		}

		// 资产相关接口
		assets := api.Group("/assets")
		{
			assets.GET("", handler.GetAssets)
			assets.GET("/summary", handler.GetAssetSummary)
		}

		// 估算历史接口
		history := api.Group("/history")
		{
			history.GET("/:code", handler.GetFundHistory)
		}

		// 配置接口
		config := api.Group("/config")
		{
			config.GET("", handler.GetConfig)
			config.PUT("", handler.UpdateConfig)
		}
	}
}

// Response 统一响应结构
type Response struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

// GetFunds 获取基金列表
func (h *FundHandler) GetFunds(c *gin.Context) {
	funds, err := h.fundService.GetAllFunds()
	if err != nil {
		c.JSON(http.StatusInternalServerError, Response{
			Code:    500,
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, Response{
		Code:    0,
		Message: "success",
		Data:    funds,
	})
}

// GetFund 获取单个基金
func (h *FundHandler) GetFund(c *gin.Context) {
	code := c.Param("code")
	fund, err := h.fundService.GetFundByCode(code)
	if err != nil {
		c.JSON(http.StatusNotFound, Response{
			Code:    404,
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, Response{
		Code:    0,
		Message: "success",
		Data:    fund,
	})
}

// GetFundEstimate 获取基金估算
func (h *FundHandler) GetFundEstimate(c *gin.Context) {
	code := c.Param("code")
	estimate, err := h.estimateService.GetEstimate(code)
	if err != nil {
		c.JSON(http.StatusInternalServerError, Response{
			Code:    500,
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, Response{
		Code:    0,
		Message: "success",
		Data:    estimate,
	})
}

// AddFundRequest 添加基金请求
type AddFundRequest struct {
	Code   string `json:"code" binding:"required"`
	Name   string `json:"name"`
	Sector string `json:"sector"`
}

// AddFund 添加基金订阅
func (h *FundHandler) AddFund(c *gin.Context) {
	var req AddFundRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, Response{
			Code:    400,
			Message: err.Error(),
		})
		return
	}

	fund, err := h.fundService.AddFund(req.Code, req.Name, req.Sector)
	if err != nil {
		c.JSON(http.StatusInternalServerError, Response{
			Code:    500,
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, Response{
		Code:    0,
		Message: "success",
		Data:    fund,
	})
}

// RemoveFund 取消基金订阅
func (h *FundHandler) RemoveFund(c *gin.Context) {
	code := c.Param("code")
	if err := h.fundService.RemoveFund(code); err != nil {
		c.JSON(http.StatusInternalServerError, Response{
			Code:    500,
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, Response{
		Code:    0,
		Message: "success",
	})
}

// UpdateFundRequest 更新基金请求
type UpdateFundRequest struct {
	Name   string `json:"name"`
	Sector string `json:"sector"`
}

// UpdateFund 更新基金信息
func (h *FundHandler) UpdateFund(c *gin.Context) {
	code := c.Param("code")
	var req UpdateFundRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, Response{
			Code:    400,
			Message: err.Error(),
		})
		return
	}

	fund, err := h.fundService.UpdateFund(code, req.Name, req.Sector)
	if err != nil {
		c.JSON(http.StatusInternalServerError, Response{
			Code:    500,
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, Response{
		Code:    0,
		Message: "success",
		Data:    fund,
	})
}

// GetSectors 获取板块列表
func (h *FundHandler) GetSectors(c *gin.Context) {
	sectors, err := h.fundService.GetAllSectors()
	if err != nil {
		c.JSON(http.StatusInternalServerError, Response{
			Code:    500,
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, Response{
		Code:    0,
		Message: "success",
		Data:    sectors,
	})
}

// CreateSectorRequest 创建板块请求
type CreateSectorRequest struct {
	Name      string `json:"name" binding:"required"`
	Color     string `json:"color"`
	SortOrder int    `json:"sort_order"`
}

// CreateSector 创建板块
func (h *FundHandler) CreateSector(c *gin.Context) {
	var req CreateSectorRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, Response{
			Code:    400,
			Message: err.Error(),
		})
		return
	}

	sector, err := h.fundService.CreateSector(req.Name, req.Color, req.SortOrder)
	if err != nil {
		c.JSON(http.StatusInternalServerError, Response{
			Code:    500,
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, Response{
		Code:    0,
		Message: "success",
		Data:    sector,
	})
}

// UpdateSectorRequest 更新板块请求
type UpdateSectorRequest struct {
	Name      string `json:"name"`
	Color     string `json:"color"`
	SortOrder int    `json:"sort_order"`
}

// UpdateSector 更新板块
func (h *FundHandler) UpdateSector(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, Response{
			Code:    400,
			Message: "invalid sector id",
		})
		return
	}

	var req UpdateSectorRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, Response{
			Code:    400,
			Message: err.Error(),
		})
		return
	}

	sector, err := h.fundService.UpdateSector(id, req.Name, req.Color, req.SortOrder)
	if err != nil {
		c.JSON(http.StatusInternalServerError, Response{
			Code:    500,
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, Response{
		Code:    0,
		Message: "success",
		Data:    sector,
	})
}

// DeleteSector 删除板块
func (h *FundHandler) DeleteSector(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, Response{
			Code:    400,
			Message: "invalid sector id",
		})
		return
	}

	if err := h.fundService.DeleteSector(id); err != nil {
		c.JSON(http.StatusInternalServerError, Response{
			Code:    500,
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, Response{
		Code:    0,
		Message: "success",
	})
}

// GetPositions 获取持仓列表
func (h *FundHandler) GetPositions(c *gin.Context) {
	positions, err := h.fundService.GetAllPositions()
	if err != nil {
		c.JSON(http.StatusInternalServerError, Response{
			Code:    500,
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, Response{
		Code:    0,
		Message: "success",
		Data:    positions,
	})
}

// AddPositionRequest 添加持仓请求
type AddPositionRequest struct {
	FundCode string  `json:"fund_code" binding:"required"`
	FundName string  `json:"fund_name"`
	Shares   float64 `json:"shares" binding:"required"`
	Cost     float64 `json:"cost" binding:"required"`
	Sector   string  `json:"sector"`
}

// AddPosition 添加持仓
func (h *FundHandler) AddPosition(c *gin.Context) {
	var req AddPositionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, Response{
			Code:    400,
			Message: err.Error(),
		})
		return
	}

	position, err := h.fundService.AddPosition(req.FundCode, req.FundName, req.Shares, req.Cost, req.Sector)
	if err != nil {
		c.JSON(http.StatusInternalServerError, Response{
			Code:    500,
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, Response{
		Code:    0,
		Message: "success",
		Data:    position,
	})
}

// UpdatePositionRequest 更新持仓请求
type UpdatePositionRequest struct {
	Shares float64 `json:"shares"`
	Cost   float64 `json:"cost"`
	Sector string  `json:"sector"`
}

// UpdatePosition 更新持仓
func (h *FundHandler) UpdatePosition(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, Response{
			Code:    400,
			Message: "invalid position id",
		})
		return
	}

	var req UpdatePositionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, Response{
			Code:    400,
			Message: err.Error(),
		})
		return
	}

	position, err := h.fundService.UpdatePosition(id, req.Shares, req.Cost, req.Sector)
	if err != nil {
		c.JSON(http.StatusInternalServerError, Response{
			Code:    500,
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, Response{
		Code:    0,
		Message: "success",
		Data:    position,
	})
}

// DeletePosition 删除持仓
func (h *FundHandler) DeletePosition(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, Response{
			Code:    400,
			Message: "invalid position id",
		})
		return
	}

	if err := h.fundService.DeletePosition(id); err != nil {
		c.JSON(http.StatusInternalServerError, Response{
			Code:    500,
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, Response{
		Code:    0,
		Message: "success",
	})
}

// GetAssets 获取资产统计
func (h *FundHandler) GetAssets(c *gin.Context) {
	assets, err := h.fundService.GetAssetStats()
	if err != nil {
		c.JSON(http.StatusInternalServerError, Response{
			Code:    500,
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, Response{
		Code:    0,
		Message: "success",
		Data:    assets,
	})
}

// GetAssetSummary 获取资产摘要
func (h *FundHandler) GetAssetSummary(c *gin.Context) {
	summary, err := h.fundService.GetAssetSummary()
	if err != nil {
		c.JSON(http.StatusInternalServerError, Response{
			Code:    500,
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, Response{
		Code:    0,
		Message: "success",
		Data:    summary,
	})
}

// GetFundHistory 获取基金历史估算
func (h *FundHandler) GetFundHistory(c *gin.Context) {
	code := c.Param("code")
	days := c.DefaultQuery("days", "7")

	dayCount, err := strconv.Atoi(days)
	if err != nil {
		dayCount = 7
	}

	history, err := h.estimateService.GetHistory(code, dayCount)
	if err != nil {
		c.JSON(http.StatusInternalServerError, Response{
			Code:    500,
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, Response{
		Code:    0,
		Message: "success",
		Data:    history,
	})
}

// GetConfig 获取配置
func (h *FundHandler) GetConfig(c *gin.Context) {
	config := h.fundService.GetConfig()
	c.JSON(http.StatusOK, Response{
		Code:    0,
		Message: "success",
		Data:    config,
	})
}

// UpdateConfigRequest 更新配置请求
type UpdateConfigRequest struct {
	RefreshInterval int    `json:"refresh_interval"`
	LogLevel        string `json:"log_level"`
}

// UpdateConfig 更新配置
func (h *FundHandler) UpdateConfig(c *gin.Context) {
	var req UpdateConfigRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, Response{
			Code:    400,
			Message: err.Error(),
		})
		return
	}

	if err := h.fundService.UpdateConfig(req.RefreshInterval, req.LogLevel); err != nil {
		c.JSON(http.StatusInternalServerError, Response{
			Code:    500,
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, Response{
		Code:    0,
		Message: "success",
	})
}
