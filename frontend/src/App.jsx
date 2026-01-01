import { useState, useEffect } from 'react'
import './App.css'

function App() {
  const [count, setCount] = useState(0)
  const [status, setStatus] = useState('loading')

  useEffect(() => {
    // 检查后端健康状态
    fetch('/api/health')
      .then(res => res.json())
      .then(data => setStatus(data.status))
      .catch(() => setStatus('error'))
  }, [])

  return (
    <div className="min-h-screen bg-gradient-to-br from-amber-50 to-amber-100">
      <header className="bg-white shadow">
        <div className="max-w-7xl mx-auto px-4 py-6">
          <h1 className="text-3xl font-bold text-amber-900">
            iNarbit - 三角套利交易机器人
          </h1>
        </div>
      </header>

      <main className="max-w-7xl mx-auto px-4 py-12">
        <div className="grid grid-cols-1 md:grid-cols-3 gap-6">
          {/* 状态卡片 */}
          <div className="bg-white rounded-lg shadow p-6">
            <h2 className="text-lg font-semibold text-amber-900 mb-4">
              系统状态
            </h2>
            <p className="text-2xl font-bold text-amber-600">
              {status === 'ok' ? '✓ 正常' : status === 'loading' ? '加载中...' : '✗ 异常'}
            </p>
          </div>

          {/* 计数器 */}
          <div className="bg-white rounded-lg shadow p-6">
            <h2 className="text-lg font-semibold text-amber-900 mb-4">
              测试计数器
            </h2>
            <button
              onClick={() => setCount(count + 1)}
              className="bg-amber-600 hover:bg-amber-700 text-white font-bold py-2 px-4 rounded"
            >
              点击次数: {count}
            </button>
          </div>

          {/* 信息卡片 */}
          <div className="bg-white rounded-lg shadow p-6">
            <h2 className="text-lg font-semibold text-amber-900 mb-4">
              项目信息
            </h2>
            <p className="text-sm text-gray-600">
              版本: 1.0.0<br/>
              环境: Production<br/>
              状态: 部署成功
            </p>
          </div>
        </div>

        {/* 功能区 */}
        <div className="mt-12 bg-white rounded-lg shadow p-8">
          <h2 className="text-2xl font-bold text-amber-900 mb-6">
            主要功能
          </h2>
          <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
            <div className="border-l-4 border-amber-600 pl-4">
              <h3 className="font-semibold text-amber-900 mb-2">
                三角套利策略
              </h3>
              <p className="text-gray-600">
                支持3/4/5角套利策略，自动识别套利机会
              </p>
            </div>
            <div className="border-l-4 border-amber-600 pl-4">
              <h3 className="font-semibold text-amber-900 mb-2">
                实时监控
              </h3>
              <p className="text-gray-600">
                实时数据更新，1秒到1小时可配置
              </p>
            </div>
            <div className="border-l-4 border-amber-600 pl-4">
              <h3 className="font-semibold text-amber-900 mb-2">
                虚拟交易
              </h3>
              <p className="text-gray-600">
                先用虚拟资金测试策略，降低风险
              </p>
            </div>
            <div className="border-l-4 border-amber-600 pl-4">
              <h3 className="font-semibold text-amber-900 mb-2">
                性能优化
              </h3>
              <p className="text-gray-600">
                针对2vCPU/4GB服务器优化
              </p>
            </div>
          </div>
        </div>
      </main>
    </div>
  )
}

export default App
