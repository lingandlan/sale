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

export interface TodoItem {
  id: string
  title: string
  description: string
  type: 'warning' | 'error' | 'info'
  icon: string
  count?: number
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

export interface RechargeTrend {
  dates: string[]
  values: number[]
}

export interface QuickAction {
  icon: string
  text: string
  route: string
  background: string
  border: string
}
