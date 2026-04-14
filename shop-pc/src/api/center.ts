import request from '@/utils/request'

// 充值中心管理接口

export interface CenterItem {
  id: string
  name: string
  code: string
  province: string
  city: string
  district: string
  address: string
  managerId: string
  managerName?: string
  managerPhone?: string
  phone: string
  balance: number
  totalRecharge?: number
  totalConsumed?: number
  status: string
  created_at?: string
}

export function getCenterList() {
  return request.get<{ data: CenterItem[] }>('/center')
}

export function createCenter(data: Partial<CenterItem>) {
  return request.post<{ data: CenterItem }>('/center', data)
}

export function updateCenter(id: string, data: Partial<CenterItem>) {
  return request.put<{ data: CenterItem }>(`/center/${id}`, data)
}

export function deleteCenter(id: string) {
  return request.delete<{ data: { success: boolean } }>(`/center/${id}`)
}
