<template>
  <div class="card-inventory">
    <el-card shadow="never">
      <template #header>
        <div class="card-header">
          <span>总卡库管理</span>
          <el-button type="primary" @click="showImportDialog = true">批量入库</el-button>
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
        <el-form-item label="起始卡号" prop="startCardNo">
          <el-input v-model="allocateForm.startCardNo" placeholder="如 TJ00000001" />
        </el-form-item>
        <el-form-item label="结束卡号" prop="endCardNo">
          <el-input v-model="allocateForm.endCardNo" placeholder="如 TJ00000010" />
        </el-form-item>
        <el-form-item>
          <el-button type="primary" @click="handleAllocate" :loading="allocateLoading">确认划拨</el-button>
        </el-form-item>
      </el-form>
    </el-card>

    <!-- 批量入库对话框 -->
    <el-dialog v-model="showImportDialog" title="批量入库" width="500px">
      <el-form :model="importForm" label-width="100px" :rules="importRules" ref="importFormRef">
        <el-form-item label="起始序号" prop="startSeq">
          <el-input-number v-model="importForm.startSeq" :min="1" :max="99999999" style="width: 100%" />
        </el-form-item>
        <el-form-item label="结束序号" prop="endSeq">
          <el-input-number v-model="importForm.endSeq" :min="1" :max="99999999" style="width: 100%" />
        </el-form-item>
        <el-form-item label="卡类型" prop="cardType">
          <el-select v-model="importForm.cardType" style="width: 100%">
            <el-option :value="1" label="实体卡" />
            <el-option :value="2" label="虚拟卡" />
          </el-select>
        </el-form-item>
        <el-form-item label="数量预览">
          <span>{{ importCount }} 张</span>
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="showImportDialog = false">取消</el-button>
        <el-button type="primary" @click="handleImport" :loading="importLoading">确认入库</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import type { FormInstance } from 'element-plus'
import { batchImportCards, allocateCards, getCardInventoryStats, type CardInventoryResponse } from '@/api/card'
import request from '@/utils/request'

const inventory = ref<CardInventoryResponse>({ totalCards: 0, issuedCards: 0, inStockCards: 0 })
const centers = ref<{ id: string; name: string }[]>([])

// 批量入库
const showImportDialog = ref(false)
const importLoading = ref(false)
const importFormRef = ref<FormInstance>()
const importForm = ref({ startSeq: 1, endSeq: 10, cardType: 1 })
const importRules = {
  startSeq: [{ required: true, message: '请输入起始序号', trigger: 'blur' }],
  endSeq: [{ required: true, message: '请输入结束序号', trigger: 'blur' }],
  cardType: [{ required: true, message: '请选择卡类型', trigger: 'change' }]
}
const importCount = computed(() => {
  const s = importForm.value.startSeq
  const e = importForm.value.endSeq
  return s > 0 && e >= s ? e - s + 1 : 0
})

// 划拨
const allocateLoading = ref(false)
const allocateFormRef = ref<FormInstance>()
const allocateForm = ref({ centerId: '', startCardNo: '', endCardNo: '' })
const allocateRules = {
  centerId: [{ required: true, message: '请选择充值中心', trigger: 'change' }],
  startCardNo: [{ required: true, message: '请输入起始卡号', trigger: 'blur' }],
  endCardNo: [{ required: true, message: '请输入结束卡号', trigger: 'blur' }]
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
  await importFormRef.value?.validate()
  await ElMessageBox.confirm(`确认入库 ${importCount.value} 张门店卡？`, '确认')
  importLoading.value = true
  const res = await batchImportCards(importForm.value).finally(() => { importLoading.value = false })
  ElMessage.success(`成功入库 ${(res.data || res).count} 张卡`)
  showImportDialog.value = false
  loadInventory()
}

async function handleAllocate() {
  await allocateFormRef.value?.validate()
  await ElMessageBox.confirm(`确认划拨卡号 ${allocateForm.value.startCardNo} ~ ${allocateForm.value.endCardNo}？`, '确认')
  allocateLoading.value = true
  const res = await allocateCards(allocateForm.value).finally(() => { allocateLoading.value = false })
  ElMessage.success(`成功划拨 ${(res.data || res).count} 张卡`)
  allocateForm.value = { centerId: '', startCardNo: '', endCardNo: '' }
  loadInventory()
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
.stats-row {
  text-align: center;
}
</style>
