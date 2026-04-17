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
        icon="👥"
        :value="statistics.memberCount"
        label="总会员数"
        :trend="statistics.memberTrend"
        :value-color="'#262626'"
      />
      <StatCard
        icon="💰"
        :value="statistics.todayRecharge"
        label="今日充值金额"
        :trend="statistics.rechargeTrend"
        :value-color="'#C00000'"
        prefix="¥ "
      />
      <StatCard
        icon="🎫"
        :value="statistics.todayConsumption"
        label="今日核销金额"
        :trend="statistics.consumptionTrend"
        :value-color="'#52C41A'"
        prefix="¥ "
      />
      <StatCard
        icon="🏢"
        :value="statistics.activeCenters"
        label="活跃中心数"
        :trend="statistics.centerTrend"
        :value-color="'#1677FF'"
      />
    </div>

    <!-- 快捷操作区域 -->
    <div class="quick-actions-section">
      <div class="section-card">
        <div class="section-header">
          <h3 class="section-title">⚡ 快捷操作</h3>
        </div>
        <el-divider />
        <div class="quick-actions-grid">
          <QuickAction
            v-for="action in quickActions"
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
      <div class="todo-section">
        <div class="section-card todo-card">
          <div class="section-header">
            <h3 class="section-title">📋 待办事项</h3>
            <span class="more-link">查看更多 &gt;</span>
          </div>
          <el-divider />
          <div class="todo-list">
            <div
              v-for="todo in todos"
              :key="todo.id"
              class="todo-item"
              :class="todo.type"
            >
              <span class="todo-icon">{{ todo.icon }}</span>
              <div class="todo-content">
                <p class="todo-title">{{ todo.title }}</p>
                <p class="todo-desc">{{ todo.count }}{{ todo.description }}</p>
              </div>
            </div>
          </div>
        </div>
      </div>

      <div class="chart-section">
        <div class="section-card chart-card">
          <div class="section-header">
            <h3 class="section-title">📈 充值趋势（最近7天）</h3>
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
  memberCount: 1234,
  memberTrend: '+12%',
  todayRecharge: 56780,
  rechargeTrend: '+8%',
  todayConsumption: 23450,
  consumptionTrend: '+15%',
  activeCenters: 5,
  centerTrend: '+3%'
})

// 待办事项
const todos = ref([
  {
    id: '1',
    title: '待审批充值申请',
    description: '3笔申请等待审批',
    type: 'warning',
    icon: '⏰',
    count: 3
  },
  {
    id: '2',
    title: '即将过期的门店卡',
    description: '12张门店卡7天内过期',
    type: 'error',
    icon: '⚠️',
    count: 12
  }
])

// 充值趋势数据
const chartData = ref([
  { label: '4/3', value: 12000, color: '#C00000' },
  { label: '4/4', value: 9000, color: '#C00000' },
  { label: '4/5', value: 15000, color: '#C00000' },
  { label: '4/6', value: 18000, color: '#C00000' },
  { label: '4/7', value: 14000, color: '#C00000' },
  { label: '4/8', value: 16000, color: '#C00000' },
  { label: '4/9', value: 20000, color: '#FFD700' }
])

// 快捷操作
const quickActions = ref([
  {
    icon: '💵',
    text: 'C端充值',
    route: '/recharge/c-entry',
    background: '#FFF7E6',
    border: '#FFD700'
  },
  {
    icon: '🎫',
    text: '门店卡核销',
    route: '/card/verify',
    background: '#E6F7FF',
    border: '#1677FF'
  },
  {
    icon: '🎁',
    text: '门店卡发放',
    route: '/card/issue',
    background: '#F6FFED',
    border: '#52C41A'
  },
  {
    icon: '📝',
    text: '充值申请',
    route: '/recharge/b-apply',
    background: '#FFF1F0',
    border: '#FF4D4F'
  }
])

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
          icon: '⏰',
          count: d.pendingApprovals?.count || 0
        },
        {
          id: '2',
          title: '即将过期的门店卡',
          description: d.expiringCards?.description || '',
          type: 'error',
          icon: '⚠️',
          count: d.expiringCards?.count || 0
        }
      ].filter(t => t.count > 0)
    }

    const trendsRes = await getRechargeTrends(7)
    if (!isUnmounted && trendsRes?.data) {
      chartData.value = trendsRes.data.dates.map((date: string, index: number) => ({
        label: date,
        value: trendsRes.data.values[index],
        color: index === trendsRes.data.values.length - 1 ? '#FFD700' : '#C00000'
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
  gap: 24px;
}

/* 欢迎区域 */
.welcome-section {
  background-color: #C00000;
  padding: 16px 32px;
  display: flex;
  justify-content: space-between;
  align-items: center;
  border-radius: 8px;
}

.welcome-left {
  display: flex;
  flex-direction: column;
  gap: 4px;
}

.welcome-title {
  font-family: 'Inter', sans-serif;
  font-size: 20px;
  font-weight: 600;
  color: #FFFFFF;
  margin: 0;
}

.welcome-subtitle {
  font-family: 'Inter', sans-serif;
  font-size: 14px;
  font-weight: 400;
  color: #FFD700;
  margin: 0;
}

.welcome-right {
  display: flex;
  flex-direction: column;
  gap: 4px;
  align-items: flex-end;
}

.current-time,
.current-user {
  font-family: 'Inter', sans-serif;
  font-size: 14px;
  color: #FFFFFF;
  margin: 0;
}

/* 统计卡片区域 */
.stats-section {
  display: flex;
  gap: 16px;
}

/* 快捷操作区域 */
.quick-actions-section {
  width: 100%;
}

.section-card {
  background-color: #FFFFFF;
  border-radius: 8px;
  border: 1px solid #E5E5E5;
  padding: 24px;
}

.section-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 16px;
}

.section-title {
  font-family: 'Inter', sans-serif;
  font-size: 16px;
  font-weight: 600;
  color: #262626;
  margin: 0;
}

.more-link {
  font-family: 'Inter', sans-serif;
  font-size: 12px;
  color: #1677FF;
  cursor: pointer;
}

.quick-actions-grid {
  display: grid;
  grid-template-columns: repeat(4, 1fr);
  gap: 16px;
}

/* 底部区域 */
.bottom-section {
  display: flex;
  gap: 16px;
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
  gap: 12px;
}

.todo-item {
  display: flex;
  gap: 12px;
  padding: 12px;
  border-radius: 4px;
  align-items: center;
}

.todo-item.warning {
  background-color: #FFF7E6;
}

.todo-item.error {
  background-color: #FFF1F0;
}

.todo-icon {
  font-size: 20px;
}

.todo-content {
  display: flex;
  flex-direction: column;
  gap: 4px;
}

.todo-title {
  font-family: 'Inter', sans-serif;
  font-size: 14px;
  font-weight: 500;
  color: #262626;
  margin: 0;
}

.todo-desc {
  font-family: 'Inter', sans-serif;
  font-size: 12px;
  color: #8C8C8C;
  margin: 0;
}

.chart-container {
  padding: 0 24px;
}
</style>
