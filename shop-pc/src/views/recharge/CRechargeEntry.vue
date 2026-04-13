<template>
  <div class="crecharge-entry">
    <div class="page-header">
      <h1 class="page-title">C端充值录入</h1>
    </div>

    <div class="content-area">
      <!-- 会员查询卡片 -->
      <div class="info-card">
        <div class="card-header">
          <span class="header-icon">🔍</span>
          <h3 class="header-title">会员信息查询</h3>
        </div>
        <el-divider />
        <div class="search-row">
          <div class="search-field">
            <span class="search-label">手机号/卡号</span>
            <el-input
              v-model="searchQuery"
              placeholder="请输入手机号或会员卡号"
              style="width: 400px"
              @keyup.enter="handleSearch"
            />
          </div>
          <el-button type="primary" class="search-btn" @click="handleSearch">
            查询
          </el-button>
        </div>

        <div v-if="memberInfo" class="member-info-box">
          <div class="info-row">
            <span class="info-label">商城ID</span>
            <span class="info-value">{{ memberInfo.id }}</span>
          </div>
          <div class="info-row">
            <span class="info-label">会员姓名</span>
            <span class="info-value">{{ memberInfo.name }}</span>
          </div>
          <div class="info-row">
            <span class="info-label">手机号码</span>
            <span class="info-value">{{ memberInfo.phone }}</span>
          </div>
          <div class="info-row">
            <span class="info-label">当前积分</span>
            <span class="info-value">{{ memberInfo.balance.toLocaleString() }} 积分</span>
          </div>
          <div class="info-row">
            <span class="info-label">会员等级</span>
            <el-tag size="small" type="primary">{{ memberInfo.level }}</el-tag>
          </div>
        </div>
      </div>

      <!-- 充值信息卡片 -->
      <div class="info-card">
        <div class="card-header">
          <span class="header-icon">💰</span>
          <h3 class="header-title">充值信息</h3>
        </div>
        <el-divider />
        <div class="amount-section">
          <div class="amount-label">充值金额（元）</div>
          <div class="amount-input-wrapper">
            <span class="currency-symbol">¥</span>
            <el-input-number
              v-model="rechargeAmount"
              :min="0"
              :precision="2"
              :step="100"
              class="amount-input"
              :disabled="!memberInfo"
              @change="calculatePoints"
            />
          </div>
          <div class="quick-amounts">
            <el-button
              v-for="amount in quickAmounts"
              :key="amount"
              class="quick-amount-btn"
              :disabled="!memberInfo"
              @click="setAmount(amount)"
            >
              {{ amount }}元
            </el-button>
          </div>
        </div>

        <div class="ratio-notice">
          <span class="ratio-icon">ℹ️</span>
          <span class="ratio-text">充值比例：1元 = 1积分</span>
        </div>

        <div class="points-calc-box">
          <div class="calc-row">
            <span class="calc-label">预计获得积分</span>
            <span class="calc-value">{{ calculatedPoints.toLocaleString() }} 积分</span>
          </div>
          <el-divider style="margin: 12px 0" />
          <div class="calc-row">
            <span class="calc-label">充值后会员余额</span>
            <span class="calc-value">{{ afterRechargeBalance.toLocaleString() }} 积分</span>
          </div>
        </div>

        <div class="store-balance">
          <span class="store-label">门店当前积分余额</span>
          <span class="store-value">{{ storeBalance.toLocaleString() }} 积分</span>
        </div>
      </div>

      <!-- 备注卡片 -->
      <div class="info-card">
        <div class="card-header">
          <span class="header-icon">📝</span>
          <h3 class="header-title">备注（可选）</h3>
        </div>
        <el-divider />
        <el-input
          v-model="remark"
          type="textarea"
          :rows="3"
          placeholder="请输入备注信息（最多200字）"
          maxlength="200"
          show-word-limit
        />
      </div>
    </div>

    <!-- 操作栏 -->
    <div class="action-bar">
      <el-button class="cancel-btn" @click="handleCancel">取消</el-button>
      <el-button type="primary" class="confirm-btn" :disabled="!canSubmit" @click="handleSubmit">
        <span class="confirm-icon">✓</span>
        确认充值
      </el-button>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed } from 'vue'
