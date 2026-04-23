import { createRouter, createWebHistory } from 'vue-router'
import type { RouteRecordRaw } from 'vue-router'
import { useUserStore } from '../stores/user'

const routes: RouteRecordRaw[] = [
  {
    path: '/login',
    name: 'Login',
    component: () => import('../views/Login.vue'),
    meta: { requiresAuth: false }
  },
  {
    path: '/',
    component: () => import('../layouts/MainLayout.vue'),
    meta: { requiresAuth: true },
    redirect: '/dashboard',
    children: [
      {
        path: 'dashboard',
        name: 'Dashboard',
        component: () => import('../views/Dashboard.vue'),
        meta: { title: '首页仪表盘', icon: '📊', permission: 'dashboard' }
      },
      {
        path: 'recharge/b-apply',
        name: 'BRechargeApply',
        component: () => import('../views/recharge/BRechargeApply.vue'),
        meta: { title: 'B端充值申请', icon: '💰', permission: 'recharge:b' }
      },
      {
        path: 'recharge/b-approval',
        name: 'BRechargeList',
        component: () => import('../views/recharge/BRechargeList.vue'),
        meta: { title: 'B端充值审批', icon: '💰', permission: 'recharge:b' }
      },
      {
        path: 'recharge/b-approval/:id',
        name: 'BRechargeDetail',
        component: () => import('../views/recharge/BRechargeDetail.vue'),
        meta: { title: 'B端充值审批详情', icon: '💰', permission: 'recharge:b' }
      },
      {
        path: 'recharge/c-entry',
        name: 'CRechargeEntry',
        component: () => import('../views/recharge/CRechargeEntry.vue'),
        meta: { title: 'C端充值录入', icon: '💰', permission: 'recharge:c' }
      },
      {
        path: 'card/inventory',
        name: 'CardInventory',
        component: () => import('../views/card/CardInventory.vue'),
        meta: { title: '总卡库管理', icon: '🎫', permission: 'card:inventory' }
      },
      {
        path: 'card/issue',
        name: 'CardIssue',
        component: () => import('../views/card/CardIssue.vue'),
        meta: { title: '绑定卡号', icon: '🎫', permission: 'card:issue' }
      },
      {
        path: 'card/verify',
        name: 'CardVerify',
        component: () => import('../views/card/CardVerify.vue'),
        meta: { title: '门店卡核销', icon: '🎫', permission: 'card:verify' }
      },
      {
        path: 'card/manage',
        name: 'CardManage',
        component: () => import('../views/card/CardManage.vue'),
        meta: { title: '门店卡管理', icon: '🎫', permission: 'card:manage' }
      },
      {
        path: 'card/detail/:cardNo',
        name: 'CardDetail',
        component: () => import('../views/card/CardDetail.vue'),
        meta: { title: '门店卡详情', icon: '🎫', permission: 'card:manage' }
      },
      {
        path: 'card/stats',
        name: 'CardStats',
        component: () => import('../views/card/CardStats.vue'),
        meta: { title: '门店卡统计', icon: '🎫', permission: 'card:stats' }
      },
      {
        path: 'recharge/records',
        name: 'RechargeRecordList',
        component: () => import('../views/recharge-record/RecordList.vue'),
        meta: { title: '充值记录', icon: '💰', permission: 'recharge:records' }
      },
      {
        path: 'recharge/records/:id',
        name: 'RechargeRecordDetail',
        component: () => import('../views/recharge-record/RecordDetail.vue'),
        meta: { title: '充值记录详情', icon: '💰', permission: 'recharge:records' }
      },
      {
        path: 'user/manage',
        name: 'UserManage',
        component: () => import('../views/user/UserManage.vue'),
        meta: { title: '用户管理', icon: '👥', permission: 'user:manage' }
      },
      {
        path: 'center/manage',
        name: 'CenterManage',
        component: () => import('../views/center/CenterManage.vue'),
        meta: { title: '充值中心管理', icon: '🏢', permission: 'center:manage' }
      },
      {
        path: 'operator/manage',
        name: 'OperatorManage',
        component: () => import('../views/operator/OperatorManage.vue'),
        meta: { title: '充值操作员管理', icon: '👔', permission: 'operator:manage' }
      },
      {
        path: 'system/config',
        name: 'SystemConfig',
        component: () => import('../views/system/SystemConfig.vue'),
        meta: { title: '系统设置', icon: '⚙️', permission: 'system:config' }
      }
    ]
  }
]

const router = createRouter({
  history: createWebHistory(),
  routes
})

// 路由守卫
router.beforeEach(async (to, _from, next) => {
  const token = localStorage.getItem('access_token')

  if (to.meta?.requiresAuth !== false && !token) {
    next('/login')
  } else if (to.path === '/login' && token) {
    next('/dashboard')
  } else {
    // 有 token 但还没加载用户信息时，先加载
    if (token) {
      const userStore = useUserStore()
      if (!userStore.userInfo) {
        try {
          await userStore.fetchUserInfo()
        } catch {
          localStorage.removeItem('access_token')
          localStorage.removeItem('refresh_token')
          next('/login')
          return
        }
      }
      // 权限检查：路由 meta.permission 不在当前角色权限中则跳转 dashboard
      const requiredPermission = to.meta?.permission as string | undefined
      if (requiredPermission && !userStore.hasPermission(requiredPermission)) {
        next('/dashboard')
        return
      }
    }
    next()
  }
})

export default router
