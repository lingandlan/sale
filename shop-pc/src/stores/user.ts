import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import { getUserInfo, type UserInfo } from '@/api/auth'

export const useUserStore = defineStore('user', () => {
  const userInfo = ref<UserInfo | null>(null)

  const isSuperAdmin = computed(() => userInfo.value?.role === 'super_admin')
  const isHQAdmin = computed(() => userInfo.value?.role === 'hq_admin')
  const isFinance = computed(() => userInfo.value?.role === 'finance')
  const userCenterId = computed(() => userInfo.value?.center_id ?? null)
  const userCenterName = computed(() => userInfo.value?.center_name ?? '')

  /** 超管/总部/财务 可以选择所有中心 */
  const canSelectAllCenters = computed(() =>
    isSuperAdmin.value || isHQAdmin.value || isFinance.value
  )

  const displayName = computed(() => {
    if (!userInfo.value) return ''
    return userInfo.value.name || userInfo.value.phone
  })

  async function fetchUserInfo() {
    const res = await getUserInfo()
    userInfo.value = res.data
  }

  function clear() {
    userInfo.value = null
  }

  return {
    userInfo,
    isSuperAdmin,
    isHQAdmin,
    isFinance,
    userCenterId,
    userCenterName,
    canSelectAllCenters,
    displayName,
    fetchUserInfo,
    clear
  }
})
