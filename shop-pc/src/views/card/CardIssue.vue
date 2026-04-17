<template>
  <div class="card-issue">
    <div class="page-header">
      <h1 class="page-title">绑定卡号和用户</h1>
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
          <el-form-item label="卡号" prop="cardNo">
            <el-input v-model="formData.cardNo" placeholder="请输入卡号（如 TJ00000001）" />
          </el-form-item>

          <el-form-item label="用户手机号" prop="userPhone">
            <el-input v-model="formData.userPhone" placeholder="请输入用户手机号" maxlength="11" />
          </el-form-item>

          <el-form-item label="充值中心" prop="rechargeCenterId">
            <el-select v-model="formData.rechargeCenterId" placeholder="选择充值中心" style="width: 100%">
              <el-option v-for="c in centers" :key="c.id" :label="c.name" :value="c.id" />
            </el-select>
          </el-form-item>

          <el-form-item label="发放原因" prop="issueReason">
            <el-select v-model="formData.issueReason" placeholder="选择发放原因" style="width: 100%">
              <el-option value="购买套餐包" label="购买套餐包" />
              <el-option value="推荐奖励" label="推荐奖励" />
              <el-option value="其他" label="其他" />
            </el-select>
          </el-form-item>

          <el-form-item label="发放方式" prop="issueType">
            <el-select v-model="formData.issueType" placeholder="选择发放方式" style="width: 100%">
              <el-option :value="1" label="实体卡（运营绑定）" />
              <el-option :value="2" label="虚拟卡（用户领取）" />
            </el-select>
          </el-form-item>

          <el-form-item v-if="formData.issueReason === '推荐奖励'" label="关联购买人" prop="relatedUserPhone">
            <el-input v-model="formData.relatedUserPhone" placeholder="被推荐人手机号" maxlength="11" />
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
              确认绑定
            </el-button>
            <el-button class="cancel-btn" @click="handleCancel">取消</el-button>
          </el-form-item>
        </el-form>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { ElMessage, ElMessageBox, type FormInstance, type FormRules } from 'element-plus'
import { extractErrorMessage } from '@/utils/request'
import { bindCard } from '@/api/card'
import request from '@/utils/request'

const formRef = ref<FormInstance>()
const centers = ref<{ id: string; name: string }[]>([])

const formData = ref({
  cardNo: '',
  userPhone: '',
  rechargeCenterId: '',
  issueReason: '',
  issueType: 1,
  relatedUserPhone: '',
  remark: ''
})

const formRules: FormRules = {
  cardNo: [{ required: true, message: '请输入卡号', trigger: 'blur' }],
  userPhone: [
    { required: true, message: '请输入手机号', trigger: 'blur' },
    { pattern: /^1[3-9]\d{9}$/, message: '请输入正确的手机号', trigger: 'blur' }
  ],
  rechargeCenterId: [{ required: true, message: '请选择充值中心', trigger: 'change' }],
  issueReason: [{ required: true, message: '请选择发放原因', trigger: 'change' }],
  issueType: [{ required: true, message: '请选择发放方式', trigger: 'change' }]
}

async function loadCenters() {
  const res = await request.get('/center')
  const data = res.data || res
  centers.value = Array.isArray(data) ? data : (data.list || [])
}

const handleSubmit = async () => {
  if (!formRef.value) return

  try {
    await formRef.value.validate()
  } catch {
    return
  }

  try {
    await bindCard(formData.value)
    ElMessage.success('绑定成功')
    formRef.value.resetFields()
  } catch (err: any) {
    ElMessage.error(extractErrorMessage(err, '绑定失败'))
  }
}

const handleCancel = () => {
  ElMessageBox.confirm('确认取消？已填写的信息将不会保存', '提示', {
    confirmButtonText: '确认',
    cancelButtonText: '继续填写',
    type: 'warning'
  }).then(() => {
    formRef.value?.resetFields()
  }).catch(() => {})
}

onMounted(() => { loadCenters() })
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
