/**
 * 用户相关 API
 */
import api from './request'

// 用户类型
export interface User {
  id: number
  username: string
  nickname: string
  email: string
  avatar?: string
  role: number
  status: number
  created_at: string
}

// 登录参数
export interface LoginParams {
  username: string
  password: string
}

// 登录响应
export interface LoginResponse {
  access_token: string
  refresh_token: string
  user: User
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

// 更新用户信息
export function updateUser(data: Partial<User>) {
  return api.put('/user/info', data)
}

// 退出登录
export function logout() {
  return api.post('/auth/logout')
}
