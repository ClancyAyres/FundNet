<template>
  <div class="settings">
    <div class="settings-header">
      <h2>设置</h2>
    </div>
    
    <div class="settings-content">
      <!-- 刷新设置 -->
      <div class="setting-section">
        <h3>刷新设置</h3>
        <div class="setting-item">
          <label>自动刷新</label>
          <div class="setting-control">
            <label class="switch">
              <input 
                type="checkbox" 
                v-model="settings.autoRefresh"
                @change="updateSettings"
              />
              <span class="slider"></span>
            </label>
          </div>
        </div>
        
        <div class="setting-item">
          <label>刷新间隔 (秒)</label>
          <div class="setting-control">
            <input 
              type="number" 
              v-model.number="settings.refreshInterval"
              min="10"
              max="300"
              @change="updateSettings"
            />
          </div>
        </div>
        
        <div class="setting-item">
          <label>最大历史天数</label>
          <div class="setting-control">
            <input 
              type="number" 
              v-model.number="settings.maxHistoryDays"
              min="7"
              max="365"
              @change="updateSettings"
            />
          </div>
        </div>
      </div>

      <!-- 数据源设置 -->
      <div class="setting-section">
        <h3>数据源设置</h3>
        <div class="setting-item">
          <label>默认数据源</label>
          <div class="setting-control">
            <select v-model="settings.dataSource" @change="updateSettings">
              <option value="tiantian">天天基金</option>
              <option value="eastmoney">东方财富</option>
              <option value="zhaoshang">招商银行</option>
            </select>
          </div>
        </div>
        
        <div class="setting-item">
          <label>API超时时间 (秒)</label>
          <div class="setting-control">
            <input 
              type="number" 
              v-model.number="settings.apiTimeout"
              min="5"
              max="60"
              @change="updateSettings"
            />
          </div>
        </div>
      </div>

      <!-- 显示设置 -->
      <div class="setting-section">
        <h3>显示设置</h3>
        <div class="setting-item">
          <label>显示涨跌颜色</label>
          <div class="setting-control">
            <label class="switch">
              <input 
                type="checkbox" 
                v-model="settings.showColors"
                @change="updateSettings"
              />
              <span class="slider"></span>
            </label>
          </div>
        </div>
        
        <div class="setting-item">
          <label>小数位数</label>
          <div class="setting-control">
            <select v-model="settings.decimalPlaces" @change="updateSettings">
              <option value="2">2位</option>
              <option value="3">3位</option>
              <option value="4">4位</option>
            </select>
          </div>
        </div>
        
        <div class="setting-item">
          <label>货币符号</label>
          <div class="setting-control">
            <select v-model="settings.currencySymbol" @change="updateSettings">
              <option value="CNY">人民币 (¥)</option>
              <option value="USD">美元 ($)</option>
              <option value="EUR">欧元 (€)</option>
            </select>
          </div>
        </div>
      </div>

      <!-- 通知设置 -->
      <div class="setting-section">
        <h3>通知设置</h3>
        <div class="setting-item">
          <label>启用桌面通知</label>
          <div class="setting-control">
            <label class="switch">
              <input 
                type="checkbox" 
                v-model="settings.enableNotifications"
                @change="updateSettings"
              />
              <span class="slider"></span>
            </label>
          </div>
        </div>
        
        <div class="setting-item">
          <label>涨跌幅提醒阈值 (%)</label>
          <div class="setting-control">
            <input 
              type="number" 
              v-model.number="settings.alertThreshold"
              min="0.1"
              max="10"
              step="0.1"
              @change="updateSettings"
            />
          </div>
        </div>
      </div>

      <!-- 备份与恢复 -->
      <div class="setting-section">
        <h3>数据管理</h3>
        <div class="setting-actions">
          <button @click="backupData">备份数据</button>
          <button @click="restoreData">恢复数据</button>
          <button @click="clearData" class="danger">清空数据</button>
        </div>
        
        <div class="setting-item">
          <label>最后备份时间</label>
          <div class="setting-value">
            {{ lastBackupTime || '从未备份' }}
          </div>
        </div>
      </div>

      <!-- 关于 -->
      <div class="setting-section">
        <h3>关于</h3>
        <div class="about-info">
          <div class="info-item">
            <span class="info-label">版本</span>
            <span class="info-value">{{ appVersion }}</span>
          </div>
          <div class="info-item">
            <span class="info-label">数据库版本</span>
            <span class="info-value">{{ dbVersion }}</span>
          </div>
          <div class="info-item">
            <span class="info-label">数据目录</span>
            <span class="info-value">{{ dataPath }}</span>
          </div>
        </div>
        
        <div class="setting-actions">
          <button @click="checkForUpdates">检查更新</button>
          <button @click="openLogs">查看日志</button>
        </div>
      </div>
    </div>
  </div>
