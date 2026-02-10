// 基金类型
export interface Fund {
  id: number;
  code: string;
  name: string;
  sector: string;
  nav: number;
  nav_date: string;
  estimate_nav: number;
  estimate_time: string;
  daily_growth: number;
  subscribed: boolean;
  subscribe_time: string;
  created_at: string;
  updated_at: string;
}

// 板块类型
export interface Sector {
  id: number;
  name: string;
  color: string;
  sort_order: number;
  created_at: string;
  updated_at: string;
}

// 持仓类型
export interface Position {
  id: number;
  fund_code: string;
  fund_name: string;
  shares: number;
  cost: number;
  cost_basis: number;
  current_value: number;
  profit_loss: number;
  profit_rate: number;
  daily_growth: number;
  sector: string;
  created_at: string;
  updated_at: string;
}

// 估算结果
export interface EstimateResult {
  code: string;
  name: string;
  nav: number;
  estimate_nav: number;
  daily_growth: number;
  estimate_time: string;
}

// 历史数据点
export interface HistoryPoint {
  time: string;
  estimate_nav: number;
  daily_growth: number;
}

// 资产摘要
export interface AssetSummary {
  sectors: SectorStat[];
  total_cost_basis: number;
  total_current_value: number;
  total_profit_loss: number;
  total_profit_rate: number;
}

// 板块统计
export interface SectorStat {
  name: string;
  cost_basis: number;
  current_value: number;
  profit_loss: number;
  profit_rate: number;
  weight: number;
}

// 配置类型
export interface Config {
  refresh_interval: number;
  log_level: string;
}
