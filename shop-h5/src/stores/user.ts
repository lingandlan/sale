/**
 * Pinia Store - 用户模块
 */
import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import type { User } from '@/api/user'
import { login as loginApi, logout as logoutApi, getUserInfo as getUserInfoApi } from '@/api/user'

export const useUserStore = defineStore('user', () => {
  // State
  const token = ref<string>(uni.getStorageSync('token') || '')
  const userInfo = ref<User | null>(null)

  // Getters
  const isLoggedIn = computed(() => !!token.value)
  const userId = computed(() => userInfo.value?.id)
  const username = computed(() => userInfo.value?.username || '')
  const role = computed(() => userInfo.value?.role || 0)

  // 判断是否是管理员
  const isAdmin = computed(() => role.value === 2)
  // 判断是否是商家
  const isMerchant = computed(() => role.value === 1)

  // Actions
  async function login(username: string, password: string) {
    const res = await loginApi({ username, password })
    token.value = res.access_token
    userInfo.value = res.user
    // 持久化
    uni.setStorageSync('token', res.access_token)
    uni.setStorageSync('refresh_token', res.refresh_token)
    return res
  }

  async function logout() {
    try {
      await logoutApi()
    } catch {
      // ignore error
    } finally {
      // 清理状态
      token.value = ''
      userInfo.value = null
      uni.removeStorageSync('token')
      uni.removeStorageSync('refresh_token')
    }
  }

  async function fetchUserInfo() {
    if (!token.value) return null
    try {
      const info = await getUserInfoApi()
      userInfo.value = info
      return info
    } catch {
      // token 失效
      await logout()
      return null
    }
  }

  return {
    // State
    token,
    userInfo,
    // Getters
    isLoggedIn,
    userId,
    username,
    role,
    isAdmin,
    isMerchant,
    // Actions
    login,
    logout,
    fetchUserInfo
  }
})
