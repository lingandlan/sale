<template>
  <div class="card-manage">
    <div class="page-header">
      <h1 class="page-title">门店卡管理</h1>
      <div class="header-actions">
        <el-button class="action-btn" @click="handleRefresh">刷新</el-button>
      </div>
    </div>

    <div class="content-area">
      <!-- 统计区域 -->
      <div class="stats-row">
        <div class="stat-card stat-blue">
          <span class="stat-value">{{ stats.totalCards }}</span>
          <span class="stat-label">总卡数</span>
        </div>
        <div class="stat-card" style="background:#F0F0F0;border-color:#D9D9D9">
          <span class="stat-value">{{ stats.inStockCards }}</span>
          <span class="stat-label">已入库</span>
        </div>
        <div class="stat-card" style="background:#E6F7FF;border-color:#91D5FF">
          <span class="stat-value">{{ stats.issuedCards }}</span>
          <span class="stat-label">已发放</span>
        </div>
        <div class="stat-card stat-green">
          <span class="stat-value">{{ stats.activeCards }}</span>
          <span class="stat-label">已激活</span>
        </div>
        <div class="stat-card stat-orange">
          <span class="stat-value">{{ stats.frozenCards }}</span>
          <span class="stat-label">已冻结</span>
        </div>
        <div class="stat-card stat-red">
          <span class="stat-value">{{ stats.expiredCards }}</span>
          <span class="stat-label">已过期</span>
        </div>
      </div>

      <!-- 筛选 -->
      <div class="filter-card">
        <div class="filter-row">
          <div class="filter-item">
            <span class="filter-label">状态：</span>
            <el-select v-model="filterStatus" placeholder="全部" clearable style="width: 120px">
              <el-option v-for="(label, key) in CardStatusMap" :key="key" :label="label" :value="Number(key)" />
            </el-select>
          </div>
          <div class="filter-item">
            <span class="filter-label">充值中心：</span>
            <el-select v-model="filterCenterId" placeholder="全部" clearable style="width: 180px" :disabled="!userStore.canSelectAllCenters">
              <el-option v-for="c in centerOptions" :key="c.id" :label="c.name" :value="c.id" />
            </el-select>
          </div>
          <div class="filter-item">
            <span class="filter-label">卡号：</span>
            <el-input v-model="filterCardNo" placeholder="输入卡号" style="width: 200px" />
          </div>
          <el-button type="primary" class="search-btn" @click="loadData">查询</el-button>
        </div>
      </div>

      <!-- 列表 -->
      <div class="list-card">
        <el-table :data="tableData" style="width: 100%">
          <el-table-column prop="cardNo" label="卡号" width="160" />
          <el-table-column prop="rechargeCenterName" label="充值中心" min-width="140" />
          <el-table-column prop="cardType" label="类型" width="80">
            <template #default="{ row }">{{ CardTypeMap[row.cardType] || '-' }}</template>
          </el-table-column>
          <el-table-column label="余额" width="100">
            <template #default="{ row }">¥{{ row.balance }}</template>
          </el-table-column>
          <el-table-column label="状态" width="90">
            <template #default="{ row }">
              <el-tag :type="CardStatusTagType[row.status] || 'info'" size="small">
                {{ CardStatusMap[row.status] || '未知' }}
              </el-tag>
            </template>
          </el-table-column>
          <el-table-column prop="batchNo" label="批次号" width="160" />
          <el-table-column label="发放日期" width="110">
            <template #default="{ row }">{{ row.issuedAt ? row.issuedAt.slice(0, 10) : '-' }}</template>
          </el-table-column>
          <el-table-column label="激活日期" width="110">
            <template #default="{ row }">{{ row.activatedAt ? row.activatedAt.slice(0, 10) : '-' }}</template>
          </el-table-column>
          <el-table-column label="过期日期" width="110">
            <template #default="{ row }">{{ row.expiredAt ? row.expiredAt.slice(0, 10) : '-' }}</template>
          </el-table-column>
          <el-table-column label="操作" fixed="right" width="200">
            <template #default="{ row }">
              <el-button type="primary" link size="small" @click="router.push(`/card/detail/${row.cardNo}`)">详情</el-button>
              <el-button v-if="row.status === 2 || row.status === 3" type="warning" link size="small" @click="handleFreeze(row)">冻结</el-button>
              <el-button v-if="row.status === 4" type="success" link size="small" @click="handleUnfreeze(row)">解冻</el-button>
              <el-button v-if="row.status !== 5 && row.status !== 6" type="danger" link size="small" @click="handleVoid(row)">作废</el-button>
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
import { ElMessage, ElMessageBox } from 'element-plus'
import { extractErrorMessage } from '@/utils/request'
import { getCardList, getCardStats, freezeCard, unfreezeCard, voidCard, CardStatusMap, CardStatusTagType, CardTypeMap, type CardListItem } from '@/api/card'
import { getCenterList } from '@/api/center'
import { useUserStore } from '@/stores/user'

