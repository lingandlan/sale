<template>
  <div class="card-manage">
    <div class="page-header">
      <h1 class="page-title">门店卡管理</h1>
      <div class="header-actions">
        <el-button class="action-btn" @click="handleRefresh">刷新</el-button>
        <el-button class="action-btn" @click="handleExport">导出</el-button>
      </div>
    </div>

    <div class="content-area">
      <!-- 统计区域 -->
      <div class="stats-row">
        <div class="stat-card stat-blue">
          <span class="stat-icon">📦</span>
          <span class="stat-value">{{ stats.totalCards }}</span>
          <span class="stat-label">总库存</span>
        </div>
        <div class="stat-card stat-green">
          <span class="stat-icon">✅</span>
          <span class="stat-value">{{ stats.activeCards }}</span>
          <span class="stat-label">已发放</span>
        </div>
        <div class="stat-card stat-orange">
          <span class="stat-icon">⚠️</span>
          <span class="stat-value">{{ stats.frozenCards }}</span>
          <span class="stat-label">已冻结</span>
        </div>
        <div class="stat-card stat-red">
          <span class="stat-icon">❌</span>
          <span class="stat-value">{{ stats.expiredCards }}</span>
          <span class="stat-label">已过期</span>
        </div>
      </div>

      <!-- 筛选卡片 -->
      <div class="filter-card">
        <div class="filter-header">🔍 筛选查询</div>
        <el-divider />
        <div class="filter-row">
          <div class="filter-item">
            <span class="filter-label">卡状态：</span>
            <el-select v-model="filterStatus" placeholder="全部" style="width: 120px">
              <el-option label="全部" value="" />
              <el-option label="已发放" value="active" />
              <el-option label="已冻结" value="inactive" />
              <el-option label="已过期" value="expired" />
            </el-select>
          </div>
          <div class="filter-item">
            <span class="filter-label">卡号：</span>
            <el-input v-model="filterCardNo" placeholder="请输入卡号" style="width: 200px" />
          </div>
          <div class="filter-item">
            <span class="filter-label">持卡人：</span>
            <el-input v-model="filterHolder" placeholder="手机号或姓名" style="width: 200px" />
          </div>
          <el-button type="primary" class="search-btn" @click="handleSearch">
            查询
          </el-button>
        </div>
      </div>

      <!-- 列表卡片 -->
      <div class="list-card">
        <div class="list-header">📋 门店卡列表</div>
        <el-divider />
        <el-table
          :data="tableData"
          style="width: 100%"
          :header-cell-style="{
            backgroundColor: '#FAFAFA',
            color: '#262626',
            fontWeight: '600',
            fontSize: '14px'
          }"
        >
          <el-table-column prop="cardNo" label="卡号" width="180" />
          <el-table-column prop="holder" label="持卡人" width="200" />
          <el-table-column label="卡余额" width="120">
            <template #default="{ row }">
              ¥{{ row.balance.toLocaleString() }}
            </template>
          </el-table-column>
          <el-table-column label="状态" width="100">
            <template #default="{ row }">
              <el-tag :type="getStatusType(row.status)" size="small">
                {{ getStatusText(row.status) }}
              </el-tag>
            </template>
          </el-table-column>
          <el-table-column prop="issueDate" label="发放日期" width="120" />
          <el-table-column prop="expiryDate" label="过期日期" width="120" />
          <el-table-column label="操作" fixed="right">
            <template #default="{ row }">
              <el-button type="primary" link size="small" @click="handleViewDetail(row)">
                查看详情
              </el-button>
              <el-button
                v-if="row.status === 'active'"
                type="warning"
                link
                size="small"
                @click="handleToggleStatus(row, 'inactive')"
              >
                冻结
              </el-button>
              <el-button
                v-if="row.status === 'inactive'"
                type="success"
                link
                size="small"
                @click="handleToggleStatus(row, 'active')"
              >
                启用
              </el-button>
            </template>
          </el-table-column>
        </el-table>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { ElMessage } from 'element-plus'
import { getCardList, getCardStats, toggleCardStatus, type CardListItem } from '@/api/card'

const router = useRouter()

const stats = ref({
  totalCards: 0,
  activeCards: 0,
  frozenCards: 0,
  expiredCards: 0
})

