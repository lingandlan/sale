<template>
  <div class="center-manage">
    <div class="page-header">
      <h1 class="page-title">充值中心管理</h1>
      <el-button type="primary" class="add-btn" @click="handleAdd">
        <span style="margin-right: 4px">+</span> 新建充值中心
      </el-button>
    </div>

    <div class="content-area">
      <!-- 筛选栏 -->
      <div class="filter-card">
        <el-input v-model="filters.keyword" placeholder="搜索中心名称" style="width: 240px" clearable />
        <el-select v-model="filters.level" placeholder="级别" style="width: 180px" clearable @change="handleSearch">
          <el-option label="子公司合伙人" value="1" />
          <el-option label="服务中心合伙人" value="2" />
        </el-select>
        <el-select v-model="filters.status" placeholder="状态" style="width: 120px" clearable @change="handleSearch">
          <el-option label="正常" value="normal" />
          <el-option label="冻结" value="frozen" />
        </el-select>
        <el-button type="primary" @click="handleSearch">查询</el-button>
        <el-button @click="handleResetFilter">重置</el-button>
      </div>

      <!-- 表格 -->
      <div class="list-card">
        <el-table
          :data="tableData"
          style="width: 100%"
        >
          <el-table-column prop="name" label="中心名称" width="200" />
          <el-table-column prop="region" label="省/市/区" width="200" />
          <el-table-column prop="level" label="级别" width="150">
            <template #default="{ row }">
              <el-tag :type="row.level === '子公司合伙人' ? 'primary' : 'info'" size="small">{{ row.level }}</el-tag>
            </template>
          </el-table-column>
          <el-table-column label="积分余额" width="140" align="right">
            <template #default="{ row }">
              <span class="balance-value">{{ row.balance.toLocaleString() }}</span>
            </template>
          </el-table-column>
          <el-table-column prop="totalIn" label="累计充值" width="140" align="right">
            <template #default="{ row }">{{ row.totalIn.toLocaleString() }}</template>
          </el-table-column>
          <el-table-column prop="totalOut" label="已消耗" width="140" align="right">
            <template #default="{ row }">{{ row.totalOut.toLocaleString() }}</template>
          </el-table-column>
          <el-table-column label="状态" width="100">
            <template #default="{ row }">
              <el-tag :type="row.status === 'normal' ? 'success' : 'danger'" size="small">
                {{ row.status === 'normal' ? '正常' : '冻结' }}
              </el-tag>
            </template>
          </el-table-column>
          <el-table-column label="操作" width="200" fixed="right">
            <template #default="{ row }">
              <el-button type="primary" link size="small" @click="handleEdit(row)">编辑</el-button>
              <el-button
                :type="row.status === 'normal' ? 'danger' : 'success'"
                link size="small"
                @click="handleToggleFreeze(row)"
              >
                {{ row.status === 'normal' ? '冻结' : '解冻' }}
              </el-button>
              <el-button type="primary" link size="small" @click="handleDetail(row)">详情</el-button>
            </template>
          </el-table-column>
        </el-table>

        <div class="pagination-row">
          <el-pagination
            v-model:current-page="pagination.page"
            v-model:page-size="pagination.size"
            :total="pagination.total"
            :page-sizes="[10, 20, 50]"
            layout="total, sizes, prev, pager, next"
          />
        </div>
      </div>
    </div>

    <!-- 新建/编辑弹窗 -->
    <el-dialog v-model="dialogVisible" :title="dialogTitle" width="600px" :close-on-click-modal="false">
      <el-form :model="formData" label-width="120px">
        <el-form-item label="中心名称" required>
          <el-input v-model="formData.name" placeholder="请输入中心名称" />
        </el-form-item>
        <el-form-item label="级别" required>
          <el-select v-model="formData.level" placeholder="请选择级别" style="width: 100%">
            <el-option label="子公司合伙人" value="子公司合伙人" />
            <el-option label="服务中心合伙人" value="服务中心合伙人" />
          </el-select>
        </el-form-item>
        <el-form-item label="省份" required>
          <el-select
            v-model="formData.province"
            placeholder="请选择省份"
            style="width: 100%"
            @change="handleProvinceChange"
          >
            <el-option v-for="p in regionData" :key="p.value" :label="p.label" :value="p.value" />
          </el-select>
        </el-form-item>
        <el-form-item label="城市" required>
          <el-select
            v-model="formData.city"
            placeholder="请选择城市"
            style="width: 100%"
            :disabled="!formData.province"
            @change="handleCityChange"
          >
            <el-option v-for="c in cityOptions" :key="c.value" :label="c.label" :value="c.value" />
          </el-select>
        </el-form-item>
        <el-form-item label="区县" required>
          <el-select
            v-model="formData.district"
            placeholder="请选择区县"
            style="width: 100%"
            :disabled="!formData.city"
          >
            <el-option v-for="d in districtOptions" :key="d.value" :label="d.label" :value="d.value" />
          </el-select>
        </el-form-item>
        <el-form-item label="具体位置" required>
          <el-input v-model="formData.address" placeholder="请输入具体位置" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="dialogVisible = false">取消</el-button>
        <el-button type="primary" class="save-btn" @click="handleSaveCenter">保存</el-button>
      </template>
    </el-dialog>

    <!-- 详情弹窗 -->
    <el-dialog v-model="detailVisible" title="充值中心详情" width="600px">
      <el-descriptions :column="2" border v-if="detailData">
        <el-descriptions-item label="中心名称">{{ detailData.name }}</el-descriptions-item>
        <el-descriptions-item label="级别">
          <el-tag :type="detailData.code?.includes('服务') ? 'info' : 'primary'" size="small">
            {{ detailData.code?.includes('服务') ? '服务中心合伙人' : '子公司合伙人' }}
          </el-tag>
        </el-descriptions-item>
        <el-descriptions-item label="省/市/区">{{ [detailData.province, detailData.city, detailData.district].filter(Boolean).join(' / ') || '-' }}</el-descriptions-item>
        <el-descriptions-item label="具体位置">{{ detailData.address || '-' }}</el-descriptions-item>
        <el-descriptions-item label="联系电话">{{ detailData.phone || '-' }}</el-descriptions-item>
        <el-descriptions-item label="积分余额">
          <span class="balance-value">{{ (detailData.balance ?? 0).toLocaleString() }}</span>
        </el-descriptions-item>
        <el-descriptions-item label="状态">
          <el-tag :type="detailData.status === 'active' ? 'success' : 'danger'" size="small">
            {{ detailData.status === 'active' ? '正常' : '冻结' }}
          </el-tag>
        </el-descriptions-item>
      </el-descriptions>
      <template #footer>
        <el-button @click="detailVisible = false">关闭</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, computed, onMounted } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { extractErrorMessage } from '@/utils/request'