const userStore = useUserStore()
const router = useRouter()

const stats = ref({
  totalCards: 0, inStockCards: 0, issuedCards: 0,
  activeCards: 0, frozenCards: 0, expiredCards: 0, voidedCards: 0
})

const filterStatus = ref<number | undefined>(undefined)
const filterCenterId = ref('')
const filterCardNo = ref('')
const centerOptions = ref<{ id: string; name: string }[]>([])
const tableData = ref<CardListItem[]>([])

const handleRefresh = () => { loadData(); ElMessage.success('已刷新') }

const handleFreeze = async (row: CardListItem) => {
  await ElMessageBox.confirm(`确认冻结卡 ${row.cardNo}？`, '冻结')
  await freezeCard(row.cardNo)
  ElMessage.success('冻结成功')
  loadData()
}

const handleUnfreeze = async (row: CardListItem) => {
  await ElMessageBox.confirm(`确认解冻卡 ${row.cardNo}？`, '解冻')
  await unfreezeCard(row.cardNo)
  ElMessage.success('解冻成功')
  loadData()
}

const handleVoid = async (row: CardListItem) => {
  await ElMessageBox.confirm(`确认作废卡 ${row.cardNo}？此操作不可恢复`, '作废', { type: 'warning' })
  await voidCard(row.cardNo)
  ElMessage.success('作废成功')
  loadData()
}

const loadData = async () => {
  try {
    const [listRes, statsRes] = await Promise.all([
      getCardList({ status: filterStatus.value, cardNo: filterCardNo.value || undefined, centerId: filterCenterId.value || undefined, page: 1, pageSize: 50 }),
      getCardStats()
    ])
    const listData = listRes?.data || listRes
    tableData.value = listData?.list || []
    const s = statsRes?.data || statsRes
    if (s) {
      stats.value = {
        totalCards: s.totalCards || 0, inStockCards: s.inStockCards || 0,
        issuedCards: s.issuedCards || 0, activeCards: s.activeCards || 0,
        frozenCards: s.frozenCards || 0, expiredCards: s.expiredCards || 0,
        voidedCards: s.voidedCards || 0
      }
    }
  } catch (err: any) {
    ElMessage.error(extractErrorMessage(err, '加载数据失败'))
  }
}

const loadCenterOptions = async () => {
  if (userStore.canSelectAllCenters) {
    try {
      const res = await getCenterList()
      centerOptions.value = (res.data || []).map((c: any) => ({ id: c.id, name: c.name }))
    } catch {
      centerOptions.value = []
    }
  } else {
    const cid = userStore.userCenterId
    const cname = userStore.userCenterName
    if (cid) {
      centerOptions.value = [{ id: String(cid), name: cname || '' }]
      filterCenterId.value = String(cid)
    }
  }
}

onMounted(() => {
  if (!userStore.userInfo) {
    userStore.fetchUserInfo().then(() => loadCenterOptions())
  } else {
    loadCenterOptions()
  }
  loadData()
})
</script>

<style scoped>
.card-manage { background: #F5F5F5; min-height: calc(100vh - 64px); }
.page-header { display: flex; justify-content: space-between; align-items: center; height: 64px; background: #FFF; border-bottom: 1px solid #E5E5E5; padding: 16px 24px; }
.page-title { font-size: 20px; font-weight: 600; color: #262626; margin: 0; }
.header-actions { display: flex; gap: 12px; }
.action-btn { width: 80px; height: 36px; }
.content-area { padding: 24px; display: flex; flex-direction: column; gap: 20px; }
.stats-row { display: grid; grid-template-columns: repeat(6, 1fr); gap: 12px; }
.stat-card { background: #FFF; border-radius: 8px; border: 1px solid; padding: 16px; display: flex; flex-direction: column; gap: 4px; align-items: center; }
.stat-value { font-size: 28px; font-weight: 600; color: #262626; }
.stat-label { font-size: 13px; color: #595959; }
.stat-blue { background: #E6F7FF; border-color: #91D5FF; }
.stat-green { background: #F6FFED; border-color: #B7EB8F; }
.stat-green .stat-value { color: #52C41A; }
.stat-orange { background: #FFF7E6; border-color: #FFD591; }
.stat-orange .stat-value { color: #FA8C16; }
.stat-red { background: #FFF1F0; border-color: #FFA39E; }
.stat-red .stat-value { color: #FF4D4F; }
.filter-card, .list-card { background: #FFF; border-radius: 8px; border: 1px solid #E5E5E5; padding: 20px; }
.filter-row { display: flex; gap: 16px; align-items: center; }
.filter-item { display: flex; gap: 8px; align-items: center; }
.filter-label { font-size: 14px; color: #595959; white-space: nowrap; }
.search-btn { background: #C00000; border-color: #C00000; }
</style>
