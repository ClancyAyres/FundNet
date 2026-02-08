<template>
  <div class="chart-container">
    <div class="chart-header">
      <h3>{{ title }}</h3>
      <div class="chart-controls">
        <select v-model="timeRange" @change="updateChart">
          <option value="7">7天</option>
          <option value="30">30天</option>
          <option value="90">90天</option>
          <option value="365">1年</option>
        </select>
        <button @click="toggleChartType">
          {{ chartType === 'line' ? '切换为柱状图' : '切换为折线图' }}
        </button>
      </div>
    </div>
    
    <div class="chart-wrapper">
      <canvas ref="chartCanvas" width="800" height="400"></canvas>
    </div>
    
    <div class="chart-legend">
      <div class="legend-item" v-for="dataset in chartData.datasets" :key="dataset.label">
        <span class="legend-color" :style="{ backgroundColor: dataset.borderColor }"></span>
        <span class="legend-label">{{ dataset.label }}</span>
      </div>
    </div>
  </div>
</template>

<script>
import Chart from 'chart.js/auto'

export default {
  name: 'Chart',
  props: {
    title: {
      type: String,
      default: '净值走势'
    },
    data: {
      type: Array,
      default: () => []
    },
    type: {
      type: String,
      default: 'line'
    }
  },
  data() {
    return {
      chartInstance: null,
      chartType: this.type,
      timeRange: 30
    }
  },
  computed: {
    chartData() {
      if (!this.data || this.data.length === 0) {
        return {
          labels: [],
          datasets: []
        }
      }

      // 按时间排序数据
      const sortedData = [...this.data].sort((a, b) => new Date(a.timestamp) - new Date(b.timestamp))
      
      // 根据时间范围过滤数据
      const cutoffDate = new Date()
      cutoffDate.setDate(cutoffDate.getDate() - this.timeRange)
      const filteredData = sortedData.filter(item => new Date(item.timestamp) >= cutoffDate)

      // 准备图表数据
      const labels = filteredData.map(item => {
        const date = new Date(item.timestamp)
        return date.toLocaleDateString('zh-CN', { month: '2-digit', day: '2-digit' })
      })

      const datasets = [
        {
          label: '净值',
          data: filteredData.map(item => item.currentPrice),
          borderColor: '#007bff',
          backgroundColor: 'rgba(0, 123, 255, 0.1)',
          borderWidth: 2,
          fill: true,
          tension: 0.4
        },
        {
          label: '估值',
          data: filteredData.map(item => item.estimateValue),
          borderColor: '#28a745',
          backgroundColor: 'rgba(40, 167, 69, 0.1)',
          borderWidth: 2,
          fill: false,
          tension: 0.4
        }
      ]

      return {
        labels,
        datasets
      }
    }
  },
  mounted() {
    this.initChart()
  },
  beforeDestroy() {
    if (this.chartInstance) {
      this.chartInstance.destroy()
    }
  },
  watch: {
    chartData: {
      handler() {
        this.updateChart()
      },
      deep: true
    },
    chartType() {
      this.updateChart()
    }
  },
  methods: {
    initChart() {
      const ctx = this.$refs.chartCanvas.getContext('2d')
      
      this.chartInstance = new Chart(ctx, {
        type: this.chartType,
        data: this.chartData,
        options: {
          responsive: true,
          maintainAspectRatio: false,
          plugins: {
            legend: {
              display: true,
              position: 'top'
            },
            tooltip: {
              mode: 'index',
              intersect: false,
              callbacks: {
                label: function(context) {
                  return `${context.dataset.label}: ${context.raw.toFixed(4)}`
                }
              }
            }
          },
          scales: {
            y: {
              beginAtZero: false,
              grid: {
                color: '#eee'
              }
            },
            x: {
              grid: {
                display: false
              }
            }
          },
          interaction: {
            mode: 'nearest',
            axis: 'x',
            intersect: false
          }
        }
      })
    },
    
    updateChart() {
      if (this.chartInstance) {
        this.chartInstance.data = this.chartData
        this.chartInstance.options.type = this.chartType
        this.chartInstance.update()
      } else {
        this.initChart()
      }
    },
    
    toggleChartType() {
      this.chartType = this.chartType === 'line' ? 'bar' : 'line'
    }
  }
}
</script>

<style scoped>
.chart-container {
  background: white;
  border: 1px solid #ddd;
  border-radius: 8px;
  padding: 20px;
  margin-bottom: 20px;
}

.chart-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 15px;
  border-bottom: 1px solid #eee;
  padding-bottom: 10px;
}

.chart-header h3 {
  margin: 0;
  color: #333;
}

.chart-controls {
  display: flex;
  gap: 10px;
  align-items: center;
}

.chart-controls select {
  padding: 6px 12px;
  border: 1px solid #ddd;
  border-radius: 4px;
  font-size: 14px;
}

.chart-controls button {
  padding: 6px 12px;
  background-color: #007bff;
  color: white;
  border: none;
  border-radius: 4px;
  cursor: pointer;
  font-size: 14px;
}

.chart-controls button:hover {
  background-color: #0056b3;
}

.chart-wrapper {
  position: relative;
  height: 400px;
  margin-bottom: 15px;
}

.chart-wrapper canvas {
  width: 100% !important;
  height: 100% !important;
}

.chart-legend {
  display: flex;
  gap: 20px;
  justify-content: center;
  padding-top: 10px;
  border-top: 1px solid #eee;
}

.legend-item {
  display: flex;
  align-items: center;
  gap: 8px;
  font-size: 14px;
  color: #666;
}

.legend-color {
  width: 16px;
  height: 4px;
  border-radius: 2px;
}

.legend-label {
  font-weight: 500;
}
</style>