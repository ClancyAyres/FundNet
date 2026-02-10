import { useEffect, useState } from 'react'
import { Row, Col, Card, Statistic, Spin } from 'antd'
import api from '../services/api'

interface SummaryData {
  total_current_value: number
  total_cost_basis: number
  total_profit_loss: number
  total_profit_rate: number
  sectors?: unknown[]
}

interface FundData {
  code: string
  name?: string
  daily_growth: number
}

function Dashboard() {
  const [loading, setLoading] = useState(true)
  const [summary, setSummary] = useState<SummaryData>({
    total_current_value: 0,
    total_cost_basis: 0,
    total_profit_loss: 0,
    total_profit_rate: 0,
  })
  const [funds, setFunds] = useState<FundData[]>([])

  useEffect(() => {
    fetchData()
  }, [])

  const fetchData = async () => {
    setLoading(true)
    try {
      const [summaryRes, fundsRes] = await Promise.all([
        api.get('/assets/summary'),
        api.get('/funds'),
      ])
      if (summaryRes.data?.data) {
        setSummary(summaryRes.data.data as SummaryData)
      }
      if (fundsRes.data?.data) {
        setFunds(fundsRes.data.data as FundData[])
      }
    } catch (error) {
      console.error('Failed to fetch data:', error)
    } finally {
      setLoading(false)
    }
  }

  if (loading) {
    return (
      <div style={{ textAlign: 'center', padding: '100px' }}>
        <Spin size="large" />
      </div>
    )
  }

  return (
    <div>
      <Row gutter={[16, 16]} style={{ marginBottom: 24 }}>
        <Col xs={24} sm={12} md={6}>
          <Card>
            <Statistic
              title="总资产"
              value={summary.total_current_value}
              prefix="¥"
              precision={2}
            />
          </Card>
        </Col>
        <Col xs={24} sm={12} md={6}>
          <Card>
            <Statistic
              title="持仓成本"
              value={summary.total_cost_basis}
              prefix="¥"
              precision={2}
            />
          </Card>
        </Col>
        <Col xs={24} sm={12} md={6}>
          <Card>
            <Statistic
              title="累计盈亏"
              value={summary.total_profit_loss}
              prefix="¥"
              precision={2}
              valueStyle={{ color: summary.total_profit_loss >= 0 ? '#52c41a' : '#ff4d4f' }}
            />
          </Card>
        </Col>
        <Col xs={24} sm={12} md={6}>
          <Card>
            <Statistic
              title="收益率"
              value={summary.total_profit_rate}
              suffix="%"
              precision={2}
              valueStyle={{ color: summary.total_profit_rate >= 0 ? '#52c41a' : '#ff4d4f' }}
            />
          </Card>
        </Col>
      </Row>

      <Row gutter={[16, 16]}>
        <Col span={24}>
          <Card title="订阅基金">
            <div style={{ maxHeight: 300, overflow: 'auto' }}>
              {funds.map((fund) => (
                <div key={fund.code} className="fund-card">
                  <div className="fund-header">
                    <div>
                      <div className="fund-name">{fund.name || fund.code}</div>
                      <div className="fund-code">{fund.code}</div>
                    </div>
                    <span style={{ color: fund.daily_growth >= 0 ? '#52c41a' : '#ff4d4f' }}>
                      {fund.daily_growth >= 0 ? '+' : ''}{fund.daily_growth.toFixed(2)}%
                    </span>
                  </div>
                </div>
              ))}
              {funds.length === 0 && (
                <div style={{ textAlign: 'center', padding: 20, color: '#999' }}>
                  暂无订阅基金
                </div>
              )}
            </div>
          </Card>
        </Col>
      </Row>
    </div>
  )
}

export default Dashboard
