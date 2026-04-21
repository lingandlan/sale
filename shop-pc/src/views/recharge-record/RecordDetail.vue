<template>
  <div class="record-detail">
    <div class="page-header">
      <div class="back-button" @click="handleBack">
        <span class="back-icon">←</span>
        <span class="back-text">返回</span>
      </div>
      <h1 class="page-title">充值记录详情</h1>
    </div>

    <div class="content-area">
      <div class="info-card">
        <div class="card-header">充值信息</div>
        <el-divider />
        <div class="info-grid">
          <div class="info-row">
            <span class="info-label">交易单号</span>
            <span class="info-value">{{ record.id || '-' }}</span>
          </div>
          <div class="info-row">
            <span class="info-label">会员姓名</span>
            <span class="info-value">{{ record.memberName || '-' }}</span>
          </div>
          <div class="info-row">
            <span class="info-label">手机号</span>
            <span class="info-value">{{ record.memberPhone || '-' }}</span>
          </div>
          <div class="info-row">
            <span class="info-label">充值金额</span>
            <span class="amount-value">{{ record.amount != null ? `¥${Number(record.amount).toLocaleString()}` : '-' }}</span>
          </div>
          <div class="info-row">
            <span class="info-label">获得积分</span>
            <span class="info-value">{{ record.points != null ? `${Number(record.points).toLocaleString()} 积分` : '-' }}</span>
          </div>
          <div class="info-row">
            <span class="info-label">充值前余额</span>
            <span class="info-value">{{ record.balanceBefore != null ? `${Number(record.balanceBefore).toLocaleString()} 积分` : '-' }}</span>
          </div>
          <div class="info-row">
            <span class="info-label">充值后余额</span>
            <span class="info-value text-green">{{ record.balanceAfter != null ? `${Number(record.balanceAfter).toLocaleString()} 积分` : '-' }}</span>
          </div>
          <div class="info-row">
            <span class="info-label">充值中心</span>
            <span class="info-value">{{ record.centerName || '-' }}</span>
          </div>
          <div class="info-row">
            <span class="info-label">操作员</span>
            <span class="info-value">{{ record.operatorName || '-' }}</span>
          </div>
          <div class="info-row">
            <span class="info-label">充值时间</span>
            <span class="info-value">{{ formatTime(record.createdAt) }}</span>
          </div>
          <div class="info-row" v-if="record.remark">
            <span class="info-label">备注</span>
            <span class="info-value">{{ record.remark }}</span>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { ElMessage } from 'element-plus'
import { extractErrorMessage } from '@/utils/request'
import { getRechargeRecordDetail } from '@/api/recharge'

const router = useRouter()
const route = useRoute()

const formatTime = (t: string) => {
  if (!t) return '-'
  return t.replace('T', ' ').slice(0, 19)
}

const record = ref<any>({})

const handleBack = () => {
  router.back()
}

const loadData = async () => {
  const id = route.params.id as string
  try {
    const res = await getRechargeRecordDetail(id)
    if (res?.data) {
      record.value = res.data
    }
  } catch (err: any) {
    ElMessage.error(extractErrorMessage(err, '加载充值记录详情失败'))
  }
}

onMounted(() => {
  loadData()
})
</script>

<style scoped>
.record-detail {
  background-color: var(--color-bg);
  min-height: calc(100vh - 64px);
}

.page-header {
  display: flex;
  gap: 16px;
  align-items: center;
  height: 64px;
  background-color: var(--color-bg-card);
  border-bottom: 1px solid var(--color-border);
  padding: 16px 24px;
}

.back-button {
  display: flex;
  gap: 8px;
  align-items: center;
  cursor: pointer;
}

.back-icon {
  font-size: 20px;
  color: var(--color-text-primary);
}

.back-text {
  font-size: 14px;
  color: var(--color-text-primary);
}

.page-title {
  flex: 1;
  font-size: 20px;
  font-weight: 600;
  color: var(--color-text-primary);
  margin: 0;
}

.content-area {
  padding: 24px;
}

.info-card {
  background-color: var(--color-bg-card);
  border-radius: var(--radius-md);
  border: 1px solid var(--color-border);
  padding: 24px;
}

.card-header {
  font-size: 16px;
  font-weight: 600;
  color: var(--color-text-primary);
}

.info-grid {
  display: grid;
  grid-template-columns: repeat(2, 1fr);
  gap: 12px;
}

.info-row {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.info-label {
  font-size: 14px;
  color: var(--color-text-secondary);
}

.info-value {
  font-size: 14px;
  font-weight: 600;
  color: var(--color-text-primary);
}

.amount-value {
  font-size: 18px;
  font-weight: 600;
  color: var(--color-primary);
}

.text-green {
  color: var(--color-success);
}
</style>
