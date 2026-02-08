<template>
  <div class="portfolio">
    <div class="portfolio-header">
      <h2>投资组合</h2>
      <div class="portfolio-actions">
        <button @click="showAddPosition = true">添加持仓</button>
        <button @click="refreshPortfolio">刷新</button>
        <button @click="exportData">导出数据</button>
      </div>
    </div>

    <!-- 总览卡片 -->
    <div class="overview-cards">
      <div class="overview-card total-value">
        <div class="card-title">总资产</div>
        <div class="card-value">{{ formatCurrency(portfolioSummary.totalValue) }}</div>
        <div :class="['card-change', portfolioSummary.totalGain >= 0 ? 'positive' : 'negative']">
          {{ portfolioSummary.totalGain >= 0 ? '+' : '' }}{{ formatCurrency(portfolioSummary.totalGain) }}
          ({{ portfolioSummary.totalGainRate >= 0 ? '+' : '' }}{{ portfolioSummary.totalGainRate.toFixed(2) }}%)
        </div>
      </div>
      
      <div class="overview-card daily-gain">
        <div class="card-title">当日收益</div>
        <div :class="['card-value', portfolioSummary.dailyGain >= 0 ? 'positive' : 'negative']">
          {{ portfolioSummary.dailyGain >= 0 ? '+' : '' }}{{ formatCurrency(portfolioSummary.dailyGain) }}
        </div>
        <div :class="['card-change', portfolioSummary.dailyGain >= 0 ? 'positive' : 'negative']">
          ({{ portfolioSummary.dailyGainRate >= 0 ? '+' : '' }}{{ portfolioSummary.dailyGainRate.toFixed(2) }}%)
        </div>
      </div>
      
      <div class="overview-card total-cost">
        <div class="card-title">总成本</div>
        <div class="card-value">{{ formatCurrency(portfolioSummary.totalCost) }}</div>
        <div class="card-change">持仓总成本</div>
      </div>
    </div>

    <!-- 分组管理 -->
    <div class="groups-section">
      <div class="section-header">
        <h3>分组管理</h3>
        <div class="group-actions">
          <input 
            v-model="newGroupName" 
            placeholder="新分组名称"
            @keyup.enter="addGroup"
          />
          <button @click="addGroup">添加分组</button>
        </div>
      </div>
      
      <div class="groups-list">
        <div 
          v-for="group in groups" 
          :key="group"
          :class="['group-card', { active: selectedGroup === group }]"
          @click="selectGroup(group)"
        >
          <div class="group-name">{{ group }}</div>
          <div class="group-stats">
            <span class="group-value">{{ formatCurrency(groupValues[group]?.value || 0) }}</span>
            <span :class="['group-gain', (groupValues[group]?.gain || 0) >= 0 ? 'positive' : 'negative']">
              {{ (groupValues[group]?.gain || 0) >= 0 ? '+' : '' }}{{ formatCurrency(groupValues[group]?.gain || 0) }}
            </span>
          </div>
          <div class="group-actions">
            <button @click.stop="renameGroup(group)">重命名</button>
            <button @click.stop="deleteGroup(group)">删除</button>
          </div>
        </div>
      </div>
    </div>

    <!-- 持仓列表 -->
    <div class="positions-section">
      <div class="section-header">
        <h3>{{ selectedGroup === '全部' ? '全部持仓' : selectedGroup + '持仓' }}</h3>
        <div class="position-actions">
          <input 
            v-model="positionSearch" 
            placeholder="搜索持仓..."
            @input="filterPositions"
          />
        </div>
      </div>
      
      <div class="positions-list">
        <div 
          v-for="position in filteredPositions" 
          :key="position.fundCode"
          class="position-card"
        >
          <div class="position-header">
            <div class="fund-info">
              <div class="fund-code">{{ position.fundCode }}</div>
              <div class="fund-name">{{ position.fundName || '未知基金' }}</div>
            </div>
            <div class="position-actions">
              <button @click="editPosition(position)">编辑</button>
              <button @click="removePosition(position.fundCode)">删除</button>
            </div>
          </div>
          
          <div class="position-details">
            <div class="detail-row">
              <span>份额: {{ position.shares }}</span>
              <span>成本价: {{ formatCurrency(position.costPrice) }}</span>
            </div>
            <div class="detail-row">
              <span>当前价格: {{ formatCurrency(position.currentPrice) }}</span>
              <span>涨跌幅: {{ position.changeRate >= 0 ? '+' : '' }}{{ position.changeRate }}%</span>
            </div>
            <div class="detail-row">
              <span>持仓价值: {{ formatCurrency(position.estimateValue) }}</span>
              <span>浮动盈亏: {{ position.estimateChange >= 0 ? '+' : '' }}{{ formatCurrency(position.estimateChange) }}</span>
            </div>
            <div class="detail-row">
              <span>盈亏率: {{ position.gainRate >= 0 ? '+' : '' }}{{ position.gainRate.toFixed(2) }}%</span>
              <span>当日盈亏: {{ position.dailyGain >= 0 ? '+' : '' }}{{ formatCurrency(position.dailyGain) }}</span>
            </div>
          </div>
        </div>
      </div>
    </div>

    <!-- 添加持仓弹窗 -->
    <div v-if="showAddPosition" class="modal-overlay" @click="showAddPosition = false">
      <div class="modal" @click.stop>
        <div class="modal-header">
          <h3>添加持仓</h3>
          <button @click="showAddPosition = false">×</button>
        </div>
        <div class="modal-body">
          <div class="form-group">
            <label>基金代码</label>
            <input v-model="newPosition.code" placeholder="请输入基金代码" />
          </div>
          <div class="form-group">
            <label>基金名称</label>
            <input v-model="newPosition.name" placeholder="请输入基金名称" />
          </div>
          <div class="form-group">
            <label>份额</label>
            <input v-model.number="newPosition.shares" type="number" placeholder="请输入持有份额" />
          </div>
          <div class="form-group">
            <label>成本价</label>
            <input v-model.number="newPosition.costPrice" type="number" placeholder="请输入成本价格" />
          </div>
          <div class="form-group">
            <label>分组</label>
            <select v-model="newPosition.groupName">
              <option v-for="group in groups" :key="group" :value="group">{{ group }}</option>
            </select>
          </div>
        </div>
        <div class="modal-footer">
          <button @click="addPosition">添加</button>
          <button @click="showAddPosition = false">取消</button>
        </div>
      </div>
    </div>
  </div>
