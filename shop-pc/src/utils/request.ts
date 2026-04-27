import axios from 'axios'
import type { AxiosInstance, AxiosResponse, InternalAxiosRequestConfig, AxiosRequestConfig } from 'axios'
import { ElMessage } from 'element-plus'

// 响应拦截器 unwrap 了 response.data，实际返回 { code, data, message }
export interface ApiResponse<T = any> {
  code: number
  data: T
  message?: string
}

// 类型覆盖：拦截器返回 ApiResponse<T> 而非 AxiosResponse<T>
export interface TypedAxiosInstance {
  <T = any>(config: AxiosRequestConfig): Promise<ApiResponse<T>>
  get<T = any>(url: string, config?: AxiosRequestConfig): Promise<ApiResponse<T>>
  post<T = any>(url: string, data?: any, config?: AxiosRequestConfig): Promise<ApiResponse<T>>
  put<T = any>(url: string, data?: any, config?: AxiosRequestConfig): Promise<ApiResponse<T>>
  delete<T = any>(url: string, config?: AxiosRequestConfig): Promise<ApiResponse<T>>
  interceptors: AxiosInstance['interceptors']
  defaults: AxiosInstance['defaults']
}

const isMock = import.meta.env.VITE_USE_MOCK === 'true'

// Mock API ID 映射表（路径+方法 → Apifox API ID）
const mockApiIdMap: Record<string, string> = {
  // Auth
  'POST /auth/login': '440994065',
  'POST /auth/logout': '440994066',
  'POST /auth/refresh': '440994067',
  // User
  'GET /user/info': '440994068',
  'PUT /user/info': '440994069',
  'POST /user/change-password': '440994070',
  // Admin
  'GET /admin/users': '440996304',
  'POST /admin/users/{id}/reset-password': '440996305',
  'PUT /admin/users/{id}/status': '440996306',
  'DELETE /admin/users/{id}': '440996307',
  // Dashboard
  'GET /dashboard/statistics': '441722169',
  'GET /dashboard/todos': '441722170',
  'GET /dashboard/recharge-trends': '441722171',
  // Recharge
  'POST /recharge/b-apply': '441722198',
  'GET /recharge/b-approval': '441722199',
  'GET /recharge/b-approval/{id}': '441722200',
  'POST /recharge/b-approval/action': '441722201',
  'POST /recharge/c-entry': '441722202',
  'GET /recharge/records': '441722203',
  'GET /recharge/records/{id}': '441722204',
  // Card
  'GET /card/verify/{cardNo}': '441722207',
  'POST /card/consume': '441722208',
  'GET /card/list': '441722209',
  'POST /card/issue': '441722210',
  'GET /card/detail/{cardNo}': '441722211',
  'POST /card/{cardNo}/status': '441722212',
  'GET /card/stats': '441722213',
}

// 创建axios实例
const service: AxiosInstance = axios.create({
  baseURL: '/api/v1',
  timeout: 15000
})

// 请求拦截器
service.interceptors.request.use(
  (config: InternalAxiosRequestConfig) => {
    const token = localStorage.getItem('access_token')
    if (token) {
      config.headers.Authorization = `Bearer ${token}`
    }

    // Mock 模式下自动注入 apifoxApiId
    if (isMock) {
      const method = (config.method || 'get').toUpperCase()
      const path = config.url || ''

      // 先尝试精确匹配
      let key = `${method} ${path}`
      let apiId = mockApiIdMap[key]

      // 精确匹配失败时，尝试模板匹配（将实际值替换为模板参数）
      if (!apiId) {
        for (const [pattern, id] of Object.entries(mockApiIdMap)) {
          const [pMethod, pPath] = pattern.split(' ', 2)
          if (pMethod !== method) continue
          // 将模板参数如 {cardNo} 转换为正则
          const regexStr = '^' + pPath.replace(/\{[^}]+\}/g, '[^/]+') + '$'
          if (new RegExp(regexStr).test(path)) {
            apiId = id
            break
          }
        }
      }

      if (apiId) {
        config.params = { ...config.params, apifoxApiId: apiId }
      }
    }

    return config
  },
  (error) => {
    return Promise.reject(error)
  }
)

