<template>
  <div class="header">
    <div class="header-left">
      <el-button
        :icon="collapsed ? 'Expand' : 'Fold'"
        type="text"
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
import { ref, computed } from 'vue'
import { ElMessageBox } from 'element-plus'
import { User, Expand, Fold } from '@element-plus/icons-vue'

interface Props {
  collapsed?: boolean
}

const props = defineProps<Props>()
const emit = defineEmits<{
  'toggle': []
}>()

// 从localStorage获取用户信息（暂时用Mock数据）
const userName = computed(() => {
  return '管理员：张三'
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
  background-color: #FFFFFF;
  border-bottom: 1px solid #E5E5E5;
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 0 24px;
  flex-shrink: 0;
}

.header-left {
  display: flex;
  align-items: center;
}

.header-right {
  display: flex;
  align-items: center;
  gap: 24px;
}

.user-info {
  display: flex;
  align-items: center;
  gap: 12px;
  cursor: pointer;
}

.user-name {
  font-family: 'Inter', sans-serif;
  font-size: 14px;
  color: #262626;
}

.user-avatar {
  width: 32px;
  height: 32px;
  border-radius: 50%;
  background-color: #F5F5F5;
  display: flex;
  align-items: center;
  justify-content: center;
  color: #8C8C8C;
  font-size: 16px;
  transition: all 0.3s;
}

.user-avatar:hover {
  background-color: #E5E5E5;
}
</style>
