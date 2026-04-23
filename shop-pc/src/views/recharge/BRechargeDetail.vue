<template>
  <div class="brecharge-detail">
    <div class="page-header">
      <div class="back-button" @click="handleBack">
        <span class="back-icon">←</span>
        <span class="back-text">返回</span>
      </div>
      <h1 class="page-title">B端充值审批详情</h1>
      <el-tag :type="getStatusType(detailData.status)" size="small">
        {{ getStatusText(detailData.status) }}
      </el-tag>
    </div>

    <div class="content-area">
      <!-- 基本信息 -->
      <div class="info-card">
        <div class="card-header">
          <span class="header-icon">📋</span>
          <h3 class="header-title">申请基本信息</h3>
        </div>
        <el-divider />
        <div class="info-grid">
          <div class="info-row">
            <span class="info-label">申请单号</span>
            <span class="info-value">{{ detailData.id }}</span>
          </div>
          <div class="info-row">
            <span class="info-label">充值中心</span>
            <span class="info-value">{{ detailData.centerName }}</span>
          </div>
          <div class="info-row">
            <span class="info-label">申请人</span>
            <span class="info-value">{{ detailData.applicant }}</span>
          </div>
          <div class="info-row">
            <span class="info-label">申请时间</span>
            <span class="info-value">{{ detailData.createdAt }}</span>
          </div>
        </div>
      </div>

      <!-- 充值详情 -->
      <div class="info-card">
        <div class="card-header">
          <span class="header-icon">💰</span>
          <h3 class="header-title">充值详情</h3>
        </div>
        <el-divider />
        <div class="info-grid">
          <div class="info-row">
            <span class="info-label">上月商城净消费</span>
            <span class="info-value">¥{{ detailData.lastMonthConsumption?.toLocaleString() || '0' }}</span>
          </div>
          <div class="info-row">
            <span class="info-label">适用返还比例</span>
            <el-tag size="small" type="warning">{{ detailData.points?.rebateRate || 0 }}%</el-tag>
          </div>
          <div class="info-row">
            <span class="info-label">充值金额</span>
            <span class="amount-value">¥{{ detailData.amount.toLocaleString() }}</span>
          </div>
          <div class="info-row with-border">
            <span class="info-label">当前中心余额</span>
            <span class="balance-value">¥{{ detailData.currentBalance?.toLocaleString() || '0' }} 积分</span>
          </div>
        </div>

        <!-- 积分计算 -->
        <div class="points-box">
          <div class="points-title">预计获得积分</div>
          <div class="points-number">
            <span class="points-value">{{ formatTotalPoints }}</span>
            <span class="points-unit">积分</span>
          </div>
          <el-divider />
          <div class="points-details">
            <div class="detail-line">• 基础积分：{{ detailData.points?.base?.toLocaleString() || '0' }}</div>
            <div class="detail-line">• 返还积分（{{ detailData.points?.rebateRate || 0 }}%）：{{ detailData.points?.rebate?.toLocaleString() || '0' }}</div>
            <div class="detail-line highlight">• 总计（向上取整）：{{ formatTotalPoints }}</div>
          </div>
        </div>
      </div>

      <!-- 付款凭证 -->
      <div class="info-card">
        <div class="card-header">
          <span class="header-icon">📎</span>
          <h3 class="header-title">付款凭证</h3>
        </div>
        <el-divider />
        <div class="payment-info">
          <div class="payment-field">
            <div class="payment-label">银行流水单号</div>
            <div class="payment-value">{{ detailData.transactionNo || '未填写' }}</div>
          </div>
          <div class="payment-field">
            <div class="payment-label">付款截图</div>
            <div v-if="detailData.screenshot" class="screenshot-preview">
              <el-image
                v-if="detailData.screenshot.startsWith('http') || detailData.screenshot.startsWith('/uploads/')"
                :src="detailData.screenshot"
                fit="cover"
                style="width: 160px; height: 100px; border-radius: var(--radius-sm)"
                :preview-src-list="[detailData.screenshot]"
              />
              <div v-else class="no-screenshot">
                <span class="no-screenshot-icon">📎</span>
                <span class="no-screenshot-text">{{ detailData.screenshot }}</span>
              </div>
            </div>
            <div v-else class="no-screenshot">
              <span class="no-screenshot-icon">📷</span>
              <span class="no-screenshot-text">未上传截图</span>
            </div>
          </div>
        </div>
      </div>

      <!-- 备注信息 -->
      <div v-if="detailData.remark" class="info-card">
        <div class="card-header">
          <span class="header-icon">📝</span>
          <h3 class="header-title">备注信息</h3>
        </div>
        <el-divider />
        <div class="remark-content">
          {{ detailData.remark }}
        </div>
      </div>

      <!-- 审批操作 -->
      <div v-if="detailData.status === 'pending'" class="approval-card">
        <div class="card-header">
          <span class="header-icon">✍️</span>
          <h3 class="header-title">审批操作</h3>
        </div>
        <el-divider />
        <el-form ref="approvalFormRef" :model="approvalForm" :rules="approvalRules">
          <el-form-item label="审批意见" prop="opinion">
            <el-input
              v-model="approvalForm.opinion"
              type="textarea"
              :rows="4"
              placeholder="请输入审批意见（必填）"
            />
          </el-form-item>
          <div class="approval-buttons">
            <el-button type="danger" class="reject-btn" @click="handleReject">
              拒绝
            </el-button>
            <el-button type="success" class="approve-btn" @click="handleApprove">
              通过
            </el-button>
          </div>
        </el-form>
      </div>

      <!-- 审批历史 -->
      <div v-if="detailData.status !== 'pending'" class="info-card">
        <div class="card-header">
          <span class="header-icon">📜</span>
          <h3 class="header-title">审批结果</h3>
        </div>
        <el-divider />
        <div class="approval-result">
          <div class="result-row">
            <span class="result-label">审批人</span>
            <span class="result-value">{{ detailData.approvedBy || '-' }}</span>
          </div>
          <div class="result-row">
            <span class="result-label">审批时间</span>
            <span class="result-value">{{ detailData.approvedAt || '-' }}</span>
          </div>
          <div v-if="detailData.approvalRemark" class="result-row">
            <span class="result-label">审批意见</span>
            <span class="result-value">{{ detailData.approvalRemark }}</span>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { ElMessage, type FormInstance, type FormRules } from 'element-plus'