import { useRouter } from 'vue-router'
import { ElMessage, ElMessageBox } from 'element-plus'
import { submitCRechargeEntry } from '@/api/recharge'

interface MemberInfo {
  id: string
  name: string
  phone: string
  balance: number
  level: string
}

const router = useRouter()

const searchQuery = ref('')
const memberInfo = ref<MemberInfo | null>(null)
const rechargeAmount = ref<number>(0)
const remark = ref('')
const storeBalance = ref(120000)

const quickAmounts = [100, 500, 1000, 5000]

const calculatedPoints = computed(() => Math.floor(rechargeAmount.value || 0))

const afterRechargeBalance = computed(() => {
  if (!memberInfo.value) return 0
  return memberInfo.value.balance + calculatedPoints.value
})

const canSubmit = computed(() => {
  return memberInfo.value && rechargeAmount.value > 0
})

const handleSearch = () => {
  if (!searchQuery.value.trim()) {
    ElMessage.warning('请输入手机号或会员卡号')
    return
  }

  // TODO: 调用API查询会员信息
  // Mock数据
  memberInfo.value = {
    id: 'M888888',
    name: '张三',
    phone: '138****5678',
    balance: 5000,
    level: '普通会员'
  }

  ElMessage.success('查询成功')
}

const setAmount = (amount: number) => {
  rechargeAmount.value = amount
  calculatePoints()
}

const calculatePoints = () => {
  // 充值比例：1元 = 1积分
  // 已在 computed 中处理
}

const handleCancel = () => {
  ElMessageBox.confirm('确认取消充值？已填写的信息将不会保存', '提示', {
    confirmButtonText: '确认',
    cancelButtonText: '继续填写',
    type: 'warning'
  }).then(() => {
    resetForm()
  }).catch(() => {
    // 用户选择继续填写
  })
}

const handleSubmit = async () => {
  if (!canSubmit.value) return

  try {
    await ElMessageBox.confirm(
      `确认充值 ¥${rechargeAmount.value.toLocaleString()} ？`,
      '确认充值',
      {
        confirmButtonText: '确认',
        cancelButtonText: '取消',
        type: 'warning'
      }
    )
    await submitCRechargeEntry({
      memberId: memberInfo.value!.id,
      memberName: memberInfo.value!.name,
      memberPhone: memberInfo.value!.phone,
      centerId: '1',
      centerName: '北京朝阳中心',
      amount: rechargeAmount.value,
      paymentMethod: 'cash',
      remark: remark.value || ''
    })
    ElMessage.success(`充值成功！已为会员 ${memberInfo.value?.name} 充值 ${calculatedPoints.value.toLocaleString()} 积分`)
    resetForm()
  } catch {
    // 用户取消
  }
}

const resetForm = () => {
  searchQuery.value = ''
  memberInfo.value = null
  rechargeAmount.value = 0
  remark.value = ''
}
</script>

<style scoped>
.crecharge-entry {
  background-color: #F5F5F5;
  min-height: calc(100vh - 64px);
  display: flex;
  flex-direction: column;
}

.page-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  height: 64px;
  background-color: #FFFFFF;
  border-bottom: 1px solid #E5E5E5;
  padding: 16px 24px;
}

.page-title {
  font-family: 'Inter', sans-serif;
  font-size: 20px;
  font-weight: 600;
  color: #262626;
  margin: 0;
}

.content-area {
  flex: 1;
  padding: 24px;
  display: flex;
  flex-direction: column;
  gap: 24px;
}

.info-card {
  background-color: #FFFFFF;
  border-radius: 8px;
  border: 1px solid #E5E5E5;
  padding: 24px;
}

.card-header {
  display: flex;
  gap: 8px;
  align-items: center;
}

.header-icon {
  font-size: 18px;
}

.header-title {
  font-family: 'Inter', sans-serif;
  font-size: 16px;
  font-weight: 600;
  color: #262626;
  margin: 0;
}

.search-row {
  display: flex;
  gap: 12px;
  align-items: center;
}

