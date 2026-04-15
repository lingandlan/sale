<template>
  <div class="card-detail">
    <div class="page-header">
      <div class="back-button" @click="router.back()">
        <span class="back-icon">←</span>
        <span class="back-text">返回</span>
      </div>
      <h1 class="page-title">门店卡详情</h1>
    </div>

    <div class="content-area">
      <!-- 卡信息 -->
      <div class="info-card">
        <div class="card-header">卡信息</div>
        <el-divider />
        <div class="info-grid">
          <div class="info-row">
            <span class="info-label">卡号</span>
            <span class="info-value">{{ card?.cardNo }}</span>
          </div>
          <div class="info-row">
            <span class="info-label">卡类型</span>
            <span class="info-value">{{ CardTypeMap[card?.cardType] || '-' }}</span>
          </div>
          <div class="info-row">
            <span class="info-label">余额</span>
            <span class="balance-value">¥{{ card?.balance }}</span>
          </div>
          <div class="info-row">
            <span class="info-label">状态</span>
            <el-tag :type="CardStatusTagType[card?.status] || 'info'">
              {{ CardStatusMap[card?.status] || '未知' }}
            </el-tag>
          </div>
          <div class="info-row">
            <span class="info-label">批次号</span>
            <span class="info-value">{{ card?.batchNo || '-' }}</span>
          </div>
          <div class="info-row">
            <span class="info-label">发放原因</span>
            <span class="info-value">{{ card?.issueReason || '-' }}</span>
          </div>
          <div class="info-row">
            <span class="info-label">发放日期</span>
            <span class="info-value">{{ card?.issuedAt ? card.issuedAt.slice(0, 10) : '-' }}</span>
          </div>
          <div class="info-row">
            <span class="info-label">激活日期</span>
            <span class="info-value">{{ card?.activatedAt ? card.activatedAt.slice(0, 10) : '-' }}</span>
          </div>
          <div class="info-row">
            <span class="info-label">过期日期</span>
            <span class="info-value">{{ card?.expiredAt ? card.expiredAt.slice(0, 10) : '-' }}</span>
          </div>
        </div>
      </div>

      <!-- 交易记录 -->
      <div class="info-card">
        <div class="card-header">交易记录</div>
        <el-divider />
        <el-table :data="transactions" style="width: 100%">
          <el-table-column label="时间" width="180">
            <template #default="{ row }">{{ row.createdAt ? row.createdAt.slice(0, 19).replace('T', ' ') : '-' }}</template>
          </el-table-column>
          <el-table-column label="类型" width="100">
            <template #default="{ row }">
              <el-tag :type="row.type === 'consume' ? 'warning' : row.type === 'issue' || row.type === 'stock_in' ? 'success' : 'info'" size="small">
                {{ txnTypeText(row.type) }}
              </el-tag>
            </template>
          </el-table-column>
          <el-table-column label="金额" width="100">
            <template #default="{ row }">
              <span :class="row.type === 'consume' ? 'text-red' : 'text-green'">
                {{ row.amount > 0 ? (row.type === 'consume' ? '-' : '+') : '' }}¥{{ row.amount }}
              </span>
            </template>
          </el-table-column>
          <el-table-column label="余额" width="100">
            <template #default="{ row }">¥{{ row.balanceAfter }}</template>
          </el-table-column>
          <el-table-column prop="remark" label="备注" />
        </el-table>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { ElMessage } from 'element-plus'
import { getCardDetail, CardStatusMap, CardStatusTagType, CardTypeMap } from '@/api/card'

const router = useRouter()
const route = useRoute()

const card = ref<any>(null)
const transactions = ref<any[]>([])

const txnTypeText = (type: string) => {
  const map: Record<string, string> = {
    stock_in: '入库', issue: '发放', consume: '消费',
    freeze: '冻结', unfreeze: '解冻', activate: '激活', void: '作废'
  }
  return map[type] || type
}

const loadDetail = async () => {
  const cardNo = route.params.cardNo as string
  try {
    const res = await getCardDetail(cardNo)
    const data = res?.data || res
    card.value = data.card
    transactions.value = data.transactions || []
  } catch {
    ElMessage.error('加载详情失败')
  }
}

onMounted(() => { loadDetail() })
</script>

<style scoped>
.card-detail { background: #F5F5F5; min-height: calc(100vh - 64px); }
.page-header { display: flex; gap: 16px; align-items: center; height: 64px; background: #FFF; border-bottom: 1px solid #E5E5E5; padding: 16px 24px; }
.back-button { display: flex; gap: 8px; align-items: center; cursor: pointer; }
.back-icon { font-size: 20px; color: #262626; }
.back-text { font-size: 14px; color: #262626; }
.page-title { flex: 1; font-size: 20px; font-weight: 600; color: #262626; margin: 0; }
.content-area { padding: 24px; display: flex; flex-direction: column; gap: 20px; }
.info-card { background: #FFF; border-radius: 8px; border: 1px solid #E5E5E5; padding: 20px; }
.card-header { font-size: 16px; font-weight: 600; color: #262626; }
.info-grid { display: grid; grid-template-columns: repeat(2, 1fr); gap: 12px; }
.info-row { display: flex; justify-content: space-between; align-items: center; }
.info-label { font-size: 14px; color: #595959; }
.info-value { font-size: 14px; font-weight: 600; color: #262626; }
.balance-value { font-size: 18px; font-weight: 600; color: #C00000; }
.text-red { color: #FF4D4F; }
.text-green { color: #52C41A; }
</style>
