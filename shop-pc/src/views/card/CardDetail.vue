<template>
  <div class="card-detail">
    <div class="page-header">
      <div class="back-button" @click="handleBack">
        <span class="back-icon">←</span>
        <span class="back-text">返回</span>
      </div>
      <h1 class="page-title">门店卡详情</h1>
    </div>

    <div class="content-area">
      <!-- 卡信息 -->
      <div class="info-card">
        <div class="card-header">📇 卡信息</div>
        <el-divider />
        <div class="info-grid">
          <div class="info-row">
            <span class="info-label">卡号</span>
            <span class="info-value">{{ cardInfo?.cardNo }}</span>
          </div>
          <div class="info-row">
            <span class="info-label">持卡人</span>
            <span class="info-value">{{ cardInfo?.holder }}</span>
          </div>
          <div class="info-row">
            <span class="info-label">手机号</span>
            <span class="info-value">{{ cardInfo?.holderPhone }}</span>
          </div>
          <div class="info-row">
            <span class="info-label">卡余额</span>
            <span class="balance-value">¥{{ cardInfo?.balance.toLocaleString() }}</span>
          </div>
          <div class="info-row">
            <span class="info-label">状态</span>
            <el-tag :type="cardInfo?.status === 'active' ? 'success' : 'danger'">
              {{ cardInfo?.status === 'active' ? '已发放' : '已停用' }}
            </el-tag>
          </div>
          <div class="info-row">
            <span class="info-label">发放日期</span>
            <span class="info-value">{{ cardInfo?.issueDate }}</span>
          </div>
          <div class="info-row">
            <span class="info-label">过期日期</span>
            <span class="info-value">{{ cardInfo?.expiryDate }}</span>
          </div>
          <div class="info-row">
            <span class="info-label">发放中心</span>
            <span class="info-value">{{ cardInfo?.issueCenter }}</span>
          </div>
        </div>
      </div>

      <!-- 交易记录 -->
      <div class="info-card">
        <div class="card-header">📜 交易记录</div>
        <el-divider />
        <el-table
          :data="cardInfo?.transactions"
          style="width: 100%"
          :header-cell-style="{
            backgroundColor: '#FAFAFA',
            color: '#262626',
            fontWeight: '600'
          }"
        >
          <el-table-column prop="time" label="交易时间" width="180" />
          <el-table-column label="交易类型" width="100">
            <template #default="{ row }">
              <el-tag :type="getTransactionType(row.type)" size="small">
                {{ getTransactionText(row.type) }}
              </el-tag>
            </template>
          </el-table-column>
          <el-table-column label="金额" width="120">
            <template #default="{ row }">
              <span :class="{ 'text-red': row.type === 'consume', 'text-green': row.type === 'issue' || row.type === 'recharge' }">
                {{ row.type === 'consume' ? '-' : '+' }}¥{{ row.amount.toLocaleString() }}
              </span>
            </template>
          </el-table-column>
          <el-table-column label="交易后余额" width="120">
            <template #default="{ row }">
              ¥{{ row.balanceAfter.toLocaleString() }}
            </template>
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
import { getCardDetail } from '@/api/card'

const router = useRouter()
const route = useRoute()

const cardInfo = ref<any>(null)

const getTransactionType = (type: string) => {
  switch (type) {
    case 'issue':
    case 'recharge':
      return 'success'
    case 'consume':
      return 'warning'
    default:
      return 'info'
  }
}

const getTransactionText = (type: string) => {
  switch (type) {
    case 'issue':
      return '发放'
    case 'consume':
      return '消费'
    case 'recharge':
      return '充值'
    default:
      return '未知'
  }
}

const handleBack = () => {
  router.back()
}

const loadDetail = async () => {
  const cardNo = route.params.cardNo as string
  try {
    const res = await getCardDetail(cardNo)
    if (res?.data) {
      cardInfo.value = res.data
    }
  } catch (error) {
    ElMessage.error('加载详情失败')
  }
}

onMounted(() => {
  loadDetail()
})
</script>

<style scoped>
.card-detail {
  background-color: #F5F5F5;
  min-height: calc(100vh - 64px);
}

.page-header {
  display: flex;
  gap: 16px;
  align-items: center;
  height: 64px;
  background-color: #FFFFFF;
  border-bottom: 1px solid #E5E5E5;
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
  color: #262626;
}

.back-text {
  font-family: 'Inter', sans-serif;
  font-size: 14px;
  color: #262626;
}

.page-title {
  flex: 1;
  font-family: 'Inter', sans-serif;
  font-size: 20px;
  font-weight: 600;
  color: #262626;
  margin: 0;
}

.content-area {
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
  font-family: 'Inter', sans-serif;
  font-size: 16px;
  font-weight: 600;
  color: #262626;
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

.text-red {
  color: #FF4D4F;
}

.text-green {
  color: #52C41A;
}
</style>
