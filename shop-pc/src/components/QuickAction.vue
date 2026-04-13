<template>
  <div class="quick-action" @click="handleClick">
    <div class="action-icon">{{ icon }}</div>
    <div class="action-text">{{ text }}</div>
  </div>
</template>

<script setup lang="ts">
import { useRouter } from 'vue-router'

interface Props {
  icon: string
  text: string
  route: string
  background: string
  border: string
}

const props = defineProps<Props>()
const router = useRouter()

const handleClick = () => {
  router.push(props.route)
}
</script>

<style scoped>
.quick-action {
  width: 140px;
  height: 100px;
  border-radius: 8px;
  border: 2px solid v-bind(background, props.border);
  padding: 16px;
  display: flex;
  flex-direction: column;
  gap: 12px;
  justify-content: center;
  align-items: center;
  cursor: pointer;
  transition: all 0.3s;
  background-color: v-bind('props.background');
}

.quick-action:hover {
  transform: translateY(-2px);
  box-shadow: 0 4px 12px rgba(0,0,0,0.1);
}

.action-icon {
  font-size: 32px;
}

.action-text {
  font-family: 'Inter', sans-serif;
  font-size: 14px;
  font-weight: 500;
  color: #262626;
}
</style>
