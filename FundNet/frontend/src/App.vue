<template>
  <div id="app">
    <div class="app-container">
      <!-- 侧边栏 -->
      <aside class="sidebar">
        <div class="sidebar-header">
          <h1>FundNet</h1>
          <p>基金净值估算</p>
        </div>
        
        <nav class="sidebar-nav">
          <button 
            v-for="tab in tabs" 
            :key="tab.id"
            :class="['nav-item', { active: currentTab === tab.id }]"
            @click="switchTab(tab.id)"
          >
            <span class="nav-icon">{{ tab.icon }}</span>
            <span class="nav-text">{{ tab.name }}</span>
          </button>
        </nav>
        
        <div class="sidebar-footer">
          <div class="status-indicator">
            <span :class="['status-dot', { online: isOnline, offline: !isOnline }]"></span>
            <span class="status-text">{{ isOnline ? '已连接' : '未连接' }}</span>
          </div>
          <div class="refresh-status">
            <span class="refresh-text">最后更新: {{ lastUpdateTime }}</span>
          </div>
        </div>
      </aside>

      <!-- 主内容区 -->
      <main class="main-content">
        <!-- 顶部工具栏 -->
        <header class="toolbar">
          <div class="toolbar-left">
            <h2>{{ currentTabName }}</h2>
            <div class="toolbar-actions">
              <button @click="refreshData" :disabled="isRefreshing">
                {{ isRefreshing ? '刷新中...' : '刷新数据' }}
              </button>
              <button @click="exportData">导出数据</button>
              <button @click="openSettings">设置</button>
            </div>
          </div>
          
          <div class="toolbar-right">
            <div class="search-box">
              <input 
                v-model="globalSearch" 
                type="text" 
                placeholder="全局搜索..."
                @input="performGlobalSearch"
              />
            </div>
            <div class="quick-stats">
              <div class="stat-item">
                <span class="stat-label">总资产</span>
                <span class="stat-value">{{ formatCurrency(portfolioSummary.totalValue) }}</span>
              </div>
              <div class="stat-item">
                <span class="stat-label">总收益</span>
                <span :class="['stat-value', portfolioSummary.totalGain >= 0 ? 'positive' : 'negative']">
                  {{ portfolioSummary.totalGain >= 0 ? '+' : '' }}{{ formatCurrency(portfolioSummary.totalGain) }}
                </span>
              </div>
            </div>
          </div>
        </header>

        <!-- 内容区域 -->
        <div class="content-area">
          <!-- 基金列表页面 -->
          <div v-if="currentTab === 'funds'" class="funds-page">
            <FundList 
              :funds="funds"
              :positions="positions"
              @select-fund="selectFund"
              @add-to-portfolio="addPosition"
              @view-details="viewFundDetails"
              @refresh-funds="refreshFunds"
            />
          </div>

          <!-- 投资组合页面 -->
          <div v-else-if="currentTab === 'portfolio'" class="portfolio-page">
            <Portfolio 
              :portfolio-summary="portfolioSummary"
              :positions="positions"
              :groups="groups"
              :group-values="groupValues"
              @add-group="addGroup"
              @rename-group="renameGroup"
              @delete-group="deleteGroup"
              @add-position="addPosition"
              @remove-position="removePosition"
              @refresh-portfolio="refreshPortfolio"
              @export-data="exportPortfolioData"
            />
          </div>

          <!-- 实时监控页面 -->
          <div v-else-if="currentTab === 'monitoring'" class="monitoring-page">
            <div class="monitoring-header">
              <h3>实时监控</h3>
              <div class="monitoring-controls">
                <button @click="toggleRealTime" :class="{ active: isRealTimeRunning }">
                  {{ isRealTimeRunning ? '停止监控' : '开始监控' }}
                </button>
                <select v-model="selectedTimeRange" @change="updateCharts">
                  <option value="7">7天</option>
                  <option value="30" selected>30天</option>
                  <option value="90">90天</option>
                  <option value="365">1年</option>
                </select>
              </div>
            </div>
            
            <div class="charts-grid">
              <Chart 
                v-for="chart in charts" 
                :key="chart.id"
                :title="chart.title"
                :data="chart.data"
                :type="chart.type"
              />
            </div>
            
            <div class="monitoring-stats">
              <div class="stat-card">
                <h4>投资组合表现</h4>
                <div class="stat-value">{{ formatCurrency(portfolioSummary.totalValue) }}</div>
                <div :class="['stat-change', portfolioSummary.totalGain >= 0 ? 'positive' : 'negative']">
                  {{ portfolioSummary.totalGain >= 0 ? '+' : '' }}{{ formatCurrency(portfolioSummary.totalGain) }}
                  ({{ portfolioSummary.totalGainRate >= 0 ? '+' : '' }}{{ portfolioSummary.totalGainRate.toFixed(2) }}%)
                </div>
              </div>
              
              <div class="stat-card">
                <h4>当日表现</h4>
                <div :class="['stat-value', portfolioSummary.dailyGain >= 0 ? 'positive' : 'negative']">
                  {{ portfolioSummary.dailyGain >= 0 ? '+' : '' }}{{ formatCurrency(portfolioSummary.dailyGain) }}
                </div>
                <div :class="['stat-change', portfolioSummary.dailyGainRate >= 0 ? 'positive' : 'negative']">
                  ({{ portfolioSummary.dailyGainRate >= 0 ? '+' : '' }}{{ portfolioSummary.dailyGainRate.toFixed(2) }}%)
                </div>
              </div>
            </div>
          </div>

          <!-- 设置页面 -->
          <div v-else-if="currentTab === 'settings'" class="settings-page">
            <Settings 
              @settings-updated="updateAppSettings"
              @data-restored="handleDataRestored"
              @data-cleared="handleDataCleared"
            />
          </div>
        </div>
      </main>
    </div>

    <!-- 加载遮罩 -->
    <div v-if="isLoading" class="loading-overlay">
      <div class="loading-spinner"></div>
      <p>{{ loadingMessage }}</p>
    </div>
  </div>