</template>

<script>
export default {
  name: 'Settings',
  data() {
    return {
      settings: {
        autoRefresh: true,
        refreshInterval: 30,
        maxHistoryDays: 30,
        dataSource: 'tiantian',
        apiTimeout: 10,
        showColors: true,
        decimalPlaces: 2,
        currencySymbol: 'CNY',
        enableNotifications: false,
        alertThreshold: 1.0
      },
      lastBackupTime: null,
      appVersion: '1.0.0',
      dbVersion: '1.0',
      dataPath: ''
    }
  },
  mounted() {
    this.loadSettings()
    this.loadAppInfo()
  },
  methods: {
    async loadSettings() {
      try {
        // 从后端加载设置
        const autoRefresh = await this.$wails.App.GetConfig('auto_refresh')
        const refreshInterval = await this.$wails.App.GetConfig('refresh_interval')
        const maxHistoryDays = await this.$wails.App.GetConfig('max_history_days')
        
        if (autoRefresh) this.settings.autoRefresh = autoRefresh === 'true'
        if (refreshInterval) this.settings.refreshInterval = parseInt(refreshInterval)
        if (maxHistoryDays) this.settings.maxHistoryDays = parseInt(maxHistoryDays)
      } catch (error) {
        console.error('加载设置失败:', error)
      }
    },
    
    async loadAppInfo() {
      try {
        // 获取应用信息
        this.dataPath = await this.$wails.App.GetConfig('data_path') || '未知'
        this.lastBackupTime = await this.$wails.App.GetConfig('last_backup') || null
      } catch (error) {
        console.error('加载应用信息失败:', error)
      }
    },
    
    async updateSettings() {
      try {
        // 保存设置到后端
        await this.$wails.App.SetConfig('auto_refresh', this.settings.autoRefresh.toString())
        await this.$wails.App.SetConfig('refresh_interval', this.settings.refreshInterval.toString())
        await this.$wails.App.SetConfig('max_history_days', this.settings.maxHistoryDays.toString())
        
        // 通知后端更新刷新间隔
        if (this.settings.autoRefresh) {
          this.$wails.App.SetRefreshInterval(this.settings.refreshInterval)
        }
        
        this.$emit('settings-updated', this.settings)
      } catch (error) {
        console.error('保存设置失败:', error)
        alert('保存设置失败，请重试')
      }
    },
    
    async backupData() {
      try {
        const result = await this.$wails.App.BackupData()
        if (result) {
          this.lastBackupTime = new Date().toLocaleString('zh-CN')
          alert('数据备份成功！')
        } else {
          alert('数据备份失败！')
        }
      } catch (error) {
        console.error('备份数据失败:', error)
        alert('备份数据失败，请重试')
      }
    },
    
    async restoreData() {
      if (!confirm('确定要恢复数据吗？这将覆盖当前所有数据！')) {
        return
      }
      
      try {
        const result = await this.$wails.App.RestoreData()
        if (result) {
          alert('数据恢复成功！')
          this.$emit('data-restored')
        } else {
          alert('数据恢复失败！')
        }
      } catch (error) {
        console.error('恢复数据失败:', error)
        alert('恢复数据失败，请重试')
      }
    },
    
    async clearData() {
      if (!confirm('确定要清空所有数据吗？此操作不可恢复！')) {
        return
      }
      
      try {
        const result = await this.$wails.App.ClearData()
        if (result) {
          alert('数据已清空！')
          this.$emit('data-cleared')
        } else {
          alert('清空数据失败！')
        }
      } catch (error) {
        console.error('清空数据失败:', error)
        alert('清空数据失败，请重试')
      }
    },
    
    async checkForUpdates() {
      try {
        const result = await this.$wails.App.CheckForUpdates()
        if (result) {
          alert('检查更新完成')
        } else {
          alert('检查更新失败')
        }
      } catch (error) {
        console.error('检查更新失败:', error)
        alert('检查更新失败，请重试')
      }
    },
    
    async openLogs() {
      try {
        await this.$wails.App.OpenLogs()
      } catch (error) {
        console.error('打开日志失败:', error)
        alert('打开日志失败，请重试')
      }
    }
  }
}
</script>