const filterStatus = ref('')
const filterCardNo = ref('')
const filterHolder = ref('')

const tableData = ref<CardListItem[]>([])

const getStatusType = (status: string) => {
  switch (status) {
    case 'active':
      return 'success'
    case 'inactive':
      return 'warning'
    case 'expired':
      return 'danger'
    default:
      return 'info'
  }
}

const getStatusText = (status: string) => {
  switch (status) {
    case 'active':
      return '已发放'
    case 'inactive':
      return '已冻结'
    case 'expired':
      return '已过期'
    default:
      return '未知'
  }
}

const handleRefresh = () => {
  loadData()
  ElMessage.success('刷新成功')
}

const handleExport = () => {
  ElMessage.success('导出功能开发中')
}

const handleSearch = () => {
  loadData()
}

const handleViewDetail = (row: CardListItem) => {
  router.push(`/card/detail/${row.cardNo}`)
}

const handleToggleStatus = async (row: CardListItem, status: 'active' | 'inactive') => {
  await toggleCardStatus(row.cardNo, status)
  loadData()
  ElMessage.success(status === 'active' ? '启用成功' : '冻结成功')
}

const loadData = async () => {
  try {
    const [listRes, statsRes] = await Promise.all([
      getCardList({ status: filterStatus.value || undefined, holderPhone: filterHolder.value || undefined, page: 1, pageSize: 50 }),
      getCardStats()
    ])
    if (listRes?.data) {
      tableData.value = listRes.data.list || []
    }
    if (statsRes?.data) {
      stats.value = {
        totalCards: statsRes.data.totalCards || 0,
        activeCards: statsRes.data.activeCards || 0,
        frozenCards: 0,
        expiredCards: 0
      }
    }
  } catch (error) {
    // fallback to empty data
  }
}

onMounted(() => {
  loadData()
})
</script>

<style scoped>
.card-manage {
  background-color: #F5F5F5;
  min-height: calc(100vh - 64px);
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

.header-actions {
  display: flex;
  gap: 12px;
}

.action-btn {
  width: 80px;
  height: 36px;
  border-radius: 4px;
}

.content-area {
  padding: 24px;
  display: flex;
  flex-direction: column;
  gap: 24px;
}

.stats-row {
  display: grid;
  grid-template-columns: repeat(4, 1fr);
  gap: 16px;
}

.stat-card {
  background-color: #FFFFFF;
  border-radius: 8px;
  border: 1px solid;
  padding: 20px;
  display: flex;
  flex-direction: column;
  gap: 8px;
  align-items: center;
  justify-content: center;
}

.stat-blue {
  background-color: #E6F7FF;
  border-color: #91D5FF;
}

.stat-green {
  background-color: #F6FFED;
  border-color: #B7EB8F;
}

.stat-orange {
  background-color: #FFF7E6;
  border-color: #FFD591;
}

.stat-red {
  background-color: #FFF1F0;
  border-color: #FFA39E;
}

.stat-icon {
  font-size: 32px;
}

.stat-value {
  font-family: 'Inter', sans-serif;
  font-size: 32px;
  font-weight: 600;
  color: #262626;
}

.stat-green .stat-value {
  color: #52C41A;
}

.stat-orange .stat-value {
  color: #FA8C16;
}

.stat-red .stat-value {
  color: #FF4D4F;
}

.stat-label {
  font-family: 'Inter', sans-serif;
  font-size: 14px;
  color: #595959;
}

.filter-card,
.list-card {
  background-color: #FFFFFF;
  border-radius: 8px;
  border: 1px solid #E5E5E5;
  padding: 24px;
}

.filter-header,
.list-header {
  font-family: 'Inter', sans-serif;
  font-size: 16px;
  font-weight: 600;
  color: #262626;
}

.filter-row {
  display: flex;
  gap: 16px;
  align-items: center;
}

.filter-item {
  display: flex;
  gap: 8px;
  align-items: center;
}

.filter-label {
  font-family: 'Inter', sans-serif;
  font-size: 14px;
  color: #595959;
  white-space: nowrap;
}

.search-btn {
  background-color: #C00000;
  border-color: #C00000;
  border-radius: 4px;
}
</style>
