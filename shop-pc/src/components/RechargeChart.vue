<template>
  <div class="recharge-chart" :style="{ height: `${height}px` }">
    <div class="chart-bars">
      <div
        v-for="(bar, index) in data"
        :key="index"
        class="chart-bar-wrapper"
      >
        <div class="bar-container">
          <div
            class="chart-bar"
            :style="{
              height: `${(bar.value / maxValue) * 100}%`,
              backgroundColor: bar.color
            }"
          >
            <span class="bar-value">¥{{ bar.value.toLocaleString() }}</span>
          </div>
        </div>
        <div class="chart-label">{{ bar.label }}</div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed } from 'vue'

interface ChartBar {
  label: string
  value: number
  color: string
}

interface Props {
  height?: number
  data: ChartBar[]
}

const props = withDefaults(defineProps<Props>(), {
  height: 200
})

const maxValue = computed(() => {
  return Math.max(...props.data.map(d => d.value), 1)
})
</script>

<style scoped>
.recharge-chart {
  display: flex;
  align-items: flex-end;
  width: 100%;
}

.chart-bars {
  display: flex;
  justify-content: space-around;
  align-items: flex-end;
  width: 100%;
  height: 100%;
  gap: 16px;
}

.chart-bar-wrapper {
  display: flex;
  flex-direction: column;
  align-items: center;
  flex: 1;
  height: 100%;
}

.bar-container {
  flex: 1;
  width: 100%;
  display: flex;
  align-items: flex-end;
  justify-content: center;
}

.chart-bar {
  width: 48px;
  border-radius: 6px 6px 0 0;
  transition: all 0.3s ease;
  position: relative;
  min-height: 4px;
}

.chart-bar:hover {
  filter: brightness(1.15);
  width: 56px;
}

.bar-value {
  position: absolute;
  top: -24px;
  left: 50%;
  transform: translateX(-50%);
  font-size: 11px;
  color: var(--color-text-secondary);
  white-space: nowrap;
  opacity: 0;
  transition: opacity 0.2s;
}

.chart-bar:hover .bar-value {
  opacity: 1;
}

.chart-label {
  margin-top: 8px;
  font-size: 12px;
  color: var(--color-text-muted);
  text-align: center;
}
</style>
