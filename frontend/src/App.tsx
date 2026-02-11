import { Routes, Route, Link, useLocation } from 'react-router-dom'
import { Layout, Menu } from 'antd'
import { HomeOutlined, FundOutlined, PieChartOutlined, SettingOutlined } from '@ant-design/icons'
import Dashboard from './pages/Dashboard'
import Funds from './pages/Funds'
import Positions from './pages/Positions'
import Settings from './pages/Settings'

const { Header, Content, Footer } = Layout

function App() {
  const location = useLocation()
  const currentKey = location.pathname

  const menuItems = [
    {
      key: '/',
      icon: <HomeOutlined />,
      label: <Link to="/">首页</Link>,
    },
    {
      key: '/funds',
      icon: <FundOutlined />,
      label: <Link to="/funds">基金</Link>,
    },
    {
      key: '/positions',
      icon: <PieChartOutlined />,
      label: <Link to="/positions">持仓</Link>,
    },
    {
      key: '/settings',
      icon: <SettingOutlined />,
      label: <Link to="/settings">设置</Link>,
    },
  ]

  return (
    <Layout className="site-layout">
      <Header className="site-header">
        <div className="site-logo">FundNet</div>
        <Menu
          mode="horizontal"
          selectedKeys={[currentKey]}
          items={menuItems}
          style={{ minWidth: 400, borderBottom: 'none' }}
        />
      </Header>
      <Content className="site-content">
        <Routes>
          <Route path="/" element={<Dashboard />} />
          <Route path="/funds" element={<Funds key="funds" />} />
          <Route path="/positions" element={<Positions key="positions" />} />
          <Route path="/settings" element={<Settings key="settings" />} />
        </Routes>
      </Content>
      <Footer style={{ textAlign: 'center' }}>
        FundNet ©{new Date().getFullYear()} 基金净值估计应用
      </Footer>
    </Layout>
  )
}

export default App
