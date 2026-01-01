import React, { useState } from 'react'
import { Play, Pause, Settings, MoreVertical, TrendingUp, Clock } from 'lucide-react'
import './BotCard.css'

/**
 * 机器人卡片组件
 * 展示单个机器人的信息和控制按钮
 */
export default function BotCard({ bot }) {
  const [showMenu, setShowMenu] = useState(false)

  const getStrategyLabel = (strategy) => {
    const labels = {
      triangular: '三角套利',
      quadrangular: '四角套利',
      pentagonal: '五角套利',
    }
    return labels[strategy] || strategy
  }

  const getStatusLabel = (status) => {
    const labels = {
      running: '运行中',
      stopped: '已停止',
      paused: '已暂停',
      error: '错误',
    }
    return labels[status] || status
  }

  const getStatusColor = (status) => {
    const colors = {
      running: 'bot-card-status-running',
      stopped: 'bot-card-status-stopped',
      paused: 'bot-card-status-paused',
      error: 'bot-card-status-error',
    }
    return colors[status] || 'bot-card-status-stopped'
  }

  return (
    <div className="bot-card">
      {/* 卡片头部 */}
      <div className="bot-card-header">
        <div className="bot-card-title-section">
          <h3 className="bot-card-title">{bot.name}</h3>
          <span className={`bot-card-strategy ${getStatusColor(bot.status)}`}>
            {getStrategyLabel(bot.strategy)}
          </span>
        </div>

        <div className="bot-card-menu">
          <button
            className="bot-card-menu-btn"
            onClick={() => setShowMenu(!showMenu)}
          >
            <MoreVertical size={18} />
          </button>

          {showMenu && (
            <div className="bot-card-dropdown">
              <button className="bot-card-dropdown-item">编辑</button>
              <button className="bot-card-dropdown-item">查看日志</button>
              <button className="bot-card-dropdown-item bot-card-dropdown-danger">
                删除
              </button>
            </div>
          )}
        </div>
      </div>

      {/* 状态指示 */}
      <div className="bot-card-status">
        <div className={`bot-card-status-indicator ${getStatusColor(bot.status)}`}></div>
        <span className="bot-card-status-text">
          {getStatusLabel(bot.status)}
        </span>
      </div>

      {/* 统计数据 */}
      <div className="bot-card-stats">
        <div className="bot-card-stat">
          <div className="bot-card-stat-icon">
            <TrendingUp size={16} />
          </div>
          <div className="bot-card-stat-content">
            <p className="bot-card-stat-label">收益</p>
            <p className="bot-card-stat-value">
              ${bot.profit.toFixed(2)}
            </p>
          </div>
        </div>

        <div className="bot-card-stat">
          <div className="bot-card-stat-icon">
            <Clock size={16} />
          </div>
          <div className="bot-card-stat-content">
            <p className="bot-card-stat-label">运行时间</p>
            <p className="bot-card-stat-value">{bot.uptime}</p>
          </div>
        </div>
      </div>

      {/* 交易信息 */}
      <div className="bot-card-trades">
        <div className="bot-card-trade-item">
          <span className="bot-card-trade-label">总交易数</span>
          <span className="bot-card-trade-value">{bot.trades}</span>
        </div>
        <div className="bot-card-trade-divider"></div>
        <div className="bot-card-trade-item">
          <span className="bot-card-trade-label">成功率</span>
          <span className="bot-card-trade-value">
            {Math.round((Math.random() * 30 + 70))}%
          </span>
        </div>
      </div>

      {/* 操作按钮 */}
      <div className="bot-card-actions">
        <button className="bot-card-action-btn bot-card-action-primary">
          {bot.status === 'running' ? (
            <>
              <Pause size={16} />
              <span>暂停</span>
            </>
          ) : (
            <>
              <Play size={16} />
              <span>启动</span>
            </>
          )}
        </button>

        <button className="bot-card-action-btn bot-card-action-secondary">
          <Settings size={16} />
          <span>配置</span>
        </button>
      </div>
    </div>
  )
}
