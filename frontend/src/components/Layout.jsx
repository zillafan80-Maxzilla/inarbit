import React, { useState } from 'react'
import { Outlet, useNavigate, useLocation } from 'react-router-dom'
import { Menu, X, LogOut, Settings, Home, Zap, TrendingUp, BarChart3 } from 'lucide-react'
import { useAuthStore } from '../store/authStore'
import './Layout.css'

/**
 * 主布局组件
 * 包含导航栏、侧边栏等
 */
export default function Layout() {
  const navigate = useNavigate()
  const location = useLocation()
  const { user, logout } = useAuthStore()
  const [sidebarOpen, setSidebarOpen] = useState(true)
  const [userMenuOpen, setUserMenuOpen] = useState(false)

  const handleLogout = () => {
    logout()
    navigate('/login')
  }

  const isActive = (path) => location.pathname === path

  const navItems = [
    { path: '/', label: '仪表板', icon: Home },
    { path: '/bots', label: '机器人', icon: Zap },
    { path: '/strategies', label: '策略', icon: TrendingUp },
    { path: '/trades', label: '交易', icon: BarChart3 },
  ]

  return (
    <div className="layout-container">
      {/* 侧边栏 */}
      <aside className={`layout-sidebar ${sidebarOpen ? 'open' : 'closed'}`}>
        <div className="layout-sidebar-header">
          <div className="layout-logo">
            <span className="layout-logo-icon">⚙️</span>
            <span className="layout-logo-text">iNarbit</span>
          </div>
          <button
            className="layout-sidebar-toggle"
            onClick={() => setSidebarOpen(!sidebarOpen)}
          >
            {sidebarOpen ? <X size={20} /> : <Menu size={20} />}
          </button>
        </div>

        <nav className="layout-nav">
          {navItems.map(item => (
            <a
              key={item.path}
              href={item.path}
              className={`layout-nav-item ${isActive(item.path) ? 'active' : ''}`}
              onClick={(e) => {
                e.preventDefault()
                navigate(item.path)
              }}
            >
              <item.icon size={20} />
              <span className="layout-nav-label">{item.label}</span>
            </a>
          ))}
        </nav>

        <div className="layout-sidebar-footer">
          <button
            className="layout-nav-item"
            onClick={() => navigate('/settings')}
          >
            <Settings size={20} />
            <span className="layout-nav-label">设置</span>
          </button>
        </div>
      </aside>

      {/* 主内容区域 */}
      <div className="layout-main">
        {/* 顶部导航栏 */}
        <header className="layout-header">
          <button
            className="layout-mobile-toggle"
            onClick={() => setSidebarOpen(!sidebarOpen)}
          >
            <Menu size={24} />
          </button>

          <div className="layout-header-spacer"></div>

          {/* 用户菜单 */}
          <div className="layout-user-menu">
            <button
              className="layout-user-button"
              onClick={() => setUserMenuOpen(!userMenuOpen)}
            >
              <div className="layout-user-avatar">
                {user?.username?.charAt(0).toUpperCase() || 'U'}
              </div>
              <span className="layout-user-name">{user?.username || '用户'}</span>
            </button>

            {userMenuOpen && (
              <div className="layout-user-dropdown">
                <div className="layout-user-info">
                  <p className="layout-user-email">{user?.email || '未设置'}</p>
                </div>
                <button
                  className="layout-user-logout"
                  onClick={handleLogout}
                >
                  <LogOut size={16} />
                  <span>登出</span>
                </button>
              </div>
            )}
          </div>
        </header>

        {/* 页面内容 */}
        <main className="layout-content">
          <Outlet />
        </main>
      </div>
    </div>
  )
}
