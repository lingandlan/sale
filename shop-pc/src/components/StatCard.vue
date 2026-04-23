<template>
  <div class="stat-card">
    <div class="card-body">
      <div class="card-header">
        <span class="trend" :class="trendClass">{{ trend }}</span>
      </div>
      <div class="value" :style="{ color: valueColor }">
        {{ prefix }}{{ formattedValue }}
      </div>
      <div class="label">{{ label }}</div>
    </div>
    <div class="icon-block" :style="{ background: iconBg }">
      <el-icon :size="28" color="#fff"><component :is="iconComp" /></el-icon>
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import { DataAnalysis, Wallet, Ticket, User, OfficeBuilding } from '@element-plus/icons-vue'

const iconMap: Record<string, any> = {
  DataAnalysis, Wallet, Ticket, User, OfficeBuilding
}

interface Props {
  icon: string
  value: number
  label: string
  trend?: string
  valueColor?: string
  prefix?: string
  color?: string
}

const props = withDefaults(defineProps<Props>(), {
  valueColor: 'var(--color-text-primary)',
  prefix: '',
  color: 'var(--color-primary)'
})

const iconComp = computed(() => iconMap[props.icon] || DataAnalysis)

const colorMap: Record<string, string> = {
  'var(--color-info)': '#1677FF',
  'var(--color-primary)': '#C00000',
  'var(--color-success)': '#52C41A',
  'var(--color-warning)': '#FAAD14'
}

const iconBg = computed(() => {
  const hex = colorMap[props.color] || '#C00000'
  return hex
})

const trendClass = computed(() => {
  if (!props.trend) return ''
  return props.trend.startsWith('+') ? 'trend-up' : 'trend-down'
})

const formattedValue = computed(() => {
  return (props.value ?? 0).toLocaleString()
})
</script>

<style scoped>
.stat-card {
  width: 280px;
  height: 120px;
  background: var(--color-bg-card);
  border-radius: var(--radius-md);
  border: 1px solid var(--color-border);
  display: flex;
  overflow: hidden;
}

.card-body {
  flex: 1;
  padding: 16px 20px;
  display: flex;
  flex-direction: column;
  justify-content: center;
  gap: 4px;
}

.card-header {
  display: flex;
  justify-content: flex-end;
  align-items: center;
  min-height: 20px;
}

.trend {
  font-size: 12px;
}

.trend-up {
  color: var(--color-success);
}

.trend-down {
  color: var(--color-danger);
}

.value {
  font-size: 28px;
  font-weight: 700;
  line-height: 1.2;
}

.label {
  font-size: 13px;
  color: var(--color-text-muted);
}

.icon-block {
  width: 80px;
  display: flex;
  align-items: center;
  justify-content: center;
  flex-shrink: 0;
}
</style>
