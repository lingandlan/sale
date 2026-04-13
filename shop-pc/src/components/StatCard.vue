<template>
  <div class="stat-card">
    <div class="card-header">
      <span class="icon">{{ icon }}</span>
      <span class="trend" :class="trendClass">{{ trend }}</span>
    </div>
    <div class="value" :style="{ color: valueColor }">
      {{ prefix }}{{ formattedValue }}
    </div>
    <div class="label">{{ label }}</div>
  </div>
</template>

<script setup lang="ts">
import { computed } from 'vue'

interface Props {
  icon: string
  value: number
  label: string
  trend?: string
  valueColor?: string
  prefix?: string
}

const props = withDefaults(defineProps<Props>(), {
  valueColor: '#262626',
  prefix: ''
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
  background: #FFFFFF;
  border-radius: 8px;
  border: 1px solid #E5E5E5;
  padding: 20px;
  display: flex;
  flex-direction: column;
  gap: 12px;
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.icon {
  font-size: 24px;
}

.trend {
  font-size: 12px;
  font-family: 'Inter', sans-serif;
}

.trend-up {
  color: #52C41A;
}

.trend-down {
  color: #FF4D4F;
}

.value {
  font-family: 'Inter', sans-serif;
  font-size: 32px;
  font-weight: 600;
}

.label {
  font-family: 'Inter', sans-serif;
  font-size: 14px;
  color: #8C8C8C;
}
</style>
