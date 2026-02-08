<template>
  <div class="fund-list">
    <div class="fund-list-header">
      <h2>基金列表</h2>
      <div class="search-box">
        <input 
          type="text" 
          v-model="searchQuery" 
          placeholder="搜索基金代码或名称..."
          @input="searchFunds"
        />
        <button @click="refreshFunds">刷新</button>
      </div>
    </div>
    
    <div class="fund-list-content">
      <div class="fund-grid">
        <div 
          v-for="fund in filteredFunds" 
          :key="fund.code"
          :class="['fund-card', { 'has-position': hasPosition(fund.code) }]"
          @click="selectFund(fund)"
        >
          <div class="fund-header">
            <div class="fund-code">{{ fund.code }}</div>
            <div class="fund-name">{{ fund.name }}</div>
          </div>
          <div class="fund-price">
            <div class="current-price">当前价格: {{ fund.currentPrice }}</div>
            <div :class="['change-rate', fund.changeRate >= 0 ? 'positive' : 'negative']">
              涨跌幅: {{ fund.changeRate >= 0 ? '+' : '' }}{{ fund.changeRate }}%
            </div>
          </div>
          <div class="fund-actions">
            <button @click.stop="addToPortfolio(fund)">添加持仓</button>
            <button @click.stop="viewDetails(fund)">查看详情</button>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script>
export default {
  name: 'FundList',
  props: {
    funds: {
      type: Array,
      default: () => []
    },
    positions: {
      type: Array,
      default: () => []
    }
  },
  data() {
    return {
      searchQuery: ''
    }
  },
  computed: {
    filteredFunds() {
      if (!this.searchQuery.trim()) {
        return this.funds
      }
      
      const query = this.searchQuery.toLowerCase()
      return this.funds.filter(fund => 
        fund.code.toLowerCase().includes(query) ||
        fund.name.toLowerCase().includes(query)
      )
    }
  },
  methods: {
    hasPosition(code) {
      return this.positions.some(pos => pos.fundCode === code)
    },
    
    selectFund(fund) {
      this.$emit('select-fund', fund)
    },
    
    addToPortfolio(fund) {
      this.$emit('add-to-portfolio', fund)
    },
    
    viewDetails(fund) {
      this.$emit('view-details', fund)
    },
    
    searchFunds() {
      // 搜索逻辑已在 computed 中处理
    },
    
    refreshFunds() {
      this.$emit('refresh-funds')
    }
  }
}
</script>

<style scoped>
.fund-list {
  padding: 20px;
}

.fund-list-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 20px;
  border-bottom: 1px solid #eee;
  padding-bottom: 15px;
}

.fund-list-header h2 {
  margin: 0;
  color: #333;
}

.search-box {
  display: flex;
  gap: 10px;
}

.search-box input {
  padding: 8px 12px;
  border: 1px solid #ddd;
  border-radius: 4px;
  width: 300px;
  font-size: 14px;
}

.search-box button {
  padding: 8px 16px;
  background-color: #007bff;
  color: white;
  border: none;
  border-radius: 4px;
  cursor: pointer;
  font-size: 14px;
}

.search-box button:hover {
  background-color: #0056b3;
}

.fund-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(300px, 1fr));
  gap: 15px;
}

.fund-card {
  border: 1px solid #ddd;
  border-radius: 8px;
  padding: 15px;
  cursor: pointer;
  transition: all 0.3s ease;
  background-color: white;
}

.fund-card:hover {
  border-color: #007bff;
  box-shadow: 0 2px 8px rgba(0, 123, 255, 0.1);
}

.fund-card.has-position {
  border-left: 4px solid #28a745;
  background-color: #f8fff9;
}

.fund-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 10px;
}

.fund-code {
  font-weight: bold;
  font-size: 16px;
  color: #333;
}

.fund-name {
  font-size: 14px;
  color: #666;
  text-align: right;
}

.fund-price {
  margin-bottom: 15px;
}

.current-price {
  font-size: 16px;
  font-weight: bold;
  color: #333;
  margin-bottom: 5px;
}

.change-rate {
  font-size: 14px;
  font-weight: bold;
}

.change-rate.positive {
  color: #28a745;
}

.change-rate.negative {
  color: #dc3545;
}

.fund-actions {
  display: flex;
  gap: 10px;
}

.fund-actions button {
  flex: 1;
  padding: 8px;
  border: 1px solid #ddd;
  border-radius: 4px;
  cursor: pointer;
  font-size: 12px;
  background-color: white;
  transition: all 0.3s ease;
}

.fund-actions button:hover {
  background-color: #f8f9fa;
}

.fund-actions button:first-child {
  background-color: #007bff;
  color: white;
  border-color: #007bff;
}

.fund-actions button:first-child:hover {
  background-color: #0056b3;
}

.fund-actions button:last-child {
  background-color: #6c757d;
  color: white;
  border-color: #6c757d;
}

.fund-actions button:last-child:hover {
  background-color: #545b62;
}
</style>