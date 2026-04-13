<template>
  <div class="main-layout">
    <!-- 左侧边栏 -->
    <Sidebar v-model:collapsed="collapsed" />

    <!-- 右侧内容区域 -->
    <div class="content-area">
      <!-- 顶部Header -->
      <Header :collapsed="collapsed" @toggle="handleToggle" />

      <!-- 主内容 -->
      <div class="main-content">
        <router-view v-slot="{ Component, route }">
          <component :is="Component" :key="route.fullPath" />
        </router-view>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import Sidebar from './Sidebar.vue'
import Header from './Header.vue'

const collapsed = ref(false)

const handleToggle = () => {
  collapsed.value = !collapsed.value
}
</script>

<style scoped>
.main-layout {
  display: flex;
  width: 100%;
  height: 100vh;
  overflow: hidden;
}

.content-area {
  flex: 1;
  display: flex;
  flex-direction: column;
  overflow: hidden;
  background-color: #F5F5F5;
}

.main-content {
  flex: 1;
  overflow-y: auto;
  padding: 24px;
}
</style>
