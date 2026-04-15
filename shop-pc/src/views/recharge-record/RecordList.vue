<template>
  <div class="record-list">
    <div class="page-header">
      <h1 class="page-title">充值记录</h1>
    </div>

    <div class="content-area">
      <!-- 筛选栏 -->
      <div class="filter-card">
        <el-input v-model="filterPhone" placeholder="会员手机号" style="width: 200px" clearable />
        <el-select v-model="filterCenter" placeholder="充值中心" style="width: 180px" clearable>
          <el-option v-for="c in centerList" :key="c.id" :label="c.name" :value="c.id" />
        </el-select>
        <el-date-picker
          v-model="filterDate"
          type="daterange"
          range-separator="至"
          start-placeholder="开始日期"
          end-placeholder="结束日期"
          value-format="YYYY-MM-DD"
          style="width: 260px"
        />
        <el-button type="primary" plain @click="handleSearch">查询</el-button>
        <el-button @click="handleReset">重置</el-button>
      </div>

      <!-- 列表 -->
      <div class="list-card">
        <el-table
          :data="tableData"
          style="width: 100%"
          :header-cell-style="{ backgroundColor: '#FAFAFA', color: '#262626', fontWeight: '600' }"
        >
          <el-table-column label="交易单号" width="240">
            <template #default="{ row }">{{ row.id || '-' }}</template>
          </el-table-column>
          <el-table-column prop="memberName" label="会员姓名" width="120" />
          <el-table-column prop="memberPhone" label="手机号" width="130" />
          <el-table-column prop="centerName" label="充值中心" width="160" />
          <el-table-column label="充值金额" width="120" align="right">
            <template #default="{ row }">
              <span class="amount-value">{{ row.amount != null ? `¥${Number(row.amount).toLocaleString()}` : '-' }}</span>
            </template>
          </el-table-column>
          <el-table-column label="支付方式" width="100">
            <template #default="{ row }">{{ paymentMethodMap[row.paymentMethod] || row.paymentMethod || '-' }}</template>
          </el-table-column>
          <el-table-column label="充值时间" width="180">
            <template #default="{ row }">{{ formatTime(row.createdAt) }}</template>
          </el-table-column>
          <el-table-column label="操作" width="100" fixed="right">
            <template #default="{ row }">
              <el-button type="primary" link size="small" @click="handleViewDetail(row)">详情</el-button>
            </template>
          </el-table-column>
        </el-table>

        <div class="pagination-row">
          <el-pagination
            v-model:current-page="currentPage"
            v-model:page-size="pageSize"
            :total="total"
            :page-sizes="[10, 20, 50]"
            layout="total, sizes, prev, pager, next"
            @size-change="loadData"
            @current-change="loadData"
          />
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { getRechargeRecordList, getCenterList } from '@/api/recharge'
import type { CenterItem } from '@/api/recharge'

const router = useRouter()
const filterPhone = ref('')
const filterCenter = ref('')
const filterDate = ref<[string, string] | null>(null)

const tableData = ref<any[]>([])
const centerList = ref<CenterItem[]>([])
const currentPage = ref(1)
const pageSize = ref(10)
const total = ref(0)

const paymentMethodMap: Record<string, string> = {
  cash: '现金',
  wechat: '微信',
  alipay: '支付宝',
  card: '门店卡',
}

const formatTime = (t: string) => {
  if (!t) return '-'
  return t.replace('T', ' ').slice(0, 19)
}

const handleSearch = () => {
  currentPage.value = 1
  loadData()
}

const handleReset = () => {
  filterPhone.value = ''
  filterCenter.value = ''
  filterDate.value = null
  currentPage.value = 1
  loadData()
}

const handleViewDetail = (row: any) => {
  router.push(`/recharge/records/${row.id}`)
}

const loadData = async () => {
  try {
    const params: any = {
      page: currentPage.value,
      pageSize: pageSize.value,
    }
    if (filterPhone.value) params.memberPhone = filterPhone.value
    if (filterCenter.value) params.centerId = filterCenter.value
    if (filterDate.value && filterDate.value[0]) {
      params.startDate = filterDate.value[0]
      params.endDate = filterDate.value[1]
    }

    const res = await getRechargeRecordList(params)
    if (res?.data) {
      tableData.value = res.data.list || []
      total.value = res.data.total || 0
    }
  } catch {
    // fallback to empty
  }
}

const loadCenters = async () => {
  try {
    const res = await getCenterList()
    if (res?.data) {
      centerList.value = res.data
    }
  } catch {
    // ignore
  }
}

onMounted(() => {
  loadCenters()
  loadData()
})
</script>

<style scoped>
.record-list {
  background-color: #F5F5F5;
  min-height: calc(100vh - 64px);
}

.page-header {
  height: 64px;
  background-color: #FFFFFF;
  border-bottom: 1px solid #E5E5E5;
  padding: 16px 24px;
  display: flex;
  justify-content: space-between;
  align-items: center;
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
  flex-direction: column;
  gap: 16px;
}

.filter-card {
  background-color: #FFFFFF;
  border-radius: 8px;
  border: 1px solid #E5E5E5;
  padding: 16px;
  display: flex;
  gap: 12px;
  align-items: center;
}

.list-card {
  background-color: #FFFFFF;
  border-radius: 8px;
  border: 1px solid #E5E5E5;
  padding: 24px;
}

.amount-value {
  font-weight: 600;
  color: #C00000;
}

.pagination-row {
  display: flex;
  justify-content: center;
  margin-top: 16px;
}
</style>
