import { useEffect, useState } from 'react'
import { Card, Table, Button, Modal, Form, Input, Select, message, Tag, Space } from 'antd'
import { PlusOutlined, DeleteOutlined } from '@ant-design/icons'
import api from '../services/api'

interface Fund {
  id: number
  code: string
  name?: string
  sector?: string
  nav: number
  estimate_nav: number
  daily_growth: number
  estimate_time?: string
}

interface Sector {
  id: number
  name: string
}

function Funds() {
  const [funds, setFunds] = useState<Fund[]>([])
  const [sectors, setSectors] = useState<Sector[]>([])
  const [loading, setLoading] = useState(true)
  const [modalVisible, setModalVisible] = useState(false)
  const [form] = Form.useForm()

  useEffect(() => {
    fetchData()
  }, [])

  const fetchData = async () => {
    setLoading(true)
    try {
      const [fundsRes, sectorsRes] = await Promise.all([
        api.get('/funds'),
        api.get('/sectors'),
      ])
      if (fundsRes.data?.data) {
        setFunds(fundsRes.data.data as Fund[])
      }
      if (sectorsRes.data?.data) {
        setSectors(sectorsRes.data.data as Sector[])
      }
    } catch (error) {
      console.error('Failed to fetch data:', error)
    } finally {
      setLoading(false)
    }
  }

  const handleAdd = async (values: { code: string; name?: string; sector?: string }) => {
    try {
      await api.post('/funds', values)
      message.success('添加成功')
      setModalVisible(false)
      form.resetFields()
      fetchData()
    } catch (error) {
      message.error('添加失败')
    }
  }

  const handleRemove = async (code: string) => {
    try {
      await api.delete(`/funds/${code}`)
      message.success('删除成功')
      fetchData()
    } catch (error) {
      message.error('删除失败')
    }
  }

  const columns = [
    { title: '基金代码', dataIndex: 'code', key: 'code' },
    {
      title: '基金名称',
      dataIndex: 'name',
      key: 'name',
      render: (_: unknown, record: Fund) => record.name || record.code,
    },
    {
      title: '板块',
      dataIndex: 'sector',
      key: 'sector',
      render: (sector: string) => (sector ? <Tag color="blue">{sector}</Tag> : '-'),
    },
    {
      title: '单位净值',
      dataIndex: 'nav',
      key: 'nav',
      render: (nav: number) => (nav ? nav.toFixed(4) : '-'),
    },
    {
      title: '估算净值',
      dataIndex: 'estimate_nav',
      key: 'estimate_nav',
      render: (nav: number) => (nav ? nav.toFixed(4) : '-'),
    },
    {
      title: '日增长率',
      dataIndex: 'daily_growth',
      key: 'daily_growth',
      render: (growth: number) => (
        <span style={{ color: (growth || 0) >= 0 ? '#52c41a' : '#ff4d4f' }}>
          {(growth || 0) >= 0 ? '+' : ''}
          {(growth || 0).toFixed(2)}%
        </span>
      ),
    },
    {
      title: '更新时间',
      dataIndex: 'estimate_time',
      key: 'estimate_time',
      render: (time: string) => (time ? time.substring(0, 19) : '-'),
    },
    {
      title: '操作',
      key: 'action',
      render: (_: unknown, record: Fund) => (
        <Space>
          <Button type="link" danger icon={<DeleteOutlined />} onClick={() => handleRemove(record.code)}>
            删除
          </Button>
        </Space>
      ),
    },
  ]

  return (
    <div>
      <Card
        title="基金订阅管理"
        extra={
          <Button type="primary" icon={<PlusOutlined />} onClick={() => setModalVisible(true)}>
            添加基金
          </Button>
        }
      >
        <Table columns={columns} dataSource={funds} rowKey="id" loading={loading} pagination={{ pageSize: 10 }} />
      </Card>

      <Modal title="添加基金订阅" open={modalVisible} onCancel={() => setModalVisible(false)} footer={null}>
        <Form form={form} layout="vertical" onFinish={handleAdd}>
          <Form.Item name="code" label="基金代码" rules={[{ required: true, message: '请输入基金代码' }]}>
            <Input placeholder="例如: 161039" />
          </Form.Item>
          <Form.Item name="name" label="基金名称">
            <Input placeholder="可选" />
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

export default Funds
