import request from '@/utils/request'

// 用户管理接口

export interface AdminUserListParams {
  keyword?: string
  role?: string
  status?: number
  page: number
  page_size: number
}

export interface AdminUserItem {
  id: number
  username: string
  phone: string
  name: string
  role: string
  center_id?: number
  center_name?: string
  status: number
  last_login_at?: string
  created_at?: string
}

export interface AdminUserListResponse {
  items: AdminUserItem[]
  total: number
  page: number
  page_size: number
}

export interface AdminUserCreateData {
  username: string
  phone: string
  name: string
  role: string
  center_id?: number
  password?: string
}

export interface AdminUserUpdateData {
  username?: string
  name?: string
  phone?: string
  role?: string
  center_id?: number
}

// 获取用户列表
export function getAdminUsers(params: AdminUserListParams) {
  return request.get<AdminUserListResponse>('/admin/users', { params })
}

// 创建用户
export function createAdminUser(data: AdminUserCreateData) {
  return request.post<AdminUserItem>('/admin/users', data)
}

// 更新用户
export function updateAdminUser(id: number, data: AdminUserUpdateData) {
  return request.put<AdminUserItem>(`/admin/users/${id}`, data)
}

// 重置用户密码
export function resetUserPassword(id: number, newPassword: string) {
  return request.post<{ success: boolean }>(`/admin/users/${id}/reset-password`, { new_password: newPassword })
}

// 启用/禁用用户
export function toggleUserStatus(id: number, status: number) {
  return request.put<{ success: boolean }>(`/admin/users/${id}/status`, { status })
}

// 删除用户
export function deleteAdminUser(id: number) {
  return request.delete<{ success: boolean }>(`/admin/users/${id}`)
}
