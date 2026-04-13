<template>
  <div class="card-verify">
    <div class="page-header">
      <h1 class="page-title">门店卡消费核销</h1>
    </div>

    <div class="content-area">
      <!-- 卡号查询卡片 -->
      <div class="info-card">
        <div class="card-header">
          <span class="header-title">🔍 卡号查询</span>
        </div>
        <el-divider />
        <div class="search-row">
          <span class="search-label">卡号</span>
          <el-input
            v-model="cardNo"
            placeholder="请输入或扫描卡号"
            style="flex: 1"
            @keyup.enter="handleSearch"
          />
          <el-button type="primary" class="search-btn" @click="handleSearch">
            查询
          </el-button>
        </div>
      </div>

      <!-- 卡信息卡片 -->
      <div v-if="cardInfo" class="info-card">
        <div class="card-header">
          <span class="header-title">📇 卡信息详情</span>
        </div>
        <el-divider />
        <div class="card-info-box">
          <div class="info-row">
            <span class="info-label">卡号</span>
            <span class="info-value">{{ cardInfo.cardNo }}</span>
          </div>
          <div class="info-row">
            <span class="info-label">持卡人</span>
            <span class="info-value">{{ cardInfo.holder }}</span>
          </div>
          <div class="info-row">
            <span class="info-label">卡余额</span>
            <span class="balance-value">¥{{ cardInfo.balance.toLocaleString() }}</span>
          </div>
          <div class="info-row">
            <span class="info-label">卡状态</span>
            <el-tag :type="cardInfo.status === 'active' ? 'success' : 'danger'" size="small">
              {{ cardInfo.status === 'active' ? '✓ 已发放' : '已停用' }}
            </el-tag>
          </div>
        </div>
      </div>

      <!-- 消费信息卡片 -->
      <div v-if="cardInfo" class="info-card">
        <div class="card-header">
          <span class="header-title">💰 消费信息</span>
        </div>
        <el-divider />
        <div class="consume-section">
          <div class="field-group">
            <div class="field-label">消费金额（元）</div>
            <el-input-number
              v-model="consumeAmount"
              :min="100"
              :step="100"
              :precision="0"
              style="width: 100%"
              :disabled="!cardInfo || cardInfo.status !== 'active'"
            />
            <div class="field-hint">（最低100元，正整数）</div>
          </div>

          <div class="time-row">
            <span class="time-label">消费时间：</span>
            <span class="time-value">{{ currentTime }}</span>
          </div>

          <div class="balance-calc-box">
            <span class="balance-calc-text">消费后余额：¥{{ afterConsumeBalance.toLocaleString() }}</span>
          </div>
        </div>
      </div>

      <!-- 备注卡片 -->
      <div v-if="cardInfo" class="info-card">
        <div class="card-header">
          <span class="header-title">📝 备注（可选）</span>
        </div>
        <el-divider />
        <el-input
          v-model="remark"
          type="textarea"
          :rows="3"
          placeholder="请输入备注信息（可选）"
        />
      </div>
    </div>

    <!-- 操作栏 -->
    <div class="action-bar">
      <el-button class="cancel-btn" @click="handleCancel">取消</el-button>
      <el-button
        type="primary"
        class="confirm-btn"
        :disabled="!canSubmit"
        @click="handleSubmit"
      >
        <span class="confirm-icon">✓</span>
        确认核销
      </el-button>
    </div>

    <!-- 核销记录 -->
    <div class="info-card" style="margin: 0 24px 24px 24px;">
      <div class="card-header">
        <span class="header-title">核销记录</span>
      </div>
      <el-divider />
      <el-table :data="consumeRecords" stripe style="width: 100%">
        <el-table-column label="序号" width="80" align="center">
          <template #default="{ $index }">{{ $index + 1 }}</template>
        </el-table-column>
        <el-table-column prop="cardNo" label="卡号" width="160" />
        <el-table-column prop="amount" label="核销金额" width="140" align="right">
          <template #default="{ row }">¥{{ row.amount.toLocaleString() }}</template>
        </el-table-column>
        <el-table-column prop="time" label="核销时间" width="200" />
        <el-table-column prop="operator" label="操作员" width="140" />
        <el-table-column prop="remark" label="备注" />
      </el-table>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted, onUnmounted } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { verifyCard, consumeCard } from '@/api/card'

interface CardInfo {
  cardNo: string
  holder: string
  balance: number
  status: 'active' | 'inactive'
}

const cardNo = ref('')
const cardInfo = ref<CardInfo | null>(null)
const consumeAmount = ref<number>(100)
const remark = ref('')
const currentTime = ref('')
const consumeRecords = ref([
  { cardNo: 'TJ00000001', amount: 200, time: '2026-04-10 14:30', operator: '张三', remark: '' },
  { cardNo: 'TJ00000003', amount: 500, time: '2026-04-10 11:20', operator: '李四', remark: '部分消费' },
  { cardNo: 'TJ00000005', amount: 1000, time: '2026-04-09 16:45', operator: '张三', remark: '' }
])

