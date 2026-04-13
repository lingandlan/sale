<template>
  <div class="brecharge-list">
    <div class="page-header">
      <h1 class="page-title">B端充值审批列表</h1>
    </div>

    <div class="tabs-container">
      <div
        class="tab-item"
        :class="{ active: activeTab === 'pending' }"
        @click="handleTabChange('pending')"
      >
        待审批
      </div>
      <div
        class="tab-item"
        :class="{ active: activeTab === 'approved' }"
        @click="handleTabChange('approved')"
      >
        已审批
      </div>
    </div>

    <div class="filter-bar">
      <div class="filter-item">
        <span class="filter-label">申请日期</span>
        <el-date-picker
          v-model="filterDate"
          type="daterange"
          range-separator="至"
          start-placeholder="开始日期"
          end-placeholder="结束日期"
          size="default"
          style="width: 240px"
        />
      </div>
      <div class="filter-item">
        <span class="filter-label">充值中心</span>
        <el-select
          v-model="filterCenter"
          placeholder="全部充值中心"
          style="width: 200px"
        >
          <el-option label="全部充值中心" value="" />
          <el-option
            v-for="center in centers"
            :key="center.id"
            :label="center.name"
            :value="center.id"
          />
        </el-select>
      </div>
      <el-button type="primary" class="search-btn" @click="handleSearch">
        查询
      </el-button>
    </div>

    <div class="table-card">
      <el-table
        ref="tableRef"
        :data="tableData"
        style="width: 100%"
        :header-cell-style="{
          backgroundColor: '#FAFAFA',
          color: '#262626',
          fontWeight: '600',
          fontSize: '14px'
        }"
        @selection-change="handleSelectionChange"
      >
        <el-table-column type="selection" width="60" align="center" />
        <el-table-column prop="centerName" label="充值中心" min-width="140" />
        <el-table-column label="充值金额" min-width="120" align="right">
          <template #default="{ row }">
            ¥{{ (row.amount ?? 0).toLocaleString() }}
          </template>
        </el-table-column>
        <el-table-column label="预计积分" min-width="120" align="right">
          <template #default="{ row }">
            {{ (row.points ?? 0).toLocaleString() }}
          </template>
        </el-table-column>
        <el-table-column prop="applicantName" label="申请人" width="100" align="center" />
        <el-table-column label="申请时间" min-width="170" align="center">
          <template #default="{ row }">
            {{ formatTime(row.createdAt) }}
          </template>
        </el-table-column>
        <el-table-column label="状态" width="90" align="center">
          <template #default="{ row }">
            <el-tag :type="getStatusType(row.status)" size="small">
              {{ getStatusText(row.status) }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column label="操作" width="160" align="center">
          <template #default="{ row }">
            <div class="action-buttons">
              <el-button
                v-if="row.status === 'pending'"
                type="success"
                size="small"
                @click="handleApprove(row)"
              >
                通过
              </el-button>
              <el-button
                v-if="row.status === 'pending'"
                type="danger"
                size="small"
                @click="handleReject(row)"
              >
                拒绝
              </el-button>
              <el-button
                v-if="row.status !== 'pending'"
                type="primary"
                size="small"
                link
                @click="handleViewDetail(row)"
              >
                查看详情
              </el-button>
            </div>
          </template>
        </el-table-column>
      </el-table>
    </div>

    <div v-if="selectedRows.length > 0" class="batch-operation-bar">
      <span class="selected-info">已选择 {{ selectedRows.length }} 项</span>
      <el-button type="success" @click="handleBatchApprove">批量通过</el-button>
      <el-button type="danger" @click="handleBatchReject">批量拒绝</el-button>
    </div>

    <div class="pagination">
      <span class="page-info">共 {{ total }} 条，每页 {{ pageSize }} 条</span>
      <el-button :disabled="currentPage === 1" @click="handlePageChange(currentPage - 1)">
        上一页
      </el-button>
      <span class="page-indicator">{{ currentPage }} / {{ totalPages }}</span>
      <el-button :disabled="currentPage === totalPages" @click="handlePageChange(currentPage + 1)">
        下一页
      </el-button>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { ElMessage, ElMessageBox } from 'element-plus'
import { getBRechargeApprovalList, approvalAction } from '@/api/recharge'
import type { BRechargeApprovalItem } from '@/api/recharge'

const router = useRouter()

interface Center {
  id: string
  name: string
}

const centers = ref<Center[]>([
  { id: '1', name: '北京朝阳中心' },
  { id: '2', name: '北京海淀中心' },
  { id: '3', name: '上海浦东中心' }
])

const activeTab = ref<'pending' | 'approved'>('pending')
const filterDate = ref<[Date, Date] | null>(null)
const filterCenter = ref('')
const selectedRows = ref<BRechargeApprovalItem[]>([])

const tableData = ref<BRechargeApprovalItem[]>([])

const currentPage = ref(1)
const pageSize = ref(10)
const total = ref(0)

const totalPages = computed(() => Math.ceil(total.value / pageSize.value))

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

const formatTime = (iso: string) => {
  if (!iso) return ''
  const d = new Date(iso)
  const pad = (n: number) => String(n).padStart(2, '0')
  return `${d.getFullYear()}-${pad(d.getMonth() + 1)}-${pad(d.getDate())} ${pad(d.getHours())}:${pad(d.getMinutes())}`
}

const loadData = async () => {
  try {
    const res = await getBRechargeApprovalList({
      status: activeTab.value === 'pending' ? 'pending' : 'approved,rejected',
      page: currentPage.value,
      pageSize: pageSize.value
    })
    if (res?.data) {
      tableData.value = res.data.list || []
      total.value = res.data.total || 0
    }
  } catch (error) {
    // fallback to empty data
  }
}

const handleTabChange = (tab: 'pending' | 'approved') => {
  activeTab.value = tab
  currentPage.value = 1
  loadData()
}

const handleSearch = () => {
  currentPage.value = 1
  loadData()
}

const handleSelectionChange = (selection: BRechargeApprovalItem[]) => {
  selectedRows.value = selection
}

const handleApprove = async (row: BRechargeApprovalItem) => {
  try {
    await ElMessageBox.confirm('确认通过该充值申请？', '审批确认', {
      confirmButtonText: '确认',
      cancelButtonText: '取消',
      type: 'success'
    })
    await approvalAction({ id: row.id, action: 'approve' })
    loadData()
    ElMessage.success('已通过')
  } catch {
    // 用户取消
  }
}

const handleReject = async (row: BRechargeApprovalItem) => {
  try {
    const { value } = await ElMessageBox.prompt('请输入拒绝原因', '拒绝申请', {
      confirmButtonText: '确认',
      cancelButtonText: '取消',
      inputPattern: /.+/,
      inputErrorMessage: '请输入拒绝原因'
    })
    await approvalAction({ id: row.id, action: 'reject', remark: value })
    loadData()
    ElMessage.success('已拒绝')
  } catch {
    // 用户取消
  }
}

const handleViewDetail = (row: BRechargeApprovalItem) => {
  router.push(`/recharge/b-approval/${row.id}`)
}

const handleBatchApprove = async () => {
  try {
    await ElMessageBox.confirm(
      `确认批量通过 ${selectedRows.value.length} 条充值申请？`,
      '批量审批确认',
      {
        confirmButtonText: '确认',
        cancelButtonText: '取消',
        type: 'success'
      }
    )
    for (const row of selectedRows.value) {
      await approvalAction({ id: row.id, action: 'approve' })
    }
    selectedRows.value = []
    loadData()
    ElMessage.success('批量通过成功')
  } catch {
    // 用户取消
  }
}

const handleBatchReject = async () => {
  try {
    const { value } = await ElMessageBox.prompt('请输入拒绝原因', '批量拒绝申请', {
      confirmButtonText: '确认',
      cancelButtonText: '取消',
      inputPattern: /.+/,
      inputErrorMessage: '请输入拒绝原因'
    })
    for (const row of selectedRows.value) {
      await approvalAction({ id: row.id, action: 'reject', remark: value })
    }
    selectedRows.value = []
    loadData()
    ElMessage.success('批量拒绝成功')
  } catch {
    // 用户取消
  }
}

const handlePageChange = (page: number) => {
  currentPage.value = page
  loadData()
}

onMounted(() => {
  loadData()
})
</script>

<style scoped>
.brecharge-list {
  padding: 24px;
  background-color: #F5F5F5;
  min-height: calc(100vh - 64px);
  display: flex;
  flex-direction: column;
  gap: 16px;
}

.page-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.page-title {
  font-family: 'Inter', sans-serif;
  font-size: 20px;
  font-weight: 600;
  color: #C00000;
  margin: 0;
}

.tabs-container {
  display: flex;
  gap: 4px;
  padding: 4px;
  background-color: #FFFFFF;
  border-radius: 4px;
  width: fit-content;
}

.tab-item {
  width: 100px;
  height: 40px;
  display: flex;
  align-items: center;
  justify-content: center;
  font-family: 'Inter', sans-serif;
  font-size: 14px;
  color: #262626;
  border-radius: 4px;
  cursor: pointer;
  transition: all 0.3s;
}

.tab-item.active {
  background-color: #C00000;
  color: #FFFFFF;
}

.tab-item:hover:not(.active) {
  background-color: #F5F5F5;
}

.filter-bar {
  display: flex;
  gap: 12px;
  align-items: center;
  background-color: #FFFFFF;
  border-radius: 8px;
  padding: 16px;
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
  height: 40px;
}

.search-btn:hover {
  background-color: #A00000;
  border-color: #A00000;
}

.table-card {
  background-color: #FFFFFF;
  border-radius: 8px;
  overflow: hidden;
}

.action-buttons {
  display: flex;
  gap: 8px;
  justify-content: center;
}

.batch-operation-bar {
  display: flex;
  gap: 12px;
  align-items: center;
  background-color: #FFFFFF;
  border-radius: 8px;
  padding: 16px;
  height: 56px;
}

.selected-info {
  font-family: 'Inter', sans-serif;
  font-size: 14px;
  color: #595959;
}

.pagination {
  display: flex;
  gap: 16px;
  align-items: center;
  justify-content: center;
  background-color: #FFFFFF;
  border-radius: 8px;
  padding: 16px;
  height: 48px;
}

.page-info,
.page-indicator {
  font-family: 'Inter', sans-serif;
  font-size: 14px;
  color: #595959;
}

.page-indicator {
  color: #262626;
}
</style>
