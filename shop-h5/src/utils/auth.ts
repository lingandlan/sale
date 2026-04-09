/**
 * 认证工具函数
 */
import { refreshTokenApi } from '@/api/user'

/**
 * 获取 Token
 */
export function getToken(): string | null {
  return uni.getStorageSync('token') || null
}

/**
 * 获取 Refresh Token
 */
export function getRefreshToken(): string | null {
  return uni.getStorageSync('refresh_token') || null
}

/**
 * 刷新 Token
 */
export async function refreshToken(): Promise<string | null> {
  const refreshToken = getRefreshToken()
  if (!refreshToken) {
    return null
  }

  try {
    const res = await refreshTokenApi(refreshToken)
    // 更新本地存储
    uni.setStorageSync('token', res.access_token)
    uni.setStorageSync('refresh_token', res.refresh_token)
    return res.access_token
  } catch {
    // 刷新失败
    return null
  }
}

/**
 * 检查是否已登录
 */
export function isAuthenticated(): boolean {
  return !!getToken()
}

/**
 * 清除认证信息
 */
export function clearAuth() {
  uni.removeStorageSync('token')
  uni.removeStorageSync('refresh_token')
}
