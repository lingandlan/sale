import { createRouter, createWebHistory } from 'vue-router'
import type { RouteRecordRaw } from 'vue-router'

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
        meta: {
          title: '首页仪表盘',
          icon: '📊'
        }
      },
      {
        path: 'recharge/b-apply',
        name: 'BRechargeApply',
        component: () => import('../views/recharge/BRechargeApply.vue'),
        meta: {
          title: 'B端充值申请',
          icon: '💰'
        }
      },
      {
        path: 'recharge/b-approval',
        name: 'BRechargeList',
        component: () => import('../views/recharge/BRechargeList.vue'),
        meta: {
          title: 'B端充值审批',
          icon: '💰'
        }
      },
      {
        path: 'recharge/b-approval/:id',
        name: 'BRechargeDetail',
        component: () => import('../views/recharge/BRechargeDetail.vue'),
        meta: {
          title: 'B端充值审批详情',
          icon: '💰'
        }
      },
      {
        path: 'recharge/c-entry',
        name: 'CRechargeEntry',
        component: () => import('../views/recharge/CRechargeEntry.vue'),
        meta: {
          title: 'C端充值录入',
          icon: '💰'
        }
      },
      {
        path: 'card/verify',
        name: 'CardVerify',
        component: () => import('../views/card/CardVerify.vue'),
        meta: {
          title: '门店卡核销',
          icon: '🎫'
        }
      },
      {
        path: 'card/manage',
        name: 'CardManage',
        component: () => import('../views/card/CardManage.vue'),
        meta: {
          title: '门店卡管理',
          icon: '🎫'
        }
      },
      {
        path: 'card/issue',
        name: 'CardIssue',
        component: () => import('../views/card/CardIssue.vue'),
        meta: {
          title: '门店卡发放',
          icon: '🎫'
        }
      },
      {
        path: 'card/detail/:cardNo',
        name: 'CardDetail',
        component: () => import('../views/card/CardDetail.vue'),
        meta: {
          title: '门店卡详情',
          icon: '🎫'
        }
      },
      {
        path: 'card/stats',
        name: 'CardStats',
        component: () => import('../views/card/CardStats.vue'),
        meta: {
          title: '门店卡统计',
          icon: '🎫'
        }
      },
      {
        path: 'recharge/records',
        name: 'RechargeRecordList',
        component: () => import('../views/recharge-record/RecordList.vue'),
        meta: {
          title: '充值记录',
          icon: '💰'
        }
      },
      {
        path: 'user/manage',
        name: 'UserManage',
        component: () => import('../views/user/UserManage.vue'),
        meta: {
          title: '用户管理',
          icon: '👥'
        }
      },
      {
        path: 'center/manage',
        name: 'CenterManage',
        component: () => import('../views/center/CenterManage.vue'),
        meta: {
          title: '充值中心管理',
          icon: '🏢'
        }
      },
      {
        path: 'operator/manage',
        name: 'OperatorManage',
        component: () => import('../views/operator/OperatorManage.vue'),
        meta: {
          title: '充值操作员管理',
          icon: '👔'
        }
      },
      {
        path: 'system/config',
        name: 'SystemConfig',
        component: () => import('../views/system/SystemConfig.vue'),
        meta: {
          title: '系统设置',
          icon: '⚙️'
        }
      }
    ]
  }
]

const router = createRouter({
  history: createWebHistory(),
  routes
})

// 路由守卫
router.beforeEach((to, from, next) => {
  const token = localStorage.getItem('access_token')

  if (to.meta?.requiresAuth !== false && !token) {
    // 需要认证但没有token，跳转到登录页
    next('/login')
  } else if (to.path === '/login' && token) {
    // 已登录用户访问登录页，跳转到首页
    next('/dashboard')
  } else {
    next()
  }
})

export default router
