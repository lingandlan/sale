<template>
  <div class="card-stats">
    <div class="page-header">
      <h1 class="page-title">门店卡统计</h1>
    </div>

    <div class="content-area">
      <!-- 核心指标 -->
      <div class="stats-row">
        <div class="stat-card">
          <span class="stat-value">{{ stats.totalCards }}</span>
          <span class="stat-label">总库存</span>
        </div>
        <div class="stat-card">
          <span class="stat-value text-green">{{ stats.activeCards }}</span>
          <span class="stat-label">已发放</span>
        </div>
        <div class="stat-card">
          <span class="stat-value text-red">¥{{ stats.totalBalance.toLocaleString() }}</span>
          <span class="stat-label">总余额（元）</span>
        </div>
        <div class="stat-card">
          <span class="stat-value text-blue">¥{{ stats.todayConsume.toLocaleString() }}</span>
          <span class="stat-label">今日消费（元）</span>
        </div>
      </div>

      <!-- 过期提醒 -->
      <div class="warning-card" v-if="stats.expireIn7Days > 0">
        <span class="warning-icon">⚠️</span>
        <span class="warning-text">有 {{ stats.expireIn7Days }} 张卡片将在7天内过期</span>
      </div>

      <!-- 图表区域：左饼图 + 右柱状图 -->
      <div class="charts-row">
        <!-- 卡状态分布 -->
        <div class="chart-card">
          <div class="chart-header">卡状态分布</div>
          <el-divider />
          <div class="chart-container">
            <v-chart :option="pieOption" autoresize style="height: 320px" />
          </div>
        </div>

        <!-- 月度趋势 -->
        <div class="chart-card" style="flex: 1.5">
          <div class="chart-header">月度发放/核销趋势</div>
          <el-divider />
          <div class="chart-container">
            <v-chart :option="barOption" autoresize style="height: 320px" />
          </div>
        </div>
      </div>

      <!-- 充值中心卡统计表格 -->
      <div class="chart-card">
        <div class="chart-header">充值中心卡统计</div>
        <el-divider />
        <el-table :data="centerStats" stripe style="width: 100%">
          <el-table-column prop="centerName" label="充值中心" />
          <el-table-column prop="totalCards" label="总卡数" width="120" align="center" />
          <el-table-column prop="issuedCards" label="已发放" width="120" align="center" />
          <el-table-column prop="activeCards" label="已激活" width="120" align="center" />
          <el-table-column prop="frozenCards" label="已冻结" width="120" align="center" />
          <el-table-column prop="expiredCards" label="已过期" width="120" align="center" />
          <el-table-column prop="totalBalance" label="总余额" width="160" align="right">
            <template #default="{ row }">¥{{ row.totalBalance.toLocaleString() }}</template>
          </el-table-column>
        </el-table>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { ElMessage } from 'element-plus'
import { extractErrorMessage } from '@/utils/request'
import VChart from 'vue-echarts'
import { use } from 'echarts/core'
import { PieChart, BarChart } from 'echarts/charts'
import {
  TitleComponent,
  TooltipComponent,
  LegendComponent,
  GridComponent
} from 'echarts/components'
import { CanvasRenderer } from 'echarts/renderers'
import { getCardStats, getMonthlyTrend, getCenterCardStats } from '@/api/card'

use([
  PieChart, BarChart,
  TitleComponent, TooltipComponent, LegendComponent, GridComponent,
  CanvasRenderer
])

const stats = ref<Record<string, number>>({
  totalCards: 0,
  inStockCards: 0,
  issuedCards: 0,
  activeCards: 0,
  frozenCards: 0,
  expiredCards: 0,
  voidedCards: 0,
  totalBalance: 0,
  todayConsume: 0,
  expireIn7Days: 0
})

const centerStats = ref<any[]>([])
const monthlyData = ref<{ month: string; issue: number; consume: number }[]>([])

// 饼图 - 从 stats 动态生成
const pieOption = computed(() => ({
  tooltip: { trigger: 'item', formatter: '{b}: {c}张 ({d}%)' },
  legend: { bottom: 0, left: 'center' },
  series: [{
    type: 'pie',
    radius: ['40%', '70%'],
    center: ['50%', '45%'],
    avoidLabelOverlap: true,
    itemStyle: { borderRadius: 4 },
    label: { show: true, formatter: '{b}\n{c}张' },
    data: [
      { value: stats.value.inStockCards || 0, name: '已入库', itemStyle: { color: '#8C8C8C' } },
      { value: stats.value.issuedCards || 0, name: '已发放', itemStyle: { color: '#1677FF' } },
      { value: stats.value.activeCards || 0, name: '已激活', itemStyle: { color: '#52C41A' } },
      { value: stats.value.frozenCards || 0, name: '已冻结', itemStyle: { color: '#FF4D4F' } },
      { value: stats.value.expiredCards || 0, name: '已过期', itemStyle: { color: '#BFBFBF' } }
    ]
  }]
}))