import { extractErrorMessage } from '@/utils/request'
import { getBRechargeApprovalDetail, approvalAction } from '@/api/recharge'

const router = useRouter()
const route = useRoute()

const approvalFormRef = ref<FormInstance>()

const detailData = ref<any>({
  id: '',
  centerName: '',
  amount: 0,
  points: {
    base: 0,
    rebate: 0,
    rebateRate: 0,
    total: 0
  },
  applicant: '',
  createdAt: '',
  transactionNo: '',
  screenshot: '',
  remark: '',
  status: 'pending',
  approvedBy: '',
  approvedAt: '',
  approvalRemark: '',
  lastMonthConsumption: 0,
  currentBalance: 0
})

const approvalForm = ref({
  opinion: ''
})

const approvalRules: FormRules = {
  opinion: [{ required: true, message: '请输入审批意见', trigger: 'blur' }]
}

const formatTotalPoints = computed(() => {
  return detailData.value.points?.total?.toLocaleString() || '0'
})

const getStatusType = (status: string) => {
  switch (status) {
    case 'pending':
      return 'warning'
    case 'approved':
      return 'success'
    case 'rejected':
      return 'danger'
    default:
      return 'info'
  }
}

const getStatusText = (status: string) => {
  switch (status) {
    case 'pending':
      return '待审批'
    case 'approved':
      return '已通过'
    case 'rejected':
      return '已拒绝'
    default:
      return '未知'
  }
}

const loadDetail = async () => {
  const id = route.params.id as string
  try {
    const res = await getBRechargeApprovalDetail(id)
    if (res?.data) {
      const d = res.data
      detailData.value = {
        ...d,
        applicant: d.applicant?.name || d.applicant || '',
        lastMonthConsumption: d.lastMonthConsumption || 0,
        currentBalance: d.currentBalance || 0
      }
    }
  } catch (err: any) {
    ElMessage.error(extractErrorMessage(err, '加载详情失败'))
  }
}

const handleBack = () => {
  router.back()
}

const handleReject = async () => {
  if (!approvalFormRef.value) return

  try {
    await approvalFormRef.value.validate()
  } catch {
    return
  }

  try {
    await approvalAction({
      id: detailData.value.id,
      action: 'reject',
      remark: approvalForm.value.opinion
    })
    ElMessage.success('已拒绝')
    router.back()
  } catch (err: any) {
    ElMessage.error(extractErrorMessage(err, '拒绝失败'))
  }
}

const handleApprove = async () => {
  if (!approvalFormRef.value) return

  try {
    await approvalFormRef.value.validate()
  } catch {
    return
  }

  try {
    await approvalAction({
      id: detailData.value.id,
      action: 'approve',
      remark: approvalForm.value.opinion
    })
    ElMessage.success('已通过')
    router.back()
  } catch (err: any) {
    ElMessage.error(extractErrorMessage(err, '审批通过失败'))
  }
}

onMounted(() => {
  loadDetail()
})
</script>

