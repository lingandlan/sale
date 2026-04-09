/**
 * API 请求封装
 * 基于 Axios + UniApp
 */

import axios, { AxiosInstance, AxiosRequestConfig, AxiosResponse } from 'axios'
import { getToken, refreshToken } from '@/utils/auth'

// API 基础配置
const BASE_URL = import.meta.env.VITE_API_BASE_URL || 'http://localhost:8080/api/v1'

// 创建 axios 实例
const instance: AxiosInstance = axios.create({
  baseURL: BASE_URL,
  timeout: 30000,
  headers: {
    'Content-Type': 'application/json'
  }
})

// 请求拦截器
instance.interceptors.request.use(
  (config) => {
    // 添加 Token
    const token = getToken()
    if (token) {
      config.headers.Authorization = `Bearer ${token}`
    }
    return config
  },
  (error) => {
    return Promise.reject(error)
  }
)

// 响应拦截器
instance.interceptors.response.use(
  (response: AxiosResponse) => {
    const res = response.data

    // 根据业务码判断
    if (res.code === 0) {
      return res.data
    }

    // Token 过期，尝试刷新
    if (res.code === 401) {
      return handleUnauthorized()
    }

    // 其他业务错误
    uni.showToast({
      title: res.message || '请求失败',
      icon: 'none'
    })
    return Promise.reject(new Error(res.message || '请求失败'))
  },
  async (error) => {
    if (error.response) {
      // 服务器返回错误状态码
      switch (error.response.status) {
        case 400:
          uni.showToast({ title: '参数错误', icon: 'none' })
          break
        case 401:
          await handleUnauthorized()
          break
        case 403:
          uni.showToast({ title: '无权限', icon: 'none' })
          break
        case 404:
          uni.showToast({ title: '资源不存在', icon: 'none' })
          break
        case 500:
          uni.showToast({ title: '服务器错误', icon: 'none' })
          break
        default:
          uni.showToast({ title: '网络错误', icon: 'none' })
      }
    } else if (error.request) {
      uni.showToast({ title: '网络连接失败', icon: 'none' })
    } else {
      uni.showToast({ title: '请求配置错误', icon: 'none' })
    }
    return Promise.reject(error)
  }
)

// 处理 Token 过期
async function handleUnauthorized() {
  try {
    const newToken = await refreshToken()
    if (newToken) {
      // 重试当前请求
      return instance.request(error.config!)
    }
  } catch {
    // 刷新失败，跳转登录
    uni.removeStorageSync('token')
    uni.removeStorageSync('refresh_token')
    uni.reLaunch({ url: '/pages/login/login' })
  }
}

// 封装请求方法
export function request<T = any>(config: AxiosRequestConfig): Promise<T> {
  return instance.request(config)
}

// 常用请求方法
export const api = {
  get<T = any>(url: string, params?: any, config?: AxiosRequestConfig): Promise<T> {
    return request({ method: 'GET', url, params, ...config })
  },

  post<T = any>(url: string, data?: any, config?: AxiosRequestConfig): Promise<T> {
    return request({ method: 'POST', url, data, ...config })
  },

  put<T = any>(url: string, data?: any, config?: AxiosRequestConfig): Promise<T> {
    return request({ method: 'PUT', url, data, ...config })
  },

  delete<T = any>(url: string, params?: any, config?: AxiosRequestConfig): Promise<T> {
    return request({ method: 'DELETE', url, params, ...config })
  }
}

export default api
