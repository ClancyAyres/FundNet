# FundNet - 基金净值估计应用

跨平台基金净值估计应用，支持 Windows、macOS 和 Web 平台。

## 技术栈

- **后端**: Go 1.21+ + Gin + Wails
- **前端**: Node.js + React + TypeScript + Vite + Ant Design
- **桌面**: Wails 2.x
- **数据库**: SQLite

## 功能特性

### 核心功能
- 基金数据实时抓取（天天基金/搜狐财经）
- 净值估算与计算
- 实时估值刷新
- 板块分组管理
- 资产统计与图表展示

### 预留功能
- AI 配置接口
- 用户认证系统
- 股票查询与购买
- Docker 容器化部署

## 快速开始

### 环境要求
- Go 1.21+
- Node.js 18+
- npm 或 yarn

### 安装依赖

```bash
# 后端依赖
cd backend
go mod tidy

# 前端依赖
cd frontend
npm install
```

### 运行应用

#### 开发模式

```bash
# 终端1: 启动后端
cd backend
go run main.go

# 终端2: 启动前端
cd frontend
npm run dev
```

#### Wails 桌面应用

```bash
# 安装 Wails
go install github.com/wailsapp/wails/v2/cmd/wails@latest

# 构建桌面应用
wails build
```

### Docker 部署

```bash
docker-compose up -d
```

## 端口配置

| 服务 | 端口 |
|------|------|
| 前端 | 2800 |
| 后端 | 3800 |

## 项目结构

```
FundNet/
├── backend/              # Go 后端服务
│   ├── cmd/
│   ├── internal/
│   │   ├── config/      # 配置管理
│   │   ├── handlers/    # HTTP 处理器
│   │   ├── models/      # 数据模型
│   │   ├── services/    # 业务逻辑
│   │   └── scrapers/    # 数据抓取
│   ├── pkg/
│   │   ├── ai/          # AI 配置接口（预留）
│   │   ├── auth/        # 认证模块（预留）
│   │   └── stock/        # 股票模块（预留）
│   ├── main.go
│   └── go.mod
├── frontend/            # Node.js 前端
│   ├── src/
│   │   ├── components/  # 组件
│   │   ├── pages/       # 页面
│   │   ├── services/    # API 服务
│   │   ├── stores/      # 状态管理
│   │   ├── types/       # 类型定义
│   │   └── utils/       # 工具函数
│   ├── public/
│   ├── package.json
│   └── vite.config.ts
├── wails.json           # Wails 配置
├── docker-compose.yml   # Docker 配置
└── README.md
```

## API 接口

### 基金相关

| 方法 | 路径 | 描述 |
|------|------|------|
| GET | /api/funds | 获取基金列表 |
| GET | /api/funds/:code | 获取基金详情 |
| POST | /api/funds | 添加基金订阅 |
| DELETE | /api/funds/:code | 取消基金订阅 |
| GET | /api/funds/:code/estimate | 获取基金净值估算 |

### 板块相关

| 方法 | 路径 | 描述 |
|------|------|------|
| GET | /api/sectors | 获取板块列表 |
| POST | /api/sectors | 创建板块 |
| PUT | /api/sectors/:id | 更新板块 |
| DELETE | /api/sectors/:id | 删除板块 |

### 资产相关

| 方法 | 路径 | 描述 |
|------|------|------|
| GET | /api/assets | 获取资产统计 |
| POST | /api/assets | 添加持仓 |
| PUT | /api/assets/:id | 更新持仓 |
| DELETE | /api/assets/:id | 删除持仓 |

## 许可证

MIT
