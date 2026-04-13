<template>
  <div class="record-list">
    <div class="page-header">
      <h1 class="page-title">充值记录</h1>
    </div>

    <div class="content-area">
      <!-- 筛选卡片 -->
      <div class="filter-card">
        <div class="filter-header">🔍 筛选查询</div>
        <el-divider />
        <div class="filter-row">
          <div class="filter-item">
            <span class="filter-label">会员手机：</span>
            <el-input v-model="filterPhone" placeholder="手机号" style="width: 200px" />
          </div>
          <div class="filter-item">
            <span class="filter-label">充值中心：</span>
            <el-select v-model="filterCenter" placeholder="全部" style="width: 150px">
              <el-option label="全部" value="" />
              <el-option label="北京朝阳中心" value="1" />
              <el-option label="北京海淀中心" value="2" />
            </el-select>
          </div>
          <div class="filter-item">
            <span class="filter-label">日期：</span>
            <el-date-picker v-model="filterDate" type="daterange" style="width: 240px" />
          </div>
          <el-button type="primary" class="search-btn" @click="handleSearch">查询</el-button>
        </div>
      </div>

      <!-- 列表卡片 -->
      <div class="list-card">
        <el-table
          :data="tableData"
          style="width: 100%"
          :header-cell-style="{ backgroundColor: '#FAFAFA', color: '#262626', fontWeight: '600' }"
        >
          <el-table-column label="交易单号" width="220">
            <template #default="{ row }">
              {{ row.id || row.transactionNo }}
            </template>
          </el-table-column>
          <el-table-column prop="memberName" label="会员姓名" width="120" />
          <el-table-column prop="memberPhone" label="手机号" width="120" />
          <el-table-column prop="centerName" label="充值中心" width="150" />
          <el-table-column label="充值金额" width="120">
            <template #default="{ row }">¥{{ row.amount.toLocaleString() }}</template>
          </el-table-column>
          <el-table-column prop="paymentMethod" label="支付方式" width="100" />
          <el-table-column prop="createdAt" label="充值时间" width="180" />
          <el-table-column label="操作" width="100" fixed="right">
            <template #default="{ row }">
              <el-button type="primary" link size="small" @click="handleViewDetail(row)">
                查看详情
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
import { getRechargeRecordList } from '@/api/recharge'

const router = useRouter()
const filterPhone = ref('')
const filterCenter = ref('')
const filterDate = ref<[Date, Date] | null>(null)

const tableData = ref<any[]>([])

const handleSearch = () => {
  loadData()
}

const handleViewDetail = (row: any) => {
  router.push(`/recharge/records/${row.id}`)
}

const loadData = async () => {
  try {
    const res = await getRechargeRecordList({
      memberPhone: filterPhone.value || undefined,
      centerId: filterCenter.value || undefined,
      page: 1,
      pageSize: 50
    })
    if (res?.data) {
      tableData.value = res.data.list || []
    }
  } catch (error) {
    // fallback to empty
  }
}

onMounted(() => {
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
  gap: 24px;
}

.filter-card,
.list-card {
  background-color: #FFFFFF;
  border-radius: 8px;
  border: 1px solid #E5E5E5;
  padding: 24px;
}

.filter-header {
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
  font-size: 14px;
  color: #595959;
}

.search-btn {
  background-color: #C00000;
  border-color: #C00000;
  border-radius: 4px;
}
</style>
