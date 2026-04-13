<template>
  <div class="card-issue">
    <div class="page-header">
      <h1 class="page-title">门店卡发放</h1>
    </div>

    <div class="content-area">
      <div class="form-card">
        <div class="form-header">🎫 发放信息</div>
        <el-divider />

        <el-form
          ref="formRef"
          :model="formData"
          :rules="formRules"
          label-width="120px"
          class="issue-form"
        >
          <el-form-item label="持卡人姓名" prop="holderName">
            <el-input v-model="formData.holderName" placeholder="请输入持卡人姓名" />
          </el-form-item>

          <el-form-item label="持卡人手机" prop="holderPhone">
            <el-input v-model="formData.holderPhone" placeholder="请输入手机号" maxlength="11" />
          </el-form-item>

          <el-form-item label="充值金额" prop="amount">
            <el-input-number
              v-model="formData.amount"
              :min="100"
              :step="100"
              :precision="0"
              style="width: 100%"
            />
            <div class="form-hint">最低100元，必须是100的整数倍</div>
          </el-form-item>

          <el-form-item label="备注" prop="remark">
            <el-input
              v-model="formData.remark"
              type="textarea"
              :rows="3"
              placeholder="请输入备注信息（可选）"
            />
          </el-form-item>

          <el-form-item>
            <el-button type="primary" class="submit-btn" @click="handleSubmit">
              确认发放
            </el-button>
            <el-button class="cancel-btn" @click="handleCancel">取消</el-button>
          </el-form-item>
        </el-form>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import { ElMessage, ElMessageBox, type FormInstance, type FormRules } from 'element-plus'
import { issueCard } from '@/api/card'

const formRef = ref<FormInstance>()

const formData = ref({
  holderName: '',
  holderPhone: '',
  amount: 100,
  remark: ''
})

const formRules: FormRules = {
  holderName: [
    { required: true, message: '请输入持卡人姓名', trigger: 'blur' }
  ],
  holderPhone: [
    { required: true, message: '请输入手机号', trigger: 'blur' },
    { pattern: /^1[3-9]\d{9}$/, message: '请输入正确的手机号', trigger: 'blur' }
  ],
  amount: [
    { required: true, message: '请输入充值金额', trigger: 'blur' },
    { type: 'number', min: 100, message: '最低100元', trigger: 'blur' }
  ]
}

const handleSubmit = async () => {
  if (!formRef.value) return

  try {
    await formRef.value.validate()
  } catch {
    return
  }

  try {
    await issueCard({
      holderName: formData.value.holderName,
      holderPhone: formData.value.holderPhone,
      amount: formData.value.amount,
      remark: formData.value.remark
    })
    ElMessage.success('发放成功')
    handleCancel()
  } catch (error) {
    ElMessage.error('发放失败')
  }
}

const handleCancel = () => {
  ElMessageBox.confirm('确认取消？已填写的信息将不会保存', '提示', {
    confirmButtonText: '确认',
    cancelButtonText: '继续填写',
    type: 'warning'
  }).then(() => {
    formRef.value?.resetFields()
  }).catch(() => {
    // 用户选择继续填写
  })
}
</script>

<style scoped>
.card-issue {
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

.content-area {
  padding: 24px;
}

.form-card {
  background-color: #FFFFFF;
  border-radius: 8px;
  border: 1px solid #E5E5E5;
  padding: 24px;
  max-width: 600px;
}

.form-header {
  font-family: 'Inter', sans-serif;
  font-size: 16px;
  font-weight: 600;
  color: #262626;
}

.issue-form {
  margin-top: 16px;
}

.form-hint {
  font-family: 'Inter', sans-serif;
  font-size: 12px;
  color: #8C8C8C;
  margin-top: 4px;
}

.submit-btn {
  width: 120px;
  height: 40px;
  border-radius: 4px;
  background-color: #C00000;
  border-color: #C00000;
}

.cancel-btn {
  width: 100px;
  height: 40px;
  border-radius: 4px;
}
</style>
