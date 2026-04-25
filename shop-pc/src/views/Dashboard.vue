<template>
  <div class="dashboard">
    <!-- 欢迎区域 -->
    <div class="welcome-section">
      <div class="welcome-left">
        <h1 class="welcome-title">欢迎使用太积堂充值与门店管理系统</h1>
        <p class="welcome-subtitle">高效管理充值业务，实时掌握运营数据</p>
      </div>
      <div class="welcome-right">
        <p class="current-time">{{ currentTime }}</p>
        <p class="current-user">{{ userName }}</p>
      </div>
    </div>

    <!-- 统计卡片区域 -->
    <div class="stats-section">
      <StatCard
        icon="Wallet"
        :value="statistics.todayRecharge"
        label="今日充值金额"
        :trend="statistics.rechargeTrend"
        :value-color="'var(--color-primary)'"
        prefix="¥ "
        color="var(--color-primary)"
      />
      <StatCard
        icon="Ticket"
        :value="statistics.todayConsumption"
        label="今日核销金额"
        :trend="statistics.consumptionTrend"
        :value-color="'var(--color-success)'"
        prefix="¥ "
        color="var(--color-success)"
      />
      <StatCard
        v-if="isHeadquarters"
        icon="OfficeBuilding"
        :value="statistics.activeCenters"
        label="活跃中心数"
        :trend="statistics.centerTrend"
        color="var(--color-warning)"
      />
    </div>

    <!-- 快捷操作区域 -->
    <div class="quick-actions-section">
      <div class="section-card">
        <div class="section-header">
          <h3 class="section-title">快捷操作</h3>
        </div>
        <el-divider />
        <div class="quick-actions-grid">
          <QuickAction
            v-for="action in filteredQuickActions"
            :key="action.route"
            :icon="action.icon"
            :text="action.text"
            :background="action.background"
            :border="action.border"
            :route="action.route"
          />
        </div>
      </div>
    </div>

    <!-- 底部区域：待办事项 + 图表 -->
    <div class="bottom-section">
      <div v-if="isHeadquarters" class="todo-section">
        <div class="section-card todo-card">
          <div class="section-header">
            <h3 class="section-title">待办事项</h3>
          </div>
          <el-divider />
          <div class="todo-list">
            <div
              v-for="todo in todos"
              :key="todo.id"
              class="todo-item"
              :class="[todo.type, { clickable: todo.route }]"
              @click="todo.route && $router.push(todo.route)"
            >
              <span class="todo-icon">{{ todo.icon }}</span>
              <div class="todo-content">
                <p class="todo-title">{{ todo.title }}<span v-if="todo.route" class="todo-link">去处理 &gt;</span></p>
                <p class="todo-desc">{{ todo.count }}{{ todo.description }}</p>
              </div>
            </div>
          </div>
        </div>
      </div>

      <div class="chart-section">
        <div class="section-card chart-card">
          <div class="section-header">
            <h3 class="section-title">充值趋势（最近7天）</h3>
          </div>
          <el-divider />
          <div class="chart-container">
            <RechargeChart :data="chartData" :height="200" />
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted, onUnmounted } from 'vue'
import { ElMessage } from 'element-plus'
import { extractErrorMessage } from '@/utils/request'
import StatCard from '../components/StatCard.vue'
import QuickAction from '../components/QuickAction.vue'
import RechargeChart from '../components/RechargeChart.vue'
import { getStatistics, getTodos, getRechargeTrends } from '../api/dashboard'
import type { Statistics as StatisticsType, Todos, RechargeTrend } from '../types/dashboard'
import { useUserStore } from '@/stores/user'

const userStore = useUserStore()
const isHeadquarters = computed(() => userStore.canSelectAllCenters)
const userName = computed(() => {
  const name = userStore.displayName
  const roleMap: Record<string, string> = {
    super_admin: '超级管理员',
    hq_admin: '总部管理员',
    finance: '财务运营',
    center_admin: '中心管理员',
    operator: '操作员'
  }
  const role = roleMap[userStore.userInfo?.role || ''] || ''
  return role ? `${role}：${name}` : name || ''
})
const currentTime = ref('')
let isUnmounted = false

// 统计数据
const statistics = ref<StatisticsType>({
  todayRecharge: 0,
  rechargeTrend: '—',
  todayConsumption: 0,
  consumptionTrend: '—',
  activeCenters: 0,
  centerTrend: '—'
})

// 待办事项
const todos = ref([] as { id: string; title: string; description: string; type: string; icon: string; count: number; route?: string }[])

// 充值趋势数据
const chartData = ref([] as { label: string; value: number; color: string }[])

// 快捷操作
const quickActions = ref([
  {
    icon: '💰',
    text: 'C端充值录入',
    route: '/recharge/c-entry',
    background: 'var(--color-warning-bg)',
    border: 'var(--color-primary-gold)'
  },
  {
    icon: '🎫',
    text: '门店卡核销',
    route: '/card/verify',
    background: 'var(--color-info-bg)',
    border: 'var(--color-info)'
  },
  {
    icon: '🎁',
    text: '绑定卡号',
    route: '/card/issue',
    background: 'var(--color-success-bg)',
    border: 'var(--color-success)'
  },
  {
    icon: '📝',
    text: 'B端充值申请',
    route: '/recharge/b-apply',
    background: 'var(--color-danger-bg)',
    border: 'var(--color-danger)'
  }
])