</template>

<script>
import FundList from './components/FundList.vue'
import Portfolio from './components/Portfolio.vue'
import Chart from './components/Chart.vue'
import Settings from './components/Settings.vue'

export default {
  name: 'App',
  components: {
    FundList,
    Portfolio,
    Chart,
    Settings
  },
  data() {
    return {
      currentTab: 'portfolio',
      isLoading: false,
      loadingMessage: '加载中...',
      isRefreshing: false,
      isOnline: true,
      lastUpdateTime: new Date().toLocaleString('zh-CN'),
      globalSearch: '',
      
      // 数据状态
      funds: [],
      positions: [],
      groups: ['全部', '科技', '医疗', '新能源', 'QDII'],
      portfolioSummary: {
        totalValue: 0,
        totalCost: 0,
        totalGain: 0,
        totalGainRate: 0,
        dailyGain: 0,
        dailyGainRate: 0
      },
      groupValues: {},
      charts: [
        { id: 'portfolio', title: '投资组合净值走势', data: [], type: 'line' },
        { id: 'daily', title: '每日收益变化', data: [], type: 'bar' }
      ],
      
      // 监控设置
      isRealTimeRunning: false,
      selectedTimeRange: 30,
      
      // 标签页配置
      tabs: [
        { id: 'funds', name: '基金列表', icon: '📊' },
        { id: 'portfolio', name: '投资组合', icon: '💼' },
        { id: 'monitoring', name: '实时监控', icon: '📈' },
        { id: 'settings', name: '设置', icon: '⚙️' }
      ]
    }
  },
  computed: {
    currentTabName() {
      const tab = this.tabs.find(t => t.id === this.currentTab)
      return tab ? tab.name : 'FundNet'
    }
  },
  async mounted() {
    await this.initializeApp()
    this.startPeriodicUpdate()
  },
  beforeDestroy() {
    this.stopPeriodicUpdate()
  },
  methods: {
    async initializeApp() {
      this.isLoading = true
      this.loadingMessage = '初始化应用...'
      
      try {
        // 初始化后端服务
        await this.$wails.App.startup()
        
        // 加载数据
        await this.loadInitialData()
        
        // 启动实时监控
        this.$wails.App.StartRealTimeService()
        this.isRealTimeRunning = true
        
        this.isLoading = false
      } catch (error) {
        console.error('应用初始化失败:', error)
        this.isLoading = false
        alert('应用初始化失败，请重试')
      }
    },
    
    async loadInitialData() {
      try {
        // 加载基金数据
        this.funds = await this.$wails.App.GetAllFunds()
        
        // 加载持仓数据
        this.positions = await this.$wails.App.GetAllPositions()
        
        // 加载投资组合概览
        this.portfolioSummary = await this.$wails.App.GetPortfolioSummary()
        
        // 加载分组数据
        this.groups = await this.$wails.App.GetAllGroups()
        
        // 加载分组价值
        this.groupValues = await this.$wails.App.GetRealTimeGroupValues()
        
        // 更新图表数据
        this.updateCharts()
        
        this.lastUpdateTime = new Date().toLocaleString('zh-CN')
      } catch (error) {
        console.error('加载数据失败:', error)
      }
    },
    
    switchTab(tabId) {
      this.currentTab = tabId
    },
    
    async refreshData() {
      this.isRefreshing = true
      this.loadingMessage = '刷新数据中...'
      
      try {
        await this.loadInitialData()
        this.isRefreshing = false
      } catch (error) {
        console.error('刷新数据失败:', error)
        this.isRefreshing = false
        alert('刷新数据失败，请重试')
      }
    },
    
    async refreshFunds() {
      try {
        await this.$wails.App.BatchUpdateFunds()
        await this.loadInitialData()
      } catch (error) {
        console.error('刷新基金数据失败:', error)
      }
    },
    
    async refreshPortfolio() {
      try {
        this.portfolioSummary = await this.$wails.App.GetPortfolioSummary()
        this.groupValues = await this.$wails.App.GetRealTimeGroupValues()
      } catch (error) {
        console.error('刷新投资组合失败:', error)
      }
    },
    
    async addPosition(position) {
      try {
        await this.$wails.App.AddPosition(position)
        await this.loadInitialData()
        alert('持仓添加成功！')
      } catch (error) {
        console.error('添加持仓失败:', error)
        alert('添加持仓失败，请重试')
      }
    },
    
    async removePosition(fundCode) {
      try {
        await this.$wails.App.DeletePosition(fundCode)
        await this.loadInitialData()
        alert('持仓删除成功！')
      } catch (error) {
        console.error('删除持仓失败:', error)
        alert('删除持仓失败，请重试')
      }
    },
    
    async addGroup(groupName) {
      // 添加分组逻辑
      this.groups.push(groupName)
    },
    
    async renameGroup(oldName, newName) {
      // 重命名分组逻辑
      const index = this.groups.indexOf(oldName)
      if (index > -1) {
        this.groups[index] = newName
      }
    },
    
    async deleteGroup(groupName) {
      // 删除分组逻辑
      this.groups = this.groups.filter(g => g !== groupName)
    },
    
    async exportData() {
      try {
        const result = await this.$wails.App.ExportData()
        if (result) {
          alert('数据导出成功！')
        } else {
          alert('数据导出失败！')
        }
      } catch (error) {
        console.error('导出数据失败:', error)
        alert('导出数据失败，请重试')
      }
    },
    
    async exportPortfolioData() {
      await this.exportData()
    },
    
    async selectFund(fund) {
      // 选择基金逻辑
      console.log('选择基金:', fund)
    },
    
    async viewFundDetails(fund) {
      // 查看基金详情逻辑
      console.log('查看基金详情:', fund)
    },
    
    async openSettings() {
      this.currentTab = 'settings'
    },
    
    async updateAppSettings(settings) {
      // 更新应用设置
      console.log('更新应用设置:', settings)
    },
    
    async handleDataRestored() {
      await this.loadInitialData()
    },
    
    async handleDataCleared() {
      this.funds = []
      this.positions = []
      this.portfolioSummary = {
        totalValue: 0,
        totalCost: 0,
        totalGain: 0,
        totalGainRate: 0,
        dailyGain: 0,
        dailyGainRate: 0
      }
    },
    
    async toggleRealTime() {
      if (this.isRealTimeRunning) {
        this.$wails.App.StopRealTimeService()
        this.isRealTimeRunning = false
      } else {
        this.$wails.App.StartRealTimeService()
        this.isRealTimeRunning = true
      }
    },
    
    updateCharts() {
      // 更新图表数据
      // 这里可以根据 selectedTimeRange 获取历史数据并更新 charts
      console.log('更新图表数据，时间范围:', this.selectedTimeRange)
    },
    
    performGlobalSearch() {
      // 全局搜索逻辑
      console.log('全局搜索:', this.globalSearch)
    },
    
    formatCurrency(value) {
      return '¥' + (value || 0).toLocaleString('zh-CN', {
        minimumFractionDigits: 2,
        maximumFractionDigits: 2
      })
    },
    
    startPeriodicUpdate() {
      // 每30秒更新一次数据
      this.updateInterval = setInterval(async () => {
        try {
          this.portfolioSummary = await this.$wails.App.GetRealTimeSummary()
          this.lastUpdateTime = new Date().toLocaleString('zh-CN')
        } catch (error) {
          console.error('定期更新失败:', error)
        }
      }, 30000)
    },
    
    stopPeriodicUpdate() {
      if (this.updateInterval) {
        clearInterval(this.updateInterval)
      }
    }
  }
}
</script>

