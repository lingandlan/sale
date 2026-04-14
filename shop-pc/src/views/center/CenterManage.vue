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
        <el-select v-model="filters.level" placeholder="级别" style="width: 180px" clearable>
          <el-option label="子公司合伙人" value="1" />
          <el-option label="服务中心合伙人" value="2" />
        </el-select>
        <el-select v-model="filters.status" placeholder="状态" style="width: 120px" clearable>
          <el-option label="正常" value="normal" />
          <el-option label="冻结" value="frozen" />
        </el-select>
        <el-button type="primary" plain @click="handleSearch">查询</el-button>
        <el-button @click="handleResetFilter">重置</el-button>
      </div>

      <!-- 表格 -->
      <div class="list-card">
        <el-table
          :data="tableData"
          style="width: 100%"
          :header-cell-style="{ backgroundColor: '#FAFAFA', color: '#262626', fontWeight: '600' }"
        >
          <el-table-column prop="name" label="中心名称" width="200" />
          <el-table-column prop="region" label="省/市/区" width="200" />
          <el-table-column prop="level" label="级别" width="150">
            <template #default="{ row }">
              <el-tag :type="row.level === '子公司合伙人' ? '' : 'info'" size="small">{{ row.level }}</el-tag>
            </template>
          </el-table-column>
          <el-table-column prop="manager" label="管理员" width="120" />
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
        <el-form-item label="管理员">
          <el-select
            v-model="formData.managerId"
            placeholder="搜索操作员姓名或手机号"
            style="width: 100%"
            filterable
            clearable
          >
            <el-option
              v-for="op in operatorList"
              :key="op.id"
              :label="`${op.name}（${op.phone}）`"
              :value="op.id"
            />
          </el-select>
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="dialogVisible = false">取消</el-button>
        <el-button type="primary" class="save-btn" @click="handleSaveCenter">保存</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, computed, onMounted } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { getCenterList, createCenter, updateCenter, type CenterItem } from '@/api/center'
import { getOperatorList, type OperatorItem } from '@/api/operator'
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
  address: '',
  managerId: ''
})

const tableData = ref<any[]>([])
const operatorList = ref<OperatorItem[]>([])

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

const loadOperators = async () => {
  try {
    const res = await getOperatorList()
    if (res?.data) {
      operatorList.value = res.data
    }
  } catch {
    // fallback to empty
  }
}

const loadData = async () => {
  try {
    const res = await getCenterList()
    if (res?.data) {
      let list: CenterItem[] = res.data || []
      if (filters.keyword) {
        list = list.filter((c: CenterItem) => c.name.includes(filters.keyword))
      }
      if (filters.status) {
        list = list.filter((c: CenterItem) => c.status === filters.status)
      }
      tableData.value = list.map((c: CenterItem) => {
        const managerOp = operatorList.value.find(op => op.id === c.managerId)
        return {
          id: c.id,
          name: c.name,
          region: [c.province, c.city, c.district].filter(Boolean).join('/') || c.address || '-',
          level: c.code?.includes('服务') ? '服务中心合伙人' : '子公司合伙人',
          manager: c.managerName ? `${c.managerName}（${c.managerPhone}）` : '-',
          balance: c.balance ?? 0,
          totalIn: c.totalRecharge ?? 0,
          totalOut: c.totalConsumed ?? 0,
          status: c.status === 'active' ? 'normal' : 'frozen'
        }
      })
      pagination.total = tableData.value.length
    }
  } catch (error) {
    // fallback to empty data
  }
}

const handleSearch = () => { loadData() }
const handleResetFilter = () => { filters.keyword = ''; filters.level = ''; filters.status = ''; loadData() }

const handleAdd = () => {
  dialogTitle.value = '新建充值中心'
  Object.assign(formData, { id: '', name: '', level: '', province: '', city: '', district: '', address: '', managerId: '' })
  dialogVisible.value = true
}

const handleEdit = (row: any) => {
  dialogTitle.value = '编辑充值中心'
  // 从 region 字段解析出省/市/区
  const regionParts = (row.region || '').split('/')
  Object.assign(formData, {
    id: row.id,
    name: row.name,
    level: row.level,
    province: regionParts[0] || '',
    city: regionParts[1] || '',
    district: regionParts[2] || '',
    address: row.address || '',
    managerId: row.managerId || ''
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
      managerId: formData.managerId,
    }
    if (formData.id) {
      await updateCenter(formData.id, { ...payload, status: 'active' })
    } else {
      await createCenter(payload)
    }
    ElMessage.success(formData.id ? '编辑成功' : '创建成功')
    dialogVisible.value = false
    loadData()
  } catch (error) {
    // error handled by interceptor
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
  } catch {}
}

const handleDetail = (row: any) => {
  ElMessage.info(`查看 ${row.name} 详情`)
}

onMounted(() => {
  loadOperators()
  loadData()
})
</script>

<style scoped>
.center-manage {
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

.add-btn {
  background-color: #C00000;
  border-color: #C00000;
  border-radius: 4px;
}

.add-btn:hover {
  background-color: #A00000;
  border-color: #A00000;
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

.balance-value {
  font-weight: 600;
  color: #C00000;
}

.pagination-row {
  display: flex;
  justify-content: center;
  margin-top: 16px;
}

.save-btn {
  background-color: #C00000;
  border-color: #C00000;
}
</style>