let timeInterval: number

const canSubmit = computed(() => {
  return cardInfo.value &&
         cardInfo.value.status === 'active' &&
         consumeAmount.value >= 100 &&
         consumeAmount.value <= cardInfo.value.balance
})

const afterConsumeBalance = computed(() => {
  if (!cardInfo.value) return 0
  return Math.max(0, cardInfo.value.balance - (consumeAmount.value || 0))
})

const updateTime = () => {
  const now = new Date()
  const year = now.getFullYear()
  const month = String(now.getMonth() + 1).padStart(2, '0')
  const day = String(now.getDate()).padStart(2, '0')
  const hour = String(now.getHours()).padStart(2, '0')
  const minute = String(now.getMinutes()).padStart(2, '0')
  currentTime.value = `${year}-${month}-${day} ${hour}:${minute}`
}

const handleSearch = async () => {
  if (!cardNo.value.trim()) {
    ElMessage.warning('请输入卡号')
    return
  }

  try {
    const res = await verifyCard(cardNo.value)
    if (res?.data) {
      cardInfo.value = res.data
    }
    consumeAmount.value = 100
    ElMessage.success('查询成功')
  } catch (error) {
    ElMessage.error('查询失败')
  }
}

const handleCancel = () => {
  if (cardInfo.value || remark.value) {
    ElMessageBox.confirm('确认取消？已填写的信息将不会保存', '提示', {
      confirmButtonText: '确认',
      cancelButtonText: '继续填写',
      type: 'warning'
    }).then(() => {
      resetForm()
    }).catch(() => {
      // 用户选择继续填写
    })
  } else {
    resetForm()
  }
}

const handleSubmit = async () => {
  if (!canSubmit.value) return

  if (consumeAmount.value > cardInfo.value!.balance) {
    ElMessage.warning('消费金额不能超过卡余额')
    return
  }

  try {
    await ElMessageBox.confirm(
      `确认核销 ¥${consumeAmount.value}？核销后余额：¥${afterConsumeBalance.value.toLocaleString()}`,
      '确认核销',
      {
        confirmButtonText: '确认',
        cancelButtonText: '取消',
        type: 'warning'
      }
    )
    await consumeCard({
      cardNo: cardNo.value,
      amount: consumeAmount.value,
      remark: remark.value
    })
    ElMessage.success(`核销成功！卡号 ${cardNo.value} 已核销 ¥${consumeAmount.value}`)
    resetForm()
  } catch {
    // 用户取消
  }
}

const resetForm = () => {
  cardNo.value = ''
  cardInfo.value = null
  consumeAmount.value = 100
  remark.value = ''
}

onMounted(() => {
  updateTime()
  timeInterval = window.setInterval(updateTime, 60000)
})

onUnmounted(() => {
  if (timeInterval) {
    clearInterval(timeInterval)
  }
})
</script>

<style scoped>
.card-verify {
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

.header-title {
  font-family: 'Inter', sans-serif;
  font-size: 16px;
  font-weight: 600;
  color: #262626;
}

.search-row {
  display: flex;
  gap: 12px;
  align-items: center;
}

.search-label {
  font-family: 'Inter', sans-serif;
  font-size: 14px;
  color: #595959;
  white-space: nowrap;
}

.search-btn {
  background-color: #C00000;
  border-color: #C00000;
  border-radius: 4px;
  height: 40px;
}

.card-info-box {
  background-color: #F9F9F9;
  border-radius: 4px;
  padding: 16px;
  display: flex;
  flex-direction: column;
  gap: 12px;
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

.balance-value {
  font-family: 'Inter', sans-serif;
  font-size: 18px;
  font-weight: 600;
  color: #C00000;
}

.consume-section {
  display: flex;
  flex-direction: column;
  gap: 12px;
}

.field-group {
  display: flex;
  flex-direction: column;
  gap: 12px;
}

.field-label {
  font-family: 'Inter', sans-serif;
  font-size: 14px;
  font-weight: 600;
  color: #262626;
}

.field-hint {
  font-family: 'Inter', sans-serif;
  font-size: 12px;
  color: #8C8C8C;
}

.time-row {
  display: flex;
  gap: 8px;
  align-items: center;
}

.time-label {
  font-family: 'Inter', sans-serif;
  font-size: 14px;
  color: #595959;
}

.time-value {
  font-family: 'Inter', sans-serif;
  font-size: 14px;
  color: #262626;
}

.balance-calc-box {
  background-color: #FFF7E6;
  border-radius: 8px;
  padding: 16px;
  display: flex;
  justify-content: center;
  align-items: center;
}

.balance-calc-text {
  font-family: 'Inter', sans-serif;
  font-size: 16px;
  font-weight: 600;
  color: #D48806;
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
  margin: 0 24px;
  height: 80px;
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
  font-size: 16px;
}

.confirm-btn:hover:not(:disabled) {
  background-color: #A00000;
  border-color: #A00000;
}

.confirm-btn:disabled {
  background-color: #D9D9D9;
  border-color: #D9D9D9;
}
</style>
