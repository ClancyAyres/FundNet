import { useEffect, useState } from 'react'
import { Card, Form, Input, InputNumber, Button, message, Divider, Table, Select } from 'antd'
import api from '../services/api'

interface Sector {
  id: number
  name: string
  color: string
}

interface Config {
  refresh_interval: number
  log_level: string
}

function Settings() {
  const [sectors, setSectors] = useState<Sector[]>([])
  const [config, setConfig] = useState<Config>({ refresh_interval: 60, log_level: 'info' })
  const [loading, setLoading] = useState(true)
  const [form] = Form.useForm()
  const [configForm] = Form.useForm()

  useEffect(() => {
    fetchData()
  }, [])

  const fetchData = async () => {
    setLoading(true)
    try {
      const [sectorsRes, configRes] = await Promise.all([api.get('/sectors'), api.get('/config')])
      if (sectorsRes.data?.data) {
        setSectors(sectorsRes.data.data as Sector[])
      }
      if (configRes.data?.data) {
        const cfg = configRes.data.data as Config
        setConfig(cfg)
        configForm.setFieldsValue({ refresh_interval: cfg.refresh_interval, log_level: cfg.log_level })
      }
    } catch (error) {
      console.error('Failed to fetch data:', error)
    } finally {
      setLoading(false)
    }
  }

  const handleAddSector = async (values: { name: string; color?: string }) => {
    try {
      await api.post('/sectors', values)
      message.success('板块创建成功')
      form.resetFields()
      fetchData()
    } catch (error) {
      message.error('创建失败')
    }
  }

  const handleDeleteSector = async (id: number) => {
    try {
      await api.delete(`/sectors/${id}`)
      message.success('删除成功')
      fetchData()
    } catch (error) {
      message.error('删除失败')
    }
  }

  const handleUpdateConfig = async (values: { refresh_interval?: number; log_level?: string }) => {
    try {
      await api.put('/config', values)
      message.success('配置更新成功')
    } catch (error) {
      message.error('更新失败')
    }
  }

  const sectorColumns = [
    { title: '板块名称', dataIndex: 'name', key: 'name' },
    {
      title: '颜色',
      dataIndex: 'color',
      key: 'color',
      render: (color: string) => (
        <div
          style={{
            width: 24,
            height: 24,
            backgroundColor: color,
            borderRadius: 4,
            border: '1px solid #d9d9d9',
          }}
        />
      ),
    },
    {
      title: '操作',
      key: 'action',
      render: (_: unknown, record: Sector) => (
        <Button type="link" danger onClick={() => handleDeleteSector(record.id)}>
          删除
        </Button>
      ),
    },
  ]

  return (
    <div>
      <Card title="板块管理" style={{ marginBottom: 24 }}>
        <Form form={form} layout="inline" onFinish={handleAddSector} style={{ marginBottom: 16 }}>
          <Form.Item name="name" rules={[{ required: true, message: '请输入板块名称' }]}>
            <Input placeholder="板块名称" />
          </Form.Item>
          <Form.Item name="color" initialValue="#1890ff">
            <Input type="color" style={{ width: 60 }} />
          </Form.Item>
          <Form.Item>
            <Button type="primary" htmlType="submit">
              添加板块
            </Button>
          </Form.Item>
        </Form>
        <Table columns={sectorColumns} dataSource={sectors} rowKey="id" loading={loading} pagination={false} />
      </Card>

      <Card title="系统配置">
        <Form form={configForm} layout="vertical" onFinish={handleUpdateConfig} initialValues={config}>
          <Form.Item name="refresh_interval" label="刷新间隔（秒）" rules={[{ required: true, message: '请输入刷新间隔' }]}>
            <InputNumber min={10} max={3600} />
          </Form.Item>
          <Form.Item name="log_level" label="日志级别">
            <Select defaultValue="info" style={{ width: 120 }}>
              <Select.Option value="debug">debug</Select.Option>
              <Select.Option value="info">info</Select.Option>
              <Select.Option value="warn">warn</Select.Option>
              <Select.Option value="error">error</Select.Option>
            </Select>
          </Form.Item>
          <Form.Item>
            <Button type="primary" htmlType="submit">
              保存配置
            </Button>
          </Form.Item>
        </Form>
      </Card>

      <Divider />

      <Card title="预留功能">
        <div style={{ color: '#999', lineHeight: 2 }}>
          <p>🤖 AI 配置 - 可配置外部 AI 模型接口</p>
          <p>🔐 用户认证 - 用户登录系统（当前无需登录）</p>
          <p>📈 股票查询 - 股票信息查询与购买接口</p>
          <p>🐳 Docker 部署 - 容器化部署配置</p>
        </div>
      </Card>
    </div>
  )
}

export default Settings