</template>

<script>
export default {
  name: 'Portfolio',
  props: {
    portfolioSummary: {
      type: Object,
      default: () => ({
        totalValue: 0,
        totalCost: 0,
        totalGain: 0,
        totalGainRate: 0,
        dailyGain: 0,
        dailyGainRate: 0
      })
    },
    positions: {
      type: Array,
      default: () => []
    },
    groups: {
      type: Array,
      default: () => ['全部']
    },
    groupValues: {
      type: Object,
      default: () => ({})
    }
  },
  data() {
    return {
      selectedGroup: '全部',
      showAddPosition: false,
      newGroupName: '',
      positionSearch: '',
      newPosition: {
        code: '',
        name: '',
        shares: 0,
        costPrice: 0,
        groupName: '全部'
      }
    }
  },
  computed: {
    filteredPositions() {
      if (!this.positionSearch.trim()) {
        return this.getPositionsByGroup(this.selectedGroup)
      }
      
      const query = this.positionSearch.toLowerCase()
      return this.getPositionsByGroup(this.selectedGroup).filter(pos => 
        pos.fundCode.toLowerCase().includes(query) ||
        pos.fundName.toLowerCase().includes(query)
      )
    }
  },
  methods: {
    formatCurrency(value) {
      return '¥' + (value || 0).toLocaleString('zh-CN', {
        minimumFractionDigits: 2,
        maximumFractionDigits: 2
      })
    },
    
    getPositionsByGroup(groupName) {
      if (groupName === '全部') {
        return this.positions
      }
      return this.positions.filter(pos => pos.groupName === groupName)
    },
    
    selectGroup(group) {
      this.selectedGroup = group
    },
    
    addGroup() {
      if (this.newGroupName.trim()) {
        this.$emit('add-group', this.newGroupName.trim())
        this.newGroupName = ''
      }
    },
    
    renameGroup(group) {
      const newName = prompt('请输入新分组名称:', group)
      if (newName && newName.trim()) {
        this.$emit('rename-group', group, newName.trim())
      }
    },
    
    deleteGroup(group) {
      if (confirm(`确定要删除分组 "${group}" 吗？`)) {
        this.$emit('delete-group', group)
        if (this.selectedGroup === group) {
          this.selectedGroup = '全部'
        }
      }
    },
    
    editPosition(position) {
      this.newPosition = { ...position }
      this.showAddPosition = true
    },
    
    addPosition() {
      if (!this.newPosition.code || !this.newPosition.name || this.newPosition.shares <= 0 || this.newPosition.costPrice <= 0) {
        alert('请填写完整信息')
        return
      }
      
      this.$emit('add-position', { ...this.newPosition })
      this.showAddPosition = false
      this.resetNewPosition()
    },
    
    removePosition(fundCode) {
      if (confirm('确定要删除这个持仓吗？')) {
        this.$emit('remove-position', fundCode)
      }
    },
    
    refreshPortfolio() {
      this.$emit('refresh-portfolio')
    },
    
    exportData() {
      this.$emit('export-data')
    },
    
    filterPositions() {
      // 搜索逻辑在 computed 中处理
    },
    
    resetNewPosition() {
      this.newPosition = {
        code: '',
        name: '',
        shares: 0,
        costPrice: 0,
        groupName: '全部'
      }
    }
  }
}
</script>

<style scoped>
.portfolio {
  padding: 20px;
}

.portfolio-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 20px;
  border-bottom: 1px solid #eee;
  padding-bottom: 15px;
}

.portfolio-header h2 {
  margin: 0;
  color: #333;
}

.portfolio-actions {
  display: flex;
  gap: 10px;
}

