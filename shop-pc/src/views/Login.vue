<template>
  <div class="login-container">
    <!-- 左侧品牌区域 -->
    <div class="brand-section">
      <img src="@/assets/logo.png" alt="太积堂" class="logo-img" />
      <div class="brand-name">太积堂</div>
      <div class="system-title">充值与门店管理系统</div>
    </div>

    <!-- 右侧表单区域 -->
    <div class="form-section">
      <div class="login-card">
        <!-- 表单标题 -->
        <div class="form-header">
          <h2>登录</h2>
        </div>

        <!-- 表单内容 -->
        <div class="form-content">
          <!-- 手机号输入 -->
          <div class="form-item">
            <label class="form-label">手机号</label>
            <el-input
              v-model="phone"
              placeholder="请输入手机号"
              size="large"
              :style="{ height: '48px' }"
            />
          </div>

          <!-- 密码输入 -->
          <div class="form-item">
            <label class="form-label">密码</label>
            <el-input
              v-model="password"
              type="password"
              placeholder="请输入密码"
              size="large"
              show-password
              :style="{ height: '48px' }"
            />
          </div>

          <!-- 操作行 -->
          <div class="form-row">
            <div class="remember-row">
              <el-checkbox v-model="remember">记住密码</el-checkbox>
            </div>
            <div class="forgot-link" @click="handleForgot">忘记密码？</div>
          </div>

          <!-- 登录按钮 -->
          <el-button
            type="primary"
            size="large"
            :loading="loading"
            :style="{
              width: '100%',
              height: '48px'
            }"
            @click="handleLogin"
          >
            登录
          </el-button>

          <!-- 用户协议 -->
          <div class="agreement-text">
            登录即表示同意《用户协议》和《隐私政策》
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import { useRouter } from 'vue-router'
import { ElMessage } from 'element-plus'
import { login } from '@/api/auth'
import { useUserStore } from '@/stores/user'

const router = useRouter()
const userStore = useUserStore()

const phone = ref('')
const password = ref('')
const remember = ref(false)
const loading = ref(false)

// 登录处理
const handleLogin = async () => {
  if (!phone.value) {
    ElMessage.warning('请输入手机号')
    return
  }

  if (!password.value) {
    ElMessage.warning('请输入密码')
    return
  }

  try {
    loading.value = true
    const res = await login({
      phone: phone.value,
      password: password.value
    })
    localStorage.setItem('access_token', res.data.access_token)
    localStorage.setItem('refresh_token', res.data.refresh_token)
    await userStore.fetchUserInfo()
    ElMessage.success('登录成功')
    router.push('/dashboard')
  } catch (error) {
    console.error('登录失败', error)
  } finally {
    loading.value = false
  }
}

const handleForgot = () => {
  ElMessage.info('忘记密码功能开发中')
}
</script>

<style scoped>
.login-container {
  display: flex;
  width: 100%;
  height: 100vh;
  background-color: var(--color-bg);
}

/* 左侧品牌区域 */
.brand-section {
  width: 560px;
  background-color: var(--color-primary);
  background-image: linear-gradient(135deg, var(--color-primary) 0%, var(--color-primary-dark) 100%);
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  padding: 80px;
  position: relative;
  overflow: hidden;
}

.brand-section::before {
  content: '';
  position: absolute;
  width: 400px;
  height: 400px;
  border: 2px solid rgba(255, 215, 0, 0.08);
  border-radius: 50%;
  top: -120px;
  right: -80px;
}

.brand-section::after {
  content: '';
  position: absolute;
  width: 280px;
  height: 280px;
  border: 1px solid rgba(255, 255, 255, 0.06);
  border-radius: 50%;
  bottom: -60px;
  left: -40px;
}

.logo-img {
  width: 120px;
  height: 120px;
  border-radius: 50%;
  object-fit: contain;
}

.brand-name {
  margin-top: 24px;
  font-family: var(--font-family);
  font-size: var(--font-size-4xl);
  font-weight: 600;
  color: var(--color-primary-gold);
}

.system-title {
  margin-top: 24px;
  font-family: var(--font-family);
  font-size: 18px;
  font-weight: 400;
  color: var(--color-text-white);
}

/* 右侧表单区域 */
.form-section {
  flex: 1;
  display: flex;
  align-items: center;
  justify-content: center;
  padding: 80px;
}

.login-card {
  width: 480px;
  background-color: var(--color-bg-card);
  border-radius: var(--radius-lg);
  border: 1px solid rgba(0, 0, 0, 0.06);
  box-shadow: var(--shadow-md);
  padding: 40px;
}

/* 表单标题 */
.form-header {
  margin-bottom: 24px;
}

.form-header h2 {
  font-family: var(--font-family);
  font-size: var(--font-size-2xl);
  font-weight: 600;
  color: var(--color-text-primary);
  margin: 0;
  text-align: left;
}

/* 表单内容 */
.form-content {
  display: flex;
  flex-direction: column;
  gap: 24px;
  text-align: left;
}

.form-item {
  display: flex;
  flex-direction: column;
  gap: var(--spacing-sm);
  text-align: left;
}

.form-label {
  font-family: var(--font-family);
  font-size: var(--font-size-base);
  font-weight: 400;
  color: var(--color-text-primary);
  text-align: left;
}

.form-row {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.remember-row {
  display: flex;
  align-items: center;
}

.forgot-link {
  font-family: var(--font-family);
  font-size: var(--font-size-base);
  font-weight: 400;
  color: var(--color-primary);
  cursor: pointer;
}

.agreement-text {
  font-family: var(--font-family);
  font-size: var(--font-size-xs);
  font-weight: 400;
  color: var(--color-text-muted);
  text-align: center;
  margin-top: var(--spacing-sm);
}

/* Element Plus样式覆盖 */
:deep(.el-input__wrapper) {
  background-color: var(--color-bg-card);
}

:deep(.el-input__inner) {
  font-family: var(--font-family);
  font-size: var(--font-size-base);
}

:deep(.el-checkbox__label) {
  font-family: var(--font-family);
  font-size: var(--font-size-base);
  font-weight: 400;
  color: var(--color-text-secondary);
}
</style>
