<template>
  <div class="recharge-chart" :style="{ height: `${height}px` }">
    <div class="chart-bars">
      <div
        v-for="(bar, index) in data"
        :key="index"
        class="chart-bar-wrapper"
      >
        <div
          class="chart-bar"
          :style="{
            height: `${(bar.value / maxValue) * 100}%`,
            backgroundColor: bar.color
          }"
        ></div>
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
  return Math.max(...props.data.map(d => d.value))
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

.chart-bar {
  width: 80px;
  border-radius: 4px;
  transition: height 0.3s ease;
}

.chart-label {
  margin-top: 8px;
  font-family: 'Inter', sans-serif;
  font-size: 12px;
  color: #8C8C8C;
  text-align: center;
}
</style>
