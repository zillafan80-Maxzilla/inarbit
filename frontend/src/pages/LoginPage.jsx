import React, { useState } from 'react'
import { useNavigate } from 'react-router-dom'
import { Lock, Mail, Eye, EyeOff, AlertCircle } from 'lucide-react'
import { useAuthStore } from '../store/authStore'
import './LoginPage.css'

/**
 * 登录页面组件
 * 设计风格：米色/棕色主题，简洁现代
 * 特点：居中对称布局、卡片式设计、平滑过渡动画
 */
export default function LoginPage() {
  const navigate = useNavigate()
  const { login } = useAuthStore()
  
  const [formData, setFormData] = useState({
    username: '',
    password: '',
  })
  
  const [showPassword, setShowPassword] = useState(false)
  const [loading, setLoading] = useState(false)
  const [error, setError] = useState('')
  const [rememberMe, setRememberMe] = useState(false)

  const handleChange = (e) => {
    const { name, value } = e.target
    setFormData(prev => ({
      ...prev,
      [name]: value
    }))
    // 清除错误信息
    if (error) setError('')
  }

  const handleSubmit = async (e) => {
    e.preventDefault()
    setLoading(true)
    setError('')

    try {
      // 验证输入
      if (!formData.username.trim()) {
        setError('请输入用户名')
        setLoading(false)
        return
      }
      if (!formData.password) {
        setError('请输入密码')
        setLoading(false)
        return
      }

      // 调用登录API
      const response = await fetch('/api/auth/login', {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify(formData),
      })

      if (!response.ok) {
        const data = await response.json()
        throw new Error(data.message || '登录失败')
      }

      const data = await response.json()
      
      // 保存token和用户信息
      login(data.token, data.user)
      
      // 记住我功能
      if (rememberMe) {
        localStorage.setItem('rememberMe', 'true')
        localStorage.setItem('username', formData.username)
      } else {
        localStorage.removeItem('rememberMe')
        localStorage.removeItem('username')
      }

      // 重定向到仪表板
      navigate('/')
    } catch (err) {
      setError(err.message || '登录失败，请检查用户名和密码')
      console.error('登录错误:', err)
    } finally {
      setLoading(false)
    }
  }

  // 从localStorage恢复记住的用户名
  React.useEffect(() => {
    const rememberMe = localStorage.getItem('rememberMe') === 'true'
    const username = localStorage.getItem('username')
    
    if (rememberMe && username) {
      setFormData(prev => ({ ...prev, username }))
      setRememberMe(true)
    }
  }, [])

  return (
    <div className="login-container">
      {/* 背景装饰 */}
      <div className="login-background">
        <div className="login-bg-shape login-bg-shape-1"></div>
        <div className="login-bg-shape login-bg-shape-2"></div>
        <div className="login-bg-shape login-bg-shape-3"></div>
      </div>

      {/* 主内容 */}
      <div className="login-content">
        {/* Logo和标题 */}
        <div className="login-header">
          <div className="login-logo">
            <div className="logo-icon">⚙️</div>
          </div>
          <h1 className="login-title">iNarbit</h1>
          <p className="login-subtitle">三角套利机器人控制系统</p>
        </div>

        {/* 登录表单卡片 */}
        <div className="login-card">
          <form onSubmit={handleSubmit} className="login-form">
            {/* 错误提示 */}
            {error && (
              <div className="login-error-alert">
                <AlertCircle size={18} />
                <span>{error}</span>
              </div>
            )}

            {/* 用户名输入框 */}
            <div className="login-form-group">
              <label htmlFor="username" className="login-label">
                用户名
              </label>
              <div className="login-input-wrapper">
                <Mail size={18} className="login-input-icon" />
                <input
                  id="username"
                  type="text"
                  name="username"
                  value={formData.username}
                  onChange={handleChange}
                  placeholder="输入用户名"
                  className="login-input"
                  disabled={loading}
                  autoComplete="username"
                />
              </div>
            </div>

            {/* 密码输入框 */}
            <div className="login-form-group">
              <label htmlFor="password" className="login-label">
                密码
              </label>
              <div className="login-input-wrapper">
                <Lock size={18} className="login-input-icon" />
                <input
                  id="password"
                  type={showPassword ? 'text' : 'password'}
                  name="password"
                  value={formData.password}
                  onChange={handleChange}
                  placeholder="输入密码"
                  className="login-input"
                  disabled={loading}
                  autoComplete="current-password"
                />
                <button
                  type="button"
                  onClick={() => setShowPassword(!showPassword)}
                  className="login-password-toggle"
                  disabled={loading}
                  aria-label={showPassword ? '隐藏密码' : '显示密码'}
                >
                  {showPassword ? (
                    <EyeOff size={18} />
                  ) : (
                    <Eye size={18} />
                  )}
                </button>
              </div>
            </div>

            {/* 记住我和忘记密码 */}
            <div className="login-form-footer">
              <label className="login-checkbox">
                <input
                  type="checkbox"
                  checked={rememberMe}
                  onChange={(e) => setRememberMe(e.target.checked)}
                  disabled={loading}
                />
                <span>记住我</span>
              </label>
              <a href="#" className="login-forgot-password">
                忘记密码？
              </a>
            </div>

            {/* 登录按钮 */}
            <button
              type="submit"
              className="login-button"
              disabled={loading}
            >
              {loading ? (
                <>
                  <span className="login-spinner"></span>
                  <span>登录中...</span>
                </>
              ) : (
                '登录'
              )}
            </button>
          </form>

          {/* 底部链接 */}
          <div className="login-footer">
            <p className="login-footer-text">
              首次使用？
              <a href="#" className="login-signup-link">
                联系管理员
              </a>
            </p>
          </div>
        </div>

        {/* 版本信息 */}
        <div className="login-version">
          <span>iNarbit v1.0.0</span>
          <span>•</span>
          <span>© 2024 All Rights Reserved</span>
        </div>
      </div>
    </div>
  )
}
