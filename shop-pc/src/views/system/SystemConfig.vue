<template>
  <div class="system-config">
    <div class="page-header">
      <h1 class="page-title">系统设置</h1>
    </div>

    <div class="content-area">
      <!-- 左侧分类导航 -->
      <div class="sidebar">
        <div
          v-for="item in categories"
          :key="item.key"
          :class="['sidebar-item', { active: activeCategory === item.key }]"
          @click="activeCategory = item.key"
        >
          <span class="sidebar-icon">{{ item.icon }}</span>
          <span class="sidebar-text">{{ item.label }}</span>
        </div>
      </div>

      <!-- 右侧内容区 -->
      <div class="main-content">
        <!-- 商城系统设置 -->
        <div v-if="activeCategory === 'mall'" class="config-card">
          <div class="config-header">商城系统设置</div>
          <el-divider />
          <el-form label-width="140px" style="max-width: 600px">
            <el-form-item label="商城维护模式">
              <el-switch
                v-model="mallConfig.maintenanceEnabled"
                active-text="开"
                inactive-text="关"
                active-color="#C00000"
              />
            </el-form-item>
            <el-form-item label="维护原因" v-if="mallConfig.maintenanceEnabled">
              <el-input v-model="mallConfig.maintenanceReason" placeholder="请输入维护原因" />
            </el-form-item>
            <el-form-item label="商城接口地址">
              <el-input v-model="mallConfig.baseUrl" placeholder="https://api.example.com" />
            </el-form-item>
            <el-form-item label="商城AppID">
              <el-input v-model="mallConfig.appId" placeholder="请输入AppID" />
            </el-form-item>
            <el-form-item label="商城AppSecret">
              <el-input v-model="mallConfig.appSecret" type="password" show-password placeholder="请输入AppSecret" />
            </el-form-item>
            <el-form-item label="客户标识">
              <el-input v-model="mallConfig.customerId" placeholder="请输入客户标识" />
            </el-form-item>
            <el-form-item label="接口超时时间">
              <el-input-number v-model="mallConfig.timeout" :min="1" :max="60" />
              <span class="form-hint">秒</span>
            </el-form-item>
            <el-form-item label="重试次数">
              <el-input-number v-model="mallConfig.maxRetry" :min="1" :max="10" />
              <span class="form-hint">次</span>
            </el-form-item>
            <el-form-item label="重试间隔">
              <el-input-number v-model="mallConfig.retryInterval" :min="10" :max="300" />
              <span class="form-hint">秒</span>
            </el-form-item>
            <el-form-item>
              <el-button type="primary" class="save-btn" @click="handleSave('mall')">保存设置</el-button>
            </el-form-item>
          </el-form>
        </div>

        <!-- 充值系统设置 -->
        <div v-if="activeCategory === 'recharge'" class="config-card">
          <div class="config-header">充值系统设置</div>
          <el-divider />
          <el-form label-width="140px" style="max-width: 600px">
            <el-form-item label="最低充值金额">
              <el-input-number v-model="rechargeConfig.minAmount" :min="100" :step="100" />
              <span class="form-hint">元</span>
            </el-form-item>
            <el-form-item label="充值倍数要求">
              <el-input-number v-model="rechargeConfig.multiple" :min="100" :step="100" />
              <span class="form-hint">的倍数</span>
            </el-form-item>
            <el-form-item label="B端充值积分比例">
              <el-input v-model="rechargeConfig.defaultRatio" disabled />
              <span class="form-hint">&ge;10万按2%，&lt;10万按1%</span>
            </el-form-item>
            <el-form-item label="异步重试开启">
              <el-switch
                v-model="rechargeConfig.asyncRetryEnabled"
                active-text="开"
                inactive-text="关"
                active-color="#C00000"
              />
            </el-form-item>
            <el-form-item label="充值通知方式">
              <el-select v-model="rechargeConfig.notifyType">
                <el-option label="同步通知" value="sync" />
                <el-option label="异步通知（MQ）" value="async" />
              </el-select>
            </el-form-item>
            <el-form-item>
              <el-button type="primary" class="save-btn" @click="handleSave('recharge')">保存设置</el-button>
            </el-form-item>
          </el-form>
        </div>

        <!-- 门店卡设置 -->
        <div v-if="activeCategory === 'card'" class="config-card">
          <div class="config-header">门店卡设置</div>
          <el-divider />
          <el-form label-width="140px" style="max-width: 600px">
            <el-form-item label="卡号前缀">
              <el-input v-model="cardConfig.cardPrefix" maxlength="4" />
            </el-form-item>
            <el-form-item label="卡面值">
              <el-input-number v-model="cardConfig.faceValue" :min="100" :step="100" disabled />
              <span class="form-hint">元（固定）</span>
            </el-form-item>
            <el-form-item label="有效期">
              <el-input-number v-model="cardConfig.validityMonths" :min="1" disabled />
              <span class="form-hint">个月（从激活日起）</span>
            </el-form-item>
            <el-form-item label="最低核销金额">
              <el-input-number v-model="cardConfig.minConsume" :min="1" />
              <span class="form-hint">元</span>
            </el-form-item>
            <el-form-item label="库存告警阈值">
              <el-input-number v-model="cardConfig.stockAlertThreshold" :min="1" />
              <span class="form-hint">张</span>
            </el-form-item>
            <el-form-item>
              <el-button type="primary" class="save-btn" @click="handleSave('card')">保存设置</el-button>
            </el-form-item>
          </el-form>
        </div>

        <!-- 操作日志设置 -->
        <div v-if="activeCategory === 'log'" class="config-card">
          <div class="config-header">操作日志设置</div>
          <el-divider />
          <el-form label-width="140px" style="max-width: 600px">
            <el-form-item label="哈希链校验">
              <el-switch
                v-model="logConfig.hashChainEnabled"
                active-text="开"
                inactive-text="关"
                active-color="#C00000"
              />
            </el-form-item>
            <el-form-item label="数据在线保留">
              <el-input-number v-model="logConfig.onlineRetentionMonths" :min="1" :max="12" />
              <span class="form-hint">个月</span>
            </el-form-item>
            <el-form-item label="归档存储">
              <el-switch
                v-model="logConfig.archiveEnabled"
                active-text="开"
                inactive-text="关"
                active-color="#C00000"
              />
            </el-form-item>
            <el-form-item>
              <el-button type="primary" class="save-btn" @click="handleSave('log')">保存设置</el-button>
            </el-form-item>
          </el-form>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import { ElMessage } from 'element-plus'

