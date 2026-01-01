import { create } from 'zustand'
import { jwtDecode } from 'jwt-decode'

/**
 * 认证状态管理
 * 使用Zustand管理全局认证状态
 */
export const useAuthStore = create((set) => ({
  // 状态
  token: localStorage.getItem('token') || null,
  user: JSON.parse(localStorage.getItem('user') || 'null'),
  isAuthenticated: !!localStorage.getItem('token'),
  loading: false,
  error: null,

  // 登录
  login: (token, user) => {
    localStorage.setItem('token', token)
    localStorage.setItem('user', JSON.stringify(user))
    
    set({
      token,
      user,
      isAuthenticated: true,
      error: null,
    })
  },

  // 登出
  logout: () => {
    localStorage.removeItem('token')
    localStorage.removeItem('user')
    localStorage.removeItem('rememberMe')
    localStorage.removeItem('username')
    
    set({
      token: null,
      user: null,
      isAuthenticated: false,
      error: null,
    })
  },

  // 检查认证状态
  checkAuth: async () => {
    const token = localStorage.getItem('token')
    
    if (!token) {
      set({
        isAuthenticated: false,
        token: null,
        user: null,
      })
      return
    }

    try {
      // 验证token是否过期
      const decoded = jwtDecode(token)
      const isExpired = decoded.exp * 1000 < Date.now()

      if (isExpired) {
        // Token已过期，清除认证信息
        localStorage.removeItem('token')
        localStorage.removeItem('user')
        set({
          isAuthenticated: false,
          token: null,
          user: null,
        })
        return
      }

      // Token有效，验证服务器
      const response = await fetch('/api/auth/verify', {
        headers: {
          'Authorization': `Bearer ${token}`,
        },
      })

      if (response.ok) {
        const user = JSON.parse(localStorage.getItem('user') || 'null')
        set({
          isAuthenticated: true,
          token,
          user,
          error: null,
        })
      } else {
        // 服务器验证失败
        localStorage.removeItem('token')
        localStorage.removeItem('user')
        set({
          isAuthenticated: false,
          token: null,
          user: null,
        })
      }
    } catch (error) {
      console.error('认证检查错误:', error)
      set({
        isAuthenticated: false,
        token: null,
        user: null,
        error: error.message,
      })
    }
  },

  // 更新用户信息
  updateUser: (user) => {
    localStorage.setItem('user', JSON.stringify(user))
    set({ user })
  },

  // 设置错误
  setError: (error) => {
    set({ error })
  },

  // 清除错误
  clearError: () => {
    set({ error: null })
  },
}))