<style scoped>
.settings {
  padding: 20px;
}

.settings-header {
  margin-bottom: 20px;
  border-bottom: 1px solid #eee;
  padding-bottom: 15px;
}

.settings-header h2 {
  margin: 0;
  color: #333;
}

.settings-content {
  display: grid;
  gap: 25px;
}

.setting-section {
  background: white;
  border: 1px solid #ddd;
  border-radius: 8px;
  padding: 20px;
}

.setting-section h3 {
  margin: 0 0 15px 0;
  color: #333;
  font-size: 16px;
  border-bottom: 1px solid #eee;
  padding-bottom: 10px;
}

.setting-item {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 10px 0;
  border-bottom: 1px solid #f5f5f5;
}

.setting-item:last-child {
  border-bottom: none;
}

.setting-item label {
  font-weight: 500;
  color: #666;
  width: 200px;
}

.setting-control {
  display: flex;
  align-items: center;
  gap: 10px;
}

.setting-control input,
.setting-control select {
  padding: 8px 12px;
  border: 1px solid #ddd;
  border-radius: 4px;
  font-size: 14px;
  min-width: 150px;
}

.setting-control input[type="number"] {
  width: 100px;
}

.setting-value {
  color: #333;
  font-weight: 500;
}

.setting-actions {
  display: flex;
  gap: 10px;
  margin-top: 15px;
}

.setting-actions button {
  padding: 8px 16px;
  border: 1px solid #ddd;
  border-radius: 4px;
  cursor: pointer;
  font-size: 14px;
  background-color: white;
}

.setting-actions button:hover {
  background-color: #f8f9fa;
}

.setting-actions button:first-child {
  background-color: #007bff;
  color: white;
  border-color: #007bff;
}

.setting-actions button:first-child:hover {
  background-color: #0056b3;
}

.setting-actions button.danger {
  background-color: #dc3545;
  color: white;
  border-color: #dc3545;
}

.setting-actions button.danger:hover {
  background-color: #c82333;
}

/* Switch 滑块样式 */
.switch {
  position: relative;
  display: inline-block;
  width: 44px;
  height: 24px;
}

.switch input {
  opacity: 0;
  width: 0;
  height: 0;
}

.slider {
  position: absolute;
  cursor: pointer;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background-color: #ccc;
  transition: .4s;
  border-radius: 24px;
}

.slider:before {
  position: absolute;
  content: "";
  height: 18px;
  width: 18px;
  left: 3px;
  bottom: 3px;
  background-color: white;
  transition: .4s;
  border-radius: 50%;
}

input:checked + .slider {
  background-color: #007bff;
}

input:checked + .slider:before {
  transform: translateX(20px);
}

.about-info {
  margin-bottom: 15px;
}

.info-item {
  display: flex;
  justify-content: space-between;
  padding: 8px 0;
  border-bottom: 1px solid #f5f5f5;
}

.info-item:last-child {
  border-bottom: none;
}

.info-label {
  color: #666;
  font-weight: 500;
}

.info-value {
  color: #333;
  font-family: monospace;
}
</style>