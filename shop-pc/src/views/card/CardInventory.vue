<template>
  <div class="card-inventory">
    <el-card shadow="never">
      <template #header>
        <div class="card-header">
          <span>总卡库管理</span>
          <div class="card-header-actions">
            <el-button @click="downloadTemplate">下载模板</el-button>
            <el-button type="primary" @click="showImportDialog = true">批量入库</el-button>
          </div>
        </div>
      </template>

      <!-- 库存统计 -->
      <el-row :gutter="20" class="stats-row">
        <el-col :span="8">
          <el-statistic title="总卡数" :value="inventory.totalCards" />
        </el-col>
        <el-col :span="8">
          <el-statistic title="已出库" :value="inventory.issuedCards" />
        </el-col>
        <el-col :span="8">
          <el-statistic title="剩余库存" :value="inventory.inStockCards" />
        </el-col>
      </el-row>
    </el-card>

    <!-- 划拨到充值中心 -->
    <el-card shadow="never" style="margin-top: 16px">
      <template #header>
        <span>划拨到充值中心</span>
      </template>
      <el-form :model="allocateForm" label-width="100px" :rules="allocateRules" ref="allocateFormRef">
        <el-form-item label="目标中心" prop="centerId">
          <el-select v-model="allocateForm.centerId" placeholder="选择充值中心" style="width: 100%">
            <el-option v-for="c in centers" :key="c.id" :label="c.name" :value="c.id" />
          </el-select>
        </el-form-item>
        <el-form-item label="划拨数量" prop="quantity">
          <el-input-number v-model="allocateForm.quantity" :min="1" :max="1000" style="width: 100%" />
        </el-form-item>
        <el-form-item>
          <el-button type="primary" @click="handleAllocate" :loading="allocateLoading">确认划拨</el-button>
        </el-form-item>
      </el-form>
    </el-card>

    <!-- 批量入库对话框 -->
    <el-dialog v-model="showImportDialog" title="批量入库" width="500px">
      <el-upload
        ref="uploadRef"
        :auto-upload="false"
        :limit="1"
        accept=".xlsx,.xls,.csv"
        :on-change="handleFileChange"
        :on-remove="handleFileRemove"
        drag
      >
        <el-icon style="font-size: 48px; color: #c0c4cc"><upload-filled /></el-icon>
        <div style="margin-top: 8px">拖拽文件到此处，或<em>点击上传</em></div>
        <template #tip>
          <div class="el-upload__tip">支持 .xlsx / .csv 文件，格式：卡号 | 卡类型 | 面值</div>
        </template>
      </el-upload>
      <template #footer>
        <el-button @click="showImportDialog = false">取消</el-button>
        <el-button type="primary" @click="handleImport" :loading="importLoading" :disabled="!uploadFile">确认入库</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { UploadFilled } from '@element-plus/icons-vue'
import type { FormInstance, UploadFile } from 'element-plus'
import { batchImportCards, allocateCards, getCardInventoryStats, type CardInventoryResponse } from '@/api/card'
import request from '@/utils/request'

const inventory = ref<CardInventoryResponse>({ totalCards: 0, issuedCards: 0, inStockCards: 0 })
const centers = ref<{ id: string; name: string }[]>([])

// 批量入库
const showImportDialog = ref(false)
const importLoading = ref(false)
const uploadRef = ref()
const uploadFile = ref<File | null>(null)

function handleFileChange(file: UploadFile) {
  if (file.raw) {
    uploadFile.value = file.raw
  }
}

function downloadTemplate() {
  // 生成 UTF-8 BOM + CSV 内容，Excel 可直接打开
  const rows = [
    '卡号,卡类型,面值',
    'xx1,实体,1000',
    'xx2,虚拟,500'
  ]
  const blob = new Blob(['\uFEFF' + rows.join('\n')], { type: 'text/csv;charset=utf-8' })
  const url = URL.createObjectURL(blob)
  const a = document.createElement('a')
  a.href = url
  a.download = '门店卡导入模板.csv'
  a.click()
  URL.revokeObjectURL(url)
}

function handleFileRemove() {
  uploadFile.value = null
}

// 划拨
const allocateLoading = ref(false)
const allocateFormRef = ref<FormInstance>()
const allocateForm = ref({ centerId: '', quantity: 1 })
const allocateRules = {
  centerId: [{ required: true, message: '请选择充值中心', trigger: 'change' }],
  quantity: [{ required: true, message: '请输入划拨数量', trigger: 'blur' }]
}

async function loadInventory() {
  const res = await getCardInventoryStats()
  inventory.value = res.data || res
}

async function loadCenters() {
  const res = await request.get('/center')
  const data = res.data || res
  centers.value = Array.isArray(data) ? data : (data.list || [])
}

async function handleImport() {
  if (!uploadFile.value) {
    ElMessage.warning('请先选择文件')
    return
  }
  await ElMessageBox.confirm('确认导入所选Excel文件中的门店卡？', '确认')
  importLoading.value = true
  const file = uploadFile.value
  try {
    const res = await batchImportCards(file)
    ElMessage.success(`成功入库 ${(res.data || res).count} 张卡`)
    showImportDialog.value = false
    uploadFile.value = null
    loadInventory()
  } catch (err: any) {
    const msg = err?.response?.data?.msg || err?.message || '导入失败'
    ElMessage.error(msg)
  } finally {
    importLoading.value = false
  }
}

async function handleAllocate() {
  await allocateFormRef.value?.validate()
  await ElMessageBox.confirm(`确认划拨 ${allocateForm.value.quantity} 张卡到所选充值中心？`, '确认')
  allocateLoading.value = true
  try {
    const res = await allocateCards(allocateForm.value)
    ElMessage.success(`成功划拨 ${(res.data || res).count} 张卡`)
    allocateForm.value = { centerId: '', quantity: 1 }
    loadInventory()
  } catch (err: any) {
    const msg = err?.response?.data?.msg || err?.message || '划拨失败'
    ElMessage.error(msg)
  } finally {
    allocateLoading.value = false
  }
}

onMounted(() => {
  loadInventory()
  loadCenters()
})
</script>

<style scoped>
.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}
.card-header-actions {
  display: flex;
  gap: 8px;
}
.stats-row {
  text-align: center;
}
</style>
