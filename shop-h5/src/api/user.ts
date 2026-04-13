/**
 * 用户相关 API
 */
import api from './request'

// 用户类型（太积堂系统）
export interface User {
  id: number
  phone: string
  name: string
  role: string  // super_admin, admin, operator
  center_id?: number
  center_name?: string
  status: number  // 0=禁用, 1=启用
  last_login_at?: string
  created_at: string
}

// 登录参数（太积堂系统）
export interface LoginParams {
  phone: string
  password: string
}

// 登录响应（太积堂系统）
export interface LoginResponse {
  access_token: string
  refresh_token: string
  expires_in: number
}

// 修改密码参数
export interface ChangePasswordParams {
  old_password: string
  new_password: string
}

// 获取用户信息
export function getUserInfo() {
  return api.get<User>('/user/info')
}

// 登录
export function login(params: LoginParams) {
  return api.post<LoginResponse>('/auth/login', params)
}

// 刷新 Token
export function refreshTokenApi(refreshToken: string) {
  return api.post<{ access_token: string; refresh_token: string }>('/auth/refresh', {
    refresh_token: refreshToken
  })
}

// 修改密码
export function changePassword(params: ChangePasswordParams) {
  return api.post('/user/change-password', params)
}

// 更新用户信息
export function updateUser(data: { name?: string }) {
  return api.put('/user/info', data)
}

// 退出登录
export function logout() {
  return api.post('/auth/logout')
}