<style scoped>
.brecharge-detail {
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
  transition: opacity 0.3s;
}

.back-button:hover {
  opacity: 0.7;
}

.back-icon {
  font-size: 20px;
  color: var(--color-text-primary);
}

.back-text {
  font-family: var(--font-family);
  font-size: 14px;
  color: var(--color-text-primary);
}

.page-title {
  flex: 1;
  font-family: var(--font-family);
  font-size: 20px;
  font-weight: 600;
  color: var(--color-text-primary);
  margin: 0;
}

.content-area {
  padding: 24px;
  display: flex;
  flex-direction: column;
  gap: 24px;
}

.info-card {
  background-color: var(--color-bg-card);
  border-radius: var(--radius-md);
  border: 1px solid var(--color-border);
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
  font-family: var(--font-family);
  font-size: 16px;
  font-weight: 600;
  color: var(--color-text-primary);
  margin: 0;
}

.info-grid {
  display: flex;
  flex-direction: column;
  gap: 12px;
}

.info-row {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.info-row.with-border {
  padding-top: 12px;
  border-top: 1px solid var(--color-border);
}

.info-label {
  font-family: var(--font-family);
  font-size: 14px;
  color: var(--color-text-secondary);
}

.info-value {
  font-family: var(--font-family);
  font-size: 14px;
  font-weight: 600;
  color: var(--color-text-primary);
}

.amount-value {
  font-family: var(--font-family);
  font-size: 24px;
  font-weight: 700;
  color: var(--color-primary);
}

.balance-value {
  font-family: var(--font-family);
  font-size: 16px;
  font-weight: 600;
  color: var(--color-success);
}

.points-box {
  background-color: var(--color-warning-bg);
  border-radius: var(--radius-md);
  padding: 16px;
  display: flex;
  flex-direction: column;
  gap: 12px;
  margin-top: 16px;
}

.points-title {
  font-family: var(--font-family);
  font-size: 14px;
  font-weight: 600;
  color: var(--color-text-primary);
}

.points-number {
  display: flex;
  gap: 8px;
  justify-content: center;
  align-items: baseline;
  padding: 12px;
}

.points-value {
  font-family: var(--font-family);
  font-size: 32px;
  font-weight: 700;
  color: var(--color-primary);
}

.points-unit {
  font-family: var(--font-family);
  font-size: 16px;
  color: var(--color-primary);
}

.points-details {
  display: flex;
  flex-direction: column;
  gap: 8px;
}

.detail-line {
  font-family: var(--font-family);
  font-size: 13px;
  color: var(--color-text-secondary);
}

.detail-line.highlight {
  font-weight: 600;
  color: var(--color-primary);
}

.payment-info {
  display: flex;
  flex-direction: column;
  gap: 12px;
}

.payment-field {
  display: flex;
  flex-direction: column;
  gap: 8px;
}

.payment-label {
  font-family: var(--font-family);
  font-size: 14px;
  color: var(--color-text-secondary);
}

.payment-value {
  font-family: var(--font-family);
  font-size: 14px;
  font-weight: 600;
  color: var(--color-text-primary);
}

.screenshot-preview {
  width: 160px;
  height: 100px;
}

.no-screenshot {
  width: 160px;
  height: 100px;
  background-color: var(--color-bg);
  border: 1px solid var(--color-border);
  border-radius: var(--radius-sm);
  display: flex;
  flex-direction: column;
  gap: 8px;
  justify-content: center;
  align-items: center;
}

.no-screenshot-icon {
  font-size: 24px;
}

.no-screenshot-text {
  font-family: var(--font-family);
  font-size: 12px;
  color: var(--color-text-muted);
}

.remark-content {
  background-color: var(--color-bg-section);
  border-radius: var(--radius-sm);
  padding: 16px;
  font-family: var(--font-family);
  font-size: 14px;
  color: var(--color-text-primary);
  line-height: 1.6;
}

.approval-card {
  background-color: var(--color-bg-card);
  border-radius: var(--radius-md);
  border: 2px solid var(--color-primary);
  padding: 24px;
}

.approval-buttons {
  display: flex;
  gap: 12px;
  justify-content: center;
  margin-top: 16px;
}

.reject-btn,
.approve-btn {
  width: 140px;
  height: 48px;
  border-radius: var(--radius-sm);
}

.approval-result {
  display: flex;
  flex-direction: column;
  gap: 12px;
}

.result-row {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.result-label {
  font-family: var(--font-family);
  font-size: 14px;
  color: var(--color-text-secondary);
}

.result-value {
  font-family: var(--font-family);
  font-size: 14px;
  font-weight: 600;
  color: var(--color-text-primary);
}
</style>
