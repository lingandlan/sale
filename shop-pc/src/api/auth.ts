import request from '@/utils/request'

// 登录接口
export interface LoginParams {
  phone: string
  password: string
}

export interface LoginResponse {
  access_token: string
  refresh_token: string
}

export const login = (data: LoginParams) => {
  return request({
    url: '/auth/login',
    method: 'POST',
    data
  })
}

export interface UserInfo {
  id: number
  phone: string
  name: string
  role: string
  center_id?: number | null
  center_name?: string | null
  status: number
}

export const getUserInfo = () => {
  return request.get('/user/info')
}