<style>
#app {
  height: 100vh;
  display: flex;
  flex-direction: column;
  background-color: #f5f5f5;
  font-family: -apple-system, BlinkMacSystemFont, "Segoe UI", Roboto, Helvetica, Arial, sans-serif;
}

.app-container {
  display: flex;
  flex: 1;
  height: 100%;
}

/* 侧边栏样式 */
.sidebar {
  width: 250px;
  background: #2c3e50;
  color: white;
  display: flex;
  flex-direction: column;
  padding: 20px;
}

.sidebar-header {
  margin-bottom: 30px;
}

.sidebar-header h1 {
  margin: 0 0 5px 0;
  font-size: 24px;
  font-weight: bold;
}

.sidebar-header p {
  margin: 0;
  font-size: 12px;
  opacity: 0.8;
}

.sidebar-nav {
  flex: 1;
}

.nav-item {
  width: 100%;
  padding: 12px 16px;
  background: transparent;
  border: none;
  color: white;
  text-align: left;
  cursor: pointer;
  border-radius: 6px;
  margin-bottom: 8px;
  display: flex;
  align-items: center;
  gap: 12px;
  transition: background-color 0.3s ease;
  font-size: 14px;
}

.nav-item:hover {
  background-color: rgba(255, 255, 255, 0.1);
}