.search-field {
  display: flex;
  gap: 8px;
  align-items: center;
  flex: 1;
}

.search-label {
  font-family: 'Inter', sans-serif;
  font-size: 14px;
  color: #262626;
  white-space: nowrap;
}

.search-btn {
  background-color: #C00000;
  border-color: #C00000;
  border-radius: 4px;
  height: 48px;
}

.member-info-box {
  background-color: #F9F9F9;
  border-radius: 4px;
  padding: 16px;
  display: flex;
  flex-direction: column;
  gap: 12px;
  margin-top: 16px;
}

.info-row {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.info-label {
  font-family: 'Inter', sans-serif;
  font-size: 14px;
  color: #595959;
}

.info-value {
  font-family: 'Inter', sans-serif;
  font-size: 14px;
  font-weight: 600;
  color: #262626;
}

.amount-section {
  display: flex;
  flex-direction: column;
  gap: 12px;
}

.amount-label {
  font-family: 'Inter', sans-serif;
  font-size: 14px;
  font-weight: 600;
  color: #262626;
}

.amount-input-wrapper {
  display: flex;
  align-items: center;
  gap: 8px;
  padding: 12px 16px;
  height: 56px;
  background-color: #FFFFFF;
  border: 2px solid #C00000;
  border-radius: 4px;
}

.currency-symbol {
  font-family: 'Inter', sans-serif;
  font-size: 24px;
  font-weight: 600;
  color: #C00000;
}

.amount-input {
  flex: 1;
}

.amount-input :deep(.el-input__inner) {
  border: none;
  padding: 0;
  font-size: 24px;
  font-weight: 600;
  color: #262626;
}

.amount-input :deep(.el-input-number__decrease),
.amount-input :deep(.el-input-number__increase) {
  display: none;
}

.quick-amounts {
  display: flex;
  gap: 8px;
}

.quick-amount-btn {
  width: 80px;
  height: 36px;
  border-radius: 4px;
  border: 1px solid #D9D9D9;
  background-color: #FFFFFF;
  color: #595959;
  font-size: 13px;
}

.quick-amount-btn:hover {
  border-color: #C00000;
  color: #C00000;
}

.ratio-notice {
  display: flex;
  gap: 8px;
  align-items: center;
  background-color: #FFF7E6;
  border-radius: 4px;
  padding: 12px;
}

.ratio-icon {
  font-size: 16px;
}

.ratio-text {
  font-family: 'Inter', sans-serif;
  font-size: 13px;
  color: #8C6000;
}

.points-calc-box {
  background-color: #F0F9FF;
  border-radius: 8px;
  padding: 16px;
  display: flex;
  flex-direction: column;
  gap: 12px;
}

.calc-row {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.calc-label {
  font-family: 'Inter', sans-serif;
  font-size: 14px;
  color: #595959;
}

.calc-value {
  font-family: 'Inter', sans-serif;
  font-size: 18px;
  font-weight: 600;
  color: #1677FF;
}

.store-balance {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding-top: 12px;
  border-top: 1px solid #E5E5E5;
}

.store-label {
  font-family: 'Inter', sans-serif;
  font-size: 14px;
  color: #595959;
}

.store-value {
  font-family: 'Inter', sans-serif;
  font-size: 14px;
  font-weight: 600;
  color: #52C41A;
}

.action-bar {
  display: flex;
  gap: 12px;
  justify-content: center;
  align-items: center;
  background-color: #FFFFFF;
  border-radius: 8px;
  border: 1px solid #E5E5E5;
  padding: 24px;
  margin: 0 24px 24px 24px;
}

.cancel-btn {
  width: 140px;
  height: 48px;
  border-radius: 4px;
}

.confirm-btn {
  width: 160px;
  height: 48px;
  border-radius: 4px;
  background-color: #C00000;
  border-color: #C00000;
  display: flex;
  gap: 8px;
  align-items: center;
}

.confirm-icon {
  font-size: 18px;
}

.confirm-btn:hover {
  background-color: #A00000;
  border-color: #A00000;
}

.confirm-btn:disabled {
  background-color: #D9D9D9;
  border-color: #D9D9D9;
}
</style>
