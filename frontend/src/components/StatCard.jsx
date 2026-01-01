import React from 'react'
import { TrendingUp, TrendingDown } from 'lucide-react'
import './StatCard.css'

/**
 * 统计卡片组件
 * 用于展示关键指标
 */
export default function StatCard({ icon, label, value, trend, trendValue, color }) {
  const colorClass = `stat-card-${color}`
  const trendIcon = trend === 'up' ? <TrendingUp size={16} /> : <TrendingDown size={16} />
  const trendClass = trend === 'up' ? 'stat-card-trend-up' : trend === 'down' ? 'stat-card-trend-down' : 'stat-card-trend-neutral'

  return (
    <div className={`stat-card ${colorClass}`}>
      <div className="stat-card-header">
        <div className="stat-card-icon">
          {icon}
        </div>
        <div className={`stat-card-trend ${trendClass}`}>
          {trend !== 'neutral' && trendIcon}
          <span>{trendValue}</span>
        </div>
      </div>

      <div className="stat-card-content">
        <p className="stat-card-label">{label}</p>
        <h3 className="stat-card-value">{value}</h3>
      </div>

      <div className="stat-card-bar"></div>
    </div>
  )
}