import { getCenterList, getCenterDetail, createCenter, updateCenter, type CenterItem } from '@/api/center'
import { regionData, getCities, getDistricts } from '@/utils/regionData'

const filters = reactive({ keyword: '', level: '', status: '' })
const pagination = reactive({ page: 1, size: 10, total: 0 })

const dialogVisible = ref(false)
const dialogTitle = ref('新建充值中心')

const formData = reactive({
  id: '' as string,
  name: '',
  level: '',
  province: '',
  city: '',
  district: '',
  address: ''
})

const tableData = ref<any[]>([])
const detailVisible = ref(false)
const detailData = ref<CenterItem | null>(null)
// 级联地区选项
const cityOptions = computed(() => getCities(formData.province))
const districtOptions = computed(() => getDistricts(formData.province, formData.city))

const handleProvinceChange = () => {
  formData.city = ''
  formData.district = ''
}

const handleCityChange = () => {
  formData.district = ''
}

const loadData = async () => {
  try {
    const res = await getCenterList()
    if (res?.data) {
      let list: CenterItem[] = res.data || []

      // 关键字过滤
      const keyword = (filters.keyword || '').trim()
      if (keyword) {
        list = list.filter((c: CenterItem) => c.name.includes(keyword))
      }

      // 级别过滤：基于显示文本匹配
      const levelVal = filters.level || ''
      if (levelVal) {
        list = list.filter((c: CenterItem) => {
          const displayLevel = c.code?.includes('服务') ? '服务中心合伙人' : '子公司合伙人'
          const targetLevel = levelVal === '1' ? '子公司合伙人' : '服务中心合伙人'
          return displayLevel === targetLevel
        })
      }

      // 状态过滤：API 用 active/frozen，前端用 normal/frozen
      const statusVal = filters.status || ''
      if (statusVal) {
        const apiStatus = statusVal === 'normal' ? 'active' : statusVal
        list = list.filter((c: CenterItem) => c.status === apiStatus)
      }

      tableData.value = list.map((c: CenterItem) => {
        return {
          id: c.id,
          name: c.name,
          region: [c.province, c.city, c.district].filter(Boolean).join('/') || '-',
          level: c.code?.includes('服务') ? '服务中心合伙人' : '子公司合伙人',
          province: c.province,
          city: c.city,
          district: c.district,
          address: c.address,
          balance: c.balance ?? 0,
          totalIn: c.totalRecharge ?? 0,
          totalOut: c.totalConsumed ?? 0,
          status: c.status === 'active' ? 'normal' : 'frozen'
        }
      })
      pagination.total = tableData.value.length
    }
  } catch (err: any) {
    ElMessage.error(extractErrorMessage(err, '加载充值中心列表失败'))
  }
}

