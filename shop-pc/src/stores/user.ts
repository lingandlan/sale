import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import { getUserInfo, type UserInfo } from '@/api/auth'

const USER_INFO_KEY = 'user_info'

/** 角色 → 权限码映射表 */
const RolePermissions: Record<string, string[]> = {
  super_admin: ['*'],
  hq_admin: [
    'dashboard', 'recharge:b', 'recharge:c', 'recharge:records',
    'card:inventory', 'card:issue', 'card:verify', 'card:manage', 'card:stats',
    'center:manage', 'operator:manage', 'user:manage',
  ],
  finance: ['dashboard', 'recharge:records', 'card:stats'],
  center_admin: [
    'dashboard', 'recharge:c', 'recharge:records',
    'operator:manage',
    'card:issue', 'card:verify', 'card:manage', 'card:stats',
  ],
  operator: [
    'dashboard', 'recharge:c', 'recharge:records',
    'card:issue', 'card:verify', 'card:manage', 'card:stats',
  ],
}

export const useUserStore = defineStore('user', () => {
  const userInfo = ref<UserInfo | null>(null)

  // 初始化：从 localStorage 恢复 userInfo
  const raw = localStorage.getItem(USER_INFO_KEY)
  if (raw) {
    try { userInfo.value = JSON.parse(raw) } catch {}
  }

  const isSuperAdmin = computed(() => userInfo.value?.role === 'super_admin')
  const isHQAdmin = computed(() => userInfo.value?.role === 'hq_admin')
  const isFinance = computed(() => userInfo.value?.role === 'finance')
  const userCenterId = computed(() => userInfo.value?.center_id ?? null)
  const userCenterName = computed(() => userInfo.value?.center_name ?? '')

  /** 当前角色的权限码列表 */
  const permissions = computed(() => {
    const role = userInfo.value?.role
    if (!role) return []
    return RolePermissions[role] ?? []
  })

  /** 判断当前用户是否拥有指定权限码 */
  function hasPermission(code: string): boolean {
    if (permissions.value.includes('*')) return true
    return permissions.value.includes(code)
  }

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
    localStorage.setItem(USER_INFO_KEY, JSON.stringify(res.data))
  }

  function clear() {
    userInfo.value = null
    localStorage.removeItem(USER_INFO_KEY)
  }

  return {
    userInfo,
    isSuperAdmin,
    isHQAdmin,
    isFinance,
    permissions,
    hasPermission,
    userCenterId,
    userCenterName,
    canSelectAllCenters,
    displayName,
    fetchUserInfo,
    clear
  }
})
