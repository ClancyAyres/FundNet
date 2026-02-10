import { useEffect, useState } from 'react'
import { Card, Table, Button, Modal, Form, Input, InputNumber, Select, message, Space } from 'antd'
import { PlusOutlined, DeleteOutlined } from '@ant-design/icons'
import api from '../services/api'

interface Position {
  id: number
  fund_code: string
  fund_name?: string
  sector?: string
  shares: number
  cost: number
  cost_basis: number
  current_value: number
  profit_loss: number
  profit_rate: number
}

interface Sector {
  id: number
  name: string
}

interface Fund {
  code: string
  name?: string
}

function Positions() {
  const [positions, setPositions] = useState<Position[]>([])
  const [sectors, setSectors] = useState<Sector[]>([])
  const [funds, setFunds] = useState<Fund[]>([])
  const [loading, setLoading] = useState(true)
  const [modalVisible, setModalVisible] = useState(false)
  const [form] = Form.useForm()

  useEffect(() => {
    fetchData()
  }, [])

  const fetchData = async () => {
    setLoading(true)
    try {
      const [positionsRes, sectorsRes, fundsRes] = await Promise.all([
        api.get('/positions'),
        api.get('/sectors'),
        api.get('/funds'),
      ])
      if (positionsRes.data?.data) {
        setPositions(positionsRes.data.data as Position[])
      }
      if (sectorsRes.data?.data) {
        setSectors(sectorsRes.data.data as Sector[])
      }
      if (fundsRes.data?.data) {
        setFunds(fundsRes.data.data as Fund[])
      }
    } catch (error) {
      console.error('Failed to fetch data:', error)
    } finally {
      setLoading(false)
    }
  }

  const handleAdd = async (values: { fundCode: string; fundName?: string; shares: number; cost: number; sector?: string }) => {
    try {
      await api.post('/positions', values)
      message.success('添加成功')
      setModalVisible(false)
      form.resetFields()
      fetchData()
    } catch (error) {
      message.error('添加失败')
    }
  }

  const handleRemove = async (id: number) => {
    try {
      await api.delete(`/positions/${id}`)
      message.success('删除成功')
      fetchData()
    } catch (error) {
      message.error('删除失败')
    }
  }

  const columns = [
    { title: '基金代码', dataIndex: 'fund_code', key: 'fund_code' },
    {
      title: '基金名称',
      dataIndex: 'fund_name',
      key: 'fund_name',
      render: (_: unknown, record: Position) => record.fund_name || record.fund_code,
    },
    {
      title: '板块',
      dataIndex: 'sector',
      key: 'sector',
      render: (sector: string) => (sector ? <span style={{ backgroundColor: '#1890ff', color: '#fff', padding: '2px 8px', borderRadius: 4 }}>{sector}</span> : '-'),
    },
    { title: '持有份额', dataIndex: 'shares', key: 'shares', render: (shares: number) => (shares ? shares.toFixed(2) : '-') },
    { title: '持仓成本', dataIndex: 'cost', key: 'cost', render: (cost: number) => `¥${(cost || 0).toFixed(4)}` },
    { title: '成本总额', dataIndex: 'cost_basis', key: 'cost_basis', render: (basis: number) => `¥${(basis || 0).toFixed(2)}` },
    { title: '当前价值', dataIndex: 'current_value', key: 'current_value', render: (value: number) => `¥${(value || 0).toFixed(2)}` },
    {
      title: '盈亏',
      dataIndex: 'profit_loss',
      key: 'profit_loss',
      render: (loss: number) => (
        <span style={{ color: (loss || 0) >= 0 ? '#52c41a' : '#ff4d4f' }}>
          {(loss || 0) >= 0 ? '+' : ''}¥{(loss || 0).toFixed(2)}
        </span>
      ),
    },
    {
      title: '收益率',
      dataIndex: 'profit_rate',
      key: 'profit_rate',
      render: (rate: number) => (
        <span style={{ color: (rate || 0) >= 0 ? '#52c41a' : '#ff4d4f' }}>
          {(rate || 0) >= 0 ? '+' : ''}{(rate || 0).toFixed(2)}%
        </span>
      ),
    },
    {
      title: '操作',
      key: 'action',
      render: (_: unknown, record: Position) => (
        <Space>
          <Button type="link" danger icon={<DeleteOutlined />} onClick={() => handleRemove(record.id)}>
            删除
          </Button>
        </Space>
      ),
    },
  ]

  return (
    <div>
      <Card
        title="持仓管理"
        extra={
          <Button type="primary" icon={<PlusOutlined />} onClick={() => setModalVisible(true)}>
            添加持仓
          </Button>
        }
      >
        <Table columns={columns} dataSource={positions} rowKey="id" loading={loading} pagination={{ pageSize: 10 }} scroll={{ x: 1200 }} />
      </Card>

      <Modal title="添加持仓" open={modalVisible} onCancel={() => setModalVisible(false)} footer={null}>
        <Form form={form} layout="vertical" onFinish={handleAdd}>
          <Form.Item name="fundCode" label="基金代码" rules={[{ required: true, message: '请选择或输入基金代码' }]}>
            <Select showSearch placeholder="选择或输入基金代码" optionFilterProp="children">
              {funds.map((f) => (
                <Select.Option key={f.code} value={f.code}>
                  {f.code} - {f.name || f.code}
                </Select.Option>
              ))}
            </Select>
          </Form.Item>
          <Form.Item name="fundName" label="基金名称">
            <Input placeholder="可选" />
          </Form.Item>
          <Form.Item name="shares" label="持有份额" rules={[{ required: true, message: '请输入持有份额' }]}>
            <InputNumber min={0} precision={2} style={{ width: '100%' }} placeholder="例如: 1000" />
          </Form.Item>
          <Form.Item name="cost" label="持仓成本" rules={[{ required: true, message: '请输入持仓成本' }]}>
            <InputNumber min={0} precision={4} style={{ width: '100%' }} placeholder="例如: 1.2345" />
          </Form.Item>
          <Form.Item name="sector" label="板块">
            <Select placeholder="选择板块">
              {sectors.map((s) => (
                <Select.Option key={s.id} value={s.name}>
                  {s.name}
                </Select.Option>
              ))}
            </Select>
          </Form.Item>
          <Form.Item>
            <Button type="primary" htmlType="submit" block>
              添加
            </Button>
          </Form.Item>
        </Form>
      </Modal>
    </div>
  )
}

export default Positions