.nav-item.active {
  background-color: rgba(255, 255, 255, 0.2);
  font-weight: bold;
}

.nav-icon {
  font-size: 18px;
}

.sidebar-footer {
  margin-top: auto;
  padding-top: 20px;
  border-top: 1px solid rgba(255, 255, 255, 0.1);
}

.status-indicator {
  display: flex;
  align-items: center;
  gap: 8px;
  margin-bottom: 10px;
  font-size: 12px;
}

.status-dot {
  width: 8px;
  height: 8px;
  border-radius: 50%;
  display: inline-block;
}

.status-dot.online {
  background-color: #28a745;
  box-shadow: 0 0 5px #28a745;
}

.status-dot.offline {
  background-color: #dc3545;
}

.status-text {
  font-weight: 500;
}

.refresh-status {
  font-size: 11px;
  opacity: 0.8;
}

/* 主内容区样式 */
.main-content {
  flex: 1;
  display: flex;
  flex-direction: column;
  overflow: hidden;
}

.toolbar {
  background: white;
  border-bottom: 1px solid #eee;
  padding: 15px 20px;
  display: flex;
  justify-content: space-between;
  align-items: center;
  height: 60px;
}

.toolbar-left {
  display: flex;
  align-items: center;
  gap: 20px;
}

.toolbar-left h2 {
  margin: 0;
  color: #333;
  font-size: 18px;
}

.toolbar-actions {
  display: flex;
  gap: 10px;
}

.toolbar-actions button {
  padding: 6px 12px;
  border: 1px solid #ddd;
  border-radius: 4px;
  cursor: pointer;
  font-size: 12px;
  background-color: white;
}

.toolbar-actions button:hover {
  background-color: #f8f9fa;
}

.toolbar-actions button:first-child {
  background-color: #007bff;
  color: white;
  border-color: #007bff;
}

.toolbar-actions button:first-child:hover {
  background-color: #0056b3;
}

.toolbar-right {
  display: flex;
  align-items: center;
  gap: 20px;
}

.search-box input {
  padding: 8px 12px;
  border: 1px solid #ddd;
  border-radius: 4px;
  width: 250px;
  font-size: 14px;
}

.quick-stats {
  display: flex;
  gap: 20px;
}

.stat-item {
  display: flex;
  flex-direction: column;
  align-items: flex-end;
}

.stat-label {
  font-size: 12px;
  color: #666;
}

.stat-value {
  font-size: 16px;
  font-weight: bold;
  color: #333;
}

.stat-value.positive {
  color: #28a745;
}

.stat-value.negative {
  color: #dc3545;
}

.content-area {
  flex: 1;
  overflow: auto;
  padding: 20px;
  background-color: #f5f5f5;
}

/* 加载遮罩 */
.loading-overlay {
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background: rgba(255, 255, 255, 0.9);
  display: flex;
  flex-direction: column;
  justify-content: center;
  align-items: center;
  z-index: 9999;
}

.loading-spinner {
  width: 40px;
  height: 40px;
  border: 4px solid #f3f3f3;
  border-top: 4px solid #007bff;
  border-radius: 50%;
  animation: spin 1s linear infinite;
  margin-bottom: 20px;
}

@keyframes spin {
  0% { transform: rotate(0deg); }
  100% { transform: rotate(360deg); }
}

/* 响应式设计 */
@media (max-width: 768px) {
  .sidebar {
    width: 60px;
    padding: 15px;
  }
  
  .sidebar-header h1,
  .sidebar-header p,
  .nav-text {
    display: none;
  }
  
  .toolbar {
    flex-direction: column;
    gap: 15px;
    height: auto;
    padding: 15px;
  }
  
  .toolbar-left,
  .toolbar-right {
    width: 100%;
    justify-content: space-between;
  }
  
  .search-box input {
    width: 100%;
  }
  
  .quick-stats {
    width: 100%;
    justify-content: space-around;
  }
  
  .content-area {
    padding: 15px;
  }
}
</style>