const filteredQuickActions = computed(() => {
  if (isHeadquarters.value) return quickActions.value
  return quickActions.value.filter(a => a.route !== '/recharge/b-apply')
})

// 加载数据
const loadDashboardData = async () => {
  try {
    const statsRes = await getStatistics()
    if (!isUnmounted && statsRes?.data) {
      statistics.value = statsRes.data
    }

    const todosRes = await getTodos()
    if (!isUnmounted && todosRes?.data) {
      const d = todosRes.data
      todos.value = [
        {
          id: '1',
          title: '待审批充值申请',
          description: d.pendingApprovals?.description || '',
          type: 'warning',
          icon: '',
          count: d.pendingApprovals?.count || 0,
          route: '/recharge/b-approval'
        },
        {
          id: '2',
          title: '即将过期的门店卡',
          description: d.expiringCards?.description || '',
          type: 'error',
          icon: '',
          count: d.expiringCards?.count || 0
        }
      ].filter(t => t.count > 0)
    }

    const trendsRes = await getRechargeTrends(7)
    if (!isUnmounted && trendsRes?.data) {
      chartData.value = trendsRes.data.dates.map((date: string, index: number) => ({
        label: date,
        value: trendsRes.data.values[index],
        color: index === trendsRes.data.values.length - 1 ? 'var(--color-primary-gold)' : 'var(--color-primary)'
      }))
    }
  } catch (err: any) {
    if (!isUnmounted) {
      ElMessage.error(extractErrorMessage(err, '加载Dashboard数据失败'))
    }
  }
}

// 更新当前时间
const updateTime = () => {
  const now = new Date()
  const options: Intl.DateTimeFormatOptions = {
    year: 'numeric',
    month: 'long',
    day: 'numeric',
    weekday: 'long'
  }
  currentTime.value = now.toLocaleDateString('zh-CN', options)
}

let timeInterval: number

onMounted(() => {
  updateTime()
  timeInterval = window.setInterval(updateTime, 60000)
  loadDashboardData()
})

onUnmounted(() => {
  isUnmounted = true
  if (timeInterval) {
    clearInterval(timeInterval)
  }
})
</script>

<style scoped>
.dashboard {
  display: flex;
  flex-direction: column;
  gap: var(--spacing-lg);
}

/* 欢迎区域 */
.welcome-section {
  background-color: var(--color-primary);
  padding: var(--spacing-base) var(--spacing-xl);
  display: flex;
  justify-content: space-between;
  align-items: center;
  border-radius: var(--radius-md);
}

.welcome-left {
  display: flex;
  flex-direction: column;
  gap: var(--spacing-xs);
}

.welcome-title {
  font-family: var(--font-family);
  font-size: var(--font-size-xl);
  font-weight: 600;
  color: var(--color-text-white);
  margin: 0;
}

.welcome-subtitle {
  font-family: var(--font-family);
  font-size: var(--font-size-base);
  font-weight: 400;
  color: var(--color-primary-gold);
  margin: 0;
}

.welcome-right {
  display: flex;
  flex-direction: column;
  gap: var(--spacing-xs);
  align-items: flex-end;
}

.current-time,
.current-user {
  font-family: var(--font-family);
  font-size: var(--font-size-base);
  color: var(--color-text-white);
  margin: 0;
}

/* 统计卡片区域 */
.stats-section {
  display: flex;
  gap: var(--spacing-base);
}

/* 快捷操作区域 */
.quick-actions-section {
  width: 100%;
}

.section-card {
  background-color: var(--color-bg-card);
  border-radius: var(--radius-md);
  border: 1px solid var(--color-border);
  padding: var(--spacing-lg);
}

.section-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: var(--spacing-base);
}

.section-title {
  font-family: var(--font-family);
  font-size: var(--font-size-lg);
  font-weight: 600;
  color: var(--color-text-primary);
  margin: 0;
}

.more-link {
  font-family: var(--font-family);
  font-size: var(--font-size-xs);
  color: var(--color-info);
  cursor: pointer;
}

.quick-actions-grid {
  display: grid;
  grid-template-columns: repeat(4, 1fr);
  gap: var(--spacing-base);
}

/* 底部区域 */
.bottom-section {
  display: flex;
  gap: var(--spacing-base);
}

.todo-card {
  width: 400px;
}

.chart-card {
  flex: 1;
}

.todo-list {
  display: flex;
  flex-direction: column;
  gap: var(--spacing-md);
}

.todo-item {
  display: flex;
  gap: var(--spacing-md);
  padding: var(--spacing-md);
  border-radius: var(--radius-sm);
  align-items: center;
}

.todo-item.clickable {
  cursor: pointer;
}

.todo-item.clickable:hover {
  opacity: 0.8;
}

.todo-item.warning {
  background-color: var(--color-warning-bg);
}

.todo-item.error {
  background-color: var(--color-danger-bg);
}

.todo-icon {
  font-size: 20px;
}

.todo-content {
  display: flex;
  flex-direction: column;
  gap: var(--spacing-xs);
}

.todo-title {
  font-family: var(--font-family);
  font-size: var(--font-size-base);
  font-weight: 500;
  color: var(--color-text-primary);
  margin: 0;
}

.todo-link {
  font-size: var(--font-size-xs);
  color: var(--color-info);
  margin-left: var(--spacing-sm);
}

.todo-desc {
  font-family: var(--font-family);
  font-size: var(--font-size-xs);
  color: var(--color-text-muted);
  margin: 0;
}

.chart-container {
  padding: 0 var(--spacing-lg);
}
</style>