// 错误码 → 中文提示兜底映射
const errorCodeMap: Record<number, string> = {
  400: '请求参数有误',
  401: '请先登录',
  403: '无权限访问',
  404: '请求的资源不存在',
  409: '数据冲突，请刷新后重试',
  500: '服务器内部错误',
}

// 响应拦截器
service.interceptors.response.use(
  (response: AxiosResponse) => {
    const res = response.data as any
    // 兼容 code: 0（前端约定）和 code: 200（Apifox 默认）
    if (res.code === 0 || res.code === 200) {
      // 统一转换为 code: 0 格式
      if (res.code === 200) {
        res.code = 0
      }
      return res
    }

    // 业务层 401：token 无效/过期，尝试刷新
    if (res.code === 401) {
      return handleTokenExpired(response.config)
    }

    const msg = res.message || errorCodeMap[res.code] || '请求失败'
    ElMessage.error(msg)
    return Promise.reject(new Error(msg))
  },
  async (error) => {
    const status = error.response?.status

    // HTTP 401：未认证，尝试刷新 token
    if (status === 401) {
      try {
        const config = error.response?.config
        return await handleTokenExpired(config)
      } catch {
        // 刷新也失败，走正常登出
      }
    }

    const data = error.response?.data
    const msg = data?.message || errorCodeMap[status] || error.message || '网络错误'
    ElMessage.error(msg)
    return Promise.reject(error)
  }
)

// ==================== Token 刷新机制 ====================

let isRefreshing = false
let pendingRequests: ((token: string) => void)[] = []

/** 用 refresh_token 换新的 access_token */
async function refreshAccessToken(): Promise<string> {
  const refreshToken = localStorage.getItem('refresh_token')
  if (!refreshToken) {
    return Promise.reject(new Error('no refresh token'))
  }

  // 用独立的 axios 实例发请求，避免触发当前实例的拦截器
  const { data } = await axios.post('/api/v1/auth/refresh', {
    refresh_token: refreshToken
  })

  const res = data as any
  if (res.code !== 0 && res.code !== 200) {
    return Promise.reject(new Error('refresh failed'))
  }

  const newAccessToken: string = res.data.access_token
  const newRefreshToken: string = res.data.refresh_token

  localStorage.setItem('access_token', newAccessToken)
  localStorage.setItem('refresh_token', newRefreshToken)

  return newAccessToken
}

/** 处理 token 过期：刷新 token 并重试原始请求 */
async function handleTokenExpired(config: any) {
  if (!config) {
    clearAuthAndRedirect()
    return Promise.reject(new Error('请先登录'))
  }

  if (isRefreshing) {
    // 正在刷新中，排队等待
    return new Promise((resolve) => {
      pendingRequests.push((token: string) => {
        config.headers.Authorization = `Bearer ${token}`
        resolve(service(config))
      })
    })
  }

  isRefreshing = true

  try {
    const newToken = await refreshAccessToken()

    // 通知排队的请求
    pendingRequests.forEach(cb => cb(newToken))
    pendingRequests = []

    // 重试原始请求
    config.headers.Authorization = `Bearer ${newToken}`
    return service(config)
  } catch {
    // 刷新失败，清除所有排队请求并跳转登录
    pendingRequests = []
    clearAuthAndRedirect()
    return Promise.reject(new Error('请先登录'))
  } finally {
    isRefreshing = false
  }
}

function clearAuthAndRedirect() {
  // 避免重复跳转
  if (window.location.pathname === '/login') return
  localStorage.removeItem('access_token')
  localStorage.removeItem('refresh_token')
  window.location.href = '/login'
}

export default service as TypedAxiosInstance

// extractErrorMessage 统一从错误对象中提取可展示的错误信息
export function extractErrorMessage(err: any, fallback = '操作失败'): string {
  return err?.response?.data?.message || err?.message || fallback
}
