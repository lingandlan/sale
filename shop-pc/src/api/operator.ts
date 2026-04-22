import request from '@/utils/request'

// 操作员管理接口

export interface OperatorItem {
  id: string
  name: string
  phone: string
  password?: string
  center_id: string
  role: string
  status: string
  created_at?: string
}

export function getOperatorList() {
  return request.get('/operator')
}

export function createOperator(data: Partial<OperatorItem>) {
  return request.post('/operator', data)
}

export function updateOperator(id: string, data: Partial<OperatorItem>) {
  return request.put(`/operator/${id}`, data)
}

export function deleteOperator(id: string) {
  return request.delete<{ data: { success: boolean } }>(`/operator/${id}`)
}
