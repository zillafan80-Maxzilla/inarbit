import React, { useEffect, useState } from 'react'
import { BrowserRouter as Router, Routes, Route, Navigate } from 'react-router-dom'
import { useAuthStore } from './store/authStore'
import { useWebSocket } from './hooks/useWebSocket'

// 页面组件
import LoginPage from './pages/LoginPage'
import DashboardPage from './pages/DashboardPage'
import BotsPage from './pages/BotsPage'
import StrategiesPage from './pages/StrategiesPage'
import TradesPage from './pages/TradesPage'
import SettingsPage from './pages/SettingsPage'
import NotFoundPage from './pages/NotFoundPage'

// 布局组件
import Layout from './components/Layout'
import ProtectedRoute from './components/ProtectedRoute'

// 样式
import './styles/index.css'

function App() {
  const { isAuthenticated, checkAuth } = useAuthStore()
  const [loading, setLoading] = useState(true)
  const { connect } = useWebSocket()

  useEffect(() => {
    // 检查认证状态
    const initAuth = async () => {
      await checkAuth()
      setLoading(false)
    }

    initAuth()
  }, [checkAuth])

  useEffect(() => {
    // 如果已认证，连接WebSocket
    if (isAuthenticated) {
      connect()
    }
  }, [isAuthenticated, connect])

  if (loading) {
    return (
      <div className="flex items-center justify-center w-full h-screen bg-gradient-to-br from-slate-900 to-slate-800">
        <div className="flex flex-col items-center gap-4">
          <div className="w-12 h-12 border-4 border-slate-600 border-t-indigo-500 rounded-full animate-spin"></div>
          <p className="text-slate-400">加载中...</p>
        </div>
      </div>
    )
  }

  return (
    <Router>
      <Routes>
        {/* 登录页面 */}
        <Route path="/login" element={<LoginPage />} />

        {/* 受保护的路由 */}
        <Route element={<ProtectedRoute />}>
          <Route element={<Layout />}>
            <Route path="/" element={<DashboardPage />} />
            <Route path="/bots" element={<BotsPage />} />
            <Route path="/strategies" element={<StrategiesPage />} />
            <Route path="/trades" element={<TradesPage />} />
            <Route path="/settings" element={<SettingsPage />} />
          </Route>
        </Route>

        {/* 404页面 */}
        <Route path="*" element={<NotFoundPage />} />
      </Routes>
    </Router>
  )
}

export default App
