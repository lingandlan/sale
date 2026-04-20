<template>
  <div class="card-verify">
    <div class="page-header">
      <h1 class="page-title">门店卡消费核销</h1>
    </div>

    <div class="content-area">
      <!-- 卡号查询 -->
      <div class="info-card">
        <div class="header-title">卡号查询</div>
        <el-divider />
        <div class="search-row">
          <span class="search-label">卡号</span>
          <el-input v-model="cardNo" placeholder="请输入卡号" style="flex: 1" @keyup.enter="handleSearch" />
          <el-button type="primary" class="search-btn" @click="handleSearch">查询</el-button>
        </div>
      </div>

      <!-- 卡信息 -->
      <div v-if="cardInfo" class="info-card">
        <div class="header-title">卡信息</div>
        <el-divider />
        <div class="card-info-box">
          <div class="info-row">
            <span class="info-label">卡号</span>
            <span class="info-value">{{ cardInfo.cardNo }}</span>
          </div>
          <div class="info-row">
            <span class="info-label">卡余额</span>
            <span class="balance-value">¥{{ cardInfo.balance }}</span>
          </div>
          <div class="info-row">
            <span class="info-label">状态</span>
            <el-tag :type="CardStatusTagType[cardInfo.status] || 'info'" size="small">
              {{ CardStatusMap[cardInfo.status] || '未知' }}
            </el-tag>
          </div>
          <div v-if="cardInfo.activatedAt" class="info-row">
            <span class="info-label">激活日期</span>
            <span class="info-value">{{ cardInfo.activatedAt.slice(0, 10) }}</span>
          </div>
          <div v-if="cardInfo.expiredAt" class="info-row">
            <span class="info-label">过期日期</span>
            <span class="info-value">{{ cardInfo.expiredAt.slice(0, 10) }}</span>
          </div>
        </div>
      </div>

      <!-- 消费信息 -->
      <div v-if="cardInfo" class="info-card">
        <div class="header-title">消费信息</div>
        <el-divider />
        <div class="consume-section">
          <div class="field-label">消费金额（元）</div>
          <el-input-number v-model="consumeAmount" :min="100" :step="100" :precision="0" style="width: 100%"
            :disabled="cardInfo.status !== 2 && cardInfo.status !== 3" />
          <div class="field-hint">最低100元，正整数</div>
          <div class="balance-calc-box">
            消费后余额：¥{{ Math.max(0, cardInfo.balance - consumeAmount) }}
          </div>
        </div>
      </div>

      <!-- 备注 -->
      <div v-if="cardInfo" class="info-card">
        <div class="header-title">备注（可选）</div>
        <el-divider />
        <el-input v-model="remark" type="textarea" :rows="3" placeholder="请输入备注信息" />
      </div>
    </div>

    <!-- 操作栏 -->
    <div v-if="cardInfo" class="action-bar">
      <el-button class="cancel-btn" @click="handleCancel">取消</el-button>
      <el-button type="primary" class="confirm-btn" :disabled="!canSubmit" @click="handleSubmit">
        确认核销
      </el-button>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted, onUnmounted } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { verifyCard, consumeCard, CardStatusMap, CardStatusTagType } from '@/api/card'

interface CardInfo {
  cardNo: string
  cardType: number
  balance: number
  status: number
  activatedAt: string | null
  expiredAt: string | null
}

const cardNo = ref('')
const cardInfo = ref<CardInfo | null>(null)
const consumeAmount = ref(100)
const remark = ref('')
let timeInterval: number

const canSubmit = computed(() => {
  return cardInfo.value &&
    (cardInfo.value.status === 2 || cardInfo.value.status === 3) &&
    consumeAmount.value >= 100 &&
    consumeAmount.value <= cardInfo.value.balance
})

const handleSearch = async () => {
  if (!cardNo.value.trim()) { ElMessage.warning('请输入卡号'); return }
  try {
    const res = await verifyCard(cardNo.value)
    cardInfo.value = res?.data || res
    consumeAmount.value = 100
    ElMessage.success('查询成功')
  } catch (error: any) {
    ElMessage.error(error?.message || '查询失败')
    cardInfo.value = null
  }
}

const handleCancel = () => {
  if (cardInfo.value || remark.value) {
    ElMessageBox.confirm('确认取消？', '提示', { type: 'warning' }).then(() => resetForm()).catch(() => {})
  } else { resetForm() }
}

const handleSubmit = async () => {
  if (!canSubmit.value) return
  try {
    await ElMessageBox.confirm(
      `确认核销 ¥${consumeAmount.value}？`,
      '确认核销', { type: 'warning' }
    )
    await consumeCard({ cardNo: cardNo.value, amount: consumeAmount.value, remark: remark.value })
    ElMessage.success(`核销成功！卡号 ${cardNo.value} 已核销 ¥${consumeAmount.value}`)
    resetForm()
  } catch { /* cancelled */ }
}

const resetForm = () => { cardNo.value = ''; cardInfo.value = null; consumeAmount.value = 100; remark.value = '' }

onMounted(() => { timeInterval = window.setInterval(() => {}, 60000) })
onUnmounted(() => { if (timeInterval) clearInterval(timeInterval) })
</script>

<style scoped>
.card-verify { background: var(--color-bg); min-height: calc(100vh - 64px); display: flex; flex-direction: column; }
.page-header { display: flex; align-items: center; height: 64px; background: var(--color-bg-card); border-bottom: 1px solid var(--color-border); padding: 16px 24px; }
.page-title { font-size: 20px; font-weight: 600; color: var(--color-text-primary); margin: 0; }
.content-area { flex: 1; padding: 24px; display: flex; flex-direction: column; gap: 20px; }
.info-card { background: var(--color-bg-card); border-radius: var(--radius-md); border: 1px solid var(--color-border); padding: 20px; }
.header-title { font-size: 16px; font-weight: 600; color: var(--color-text-primary); }
.search-row { display: flex; gap: 12px; align-items: center; }
.search-label { font-size: 14px; color: var(--color-text-secondary); white-space: nowrap; }
.search-btn { background: var(--color-primary); border-color: var(--color-primary); height: 40px; }
.card-info-box { background: var(--color-bg-section); border-radius: var(--radius-sm); padding: 16px; display: flex; flex-direction: column; gap: 12px; }
.info-row { display: flex; justify-content: space-between; align-items: center; }
.info-label { font-size: 14px; color: var(--color-text-secondary); }
.info-value { font-size: 14px; font-weight: 600; color: var(--color-text-primary); }
.balance-value { font-size: 18px; font-weight: 600; color: var(--color-primary); }
.consume-section { display: flex; flex-direction: column; gap: 8px; }
.field-label { font-size: 14px; font-weight: 600; color: var(--color-text-primary); }
.field-hint { font-size: 12px; color: var(--color-text-muted); }
.balance-calc-box { background: var(--color-warning-bg); border-radius: var(--radius-md); padding: 16px; text-align: center; font-size: 16px; font-weight: 600; color: #D48806; }
.action-bar { display: flex; gap: 12px; justify-content: center; background: var(--color-bg-card); border-top: 1px solid var(--color-border); padding: 20px; }
.cancel-btn { width: 140px; height: 48px; }
.confirm-btn { width: 160px; height: 48px; background: var(--color-primary); border-color: var(--color-primary); }
.confirm-btn:disabled { background: var(--color-border); border-color: var(--color-border); }
</style>
