<template>
  <div class="header">
    <div class="header-left">
      <el-button
        :icon="collapsed ? 'Expand' : 'Fold'"
        link
        @click="handleToggle"
      />
    </div>

    <div class="header-right">
      <div class="user-info">
        <span class="user-name">{{ userName }}</span>
        <el-dropdown trigger="click" @command="handleCommand">
          <span class="user-avatar">
            <el-icon><User /></el-icon>
          </span>
          <template #dropdown>
            <el-dropdown-menu>
              <el-dropdown-item command="profile">个人信息</el-dropdown-item>
              <el-dropdown-item command="logout">退出登录</el-dropdown-item>
            </el-dropdown-menu>
          </template>
        </el-dropdown>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import { ElMessageBox } from 'element-plus'
import { User } from '@element-plus/icons-vue'
import { useUserStore } from '@/stores/user'

interface Props {
  collapsed?: boolean
}

const props = defineProps<Props>()
const emit = defineEmits<{
  'toggle': []
}>()

const userStore = useUserStore()
const userName = computed(() => {
  return userStore.displayName || '未登录'
})

const handleToggle = () => {
  emit('toggle')
}

const handleCommand = (command: string) => {
  if (command === 'logout') {
    ElMessageBox.confirm('确定要退出登录吗？', '提示', {
      confirmButtonText: '确定',
      cancelButtonText: '取消',
      type: 'warning'
    }).then(() => {
      localStorage.removeItem('access_token')
      localStorage.removeItem('refresh_token')
      userStore.clear()
      window.location.href = '/login'
    }).catch(() => {
      // 取消操作
    })
  } else if (command === 'profile') {
    // TODO: 跳转到个人信息页面
    console.log('个人信息')
  }
}
</script>

<style scoped>
.header {
  height: 64px;
  background-color: var(--color-bg-card);
  border-bottom: 1px solid var(--color-border);
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 0 var(--spacing-lg);
  flex-shrink: 0;
}

.header-left {
  display: flex;
  align-items: center;
}

.header-right {
  display: flex;
  align-items: center;
  gap: var(--spacing-lg);
}

.user-info {
  display: flex;
  align-items: center;
  gap: var(--spacing-md);
  cursor: pointer;
}

.user-name {
  font-family: var(--font-family);
  font-size: var(--font-size-base);
  color: var(--color-text-primary);
}

.user-avatar {
  width: 32px;
  height: 32px;
  border-radius: 50%;
  background-color: var(--color-bg);
  display: flex;
  align-items: center;
  justify-content: center;
  color: var(--color-text-muted);
  font-size: var(--font-size-lg);
  transition: all 0.3s;
}

.user-avatar:hover {
  background-color: var(--color-border);
}
</style>