.portfolio-actions button {
  padding: 8px 16px;
  background-color: #007bff;
  color: white;
  border: none;
  border-radius: 4px;
  cursor: pointer;
  font-size: 14px;
}

.portfolio-actions button:hover {
  background-color: #0056b3;
}

.overview-cards {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(250px, 1fr));
  gap: 15px;
  margin-bottom: 30px;
}

.overview-card {
  background: white;
  border: 1px solid #ddd;
  border-radius: 8px;
  padding: 20px;
  box-shadow: 0 2px 4px rgba(0,0,0,0.1);
}

.card-title {
  font-size: 14px;
  color: #666;
  margin-bottom: 10px;
}

.card-value {
  font-size: 24px;
  font-weight: bold;
  color: #333;
  margin-bottom: 5px;
}

.card-change {
  font-size: 14px;
  font-weight: bold;
}

.card-change.positive {
  color: #28a745;
}

.card-change.negative {
  color: #dc3545;
}

.total-value .card-title { color: #007bff; }
.daily-gain .card-title { color: #28a745; }
.total-cost .card-title { color: #6c757d; }

.groups-section {
  margin-bottom: 30px;
}

.section-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 15px;
}

.section-header h3 {
  margin: 0;
  color: #333;
}

.group-actions {
  display: flex;
  gap: 10px;
}

.group-actions input {
  padding: 8px 12px;
  border: 1px solid #ddd;
  border-radius: 4px;
  width: 200px;
}

.groups-list {
  display: flex;
  gap: 15px;
  flex-wrap: wrap;
}

.group-card {
  background: white;
  border: 1px solid #ddd;
  border-radius: 8px;
  padding: 15px;
  cursor: pointer;
  transition: all 0.3s ease;
  width: 200px;
}

.group-card:hover {
  border-color: #007bff;
  box-shadow: 0 2px 8px rgba(0, 123, 255, 0.1);
}

.group-card.active {
  border-color: #007bff;
  background-color: #f8f9ff;
}

.group-name {
  font-weight: bold;
  margin-bottom: 10px;
  color: #333;
}

.group-stats {
  display: flex;
  justify-content: space-between;
  margin-bottom: 10px;
}

.group-value {
  font-weight: bold;
  color: #333;
}

.group-gain {
  font-weight: bold;
}

.group-gain.positive {
  color: #28a745;
}

.group-gain.negative {
  color: #dc3545;
}

.group-actions {
  display: flex;
  gap: 5px;
}

.group-actions button {
  padding: 4px 8px;
  font-size: 12px;
  border: 1px solid #ddd;
  border-radius: 4px;
  cursor: pointer;
  background-color: white;
}

.group-actions button:hover {
  background-color: #f8f9fa;
}

.positions-section {
  background: white;
  border: 1px solid #ddd;
  border-radius: 8px;
  padding: 20px;
}

.positions-list {
  display: grid;
  gap: 15px;
}

.position-card {
  border: 1px solid #eee;
  border-radius: 8px;
  padding: 15px;
  background-color: #fafafa;
}

.position-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 10px;
}

.fund-info {
  flex: 1;
}

.fund-code {
  font-weight: bold;
  font-size: 16px;
  color: #333;
  margin-bottom: 5px;
}

.fund-name {
  font-size: 14px;
  color: #666;
}

.position-actions button {
  padding: 6px 12px;
  font-size: 12px;
  border: 1px solid #ddd;
  border-radius: 4px;
  cursor: pointer;
  background-color: white;
  margin-left: 5px;
}

.position-actions button:hover {
  background-color: #f8f9fa;
}

.position-details {
  display: grid;
  grid-template-columns: repeat(2, 1fr);
  gap: 10px;
}

.detail-row {
  display: flex;
  justify-content: space-between;
  font-size: 14px;
  color: #666;
}

.modal-overlay {
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background-color: rgba(0, 0, 0, 0.5);
  display: flex;
  justify-content: center;
  align-items: center;
  z-index: 1000;
}

.modal {
  background: white;
  border-radius: 8px;
  padding: 20px;
  width: 400px;
  max-width: 90%;
}

.modal-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 20px;
}

.modal-header h3 {
  margin: 0;
}

.modal-header button {
  background: none;
  border: none;
  font-size: 24px;
  cursor: pointer;
  color: #666;
}

.form-group {
  margin-bottom: 15px;
}

.form-group label {
  display: block;
  margin-bottom: 5px;
  font-weight: bold;
  color: #333;
}

.form-group input, .form-group select {
  width: 100%;
  padding: 8px 12px;
  border: 1px solid #ddd;
  border-radius: 4px;
  font-size: 14px;
}

.modal-footer {
  display: flex;
  justify-content: flex-end;
  gap: 10px;
  margin-top: 20px;
}

.modal-footer button {
  padding: 8px 16px;
  border: none;
  border-radius: 4px;
  cursor: pointer;
  font-size: 14px;
}

.modal-footer button:first-child {
  background-color: #007bff;
  color: white;
}

.modal-footer button:last-child {
  background-color: #6c757d;
  color: white;
}
</style>