const handleSearch = () => { loadData() }
const handleResetFilter = () => { filters.keyword = ''; filters.level = ''; filters.status = ''; loadData() }

const handleAdd = () => {
  dialogTitle.value = '新建充值中心'
  Object.assign(formData, { id: '', name: '', level: '', province: '', city: '', district: '', address: '' })
  dialogVisible.value = true
}

const handleEdit = (row: any) => {
  dialogTitle.value = '编辑充值中心'
  Object.assign(formData, {
    id: row.id,
    name: row.name,
    level: row.level,
    province: row.province || '',
    city: row.city || '',
    district: row.district || '',
    address: row.address || '',
  })
  dialogVisible.value = true
}

const handleSaveCenter = async () => {
  if (!formData.name || !formData.province || !formData.city || !formData.district) {
    ElMessage.warning('请填写完整信息')
    return
  }
  try {
    const payload: any = {
      name: formData.name,
      code: formData.level || '子公司合伙人',
      province: formData.province,
      city: formData.city,
      district: formData.district,
      address: formData.address,
    }
    if (formData.id) {
      await updateCenter(formData.id, { ...payload, status: 'active' })
    } else {
      await createCenter(payload)
    }
    ElMessage.success(formData.id ? '编辑成功' : '创建成功')
    dialogVisible.value = false
    loadData()
  } catch (err: any) {
    ElMessage.error(extractErrorMessage(err, formData.id ? '编辑充值中心失败' : '创建充值中心失败'))
  }
}

const handleToggleFreeze = async (row: any) => {
  const action = row.status === 'normal' ? '冻结' : '解冻'
  try {
    await ElMessageBox.confirm(`确认${action}充值中心 "${row.name}"？`, `确认${action}`, {
      confirmButtonText: '确认',
      cancelButtonText: '取消',
      type: 'warning'
    })
    const newStatus = row.status === 'normal' ? 'frozen' : 'active'
    await updateCenter(row.id, { status: newStatus })
    loadData()
    ElMessage.success(`已${action}`)
  } catch (err: any) {
    if (err === 'cancel' || err?.toString?.().includes('cancel')) return
    ElMessage.error(extractErrorMessage(err, `${action}失败`))
  }
}

const handleDetail = async (row: any) => {
  try {
    const res = await getCenterDetail(row.id)
    if (res?.data) {
      detailData.value = res.data
      detailVisible.value = true
    }
  } catch (err: any) {
    ElMessage.error(extractErrorMessage(err, '加载详情失败'))
  }
}

onMounted(() => {
  loadData()
})
</script>

<style scoped>
.center-manage {
  background-color: var(--color-bg);
  min-height: calc(100vh - 64px);
}

.page-header {
  height: 64px;
  background-color: var(--color-bg-card);
  border-bottom: 1px solid var(--color-border);
  padding: 16px 24px;
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.page-title {
  font-size: 20px;
  font-weight: 600;
  color: var(--color-text-primary);
  margin: 0;
}

.add-btn {
  background-color: var(--color-primary);
  border-color: var(--color-primary);
  border-radius: var(--radius-sm);
}

.add-btn:hover {
  background-color: var(--color-primary-hover);
  border-color: var(--color-primary-hover);
}

.content-area {
  padding: 24px;
  display: flex;
  flex-direction: column;
  gap: 16px;
}

.filter-card {
  background-color: var(--color-bg-card);
  border-radius: var(--radius-md);
  border: 1px solid var(--color-border);
  padding: 16px;
  display: flex;
  gap: 12px;
  align-items: center;
}

.list-card {
  background-color: var(--color-bg-card);
  border-radius: var(--radius-md);
  border: 1px solid var(--color-border);
  padding: 24px;
}

.balance-value {
  font-weight: 600;
  color: var(--color-primary);
}

.pagination-row {
  display: flex;
  justify-content: center;
  margin-top: 16px;
}

.save-btn {
  background-color: var(--color-primary);
  border-color: var(--color-primary);
}
</style>
