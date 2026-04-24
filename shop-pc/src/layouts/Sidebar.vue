<template>
  <div class="sidebar" :class="{ collapsed }">
    <!-- Logo区域 -->
    <div class="sidebar-header">
      <img src="@/assets/logo.png" alt="太积堂" class="logo-img" />
      <div v-show="!collapsed" class="logo-text">太积堂管理系统</div>
    </div>

    <!-- 菜单区域 -->
    <el-menu
      :default-active="activeMenu"
      :collapse="collapsed"
      :unique-opened="false"
      class="sidebar-menu"
      background-color="#001529"
      text-color="rgba(255, 255, 255, 0.65)"
      active-text-color="#FFFFFF"
      @select="handleMenuSelect"
    >
      <template v-for="group in menuGroups" :key="group.title">
        <el-sub-menu v-if="group.items && group.items.length > 0" :index="group.key">
          <template #title>
            <el-icon><component :is="group.icon" /></el-icon>
            <span>{{ group.title }}</span>
          </template>
          <el-menu-item
            v-for="item in group.items"
            :key="item.key"
            :index="item.key"
          >
            {{ item.title }}
          </el-menu-item>
        </el-sub-menu>
        <el-menu-item v-else :index="group.key">
          <el-icon><component :is="group.icon" /></el-icon>
          <template #title>{{ group.title }}</template>
        </el-menu-item>
      </template>
    </el-menu>

    <!-- 退出登录 -->
    <div class="sidebar-footer">
      <div class="logout-btn" @click="handleLogout">
        <el-icon><SwitchButton /></el-icon>
        <span v-show="!collapsed" class="logout-text">退出登录</span>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, watch } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { ElMessageBox } from 'element-plus'
import { useUserStore } from '@/stores/user'

interface Props {
  collapsed?: boolean
}

const props = defineProps<Props>()
const emit = defineEmits<{
  'update:modelValue': [value: boolean]
}>()

const router = useRouter()
const route = useRoute()
const userStore = useUserStore()

const collapsed = computed({
  get: () => props.collapsed || false,
  set: (value) => emit('update:modelValue', value)
})

// 菜单定义（图标使用 Element Plus Icons 组件名）
const allMenuGroups = [
  {
    key: 'dashboard',
    title: '数据概览',
    icon: 'DataAnalysis',
    items: [
      { key: '/dashboard', title: '首页仪表盘', permission: 'dashboard' }
    ]
  },
  {
    key: 'recharge',
    title: '充值管理',
    icon: 'Wallet',
    items: [
      { key: '/recharge/b-apply', title: 'B端充值申请', permission: 'recharge:b' },
      { key: '/recharge/b-approval', title: 'B端充值审批', permission: 'recharge:b' },
      { key: '/recharge/c-entry', title: 'C端充值录入', permission: 'recharge:c' },
      { key: '/recharge/records', title: '充值记录', permission: 'recharge:records' }
    ]
  },
  {
    key: 'center',
    title: '充值中心',
    icon: 'OfficeBuilding',
    items: [
      { key: '/center/manage', title: '中心列表', permission: 'center:manage' },
      { key: '/operator/manage', title: '操作员管理', permission: 'operator:manage' }
    ]
  },
  {
    key: 'card',
    title: '门店卡',
    icon: 'Ticket',
    items: [
      { key: '/card/inventory', title: '总卡库管理', permission: 'card:inventory' },
      { key: '/card/issue', title: '绑定卡号', permission: 'card:issue' },
      { key: '/card/verify', title: '门店卡核销', permission: 'card:verify' },
      { key: '/card/manage', title: '门店卡管理', permission: 'card:manage' },
      { key: '/card/stats', title: '门店卡统计', permission: 'card:stats' }
    ]
  },
  {
    key: 'user',
    title: '用户管理',
    icon: 'User',
    items: [
      { key: '/user/manage', title: '用户列表', permission: 'user:manage' }
    ]
  },
  {
    key: 'system',
    title: '系统设置',
    icon: 'Setting',
    items: [
      { key: '/system/config', title: '系统配置', permission: 'system:config' }
    ]
  }
]

// 根据权限过滤菜单
const menuGroups = computed(() =>
  allMenuGroups
    .map(group => ({
      ...group,
      items: group.items.filter(item => userStore.hasPermission(item.permission))
    }))
    .filter(group => group.items.length > 0)
)

// 当前激活的菜单
const activeMenu = computed(() => {
  return route.path
})

// 菜单选择
const handleMenuSelect = (key: string) => {
  router.push(key)
}

// 退出登录
const handleLogout = () => {
  ElMessageBox.confirm('确定要退出登录吗？', '提示', {
    confirmButtonText: '确定',
    cancelButtonText: '取消',
    type: 'warning'
  }).then(() => {
    localStorage.removeItem('access_token')
    localStorage.removeItem('refresh_token')
    window.location.href = '/login'
  }).catch(() => {
    // 取消操作
  })
}

// 默认展开"数据概览"和"充值管理"
const defaultOpeneds = ref(['dashboard', 'recharge'])
</script>

<style scoped>
.sidebar {
  width: 240px;
  height: 100vh;
  background-color: var(--color-bg-dark);
  display: flex;
  flex-direction: column;
  transition: width 0.3s;
  flex-shrink: 0;
}

.sidebar.collapsed {
  width: 64px;
}

.sidebar-header {
  height: 64px;
  display: flex;
  align-items: center;
  padding: 0 var(--spacing-lg);
  border-bottom: 1px solid var(--color-border-light);
}

.logo-img {
  width: 32px;
  height: 32px;
  border-radius: var(--radius-sm);
  object-fit: contain;
  flex-shrink: 0;
}

.logo-text {
  margin-left: var(--spacing-md);
  font-family: var(--font-family);
  font-size: var(--font-size-lg);
  font-weight: 600;
  color: var(--color-text-white);
  white-space: nowrap;
}

.sidebar-menu {
  flex: 1;
  border: none;
  padding: var(--spacing-base) 0;
}

:deep(.el-menu) {
  border: none;
}

:deep(.el-sub-menu__title) {
  height: 48px;
  line-height: 48px;
  padding: 0 var(--spacing-lg) !important;
  color: var(--color-text-white-secondary) !important;
}

:deep(.el-sub-menu__title:hover) {
  background-color: rgba(255, 255, 255, 0.08) !important;
}

:deep(.el-menu-item) {
  height: 40px;
  line-height: 40px;
  padding: 0 var(--spacing-lg) 0 48px !important;
  color: var(--color-text-white-secondary) !important;
}

:deep(.el-menu-item:hover) {
  background-color: rgba(255, 255, 255, 0.08) !important;
}

:deep(.el-menu-item.is-active) {
  background-color: var(--color-primary) !important;
  color: var(--color-text-white) !important;
}

.sidebar-footer {
  padding: var(--spacing-base) var(--spacing-lg);
  border-top: 1px solid var(--color-border-light);
}

.logout-btn {
  display: flex;
  align-items: center;
  cursor: pointer;
  transition: all 0.3s;
  color: var(--color-text-white-secondary);
}

.logout-btn:hover {
  opacity: 0.8;
}

.logout-text {
  margin-left: var(--spacing-sm);
  font-family: var(--font-family);
  font-size: var(--font-size-base);
  color: var(--color-text-white-secondary);
}
</style>
