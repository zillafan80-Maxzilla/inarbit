import React, { useState, useEffect } from 'react'
import { BarChart, Bar, LineChart, Line, XAxis, YAxis, CartesianGrid, Tooltip, Legend, ResponsiveContainer, PieChart, Pie, Cell } from 'recharts'
import { TrendingUp, Activity, Zap, DollarSign, ArrowUpRight, ArrowDownLeft, Plus, Settings } from 'lucide-react'
import { useWebSocket } from '../hooks/useWebSocket'
import StatCard from '../components/StatCard'
import BotCard from '../components/BotCard'
import './DashboardPage.css'

/**
 * 仪表板页面组件
 * 展示机器人统计数据、交易历史、实时数据等
 * 设计风格：米色/棕色主题，现代卡片式布局
 */
export default function DashboardPage() {
  const { data: wsData } = useWebSocket()
  
  const [stats, setStats] = useState({
    totalProfit: 0,
    totalTrades: 0,
    activeBots: 0,
    winRate: 0,
  })
  
  const [bots, setBots] = useState([])
  const [chartData, setChartData] = useState([])
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState('')

  // 获取仪表板数据
  useEffect(() => {
    fetchDashboardData()
  }, [])

  // 监听WebSocket数据更新
  useEffect(() => {
    if (wsData) {
      updateStatsFromWebSocket(wsData)
    }
  }, [wsData])

  const fetchDashboardData = async () => {
    try {
      setLoading(true)
      
      // 获取统计数据
      const statsRes = await fetch('/api/dashboard/stats')
      if (!statsRes.ok) throw new Error('获取统计数据失败')
      const statsData = await statsRes.json()
      setStats(statsData)

      // 获取机器人列表
      const botsRes = await fetch('/api/bots')
      if (!botsRes.ok) throw new Error('获取机器人列表失败')
      const botsData = await botsRes.json()
      setBots(botsData.bots || [])

      // 获取图表数据
      const chartRes = await fetch('/api/dashboard/chart-data')
      if (!chartRes.ok) throw new Error('获取图表数据失败')
      const chartDataRes = await chartRes.json()
      setChartData(chartDataRes.data || [])

      setError('')
    } catch (err) {
      setError(err.message || '加载数据失败')
      console.error('获取仪表板数据错误:', err)
    } finally {
      setLoading(false)
    }
  }

  const updateStatsFromWebSocket = (data) => {
    if (data.type === 'stats_update') {
      setStats(prev => ({
        ...prev,
        ...data.payload
      }))
    }
  }

  // 模拟图表数据
  const mockChartData = [
    { date: '2024-01-01', profit: 120, trades: 8 },
    { date: '2024-01-02', profit: 290, trades: 12 },
    { date: '2024-01-03', profit: 200, trades: 10 },
    { date: '2024-01-04', profit: 380, trades: 15 },
    { date: '2024-01-05', profit: 350, trades: 14 },
    { date: '2024-01-06', profit: 420, trades: 18 },
    { date: '2024-01-07', profit: 510, trades: 20 },
  ]

  const mockPieData = [
    { name: '成功交易', value: 72, color: '#9d8b7e' },
    { name: '失败交易', value: 18, color: '#ddd3c3' },
    { name: '待处理', value: 10, color: '#e5ddd0' },
  ]

  const mockBots = [
    {
      id: 1,
      name: '三角套利机器人-1',
      strategy: 'triangular',
      status: 'running',
      profit: 1250.50,
      trades: 45,
      uptime: '12h 34m',
    },
    {
      id: 2,
      name: '四角套利机器人-1',
      strategy: 'quadrangular',
      status: 'running',
      profit: 2100.75,
      trades: 68,
      uptime: '8h 20m',
    },
    {
      id: 3,
      name: '五角套利机器人-1',
      strategy: 'pentagonal',
      status: 'stopped',
      profit: 850.25,
      trades: 32,
      uptime: '0h 0m',
    },
  ]

  return (
    <div className="dashboard-container">
      {/* 页面头部 */}
      <div className="dashboard-header">
        <div className="dashboard-title-section">
          <h1 className="dashboard-title">仪表板</h1>
          <p className="dashboard-subtitle">实时监控机器人运行状态和交易数据</p>
        </div>
        
        <div className="dashboard-header-actions">
          <button className="dashboard-btn dashboard-btn-secondary">
            <Settings size={18} />
            <span>设置</span>
          </button>
          <button className="dashboard-btn dashboard-btn-primary">
            <Plus size={18} />
            <span>新建机器人</span>
          </button>
        </div>
      </div>

      {/* 错误提示 */}
      {error && (
        <div className="dashboard-error-alert">
          <span>⚠️ {error}</span>
          <button onClick={fetchDashboardData} className="dashboard-retry-btn">
            重试
          </button>
        </div>
      )}

      {/* 统计卡片 */}
      <div className="dashboard-stats-grid">
        <StatCard
          icon={<DollarSign size={24} />}
          label="总收益"
          value={`$${stats.totalProfit?.toFixed(2) || '0.00'}`}
          trend={stats.totalProfit > 0 ? 'up' : 'down'}
          trendValue={`${Math.abs(stats.totalProfit || 0).toFixed(2)}`}
          color="primary"
        />
        
        <StatCard
          icon={<Activity size={24} />}
          label="总交易数"
          value={stats.totalTrades || '0'}
          trend="up"
          trendValue={`${Math.floor(Math.random() * 10) + 1} 次`}
          color="secondary"
        />
        
        <StatCard
          icon={<Zap size={24} />}
          label="活跃机器人"
          value={stats.activeBots || '0'}
          trend="neutral"
          trendValue="3 个"
          color="tertiary"
        />
        
        <StatCard
          icon={<TrendingUp size={24} />}
          label="胜率"
          value={`${stats.winRate || '0'}%`}
          trend="up"
          trendValue="↑ 5.2%"
          color="success"
        />
      </div>

      {/* 主要内容区域 */}
      <div className="dashboard-main-grid">
        {/* 左侧：图表区域 */}
        <div className="dashboard-charts-section">
          {/* 收益趋势图 */}
          <div className="dashboard-card">
            <div className="dashboard-card-header">
              <h2 className="dashboard-card-title">收益趋势</h2>
              <select className="dashboard-select">
                <option>最近7天</option>
                <option>最近30天</option>
                <option>最近90天</option>
              </select>
            </div>
            <div className="dashboard-card-body">
              <ResponsiveContainer width="100%" height={300}>
                <LineChart data={mockChartData}>
                  <CartesianGrid strokeDasharray="3 3" stroke="#e5ddd0" />
                  <XAxis dataKey="date" stroke="#9d8b7e" />
                  <YAxis stroke="#9d8b7e" />
                  <Tooltip 
                    contentStyle={{
                      backgroundColor: '#fff',
                      border: '1px solid #e5ddd0',
                      borderRadius: '8px',
                    }}
                  />
                  <Legend />
                  <Line 
                    type="monotone" 
                    dataKey="profit" 
                    stroke="#9d8b7e" 
                    strokeWidth={2}
                    dot={{ fill: '#9d8b7e', r: 4 }}
                    activeDot={{ r: 6 }}
                    name="收益 ($)"
                  />
                </LineChart>
              </ResponsiveContainer>
            </div>
          </div>

          {/* 交易统计图 */}
          <div className="dashboard-card">
            <div className="dashboard-card-header">
              <h2 className="dashboard-card-title">交易统计</h2>
            </div>
            <div className="dashboard-card-body">
              <ResponsiveContainer width="100%" height={300}>
                <BarChart data={mockChartData}>
                  <CartesianGrid strokeDasharray="3 3" stroke="#e5ddd0" />
                  <XAxis dataKey="date" stroke="#9d8b7e" />
                  <YAxis stroke="#9d8b7e" />
                  <Tooltip 
                    contentStyle={{
                      backgroundColor: '#fff',
                      border: '1px solid #e5ddd0',
                      borderRadius: '8px',
                    }}
                  />
                  <Legend />
                  <Bar dataKey="trades" fill="#9d8b7e" name="交易数" radius={[8, 8, 0, 0]} />
                </BarChart>
              </ResponsiveContainer>
            </div>
          </div>
        </div>

        {/* 右侧：机器人和统计 */}
        <div className="dashboard-sidebar">
          {/* 交易成功率 */}
          <div className="dashboard-card">
            <div className="dashboard-card-header">
              <h2 className="dashboard-card-title">交易成功率</h2>
            </div>
            <div className="dashboard-card-body dashboard-pie-container">
              <ResponsiveContainer width="100%" height={250}>
                <PieChart>
                  <Pie
                    data={mockPieData}
                    cx="50%"
                    cy="50%"
                    innerRadius={60}
                    outerRadius={90}
                    paddingAngle={2}
                    dataKey="value"
                  >
                    {mockPieData.map((entry, index) => (
                      <Cell key={`cell-${index}`} fill={entry.color} />
                    ))}
                  </Pie>
                </PieChart>
              </ResponsiveContainer>
              <div className="dashboard-pie-legend">
                {mockPieData.map((item, index) => (
                  <div key={index} className="dashboard-pie-item">
                    <span className="dashboard-pie-dot" style={{ backgroundColor: item.color }}></span>
                    <span className="dashboard-pie-label">{item.name}</span>
                    <span className="dashboard-pie-value">{item.value}%</span>
                  </div>
                ))}
              </div>
            </div>
          </div>

          {/* 最近活动 */}
          <div className="dashboard-card">
            <div className="dashboard-card-header">
              <h2 className="dashboard-card-title">最近活动</h2>
            </div>
            <div className="dashboard-activity-list">
              <div className="dashboard-activity-item">
                <div className="dashboard-activity-icon dashboard-activity-success">
                  <ArrowUpRight size={16} />
                </div>
                <div className="dashboard-activity-content">
                  <p className="dashboard-activity-title">成功交易</p>
                  <p className="dashboard-activity-time">2分钟前</p>
                </div>
                <span className="dashboard-activity-amount">+$125.50</span>
              </div>

              <div className="dashboard-activity-item">
                <div className="dashboard-activity-icon dashboard-activity-warning">
                  <ArrowDownLeft size={16} />
                </div>
                <div className="dashboard-activity-content">
                  <p className="dashboard-activity-title">交易失败</p>
                  <p className="dashboard-activity-time">8分钟前</p>
                </div>
                <span className="dashboard-activity-amount">-$45.20</span>
              </div>

              <div className="dashboard-activity-item">
                <div className="dashboard-activity-icon dashboard-activity-info">
                  <Zap size={16} />
                </div>
                <div className="dashboard-activity-content">
                  <p className="dashboard-activity-title">机器人启动</p>
                  <p className="dashboard-activity-time">15分钟前</p>
                </div>
                <span className="dashboard-activity-amount">-</span>
              </div>
            </div>
          </div>
        </div>
      </div>

      {/* 机器人列表 */}
      <div className="dashboard-bots-section">
        <div className="dashboard-section-header">
          <h2 className="dashboard-section-title">活跃机器人</h2>
          <a href="/bots" className="dashboard-view-all">
            查看全部 →
          </a>
        </div>
        
        <div className="dashboard-bots-grid">
          {mockBots.map(bot => (
            <BotCard key={bot.id} bot={bot} />
          ))}
        </div>
      </div>
    </div>
  )
}