// 柱状图 - 从 monthlyData 动态生成
const barOption = computed(() => ({
  tooltip: { trigger: 'axis' },
  legend: { data: ['发放', '核销'], bottom: 0 },
  grid: { left: 50, right: 20, top: 20, bottom: 40 },
  xAxis: {
    type: 'category',
    data: monthlyData.value.map(d => d.month)
  },
  yAxis: { type: 'value', name: '张数' },
  series: [
    {
      name: '发放',
      type: 'bar',
      data: monthlyData.value.map(d => d.issue),
      itemStyle: { color: '#C00000', borderRadius: [4, 4, 0, 0] }
    },
    {
      name: '核销',
      type: 'bar',
      data: monthlyData.value.map(d => d.consume),
      itemStyle: { color: '#FFD700', borderRadius: [4, 4, 0, 0] }
    }
  ]
}))

const loadData = async () => {
  try {
    const [statsRes, trendRes, centerRes] = await Promise.all([
      getCardStats(),
      getMonthlyTrend(),
      getCenterCardStats()
    ])

    if (statsRes?.data) {
      const d = statsRes.data as any
      stats.value = {
        totalCards: d.totalCards || 0,
        inStockCards: d.inStockCards || 0,
        issuedCards: d.issuedCards || 0,
        activeCards: d.activeCards || 0,
        frozenCards: d.frozenCards || 0,
        expiredCards: d.expiredCards || 0,
        voidedCards: d.voidedCards || 0,
        totalBalance: d.totalBalance || 0,
        todayConsume: d.todayConsume || 0,
        expireIn7Days: d.expireIn7Days || 0
      }
    }

    if (trendRes?.data) {
      monthlyData.value = (trendRes.data as any[]) || []
    }

    if (centerRes?.data) {
      centerStats.value = (centerRes.data as any[]) || []
    }
  } catch (err: any) {
    ElMessage.error(extractErrorMessage(err, '加载门店卡统计失败'))
  }
}

onMounted(() => { loadData() })
</script>

<style scoped>
.card-stats {
  background-color: var(--color-bg);
  min-height: calc(100vh - 64px);
}

.page-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  height: 64px;
  background-color: var(--color-bg-card);
  border-bottom: 1px solid var(--color-border);
  padding: 16px 24px;
}

.page-title {
  font-family: var(--font-family);
  font-size: 20px;
  font-weight: 600;
  color: var(--color-text-primary);
  margin: 0;
}

.content-area {
  padding: 24px;
  display: flex;
  flex-direction: column;
  gap: 24px;
}

.stats-row {
  display: grid;
  grid-template-columns: repeat(4, 1fr);
  gap: 16px;
}

.stat-card {
  background-color: var(--color-bg-card);
  border-radius: var(--radius-md);
  border: 1px solid var(--color-border);
  padding: 24px;
  display: flex;
  flex-direction: column;
  gap: 8px;
  align-items: center;
}

.stat-value {
  font-family: var(--font-family);
  font-size: 32px;
  font-weight: 600;
  color: var(--color-text-primary);
}

.text-green {
  color: var(--color-success);
}

.text-red {
  color: var(--color-danger);
}

.text-blue {
  color: var(--color-info);
}

.stat-label {
  font-family: var(--font-family);
  font-size: 14px;
  color: var(--color-text-secondary);
}

.warning-card {
  background-color: var(--color-warning-bg);
  border-radius: var(--radius-md);
  border: 1px solid #FFD591;
  padding: 16px 24px;
  display: flex;
  gap: 12px;
  align-items: center;
}

.warning-icon {
  font-size: 20px;
}

.warning-text {
  font-family: var(--font-family);
  font-size: 14px;
  color: #8C6000;
}

.charts-row {
  display: flex;
  gap: 24px;
}

.chart-card {
  background-color: var(--color-bg-card);
  border-radius: var(--radius-md);
  border: 1px solid var(--color-border);
  padding: 24px;
  flex: 1;
}

.chart-header {
  font-family: var(--font-family);
  font-size: 16px;
  font-weight: 600;
  color: var(--color-text-primary);
}

.chart-container {
  width: 100%;
}
</style>
