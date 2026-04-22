import request from '@/utils/request'

// 统计数据接口
export interface Statistics {
  memberCount: number
  memberTrend: string
  todayRecharge: number
  rechargeTrend: string
  todayConsumption: number
  consumptionTrend: string
  activeCenters: number
  centerTrend: string
}

// 待办事项接口
export interface TodoItem {
  id: string
  title: string
  description: string
  type: 'warning' | 'error' | 'info'
  icon: string
}

export interface Todos {
  pendingApprovals: {
    count: number
    description: string
  }
  expiringCards: {
    count: number
    description: string
  }
}

// 充值趋势接口
export interface RechargeTrend {
  dates: string[]
  values: number[]
}

export interface ChartDataPoint {
  date: string
  value: number
}

// 获取统计数据
export function getStatistics() {
  return request.get('/dashboard/statistics')
}

// 获取待办事项
export function getTodos() {
  return request.get('/dashboard/todos')
}

// 获取充值趋势
export function getRechargeTrends(days: number = 7) {
  return request.get(`/dashboard/recharge-trends?days=${days}`)
}
