# FundNet - 基金净值估算应用

一个功能完整的跨平台基金投资管理应用，支持实时净值监控、投资组合管理和数据可视化。

## 🚀 快速开始

### 系统要求
- Go 1.20+
- Node.js 16+
- Wails CLI

### 安装依赖

1. 安装 Wails CLI：
```bash
go install github.com/wailsapp/wails/v2/cmd/wails@latest
```

2. 安装前端依赖：
```bash
cd FundNet/frontend
npm install
```

### 构建应用

1. 开发模式构建：
```bash
cd FundNet
wails build -debug
```

2. 生产模式构建：
```bash
cd FundNet
wails build
```

### 启动应用

构建完成后，应用会生成在 `build/bin/` 目录中：

**macOS:**
```bash
open FundNet/build/bin/FundNet.app
```

**Windows:**
```bash
start FundNet/build/bin/FundNet.exe
```

**Linux:**
```bash
./FundNet/build/bin/FundNet
```

## 📱 功能特性

### 📊 基金数据管理
- 实时基金净值获取
- 基金信息本地存储
- 历史数据下载和管理
- 基金列表搜索和筛选

### 💼 投资组合管理
- 持仓添加、编辑、删除
- 多分组管理（科技、医疗、新能源等）
- 实时资产计算和收益统计
- 当日收益和累计收益跟踪

### 📈 实时监控
- 自动刷新机制（可配置间隔）
- 实时净值更新
- 涨跌幅监控
- 投资组合表现追踪

### 📊 数据可视化
- 净值走势曲线图
- 收益统计图表
- 投资组合分布展示
- 实时数据仪表板

### ⚙️ 系统设置
- 刷新间隔配置
- 数据源选择
- 显示偏好设置
- 数据备份和恢复

## 🏗️ 项目结构

```
FundNet/
├── FundNet/                   # Go后端项目
│   ├── main.go                # 应用入口
│   ├── app.go                 # 应用逻辑
│   ├── models/                # 数据模型
│   │   ├── fund.go           # 基金模型
│   │   └── portfolio.go      # 投资组合模型
│   ├── services/             # 业务服务
│   │   ├── database.go       # 数据库服务
│   │   ├── portfolio.go      # 投资组合服务
│   │   ├── realtime.go       # 实时监控服务
│   │   └── fund_data.go      # 基金数据服务
│   └── wails.json            # Wails配置
├── frontend/                 # Vue.js前端项目
│   ├── src/
│   │   ├── App.vue          # 主应用组件
│   │   ├── components/      # 组件目录
│   │   │   ├── FundList.vue     # 基金列表组件
│   │   │   ├── Portfolio.vue    # 投资组合组件
│   │   │   ├── Chart.vue        # 图表组件
│   │   │   └── Settings.vue     # 设置组件
│   │   └── main.js          # 前端入口
│   └── package.json         # 前端依赖
└── README.md                # 项目说明
```

## 🔧 开发指南

### 后端开发
后端使用 Go 语言开发，主要包含以下模块：

- **models**: 数据模型定义
- **services**: 业务逻辑服务
- **database**: 数据库操作
- **realtime**: 实时数据处理

### 前端开发
前端使用 Vue.js 3 开发，包含以下组件：

- **FundList**: 基金列表和搜索
- **Portfolio**: 投资组合管理
- **Chart**: 数据可视化
- **Settings**: 应用设置

### 数据库
应用使用 SQLite 作为本地数据库，数据文件位于：
- macOS: `~/Library/Application Support/FundNet/fundnet.db`
- Windows: `%APPDATA%\FundNet\fundnet.db`
- Linux: `~/.config/FundNet/fundnet.db`

## 🌐 平台支持

- ✅ Windows
- ✅ macOS
- ✅ Linux
- ✅ Web（通过Wails支持）

## 📝 使用说明

1. **启动应用**: 双击生成的应用程序文件
2. **添加基金**: 在基金列表中搜索并添加持仓
3. **管理投资组合**: 在投资组合页面查看和管理持仓
4. **实时监控**: 在监控页面查看净值走势和收益变化
5. **配置设置**: 在设置页面调整应用参数

## 🐛 故障排除

### 常见问题

1. **构建失败**: 确保 Go 和 Node.js 版本符合要求
2. **数据获取失败**: 检查网络连接和API可用性
3. **应用启动失败**: 检查系统权限和依赖

### 日志查看
应用日志位于：
- macOS: `~/Library/Logs/FundNet/`
- Windows: `%APPDATA%\FundNet\logs\`
- Linux: `~/.config/FundNet/logs/`

## 🤝 贡献

欢迎提交 Issue 和 Pull Request！

## 📄 许可证

MIT License

## 🙏 致谢

感谢以下项目的支持：
- [Wails](https://wails.io/) - 跨平台桌面应用框架
- [Vue.js](https://vuejs.org/) - 前端框架
- [SQLite](https://www.sqlite.org/) - 数据库引擎