const activeCategory = ref('mall')

const categories = [
  { key: 'mall', label: '商城系统设置', icon: '🏪' },
  { key: 'recharge', label: '充值系统设置', icon: '💰' },
  { key: 'card', label: '门店卡设置', icon: '🎫' },
  { key: 'log', label: '操作日志设置', icon: '📋' }
]

const mallConfig = ref({
  maintenanceEnabled: false,
  maintenanceReason: '',
  baseUrl: '',
  appId: 'wx5192833637f98610',
  appSecret: '',
  customerId: 'taijido',
  timeout: 10,
  maxRetry: 5,
  retryInterval: 60
})

const rechargeConfig = ref({
  minAmount: 1000,
  multiple: 1000,
  defaultRatio: '自动计算',
  asyncRetryEnabled: true,
  notifyType: 'async'
})

const cardConfig = ref({
  cardPrefix: 'TJ',
  faceValue: 1000,
  validityMonths: 12,
  minConsume: 100,
  stockAlertThreshold: 100
})

const logConfig = ref({
  hashChainEnabled: true,
  onlineRetentionMonths: 2,
  archiveEnabled: true
})

const handleSave = (category: string) => {
  ElMessage.success('设置已保存')
}
</script>

<style scoped>
.system-config {
  background-color: #F5F5F5;
  min-height: calc(100vh - 64px);
}

.page-header {
  height: 64px;
  background-color: #FFFFFF;
  border-bottom: 1px solid #E5E5E5;
  padding: 16px 24px;
}

.page-title {
  font-size: 20px;
  font-weight: 600;
  color: #262626;
  margin: 0;
}

.content-area {
  padding: 24px;
  display: flex;
  gap: 24px;
}

.sidebar {
  width: 240px;
  background-color: #FFFFFF;
  border-radius: 8px;
  border: 1px solid #E5E5E5;
  padding: 16px 0;
  flex-shrink: 0;
}

.sidebar-item {
  display: flex;
  align-items: center;
  gap: 8px;
  padding: 12px 24px;
  cursor: pointer;
  font-size: 14px;
  color: #595959;
  transition: all 0.2s;
}

.sidebar-item:hover {
  background-color: #F5F5F5;
  color: #C00000;
}

.sidebar-item.active {
  background-color: #FFF1F0;
  color: #C00000;
  font-weight: 600;
  border-right: 3px solid #C00000;
}

.sidebar-icon {
  font-size: 16px;
}

.sidebar-text {
  font-family: 'Inter', sans-serif;
}

.main-content {
  flex: 1;
}

.config-card {
  background-color: #FFFFFF;
  border-radius: 8px;
  border: 1px solid #E5E5E5;
  padding: 24px;
}

.config-header {
  font-size: 16px;
  font-weight: 600;
  color: #262626;
}

.form-hint {
  margin-left: 8px;
  font-size: 12px;
  color: #8C8C8C;
}

.save-btn {
  width: 120px;
  background-color: #C00000;
  border-color: #C00000;
  border-radius: 4px;
}

.save-btn:hover {
  background-color: #A00000;
  border-color: #A00000;
}
</style>
