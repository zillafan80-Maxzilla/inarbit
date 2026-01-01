import { useEffect, useRef, useState, useCallback } from 'react'
import { useAuthStore } from '../store/authStore'

/**
 * WebSocket Hook
 * 用于实时数据推送和双向通信
 */
export const useWebSocket = () => {
  const { token } = useAuthStore()
  const wsRef = useRef(null)
  const reconnectTimeoutRef = useRef(null)
  const [data, setData] = useState(null)
  const [connected, setConnected] = useState(false)
  const [error, setError] = useState(null)

  // 获取WebSocket URL
  const getWebSocketUrl = useCallback(() => {
    const protocol = window.location.protocol === 'https:' ? 'wss:' : 'ws:'
    const host = window.location.host
    return `${protocol}//${host}/ws`
  }, [])

  // 连接WebSocket
  const connect = useCallback(() => {
    if (!token) return

    try {
      const url = getWebSocketUrl()
      const ws = new WebSocket(url)

      ws.onopen = () => {
        console.log('WebSocket已连接')
        setConnected(true)
        setError(null)

        // 发送认证信息
        ws.send(JSON.stringify({
          type: 'auth',
          token,
        }))
      }

      ws.onmessage = (event) => {
        try {
          const message = JSON.parse(event.data)
          setData(message)
        } catch (err) {
          console.error('解析WebSocket消息失败:', err)
        }
      }

      ws.onerror = (error) => {
        console.error('WebSocket错误:', error)
        setError('WebSocket连接错误')
        setConnected(false)
      }

      ws.onclose = () => {
        console.log('WebSocket已断开')
        setConnected(false)

        // 自动重连（延迟3秒）
        reconnectTimeoutRef.current = setTimeout(() => {
          console.log('正在重新连接WebSocket...')
          connect()
        }, 3000)
      }

      wsRef.current = ws
    } catch (err) {
      console.error('WebSocket连接失败:', err)
      setError(err.message)
      setConnected(false)
    }
  }, [token, getWebSocketUrl])

  // 断开连接
  const disconnect = useCallback(() => {
    if (reconnectTimeoutRef.current) {
      clearTimeout(reconnectTimeoutRef.current)
    }

    if (wsRef.current) {
      wsRef.current.close()
      wsRef.current = null
    }

    setConnected(false)
  }, [])

  // 发送消息
  const send = useCallback((message) => {
    if (wsRef.current && wsRef.current.readyState === WebSocket.OPEN) {
      wsRef.current.send(JSON.stringify(message))
    } else {
      console.warn('WebSocket未连接')
    }
  }, [])

  // 订阅特定事件
  const subscribe = useCallback((eventType, callback) => {
    const handleMessage = (event) => {
      try {
        const message = JSON.parse(event.data)
        if (message.type === eventType) {
          callback(message)
        }
      } catch (err) {
        console.error('处理WebSocket消息失败:', err)
      }
    }

    if (wsRef.current) {
      wsRef.current.addEventListener('message', handleMessage)

      // 返回取消订阅函数
      return () => {
        if (wsRef.current) {
          wsRef.current.removeEventListener('message', handleMessage)
        }
      }
    }

    return () => {}
  }, [])

  // 自动连接和清理
  useEffect(() => {
    if (token) {
      connect()
    }

    return () => {
      disconnect()
    }
  }, [token, connect, disconnect])

  return {
    data,
    connected,
    error,
    send,
    subscribe,
    connect,
    disconnect,
  }
}
